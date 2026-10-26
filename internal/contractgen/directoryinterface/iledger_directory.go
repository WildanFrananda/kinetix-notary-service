// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package directoryinterface

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

// Generation is an auto generated low-level Go binding around an user-defined struct.
type Generation struct {
	ChainId       uint64
	Ledger        common.Address
	Court         common.Address
	Withholding   common.Address
	GenesisEpoch  uint64
	GenesisDigest [32]byte
	Scheme        uint16
	RegisteredAt  uint64
}

// PendingGeneration is an auto generated low-level Go binding around an user-defined struct.
type PendingGeneration struct {
	Generation  Generation
	ActivatesAt uint64
}

// ILedgerDirectoryMetaData contains all meta data concerning the ILedgerDirectory contract.
var ILedgerDirectoryMetaData = bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"activateForeignSuccessor\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"cancelForeignSuccessor\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"current\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structGeneration\",\"components\":[{\"name\":\"chainId\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"ledger\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"court\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"withholding\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"genesisEpoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"genesisDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"scheme\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"registeredAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"foreignSuccessorDelay\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"generation\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structGeneration\",\"components\":[{\"name\":\"chainId\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"ledger\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"court\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"withholding\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"genesisEpoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"genesisDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"scheme\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"registeredAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"generationCount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"generationOf\",\"inputs\":[{\"name\":\"epoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structGeneration\",\"components\":[{\"name\":\"chainId\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"ledger\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"court\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"withholding\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"genesisEpoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"genesisDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"scheme\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"registeredAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingSuccessor\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structPendingGeneration\",\"components\":[{\"name\":\"generation\",\"type\":\"tuple\",\"internalType\":\"structGeneration\",\"components\":[{\"name\":\"chainId\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"ledger\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"court\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"withholding\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"genesisEpoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"genesisDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"scheme\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"registeredAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"activatesAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proposeForeignSuccessor\",\"inputs\":[{\"name\":\"chainId\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"ledger\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"court\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"withholding\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"scheme\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registerSuccessor\",\"inputs\":[{\"name\":\"ledger\",\"type\":\"address\",\"internalType\":\"contractILiabilityLedger\"},{\"name\":\"court\",\"type\":\"address\",\"internalType\":\"contractIStatementCourt\"},{\"name\":\"withholding\",\"type\":\"address\",\"internalType\":\"contractIWithholdingCourt\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"ForeignSuccessorCancelled\",\"inputs\":[{\"name\":\"chainId\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"ledger\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ForeignSuccessorProposed\",\"inputs\":[{\"name\":\"chainId\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"ledger\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"court\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"withholding\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"scheme\",\"type\":\"uint16\",\"indexed\":false,\"internalType\":\"uint16\"},{\"name\":\"activatesAt\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"GenerationRegistered\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"chainId\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"ledger\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"court\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"withholding\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"genesisEpoch\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"genesisDigest\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"scheme\",\"type\":\"uint16\",\"indexed\":false,\"internalType\":\"uint16\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"CourtOfAnotherLedger\",\"inputs\":[{\"name\":\"court\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"ledger\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"DirectoryMoved\",\"inputs\":[{\"name\":\"chainId\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"EpochBeforeDirectory\",\"inputs\":[{\"name\":\"epoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"genesisEpoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"GenesisMismatch\",\"inputs\":[{\"name\":\"expectedEpoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"expectedDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"givenEpoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"givenDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"NoSuccessorPending\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotALedger\",\"inputs\":[{\"name\":\"ledger\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"NotAStatementCourt\",\"inputs\":[{\"name\":\"court\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"NotAWithholdingCourt\",\"inputs\":[{\"name\":\"withholding\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"NotForeign\",\"inputs\":[{\"name\":\"chainId\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"PredecessorNotSealed\",\"inputs\":[{\"name\":\"ledger\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"SuccessorPending\",\"inputs\":[{\"name\":\"activatesAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"SuccessorSealed\",\"inputs\":[{\"name\":\"ledger\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"SuccessorTooEarly\",\"inputs\":[{\"name\":\"activatesAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"blockTime\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"UnknownGeneration\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"count\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"WithholdingOfAnotherCourt\",\"inputs\":[{\"name\":\"withholding\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"court\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ZeroAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroDelay\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroScheme\",\"inputs\":[]}]",
	ID:  "ILedgerDirectory",
}

// ILedgerDirectory is an auto generated Go binding around an Ethereum contract.
type ILedgerDirectory struct {
	abi abi.ABI
}

// GetABI returns the ABI associated with this contract binding.
func (c *ILedgerDirectory) GetABI() abi.ABI {
	return c.abi
}

