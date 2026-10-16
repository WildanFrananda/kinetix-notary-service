package notary

import "math/big"

type PaymasterIdentity struct {
	ChainID     uint64
	Address     Address
	EntryPoint  Address
	Account     Address
	Court       Address
	Withholding Address
	Signer      Address
	Ceiling     *big.Int
}
