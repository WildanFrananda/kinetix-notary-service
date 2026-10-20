package delivering

import (
	"cmp"
	"context"
	"log/slog"
	"slices"
	"time"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type Planner struct {
	store  notary.StatementStore
	clock  func() time.Time
	log    *slog.Logger
	policy Policy
}

func NewPlanner(store notary.StatementStore, clock func() time.Time, log *slog.Logger, policy Policy) *Planner {
	return &Planner{store: store, clock: clock, log: log, policy: policy}
}

func (p *Planner) Plan(ctx context.Context) error {
	progress, started, err := p.store.StatementProgress(ctx)
	if err != nil {
		return err
	}

	epochs, err := p.candidates(ctx, progress, started)
	if err != nil {
		return err
	}

	for _, epoch := range epochs {
		next := notary.StatementProgress{Epoch: epoch.Header.Number, Day: p.dayOf(epoch.Header.AsOf)}

		var planned []notary.PlannedStatement

		if started && !next.Day.After(progress.Day) {
			next.Day = progress.Day
		} else {
			planned, err = p.planned(ctx, epoch.Header.Number)
			if err != nil {
				return err
			}
		}

		advanced, err := p.store.AdvanceStatements(ctx, next, planned, p.clock())
		if err != nil {
			return err
		}

		if !advanced {
			return nil
		}

		if len(planned) > 0 {
			p.log.InfoContext(
				ctx,
				"statements planned",
				"epoch",
				next.Epoch,
				"day",
				next.Day.Format(time.DateOnly),
				"holders",
				len(planned),
			)
		}

		progress, started = next, true
	}

	return nil
}

func (p *Planner) candidates(
	ctx context.Context,
	progress notary.StatementProgress,
	started bool,
) ([]notary.EpochRecord, error) {
	if started {
		return p.store.FinalizedAfter(ctx, progress.Epoch, p.policy.EpochsPerTick)
	}

	newest, found, err := p.store.LastFinalized(ctx)
	if err != nil || !found {
		return nil, err
	}

	return []notary.EpochRecord{newest}, nil
}

func (p *Planner) planned(ctx context.Context, epoch uint64) ([]notary.PlannedStatement, error) {
	accounts, err := p.store.Accounts(ctx, epoch)
	if err != nil {
		return nil, err
	}

	held := make(map[notary.Holder][]notary.Liability)

	for _, account := range accounts {
		if !account.Liability.Kind.HeldByPrincipal() {
			continue
		}

		held[account.Liability.Holder] = append(held[account.Liability.Holder], account.Liability)
	}

	if len(held) == 0 {
		return nil, nil
	}

	holders := make([]notary.Holder, 0, len(held))
	for holder := range held {
		holders = append(holders, holder)
	}

	slices.SortFunc(holders, func(a, b notary.Holder) int {
		return cmp.Compare(a.String(), b.String())
	})

	last, err := p.store.StatementFingerprints(ctx, holders)
	if err != nil {
		return nil, err
	}

	planned := make([]notary.PlannedStatement, 0, len(holders))
	for _, holder := range holders {
		fingerprint := notary.StatementFingerprint(held[holder])
		if previous, ok := last[holder]; ok && previous == fingerprint {
			continue
		}

		planned = append(planned, notary.PlannedStatement{Holder: holder, Fingerprint: fingerprint})
	}

	return planned, nil
}

func (p *Planner) dayOf(at notary.UnixTime) time.Time {
	local := at.Time().Add(p.policy.UTCOffset)

	return time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, time.UTC)
}
