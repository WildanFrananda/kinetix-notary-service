package wiring

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/dig"

	"github.com/WildanFrananda/kinetix-notary-service/internal/config"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/postgres"
)

type databaseParams struct {
	dig.In
	Config config.Daemon
}

func provideDatabase(p databaseParams) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return postgres.Open(ctx, p.Config.DatabaseURL.Reveal())
}

type storeParams struct {
	dig.In
	Pool   *pgxpool.Pool
	Scheme notary.Scheme
}

type storeResult struct {
	dig.Out
	Store          notary.EpochStore
	Proofs         notary.ProofSource
	Registrations  notary.RegistrationStore
	Evidence       notary.EvidenceStore
	EvidenceProofs notary.EvidenceProofSource
	Statements     notary.StatementStore
	Withholding    notary.WithholdingScans
}

func provideStore(p storeParams) storeResult {
	store := postgres.NewStore(p.Pool, p.Scheme)

	return storeResult{
		Store:          store,
		Proofs:         store,
		Registrations:  store,
		Evidence:       store,
		EvidenceProofs: store,
		Statements:     store,
		Withholding:    store,
	}
}
