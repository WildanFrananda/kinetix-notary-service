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
