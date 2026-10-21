package ethereum

import "time"

type RPCObserver interface {
	RPCRequest()
	RPCThrottled()
	RPCPaced(waited time.Duration)
}
