package notary

import "context"

type LedgerHeads interface {
	Head(ctx context.Context) (LedgerHead, error)
}
