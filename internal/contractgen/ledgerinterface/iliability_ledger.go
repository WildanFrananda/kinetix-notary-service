// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ledgerinterface

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

// Epoch is an auto generated low-level Go binding around an user-defined struct.
type Epoch struct {
	Root       [32]byte
	Digest     [32]byte
	Total      uint64
	LeafCount  uint64
	AsOf       uint64
	AnchoredAt uint64
}

// EpochHeader is an auto generated low-level Go binding around an user-defined struct.
type EpochHeader struct {
	Number         uint64
	PreviousDigest [32]byte
	Root           [32]byte
	Total          uint64
	LeafCount      uint64
	AsOf           uint64
}

// SumNode is an auto generated low-level Go binding around an user-defined struct.
type SumNode struct {
	Hash [32]byte
	Sum  uint64
}

// ILiabilityLedgerMetaData contains all meta data concerning the ILiabilityLedger contract.
var ILiabilityLedgerMetaData = bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"anchor\",\"inputs\":[{\"name\":\"header\",\"type\":\"tuple\",\"internalType\":\"structEpochHeader\",\"components\":[{\"name\":\"number\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"previousDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"total\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"leafCount\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"asOf\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"epochAt\",\"inputs\":[{\"name\":\"number\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structEpoch\",\"components\":[{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"digest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"total\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"leafCount\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"asOf\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"anchoredAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"genesisDigest\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"genesisEpoch\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"head\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isSealed\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"lastAsOf\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"nextEpoch\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"scheme\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"seal\",\"inputs\":[{\"name\":\"expectedHead\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"verifyInclusion\",\"inputs\":[{\"name\":\"number\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"commitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"amount\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"index\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"siblings\",\"type\":\"tuple[]\",\"internalType\":\"structSumNode[]\",\"components\":[{\"name\":\"hash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sum\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"EpochAnchored\",\"inputs\":[{\"name\":\"number\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"digest\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"previousDigest\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"root\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"total\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"leafCount\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"asOf\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"LedgerSealed\",\"inputs\":[{\"name\":\"head\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"nextEpoch\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"lastAsOf\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"EmptyEpochMismatch\",\"inputs\":[{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"total\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"EpochAlreadyAnchored\",\"inputs\":[{\"name\":\"number\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"digest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"EpochBeforeGenesis\",\"inputs\":[{\"name\":\"genesisEpoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"number\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"EpochConflict\",\"inputs\":[{\"name\":\"number\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"anchored\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"given\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"EpochOutOfSequence\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"given\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"HeadMismatch\",\"inputs\":[{\"name\":\"head\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"InvalidGenesis\",\"inputs\":[{\"name\":\"genesisEpoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"genesisDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"LedgerIsSealed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PreviousDigestMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"given\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"SnapshotInFuture\",\"inputs\":[{\"name\":\"asOf\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"blockTime\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"SnapshotNotAfterPrevious\",\"inputs\":[{\"name\":\"asOf\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"previous\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"TooManyLeaves\",\"inputs\":[{\"name\":\"leafCount\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"maximum\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"UnknownEpoch\",\"inputs\":[{\"name\":\"number\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"ZeroAnchorer\",\"inputs\":[]}]",
	ID:  "ILiabilityLedger",
}

// ILiabilityLedger is an auto generated Go binding around an Ethereum contract.
type ILiabilityLedger struct {
	abi abi.ABI
}

// GetABI returns the ABI associated with this contract binding.
func (c *ILiabilityLedger) GetABI() abi.ABI {
	return c.abi
}

