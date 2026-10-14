package postgres

import (
	"math"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/storegen"
)

func signed(op string, value uint64) (int64, error) {
	if value > math.MaxInt64 {
		return 0, notary.Errf(notary.KindOverflow, op, nil, "%d does not fit a BIGINT column", value)
	}

	return int64(value), nil
}

func unsigned(op string, value int64) (uint64, error) {
	if value < 0 {
		return 0, notary.Errf(notary.KindStoreUnavailable, op, nil, "the database holds a negative %d", value)
	}

	return uint64(value), nil
}

func digestOf(op string, raw []byte) (notary.Digest, error) {
	var d notary.Digest
	if len(raw) != len(d) {
		return d, notary.Errf(notary.KindStoreUnavailable, op, nil, "the database holds a %d-byte digest", len(raw))
	}

	copy(d[:], raw)

	return d, nil
}

func addressOf(op string, raw []byte) (notary.Address, error) {
	var a notary.Address
	if len(raw) != len(a) {
		return a, notary.Errf(notary.KindStoreUnavailable, op, nil, "the database holds a %d-byte address", len(raw))
	}

	copy(a[:], raw)

	return a, nil
}

func recordOf(op string, row storegen.Epoch) (notary.EpochRecord, error) {
	state, err := notary.ParseEpochState(row.State)
	if err != nil {
		return notary.EpochRecord{}, notary.Errf(notary.KindStoreUnavailable, op, err, "")
	}

	values := make([]uint64, 0, 4)

	for _, v := range []int64{row.Number, row.Total, row.LeafCount, row.ChainID} {
		u, err := unsigned(op, v)
		if err != nil {
			return notary.EpochRecord{}, err
		}

		values = append(values, u)
	}

	asOfSeconds, err := unsigned(op, row.AsOf)
	if err != nil {
		return notary.EpochRecord{}, err
	}

	asOf, err := notary.UnixTimeFromSeconds(asOfSeconds)
	if err != nil {
		return notary.EpochRecord{}, notary.Errf(notary.KindStoreUnavailable, op, err, "")
	}

	address, err := addressOf(op, row.LedgerAddress)
	if err != nil {
		return notary.EpochRecord{}, err
	}

	digests := make([]notary.Digest, 0, 3)

	for _, raw := range [][]byte{row.PreviousDigest, row.Root, row.Digest} {
		d, err := digestOf(op, raw)
		if err != nil {
			return notary.EpochRecord{}, err
		}

		digests = append(digests, d)
	}

	reason := ""
	if row.AbandonedReason != nil {
		reason = *row.AbandonedReason
	}

	return notary.EpochRecord{
		Ledger: notary.LedgerIdentity{ChainID: values[3], Address: address},
		Header: notary.EpochHeader{
			Number:         values[0],
			PreviousDigest: digests[0],
			Root:           digests[1],
			Total:          values[1],
			LeafCount:      values[2],
			AsOf:           asOf,
		},
		Digest:          digests[2],
		Currency:        row.Currency,
		State:           state,
		AbandonedReason: reason,
	}, nil
}

func receiptOf(op string, hash notary.Digest, block, gasUsed, price, l1Fee *int64) (*notary.AnchorReceipt, error) {
	if block == nil || gasUsed == nil || price == nil || l1Fee == nil {
		return nil, nil
	}

	values := make([]uint64, 0, 4)

	for _, v := range []int64{*block, *gasUsed, *price, *l1Fee} {
		u, err := unsigned(op, v)
		if err != nil {
			return nil, err
		}

		values = append(values, u)
	}

	return &notary.AnchorReceipt{
		Hash:              hash,
		BlockNumber:       values[0],
		GasUsed:           values[1],
		EffectiveGasPrice: values[2],
		L1FeeWei:          values[3],
	}, nil
}

func accountOf(op string, row storegen.EpochAccountsRow) (notary.Account, error) {
	position, err := unsigned(op, row.Position)
	if err != nil {
		return notary.Account{}, err
	}

	amount, err := unsigned(op, row.Amount)
	if err != nil {
		return notary.Account{}, err
	}

	holder, err := notary.NewHolder(row.Holder)
	if err != nil {
		return notary.Account{}, notary.Errf(notary.KindStoreUnavailable, op, err, "")
	}

	if row.Kind < 0 || row.Kind > 255 {
		return notary.Account{}, notary.Errf(
			notary.KindStoreUnavailable,
			op,
			nil,
			"the database holds kind %d",
			row.Kind,
		)
	}

	kind, err := notary.AccountKindOf(uint8(row.Kind))
	if err != nil {
		return notary.Account{}, notary.Errf(notary.KindStoreUnavailable, op, err, "")
	}

	salt, err := digestOf(op, row.Salt)
	if err != nil {
		return notary.Account{}, err
	}

	commitment, err := digestOf(op, row.Commitment)
	if err != nil {
		return notary.Account{}, err
	}

	account := notary.Account{
		Position:   position,
		Liability:  notary.Liability{Holder: holder, Kind: kind, Amount: amount},
		Salt:       notary.Salt(salt),
		Commitment: commitment,
	}

	if (row.Account == nil) != (row.Secret == nil) {
		return notary.Account{}, notary.Errf(
			notary.KindDisagreement,
			op,
			nil,
			"position %d names an account without its secret",
			position,
		)
	}

	if row.Account == nil {
		return account, nil
	}

	if account.ID, err = digestOf(op, row.Account); err != nil {
		return notary.Account{}, err
	}

	secret, err := digestOf(op, row.Secret)
	if err != nil {
		return notary.Account{}, err
	}

	account.Secret = notary.AccountSecret(secret)

	return account, nil
}

func secretOf(op string, row storegen.AccountSecretsForRow) (notary.AccountKey, notary.AccountSecret, error) {
	holder, err := notary.NewHolder(row.Holder)
	if err != nil {
		return notary.AccountKey{}, notary.AccountSecret{}, notary.Errf(notary.KindStoreUnavailable, op, err, "")
	}

	if row.Kind < 0 || row.Kind > 255 {
		return notary.AccountKey{}, notary.AccountSecret{}, notary.Errf(
			notary.KindStoreUnavailable,
			op,
			nil,
			"the database holds kind %d",
			row.Kind,
		)
	}

	kind, err := notary.AccountKindOf(uint8(row.Kind))
	if err != nil {
		return notary.AccountKey{}, notary.AccountSecret{}, notary.Errf(notary.KindStoreUnavailable, op, err, "")
	}

	secret, err := digestOf(op, row.Secret)
	if err != nil {
		return notary.AccountKey{}, notary.AccountSecret{}, err
	}

	return notary.AccountKey{Holder: holder, Kind: kind}, notary.AccountSecret(secret), nil
}
