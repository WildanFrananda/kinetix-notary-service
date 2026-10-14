package notary

import "context"

type EpochWatch interface {
	Approve(ctx context.Context, proposal EpochProposal) (Approval, error)
}
