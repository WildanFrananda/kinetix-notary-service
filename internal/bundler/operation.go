package bundler

import (
	"encoding/hex"
	"math/big"
	"strconv"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

const delegationMarker = "0x7702"

type operation struct {
	Sender                        string         `json:"sender"`
	Nonce                         string         `json:"nonce"`
	Factory                       *string        `json:"factory,omitempty"`
	FactoryData                   *string        `json:"factoryData,omitempty"`
	CallData                      string         `json:"callData"`
	CallGasLimit                  string         `json:"callGasLimit"`
	VerificationGasLimit          string         `json:"verificationGasLimit"`
	PreVerificationGas            string         `json:"preVerificationGas"`
	MaxFeePerGas                  string         `json:"maxFeePerGas"`
	MaxPriorityFeePerGas          string         `json:"maxPriorityFeePerGas"`
	Paymaster                     string         `json:"paymaster"`
	PaymasterVerificationGasLimit string         `json:"paymasterVerificationGasLimit"`
	PaymasterPostOpGasLimit       string         `json:"paymasterPostOpGasLimit"`
	PaymasterData                 string         `json:"paymasterData"`
	Signature                     string         `json:"signature"`
	Authorization                 *authorization `json:"eip7702Auth,omitempty"`
}

func operationOf(op notary.UserOperation, delegation *notary.Delegation) operation {
	out := operation{
		Sender:                        op.Sender.Hex(),
		Nonce:                         hexQuantity(op.Nonce),
		CallData:                      hexBytes(op.CallData),
		CallGasLimit:                  hexQuantity(op.CallGasLimit),
		VerificationGasLimit:          hexQuantity(op.VerificationGasLimit),
		PreVerificationGas:            hexQuantity(op.PreVerificationGas),
		MaxFeePerGas:                  hexQuantity(op.MaxFeePerGas),
		MaxPriorityFeePerGas:          hexQuantity(op.MaxPriorityFeePerGas),
		Paymaster:                     op.Paymaster.Hex(),
		PaymasterVerificationGasLimit: hexQuantity(op.PaymasterVerificationGasLimit),
		PaymasterPostOpGasLimit:       hexQuantity(op.PaymasterPostOpGasLimit),
		PaymasterData:                 hexBytes(op.PaymasterData),
		Signature:                     hexBytes(op.Signature),
	}

	if !op.Delegate.IsZero() {
		marker, empty := delegationMarker, "0x"
		out.Factory = &marker
		out.FactoryData = &empty
	}

	if delegation != nil {
		out.Authorization = &authorization{
			ChainID: hexQuantity(delegation.ChainID),
			Address: delegation.Delegate.Hex(),
			Nonce:   hexQuantity(delegation.Nonce),
			YParity: hexQuantity(uint64(delegation.YParity)),
			R:       wordQuantity(delegation.R),
			S:       wordQuantity(delegation.S),
		}
	}

	return out
}

func hexQuantity(value uint64) string {
	return "0x" + strconv.FormatUint(value, 16)
}

func wordQuantity(word notary.Digest) string {
	return "0x" + new(big.Int).SetBytes(word[:]).Text(16)
}

func hexBytes(raw []byte) string {
	return "0x" + hex.EncodeToString(raw)
}
