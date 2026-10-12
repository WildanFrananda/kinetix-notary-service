package ethereum

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type rawReceipt struct {
	Status            *hexutil.Uint64 `json:"status"`
	BlockNumber       *hexutil.Big    `json:"blockNumber"`
	GasUsed           *hexutil.Big    `json:"gasUsed"`
	EffectiveGasPrice *hexutil.Big    `json:"effectiveGasPrice"`
	L1Fee             *hexutil.Big    `json:"l1Fee"`
}

func (r *rawReceipt) decode(hash notary.Digest) (notary.AnchorReceipt, error) {
	if r.Status == nil || r.BlockNumber == nil || r.GasUsed == nil || r.EffectiveGasPrice == nil {
		return notary.AnchorReceipt{}, errors.New("the receipt lacks status, block, gas used or price")
	}

	for name, value := range map[string]*hexutil.Big{"block": r.BlockNumber, "gas used": r.GasUsed, "price": r.EffectiveGasPrice} {
		if !value.ToInt().IsUint64() {
			return notary.AnchorReceipt{}, fmt.Errorf("the receipt's %s does not fit 64 bits", name)
		}
	}

	var l1Fee uint64

	if r.L1Fee != nil {
		if !r.L1Fee.ToInt().IsUint64() {
			return notary.AnchorReceipt{}, errors.New("the receipt's L1 fee does not fit 64 bits")
		}

		l1Fee = r.L1Fee.ToInt().Uint64()
	}

	return notary.AnchorReceipt{
		Hash:              hash,
		Succeeded:         uint64(*r.Status) == 1,
		BlockNumber:       r.BlockNumber.ToInt().Uint64(),
		GasUsed:           r.GasUsed.ToInt().Uint64(),
		EffectiveGasPrice: r.EffectiveGasPrice.ToInt().Uint64(),
		L1FeeWei:          l1Fee,
	}, nil
}
