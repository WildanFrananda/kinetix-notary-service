// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package evidenceinterface

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

// Batch is an auto generated low-level Go binding around an user-defined struct.
type Batch struct {
	Root               [32]byte
	Digest             [32]byte
	LeafCount          uint64
	ThroughTransaction uint64
	ThroughSequence    uint64
	ClosedAt           uint64
	AnchoredAt         uint64
}

// BatchHeader is an auto generated low-level Go binding around an user-defined struct.
type BatchHeader struct {
	Stream             uint8
	Number             uint64
	PreviousDigest     [32]byte
	Root               [32]byte
	LeafCount          uint64
	ThroughTransaction uint64
	ThroughSequence    uint64
	ClosedAt           uint64
}

// StreamState is an auto generated low-level Go binding around an user-defined struct.
type StreamState struct {
	Known              bool
	Head               [32]byte
	GenesisBatch       uint64
	NextBatch          uint64
	LastClosedAt       uint64
	ThroughTransaction uint64
	ThroughSequence    uint64
}

// IEvidenceRegistryMetaData contains all meta data concerning the IEvidenceRegistry contract.
var IEvidenceRegistryMetaData = bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"anchor\",\"inputs\":[{\"name\":\"header\",\"type\":\"tuple\",\"internalType\":\"structBatchHeader\",\"components\":[{\"name\":\"stream\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"number\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"previousDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"leafCount\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"throughTransaction\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"throughSequence\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"closedAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"batchAt\",\"inputs\":[{\"name\":\"stream\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"number\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structBatch\",\"components\":[{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"digest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"leafCount\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"throughTransaction\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"throughSequence\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"closedAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"anchoredAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isSealed\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"scheme\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"seal\",\"inputs\":[{\"name\":\"expectedHeads\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"stateOf\",\"inputs\":[{\"name\":\"stream\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structStreamState\",\"components\":[{\"name\":\"known\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"head\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"genesisBatch\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"nextBatch\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"lastClosedAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"throughTransaction\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"throughSequence\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"streams\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8[]\",\"internalType\":\"uint8[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyInclusion\",\"inputs\":[{\"name\":\"stream\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"number\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"recordDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"index\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"siblings\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"BatchAnchored\",\"inputs\":[{\"name\":\"stream\",\"type\":\"uint8\",\"indexed\":true,\"internalType\":\"uint8\"},{\"name\":\"number\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"digest\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"header\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structBatchHeader\",\"components\":[{\"name\":\"stream\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"number\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"previousDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"leafCount\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"throughTransaction\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"throughSequence\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"closedAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StreamSealed\",\"inputs\":[{\"name\":\"stream\",\"type\":\"uint8\",\"indexed\":true,\"internalType\":\"uint8\"},{\"name\":\"head\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"nextBatch\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"lastClosedAt\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"throughTransaction\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"throughSequence\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AnchorerNotGranted\",\"inputs\":[{\"name\":\"anchorer\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"BatchAlreadyAnchored\",\"inputs\":[{\"name\":\"stream\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"number\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"digest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"BatchBeforeGenesis\",\"inputs\":[{\"name\":\"stream\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"genesisBatch\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"number\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"BatchConflict\",\"inputs\":[{\"name\":\"stream\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"number\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"anchored\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"offered\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"BatchOutOfSequence\",\"inputs\":[{\"name\":\"stream\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"expected\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"offered\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"CloseInFuture\",\"inputs\":[{\"name\":\"stream\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"closedAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"blockTime\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"CloseNotAfterPrevious\",\"inputs\":[{\"name\":\"stream\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"closedAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"lastClosedAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"EmptyBatch\",\"inputs\":[{\"name\":\"stream\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"HeadCountMismatch\",\"inputs\":[{\"name\":\"streams\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"offered\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"HeadMismatch\",\"inputs\":[{\"name\":\"stream\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"head\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"InvalidGenesis\",\"inputs\":[{\"name\":\"stream\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"batch\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"digest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"NoStreams\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PositionNotAfterPrevious\",\"inputs\":[{\"name\":\"stream\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"throughTransaction\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"throughSequence\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"previousTransaction\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"previousSequence\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"PreviousDigestMismatch\",\"inputs\":[{\"name\":\"stream\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"head\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"offered\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"RegistryIsSealed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"StreamsNotAscending\",\"inputs\":[{\"name\":\"previous\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"stream\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"TooManyLeaves\",\"inputs\":[{\"name\":\"leafCount\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"maximum\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"TooManyStreams\",\"inputs\":[{\"name\":\"streams\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maximum\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"UnknownBatch\",\"inputs\":[{\"name\":\"stream\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"number\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"UnknownStream\",\"inputs\":[{\"name\":\"stream\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"ZeroAnchorer\",\"inputs\":[]}]",
	ID:  "IEvidenceRegistry",
}

// IEvidenceRegistry is an auto generated Go binding around an Ethereum contract.
type IEvidenceRegistry struct {
	abi abi.ABI
}

// GetABI returns the ABI associated with this contract binding.
func (c *IEvidenceRegistry) GetABI() abi.ABI {
	return c.abi
}

