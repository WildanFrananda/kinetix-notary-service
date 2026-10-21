package tokens

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type keySetHealth struct {
	next      http.RoundTripper
	reachable atomic.Bool
}

func (h *keySetHealth) RoundTrip(request *http.Request) (*http.Response, error) {
	response, err := h.next.RoundTrip(request)
	if err != nil {
		h.reachable.Store(false)

		return nil, fmt.Errorf("asking identity for its signing keys: %w", err)
	}

	h.reachable.Store(response.StatusCode == http.StatusOK)

	return response, nil
}
