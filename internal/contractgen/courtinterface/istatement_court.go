// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package courtinterface

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

// EpochStatement is an auto generated low-level Go binding around an user-defined struct.
type EpochStatement struct {
	Ledger       common.Address
	Epoch        uint64
	Commitment   [32]byte
	Amount       uint64
	Index        uint64
	SiblingsHash [32]byte
}

// Fault is an auto generated low-level Go binding around an user-defined struct.
type Fault struct {
	Kind     uint8
	Evidence [32]byte
	Prover   common.Address
	ProvenAt uint64
}

// PendingSigner is an auto generated low-level Go binding around an user-defined struct.
type PendingSigner struct {
	Signer      common.Address
	FromEpoch   uint64
	ActivatesAt uint64
}

// SignerTerm is an auto generated low-level Go binding around an user-defined struct.
type SignerTerm struct {
	Signer    common.Address
	FromEpoch uint64
}

// SumNode is an auto generated low-level Go binding around an user-defined struct.
type SumNode struct {
	Hash [32]byte
	Sum  uint64
}

// IStatementCourtMetaData contains all meta data concerning the IStatementCourt contract.
var IStatementCourtMetaData = bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"activateSigner\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"cancelSigner\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"faultCount\",\"inputs\":[{\"name\":\"epoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"faultsOf\",\"inputs\":[{\"name\":\"epoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple[]\",\"internalType\":\"structFault[]\",\"components\":[{\"name\":\"kind\",\"type\":\"uint8\",\"internalType\":\"enumFaultKind\"},{\"name\":\"evidence\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"prover\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"provenAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isProven\",\"inputs\":[{\"name\":\"evidence\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ledger\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractILiabilityLedger\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingSigner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structPendingSigner\",\"components\":[{\"name\":\"signer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"fromEpoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"activatesAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proveBadProof\",\"inputs\":[{\"name\":\"statement\",\"type\":\"tuple\",\"internalType\":\"structEpochStatement\",\"components\":[{\"name\":\"ledger\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"epoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"commitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"amount\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"index\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"siblingsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"siblings\",\"type\":\"tuple[]\",\"internalType\":\"structSumNode[]\",\"components\":[{\"name\":\"hash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sum\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"proveEquivocation\",\"inputs\":[{\"name\":\"first\",\"type\":\"tuple\",\"internalType\":\"structEpochStatement\",\"components\":[{\"name\":\"ledger\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"epoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"commitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"amount\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"index\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"siblingsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"firstSignature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"second\",\"type\":\"tuple\",\"internalType\":\"structEpochStatement\",\"components\":[{\"name\":\"ledger\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"epoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"commitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"amount\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"index\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"siblingsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"secondSignature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"recoverSigner\",\"inputs\":[{\"name\":\"statement\",\"type\":\"tuple\",\"internalType\":\"structEpochStatement\",\"components\":[{\"name\":\"ledger\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"epoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"commitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"amount\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"index\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"siblingsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"scheduleSigner\",\"inputs\":[{\"name\":\"signer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"fromEpoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"siblingsHash\",\"inputs\":[{\"name\":\"siblings\",\"type\":\"tuple[]\",\"internalType\":\"structSumNode[]\",\"components\":[{\"name\":\"hash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sum\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"signerChangeDelay\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"signerOf\",\"inputs\":[{\"name\":\"epoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"signerTerms\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple[]\",\"internalType\":\"structSignerTerm[]\",\"components\":[{\"name\":\"signer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"fromEpoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"statementDigest\",\"inputs\":[{\"name\":\"statement\",\"type\":\"tuple\",\"internalType\":\"structEpochStatement\",\"components\":[{\"name\":\"ledger\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"epoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"commitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"amount\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"index\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"siblingsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalFaults\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"FaultProven\",\"inputs\":[{\"name\":\"epoch\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"kind\",\"type\":\"uint8\",\"indexed\":true,\"internalType\":\"enumFaultKind\"},{\"name\":\"evidence\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"prover\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SignerActivated\",\"inputs\":[{\"name\":\"signer\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"fromEpoch\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SignerCancelled\",\"inputs\":[{\"name\":\"signer\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"fromEpoch\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SignerScheduled\",\"inputs\":[{\"name\":\"signer\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"fromEpoch\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"activatesAt\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AlreadyProven\",\"inputs\":[{\"name\":\"evidence\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"DifferentEpochs\",\"inputs\":[{\"name\":\"first\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"second\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"EpochNotAnchored\",\"inputs\":[{\"name\":\"epoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"FromEpochTooLow\",\"inputs\":[{\"name\":\"fromEpoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"least\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"NoConflict\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NoSignerChangePending\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotALedger\",\"inputs\":[{\"name\":\"ledger\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"NotSignedBySigner\",\"inputs\":[{\"name\":\"epoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"signer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"recovered\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ProofHolds\",\"inputs\":[{\"name\":\"epoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"SameSigner\",\"inputs\":[{\"name\":\"signer\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"SiblingsMismatch\",\"inputs\":[{\"name\":\"signed\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"given\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"SignerChangePending\",\"inputs\":[{\"name\":\"signer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"fromEpoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"SignerChangeStale\",\"inputs\":[{\"name\":\"fromEpoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"nextEpoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"SignerChangeTooEarly\",\"inputs\":[{\"name\":\"activatesAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"blockTime\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"WrongLedger\",\"inputs\":[{\"name\":\"expected\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"given\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ZeroDelay\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroSigner\",\"inputs\":[]}]",
	ID:  "IStatementCourt",
}

