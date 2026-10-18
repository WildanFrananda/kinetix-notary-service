package integration_test

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/config"
	"github.com/WildanFrananda/kinetix-notary-service/internal/ethereum"
	"github.com/WildanFrananda/kinetix-notary-service/internal/keccak"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/verifycli"
	"github.com/WildanFrananda/kinetix-notary-service/internal/verifying"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

func (e *estate) evidenceState(t *testing.T, stream notary.EvidenceStream, number uint64) string {
	t.Helper()

	var state string

	err := e.pool.QueryRow(
		context.Background(),
		"SELECT coalesce(max(state), '') FROM evidence_batches WHERE stream = $1 AND number = $2",
		int16(stream),
		int64(number),
	).Scan(&state)
	require.NoError(t, err)

	return state
}

func (e *estate) evidenceUntil(t *testing.T, r running, stream notary.EvidenceStream, number uint64, want string) {
	t.Helper()
	require.Eventually(
		t,
		func() bool {
			if err := r.Evidencer.Tick(context.Background()); err != nil {
				t.Logf("evidence tick: %v", err)
			}

			return e.evidenceState(t, stream, number) == want
		},
		30*time.Second,
		200*time.Millisecond,
		"%s batch %d never became %s",
		stream,
		number,
		want,
	)
}

func (e *estate) anchorEvidence(t *testing.T, r running) {
	t.Helper()

	for _, stream := range notary.EvidenceStreams() {
		e.evidenceUntil(t, r, stream, 0, "mined")
		e.chain.finalize(t)
		e.evidenceUntil(t, r, stream, 0, "finalized")
	}
}

func (e *estate) evidenceProofs(t *testing.T, r running, path, subject string) *httptest.ResponseRecorder {
	t.Helper()

	request := httptest.NewRequest(http.MethodGet, path, nil)
	request.Header.Set("Authorization", "Bearer "+e.identity.Token(t, subject))

	response := httptest.NewRecorder()
	r.Handler.ServeHTTP(response, request)

	return response
}

func (e *estate) evidenceVerifier() *verifying.EvidenceVerifier {
	scheme := keccak.NewEvidence()
	registry := ethereum.NewEvidenceRegistry(e.chain.client, e.registry, scheme, 10*time.Second)

	return verifying.NewEvidenceVerifier(registry, scheme)
}

func (e *estate) feedEvidence(t *testing.T) {
	t.Helper()

	for _, event := range notarytest.Lifecycle(t, "ORD-1", "buyer-1", "merchant-1", 100, 3) {
		e.evidence.AddEvents(notarytest.WireEvent(event))
	}

	delivery := notarytest.Deliveries(t, 1)[0]
	delivery.OrderNumber, delivery.DriverPrincipalID = "ORD-1", "courier-1"
	e.evidence.AddDeliveries(notarytest.WireDelivery(delivery))

	review := notarytest.Reviews(t, 1)[0]
	review.ReviewID, review.AuthorPrincipalID, review.MerchantPrincipalID = "review-1", "buyer-1", "merchant-1"
	e.evidence.AddReviews(notarytest.WireReview(review))
}

