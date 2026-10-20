package integration_test

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/config"
	"github.com/WildanFrananda/kinetix-notary-service/internal/contractgen/withholding"
	"github.com/WildanFrananda/kinetix-notary-service/internal/ethereum"
	"github.com/WildanFrananda/kinetix-notary-service/internal/keccak"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/poseidon"
	"github.com/WildanFrananda/kinetix-notary-service/internal/verifycli"
	"github.com/WildanFrananda/kinetix-notary-service/internal/verifying"
)

func (c *chain) genesisAfter(t *testing.T, predecessor common.Address) (uint64, notary.Digest, uint64) {
	t.Helper()

	contract := c.binding.Instance(c.client, predecessor)
	options := &bind.CallOpts{Context: context.Background()}
	next, err := bind.Call(contract, options, c.binding.PackNextEpoch(), c.binding.UnpackNextEpoch)
	require.NoError(t, err)
	head, err := bind.Call(contract, options, c.binding.PackHead(), c.binding.UnpackHead)
	require.NoError(t, err)
	lastAsOf, err := bind.Call(contract, options, c.binding.PackLastAsOf(), c.binding.UnpackLastAsOf)
	require.NoError(t, err)

	return next, notary.Digest(head), lastAsOf
}

func (e *estate) succeed(t *testing.T) {
	t.Helper()

	c := e.chain
	predecessor := common.Address(e.ledger.Address)
	c.seal(t, predecessor)

	next, head, lastAsOf := c.genesisAfter(t, predecessor)
	setup := testConfidential()
	two := c.deployConfidentialLedger(t, setup.verifier(c, t), setup.depth, c.from, next, head, lastAsOf)
	council := c.deployConfidentialCouncil(t, two, notary.Address(c.from), e.watchMember)
	c.handOver(t, two, council)
	court := c.deployCourt(t, two, e.signer)
	withholdingCourt := c.deployWithholding(t, court)
	c.registerSuccessor(t, common.Address(e.directory), two, court, withholdingCourt)

	e.ledger = c.identity(two)
	e.council = notary.Address(council)
	e.court = notary.Address(court)
	e.withholding = notary.Address(withholdingCourt)
	e.paymaster = notary.Address(c.sponsorship(t, court, withholdingCourt).paymaster)

	path, sha := setup.srs(t)

	values := map[string]string{
		config.LedgerAddressKey:    e.ledger.Address.Hex(),
		config.CourtAddressKey:     e.court.Hex(),
		config.WithholdingKey:      e.withholding.Hex(),
		config.CouncilAddressKey:   e.council.Hex(),
		config.PaymasterAddressKey: e.paymaster.Hex(),
		config.WatchAddressKey:     e.serveWatch(t, e.pki, e.paymentAddr, e.watchKey),
		config.SchemeKey:           "2",
		config.SRSPathKey:          path,
		config.SRSSHA256Key:        sha,
	}
	for key, value := range values {
		e.values[key] = value
	}

	putProverOnPath(t)
}

func (e *estate) verifyThroughTheDirectory(t *testing.T, r running, epoch, holder string) map[string]string {
	t.Helper()

	ctx := context.Background()
	response := e.proofs(t, r, epoch, holder)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())

	path := filepath.Join(t.TempDir(), "proofs.json")
	require.NoError(t, os.WriteFile(path, response.Body.Bytes(), 0o600))

	settings := ethereum.Settings{RPCURL: e.chain.url, Ledger: e.ledger, Timeout: 10 * time.Second}
	directory := verifying.NewDirectoryVerifier(
		ethereum.NewDirectory(e.chain.client, ethereum.Settings{
			RPCURL:    e.chain.url,
			Ledger:    e.ledger,
			Directory: notary.DirectoryIdentity{ChainID: e.ledger.ChainID, Address: e.directory},
			Timeout:   10 * time.Second,
		}),
		ethereum.NewReaders(e.chain.client, settings),
		keccak.NewScheme(),
		poseidon.NewScheme(),
		keccak.NewScheme(),
	)

	var out bytes.Buffer

	code := verifycli.Run(ctx, path, poseidon.NewScheme().Version(), directory, &out)
	require.Equal(t, verifycli.ExitIncluded, code, out.String())

	var report map[string]string
	require.NoError(t, json.Unmarshal(bytes.TrimSpace(out.Bytes()), &report))

	return report
}

