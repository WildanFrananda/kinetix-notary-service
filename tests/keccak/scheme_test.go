package keccak_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/keccak"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

func digestOf(t *testing.T, text string) notary.Digest {
	t.Helper()
	d, err := notary.ParseDigest(text)
	require.NoError(t, err)
	return d
}

func repeated(b byte) notary.Digest {
	var d notary.Digest
	copy(d[:], bytes.Repeat([]byte{b}, len(d)))
	return d
}

func TestTheEmptyLeafMatchesCast(t *testing.T) {
	got := keccak.NewScheme().EmptyHash()

	require.Equal(t, digestOf(t, "0xf2ee15ea639b73fa3db9b34a245bdfa015c260c598b211bf05a1ecc4b3e3b4f2"), got)
}

func TestALeafMatchesCast(t *testing.T) {
	got := keccak.NewScheme().LeafHash(7, notary.Entry{Commitment: repeated(0x11), Amount: 1000})

	require.Equal(t, digestOf(t, "0x202b672155173a450d47f4e0e2e16a8a26b9160b1d0eec2bb5f8e3bb36fbb7de"), got)
}

func TestANodeMatchesCast(t *testing.T) {
	got := keccak.NewScheme().NodeHash(
		notary.Node{Hash: repeated(0xaa), Sum: 1},
		notary.Node{Hash: repeated(0xbb), Sum: 2},
	)

	require.Equal(t, digestOf(t, "0xabf64866904f3af2d285e566fb728a984218b3da46ea0bc784beff8a67ce7b74"), got)
}

func TestAnEpochDigestMatchesCast(t *testing.T) {
	got := keccak.NewScheme().EpochDigest(notary.EpochHeader{
		Number:         5,
		PreviousDigest: repeated(0x22),
		Root:           repeated(0x33),
		Total:          123456789,
		LeafCount:      42,
		AsOf:           notarytest.Unix(t, 1760000000),
	})

	require.Equal(t, digestOf(t, "0x9f4e54b0d366f7dfbe07c3641e4137ff738e612d4df1d01fd35dac7622e7915b"), got)
}

func TestTheSchemeIsVersionOne(t *testing.T) {
	require.Equal(t, uint16(1), keccak.NewScheme().Version())
}
