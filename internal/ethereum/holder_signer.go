package ethereum

import (
	"crypto/ecdsa"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type HolderSigner struct {
	key    *ecdsa.PrivateKey
	holder notary.Address
}

var _ notary.HolderSigner = (*HolderSigner)(nil)

func OpenHolderKeystore(keystorePath, passwordPath string) (*HolderSigner, error) {
	const op = "ethereum.OpenHolderKeystore"

	encrypted, err := os.ReadFile(filepath.Clean(keystorePath))
	if err != nil {
		return nil, notary.Errf(notary.KindInvalidArgument, op, err, "the holder's keystore could not be read")
	}

	password, err := os.ReadFile(filepath.Clean(passwordPath))
	if err != nil {
		return nil, notary.Errf(notary.KindInvalidArgument, op, err, "the holder's password file could not be read")
	}

	key, err := keystore.DecryptKey(encrypted, trimNewline(string(password)))
	if err != nil {
		return nil, notary.Errf(
			notary.KindInvalidArgument,
			op,
			nil,
			"the holder's keystore does not open with that password",
		)
	}

	return NewHolderSigner(key.PrivateKey), nil
}

func NewHolderSigner(key *ecdsa.PrivateKey) *HolderSigner {
	return &HolderSigner{key: key, holder: notary.Address(crypto.PubkeyToAddress(key.PublicKey))}
}

func (s *HolderSigner) Holder() notary.Address {
	return s.holder
}

func (s *HolderSigner) SignOperation(hash notary.Digest) (notary.Signature, error) {
	signature, err := crypto.Sign(hash[:], s.key)
	if err != nil {
		return notary.Signature{}, notary.Errf(
			notary.KindInvalidArgument,
			"ethereum.HolderSigner.SignOperation",
			err,
			"",
		)
	}

	signature[64] += 27

	return notary.Signature(signature), nil
}

func (s *HolderSigner) Authorize(chainID uint64, delegate notary.Address, nonce uint64) (notary.Delegation, error) {
	signed, err := types.SignSetCode(s.key, types.SetCodeAuthorization{
		ChainID: *uint256.NewInt(chainID),
		Address: common.Address(delegate),
		Nonce:   nonce,
	})
	if err != nil {
		return notary.Delegation{}, notary.Errf(notary.KindInvalidArgument, "ethereum.HolderSigner.Authorize", err, "")
	}

	return notary.Delegation{
		ChainID:  chainID,
		Delegate: delegate,
		Nonce:    nonce,
		YParity:  signed.V,
		R:        signed.R.Bytes32(),
		S:        signed.S.Bytes32(),
	}, nil
}

func trimNewline(text string) string {
	for len(text) > 0 && (text[len(text)-1] == '\n' || text[len(text)-1] == '\r') {
		text = text[:len(text)-1]
	}

	return text
}
