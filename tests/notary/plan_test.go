package notary_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/keccak"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

func TestAPlannedEpochCommitsToEveryAccountAndProvesEachOne(t *testing.T) {
	scheme := keccak.NewScheme()
	snapshot := notarytest.Snapshot(t, 1760000000, notarytest.Liabilities(t, 37))
	previous := notary.Digest{0x42}

	plan, err := notary.PlanEpoch(scheme, &notarytest.Salts{Seed: 1}, 9, previous, snapshot)
	require.NoError(t, err)

	require.Equal(t, uint64(9), plan.Header.Number)
	require.Equal(t, previous, plan.Header.PreviousDigest)
	require.Equal(t, snapshot.Total(), plan.Header.Total)
	require.Equal(t, uint64(37), plan.Header.LeafCount)
	require.Equal(t, snapshot.AsOf(), plan.Header.AsOf)
	require.Equal(t, scheme.EpochDigest(plan.Header), plan.Digest)
	require.Equal(t, "IDR", plan.Currency)
	require.Len(t, plan.Accounts, 37)

	entries := make([]notary.Entry, len(plan.Accounts))
	salts := map[notary.Salt]bool{}

	for i, account := range plan.Accounts {
		require.Equal(t, uint64(i), account.Position)
		require.Equal(
			t,
			scheme.Commitment(account.Liability.Holder, account.Liability.Kind, account.Salt),
			account.Commitment,
		)

		if i > 0 {
			require.Equal(t, -1, bytes.Compare(plan.Accounts[i-1].Commitment[:], account.Commitment[:]))
		}

		require.False(t, salts[account.Salt], "salt reused")
		salts[account.Salt] = true
		entries[i] = notary.Entry{Commitment: account.Commitment, Amount: account.Liability.Amount}
	}

	tree, err := notary.BuildTree(scheme, 9, entries)
	require.NoError(t, err)
	require.Equal(t, plan.Header.Root, tree.Root().Hash)

	anchored := notary.AnchoredEpoch{
		Number:    9,
		Root:      plan.Header.Root,
		Total:     plan.Header.Total,
		LeafCount: plan.Header.LeafCount,
	}
	for i, account := range plan.Accounts {
		proof, err := tree.Prove(uint64(i))
		require.NoError(t, err)
		require.NoError(t, notary.VerifyInclusion(scheme, anchored, entries[i], proof), "account %d", account.Position)
	}
}

func TestTheLeafOrderIsNotTheSourcesOrder(t *testing.T) {
	snapshot := notarytest.Snapshot(t, 1760000000, notarytest.Liabilities(t, 64))

	first, err := notary.PlanEpoch(keccak.NewScheme(), &notarytest.Salts{Seed: 1}, 0, notary.Digest{}, snapshot)
	require.NoError(t, err)
	second, err := notary.PlanEpoch(keccak.NewScheme(), &notarytest.Salts{Seed: 2}, 0, notary.Digest{}, snapshot)
	require.NoError(t, err)

	inSourceOrder := 0
	samePosition := 0

	for i := range first.Accounts {
		if first.Accounts[i].Liability == snapshot.Liabilities()[i] {
			inSourceOrder++
		}

		if first.Accounts[i].Liability == second.Accounts[i].Liability {
			samePosition++
		}
	}

	require.Less(t, inSourceOrder, 8)
	require.Less(t, samePosition, 8)
	require.NotEqual(t, first.Header.Root, second.Header.Root)
	require.Equal(t, first.Header.Total, second.Header.Total)
}

func TestAnEmptySnapshotIsTheCanonicalEmptyEpoch(t *testing.T) {
	scheme := keccak.NewScheme()

	plan, err := notary.PlanEpoch(
		scheme,
		&notarytest.Salts{},
		3,
		notary.Digest{1},
		notarytest.Snapshot(t, 1760000000, nil),
	)

	require.NoError(t, err)
	require.Equal(t, scheme.EmptyHash(), plan.Header.Root)
	require.Zero(t, plan.Header.Total)
	require.Zero(t, plan.Header.LeafCount)
	require.Empty(t, plan.Accounts)
}

func TestASaltSourceThatFailsOrReturnsZerosStopsThePlan(t *testing.T) {
	snapshot := notarytest.Snapshot(t, 1760000000, notarytest.Liabilities(t, 3))
	broken := errors.New("entropy unavailable")

	_, err := notary.PlanEpoch(keccak.NewScheme(), &notarytest.Salts{Err: broken}, 0, notary.Digest{}, snapshot)
	require.ErrorIs(t, err, broken)

	_, err = notary.PlanEpoch(keccak.NewScheme(), &notarytest.Salts{Zero: true}, 0, notary.Digest{}, snapshot)
	require.Equal(t, notary.KindInvalidArgument, notary.KindOf(err))
}
