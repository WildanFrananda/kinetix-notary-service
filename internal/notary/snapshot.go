package notary

import "math"

type Snapshot struct {
	asOf        UnixTime
	currency    string
	liabilities []Liability
	total       uint64
}

func NewSnapshot(asOf UnixTime, currency string, liabilities []Liability) (Snapshot, error) {
	const op = "notary.NewSnapshot"
	if asOf.IsZero() {
		return Snapshot{}, Errf(KindInvalidArgument, op, nil, "a snapshot is taken at a time")
	}

	if !isISOCurrency(currency) {
		return Snapshot{}, Errf(KindInvalidArgument, op, nil, "%q is not an ISO 4217 currency", currency)
	}

	if uint64(len(liabilities)) > MaxLeaves {
		return Snapshot{}, Errf(
			KindTooManyLeaves,
			op,
			nil,
			"%d accounts exceed the scheme's %d",
			len(liabilities),
			MaxLeaves,
		)
	}

	seen := make(map[string]struct{}, len(liabilities))

	var total uint64

	for _, liability := range liabilities {
		if liability.Holder.IsZero() || !liability.Kind.Valid() {
			return Snapshot{}, Errf(KindInvalidArgument, op, nil, "an account needs a holder and a valid kind")
		}

		k := string([]byte{byte(liability.Kind)}) + liability.Holder.String()
		if _, duplicate := seen[k]; duplicate {
			return Snapshot{}, Errf(
				KindInvalidArgument,
				op,
				nil,
				"%s of %s appears twice",
				liability.Kind,
				liability.Holder,
			)
		}

		seen[k] = struct{}{}

		if liability.Amount > math.MaxInt64-total {
			return Snapshot{}, Errf(KindOverflow, op, nil, "the total passes %d minor units", uint64(math.MaxInt64))
		}

		total += liability.Amount
	}

	owned := make([]Liability, len(liabilities))
	copy(owned, liabilities)

	return Snapshot{asOf: asOf, currency: currency, liabilities: owned, total: total}, nil
}

func (s Snapshot) AsOf() UnixTime {
	return s.asOf
}

func (s Snapshot) Currency() string {
	return s.currency
}

func (s Snapshot) Liabilities() []Liability {
	out := make([]Liability, len(s.liabilities))
	copy(out, s.liabilities)

	return out
}

func (s Snapshot) Total() uint64 {
	return s.total
}

func isISOCurrency(code string) bool {
	if len(code) != 3 {
		return false
	}

	for i := range len(code) {
		if code[i] < 'A' || code[i] > 'Z' {
			return false
		}
	}

	return true
}
