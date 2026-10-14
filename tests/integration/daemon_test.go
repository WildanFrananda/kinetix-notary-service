package integration_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"go.uber.org/dig"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/WildanFrananda/kinetix-notary-service/internal/anchoring"
	"github.com/WildanFrananda/kinetix-notary-service/internal/answering"
	"github.com/WildanFrananda/kinetix-notary-service/internal/config"
	"github.com/WildanFrananda/kinetix-notary-service/internal/keccak"
	"github.com/WildanFrananda/kinetix-notary-service/internal/metrics"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/postgres"
	commonv1 "github.com/WildanFrananda/kinetix-notary-service/internal/protogen/common/v1"
	paymentv1 "github.com/WildanFrananda/kinetix-notary-service/internal/protogen/payment/v1"
	"github.com/WildanFrananda/kinetix-notary-service/internal/proving"
	"github.com/WildanFrananda/kinetix-notary-service/internal/registering"
	"github.com/WildanFrananda/kinetix-notary-service/internal/verifycli"
	"github.com/WildanFrananda/kinetix-notary-service/internal/verifying"
	"github.com/WildanFrananda/kinetix-notary-service/internal/wiring"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

const databaseURLKey = "KINETIX_NOTARY_TEST_DATABASE_URL"

type estate struct {
	chain       *chain
	ledger      notary.LedgerIdentity
	court       notary.Address
	withholding notary.Address
	dbURL       string
	pool        *pgxpool.Pool
	payment     *notarytest.LiabilityServer
	identity    *notarytest.Identity
	values      map[string]string
	omit        string
	mu          sync.Mutex
	now         time.Time
}

type running struct {
	dig.In
	Anchorer  *anchoring.Anchorer
	Handler   http.Handler
	Verifier  *verifying.Verifier
	Prover    *proving.Prover
	Registrar *registering.Registrar
	Answerer  *answering.Answerer
}

func newEstate(t *testing.T) *estate {
	t.Helper()

	base := os.Getenv(databaseURLKey)
	if base == "" {
		t.Skip(databaseURLKey + " is not set; this test needs Postgres, and CI must set it")
	}

	c := connect(t)
	ctx := context.Background()
	latest, err := c.client.HeaderByNumber(ctx, nil)
	require.NoError(t, err)
	e := &estate{chain: c, ledger: c.identity(c.deploy(t)), now: time.Unix(int64(latest.Time), 0).UTC()}
	statementKey, statementSigner := notarytest.StatementKey(t, "estate-signer")
	e.court = notary.Address(c.deployCourt(t, common.Address(e.ledger.Address), statementSigner))
	e.withholding = notary.Address(c.deployWithholding(t, common.Address(e.court)))
	answererKey, _ := c.funded(t)
	e.dbURL = schema(t, base)
	pool, err := pgxpool.New(ctx, e.dbURL)
	require.NoError(t, err)
	t.Cleanup(pool.Close)
	_, _, err = postgres.Migrate(ctx, pool)
	require.NoError(t, err)

	e.pool = pool

	e.identity = notarytest.NewIdentity(t)

	pki := notarytest.NewPKI(t)
	e.payment = &notarytest.LiabilityServer{Script: e.snapshot}
	address := notarytest.ServeLiabilities(t, pki.Issue(t, "kinetix.local", "payment"), e.payment)
	e.values = map[string]string{
		config.RPCURLKey:          c.url,
		config.ChainIDKey:         strconv.FormatUint(e.ledger.ChainID, 10),
		config.LedgerAddressKey:   e.ledger.Address.Hex(),
		config.CourtAddressKey:    e.court.Hex(),
		config.StatementKeysKey:   "0x" + hex.EncodeToString(crypto.FromECDSA(statementKey)),
		config.WithholdingKey:     e.withholding.Hex(),
		config.WithholdingFromKey: "0",
		config.LogRangeKey:        "200",
		config.AnswererKeyKey:     "0x" + hex.EncodeToString(crypto.FromECDSA(answererKey)),
		config.RPCTimeoutKey:      "10s",
		config.DatabaseURLKey:     e.dbURL,
		config.AnchorerKeyKey:     "0x" + hex.EncodeToString(crypto.FromECDSA(c.key)),
		config.PaymentAddressKey:  address,
		config.PKIDirKey:          pki.Issue(t, "kinetix.local", "notary"),
		config.CurrencyKey:        "IDR",
		config.MaxFeeKey:          "1000000000000",
		config.HTTPPortKey:        "18010",
		config.JWKSURLKey:         e.identity.URL,
		config.JWTIssuerKey:       notarytest.IdentityIssuer,
		config.JWTAudienceKey:     notarytest.IdentityAudience,
	}

	return e
}

