package verifying

import (
	"context"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type DirectoryVerifier struct {
	directory    notary.Directory
	readers      notary.GenerationReaders
	scheme       notary.Scheme
	confidential notary.ConfidentialScheme
	statements   notary.StatementScheme
}

func NewDirectoryVerifier(
	directory notary.Directory,
	readers notary.GenerationReaders,
	scheme notary.Scheme,
	confidential notary.ConfidentialScheme,
	statements notary.StatementScheme,
) *DirectoryVerifier {
	return &DirectoryVerifier{
		directory:    directory,
		readers:      readers,
		scheme:       scheme,
		confidential: confidential,
		statements:   statements,
	}
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

	scheme := claim.Scheme
	if scheme == 0 {
		scheme = v.scheme.Version()
	}

	if scheme != generation.Scheme {
		return notary.Verification{}, notary.Errf(
			notary.KindWrongLedger,
			op,
			nil,
			"the proof is scheme %d, and epoch %d is held by a scheme-%d generation",
			scheme,
			claim.Epoch,
			generation.Scheme,
		)
	}

	ledger, court := v.readers.Ledger(generation), v.readers.Court(generation)
	if v.confidential != nil && generation.Scheme == v.confidential.Version() {
		verifier := NewConfidentialVerifier(
			ledger,
			v.readers.Confidential(generation),
			v.scheme,
			v.confidential,
			v.statements,
			court,
		)

		return verifier.Verify(ctx, claim)
	}

	return NewVerifier(ledger, v.scheme, v.statements, court).Verify(ctx, claim)
}
