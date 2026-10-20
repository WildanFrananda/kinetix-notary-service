package delivering

import (
	"context"
	"log/slog"
	"time"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type Deliverer struct {
	store    notary.StatementStore
	notifier notary.StatementNotifier
	observer Observer
	clock    func() time.Time
	log      *slog.Logger
	policy   Policy
}

func NewDeliverer(
	store notary.StatementStore,
	notifier notary.StatementNotifier,
	observer Observer,
	clock func() time.Time,
	log *slog.Logger,
	policy Policy,
) *Deliverer {
	return &Deliverer{
		store:    store,
		notifier: notifier,
		observer: observer,
		clock:    clock,
		log:      log,
		policy:   policy,
	}
}

func (d *Deliverer) Deliver(ctx context.Context) error {
	due, err := d.store.DueStatements(ctx, d.clock(), d.policy.Batch)
	if err != nil {
		return err
	}

	for _, delivery := range due {
		answered, err := d.deliver(ctx, delivery)
		if err != nil {
			return err
		}

		if !answered {
			break
		}
	}

	backlog, err := d.store.StatementBacklog(ctx, d.clock().Add(-d.policy.OverdueAfter))
	if err != nil {
		return err
	}

	d.observer.StatementBacklog(backlog)

	return nil
}

func (d *Deliverer) deliver(ctx context.Context, delivery notary.StatementDelivery) (bool, error) {
	notification, err := d.notifier.NotifyStatement(ctx, delivery.Holder, delivery.Epoch)
	if err == nil {
		if err := d.store.StatementSent(ctx, delivery, notification); err != nil {
			return false, err
		}

		d.observer.Statement(OutcomeSent)
		d.log.InfoContext(ctx, "statement sent", "epoch", delivery.Epoch, "notification", notification)

		return true, nil
	}

	kind := notary.KindOf(err)
	if kind == notary.KindRecipientUnreachable {
		if err := d.store.StatementRefused(ctx, delivery, err.Error()); err != nil {
			return false, err
		}

		d.observer.Statement(OutcomeRefused)
		d.log.WarnContext(ctx, "statement has nowhere to go", "epoch", delivery.Epoch, "err", err.Error())

		return true, nil
	}

	next := d.clock().Add(d.backoff(delivery.Attempts))
	if err := d.store.StatementPostponed(ctx, delivery, next); err != nil {
		return false, err
	}

	d.observer.Statement(OutcomePostponed)

	if kind == notary.KindNotifierUnavailable {
		d.log.WarnContext(
			ctx,
			"statement delivery will retry",
			"epoch",
			delivery.Epoch,
			"next_attempt",
			next.Format(time.RFC3339),
			"err",
			err.Error(),
		)

		return false, nil
	}

	d.log.ErrorContext(
		ctx,
		"notification refused a statement; it is retried until the refusal is fixed",
		"epoch",
		delivery.Epoch,
		"kind",
		kind.String(),
		"next_attempt",
		next.Format(time.RFC3339),
		"err",
		err.Error(),
	)

	return true, nil
}

func (d *Deliverer) backoff(attempts uint32) time.Duration {
	wait := d.policy.FirstRetry
	for range attempts {
		if wait >= d.policy.LastRetry/2 {
			return d.policy.LastRetry
		}

		wait *= 2
	}

	return min(wait, d.policy.LastRetry)
}
