package notary

import "math"

type FeeQuote struct {
	MaxFeePerGas         uint64
	MaxPriorityFeePerGas uint64
}

func QuoteFees(baseFee, tip, ceiling uint64) (FeeQuote, error) {
	const op = "notary.QuoteFees"
	if ceiling == 0 {
		return FeeQuote{}, Errf(KindInvalidArgument, op, nil, "a fee ceiling of zero admits no transaction")
	}

	if baseFee > ceiling {
		return FeeQuote{}, Errf(
			KindAnchorDeferred,
			op,
			nil,
			"the base fee %d wei is above the ceiling of %d wei per gas",
			baseFee,
			ceiling,
		)
	}

	wanted := uint64(math.MaxUint64)
	if baseFee <= (math.MaxUint64-tip)/2 {
		wanted = 2*baseFee + tip
	}

	maxFee := min(wanted, ceiling)

	return FeeQuote{MaxFeePerGas: maxFee, MaxPriorityFeePerGas: min(tip, maxFee-baseFee)}, nil
}

func (q FeeQuote) Bumped(ceiling uint64) (FeeQuote, bool) {
	next := FeeQuote{
		MaxFeePerGas:         raiseByEighth(q.MaxFeePerGas),
		MaxPriorityFeePerGas: raiseByEighth(q.MaxPriorityFeePerGas),
	}
	if next.MaxFeePerGas > ceiling {
		return q, false
	}

	return next, true
}

func raiseByEighth(value uint64) uint64 {
	step := value/8 + 1
	if value > math.MaxUint64-step {
		return math.MaxUint64
	}

	return value + step
}
