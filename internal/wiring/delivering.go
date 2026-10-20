package wiring

import (
	"log/slog"
	"time"

	"go.uber.org/dig"

	"github.com/WildanFrananda/kinetix-notary-service/internal/config"
	"github.com/WildanFrananda/kinetix-notary-service/internal/delivering"
	"github.com/WildanFrananda/kinetix-notary-service/internal/mesh"
	"github.com/WildanFrananda/kinetix-notary-service/internal/metrics"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notification"
)

const (
	statementEpochsPerTick = 50
	statementBatch         = 20
	statementFirstRetry    = time.Minute
	statementLastRetry     = 6 * time.Hour
	statementOverdueAfter  = 6 * time.Hour
)

type notifierParams struct {
	dig.In
	Config  config.Daemon
	Metrics *metrics.Metrics
}

type notifierResult struct {
	dig.Out
	Notifier   *notification.Notifier
	Statements notary.StatementNotifier
}

func provideNotifier(p notifierParams) (notifierResult, error) {
	notifier, err := notification.Dial(
		notification.Settings{
			Endpoint:  p.Config.Notification,
			Identity:  mesh.Identity{PKIDir: p.Config.PKIDir, TrustDomains: p.Config.TrustDomains},
			Timeout:   p.Config.NotifyTimeout,
			VerifyURL: p.Config.VerifyURL,
		},
		p.Metrics.ClientInterceptors("notification")...,
	)
	if err != nil {
		return notifierResult{}, err
	}

	return notifierResult{Notifier: notifier, Statements: notifier}, nil
}

type statementsParams struct {
	dig.In
	Store    notary.StatementStore
	Notifier notary.StatementNotifier
	Observer delivering.Observer
	Clock    func() time.Time
	Log      *slog.Logger
	Config   config.Daemon
}

func provideStatements(p statementsParams) *delivering.Runner {
	policy := delivering.Policy{
		UTCOffset:     p.Config.UTCOffset,
		EpochsPerTick: statementEpochsPerTick,
		Batch:         statementBatch,
		FirstRetry:    statementFirstRetry,
		LastRetry:     statementLastRetry,
		OverdueAfter:  statementOverdueAfter,
	}

	return delivering.NewRunner(
		delivering.NewPlanner(p.Store, p.Clock, p.Log, policy),
		delivering.NewDeliverer(p.Store, p.Notifier, p.Observer, p.Clock, p.Log, policy),
		p.Config.Tick,
		p.Clock,
		p.Log,
	)
}
