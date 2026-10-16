package notary

import "math/big"

type OperationReceipt struct {
	Hash            Digest
	Success         bool
	TransactionHash Digest
	ActualGasCost   *big.Int
	Reason          string
}
