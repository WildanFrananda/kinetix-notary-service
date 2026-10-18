package notarytest

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/postgres"
)

const DatabaseURLKey = "KINETIX_NOTARY_TEST_DATABASE_URL"

func Database(t *testing.T) *pgxpool.Pool {
	t.Helper()

	url := os.Getenv(DatabaseURLKey)
	if url == "" {
		t.Skip(DatabaseURLKey + " is not set; these tests need Postgres, and CI must set it")
	}

	ctx := context.Background()
	suffix := make([]byte, 6)
	_, err := rand.Read(suffix)
	require.NoError(t, err)

	schema := "e_" + hex.EncodeToString(suffix)
	admin, err := pgx.Connect(ctx, url)
	require.NoError(t, err)
	_, err = admin.Exec(ctx, "CREATE SCHEMA "+schema)
	require.NoError(t, err)
	t.Cleanup(func() {
		_, _ = admin.Exec(context.Background(), "DROP SCHEMA "+schema+" CASCADE")
		_ = admin.Close(context.Background())
	})

	config, err := pgxpool.ParseConfig(url)
	require.NoError(t, err)

	config.ConnConfig.RuntimeParams["search_path"] = schema
	pool, err := pgxpool.NewWithConfig(ctx, config)
	require.NoError(t, err)
	t.Cleanup(pool.Close)
	_, _, err = postgres.Migrate(ctx, pool)
	require.NoError(t, err)

	return pool
}
