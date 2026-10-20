package plonk_test

import (
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/consensys/gnark/test/unsafekzg"
	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/plonk"
	"github.com/WildanFrananda/kinetix-notary-service/internal/poseidon"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

func TestMeasureProvingWhenAsked(t *testing.T) {
	depths := os.Getenv("KINETIX_NOTARY_MEASURE_DEPTHS")
	if depths == "" {
		t.Skip("set KINETIX_NOTARY_MEASURE_DEPTHS, for example 6,7,8, to measure")
	}

	for _, field := range strings.Split(depths, ",") {
		depth, err := strconv.Atoi(strings.TrimSpace(field))
		require.NoError(t, err)

		started := time.Now()
		system, err := plonk.Compile(depth)
		require.NoError(t, err)

		compiled := time.Since(started)

		canonical, lagrange, err := unsafekzg.NewSRS(system)
		require.NoError(t, err)

		started = time.Now()
		prover, err := plonk.NewProver(depth, system, canonical, lagrange)
		require.NoError(t, err)

		setUp := time.Since(started)
		assignment := honest(t, poseidon.NewScheme(), 12, depth, notarytest.Entries(9, 1<<depth))

		started = time.Now()
		_, err = prover.Prove(assignment)
		require.NoError(t, err)

		t.Logf(
			"depth %d: %d constraints, compile %s, setup %s, prove+verify %s, peak RSS so far %s",
			depth,
			system.GetNbConstraints(),
			compiled.Round(time.Millisecond),
			setUp.Round(time.Millisecond),
			time.Since(started).Round(time.Millisecond),
			peakRSS(t),
		)
	}
}

func peakRSS(t *testing.T) string {
	t.Helper()

	status, err := os.ReadFile("/proc/self/status")
	if err != nil {
		return "unknown"
	}

	for _, line := range strings.Split(string(status), "\n") {
		if strings.HasPrefix(line, "VmHWM:") {
			return strings.TrimSpace(strings.TrimPrefix(line, "VmHWM:"))
		}
	}

	return "unknown"
}
