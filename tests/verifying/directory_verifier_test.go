package verifying_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/keccak"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/poseidon"
	"github.com/WildanFrananda/kinetix-notary-service/internal/verifying"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

func generationOf(ledger notary.LedgerIdentity, court notary.Address, genesis uint64) notary.Generation {
	return notary.Generation{
		ChainID:      ledger.ChainID,
		Ledger:       ledger.Address,
		Court:        court,
		GenesisEpoch: genesis,
		Scheme:       1,
	}
}

func directoryVerifier(directory notary.Directory, readers notary.GenerationReaders) *verifying.DirectoryVerifier {
	return verifying.NewDirectoryVerifier(
		directory,
		readers,
		keccak.NewScheme(),
		poseidon.NewScheme(),
		keccak.NewScheme(),
	)
}

func TestTheDirectoryPicksTheLedgerThatHoldsTheEpoch(t *testing.T) {
	epoch, claim := notarytest.AnchoredClaim(t, 4, 6, 3)
	predecessor := notary.LedgerIdentity{ChainID: 84532, Address: notary.Address{0x01}}
	directory := &notarytest.Directory{
		IdentityValue: notarytest.DirectoryIdentity,
		Generations: []notary.Generation{
			generationOf(predecessor, notary.Address{0x02}, 0),
			generationOf(notarytest.LedgerIdentity, notarytest.CourtIdentity.Address, epoch.Number),
		},
	}
	held := ledgerFor(epoch, true)
	sealed := &notarytest.Ledger{IdentityValue: predecessor}
	readers := &notarytest.Readers{
		Ledgers: map[notary.Address]notary.Ledger{
			predecessor.Address:               sealed,
			notarytest.LedgerIdentity.Address: held,
		},
		Courts: map[notary.Address]notary.Court{notarytest.CourtIdentity.Address: courtFor()},
	}

	got, err := directoryVerifier(directory, readers).Verify(context.Background(), claim)

	require.NoError(t, err)
	require.Equal(t, notary.VerdictIncluded, got.Verdict)
	require.Equal(t, notarytest.LedgerIdentity, got.Ledger)
	require.Equal(t, notarytest.CourtIdentity, got.Court)
	require.Empty(t, sealed.Asked)
	require.Equal(t, []string{"identity", "epoch", "confirm"}, held.Asked)
}

func TestAProofIssuedAgainstAnotherLedgerThanTheOneHoldingTheEpochIsTheWrongLedger(t *testing.T) {
	_, claim := notarytest.AnchoredClaim(t, 4, 6, 3)
	holder := notary.LedgerIdentity{ChainID: 84532, Address: notary.Address{0x0e}}
	directory := &notarytest.Directory{
		IdentityValue: notarytest.DirectoryIdentity,
		Generations:   []notary.Generation{generationOf(holder, notarytest.CourtIdentity.Address, 0)},
	}
	readers := &notarytest.Readers{
		Ledgers: map[notary.Address]notary.Ledger{holder.Address: &notarytest.Ledger{IdentityValue: holder}},
		Courts:  map[notary.Address]notary.Court{notarytest.CourtIdentity.Address: courtFor()},
	}

	_, err := directoryVerifier(directory, readers).Verify(context.Background(), claim)

	require.Equal(t, notary.KindWrongLedger, notary.KindOf(err))
}

func TestAnEpochHeldOnAnotherChainIsMovedNeverNotFound(t *testing.T) {
	_, claim := notarytest.AnchoredClaim(t, 4, 6, 3)
	elsewhere := notary.LedgerIdentity{ChainID: 11155111, Address: notary.Address{0x0f}}
	directory := &notarytest.Directory{
		IdentityValue: notarytest.DirectoryIdentity,
		Generations: []notary.Generation{
			generationOf(notarytest.LedgerIdentity, notarytest.CourtIdentity.Address, 0),
			generationOf(elsewhere, notary.Address{0x10}, 0),
		},
	}

	_, err := directoryVerifier(directory, &notarytest.Readers{}).Verify(context.Background(), claim)

	require.Equal(t, notary.KindMoved, notary.KindOf(err))
	require.ErrorContains(t, err, "chain 11155111")
}

func TestADirectoryThatCannotBeTrustedAnswersNothingAboutTheEpoch(t *testing.T) {
	_, claim := notarytest.AnchoredClaim(t, 4, 6, 3)
	refusal := notary.Errf(notary.KindWrongDirectory, "test", nil, "not a directory")
	outage := notary.Errf(notary.KindLedgerUnavailable, "test", nil, "no answer")

	for kind, directory := range map[notary.Kind]*notarytest.Directory{
		notary.KindWrongDirectory:    {IdentityErr: refusal},
		notary.KindLedgerUnavailable: {IdentityValue: notarytest.DirectoryIdentity, Err: outage},
	} {
		_, err := directoryVerifier(directory, &notarytest.Readers{}).Verify(context.Background(), claim)

		require.Equal(t, kind, notary.KindOf(err))
	}
}
