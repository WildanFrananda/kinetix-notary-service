package httpapi_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/httpapi"
	"github.com/WildanFrananda/kinetix-notary-service/internal/keccak"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/prooffile"
	"github.com/WildanFrananda/kinetix-notary-service/internal/proving"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

type auth struct {
	holder notary.Holder
	err    error
	seen   string
}

func (a *auth) Principal(authorization string) (notary.Holder, error) {
	a.seen = authorization

	return a.holder, a.err
}

type issuer struct {
	bundle proving.Bundle
	err    error
	asked  []proving.EpochRef
}

func (i *issuer) Proofs(_ context.Context, _ notary.Holder, ref proving.EpochRef) (proving.Bundle, error) {
	i.asked = append(i.asked, ref)

	return i.bundle, i.err
}

func router(a *auth, i *issuer) http.Handler {
	pass := func(_ string, next http.Handler) http.Handler { return next }

	return httpapi.NewRouter(
		pass,
		http.NotFoundHandler(),
		func(*http.Request) error { return nil },
		httpapi.NewProofHandler(a, i, keccak.NewScheme(), slog.New(slog.NewTextHandler(io.Discard, nil))),
	)
}

func get(handler http.Handler, path string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(http.MethodGet, path, nil)
	request.Header.Set("Authorization", "Bearer token")

	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	return response
}

func bundle(t *testing.T) proving.Bundle {
	t.Helper()

	scheme := keccak.NewScheme()
	holder := notarytest.Holder(t, "principal-0000")
	liabilities := []notary.Liability{
		{Holder: holder, Kind: notary.AccountMerchantAvailable, Amount: 10},
		{Holder: holder, Kind: notary.AccountMerchantPendingEscrow, Amount: 20},
	}
	plan, err := notary.PlanEpoch(
		scheme,
		&notarytest.Salts{Seed: 3},
		3,
		notary.Digest{},
		notarytest.Snapshot(t, 1760000000, liabilities),
	)
	require.NoError(t, err)

	entries := []notary.Entry{
		{Commitment: plan.Accounts[0].Commitment, Amount: plan.Accounts[0].Liability.Amount},
		{Commitment: plan.Accounts[1].Commitment, Amount: plan.Accounts[1].Liability.Amount},
	}
	tree, err := notary.BuildTree(scheme, 3, entries)
	require.NoError(t, err)

	record := notary.EpochRecord{
		Ledger: notary.LedgerIdentity{ChainID: 84532, Address: notary.Address{0x4e}},
		Header: plan.Header,
		Digest: plan.Digest,
		State:  notary.StateFinalized,
	}
	claims := make([]notary.Claim, 0, 2)

	for i, account := range plan.Accounts {
		proof, err := tree.Prove(uint64(i))
		require.NoError(t, err)

		claims = append(claims, notary.Claim{
			Ledger: record.Ledger,
			Epoch:  3,
			Entry:  entries[i],
			Proof:  proof,
			Opening: &notary.Opening{
				Holder: account.Liability.Holder,
				Kind:   account.Liability.Kind,
				Salt:   account.Salt,
			},
			Statement: &notary.StatementSignature{
				Court:     notarytest.CourtIdentity.Address,
				Signature: notary.Signature{0: byte(i + 1), 64: 28},
			},
		})
	}

	return proving.Bundle{Epoch: record, Claims: claims}
}

func TestTheProofsComeBackAsDocumentsTheVerifierReads(t *testing.T) {
	a := &auth{holder: notarytest.Holder(t, "principal-0000")}
	i := &issuer{bundle: bundle(t)}

	response := get(router(a, i), "/api/v1/proofs/liabilities/latest")

	require.Equal(t, http.StatusOK, response.Code)
	require.Equal(t, "application/json", response.Header().Get("Content-Type"))
	require.Equal(t, "no-store", response.Header().Get("Cache-Control"))
	require.Equal(t, "Bearer token", a.seen)
	require.Equal(t, []proving.EpochRef{proving.Latest()}, i.asked)

	var body struct {
		Epoch  string `json:"epoch"`
		AsOf   string `json:"asOf"`
		Digest string `json:"digest"`
		Ledger struct {
			ChainID string `json:"chainId"`
			Address string `json:"address"`
		} `json:"ledger"`
	}
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &body))
	require.Equal(t, "3", body.Epoch)
	require.Equal(t, "84532", body.Ledger.ChainID)
	require.Equal(t, i.bundle.Epoch.Digest.Hex(), body.Digest)

	claims, err := prooffile.DecodeAll(strings.NewReader(response.Body.String()), 1)
	require.NoError(t, err)
	require.Equal(t, i.bundle.Claims, claims)
}