// NewILedgerDirectory creates a new instance of ILedgerDirectory.
func NewILedgerDirectory() *ILedgerDirectory {
	parsed, err := ILedgerDirectoryMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &ILedgerDirectory{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *ILedgerDirectory) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackActivateForeignSuccessor is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6331e0a0.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function activateForeignSuccessor() returns(uint256)
func (iLedgerDirectory *ILedgerDirectory) PackActivateForeignSuccessor() []byte {
	enc, err := iLedgerDirectory.abi.Pack("activateForeignSuccessor")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackActivateForeignSuccessor is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6331e0a0.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function activateForeignSuccessor() returns(uint256)
func (iLedgerDirectory *ILedgerDirectory) TryPackActivateForeignSuccessor() ([]byte, error) {
	return iLedgerDirectory.abi.Pack("activateForeignSuccessor")
}

// UnpackActivateForeignSuccessor is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x6331e0a0.
//
// Solidity: function activateForeignSuccessor() returns(uint256)
func (iLedgerDirectory *ILedgerDirectory) UnpackActivateForeignSuccessor(data []byte) (*big.Int, error) {
	out, err := iLedgerDirectory.abi.Unpack("activateForeignSuccessor", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackCancelForeignSuccessor is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x084d1120.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function cancelForeignSuccessor() returns()
func (iLedgerDirectory *ILedgerDirectory) PackCancelForeignSuccessor() []byte {
	enc, err := iLedgerDirectory.abi.Pack("cancelForeignSuccessor")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCancelForeignSuccessor is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x084d1120.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function cancelForeignSuccessor() returns()
func (iLedgerDirectory *ILedgerDirectory) TryPackCancelForeignSuccessor() ([]byte, error) {
	return iLedgerDirectory.abi.Pack("cancelForeignSuccessor")
}

// PackCurrent is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9fa6a6e3.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function current() view returns((uint64,address,address,address,uint64,bytes32,uint16,uint64))
func (iLedgerDirectory *ILedgerDirectory) PackCurrent() []byte {
	enc, err := iLedgerDirectory.abi.Pack("current")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCurrent is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9fa6a6e3.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function current() view returns((uint64,address,address,address,uint64,bytes32,uint16,uint64))
func (iLedgerDirectory *ILedgerDirectory) TryPackCurrent() ([]byte, error) {
	return iLedgerDirectory.abi.Pack("current")
}

// UnpackCurrent is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x9fa6a6e3.
//
// Solidity: function current() view returns((uint64,address,address,address,uint64,bytes32,uint16,uint64))
func (iLedgerDirectory *ILedgerDirectory) UnpackCurrent(data []byte) (Generation, error) {
	out, err := iLedgerDirectory.abi.Unpack("current", data)
	if err != nil {
		return *new(Generation), err
	}
	out0 := *abi.ConvertType(out[0], new(Generation)).(*Generation)
	return out0, nil
}

// PackForeignSuccessorDelay is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xaaeabc25.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function foreignSuccessorDelay() view returns(uint64)
func (iLedgerDirectory *ILedgerDirectory) PackForeignSuccessorDelay() []byte {
	enc, err := iLedgerDirectory.abi.Pack("foreignSuccessorDelay")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackForeignSuccessorDelay is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xaaeabc25.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function foreignSuccessorDelay() view returns(uint64)
func (iLedgerDirectory *ILedgerDirectory) TryPackForeignSuccessorDelay() ([]byte, error) {
	return iLedgerDirectory.abi.Pack("foreignSuccessorDelay")
}

// UnpackForeignSuccessorDelay is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xaaeabc25.
//
// Solidity: function foreignSuccessorDelay() view returns(uint64)
func (iLedgerDirectory *ILedgerDirectory) UnpackForeignSuccessorDelay(data []byte) (uint64, error) {
	out, err := iLedgerDirectory.abi.Unpack("foreignSuccessorDelay", data)
	if err != nil {
		return *new(uint64), err
	}
	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)
	return out0, nil
}

// PackGeneration is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7d71dc35.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function generation(uint256 index) view returns((uint64,address,address,address,uint64,bytes32,uint16,uint64))
func (iLedgerDirectory *ILedgerDirectory) PackGeneration(index *big.Int) []byte {
	enc, err := iLedgerDirectory.abi.Pack("generation", index)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGeneration is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7d71dc35.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function generation(uint256 index) view returns((uint64,address,address,address,uint64,bytes32,uint16,uint64))
func (iLedgerDirectory *ILedgerDirectory) TryPackGeneration(index *big.Int) ([]byte, error) {
	return iLedgerDirectory.abi.Pack("generation", index)
}

// UnpackGeneration is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x7d71dc35.
//
// Solidity: function generation(uint256 index) view returns((uint64,address,address,address,uint64,bytes32,uint16,uint64))
func (iLedgerDirectory *ILedgerDirectory) UnpackGeneration(data []byte) (Generation, error) {
	out, err := iLedgerDirectory.abi.Unpack("generation", data)
	if err != nil {
		return *new(Generation), err
	}
	out0 := *abi.ConvertType(out[0], new(Generation)).(*Generation)
	return out0, nil
}

// PackGenerationCount is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x17b67acc.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function generationCount() view returns(uint256)
func (iLedgerDirectory *ILedgerDirectory) PackGenerationCount() []byte {
	enc, err := iLedgerDirectory.abi.Pack("generationCount")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGenerationCount is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x17b67acc.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function generationCount() view returns(uint256)
func (iLedgerDirectory *ILedgerDirectory) TryPackGenerationCount() ([]byte, error) {
	return iLedgerDirectory.abi.Pack("generationCount")
}

// UnpackGenerationCount is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x17b67acc.
//
// Solidity: function generationCount() view returns(uint256)
func (iLedgerDirectory *ILedgerDirectory) UnpackGenerationCount(data []byte) (*big.Int, error) {
	out, err := iLedgerDirectory.abi.Unpack("generationCount", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackGenerationOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc9a2133d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function generationOf(uint64 epoch) view returns(uint256, (uint64,address,address,address,uint64,bytes32,uint16,uint64))
func (iLedgerDirectory *ILedgerDirectory) PackGenerationOf(epoch uint64) []byte {
	enc, err := iLedgerDirectory.abi.Pack("generationOf", epoch)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGenerationOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc9a2133d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function generationOf(uint64 epoch) view returns(uint256, (uint64,address,address,address,uint64,bytes32,uint16,uint64))
func (iLedgerDirectory *ILedgerDirectory) TryPackGenerationOf(epoch uint64) ([]byte, error) {
	return iLedgerDirectory.abi.Pack("generationOf", epoch)
}

// GenerationOfOutput serves as a container for the return parameters of contract
// method GenerationOf.
type GenerationOfOutput struct {
	Arg0 *big.Int
	Arg1 Generation
}

// UnpackGenerationOf is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc9a2133d.
//
// Solidity: function generationOf(uint64 epoch) view returns(uint256, (uint64,address,address,address,uint64,bytes32,uint16,uint64))
func (iLedgerDirectory *ILedgerDirectory) UnpackGenerationOf(data []byte) (GenerationOfOutput, error) {
	out, err := iLedgerDirectory.abi.Unpack("generationOf", data)
	outstruct := new(GenerationOfOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Arg0 = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.Arg1 = *abi.ConvertType(out[1], new(Generation)).(*Generation)
	return *outstruct, nil
}

// PackPendingSuccessor is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xed60ed63.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function pendingSuccessor() view returns(((uint64,address,address,address,uint64,bytes32,uint16,uint64),uint64))
func (iLedgerDirectory *ILedgerDirectory) PackPendingSuccessor() []byte {
	enc, err := iLedgerDirectory.abi.Pack("pendingSuccessor")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPendingSuccessor is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xed60ed63.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function pendingSuccessor() view returns(((uint64,address,address,address,uint64,bytes32,uint16,uint64),uint64))
func (iLedgerDirectory *ILedgerDirectory) TryPackPendingSuccessor() ([]byte, error) {
	return iLedgerDirectory.abi.Pack("pendingSuccessor")
}

// UnpackPendingSuccessor is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xed60ed63.
//
// Solidity: function pendingSuccessor() view returns(((uint64,address,address,address,uint64,bytes32,uint16,uint64),uint64))
func (iLedgerDirectory *ILedgerDirectory) UnpackPendingSuccessor(data []byte) (PendingGeneration, error) {
	out, err := iLedgerDirectory.abi.Unpack("pendingSuccessor", data)
	if err != nil {
		return *new(PendingGeneration), err
	}
	out0 := *abi.ConvertType(out[0], new(PendingGeneration)).(*PendingGeneration)
	return out0, nil
}

// PackProposeForeignSuccessor is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb209b04a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function proposeForeignSuccessor(uint64 chainId, address ledger, address court, address withholding, uint16 scheme) returns()
func (iLedgerDirectory *ILedgerDirectory) PackProposeForeignSuccessor(chainId uint64, ledger common.Address, court common.Address, withholding common.Address, scheme uint16) []byte {
	enc, err := iLedgerDirectory.abi.Pack("proposeForeignSuccessor", chainId, ledger, court, withholding, scheme)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackProposeForeignSuccessor is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb209b04a.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function proposeForeignSuccessor(uint64 chainId, address ledger, address court, address withholding, uint16 scheme) returns()
func (iLedgerDirectory *ILedgerDirectory) TryPackProposeForeignSuccessor(chainId uint64, ledger common.Address, court common.Address, withholding common.Address, scheme uint16) ([]byte, error) {
	return iLedgerDirectory.abi.Pack("proposeForeignSuccessor", chainId, ledger, court, withholding, scheme)
}

// PackRegisterSuccessor is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1949190f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function registerSuccessor(address ledger, address court, address withholding) returns(uint256)
func (iLedgerDirectory *ILedgerDirectory) PackRegisterSuccessor(ledger common.Address, court common.Address, withholding common.Address) []byte {
	enc, err := iLedgerDirectory.abi.Pack("registerSuccessor", ledger, court, withholding)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRegisterSuccessor is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1949190f.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function registerSuccessor(address ledger, address court, address withholding) returns(uint256)
func (iLedgerDirectory *ILedgerDirectory) TryPackRegisterSuccessor(ledger common.Address, court common.Address, withholding common.Address) ([]byte, error) {
	return iLedgerDirectory.abi.Pack("registerSuccessor", ledger, court, withholding)
}

// UnpackRegisterSuccessor is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x1949190f.
//
// Solidity: function registerSuccessor(address ledger, address court, address withholding) returns(uint256)
func (iLedgerDirectory *ILedgerDirectory) UnpackRegisterSuccessor(data []byte) (*big.Int, error) {
	out, err := iLedgerDirectory.abi.Unpack("registerSuccessor", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// ILedgerDirectoryForeignSuccessorCancelled represents a ForeignSuccessorCancelled event raised by the ILedgerDirectory contract.
type ILedgerDirectoryForeignSuccessorCancelled struct {
	ChainId uint64
	Ledger  common.Address
	Raw     *types.Log // Blockchain specific contextual infos
}

const ILedgerDirectoryForeignSuccessorCancelledEventName = "ForeignSuccessorCancelled"

// ContractEventName returns the user-defined event name.
func (ILedgerDirectoryForeignSuccessorCancelled) ContractEventName() string {
	return ILedgerDirectoryForeignSuccessorCancelledEventName
}

// UnpackForeignSuccessorCancelledEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event ForeignSuccessorCancelled(uint64 indexed chainId, address indexed ledger)
func (iLedgerDirectory *ILedgerDirectory) UnpackForeignSuccessorCancelledEvent(log *types.Log) (*ILedgerDirectoryForeignSuccessorCancelled, error) {
	event := "ForeignSuccessorCancelled"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != iLedgerDirectory.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(ILedgerDirectoryForeignSuccessorCancelled)
	if len(log.Data) > 0 {
		if err := iLedgerDirectory.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iLedgerDirectory.abi.Events[event].Inputs {
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

// ILedgerDirectoryForeignSuccessorProposed represents a ForeignSuccessorProposed event raised by the ILedgerDirectory contract.
type ILedgerDirectoryForeignSuccessorProposed struct {
	ChainId     uint64
	Ledger      common.Address
	Court       common.Address
	Withholding common.Address
	Scheme      uint16
	ActivatesAt uint64
	Raw         *types.Log // Blockchain specific contextual infos
}

const ILedgerDirectoryForeignSuccessorProposedEventName = "ForeignSuccessorProposed"

// ContractEventName returns the user-defined event name.
func (ILedgerDirectoryForeignSuccessorProposed) ContractEventName() string {
	return ILedgerDirectoryForeignSuccessorProposedEventName
}

// UnpackForeignSuccessorProposedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event ForeignSuccessorProposed(uint64 indexed chainId, address indexed ledger, address court, address withholding, uint16 scheme, uint64 activatesAt)
func (iLedgerDirectory *ILedgerDirectory) UnpackForeignSuccessorProposedEvent(log *types.Log) (*ILedgerDirectoryForeignSuccessorProposed, error) {
	event := "ForeignSuccessorProposed"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != iLedgerDirectory.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(ILedgerDirectoryForeignSuccessorProposed)
	if len(log.Data) > 0 {
		if err := iLedgerDirectory.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iLedgerDirectory.abi.Events[event].Inputs {
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

// ILedgerDirectoryGenerationRegistered represents a GenerationRegistered event raised by the ILedgerDirectory contract.
type ILedgerDirectoryGenerationRegistered struct {
	Index         *big.Int
	ChainId       uint64
	Ledger        common.Address
	Court         common.Address
	Withholding   common.Address
	GenesisEpoch  uint64
	GenesisDigest [32]byte
	Scheme        uint16
	Raw           *types.Log // Blockchain specific contextual infos
}

const ILedgerDirectoryGenerationRegisteredEventName = "GenerationRegistered"

// ContractEventName returns the user-defined event name.
func (ILedgerDirectoryGenerationRegistered) ContractEventName() string {
	return ILedgerDirectoryGenerationRegisteredEventName
}

// UnpackGenerationRegisteredEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event GenerationRegistered(uint256 indexed index, uint64 indexed chainId, address indexed ledger, address court, address withholding, uint64 genesisEpoch, bytes32 genesisDigest, uint16 scheme)
func (iLedgerDirectory *ILedgerDirectory) UnpackGenerationRegisteredEvent(log *types.Log) (*ILedgerDirectoryGenerationRegistered, error) {
	event := "GenerationRegistered"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != iLedgerDirectory.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(ILedgerDirectoryGenerationRegistered)
	if len(log.Data) > 0 {
		if err := iLedgerDirectory.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iLedgerDirectory.abi.Events[event].Inputs {
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
func (iLedgerDirectory *ILedgerDirectory) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], iLedgerDirectory.abi.Errors["CourtOfAnotherLedger"].ID.Bytes()[:4]) {
		return iLedgerDirectory.UnpackCourtOfAnotherLedgerError(raw[4:])
	}
	if bytes.Equal(raw[:4], iLedgerDirectory.abi.Errors["DirectoryMoved"].ID.Bytes()[:4]) {
		return iLedgerDirectory.UnpackDirectoryMovedError(raw[4:])
	}
	if bytes.Equal(raw[:4], iLedgerDirectory.abi.Errors["EpochBeforeDirectory"].ID.Bytes()[:4]) {
		return iLedgerDirectory.UnpackEpochBeforeDirectoryError(raw[4:])
	}
	if bytes.Equal(raw[:4], iLedgerDirectory.abi.Errors["GenesisMismatch"].ID.Bytes()[:4]) {
		return iLedgerDirectory.UnpackGenesisMismatchError(raw[4:])
	}
	if bytes.Equal(raw[:4], iLedgerDirectory.abi.Errors["NoSuccessorPending"].ID.Bytes()[:4]) {
		return iLedgerDirectory.UnpackNoSuccessorPendingError(raw[4:])
	}
	if bytes.Equal(raw[:4], iLedgerDirectory.abi.Errors["NotALedger"].ID.Bytes()[:4]) {
		return iLedgerDirectory.UnpackNotALedgerError(raw[4:])
	}
	if bytes.Equal(raw[:4], iLedgerDirectory.abi.Errors["NotAStatementCourt"].ID.Bytes()[:4]) {
		return iLedgerDirectory.UnpackNotAStatementCourtError(raw[4:])
	}
	if bytes.Equal(raw[:4], iLedgerDirectory.abi.Errors["NotAWithholdingCourt"].ID.Bytes()[:4]) {
		return iLedgerDirectory.UnpackNotAWithholdingCourtError(raw[4:])
	}
	if bytes.Equal(raw[:4], iLedgerDirectory.abi.Errors["NotForeign"].ID.Bytes()[:4]) {
		return iLedgerDirectory.UnpackNotForeignError(raw[4:])
	}
	if bytes.Equal(raw[:4], iLedgerDirectory.abi.Errors["PredecessorNotSealed"].ID.Bytes()[:4]) {
		return iLedgerDirectory.UnpackPredecessorNotSealedError(raw[4:])
	}
	if bytes.Equal(raw[:4], iLedgerDirectory.abi.Errors["SuccessorPending"].ID.Bytes()[:4]) {
		return iLedgerDirectory.UnpackSuccessorPendingError(raw[4:])
	}
	if bytes.Equal(raw[:4], iLedgerDirectory.abi.Errors["SuccessorSealed"].ID.Bytes()[:4]) {
		return iLedgerDirectory.UnpackSuccessorSealedError(raw[4:])
	}
	if bytes.Equal(raw[:4], iLedgerDirectory.abi.Errors["SuccessorTooEarly"].ID.Bytes()[:4]) {
		return iLedgerDirectory.UnpackSuccessorTooEarlyError(raw[4:])
	}
	if bytes.Equal(raw[:4], iLedgerDirectory.abi.Errors["UnknownGeneration"].ID.Bytes()[:4]) {
		return iLedgerDirectory.UnpackUnknownGenerationError(raw[4:])
	}
	if bytes.Equal(raw[:4], iLedgerDirectory.abi.Errors["WithholdingOfAnotherCourt"].ID.Bytes()[:4]) {
		return iLedgerDirectory.UnpackWithholdingOfAnotherCourtError(raw[4:])
	}
	if bytes.Equal(raw[:4], iLedgerDirectory.abi.Errors["ZeroAddress"].ID.Bytes()[:4]) {
		return iLedgerDirectory.UnpackZeroAddressError(raw[4:])
	}
	if bytes.Equal(raw[:4], iLedgerDirectory.abi.Errors["ZeroDelay"].ID.Bytes()[:4]) {
		return iLedgerDirectory.UnpackZeroDelayError(raw[4:])
	}
	if bytes.Equal(raw[:4], iLedgerDirectory.abi.Errors["ZeroScheme"].ID.Bytes()[:4]) {
		return iLedgerDirectory.UnpackZeroSchemeError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// ILedgerDirectoryCourtOfAnotherLedger represents a CourtOfAnotherLedger error raised by the ILedgerDirectory contract.
type ILedgerDirectoryCourtOfAnotherLedger struct {
	Court  common.Address
	Ledger common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error CourtOfAnotherLedger(address court, address ledger)
func ILedgerDirectoryCourtOfAnotherLedgerErrorID() common.Hash {
	return common.HexToHash("0x8f6b194dbeeaffcf237320ac6f7e0fc8c84615fe1f375c836c737f8a7e59f58a")
}

// UnpackCourtOfAnotherLedgerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error CourtOfAnotherLedger(address court, address ledger)
func (iLedgerDirectory *ILedgerDirectory) UnpackCourtOfAnotherLedgerError(raw []byte) (*ILedgerDirectoryCourtOfAnotherLedger, error) {
	out := new(ILedgerDirectoryCourtOfAnotherLedger)
	if err := iLedgerDirectory.abi.UnpackIntoInterface(out, "CourtOfAnotherLedger", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ILedgerDirectoryDirectoryMoved represents a DirectoryMoved error raised by the ILedgerDirectory contract.
type ILedgerDirectoryDirectoryMoved struct {
	ChainId uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error DirectoryMoved(uint64 chainId)
func ILedgerDirectoryDirectoryMovedErrorID() common.Hash {
	return common.HexToHash("0x0247d34bc24396ffb5520e6fd37d229b1d6bcb58192cfcec9b579b6ad43b0107")
}

// UnpackDirectoryMovedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error DirectoryMoved(uint64 chainId)
func (iLedgerDirectory *ILedgerDirectory) UnpackDirectoryMovedError(raw []byte) (*ILedgerDirectoryDirectoryMoved, error) {
	out := new(ILedgerDirectoryDirectoryMoved)
	if err := iLedgerDirectory.abi.UnpackIntoInterface(out, "DirectoryMoved", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ILedgerDirectoryEpochBeforeDirectory represents a EpochBeforeDirectory error raised by the ILedgerDirectory contract.
type ILedgerDirectoryEpochBeforeDirectory struct {
	Epoch        uint64
	GenesisEpoch uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error EpochBeforeDirectory(uint64 epoch, uint64 genesisEpoch)
func ILedgerDirectoryEpochBeforeDirectoryErrorID() common.Hash {
	return common.HexToHash("0x00ecd8bff74bf151c4ab76b41042988985631b9f910ffaaa474f5bd68242bd32")
}

// UnpackEpochBeforeDirectoryError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error EpochBeforeDirectory(uint64 epoch, uint64 genesisEpoch)
func (iLedgerDirectory *ILedgerDirectory) UnpackEpochBeforeDirectoryError(raw []byte) (*ILedgerDirectoryEpochBeforeDirectory, error) {
	out := new(ILedgerDirectoryEpochBeforeDirectory)
	if err := iLedgerDirectory.abi.UnpackIntoInterface(out, "EpochBeforeDirectory", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ILedgerDirectoryGenesisMismatch represents a GenesisMismatch error raised by the ILedgerDirectory contract.
type ILedgerDirectoryGenesisMismatch struct {
	ExpectedEpoch  uint64
	ExpectedDigest [32]byte
	GivenEpoch     uint64
	GivenDigest    [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error GenesisMismatch(uint64 expectedEpoch, bytes32 expectedDigest, uint64 givenEpoch, bytes32 givenDigest)
func ILedgerDirectoryGenesisMismatchErrorID() common.Hash {
	return common.HexToHash("0x1a73749707f69ccdcbbeeb2b66c5256863086268a09e247754e79a6118d8a173")
}

// UnpackGenesisMismatchError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error GenesisMismatch(uint64 expectedEpoch, bytes32 expectedDigest, uint64 givenEpoch, bytes32 givenDigest)
func (iLedgerDirectory *ILedgerDirectory) UnpackGenesisMismatchError(raw []byte) (*ILedgerDirectoryGenesisMismatch, error) {
	out := new(ILedgerDirectoryGenesisMismatch)
	if err := iLedgerDirectory.abi.UnpackIntoInterface(out, "GenesisMismatch", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ILedgerDirectoryNoSuccessorPending represents a NoSuccessorPending error raised by the ILedgerDirectory contract.
type ILedgerDirectoryNoSuccessorPending struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NoSuccessorPending()
func ILedgerDirectoryNoSuccessorPendingErrorID() common.Hash {
	return common.HexToHash("0xc9fe10b5cd27de39964dc50095bfa0eb7a3f8842427f606f7997a9f43572de73")
}

// UnpackNoSuccessorPendingError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NoSuccessorPending()
func (iLedgerDirectory *ILedgerDirectory) UnpackNoSuccessorPendingError(raw []byte) (*ILedgerDirectoryNoSuccessorPending, error) {
	out := new(ILedgerDirectoryNoSuccessorPending)
	if err := iLedgerDirectory.abi.UnpackIntoInterface(out, "NoSuccessorPending", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ILedgerDirectoryNotALedger represents a NotALedger error raised by the ILedgerDirectory contract.
type ILedgerDirectoryNotALedger struct {
	Ledger common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotALedger(address ledger)
func ILedgerDirectoryNotALedgerErrorID() common.Hash {
	return common.HexToHash("0x44a5b4b2f2697952fce85b89e8be359271d0792b4f0be489d664b4771187455b")
}

// UnpackNotALedgerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotALedger(address ledger)
func (iLedgerDirectory *ILedgerDirectory) UnpackNotALedgerError(raw []byte) (*ILedgerDirectoryNotALedger, error) {
	out := new(ILedgerDirectoryNotALedger)
	if err := iLedgerDirectory.abi.UnpackIntoInterface(out, "NotALedger", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ILedgerDirectoryNotAStatementCourt represents a NotAStatementCourt error raised by the ILedgerDirectory contract.
type ILedgerDirectoryNotAStatementCourt struct {
	Court common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotAStatementCourt(address court)
func ILedgerDirectoryNotAStatementCourtErrorID() common.Hash {
	return common.HexToHash("0x8e35ed64fa846e83af245ae076efcecd211c173d4d90de52d3daab09ef787aff")
}

// UnpackNotAStatementCourtError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotAStatementCourt(address court)
func (iLedgerDirectory *ILedgerDirectory) UnpackNotAStatementCourtError(raw []byte) (*ILedgerDirectoryNotAStatementCourt, error) {
	out := new(ILedgerDirectoryNotAStatementCourt)
	if err := iLedgerDirectory.abi.UnpackIntoInterface(out, "NotAStatementCourt", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ILedgerDirectoryNotAWithholdingCourt represents a NotAWithholdingCourt error raised by the ILedgerDirectory contract.
type ILedgerDirectoryNotAWithholdingCourt struct {
	Withholding common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotAWithholdingCourt(address withholding)
func ILedgerDirectoryNotAWithholdingCourtErrorID() common.Hash {
	return common.HexToHash("0x6e9637d8ce542711d9e0648e507cdb2eaecb6b823047364c9c5ed9e719d00bac")
}

// UnpackNotAWithholdingCourtError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotAWithholdingCourt(address withholding)
func (iLedgerDirectory *ILedgerDirectory) UnpackNotAWithholdingCourtError(raw []byte) (*ILedgerDirectoryNotAWithholdingCourt, error) {
	out := new(ILedgerDirectoryNotAWithholdingCourt)
	if err := iLedgerDirectory.abi.UnpackIntoInterface(out, "NotAWithholdingCourt", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ILedgerDirectoryNotForeign represents a NotForeign error raised by the ILedgerDirectory contract.
type ILedgerDirectoryNotForeign struct {
	ChainId uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotForeign(uint64 chainId)
func ILedgerDirectoryNotForeignErrorID() common.Hash {
	return common.HexToHash("0x66b10aaf4d91050b3af5f7ec04493d8d816cd0c4b5b1578abb8949dc8694c692")
}

// UnpackNotForeignError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotForeign(uint64 chainId)
func (iLedgerDirectory *ILedgerDirectory) UnpackNotForeignError(raw []byte) (*ILedgerDirectoryNotForeign, error) {
	out := new(ILedgerDirectoryNotForeign)
	if err := iLedgerDirectory.abi.UnpackIntoInterface(out, "NotForeign", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ILedgerDirectoryPredecessorNotSealed represents a PredecessorNotSealed error raised by the ILedgerDirectory contract.
type ILedgerDirectoryPredecessorNotSealed struct {
	Ledger common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error PredecessorNotSealed(address ledger)
func ILedgerDirectoryPredecessorNotSealedErrorID() common.Hash {
	return common.HexToHash("0x92422182b3494808df1a64555a9d86cdc7396f69d34f57e489bb0a7a416a7a1c")
}

// UnpackPredecessorNotSealedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error PredecessorNotSealed(address ledger)
func (iLedgerDirectory *ILedgerDirectory) UnpackPredecessorNotSealedError(raw []byte) (*ILedgerDirectoryPredecessorNotSealed, error) {
	out := new(ILedgerDirectoryPredecessorNotSealed)
	if err := iLedgerDirectory.abi.UnpackIntoInterface(out, "PredecessorNotSealed", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ILedgerDirectorySuccessorPending represents a SuccessorPending error raised by the ILedgerDirectory contract.
type ILedgerDirectorySuccessorPending struct {
	ActivatesAt uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SuccessorPending(uint64 activatesAt)
func ILedgerDirectorySuccessorPendingErrorID() common.Hash {
	return common.HexToHash("0xbda32a82d6c520c0cd1b997f0b6b9545163906e859b28decce9df08e009a6714")
}

// UnpackSuccessorPendingError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SuccessorPending(uint64 activatesAt)
func (iLedgerDirectory *ILedgerDirectory) UnpackSuccessorPendingError(raw []byte) (*ILedgerDirectorySuccessorPending, error) {
	out := new(ILedgerDirectorySuccessorPending)
	if err := iLedgerDirectory.abi.UnpackIntoInterface(out, "SuccessorPending", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ILedgerDirectorySuccessorSealed represents a SuccessorSealed error raised by the ILedgerDirectory contract.
type ILedgerDirectorySuccessorSealed struct {
	Ledger common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SuccessorSealed(address ledger)
func ILedgerDirectorySuccessorSealedErrorID() common.Hash {
	return common.HexToHash("0xdbf4bf962732730b54b17c3a9bd9140f4f63c513815426127df2be18d57686da")
}

// UnpackSuccessorSealedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SuccessorSealed(address ledger)
func (iLedgerDirectory *ILedgerDirectory) UnpackSuccessorSealedError(raw []byte) (*ILedgerDirectorySuccessorSealed, error) {
	out := new(ILedgerDirectorySuccessorSealed)
	if err := iLedgerDirectory.abi.UnpackIntoInterface(out, "SuccessorSealed", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ILedgerDirectorySuccessorTooEarly represents a SuccessorTooEarly error raised by the ILedgerDirectory contract.
type ILedgerDirectorySuccessorTooEarly struct {
	ActivatesAt uint64
	BlockTime   uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SuccessorTooEarly(uint64 activatesAt, uint64 blockTime)
func ILedgerDirectorySuccessorTooEarlyErrorID() common.Hash {
	return common.HexToHash("0xa32346b6ab676dda3605c57ad7986777db6179ffa656057f42a0e9141cf0d3b3")
}

// UnpackSuccessorTooEarlyError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SuccessorTooEarly(uint64 activatesAt, uint64 blockTime)
func (iLedgerDirectory *ILedgerDirectory) UnpackSuccessorTooEarlyError(raw []byte) (*ILedgerDirectorySuccessorTooEarly, error) {
	out := new(ILedgerDirectorySuccessorTooEarly)
	if err := iLedgerDirectory.abi.UnpackIntoInterface(out, "SuccessorTooEarly", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ILedgerDirectoryUnknownGeneration represents a UnknownGeneration error raised by the ILedgerDirectory contract.
type ILedgerDirectoryUnknownGeneration struct {
	Index *big.Int
	Count *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error UnknownGeneration(uint256 index, uint256 count)
func ILedgerDirectoryUnknownGenerationErrorID() common.Hash {
	return common.HexToHash("0x23eb3cd78e103c67e711b97a8c18bafd80842e97bc80ee6c0c9c65b79d4c38e9")
}

// UnpackUnknownGenerationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error UnknownGeneration(uint256 index, uint256 count)
func (iLedgerDirectory *ILedgerDirectory) UnpackUnknownGenerationError(raw []byte) (*ILedgerDirectoryUnknownGeneration, error) {
	out := new(ILedgerDirectoryUnknownGeneration)
	if err := iLedgerDirectory.abi.UnpackIntoInterface(out, "UnknownGeneration", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ILedgerDirectoryWithholdingOfAnotherCourt represents a WithholdingOfAnotherCourt error raised by the ILedgerDirectory contract.
type ILedgerDirectoryWithholdingOfAnotherCourt struct {
	Withholding common.Address
	Court       common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error WithholdingOfAnotherCourt(address withholding, address court)
func ILedgerDirectoryWithholdingOfAnotherCourtErrorID() common.Hash {
	return common.HexToHash("0x971e79ac5d36b6454b83e1b48cbaee8b0a3bfb145c25170cd4b9613099fb03aa")
}

// UnpackWithholdingOfAnotherCourtError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error WithholdingOfAnotherCourt(address withholding, address court)
func (iLedgerDirectory *ILedgerDirectory) UnpackWithholdingOfAnotherCourtError(raw []byte) (*ILedgerDirectoryWithholdingOfAnotherCourt, error) {
	out := new(ILedgerDirectoryWithholdingOfAnotherCourt)
	if err := iLedgerDirectory.abi.UnpackIntoInterface(out, "WithholdingOfAnotherCourt", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ILedgerDirectoryZeroAddress represents a ZeroAddress error raised by the ILedgerDirectory contract.
type ILedgerDirectoryZeroAddress struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ZeroAddress()
func ILedgerDirectoryZeroAddressErrorID() common.Hash {
	return common.HexToHash("0xd92e233df2717d4a40030e20904abd27b68fcbeede117eaaccbbdac9618c8c73")
}

// UnpackZeroAddressError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ZeroAddress()
func (iLedgerDirectory *ILedgerDirectory) UnpackZeroAddressError(raw []byte) (*ILedgerDirectoryZeroAddress, error) {
	out := new(ILedgerDirectoryZeroAddress)
	if err := iLedgerDirectory.abi.UnpackIntoInterface(out, "ZeroAddress", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ILedgerDirectoryZeroDelay represents a ZeroDelay error raised by the ILedgerDirectory contract.
type ILedgerDirectoryZeroDelay struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ZeroDelay()
func ILedgerDirectoryZeroDelayErrorID() common.Hash {
	return common.HexToHash("0x34c2b9dab051f3b49e0b9121d90ae110f06df318864d2741799dc06f9b92cebc")
}

// UnpackZeroDelayError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ZeroDelay()
func (iLedgerDirectory *ILedgerDirectory) UnpackZeroDelayError(raw []byte) (*ILedgerDirectoryZeroDelay, error) {
	out := new(ILedgerDirectoryZeroDelay)
	if err := iLedgerDirectory.abi.UnpackIntoInterface(out, "ZeroDelay", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ILedgerDirectoryZeroScheme represents a ZeroScheme error raised by the ILedgerDirectory contract.
type ILedgerDirectoryZeroScheme struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ZeroScheme()
func ILedgerDirectoryZeroSchemeErrorID() common.Hash {
	return common.HexToHash("0xb3db53140c5c18ae31b9ec7c97dd71357e34672c7f74f824621c54e674fcf0ce")
}

// UnpackZeroSchemeError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ZeroScheme()
func (iLedgerDirectory *ILedgerDirectory) UnpackZeroSchemeError(raw []byte) (*ILedgerDirectoryZeroScheme, error) {
	out := new(ILedgerDirectoryZeroScheme)
	if err := iLedgerDirectory.abi.UnpackIntoInterface(out, "ZeroScheme", raw); err != nil {
		return nil, err
	}
	return out, nil
}
