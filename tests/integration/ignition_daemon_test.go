package integration_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/plonk"
)

const ignitionSRSKey = "KINETIX_NOTARY_IGNITION_SRS"

func TestTheDaemonAnchorsAnEpochProvenOnTheIgnitionSRSAtProductionDepth(t *testing.T) {
	path := os.Getenv(ignitionSRSKey)
	if path == "" {
		t.Skip(ignitionSRSKey + " is not set; this test proves with the SRS notary-srs builds from Aztec Ignition")
	}

	e := newEstateOf(t, &confidential{
		depth:    plonk.ProductionDepth,
		verifier: (*chain).deployProductionVerifier,
		srs:      func(*testing.T) (string, string) { return path, plonk.ProductionSRSSHA256 },
	})
	r := e.start(t)
	ctx := context.Background()

	require.NoError(t, r.Prover.CheckCourt(ctx))
	require.NoError(t, r.Directory.Check(ctx))

	started := time.Now()

	e.tickUntilWithin(t, r, 0, "mined", 5*time.Minute)
	t.Logf("epoch 0 proven, its proof landed and its anchor mined in %s", time.Since(started).Round(time.Second))
	e.chain.finalize(t)
	e.tickUntil(t, r, 0, "finalized")

	for _, holder := range []string{"principal-0", "principal-5"} {
		report := e.verifyThroughTheDirectory(t, r, "latest", holder)
		require.Equal(t, "included", report["verdict"], holder)
		require.Equal(t, "attributable", report["statement"], holder)
		require.Equal(t, "0", report["total"], holder)
	}

	scraped := httptest.NewRecorder()
	r.Handler.ServeHTTP(scraped, httptest.NewRequest(http.MethodGet, "/metrics", nil))
	body := scraped.Body.String()
	require.Contains(t, body, "kinetix_notary_account_capacity 256")
	require.Contains(t, body, "kinetix_notary_proving_failures_total 0")
}
