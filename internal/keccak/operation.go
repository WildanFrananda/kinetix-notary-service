package keccak

import (
	"encoding/binary"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

const (
	EntryPointName    = "ERC4337"
	EntryPointVersion = "1"

	operationType = "PackedUserOperation(address sender,uint256 nonce,bytes initCode,bytes callData," +
		"bytes32 accountGasLimits,uint256 preVerificationGas,bytes32 gasFees,bytes paymasterAndData)"
)

var _ notary.OperationScheme = Scheme{}

var entryPointDomain = domain{name: EntryPointName, version: EntryPointVersion}

func (Scheme) OperationHash(entryPoint notary.EntryPointIdentity, op notary.UserOperation) notary.Digest {
	initCode := digest(nil)
	if !op.Delegate.IsZero() {
		initCode = digest(op.Delegate[:])
	}

	structure := words(
		digest([]byte(operationType)),
		addressWord(op.Sender),
		uintWord(op.Nonce),
		initCode,
		digest(op.CallData),
		pairWord(op.VerificationGasLimit, op.CallGasLimit),
		uintWord(op.PreVerificationGas),
		pairWord(op.MaxPriorityFeePerGas, op.MaxFeePerGas),
		digest(op.PaymasterAndData()),
	)

	return typedData(entryPointDomain, entryPoint.ChainID, entryPoint.Address, structure)
}

func pairWord(high, low uint64) notary.Digest {
	var word notary.Digest
	binary.BigEndian.PutUint64(word[8:16], high)
	binary.BigEndian.PutUint64(word[24:], low)

	return word
}
