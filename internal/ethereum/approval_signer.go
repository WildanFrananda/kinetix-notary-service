package ethereum

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type ApprovalSigner struct {
	key    *ecdsa.PrivateKey
	member notary.Address
}

var _ notary.ApprovalSigner = (*ApprovalSigner)(nil)

func NewApprovalSigner(key *ecdsa.PrivateKey) *ApprovalSigner {
	return &ApprovalSigner{key: key, member: notary.Address(crypto.PubkeyToAddress(key.PublicKey))}
}

func ParseWatchKey(text string) (*ecdsa.PrivateKey, error) {
	return parseKey("ethereum.ParseWatchKey", "watch", text)
}

func (s *ApprovalSigner) Member() notary.Address {
	return s.member
}

func (s *ApprovalSigner) Sign(digest notary.Digest) (notary.Signature, error) {
	raw, err := crypto.Sign(digest[:], s.key)
	if err != nil {
		return notary.Signature{}, notary.Errf(
			notary.KindSignerUnavailable,
			"ethereum.ApprovalSigner.Sign",
			err,
			"signing an approval failed",
		)
	}

	var signature notary.Signature
	copy(signature[:], raw)
	signature[64] += 27

	return signature, nil
}
