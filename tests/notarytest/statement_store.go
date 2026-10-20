package notarytest

import (
	"context"
	"slices"
	"sync"
	"time"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type StatementStore struct {
	mu        sync.Mutex
	Finalized []notary.EpochRecord
	Holdings  map[uint64][]notary.Account
	Progress  *notary.StatementProgress
	Rows      []StatementRow
	Raced     bool
	Fail      error
}

var _ notary.StatementStore = (*StatementStore)(nil)

func (s *StatementStore) StatementProgress(context.Context) (notary.StatementProgress, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Fail != nil {
		return notary.StatementProgress{}, false, s.Fail
	}

	if s.Progress == nil {
		return notary.StatementProgress{}, false, nil
	}

	return *s.Progress, true, nil
}

func (s *StatementStore) LastFinalized(context.Context) (notary.EpochRecord, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.Finalized) == 0 {
		return notary.EpochRecord{}, false, nil
	}

	return s.Finalized[len(s.Finalized)-1], true, nil
}

func (s *StatementStore) FinalizedAfter(_ context.Context, after uint64, limit int32) ([]notary.EpochRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var out []notary.EpochRecord

	for _, record := range s.Finalized {
		if record.Header.Number > after && int32(len(out)) < limit {
			out = append(out, record)
		}
	}

	return out, nil
}

func (s *StatementStore) Accounts(_ context.Context, epoch uint64) ([]notary.Account, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return slices.Clone(s.Holdings[epoch]), nil
}

func (s *StatementStore) StatementFingerprints(
	_ context.Context,
	holders []notary.Holder,
) (map[notary.Holder]notary.Digest, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	latest := map[notary.Holder]StatementRow{}

	for _, row := range s.Rows {
		if row.State == "superseded" || !slices.Contains(holders, row.Holder) {
			continue
		}

		if seen, ok := latest[row.Holder]; !ok || row.Epoch > seen.Epoch {
			latest[row.Holder] = row
		}
	}

	out := make(map[notary.Holder]notary.Digest, len(latest))
	for holder, row := range latest {
		out[holder] = row.Fingerprint
	}

	return out, nil
}

func (s *StatementStore) AdvanceStatements(
	_ context.Context,
	progress notary.StatementProgress,
	planned []notary.PlannedStatement,
	at time.Time,
) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Raced || (s.Progress != nil && s.Progress.Epoch >= progress.Epoch) {
		return false, nil
	}

	next := progress
	s.Progress = &next

	for _, statement := range planned {
		for i := range s.Rows {
			row := s.Rows[i]
			if row.Holder == statement.Holder && row.State == "pending" && row.Day.Before(progress.Day) {
				s.Rows[i].State = "superseded"
			}
		}

		if s.has(statement.Holder, progress) {
			continue
		}

		s.Rows = append(s.Rows, StatementRow{
			Holder:      statement.Holder,
			Epoch:       progress.Epoch,
			Day:         progress.Day,
			Fingerprint: statement.Fingerprint,
			State:       "pending",
			Next:        at,
			Planned:     at,
		})
	}

	return true, nil
}

func (s *StatementStore) DueStatements(
	_ context.Context,
	now time.Time,
	limit int32,
) ([]notary.StatementDelivery, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var out []notary.StatementDelivery

	for _, row := range s.Rows {
		if row.State == "pending" && !row.Next.After(now) && int32(len(out)) < limit {
			out = append(out, notary.StatementDelivery{Holder: row.Holder, Epoch: row.Epoch, Attempts: row.Attempts})
		}
	}

	return out, nil
}

func (s *StatementStore) StatementSent(
	_ context.Context,
	delivery notary.StatementDelivery,
	notification string,
) error {
	return s.settle(delivery, func(row *StatementRow) {
		row.State, row.Notification = "sent", notification
	})
}

func (s *StatementStore) StatementRefused(_ context.Context, delivery notary.StatementDelivery, reason string) error {
	return s.settle(delivery, func(row *StatementRow) {
		row.State, row.Reason = "refused", reason
	})
}

func (s *StatementStore) StatementPostponed(
	_ context.Context,
	delivery notary.StatementDelivery,
	next time.Time,
) error {
	return s.settle(delivery, func(row *StatementRow) {
		row.Next = next
	})
}

func (s *StatementStore) StatementBacklog(_ context.Context, plannedBefore time.Time) (notary.StatementBacklog, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var backlog notary.StatementBacklog

	for _, row := range s.Rows {
		if row.State != "pending" {
			continue
		}

		backlog.Pending++

		if row.Planned.Before(plannedBefore) {
			backlog.Overdue++
		}
	}

	return backlog, nil
}

func (s *StatementStore) Snapshot() []StatementRow {
	s.mu.Lock()
	defer s.mu.Unlock()

	return slices.Clone(s.Rows)
}

func (s *StatementStore) settle(delivery notary.StatementDelivery, change func(*StatementRow)) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.Rows {
		row := &s.Rows[i]
		if row.Holder == delivery.Holder && row.Epoch == delivery.Epoch && row.State == "pending" {
			row.Attempts++
			change(row)
		}
	}

	return nil
}

func (s *StatementStore) has(holder notary.Holder, progress notary.StatementProgress) bool {
	for _, row := range s.Rows {
		if row.Holder == holder && (row.Epoch == progress.Epoch || row.Day.Equal(progress.Day)) {
			return true
		}
	}

	return false
}
