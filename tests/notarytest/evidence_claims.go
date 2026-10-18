package notarytest

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type EvidenceBatch struct {
	Header notary.EvidenceBatchHeader
	Digest notary.Digest
	Claims []notary.EvidenceClaim
}

func BatchOf(
	t *testing.T,
	scheme notary.EvidenceScheme,
	registry notary.RegistryIdentity,
	number uint64,
	previous notary.Digest,
	closedAt uint64,
	records []notary.EvidenceRecord,
) EvidenceBatch {
	t.Helper()

	salts := &Salts{Seed: 41}
	entries := make([]notary.EvidenceEntry, len(records))

	for i, record := range records {
		entry, err := notary.EnterEvidence(scheme, salts, record)
		require.NoError(t, err)

		entries[i] = entry
	}

	closed, err := notary.UnixTimeFromSeconds(closedAt)
	require.NoError(t, err)

	plan, err := notary.PlanEvidenceBatch(scheme, records[0].Stream(), number, previous, closed, entries)
	require.NoError(t, err)

	leaves := make([]notary.Digest, len(entries))
	for i, entry := range entries {
		leaves[i] = entry.Leaf
	}

	tree, err := notary.BuildEvidenceTree(scheme, records[0].Stream(), leaves)
	require.NoError(t, err)

	batch := EvidenceBatch{Header: plan.Header, Digest: plan.Digest, Claims: make([]notary.EvidenceClaim, len(entries))}
	for i, entry := range entries {
		proof, err := tree.Prove(uint64(i))
		require.NoError(t, err)

		batch.Claims[i] = notary.EvidenceClaim{
			Registry: registry,
			Scheme:   scheme.EvidenceVersion(),
			Header:   plan.Header,
			Digest:   plan.Digest,
			Record:   entry.Record,
			Salt:     entry.Salt,
			Proof:    proof,
		}
	}

	return batch
}
