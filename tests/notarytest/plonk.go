package notarytest

import (
	"crypto/sha256"
	"math/big"
	"testing"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark/test/unsafekzg"
	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/plonk"
)

func TestSRSToxicValue() *big.Int {
	seed := sha256.Sum256([]byte("kinetix.notary.plonk.test-srs"))

	return new(big.Int).Mod(new(big.Int).SetBytes(seed[:]), fr.Modulus())
}

func TestProver(t *testing.T, depth int) *plonk.Prover {
	t.Helper()

	system, err := plonk.Compile(depth)
	require.NoError(t, err)

	canonical, lagrange, err := unsafekzg.NewSRS(system, unsafekzg.WithToxicValue(TestSRSToxicValue()))
	require.NoError(t, err)

	prover, err := plonk.NewProver(depth, system, canonical, lagrange)
	require.NoError(t, err)

	return prover
}
