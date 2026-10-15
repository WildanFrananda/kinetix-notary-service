package wiring

import (
	"fmt"
	"log/slog"
	"time"

	"go.uber.org/dig"

	"github.com/WildanFrananda/kinetix-notary-service/internal/config"
	"github.com/WildanFrananda/kinetix-notary-service/internal/metrics"
)

func Daemon(lookup config.Lookup, build metrics.Build, log *slog.Logger) (*dig.Container, error) {
	container := dig.New()

	providers := []any{
		func() config.Lookup { return lookup },
		func() metrics.Build { return build },
		func() *slog.Logger { return log },
		func() func() time.Time { return time.Now },
		provideDaemonConfig,
		provideDaemonLedgerSettings,
		provideChainClient,
		provideScheme,
		provideLedger,
		provideCourt,
		provideDirectory,
		provideDirectoryCheck,
		provideStatementSigner,
		provideAnchorChain,
		provideCouncil,
		provideEpochWatch,
		provideQuorum,
		provideDatabase,
		provideStore,
		provideMetrics,
		provideSnapshotSource,
		provideSalts,
		providePolicy,
		provideAnchorer,
		provideRunner,
		provideProver,
		provideWatcher,
		provideAuthenticator,
		provideProofHandler,
		provideWithholdingCourt,
		provideRegistrar,
		provideRegistrationHandler,
		provideAnswerChain,
		provideAnswerer,
		provideHTTP,
	}
	for _, provider := range providers {
		if err := container.Provide(provider); err != nil {
			return nil, fmt.Errorf("provide %T: %w", provider, err)
		}
	}

	return container, nil
}
