package integration_test

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/contractgen/withholding"
	"github.com/WildanFrananda/kinetix-notary-service/internal/keccak"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

type issuedBody struct {
	Court struct {
		ChainID string `json:"chainId"`
		Address string `json:"address"`
	} `json:"court"`
	Epoch         string `json:"epoch"`
	Registrations []struct {
		Holder    string `json:"holder"`
		Account   string `json:"account"`
		FromEpoch string `json:"fromEpoch"`
		Kind      string `json:"kind"`
		Digest    string `json:"digest"`
		Signature string `json:"signature"`
	} `json:"registrations"`
}

func (c *chain) funded(t *testing.T) (*ecdsa.PrivateKey, common.Address) {
	t.Helper()

	key, err := crypto.GenerateKey()
	require.NoError(t, err)

	address := crypto.PubkeyToAddress(key.PublicKey)
	funds := new(big.Int).Mul(big.NewInt(1), big.NewInt(1e18))
	require.NoError(
		t,
		c.rpc.CallContext(context.Background(), nil, "anvil_setBalance", address, hexutil.EncodeBig(funds)),
	)

	return key, address
}

func (e *estate) register(
	t *testing.T,
	r running,
	subject string,
	key *ecdsa.PrivateKey,
	issuedAt time.Time,
) *httptest.ResponseRecorder {
	t.Helper()

	holder := notary.Address(crypto.PubkeyToAddress(key.PublicKey))
	request := notary.RegistrationRequest{
		Holder:    holder,
		Principal: notarytest.Holder(t, subject),
		IssuedAt:  notarytest.Unix(t, uint64(issuedAt.Unix())),
	}
	court := notary.CourtIdentity{ChainID: e.ledger.ChainID, Address: e.withholding}
	digest := keccak.NewScheme().RegistrationRequestDigest(court, request)
	raw, err := crypto.Sign(digest[:], key)
	require.NoError(t, err)

	raw[64] += 27
	body := `{"holder":"` + holder.Hex() + `","issuedAt":"` + strconv.FormatInt(issuedAt.Unix(), 10) +
		`","signature":"` + hexutil.Encode(raw) + `"}`

	post := httptest.NewRequest(http.MethodPost, "/api/v1/proofs/registrations", strings.NewReader(body))
	post.Header.Set("Authorization", "Bearer "+e.identity.Token(t, subject))

	response := httptest.NewRecorder()
	r.Handler.ServeHTTP(response, post)

	return response
}

func TestAHolderRegistersThroughTheAPIAndTheCourtAcceptsTheirRequest(t *testing.T) {
	e := newEstate(t)
	r := e.start(t)
	ctx := context.Background()
	key, address := e.chain.funded(t)

	require.NoError(t, r.Registrar.CheckWithholding(ctx))

	early := e.register(t, r, "principal-3", key, e.clock())
	require.Equal(t, http.StatusServiceUnavailable, early.Code, early.Body.String())
	require.Contains(t, early.Body.String(), `"error":"registration_not_open"`)

	e.tickUntil(t, r, 0, "mined")
	e.chain.finalize(t)
	e.tickUntil(t, r, 0, "finalized")

	response := e.register(t, r, "principal-3", key, e.clock())
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())

	var issued issuedBody
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &issued))
	require.Equal(t, e.withholding.Hex(), issued.Court.Address)
	require.Equal(t, "0", issued.Epoch)
	require.Len(t, issued.Registrations, 1)

	one := issued.Registrations[0]
	require.Equal(t, notary.Address(address).Hex(), one.Holder)
	require.Equal(t, "driver_available", one.Kind)

	registration := withholding.Registration{
		Holder:    address,
		Account:   common.HexToHash(one.Account),
		FromEpoch: 0,
	}
	binding := withholding.NewWithholdingCourt()
	instance := binding.Instance(e.chain.client, common.Address(e.withholding))

	opts := bind.NewKeyedTransactor(key, e.chain.chainID)
	tx, err := bind.Transact(
		instance,
		opts,
		binding.PackRequestStatement(registration, hexutil.MustDecode(one.Signature), 0),
	)
	require.NoError(t, err)

	receipt, err := bind.WaitMined(ctx, e.chain.client, tx.Hash())
	require.NoError(t, err)
	require.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)

	request, err := bind.Call(
		instance,
		&bind.CallOpts{Context: ctx},
		binding.PackRequestOf(common.HexToHash(one.Digest), 0),
		binding.UnpackRequestOf,
	)
	require.NoError(t, err)
	require.Equal(t, uint8(1), request.State)

	stale := e.register(t, r, "principal-3", key, e.clock().Add(-10*time.Minute))
	require.Equal(t, http.StatusForbidden, stale.Code, stale.Body.String())

	nobody := e.register(t, r, "nobody-at-all", key, e.clock())
	require.Equal(t, http.StatusNotFound, nobody.Code, nobody.Body.String())
}
