package evidencing

import (
	"context"
	"log/slog"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type Evidencer struct {
	store     notary.EvidenceStore
	collector *Collector
	closer    *Closer
	log       *slog.Logger
	policy    Policy
}

func NewEvidencer(
	store notary.EvidenceStore,
	collector *Collector,
	closer *Closer,
	log *slog.Logger,
	policy Policy,
) *Evidencer {
	return &Evidencer{store: store, collector: collector, closer: closer, log: log, policy: policy}
}

func (e *Evidencer) Tick(ctx context.Context) error {
	lease, held, err := e.store.Lease(ctx, e.policy.Instance, e.policy.LeaseTerm)
	if err != nil {
		return err
	}

	if !held {
		return nil
	}

	if err := e.collector.Collect(ctx, lease); err != nil {
		return err
	}

	work := e.closer.Close(ctx, lease)

	reported := e.closer.Report(ctx)
	if work == nil {
		return reported
	}

	if reported != nil {
		e.log.WarnContext(ctx, "evidence reporting failed after a failed step", "err", reported.Error())
	}

	return work
}
