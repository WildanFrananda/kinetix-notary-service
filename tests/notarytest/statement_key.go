package notarytest

import (
	"crypto/ecdsa"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

var CourtIdentity = notary.CourtIdentity{ChainID: 84532, Address: notary.Address{0xc0, 0x07}}

func StatementKey(t *testing.T, label string) (*ecdsa.PrivateKey, notary.Address) {
	t.Helper()

	key, err := crypto.ToECDSA(crypto.Keccak256([]byte("kinetix.notary.test." + label)))
	require.NoError(t, err)

	return key, notary.Address(crypto.PubkeyToAddress(key.PublicKey))
}
