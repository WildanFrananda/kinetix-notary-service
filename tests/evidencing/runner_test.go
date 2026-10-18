package evidencing_test

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/evidencing"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

func TestTheRunnerTicksAndReportsWhetherItIsAlive(t *testing.T) {
	r := newRig(t)
	runner := evidencing.NewRunner(r.ticker, time.Hour, r.clock, slog.New(slog.NewTextHandler(io.Discard, nil)))

	alive, _ := runner.Alive()
	require.False(t, alive, "alive before it ever ran")

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})

	go func() {
		runner.Run(ctx)
		close(done)
	}()

	require.Eventually(
		t,
		func() bool { return r.state(t, notary.StreamOrderLifecycle, 0) == "submitted" },
		10*time.Second,
		20*time.Millisecond,
	)
	cancel()
	<-done

	alive, since := runner.Alive()
	require.True(t, alive)
	require.Zero(t, since)

	r.advance(4 * time.Hour)

	alive, _ = runner.Alive()
	require.False(t, alive)
}
