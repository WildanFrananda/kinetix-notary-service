package notarytest

import (
	"context"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type ConfidentialLedger struct {
	Statements map[uint64]notary.ConfidentialEpoch
	Err        error
}

var _ notary.ConfidentialLedger = (*ConfidentialLedger)(nil)

func (l *ConfidentialLedger) ConfidentialEpoch(_ context.Context, number uint64) (notary.ConfidentialEpoch, error) {
	if l.Err != nil {
		return notary.ConfidentialEpoch{}, l.Err
	}

	statement, ok := l.Statements[number]
	if !ok {
		return notary.ConfidentialEpoch{}, notary.Errf(
			notary.KindEpochNotAnchored,
			"notarytest.ConfidentialLedger",
			nil,
			"no epoch %d",
			number,
		)
	}

	return statement, nil
}