func TestEvidenceIsAnchoredAndEveryPartyProvesItOnChain(t *testing.T) {
	e := newEstate(t)
	e.feedEvidence(t)
	r := e.start(t)
	require.NoError(t, r.Registry.Check(context.Background()))

	pending := e.evidenceProofs(t, r, "/api/v1/proofs/evidence/orders/ORD-1", "buyer-1")
	require.Equal(t, http.StatusNotFound, pending.Code, "nothing is collected before the first tick")

	e.anchorEvidence(t, r)

	response := e.evidenceProofs(t, r, "/api/v1/proofs/evidence/orders/ORD-1", "buyer-1")
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())

	var body struct {
		Proofs  []json.RawMessage `json:"proofs"`
		Pending int               `json:"pending"`
	}
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &body))
	require.Len(t, body.Proofs, 4)
	require.Zero(t, body.Pending)

	var out bytes.Buffer

	code := verifycli.RunEvidence(
		context.Background(),
		bytes.NewReader(response.Body.Bytes()),
		1,
		e.evidenceVerifier(),
		&out,
	)
	require.Equal(t, verifycli.ExitIncluded, code, out.String())
	require.Equal(t, 4, strings.Count(out.String(), `"verdict":"included"`), out.String())

	forged := strings.Replace(response.Body.String(), "courier-1", "courier-x", 1)

	out.Reset()
	code = verifycli.RunEvidence(context.Background(), strings.NewReader(forged), 1, e.evidenceVerifier(), &out)
	require.Equal(t, verifycli.ExitRejected, code, out.String())

	courier := e.evidenceProofs(t, r, "/api/v1/proofs/evidence/orders/ORD-1", "courier-1")
	require.Equal(t, http.StatusOK, courier.Code)
	require.NoError(t, json.Unmarshal(courier.Body.Bytes(), &body))
	require.Len(t, body.Proofs, 1)

	require.Equal(
		t,
		http.StatusNotFound,
		e.evidenceProofs(t, r, "/api/v1/proofs/evidence/orders/ORD-1", "stranger").Code,
	)

	review := e.evidenceProofs(t, r, "/api/v1/proofs/evidence/reviews/review-1", "merchant-1")
	require.Equal(t, http.StatusOK, review.Code)
	out.Reset()
	code = verifycli.RunEvidence(
		context.Background(),
		bytes.NewReader(review.Body.Bytes()),
		1,
		e.evidenceVerifier(),
		&out,
	)
	require.Equal(t, verifycli.ExitIncluded, code, out.String())

	scrape := httptest.NewRecorder()
	r.Handler.ServeHTTP(scrape, httptest.NewRequest(http.MethodGet, "/metrics", nil))
	require.Contains(t, scrape.Body.String(), `kinetix_notary_evidence_last_finalized_batch{stream="review"} 0`)
	require.Contains(t, scrape.Body.String(), `kinetix_notary_evidence_feed_up{stream="delivery"} 1`)
	require.Contains(t, scrape.Body.String(), `kinetix_notary_evidence_collected_total{stream="order_lifecycle"} 3`)
}

func TestTheRegistryAnswersNotFinalUntilABatchIsFinal(t *testing.T) {
	e := newEstate(t)
	e.feedEvidence(t)
	r := e.start(t)
	registry := ethereum.NewEvidenceRegistry(e.chain.client, e.registry, keccak.NewEvidence(), 10*time.Second)

	identity, err := registry.Identity(context.Background())
	require.NoError(t, err)
	require.Equal(t, e.registry, identity)

	e.evidenceUntil(t, r, notary.StreamOrderLifecycle, 0, "mined")

	_, err = registry.Batch(context.Background(), notary.StreamOrderLifecycle, 0)
	require.Equal(t, notary.KindEpochNotFinal, notary.KindOf(err))

	_, err = registry.Batch(context.Background(), notary.StreamOrderLifecycle, 5)
	require.Equal(t, notary.KindEpochNotAnchored, notary.KindOf(err))

	e.chain.finalize(t)

	batch, err := registry.Batch(context.Background(), notary.StreamOrderLifecycle, 0)
	require.NoError(t, err)
	require.Equal(t, uint64(3), batch.LeafCount)
	require.Equal(t, notary.LogPosition{Transaction: 102, Sequence: 3}, batch.Through)

	notALedger := ethereum.NewEvidenceRegistry(
		e.chain.client,
		notary.RegistryIdentity{ChainID: e.ledger.ChainID, Address: e.ledger.Address},
		keccak.NewEvidence(),
		10*time.Second,
	)
	_, err = notALedger.Identity(context.Background())
	require.Equal(t, notary.KindWrongLedger, notary.KindOf(err))
}

func TestTheDaemonRefusesAnEvidenceKeyThatIsNotItsOwnAnchorer(t *testing.T) {
	e := newEstate(t)
	stranger, _ := e.chain.funded(t)
	e.values[config.EvidenceKeyKey] = "0x" + hex.EncodeToString(crypto.FromECDSA(stranger))

	err := e.start(t).Registry.Check(context.Background())
	require.Equal(t, notary.KindNotAnchorer, notary.KindOf(err))

	e.values[config.EvidenceKeyKey] = e.values[config.AnchorerKeyKey]

	err = e.start(t).Registry.Check(context.Background())
	require.Equal(t, notary.KindInvalidArgument, notary.KindOf(err))

	e.values[config.EvidenceKeyKey] = "0x" + hex.EncodeToString(crypto.FromECDSA(stranger))
	e.values[config.EvidenceRegistryKey] = e.ledger.Address.Hex()

	err = e.start(t).Registry.Check(context.Background())
	require.Equal(t, notary.KindWrongLedger, notary.KindOf(err))
}
