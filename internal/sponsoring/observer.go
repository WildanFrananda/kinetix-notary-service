package sponsoring

import "math/big"

type Observer interface {
	PaymasterDeposit(wei *big.Int)
	PaymasterStake(wei *big.Int)
	PaymasterCurrent(current bool)
}
