package notary

import (
	"math"
	"math/big"
	"time"
)

func Runway(balance *big.Int, costPerEpoch uint64, interval time.Duration) (time.Duration, bool) {
	if balance == nil || balance.Sign() < 0 || costPerEpoch == 0 || interval <= 0 {
		return 0, false
	}

	epochs := new(big.Int).Quo(balance, new(big.Int).SetUint64(costPerEpoch))

	limit := big.NewInt(int64(math.MaxInt64 / interval))
	if epochs.Cmp(limit) > 0 {
		return time.Duration(math.MaxInt64), true
	}

	return time.Duration(epochs.Int64()) * interval, true
}
