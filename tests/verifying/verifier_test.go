package verifying_test

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/ethereum"
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

func courtFor() *notarytest.Court {
	return &notarytest.Court{IdentityValue: notarytest.CourtIdentity, Statements: keccak.NewScheme()}
}

func verifierOn(ledger notary.Ledger, scheme notary.Scheme, court notary.Court) *verifying.Verifier {
	return verifying.NewVerifier(ledger, scheme, keccak.NewScheme(), court)
}

func signed(t *testing.T, claim notary.Claim, label string) (notary.Claim, notary.Address) {
	t.Helper()

	key, address := notarytest.StatementKey(t, label)
	scheme := keccak.NewScheme()
	digest := scheme.StatementDigest(notarytest.CourtIdentity, notary.StatementOf(scheme, claim))

	signature, err := ethereum.NewStatementSigner([]*ecdsa.PrivateKey{key}).Sign(address, digest)
	require.NoError(t, err)

	claim.Statement = &notary.StatementSignature{Court: notarytest.CourtIdentity.Address, Signature: signature}

	return claim, address
}

func TestAClaimBothSidesConfirmIsIncluded(t *testing.T) {
	epoch, claim := notarytest.AnchoredClaim(t, 4, 6, 3)
	ledger := ledgerFor(epoch, true)

	got, err := verifierOn(ledger, keccak.NewScheme(), courtFor()).Verify(context.Background(), claim)

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

	got, err := verifierOn(ledgerFor(epoch, false), keccak.NewScheme(), courtFor()).Verify(context.Background(), claim)

	require.NoError(t, err)
	require.Equal(t, notary.VerdictRejected, got.Verdict)
	require.Contains(t, got.Reason, "the proof sums to")
}

func TestTheLedgerSayingYesWhereTheSchemeSaysNoIsADisagreement(t *testing.T) {
	epoch, claim := notarytest.AnchoredClaim(t, 4, 6, 3)
	claim.Entry.Amount++

	_, err := verifierOn(ledgerFor(epoch, true), keccak.NewScheme(), courtFor()).Verify(context.Background(), claim)

	require.Equal(t, notary.KindDisagreement, notary.KindOf(err))
}

func TestTheLedgerSayingNoWhereTheSchemeSaysYesIsADisagreement(t *testing.T) {
	epoch, claim := notarytest.AnchoredClaim(t, 4, 6, 3)

	_, err := verifierOn(ledgerFor(epoch, false), keccak.NewScheme(), courtFor()).Verify(context.Background(), claim)

	require.Equal(t, notary.KindDisagreement, notary.KindOf(err))
}

func TestAProofIssuedAgainstAnotherLedgerIsRefusedBeforeAnyEpochIsRead(t *testing.T) {
	epoch, claim := notarytest.AnchoredClaim(t, 4, 6, 3)
	claim.Ledger.Address = notary.Address{0x01}
	ledger := ledgerFor(epoch, true)

	_, err := verifierOn(ledger, keccak.NewScheme(), courtFor()).Verify(context.Background(), claim)

	require.Equal(t, notary.KindWrongLedger, notary.KindOf(err))
	require.Equal(t, []string{"identity"}, ledger.Asked)
}

func TestAProofIssuedOnAnotherChainIsRefused(t *testing.T) {
	epoch, claim := notarytest.AnchoredClaim(t, 4, 6, 3)
	claim.Ledger.ChainID = 1

	_, err := verifierOn(ledgerFor(epoch, true), keccak.NewScheme(), courtFor()).Verify(context.Background(), claim)

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

			got, err := verifierOn(ledger, keccak.NewScheme(), courtFor()).Verify(context.Background(), claim)

			require.Equal(t, notary.KindLedgerUnavailable, notary.KindOf(err))
			require.Equal(t, notary.Verification{}, got)
		})
	}
}

func TestAnEpochTheLedgerNeverAnchoredIsReportedAsSuch(t *testing.T) {
	epoch, claim := notarytest.AnchoredClaim(t, 4, 6, 3)
	claim.Epoch = 9

	_, err := verifierOn(ledgerFor(epoch, true), keccak.NewScheme(), courtFor()).Verify(context.Background(), claim)

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

	verification, err := verifierOn(ledger, keccak.NewScheme(), courtFor()).Verify(context.Background(), claim)

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

	verification, err := verifierOn(ledger, scheme, courtFor()).Verify(context.Background(), claim)

	require.NoError(t, err)
	require.Equal(t, notary.VerdictIncluded, verification.Verdict)
	require.Equal(t, opening, verification.Opening)
}

func TestAProofWithoutAStatementSaysSoAndStillCountsItsEpochsFaults(t *testing.T) {
	epoch, claim := notarytest.AnchoredClaim(t, 4, 6, 3)
	court := courtFor()
	court.Faults = map[uint64]uint64{4: 2}

	got, err := verifierOn(ledgerFor(epoch, true), keccak.NewScheme(), court).Verify(context.Background(), claim)

	require.NoError(t, err)
	require.Equal(t, notary.VerdictIncluded, got.Verdict)
	require.Equal(t, notary.StatementAbsent, got.Check)
	require.Equal(t, uint64(2), got.Faults)
	require.Equal(t, notarytest.CourtIdentity, got.Court)
}

