package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	"github.com/WildanFrananda/kinetix-notary-service/internal/migrations"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

func Open(ctx context.Context, url string) (*pgxpool.Pool, error) {
	const op = "postgres.Open"

	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, notary.Errf(notary.KindInvalidArgument, op, nil, "the database URL could not be parsed")
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, notary.Errf(notary.KindStoreUnavailable, op, err, "reaching the database")
	}

	return pool, nil
}

func Migrate(ctx context.Context, pool *pgxpool.Pool) (int64, int64, error) {
	const op = "postgres.Migrate"

	db := stdlib.OpenDBFromPool(pool)
	defer func() { _ = db.Close() }()

	provider, err := goose.NewProvider(goose.DialectPostgres, db, migrations.FS)
	if err != nil {
		return 0, 0, notary.Errf(notary.KindStoreUnavailable, op, err, "loading migrations")
	}

	before, err := provider.GetDBVersion(ctx)
	if err != nil {
		return 0, 0, notary.Errf(notary.KindStoreUnavailable, op, err, "reading the schema version")
	}

	if _, err := provider.Up(ctx); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return 0, 0, notary.Errf(notary.KindStoreUnavailable, op, err, "applying migrations")
	}

	after, err := provider.GetDBVersion(ctx)
	if err != nil {
		return 0, 0, notary.Errf(notary.KindStoreUnavailable, op, err, "reading the schema version")
	}

	return before, after, nil
}
