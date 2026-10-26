// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package withholdinginterface

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

// OpenRequest is an auto generated low-level Go binding around an user-defined struct.
type OpenRequest struct {
	Active bool
	Epoch  uint64
}

// Registration is an auto generated low-level Go binding around an user-defined struct.
type Registration struct {
	Holder    common.Address
	Account   [32]byte
	FromEpoch uint64
}

// Silence is an auto generated low-level Go binding around an user-defined struct.
type Silence struct {
	Registration [32]byte
	Account      [32]byte
	Holder       common.Address
	DueAt        uint64
	Recorder     common.Address
	RecordedAt   uint64
}

// StatementRequest is an auto generated low-level Go binding around an user-defined struct.
type StatementRequest struct {
	State uint8
	DueAt uint64
}

// SumNode is an auto generated low-level Go binding around an user-defined struct.
type SumNode struct {
	Hash [32]byte
	Sum  uint64
}

// IWithholdingCourtMetaData contains all meta data concerning the IWithholdingCourt contract.
var IWithholdingCourtMetaData = bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"answer\",\"inputs\":[{\"name\":\"registration\",\"type\":\"tuple\",\"internalType\":\"structRegistration\",\"components\":[{\"name\":\"holder\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"account\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"fromEpoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"epoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"amount\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"index\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"siblings\",\"type\":\"tuple[]\",\"internalType\":\"structSumNode[]\",\"components\":[{\"name\":\"hash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sum\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"answerPeriod\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"commitmentOf\",\"inputs\":[{\"name\":\"account\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"court\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIStatementCourt\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ledger\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractILiabilityLedger\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"openRequestOf\",\"inputs\":[{\"name\":\"registration\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structOpenRequest\",\"components\":[{\"name\":\"active\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"epoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"recordUnanswered\",\"inputs\":[{\"name\":\"registration\",\"type\":\"tuple\",\"internalType\":\"structRegistration\",\"components\":[{\"name\":\"holder\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"account\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"fromEpoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"epoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"recoverSigner\",\"inputs\":[{\"name\":\"registration\",\"type\":\"tuple\",\"internalType\":\"structRegistration\",\"components\":[{\"name\":\"holder\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"account\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"fromEpoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registrationDigest\",\"inputs\":[{\"name\":\"registration\",\"type\":\"tuple\",\"internalType\":\"structRegistration\",\"components\":[{\"name\":\"holder\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"account\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"fromEpoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"requestOf\",\"inputs\":[{\"name\":\"registration\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"epoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structStatementRequest\",\"components\":[{\"name\":\"state\",\"type\":\"uint8\",\"internalType\":\"enumRequestState\"},{\"name\":\"dueAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"requestStatement\",\"inputs\":[{\"name\":\"registration\",\"type\":\"tuple\",\"internalType\":\"structRegistration\",\"components\":[{\"name\":\"holder\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"account\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"fromEpoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"epoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"requestWindow\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"silenceCount\",\"inputs\":[{\"name\":\"epoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"silencesOf\",\"inputs\":[{\"name\":\"epoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple[]\",\"internalType\":\"structSilence[]\",\"components\":[{\"name\":\"registration\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"holder\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"dueAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"recorder\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"recordedAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalSilences\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"SilenceRecorded\",\"inputs\":[{\"name\":\"registration\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"epoch\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"recorder\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StatementAnswered\",\"inputs\":[{\"name\":\"registration\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"epoch\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"commitment\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"amount\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"answerer\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StatementRequested\",\"inputs\":[{\"name\":\"registration\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"epoch\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"holder\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"dueAt\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AlreadyRequested\",\"inputs\":[{\"name\":\"registration\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"epoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"AnswerDoesNotHold\",\"inputs\":[{\"name\":\"epoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"AnswerTooLate\",\"inputs\":[{\"name\":\"dueAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"blockTime\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"BeforeRegistration\",\"inputs\":[{\"name\":\"epoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"fromEpoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"DurationOutOfRange\",\"inputs\":[{\"name\":\"given\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"maximum\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"EpochNotAnchored\",\"inputs\":[{\"name\":\"epoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"NoOpenRequest\",\"inputs\":[{\"name\":\"registration\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"epoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"NotACourt\",\"inputs\":[{\"name\":\"court\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"NotOverdue\",\"inputs\":[{\"name\":\"dueAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"blockTime\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"NotSignedBySigner\",\"inputs\":[{\"name\":\"fromEpoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"signer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"recovered\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"NotTheHolder\",\"inputs\":[{\"name\":\"holder\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"RequestStillOpen\",\"inputs\":[{\"name\":\"registration\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"epoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"RequestWindowClosed\",\"inputs\":[{\"name\":\"epoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"closedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]",
	ID:  "IWithholdingCourt",
}

// IWithholdingCourt is an auto generated Go binding around an Ethereum contract.
type IWithholdingCourt struct {
	abi abi.ABI
}

