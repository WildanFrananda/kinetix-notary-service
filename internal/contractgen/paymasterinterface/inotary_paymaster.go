// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package paymasterinterface

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

// IStakeManagerDepositInfo is an auto generated low-level Go binding around an user-defined struct.
type IStakeManagerDepositInfo struct {
	Deposit         *big.Int
	Staked          bool
	Stake           *big.Int
	UnstakeDelaySec uint32
	WithdrawTime    *big.Int
}

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

// INotaryPaymasterMetaData contains all meta data concerning the INotaryPaymaster contract.
var INotaryPaymasterMetaData = bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"account\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"addStake\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"court\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIStatementCourt\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"entryPoint\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEntryPoint\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDeposit\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDepositInfo\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIStakeManager.DepositInfo\",\"components\":[{\"name\":\"deposit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"staked\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"stake\",\"type\":\"uint112\",\"internalType\":\"uint112\"},{\"name\":\"unstakeDelaySec\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"withdrawTime\",\"type\":\"uint48\",\"internalType\":\"uint48\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"maxCostPerOperation\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"postOp\",\"inputs\":[{\"name\":\"mode\",\"type\":\"uint8\",\"internalType\":\"enumIPaymaster.PostOpMode\"},{\"name\":\"context\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"actualGasCost\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actualUserOpFeePerGas\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"sponsorship\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"callData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"paymasterData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"validAfter\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"validUntil\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"statementSigner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unlockStake\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unstakeDelay\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"validatePaymasterUserOp\",\"inputs\":[{\"name\":\"userOp\",\"type\":\"tuple\",\"internalType\":\"structPackedUserOperation\",\"components\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"nonce\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"initCode\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"callData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"accountGasLimits\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"preVerificationGas\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"gasFees\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"paymasterAndData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"userOpHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"maxCost\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"context\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"validationData\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawBalance\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawStake\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"addresspayable\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawTo\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"addresspayable\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withholding\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIWithholdingCourt\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"Withdrawn\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AlreadyProven\",\"inputs\":[{\"name\":\"evidence\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"AlreadyRequested\",\"inputs\":[{\"name\":\"registration\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"epoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"AnotherSignersEpoch\",\"inputs\":[{\"name\":\"epoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"signer\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"BalanceNotSent\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"BeforeRegistration\",\"inputs\":[{\"name\":\"epoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"fromEpoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"CallNotSponsored\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"selector\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}]},{\"type\":\"error\",\"name\":\"CostTooHigh\",\"inputs\":[{\"name\":\"maxCost\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"ceiling\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"DifferentEpochs\",\"inputs\":[{\"name\":\"first\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"second\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"DomainMismatch\",\"inputs\":[{\"name\":\"court\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"EpochNotAnchored\",\"inputs\":[{\"name\":\"epoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"NoConflict\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NoContract\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"NoOpenRequest\",\"inputs\":[{\"name\":\"registration\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"epoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"NoPostOp\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NoStatementSigner\",\"inputs\":[{\"name\":\"court\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"NotAStatementCourt\",\"inputs\":[{\"name\":\"court\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"NotAWithholdingCourt\",\"inputs\":[{\"name\":\"withholding\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"NotAnExecute\",\"inputs\":[{\"name\":\"selector\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}]},{\"type\":\"error\",\"name\":\"NotDelegated\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"NotSignedByTheNotary\",\"inputs\":[{\"name\":\"recovered\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"NotTheEntryPoint\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"NotTheHolder\",\"inputs\":[{\"name\":\"holder\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ProofHolds\",\"inputs\":[{\"name\":\"epoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"RequestStillOpen\",\"inputs\":[{\"name\":\"registration\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"epoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"SiblingsMismatch\",\"inputs\":[{\"name\":\"signed\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"given\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"TargetNotSponsored\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"UnstakeDelayTooShort\",\"inputs\":[{\"name\":\"given\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"least\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"ValueNotZero\",\"inputs\":[{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"WithholdingOfAnotherCourt\",\"inputs\":[{\"name\":\"withholding\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"court\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"WrongLedger\",\"inputs\":[{\"name\":\"expected\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"given\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ZeroCostCeiling\",\"inputs\":[]}]",
	ID:  "INotaryPaymaster",
}

// INotaryPaymaster is an auto generated Go binding around an Ethereum contract.
type INotaryPaymaster struct {
	abi abi.ABI
}

// GetABI returns the ABI associated with this contract binding.
func (c *INotaryPaymaster) GetABI() abi.ABI {
	return c.abi
}

