package plonk_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/plonk"
	"github.com/WildanFrananda/kinetix-notary-service/internal/poseidon"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

func TestAnEpochProverProvesAPlannedEpoch(t *testing.T) {
	prover := plonk.NewEpochProver(testProver(t, 3))

	proof, err := prover.Prove(context.Background(), proofRequest(t, 3, 6))
	require.NoError(t, err)
	require.NotEmpty(t, proof)
}

func TestAnEpochProverRefusesAnotherDepth(t *testing.T) {
	prover := plonk.NewEpochProver(testProver(t, 3))

	_, err := prover.Prove(context.Background(), proofRequest(t, 2, 3))
	require.Equal(t, notary.KindInvalidArgument, notary.KindOf(err))
}

func TestAnEpochProverRefusesAStatementTheEntriesDoNotMake(t *testing.T) {
	prover := plonk.NewEpochProver(testProver(t, 3))
	request := proofRequest(t, 3, 6)
	request.Entries[0].Amount++

	_, err := prover.Prove(context.Background(), request)
	require.Equal(t, notary.KindInvalidArgument, notary.KindOf(err))
}

func TestAnEpochProverStopsWhenItsContextDoes(t *testing.T) {
	prover := plonk.NewEpochProver(testProver(t, 3))
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := prover.Prove(ctx, proofRequest(t, 3, 6))
	require.Equal(t, notary.KindUnknown, notary.KindOf(err))
}

func proofRequest(t *testing.T, depth, count int) notary.ProofRequest {
	t.Helper()

	scheme := poseidon.NewScheme()
	entries := notarytest.Entries(31, count)

	tree, err := notary.BuildConfidentialTree(scheme, 12, depth, entries)
	require.NoError(t, err)

	blinding := notary.Blinding{0x33}

	header, err := tree.Header(scheme, notary.Digest{}, notarytest.Unix(t, 1_790_000_000), blinding)
	require.NoError(t, err)

	return notary.ProofRequest{
		Number:          12,
		Depth:           depth,
		Entries:         entries,
		TreeRoot:        header.TreeRoot,
		TotalCommitment: header.TotalCommitment,
		Blinding:        blinding,
	}
}
