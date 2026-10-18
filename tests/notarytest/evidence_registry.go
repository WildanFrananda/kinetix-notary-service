package notarytest

import (
	"context"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type EvidenceRegistry struct {
	Chain    *EvidenceChain
	Answers  notary.RegistryIdentity
	Err      error
	Disagree bool
}

var _ notary.EvidenceRegistry = (*EvidenceRegistry)(nil)

func (r *EvidenceRegistry) Identity(context.Context) (notary.RegistryIdentity, error) {
	return r.Answers, r.Err
}

func (r *EvidenceRegistry) Batch(
	ctx context.Context,
	stream notary.EvidenceStream,
	number uint64,
) (notary.AnchoredBatch, error) {
	const op = "notarytest.EvidenceRegistry.Batch"

	final, found, _ := r.Chain.BatchAtFinalized(ctx, stream, number)
	if found {
		return final, nil
	}

	if _, found, _ := r.Chain.BatchAtLatest(ctx, stream, number); found {
		return notary.AnchoredBatch{}, notary.Errf(
			notary.KindEpochNotFinal,
			op,
			nil,
			"%s batch %d is not final",
			stream,
			number,
		)
	}

	return notary.AnchoredBatch{}, notary.Errf(notary.KindEpochNotAnchored, op, nil, "no %s batch %d", stream, number)
}

func (r *EvidenceRegistry) ConfirmsEvidence(
	_ context.Context,
	batch notary.AnchoredBatch,
	salt notary.Salt,
	record notary.Digest,
	proof notary.EvidenceProof,
) (bool, error) {
	confirmed := notary.VerifyEvidence(r.Chain.Scheme, batch, salt, record, proof) == nil
	if r.Disagree {
		return !confirmed, nil
	}

	return confirmed, nil
}
