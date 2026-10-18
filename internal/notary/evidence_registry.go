package notary

import "context"

type EvidenceRegistry interface {
	Identity(ctx context.Context) (RegistryIdentity, error)
	Batch(ctx context.Context, stream EvidenceStream, number uint64) (AnchoredBatch, error)
	ConfirmsEvidence(
		ctx context.Context,
		batch AnchoredBatch,
		salt Salt,
		record Digest,
		proof EvidenceProof,
	) (bool, error)
}
