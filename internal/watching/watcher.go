package watching

import (
	"context"
	"log/slog"
	"time"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type Watcher struct {
	court    notary.Court
	observer Observer
	log      *slog.Logger
	interval time.Duration
	seen     uint64
}

func NewWatcher(court notary.Court, observer Observer, log *slog.Logger, interval time.Duration) *Watcher {
	return &Watcher{court: court, observer: observer, log: log, interval: interval}
}

func (w *Watcher) Run(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		w.Look(ctx)

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func (w *Watcher) Look(ctx context.Context) {
	total, err := w.court.TotalFaults(ctx)
	if err != nil {
		w.log.WarnContext(ctx, "the court could not be read", "kind", notary.KindOf(err).String(), "err", err.Error())
		return
	}

	w.observer.ProvenFaults(total)

	if total > w.seen {
		w.log.ErrorContext(
			ctx,
			"the court holds faults proven against this notary's statements",
			"total",
			total,
			"new",
			total-w.seen,
		)
	}

	w.seen = total
}
