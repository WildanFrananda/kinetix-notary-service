package verifying_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/keccak"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/verifying"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

func ledgerFor(epoch notary.AnchoredEpoch, confirmed bool) *notarytest.Ledger {
	return &notarytest.Ledger{
		IdentityValue: notarytest.LedgerIdentity,
		Epochs:        map[uint64]notary.AnchoredEpoch{epoch.Number: epoch},
		Confirmed:     confirmed,
	}
}

func TestAClaimBothSidesConfirmIsIncluded(t *testing.T) {
	epoch, claim := notarytest.AnchoredClaim(t, 4, 6, 3)
	ledger := ledgerFor(epoch, true)

	got, err := verifying.NewVerifier(ledger, keccak.NewScheme()).Verify(context.Background(), claim)

	require.NoError(t, err)
	require.Equal(t, notary.VerdictIncluded, got.Verdict)
	require.Equal(t, epoch, got.Epoch)
	require.Equal(t, notarytest.LedgerIdentity, got.Ledger)
	require.Empty(t, got.Reason)
	require.Equal(t, []string{"identity", "epoch", "confirm"}, ledger.Asked)
}

func TestAClaimBothSidesRefuseIsRejectedWithTheReason(t *testing.T) {
	epoch, claim := notarytest.AnchoredClaim(t, 4, 6, 3)
	claim.Entry.Amount++

	got, err := verifying.NewVerifier(ledgerFor(epoch, false), keccak.NewScheme()).Verify(context.Background(), claim)

	require.NoError(t, err)
	require.Equal(t, notary.VerdictRejected, got.Verdict)
	require.Contains(t, got.Reason, "the proof sums to")
}

func TestTheLedgerSayingYesWhereTheSchemeSaysNoIsADisagreement(t *testing.T) {
	epoch, claim := notarytest.AnchoredClaim(t, 4, 6, 3)
	claim.Entry.Amount++

	_, err := verifying.NewVerifier(ledgerFor(epoch, true), keccak.NewScheme()).Verify(context.Background(), claim)

	require.Equal(t, notary.KindDisagreement, notary.KindOf(err))
}

func TestTheLedgerSayingNoWhereTheSchemeSaysYesIsADisagreement(t *testing.T) {
	epoch, claim := notarytest.AnchoredClaim(t, 4, 6, 3)

	_, err := verifying.NewVerifier(ledgerFor(epoch, false), keccak.NewScheme()).Verify(context.Background(), claim)

	require.Equal(t, notary.KindDisagreement, notary.KindOf(err))
}

func TestAProofIssuedAgainstAnotherLedgerIsRefusedBeforeAnyEpochIsRead(t *testing.T) {
	epoch, claim := notarytest.AnchoredClaim(t, 4, 6, 3)
	claim.Ledger.Address = notary.Address{0x01}
	ledger := ledgerFor(epoch, true)

	_, err := verifying.NewVerifier(ledger, keccak.NewScheme()).Verify(context.Background(), claim)

	require.Equal(t, notary.KindWrongLedger, notary.KindOf(err))
	require.Equal(t, []string{"identity"}, ledger.Asked)
}

func TestAProofIssuedOnAnotherChainIsRefused(t *testing.T) {
	epoch, claim := notarytest.AnchoredClaim(t, 4, 6, 3)
	claim.Ledger.ChainID = 1

	_, err := verifying.NewVerifier(ledgerFor(epoch, true), keccak.NewScheme()).Verify(context.Background(), claim)

	require.Equal(t, notary.KindWrongLedger, notary.KindOf(err))
}

