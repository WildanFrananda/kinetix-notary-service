package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/WildanFrananda/kinetix-notary-service/internal/config"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/postgres"
)

func main() {
	os.Exit(run())
}

func run() int {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := migrate(ctx, log); err != nil {
		log.Error("migration failed", "kind", notary.KindOf(err).String(), "err", err.Error())
		return 1
	}

	return 0
}

func migrate(ctx context.Context, log *slog.Logger) error {
	url, err := config.LoadMigrator(os.LookupEnv)
	if err != nil {
		return notary.Errf(notary.KindInvalidArgument, "notary-migrate", err, "")
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	pool, err := postgres.Open(ctx, url.Reveal())
	if err != nil {
		return err
	}
	defer pool.Close()

	before, after, err := postgres.Migrate(ctx, pool)
	if err != nil {
		return err
	}

	log.InfoContext(ctx, "migrations applied", "from", before, "to", after)

	return nil
}
