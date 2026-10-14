package notary

import (
	"context"
	"math/big"
)

type AnchorChain interface {
	Anchorer() Address
	Head(ctx context.Context) (LedgerHead, error)
	EpochAtLatest(ctx context.Context, number uint64) (AnchoredEpoch, bool, error)
	EpochAtFinalized(ctx context.Context, number uint64) (AnchoredEpoch, bool, error)
	Nonce(ctx context.Context) (uint64, error)
	Fees(ctx context.Context) (baseFee uint64, tip uint64, err error)
	Prepare(
		ctx context.Context,
		header EpochHeader,
		approvals []Signature,
		nonce uint64,
		fees FeeQuote,
	) (PreparedAnchor, error)
	Resign(previous PreparedAnchor, fees FeeQuote) (PreparedAnchor, error)
	Broadcast(ctx context.Context, anchor PreparedAnchor) error
	Known(ctx context.Context, hash Digest) (bool, error)
	Receipt(ctx context.Context, hash Digest) (AnchorReceipt, bool, error)
	Balance(ctx context.Context) (*big.Int, error)
}
