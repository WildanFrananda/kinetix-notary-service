package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/storegen"
)

var _ notary.WithholdingScans = (*Store)(nil)

func (s *Store) ScanFrom(ctx context.Context, court notary.CourtIdentity) (uint64, bool, error) {
	const op = "postgres.Store.ScanFrom"

	chainID, err := signed(op, court.ChainID)
	if err != nil {
		return 0, false, err
	}

	next, err := s.queries.WithholdingScanFrom(
		ctx,
		storegen.WithholdingScanFromParams{ChainID: chainID, Court: court.Address[:]},
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, false, nil
	}

	if err != nil {
		return 0, false, unavailable(op, err)
	}

	block, err := unsigned(op, next)
	if err != nil {
		return 0, false, err
	}

	return block, true, nil
}

func (s *Store) OpenRequests(ctx context.Context, court notary.CourtIdentity) ([]notary.WithholdingRequest, error) {
	const op = "postgres.Store.OpenRequests"

	chainID, err := signed(op, court.ChainID)
	if err != nil {
		return nil, err
	}

	rows, err := s.queries.OpenWithholdingRequests(
		ctx,
		storegen.OpenWithholdingRequestsParams{ChainID: chainID, Court: court.Address[:]},
	)
	if err != nil {
		return nil, unavailable(op, err)
	}

	requests := make([]notary.WithholdingRequest, 0, len(rows))

	for _, row := range rows {
		request, err := withholdingRequestOf(op, row)
		if err != nil {
			return nil, err
		}

		requests = append(requests, request)
	}

	return requests, nil
}

func (s *Store) RecordScan(
	ctx context.Context,
	court notary.CourtIdentity,
	next uint64,
	found []notary.WithholdingRequest,
) error {
	const op = "postgres.Store.RecordScan"

	chainID, err := signed(op, court.ChainID)
	if err != nil {
		return err
	}

	block, err := signed(op, next)
	if err != nil {
		return err
	}

	params := storegen.InsertWithholdingRequestsParams{ChainID: chainID, Court: court.Address[:]}

	for _, request := range found {
		epoch, err := signed(op, request.Epoch)
		if err != nil {
			return err
		}

		dueAt, err := signed(op, request.DueAt.Seconds())
		if err != nil {
			return err
		}

		params.Registrations = append(params.Registrations, request.Registration[:])
		params.Epochs = append(params.Epochs, epoch)
		params.Accounts = append(params.Accounts, request.Account[:])
		params.Holders = append(params.Holders, request.Holder[:])
		params.DueAts = append(params.DueAts, dueAt)
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return unavailable(op, err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	q := s.queries.WithTx(tx)

	if err := q.AdvanceWithholdingScan(
		ctx,
		storegen.AdvanceWithholdingScanParams{ChainID: chainID, Court: court.Address[:], NextBlock: block},
	); err != nil {
		return unavailable(op, err)
	}

	if len(found) > 0 {
		if err := q.InsertWithholdingRequests(ctx, params); err != nil {
			return unavailable(op, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return unavailable(op, err)
	}

	return nil
}

func (s *Store) CloseRequest(
	ctx context.Context,
	court notary.CourtIdentity,
	registration notary.Digest,
	epoch uint64,
) error {
	const op = "postgres.Store.CloseRequest"

	chainID, err := signed(op, court.ChainID)
	if err != nil {
		return err
	}

	number, err := signed(op, epoch)
	if err != nil {
		return err
	}

	if err := s.queries.CloseWithholdingRequest(
		ctx,
		storegen.CloseWithholdingRequestParams{
			ChainID:      chainID,
			Court:        court.Address[:],
			Registration: registration[:],
			Epoch:        number,
		},
	); err != nil {
		return unavailable(op, err)
	}

	return nil
}

func withholdingRequestOf(op string, row storegen.OpenWithholdingRequestsRow) (notary.WithholdingRequest, error) {
	registration, err := digestOf(op, row.Registration)
	if err != nil {
		return notary.WithholdingRequest{}, err
	}

	account, err := digestOf(op, row.Account)
	if err != nil {
		return notary.WithholdingRequest{}, err
	}

	holder, err := addressOf(op, row.Holder)
	if err != nil {
		return notary.WithholdingRequest{}, err
	}

	epoch, err := unsigned(op, row.EpochNumber)
	if err != nil {
		return notary.WithholdingRequest{}, err
	}

	seconds, err := unsigned(op, row.DueAt)
	if err != nil {
		return notary.WithholdingRequest{}, err
	}

	dueAt, err := notary.UnixTimeFromSeconds(seconds)
	if err != nil {
		return notary.WithholdingRequest{}, notary.Errf(notary.KindStoreUnavailable, op, err, "")
	}

	return notary.WithholdingRequest{
		Registration: registration,
		Account:      account,
		Epoch:        epoch,
		Holder:       holder,
		DueAt:        dueAt,
	}, nil
}
