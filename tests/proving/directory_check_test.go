package proving_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/proving"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

func notaryGeneration() notary.Generation {
	return notary.Generation{
		ChainID:     84532,
		Ledger:      notary.Address{0x1e},
		Court:       notary.Address{0xc0},
		Withholding: notary.Address{0x3d},
	}
}

func TestTheNotaryStartsOnlyWhenTheDirectoryCurrentlyNamesItsLedgerAndCourts(t *testing.T) {
	expected := notaryGeneration()
	current := expected
	current.GenesisEpoch = 7
	directory := &notarytest.Directory{
		IdentityValue: notarytest.DirectoryIdentity,
		Generations:   []notary.Generation{current},
	}

	require.NoError(t, proving.NewDirectoryCheck(directory, expected).Check(context.Background()))

	for name, change := range map[string]func(*notary.Generation){
		"another chain":             func(g *notary.Generation) { g.ChainID = 1 },
		"another ledger":            func(g *notary.Generation) { g.Ledger = notary.Address{0x1f} },
		"another court":             func(g *notary.Generation) { g.Court = notary.Address{0xc1} },
		"another withholding court": func(g *notary.Generation) { g.Withholding = notary.Address{0x3e} },
	} {
		t.Run(name, func(t *testing.T) {
			other := current
			change(&other)
			directory := &notarytest.Directory{
				IdentityValue: notarytest.DirectoryIdentity,
				Generations:   []notary.Generation{other},
			}

			err := proving.NewDirectoryCheck(directory, expected).Check(context.Background())

			require.Equal(t, notary.KindWrongDirectory, notary.KindOf(err))
		})
	}
}

func TestADirectoryThatCannotBeReadStopsTheNotaryWithItsOwnKind(t *testing.T) {
	outage := notary.Errf(notary.KindLedgerUnavailable, "test", nil, "no answer")
	refusal := notary.Errf(notary.KindWrongDirectory, "test", nil, "not a directory")

	for kind, directory := range map[notary.Kind]*notarytest.Directory{
		notary.KindWrongDirectory:    {IdentityErr: refusal},
		notary.KindLedgerUnavailable: {IdentityValue: notarytest.DirectoryIdentity, Err: outage},
	} {
		err := proving.NewDirectoryCheck(directory, notaryGeneration()).Check(context.Background())

		require.Equal(t, kind, notary.KindOf(err))
	}
}
