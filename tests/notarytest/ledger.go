package notarytest

import (
	"context"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type Ledger struct {
	IdentityValue notary.LedgerIdentity
	IdentityErr   error
	Epochs        map[uint64]notary.AnchoredEpoch
	EpochErr      error
	Confirmed     bool
	ConfirmErr    error
	Asked         []string
}

var _ notary.Ledger = (*Ledger)(nil)

func (l *Ledger) Identity(context.Context) (notary.LedgerIdentity, error) {
	l.Asked = append(l.Asked, "identity")
	return l.IdentityValue, l.IdentityErr
}

func (l *Ledger) Epoch(_ context.Context, number uint64) (notary.AnchoredEpoch, error) {
	l.Asked = append(l.Asked, "epoch")
	if l.EpochErr != nil {
		return notary.AnchoredEpoch{}, l.EpochErr
	}

	epoch, ok := l.Epochs[number]
	if !ok {
		return notary.AnchoredEpoch{}, notary.Errf(
			notary.KindEpochNotAnchored,
			"notarytest.Ledger.Epoch",
			nil,
			"no epoch %d",
			number,
		)
	}

	return epoch, nil
}

func (l *Ledger) ConfirmsInclusion(context.Context, notary.Claim) (bool, error) {
	l.Asked = append(l.Asked, "confirm")
	return l.Confirmed, l.ConfirmErr
}
