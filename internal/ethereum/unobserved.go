package ethereum

import "time"

type unobserved struct{}

func (unobserved) RPCRequest() {}

func (unobserved) RPCThrottled() {}

func (unobserved) RPCPaced(time.Duration) {}
