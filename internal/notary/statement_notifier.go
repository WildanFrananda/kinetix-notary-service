package notary

import "context"

type StatementNotifier interface {
	NotifyStatement(ctx context.Context, holder Holder, epoch uint64) (string, error)
}
