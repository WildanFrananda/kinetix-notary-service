package proving_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/keccak"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/proving"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

var ledgerID = notary.LedgerIdentity{ChainID: 84532, Address: notary.Address{0x4e, 0x86}}

func epoch(t *testing.T, number uint64, state notary.EpochState) (notary.EpochRecord, []notary.Account) {
	t.Helper()

	snapshot := notarytest.Snapshot(t, 1760000000+number, notarytest.Liabilities(t, 15))
	plan, err := notary.PlanEpoch(
		keccak.NewScheme(),
		&notarytest.Salts{Seed: number},
		number,
		notary.Digest{},
		snapshot,
	)
	require.NoError(t, err)

	record := notary.EpochRecord{
		Ledger:   ledgerID,
		Header:   plan.Header,
		Digest:   plan.Digest,
		Currency: "IDR",
		State:    state,
	}

	return record, plan.Accounts
}

func source(t *testing.T, states ...notary.EpochState) *notarytest.ProofSource {
	t.Helper()

	s := &notarytest.ProofSource{Records: map[uint64]notary.EpochRecord{}, Stored: map[uint64][]notary.Account{}}

	for i, state := range states {
		record, accounts := epoch(t, uint64(i), state)
		s.Records[uint64(i)] = record
		s.Stored[uint64(i)] = accounts
	}

	return s
}

func TestAPrincipalGetsAProofForEachOfItsOwnAccountsAndNothingElse(t *testing.T) {
	s := source(t, notary.StateFinalized)
	prover := proving.NewProver(s, keccak.NewScheme())
	holder := notarytest.Holder(t, "principal-0001")

	bundle, err := prover.Proofs(context.Background(), holder, proving.Number(0))

	require.NoError(t, err)
	require.Equal(t, s.Records[0], bundle.Epoch)
	require.Len(t, bundle.Claims, 7)

	anchored := notary.AnchoredEpoch{
		Number:    0,
		Root:      bundle.Epoch.Header.Root,
		Total:     bundle.Epoch.Header.Total,
		LeafCount: bundle.Epoch.Header.LeafCount,
	}
	scheme := keccak.NewScheme()
	kinds := map[notary.AccountKind]bool{}

	for _, claim := range bundle.Claims {
		require.Equal(t, ledgerID, claim.Ledger)
		require.Equal(t, holder, claim.Opening.Holder)
		require.Equal(
			t,
			claim.Entry.Commitment,
			scheme.Commitment(claim.Opening.Holder, claim.Opening.Kind, claim.Opening.Salt),
		)
		require.NoError(t, notary.VerifyInclusion(scheme, anchored, claim.Entry, claim.Proof))

		kinds[claim.Opening.Kind] = true
	}

	require.Len(t, kinds, 7)
}

func TestAPrincipalWithNoAccountGetsAnEmptyAnswerNotAnError(t *testing.T) {
	prover := proving.NewProver(source(t, notary.StateFinalized), keccak.NewScheme())

	bundle, err := prover.Proofs(context.Background(), notarytest.Holder(t, "nobody"), proving.Latest())

	require.NoError(t, err)
	require.Empty(t, bundle.Claims)
}

func TestLatestIsTheNewestFinalEpoch(t *testing.T) {
	s := source(t, notary.StateFinalized, notary.StateFinalized, notary.StateSubmitted)
	prover := proving.NewProver(s, keccak.NewScheme())

	bundle, err := prover.Proofs(context.Background(), notarytest.Holder(t, "principal-0000"), proving.Latest())

	require.NoError(t, err)
	require.Equal(t, uint64(1), bundle.Epoch.Header.Number)
}