func schema(t *testing.T, base string) string {
	t.Helper()

	ctx := context.Background()
	suffix := make([]byte, 6)
	_, err := rand.Read(suffix)
	require.NoError(t, err)

	name := "d_" + hex.EncodeToString(suffix)
	admin, err := pgx.Connect(ctx, base)
	require.NoError(t, err)
	_, err = admin.Exec(ctx, "CREATE SCHEMA "+name)
	require.NoError(t, err)
	t.Cleanup(func() {
		_, _ = admin.Exec(context.Background(), "DROP SCHEMA "+name+" CASCADE")
		_ = admin.Close(context.Background())
	})

	separator := "?"
	if strings.Contains(base, "?") {
		separator = "&"
	}

	return base + separator + "search_path=" + name
}

func (e *estate) snapshot() []*paymentv1.StreamLiabilitySnapshotResponse {
	head, err := e.chain.client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return nil
	}

	accounts := make([]*paymentv1.LiabilityAccount, 0, 9)

	var total int64

	e.mu.Lock()
	omit := e.omit
	e.mu.Unlock()

	for i := range 9 {
		if fmt.Sprintf("principal-%d", i) == omit {
			continue
		}

		amount := int64(1_000_000 * (i + 1))
		total += amount
		accounts = append(accounts, &paymentv1.LiabilityAccount{
			Kind:    paymentv1.LiabilityKind(i%5 + 1),
			Holder:  &paymentv1.LiabilityAccount_PrincipalId{PrincipalId: fmt.Sprintf("principal-%d", i)},
			Balance: &commonv1.Money{AmountMinor: amount, Currency: "IDR"},
		})
	}

	return []*paymentv1.StreamLiabilitySnapshotResponse{
		{Part: &paymentv1.StreamLiabilitySnapshotResponse_Header{Header: &paymentv1.LiabilitySnapshotHeader{
			AsOf:     timestamppb.New(time.Unix(int64(head.Time)-5, 0)),
			Currency: "IDR",
		}}},
		{
			Part: &paymentv1.StreamLiabilitySnapshotResponse_Accounts{
				Accounts: &paymentv1.LiabilityAccountBatch{Accounts: accounts},
			},
		},
		{Part: &paymentv1.StreamLiabilitySnapshotResponse_Trailer{Trailer: &paymentv1.LiabilitySnapshotTrailer{
			AccountCount: int64(len(accounts)),
			Total:        &commonv1.Money{AmountMinor: total, Currency: "IDR"},
		}}},
	}
}

func (e *estate) clock() time.Time {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.now
}

func (e *estate) advance(d time.Duration) {
	e.mu.Lock()
	e.now = e.now.Add(d)
	e.mu.Unlock()
}

func (e *estate) start(t *testing.T) running {
	t.Helper()

	lookup := func(key string) (string, bool) {
		value, ok := e.values[key]
		return value, ok
	}
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	container, err := wiring.Daemon(lookup, metrics.Build{Service: "kinetix-notary-service", Version: "test"}, log)
	require.NoError(t, err)
	require.NoError(t, container.Decorate(func(func() time.Time) func() time.Time { return e.clock }))
	require.NoError(t, container.Provide(verifying.NewVerifier))

	var r running

	require.NoError(t, container.Invoke(func(in running) { r = in }))

	return r
}

func (e *estate) state(t *testing.T, epoch uint64) string {
	t.Helper()

	var state string
	require.NoError(
		t,
		e.pool.QueryRow(context.Background(), "SELECT state FROM epochs WHERE number = $1", int64(epoch)).Scan(&state),
	)

	return state
}

