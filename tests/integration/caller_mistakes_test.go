package integration_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWhatTheNotaryCouldNeverKeepIsTheCallersMistakeNotAnOutage(t *testing.T) {
	e := newEstate(t)
	r := e.start(t)

	for _, epoch := range []string{"9223372036854775808", "18446744073709551615"} {
		response := e.proofs(t, r, epoch, "principal-3")
		require.Equal(t, http.StatusBadRequest, response.Code, response.Body.String())

		var body map[string]string
		require.NoError(t, json.Unmarshal(response.Body.Bytes(), &body))
		require.Equal(t, "invalid_epoch", body["error"])
	}

	for _, path := range []string{
		"/api/v1/proofs/evidence/orders/ORD%001",
		"/api/v1/proofs/evidence/orders/ORD-%FF",
		"/api/v1/proofs/evidence/reviews/review%00",
		"/api/v1/proofs/evidence/reviews/%C3%28",
	} {
		response := e.evidenceProofs(t, r, path, "buyer-1")
		require.Equal(t, http.StatusBadRequest, response.Code, path+" "+response.Body.String())

		var body map[string]string
		require.NoError(t, json.Unmarshal(response.Body.Bytes(), &body))
		require.Equal(t, "invalid_reference", body["error"])
	}
}
