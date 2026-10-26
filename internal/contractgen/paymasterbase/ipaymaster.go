// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package paymasterbase

import (
	"bytes"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = bytes.Equal
	_ = errors.New
	_ = big.NewInt
	_ = common.Big1
	_ = types.BloomLookup
	_ = abi.ConvertType
)

// PackedUserOperation is an auto generated low-level Go binding around an user-defined struct.
type PackedUserOperation struct {
	Sender             common.Address
	Nonce              *big.Int
	InitCode           []byte
	CallData           []byte
	AccountGasLimits   [32]byte
	PreVerificationGas *big.Int
	GasFees            [32]byte
	PaymasterAndData   []byte
	Signature          []byte
}

// IPaymasterMetaData contains all meta data concerning the IPaymaster contract.
var IPaymasterMetaData = bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"postOp\",\"inputs\":[{\"name\":\"mode\",\"type\":\"uint8\",\"internalType\":\"enumIPaymaster.PostOpMode\"},{\"name\":\"context\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"actualGasCost\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actualUserOpFeePerGas\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"validatePaymasterUserOp\",\"inputs\":[{\"name\":\"userOp\",\"type\":\"tuple\",\"internalType\":\"structPackedUserOperation\",\"components\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"nonce\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"initCode\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"callData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"accountGasLimits\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"preVerificationGas\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"gasFees\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"paymasterAndData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"userOpHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"maxCost\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"context\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"validationData\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"}]",
	ID:  "IPaymaster",
}

// IPaymaster is an auto generated Go binding around an Ethereum contract.
type IPaymaster struct {
	abi abi.ABI
}

// GetABI returns the ABI associated with this contract binding.
func (c *IPaymaster) GetABI() abi.ABI {
	return c.abi
}

// NewIPaymaster creates a new instance of IPaymaster.
func NewIPaymaster() *IPaymaster {
	parsed, err := IPaymasterMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &IPaymaster{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *IPaymaster) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackPostOp is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7c627b21.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function postOp(uint8 mode, bytes context, uint256 actualGasCost, uint256 actualUserOpFeePerGas) returns()
func (iPaymaster *IPaymaster) PackPostOp(mode uint8, context []byte, actualGasCost *big.Int, actualUserOpFeePerGas *big.Int) []byte {
	enc, err := iPaymaster.abi.Pack("postOp", mode, context, actualGasCost, actualUserOpFeePerGas)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPostOp is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7c627b21.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function postOp(uint8 mode, bytes context, uint256 actualGasCost, uint256 actualUserOpFeePerGas) returns()
func (iPaymaster *IPaymaster) TryPackPostOp(mode uint8, context []byte, actualGasCost *big.Int, actualUserOpFeePerGas *big.Int) ([]byte, error) {
	return iPaymaster.abi.Pack("postOp", mode, context, actualGasCost, actualUserOpFeePerGas)
}

// PackValidatePaymasterUserOp is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x52b7512c.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function validatePaymasterUserOp((address,uint256,bytes,bytes,bytes32,uint256,bytes32,bytes,bytes) userOp, bytes32 userOpHash, uint256 maxCost) returns(bytes context, uint256 validationData)
func (iPaymaster *IPaymaster) PackValidatePaymasterUserOp(userOp PackedUserOperation, userOpHash [32]byte, maxCost *big.Int) []byte {
	enc, err := iPaymaster.abi.Pack("validatePaymasterUserOp", userOp, userOpHash, maxCost)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackValidatePaymasterUserOp is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x52b7512c.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function validatePaymasterUserOp((address,uint256,bytes,bytes,bytes32,uint256,bytes32,bytes,bytes) userOp, bytes32 userOpHash, uint256 maxCost) returns(bytes context, uint256 validationData)
func (iPaymaster *IPaymaster) TryPackValidatePaymasterUserOp(userOp PackedUserOperation, userOpHash [32]byte, maxCost *big.Int) ([]byte, error) {
	return iPaymaster.abi.Pack("validatePaymasterUserOp", userOp, userOpHash, maxCost)
}

// ValidatePaymasterUserOpOutput serves as a container for the return parameters of contract
// method ValidatePaymasterUserOp.
type ValidatePaymasterUserOpOutput struct {
	Context        []byte
	ValidationData *big.Int
}

// UnpackValidatePaymasterUserOp is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x52b7512c.
//
// Solidity: function validatePaymasterUserOp((address,uint256,bytes,bytes,bytes32,uint256,bytes32,bytes,bytes) userOp, bytes32 userOpHash, uint256 maxCost) returns(bytes context, uint256 validationData)
func (iPaymaster *IPaymaster) UnpackValidatePaymasterUserOp(data []byte) (ValidatePaymasterUserOpOutput, error) {
	out, err := iPaymaster.abi.Unpack("validatePaymasterUserOp", data)
	outstruct := new(ValidatePaymasterUserOpOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Context = *abi.ConvertType(out[0], new([]byte)).(*[]byte)
	outstruct.ValidationData = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	return *outstruct, nil
}