// NewILiabilityLedger creates a new instance of ILiabilityLedger.
func NewILiabilityLedger() *ILiabilityLedger {
	parsed, err := ILiabilityLedgerMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &ILiabilityLedger{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *ILiabilityLedger) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackAnchor is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1ffcf5d4.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function anchor((uint64,bytes32,bytes32,uint64,uint64,uint64) header) returns(bytes32)
func (iLiabilityLedger *ILiabilityLedger) PackAnchor(header EpochHeader) []byte {
	enc, err := iLiabilityLedger.abi.Pack("anchor", header)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAnchor is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1ffcf5d4.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function anchor((uint64,bytes32,bytes32,uint64,uint64,uint64) header) returns(bytes32)
func (iLiabilityLedger *ILiabilityLedger) TryPackAnchor(header EpochHeader) ([]byte, error) {
	return iLiabilityLedger.abi.Pack("anchor", header)
}

// UnpackAnchor is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x1ffcf5d4.
//
// Solidity: function anchor((uint64,bytes32,bytes32,uint64,uint64,uint64) header) returns(bytes32)
func (iLiabilityLedger *ILiabilityLedger) UnpackAnchor(data []byte) ([32]byte, error) {
	out, err := iLiabilityLedger.abi.Unpack("anchor", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackEpochAt is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7d60d355.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function epochAt(uint64 number) view returns((bytes32,bytes32,uint64,uint64,uint64,uint64))
func (iLiabilityLedger *ILiabilityLedger) PackEpochAt(number uint64) []byte {
	enc, err := iLiabilityLedger.abi.Pack("epochAt", number)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackEpochAt is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7d60d355.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function epochAt(uint64 number) view returns((bytes32,bytes32,uint64,uint64,uint64,uint64))
func (iLiabilityLedger *ILiabilityLedger) TryPackEpochAt(number uint64) ([]byte, error) {
	return iLiabilityLedger.abi.Pack("epochAt", number)
}

// UnpackEpochAt is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x7d60d355.
//
// Solidity: function epochAt(uint64 number) view returns((bytes32,bytes32,uint64,uint64,uint64,uint64))
func (iLiabilityLedger *ILiabilityLedger) UnpackEpochAt(data []byte) (Epoch, error) {
	out, err := iLiabilityLedger.abi.Unpack("epochAt", data)
	if err != nil {
		return *new(Epoch), err
	}
	out0 := *abi.ConvertType(out[0], new(Epoch)).(*Epoch)
	return out0, nil
}

// PackGenesisDigest is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x79ac2661.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function genesisDigest() view returns(bytes32)
func (iLiabilityLedger *ILiabilityLedger) PackGenesisDigest() []byte {
	enc, err := iLiabilityLedger.abi.Pack("genesisDigest")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGenesisDigest is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x79ac2661.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function genesisDigest() view returns(bytes32)
func (iLiabilityLedger *ILiabilityLedger) TryPackGenesisDigest() ([]byte, error) {
	return iLiabilityLedger.abi.Pack("genesisDigest")
}

// UnpackGenesisDigest is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x79ac2661.
//
// Solidity: function genesisDigest() view returns(bytes32)
func (iLiabilityLedger *ILiabilityLedger) UnpackGenesisDigest(data []byte) ([32]byte, error) {
	out, err := iLiabilityLedger.abi.Unpack("genesisDigest", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackGenesisEpoch is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb70e6be6.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function genesisEpoch() view returns(uint64)
func (iLiabilityLedger *ILiabilityLedger) PackGenesisEpoch() []byte {
	enc, err := iLiabilityLedger.abi.Pack("genesisEpoch")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGenesisEpoch is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb70e6be6.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function genesisEpoch() view returns(uint64)
func (iLiabilityLedger *ILiabilityLedger) TryPackGenesisEpoch() ([]byte, error) {
	return iLiabilityLedger.abi.Pack("genesisEpoch")
}

// UnpackGenesisEpoch is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb70e6be6.
//
// Solidity: function genesisEpoch() view returns(uint64)
func (iLiabilityLedger *ILiabilityLedger) UnpackGenesisEpoch(data []byte) (uint64, error) {
	out, err := iLiabilityLedger.abi.Unpack("genesisEpoch", data)
	if err != nil {
		return *new(uint64), err
	}
	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)
	return out0, nil
}

// PackHead is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8f7dcfa3.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function head() view returns(bytes32)
func (iLiabilityLedger *ILiabilityLedger) PackHead() []byte {
	enc, err := iLiabilityLedger.abi.Pack("head")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackHead is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8f7dcfa3.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function head() view returns(bytes32)
func (iLiabilityLedger *ILiabilityLedger) TryPackHead() ([]byte, error) {
	return iLiabilityLedger.abi.Pack("head")
}

// UnpackHead is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8f7dcfa3.
//
// Solidity: function head() view returns(bytes32)
func (iLiabilityLedger *ILiabilityLedger) UnpackHead(data []byte) ([32]byte, error) {
	out, err := iLiabilityLedger.abi.Unpack("head", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackIsSealed is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x631f9852.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function isSealed() view returns(bool)
func (iLiabilityLedger *ILiabilityLedger) PackIsSealed() []byte {
	enc, err := iLiabilityLedger.abi.Pack("isSealed")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackIsSealed is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x631f9852.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function isSealed() view returns(bool)
func (iLiabilityLedger *ILiabilityLedger) TryPackIsSealed() ([]byte, error) {
	return iLiabilityLedger.abi.Pack("isSealed")
}

// UnpackIsSealed is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x631f9852.
//
// Solidity: function isSealed() view returns(bool)
func (iLiabilityLedger *ILiabilityLedger) UnpackIsSealed(data []byte) (bool, error) {
	out, err := iLiabilityLedger.abi.Unpack("isSealed", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackLastAsOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa89210cb.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function lastAsOf() view returns(uint64)
func (iLiabilityLedger *ILiabilityLedger) PackLastAsOf() []byte {
	enc, err := iLiabilityLedger.abi.Pack("lastAsOf")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackLastAsOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa89210cb.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function lastAsOf() view returns(uint64)
func (iLiabilityLedger *ILiabilityLedger) TryPackLastAsOf() ([]byte, error) {
	return iLiabilityLedger.abi.Pack("lastAsOf")
}

// UnpackLastAsOf is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa89210cb.
//
// Solidity: function lastAsOf() view returns(uint64)
func (iLiabilityLedger *ILiabilityLedger) UnpackLastAsOf(data []byte) (uint64, error) {
	out, err := iLiabilityLedger.abi.Unpack("lastAsOf", data)
	if err != nil {
		return *new(uint64), err
	}
	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)
	return out0, nil
}

// PackNextEpoch is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xaea0e78b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function nextEpoch() view returns(uint64)
func (iLiabilityLedger *ILiabilityLedger) PackNextEpoch() []byte {
	enc, err := iLiabilityLedger.abi.Pack("nextEpoch")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackNextEpoch is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xaea0e78b.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function nextEpoch() view returns(uint64)
func (iLiabilityLedger *ILiabilityLedger) TryPackNextEpoch() ([]byte, error) {
	return iLiabilityLedger.abi.Pack("nextEpoch")
}

// UnpackNextEpoch is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xaea0e78b.
//
// Solidity: function nextEpoch() view returns(uint64)
func (iLiabilityLedger *ILiabilityLedger) UnpackNextEpoch(data []byte) (uint64, error) {
	out, err := iLiabilityLedger.abi.Unpack("nextEpoch", data)
	if err != nil {
		return *new(uint64), err
	}
	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)
	return out0, nil
}

// PackScheme is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8518b0ad.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function scheme() pure returns(uint16)
func (iLiabilityLedger *ILiabilityLedger) PackScheme() []byte {
	enc, err := iLiabilityLedger.abi.Pack("scheme")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackScheme is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8518b0ad.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function scheme() pure returns(uint16)
func (iLiabilityLedger *ILiabilityLedger) TryPackScheme() ([]byte, error) {
	return iLiabilityLedger.abi.Pack("scheme")
}

// UnpackScheme is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8518b0ad.
//
// Solidity: function scheme() pure returns(uint16)
func (iLiabilityLedger *ILiabilityLedger) UnpackScheme(data []byte) (uint16, error) {
	out, err := iLiabilityLedger.abi.Unpack("scheme", data)
	if err != nil {
		return *new(uint16), err
	}
	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)
	return out0, nil
}

// PackSeal is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb07eeda8.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function seal(bytes32 expectedHead) returns()
func (iLiabilityLedger *ILiabilityLedger) PackSeal(expectedHead [32]byte) []byte {
	enc, err := iLiabilityLedger.abi.Pack("seal", expectedHead)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSeal is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb07eeda8.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function seal(bytes32 expectedHead) returns()
func (iLiabilityLedger *ILiabilityLedger) TryPackSeal(expectedHead [32]byte) ([]byte, error) {
	return iLiabilityLedger.abi.Pack("seal", expectedHead)
}

// PackVerifyInclusion is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc0566abe.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function verifyInclusion(uint64 number, bytes32 commitment, uint64 amount, uint64 index, (bytes32,uint64)[] siblings) view returns(bool)
func (iLiabilityLedger *ILiabilityLedger) PackVerifyInclusion(number uint64, commitment [32]byte, amount uint64, index uint64, siblings []SumNode) []byte {
	enc, err := iLiabilityLedger.abi.Pack("verifyInclusion", number, commitment, amount, index, siblings)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackVerifyInclusion is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc0566abe.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function verifyInclusion(uint64 number, bytes32 commitment, uint64 amount, uint64 index, (bytes32,uint64)[] siblings) view returns(bool)
func (iLiabilityLedger *ILiabilityLedger) TryPackVerifyInclusion(number uint64, commitment [32]byte, amount uint64, index uint64, siblings []SumNode) ([]byte, error) {
	return iLiabilityLedger.abi.Pack("verifyInclusion", number, commitment, amount, index, siblings)
}

// UnpackVerifyInclusion is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc0566abe.
//
// Solidity: function verifyInclusion(uint64 number, bytes32 commitment, uint64 amount, uint64 index, (bytes32,uint64)[] siblings) view returns(bool)
func (iLiabilityLedger *ILiabilityLedger) UnpackVerifyInclusion(data []byte) (bool, error) {
	out, err := iLiabilityLedger.abi.Unpack("verifyInclusion", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// ILiabilityLedgerEpochAnchored represents a EpochAnchored event raised by the ILiabilityLedger contract.
type ILiabilityLedgerEpochAnchored struct {
	Number         uint64
	Digest         [32]byte
	PreviousDigest [32]byte
	Root           [32]byte
	Total          uint64
	LeafCount      uint64
	AsOf           uint64
	Raw            *types.Log // Blockchain specific contextual infos
}

const ILiabilityLedgerEpochAnchoredEventName = "EpochAnchored"

// ContractEventName returns the user-defined event name.
func (ILiabilityLedgerEpochAnchored) ContractEventName() string {
	return ILiabilityLedgerEpochAnchoredEventName
}

// UnpackEpochAnchoredEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event EpochAnchored(uint64 indexed number, bytes32 indexed digest, bytes32 previousDigest, bytes32 root, uint64 total, uint64 leafCount, uint64 asOf)
func (iLiabilityLedger *ILiabilityLedger) UnpackEpochAnchoredEvent(log *types.Log) (*ILiabilityLedgerEpochAnchored, error) {
	event := "EpochAnchored"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != iLiabilityLedger.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(ILiabilityLedgerEpochAnchored)
	if len(log.Data) > 0 {
		if err := iLiabilityLedger.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iLiabilityLedger.abi.Events[event].Inputs {
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

// ILiabilityLedgerLedgerSealed represents a LedgerSealed event raised by the ILiabilityLedger contract.
type ILiabilityLedgerLedgerSealed struct {
	Head      [32]byte
	NextEpoch uint64
	LastAsOf  uint64
	Raw       *types.Log // Blockchain specific contextual infos
}

const ILiabilityLedgerLedgerSealedEventName = "LedgerSealed"

// ContractEventName returns the user-defined event name.
func (ILiabilityLedgerLedgerSealed) ContractEventName() string {
	return ILiabilityLedgerLedgerSealedEventName
}

// UnpackLedgerSealedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event LedgerSealed(bytes32 indexed head, uint64 nextEpoch, uint64 lastAsOf)
func (iLiabilityLedger *ILiabilityLedger) UnpackLedgerSealedEvent(log *types.Log) (*ILiabilityLedgerLedgerSealed, error) {
	event := "LedgerSealed"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != iLiabilityLedger.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(ILiabilityLedgerLedgerSealed)
	if len(log.Data) > 0 {
		if err := iLiabilityLedger.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iLiabilityLedger.abi.Events[event].Inputs {
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
func (iLiabilityLedger *ILiabilityLedger) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], iLiabilityLedger.abi.Errors["EmptyEpochMismatch"].ID.Bytes()[:4]) {
		return iLiabilityLedger.UnpackEmptyEpochMismatchError(raw[4:])
	}
	if bytes.Equal(raw[:4], iLiabilityLedger.abi.Errors["EpochAlreadyAnchored"].ID.Bytes()[:4]) {
		return iLiabilityLedger.UnpackEpochAlreadyAnchoredError(raw[4:])
	}
	if bytes.Equal(raw[:4], iLiabilityLedger.abi.Errors["EpochBeforeGenesis"].ID.Bytes()[:4]) {
		return iLiabilityLedger.UnpackEpochBeforeGenesisError(raw[4:])
	}
	if bytes.Equal(raw[:4], iLiabilityLedger.abi.Errors["EpochConflict"].ID.Bytes()[:4]) {
		return iLiabilityLedger.UnpackEpochConflictError(raw[4:])
	}
	if bytes.Equal(raw[:4], iLiabilityLedger.abi.Errors["EpochOutOfSequence"].ID.Bytes()[:4]) {
		return iLiabilityLedger.UnpackEpochOutOfSequenceError(raw[4:])
	}
	if bytes.Equal(raw[:4], iLiabilityLedger.abi.Errors["HeadMismatch"].ID.Bytes()[:4]) {
		return iLiabilityLedger.UnpackHeadMismatchError(raw[4:])
	}
	if bytes.Equal(raw[:4], iLiabilityLedger.abi.Errors["InvalidGenesis"].ID.Bytes()[:4]) {
		return iLiabilityLedger.UnpackInvalidGenesisError(raw[4:])
	}
	if bytes.Equal(raw[:4], iLiabilityLedger.abi.Errors["LedgerIsSealed"].ID.Bytes()[:4]) {
		return iLiabilityLedger.UnpackLedgerIsSealedError(raw[4:])
	}
	if bytes.Equal(raw[:4], iLiabilityLedger.abi.Errors["PreviousDigestMismatch"].ID.Bytes()[:4]) {
		return iLiabilityLedger.UnpackPreviousDigestMismatchError(raw[4:])
	}
	if bytes.Equal(raw[:4], iLiabilityLedger.abi.Errors["SnapshotInFuture"].ID.Bytes()[:4]) {
		return iLiabilityLedger.UnpackSnapshotInFutureError(raw[4:])
	}
	if bytes.Equal(raw[:4], iLiabilityLedger.abi.Errors["SnapshotNotAfterPrevious"].ID.Bytes()[:4]) {
		return iLiabilityLedger.UnpackSnapshotNotAfterPreviousError(raw[4:])
	}
	if bytes.Equal(raw[:4], iLiabilityLedger.abi.Errors["TooManyLeaves"].ID.Bytes()[:4]) {
		return iLiabilityLedger.UnpackTooManyLeavesError(raw[4:])
	}
	if bytes.Equal(raw[:4], iLiabilityLedger.abi.Errors["UnknownEpoch"].ID.Bytes()[:4]) {
		return iLiabilityLedger.UnpackUnknownEpochError(raw[4:])
	}
	if bytes.Equal(raw[:4], iLiabilityLedger.abi.Errors["ZeroAnchorer"].ID.Bytes()[:4]) {
		return iLiabilityLedger.UnpackZeroAnchorerError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// ILiabilityLedgerEmptyEpochMismatch represents a EmptyEpochMismatch error raised by the ILiabilityLedger contract.
type ILiabilityLedgerEmptyEpochMismatch struct {
	Root  [32]byte
	Total uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error EmptyEpochMismatch(bytes32 root, uint64 total)
func ILiabilityLedgerEmptyEpochMismatchErrorID() common.Hash {
	return common.HexToHash("0xde718304bb7fe1b2763221798581b76b906780448b2dcd8bad16597d786f1080")
}

// UnpackEmptyEpochMismatchError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error EmptyEpochMismatch(bytes32 root, uint64 total)
func (iLiabilityLedger *ILiabilityLedger) UnpackEmptyEpochMismatchError(raw []byte) (*ILiabilityLedgerEmptyEpochMismatch, error) {
	out := new(ILiabilityLedgerEmptyEpochMismatch)
	if err := iLiabilityLedger.abi.UnpackIntoInterface(out, "EmptyEpochMismatch", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ILiabilityLedgerEpochAlreadyAnchored represents a EpochAlreadyAnchored error raised by the ILiabilityLedger contract.
type ILiabilityLedgerEpochAlreadyAnchored struct {
	Number uint64
	Digest [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error EpochAlreadyAnchored(uint64 number, bytes32 digest)
func ILiabilityLedgerEpochAlreadyAnchoredErrorID() common.Hash {
	return common.HexToHash("0xab455a25015fe0023a0a76c842f75125738bc3b7b3167d58a41d09451f0baf74")
}

// UnpackEpochAlreadyAnchoredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error EpochAlreadyAnchored(uint64 number, bytes32 digest)
func (iLiabilityLedger *ILiabilityLedger) UnpackEpochAlreadyAnchoredError(raw []byte) (*ILiabilityLedgerEpochAlreadyAnchored, error) {
	out := new(ILiabilityLedgerEpochAlreadyAnchored)
	if err := iLiabilityLedger.abi.UnpackIntoInterface(out, "EpochAlreadyAnchored", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ILiabilityLedgerEpochBeforeGenesis represents a EpochBeforeGenesis error raised by the ILiabilityLedger contract.
type ILiabilityLedgerEpochBeforeGenesis struct {
	GenesisEpoch uint64
	Number       uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error EpochBeforeGenesis(uint64 genesisEpoch, uint64 number)
func ILiabilityLedgerEpochBeforeGenesisErrorID() common.Hash {
	return common.HexToHash("0xacc3e9cb2b490f2ca39781764307e7ead97b1740e688ec48ce82eb74ab3cad45")
}

// UnpackEpochBeforeGenesisError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error EpochBeforeGenesis(uint64 genesisEpoch, uint64 number)
func (iLiabilityLedger *ILiabilityLedger) UnpackEpochBeforeGenesisError(raw []byte) (*ILiabilityLedgerEpochBeforeGenesis, error) {
	out := new(ILiabilityLedgerEpochBeforeGenesis)
	if err := iLiabilityLedger.abi.UnpackIntoInterface(out, "EpochBeforeGenesis", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ILiabilityLedgerEpochConflict represents a EpochConflict error raised by the ILiabilityLedger contract.
type ILiabilityLedgerEpochConflict struct {
	Number   uint64
	Anchored [32]byte
	Given    [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error EpochConflict(uint64 number, bytes32 anchored, bytes32 given)
func ILiabilityLedgerEpochConflictErrorID() common.Hash {
	return common.HexToHash("0x6d613c397e2d7ea70b114e132c49a28c0fa5de67dbb8170276880d74174393a9")
}

// UnpackEpochConflictError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error EpochConflict(uint64 number, bytes32 anchored, bytes32 given)
func (iLiabilityLedger *ILiabilityLedger) UnpackEpochConflictError(raw []byte) (*ILiabilityLedgerEpochConflict, error) {
	out := new(ILiabilityLedgerEpochConflict)
	if err := iLiabilityLedger.abi.UnpackIntoInterface(out, "EpochConflict", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ILiabilityLedgerEpochOutOfSequence represents a EpochOutOfSequence error raised by the ILiabilityLedger contract.
type ILiabilityLedgerEpochOutOfSequence struct {
	Expected uint64
	Given    uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error EpochOutOfSequence(uint64 expected, uint64 given)
func ILiabilityLedgerEpochOutOfSequenceErrorID() common.Hash {
	return common.HexToHash("0x7a8d061d59616552e34e07c4e10cd96598f8a51489f859f21055dba983f36722")
}

// UnpackEpochOutOfSequenceError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error EpochOutOfSequence(uint64 expected, uint64 given)
func (iLiabilityLedger *ILiabilityLedger) UnpackEpochOutOfSequenceError(raw []byte) (*ILiabilityLedgerEpochOutOfSequence, error) {
	out := new(ILiabilityLedgerEpochOutOfSequence)
	if err := iLiabilityLedger.abi.UnpackIntoInterface(out, "EpochOutOfSequence", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ILiabilityLedgerHeadMismatch represents a HeadMismatch error raised by the ILiabilityLedger contract.
type ILiabilityLedgerHeadMismatch struct {
	Head     [32]byte
	Expected [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error HeadMismatch(bytes32 head, bytes32 expected)
func ILiabilityLedgerHeadMismatchErrorID() common.Hash {
	return common.HexToHash("0x6f2de3c75673389a279dc32b221eac892abdabf6b2d716499a572ff532f5ec50")
}

// UnpackHeadMismatchError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error HeadMismatch(bytes32 head, bytes32 expected)
func (iLiabilityLedger *ILiabilityLedger) UnpackHeadMismatchError(raw []byte) (*ILiabilityLedgerHeadMismatch, error) {
	out := new(ILiabilityLedgerHeadMismatch)
	if err := iLiabilityLedger.abi.UnpackIntoInterface(out, "HeadMismatch", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ILiabilityLedgerInvalidGenesis represents a InvalidGenesis error raised by the ILiabilityLedger contract.
type ILiabilityLedgerInvalidGenesis struct {
	GenesisEpoch  uint64
	GenesisDigest [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidGenesis(uint64 genesisEpoch, bytes32 genesisDigest)
func ILiabilityLedgerInvalidGenesisErrorID() common.Hash {
	return common.HexToHash("0x111f019fa2344bb16ed38186e3c5142be0c96d972db2a225af2e8f4835689530")
}

// UnpackInvalidGenesisError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidGenesis(uint64 genesisEpoch, bytes32 genesisDigest)
func (iLiabilityLedger *ILiabilityLedger) UnpackInvalidGenesisError(raw []byte) (*ILiabilityLedgerInvalidGenesis, error) {
	out := new(ILiabilityLedgerInvalidGenesis)
	if err := iLiabilityLedger.abi.UnpackIntoInterface(out, "InvalidGenesis", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ILiabilityLedgerLedgerIsSealed represents a LedgerIsSealed error raised by the ILiabilityLedger contract.
type ILiabilityLedgerLedgerIsSealed struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error LedgerIsSealed()
func ILiabilityLedgerLedgerIsSealedErrorID() common.Hash {
	return common.HexToHash("0x712a4e2f2cfa8aa62d1c3b9f2e3bce545dad83b3813e02d35f0b2a51811808e4")
}

// UnpackLedgerIsSealedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error LedgerIsSealed()
func (iLiabilityLedger *ILiabilityLedger) UnpackLedgerIsSealedError(raw []byte) (*ILiabilityLedgerLedgerIsSealed, error) {
	out := new(ILiabilityLedgerLedgerIsSealed)
	if err := iLiabilityLedger.abi.UnpackIntoInterface(out, "LedgerIsSealed", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ILiabilityLedgerPreviousDigestMismatch represents a PreviousDigestMismatch error raised by the ILiabilityLedger contract.
type ILiabilityLedgerPreviousDigestMismatch struct {
	Expected [32]byte
	Given    [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error PreviousDigestMismatch(bytes32 expected, bytes32 given)
func ILiabilityLedgerPreviousDigestMismatchErrorID() common.Hash {
	return common.HexToHash("0x0796b55ac9b13adadc488175cf6ef6fbf1565c2537f4e021d7268b5ab5a335d9")
}

// UnpackPreviousDigestMismatchError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error PreviousDigestMismatch(bytes32 expected, bytes32 given)
func (iLiabilityLedger *ILiabilityLedger) UnpackPreviousDigestMismatchError(raw []byte) (*ILiabilityLedgerPreviousDigestMismatch, error) {
	out := new(ILiabilityLedgerPreviousDigestMismatch)
	if err := iLiabilityLedger.abi.UnpackIntoInterface(out, "PreviousDigestMismatch", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ILiabilityLedgerSnapshotInFuture represents a SnapshotInFuture error raised by the ILiabilityLedger contract.
type ILiabilityLedgerSnapshotInFuture struct {
	AsOf      uint64
	BlockTime uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SnapshotInFuture(uint64 asOf, uint64 blockTime)
func ILiabilityLedgerSnapshotInFutureErrorID() common.Hash {
	return common.HexToHash("0x18aaced2679cfd7d188264841ef0e83939ee734887344b8d78b84cd47abff282")
}

// UnpackSnapshotInFutureError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SnapshotInFuture(uint64 asOf, uint64 blockTime)
func (iLiabilityLedger *ILiabilityLedger) UnpackSnapshotInFutureError(raw []byte) (*ILiabilityLedgerSnapshotInFuture, error) {
	out := new(ILiabilityLedgerSnapshotInFuture)
	if err := iLiabilityLedger.abi.UnpackIntoInterface(out, "SnapshotInFuture", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ILiabilityLedgerSnapshotNotAfterPrevious represents a SnapshotNotAfterPrevious error raised by the ILiabilityLedger contract.
type ILiabilityLedgerSnapshotNotAfterPrevious struct {
	AsOf     uint64
	Previous uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SnapshotNotAfterPrevious(uint64 asOf, uint64 previous)
func ILiabilityLedgerSnapshotNotAfterPreviousErrorID() common.Hash {
	return common.HexToHash("0x81ed60e995f8fcd4753ca10053fc7f547e1dd0c250e9bc06e427fee7bea5773f")
}

// UnpackSnapshotNotAfterPreviousError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SnapshotNotAfterPrevious(uint64 asOf, uint64 previous)
func (iLiabilityLedger *ILiabilityLedger) UnpackSnapshotNotAfterPreviousError(raw []byte) (*ILiabilityLedgerSnapshotNotAfterPrevious, error) {
	out := new(ILiabilityLedgerSnapshotNotAfterPrevious)
	if err := iLiabilityLedger.abi.UnpackIntoInterface(out, "SnapshotNotAfterPrevious", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ILiabilityLedgerTooManyLeaves represents a TooManyLeaves error raised by the ILiabilityLedger contract.
type ILiabilityLedgerTooManyLeaves struct {
	LeafCount uint64
	Maximum   uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TooManyLeaves(uint64 leafCount, uint64 maximum)
func ILiabilityLedgerTooManyLeavesErrorID() common.Hash {
	return common.HexToHash("0x09e51df967feec9808a68989cf05b6cd9c8c3c92676ae3ea00eda6a12380a9af")
}

// UnpackTooManyLeavesError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TooManyLeaves(uint64 leafCount, uint64 maximum)
func (iLiabilityLedger *ILiabilityLedger) UnpackTooManyLeavesError(raw []byte) (*ILiabilityLedgerTooManyLeaves, error) {
	out := new(ILiabilityLedgerTooManyLeaves)
	if err := iLiabilityLedger.abi.UnpackIntoInterface(out, "TooManyLeaves", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ILiabilityLedgerUnknownEpoch represents a UnknownEpoch error raised by the ILiabilityLedger contract.
type ILiabilityLedgerUnknownEpoch struct {
	Number uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error UnknownEpoch(uint64 number)
func ILiabilityLedgerUnknownEpochErrorID() common.Hash {
	return common.HexToHash("0xa4e21e6c8f881eff51b5764f98a67304231c1cb1ba8af636f941098552d2bd6c")
}

// UnpackUnknownEpochError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error UnknownEpoch(uint64 number)
func (iLiabilityLedger *ILiabilityLedger) UnpackUnknownEpochError(raw []byte) (*ILiabilityLedgerUnknownEpoch, error) {
	out := new(ILiabilityLedgerUnknownEpoch)
	if err := iLiabilityLedger.abi.UnpackIntoInterface(out, "UnknownEpoch", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ILiabilityLedgerZeroAnchorer represents a ZeroAnchorer error raised by the ILiabilityLedger contract.
type ILiabilityLedgerZeroAnchorer struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ZeroAnchorer()
func ILiabilityLedgerZeroAnchorerErrorID() common.Hash {
	return common.HexToHash("0x7110c4a1e207dd190eaa2e6c40aeb0e70affea5eb498857bf51f0d9d85807115")
}

// UnpackZeroAnchorerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ZeroAnchorer()
func (iLiabilityLedger *ILiabilityLedger) UnpackZeroAnchorerError(raw []byte) (*ILiabilityLedgerZeroAnchorer, error) {
	out := new(ILiabilityLedgerZeroAnchorer)
	if err := iLiabilityLedger.abi.UnpackIntoInterface(out, "ZeroAnchorer", raw); err != nil {
		return nil, err
	}
	return out, nil
}
