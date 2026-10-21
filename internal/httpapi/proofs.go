package httpapi

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/prooffile"
	"github.com/WildanFrananda/kinetix-notary-service/internal/proving"
)

const ProofsRoute = "/api/v1/proofs/liabilities/{epoch}"

type ProofHandler struct {
	auth   Authenticator
	issuer ProofIssuer
	scheme uint16
	log    *slog.Logger
}

func NewProofHandler(auth Authenticator, issuer ProofIssuer, scheme notary.Scheme, log *slog.Logger) *ProofHandler {
	return &ProofHandler{auth: auth, issuer: issuer, scheme: scheme.Version(), log: log}
}

func (h *ProofHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	holder, err := h.auth.Principal(r.Header.Get("Authorization"))
	if err != nil {
		h.fail(w, r, err)

		return
	}

	ref, err := proving.ParseEpochRef(r.PathValue("epoch"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_epoch", notary.MessageOf(err))

		return
	}

	bundle, err := h.issuer.Proofs(r.Context(), holder, ref)
	if err != nil {
		h.fail(w, r, err)

		return
	}

	response := proofResponse{
		Epoch:  strconv.FormatUint(bundle.Epoch.Header.Number, 10),
		AsOf:   bundle.Epoch.Header.AsOf.Time().Format(time.RFC3339),
		Digest: bundle.Epoch.Digest.Hex(),
		Ledger: ledgerResponse{
			ChainID: strconv.FormatUint(bundle.Epoch.Ledger.ChainID, 10),
			Address: bundle.Epoch.Ledger.Address.Hex(),
		},
		Proofs: make([]json.RawMessage, 0, len(bundle.Claims)),
	}

	for _, claim := range bundle.Claims {
		scheme := h.scheme
		if claim.Scheme != 0 {
			scheme = claim.Scheme
		}

		document, err := prooffile.Marshal(claim, scheme)
		if err != nil {
			h.fail(w, r, err)

			return
		}

		response.Proofs = append(response.Proofs, document)
	}

	writeJSON(w, http.StatusOK, response)
}

func (h *ProofHandler) fail(w http.ResponseWriter, r *http.Request, err error) {
	status, code := statusOf(notary.KindOf(err))
	message := notary.MessageOf(err)

	if message == "" {
		message = http.StatusText(status)
	}

	if status >= http.StatusInternalServerError {
		h.log.ErrorContext(
			r.Context(),
			"a proof request failed",
			"kind",
			notary.KindOf(err).String(),
			"err",
			err.Error(),
		)

		message = "the proofs could not be produced"

		switch code {
		case "unavailable":
			message = "the notary cannot read its records right now"
		case "signer_unavailable":
			message = "the notary cannot sign a statement for this epoch right now"
		case "chain_unavailable":
			message = "the notary cannot reach the chain right now"
		case "identity_unavailable":
			message = identityUnavailableMessage
		}
	}

	writeError(w, status, code, message)
}

func statusOf(kind notary.Kind) (int, string) {
	switch kind {
	case notary.KindUnauthenticated:
		return http.StatusUnauthorized, "unauthenticated"
	case notary.KindIdentityUnavailable:
		return http.StatusServiceUnavailable, "identity_unavailable"
	case notary.KindEpochNotAnchored:
		return http.StatusNotFound, "epoch_not_found"
	case notary.KindEpochNotFinal:
		return http.StatusConflict, "epoch_not_final"
	case notary.KindEpochAbandoned:
		return http.StatusGone, "epoch_abandoned"
	case notary.KindStoreUnavailable:
		return http.StatusServiceUnavailable, "unavailable"
	case notary.KindSignerUnavailable:
		return http.StatusServiceUnavailable, "signer_unavailable"
	case notary.KindLedgerUnavailable:
		return http.StatusServiceUnavailable, "chain_unavailable"
	case notary.KindDisagreement:
		return http.StatusInternalServerError, "disagreement"
	default:
		return http.StatusInternalServerError, "internal"
	}
}
