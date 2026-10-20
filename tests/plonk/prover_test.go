package plonk_test

import (
	"bytes"
	"testing"

	"github.com/consensys/gnark/test/unsafekzg"
	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/plonk"
	"github.com/WildanFrananda/kinetix-notary-service/internal/poseidon"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

func TestAnHonestEpochProvesAndVerifies(t *testing.T) {
	prover := testProver(t, 3)

	proof, err := prover.Prove(honest(t, poseidon.NewScheme(), 11, 3, notarytest.Entries(4, 6)))
	require.NoError(t, err)
	require.NotEmpty(t, proof)
}

func TestAForgedEpochIsNeverProven(t *testing.T) {
	prover := testProver(t, 3)

	forged := honest(t, poseidon.NewScheme(), 11, 3, notarytest.Entries(4, 6))
	forged.LeafCount = 5

	_, err := prover.Prove(forged)
	require.Equal(t, notary.KindInvalidArgument, notary.KindOf(err))
}

func TestTheProverExportsASolidityVerifier(t *testing.T) {
	var out bytes.Buffer

	require.NoError(t, testProver(t, 2).ExportVerifier(&out))
	require.Contains(t, out.String(), "pragma solidity 0.8.37;")
	require.Contains(t, out.String(), "function Verify(")
}

func testProver(t *testing.T, depth int) *plonk.Prover {
	t.Helper()

	system, err := plonk.Compile(depth)
	require.NoError(t, err)

	canonical, lagrange, err := unsafekzg.NewSRS(system)
	require.NoError(t, err)

	prover, err := plonk.NewProver(depth, system, canonical, lagrange)
	require.NoError(t, err)

	return prover
}
