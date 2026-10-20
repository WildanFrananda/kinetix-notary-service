package notary

import "context"

type ConfidentialLedger interface {
	ConfidentialEpoch(ctx context.Context, number uint64) (ConfidentialEpoch, error)
}
