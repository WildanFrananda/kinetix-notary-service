package verifying_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/keccak"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/poseidon"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

const confidentialDepth = 4

type confidentialCase struct {
	directory *notarytest.Directory
	readers   *notarytest.Readers
	ledger    *notarytest.Ledger
	claim     notary.Claim
}

func confidentialCaseOf(t *testing.T, confirmed bool) confidentialCase {
	t.Helper()

	scheme := poseidon.NewScheme()
	entries := notarytest.Entries(17, 6)

	tree, err := notary.BuildConfidentialTree(scheme, 9, confidentialDepth, entries)
	require.NoError(t, err)

	header, err := tree.Header(scheme, notary.Digest{0x01}, notarytest.Unix(t, 1760000000), notary.Blinding{0x02})
	require.NoError(t, err)

	proof, err := tree.Prove(2)
	require.NoError(t, err)

	ledger := &notarytest.Ledger{
		IdentityValue: notarytest.LedgerIdentity,
		Epochs: map[uint64]notary.AnchoredEpoch{9: {
			Number:     9,
			Root:       header.Header.Root,
			LeafCount:  6,
			AsOf:       notarytest.Unix(t, 1760000000),
			AnchoredAt: notarytest.Unix(t, 1760000060),
		}},
		Confirmed: confirmed,
	}
	generation := generationOf(notarytest.LedgerIdentity, notarytest.CourtIdentity.Address, 0)
	generation.Scheme = 2

	return confidentialCase{
		directory: &notarytest.Directory{
			IdentityValue: notarytest.DirectoryIdentity,
			Generations:   []notary.Generation{generation},
		},
		readers: &notarytest.Readers{
			Ledgers: map[notary.Address]notary.Ledger{notarytest.LedgerIdentity.Address: ledger},
			Confidentials: map[notary.Address]notary.ConfidentialLedger{
				notarytest.LedgerIdentity.Address: &notarytest.ConfidentialLedger{
					Statements: map[uint64]notary.ConfidentialEpoch{9: {
						Number:    9,
						TreeRoot:  tree.TreeRoot(),
						LeafCount: 6,
						Depth:     confidentialDepth,
					}},
				},
			},
			Courts: map[notary.Address]notary.Court{notarytest.CourtIdentity.Address: courtFor()},
		},
		ledger: ledger,
		claim: notary.Claim{
			Ledger: notarytest.LedgerIdentity,
			Scheme: 2,
			Epoch:  9,
			Entry:  entries[2],
			Proof:  proof,
		},
	}
}

func TestASchemeTwoProofIsIncludedWhenTheTreeAndTheLedgerAgree(t *testing.T) {
	c := confidentialCaseOf(t, true)

	got, err := directoryVerifier(c.directory, c.readers).Verify(context.Background(), c.claim)
	require.NoError(t, err)
	require.Equal(t, notary.VerdictIncluded, got.Verdict)
	require.Equal(t, uint64(0), got.Epoch.Total)
	require.Equal(t, []string{"identity", "epoch", "confirm"}, c.ledger.Asked)
}

func TestASchemeTwoProofForAnotherAmountIsRejectedByBoth(t *testing.T) {
	c := confidentialCaseOf(t, false)
	c.claim.Entry.Amount++

	got, err := directoryVerifier(c.directory, c.readers).Verify(context.Background(), c.claim)
	require.NoError(t, err)
	require.Equal(t, notary.VerdictRejected, got.Verdict)
}

func TestTheVerifierRefusesAVerdictWhenTheLedgerAndPoseidonDisagree(t *testing.T) {
	c := confidentialCaseOf(t, false)

	_, err := directoryVerifier(c.directory, c.readers).Verify(context.Background(), c.claim)
	require.Equal(t, notary.KindDisagreement, notary.KindOf(err))
}

func TestAProofOfOneSchemeForAnEpochOfTheOtherIsTheWrongLedger(t *testing.T) {
	c := confidentialCaseOf(t, true)
	c.claim.Scheme = 1

	_, err := directoryVerifier(c.directory, c.readers).Verify(context.Background(), c.claim)
	require.Equal(t, notary.KindWrongLedger, notary.KindOf(err))
	require.Empty(t, c.ledger.Asked)
}

func TestASchemeOneVerifierStillVerifiesASchemeOneGeneration(t *testing.T) {
	epoch, claim := notarytest.AnchoredClaim(t, 4, 6, 3)
	directory := &notarytest.Directory{
		IdentityValue: notarytest.DirectoryIdentity,
		Generations: []notary.Generation{
			generationOf(notarytest.LedgerIdentity, notarytest.CourtIdentity.Address, 0),
		},
	}
	readers := &notarytest.Readers{
		Ledgers: map[notary.Address]notary.Ledger{notarytest.LedgerIdentity.Address: ledgerFor(epoch, true)},
		Courts:  map[notary.Address]notary.Court{notarytest.CourtIdentity.Address: courtFor()},
	}

	got, err := directoryVerifier(directory, readers).Verify(context.Background(), claim)
	require.NoError(t, err)
	require.Equal(t, notary.VerdictIncluded, got.Verdict)
	require.Equal(t, keccak.NewScheme().Version(), claim.Scheme)
}
