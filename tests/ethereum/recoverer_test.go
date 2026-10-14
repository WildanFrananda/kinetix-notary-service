package ethereum_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/ethereum"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

func signedBy(t *testing.T, label string, digest notary.Digest) (notary.Signature, notary.Address) {
	t.Helper()

	key, address := notarytest.StatementKey(t, label)
	raw, err := crypto.Sign(digest[:], key)
	require.NoError(t, err)

	var signature notary.Signature
	copy(signature[:], raw)

	return signature, address
}

func TestTheRecovererAcceptsBothRecoveryIdConventions(t *testing.T) {
	digest := notary.Digest{0x42}
	signature, address := signedBy(t, "recover", digest)

	recovered, err := ethereum.NewRecoverer().Recover(digest, signature)
	require.NoError(t, err)
	require.Equal(t, address, recovered)

	signature[64] += 27
	recovered, err = ethereum.NewRecoverer().Recover(digest, signature)
	require.NoError(t, err)
	require.Equal(t, address, recovered)
}

func TestTheRecovererRefusesAMalleableOrBrokenSignature(t *testing.T) {
	digest := notary.Digest{0x42}
	signature, _ := signedBy(t, "recover", digest)

	high := signature
	n := crypto.S256().Params().N
	s := new(big.Int).SetBytes(signature[32:64])
	new(big.Int).Sub(n, s).FillBytes(high[32:64])
	high[64] ^= 1

	wrongV := signature
	wrongV[64] = 29

	for name, broken := range map[string]notary.Signature{
		"high s":      high,
		"recovery 29": wrongV,
		"all zero":    {},
	} {
		t.Run(name, func(t *testing.T) {
			_, err := ethereum.NewRecoverer().Recover(digest, broken)

			require.Equal(t, notary.KindRegistrationRefused, notary.KindOf(err))
		})
	}
}
