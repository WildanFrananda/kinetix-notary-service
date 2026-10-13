package watching_test

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/watching"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

type recorder struct {
	seen []uint64
}

func (r *recorder) ProvenFaults(total uint64) {
	r.seen = append(r.seen, total)
}

func TestTheGaugeFollowsTheCourtAndANewFaultIsAnError(t *testing.T) {
	court := &notarytest.Court{}
	observer := &recorder{}

	var logs bytes.Buffer

	watcher := watching.NewWatcher(court, observer, slog.New(slog.NewJSONHandler(&logs, nil)), 0)

	watcher.Look(context.Background())
	require.NotContains(t, logs.String(), `"level":"ERROR"`)

	court.Total = 2

	watcher.Look(context.Background())
	require.Contains(t, logs.String(), `"level":"ERROR"`)
	require.Contains(t, logs.String(), `"new":2`)

	logs.Reset()
	watcher.Look(context.Background())
	require.NotContains(t, logs.String(), `"level":"ERROR"`)

	require.Equal(t, []uint64{0, 2, 2}, observer.seen)
}

func TestACourtThatCannotBeReadLeavesTheGaugeAlone(t *testing.T) {
	court := &notarytest.Court{TotalErr: notary.Errf(notary.KindLedgerUnavailable, "test", nil, "down")}
	observer := &recorder{}

	var logs bytes.Buffer

	watching.NewWatcher(court, observer, slog.New(slog.NewJSONHandler(&logs, nil)), 0).Look(context.Background())

	require.Empty(t, observer.seen)
	require.Contains(t, logs.String(), `"level":"WARN"`)
	require.Contains(t, logs.String(), "ledger_unavailable")
}
