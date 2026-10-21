package integration_test

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net"
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
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/WildanFrananda/kinetix-notary-service/internal/anchoring"
	"github.com/WildanFrananda/kinetix-notary-service/internal/answering"
	"github.com/WildanFrananda/kinetix-notary-service/internal/approving"
	"github.com/WildanFrananda/kinetix-notary-service/internal/config"
	"github.com/WildanFrananda/kinetix-notary-service/internal/delivering"
	"github.com/WildanFrananda/kinetix-notary-service/internal/evidencing"
	"github.com/WildanFrananda/kinetix-notary-service/internal/keccak"
	"github.com/WildanFrananda/kinetix-notary-service/internal/metrics"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/postgres"
	commonv1 "github.com/WildanFrananda/kinetix-notary-service/internal/protogen/common/v1"
	notificationv1 "github.com/WildanFrananda/kinetix-notary-service/internal/protogen/notification/v1"
	paymentv1 "github.com/WildanFrananda/kinetix-notary-service/internal/protogen/payment/v1"
	"github.com/WildanFrananda/kinetix-notary-service/internal/proving"
	"github.com/WildanFrananda/kinetix-notary-service/internal/registering"
	"github.com/WildanFrananda/kinetix-notary-service/internal/sponsoring"
	"github.com/WildanFrananda/kinetix-notary-service/internal/verifycli"
	"github.com/WildanFrananda/kinetix-notary-service/internal/verifying"
	"github.com/WildanFrananda/kinetix-notary-service/internal/wiring"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

const (
	databaseURLKey     = "KINETIX_NOTARY_TEST_DATABASE_URL"
	statementVerifyURL = "https://docs.example.test/id/verify/proof-of-liabilities/"
)

type estate struct {
	t           *testing.T
	chain       *chain
	ledger      notary.LedgerIdentity
	court       notary.Address
	withholding notary.Address
	directory   notary.Address
	paymaster   notary.Address
	council     notary.Address
	watch       *approving.Approver
	watchHTTP   http.Handler
	watchKey    *ecdsa.PrivateKey
	watchMember notary.Address
	signer      notary.Address
	pki         *notarytest.PKI
	paymentAddr string
	snapshots   uint64
	dbURL       string
	pool        *pgxpool.Pool
	payment     *notarytest.LiabilityServer
	identity    *notarytest.Identity
	registry    notary.RegistryIdentity
	evidence    *notarytest.EvidenceServers
	notified    *notarytest.NotificationServer
	values      map[string]string
	omit        string
	mu          sync.Mutex
	now         time.Time
	count       int
}

type running struct {
	dig.In
	Evidencer  *evidencing.Evidencer
	Registry   *evidencing.RegistryCheck
	Anchorer   *anchoring.Anchorer
	Handler    http.Handler
	Verifier   *verifying.Verifier
	Prover     *proving.Prover
	Directory  *proving.DirectoryCheck
	Sponsors   *sponsoring.Watcher
	Registrar  *registering.Registrar
	Answerer   *answering.Answerer
	Statements *delivering.Runner
}

type confidential struct {
	depth    uint8
	verifier func(*chain, *testing.T) common.Address
	srs      func(*testing.T) (string, string)
}

func testConfidential() *confidential {
	return &confidential{
		depth:    confidentialDepth,
		verifier: (*chain).deployTestVerifier,
		srs:      func(t *testing.T) (string, string) { return testSRSFile(t, confidentialDepth) },
	}
}

func newEstate(t *testing.T) *estate {
	t.Helper()

	return newEstateOf(t, nil)
}

