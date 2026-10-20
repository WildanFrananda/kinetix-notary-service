package plonk_test

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/plonk"
)

var productionVerifierPath = filepath.Join("..", "..", "contracts", "src", "confidential", "PlonkVerifier.sol")

func TestTheCommittedVerifierIsThePinnedOne(t *testing.T) {
	source, err := os.ReadFile(productionVerifierPath)
	require.NoError(t, err)

	digest := sha256.Sum256(source)
	require.Equal(t, plonk.ProductionVerifierSHA256, hex.EncodeToString(digest[:]))
}

func TestTheProductionCircuitIsTheOneItsVerifierWasMadeFor(t *testing.T) {
	first, err := plonk.Compile(plonk.ProductionDepth)
	require.NoError(t, err)

	second, err := plonk.Compile(plonk.ProductionDepth)
	require.NoError(t, err)

	digest, err := plonk.CircuitSHA256(first)
	require.NoError(t, err)

	again, err := plonk.CircuitSHA256(second)
	require.NoError(t, err)
	require.Equal(t, digest, again)
	require.Equal(
		t,
		plonk.ProductionCircuitSHA256,
		digest,
		"the circuit changed: derive the verifier again with notary-srs -from and update every production pin",
	)
}