// IStatementCourt is an auto generated Go binding around an Ethereum contract.
type IStatementCourt struct {
	abi abi.ABI
}

// GetABI returns the ABI associated with this contract binding.
func (c *IStatementCourt) GetABI() abi.ABI {
	return c.abi
}

// NewIStatementCourt creates a new instance of IStatementCourt.
func NewIStatementCourt() *IStatementCourt {
	parsed, err := IStatementCourtMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &IStatementCourt{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *IStatementCourt) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackActivateSigner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x97326b4a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function activateSigner() returns()
func (iStatementCourt *IStatementCourt) PackActivateSigner() []byte {
	enc, err := iStatementCourt.abi.Pack("activateSigner")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackActivateSigner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x97326b4a.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function activateSigner() returns()
func (iStatementCourt *IStatementCourt) TryPackActivateSigner() ([]byte, error) {
	return iStatementCourt.abi.Pack("activateSigner")
}

// PackCancelSigner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x21152eb0.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function cancelSigner() returns()
func (iStatementCourt *IStatementCourt) PackCancelSigner() []byte {
	enc, err := iStatementCourt.abi.Pack("cancelSigner")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCancelSigner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x21152eb0.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function cancelSigner() returns()
func (iStatementCourt *IStatementCourt) TryPackCancelSigner() ([]byte, error) {
	return iStatementCourt.abi.Pack("cancelSigner")
}

// PackFaultCount is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x36772c82.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function faultCount(uint64 epoch) view returns(uint256)
func (iStatementCourt *IStatementCourt) PackFaultCount(epoch uint64) []byte {
	enc, err := iStatementCourt.abi.Pack("faultCount", epoch)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackFaultCount is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x36772c82.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function faultCount(uint64 epoch) view returns(uint256)
func (iStatementCourt *IStatementCourt) TryPackFaultCount(epoch uint64) ([]byte, error) {
	return iStatementCourt.abi.Pack("faultCount", epoch)
}

// UnpackFaultCount is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x36772c82.
//
// Solidity: function faultCount(uint64 epoch) view returns(uint256)
func (iStatementCourt *IStatementCourt) UnpackFaultCount(data []byte) (*big.Int, error) {
	out, err := iStatementCourt.abi.Unpack("faultCount", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackFaultsOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb745c87b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function faultsOf(uint64 epoch) view returns((uint8,bytes32,address,uint64)[])
func (iStatementCourt *IStatementCourt) PackFaultsOf(epoch uint64) []byte {
	enc, err := iStatementCourt.abi.Pack("faultsOf", epoch)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackFaultsOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb745c87b.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function faultsOf(uint64 epoch) view returns((uint8,bytes32,address,uint64)[])
func (iStatementCourt *IStatementCourt) TryPackFaultsOf(epoch uint64) ([]byte, error) {
	return iStatementCourt.abi.Pack("faultsOf", epoch)
}

// UnpackFaultsOf is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb745c87b.
//
// Solidity: function faultsOf(uint64 epoch) view returns((uint8,bytes32,address,uint64)[])
func (iStatementCourt *IStatementCourt) UnpackFaultsOf(data []byte) ([]Fault, error) {
	out, err := iStatementCourt.abi.Unpack("faultsOf", data)
	if err != nil {
		return *new([]Fault), err
	}
	out0 := *abi.ConvertType(out[0], new([]Fault)).(*[]Fault)
	return out0, nil
}

// PackIsProven is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa44256c6.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function isProven(bytes32 evidence) view returns(bool)
func (iStatementCourt *IStatementCourt) PackIsProven(evidence [32]byte) []byte {
	enc, err := iStatementCourt.abi.Pack("isProven", evidence)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackIsProven is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa44256c6.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function isProven(bytes32 evidence) view returns(bool)
func (iStatementCourt *IStatementCourt) TryPackIsProven(evidence [32]byte) ([]byte, error) {
	return iStatementCourt.abi.Pack("isProven", evidence)
}

// UnpackIsProven is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa44256c6.
//
// Solidity: function isProven(bytes32 evidence) view returns(bool)
func (iStatementCourt *IStatementCourt) UnpackIsProven(data []byte) (bool, error) {
	out, err := iStatementCourt.abi.Unpack("isProven", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackLedger is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x56397c35.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function ledger() view returns(address)
func (iStatementCourt *IStatementCourt) PackLedger() []byte {
	enc, err := iStatementCourt.abi.Pack("ledger")
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
func (iStatementCourt *IStatementCourt) TryPackLedger() ([]byte, error) {
	return iStatementCourt.abi.Pack("ledger")
}

// UnpackLedger is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x56397c35.
//
// Solidity: function ledger() view returns(address)
func (iStatementCourt *IStatementCourt) UnpackLedger(data []byte) (common.Address, error) {
	out, err := iStatementCourt.abi.Unpack("ledger", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackPendingSigner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x14d90c45.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function pendingSigner() view returns((address,uint64,uint64))
func (iStatementCourt *IStatementCourt) PackPendingSigner() []byte {
	enc, err := iStatementCourt.abi.Pack("pendingSigner")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPendingSigner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x14d90c45.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function pendingSigner() view returns((address,uint64,uint64))
func (iStatementCourt *IStatementCourt) TryPackPendingSigner() ([]byte, error) {
	return iStatementCourt.abi.Pack("pendingSigner")
}

// UnpackPendingSigner is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x14d90c45.
//
// Solidity: function pendingSigner() view returns((address,uint64,uint64))
func (iStatementCourt *IStatementCourt) UnpackPendingSigner(data []byte) (PendingSigner, error) {
	out, err := iStatementCourt.abi.Unpack("pendingSigner", data)
	if err != nil {
		return *new(PendingSigner), err
	}
	out0 := *abi.ConvertType(out[0], new(PendingSigner)).(*PendingSigner)
	return out0, nil
}

// PackProveBadProof is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xca9a360d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function proveBadProof((address,uint64,bytes32,uint64,uint64,bytes32) statement, bytes signature, (bytes32,uint64)[] siblings) returns(bytes32)
func (iStatementCourt *IStatementCourt) PackProveBadProof(statement EpochStatement, signature []byte, siblings []SumNode) []byte {
	enc, err := iStatementCourt.abi.Pack("proveBadProof", statement, signature, siblings)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackProveBadProof is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xca9a360d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function proveBadProof((address,uint64,bytes32,uint64,uint64,bytes32) statement, bytes signature, (bytes32,uint64)[] siblings) returns(bytes32)
func (iStatementCourt *IStatementCourt) TryPackProveBadProof(statement EpochStatement, signature []byte, siblings []SumNode) ([]byte, error) {
	return iStatementCourt.abi.Pack("proveBadProof", statement, signature, siblings)
}

// UnpackProveBadProof is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xca9a360d.
//
// Solidity: function proveBadProof((address,uint64,bytes32,uint64,uint64,bytes32) statement, bytes signature, (bytes32,uint64)[] siblings) returns(bytes32)
func (iStatementCourt *IStatementCourt) UnpackProveBadProof(data []byte) ([32]byte, error) {
	out, err := iStatementCourt.abi.Unpack("proveBadProof", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackProveEquivocation is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc1c23c65.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function proveEquivocation((address,uint64,bytes32,uint64,uint64,bytes32) first, bytes firstSignature, (address,uint64,bytes32,uint64,uint64,bytes32) second, bytes secondSignature) returns(bytes32)
func (iStatementCourt *IStatementCourt) PackProveEquivocation(first EpochStatement, firstSignature []byte, second EpochStatement, secondSignature []byte) []byte {
	enc, err := iStatementCourt.abi.Pack("proveEquivocation", first, firstSignature, second, secondSignature)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackProveEquivocation is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc1c23c65.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function proveEquivocation((address,uint64,bytes32,uint64,uint64,bytes32) first, bytes firstSignature, (address,uint64,bytes32,uint64,uint64,bytes32) second, bytes secondSignature) returns(bytes32)
func (iStatementCourt *IStatementCourt) TryPackProveEquivocation(first EpochStatement, firstSignature []byte, second EpochStatement, secondSignature []byte) ([]byte, error) {
	return iStatementCourt.abi.Pack("proveEquivocation", first, firstSignature, second, secondSignature)
}

// UnpackProveEquivocation is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc1c23c65.
//
// Solidity: function proveEquivocation((address,uint64,bytes32,uint64,uint64,bytes32) first, bytes firstSignature, (address,uint64,bytes32,uint64,uint64,bytes32) second, bytes secondSignature) returns(bytes32)
func (iStatementCourt *IStatementCourt) UnpackProveEquivocation(data []byte) ([32]byte, error) {
	out, err := iStatementCourt.abi.Unpack("proveEquivocation", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackRecoverSigner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x181a36a4.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function recoverSigner((address,uint64,bytes32,uint64,uint64,bytes32) statement, bytes signature) view returns(address)
func (iStatementCourt *IStatementCourt) PackRecoverSigner(statement EpochStatement, signature []byte) []byte {
	enc, err := iStatementCourt.abi.Pack("recoverSigner", statement, signature)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRecoverSigner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x181a36a4.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function recoverSigner((address,uint64,bytes32,uint64,uint64,bytes32) statement, bytes signature) view returns(address)
func (iStatementCourt *IStatementCourt) TryPackRecoverSigner(statement EpochStatement, signature []byte) ([]byte, error) {
	return iStatementCourt.abi.Pack("recoverSigner", statement, signature)
}

// UnpackRecoverSigner is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x181a36a4.
//
// Solidity: function recoverSigner((address,uint64,bytes32,uint64,uint64,bytes32) statement, bytes signature) view returns(address)
func (iStatementCourt *IStatementCourt) UnpackRecoverSigner(data []byte) (common.Address, error) {
	out, err := iStatementCourt.abi.Unpack("recoverSigner", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackScheduleSigner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5e8a5604.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function scheduleSigner(address signer, uint64 fromEpoch) returns()
func (iStatementCourt *IStatementCourt) PackScheduleSigner(signer common.Address, fromEpoch uint64) []byte {
	enc, err := iStatementCourt.abi.Pack("scheduleSigner", signer, fromEpoch)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackScheduleSigner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5e8a5604.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function scheduleSigner(address signer, uint64 fromEpoch) returns()
func (iStatementCourt *IStatementCourt) TryPackScheduleSigner(signer common.Address, fromEpoch uint64) ([]byte, error) {
	return iStatementCourt.abi.Pack("scheduleSigner", signer, fromEpoch)
}

// PackSiblingsHash is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x51ba300e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function siblingsHash((bytes32,uint64)[] siblings) pure returns(bytes32)
func (iStatementCourt *IStatementCourt) PackSiblingsHash(siblings []SumNode) []byte {
	enc, err := iStatementCourt.abi.Pack("siblingsHash", siblings)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSiblingsHash is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x51ba300e.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function siblingsHash((bytes32,uint64)[] siblings) pure returns(bytes32)
func (iStatementCourt *IStatementCourt) TryPackSiblingsHash(siblings []SumNode) ([]byte, error) {
	return iStatementCourt.abi.Pack("siblingsHash", siblings)
}

// UnpackSiblingsHash is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x51ba300e.
//
// Solidity: function siblingsHash((bytes32,uint64)[] siblings) pure returns(bytes32)
func (iStatementCourt *IStatementCourt) UnpackSiblingsHash(data []byte) ([32]byte, error) {
	out, err := iStatementCourt.abi.Unpack("siblingsHash", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackSignerChangeDelay is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x357c0fb0.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function signerChangeDelay() view returns(uint64)
func (iStatementCourt *IStatementCourt) PackSignerChangeDelay() []byte {
	enc, err := iStatementCourt.abi.Pack("signerChangeDelay")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSignerChangeDelay is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x357c0fb0.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function signerChangeDelay() view returns(uint64)
func (iStatementCourt *IStatementCourt) TryPackSignerChangeDelay() ([]byte, error) {
	return iStatementCourt.abi.Pack("signerChangeDelay")
}

// UnpackSignerChangeDelay is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x357c0fb0.
//
// Solidity: function signerChangeDelay() view returns(uint64)
func (iStatementCourt *IStatementCourt) UnpackSignerChangeDelay(data []byte) (uint64, error) {
	out, err := iStatementCourt.abi.Unpack("signerChangeDelay", data)
	if err != nil {
		return *new(uint64), err
	}
	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)
	return out0, nil
}

// PackSignerOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x22a3808f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function signerOf(uint64 epoch) view returns(address)
func (iStatementCourt *IStatementCourt) PackSignerOf(epoch uint64) []byte {
	enc, err := iStatementCourt.abi.Pack("signerOf", epoch)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSignerOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x22a3808f.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function signerOf(uint64 epoch) view returns(address)
func (iStatementCourt *IStatementCourt) TryPackSignerOf(epoch uint64) ([]byte, error) {
	return iStatementCourt.abi.Pack("signerOf", epoch)
}

// UnpackSignerOf is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x22a3808f.
//
// Solidity: function signerOf(uint64 epoch) view returns(address)
func (iStatementCourt *IStatementCourt) UnpackSignerOf(data []byte) (common.Address, error) {
	out, err := iStatementCourt.abi.Unpack("signerOf", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackSignerTerms is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x854209e8.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function signerTerms() view returns((address,uint64)[])
func (iStatementCourt *IStatementCourt) PackSignerTerms() []byte {
	enc, err := iStatementCourt.abi.Pack("signerTerms")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSignerTerms is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x854209e8.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function signerTerms() view returns((address,uint64)[])
func (iStatementCourt *IStatementCourt) TryPackSignerTerms() ([]byte, error) {
	return iStatementCourt.abi.Pack("signerTerms")
}

// UnpackSignerTerms is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x854209e8.
//
// Solidity: function signerTerms() view returns((address,uint64)[])
func (iStatementCourt *IStatementCourt) UnpackSignerTerms(data []byte) ([]SignerTerm, error) {
	out, err := iStatementCourt.abi.Unpack("signerTerms", data)
	if err != nil {
		return *new([]SignerTerm), err
	}
	out0 := *abi.ConvertType(out[0], new([]SignerTerm)).(*[]SignerTerm)
	return out0, nil
}

// PackStatementDigest is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x941c6d0c.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function statementDigest((address,uint64,bytes32,uint64,uint64,bytes32) statement) view returns(bytes32)
func (iStatementCourt *IStatementCourt) PackStatementDigest(statement EpochStatement) []byte {
	enc, err := iStatementCourt.abi.Pack("statementDigest", statement)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackStatementDigest is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x941c6d0c.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function statementDigest((address,uint64,bytes32,uint64,uint64,bytes32) statement) view returns(bytes32)
func (iStatementCourt *IStatementCourt) TryPackStatementDigest(statement EpochStatement) ([]byte, error) {
	return iStatementCourt.abi.Pack("statementDigest", statement)
}

// UnpackStatementDigest is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x941c6d0c.
//
// Solidity: function statementDigest((address,uint64,bytes32,uint64,uint64,bytes32) statement) view returns(bytes32)
func (iStatementCourt *IStatementCourt) UnpackStatementDigest(data []byte) ([32]byte, error) {
	out, err := iStatementCourt.abi.Unpack("statementDigest", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackTotalFaults is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa663d04a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function totalFaults() view returns(uint256)
func (iStatementCourt *IStatementCourt) PackTotalFaults() []byte {
	enc, err := iStatementCourt.abi.Pack("totalFaults")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTotalFaults is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa663d04a.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function totalFaults() view returns(uint256)
func (iStatementCourt *IStatementCourt) TryPackTotalFaults() ([]byte, error) {
	return iStatementCourt.abi.Pack("totalFaults")
}

// UnpackTotalFaults is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa663d04a.
//
// Solidity: function totalFaults() view returns(uint256)
func (iStatementCourt *IStatementCourt) UnpackTotalFaults(data []byte) (*big.Int, error) {
	out, err := iStatementCourt.abi.Unpack("totalFaults", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// IStatementCourtFaultProven represents a FaultProven event raised by the IStatementCourt contract.
type IStatementCourtFaultProven struct {
	Epoch    uint64
	Kind     uint8
	Evidence [32]byte
	Prover   common.Address
	Raw      *types.Log // Blockchain specific contextual infos
}

const IStatementCourtFaultProvenEventName = "FaultProven"

// ContractEventName returns the user-defined event name.
func (IStatementCourtFaultProven) ContractEventName() string {
	return IStatementCourtFaultProvenEventName
}

// UnpackFaultProvenEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event FaultProven(uint64 indexed epoch, uint8 indexed kind, bytes32 indexed evidence, address prover)
func (iStatementCourt *IStatementCourt) UnpackFaultProvenEvent(log *types.Log) (*IStatementCourtFaultProven, error) {
	event := "FaultProven"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != iStatementCourt.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(IStatementCourtFaultProven)
	if len(log.Data) > 0 {
		if err := iStatementCourt.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iStatementCourt.abi.Events[event].Inputs {
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

// IStatementCourtSignerActivated represents a SignerActivated event raised by the IStatementCourt contract.
type IStatementCourtSignerActivated struct {
	Signer    common.Address
	FromEpoch uint64
	Raw       *types.Log // Blockchain specific contextual infos
}

const IStatementCourtSignerActivatedEventName = "SignerActivated"

// ContractEventName returns the user-defined event name.
func (IStatementCourtSignerActivated) ContractEventName() string {
	return IStatementCourtSignerActivatedEventName
}

// UnpackSignerActivatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event SignerActivated(address indexed signer, uint64 fromEpoch)
func (iStatementCourt *IStatementCourt) UnpackSignerActivatedEvent(log *types.Log) (*IStatementCourtSignerActivated, error) {
	event := "SignerActivated"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != iStatementCourt.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(IStatementCourtSignerActivated)
	if len(log.Data) > 0 {
		if err := iStatementCourt.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iStatementCourt.abi.Events[event].Inputs {
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

// IStatementCourtSignerCancelled represents a SignerCancelled event raised by the IStatementCourt contract.
type IStatementCourtSignerCancelled struct {
	Signer    common.Address
	FromEpoch uint64
	Raw       *types.Log // Blockchain specific contextual infos
}

const IStatementCourtSignerCancelledEventName = "SignerCancelled"

// ContractEventName returns the user-defined event name.
func (IStatementCourtSignerCancelled) ContractEventName() string {
	return IStatementCourtSignerCancelledEventName
}

// UnpackSignerCancelledEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event SignerCancelled(address indexed signer, uint64 fromEpoch)
func (iStatementCourt *IStatementCourt) UnpackSignerCancelledEvent(log *types.Log) (*IStatementCourtSignerCancelled, error) {
	event := "SignerCancelled"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != iStatementCourt.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(IStatementCourtSignerCancelled)
	if len(log.Data) > 0 {
		if err := iStatementCourt.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iStatementCourt.abi.Events[event].Inputs {
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

// IStatementCourtSignerScheduled represents a SignerScheduled event raised by the IStatementCourt contract.
type IStatementCourtSignerScheduled struct {
	Signer      common.Address
	FromEpoch   uint64
	ActivatesAt uint64
	Raw         *types.Log // Blockchain specific contextual infos
}

const IStatementCourtSignerScheduledEventName = "SignerScheduled"

// ContractEventName returns the user-defined event name.
func (IStatementCourtSignerScheduled) ContractEventName() string {
	return IStatementCourtSignerScheduledEventName
}

// UnpackSignerScheduledEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event SignerScheduled(address indexed signer, uint64 fromEpoch, uint64 activatesAt)
func (iStatementCourt *IStatementCourt) UnpackSignerScheduledEvent(log *types.Log) (*IStatementCourtSignerScheduled, error) {
	event := "SignerScheduled"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != iStatementCourt.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(IStatementCourtSignerScheduled)
	if len(log.Data) > 0 {
		if err := iStatementCourt.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iStatementCourt.abi.Events[event].Inputs {
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
func (iStatementCourt *IStatementCourt) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], iStatementCourt.abi.Errors["AlreadyProven"].ID.Bytes()[:4]) {
		return iStatementCourt.UnpackAlreadyProvenError(raw[4:])
	}
	if bytes.Equal(raw[:4], iStatementCourt.abi.Errors["DifferentEpochs"].ID.Bytes()[:4]) {
		return iStatementCourt.UnpackDifferentEpochsError(raw[4:])
	}
	if bytes.Equal(raw[:4], iStatementCourt.abi.Errors["EpochNotAnchored"].ID.Bytes()[:4]) {
		return iStatementCourt.UnpackEpochNotAnchoredError(raw[4:])
	}
	if bytes.Equal(raw[:4], iStatementCourt.abi.Errors["FromEpochTooLow"].ID.Bytes()[:4]) {
		return iStatementCourt.UnpackFromEpochTooLowError(raw[4:])
	}
	if bytes.Equal(raw[:4], iStatementCourt.abi.Errors["NoConflict"].ID.Bytes()[:4]) {
		return iStatementCourt.UnpackNoConflictError(raw[4:])
	}
	if bytes.Equal(raw[:4], iStatementCourt.abi.Errors["NoSignerChangePending"].ID.Bytes()[:4]) {
		return iStatementCourt.UnpackNoSignerChangePendingError(raw[4:])
	}
	if bytes.Equal(raw[:4], iStatementCourt.abi.Errors["NotALedger"].ID.Bytes()[:4]) {
		return iStatementCourt.UnpackNotALedgerError(raw[4:])
	}
	if bytes.Equal(raw[:4], iStatementCourt.abi.Errors["NotSignedBySigner"].ID.Bytes()[:4]) {
		return iStatementCourt.UnpackNotSignedBySignerError(raw[4:])
	}
	if bytes.Equal(raw[:4], iStatementCourt.abi.Errors["ProofHolds"].ID.Bytes()[:4]) {
		return iStatementCourt.UnpackProofHoldsError(raw[4:])
	}
	if bytes.Equal(raw[:4], iStatementCourt.abi.Errors["SameSigner"].ID.Bytes()[:4]) {
		return iStatementCourt.UnpackSameSignerError(raw[4:])
	}
	if bytes.Equal(raw[:4], iStatementCourt.abi.Errors["SiblingsMismatch"].ID.Bytes()[:4]) {
		return iStatementCourt.UnpackSiblingsMismatchError(raw[4:])
	}
	if bytes.Equal(raw[:4], iStatementCourt.abi.Errors["SignerChangePending"].ID.Bytes()[:4]) {
		return iStatementCourt.UnpackSignerChangePendingError(raw[4:])
	}
	if bytes.Equal(raw[:4], iStatementCourt.abi.Errors["SignerChangeStale"].ID.Bytes()[:4]) {
		return iStatementCourt.UnpackSignerChangeStaleError(raw[4:])
	}
	if bytes.Equal(raw[:4], iStatementCourt.abi.Errors["SignerChangeTooEarly"].ID.Bytes()[:4]) {
		return iStatementCourt.UnpackSignerChangeTooEarlyError(raw[4:])
	}
	if bytes.Equal(raw[:4], iStatementCourt.abi.Errors["WrongLedger"].ID.Bytes()[:4]) {
		return iStatementCourt.UnpackWrongLedgerError(raw[4:])
	}
	if bytes.Equal(raw[:4], iStatementCourt.abi.Errors["ZeroDelay"].ID.Bytes()[:4]) {
		return iStatementCourt.UnpackZeroDelayError(raw[4:])
	}
	if bytes.Equal(raw[:4], iStatementCourt.abi.Errors["ZeroSigner"].ID.Bytes()[:4]) {
		return iStatementCourt.UnpackZeroSignerError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// IStatementCourtAlreadyProven represents a AlreadyProven error raised by the IStatementCourt contract.
type IStatementCourtAlreadyProven struct {
	Evidence [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AlreadyProven(bytes32 evidence)
func IStatementCourtAlreadyProvenErrorID() common.Hash {
	return common.HexToHash("0xea033684b2eb17b3fc70d5009bef3efe803909cba52e3380dd0e081697f487ad")
}

// UnpackAlreadyProvenError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AlreadyProven(bytes32 evidence)
func (iStatementCourt *IStatementCourt) UnpackAlreadyProvenError(raw []byte) (*IStatementCourtAlreadyProven, error) {
	out := new(IStatementCourtAlreadyProven)
	if err := iStatementCourt.abi.UnpackIntoInterface(out, "AlreadyProven", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IStatementCourtDifferentEpochs represents a DifferentEpochs error raised by the IStatementCourt contract.
type IStatementCourtDifferentEpochs struct {
	First  uint64
	Second uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error DifferentEpochs(uint64 first, uint64 second)
func IStatementCourtDifferentEpochsErrorID() common.Hash {
	return common.HexToHash("0xbfb240a28f6d744e2666b046f0954afb0cb85dfbc6d6bd027717efc63eed50f4")
}

// UnpackDifferentEpochsError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error DifferentEpochs(uint64 first, uint64 second)
func (iStatementCourt *IStatementCourt) UnpackDifferentEpochsError(raw []byte) (*IStatementCourtDifferentEpochs, error) {
	out := new(IStatementCourtDifferentEpochs)
	if err := iStatementCourt.abi.UnpackIntoInterface(out, "DifferentEpochs", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IStatementCourtEpochNotAnchored represents a EpochNotAnchored error raised by the IStatementCourt contract.
type IStatementCourtEpochNotAnchored struct {
	Epoch uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error EpochNotAnchored(uint64 epoch)
func IStatementCourtEpochNotAnchoredErrorID() common.Hash {
	return common.HexToHash("0xfdd66f09e3e0335148d05ffa812e23b41ab1c93ea73617ca6ddd26710f05879d")
}

// UnpackEpochNotAnchoredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error EpochNotAnchored(uint64 epoch)
func (iStatementCourt *IStatementCourt) UnpackEpochNotAnchoredError(raw []byte) (*IStatementCourtEpochNotAnchored, error) {
	out := new(IStatementCourtEpochNotAnchored)
	if err := iStatementCourt.abi.UnpackIntoInterface(out, "EpochNotAnchored", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IStatementCourtFromEpochTooLow represents a FromEpochTooLow error raised by the IStatementCourt contract.
type IStatementCourtFromEpochTooLow struct {
	FromEpoch uint64
	Least     uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error FromEpochTooLow(uint64 fromEpoch, uint64 least)
func IStatementCourtFromEpochTooLowErrorID() common.Hash {
	return common.HexToHash("0x2a9c655bfc90a6b2ff4c7c129889535ba67dd8cdd7b2e2387b0e2e444a2035fd")
}

// UnpackFromEpochTooLowError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error FromEpochTooLow(uint64 fromEpoch, uint64 least)
func (iStatementCourt *IStatementCourt) UnpackFromEpochTooLowError(raw []byte) (*IStatementCourtFromEpochTooLow, error) {
	out := new(IStatementCourtFromEpochTooLow)
	if err := iStatementCourt.abi.UnpackIntoInterface(out, "FromEpochTooLow", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IStatementCourtNoConflict represents a NoConflict error raised by the IStatementCourt contract.
type IStatementCourtNoConflict struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NoConflict()
func IStatementCourtNoConflictErrorID() common.Hash {
	return common.HexToHash("0x7412124e23106f671f0c15e99815ae376e602111b93496444e63ed6071334091")
}

// UnpackNoConflictError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NoConflict()
func (iStatementCourt *IStatementCourt) UnpackNoConflictError(raw []byte) (*IStatementCourtNoConflict, error) {
	out := new(IStatementCourtNoConflict)
	if err := iStatementCourt.abi.UnpackIntoInterface(out, "NoConflict", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IStatementCourtNoSignerChangePending represents a NoSignerChangePending error raised by the IStatementCourt contract.
type IStatementCourtNoSignerChangePending struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NoSignerChangePending()
func IStatementCourtNoSignerChangePendingErrorID() common.Hash {
	return common.HexToHash("0x16d48c849add6c39b79e80fa155b4bf87a0a29d7ad1b668b5e80f101e66238f8")
}

// UnpackNoSignerChangePendingError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NoSignerChangePending()
func (iStatementCourt *IStatementCourt) UnpackNoSignerChangePendingError(raw []byte) (*IStatementCourtNoSignerChangePending, error) {
	out := new(IStatementCourtNoSignerChangePending)
	if err := iStatementCourt.abi.UnpackIntoInterface(out, "NoSignerChangePending", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IStatementCourtNotALedger represents a NotALedger error raised by the IStatementCourt contract.
type IStatementCourtNotALedger struct {
	Ledger common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotALedger(address ledger)
func IStatementCourtNotALedgerErrorID() common.Hash {
	return common.HexToHash("0x44a5b4b2f2697952fce85b89e8be359271d0792b4f0be489d664b4771187455b")
}

// UnpackNotALedgerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotALedger(address ledger)
func (iStatementCourt *IStatementCourt) UnpackNotALedgerError(raw []byte) (*IStatementCourtNotALedger, error) {
	out := new(IStatementCourtNotALedger)
	if err := iStatementCourt.abi.UnpackIntoInterface(out, "NotALedger", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IStatementCourtNotSignedBySigner represents a NotSignedBySigner error raised by the IStatementCourt contract.
type IStatementCourtNotSignedBySigner struct {
	Epoch     uint64
	Signer    common.Address
	Recovered common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotSignedBySigner(uint64 epoch, address signer, address recovered)
func IStatementCourtNotSignedBySignerErrorID() common.Hash {
	return common.HexToHash("0xd7eecde17baa6c3c5da78ec7d0f325a551b50706d7779d410753e9ea173f5813")
}

// UnpackNotSignedBySignerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotSignedBySigner(uint64 epoch, address signer, address recovered)
func (iStatementCourt *IStatementCourt) UnpackNotSignedBySignerError(raw []byte) (*IStatementCourtNotSignedBySigner, error) {
	out := new(IStatementCourtNotSignedBySigner)
	if err := iStatementCourt.abi.UnpackIntoInterface(out, "NotSignedBySigner", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IStatementCourtProofHolds represents a ProofHolds error raised by the IStatementCourt contract.
type IStatementCourtProofHolds struct {
	Epoch uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ProofHolds(uint64 epoch)
func IStatementCourtProofHoldsErrorID() common.Hash {
	return common.HexToHash("0x64212abdba47381a8e97a2ef1859b440605c270d7fc490e549b78acd72e28e3d")
}

// UnpackProofHoldsError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ProofHolds(uint64 epoch)
func (iStatementCourt *IStatementCourt) UnpackProofHoldsError(raw []byte) (*IStatementCourtProofHolds, error) {
	out := new(IStatementCourtProofHolds)
	if err := iStatementCourt.abi.UnpackIntoInterface(out, "ProofHolds", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IStatementCourtSameSigner represents a SameSigner error raised by the IStatementCourt contract.
type IStatementCourtSameSigner struct {
	Signer common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SameSigner(address signer)
func IStatementCourtSameSignerErrorID() common.Hash {
	return common.HexToHash("0x684a329eced21ee6c2d72bdaef91eb1248ae7b6be68e0850361db232a4782ee7")
}

// UnpackSameSignerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SameSigner(address signer)
func (iStatementCourt *IStatementCourt) UnpackSameSignerError(raw []byte) (*IStatementCourtSameSigner, error) {
	out := new(IStatementCourtSameSigner)
	if err := iStatementCourt.abi.UnpackIntoInterface(out, "SameSigner", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IStatementCourtSiblingsMismatch represents a SiblingsMismatch error raised by the IStatementCourt contract.
type IStatementCourtSiblingsMismatch struct {
	Signed [32]byte
	Given  [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SiblingsMismatch(bytes32 signed, bytes32 given)
func IStatementCourtSiblingsMismatchErrorID() common.Hash {
	return common.HexToHash("0x7e68dade1f921625838c9a65f0505353f5cae857ab8fac0fc3a0d07f29ef5fe1")
}

// UnpackSiblingsMismatchError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SiblingsMismatch(bytes32 signed, bytes32 given)
func (iStatementCourt *IStatementCourt) UnpackSiblingsMismatchError(raw []byte) (*IStatementCourtSiblingsMismatch, error) {
	out := new(IStatementCourtSiblingsMismatch)
	if err := iStatementCourt.abi.UnpackIntoInterface(out, "SiblingsMismatch", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IStatementCourtSignerChangePending represents a SignerChangePending error raised by the IStatementCourt contract.
type IStatementCourtSignerChangePending struct {
	Signer    common.Address
	FromEpoch uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SignerChangePending(address signer, uint64 fromEpoch)
func IStatementCourtSignerChangePendingErrorID() common.Hash {
	return common.HexToHash("0xbe977730e37b951325d90256c2cad47e6616e166013e1587bc36024d82a30178")
}

// UnpackSignerChangePendingError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SignerChangePending(address signer, uint64 fromEpoch)
func (iStatementCourt *IStatementCourt) UnpackSignerChangePendingError(raw []byte) (*IStatementCourtSignerChangePending, error) {
	out := new(IStatementCourtSignerChangePending)
	if err := iStatementCourt.abi.UnpackIntoInterface(out, "SignerChangePending", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IStatementCourtSignerChangeStale represents a SignerChangeStale error raised by the IStatementCourt contract.
type IStatementCourtSignerChangeStale struct {
	FromEpoch uint64
	NextEpoch uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SignerChangeStale(uint64 fromEpoch, uint64 nextEpoch)
func IStatementCourtSignerChangeStaleErrorID() common.Hash {
	return common.HexToHash("0x6d5c81aaef3921214d8a666dce0b18be17a1443922f6414672a46c7cfcff9770")
}

// UnpackSignerChangeStaleError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SignerChangeStale(uint64 fromEpoch, uint64 nextEpoch)
func (iStatementCourt *IStatementCourt) UnpackSignerChangeStaleError(raw []byte) (*IStatementCourtSignerChangeStale, error) {
	out := new(IStatementCourtSignerChangeStale)
	if err := iStatementCourt.abi.UnpackIntoInterface(out, "SignerChangeStale", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IStatementCourtSignerChangeTooEarly represents a SignerChangeTooEarly error raised by the IStatementCourt contract.
type IStatementCourtSignerChangeTooEarly struct {
	ActivatesAt uint64
	BlockTime   uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SignerChangeTooEarly(uint64 activatesAt, uint64 blockTime)
func IStatementCourtSignerChangeTooEarlyErrorID() common.Hash {
	return common.HexToHash("0x2a2d47109124377d337480ded5b5da8a818f42b62f77ba12bace133b3888c5ed")
}

// UnpackSignerChangeTooEarlyError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SignerChangeTooEarly(uint64 activatesAt, uint64 blockTime)
func (iStatementCourt *IStatementCourt) UnpackSignerChangeTooEarlyError(raw []byte) (*IStatementCourtSignerChangeTooEarly, error) {
	out := new(IStatementCourtSignerChangeTooEarly)
	if err := iStatementCourt.abi.UnpackIntoInterface(out, "SignerChangeTooEarly", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IStatementCourtWrongLedger represents a WrongLedger error raised by the IStatementCourt contract.
type IStatementCourtWrongLedger struct {
	Expected common.Address
	Given    common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error WrongLedger(address expected, address given)
func IStatementCourtWrongLedgerErrorID() common.Hash {
	return common.HexToHash("0x185f957c1ab2bcd607a1ccf4a1aa06b24f2dc77f61269eb02e5cc03658428df6")
}

// UnpackWrongLedgerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error WrongLedger(address expected, address given)
func (iStatementCourt *IStatementCourt) UnpackWrongLedgerError(raw []byte) (*IStatementCourtWrongLedger, error) {
	out := new(IStatementCourtWrongLedger)
	if err := iStatementCourt.abi.UnpackIntoInterface(out, "WrongLedger", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IStatementCourtZeroDelay represents a ZeroDelay error raised by the IStatementCourt contract.
type IStatementCourtZeroDelay struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ZeroDelay()
func IStatementCourtZeroDelayErrorID() common.Hash {
	return common.HexToHash("0x34c2b9dab051f3b49e0b9121d90ae110f06df318864d2741799dc06f9b92cebc")
}

// UnpackZeroDelayError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ZeroDelay()
func (iStatementCourt *IStatementCourt) UnpackZeroDelayError(raw []byte) (*IStatementCourtZeroDelay, error) {
	out := new(IStatementCourtZeroDelay)
	if err := iStatementCourt.abi.UnpackIntoInterface(out, "ZeroDelay", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IStatementCourtZeroSigner represents a ZeroSigner error raised by the IStatementCourt contract.
type IStatementCourtZeroSigner struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ZeroSigner()
func IStatementCourtZeroSignerErrorID() common.Hash {
	return common.HexToHash("0xe5c48ac578e050f5a863729c359600666b17bce09ec64f88719c5ce6fc799f6a")
}

// UnpackZeroSignerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ZeroSigner()
func (iStatementCourt *IStatementCourt) UnpackZeroSignerError(raw []byte) (*IStatementCourtZeroSigner, error) {
	out := new(IStatementCourtZeroSigner)
	if err := iStatementCourt.abi.UnpackIntoInterface(out, "ZeroSigner", raw); err != nil {
		return nil, err
	}
	return out, nil
}