func TestAStatementSignedByTheEpochsSignerIsAttributable(t *testing.T) {
	epoch, claim := notarytest.AnchoredClaim(t, 4, 6, 3)
	claim, signer := signed(t, claim, "signer")
	court := courtFor()
	court.Signer = signer

	got, err := verifierOn(ledgerFor(epoch, true), keccak.NewScheme(), court).Verify(context.Background(), claim)

	require.NoError(t, err)
	require.Equal(t, notary.StatementAttributable, got.Check)
	require.Equal(t, signer, got.Signer)
	require.False(t, got.Provable)
}

func TestAStatementSignedByAnyoneElseIsNotAttributable(t *testing.T) {
	epoch, claim := notarytest.AnchoredClaim(t, 4, 6, 3)
	claim, _ = signed(t, claim, "impostor")
	_, signer := notarytest.StatementKey(t, "signer")
	court := courtFor()
	court.Signer = signer

	got, err := verifierOn(ledgerFor(epoch, true), keccak.NewScheme(), court).Verify(context.Background(), claim)

	require.NoError(t, err)
	require.Equal(t, notary.VerdictIncluded, got.Verdict)
	require.Equal(t, notary.StatementUnattributable, got.Check)
}

func TestASignedProofTheLedgerRefusesCanBeTakenToTheCourt(t *testing.T) {
	epoch, claim := notarytest.AnchoredClaim(t, 4, 6, 3)
	claim.Entry.Amount++
	claim, signer := signed(t, claim, "signer")
	court := courtFor()
	court.Signer = signer

	got, err := verifierOn(ledgerFor(epoch, false), keccak.NewScheme(), court).Verify(context.Background(), claim)

	require.NoError(t, err)
	require.Equal(t, notary.VerdictRejected, got.Verdict)
	require.Equal(t, notary.StatementAttributable, got.Check)
	require.True(t, got.Provable)
}

func TestAnAccountThatDoesNotOpenItsCommitmentIsNoCaseForTheCourt(t *testing.T) {
	epoch, claim := notarytest.AnchoredClaim(t, 4, 6, 3)
	claim.Opening = &notary.Opening{
		Holder: notarytest.Holder(t, "someone-else"),
		Kind:   notary.AccountCustomerBalance,
		Salt:   notary.Salt{0x01},
	}
	claim, signer := signed(t, claim, "signer")
	court := courtFor()
	court.Signer = signer

	got, err := verifierOn(ledgerFor(epoch, true), keccak.NewScheme(), court).Verify(context.Background(), claim)

	require.NoError(t, err)
	require.Equal(t, notary.VerdictRejected, got.Verdict)
	require.Equal(t, notary.StatementAttributable, got.Check)
	require.False(t, got.Provable)
}

func TestAStatementSignedForAnotherCourtIsTheWrongCourt(t *testing.T) {
	epoch, claim := notarytest.AnchoredClaim(t, 4, 6, 3)
	claim, _ = signed(t, claim, "signer")
	claim.Statement.Court = notary.Address{0xee}

	_, err := verifierOn(ledgerFor(epoch, true), keccak.NewScheme(), courtFor()).Verify(context.Background(), claim)

	require.Equal(t, notary.KindWrongCourt, notary.KindOf(err))
}

func TestACourtThatHashesTheStatementDifferentlyIsADisagreement(t *testing.T) {
	epoch, claim := notarytest.AnchoredClaim(t, 4, 6, 3)
	claim, _ = signed(t, claim, "signer")
	court := courtFor()
	court.Disagree = true

	_, err := verifierOn(ledgerFor(epoch, true), keccak.NewScheme(), court).Verify(context.Background(), claim)

	require.Equal(t, notary.KindDisagreement, notary.KindOf(err))
}

func TestACourtThatCannotBeReadIsNeverAVerdict(t *testing.T) {
	down := notary.Errf(notary.KindLedgerUnavailable, "test", nil, "down")

	for name, spoil := range map[string]func(*notarytest.Court){
		"identity": func(c *notarytest.Court) { c.IdentityErr = down },
		"faults":   func(c *notarytest.Court) { c.FaultErr = down },
		"signer":   func(c *notarytest.Court) { c.SignerErr = down },
	} {
		t.Run(name, func(t *testing.T) {
			epoch, claim := notarytest.AnchoredClaim(t, 4, 6, 3)
			claim, _ = signed(t, claim, "signer")
			court := courtFor()
			spoil(court)

			got, err := verifierOn(
				ledgerFor(epoch, true),
				keccak.NewScheme(),
				court,
			).Verify(context.Background(), claim)

			require.Equal(t, notary.KindLedgerUnavailable, notary.KindOf(err))
			require.Equal(t, notary.Verification{}, got)
		})
	}
}
