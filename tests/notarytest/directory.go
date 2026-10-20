package notarytest

import (
	"context"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

var DirectoryIdentity = notary.DirectoryIdentity{ChainID: 84532, Address: notary.Address{0xd1, 0x7e}}

type Directory struct {
	IdentityValue notary.DirectoryIdentity
	IdentityErr   error
	Generations   []notary.Generation
	Err           error
}

var _ notary.Directory = (*Directory)(nil)

func (d *Directory) Identity(context.Context) (notary.DirectoryIdentity, error) {
	return d.IdentityValue, d.IdentityErr
}

func (d *Directory) Current(context.Context) (notary.Generation, error) {
	if d.Err != nil {
		return notary.Generation{}, d.Err
	}

	return d.Generations[len(d.Generations)-1], nil
}

func (d *Directory) Lineage(context.Context) ([]notary.Generation, error) {
	if d.Err != nil {
		return nil, d.Err
	}

	return append([]notary.Generation(nil), d.Generations...), nil
}

func (d *Directory) GenerationOf(_ context.Context, epoch uint64) (notary.Generation, error) {
	if d.Err != nil {
		return notary.Generation{}, d.Err
	}

	for i := len(d.Generations) - 1; i >= 0; i-- {
		if d.Generations[i].GenesisEpoch <= epoch {
			return d.Generations[i], nil
		}
	}

	return notary.Generation{}, notary.Errf(
		notary.KindEpochNotAnchored,
		"notarytest.Directory.GenerationOf",
		nil,
		"epoch %d comes before the directory",
		epoch,
	)
}
