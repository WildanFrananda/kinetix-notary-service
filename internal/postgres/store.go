package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/storegen"
)

const uniqueViolation = "23505"

type Store struct {
	pool    *pgxpool.Pool
	queries *storegen.Queries
	scheme  uint16
}

var (
	_ notary.EpochStore  = (*Store)(nil)
	_ notary.ProofSource = (*Store)(nil)
)

func NewStore(pool *pgxpool.Pool, scheme notary.Scheme) *Store {
	return &Store{pool: pool, queries: storegen.New(pool), scheme: scheme.Version()}
}

func (s *Store) Lease(ctx context.Context, instance string, ttl time.Duration) (notary.Lease, bool, error) {
	const op = "postgres.Store.Lease"
	if instance == "" || ttl <= 0 {
		return notary.Lease{}, false, notary.Errf(
			notary.KindInvalidArgument,
			op,
			nil,
			"a lease needs a holder and a positive term",
		)
	}

	row, err := s.queries.AcquireLease(
		ctx,
		storegen.AcquireLeaseParams{Instance: instance, TtlMillis: ttl.Milliseconds()},
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return notary.Lease{}, false, nil
	}

	if err != nil {
		return notary.Lease{}, false, unavailable(op, err)
	}

	return notary.Lease{Instance: row.Holder, Token: row.Token, ExpiresAt: row.ExpiresAt.Time}, true, nil
}

func (s *Store) InFlight(ctx context.Context) (notary.EpochRecord, bool, error) {
	const op = "postgres.Store.InFlight"

	row, err := s.queries.InFlightEpoch(ctx)

	return s.record(op, row, err)
}

func (s *Store) LastFinalized(ctx context.Context) (notary.EpochRecord, bool, error) {
	const op = "postgres.Store.LastFinalized"

	row, err := s.queries.LastFinalizedEpoch(ctx)

	return s.record(op, row, err)
}

func (s *Store) Epoch(ctx context.Context, number uint64) (notary.EpochRecord, bool, error) {
	const op = "postgres.Store.Epoch"

	key, err := signed(op, number)
	if err != nil {
		return notary.EpochRecord{}, false, err
	}

	row, err := s.queries.EpochByNumber(ctx, key)

	return s.record(op, row, err)
}