func TestAnEpochThatCannotBeProvenSaysWhy(t *testing.T) {
	s := source(
		t,
		notary.StateFinalized,
		notary.StateBuilt,
		notary.StateSubmitted,
		notary.StateMined,
		notary.StateAbandoned,
	)
	s.Records[4] = notary.EpochRecord{
		Ledger:          s.Records[4].Ledger,
		Header:          s.Records[4].Header,
		Digest:          s.Records[4].Digest,
		State:           notary.StateAbandoned,
		AbandonedReason: "EpochConflict",
	}
	prover := proving.NewProver(s, keccak.NewScheme())
	holder := notarytest.Holder(t, "principal-0000")

	cases := map[uint64]notary.Kind{
		1:  notary.KindEpochNotFinal,
		2:  notary.KindEpochNotFinal,
		3:  notary.KindEpochNotFinal,
		4:  notary.KindEpochAbandoned,
		99: notary.KindEpochNotAnchored,
	}
	for number, kind := range cases {
		_, err := prover.Proofs(context.Background(), holder, proving.Number(number))

		require.Equal(t, kind, notary.KindOf(err), "epoch %d", number)
	}

	empty := &notarytest.ProofSource{Records: map[uint64]notary.EpochRecord{}}
	_, err := proving.NewProver(empty, keccak.NewScheme()).Proofs(context.Background(), holder, proving.Latest())
	require.Equal(t, notary.KindEpochNotAnchored, notary.KindOf(err))

	broken := errors.New("connection reset")
	s.Err = notary.Errf(notary.KindStoreUnavailable, "test", broken, "")
	_, err = prover.Proofs(context.Background(), holder, proving.Number(0))
	require.ErrorIs(t, err, broken)
}

func TestStoredAccountsThatDoNotRebuildTheAnchoredRootAreNeverServed(t *testing.T) {
	holder := notarytest.Holder(t, "principal-0000")

	tampering := map[string]func([]notary.Account) []notary.Account{
		"an amount changed": func(a []notary.Account) []notary.Account {
			a[3].Liability.Amount++

			return a
		},
		"an account missing": func(a []notary.Account) []notary.Account { return a[:len(a)-1] },
		"two accounts swapped": func(a []notary.Account) []notary.Account {
			a[0], a[1] = a[1], a[0]

			return a
		},
		"a commitment changed": func(a []notary.Account) []notary.Account {
			a[5].Commitment[0] ^= 0xff

			return a
		},
	}
	for name, tamper := range tampering {
		t.Run(name, func(t *testing.T) {
			s := source(t, notary.StateFinalized)
			s.Stored[0] = tamper(s.Stored[0])

			_, err := proving.NewProver(s, keccak.NewScheme()).Proofs(context.Background(), holder, proving.Number(0))

			require.Equal(t, notary.KindDisagreement, notary.KindOf(err))
		})
	}
}

func TestAProvenEpochIsBuiltOnceAndRecentOnesStayCached(t *testing.T) {
	s := source(
		t,
		notary.StateFinalized,
		notary.StateFinalized,
		notary.StateFinalized,
		notary.StateFinalized,
		notary.StateFinalized,
		notary.StateFinalized,
	)
	prover := proving.NewProver(s, keccak.NewScheme())
	holder := notarytest.Holder(t, "principal-0000")
	ask := func(number uint64) {
		_, err := prover.Proofs(context.Background(), holder, proving.Number(number))
		require.NoError(t, err)
	}

	ask(0)
	ask(0)
	require.Equal(t, 1, s.AccountsReads)

	for number := range uint64(6) {
		ask(number)
	}

	require.Equal(t, 6, s.AccountsReads)

	ask(5)
	require.Equal(t, 6, s.AccountsReads)

	ask(0)
	require.Equal(t, 7, s.AccountsReads)
}

func TestAnEpochReferenceIsANumberOrLatest(t *testing.T) {
	for _, valid := range []string{"latest", "0", "42", "18446744073709551615"} {
		_, err := proving.ParseEpochRef(valid)
		require.NoError(t, err, valid)
	}

	for _, invalid := range []string{"", "Latest", "-1", "1.5", "0x10", "18446744073709551616"} {
		_, err := proving.ParseEpochRef(invalid)
		require.Equal(t, notary.KindInvalidArgument, notary.KindOf(err), invalid)
	}
}
