package payment

import (
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	paymentv1 "github.com/WildanFrananda/kinetix-notary-service/internal/protogen/payment/v1"
)

type snapshotReader struct {
	currency    string
	asOf        notary.UnixTime
	headed      bool
	trailed     bool
	liabilities []notary.Liability
	sum         uint64
	trailer     *paymentv1.LiabilitySnapshotTrailer
}

func newSnapshotReader(currency string) *snapshotReader {
	return &snapshotReader{currency: currency}
}

func (r *snapshotReader) accept(op string, part *paymentv1.StreamLiabilitySnapshotResponse) error {
	if r.trailed {
		return refused(op, "payment sent more after the trailer")
	}

	switch body := part.GetPart().(type) {
	case *paymentv1.StreamLiabilitySnapshotResponse_Header:
		return r.header(op, body.Header)
	case *paymentv1.StreamLiabilitySnapshotResponse_Accounts:
		return r.accounts(op, body.Accounts)
	case *paymentv1.StreamLiabilitySnapshotResponse_Trailer:
		if !r.headed {
			return refused(op, "payment sent a trailer before a header")
		}

		r.trailed = true
		r.trailer = body.Trailer

		return nil
	default:
		return refused(op, "payment sent an empty part")
	}
}

func (r *snapshotReader) header(op string, header *paymentv1.LiabilitySnapshotHeader) error {
	if r.headed {
		return refused(op, "payment sent a second header")
	}

	if header.GetCurrency() != r.currency {
		return refused(op, "asked for %s, payment answered in %q", r.currency, header.GetCurrency())
	}

	stamp := header.GetAsOf()

	seconds := stamp.GetSeconds()
	if stamp == nil || seconds <= 0 || !stamp.IsValid() {
		return refused(op, "the snapshot carries no valid time")
	}

	asOf, err := notary.UnixTimeFromSeconds(uint64(seconds))
	if err != nil {
		return refused(op, "the snapshot time %v is unusable", stamp.AsTime())
	}

	r.asOf = asOf
	r.headed = true

	return nil
}

func (r *snapshotReader) accounts(op string, batch *paymentv1.LiabilityAccountBatch) error {
	if !r.headed {
		return refused(op, "payment sent accounts before a header")
	}

	for _, account := range batch.GetAccounts() {
		liability, err := r.liability(op, account)
		if err != nil {
			return err
		}

		sum, ok := addAmount(r.sum, liability.Amount)
		if !ok {
			return refused(op, "the accounts sum past 64 bits")
		}

		r.sum = sum
		r.liabilities = append(r.liabilities, liability)
	}

	return nil
}

func (r *snapshotReader) liability(op string, account *paymentv1.LiabilityAccount) (notary.Liability, error) {
	kind, suspense, ok := kindOf(account.GetKind())
	if !ok {
		return notary.Liability{}, refused(op, "payment sent an account of kind %v", account.GetKind())
	}

	var (
		name         string
		namesPurpose bool
	)

	switch holder := account.GetHolder().(type) {
	case *paymentv1.LiabilityAccount_PrincipalId:
		name = holder.PrincipalId
	case *paymentv1.LiabilityAccount_SuspensePurpose:
		name = holder.SuspensePurpose
		namesPurpose = true
	default:
		return notary.Liability{}, refused(op, "a %s account names no holder", kind)
	}

	if namesPurpose != suspense {
		return notary.Liability{}, refused(op, "a %s account names the wrong kind of holder", kind)
	}

	holder, err := notary.NewHolder(name)
	if err != nil {
		return notary.Liability{}, refused(op, "a %s account has an unusable holder", kind)
	}

	balance := account.GetBalance()
	if balance.GetCurrency() != r.currency {
		return notary.Liability{}, refused(op, "a %s balance of %s is in %q", kind, holder, balance.GetCurrency())
	}

	minor := balance.GetAmountMinor()
	if minor < 0 {
		return notary.Liability{}, refused(op, "a %s balance of %s is negative", kind, holder)
	}

	return notary.Liability{Holder: holder, Kind: kind, Amount: uint64(minor)}, nil
}

func (r *snapshotReader) finish(op string) (notary.Snapshot, error) {
	if !r.headed || !r.trailed {
		return notary.Snapshot{}, refused(op, "payment ended the stream without a header and a trailer")
	}

	count := r.trailer.GetAccountCount()
	if count < 0 || uint64(count) != uint64(len(r.liabilities)) {
		return notary.Snapshot{}, refused(op, "payment counted %d accounts and sent %d", count, len(r.liabilities))
	}

	total := r.trailer.GetTotal()

	minor := total.GetAmountMinor()
	if total.GetCurrency() != r.currency || minor < 0 || uint64(minor) != r.sum {
		return notary.Snapshot{}, refused(
			op,
			"payment's total %d %s is not the %d %s its accounts add up to",
			total.GetAmountMinor(),
			total.GetCurrency(),
			r.sum,
			r.currency,
		)
	}

	snapshot, err := notary.NewSnapshot(r.asOf, r.currency, r.liabilities)
	if err != nil {
		return notary.Snapshot{}, notary.Errf(notary.KindSourceRefused, op, err, "")
	}

	return snapshot, nil
}

func kindOf(kind paymentv1.LiabilityKind) (notary.AccountKind, bool, bool) {
	switch kind {
	case paymentv1.LiabilityKind_LIABILITY_KIND_CUSTOMER_BALANCE:
		return notary.AccountCustomerBalance, false, true
	case paymentv1.LiabilityKind_LIABILITY_KIND_MERCHANT_AVAILABLE:
		return notary.AccountMerchantAvailable, false, true
	case paymentv1.LiabilityKind_LIABILITY_KIND_MERCHANT_PENDING_ESCROW:
		return notary.AccountMerchantPendingEscrow, false, true
	case paymentv1.LiabilityKind_LIABILITY_KIND_DRIVER_AVAILABLE:
		return notary.AccountDriverAvailable, false, true
	case paymentv1.LiabilityKind_LIABILITY_KIND_DRIVER_PENDING_ESCROW:
		return notary.AccountDriverPendingEscrow, false, true
	case paymentv1.LiabilityKind_LIABILITY_KIND_SUSPENSE_AVAILABLE:
		return notary.AccountSuspenseAvailable, true, true
	case paymentv1.LiabilityKind_LIABILITY_KIND_SUSPENSE_PENDING_ESCROW:
		return notary.AccountSuspensePendingEscrow, true, true
	default:
		return 0, false, false
	}
}
