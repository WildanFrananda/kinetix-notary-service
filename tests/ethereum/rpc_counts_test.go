package ethereum_test

import (
	"sync"
	"time"
)

type rpcCounts struct {
	mu        sync.Mutex
	requests  int
	throttled int
	paced     time.Duration
}

func (c *rpcCounts) RPCRequest() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.requests++
}

func (c *rpcCounts) RPCThrottled() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.throttled++
}

func (c *rpcCounts) RPCPaced(waited time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.paced += waited
}

func (c *rpcCounts) snapshot() (int, int, time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.requests, c.throttled, c.paced
}