// GetABI returns the ABI associated with this contract binding.
func (c *IWithholdingCourt) GetABI() abi.ABI {
	return c.abi
}

// NewIWithholdingCourt creates a new instance of IWithholdingCourt.
func NewIWithholdingCourt() *IWithholdingCourt {
	parsed, err := IWithholdingCourtMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &IWithholdingCourt{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *IWithholdingCourt) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackAnswer is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x53622e82.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function answer((address,bytes32,uint64) registration, uint64 epoch, bytes32 salt, uint64 amount, uint64 index, (bytes32,uint64)[] siblings) returns(bytes32)
func (iWithholdingCourt *IWithholdingCourt) PackAnswer(registration Registration, epoch uint64, salt [32]byte, amount uint64, index uint64, siblings []SumNode) []byte {
	enc, err := iWithholdingCourt.abi.Pack("answer", registration, epoch, salt, amount, index, siblings)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAnswer is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x53622e82.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function answer((address,bytes32,uint64) registration, uint64 epoch, bytes32 salt, uint64 amount, uint64 index, (bytes32,uint64)[] siblings) returns(bytes32)
func (iWithholdingCourt *IWithholdingCourt) TryPackAnswer(registration Registration, epoch uint64, salt [32]byte, amount uint64, index uint64, siblings []SumNode) ([]byte, error) {
	return iWithholdingCourt.abi.Pack("answer", registration, epoch, salt, amount, index, siblings)
}

// UnpackAnswer is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x53622e82.
//
// Solidity: function answer((address,bytes32,uint64) registration, uint64 epoch, bytes32 salt, uint64 amount, uint64 index, (bytes32,uint64)[] siblings) returns(bytes32)
func (iWithholdingCourt *IWithholdingCourt) UnpackAnswer(data []byte) ([32]byte, error) {
	out, err := iWithholdingCourt.abi.Unpack("answer", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackAnswerPeriod is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd7e1ff0f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function answerPeriod() view returns(uint64)
func (iWithholdingCourt *IWithholdingCourt) PackAnswerPeriod() []byte {
	enc, err := iWithholdingCourt.abi.Pack("answerPeriod")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAnswerPeriod is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd7e1ff0f.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function answerPeriod() view returns(uint64)
func (iWithholdingCourt *IWithholdingCourt) TryPackAnswerPeriod() ([]byte, error) {
	return iWithholdingCourt.abi.Pack("answerPeriod")
}

// UnpackAnswerPeriod is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xd7e1ff0f.
//
// Solidity: function answerPeriod() view returns(uint64)
func (iWithholdingCourt *IWithholdingCourt) UnpackAnswerPeriod(data []byte) (uint64, error) {
	out, err := iWithholdingCourt.abi.Unpack("answerPeriod", data)
	if err != nil {
		return *new(uint64), err
	}
	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)
	return out0, nil
}

// PackCommitmentOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9060e949.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function commitmentOf(bytes32 account, bytes32 salt) pure returns(bytes32)
func (iWithholdingCourt *IWithholdingCourt) PackCommitmentOf(account [32]byte, salt [32]byte) []byte {
	enc, err := iWithholdingCourt.abi.Pack("commitmentOf", account, salt)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCommitmentOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9060e949.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function commitmentOf(bytes32 account, bytes32 salt) pure returns(bytes32)
func (iWithholdingCourt *IWithholdingCourt) TryPackCommitmentOf(account [32]byte, salt [32]byte) ([]byte, error) {
	return iWithholdingCourt.abi.Pack("commitmentOf", account, salt)
}

// UnpackCommitmentOf is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x9060e949.
//
// Solidity: function commitmentOf(bytes32 account, bytes32 salt) pure returns(bytes32)
func (iWithholdingCourt *IWithholdingCourt) UnpackCommitmentOf(data []byte) ([32]byte, error) {
	out, err := iWithholdingCourt.abi.Unpack("commitmentOf", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackCourt is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x00971a30.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function court() view returns(address)
func (iWithholdingCourt *IWithholdingCourt) PackCourt() []byte {
	enc, err := iWithholdingCourt.abi.Pack("court")
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
func (iWithholdingCourt *IWithholdingCourt) TryPackCourt() ([]byte, error) {
	return iWithholdingCourt.abi.Pack("court")
}

// UnpackCourt is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x00971a30.
//
// Solidity: function court() view returns(address)
func (iWithholdingCourt *IWithholdingCourt) UnpackCourt(data []byte) (common.Address, error) {
	out, err := iWithholdingCourt.abi.Unpack("court", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackLedger is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x56397c35.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function ledger() view returns(address)
func (iWithholdingCourt *IWithholdingCourt) PackLedger() []byte {
	enc, err := iWithholdingCourt.abi.Pack("ledger")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackLedger is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x56397c35.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function ledger() view returns(address)
func (iWithholdingCourt *IWithholdingCourt) TryPackLedger() ([]byte, error) {
	return iWithholdingCourt.abi.Pack("ledger")
}

// UnpackLedger is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x56397c35.
//
// Solidity: function ledger() view returns(address)
func (iWithholdingCourt *IWithholdingCourt) UnpackLedger(data []byte) (common.Address, error) {
	out, err := iWithholdingCourt.abi.Unpack("ledger", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackOpenRequestOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe3ff236d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function openRequestOf(bytes32 registration) view returns((bool,uint64))
func (iWithholdingCourt *IWithholdingCourt) PackOpenRequestOf(registration [32]byte) []byte {
	enc, err := iWithholdingCourt.abi.Pack("openRequestOf", registration)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackOpenRequestOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe3ff236d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function openRequestOf(bytes32 registration) view returns((bool,uint64))
func (iWithholdingCourt *IWithholdingCourt) TryPackOpenRequestOf(registration [32]byte) ([]byte, error) {
	return iWithholdingCourt.abi.Pack("openRequestOf", registration)
}

// UnpackOpenRequestOf is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe3ff236d.
//
// Solidity: function openRequestOf(bytes32 registration) view returns((bool,uint64))
func (iWithholdingCourt *IWithholdingCourt) UnpackOpenRequestOf(data []byte) (OpenRequest, error) {
	out, err := iWithholdingCourt.abi.Unpack("openRequestOf", data)
	if err != nil {
		return *new(OpenRequest), err
	}
	out0 := *abi.ConvertType(out[0], new(OpenRequest)).(*OpenRequest)
	return out0, nil
}

// PackRecordUnanswered is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x404bd102.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function recordUnanswered((address,bytes32,uint64) registration, uint64 epoch) returns(bytes32)
func (iWithholdingCourt *IWithholdingCourt) PackRecordUnanswered(registration Registration, epoch uint64) []byte {
	enc, err := iWithholdingCourt.abi.Pack("recordUnanswered", registration, epoch)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRecordUnanswered is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x404bd102.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function recordUnanswered((address,bytes32,uint64) registration, uint64 epoch) returns(bytes32)
func (iWithholdingCourt *IWithholdingCourt) TryPackRecordUnanswered(registration Registration, epoch uint64) ([]byte, error) {
	return iWithholdingCourt.abi.Pack("recordUnanswered", registration, epoch)
}

// UnpackRecordUnanswered is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x404bd102.
//
// Solidity: function recordUnanswered((address,bytes32,uint64) registration, uint64 epoch) returns(bytes32)
func (iWithholdingCourt *IWithholdingCourt) UnpackRecordUnanswered(data []byte) ([32]byte, error) {
	out, err := iWithholdingCourt.abi.Unpack("recordUnanswered", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackRecoverSigner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdebd8cff.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function recoverSigner((address,bytes32,uint64) registration, bytes signature) view returns(address)
func (iWithholdingCourt *IWithholdingCourt) PackRecoverSigner(registration Registration, signature []byte) []byte {
	enc, err := iWithholdingCourt.abi.Pack("recoverSigner", registration, signature)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRecoverSigner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdebd8cff.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function recoverSigner((address,bytes32,uint64) registration, bytes signature) view returns(address)
func (iWithholdingCourt *IWithholdingCourt) TryPackRecoverSigner(registration Registration, signature []byte) ([]byte, error) {
	return iWithholdingCourt.abi.Pack("recoverSigner", registration, signature)
}

// UnpackRecoverSigner is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xdebd8cff.
//
// Solidity: function recoverSigner((address,bytes32,uint64) registration, bytes signature) view returns(address)
func (iWithholdingCourt *IWithholdingCourt) UnpackRecoverSigner(data []byte) (common.Address, error) {
	out, err := iWithholdingCourt.abi.Unpack("recoverSigner", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackRegistrationDigest is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x86579e95.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function registrationDigest((address,bytes32,uint64) registration) view returns(bytes32)
func (iWithholdingCourt *IWithholdingCourt) PackRegistrationDigest(registration Registration) []byte {
	enc, err := iWithholdingCourt.abi.Pack("registrationDigest", registration)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRegistrationDigest is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x86579e95.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function registrationDigest((address,bytes32,uint64) registration) view returns(bytes32)
func (iWithholdingCourt *IWithholdingCourt) TryPackRegistrationDigest(registration Registration) ([]byte, error) {
	return iWithholdingCourt.abi.Pack("registrationDigest", registration)
}

// UnpackRegistrationDigest is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x86579e95.
//
// Solidity: function registrationDigest((address,bytes32,uint64) registration) view returns(bytes32)
func (iWithholdingCourt *IWithholdingCourt) UnpackRegistrationDigest(data []byte) ([32]byte, error) {
	out, err := iWithholdingCourt.abi.Unpack("registrationDigest", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackRequestOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7d917512.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function requestOf(bytes32 registration, uint64 epoch) view returns((uint8,uint64))
func (iWithholdingCourt *IWithholdingCourt) PackRequestOf(registration [32]byte, epoch uint64) []byte {
	enc, err := iWithholdingCourt.abi.Pack("requestOf", registration, epoch)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRequestOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7d917512.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function requestOf(bytes32 registration, uint64 epoch) view returns((uint8,uint64))
func (iWithholdingCourt *IWithholdingCourt) TryPackRequestOf(registration [32]byte, epoch uint64) ([]byte, error) {
	return iWithholdingCourt.abi.Pack("requestOf", registration, epoch)
}

// UnpackRequestOf is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x7d917512.
//
// Solidity: function requestOf(bytes32 registration, uint64 epoch) view returns((uint8,uint64))
func (iWithholdingCourt *IWithholdingCourt) UnpackRequestOf(data []byte) (StatementRequest, error) {
	out, err := iWithholdingCourt.abi.Unpack("requestOf", data)
	if err != nil {
		return *new(StatementRequest), err
	}
	out0 := *abi.ConvertType(out[0], new(StatementRequest)).(*StatementRequest)
	return out0, nil
}

// PackRequestStatement is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x42b2f6ed.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function requestStatement((address,bytes32,uint64) registration, bytes signature, uint64 epoch) returns(bytes32)
func (iWithholdingCourt *IWithholdingCourt) PackRequestStatement(registration Registration, signature []byte, epoch uint64) []byte {
	enc, err := iWithholdingCourt.abi.Pack("requestStatement", registration, signature, epoch)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRequestStatement is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x42b2f6ed.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function requestStatement((address,bytes32,uint64) registration, bytes signature, uint64 epoch) returns(bytes32)
func (iWithholdingCourt *IWithholdingCourt) TryPackRequestStatement(registration Registration, signature []byte, epoch uint64) ([]byte, error) {
	return iWithholdingCourt.abi.Pack("requestStatement", registration, signature, epoch)
}

// UnpackRequestStatement is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x42b2f6ed.
//
// Solidity: function requestStatement((address,bytes32,uint64) registration, bytes signature, uint64 epoch) returns(bytes32)
func (iWithholdingCourt *IWithholdingCourt) UnpackRequestStatement(data []byte) ([32]byte, error) {
	out, err := iWithholdingCourt.abi.Unpack("requestStatement", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackRequestWindow is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1e636344.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function requestWindow() view returns(uint64)
func (iWithholdingCourt *IWithholdingCourt) PackRequestWindow() []byte {
	enc, err := iWithholdingCourt.abi.Pack("requestWindow")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRequestWindow is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1e636344.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function requestWindow() view returns(uint64)
func (iWithholdingCourt *IWithholdingCourt) TryPackRequestWindow() ([]byte, error) {
	return iWithholdingCourt.abi.Pack("requestWindow")
}

// UnpackRequestWindow is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x1e636344.
//
// Solidity: function requestWindow() view returns(uint64)
func (iWithholdingCourt *IWithholdingCourt) UnpackRequestWindow(data []byte) (uint64, error) {
	out, err := iWithholdingCourt.abi.Unpack("requestWindow", data)
	if err != nil {
		return *new(uint64), err
	}
	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)
	return out0, nil
}

// PackSilenceCount is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xba49f58d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function silenceCount(uint64 epoch) view returns(uint256)
func (iWithholdingCourt *IWithholdingCourt) PackSilenceCount(epoch uint64) []byte {
	enc, err := iWithholdingCourt.abi.Pack("silenceCount", epoch)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSilenceCount is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xba49f58d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function silenceCount(uint64 epoch) view returns(uint256)
func (iWithholdingCourt *IWithholdingCourt) TryPackSilenceCount(epoch uint64) ([]byte, error) {
	return iWithholdingCourt.abi.Pack("silenceCount", epoch)
}

// UnpackSilenceCount is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xba49f58d.
//
// Solidity: function silenceCount(uint64 epoch) view returns(uint256)
func (iWithholdingCourt *IWithholdingCourt) UnpackSilenceCount(data []byte) (*big.Int, error) {
	out, err := iWithholdingCourt.abi.Unpack("silenceCount", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackSilencesOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x51722cad.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function silencesOf(uint64 epoch) view returns((bytes32,bytes32,address,uint64,address,uint64)[])
func (iWithholdingCourt *IWithholdingCourt) PackSilencesOf(epoch uint64) []byte {
	enc, err := iWithholdingCourt.abi.Pack("silencesOf", epoch)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSilencesOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x51722cad.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function silencesOf(uint64 epoch) view returns((bytes32,bytes32,address,uint64,address,uint64)[])
func (iWithholdingCourt *IWithholdingCourt) TryPackSilencesOf(epoch uint64) ([]byte, error) {
	return iWithholdingCourt.abi.Pack("silencesOf", epoch)
}

// UnpackSilencesOf is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x51722cad.
//
// Solidity: function silencesOf(uint64 epoch) view returns((bytes32,bytes32,address,uint64,address,uint64)[])
func (iWithholdingCourt *IWithholdingCourt) UnpackSilencesOf(data []byte) ([]Silence, error) {
	out, err := iWithholdingCourt.abi.Unpack("silencesOf", data)
	if err != nil {
		return *new([]Silence), err
	}
	out0 := *abi.ConvertType(out[0], new([]Silence)).(*[]Silence)
	return out0, nil
}

// PackTotalSilences is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x74f9214e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function totalSilences() view returns(uint256)
func (iWithholdingCourt *IWithholdingCourt) PackTotalSilences() []byte {
	enc, err := iWithholdingCourt.abi.Pack("totalSilences")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTotalSilences is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x74f9214e.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function totalSilences() view returns(uint256)
func (iWithholdingCourt *IWithholdingCourt) TryPackTotalSilences() ([]byte, error) {
	return iWithholdingCourt.abi.Pack("totalSilences")
}

// UnpackTotalSilences is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x74f9214e.
//
// Solidity: function totalSilences() view returns(uint256)
func (iWithholdingCourt *IWithholdingCourt) UnpackTotalSilences(data []byte) (*big.Int, error) {
	out, err := iWithholdingCourt.abi.Unpack("totalSilences", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// IWithholdingCourtSilenceRecorded represents a SilenceRecorded event raised by the IWithholdingCourt contract.
type IWithholdingCourtSilenceRecorded struct {
	Registration [32]byte
	Account      [32]byte
	Epoch        uint64
	Recorder     common.Address
	Raw          *types.Log // Blockchain specific contextual infos
}

const IWithholdingCourtSilenceRecordedEventName = "SilenceRecorded"

// ContractEventName returns the user-defined event name.
func (IWithholdingCourtSilenceRecorded) ContractEventName() string {
	return IWithholdingCourtSilenceRecordedEventName
}

// UnpackSilenceRecordedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event SilenceRecorded(bytes32 indexed registration, bytes32 indexed account, uint64 indexed epoch, address recorder)
func (iWithholdingCourt *IWithholdingCourt) UnpackSilenceRecordedEvent(log *types.Log) (*IWithholdingCourtSilenceRecorded, error) {
	event := "SilenceRecorded"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != iWithholdingCourt.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(IWithholdingCourtSilenceRecorded)
	if len(log.Data) > 0 {
		if err := iWithholdingCourt.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iWithholdingCourt.abi.Events[event].Inputs {
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

// IWithholdingCourtStatementAnswered represents a StatementAnswered event raised by the IWithholdingCourt contract.
type IWithholdingCourtStatementAnswered struct {
	Registration [32]byte
	Account      [32]byte
	Epoch        uint64
	Commitment   [32]byte
	Amount       uint64
	Answerer     common.Address
	Raw          *types.Log // Blockchain specific contextual infos
}

const IWithholdingCourtStatementAnsweredEventName = "StatementAnswered"

// ContractEventName returns the user-defined event name.
func (IWithholdingCourtStatementAnswered) ContractEventName() string {
	return IWithholdingCourtStatementAnsweredEventName
}

// UnpackStatementAnsweredEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event StatementAnswered(bytes32 indexed registration, bytes32 indexed account, uint64 indexed epoch, bytes32 commitment, uint64 amount, address answerer)
func (iWithholdingCourt *IWithholdingCourt) UnpackStatementAnsweredEvent(log *types.Log) (*IWithholdingCourtStatementAnswered, error) {
	event := "StatementAnswered"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != iWithholdingCourt.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(IWithholdingCourtStatementAnswered)
	if len(log.Data) > 0 {
		if err := iWithholdingCourt.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iWithholdingCourt.abi.Events[event].Inputs {
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

// IWithholdingCourtStatementRequested represents a StatementRequested event raised by the IWithholdingCourt contract.
type IWithholdingCourtStatementRequested struct {
	Registration [32]byte
	Account      [32]byte
	Epoch        uint64
	Holder       common.Address
	DueAt        uint64
	Raw          *types.Log // Blockchain specific contextual infos
}

const IWithholdingCourtStatementRequestedEventName = "StatementRequested"

// ContractEventName returns the user-defined event name.
func (IWithholdingCourtStatementRequested) ContractEventName() string {
	return IWithholdingCourtStatementRequestedEventName
}

// UnpackStatementRequestedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event StatementRequested(bytes32 indexed registration, bytes32 indexed account, uint64 indexed epoch, address holder, uint64 dueAt)
func (iWithholdingCourt *IWithholdingCourt) UnpackStatementRequestedEvent(log *types.Log) (*IWithholdingCourtStatementRequested, error) {
	event := "StatementRequested"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != iWithholdingCourt.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(IWithholdingCourtStatementRequested)
	if len(log.Data) > 0 {
		if err := iWithholdingCourt.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iWithholdingCourt.abi.Events[event].Inputs {
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
func (iWithholdingCourt *IWithholdingCourt) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], iWithholdingCourt.abi.Errors["AlreadyRequested"].ID.Bytes()[:4]) {
		return iWithholdingCourt.UnpackAlreadyRequestedError(raw[4:])
	}
	if bytes.Equal(raw[:4], iWithholdingCourt.abi.Errors["AnswerDoesNotHold"].ID.Bytes()[:4]) {
		return iWithholdingCourt.UnpackAnswerDoesNotHoldError(raw[4:])
	}
	if bytes.Equal(raw[:4], iWithholdingCourt.abi.Errors["AnswerTooLate"].ID.Bytes()[:4]) {
		return iWithholdingCourt.UnpackAnswerTooLateError(raw[4:])
	}
	if bytes.Equal(raw[:4], iWithholdingCourt.abi.Errors["BeforeRegistration"].ID.Bytes()[:4]) {
		return iWithholdingCourt.UnpackBeforeRegistrationError(raw[4:])
	}
	if bytes.Equal(raw[:4], iWithholdingCourt.abi.Errors["DurationOutOfRange"].ID.Bytes()[:4]) {
		return iWithholdingCourt.UnpackDurationOutOfRangeError(raw[4:])
	}
	if bytes.Equal(raw[:4], iWithholdingCourt.abi.Errors["EpochNotAnchored"].ID.Bytes()[:4]) {
		return iWithholdingCourt.UnpackEpochNotAnchoredError(raw[4:])
	}
	if bytes.Equal(raw[:4], iWithholdingCourt.abi.Errors["NoOpenRequest"].ID.Bytes()[:4]) {
		return iWithholdingCourt.UnpackNoOpenRequestError(raw[4:])
	}
	if bytes.Equal(raw[:4], iWithholdingCourt.abi.Errors["NotACourt"].ID.Bytes()[:4]) {
		return iWithholdingCourt.UnpackNotACourtError(raw[4:])
	}
	if bytes.Equal(raw[:4], iWithholdingCourt.abi.Errors["NotOverdue"].ID.Bytes()[:4]) {
		return iWithholdingCourt.UnpackNotOverdueError(raw[4:])
	}
	if bytes.Equal(raw[:4], iWithholdingCourt.abi.Errors["NotSignedBySigner"].ID.Bytes()[:4]) {
		return iWithholdingCourt.UnpackNotSignedBySignerError(raw[4:])
	}
	if bytes.Equal(raw[:4], iWithholdingCourt.abi.Errors["NotTheHolder"].ID.Bytes()[:4]) {
		return iWithholdingCourt.UnpackNotTheHolderError(raw[4:])
	}
	if bytes.Equal(raw[:4], iWithholdingCourt.abi.Errors["RequestStillOpen"].ID.Bytes()[:4]) {
		return iWithholdingCourt.UnpackRequestStillOpenError(raw[4:])
	}
	if bytes.Equal(raw[:4], iWithholdingCourt.abi.Errors["RequestWindowClosed"].ID.Bytes()[:4]) {
		return iWithholdingCourt.UnpackRequestWindowClosedError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// IWithholdingCourtAlreadyRequested represents a AlreadyRequested error raised by the IWithholdingCourt contract.
type IWithholdingCourtAlreadyRequested struct {
	Registration [32]byte
	Epoch        uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AlreadyRequested(bytes32 registration, uint64 epoch)
func IWithholdingCourtAlreadyRequestedErrorID() common.Hash {
	return common.HexToHash("0x46cdd0fb7d2db3112c73591b2011dc63aca93f893280dca515f60dd290122c34")
}

// UnpackAlreadyRequestedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AlreadyRequested(bytes32 registration, uint64 epoch)
func (iWithholdingCourt *IWithholdingCourt) UnpackAlreadyRequestedError(raw []byte) (*IWithholdingCourtAlreadyRequested, error) {
	out := new(IWithholdingCourtAlreadyRequested)
	if err := iWithholdingCourt.abi.UnpackIntoInterface(out, "AlreadyRequested", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IWithholdingCourtAnswerDoesNotHold represents a AnswerDoesNotHold error raised by the IWithholdingCourt contract.
type IWithholdingCourtAnswerDoesNotHold struct {
	Epoch uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AnswerDoesNotHold(uint64 epoch)
func IWithholdingCourtAnswerDoesNotHoldErrorID() common.Hash {
	return common.HexToHash("0x2fa1f48b39b9b6b6b9cfebd069c79075983f2ff2db1d8c7dcfd897dc5429dd28")
}

// UnpackAnswerDoesNotHoldError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AnswerDoesNotHold(uint64 epoch)
func (iWithholdingCourt *IWithholdingCourt) UnpackAnswerDoesNotHoldError(raw []byte) (*IWithholdingCourtAnswerDoesNotHold, error) {
	out := new(IWithholdingCourtAnswerDoesNotHold)
	if err := iWithholdingCourt.abi.UnpackIntoInterface(out, "AnswerDoesNotHold", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IWithholdingCourtAnswerTooLate represents a AnswerTooLate error raised by the IWithholdingCourt contract.
type IWithholdingCourtAnswerTooLate struct {
	DueAt     uint64
	BlockTime uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AnswerTooLate(uint64 dueAt, uint64 blockTime)
func IWithholdingCourtAnswerTooLateErrorID() common.Hash {
	return common.HexToHash("0x5c3cd680ab4930d0fa30c8646c992355155b81cc12e62635206ad8881bfcfc6c")
}

// UnpackAnswerTooLateError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AnswerTooLate(uint64 dueAt, uint64 blockTime)
func (iWithholdingCourt *IWithholdingCourt) UnpackAnswerTooLateError(raw []byte) (*IWithholdingCourtAnswerTooLate, error) {
	out := new(IWithholdingCourtAnswerTooLate)
	if err := iWithholdingCourt.abi.UnpackIntoInterface(out, "AnswerTooLate", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IWithholdingCourtBeforeRegistration represents a BeforeRegistration error raised by the IWithholdingCourt contract.
type IWithholdingCourtBeforeRegistration struct {
	Epoch     uint64
	FromEpoch uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error BeforeRegistration(uint64 epoch, uint64 fromEpoch)
func IWithholdingCourtBeforeRegistrationErrorID() common.Hash {
	return common.HexToHash("0xd30d092f5afe15b07c027869b0e62d560be8c5e7d28f175dd1f43a5bbecfb25b")
}

// UnpackBeforeRegistrationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error BeforeRegistration(uint64 epoch, uint64 fromEpoch)
func (iWithholdingCourt *IWithholdingCourt) UnpackBeforeRegistrationError(raw []byte) (*IWithholdingCourtBeforeRegistration, error) {
	out := new(IWithholdingCourtBeforeRegistration)
	if err := iWithholdingCourt.abi.UnpackIntoInterface(out, "BeforeRegistration", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IWithholdingCourtDurationOutOfRange represents a DurationOutOfRange error raised by the IWithholdingCourt contract.
type IWithholdingCourtDurationOutOfRange struct {
	Given   uint64
	Maximum uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error DurationOutOfRange(uint64 given, uint64 maximum)
func IWithholdingCourtDurationOutOfRangeErrorID() common.Hash {
	return common.HexToHash("0xf4bc09016fa12667f401ba3f9820d43ded35635f99b8b39a6a755a97e47db3bb")
}

// UnpackDurationOutOfRangeError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error DurationOutOfRange(uint64 given, uint64 maximum)
func (iWithholdingCourt *IWithholdingCourt) UnpackDurationOutOfRangeError(raw []byte) (*IWithholdingCourtDurationOutOfRange, error) {
	out := new(IWithholdingCourtDurationOutOfRange)
	if err := iWithholdingCourt.abi.UnpackIntoInterface(out, "DurationOutOfRange", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IWithholdingCourtEpochNotAnchored represents a EpochNotAnchored error raised by the IWithholdingCourt contract.
type IWithholdingCourtEpochNotAnchored struct {
	Epoch uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error EpochNotAnchored(uint64 epoch)
func IWithholdingCourtEpochNotAnchoredErrorID() common.Hash {
	return common.HexToHash("0xfdd66f09e3e0335148d05ffa812e23b41ab1c93ea73617ca6ddd26710f05879d")
}

// UnpackEpochNotAnchoredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error EpochNotAnchored(uint64 epoch)
func (iWithholdingCourt *IWithholdingCourt) UnpackEpochNotAnchoredError(raw []byte) (*IWithholdingCourtEpochNotAnchored, error) {
	out := new(IWithholdingCourtEpochNotAnchored)
	if err := iWithholdingCourt.abi.UnpackIntoInterface(out, "EpochNotAnchored", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IWithholdingCourtNoOpenRequest represents a NoOpenRequest error raised by the IWithholdingCourt contract.
type IWithholdingCourtNoOpenRequest struct {
	Registration [32]byte
	Epoch        uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NoOpenRequest(bytes32 registration, uint64 epoch)
func IWithholdingCourtNoOpenRequestErrorID() common.Hash {
	return common.HexToHash("0x12ef5f2cac24af69d330f4a7ae9d2d98b14d43bae36fbacf6df51b7f36fd3c8d")
}

// UnpackNoOpenRequestError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NoOpenRequest(bytes32 registration, uint64 epoch)
func (iWithholdingCourt *IWithholdingCourt) UnpackNoOpenRequestError(raw []byte) (*IWithholdingCourtNoOpenRequest, error) {
	out := new(IWithholdingCourtNoOpenRequest)
	if err := iWithholdingCourt.abi.UnpackIntoInterface(out, "NoOpenRequest", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IWithholdingCourtNotACourt represents a NotACourt error raised by the IWithholdingCourt contract.
type IWithholdingCourtNotACourt struct {
	Court common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotACourt(address court)
func IWithholdingCourtNotACourtErrorID() common.Hash {
	return common.HexToHash("0xff4c4161b866d228828099f46f8d1024545fc32c86bb16e694c95f55684edd9e")
}

// UnpackNotACourtError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotACourt(address court)
func (iWithholdingCourt *IWithholdingCourt) UnpackNotACourtError(raw []byte) (*IWithholdingCourtNotACourt, error) {
	out := new(IWithholdingCourtNotACourt)
	if err := iWithholdingCourt.abi.UnpackIntoInterface(out, "NotACourt", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IWithholdingCourtNotOverdue represents a NotOverdue error raised by the IWithholdingCourt contract.
type IWithholdingCourtNotOverdue struct {
	DueAt     uint64
	BlockTime uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotOverdue(uint64 dueAt, uint64 blockTime)
func IWithholdingCourtNotOverdueErrorID() common.Hash {
	return common.HexToHash("0xbd173e570140cc2e12f7e46a73691e55b8e11f9cd24654d209d94e6d000e0497")
}

// UnpackNotOverdueError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotOverdue(uint64 dueAt, uint64 blockTime)
func (iWithholdingCourt *IWithholdingCourt) UnpackNotOverdueError(raw []byte) (*IWithholdingCourtNotOverdue, error) {
	out := new(IWithholdingCourtNotOverdue)
	if err := iWithholdingCourt.abi.UnpackIntoInterface(out, "NotOverdue", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IWithholdingCourtNotSignedBySigner represents a NotSignedBySigner error raised by the IWithholdingCourt contract.
type IWithholdingCourtNotSignedBySigner struct {
	FromEpoch uint64
	Signer    common.Address
	Recovered common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotSignedBySigner(uint64 fromEpoch, address signer, address recovered)
func IWithholdingCourtNotSignedBySignerErrorID() common.Hash {
	return common.HexToHash("0xd7eecde17baa6c3c5da78ec7d0f325a551b50706d7779d410753e9ea173f5813")
}

// UnpackNotSignedBySignerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotSignedBySigner(uint64 fromEpoch, address signer, address recovered)
func (iWithholdingCourt *IWithholdingCourt) UnpackNotSignedBySignerError(raw []byte) (*IWithholdingCourtNotSignedBySigner, error) {
	out := new(IWithholdingCourtNotSignedBySigner)
	if err := iWithholdingCourt.abi.UnpackIntoInterface(out, "NotSignedBySigner", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IWithholdingCourtNotTheHolder represents a NotTheHolder error raised by the IWithholdingCourt contract.
type IWithholdingCourtNotTheHolder struct {
	Holder common.Address
	Caller common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotTheHolder(address holder, address caller)
func IWithholdingCourtNotTheHolderErrorID() common.Hash {
	return common.HexToHash("0xb6bd8e831b87133a951788a62f78aee0a480ce4f96f3cbdb83eec5e8697c61b7")
}

// UnpackNotTheHolderError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotTheHolder(address holder, address caller)
func (iWithholdingCourt *IWithholdingCourt) UnpackNotTheHolderError(raw []byte) (*IWithholdingCourtNotTheHolder, error) {
	out := new(IWithholdingCourtNotTheHolder)
	if err := iWithholdingCourt.abi.UnpackIntoInterface(out, "NotTheHolder", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IWithholdingCourtRequestStillOpen represents a RequestStillOpen error raised by the IWithholdingCourt contract.
type IWithholdingCourtRequestStillOpen struct {
	Registration [32]byte
	Epoch        uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RequestStillOpen(bytes32 registration, uint64 epoch)
func IWithholdingCourtRequestStillOpenErrorID() common.Hash {
	return common.HexToHash("0xd63a0a6b70d31a6aa96ef62b378a9e7718d08212759a72a47918dfdb97106c00")
}

// UnpackRequestStillOpenError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RequestStillOpen(bytes32 registration, uint64 epoch)
func (iWithholdingCourt *IWithholdingCourt) UnpackRequestStillOpenError(raw []byte) (*IWithholdingCourtRequestStillOpen, error) {
	out := new(IWithholdingCourtRequestStillOpen)
	if err := iWithholdingCourt.abi.UnpackIntoInterface(out, "RequestStillOpen", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IWithholdingCourtRequestWindowClosed represents a RequestWindowClosed error raised by the IWithholdingCourt contract.
type IWithholdingCourtRequestWindowClosed struct {
	Epoch    uint64
	ClosedAt *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RequestWindowClosed(uint64 epoch, uint256 closedAt)
func IWithholdingCourtRequestWindowClosedErrorID() common.Hash {
	return common.HexToHash("0xa87a071a48c134dae493271860218ee72ff911a816b5bc0dbdc9aa522b7979bd")
}

// UnpackRequestWindowClosedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RequestWindowClosed(uint64 epoch, uint256 closedAt)
func (iWithholdingCourt *IWithholdingCourt) UnpackRequestWindowClosedError(raw []byte) (*IWithholdingCourtRequestWindowClosed, error) {
	out := new(IWithholdingCourtRequestWindowClosed)
	if err := iWithholdingCourt.abi.UnpackIntoInterface(out, "RequestWindowClosed", raw); err != nil {
		return nil, err
	}
	return out, nil
}
