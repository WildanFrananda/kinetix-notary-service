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

type proofResponse struct {
	Epoch  string            `json:"epoch"`
	AsOf   string            `json:"asOf"`
	Digest string            `json:"digest"`
	Ledger ledgerResponse    `json:"ledger"`
	Proofs []json.RawMessage `json:"proofs"`
}

type ledgerResponse struct {
	ChainID string `json:"chainId"`
	Address string `json:"address"`
}

func NewProofHandler(auth Authenticator, issuer ProofIssuer, scheme notary.Scheme, log *slog.Logger) *ProofHandler {
	return &ProofHandler{auth: auth, issuer: issuer, scheme: scheme.Version(), log: log}
}

func (h *ProofHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	holder, err := h.auth.Principal(r.Header.Get("Authorization"))
	if err != nil {
		w.Header().Set("WWW-Authenticate", `Bearer realm="kinetix"`)
		h.fail(w, r, err)

		return
	}

	ref, err := proving.ParseEpochRef(r.PathValue("epoch"))
	if err != nil {
		h.fail(w, r, err)

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
		document, err := prooffile.Marshal(claim, h.scheme)
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
	message := err.Error()

	if status >= http.StatusInternalServerError {
		h.log.ErrorContext(r.Context(), "a proof request failed", "kind", notary.KindOf(err).String(), "err", message)

		message = "the proofs could not be produced"
		if status == http.StatusServiceUnavailable {
			message = "the notary cannot read its records right now"
		}
	}

	writeJSON(w, status, map[string]string{"error": code, "message": message})
}

func statusOf(kind notary.Kind) (int, string) {
	switch kind {
	case notary.KindUnauthenticated:
		return http.StatusUnauthorized, "unauthenticated"
	case notary.KindInvalidArgument:
		return http.StatusBadRequest, "invalid_epoch"
	case notary.KindEpochNotAnchored:
		return http.StatusNotFound, "epoch_not_found"
	case notary.KindEpochNotFinal:
		return http.StatusConflict, "epoch_not_final"
	case notary.KindEpochAbandoned:
		return http.StatusGone, "epoch_abandoned"
	case notary.KindStoreUnavailable:
		return http.StatusServiceUnavailable, "unavailable"
	case notary.KindDisagreement:
		return http.StatusInternalServerError, "disagreement"
	default:
		return http.StatusInternalServerError, "internal"
	}
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