func (e *estate) tickUntil(t *testing.T, r running, epoch uint64, want string) {
	t.Helper()
	require.Eventually(
		t,
		func() bool {
			_ = r.Anchorer.Tick(context.Background())
			return e.state(t, epoch) == want
		},
		15*time.Second,
		100*time.Millisecond,
		"epoch %d never became %s",
		epoch,
		want,
	)
}

func (e *estate) accounts(t *testing.T, epoch uint64) []notary.Entry {
	t.Helper()

	rows, err := e.pool.Query(
		context.Background(),
		"SELECT commitment, amount FROM accounts WHERE epoch_number = $1 ORDER BY position",
		int64(epoch),
	)
	require.NoError(t, err)

	defer rows.Close()

	var entries []notary.Entry

	for rows.Next() {
		var (
			commitment []byte
			amount     int64
		)
		require.NoError(t, rows.Scan(&commitment, &amount))

		var entry notary.Entry
		copy(entry.Commitment[:], commitment)
		entry.Amount = uint64(amount)
		entries = append(entries, entry)
	}

	require.NoError(t, rows.Err())

	return entries
}

func TestTheDaemonAnchorsPaymentsSnapshotAndEveryAccountVerifiesOnChain(t *testing.T) {
	e := newEstate(t)
	r := e.start(t)
	ctx := context.Background()

	e.tickUntil(t, r, 0, "mined")
	e.chain.finalize(t)
	e.tickUntil(t, r, 0, "finalized")

	entries := e.accounts(t, 0)
	require.Len(t, entries, 9)
	tree, err := notary.BuildTree(keccak.NewScheme(), 0, entries)
	require.NoError(t, err)

	for i, entry := range entries {
		proof, err := tree.Prove(uint64(i))
		require.NoError(t, err)
		verification, err := r.Verifier.Verify(
			ctx,
			notary.Claim{Ledger: e.ledger, Epoch: 0, Entry: entry, Proof: proof},
		)
		require.NoError(t, err)
		require.Equal(t, notary.VerdictIncluded, verification.Verdict, "account %d", i)
	}

	require.Equal(t, []string{"notary"}, e.payment.Callers)

	response := httptest.NewRecorder()
	r.Handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/metrics", nil))
	body := response.Body.String()
	require.Equal(t, http.StatusOK, response.Code)
	require.Contains(t, body, `kinetix_build_info{service="kinetix-notary-service",version="test"} 1`)
	require.Contains(t, body, "kinetix_notary_last_finalized_epoch 0")
	require.Contains(
		t,
		body,
		`kinetix_grpc_client_calls_total{grpc_code="OK",grpc_method="/payment.v1.LiabilityService/StreamLiabilitySnapshot",peer="payment"} 1`,
	)
	require.Contains(t, body, "kinetix_notary_anchorer_balance_wei")
	require.Contains(t, body, "kinetix_http_requests_total")
	require.NotContains(t, body, "principal-")

	ready := httptest.NewRecorder()
	r.Handler.ServeHTTP(ready, httptest.NewRequest(http.MethodGet, "/health/ready", nil))
	require.Equal(t, http.StatusServiceUnavailable, ready.Code, "not ready before the loop has run")

	live := httptest.NewRecorder()
	r.Handler.ServeHTTP(live, httptest.NewRequest(http.MethodGet, "/health/live", nil))
	require.Equal(t, http.StatusOK, live.Code)
}

