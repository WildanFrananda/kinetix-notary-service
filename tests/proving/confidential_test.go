package proving_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/keccak"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/poseidon"
	"github.com/WildanFrananda/kinetix-notary-service/internal/proving"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

func confidentialSource(t *testing.T, depth int) *notarytest.ProofSource {
	t.Helper()

	snapshot := notarytest.Snapshot(t, 1760000000, notarytest.Liabilities(t, 15))
	plan, err := notary.PlanConfidentialEpoch(
		keccak.NewScheme(),
		poseidon.NewScheme(),
		depth,
		&notarytest.Salts{Seed: 4},
		&notarytest.Secrets{},
		&notarytest.Blindings{Seed: 4},
		nil,
		0,
		notary.Digest{},
		snapshot,
	)
	require.NoError(t, err)

	plan.Confidential.Proof = []byte{0x01}
	record := notary.EpochRecord{
		Ledger:       ledgerID,
		Header:       plan.Header,
		Digest:       plan.Digest,
		Currency:     "IDR",
		State:        notary.StateFinalized,
		Confidential: plan.Confidential,
	}

	return &notarytest.ProofSource{
		Records: map[uint64]notary.EpochRecord{0: record},
		Stored:  map[uint64][]notary.Account{0: plan.Accounts},
	}
}

func TestASchemeTwoEpochGivesEachHolderAProofWithNoSums(t *testing.T) {
	s := confidentialSource(t, 5)
	holder := notarytest.Holder(t, "principal-0001")

	bundle, err := newProver(t, s).Proofs(context.Background(), holder, proving.Number(0))
	require.NoError(t, err)
	require.Len(t, bundle.Claims, 7)

	record := s.Records[0]
	epoch := notary.ConfidentialEpoch{
		Number:    0,
		TreeRoot:  record.Confidential.TreeRoot,
		LeafCount: record.Header.LeafCount,
		Depth:     record.Confidential.Depth,
	}

	for _, claim := range bundle.Claims {
		require.Len(t, claim.Proof.Siblings, 5)

		for _, sibling := range claim.Proof.Siblings {
			require.Zero(t, sibling.Sum)
		}

		require.NoError(t, notary.VerifyConfidentialInclusion(poseidon.NewScheme(), epoch, claim.Entry, claim.Proof))
		require.NotNil(t, claim.Statement)
	}
}

func TestASchemeTwoEpochWhoseAccountsDoNotRebuildItIsADisagreement(t *testing.T) {
	s := confidentialSource(t, 5)
	s.Stored[0][3].Liability.Amount++

	_, err := newProver(t, s).Proofs(context.Background(), notarytest.Holder(t, "principal-0001"), proving.Number(0))
	require.Equal(t, notary.KindDisagreement, notary.KindOf(err))
}
