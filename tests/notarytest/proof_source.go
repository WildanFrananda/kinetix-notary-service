package notarytest

import (
	"context"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type ProofSource struct {
	Records       map[uint64]notary.EpochRecord
	Stored        map[uint64][]notary.Account
	Err           error
	AccountsReads int
}

var _ notary.ProofSource = (*ProofSource)(nil)

func (s *ProofSource) Epoch(_ context.Context, number uint64) (notary.EpochRecord, bool, error) {
	if s.Err != nil {
		return notary.EpochRecord{}, false, s.Err
	}

	record, ok := s.Records[number]

	return record, ok, nil
}

func (s *ProofSource) LastFinalized(context.Context) (notary.EpochRecord, bool, error) {
	if s.Err != nil {
		return notary.EpochRecord{}, false, s.Err
	}

	var best notary.EpochRecord

	found := false

	for _, record := range s.Records {
		if record.State == notary.StateFinalized && (!found || record.Header.Number > best.Header.Number) {
			best, found = record, true
		}
	}

	return best, found, nil
}

func (s *ProofSource) Accounts(_ context.Context, epoch uint64) ([]notary.Account, error) {
	s.AccountsReads++

	if s.Err != nil {
		return nil, s.Err
	}

	return append([]notary.Account(nil), s.Stored[epoch]...), nil
}
