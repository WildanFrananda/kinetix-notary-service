package notary

import "math/big"

type PaymasterFunds struct {
	Deposit *big.Int
	Stake   *big.Int
	Staked  bool
}

func (f PaymasterFunds) LockedStake() *big.Int {
	if !f.Staked || f.Stake == nil {
		return new(big.Int)
	}

	return f.Stake
}
