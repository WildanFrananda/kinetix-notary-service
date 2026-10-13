package ethereum

import (
	"crypto/ecdsa"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type StatementSigner struct {
	keys map[notary.Address]*ecdsa.PrivateKey
}

var _ notary.StatementSigner = (*StatementSigner)(nil)

func NewStatementSigner(keys []*ecdsa.PrivateKey) *StatementSigner {
	signer := &StatementSigner{keys: make(map[notary.Address]*ecdsa.PrivateKey, len(keys))}
	for _, key := range keys {
		signer.keys[notary.Address(crypto.PubkeyToAddress(key.PublicKey))] = key
	}

	return signer
}

func ParseStatementKeys(text string) ([]*ecdsa.PrivateKey, error) {
	const op = "ethereum.ParseStatementKeys"

	var keys []*ecdsa.PrivateKey

	for i, part := range strings.Split(text, ",") {
		digits, found := strings.CutPrefix(strings.TrimSpace(part), "0x")
		if !found || len(digits) != 64 {
			return nil, notary.Errf(
				notary.KindInvalidArgument,
				op,
				nil,
				"statement key %d is not 0x and 64 hex digits",
				i,
			)
		}

		key, err := crypto.HexToECDSA(digits)
		if err != nil {
			return nil, notary.Errf(
				notary.KindInvalidArgument,
				op,
				nil,
				"statement key %d is not a secp256k1 private key",
				i,
			)
		}

		keys = append(keys, key)
	}

	return keys, nil
}

func (s *StatementSigner) Signers() []notary.Address {
	signers := make([]notary.Address, 0, len(s.keys))
	for address := range s.keys {
		signers = append(signers, address)
	}

	return signers
}

func (s *StatementSigner) Holds(signer notary.Address) bool {
	_, held := s.keys[signer]

	return held
}

func (s *StatementSigner) Sign(signer notary.Address, digest notary.Digest) (notary.Signature, error) {
	const op = "ethereum.StatementSigner.Sign"

	key, held := s.keys[signer]
	if !held {
		return notary.Signature{}, notary.Errf(
			notary.KindSignerUnavailable,
			op,
			nil,
			"the court names %s as the signer and this notary holds no key for it",
			signer.Hex(),
		)
	}

	raw, err := crypto.Sign(digest[:], key)
	if err != nil {
		return notary.Signature{}, notary.Errf(notary.KindSignerUnavailable, op, err, "signing failed")
	}

	var signature notary.Signature
	copy(signature[:], raw)
	signature[64] += 27

	return signature, nil
}
