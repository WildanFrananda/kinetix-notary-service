package httpapi_test

import (
	"encoding/json"
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
	"github.com/WildanFrananda/kinetix-notary-service/internal/registering"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

const (
	registrationHolder    = "0x00000000000000000000000000000000000000a1"
	registrationSignature = "0x" + "11111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111" + "1b"
)

func registrationRouter(a *auth, r *registrant) http.Handler {
	pass := func(_ string, next http.Handler) http.Handler { return next }
	discard := slog.New(slog.NewTextHandler(io.Discard, nil))

	return httpapi.NewRouter(
		pass,
		http.NotFoundHandler(),
		func(*http.Request) error { return nil },
		httpapi.NewProofHandler(a, &issuer{}, keccak.NewScheme(), discard),
		httpapi.NewRegistrationHandler(a, r, discard),
	)
}

func post(handler http.Handler, body string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(http.MethodPost, httpapi.RegistrationsRoute, strings.NewReader(body))
	request.Header.Set("Authorization", "Bearer token")

	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	return response
}

func validBody() string {
	return `{"holder":"` + registrationHolder + `","issuedAt":"1760000000","signature":"` + registrationSignature + `"}`
}

func TestARegistrationIsIssuedForTheTokensPrincipal(t *testing.T) {
	principal := notarytest.Holder(t, "principal-0007")
	court := notary.CourtIdentity{ChainID: 84532, Address: notary.Address{0xd0}}
	issued := notary.IssuedRegistration{
		Registration: notary.Registration{Holder: notary.Address{0xa1}, Account: notary.Digest{0x0a}, FromEpoch: 21},
		Kind:         notary.AccountCustomerBalance,
		Digest:       notary.Digest{0xdd},
		Signature:    notary.Signature{0x5e},
	}
	r := &registrant{
		issued: registering.Issued{Court: court, Epoch: 21, Registrations: []notary.IssuedRegistration{issued}},
	}

	response := post(registrationRouter(&auth{holder: principal}, r), validBody())

	require.Equal(t, http.StatusOK, response.Code)
	require.Equal(t, principal, r.principal)
	require.Equal(t, principal, r.request.Principal)
	require.Equal(t, registrationHolder, r.request.Holder.Hex())
	require.Equal(t, uint64(1760000000), r.request.IssuedAt.Seconds())
	require.Equal(t, registrationSignature, r.signature.Hex())

	var body struct {
		Court struct {
			ChainID string `json:"chainId"`
			Address string `json:"address"`
		} `json:"court"`
		Epoch         string              `json:"epoch"`
		Registrations []map[string]string `json:"registrations"`
	}
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &body))
	require.Equal(t, "84532", body.Court.ChainID)
	require.Equal(t, court.Address.Hex(), body.Court.Address)
	require.Equal(t, "21", body.Epoch)
	require.Equal(
		t,
		[]map[string]string{{
			"holder":    issued.Registration.Holder.Hex(),
			"account":   issued.Registration.Account.Hex(),
			"fromEpoch": "21",
			"kind":      "customer_balance",
			"digest":    issued.Digest.Hex(),
			"signature": issued.Signature.Hex(),
		}},
		body.Registrations,
	)
}

func TestARegistrationWithoutAValidTokenNeverReachesTheRegistrar(t *testing.T) {
	r := &registrant{}
	a := &auth{err: notary.Errf(notary.KindUnauthenticated, "test", nil, "token expired")}

	response := post(registrationRouter(a, r), validBody())

	require.Equal(t, http.StatusUnauthorized, response.Code)
	require.Equal(t, `Bearer realm="kinetix"`, response.Header().Get("WWW-Authenticate"))
	require.Zero(t, r.calls)
}

func TestABrokenRegistrationBodyIsABadRequest(t *testing.T) {
	bodies := map[string]string{
		"not json":         `{`,
		"an unknown field": `{"holder":"` + registrationHolder + `","issuedAt":"1","signature":"` + registrationSignature + `","principal":"x"}`,
		"a bad holder":     `{"holder":"0x12","issuedAt":"1760000000","signature":"` + registrationSignature + `"}`,
		"a bad time":       `{"holder":"` + registrationHolder + `","issuedAt":"soon","signature":"` + registrationSignature + `"}`,
		"a bad signature":  `{"holder":"` + registrationHolder + `","issuedAt":"1760000000","signature":"0x1b"}`,
		"two documents":    validBody() + validBody(),
		"too large":        `{"holder":"` + strings.Repeat(" ", 5<<10) + `"}`,
	}
	for name, body := range bodies {
		t.Run(name, func(t *testing.T) {
			r := &registrant{}

			response := post(registrationRouter(&auth{holder: notarytest.Holder(t, "p")}, r), body)

			require.Equal(t, http.StatusBadRequest, response.Code)
			require.Contains(t, response.Body.String(), `"error":"invalid_request"`)
			require.Zero(t, r.calls)
		})
	}
}

func TestEachRefusalHasItsOwnStatusAndNoInternalsLeak(t *testing.T) {
	cases := []struct {
		kind    notary.Kind
		status  int
		code    string
		message string
	}{
		{notary.KindRegistrationRefused, http.StatusForbidden, "registration_refused", "detail for the holder"},
		{notary.KindNoAccount, http.StatusNotFound, "no_account", "detail for the holder"},
		{
			notary.KindRegistrationNotOpen,
			http.StatusServiceUnavailable,
			"registration_not_open",
			"detail for the holder",
		},
		{
			notary.KindStoreUnavailable,
			http.StatusServiceUnavailable,
			"unavailable",
			"the notary cannot read its records right now",
		},
		{
			notary.KindSignerUnavailable,
			http.StatusServiceUnavailable,
			"signer_unavailable",
			"the notary cannot sign a registration right now",
		},
		{
			notary.KindLedgerUnavailable,
			http.StatusServiceUnavailable,
			"chain_unavailable",
			"the notary cannot reach the chain right now",
		},
		{
			notary.KindDisagreement,
			http.StatusInternalServerError,
			"disagreement",
			"the registration could not be issued",
		},
	}
	for _, c := range cases {
		t.Run(c.code, func(t *testing.T) {
			r := &registrant{err: notary.Errf(c.kind, "registering.secret.op", nil, "detail for the holder")}

			response := post(registrationRouter(&auth{holder: notarytest.Holder(t, "p")}, r), validBody())

			require.Equal(t, c.status, response.Code)
			require.JSONEq(t, `{"error":"`+c.code+`","message":"`+c.message+`"}`, response.Body.String())
			require.NotContains(t, response.Body.String(), "registering.secret.op")
		})
	}
}

func TestRegistrationsAreOnlyPosted(t *testing.T) {
	r := &registrant{}
	request := httptest.NewRequest(http.MethodGet, httpapi.RegistrationsRoute, nil)
	response := httptest.NewRecorder()

	registrationRouter(&auth{}, r).ServeHTTP(response, request)

	require.Equal(t, http.StatusNotFound, response.Code)
	require.Zero(t, r.calls)
}
