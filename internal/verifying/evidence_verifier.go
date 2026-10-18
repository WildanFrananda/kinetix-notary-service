package verifying

import (
	"context"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type EvidenceVerifier struct {
	registry notary.EvidenceRegistry
	scheme   notary.EvidenceScheme
}

func NewEvidenceVerifier(registry notary.EvidenceRegistry, scheme notary.EvidenceScheme) *EvidenceVerifier {
	return &EvidenceVerifier{registry: registry, scheme: scheme}
}

func (v *EvidenceVerifier) Verify(
	ctx context.Context,
	claim notary.EvidenceClaim,
) (notary.EvidenceVerification, error) {
	const op = "verifying.EvidenceVerifier.Verify"

	identity, err := v.registry.Identity(ctx)
	if err != nil {
		return notary.EvidenceVerification{}, err
	}

	if identity != claim.Registry {
		return notary.EvidenceVerification{}, notary.Errf(
			notary.KindWrongLedger,
			op,
			nil,
			"the proof was issued against registry %s on chain %d, this verifier reads %s on chain %d",
			claim.Registry.Address.Hex(),
			claim.Registry.ChainID,
			identity.Address.Hex(),
			identity.ChainID,
		)
	}

	if claim.Scheme != v.scheme.EvidenceVersion() {
		return notary.EvidenceVerification{}, notary.Errf(
			notary.KindInvalidArgument,
			op,
			nil,
			"the proof uses evidence scheme %d, this verifier speaks scheme %d",
			claim.Scheme,
			v.scheme.EvidenceVersion(),
		)
	}

	stream, number := claim.Header.Stream, claim.Header.Number

	batch, err := v.registry.Batch(ctx, stream, number)
	if err != nil {
		return notary.EvidenceVerification{}, err
	}

	rejected := func(reason string) notary.EvidenceVerification {
		return notary.EvidenceVerification{
			Verdict:  notary.VerdictRejected,
			Registry: identity,
			Batch:    batch,
			Record:   claim.Record,
			Reason:   reason,
		}
	}

	if claim.Record.Stream() != stream {
		return rejected("the record belongs to another stream than its batch"), nil
	}

	if v.scheme.BatchDigest(claim.Header) != claim.Digest {
		return rejected("the batch header in the proof does not hash to the digest it names"), nil
	}

	if batch.Digest != claim.Digest || batch.Root != claim.Header.Root {
		return rejected(
			"the registry holds " + stream.String() + " batch with digest " + batch.Digest.Hex() +
				", the proof names " + claim.Digest.Hex(),
		), nil
	}

	record, err := v.scheme.RecordDigest(claim.Record)
	if err != nil {
		return rejected(notary.MessageOf(err)), nil
	}

	local := notary.VerifyEvidence(v.scheme, batch, claim.Salt, record, claim.Proof)
	if local != nil && notary.KindOf(local) != notary.KindProofRejected {
		return notary.EvidenceVerification{}, local
	}

	confirmed, err := v.registry.ConfirmsEvidence(ctx, batch, claim.Salt, record, claim.Proof)
	if err != nil {
		return notary.EvidenceVerification{}, err
	}

	if confirmed != (local == nil) {
		return notary.EvidenceVerification{}, notary.Errf(
			notary.KindDisagreement,
			op,
			local,
			"the registry answers %t for %s batch %d, this verifier's scheme answers %t",
			confirmed,
			stream,
			number,
			local == nil,
		)
	}

	if local != nil {
		return rejected(notary.MessageOf(local)), nil
	}

	return notary.EvidenceVerification{
		Verdict:  notary.VerdictIncluded,
		Registry: identity,
		Batch:    batch,
		Record:   claim.Record,
	}, nil
}
