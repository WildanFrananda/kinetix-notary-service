package notary

import (
	"context"
	"math/big"
)

type EvidenceChain interface {
	Anchorer() Address
	HoldsAnchorerRole(ctx context.Context) (bool, error)
	StreamHead(ctx context.Context, stream EvidenceStream) (EvidenceStreamHead, error)
	BatchAtLatest(ctx context.Context, stream EvidenceStream, number uint64) (AnchoredBatch, bool, error)
	BatchAtFinalized(ctx context.Context, stream EvidenceStream, number uint64) (AnchoredBatch, bool, error)
	Nonce(ctx context.Context) (uint64, error)
	Fees(ctx context.Context) (baseFee uint64, tip uint64, err error)
	PrepareBatch(ctx context.Context, header EvidenceBatchHeader, nonce uint64, fees FeeQuote) (PreparedAnchor, error)
	Resign(previous PreparedAnchor, fees FeeQuote) (PreparedAnchor, error)
	Broadcast(ctx context.Context, anchor PreparedAnchor) error
	Known(ctx context.Context, hash Digest) (bool, error)
	Receipt(ctx context.Context, hash Digest) (AnchorReceipt, bool, error)
	Balance(ctx context.Context) (*big.Int, error)
}
