package notarytest

import (
	"context"
	"sort"
	"time"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type Store struct {
	Now      func() time.Time
	Records  map[uint64]notary.EpochRecord
	Accounts map[uint64][]notary.Account
	attempts map[uint64][]notary.Attempt
	holder   string
	token    int64
	expires  time.Time
	Writes   int
}

var _ notary.EpochStore = (*Store)(nil)

func NewStore(now func() time.Time) *Store {
	return &Store{
		Now:      now,
		Records:  map[uint64]notary.EpochRecord{},
		Accounts: map[uint64][]notary.Account{},
		attempts: map[uint64][]notary.Attempt{},
	}
}

func (s *Store) Lease(_ context.Context, instance string, ttl time.Duration) (notary.Lease, bool, error) {
	now := s.Now()
	if s.holder != instance && now.Before(s.expires) {
		return notary.Lease{}, false, nil
	}

	if s.holder != instance || !now.Before(s.expires) {
		s.token++
	}

	s.holder, s.expires = instance, now.Add(ttl)

	return notary.Lease{Instance: instance, Token: s.token, ExpiresAt: s.expires}, true, nil
}

func (s *Store) fence(lease notary.Lease) error {
	if lease.Instance != s.holder || lease.Token != s.token || !s.Now().Before(s.expires) {
		return notary.Errf(notary.KindLeaseLost, "notarytest.Store", nil, "lease lost")
	}

	s.Writes++

	return nil
}

func (s *Store) InFlight(context.Context) (notary.EpochRecord, bool, error) {
	for _, r := range s.Records {
		if r.State.InFlight() {
			return r, true, nil
		}
	}

	return notary.EpochRecord{}, false, nil
}

func (s *Store) LastFinalized(context.Context) (notary.EpochRecord, bool, error) {
	var best notary.EpochRecord

	found := false
	for _, r := range s.Records {
		if r.State == notary.StateFinalized && (!found || r.Header.Number > best.Header.Number) {
			best, found = r, true
		}
	}

	return best, found, nil
}

func (s *Store) SaveBuilt(
	ctx context.Context,
	lease notary.Lease,
	ledger notary.LedgerIdentity,
	plan notary.EpochPlan,
) (notary.EpochRecord, error) {
	if err := s.fence(lease); err != nil {
		return notary.EpochRecord{}, err
	}

	if _, exists := s.Records[plan.Header.Number]; exists {
		return notary.EpochRecord{}, notary.Errf(notary.KindLedgerDiverged, "notarytest.Store", nil, "epoch exists")
	}

	if _, busy, _ := s.InFlight(ctx); busy {
		return notary.EpochRecord{}, notary.Errf(notary.KindLedgerDiverged, "notarytest.Store", nil, "one in flight")
	}

	record := notary.EpochRecord{
		Ledger:   ledger,
		Header:   plan.Header,
		Digest:   plan.Digest,
		Currency: plan.Currency,
		State:    notary.StateBuilt,
	}
	s.Records[plan.Header.Number] = record
	s.Accounts[plan.Header.Number] = plan.Accounts

	return record, nil
}

func (s *Store) RecordAttempt(_ context.Context, lease notary.Lease, attempt notary.Attempt) error {
	if err := s.fence(lease); err != nil {
		return err
	}

	s.attempts[attempt.Epoch] = append(s.attempts[attempt.Epoch], attempt)
	if r := s.Records[attempt.Epoch]; r.State == notary.StateBuilt {
		r.State = notary.StateSubmitted
		s.Records[attempt.Epoch] = r
	}

	return nil
}

func (s *Store) Attempts(_ context.Context, epoch uint64) ([]notary.Attempt, error) {
	out := append([]notary.Attempt(nil), s.attempts[epoch]...)
	sort.SliceStable(out, func(a, b int) bool { return out[a].SentAt.Before(out[b].SentAt) })

	return out, nil
}

func (s *Store) SettleAttempt(
	_ context.Context,
	lease notary.Lease,
	hash notary.Digest,
	outcome notary.AttemptOutcome,
	receipt *notary.AnchorReceipt,
) error {
	if err := s.fence(lease); err != nil {
		return err
	}

	for epoch, list := range s.attempts {
		for i := range list {
			if list[i].Anchor.Hash == hash {
				list[i].Outcome, list[i].Receipt = outcome, receipt
				s.attempts[epoch] = list

				return nil
			}
		}
	}

	return notary.Errf(notary.KindInvalidArgument, "notarytest.Store", nil, "no attempt %s", hash.Hex())
}

func (s *Store) Transition(
	_ context.Context,
	lease notary.Lease,
	epoch uint64,
	from, to notary.EpochState,
	reason string,
) error {
	if !from.CanBecome(to) || (to == notary.StateAbandoned) != (reason != "") {
		return notary.Errf(notary.KindInvalidArgument, "notarytest.Store", nil, "%s -> %s", from, to)
	}

	if err := s.fence(lease); err != nil {
		return err
	}

	r, ok := s.Records[epoch]
	if !ok || r.State != from {
		return notary.Errf(notary.KindLeaseLost, "notarytest.Store", nil, "epoch %d is not %s", epoch, from)
	}

	r.State, r.AbandonedReason = to, reason
	s.Records[epoch] = r

	return nil
}

func (s *Store) AbandonedCount(context.Context) (int64, error) {
	var n int64

	for _, r := range s.Records {
		if r.State == notary.StateAbandoned {
			n++
		}
	}

	return n, nil
}

func (s *Store) RecentCosts(_ context.Context, limit int32) ([]uint64, error) {
	var costs []uint64
	for number := uint64(len(s.Records)); number > 0 && int32(len(costs)) < limit; number-- {
		if s.Records[number-1].State != notary.StateFinalized {
			continue
		}

		for _, a := range s.attempts[number-1] {
			if a.Outcome == notary.AttemptIncluded && a.Receipt != nil {
				cost, _ := a.Receipt.CostWei()
				costs = append(costs, cost)
			}
		}
	}

	return costs, nil
}

func (s *Store) AttemptsOf(epoch uint64) []notary.Attempt {
	return append([]notary.Attempt(nil), s.attempts[epoch]...)
}
