package wiring

import (
	"fmt"
	"log/slog"

	"go.uber.org/dig"

	"github.com/WildanFrananda/kinetix-notary-service/internal/config"
	"github.com/WildanFrananda/kinetix-notary-service/internal/metrics"
)

func Watch(lookup config.Lookup, build metrics.Build, log *slog.Logger) (*dig.Container, error) {
	container := dig.New()

	providers := []any{
		func() config.Lookup { return lookup },
		func() metrics.Build { return build },
		func() *slog.Logger { return log },
		provideWatchConfig,
		provideWatchLedgerSettings,
		provideChainClient,
		provideScheme,
		provideWatchMetrics,
		provideHeadReader,
		provideCouncil,
		provideWatchSigner,
		provideSnapshotRecords,
		provideApprover,
		provideWatchHTTP,
		provideWatchGRPC,
	}
	for _, provider := range providers {
		if err := container.Provide(provider); err != nil {
			return nil, fmt.Errorf("provide %T: %w", provider, err)
		}
	}

	return container, nil
}
