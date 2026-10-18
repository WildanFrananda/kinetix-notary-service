package httpapi

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/WildanFrananda/kinetix-notary-service/internal/evidencing"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/prooffile"
)

const (
	OrderEvidenceRoute  = "/api/v1/proofs/evidence/orders/{order}"
	ReviewEvidenceRoute = "/api/v1/proofs/evidence/reviews/{review}"

	maxReferenceBytes = 255
)

type EvidenceHandler struct {
	auth      Authenticator
	parameter string
	issue     func(ctx context.Context, principal, reference string) (evidencing.Bundle, error)
	log       *slog.Logger
}

func NewOrderEvidenceHandler(auth Authenticator, issuer EvidenceIssuer, log *slog.Logger) *EvidenceHandler {
	return &EvidenceHandler{auth: auth, parameter: "order", issue: issuer.OrderProofs, log: log}
}

func NewReviewEvidenceHandler(auth Authenticator, issuer EvidenceIssuer, log *slog.Logger) *EvidenceHandler {
	return &EvidenceHandler{auth: auth, parameter: "review", issue: issuer.ReviewProofs, log: log}
}

func (h *EvidenceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	holder, err := h.auth.Principal(r.Header.Get("Authorization"))
	if err != nil {
		w.Header().Set("WWW-Authenticate", `Bearer realm="kinetix"`)
		h.fail(w, r, err)

		return
	}

	reference := r.PathValue(h.parameter)
	if reference == "" || len(reference) > maxReferenceBytes {
		h.fail(
			w,
			r,
			notary.Errf(notary.KindInvalidArgument, "httpapi.EvidenceHandler", nil, "the %s is not named", h.parameter),
		)

		return
	}

	bundle, err := h.issue(r.Context(), holder.String(), reference)
	if err != nil {
		h.fail(w, r, err)

		return
	}

	response := evidenceResponse{Proofs: make([]json.RawMessage, 0, len(bundle.Claims)), Pending: bundle.Pending}

	for _, claim := range bundle.Claims {
		document, err := prooffile.MarshalEvidence(claim)
		if err != nil {
			h.fail(w, r, err)

			return
		}

		response.Proofs = append(response.Proofs, document)
	}

	writeJSON(w, http.StatusOK, response)
}

func (h *EvidenceHandler) fail(w http.ResponseWriter, r *http.Request, err error) {
	status, code := evidenceStatusOf(notary.KindOf(err))
	message := notary.MessageOf(err)

	if message == "" {
		message = http.StatusText(status)
	}

	if status >= http.StatusInternalServerError {
		h.log.ErrorContext(
			r.Context(),
			"an evidence proof request failed",
			"kind",
			notary.KindOf(err).String(),
			"err",
			err.Error(),
		)

		message = "the evidence proofs could not be produced"
		if code == "unavailable" {
			message = "the notary cannot read its records right now"
		}
	}

	writeJSON(w, status, map[string]string{"error": code, "message": message})
}

func evidenceStatusOf(kind notary.Kind) (int, string) {
	switch kind {
	case notary.KindUnauthenticated:
		return http.StatusUnauthorized, "unauthenticated"
	case notary.KindInvalidArgument:
		return http.StatusBadRequest, "invalid_reference"
	case notary.KindNoEvidence:
		return http.StatusNotFound, "no_evidence"
	case notary.KindStoreUnavailable:
		return http.StatusServiceUnavailable, "unavailable"
	case notary.KindDisagreement:
		return http.StatusInternalServerError, "disagreement"
	default:
		return http.StatusInternalServerError, "internal"
	}
}
