package notary

import (
	"encoding/binary"
	"math/big"
)

type UserOperation struct {
	Sender                        Address
	Nonce                         uint64
	Delegate                      Address
	CallData                      []byte
	CallGasLimit                  uint64
	VerificationGasLimit          uint64
	PreVerificationGas            uint64
	MaxFeePerGas                  uint64
	MaxPriorityFeePerGas          uint64
	Paymaster                     Address
	PaymasterVerificationGasLimit uint64
	PaymasterPostOpGasLimit       uint64
	PaymasterData                 []byte
	Signature                     []byte
}

func (op UserOperation) MaxCost() *big.Int {
	gas := new(big.Int)
	for _, limit := range []uint64{
		op.CallGasLimit,
		op.VerificationGasLimit,
		op.PreVerificationGas,
		op.PaymasterVerificationGasLimit,
		op.PaymasterPostOpGasLimit,
	} {
		gas.Add(gas, new(big.Int).SetUint64(limit))
	}

	return gas.Mul(gas, new(big.Int).SetUint64(op.MaxFeePerGas))
}

func (op UserOperation) PaymasterAndData() []byte {
	if op.Paymaster.IsZero() {
		return nil
	}

	const half = 16

	packed := make([]byte, 0, len(op.Paymaster)+2*half+len(op.PaymasterData))
	packed = append(packed, op.Paymaster[:]...)
	packed = append(packed, make([]byte, half-8)...)
	packed = binary.BigEndian.AppendUint64(packed, op.PaymasterVerificationGasLimit)
	packed = append(packed, make([]byte, half-8)...)
	packed = binary.BigEndian.AppendUint64(packed, op.PaymasterPostOpGasLimit)

	return append(packed, op.PaymasterData...)
}
