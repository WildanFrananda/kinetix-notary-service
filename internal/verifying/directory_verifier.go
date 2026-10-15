package verifying

import (
	"context"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type DirectoryVerifier struct {
	directory  notary.Directory
	readers    notary.GenerationReaders
	scheme     notary.Scheme
	statements notary.StatementScheme
}

func NewDirectoryVerifier(
	directory notary.Directory,
	readers notary.GenerationReaders,
	scheme notary.Scheme,
	statements notary.StatementScheme,
) *DirectoryVerifier {
	return &DirectoryVerifier{directory: directory, readers: readers, scheme: scheme, statements: statements}
}

func (v *DirectoryVerifier) Verify(ctx context.Context, claim notary.Claim) (notary.Verification, error) {
	const op = "verifying.DirectoryVerifier.Verify"

	identity, err := v.directory.Identity(ctx)
	if err != nil {
		return notary.Verification{}, err
	}

	generation, err := v.directory.GenerationOf(ctx, claim.Epoch)
	if err != nil {
		return notary.Verification{}, err
	}

	if generation.ChainID != identity.ChainID {
		return notary.Verification{}, notary.Errf(
			notary.KindMoved,
			op,
			nil,
			"epoch %d is held by the ledger %s on chain %d; this verifier reads chain %d",
			claim.Epoch,
			generation.Ledger.Hex(),
			generation.ChainID,
			identity.ChainID,
		)
	}

	verifier := NewVerifier(v.readers.Ledger(generation), v.scheme, v.statements, v.readers.Court(generation))

	return verifier.Verify(ctx, claim)
}