func TestEveryFailureHasItsOwnStatus(t *testing.T) {
	cases := map[string]struct {
		path    string
		err     error
		auth    error
		status  int
		code    string
		message string
	}{
		"no token": {
			path: "latest",
			auth: notary.Errf(
				notary.KindUnauthenticated,
				"tokens.Authenticator.Principal",
				errors.New("token has invalid claims: token is expired"),
				"the token was refused",
			),
			status:  401,
			code:    "unauthenticated",
			message: "the token was refused",
		},
		"bad epoch": {
			path:    "first",
			status:  400,
			code:    "invalid_epoch",
			message: `"first" is neither an epoch number nor latest`,
		},
		"never built": {
			path:    "7",
			err:     notary.Errf(notary.KindEpochNotAnchored, "proving.Prover.resolve", nil, "epoch 7 was never built"),
			status:  404,
			code:    "epoch_not_found",
			message: "epoch 7 was never built",
		},
		"not final": {
			path: "7",
			err: notary.Errf(
				notary.KindEpochNotFinal,
				"proving.Prover.Proofs",
				nil,
				"epoch 7 is submitted, not final yet",
			),
			status:  409,
			code:    "epoch_not_final",
			message: "epoch 7 is submitted, not final yet",
		},
		"abandoned": {
			path:    "7",
			err:     notary.Errf(notary.KindEpochAbandoned, "proving.Prover.Proofs", nil, "epoch 7 was abandoned"),
			status:  410,
			code:    "epoch_abandoned",
			message: "epoch 7 was abandoned",
		},
		"store down": {
			path:   "7",
			err:    notary.Errf(notary.KindStoreUnavailable, "t", errors.New("dial tcp 10.0.0.9"), ""),
			status: 503,
			code:   "unavailable",
		},
		"no key for the epoch's signer": {
			path:   "7",
			err:    notary.Errf(notary.KindSignerUnavailable, "proving.Prover.Proofs", nil, "no key for 0xab"),
			status: 503,
			code:   "signer_unavailable",
		},
		"chain down": {
			path: "7",
			err: notary.Errf(
				notary.KindLedgerUnavailable,
				"ethereum.Court.SignerOf",
				errors.New("dial tcp 10.0.0.9"),
				"",
			),
			status: 503,
			code:   "chain_unavailable",
		},
		"disagreement": {
			path:   "7",
			err:    notary.Errf(notary.KindDisagreement, "t", nil, "root differs"),
			status: 500,
			code:   "disagreement",
		},
		"anything else": {path: "7", err: errors.New("secret internal detail"), status: 500, code: "internal"},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			response := get(router(&auth{err: c.auth}, &issuer{err: c.err}), "/api/v1/proofs/liabilities/"+c.path)

			require.Equal(t, c.status, response.Code)

			var body map[string]string
			require.NoError(t, json.Unmarshal(response.Body.Bytes(), &body))
			require.Equal(t, c.code, body["error"])

			if c.status == http.StatusUnauthorized {
				require.Contains(t, response.Header().Get("WWW-Authenticate"), "Bearer")
			}

			if c.message != "" {
				require.Equal(t, c.message, body["message"])
			}

			if c.status >= http.StatusInternalServerError {
				require.NotContains(t, body["message"], "10.0.0.9")
				require.NotContains(t, body["message"], "secret internal detail")
				require.NotContains(t, body["message"], "root differs")
			}
		})
	}
}

func TestOnlyGETReachesTheProofs(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/api/v1/proofs/liabilities/latest", nil)
	response := httptest.NewRecorder()

	router(&auth{}, &issuer{}).ServeHTTP(response, request)

	require.Equal(t, http.StatusMethodNotAllowed, response.Code)
}