// NewIEvidenceRegistry creates a new instance of IEvidenceRegistry.
func NewIEvidenceRegistry() *IEvidenceRegistry {
	parsed, err := IEvidenceRegistryMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &IEvidenceRegistry{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *IEvidenceRegistry) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackAnchor is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x97413ab0.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function anchor((uint8,uint64,bytes32,bytes32,uint64,uint64,uint64,uint64) header) returns(bytes32)
func (iEvidenceRegistry *IEvidenceRegistry) PackAnchor(header BatchHeader) []byte {
	enc, err := iEvidenceRegistry.abi.Pack("anchor", header)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAnchor is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x97413ab0.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function anchor((uint8,uint64,bytes32,bytes32,uint64,uint64,uint64,uint64) header) returns(bytes32)
func (iEvidenceRegistry *IEvidenceRegistry) TryPackAnchor(header BatchHeader) ([]byte, error) {
	return iEvidenceRegistry.abi.Pack("anchor", header)
}

// UnpackAnchor is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x97413ab0.
//
// Solidity: function anchor((uint8,uint64,bytes32,bytes32,uint64,uint64,uint64,uint64) header) returns(bytes32)
func (iEvidenceRegistry *IEvidenceRegistry) UnpackAnchor(data []byte) ([32]byte, error) {
	out, err := iEvidenceRegistry.abi.Unpack("anchor", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackBatchAt is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf53f43b0.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function batchAt(uint8 stream, uint64 number) view returns((bytes32,bytes32,uint64,uint64,uint64,uint64,uint64))
func (iEvidenceRegistry *IEvidenceRegistry) PackBatchAt(stream uint8, number uint64) []byte {
	enc, err := iEvidenceRegistry.abi.Pack("batchAt", stream, number)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackBatchAt is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf53f43b0.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function batchAt(uint8 stream, uint64 number) view returns((bytes32,bytes32,uint64,uint64,uint64,uint64,uint64))
func (iEvidenceRegistry *IEvidenceRegistry) TryPackBatchAt(stream uint8, number uint64) ([]byte, error) {
	return iEvidenceRegistry.abi.Pack("batchAt", stream, number)
}

// UnpackBatchAt is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xf53f43b0.
//
// Solidity: function batchAt(uint8 stream, uint64 number) view returns((bytes32,bytes32,uint64,uint64,uint64,uint64,uint64))
func (iEvidenceRegistry *IEvidenceRegistry) UnpackBatchAt(data []byte) (Batch, error) {
	out, err := iEvidenceRegistry.abi.Unpack("batchAt", data)
	if err != nil {
		return *new(Batch), err
	}
	out0 := *abi.ConvertType(out[0], new(Batch)).(*Batch)
	return out0, nil
}

// PackIsSealed is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x631f9852.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function isSealed() view returns(bool)
func (iEvidenceRegistry *IEvidenceRegistry) PackIsSealed() []byte {
	enc, err := iEvidenceRegistry.abi.Pack("isSealed")
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
func (iEvidenceRegistry *IEvidenceRegistry) TryPackIsSealed() ([]byte, error) {
	return iEvidenceRegistry.abi.Pack("isSealed")
}

// UnpackIsSealed is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x631f9852.
//
// Solidity: function isSealed() view returns(bool)
func (iEvidenceRegistry *IEvidenceRegistry) UnpackIsSealed(data []byte) (bool, error) {
	out, err := iEvidenceRegistry.abi.Unpack("isSealed", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackScheme is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8518b0ad.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function scheme() pure returns(uint16)
func (iEvidenceRegistry *IEvidenceRegistry) PackScheme() []byte {
	enc, err := iEvidenceRegistry.abi.Pack("scheme")
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
func (iEvidenceRegistry *IEvidenceRegistry) TryPackScheme() ([]byte, error) {
	return iEvidenceRegistry.abi.Pack("scheme")
}

// UnpackScheme is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8518b0ad.
//
// Solidity: function scheme() pure returns(uint16)
func (iEvidenceRegistry *IEvidenceRegistry) UnpackScheme(data []byte) (uint16, error) {
	out, err := iEvidenceRegistry.abi.Unpack("scheme", data)
	if err != nil {
		return *new(uint16), err
	}
	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)
	return out0, nil
}

// PackSeal is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8960b085.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function seal(bytes32[] expectedHeads) returns()
func (iEvidenceRegistry *IEvidenceRegistry) PackSeal(expectedHeads [][32]byte) []byte {
	enc, err := iEvidenceRegistry.abi.Pack("seal", expectedHeads)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSeal is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8960b085.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function seal(bytes32[] expectedHeads) returns()
func (iEvidenceRegistry *IEvidenceRegistry) TryPackSeal(expectedHeads [][32]byte) ([]byte, error) {
	return iEvidenceRegistry.abi.Pack("seal", expectedHeads)
}

// PackStateOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc7991ec7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function stateOf(uint8 stream) view returns((bool,bytes32,uint64,uint64,uint64,uint64,uint64))
func (iEvidenceRegistry *IEvidenceRegistry) PackStateOf(stream uint8) []byte {
	enc, err := iEvidenceRegistry.abi.Pack("stateOf", stream)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackStateOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc7991ec7.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function stateOf(uint8 stream) view returns((bool,bytes32,uint64,uint64,uint64,uint64,uint64))
func (iEvidenceRegistry *IEvidenceRegistry) TryPackStateOf(stream uint8) ([]byte, error) {
	return iEvidenceRegistry.abi.Pack("stateOf", stream)
}

// UnpackStateOf is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc7991ec7.
//
// Solidity: function stateOf(uint8 stream) view returns((bool,bytes32,uint64,uint64,uint64,uint64,uint64))
func (iEvidenceRegistry *IEvidenceRegistry) UnpackStateOf(data []byte) (StreamState, error) {
	out, err := iEvidenceRegistry.abi.Unpack("stateOf", data)
	if err != nil {
		return *new(StreamState), err
	}
	out0 := *abi.ConvertType(out[0], new(StreamState)).(*StreamState)
	return out0, nil
}

// PackStreams is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x32403040.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function streams() view returns(uint8[])
func (iEvidenceRegistry *IEvidenceRegistry) PackStreams() []byte {
	enc, err := iEvidenceRegistry.abi.Pack("streams")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackStreams is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x32403040.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function streams() view returns(uint8[])
func (iEvidenceRegistry *IEvidenceRegistry) TryPackStreams() ([]byte, error) {
	return iEvidenceRegistry.abi.Pack("streams")
}

// UnpackStreams is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x32403040.
//
// Solidity: function streams() view returns(uint8[])
func (iEvidenceRegistry *IEvidenceRegistry) UnpackStreams(data []byte) ([]uint8, error) {
	out, err := iEvidenceRegistry.abi.Unpack("streams", data)
	if err != nil {
		return *new([]uint8), err
	}
	out0 := *abi.ConvertType(out[0], new([]uint8)).(*[]uint8)
	return out0, nil
}

// PackVerifyInclusion is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe283c044.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function verifyInclusion(uint8 stream, uint64 number, bytes32 salt, bytes32 recordDigest, uint64 index, bytes32[] siblings) view returns(bool)
func (iEvidenceRegistry *IEvidenceRegistry) PackVerifyInclusion(stream uint8, number uint64, salt [32]byte, recordDigest [32]byte, index uint64, siblings [][32]byte) []byte {
	enc, err := iEvidenceRegistry.abi.Pack("verifyInclusion", stream, number, salt, recordDigest, index, siblings)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackVerifyInclusion is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe283c044.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function verifyInclusion(uint8 stream, uint64 number, bytes32 salt, bytes32 recordDigest, uint64 index, bytes32[] siblings) view returns(bool)
func (iEvidenceRegistry *IEvidenceRegistry) TryPackVerifyInclusion(stream uint8, number uint64, salt [32]byte, recordDigest [32]byte, index uint64, siblings [][32]byte) ([]byte, error) {
	return iEvidenceRegistry.abi.Pack("verifyInclusion", stream, number, salt, recordDigest, index, siblings)
}

// UnpackVerifyInclusion is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe283c044.
//
// Solidity: function verifyInclusion(uint8 stream, uint64 number, bytes32 salt, bytes32 recordDigest, uint64 index, bytes32[] siblings) view returns(bool)
func (iEvidenceRegistry *IEvidenceRegistry) UnpackVerifyInclusion(data []byte) (bool, error) {
	out, err := iEvidenceRegistry.abi.Unpack("verifyInclusion", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// IEvidenceRegistryBatchAnchored represents a BatchAnchored event raised by the IEvidenceRegistry contract.
type IEvidenceRegistryBatchAnchored struct {
	Stream uint8
	Number uint64
	Digest [32]byte
	Header BatchHeader
	Raw    *types.Log // Blockchain specific contextual infos
}

const IEvidenceRegistryBatchAnchoredEventName = "BatchAnchored"

// ContractEventName returns the user-defined event name.
func (IEvidenceRegistryBatchAnchored) ContractEventName() string {
	return IEvidenceRegistryBatchAnchoredEventName
}

// UnpackBatchAnchoredEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event BatchAnchored(uint8 indexed stream, uint64 indexed number, bytes32 digest, (uint8,uint64,bytes32,bytes32,uint64,uint64,uint64,uint64) header)
func (iEvidenceRegistry *IEvidenceRegistry) UnpackBatchAnchoredEvent(log *types.Log) (*IEvidenceRegistryBatchAnchored, error) {
	event := "BatchAnchored"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != iEvidenceRegistry.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(IEvidenceRegistryBatchAnchored)
	if len(log.Data) > 0 {
		if err := iEvidenceRegistry.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iEvidenceRegistry.abi.Events[event].Inputs {
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

// IEvidenceRegistryStreamSealed represents a StreamSealed event raised by the IEvidenceRegistry contract.
type IEvidenceRegistryStreamSealed struct {
	Stream             uint8
	Head               [32]byte
	NextBatch          uint64
	LastClosedAt       uint64
	ThroughTransaction uint64
	ThroughSequence    uint64
	Raw                *types.Log // Blockchain specific contextual infos
}

const IEvidenceRegistryStreamSealedEventName = "StreamSealed"

// ContractEventName returns the user-defined event name.
func (IEvidenceRegistryStreamSealed) ContractEventName() string {
	return IEvidenceRegistryStreamSealedEventName
}

// UnpackStreamSealedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event StreamSealed(uint8 indexed stream, bytes32 head, uint64 nextBatch, uint64 lastClosedAt, uint64 throughTransaction, uint64 throughSequence)
func (iEvidenceRegistry *IEvidenceRegistry) UnpackStreamSealedEvent(log *types.Log) (*IEvidenceRegistryStreamSealed, error) {
	event := "StreamSealed"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != iEvidenceRegistry.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(IEvidenceRegistryStreamSealed)
	if len(log.Data) > 0 {
		if err := iEvidenceRegistry.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iEvidenceRegistry.abi.Events[event].Inputs {
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
func (iEvidenceRegistry *IEvidenceRegistry) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], iEvidenceRegistry.abi.Errors["AnchorerNotGranted"].ID.Bytes()[:4]) {
		return iEvidenceRegistry.UnpackAnchorerNotGrantedError(raw[4:])
	}
	if bytes.Equal(raw[:4], iEvidenceRegistry.abi.Errors["BatchAlreadyAnchored"].ID.Bytes()[:4]) {
		return iEvidenceRegistry.UnpackBatchAlreadyAnchoredError(raw[4:])
	}
	if bytes.Equal(raw[:4], iEvidenceRegistry.abi.Errors["BatchBeforeGenesis"].ID.Bytes()[:4]) {
		return iEvidenceRegistry.UnpackBatchBeforeGenesisError(raw[4:])
	}
	if bytes.Equal(raw[:4], iEvidenceRegistry.abi.Errors["BatchConflict"].ID.Bytes()[:4]) {
		return iEvidenceRegistry.UnpackBatchConflictError(raw[4:])
	}
	if bytes.Equal(raw[:4], iEvidenceRegistry.abi.Errors["BatchOutOfSequence"].ID.Bytes()[:4]) {
		return iEvidenceRegistry.UnpackBatchOutOfSequenceError(raw[4:])
	}
	if bytes.Equal(raw[:4], iEvidenceRegistry.abi.Errors["CloseInFuture"].ID.Bytes()[:4]) {
		return iEvidenceRegistry.UnpackCloseInFutureError(raw[4:])
	}
	if bytes.Equal(raw[:4], iEvidenceRegistry.abi.Errors["CloseNotAfterPrevious"].ID.Bytes()[:4]) {
		return iEvidenceRegistry.UnpackCloseNotAfterPreviousError(raw[4:])
	}
	if bytes.Equal(raw[:4], iEvidenceRegistry.abi.Errors["EmptyBatch"].ID.Bytes()[:4]) {
		return iEvidenceRegistry.UnpackEmptyBatchError(raw[4:])
	}
	if bytes.Equal(raw[:4], iEvidenceRegistry.abi.Errors["HeadCountMismatch"].ID.Bytes()[:4]) {
		return iEvidenceRegistry.UnpackHeadCountMismatchError(raw[4:])
	}
	if bytes.Equal(raw[:4], iEvidenceRegistry.abi.Errors["HeadMismatch"].ID.Bytes()[:4]) {
		return iEvidenceRegistry.UnpackHeadMismatchError(raw[4:])
	}
	if bytes.Equal(raw[:4], iEvidenceRegistry.abi.Errors["InvalidGenesis"].ID.Bytes()[:4]) {
		return iEvidenceRegistry.UnpackInvalidGenesisError(raw[4:])
	}
	if bytes.Equal(raw[:4], iEvidenceRegistry.abi.Errors["NoStreams"].ID.Bytes()[:4]) {
		return iEvidenceRegistry.UnpackNoStreamsError(raw[4:])
	}
	if bytes.Equal(raw[:4], iEvidenceRegistry.abi.Errors["PositionNotAfterPrevious"].ID.Bytes()[:4]) {
		return iEvidenceRegistry.UnpackPositionNotAfterPreviousError(raw[4:])
	}
	if bytes.Equal(raw[:4], iEvidenceRegistry.abi.Errors["PreviousDigestMismatch"].ID.Bytes()[:4]) {
		return iEvidenceRegistry.UnpackPreviousDigestMismatchError(raw[4:])
	}
	if bytes.Equal(raw[:4], iEvidenceRegistry.abi.Errors["RegistryIsSealed"].ID.Bytes()[:4]) {
		return iEvidenceRegistry.UnpackRegistryIsSealedError(raw[4:])
	}
	if bytes.Equal(raw[:4], iEvidenceRegistry.abi.Errors["StreamsNotAscending"].ID.Bytes()[:4]) {
		return iEvidenceRegistry.UnpackStreamsNotAscendingError(raw[4:])
	}
	if bytes.Equal(raw[:4], iEvidenceRegistry.abi.Errors["TooManyLeaves"].ID.Bytes()[:4]) {
		return iEvidenceRegistry.UnpackTooManyLeavesError(raw[4:])
	}
	if bytes.Equal(raw[:4], iEvidenceRegistry.abi.Errors["TooManyStreams"].ID.Bytes()[:4]) {
		return iEvidenceRegistry.UnpackTooManyStreamsError(raw[4:])
	}
	if bytes.Equal(raw[:4], iEvidenceRegistry.abi.Errors["UnknownBatch"].ID.Bytes()[:4]) {
		return iEvidenceRegistry.UnpackUnknownBatchError(raw[4:])
	}
	if bytes.Equal(raw[:4], iEvidenceRegistry.abi.Errors["UnknownStream"].ID.Bytes()[:4]) {
		return iEvidenceRegistry.UnpackUnknownStreamError(raw[4:])
	}
	if bytes.Equal(raw[:4], iEvidenceRegistry.abi.Errors["ZeroAnchorer"].ID.Bytes()[:4]) {
		return iEvidenceRegistry.UnpackZeroAnchorerError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// IEvidenceRegistryAnchorerNotGranted represents a AnchorerNotGranted error raised by the IEvidenceRegistry contract.
type IEvidenceRegistryAnchorerNotGranted struct {
	Anchorer common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AnchorerNotGranted(address anchorer)
func IEvidenceRegistryAnchorerNotGrantedErrorID() common.Hash {
	return common.HexToHash("0x3012fc8bb5dd75c0c034d7c58aa1982d269d598072100fabc24d4de862d1e083")
}

// UnpackAnchorerNotGrantedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AnchorerNotGranted(address anchorer)
func (iEvidenceRegistry *IEvidenceRegistry) UnpackAnchorerNotGrantedError(raw []byte) (*IEvidenceRegistryAnchorerNotGranted, error) {
	out := new(IEvidenceRegistryAnchorerNotGranted)
	if err := iEvidenceRegistry.abi.UnpackIntoInterface(out, "AnchorerNotGranted", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IEvidenceRegistryBatchAlreadyAnchored represents a BatchAlreadyAnchored error raised by the IEvidenceRegistry contract.
type IEvidenceRegistryBatchAlreadyAnchored struct {
	Stream uint8
	Number uint64
	Digest [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error BatchAlreadyAnchored(uint8 stream, uint64 number, bytes32 digest)
func IEvidenceRegistryBatchAlreadyAnchoredErrorID() common.Hash {
	return common.HexToHash("0x6035ccca362de798c967fc5ddf5a627f48878f68f274d1ac779fe416246930c3")
}

// UnpackBatchAlreadyAnchoredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error BatchAlreadyAnchored(uint8 stream, uint64 number, bytes32 digest)
func (iEvidenceRegistry *IEvidenceRegistry) UnpackBatchAlreadyAnchoredError(raw []byte) (*IEvidenceRegistryBatchAlreadyAnchored, error) {
	out := new(IEvidenceRegistryBatchAlreadyAnchored)
	if err := iEvidenceRegistry.abi.UnpackIntoInterface(out, "BatchAlreadyAnchored", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IEvidenceRegistryBatchBeforeGenesis represents a BatchBeforeGenesis error raised by the IEvidenceRegistry contract.
type IEvidenceRegistryBatchBeforeGenesis struct {
	Stream       uint8
	GenesisBatch uint64
	Number       uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error BatchBeforeGenesis(uint8 stream, uint64 genesisBatch, uint64 number)
func IEvidenceRegistryBatchBeforeGenesisErrorID() common.Hash {
	return common.HexToHash("0x0d0878958571059a7616862a03367cf877644fac9d790878e6167d68f2f289d7")
}

// UnpackBatchBeforeGenesisError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error BatchBeforeGenesis(uint8 stream, uint64 genesisBatch, uint64 number)
func (iEvidenceRegistry *IEvidenceRegistry) UnpackBatchBeforeGenesisError(raw []byte) (*IEvidenceRegistryBatchBeforeGenesis, error) {
	out := new(IEvidenceRegistryBatchBeforeGenesis)
	if err := iEvidenceRegistry.abi.UnpackIntoInterface(out, "BatchBeforeGenesis", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IEvidenceRegistryBatchConflict represents a BatchConflict error raised by the IEvidenceRegistry contract.
type IEvidenceRegistryBatchConflict struct {
	Stream   uint8
	Number   uint64
	Anchored [32]byte
	Offered  [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error BatchConflict(uint8 stream, uint64 number, bytes32 anchored, bytes32 offered)
func IEvidenceRegistryBatchConflictErrorID() common.Hash {
	return common.HexToHash("0xc071754cfdbb9c05bd1a514ef76b5572b760af6ef162ed46016afd316167776a")
}

// UnpackBatchConflictError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error BatchConflict(uint8 stream, uint64 number, bytes32 anchored, bytes32 offered)
func (iEvidenceRegistry *IEvidenceRegistry) UnpackBatchConflictError(raw []byte) (*IEvidenceRegistryBatchConflict, error) {
	out := new(IEvidenceRegistryBatchConflict)
	if err := iEvidenceRegistry.abi.UnpackIntoInterface(out, "BatchConflict", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IEvidenceRegistryBatchOutOfSequence represents a BatchOutOfSequence error raised by the IEvidenceRegistry contract.
type IEvidenceRegistryBatchOutOfSequence struct {
	Stream   uint8
	Expected uint64
	Offered  uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error BatchOutOfSequence(uint8 stream, uint64 expected, uint64 offered)
func IEvidenceRegistryBatchOutOfSequenceErrorID() common.Hash {
	return common.HexToHash("0xb39c3985f6b4ab5807374509d65a2c99793b70bd7acb0641f36bce34c38b3faa")
}

// UnpackBatchOutOfSequenceError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error BatchOutOfSequence(uint8 stream, uint64 expected, uint64 offered)
func (iEvidenceRegistry *IEvidenceRegistry) UnpackBatchOutOfSequenceError(raw []byte) (*IEvidenceRegistryBatchOutOfSequence, error) {
	out := new(IEvidenceRegistryBatchOutOfSequence)
	if err := iEvidenceRegistry.abi.UnpackIntoInterface(out, "BatchOutOfSequence", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IEvidenceRegistryCloseInFuture represents a CloseInFuture error raised by the IEvidenceRegistry contract.
type IEvidenceRegistryCloseInFuture struct {
	Stream    uint8
	ClosedAt  uint64
	BlockTime uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error CloseInFuture(uint8 stream, uint64 closedAt, uint64 blockTime)
func IEvidenceRegistryCloseInFutureErrorID() common.Hash {
	return common.HexToHash("0xe9f968e7bdcd7fc5b23f91fd07ad9a18f73c622fa09f8198abdfdee58064dab5")
}

// UnpackCloseInFutureError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error CloseInFuture(uint8 stream, uint64 closedAt, uint64 blockTime)
func (iEvidenceRegistry *IEvidenceRegistry) UnpackCloseInFutureError(raw []byte) (*IEvidenceRegistryCloseInFuture, error) {
	out := new(IEvidenceRegistryCloseInFuture)
	if err := iEvidenceRegistry.abi.UnpackIntoInterface(out, "CloseInFuture", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IEvidenceRegistryCloseNotAfterPrevious represents a CloseNotAfterPrevious error raised by the IEvidenceRegistry contract.
type IEvidenceRegistryCloseNotAfterPrevious struct {
	Stream       uint8
	ClosedAt     uint64
	LastClosedAt uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error CloseNotAfterPrevious(uint8 stream, uint64 closedAt, uint64 lastClosedAt)
func IEvidenceRegistryCloseNotAfterPreviousErrorID() common.Hash {
	return common.HexToHash("0x5999344d325d28b495bd9ad3396675712366c905f2fbbcb3621b879e92eb13fc")
}

// UnpackCloseNotAfterPreviousError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error CloseNotAfterPrevious(uint8 stream, uint64 closedAt, uint64 lastClosedAt)
func (iEvidenceRegistry *IEvidenceRegistry) UnpackCloseNotAfterPreviousError(raw []byte) (*IEvidenceRegistryCloseNotAfterPrevious, error) {
	out := new(IEvidenceRegistryCloseNotAfterPrevious)
	if err := iEvidenceRegistry.abi.UnpackIntoInterface(out, "CloseNotAfterPrevious", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IEvidenceRegistryEmptyBatch represents a EmptyBatch error raised by the IEvidenceRegistry contract.
type IEvidenceRegistryEmptyBatch struct {
	Stream uint8
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error EmptyBatch(uint8 stream)
func IEvidenceRegistryEmptyBatchErrorID() common.Hash {
	return common.HexToHash("0x2aa62fa976a6325b695bf70c580caa4fea6f6fbc4ac0bca15c9492be85d38acf")
}

// UnpackEmptyBatchError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error EmptyBatch(uint8 stream)
func (iEvidenceRegistry *IEvidenceRegistry) UnpackEmptyBatchError(raw []byte) (*IEvidenceRegistryEmptyBatch, error) {
	out := new(IEvidenceRegistryEmptyBatch)
	if err := iEvidenceRegistry.abi.UnpackIntoInterface(out, "EmptyBatch", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IEvidenceRegistryHeadCountMismatch represents a HeadCountMismatch error raised by the IEvidenceRegistry contract.
type IEvidenceRegistryHeadCountMismatch struct {
	Streams *big.Int
	Offered *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error HeadCountMismatch(uint256 streams, uint256 offered)
func IEvidenceRegistryHeadCountMismatchErrorID() common.Hash {
	return common.HexToHash("0xbe8ea0e772df6ee158b79c8f7156d97dcac0737b032a43f94f9b385463afe2f9")
}

// UnpackHeadCountMismatchError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error HeadCountMismatch(uint256 streams, uint256 offered)
func (iEvidenceRegistry *IEvidenceRegistry) UnpackHeadCountMismatchError(raw []byte) (*IEvidenceRegistryHeadCountMismatch, error) {
	out := new(IEvidenceRegistryHeadCountMismatch)
	if err := iEvidenceRegistry.abi.UnpackIntoInterface(out, "HeadCountMismatch", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IEvidenceRegistryHeadMismatch represents a HeadMismatch error raised by the IEvidenceRegistry contract.
type IEvidenceRegistryHeadMismatch struct {
	Stream   uint8
	Head     [32]byte
	Expected [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error HeadMismatch(uint8 stream, bytes32 head, bytes32 expected)
func IEvidenceRegistryHeadMismatchErrorID() common.Hash {
	return common.HexToHash("0xf18e8ed0e05571020e598594dc2d046464d499ec16c478bf278c91d7362bde85")
}

// UnpackHeadMismatchError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error HeadMismatch(uint8 stream, bytes32 head, bytes32 expected)
func (iEvidenceRegistry *IEvidenceRegistry) UnpackHeadMismatchError(raw []byte) (*IEvidenceRegistryHeadMismatch, error) {
	out := new(IEvidenceRegistryHeadMismatch)
	if err := iEvidenceRegistry.abi.UnpackIntoInterface(out, "HeadMismatch", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IEvidenceRegistryInvalidGenesis represents a InvalidGenesis error raised by the IEvidenceRegistry contract.
type IEvidenceRegistryInvalidGenesis struct {
	Stream uint8
	Batch  uint64
	Digest [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidGenesis(uint8 stream, uint64 batch, bytes32 digest)
func IEvidenceRegistryInvalidGenesisErrorID() common.Hash {
	return common.HexToHash("0xad14a26107aa7d20646398124017110e79feed4739667c3bf81dfec51aceb429")
}

// UnpackInvalidGenesisError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidGenesis(uint8 stream, uint64 batch, bytes32 digest)
func (iEvidenceRegistry *IEvidenceRegistry) UnpackInvalidGenesisError(raw []byte) (*IEvidenceRegistryInvalidGenesis, error) {
	out := new(IEvidenceRegistryInvalidGenesis)
	if err := iEvidenceRegistry.abi.UnpackIntoInterface(out, "InvalidGenesis", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IEvidenceRegistryNoStreams represents a NoStreams error raised by the IEvidenceRegistry contract.
type IEvidenceRegistryNoStreams struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NoStreams()
func IEvidenceRegistryNoStreamsErrorID() common.Hash {
	return common.HexToHash("0x667ada5b150a16e1d46a52935ccb87bcb03cecee4f101cba3897ac4f9a77143f")
}

// UnpackNoStreamsError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NoStreams()
func (iEvidenceRegistry *IEvidenceRegistry) UnpackNoStreamsError(raw []byte) (*IEvidenceRegistryNoStreams, error) {
	out := new(IEvidenceRegistryNoStreams)
	if err := iEvidenceRegistry.abi.UnpackIntoInterface(out, "NoStreams", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IEvidenceRegistryPositionNotAfterPrevious represents a PositionNotAfterPrevious error raised by the IEvidenceRegistry contract.
type IEvidenceRegistryPositionNotAfterPrevious struct {
	Stream              uint8
	ThroughTransaction  uint64
	ThroughSequence     uint64
	PreviousTransaction uint64
	PreviousSequence    uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error PositionNotAfterPrevious(uint8 stream, uint64 throughTransaction, uint64 throughSequence, uint64 previousTransaction, uint64 previousSequence)
func IEvidenceRegistryPositionNotAfterPreviousErrorID() common.Hash {
	return common.HexToHash("0x9e8b48b4519e7edd0cad5dc17e4e2df6d8c29b0a42851e5994ce09adc55f67b5")
}

// UnpackPositionNotAfterPreviousError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error PositionNotAfterPrevious(uint8 stream, uint64 throughTransaction, uint64 throughSequence, uint64 previousTransaction, uint64 previousSequence)
func (iEvidenceRegistry *IEvidenceRegistry) UnpackPositionNotAfterPreviousError(raw []byte) (*IEvidenceRegistryPositionNotAfterPrevious, error) {
	out := new(IEvidenceRegistryPositionNotAfterPrevious)
	if err := iEvidenceRegistry.abi.UnpackIntoInterface(out, "PositionNotAfterPrevious", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IEvidenceRegistryPreviousDigestMismatch represents a PreviousDigestMismatch error raised by the IEvidenceRegistry contract.
type IEvidenceRegistryPreviousDigestMismatch struct {
	Stream  uint8
	Head    [32]byte
	Offered [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error PreviousDigestMismatch(uint8 stream, bytes32 head, bytes32 offered)
func IEvidenceRegistryPreviousDigestMismatchErrorID() common.Hash {
	return common.HexToHash("0x77a6bbb39f36ad092c0e2b13b8874b5c34d92fe800a43539f1094ea12480bff2")
}

// UnpackPreviousDigestMismatchError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error PreviousDigestMismatch(uint8 stream, bytes32 head, bytes32 offered)
func (iEvidenceRegistry *IEvidenceRegistry) UnpackPreviousDigestMismatchError(raw []byte) (*IEvidenceRegistryPreviousDigestMismatch, error) {
	out := new(IEvidenceRegistryPreviousDigestMismatch)
	if err := iEvidenceRegistry.abi.UnpackIntoInterface(out, "PreviousDigestMismatch", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IEvidenceRegistryRegistryIsSealed represents a RegistryIsSealed error raised by the IEvidenceRegistry contract.
type IEvidenceRegistryRegistryIsSealed struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RegistryIsSealed()
func IEvidenceRegistryRegistryIsSealedErrorID() common.Hash {
	return common.HexToHash("0xf3e6ff318836f081aa74168a8619b76f55c7d06015c7eb8be8b0b2fe870504c0")
}

// UnpackRegistryIsSealedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RegistryIsSealed()
func (iEvidenceRegistry *IEvidenceRegistry) UnpackRegistryIsSealedError(raw []byte) (*IEvidenceRegistryRegistryIsSealed, error) {
	out := new(IEvidenceRegistryRegistryIsSealed)
	if err := iEvidenceRegistry.abi.UnpackIntoInterface(out, "RegistryIsSealed", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IEvidenceRegistryStreamsNotAscending represents a StreamsNotAscending error raised by the IEvidenceRegistry contract.
type IEvidenceRegistryStreamsNotAscending struct {
	Previous uint8
	Stream   uint8
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error StreamsNotAscending(uint8 previous, uint8 stream)
func IEvidenceRegistryStreamsNotAscendingErrorID() common.Hash {
	return common.HexToHash("0x43f455cb7c6b84e12de47ffb24d6d2c6328f1d95892d72b2da12dc08c2ead064")
}

// UnpackStreamsNotAscendingError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error StreamsNotAscending(uint8 previous, uint8 stream)
func (iEvidenceRegistry *IEvidenceRegistry) UnpackStreamsNotAscendingError(raw []byte) (*IEvidenceRegistryStreamsNotAscending, error) {
	out := new(IEvidenceRegistryStreamsNotAscending)
	if err := iEvidenceRegistry.abi.UnpackIntoInterface(out, "StreamsNotAscending", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IEvidenceRegistryTooManyLeaves represents a TooManyLeaves error raised by the IEvidenceRegistry contract.
type IEvidenceRegistryTooManyLeaves struct {
	LeafCount uint64
	Maximum   uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TooManyLeaves(uint64 leafCount, uint64 maximum)
func IEvidenceRegistryTooManyLeavesErrorID() common.Hash {
	return common.HexToHash("0x09e51df967feec9808a68989cf05b6cd9c8c3c92676ae3ea00eda6a12380a9af")
}

// UnpackTooManyLeavesError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TooManyLeaves(uint64 leafCount, uint64 maximum)
func (iEvidenceRegistry *IEvidenceRegistry) UnpackTooManyLeavesError(raw []byte) (*IEvidenceRegistryTooManyLeaves, error) {
	out := new(IEvidenceRegistryTooManyLeaves)
	if err := iEvidenceRegistry.abi.UnpackIntoInterface(out, "TooManyLeaves", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IEvidenceRegistryTooManyStreams represents a TooManyStreams error raised by the IEvidenceRegistry contract.
type IEvidenceRegistryTooManyStreams struct {
	Streams *big.Int
	Maximum *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TooManyStreams(uint256 streams, uint256 maximum)
func IEvidenceRegistryTooManyStreamsErrorID() common.Hash {
	return common.HexToHash("0x7e1d37064d99a880f4dc620cfffcf10ad387324ee240ad76ca4bb0609b1ba381")
}

// UnpackTooManyStreamsError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TooManyStreams(uint256 streams, uint256 maximum)
func (iEvidenceRegistry *IEvidenceRegistry) UnpackTooManyStreamsError(raw []byte) (*IEvidenceRegistryTooManyStreams, error) {
	out := new(IEvidenceRegistryTooManyStreams)
	if err := iEvidenceRegistry.abi.UnpackIntoInterface(out, "TooManyStreams", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IEvidenceRegistryUnknownBatch represents a UnknownBatch error raised by the IEvidenceRegistry contract.
type IEvidenceRegistryUnknownBatch struct {
	Stream uint8
	Number uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error UnknownBatch(uint8 stream, uint64 number)
func IEvidenceRegistryUnknownBatchErrorID() common.Hash {
	return common.HexToHash("0x2293f2100b7b5077c8aef5599409d2ed688c284464fee8b993b5a163196812b1")
}

// UnpackUnknownBatchError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error UnknownBatch(uint8 stream, uint64 number)
func (iEvidenceRegistry *IEvidenceRegistry) UnpackUnknownBatchError(raw []byte) (*IEvidenceRegistryUnknownBatch, error) {
	out := new(IEvidenceRegistryUnknownBatch)
	if err := iEvidenceRegistry.abi.UnpackIntoInterface(out, "UnknownBatch", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IEvidenceRegistryUnknownStream represents a UnknownStream error raised by the IEvidenceRegistry contract.
type IEvidenceRegistryUnknownStream struct {
	Stream uint8
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error UnknownStream(uint8 stream)
func IEvidenceRegistryUnknownStreamErrorID() common.Hash {
	return common.HexToHash("0x8b7e1f972d87d002a7624f1a62232147b6c344d8b0ea791a8a21b9a8044abca8")
}

// UnpackUnknownStreamError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error UnknownStream(uint8 stream)
func (iEvidenceRegistry *IEvidenceRegistry) UnpackUnknownStreamError(raw []byte) (*IEvidenceRegistryUnknownStream, error) {
	out := new(IEvidenceRegistryUnknownStream)
	if err := iEvidenceRegistry.abi.UnpackIntoInterface(out, "UnknownStream", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IEvidenceRegistryZeroAnchorer represents a ZeroAnchorer error raised by the IEvidenceRegistry contract.
type IEvidenceRegistryZeroAnchorer struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ZeroAnchorer()
func IEvidenceRegistryZeroAnchorerErrorID() common.Hash {
	return common.HexToHash("0x7110c4a1e207dd190eaa2e6c40aeb0e70affea5eb498857bf51f0d9d85807115")
}

// UnpackZeroAnchorerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ZeroAnchorer()
func (iEvidenceRegistry *IEvidenceRegistry) UnpackZeroAnchorerError(raw []byte) (*IEvidenceRegistryZeroAnchorer, error) {
	out := new(IEvidenceRegistryZeroAnchorer)
	if err := iEvidenceRegistry.abi.UnpackIntoInterface(out, "ZeroAnchorer", raw); err != nil {
		return nil, err
	}
	return out, nil
}
