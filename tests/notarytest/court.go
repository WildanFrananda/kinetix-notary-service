package notarytest

import (
	"context"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type Court struct {
	IdentityValue notary.CourtIdentity
	IdentityErr   error
	Statements    notary.StatementScheme
	Disagree      bool
	Signer        notary.Address
	Signers       map[uint64]notary.Address
	SignerErr     error
	Faults        map[uint64]uint64
	FaultErr      error
	Total         uint64
	TotalErr      error
}

var _ notary.Court = (*Court)(nil)

func (c *Court) Identity(context.Context) (notary.CourtIdentity, error) {
	return c.IdentityValue, c.IdentityErr
}

func (c *Court) SignerOf(_ context.Context, epoch uint64) (notary.Address, error) {
	if c.SignerErr != nil {
		return notary.Address{}, c.SignerErr
	}

	if signer, ok := c.Signers[epoch]; ok {
		return signer, nil
	}

	return c.Signer, nil
}

func (c *Court) StatementDigest(_ context.Context, statement notary.Statement) (notary.Digest, error) {
	digest := c.Statements.StatementDigest(c.IdentityValue, statement)
	if c.Disagree {
		digest[0] ^= 0xff
	}

	return digest, nil
}

func (c *Court) RecoverSigner(
	_ context.Context,
	statement notary.Statement,
	signature notary.Signature,
) (notary.Address, error) {
	digest := c.Statements.StatementDigest(c.IdentityValue, statement)
	raw := make([]byte, len(signature))
	copy(raw, signature[:])
	raw[64] -= 27

	public, err := crypto.SigToPub(digest[:], raw)
	if err == nil {
		return notary.Address(crypto.PubkeyToAddress(*public)), nil
	}

	return notary.Address{}, nil
}

func (c *Court) FaultCount(_ context.Context, epoch uint64) (uint64, error) {
	return c.Faults[epoch], c.FaultErr
}

func (c *Court) TotalFaults(context.Context) (uint64, error) {
	return c.Total, c.TotalErr
}
