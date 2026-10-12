package anchoring

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type Runner struct {
	anchorer *Anchorer
	every    time.Duration
	clock    func() time.Time
	log      *slog.Logger
	mu       sync.Mutex
	lastTick time.Time
	started  time.Time
}

func NewRunner(anchorer *Anchorer, every time.Duration, clock func() time.Time, log *slog.Logger) *Runner {
	return &Runner{anchorer: anchorer, every: every, clock: clock, log: log}
}

func (r *Runner) Run(ctx context.Context) {
	r.mu.Lock()
	r.started = r.clock()
	r.mu.Unlock()

	ticker := time.NewTicker(r.every)
	defer ticker.Stop()

	for {
		r.once(ctx)

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func (r *Runner) once(ctx context.Context) {
	err := r.anchorer.Tick(ctx)
	r.mu.Lock()
	r.lastTick = r.clock()
	r.mu.Unlock()

	if err == nil || ctx.Err() != nil {
		return
	}

	kind := notary.KindOf(err)
	switch kind {
	case notary.KindAnchorDeferred:
		r.log.InfoContext(ctx, "anchoring deferred", "reason", err.Error())
	case notary.KindSourceUnavailable, notary.KindLedgerUnavailable, notary.KindStoreUnavailable, notary.KindLeaseLost:
		r.log.WarnContext(ctx, "anchoring will retry", "kind", kind.String(), "err", err.Error())
	default:
		r.log.ErrorContext(ctx, "anchoring stopped short", "kind", kind.String(), "err", err.Error())
	}
}

func (r *Runner) Alive() (bool, time.Duration) {
	r.mu.Lock()
	defer r.mu.Unlock()

	last := r.lastTick
	if last.IsZero() {
		last = r.started
	}

	if last.IsZero() {
		return false, 0
	}

	since := r.clock().Sub(last)

	return since <= 3*r.every+time.Minute, since
}
