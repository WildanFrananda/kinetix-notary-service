package httpapi_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/evidencing"
	"github.com/WildanFrananda/kinetix-notary-service/internal/httpapi"
	"github.com/WildanFrananda/kinetix-notary-service/internal/keccak"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/prooffile"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

type evidenceIssuer struct {
	bundle evidencing.Bundle
	err    error
	asked  []string
}

func (e *evidenceIssuer) OrderProofs(_ context.Context, principal, order string) (evidencing.Bundle, error) {
	e.asked = append(e.asked, "order:"+principal+":"+order)

	return e.bundle, e.err
}

func (e *evidenceIssuer) ReviewProofs(_ context.Context, principal, review string) (evidencing.Bundle, error) {
	e.asked = append(e.asked, "review:"+principal+":"+review)

	return e.bundle, e.err
}

func evidenceRouter(a *auth, e *evidenceIssuer) http.Handler {
	pass := func(_ string, next http.Handler) http.Handler { return next }
	discard := slog.New(slog.NewTextHandler(io.Discard, nil))

	return httpapi.NewRouter(
		pass,
		http.NotFoundHandler(),
		func(*http.Request) error { return nil },
		httpapi.NewProofHandler(a, &issuer{}, keccak.NewScheme(), discard),
		httpapi.NewRegistrationHandler(a, &registrant{}, discard),
		httpapi.NewOrderEvidenceHandler(a, e, discard),
		httpapi.NewReviewEvidenceHandler(a, e, discard),
	)
}

func evidenceBundle(t *testing.T) evidencing.Bundle {
	t.Helper()

	deliveries := notarytest.Deliveries(t, 3)
	batch := notarytest.BatchOf(
		t,
		keccak.NewEvidence(),
		notary.RegistryIdentity{ChainID: 84532, Address: notary.Address{0xe0}},
		4,
		notary.Digest{0x01},
		1_760_000_900,
		notarytest.Records(deliveries),
	)

	return evidencing.Bundle{Claims: batch.Claims[:2], Pending: 2}
}

func TestEvidenceGoesOnlyToAnAuthenticatedCaller(t *testing.T) {
	a := &auth{err: notary.Errf(notary.KindUnauthenticated, "test", nil, "no token")}
	e := &evidenceIssuer{}

	response := get(evidenceRouter(a, e), "/api/v1/proofs/evidence/orders/ORD-1")

	require.Equal(t, http.StatusUnauthorized, response.Code)
	require.NotEmpty(t, response.Header().Get("WWW-Authenticate"))
	require.Empty(t, e.asked)
}

func TestAnOrdersEvidenceIsServedAsDocumentsTheVerifierReads(t *testing.T) {
	a := &auth{holder: notarytest.Holder(t, "principal-0001")}
	bundle := evidenceBundle(t)
	e := &evidenceIssuer{bundle: bundle}

	response := get(evidenceRouter(a, e), "/api/v1/proofs/evidence/orders/ORD-1")

	require.Equal(t, http.StatusOK, response.Code)
	require.Equal(t, []string{"order:principal-0001:ORD-1"}, e.asked)

	var body struct {
		Proofs  []json.RawMessage `json:"proofs"`
		Pending int               `json:"pending"`
	}
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &body))
	require.Equal(t, 2, body.Pending)
	require.Len(t, body.Proofs, 2)

	for i, document := range body.Proofs {
		claim, err := prooffile.DecodeEvidence(bytes.NewReader(document), keccak.NewEvidence().EvidenceVersion())
		require.NoError(t, err)
		require.Equal(t, bundle.Claims[i], claim)
	}

	claims, err := prooffile.DecodeEvidenceAll(bytes.NewReader(response.Body.Bytes()), 1)
	require.NoError(t, err)
	require.Len(t, claims, 2)
	require.True(t, prooffile.IsEvidence(response.Body.Bytes()))
}

func TestAReviewsEvidenceIsAskedByItsID(t *testing.T) {
	a := &auth{holder: notarytest.Holder(t, "principal-0002")}
	e := &evidenceIssuer{bundle: evidencing.Bundle{Pending: 1}}

	response := get(evidenceRouter(a, e), "/api/v1/proofs/evidence/reviews/review-7")

	require.Equal(t, http.StatusOK, response.Code)
	require.Equal(t, []string{"review:principal-0002:review-7"}, e.asked)
	require.JSONEq(t, `{"proofs":[],"pending":1}`, response.Body.String())
}

func TestEvidenceRequestsFailWithoutSayingWhatIsKept(t *testing.T) {
	a := &auth{holder: notarytest.Holder(t, "principal-0003")}

	for _, c := range []struct {
		err    error
		status int
		code   string
	}{
		{notary.Errf(notary.KindNoEvidence, "test", nil, "none"), http.StatusNotFound, "no_evidence"},
		{notary.Errf(notary.KindStoreUnavailable, "test", nil, "db down"), http.StatusServiceUnavailable, "unavailable"},
		{notary.Errf(notary.KindDisagreement, "test", nil, "tree"), http.StatusInternalServerError, "disagreement"},
		{notary.Errf(notary.KindInvalidArgument, "test", nil, "tree"), http.StatusInternalServerError, "internal"},
	} {
		response := get(evidenceRouter(a, &evidenceIssuer{err: c.err}), "/api/v1/proofs/evidence/orders/ORD-1")
		require.Equal(t, c.status, response.Code)

		var body map[string]string
		require.NoError(t, json.Unmarshal(response.Body.Bytes(), &body))
		require.Equal(t, c.code, body["error"])

		if c.status >= http.StatusInternalServerError {
			require.NotContains(t, body["message"], "tree")
		}
	}
}

func TestAnOverlongReferenceIsRefusedBeforeAnyLookup(t *testing.T) {
	a := &auth{holder: notarytest.Holder(t, "principal-0004")}
	e := &evidenceIssuer{}

	response := get(evidenceRouter(a, e), "/api/v1/proofs/evidence/orders/"+strings.Repeat("A", 256))

	require.Equal(t, http.StatusBadRequest, response.Code)
	require.Empty(t, e.asked)
}

func TestAProofTheNotaryCannotEncodeIsItsOwnFaultNotTheCallers(t *testing.T) {
	a := &auth{holder: notarytest.Holder(t, "principal-0005")}
	claims := evidenceBundle(t).Claims
	claims[1].Record = strangeRecord{}
	e := &evidenceIssuer{bundle: evidencing.Bundle{Claims: claims}}

	response := get(evidenceRouter(a, e), "/api/v1/proofs/evidence/orders/ORD-1")

	require.Equal(t, http.StatusInternalServerError, response.Code)
	require.JSONEq(
		t,
		`{"error":"internal","message":"the evidence proofs could not be produced"}`,
		response.Body.String(),
	)
}

func TestAnIdentityOutageIsNotARefusedToken(t *testing.T) {
	a := &auth{err: notary.Errf(notary.KindIdentityUnavailable, "test", nil, "identity's signing keys cannot be read")}
	e := &evidenceIssuer{}

	response := get(evidenceRouter(a, e), "/api/v1/proofs/evidence/reviews/review-7")

	require.Equal(t, http.StatusServiceUnavailable, response.Code)
	require.Empty(t, response.Header().Get("WWW-Authenticate"))
	require.JSONEq(
		t,
		`{"error":"identity_unavailable","message":"the notary cannot check tokens right now"}`,
		response.Body.String(),
	)
	require.Empty(t, e.asked)
}
