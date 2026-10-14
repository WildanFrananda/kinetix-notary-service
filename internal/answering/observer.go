package answering

import "math/big"

type Observer interface {
	Silences(total uint64)
	OpenRequests(count int)
	Unanswerable(count int)
	AnswererBalance(wei *big.Int)
}