// NewINotaryPaymaster creates a new instance of INotaryPaymaster.
func NewINotaryPaymaster() *INotaryPaymaster {
	parsed, err := INotaryPaymasterMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &INotaryPaymaster{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *INotaryPaymaster) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackAccount is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5dab2420.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function account() view returns(address)
func (iNotaryPaymaster *INotaryPaymaster) PackAccount() []byte {
	enc, err := iNotaryPaymaster.abi.Pack("account")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAccount is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5dab2420.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function account() view returns(address)
func (iNotaryPaymaster *INotaryPaymaster) TryPackAccount() ([]byte, error) {
	return iNotaryPaymaster.abi.Pack("account")
}

// UnpackAccount is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5dab2420.
//
// Solidity: function account() view returns(address)
func (iNotaryPaymaster *INotaryPaymaster) UnpackAccount(data []byte) (common.Address, error) {
	out, err := iNotaryPaymaster.abi.Unpack("account", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackAddStake is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5a627dbc.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function addStake() payable returns()
func (iNotaryPaymaster *INotaryPaymaster) PackAddStake() []byte {
	enc, err := iNotaryPaymaster.abi.Pack("addStake")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAddStake is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5a627dbc.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function addStake() payable returns()
func (iNotaryPaymaster *INotaryPaymaster) TryPackAddStake() ([]byte, error) {
	return iNotaryPaymaster.abi.Pack("addStake")
}

// PackCourt is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x00971a30.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function court() view returns(address)
func (iNotaryPaymaster *INotaryPaymaster) PackCourt() []byte {
	enc, err := iNotaryPaymaster.abi.Pack("court")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCourt is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x00971a30.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function court() view returns(address)
func (iNotaryPaymaster *INotaryPaymaster) TryPackCourt() ([]byte, error) {
	return iNotaryPaymaster.abi.Pack("court")
}

// UnpackCourt is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x00971a30.
//
// Solidity: function court() view returns(address)
func (iNotaryPaymaster *INotaryPaymaster) UnpackCourt(data []byte) (common.Address, error) {
	out, err := iNotaryPaymaster.abi.Unpack("court", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackEntryPoint is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb0d691fe.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function entryPoint() view returns(address)
func (iNotaryPaymaster *INotaryPaymaster) PackEntryPoint() []byte {
	enc, err := iNotaryPaymaster.abi.Pack("entryPoint")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackEntryPoint is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb0d691fe.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function entryPoint() view returns(address)
func (iNotaryPaymaster *INotaryPaymaster) TryPackEntryPoint() ([]byte, error) {
	return iNotaryPaymaster.abi.Pack("entryPoint")
}

// UnpackEntryPoint is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb0d691fe.
//
// Solidity: function entryPoint() view returns(address)
func (iNotaryPaymaster *INotaryPaymaster) UnpackEntryPoint(data []byte) (common.Address, error) {
	out, err := iNotaryPaymaster.abi.Unpack("entryPoint", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackGetDeposit is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc399ec88.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getDeposit() view returns(uint256)
func (iNotaryPaymaster *INotaryPaymaster) PackGetDeposit() []byte {
	enc, err := iNotaryPaymaster.abi.Pack("getDeposit")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetDeposit is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc399ec88.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getDeposit() view returns(uint256)
func (iNotaryPaymaster *INotaryPaymaster) TryPackGetDeposit() ([]byte, error) {
	return iNotaryPaymaster.abi.Pack("getDeposit")
}

// UnpackGetDeposit is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc399ec88.
//
// Solidity: function getDeposit() view returns(uint256)
func (iNotaryPaymaster *INotaryPaymaster) UnpackGetDeposit(data []byte) (*big.Int, error) {
	out, err := iNotaryPaymaster.abi.Unpack("getDeposit", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackGetDepositInfo is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8969e0f9.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getDepositInfo() view returns((uint256,bool,uint112,uint32,uint48))
func (iNotaryPaymaster *INotaryPaymaster) PackGetDepositInfo() []byte {
	enc, err := iNotaryPaymaster.abi.Pack("getDepositInfo")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetDepositInfo is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8969e0f9.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getDepositInfo() view returns((uint256,bool,uint112,uint32,uint48))
func (iNotaryPaymaster *INotaryPaymaster) TryPackGetDepositInfo() ([]byte, error) {
	return iNotaryPaymaster.abi.Pack("getDepositInfo")
}

// UnpackGetDepositInfo is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8969e0f9.
//
// Solidity: function getDepositInfo() view returns((uint256,bool,uint112,uint32,uint48))
func (iNotaryPaymaster *INotaryPaymaster) UnpackGetDepositInfo(data []byte) (IStakeManagerDepositInfo, error) {
	out, err := iNotaryPaymaster.abi.Unpack("getDepositInfo", data)
	if err != nil {
		return *new(IStakeManagerDepositInfo), err
	}
	out0 := *abi.ConvertType(out[0], new(IStakeManagerDepositInfo)).(*IStakeManagerDepositInfo)
	return out0, nil
}

// PackMaxCostPerOperation is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xcf1a07dc.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function maxCostPerOperation() view returns(uint256)
func (iNotaryPaymaster *INotaryPaymaster) PackMaxCostPerOperation() []byte {
	enc, err := iNotaryPaymaster.abi.Pack("maxCostPerOperation")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMaxCostPerOperation is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xcf1a07dc.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function maxCostPerOperation() view returns(uint256)
func (iNotaryPaymaster *INotaryPaymaster) TryPackMaxCostPerOperation() ([]byte, error) {
	return iNotaryPaymaster.abi.Pack("maxCostPerOperation")
}

// UnpackMaxCostPerOperation is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xcf1a07dc.
//
// Solidity: function maxCostPerOperation() view returns(uint256)
func (iNotaryPaymaster *INotaryPaymaster) UnpackMaxCostPerOperation(data []byte) (*big.Int, error) {
	out, err := iNotaryPaymaster.abi.Unpack("maxCostPerOperation", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackPostOp is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7c627b21.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function postOp(uint8 mode, bytes context, uint256 actualGasCost, uint256 actualUserOpFeePerGas) returns()
func (iNotaryPaymaster *INotaryPaymaster) PackPostOp(mode uint8, context []byte, actualGasCost *big.Int, actualUserOpFeePerGas *big.Int) []byte {
	enc, err := iNotaryPaymaster.abi.Pack("postOp", mode, context, actualGasCost, actualUserOpFeePerGas)
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
func (iNotaryPaymaster *INotaryPaymaster) TryPackPostOp(mode uint8, context []byte, actualGasCost *big.Int, actualUserOpFeePerGas *big.Int) ([]byte, error) {
	return iNotaryPaymaster.abi.Pack("postOp", mode, context, actualGasCost, actualUserOpFeePerGas)
}

// PackSponsorship is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xada0bfa4.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function sponsorship(address sender, bytes callData, bytes paymasterData) view returns(uint48 validAfter, uint48 validUntil)
func (iNotaryPaymaster *INotaryPaymaster) PackSponsorship(sender common.Address, callData []byte, paymasterData []byte) []byte {
	enc, err := iNotaryPaymaster.abi.Pack("sponsorship", sender, callData, paymasterData)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSponsorship is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xada0bfa4.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function sponsorship(address sender, bytes callData, bytes paymasterData) view returns(uint48 validAfter, uint48 validUntil)
func (iNotaryPaymaster *INotaryPaymaster) TryPackSponsorship(sender common.Address, callData []byte, paymasterData []byte) ([]byte, error) {
	return iNotaryPaymaster.abi.Pack("sponsorship", sender, callData, paymasterData)
}

// SponsorshipOutput serves as a container for the return parameters of contract
// method Sponsorship.
type SponsorshipOutput struct {
	ValidAfter *big.Int
	ValidUntil *big.Int
}

// UnpackSponsorship is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xada0bfa4.
//
// Solidity: function sponsorship(address sender, bytes callData, bytes paymasterData) view returns(uint48 validAfter, uint48 validUntil)
func (iNotaryPaymaster *INotaryPaymaster) UnpackSponsorship(data []byte) (SponsorshipOutput, error) {
	out, err := iNotaryPaymaster.abi.Unpack("sponsorship", data)
	outstruct := new(SponsorshipOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.ValidAfter = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.ValidUntil = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	return *outstruct, nil
}

// PackStatementSigner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe0310972.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function statementSigner() view returns(address)
func (iNotaryPaymaster *INotaryPaymaster) PackStatementSigner() []byte {
	enc, err := iNotaryPaymaster.abi.Pack("statementSigner")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackStatementSigner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe0310972.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function statementSigner() view returns(address)
func (iNotaryPaymaster *INotaryPaymaster) TryPackStatementSigner() ([]byte, error) {
	return iNotaryPaymaster.abi.Pack("statementSigner")
}

// UnpackStatementSigner is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe0310972.
//
// Solidity: function statementSigner() view returns(address)
func (iNotaryPaymaster *INotaryPaymaster) UnpackStatementSigner(data []byte) (common.Address, error) {
	out, err := iNotaryPaymaster.abi.Unpack("statementSigner", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackUnlockStake is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbb9fe6bf.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function unlockStake() returns()
func (iNotaryPaymaster *INotaryPaymaster) PackUnlockStake() []byte {
	enc, err := iNotaryPaymaster.abi.Pack("unlockStake")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackUnlockStake is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbb9fe6bf.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function unlockStake() returns()
func (iNotaryPaymaster *INotaryPaymaster) TryPackUnlockStake() ([]byte, error) {
	return iNotaryPaymaster.abi.Pack("unlockStake")
}

// PackUnstakeDelay is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3222012f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function unstakeDelay() view returns(uint32)
func (iNotaryPaymaster *INotaryPaymaster) PackUnstakeDelay() []byte {
	enc, err := iNotaryPaymaster.abi.Pack("unstakeDelay")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackUnstakeDelay is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3222012f.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function unstakeDelay() view returns(uint32)
func (iNotaryPaymaster *INotaryPaymaster) TryPackUnstakeDelay() ([]byte, error) {
	return iNotaryPaymaster.abi.Pack("unstakeDelay")
}

// UnpackUnstakeDelay is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x3222012f.
//
// Solidity: function unstakeDelay() view returns(uint32)
func (iNotaryPaymaster *INotaryPaymaster) UnpackUnstakeDelay(data []byte) (uint32, error) {
	out, err := iNotaryPaymaster.abi.Unpack("unstakeDelay", data)
	if err != nil {
		return *new(uint32), err
	}
	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)
	return out0, nil
}

// PackValidatePaymasterUserOp is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x52b7512c.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function validatePaymasterUserOp((address,uint256,bytes,bytes,bytes32,uint256,bytes32,bytes,bytes) userOp, bytes32 userOpHash, uint256 maxCost) returns(bytes context, uint256 validationData)
func (iNotaryPaymaster *INotaryPaymaster) PackValidatePaymasterUserOp(userOp PackedUserOperation, userOpHash [32]byte, maxCost *big.Int) []byte {
	enc, err := iNotaryPaymaster.abi.Pack("validatePaymasterUserOp", userOp, userOpHash, maxCost)
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
func (iNotaryPaymaster *INotaryPaymaster) TryPackValidatePaymasterUserOp(userOp PackedUserOperation, userOpHash [32]byte, maxCost *big.Int) ([]byte, error) {
	return iNotaryPaymaster.abi.Pack("validatePaymasterUserOp", userOp, userOpHash, maxCost)
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
func (iNotaryPaymaster *INotaryPaymaster) UnpackValidatePaymasterUserOp(data []byte) (ValidatePaymasterUserOpOutput, error) {
	out, err := iNotaryPaymaster.abi.Unpack("validatePaymasterUserOp", data)
	outstruct := new(ValidatePaymasterUserOpOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Context = *abi.ConvertType(out[0], new([]byte)).(*[]byte)
	outstruct.ValidationData = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	return *outstruct, nil
}

// PackWithdrawBalance is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5fd8c710.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function withdrawBalance() returns()
func (iNotaryPaymaster *INotaryPaymaster) PackWithdrawBalance() []byte {
	enc, err := iNotaryPaymaster.abi.Pack("withdrawBalance")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackWithdrawBalance is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5fd8c710.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function withdrawBalance() returns()
func (iNotaryPaymaster *INotaryPaymaster) TryPackWithdrawBalance() ([]byte, error) {
	return iNotaryPaymaster.abi.Pack("withdrawBalance")
}

// PackWithdrawStake is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc23a5cea.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function withdrawStake(address to) returns()
func (iNotaryPaymaster *INotaryPaymaster) PackWithdrawStake(to common.Address) []byte {
	enc, err := iNotaryPaymaster.abi.Pack("withdrawStake", to)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackWithdrawStake is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc23a5cea.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function withdrawStake(address to) returns()
func (iNotaryPaymaster *INotaryPaymaster) TryPackWithdrawStake(to common.Address) ([]byte, error) {
	return iNotaryPaymaster.abi.Pack("withdrawStake", to)
}

// PackWithdrawTo is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x205c2878.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function withdrawTo(address to, uint256 amount) returns()
func (iNotaryPaymaster *INotaryPaymaster) PackWithdrawTo(to common.Address, amount *big.Int) []byte {
	enc, err := iNotaryPaymaster.abi.Pack("withdrawTo", to, amount)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackWithdrawTo is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x205c2878.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function withdrawTo(address to, uint256 amount) returns()
func (iNotaryPaymaster *INotaryPaymaster) TryPackWithdrawTo(to common.Address, amount *big.Int) ([]byte, error) {
	return iNotaryPaymaster.abi.Pack("withdrawTo", to, amount)
}

// PackWithholding is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8297d665.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function withholding() view returns(address)
func (iNotaryPaymaster *INotaryPaymaster) PackWithholding() []byte {
	enc, err := iNotaryPaymaster.abi.Pack("withholding")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackWithholding is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8297d665.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function withholding() view returns(address)
func (iNotaryPaymaster *INotaryPaymaster) TryPackWithholding() ([]byte, error) {
	return iNotaryPaymaster.abi.Pack("withholding")
}

// UnpackWithholding is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8297d665.
//
// Solidity: function withholding() view returns(address)
func (iNotaryPaymaster *INotaryPaymaster) UnpackWithholding(data []byte) (common.Address, error) {
	out, err := iNotaryPaymaster.abi.Unpack("withholding", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// INotaryPaymasterWithdrawn represents a Withdrawn event raised by the INotaryPaymaster contract.
type INotaryPaymasterWithdrawn struct {
	To     common.Address
	Amount *big.Int
	Raw    *types.Log // Blockchain specific contextual infos
}

const INotaryPaymasterWithdrawnEventName = "Withdrawn"

// ContractEventName returns the user-defined event name.
func (INotaryPaymasterWithdrawn) ContractEventName() string {
	return INotaryPaymasterWithdrawnEventName
}

// UnpackWithdrawnEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Withdrawn(address indexed to, uint256 amount)
func (iNotaryPaymaster *INotaryPaymaster) UnpackWithdrawnEvent(log *types.Log) (*INotaryPaymasterWithdrawn, error) {
	event := "Withdrawn"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != iNotaryPaymaster.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(INotaryPaymasterWithdrawn)
	if len(log.Data) > 0 {
		if err := iNotaryPaymaster.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iNotaryPaymaster.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// UnpackError attempts to decode the provided error data using user-defined
// error definitions.
func (iNotaryPaymaster *INotaryPaymaster) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], iNotaryPaymaster.abi.Errors["AlreadyProven"].ID.Bytes()[:4]) {
		return iNotaryPaymaster.UnpackAlreadyProvenError(raw[4:])
	}
	if bytes.Equal(raw[:4], iNotaryPaymaster.abi.Errors["AlreadyRequested"].ID.Bytes()[:4]) {
		return iNotaryPaymaster.UnpackAlreadyRequestedError(raw[4:])
	}
	if bytes.Equal(raw[:4], iNotaryPaymaster.abi.Errors["AnotherSignersEpoch"].ID.Bytes()[:4]) {
		return iNotaryPaymaster.UnpackAnotherSignersEpochError(raw[4:])
	}
	if bytes.Equal(raw[:4], iNotaryPaymaster.abi.Errors["BalanceNotSent"].ID.Bytes()[:4]) {
		return iNotaryPaymaster.UnpackBalanceNotSentError(raw[4:])
	}
	if bytes.Equal(raw[:4], iNotaryPaymaster.abi.Errors["BeforeRegistration"].ID.Bytes()[:4]) {
		return iNotaryPaymaster.UnpackBeforeRegistrationError(raw[4:])
	}
	if bytes.Equal(raw[:4], iNotaryPaymaster.abi.Errors["CallNotSponsored"].ID.Bytes()[:4]) {
		return iNotaryPaymaster.UnpackCallNotSponsoredError(raw[4:])
	}
	if bytes.Equal(raw[:4], iNotaryPaymaster.abi.Errors["CostTooHigh"].ID.Bytes()[:4]) {
		return iNotaryPaymaster.UnpackCostTooHighError(raw[4:])
	}
	if bytes.Equal(raw[:4], iNotaryPaymaster.abi.Errors["DifferentEpochs"].ID.Bytes()[:4]) {
		return iNotaryPaymaster.UnpackDifferentEpochsError(raw[4:])
	}
	if bytes.Equal(raw[:4], iNotaryPaymaster.abi.Errors["DomainMismatch"].ID.Bytes()[:4]) {
		return iNotaryPaymaster.UnpackDomainMismatchError(raw[4:])
	}
	if bytes.Equal(raw[:4], iNotaryPaymaster.abi.Errors["EpochNotAnchored"].ID.Bytes()[:4]) {
		return iNotaryPaymaster.UnpackEpochNotAnchoredError(raw[4:])
	}
	if bytes.Equal(raw[:4], iNotaryPaymaster.abi.Errors["NoConflict"].ID.Bytes()[:4]) {
		return iNotaryPaymaster.UnpackNoConflictError(raw[4:])
	}
	if bytes.Equal(raw[:4], iNotaryPaymaster.abi.Errors["NoContract"].ID.Bytes()[:4]) {
		return iNotaryPaymaster.UnpackNoContractError(raw[4:])
	}
	if bytes.Equal(raw[:4], iNotaryPaymaster.abi.Errors["NoOpenRequest"].ID.Bytes()[:4]) {
		return iNotaryPaymaster.UnpackNoOpenRequestError(raw[4:])
	}
	if bytes.Equal(raw[:4], iNotaryPaymaster.abi.Errors["NoPostOp"].ID.Bytes()[:4]) {
		return iNotaryPaymaster.UnpackNoPostOpError(raw[4:])
	}
	if bytes.Equal(raw[:4], iNotaryPaymaster.abi.Errors["NoStatementSigner"].ID.Bytes()[:4]) {
		return iNotaryPaymaster.UnpackNoStatementSignerError(raw[4:])
	}
	if bytes.Equal(raw[:4], iNotaryPaymaster.abi.Errors["NotAStatementCourt"].ID.Bytes()[:4]) {
		return iNotaryPaymaster.UnpackNotAStatementCourtError(raw[4:])
	}
	if bytes.Equal(raw[:4], iNotaryPaymaster.abi.Errors["NotAWithholdingCourt"].ID.Bytes()[:4]) {
		return iNotaryPaymaster.UnpackNotAWithholdingCourtError(raw[4:])
	}
	if bytes.Equal(raw[:4], iNotaryPaymaster.abi.Errors["NotAnExecute"].ID.Bytes()[:4]) {
		return iNotaryPaymaster.UnpackNotAnExecuteError(raw[4:])
	}
	if bytes.Equal(raw[:4], iNotaryPaymaster.abi.Errors["NotDelegated"].ID.Bytes()[:4]) {
		return iNotaryPaymaster.UnpackNotDelegatedError(raw[4:])
	}
	if bytes.Equal(raw[:4], iNotaryPaymaster.abi.Errors["NotSignedByTheNotary"].ID.Bytes()[:4]) {
		return iNotaryPaymaster.UnpackNotSignedByTheNotaryError(raw[4:])
	}
	if bytes.Equal(raw[:4], iNotaryPaymaster.abi.Errors["NotTheEntryPoint"].ID.Bytes()[:4]) {
		return iNotaryPaymaster.UnpackNotTheEntryPointError(raw[4:])
	}
	if bytes.Equal(raw[:4], iNotaryPaymaster.abi.Errors["NotTheHolder"].ID.Bytes()[:4]) {
		return iNotaryPaymaster.UnpackNotTheHolderError(raw[4:])
	}
	if bytes.Equal(raw[:4], iNotaryPaymaster.abi.Errors["ProofHolds"].ID.Bytes()[:4]) {
		return iNotaryPaymaster.UnpackProofHoldsError(raw[4:])
	}
	if bytes.Equal(raw[:4], iNotaryPaymaster.abi.Errors["RequestStillOpen"].ID.Bytes()[:4]) {
		return iNotaryPaymaster.UnpackRequestStillOpenError(raw[4:])
	}
	if bytes.Equal(raw[:4], iNotaryPaymaster.abi.Errors["SiblingsMismatch"].ID.Bytes()[:4]) {
		return iNotaryPaymaster.UnpackSiblingsMismatchError(raw[4:])
	}
	if bytes.Equal(raw[:4], iNotaryPaymaster.abi.Errors["TargetNotSponsored"].ID.Bytes()[:4]) {
		return iNotaryPaymaster.UnpackTargetNotSponsoredError(raw[4:])
	}
	if bytes.Equal(raw[:4], iNotaryPaymaster.abi.Errors["UnstakeDelayTooShort"].ID.Bytes()[:4]) {
		return iNotaryPaymaster.UnpackUnstakeDelayTooShortError(raw[4:])
	}
	if bytes.Equal(raw[:4], iNotaryPaymaster.abi.Errors["ValueNotZero"].ID.Bytes()[:4]) {
		return iNotaryPaymaster.UnpackValueNotZeroError(raw[4:])
	}
	if bytes.Equal(raw[:4], iNotaryPaymaster.abi.Errors["WithholdingOfAnotherCourt"].ID.Bytes()[:4]) {
		return iNotaryPaymaster.UnpackWithholdingOfAnotherCourtError(raw[4:])
	}
	if bytes.Equal(raw[:4], iNotaryPaymaster.abi.Errors["WrongLedger"].ID.Bytes()[:4]) {
		return iNotaryPaymaster.UnpackWrongLedgerError(raw[4:])
	}
	if bytes.Equal(raw[:4], iNotaryPaymaster.abi.Errors["ZeroCostCeiling"].ID.Bytes()[:4]) {
		return iNotaryPaymaster.UnpackZeroCostCeilingError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// INotaryPaymasterAlreadyProven represents a AlreadyProven error raised by the INotaryPaymaster contract.
type INotaryPaymasterAlreadyProven struct {
	Evidence [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AlreadyProven(bytes32 evidence)
func INotaryPaymasterAlreadyProvenErrorID() common.Hash {
	return common.HexToHash("0xea033684b2eb17b3fc70d5009bef3efe803909cba52e3380dd0e081697f487ad")
}

// UnpackAlreadyProvenError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AlreadyProven(bytes32 evidence)
func (iNotaryPaymaster *INotaryPaymaster) UnpackAlreadyProvenError(raw []byte) (*INotaryPaymasterAlreadyProven, error) {
	out := new(INotaryPaymasterAlreadyProven)
	if err := iNotaryPaymaster.abi.UnpackIntoInterface(out, "AlreadyProven", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// INotaryPaymasterAlreadyRequested represents a AlreadyRequested error raised by the INotaryPaymaster contract.
type INotaryPaymasterAlreadyRequested struct {
	Registration [32]byte
	Epoch        uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AlreadyRequested(bytes32 registration, uint64 epoch)
func INotaryPaymasterAlreadyRequestedErrorID() common.Hash {
	return common.HexToHash("0x46cdd0fb7d2db3112c73591b2011dc63aca93f893280dca515f60dd290122c34")
}

// UnpackAlreadyRequestedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AlreadyRequested(bytes32 registration, uint64 epoch)
func (iNotaryPaymaster *INotaryPaymaster) UnpackAlreadyRequestedError(raw []byte) (*INotaryPaymasterAlreadyRequested, error) {
	out := new(INotaryPaymasterAlreadyRequested)
	if err := iNotaryPaymaster.abi.UnpackIntoInterface(out, "AlreadyRequested", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// INotaryPaymasterAnotherSignersEpoch represents a AnotherSignersEpoch error raised by the INotaryPaymaster contract.
type INotaryPaymasterAnotherSignersEpoch struct {
	Epoch  uint64
	Signer common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AnotherSignersEpoch(uint64 epoch, address signer)
func INotaryPaymasterAnotherSignersEpochErrorID() common.Hash {
	return common.HexToHash("0xa637dc4e71e580f5c964f91860e49bcbd54230efbcca39beb3ea9e2ee4683368")
}

// UnpackAnotherSignersEpochError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AnotherSignersEpoch(uint64 epoch, address signer)
func (iNotaryPaymaster *INotaryPaymaster) UnpackAnotherSignersEpochError(raw []byte) (*INotaryPaymasterAnotherSignersEpoch, error) {
	out := new(INotaryPaymasterAnotherSignersEpoch)
	if err := iNotaryPaymaster.abi.UnpackIntoInterface(out, "AnotherSignersEpoch", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// INotaryPaymasterBalanceNotSent represents a BalanceNotSent error raised by the INotaryPaymaster contract.
type INotaryPaymasterBalanceNotSent struct {
	To     common.Address
	Amount *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error BalanceNotSent(address to, uint256 amount)
func INotaryPaymasterBalanceNotSentErrorID() common.Hash {
	return common.HexToHash("0xe6839b74abd0fa57c33bde4f7ccf16e495bac1ed489eb343ebdccda90606562e")
}

// UnpackBalanceNotSentError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error BalanceNotSent(address to, uint256 amount)
func (iNotaryPaymaster *INotaryPaymaster) UnpackBalanceNotSentError(raw []byte) (*INotaryPaymasterBalanceNotSent, error) {
	out := new(INotaryPaymasterBalanceNotSent)
	if err := iNotaryPaymaster.abi.UnpackIntoInterface(out, "BalanceNotSent", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// INotaryPaymasterBeforeRegistration represents a BeforeRegistration error raised by the INotaryPaymaster contract.
type INotaryPaymasterBeforeRegistration struct {
	Epoch     uint64
	FromEpoch uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error BeforeRegistration(uint64 epoch, uint64 fromEpoch)
func INotaryPaymasterBeforeRegistrationErrorID() common.Hash {
	return common.HexToHash("0xd30d092f5afe15b07c027869b0e62d560be8c5e7d28f175dd1f43a5bbecfb25b")
}

// UnpackBeforeRegistrationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error BeforeRegistration(uint64 epoch, uint64 fromEpoch)
func (iNotaryPaymaster *INotaryPaymaster) UnpackBeforeRegistrationError(raw []byte) (*INotaryPaymasterBeforeRegistration, error) {
	out := new(INotaryPaymasterBeforeRegistration)
	if err := iNotaryPaymaster.abi.UnpackIntoInterface(out, "BeforeRegistration", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// INotaryPaymasterCallNotSponsored represents a CallNotSponsored error raised by the INotaryPaymaster contract.
type INotaryPaymasterCallNotSponsored struct {
	Target   common.Address
	Selector [4]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error CallNotSponsored(address target, bytes4 selector)
func INotaryPaymasterCallNotSponsoredErrorID() common.Hash {
	return common.HexToHash("0x5b2c5475be2f0d59defc3a1ea0edc26cabdc2eb1374a34314b9fda00f66ac87e")
}

// UnpackCallNotSponsoredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error CallNotSponsored(address target, bytes4 selector)
func (iNotaryPaymaster *INotaryPaymaster) UnpackCallNotSponsoredError(raw []byte) (*INotaryPaymasterCallNotSponsored, error) {
	out := new(INotaryPaymasterCallNotSponsored)
	if err := iNotaryPaymaster.abi.UnpackIntoInterface(out, "CallNotSponsored", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// INotaryPaymasterCostTooHigh represents a CostTooHigh error raised by the INotaryPaymaster contract.
type INotaryPaymasterCostTooHigh struct {
	MaxCost *big.Int
	Ceiling *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error CostTooHigh(uint256 maxCost, uint256 ceiling)
func INotaryPaymasterCostTooHighErrorID() common.Hash {
	return common.HexToHash("0x0de362474d48e3eae8acb9881b2acab260a363515daa01c9e6c9543d4e638aec")
}

// UnpackCostTooHighError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error CostTooHigh(uint256 maxCost, uint256 ceiling)
func (iNotaryPaymaster *INotaryPaymaster) UnpackCostTooHighError(raw []byte) (*INotaryPaymasterCostTooHigh, error) {
	out := new(INotaryPaymasterCostTooHigh)
	if err := iNotaryPaymaster.abi.UnpackIntoInterface(out, "CostTooHigh", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// INotaryPaymasterDifferentEpochs represents a DifferentEpochs error raised by the INotaryPaymaster contract.
type INotaryPaymasterDifferentEpochs struct {
	First  uint64
	Second uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error DifferentEpochs(uint64 first, uint64 second)
func INotaryPaymasterDifferentEpochsErrorID() common.Hash {
	return common.HexToHash("0xbfb240a28f6d744e2666b046f0954afb0cb85dfbc6d6bd027717efc63eed50f4")
}

// UnpackDifferentEpochsError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error DifferentEpochs(uint64 first, uint64 second)
func (iNotaryPaymaster *INotaryPaymaster) UnpackDifferentEpochsError(raw []byte) (*INotaryPaymasterDifferentEpochs, error) {
	out := new(INotaryPaymasterDifferentEpochs)
	if err := iNotaryPaymaster.abi.UnpackIntoInterface(out, "DifferentEpochs", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// INotaryPaymasterDomainMismatch represents a DomainMismatch error raised by the INotaryPaymaster contract.
type INotaryPaymasterDomainMismatch struct {
	Court common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error DomainMismatch(address court)
func INotaryPaymasterDomainMismatchErrorID() common.Hash {
	return common.HexToHash("0xd3069578856ab5fa6ed66a0a701b1fa1d0722f84ec2d8b2adcb4a9643b151fc4")
}

// UnpackDomainMismatchError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error DomainMismatch(address court)
func (iNotaryPaymaster *INotaryPaymaster) UnpackDomainMismatchError(raw []byte) (*INotaryPaymasterDomainMismatch, error) {
	out := new(INotaryPaymasterDomainMismatch)
	if err := iNotaryPaymaster.abi.UnpackIntoInterface(out, "DomainMismatch", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// INotaryPaymasterEpochNotAnchored represents a EpochNotAnchored error raised by the INotaryPaymaster contract.
type INotaryPaymasterEpochNotAnchored struct {
	Epoch uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error EpochNotAnchored(uint64 epoch)
func INotaryPaymasterEpochNotAnchoredErrorID() common.Hash {
	return common.HexToHash("0xfdd66f09e3e0335148d05ffa812e23b41ab1c93ea73617ca6ddd26710f05879d")
}

// UnpackEpochNotAnchoredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error EpochNotAnchored(uint64 epoch)
func (iNotaryPaymaster *INotaryPaymaster) UnpackEpochNotAnchoredError(raw []byte) (*INotaryPaymasterEpochNotAnchored, error) {
	out := new(INotaryPaymasterEpochNotAnchored)
	if err := iNotaryPaymaster.abi.UnpackIntoInterface(out, "EpochNotAnchored", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// INotaryPaymasterNoConflict represents a NoConflict error raised by the INotaryPaymaster contract.
type INotaryPaymasterNoConflict struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NoConflict()
func INotaryPaymasterNoConflictErrorID() common.Hash {
	return common.HexToHash("0x7412124e23106f671f0c15e99815ae376e602111b93496444e63ed6071334091")
}

// UnpackNoConflictError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NoConflict()
func (iNotaryPaymaster *INotaryPaymaster) UnpackNoConflictError(raw []byte) (*INotaryPaymasterNoConflict, error) {
	out := new(INotaryPaymasterNoConflict)
	if err := iNotaryPaymaster.abi.UnpackIntoInterface(out, "NoConflict", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// INotaryPaymasterNoContract represents a NoContract error raised by the INotaryPaymaster contract.
type INotaryPaymasterNoContract struct {
	Target common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NoContract(address target)
func INotaryPaymasterNoContractErrorID() common.Hash {
	return common.HexToHash("0x5f15d672b6235f8600ffc72925d8d2f9dcea14be067296327891153847185a5c")
}

// UnpackNoContractError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NoContract(address target)
func (iNotaryPaymaster *INotaryPaymaster) UnpackNoContractError(raw []byte) (*INotaryPaymasterNoContract, error) {
	out := new(INotaryPaymasterNoContract)
	if err := iNotaryPaymaster.abi.UnpackIntoInterface(out, "NoContract", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// INotaryPaymasterNoOpenRequest represents a NoOpenRequest error raised by the INotaryPaymaster contract.
type INotaryPaymasterNoOpenRequest struct {
	Registration [32]byte
	Epoch        uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NoOpenRequest(bytes32 registration, uint64 epoch)
func INotaryPaymasterNoOpenRequestErrorID() common.Hash {
	return common.HexToHash("0x12ef5f2cac24af69d330f4a7ae9d2d98b14d43bae36fbacf6df51b7f36fd3c8d")
}

// UnpackNoOpenRequestError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NoOpenRequest(bytes32 registration, uint64 epoch)
func (iNotaryPaymaster *INotaryPaymaster) UnpackNoOpenRequestError(raw []byte) (*INotaryPaymasterNoOpenRequest, error) {
	out := new(INotaryPaymasterNoOpenRequest)
	if err := iNotaryPaymaster.abi.UnpackIntoInterface(out, "NoOpenRequest", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// INotaryPaymasterNoPostOp represents a NoPostOp error raised by the INotaryPaymaster contract.
type INotaryPaymasterNoPostOp struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NoPostOp()
func INotaryPaymasterNoPostOpErrorID() common.Hash {
	return common.HexToHash("0xd7fa401919519aabd49e863de4c3909a1380ea163237e6e30effe7b62b4c367a")
}

// UnpackNoPostOpError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NoPostOp()
func (iNotaryPaymaster *INotaryPaymaster) UnpackNoPostOpError(raw []byte) (*INotaryPaymasterNoPostOp, error) {
	out := new(INotaryPaymasterNoPostOp)
	if err := iNotaryPaymaster.abi.UnpackIntoInterface(out, "NoPostOp", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// INotaryPaymasterNoStatementSigner represents a NoStatementSigner error raised by the INotaryPaymaster contract.
type INotaryPaymasterNoStatementSigner struct {
	Court common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NoStatementSigner(address court)
func INotaryPaymasterNoStatementSignerErrorID() common.Hash {
	return common.HexToHash("0xfdd8641f39c1d15b423a0faf7e17300ed8dfd692c63500bba8d77d62b1a05212")
}

// UnpackNoStatementSignerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NoStatementSigner(address court)
func (iNotaryPaymaster *INotaryPaymaster) UnpackNoStatementSignerError(raw []byte) (*INotaryPaymasterNoStatementSigner, error) {
	out := new(INotaryPaymasterNoStatementSigner)
	if err := iNotaryPaymaster.abi.UnpackIntoInterface(out, "NoStatementSigner", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// INotaryPaymasterNotAStatementCourt represents a NotAStatementCourt error raised by the INotaryPaymaster contract.
type INotaryPaymasterNotAStatementCourt struct {
	Court common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotAStatementCourt(address court)
func INotaryPaymasterNotAStatementCourtErrorID() common.Hash {
	return common.HexToHash("0x8e35ed64fa846e83af245ae076efcecd211c173d4d90de52d3daab09ef787aff")
}

// UnpackNotAStatementCourtError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotAStatementCourt(address court)
func (iNotaryPaymaster *INotaryPaymaster) UnpackNotAStatementCourtError(raw []byte) (*INotaryPaymasterNotAStatementCourt, error) {
	out := new(INotaryPaymasterNotAStatementCourt)
	if err := iNotaryPaymaster.abi.UnpackIntoInterface(out, "NotAStatementCourt", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// INotaryPaymasterNotAWithholdingCourt represents a NotAWithholdingCourt error raised by the INotaryPaymaster contract.
type INotaryPaymasterNotAWithholdingCourt struct {
	Withholding common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotAWithholdingCourt(address withholding)
func INotaryPaymasterNotAWithholdingCourtErrorID() common.Hash {
	return common.HexToHash("0x6e9637d8ce542711d9e0648e507cdb2eaecb6b823047364c9c5ed9e719d00bac")
}

// UnpackNotAWithholdingCourtError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotAWithholdingCourt(address withholding)
func (iNotaryPaymaster *INotaryPaymaster) UnpackNotAWithholdingCourtError(raw []byte) (*INotaryPaymasterNotAWithholdingCourt, error) {
	out := new(INotaryPaymasterNotAWithholdingCourt)
	if err := iNotaryPaymaster.abi.UnpackIntoInterface(out, "NotAWithholdingCourt", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// INotaryPaymasterNotAnExecute represents a NotAnExecute error raised by the INotaryPaymaster contract.
type INotaryPaymasterNotAnExecute struct {
	Selector [4]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotAnExecute(bytes4 selector)
func INotaryPaymasterNotAnExecuteErrorID() common.Hash {
	return common.HexToHash("0x9574f9a61997270fffc6b1152d44a0dcce8ac439574cc2268693d17bfb78e042")
}

// UnpackNotAnExecuteError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotAnExecute(bytes4 selector)
func (iNotaryPaymaster *INotaryPaymaster) UnpackNotAnExecuteError(raw []byte) (*INotaryPaymasterNotAnExecute, error) {
	out := new(INotaryPaymasterNotAnExecute)
	if err := iNotaryPaymaster.abi.UnpackIntoInterface(out, "NotAnExecute", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// INotaryPaymasterNotDelegated represents a NotDelegated error raised by the INotaryPaymaster contract.
type INotaryPaymasterNotDelegated struct {
	Sender common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotDelegated(address sender)
func INotaryPaymasterNotDelegatedErrorID() common.Hash {
	return common.HexToHash("0xd48bb6b6afd364bbfed497c4390a28f98753f71589e43adf8e563536344d54dc")
}

// UnpackNotDelegatedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotDelegated(address sender)
func (iNotaryPaymaster *INotaryPaymaster) UnpackNotDelegatedError(raw []byte) (*INotaryPaymasterNotDelegated, error) {
	out := new(INotaryPaymasterNotDelegated)
	if err := iNotaryPaymaster.abi.UnpackIntoInterface(out, "NotDelegated", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// INotaryPaymasterNotSignedByTheNotary represents a NotSignedByTheNotary error raised by the INotaryPaymaster contract.
type INotaryPaymasterNotSignedByTheNotary struct {
	Recovered common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotSignedByTheNotary(address recovered)
func INotaryPaymasterNotSignedByTheNotaryErrorID() common.Hash {
	return common.HexToHash("0xe34581a74d3a92701872d6b9ed4cd80b13f6790d2555db1eaeaddd7a88ab8560")
}

// UnpackNotSignedByTheNotaryError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotSignedByTheNotary(address recovered)
func (iNotaryPaymaster *INotaryPaymaster) UnpackNotSignedByTheNotaryError(raw []byte) (*INotaryPaymasterNotSignedByTheNotary, error) {
	out := new(INotaryPaymasterNotSignedByTheNotary)
	if err := iNotaryPaymaster.abi.UnpackIntoInterface(out, "NotSignedByTheNotary", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// INotaryPaymasterNotTheEntryPoint represents a NotTheEntryPoint error raised by the INotaryPaymaster contract.
type INotaryPaymasterNotTheEntryPoint struct {
	Caller common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotTheEntryPoint(address caller)
func INotaryPaymasterNotTheEntryPointErrorID() common.Hash {
	return common.HexToHash("0x3e47173134a0d69cf692feb5bdb3d0cc73821aec263a590c1f135e32442fe01d")
}

// UnpackNotTheEntryPointError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotTheEntryPoint(address caller)
func (iNotaryPaymaster *INotaryPaymaster) UnpackNotTheEntryPointError(raw []byte) (*INotaryPaymasterNotTheEntryPoint, error) {
	out := new(INotaryPaymasterNotTheEntryPoint)
	if err := iNotaryPaymaster.abi.UnpackIntoInterface(out, "NotTheEntryPoint", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// INotaryPaymasterNotTheHolder represents a NotTheHolder error raised by the INotaryPaymaster contract.
type INotaryPaymasterNotTheHolder struct {
	Holder common.Address
	Sender common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotTheHolder(address holder, address sender)
func INotaryPaymasterNotTheHolderErrorID() common.Hash {
	return common.HexToHash("0xb6bd8e831b87133a951788a62f78aee0a480ce4f96f3cbdb83eec5e8697c61b7")
}

// UnpackNotTheHolderError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotTheHolder(address holder, address sender)
func (iNotaryPaymaster *INotaryPaymaster) UnpackNotTheHolderError(raw []byte) (*INotaryPaymasterNotTheHolder, error) {
	out := new(INotaryPaymasterNotTheHolder)
	if err := iNotaryPaymaster.abi.UnpackIntoInterface(out, "NotTheHolder", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// INotaryPaymasterProofHolds represents a ProofHolds error raised by the INotaryPaymaster contract.
type INotaryPaymasterProofHolds struct {
	Epoch uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ProofHolds(uint64 epoch)
func INotaryPaymasterProofHoldsErrorID() common.Hash {
	return common.HexToHash("0x64212abdba47381a8e97a2ef1859b440605c270d7fc490e549b78acd72e28e3d")
}

// UnpackProofHoldsError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ProofHolds(uint64 epoch)
func (iNotaryPaymaster *INotaryPaymaster) UnpackProofHoldsError(raw []byte) (*INotaryPaymasterProofHolds, error) {
	out := new(INotaryPaymasterProofHolds)
	if err := iNotaryPaymaster.abi.UnpackIntoInterface(out, "ProofHolds", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// INotaryPaymasterRequestStillOpen represents a RequestStillOpen error raised by the INotaryPaymaster contract.
type INotaryPaymasterRequestStillOpen struct {
	Registration [32]byte
	Epoch        uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RequestStillOpen(bytes32 registration, uint64 epoch)
func INotaryPaymasterRequestStillOpenErrorID() common.Hash {
	return common.HexToHash("0xd63a0a6b70d31a6aa96ef62b378a9e7718d08212759a72a47918dfdb97106c00")
}

// UnpackRequestStillOpenError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RequestStillOpen(bytes32 registration, uint64 epoch)
func (iNotaryPaymaster *INotaryPaymaster) UnpackRequestStillOpenError(raw []byte) (*INotaryPaymasterRequestStillOpen, error) {
	out := new(INotaryPaymasterRequestStillOpen)
	if err := iNotaryPaymaster.abi.UnpackIntoInterface(out, "RequestStillOpen", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// INotaryPaymasterSiblingsMismatch represents a SiblingsMismatch error raised by the INotaryPaymaster contract.
type INotaryPaymasterSiblingsMismatch struct {
	Signed [32]byte
	Given  [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SiblingsMismatch(bytes32 signed, bytes32 given)
func INotaryPaymasterSiblingsMismatchErrorID() common.Hash {
	return common.HexToHash("0x7e68dade1f921625838c9a65f0505353f5cae857ab8fac0fc3a0d07f29ef5fe1")
}

// UnpackSiblingsMismatchError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SiblingsMismatch(bytes32 signed, bytes32 given)
func (iNotaryPaymaster *INotaryPaymaster) UnpackSiblingsMismatchError(raw []byte) (*INotaryPaymasterSiblingsMismatch, error) {
	out := new(INotaryPaymasterSiblingsMismatch)
	if err := iNotaryPaymaster.abi.UnpackIntoInterface(out, "SiblingsMismatch", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// INotaryPaymasterTargetNotSponsored represents a TargetNotSponsored error raised by the INotaryPaymaster contract.
type INotaryPaymasterTargetNotSponsored struct {
	Target common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TargetNotSponsored(address target)
func INotaryPaymasterTargetNotSponsoredErrorID() common.Hash {
	return common.HexToHash("0xd12eb37eeb15cfb9c3994eeb23e0b5b833422adde8b57d58f14c550174359c1a")
}

// UnpackTargetNotSponsoredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TargetNotSponsored(address target)
func (iNotaryPaymaster *INotaryPaymaster) UnpackTargetNotSponsoredError(raw []byte) (*INotaryPaymasterTargetNotSponsored, error) {
	out := new(INotaryPaymasterTargetNotSponsored)
	if err := iNotaryPaymaster.abi.UnpackIntoInterface(out, "TargetNotSponsored", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// INotaryPaymasterUnstakeDelayTooShort represents a UnstakeDelayTooShort error raised by the INotaryPaymaster contract.
type INotaryPaymasterUnstakeDelayTooShort struct {
	Given uint32
	Least uint32
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error UnstakeDelayTooShort(uint32 given, uint32 least)
func INotaryPaymasterUnstakeDelayTooShortErrorID() common.Hash {
	return common.HexToHash("0x3698d575e66890e2047a43848531a08bf85ab74bca58455e7013e36b42ce0866")
}

// UnpackUnstakeDelayTooShortError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error UnstakeDelayTooShort(uint32 given, uint32 least)
func (iNotaryPaymaster *INotaryPaymaster) UnpackUnstakeDelayTooShortError(raw []byte) (*INotaryPaymasterUnstakeDelayTooShort, error) {
	out := new(INotaryPaymasterUnstakeDelayTooShort)
	if err := iNotaryPaymaster.abi.UnpackIntoInterface(out, "UnstakeDelayTooShort", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// INotaryPaymasterValueNotZero represents a ValueNotZero error raised by the INotaryPaymaster contract.
type INotaryPaymasterValueNotZero struct {
	Value *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ValueNotZero(uint256 value)
func INotaryPaymasterValueNotZeroErrorID() common.Hash {
	return common.HexToHash("0x2ee44ad4a7db9c6a0b6b17f4043c8cb20b74d223389baa5ff13802a72272751f")
}

// UnpackValueNotZeroError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ValueNotZero(uint256 value)
func (iNotaryPaymaster *INotaryPaymaster) UnpackValueNotZeroError(raw []byte) (*INotaryPaymasterValueNotZero, error) {
	out := new(INotaryPaymasterValueNotZero)
	if err := iNotaryPaymaster.abi.UnpackIntoInterface(out, "ValueNotZero", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// INotaryPaymasterWithholdingOfAnotherCourt represents a WithholdingOfAnotherCourt error raised by the INotaryPaymaster contract.
type INotaryPaymasterWithholdingOfAnotherCourt struct {
	Withholding common.Address
	Court       common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error WithholdingOfAnotherCourt(address withholding, address court)
func INotaryPaymasterWithholdingOfAnotherCourtErrorID() common.Hash {
	return common.HexToHash("0x971e79ac5d36b6454b83e1b48cbaee8b0a3bfb145c25170cd4b9613099fb03aa")
}

// UnpackWithholdingOfAnotherCourtError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error WithholdingOfAnotherCourt(address withholding, address court)
func (iNotaryPaymaster *INotaryPaymaster) UnpackWithholdingOfAnotherCourtError(raw []byte) (*INotaryPaymasterWithholdingOfAnotherCourt, error) {
	out := new(INotaryPaymasterWithholdingOfAnotherCourt)
	if err := iNotaryPaymaster.abi.UnpackIntoInterface(out, "WithholdingOfAnotherCourt", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// INotaryPaymasterWrongLedger represents a WrongLedger error raised by the INotaryPaymaster contract.
type INotaryPaymasterWrongLedger struct {
	Expected common.Address
	Given    common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error WrongLedger(address expected, address given)
func INotaryPaymasterWrongLedgerErrorID() common.Hash {
	return common.HexToHash("0x185f957c1ab2bcd607a1ccf4a1aa06b24f2dc77f61269eb02e5cc03658428df6")
}

// UnpackWrongLedgerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error WrongLedger(address expected, address given)
func (iNotaryPaymaster *INotaryPaymaster) UnpackWrongLedgerError(raw []byte) (*INotaryPaymasterWrongLedger, error) {
	out := new(INotaryPaymasterWrongLedger)
	if err := iNotaryPaymaster.abi.UnpackIntoInterface(out, "WrongLedger", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// INotaryPaymasterZeroCostCeiling represents a ZeroCostCeiling error raised by the INotaryPaymaster contract.
type INotaryPaymasterZeroCostCeiling struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ZeroCostCeiling()
func INotaryPaymasterZeroCostCeilingErrorID() common.Hash {
	return common.HexToHash("0x2378ed1645e65320516c3e7594b00e8066e9cd72b03bd818364b8e0a18e27480")
}

// UnpackZeroCostCeilingError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ZeroCostCeiling()
func (iNotaryPaymaster *INotaryPaymaster) UnpackZeroCostCeilingError(raw []byte) (*INotaryPaymasterZeroCostCeiling, error) {
	out := new(INotaryPaymasterZeroCostCeiling)
	if err := iNotaryPaymaster.abi.UnpackIntoInterface(out, "ZeroCostCeiling", raw); err != nil {
		return nil, err
	}
	return out, nil
}
