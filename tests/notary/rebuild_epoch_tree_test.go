package notary_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/keccak"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/poseidon"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

func TestAnEpochOfEitherSchemeRebuildsAndProvesEveryAccount(t *testing.T) {
	snapshot := notarytest.Snapshot(t, 1760000000, notarytest.Liabilities(t, 9))

	one, err := notary.PlanEpoch(
		keccak.NewScheme(),
		&notarytest.Salts{Seed: 1},
		&notarytest.Secrets{},
		nil,
		3,
		notary.Digest{},
		snapshot,
	)
	require.NoError(t, err)

	two, err := planConfidential(t, &notarytest.Blindings{Seed: 1}, snapshot, 3)
	require.NoError(t, err)

	for name, plan := range map[string]notary.EpochPlan{"scheme 1": one, "scheme 2": two} {
		record := recordOf(plan)

		tree, err := notary.RebuildEpochTree(keccak.NewScheme(), poseidon.NewScheme(), record, entriesOf(plan))
		require.NoError(t, err, name)

		for i := range plan.Accounts {
			proof, err := tree.Prove(uint64(i))
			require.NoError(t, err, name)
			require.Equal(t, uint64(i), proof.Index, name)
		}
	}
}

func TestARebuildThatDoesNotReachTheAnchoredStatementIsADisagreement(t *testing.T) {
	snapshot := notarytest.Snapshot(t, 1760000000, notarytest.Liabilities(t, 9))

	plan, err := planConfidential(t, &notarytest.Blindings{Seed: 2}, snapshot, 3)
	require.NoError(t, err)

	cases := map[string]func(*notary.EpochRecord, []notary.Entry){
		"another amount": func(_ *notary.EpochRecord, entries []notary.Entry) { entries[0].Amount++ },
		"another blinding": func(record *notary.EpochRecord, _ []notary.Entry) {
			changed := *record.Confidential
			changed.Blinding[0] ^= 0x01
			record.Confidential = &changed
		},
		"another root": func(record *notary.EpochRecord, _ []notary.Entry) { record.Header.Root[0] ^= 0x01 },
	}

	for name, tamper := range cases {
		record := recordOf(plan)
		entries := entriesOf(plan)
		tamper(&record, entries)

		_, err := notary.RebuildEpochTree(keccak.NewScheme(), poseidon.NewScheme(), record, entries)
		require.Equal(t, notary.KindDisagreement, notary.KindOf(err), name)
	}
}

func TestANotaryWithoutSchemeTwoRefusesASchemeTwoEpoch(t *testing.T) {
	snapshot := notarytest.Snapshot(t, 1760000000, notarytest.Liabilities(t, 3))

	plan, err := planConfidential(t, &notarytest.Blindings{Seed: 3}, snapshot, 1)
	require.NoError(t, err)

	_, err = notary.RebuildEpochTree(keccak.NewScheme(), nil, recordOf(plan), entriesOf(plan))
	require.Equal(t, notary.KindWrongLedger, notary.KindOf(err))
}

func recordOf(plan notary.EpochPlan) notary.EpochRecord {
	return notary.EpochRecord{Header: plan.Header, Digest: plan.Digest, Confidential: plan.Confidential}
}

func entriesOf(plan notary.EpochPlan) []notary.Entry {
	entries := make([]notary.Entry, len(plan.Accounts))
	for i, account := range plan.Accounts {
		entries[i] = notary.Entry{Commitment: account.Commitment, Amount: account.Liability.Amount}
	}

	return entries
}
