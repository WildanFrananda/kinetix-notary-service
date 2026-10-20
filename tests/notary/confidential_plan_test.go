package notary_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/keccak"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/poseidon"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

const planDepth = 6

func TestAConfidentialPlanHidesTheTotalAndProvesEveryAccount(t *testing.T) {
	accounts := keccak.NewScheme()
	scheme := poseidon.NewScheme()
	snapshot := notarytest.Snapshot(t, 1760000000, notarytest.Liabilities(t, 37))
	previous := notary.Digest{0x42}

	plan, err := notary.PlanConfidentialEpoch(
		accounts,
		scheme,
		planDepth,
		&notarytest.Salts{Seed: 1},
		&notarytest.Secrets{},
		&notarytest.Blindings{Seed: 2},
		nil,
		9,
		previous,
		snapshot,
	)
	require.NoError(t, err)

	confidential := plan.Confidential
	require.NotNil(t, confidential)
	require.Equal(t, planDepth, confidential.Depth)
	require.False(t, confidential.Blinding.IsZero())
	require.Equal(t, uint64(0), plan.Header.Total)
	require.Equal(t, uint64(37), plan.Header.LeafCount)
	require.Equal(t, previous, plan.Header.PreviousDigest)
	require.Equal(t, scheme.BindingRoot(confidential.TreeRoot, confidential.TotalCommitment), plan.Header.Root)
	require.Equal(t, scheme.TotalCommitment(snapshot.Total(), confidential.Blinding), confidential.TotalCommitment)
	require.Equal(t, scheme.EpochDigest(plan.Header), plan.Digest)

	entries := make([]notary.Entry, len(plan.Accounts))
	for i, account := range plan.Accounts {
		require.Equal(t, uint64(i), account.Position)
		require.Equal(t, accounts.AccountCommitment(account.ID, account.Salt), account.Commitment)

		if i > 0 {
			require.Equal(t, -1, bytes.Compare(plan.Accounts[i-1].Commitment[:], account.Commitment[:]))
		}

		entries[i] = notary.Entry{Commitment: account.Commitment, Amount: account.Liability.Amount}
	}

	tree, err := notary.BuildConfidentialTree(scheme, 9, planDepth, entries)
	require.NoError(t, err)
	require.Equal(t, confidential.TreeRoot, tree.TreeRoot())

	epoch := notary.ConfidentialEpoch{Number: 9, TreeRoot: confidential.TreeRoot, LeafCount: 37, Depth: planDepth}

	for i := range plan.Accounts {
		proof, err := tree.Prove(uint64(i))
		require.NoError(t, err)
		require.NoError(t, notary.VerifyConfidentialInclusion(scheme, epoch, entries[i], proof), "account %d", i)
	}
}

func TestAConfidentialPlanDrawsAFreshBlindingEachEpoch(t *testing.T) {
	snapshot := notarytest.Snapshot(t, 1760000000, notarytest.Liabilities(t, 3))
	blindings := &notarytest.Blindings{Seed: 3}

	first, err := planConfidential(t, blindings, snapshot, 4)
	require.NoError(t, err)

	second, err := planConfidential(t, blindings, snapshot, 5)
	require.NoError(t, err)

	require.NotEqual(t, first.Confidential.Blinding, second.Confidential.Blinding)
	require.NotEqual(t, first.Confidential.TotalCommitment, second.Confidential.TotalCommitment)
}

func TestAConfidentialPlanRefusesWhatItCannotCommitTo(t *testing.T) {
	small := notarytest.Snapshot(t, 1760000000, notarytest.Liabilities(t, 3))

	_, err := planConfidential(t, &notarytest.Blindings{Zero: true}, small, 1)
	require.Equal(t, notary.KindInvalidArgument, notary.KindOf(err))

	broken := errors.New("entropy unavailable")
	_, err = planConfidential(t, &notarytest.Blindings{Err: broken}, small, 1)
	require.ErrorIs(t, err, broken)

	crowded := notarytest.Snapshot(t, 1760000000, notarytest.Liabilities(t, 1<<planDepth+1))
	_, err = planConfidential(t, &notarytest.Blindings{Seed: 1}, crowded, 1)
	require.Equal(t, notary.KindTooManyLeaves, notary.KindOf(err))
}

func planConfidential(
	t *testing.T,
	blindings notary.BlindingSource,
	snapshot notary.Snapshot,
	number uint64,
) (notary.EpochPlan, error) {
	t.Helper()

	return notary.PlanConfidentialEpoch(
		keccak.NewScheme(),
		poseidon.NewScheme(),
		planDepth,
		&notarytest.Salts{Seed: number},
		&notarytest.Secrets{},
		blindings,
		nil,
		number,
		notary.Digest{},
		snapshot,
	)
}