func TestADaemonThatDiesMidEpochIsFinishedByTheNextWithoutASecondTransaction(t *testing.T) {
	e := newEstate(t)
	ctx := context.Background()
	first := e.start(t)
	e.tickUntil(t, first, 0, "mined")
	e.chain.finalize(t)
	e.tickUntil(t, first, 0, "finalized")

	require.NoError(t, e.chain.rpc.CallContext(ctx, nil, "evm_increaseTime", 7200))
	require.NoError(t, e.chain.rpc.CallContext(ctx, nil, "evm_mine"))
	e.advance(2 * time.Hour)
	require.NoError(t, e.chain.rpc.CallContext(ctx, nil, "evm_setAutomine", false))
	t.Cleanup(func() { _ = e.chain.rpc.CallContext(context.Background(), nil, "evm_setAutomine", true) })
	require.NoError(t, first.Anchorer.Tick(ctx))
	require.Equal(t, "submitted", e.state(t, 1))

	_, err := e.pool.Exec(ctx, "UPDATE writer_lease SET expires_at = now()")
	require.NoError(t, err)
	second := e.start(t)
	require.NoError(t, second.Anchorer.Tick(ctx))
	require.NoError(t, e.chain.rpc.CallContext(ctx, nil, "evm_setAutomine", true))
	require.NoError(t, e.chain.rpc.CallContext(ctx, nil, "evm_mine"))
	e.tickUntil(t, second, 1, "mined")
	e.chain.finalize(t)
	e.tickUntil(t, second, 1, "finalized")

	var attempts, included int64
	require.NoError(
		t,
		e.pool.QueryRow(ctx, "SELECT count(*), count(*) FILTER (WHERE outcome = 'included') FROM attempts WHERE epoch_number = 1").
			Scan(&attempts, &included),
	)
	require.Equal(t, int64(1), attempts)
	require.Equal(t, int64(1), included)

	err = first.Anchorer.Tick(ctx)
	require.NoError(t, err)

	var holder string
	require.NoError(t, e.pool.QueryRow(ctx, "SELECT holder FROM writer_lease").Scan(&holder))
	require.NotEmpty(t, holder)
}

func (e *estate) proofs(t *testing.T, r running, epoch, subject string) *httptest.ResponseRecorder {
	t.Helper()

	request := httptest.NewRequest(http.MethodGet, "/api/v1/proofs/liabilities/"+epoch, nil)
	request.Header.Set("Authorization", "Bearer "+e.identity.Token(t, subject))

	response := httptest.NewRecorder()
	r.Handler.ServeHTTP(response, request)

	return response
}

func TestAnAccountHolderFetchesTheirProofAndTheChainConfirmsIt(t *testing.T) {
	e := newEstate(t)
	r := e.start(t)

	require.NoError(t, r.Prover.CheckCourt(context.Background()))
	require.Equal(t, http.StatusNotFound, e.proofs(t, r, "latest", "principal-3").Code)

	e.tickUntil(t, r, 0, "mined")
	require.Equal(t, http.StatusConflict, e.proofs(t, r, "0", "principal-3").Code)

	e.chain.finalize(t)
	e.tickUntil(t, r, 0, "finalized")

	response := e.proofs(t, r, "latest", "principal-3")
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())

	path := filepath.Join(t.TempDir(), "proofs.json")
	require.NoError(t, os.WriteFile(path, response.Body.Bytes(), 0o600))

	var out bytes.Buffer

	code := verifycli.Run(context.Background(), path, keccak.NewScheme().Version(), r.Verifier, &out)
	require.Equal(t, verifycli.ExitIncluded, code, out.String())

	var report map[string]string
	require.NoError(t, json.Unmarshal(bytes.TrimSpace(out.Bytes()), &report))
	require.Equal(t, "included", report["verdict"])
	require.Equal(t, "principal-3", report["holder"])
	require.Equal(t, "driver_available", report["kind"])
	require.Equal(t, "4000000", report["amount"])
	require.Equal(t, "attributable", report["statement"])
	require.Equal(t, "0", report["faults"])
	require.Equal(t, e.court.Hex(), report["court"])

	stranger := e.proofs(t, r, "0", "nobody-at-all")
	require.Equal(t, http.StatusOK, stranger.Code)
	require.Contains(t, stranger.Body.String(), `"proofs":[]`)

	require.Equal(t, http.StatusNotFound, e.proofs(t, r, "1", "principal-3").Code)

	forged := httptest.NewRequest(http.MethodGet, "/api/v1/proofs/liabilities/latest", nil)
	forged.Header.Set("Authorization", "Bearer "+e.identity.Stranger(t, "principal-3"))

	refused := httptest.NewRecorder()
	r.Handler.ServeHTTP(refused, forged)
	require.Equal(t, http.StatusUnauthorized, refused.Code)
}