func requestStatement(
	t *testing.T,
	e *estate,
	key *ecdsa.PrivateKey,
	court notary.Address,
	issued issuedBody,
	epoch uint64,
) common.Hash {
	t.Helper()

	ctx := context.Background()
	one := issued.Registrations[0]
	registration := withholding.Registration{
		Holder:    common.HexToAddress(one.Holder),
		Account:   common.HexToHash(one.Account),
		FromEpoch: 0,
	}
	binding := withholding.NewWithholdingCourt()
	tx, err := bind.Transact(
		binding.Instance(e.chain.client, common.Address(court)),
		bind.NewKeyedTransactor(key, e.chain.chainID),
		binding.PackRequestStatement(registration, hexutil.MustDecode(one.Signature), epoch),
	)
	require.NoError(t, err)

	receipt, err := bind.WaitMined(ctx, e.chain.client, tx.Hash())
	require.NoError(t, err)
	require.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)

	return common.HexToHash(one.Digest)
}

func requestState(t *testing.T, e *estate, court notary.Address, digest common.Hash, epoch uint64) uint8 {
	t.Helper()

	binding := withholding.NewWithholdingCourt()
	request, err := bind.Call(
		binding.Instance(e.chain.client, common.Address(court)),
		&bind.CallOpts{Context: context.Background()},
		binding.PackRequestOf(digest, epoch),
		binding.UnpackRequestOf,
	)
	require.NoError(t, err)

	return request.State
}

func TestTheDaemonCarriesItsHistoryIntoAConfidentialSuccessor(t *testing.T) {
	e := newEstate(t)
	e.count = 6
	r := e.start(t)
	ctx := context.Background()

	e.tickUntil(t, r, 0, "mined")
	e.chain.finalize(t)
	e.tickUntil(t, r, 0, "finalized")

	key, _ := e.chain.funded(t)
	registered := e.register(t, r, "principal-3", key, e.clock())
	require.Equal(t, http.StatusOK, registered.Code, registered.Body.String())

	var before issuedBody
	require.NoError(t, json.Unmarshal(registered.Body.Bytes(), &before))

	first := e.ledger
	firstCourt := e.court
	firstWithholding := e.withholding

	e.succeed(t)
	e.advance(2 * time.Hour)

	_, err := e.pool.Exec(ctx, "UPDATE writer_lease SET expires_at = now()")
	require.NoError(t, err)

	next := e.start(t)
	require.NoError(t, next.Directory.Check(ctx))
	require.NoError(t, next.Prover.CheckCourt(ctx))
	require.NoError(t, next.Registrar.CheckWithholding(ctx))

	e.tickUntil(t, next, 1, "mined")
	e.chain.finalize(t)
	e.tickUntil(t, next, 1, "finalized")

	for number, want := range map[int64]notary.LedgerIdentity{0: first, 1: e.ledger} {
		var (
			address []byte
			scheme  int
		)
		require.NoError(
			t,
			e.pool.QueryRow(ctx, "SELECT ledger_address, scheme FROM epochs WHERE number = $1", number).
				Scan(&address, &scheme),
		)
		require.Equal(t, want.Address[:], address, "epoch %d", number)
		require.Equal(t, int(number)+1, scheme, "epoch %d", number)
	}

	for epoch, total := range map[string]string{"0": "21000000", "1": "0"} {
		report := e.verifyThroughTheDirectory(t, next, epoch, "principal-2")
		require.Equal(t, "included", report["verdict"], "epoch %s", epoch)
		require.Equal(t, "attributable", report["statement"], "epoch %s", epoch)
		require.Equal(t, total, report["total"], "epoch %s", epoch)
	}

	asked := requestStatement(t, e, key, firstWithholding, before, 0)
	require.Equal(t, uint8(1), requestState(t, e, firstWithholding, asked, 0))
	require.NoError(t, next.Answerer.Tick(ctx))
	require.Equal(t, uint8(2), requestState(t, e, firstWithholding, asked, 0))

	again := e.register(t, next, "principal-3", key, e.clock())
	require.Equal(t, http.StatusOK, again.Code, again.Body.String())

	var after issuedBody
	require.NoError(t, json.Unmarshal(again.Body.Bytes(), &after))
	require.Equal(t, e.withholding.Hex(), after.Court.Address)
	require.Equal(t, "1", after.Epoch)
	require.NotEqual(t, firstCourt, e.court)

	scraped := httptest.NewRecorder()
	next.Handler.ServeHTTP(scraped, httptest.NewRequest(http.MethodGet, "/metrics", nil))
	require.Contains(t, scraped.Body.String(), "kinetix_notary_accounts 6")
	require.Contains(t, scraped.Body.String(), "kinetix_notary_account_capacity 8")
	require.Contains(t, scraped.Body.String(), "kinetix_notary_last_finalized_epoch 1")
}
