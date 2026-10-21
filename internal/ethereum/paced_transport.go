package ethereum

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

const (
	throttleRetries = 3
	throttlePause   = time.Second
)

type pacedTransport struct {
	next     http.RoundTripper
	pace     *rate.Limiter
	observer RPCObserver
}

func (t *pacedTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	ctx := request.Context()
	attempt := request

	for retry := 0; ; retry++ {
		if err := t.wait(ctx); err != nil {
			return nil, err
		}

		t.observer.RPCRequest()

		response, err := t.next.RoundTrip(attempt)
		if err != nil {
			return nil, fmt.Errorf("rpc request: %w", err)
		}

		if response.StatusCode != http.StatusTooManyRequests {
			return response, nil
		}

		t.observer.RPCThrottled()

		if retry == throttleRetries || request.GetBody == nil {
			return response, nil
		}

		if deadline, ok := ctx.Deadline(); ok && time.Until(deadline) < throttlePause {
			return response, nil
		}

		_, _ = io.Copy(io.Discard, response.Body)
		_ = response.Body.Close()

		if !pause(ctx, throttlePause) {
			return nil, fmt.Errorf("rpc request after a 429: %w", ctx.Err())
		}

		body, err := request.GetBody()
		if err != nil {
			return nil, fmt.Errorf("rpc request body: %w", err)
		}

		attempt = request.Clone(ctx)
		attempt.Body = body
	}
}

func (t *pacedTransport) wait(ctx context.Context) error {
	if t.pace == nil {
		return nil
	}

	started := time.Now()

	if err := t.pace.Wait(ctx); err != nil {
		return fmt.Errorf("rpc request rate: %w", err)
	}

	t.observer.RPCPaced(time.Since(started))

	return nil
}

func pause(ctx context.Context, d time.Duration) bool {
	timer := time.NewTimer(d)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return false
	case <-timer.C:
		return true
	}
}