func (s *Store) Accounts(ctx context.Context, epoch uint64) ([]notary.Account, error) {
	const op = "postgres.Store.Accounts"

	key, err := signed(op, epoch)
	if err != nil {
		return nil, err
	}

	rows, err := s.queries.EpochAccounts(ctx, key)
	if err != nil {
		return nil, unavailable(op, err)
	}

	accounts := make([]notary.Account, 0, len(rows))
	for _, row := range rows {
		account, err := accountOf(op, row)
		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (s *Store) AccountSecrets(
	ctx context.Context,
	keys []notary.AccountKey,
) (map[notary.AccountKey]notary.AccountSecret, error) {
	const op = "postgres.Store.AccountSecrets"

	return secretsFor(ctx, op, s.queries, keys)
}

func secretsFor(
	ctx context.Context,
	op string,
	q *storegen.Queries,
	keys []notary.AccountKey,
) (map[notary.AccountKey]notary.AccountSecret, error) {
	holders := make([]string, len(keys))
	kinds := make([]int16, len(keys))

	for i, key := range keys {
		holders[i] = key.Holder.String()
		kinds[i] = int16(key.Kind)
	}

	rows, err := q.AccountSecretsFor(ctx, storegen.AccountSecretsForParams{Holders: holders, Kinds: kinds})
	if err != nil {
		return nil, unavailable(op, err)
	}

	secrets := make(map[notary.AccountKey]notary.AccountSecret, len(rows))
	for _, row := range rows {
		key, secret, err := secretOf(op, row)
		if err != nil {
			return nil, err
		}

		secrets[key] = secret
	}

	return secrets, nil
}

func (s *Store) record(op string, row storegen.Epoch, err error) (notary.EpochRecord, bool, error) {
	if errors.Is(err, pgx.ErrNoRows) {
		return notary.EpochRecord{}, false, nil
	}

	if err != nil {
		return notary.EpochRecord{}, false, unavailable(op, err)
	}

	record, err := recordOf(op, row)
	if err != nil {
		return notary.EpochRecord{}, false, err
	}

	return record, true, nil
}

func (s *Store) SaveBuilt(
	ctx context.Context,
	lease notary.Lease,
	ledger notary.LedgerIdentity,
	plan notary.EpochPlan,
) (notary.EpochRecord, error) {
	const op = "postgres.Store.SaveBuilt"

	header := plan.Header
	ints := make([]int64, 0, 5)

	for _, v := range []uint64{header.Number, ledger.ChainID, header.Total, header.LeafCount, header.AsOf.Seconds()} {
		i, err := signed(op, v)
		if err != nil {
			return notary.EpochRecord{}, err
		}

		ints = append(ints, i)
	}

	rows := make([][]any, len(plan.Accounts))
	for i, account := range plan.Accounts {
		position, err := signed(op, account.Position)
		if err != nil {
			return notary.EpochRecord{}, err
		}

		amount, err := signed(op, account.Liability.Amount)
		if err != nil {
			return notary.EpochRecord{}, err
		}

		var id []byte
		if !account.Secret.IsZero() {
			id = account.ID[:]
		}

		rows[i] = []any{
			ints[0],
			position,
			account.Liability.Holder.String(),
			int16(account.Liability.Kind),
			amount,
			account.Salt[:],
			account.Commitment[:],
			id,
		}
	}

	err := s.fenced(ctx, op, lease, func(tx pgx.Tx, q *storegen.Queries) error {
		if err := keepSecrets(ctx, op, q, plan.Accounts); err != nil {
			return err
		}

		if err := q.InsertEpoch(ctx, storegen.InsertEpochParams{
			Number:         ints[0],
			ChainID:        ints[1],
			LedgerAddress:  ledger.Address[:],
			Scheme:         int32(s.scheme),
			PreviousDigest: header.PreviousDigest[:],
			Root:           header.Root[:],
			Digest:         plan.Digest[:],
			Total:          ints[2],
			LeafCount:      ints[3],
			AsOf:           ints[4],
			Currency:       plan.Currency,
		}); err != nil {
			return err
		}

		copied, err := tx.CopyFrom(
			ctx,
			pgx.Identifier{"accounts"},
			[]string{"epoch_number", "position", "holder", "kind", "amount", "salt", "commitment", "account"},
			pgx.CopyFromRows(rows),
		)
		if err != nil {
			return notary.Errf(notary.KindStoreUnavailable, op, err, "copying the accounts of epoch %d", header.Number)
		}

		if copied != int64(len(rows)) {
			return notary.Errf(notary.KindStoreUnavailable, op, nil, "wrote %d of %d accounts", copied, len(rows))
		}

		return nil
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolation {
			return notary.EpochRecord{}, notary.Errf(
				notary.KindLedgerDiverged,
				op,
				err,
				"epoch %d collides with one this database already holds",
				header.Number,
			)
		}

		return notary.EpochRecord{}, err
	}

	return notary.EpochRecord{
		Ledger:   ledger,
		Header:   header,
		Digest:   plan.Digest,
		Currency: plan.Currency,
		State:    notary.StateBuilt,
	}, nil
}

func keepSecrets(ctx context.Context, op string, q *storegen.Queries, accounts []notary.Account) error {
	params := storegen.InsertAccountSecretsParams{}
	keys := make([]notary.AccountKey, 0, len(accounts))

	for _, account := range accounts {
		if account.Secret.IsZero() {
			continue
		}

		params.Holders = append(params.Holders, account.Liability.Holder.String())
		params.Kinds = append(params.Kinds, int16(account.Liability.Kind))
		params.Secrets = append(params.Secrets, account.Secret[:])
		params.Accounts = append(params.Accounts, account.ID[:])
		keys = append(keys, notary.AccountKey{Holder: account.Liability.Holder, Kind: account.Liability.Kind})
	}

	if len(keys) == 0 {
		return nil
	}

	if err := q.InsertAccountSecrets(ctx, params); err != nil {
		return notary.Errf(notary.KindStoreUnavailable, op, err, "keeping the account secrets")
	}

	stored, err := secretsFor(ctx, op, q, keys)
	if err != nil {
		return err
	}

	for _, account := range accounts {
		if account.Secret.IsZero() {
			continue
		}

		key := notary.AccountKey{Holder: account.Liability.Holder, Kind: account.Liability.Kind}
		if stored[key] != account.Secret {
			return notary.Errf(
				notary.KindDisagreement,
				op,
				nil,
				"the database keeps another secret for %s of %s than the epoch was built with",
				key.Kind,
				key.Holder,
			)
		}
	}

	return nil
}

func (s *Store) RecordAttempt(ctx context.Context, lease notary.Lease, attempt notary.Attempt) error {
	const op = "postgres.Store.RecordAttempt"

	anchor := attempt.Anchor
	ints := make([]int64, 0, 5)

	for _, v := range []uint64{attempt.Epoch, anchor.Nonce, anchor.GasLimit, anchor.Fees.MaxFeePerGas, anchor.Fees.MaxPriorityFeePerGas} {
		i, err := signed(op, v)
		if err != nil {
			return err
		}

		ints = append(ints, i)
	}

	return s.fenced(ctx, op, lease, func(_ pgx.Tx, q *storegen.Queries) error {
		if err := q.InsertAttempt(ctx, storegen.InsertAttemptParams{
			TxHash:               anchor.Hash[:],
			EpochNumber:          ints[0],
			Nonce:                ints[1],
			RawTx:                anchor.Raw,
			GasLimit:             ints[2],
			MaxFeePerGas:         ints[3],
			MaxPriorityFeePerGas: ints[4],
			SentAt:               pgtype.Timestamptz{Time: attempt.SentAt, Valid: true},
		}); err != nil {
			return err
		}

		_, err := q.TransitionEpoch(ctx, storegen.TransitionEpochParams{
			ToState:   notary.StateSubmitted.String(),
			Number:    ints[0],
			FromState: notary.StateBuilt.String(),
		})

		return err
	})
}

func (s *Store) Attempts(ctx context.Context, epoch uint64) ([]notary.Attempt, error) {
	const op = "postgres.Store.Attempts"

	number, err := signed(op, epoch)
	if err != nil {
		return nil, err
	}

	rows, err := s.queries.EpochAttempts(ctx, number)
	if err != nil {
		return nil, unavailable(op, err)
	}

	attempts := make([]notary.Attempt, 0, len(rows))
	for _, row := range rows {
		attempt, err := attemptOf(op, epoch, row)
		if err != nil {
			return nil, err
		}

		attempts = append(attempts, attempt)
	}

	return attempts, nil
}

func attemptOf(op string, epoch uint64, row storegen.Attempt) (notary.Attempt, error) {
	hash, err := digestOf(op, row.TxHash)
	if err != nil {
		return notary.Attempt{}, err
	}

	outcome, err := notary.ParseAttemptOutcome(row.Outcome)
	if err != nil {
		return notary.Attempt{}, notary.Errf(notary.KindStoreUnavailable, op, err, "")
	}

	values := make([]uint64, 0, 4)

	for _, v := range []int64{row.Nonce, row.GasLimit, row.MaxFeePerGas, row.MaxPriorityFeePerGas} {
		u, err := unsigned(op, v)
		if err != nil {
			return notary.Attempt{}, err
		}

		values = append(values, u)
	}

	receipt, err := receiptOf(op, hash, row.BlockNumber, row.GasUsed, row.EffectiveGasPrice, row.L1FeeWei)
	if err != nil {
		return notary.Attempt{}, err
	}

	if receipt != nil {
		receipt.Succeeded = outcome == notary.AttemptIncluded
	}

	return notary.Attempt{
		Epoch: epoch,
		Anchor: notary.PreparedAnchor{
			Hash:     hash,
			Raw:      row.RawTx,
			Nonce:    values[0],
			GasLimit: values[1],
			Fees:     notary.FeeQuote{MaxFeePerGas: values[2], MaxPriorityFeePerGas: values[3]},
		},
		Outcome: outcome,
		SentAt:  row.SentAt.Time,
		Receipt: receipt,
	}, nil
}

func (s *Store) SettleAttempt(
	ctx context.Context,
	lease notary.Lease,
	hash notary.Digest,
	outcome notary.AttemptOutcome,
	receipt *notary.AnchorReceipt,
) error {
	const op = "postgres.Store.SettleAttempt"

	settled := outcome == notary.AttemptIncluded || outcome == notary.AttemptReverted
	if settled != (receipt != nil) {
		return notary.Errf(
			notary.KindInvalidArgument,
			op,
			nil,
			"an %s attempt carries a receipt exactly when it was mined",
			outcome,
		)
	}

	params := storegen.SettleAttemptParams{Outcome: outcome.String(), TxHash: hash[:]}

	if receipt != nil {
		fields := make([]*int64, 0, 4)

		for _, v := range []uint64{receipt.BlockNumber, receipt.GasUsed, receipt.EffectiveGasPrice, receipt.L1FeeWei} {
			i, err := signed(op, v)
			if err != nil {
				return err
			}

			fields = append(fields, &i)
		}

		params.BlockNumber, params.GasUsed, params.EffectiveGasPrice, params.L1FeeWei = fields[0], fields[1], fields[2], fields[3]
	}

	return s.fenced(ctx, op, lease, func(_ pgx.Tx, q *storegen.Queries) error {
		changed, err := q.SettleAttempt(ctx, params)
		if err != nil {
			return err
		}

		if changed != 1 {
			return notary.Errf(notary.KindInvalidArgument, op, nil, "no attempt %s is recorded", hash.Hex())
		}

		return nil
	})
}

func (s *Store) Transition(
	ctx context.Context,
	lease notary.Lease,
	epoch uint64,
	from, to notary.EpochState,
	reason string,
) error {
	const op = "postgres.Store.Transition"
	if !from.CanBecome(to) {
		return notary.Errf(notary.KindInvalidArgument, op, nil, "an epoch does not go from %s to %s", from, to)
	}

	if (to == notary.StateAbandoned) != (reason != "") {
		return notary.Errf(notary.KindInvalidArgument, op, nil, "a reason is given exactly when an epoch is abandoned")
	}

	number, err := signed(op, epoch)
	if err != nil {
		return err
	}

	var why *string
	if reason != "" {
		why = &reason
	}

	return s.fenced(ctx, op, lease, func(_ pgx.Tx, q *storegen.Queries) error {
		changed, err := q.TransitionEpoch(ctx, storegen.TransitionEpochParams{
			ToState:   to.String(),
			Reason:    why,
			Number:    number,
			FromState: from.String(),
		})
		if err != nil {
			return err
		}

		if changed != 1 {
			return notary.Errf(notary.KindLeaseLost, op, nil, "epoch %d is no longer %s", epoch, from)
		}

		return nil
	})
}

func (s *Store) AbandonedCount(ctx context.Context) (int64, error) {
	count, err := s.queries.AbandonedCount(ctx)
	if err != nil {
		return 0, unavailable("postgres.Store.AbandonedCount", err)
	}

	return count, nil
}

func (s *Store) RecentCosts(ctx context.Context, limit int32) ([]uint64, error) {
	const op = "postgres.Store.RecentCosts"

	rows, err := s.queries.RecentIncludedReceipts(ctx, limit)
	if err != nil {
		return nil, unavailable(op, err)
	}

	costs := make([]uint64, 0, len(rows))
	for _, row := range rows {
		hash, err := digestOf(op, row.TxHash)
		if err != nil {
			return nil, err
		}

		receipt, err := receiptOf(op, hash, row.BlockNumber, row.GasUsed, row.EffectiveGasPrice, row.L1FeeWei)
		if err != nil || receipt == nil {
			return nil, notary.Errf(notary.KindStoreUnavailable, op, err, "an included attempt has no receipt")
		}

		cost, err := receipt.CostWei()
		if err != nil {
			return nil, err
		}

		costs = append(costs, cost)
	}

	return costs, nil
}

func (s *Store) fenced(
	ctx context.Context,
	op string,
	lease notary.Lease,
	work func(tx pgx.Tx, q *storegen.Queries) error,
) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return unavailable(op, err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	q := s.queries.WithTx(tx)
	if _, err := q.FenceLease(
		ctx,
		storegen.FenceLeaseParams{Instance: lease.Instance, Token: lease.Token},
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return notary.Errf(
				notary.KindLeaseLost,
				op,
				nil,
				"%s no longer holds lease %d",
				lease.Instance,
				lease.Token,
			)
		}

		return unavailable(op, err)
	}

	if err := work(tx, q); err != nil {
		if notary.KindOf(err) != notary.KindUnknown {
			return err
		}

		return unavailable(op, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return unavailable(op, err)
	}

	return nil
}

func unavailable(op string, err error) error {
	return notary.Errf(notary.KindStoreUnavailable, op, err, "")
}
