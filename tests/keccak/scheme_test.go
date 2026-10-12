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

func TestAnAccountCommitmentMatchesCast(t *testing.T) {
	scheme := keccak.NewScheme()

	principal := scheme.Commitment(
		notarytest.Holder(t, "9f1d4a3e-1c62-4d0a-9a7b-2f5c8e0b41d7"),
		notary.AccountCustomerBalance,
		notary.Salt(repeated(0x77)),
	)
	suspense := scheme.Commitment(
		notarytest.Holder(t, "UNASSIGNED_DRIVER_SHIPPING_FEE"),
		notary.AccountSuspensePendingEscrow,
		notary.Salt(repeated(0xaa)),
	)

	require.Equal(t, digestOf(t, "0xc290cdcd232a670b35896804c181134c3ccee6c0b1d55784b74e3d24a64e2506"), principal)
	require.Equal(t, digestOf(t, "0xde59662c9c9c4d50cd96b5e5e717f46d6f2a9cf68d05d3149d0601d05b332e46"), suspense)
}

func TestEveryPartOfACommitmentChangesIt(t *testing.T) {
	scheme := keccak.NewScheme()
	holder := notarytest.Holder(t, "p-1")
	base := scheme.Commitment(holder, notary.AccountMerchantAvailable, notary.Salt(repeated(1)))

	require.NotEqual(
		t,
		base,
		scheme.Commitment(notarytest.Holder(t, "p-2"), notary.AccountMerchantAvailable, notary.Salt(repeated(1))),
	)
	require.NotEqual(t, base, scheme.Commitment(holder, notary.AccountMerchantPendingEscrow, notary.Salt(repeated(1))))
	require.NotEqual(t, base, scheme.Commitment(holder, notary.AccountMerchantAvailable, notary.Salt(repeated(2))))
}
