package delivering_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/delivering"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

type clock struct {
	mu  sync.Mutex
	now time.Time
}

func (c *clock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.now
}

func (c *clock) Advance(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.now = c.now.Add(d)
}

type deliveryRig struct {
	clock    *clock
	store    *notarytest.StatementStore
	notifier *notarytest.StatementNotifier
	observer *notarytest.StatementObserver
	deliver  *delivering.Deliverer
}

func newDeliveryRig(t *testing.T, holders ...string) *deliveryRig {
	t.Helper()

	c := &clock{now: time.Date(2026, 10, 21, 0, 45, 0, 0, time.UTC)}
	store := &notarytest.StatementStore{}

	for _, holder := range holders {
		store.Rows = append(store.Rows, notarytest.StatementRow{
			Holder:  notarytest.Holder(t, holder),
			Epoch:   187,
			Day:     day(2026, 10, 21),
			State:   "pending",
			Next:    c.Now(),
			Planned: c.Now(),
		})
	}

	notifier := &notarytest.StatementNotifier{}
	observer := &notarytest.StatementObserver{}

	return &deliveryRig{
		clock:    c,
		store:    store,
		notifier: notifier,
		observer: observer,
		deliver:  delivering.NewDeliverer(store, notifier, observer, c.Now, quiet(), policy),
	}
}

func (r *deliveryRig) row(t *testing.T, holder string) notarytest.StatementRow {
	t.Helper()

	for _, row := range r.store.Snapshot() {
		if row.Holder.String() == holder {
			return row
		}
	}

	t.Fatalf("no statement for %s", holder)

	return notarytest.StatementRow{}
}

func TestAStatementThatNotificationAcceptsIsMarkedSentWithItsNotification(t *testing.T) {
	r := newDeliveryRig(t, "principal-a", "principal-b")

	require.NoError(t, r.deliver.Deliver(context.Background()))

	require.Equal(t, "sent", r.row(t, "principal-a").State)
	require.Equal(t, "NTF-1", r.row(t, "principal-a").Notification)
	require.Equal(t, "NTF-2", r.row(t, "principal-b").Notification)
	require.Equal(t, []delivering.Outcome{delivering.OutcomeSent, delivering.OutcomeSent}, r.observer.Outcomes())
	require.Equal(
		t,
		[]notarytest.StatementCall{
			{Holder: notarytest.Holder(t, "principal-a"), Epoch: 187},
			{Holder: notarytest.Holder(t, "principal-b"), Epoch: 187},
		},
		r.notifier.Calls(),
	)

	backlog, seen := r.observer.LastBacklog()
	require.True(t, seen)
	require.Equal(t, notary.StatementBacklog{}, backlog)
}

func TestAHolderWithNowhereToReceiveAStatementIsRefusedNotRetried(t *testing.T) {
	r := newDeliveryRig(t, "principal-a")
	r.notifier.Answer = func(notary.Holder, uint64) (string, error) {
		return "", notary.Errf(notary.KindRecipientUnreachable, "test", nil, "NO_CHANNEL_AVAILABLE: no address")
	}

	require.NoError(t, r.deliver.Deliver(context.Background()))
	r.clock.Advance(24 * time.Hour)
	require.NoError(t, r.deliver.Deliver(context.Background()))

	row := r.row(t, "principal-a")
	require.Equal(t, "refused", row.State)
	require.Contains(t, row.Reason, "NO_CHANNEL_AVAILABLE")
	require.Len(t, r.notifier.Calls(), 1)
	require.Equal(t, []delivering.Outcome{delivering.OutcomeRefused}, r.observer.Outcomes())
}

func TestAnUnansweredStatementWaitsLongerEachTimeUpToTheLastRetry(t *testing.T) {
	r := newDeliveryRig(t, "principal-a")
	r.notifier.Answer = func(notary.Holder, uint64) (string, error) {
		return "", notary.Errf(notary.KindNotifierUnavailable, "test", nil, "notification did not answer")
	}

	start := r.clock.Now()
	require.NoError(t, r.deliver.Deliver(context.Background()))
	require.Equal(t, start.Add(time.Minute), r.row(t, "principal-a").Next)

	require.NoError(t, r.deliver.Deliver(context.Background()))
	require.Len(t, r.notifier.Calls(), 1, "a postponed statement is not due before its next attempt")

	r.clock.Advance(time.Minute)
	require.NoError(t, r.deliver.Deliver(context.Background()))
	require.Equal(t, r.clock.Now().Add(2*time.Minute), r.row(t, "principal-a").Next)

	for range 20 {
		r.clock.Advance(7 * time.Hour)
		require.NoError(t, r.deliver.Deliver(context.Background()))
	}

	row := r.row(t, "principal-a")
	require.Equal(t, "pending", row.State)
	require.Equal(t, r.clock.Now().Add(6*time.Hour), row.Next)

	backlog, _ := r.observer.LastBacklog()
	require.Equal(t, notary.StatementBacklog{Pending: 1, Overdue: 1}, backlog)
}

func TestARefusalFromNotificationIsRetriedRatherThanTakenAsDelivered(t *testing.T) {
	r := newDeliveryRig(t, "principal-a")
	refusing := true
	r.notifier.Answer = func(notary.Holder, uint64) (string, error) {
		if refusing {
			return "", notary.Errf(
				notary.KindNotifierRefused,
				"test",
				nil,
				"TEMPLATE_NOT_GRANTED: notary is not granted",
			)
		}

		return "NTF-7", nil
	}

	require.NoError(t, r.deliver.Deliver(context.Background()))
	require.Equal(t, "pending", r.row(t, "principal-a").State)

	refusing = false

	r.clock.Advance(time.Minute)
	require.NoError(t, r.deliver.Deliver(context.Background()))

	require.Equal(t, "sent", r.row(t, "principal-a").State)
	require.Equal(t, "NTF-7", r.row(t, "principal-a").Notification)
	require.Equal(t, []delivering.Outcome{delivering.OutcomePostponed, delivering.OutcomeSent}, r.observer.Outcomes())
}

func TestATickStopsAtTheFirstStatementNotificationDoesNotAnswer(t *testing.T) {
	r := newDeliveryRig(t, "principal-a", "principal-b", "principal-c")
	r.notifier.Answer = func(notary.Holder, uint64) (string, error) {
		return "", notary.Errf(notary.KindNotifierUnavailable, "test", nil, "notification did not answer")
	}

	require.NoError(t, r.deliver.Deliver(context.Background()))

	require.Len(t, r.notifier.Calls(), 1, "an unanswered call is not followed by more in the same tick")
	require.Equal(t, []delivering.Outcome{delivering.OutcomePostponed}, r.observer.Outcomes())

	backlog, seen := r.observer.LastBacklog()
	require.True(t, seen)
	require.Equal(t, int64(3), backlog.Pending)
}

func TestARefusalAboutOneStatementDoesNotHoldBackTheRest(t *testing.T) {
	r := newDeliveryRig(t, "principal-a", "principal-b")
	r.notifier.Answer = func(holder notary.Holder, _ uint64) (string, error) {
		if holder.String() == "principal-a" {
			return "", notary.Errf(notary.KindNotifierRefused, "test", nil, "TEMPLATE_PARAM_MISSING: epoch")
		}

		return "NTF-b", nil
	}

	require.NoError(t, r.deliver.Deliver(context.Background()))

	require.Equal(t, "pending", r.row(t, "principal-a").State)
	require.Equal(t, "sent", r.row(t, "principal-b").State)
}
