package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/storegen"
)

var _ notary.StatementStore = (*Store)(nil)

func (s *Store) StatementProgress(ctx context.Context) (notary.StatementProgress, bool, error) {
	const op = "postgres.Store.StatementProgress"

	row, err := s.queries.StatementProgress(ctx)
	if errors.Is(err, pgx.ErrNoRows) {
		return notary.StatementProgress{}, false, nil
	}

	if err != nil {
		return notary.StatementProgress{}, false, unavailable(op, err)
	}

	epoch, err := unsigned(op, row.LastEpoch)
	if err != nil {
		return notary.StatementProgress{}, false, err
	}

	if !row.LastDay.Valid {
		return notary.StatementProgress{}, false, notary.Errf(
			notary.KindStoreUnavailable,
			op,
			nil,
			"the statement progress names no day",
		)
	}

	return notary.StatementProgress{Epoch: epoch, Day: row.LastDay.Time}, true, nil
}

func (s *Store) FinalizedAfter(ctx context.Context, after uint64, limit int32) ([]notary.EpochRecord, error) {
	const op = "postgres.Store.FinalizedAfter"

	number, err := signed(op, after)
	if err != nil {
		return nil, err
	}

	rows, err := s.queries.FinalizedEpochsAfter(
		ctx,
		storegen.FinalizedEpochsAfterParams{After: number, RowLimit: limit},
	)
	if err != nil {
		return nil, unavailable(op, err)
	}

	records := make([]notary.EpochRecord, 0, len(rows))
	for _, row := range rows {
		record, err := recordOf(op, row)
		if err != nil {
			return nil, err
		}

		records = append(records, record)
	}

	return records, nil
}

func (s *Store) StatementFingerprints(
	ctx context.Context,
	holders []notary.Holder,
) (map[notary.Holder]notary.Digest, error) {
	const op = "postgres.Store.StatementFingerprints"

	rows, err := s.queries.StatementFingerprints(ctx, holderNames(holders))
	if err != nil {
		return nil, unavailable(op, err)
	}

	fingerprints := make(map[notary.Holder]notary.Digest, len(rows))
	for _, row := range rows {
		holder, err := notary.NewHolder(row.Holder)
		if err != nil {
			return nil, notary.Errf(notary.KindStoreUnavailable, op, err, "a statement names an unusable holder")
		}

		fingerprint, err := digestOf(op, row.Fingerprint)
		if err != nil {
			return nil, err
		}

		fingerprints[holder] = fingerprint
	}

	return fingerprints, nil
}

func (s *Store) AdvanceStatements(
	ctx context.Context,
	progress notary.StatementProgress,
	planned []notary.PlannedStatement,
	at time.Time,
) (bool, error) {
	const op = "postgres.Store.AdvanceStatements"

	epoch, err := signed(op, progress.Epoch)
	if err != nil {
		return false, err
	}

	day := pgtype.Date{Time: progress.Day, Valid: true}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return false, unavailable(op, err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	q := s.queries.WithTx(tx)

	if _, err := q.AdvanceStatementProgress(
		ctx,
		storegen.AdvanceStatementProgressParams{Epoch: epoch, Day: day},
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}

		return false, unavailable(op, err)
	}

	if len(planned) > 0 {
		holders := make([]notary.Holder, len(planned))
		fingerprints := make([][]byte, len(planned))

		for i, statement := range planned {
			holders[i] = statement.Holder
			fingerprints[i] = statement.Fingerprint[:]
		}

		names := holderNames(holders)

		if err := q.SupersedeStatements(
			ctx,
			storegen.SupersedeStatementsParams{Holders: names, Day: day},
		); err != nil {
			return false, unavailable(op, err)
		}

		if err := q.InsertStatements(ctx, storegen.InsertStatementsParams{
			Holders:      names,
			Epoch:        epoch,
			Day:          day,
			Fingerprints: fingerprints,
			At:           pgtype.Timestamptz{Time: at, Valid: true},
		}); err != nil {
			return false, unavailable(op, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return false, unavailable(op, err)
	}

	return true, nil
}

func (s *Store) DueStatements(ctx context.Context, now time.Time, limit int32) ([]notary.StatementDelivery, error) {
	const op = "postgres.Store.DueStatements"

	rows, err := s.queries.DueStatements(ctx, storegen.DueStatementsParams{
		Now:      pgtype.Timestamptz{Time: now, Valid: true},
		RowLimit: limit,
	})
	if err != nil {
		return nil, unavailable(op, err)
	}

	deliveries := make([]notary.StatementDelivery, 0, len(rows))
	for _, row := range rows {
		holder, err := notary.NewHolder(row.Holder)
		if err != nil {
			return nil, notary.Errf(notary.KindStoreUnavailable, op, err, "a statement names an unusable holder")
		}

		epoch, err := unsigned(op, row.EpochNumber)
		if err != nil {
			return nil, err
		}

		if row.Attempts < 0 {
			return nil, notary.Errf(
				notary.KindStoreUnavailable,
				op,
				nil,
				"a statement counts %d attempts",
				row.Attempts,
			)
		}

		deliveries = append(deliveries, notary.StatementDelivery{
			Holder:   holder,
			Epoch:    epoch,
			Attempts: uint32(row.Attempts),
		})
	}

	return deliveries, nil
}

func (s *Store) StatementSent(ctx context.Context, delivery notary.StatementDelivery, notification string) error {
	const op = "postgres.Store.StatementSent"

	epoch, err := signed(op, delivery.Epoch)
	if err != nil {
		return err
	}

	if err := s.queries.MarkStatementSent(ctx, storegen.MarkStatementSentParams{
		NotificationID: notification,
		Holder:         delivery.Holder.String(),
		Epoch:          epoch,
	}); err != nil {
		return unavailable(op, err)
	}

	return nil
}

func (s *Store) StatementRefused(ctx context.Context, delivery notary.StatementDelivery, reason string) error {
	const op = "postgres.Store.StatementRefused"

	epoch, err := signed(op, delivery.Epoch)
	if err != nil {
		return err
	}

	if err := s.queries.MarkStatementRefused(ctx, storegen.MarkStatementRefusedParams{
		Reason: reason,
		Holder: delivery.Holder.String(),
		Epoch:  epoch,
	}); err != nil {
		return unavailable(op, err)
	}

	return nil
}

func (s *Store) StatementPostponed(ctx context.Context, delivery notary.StatementDelivery, next time.Time) error {
	const op = "postgres.Store.StatementPostponed"

	epoch, err := signed(op, delivery.Epoch)
	if err != nil {
		return err
	}

	if err := s.queries.PostponeStatement(ctx, storegen.PostponeStatementParams{
		NextAttemptAt: pgtype.Timestamptz{Time: next, Valid: true},
		Holder:        delivery.Holder.String(),
		Epoch:         epoch,
	}); err != nil {
		return unavailable(op, err)
	}

	return nil
}

func (s *Store) StatementBacklog(ctx context.Context, plannedBefore time.Time) (notary.StatementBacklog, error) {
	const op = "postgres.Store.StatementBacklog"

	row, err := s.queries.StatementBacklog(ctx, pgtype.Timestamptz{Time: plannedBefore, Valid: true})
	if err != nil {
		return notary.StatementBacklog{}, unavailable(op, err)
	}

	return notary.StatementBacklog{Pending: row.Pending, Overdue: row.Overdue}, nil
}

func holderNames(holders []notary.Holder) []string {
	names := make([]string, len(holders))
	for i, holder := range holders {
		names[i] = holder.String()
	}

	return names
}
