package notarytest

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/keccak"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

var LedgerIdentity = notary.LedgerIdentity{ChainID: 84532, Address: notary.Address{0xab, 0xcd}}

func AnchoredClaim(t *testing.T, epochNumber uint64, count int, index uint64) (notary.AnchoredEpoch, notary.Claim) {
	t.Helper()

	entries := Entries(epochNumber+1000, count)
	tree, err := notary.BuildTree(keccak.NewScheme(), epochNumber, entries)
	require.NoError(t, err)
	proof, err := tree.Prove(index)
	require.NoError(t, err)

	root := tree.Root()
	epoch := notary.AnchoredEpoch{
		Number:     epochNumber,
		Root:       root.Hash,
		Total:      root.Sum,
		LeafCount:  tree.LeafCount(),
		AsOf:       Unix(t, 1760000000),
		AnchoredAt: Unix(t, 1760000060),
	}
	claim := notary.Claim{Ledger: LedgerIdentity, Epoch: epochNumber, Entry: entries[index], Proof: proof}

	return epoch, claim
}
