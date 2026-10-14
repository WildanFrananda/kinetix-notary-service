package ethereum

import (
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type Recoverer struct{}

var _ notary.SignatureRecoverer = Recoverer{}

func NewRecoverer() Recoverer {
	return Recoverer{}
}

func (Recoverer) Recover(digest notary.Digest, signature notary.Signature) (notary.Address, error) {
	const op = "ethereum.Recoverer.Recover"

	v := signature[64]
	if v >= 27 {
		v -= 27
	}

	r := new(big.Int).SetBytes(signature[:32])

	s := new(big.Int).SetBytes(signature[32:64])
	if v > 1 || !crypto.ValidateSignatureValues(v, r, s, true) {
		return notary.Address{}, notary.Errf(
			notary.KindRegistrationRefused,
			op,
			nil,
			"the signature is not a canonical secp256k1 signature",
		)
	}

	raw := make([]byte, 65)
	copy(raw, signature[:64])
	raw[64] = v

	public, err := crypto.SigToPub(digest[:], raw)
	if err != nil {
		return notary.Address{}, notary.Errf(notary.KindRegistrationRefused, op, err, "the signature recovers no key")
	}

	return notary.Address(crypto.PubkeyToAddress(*public)), nil
}