func TestAnUnreachableLedgerIsNeverAVerdict(t *testing.T) {
	epoch, claim := notarytest.AnchoredClaim(t, 4, 6, 3)

	cases := map[string]func(*notarytest.Ledger){
		"identity": func(l *notarytest.Ledger) {
			l.IdentityErr = notary.Errf(notary.KindLedgerUnavailable, "fake", errors.New("refused"), "")
		},
		"epoch": func(l *notarytest.Ledger) {
			l.EpochErr = notary.Errf(notary.KindLedgerUnavailable, "fake", errors.New("timeout"), "")
		},
		"confirm": func(l *notarytest.Ledger) {
			l.ConfirmErr = notary.Errf(notary.KindLedgerUnavailable, "fake", errors.New("reset"), "")
		},
	}
	for name, breakIt := range cases {
		t.Run(name, func(t *testing.T) {
			ledger := ledgerFor(epoch, true)
			breakIt(ledger)

			got, err := verifying.NewVerifier(ledger, keccak.NewScheme()).Verify(context.Background(), claim)

			require.Equal(t, notary.KindLedgerUnavailable, notary.KindOf(err))
			require.Equal(t, notary.Verification{}, got)
		})
	}
}

func TestAnEpochTheLedgerNeverAnchoredIsReportedAsSuch(t *testing.T) {
	epoch, claim := notarytest.AnchoredClaim(t, 4, 6, 3)
	claim.Epoch = 9

	_, err := verifying.NewVerifier(ledgerFor(epoch, true), keccak.NewScheme()).Verify(context.Background(), claim)

	require.Equal(t, notary.KindEpochNotAnchored, notary.KindOf(err))
}

func TestAProofWhoseAccountDoesNotOpenItsCommitmentIsRejectedWithoutAskingTheChain(t *testing.T) {
	epoch, claim := notarytest.AnchoredClaim(t, 4, 6, 1)
	ledger := &notarytest.Ledger{
		IdentityValue: notarytest.LedgerIdentity,
		Epochs:        map[uint64]notary.AnchoredEpoch{4: epoch},
		Confirmed:     true,
	}
	claim.Opening = &notary.Opening{
		Holder: notarytest.Holder(t, "someone-else"),
		Kind:   notary.AccountCustomerBalance,
		Salt:   notary.Salt{0x01},
	}

	verification, err := verifying.NewVerifier(ledger, keccak.NewScheme()).Verify(context.Background(), claim)

	require.NoError(t, err)
	require.Equal(t, notary.VerdictRejected, verification.Verdict)
	require.Contains(t, verification.Reason, "does not open its commitment")
	require.Equal(t, claim.Opening, verification.Opening)
	require.NotContains(t, ledger.Asked, "confirm")
}

func TestAProofWhoseAccountOpensItsCommitmentIsIncludedAndNamesTheAccount(t *testing.T) {
	scheme := keccak.NewScheme()
	holder := notarytest.Holder(t, "principal-0042")
	opening := &notary.Opening{Holder: holder, Kind: notary.AccountMerchantAvailable, Salt: notary.Salt{0x09, 0x08}}
	entries := notarytest.Entries(77, 4)
	entries[2].Commitment = scheme.Commitment(opening.Holder, opening.Kind, opening.Salt)

	tree, err := notary.BuildTree(scheme, 4, entries)
	require.NoError(t, err)

	proof, err := tree.Prove(2)
	require.NoError(t, err)

	root := tree.Root()
	ledger := &notarytest.Ledger{
		IdentityValue: notarytest.LedgerIdentity,
		Epochs: map[uint64]notary.AnchoredEpoch{4: {
			Number:     4,
			Root:       root.Hash,
			Total:      root.Sum,
			LeafCount:  4,
			AsOf:       notarytest.Unix(t, 1760000000),
			AnchoredAt: notarytest.Unix(t, 1760000060),
		}},
		Confirmed: true,
	}
	claim := notary.Claim{
		Ledger:  notarytest.LedgerIdentity,
		Epoch:   4,
		Entry:   entries[2],
		Proof:   proof,
		Opening: opening,
	}

	verification, err := verifying.NewVerifier(ledger, scheme).Verify(context.Background(), claim)

	require.NoError(t, err)
	require.Equal(t, notary.VerdictIncluded, verification.Verdict)
	require.Equal(t, opening, verification.Opening)
}
