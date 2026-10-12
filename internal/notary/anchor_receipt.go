package notary

import "math/bits"

type AnchorReceipt struct {
	Hash              Digest
	Succeeded         bool
	BlockNumber       uint64
	GasUsed           uint64
	EffectiveGasPrice uint64
	L1FeeWei          uint64
}

func (r AnchorReceipt) CostWei() (uint64, error) {
	high, execution := bits.Mul64(r.GasUsed, r.EffectiveGasPrice)

	total, carry := bits.Add64(execution, r.L1FeeWei, 0)
	if high != 0 || carry != 0 {
		return 0, Errf(
			KindOverflow,
			"notary.AnchorReceipt.CostWei",
			nil,
			"transaction %s cost more wei than 64 bits hold",
			r.Hash.Hex(),
		)
	}

	return total, nil
}
