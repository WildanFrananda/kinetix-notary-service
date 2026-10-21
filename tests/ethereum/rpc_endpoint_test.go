package ethereum_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

const answerWord = "0x0000000000000000000000000000000000000000000000000000000000000001"

type rpcEndpoint struct {
	URL string

	mu       sync.Mutex
	served   int
	refuse   int
	fail     int
	revert   bool
	chainHex string
}

func newRPCEndpoint(t *testing.T) *rpcEndpoint {
	t.Helper()

	e := &rpcEndpoint{chainHex: "0x14a34"}
	server := httptest.NewServer(http.HandlerFunc(e.serve))
	t.Cleanup(server.Close)
	e.URL = server.URL

	return e
}

func (e *rpcEndpoint) Served() int {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.served
}

func (e *rpcEndpoint) Refuse(n int) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.refuse = n
}

func (e *rpcEndpoint) Fail(n int) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.fail = n
}

func (e *rpcEndpoint) Revert() {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.revert = true
}

func (e *rpcEndpoint) serve(w http.ResponseWriter, r *http.Request) {
	var call struct {
		ID     json.RawMessage `json:"id"`
		Method string          `json:"method"`
	}
	if err := json.NewDecoder(r.Body).Decode(&call); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	e.mu.Lock()
	e.served++
	refused := e.refuse > 0
	failed := !refused && e.fail > 0
	reverted := e.revert

	switch {
	case refused:
		e.refuse--
	case failed:
		e.fail--
	}
	e.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")

	switch {
	case refused:
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = fmt.Fprintf(
			w,
			`{"jsonrpc":"2.0","id":%s,"error":{"code":-32007,"message":"25/second request limit reached"}}`,
			call.ID,
		)
	case failed:
		w.WriteHeader(http.StatusBadGateway)
	case reverted && call.Method == "eth_call":
		_, _ = fmt.Fprintf(
			w,
			`{"jsonrpc":"2.0","id":%s,"error":{"code":3,"message":"execution reverted","data":"0x"}}`,
			call.ID,
		)
	case call.Method == "eth_chainId":
		_, _ = fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%s,"result":"%s"}`, call.ID, e.chainHex)
	case call.Method == "eth_call":
		_, _ = fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%s,"result":"%s"}`, call.ID, answerWord)
	default:
		_, _ = fmt.Fprintf(
			w,
			`{"jsonrpc":"2.0","id":%s,"error":{"code":-32601,"message":"no method %s"}}`,
			call.ID,
			call.Method,
		)
	}
}
