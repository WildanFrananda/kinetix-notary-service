package notary

import "context"

type EvidenceProofSource interface {
	OrderEvidence(ctx context.Context, orderNumber string) ([]StoredEvidence, error)
	ReviewEvidence(ctx context.Context, reviewID string) ([]StoredEvidence, error)
	EvidenceBatch(ctx context.Context, stream EvidenceStream, number uint64) (EvidenceBatchRecord, bool, error)
	EvidenceBatchLeaves(ctx context.Context, stream EvidenceStream, number uint64) ([]Digest, error)
}
