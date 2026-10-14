package integration_test

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/contractgen/withholding"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type withholdingRig struct {
	e        *estate
	binding  *withholding.WithholdingCourt
	instance *bind.BoundContract
}

func (w withholdingRig) send(t *testing.T, key *ecdsa.PrivateKey, calldata []byte) {
	t.Helper()

	opts := bind.NewKeyedTransactor(key, w.e.chain.chainID)
	tx, err := bind.Transact(w.instance, opts, calldata)
	require.NoError(t, err)

	receipt, err := bind.WaitMined(context.Background(), w.e.chain.client, tx.Hash())
	require.NoError(t, err)
	require.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)
}

func (w withholdingRig) state(t *testing.T, digest common.Hash, epoch uint64) notary.RequestState {
	t.Helper()

	request, err := bind.Call(
		w.instance,
		&bind.CallOpts{Context: context.Background()},
		w.binding.PackRequestOf(digest, epoch),
		w.binding.UnpackRequestOf,
	)
	require.NoError(t, err)

	return notary.RequestState(request.State)
}

func (e *estate) metric(t *testing.T, r running) string {
	t.Helper()

	response := httptest.NewRecorder()
	r.Handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/metrics", nil))
	require.Equal(t, http.StatusOK, response.Code)

	return response.Body.String()
}

func TestTheNotaryAnswersWhatItHoldsAndAnOmissionEndsInARecordedSilence(t *testing.T) {
	e := newEstate(t)
	r := e.start(t)
	ctx := context.Background()
	binding := withholding.NewWithholdingCourt()
	w := withholdingRig{
		e:        e,
		binding:  binding,
		instance: binding.Instance(e.chain.client, common.Address(e.withholding)),
	}

	e.tickUntil(t, r, 0, "mined")
	e.chain.finalize(t)
	e.tickUntil(t, r, 0, "finalized")

	key, address := e.chain.funded(t)
	response := e.register(t, r, "principal-3", key, e.clock())
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())

	var issued issuedBody
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &issued))
	require.Len(t, issued.Registrations, 1)

	one := issued.Registrations[0]
	digest := common.HexToHash(one.Digest)
	registration := withholding.Registration{Holder: address, Account: common.HexToHash(one.Account), FromEpoch: 0}
	signature := hexutil.MustDecode(one.Signature)

	w.send(t, key, binding.PackRequestStatement(registration, signature, 0))
	require.Equal(t, notary.RequestOpen, w.state(t, digest, 0))

	require.NoError(t, r.Answerer.Tick(ctx))
	require.Equal(t, notary.RequestAnswered, w.state(t, digest, 0))
	require.Contains(t, e.metric(t, r), "kinetix_notary_withholding_open_requests 0")

	e.mu.Lock()
	e.omit = "principal-3"
	e.mu.Unlock()
	require.NoError(t, e.chain.rpc.CallContext(ctx, nil, "evm_increaseTime", 7200))
	require.NoError(t, e.chain.rpc.CallContext(ctx, nil, "evm_mine"))
	e.advance(2 * time.Hour)
	e.tickUntil(t, r, 1, "mined")
	e.chain.finalize(t)
	e.tickUntil(t, r, 1, "finalized")

	w.send(t, key, binding.PackRequestStatement(registration, signature, 1))
	require.NoError(t, r.Answerer.Tick(ctx))
	require.Equal(t, notary.RequestOpen, w.state(t, digest, 1), "an account the epoch left out cannot be answered")

	metrics := e.metric(t, r)
	require.Contains(t, metrics, "kinetix_notary_withholding_unanswerable_requests 1")
	require.Contains(t, metrics, "kinetix_notary_withholding_silences 0")

	require.NoError(t, e.chain.rpc.CallContext(ctx, nil, "evm_increaseTime", 24*3600+60))
	require.NoError(t, e.chain.rpc.CallContext(ctx, nil, "evm_mine"))

	stranger, _ := e.chain.funded(t)
	w.send(t, stranger, binding.PackRecordUnanswered(registration, 1))
	require.Equal(t, notary.RequestUnanswered, w.state(t, digest, 1))

	require.NoError(t, r.Answerer.Tick(ctx))

	metrics = e.metric(t, r)
	require.Contains(t, metrics, "kinetix_notary_withholding_silences 1")
	require.Contains(t, metrics, "kinetix_notary_withholding_unanswerable_requests 0")
}