func newEstateOf(t *testing.T, two *confidential) *estate {
	t.Helper()

	base := os.Getenv(databaseURLKey)
	if base == "" {
		t.Skip(databaseURLKey + " is not set; this test needs Postgres, and CI must set it")
	}

	c := connect(t)
	ctx := context.Background()
	latest, err := c.client.HeaderByNumber(ctx, nil)
	require.NoError(t, err)
	e := &estate{t: t, chain: c, now: time.Unix(int64(latest.Time), 0).UTC(), count: 9}
	watchKey, watch := memberKey(t)

	deployCouncil := c.deployCouncil

	if two != nil {
		e.count = 6
		verifier := two.verifier(c, t)
		e.ledger = c.identity(c.deployConfidentialLedger(t, verifier, two.depth, c.from, 0, notary.Digest{}, 0))
		deployCouncil = c.deployConfidentialCouncil
	} else {
		e.ledger = c.identity(c.deploy(t))
	}

	council := deployCouncil(t, common.Address(e.ledger.Address), notary.Address(c.from), watch)
	c.handOver(t, common.Address(e.ledger.Address), council)
	e.council = notary.Address(council)
	e.watchKey, e.watchMember = watchKey, watch
	statementKey, statementSigner := notarytest.StatementKey(t, "estate-signer")
	e.signer = statementSigner
	e.court = notary.Address(c.deployCourt(t, common.Address(e.ledger.Address), statementSigner))
	e.withholding = notary.Address(c.deployWithholding(t, common.Address(e.court)))
	e.directory = notary.Address(
		c.deployDirectory(t, common.Address(e.ledger.Address), common.Address(e.court), common.Address(e.withholding)),
	)
	e.paymaster = notary.Address(c.sponsorship(t, common.Address(e.court), common.Address(e.withholding)).paymaster)
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
	e.pki = pki
	e.payment = &notarytest.LiabilityServer{Script: e.snapshot}
	address := notarytest.ServeLiabilities(t, pki.Issue(t, "kinetix.local", "payment"), e.payment)
	e.paymentAddr = address
	watchAddress := e.serveWatch(t, pki, address, watchKey)
	evidenceKey, evidenceAnchorer := c.funded(t)
	e.registry = notary.RegistryIdentity{
		ChainID: e.ledger.ChainID,
		Address: notary.Address(c.deployEvidenceRegistry(t, evidenceAnchorer)),
	}
	e.evidence = &notarytest.EvidenceServers{}
	feeds := map[string]string{}

	for _, service := range []string{"order", "matching", "review"} {
		feeds[service] = notarytest.ServeMesh(t, pki.Issue(t, "kinetix.local", service), e.registerFeeds)
	}

	e.notified = &notarytest.NotificationServer{}
	notificationAddress := notarytest.ServeMesh(
		t,
		pki.Issue(t, "kinetix.local", "notification"),
		func(s *grpc.Server) { notificationv1.RegisterNotificationServiceServer(s, e.notified) },
	)

	e.values = map[string]string{
		config.RPCURLKey:           c.url,
		config.ChainIDKey:          strconv.FormatUint(e.ledger.ChainID, 10),
		config.LedgerAddressKey:    e.ledger.Address.Hex(),
		config.CourtAddressKey:     e.court.Hex(),
		config.StatementKeysKey:    "0x" + hex.EncodeToString(crypto.FromECDSA(statementKey)),
		config.WithholdingKey:      e.withholding.Hex(),
		config.WithholdingFromKey:  "0",
		config.LogRangeKey:         "200",
		config.AnswererKeyKey:      "0x" + hex.EncodeToString(crypto.FromECDSA(answererKey)),
		config.CouncilAddressKey:   e.council.Hex(),
		config.DirectoryAddressKey: e.directory.Hex(),
		config.PaymasterAddressKey: e.paymaster.Hex(),
		config.WatchAddressKey:     watchAddress,
		config.RPCTimeoutKey:       "10s",
		config.RPCRateKey:          "200",
		config.DatabaseURLKey:      e.dbURL,
		config.AnchorerKeyKey:      "0x" + hex.EncodeToString(crypto.FromECDSA(c.key)),
		config.PaymentAddressKey:   address,
		config.PKIDirKey:           pki.Issue(t, "kinetix.local", "notary"),
		config.CurrencyKey:         "IDR",
		config.MaxFeeKey:           "1000000000000",
		config.HTTPPortKey:         "18010",
		config.JWKSURLKey:          e.identity.URL,
		config.JWTIssuerKey:        notarytest.IdentityIssuer,
		config.JWTAudienceKey:      notarytest.IdentityAudience,
		config.OrderAddressKey:     feeds["order"],
		config.MatchingAddressKey:  feeds["matching"],
		config.ReviewAddressKey:    feeds["review"],
		config.EvidenceRegistryKey: e.registry.Address.Hex(),
		config.EvidenceKeyKey:      "0x" + hex.EncodeToString(crypto.FromECDSA(evidenceKey)),
		config.NotificationKey:     notificationAddress,
		config.VerifyURLKey:        statementVerifyURL,
		config.StatementDayKey:     "+07:00",
	}

	if two != nil {
		path, sha := two.srs(t)
		e.values[config.SchemeKey] = "2"
		e.values[config.SRSPathKey] = path
		e.values[config.SRSSHA256Key] = sha

		putProverOnPath(t)
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

func (e *estate) serveWatch(t *testing.T, pki *notarytest.PKI, payment string, key *ecdsa.PrivateKey) string {
	t.Helper()

	values := map[string]string{
		config.RPCURLKey:         e.chain.url,
		config.ChainIDKey:        strconv.FormatUint(e.ledger.ChainID, 10),
		config.LedgerAddressKey:  e.ledger.Address.Hex(),
		config.CourtAddressKey:   e.court.Hex(),
		config.RPCTimeoutKey:     "10s",
		config.RPCRateKey:        "200",
		config.CouncilAddressKey: e.council.Hex(),
		config.WatchKeyKey:       "0x" + hex.EncodeToString(crypto.FromECDSA(key)),
		config.PaymentAddressKey: payment,
		config.PKIDirKey:         pki.Issue(t, "kinetix.local", "notary-watch"),
		config.WatchGRPCPortKey:  "50070",
		config.HTTPPortKey:       "18011",
	}
	lookup := func(key string) (string, bool) {
		value, ok := values[key]
		return value, ok
	}
	container, err := wiring.Watch(
		lookup,
		metrics.Build{Service: "kinetix-notary-watch-service", Version: "test"},
		slog.New(slog.NewTextHandler(io.Discard, nil)),
	)
	require.NoError(t, err)

	var (
		approver *approving.Approver
		server   *grpc.Server
		handler  http.Handler
	)

	invoked := container.Invoke(func(a *approving.Approver, s *grpc.Server, h http.Handler) {
		approver, server, handler = a, s, h
	})
	require.NoError(t, invoked)
	require.NoError(t, approver.CheckCouncil(context.Background()))

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	go func() { _ = server.Serve(listener) }()

	t.Cleanup(server.Stop)

	e.watch = approver
	e.watchHTTP = handler

	return listener.Addr().String()
}

func (e *estate) scrapeWatch(t *testing.T) string {
	t.Helper()

	response := httptest.NewRecorder()
	e.watchHTTP.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/metrics", nil))
	require.Equal(t, http.StatusOK, response.Code)

	return response.Body.String()
}

func (e *estate) snapshot() []*paymentv1.StreamLiabilitySnapshotResponse {
	head, err := e.chain.client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return nil
	}

	e.mu.Lock()
	omit := e.omit
	count := e.count
	e.snapshots++
	sequence := e.snapshots
	e.mu.Unlock()

	accounts := make([]*paymentv1.LiabilityAccount, 0, count)

	for i := range count {
		if fmt.Sprintf("principal-%d", i) == omit {
			continue
		}

		accounts = append(accounts, &paymentv1.LiabilityAccount{
			Kind:    paymentv1.LiabilityKind(i%5 + 1),
			Holder:  &paymentv1.LiabilityAccount_PrincipalId{PrincipalId: fmt.Sprintf("principal-%d", i)},
			Balance: &commonv1.Money{AmountMinor: int64(1_000_000 * (i + 1)), Currency: "IDR"},
		})
	}

	parts, _ := notarytest.SnapshotStream(
		e.t,
		notarytest.SnapshotID(e.t, sequence),
		timestamppb.New(time.Unix(int64(head.Time)-5, 0)),
		"IDR",
		accounts,
	)

	return parts
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
	e.tickUntilWithin(t, r, epoch, want, 15*time.Second)
}

func (e *estate) tickUntilWithin(t *testing.T, r running, epoch uint64, want string, within time.Duration) {
	t.Helper()
	require.Eventually(
		t,
		func() bool {
			_ = r.Anchorer.Tick(context.Background())
			return e.state(t, epoch) == want
		},
		within,
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

	require.NoError(t, r.Directory.Check(ctx))
	require.NoError(t, r.Sponsors.Check(ctx))
	r.Sponsors.Look(ctx)

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
	require.Contains(
		t,
		body,
		`kinetix_grpc_client_calls_total{grpc_code="OK",grpc_method="/notary.v1.EpochWatchService/ApproveEpoch",peer="notary-watch"} 1`,
	)
	require.Contains(t, body, "kinetix_notary_anchorer_balance_wei")
	require.Contains(t, body, "kinetix_notary_paymaster_deposit_wei 1e+18")
	require.Contains(t, body, "kinetix_notary_paymaster_stake_wei 1e+18")
	require.Contains(t, body, "kinetix_notary_paymaster_current 1")
	require.Contains(t, body, "kinetix_http_requests_total")
	require.NotContains(t, body, "principal-")

	watched := e.scrapeWatch(t)
	require.Contains(t, watched, `kinetix_build_info{service="kinetix-notary-watch-service",version="test"} 1`)
	require.Contains(
		t,
		watched,
		`kinetix_grpc_server_calls_total{grpc_code="OK",grpc_method="/notary.v1.EpochWatchService/ApproveEpoch"} 1`,
	)
	require.Contains(t, watched, `kinetix_notary_watch_approvals_total{outcome="approved"} 1`)

	ready := httptest.NewRecorder()
	r.Handler.ServeHTTP(ready, httptest.NewRequest(http.MethodGet, "/health/ready", nil))
	require.Equal(t, http.StatusServiceUnavailable, ready.Code, "not ready before the loop has run")

	live := httptest.NewRecorder()
	r.Handler.ServeHTTP(live, httptest.NewRequest(http.MethodGet, "/health/live", nil))
	require.Equal(t, http.StatusOK, live.Code)
}

func TestTheDaemonRefusesADirectoryWhoseCurrentLedgerIsNotItsOwn(t *testing.T) {
	e := newEstate(t)
	c := e.chain
	other := c.deploy(t)
	_, signer := notarytest.StatementKey(t, "other-signer")
	otherCourt := c.deployCourt(t, other, signer)
	otherWithholding := c.deployWithholding(t, otherCourt)
	directory := c.deployDirectory(t, other, otherCourt, otherWithholding)
	e.values[config.DirectoryAddressKey] = notary.Address(directory).Hex()
	r := e.start(t)

	err := r.Directory.Check(context.Background())

	require.Equal(t, notary.KindWrongDirectory, notary.KindOf(err), "%v", err)
	require.ErrorContains(t, err, e.ledger.Address.Hex())
}

func TestTheDaemonRefusesAPaymasterThatSponsorsAnotherNotarysCourts(t *testing.T) {
	e := newEstate(t)
	c := e.chain
	other := c.deploy(t)
	_, signer := notarytest.StatementKey(t, "other-signer")
	otherCourt := c.deployCourt(t, other, signer)
	otherWithholding := c.deployWithholding(t, otherCourt)
	otherPaymaster := c.sponsorship(t, otherCourt, otherWithholding).paymaster
	e.values[config.PaymasterAddressKey] = notary.Address(otherPaymaster).Hex()
	r := e.start(t)
	ctx := context.Background()

	err := r.Sponsors.Check(ctx)

	require.Equal(t, notary.KindNotSponsored, notary.KindOf(err), "%v", err)
	require.ErrorContains(t, err, e.withholding.Hex())

	r.Sponsors.Look(ctx)
	require.Contains(t, e.metric(t, r), "kinetix_notary_paymaster_current 0")
}

func TestAnEpochTheWatchCannotMatchToPaymentsRecordIsAbandoned(t *testing.T) {
	e := newEstate(t)
	r := e.start(t)
	e.payment.Tamper = func(record *paymentv1.GetLiabilitySnapshotRecordResponse) {
		record.Digest[0] ^= 0xff
	}

	err := r.Anchorer.Tick(context.Background())

	require.Equal(t, notary.KindApprovalRefused, notary.KindOf(err), err)
	require.Equal(t, "abandoned", e.state(t, 0))

	var reason string
	require.NoError(
		t,
		e.pool.QueryRow(context.Background(), "SELECT abandoned_reason FROM epochs WHERE number = 0").Scan(&reason),
	)
	require.Contains(t, reason, "payment recorded snapshot")

	var attempts int64
	require.NoError(t, e.pool.QueryRow(context.Background(), "SELECT count(*) FROM attempts").Scan(&attempts))
	require.Zero(t, attempts)

	watched := e.scrapeWatch(t)
	require.Contains(
		t,
		watched,
		`kinetix_grpc_server_calls_total{grpc_code="FailedPrecondition",grpc_method="/notary.v1.EpochWatchService/ApproveEpoch"} 1`,
	)
	require.Contains(t, watched, `kinetix_notary_watch_approvals_total{outcome="refused"} 1`)
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
