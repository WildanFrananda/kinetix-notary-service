package notarytest

import (
	"context"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type Heads struct {
	Value notary.LedgerHead
	Err   error
}

var _ notary.LedgerHeads = (*Heads)(nil)

func (h *Heads) Head(context.Context) (notary.LedgerHead, error) {
	return h.Value, h.Err
}
