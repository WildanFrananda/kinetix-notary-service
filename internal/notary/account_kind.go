package notary

type AccountKind uint8

const (
	AccountCustomerBalance AccountKind = iota + 1
	AccountMerchantAvailable
	AccountMerchantPendingEscrow
	AccountDriverAvailable
	AccountDriverPendingEscrow
	AccountSuspenseAvailable
	AccountSuspensePendingEscrow
)

func AccountKindOf(code uint8) (AccountKind, error) {
	kind := AccountKind(code)
	if !kind.Valid() {
		return 0, Errf(KindInvalidArgument, "notary.AccountKindOf", nil, "%d is not an account kind", code)
	}

	return kind, nil
}

func (k AccountKind) Valid() bool {
	return k >= AccountCustomerBalance && k <= AccountSuspensePendingEscrow
}

func (k AccountKind) String() string {
	switch k {
	case AccountCustomerBalance:
		return "customer_balance"
	case AccountMerchantAvailable:
		return "merchant_available"
	case AccountMerchantPendingEscrow:
		return "merchant_pending_escrow"
	case AccountDriverAvailable:
		return "driver_available"
	case AccountDriverPendingEscrow:
		return "driver_pending_escrow"
	case AccountSuspenseAvailable:
		return "suspense_available"
	case AccountSuspensePendingEscrow:
		return "suspense_pending_escrow"
	default:
		return "invalid"
	}
}
