package delivering

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type Runner struct {
	planner   *Planner
	deliverer *Deliverer
	every     time.Duration
	clock     func() time.Time
	log       *slog.Logger
	mu        sync.Mutex
	lastTick  time.Time
	started   time.Time
}

func NewRunner(
	planner *Planner,
	deliverer *Deliverer,
	every time.Duration,
	clock func() time.Time,
	log *slog.Logger,
) *Runner {
	return &Runner{planner: planner, deliverer: deliverer, every: every, clock: clock, log: log}
}

func (r *Runner) Run(ctx context.Context) {
	r.mu.Lock()
	r.started = r.clock()
	r.mu.Unlock()

	ticker := time.NewTicker(r.every)
	defer ticker.Stop()

	for {
		r.Tick(ctx)

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func (r *Runner) Tick(ctx context.Context) {
	r.report(ctx, "statement planning", r.planner.Plan(ctx))
	r.report(ctx, "statement delivery", r.deliverer.Deliver(ctx))

	r.mu.Lock()
	r.lastTick = r.clock()
	r.mu.Unlock()
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

func (r *Runner) report(ctx context.Context, step string, err error) {
	if err == nil || ctx.Err() != nil {
		return
	}

	kind := notary.KindOf(err)
	if kind == notary.KindStoreUnavailable {
		r.log.WarnContext(ctx, step+" will retry", "kind", kind.String(), "err", err.Error())
		return
	}

	r.log.ErrorContext(ctx, step+" stopped short", "kind", kind.String(), "err", err.Error())
}
