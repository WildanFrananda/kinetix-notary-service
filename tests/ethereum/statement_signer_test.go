package ethereum_test

import (
	"crypto/ecdsa"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/ethereum"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

func TestTheSignerSignsOnlyAsTheKeysItHolds(t *testing.T) {
	key, address := notarytest.StatementKey(t, "held")
	_, stranger := notarytest.StatementKey(t, "not held")
	signer := ethereum.NewStatementSigner([]*ecdsa.PrivateKey{key})
	digest := notary.Digest{0x42}

	require.True(t, signer.Holds(address))
	require.False(t, signer.Holds(stranger))

	signature, err := signer.Sign(address, digest)
	require.NoError(t, err)
	require.Contains(t, []byte{27, 28}, signature[64])

	public, err := crypto.SigToPub(digest[:], append(signature[:64:64], signature[64]-27))
	require.NoError(t, err)
	require.Equal(t, address, notary.Address(crypto.PubkeyToAddress(*public)))

	_, err = signer.Sign(stranger, digest)
	require.Equal(t, notary.KindSignerUnavailable, notary.KindOf(err))
}

func TestStatementKeysAreAListAndNothingElse(t *testing.T) {
	first, _ := notarytest.StatementKey(t, "first")
	second, _ := notarytest.StatementKey(t, "second")
	text := hexOf(first) + ", " + hexOf(second)

	keys, err := ethereum.ParseStatementKeys(text)

	require.NoError(t, err)
	require.Len(t, keys, 2)

	for _, broken := range []string{"", "0x12", hexOf(first) + ",", "not-hex" + strings.Repeat("0", 57)} {
		_, err := ethereum.ParseStatementKeys(broken)
		require.Equal(t, notary.KindInvalidArgument, notary.KindOf(err), "%q", broken)
		require.NotContains(t, err.Error(), hexOf(first)[2:])
	}
}

func hexOf(key *ecdsa.PrivateKey) string {
	return "0x" + hex.EncodeToString(crypto.FromECDSA(key))
}
