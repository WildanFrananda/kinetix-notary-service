// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package councilinterface

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

// EpochHeader is an auto generated low-level Go binding around an user-defined struct.
type EpochHeader struct {
	Number         uint64
	PreviousDigest [32]byte
	Root           [32]byte
	Total          uint64
	LeafCount      uint64
	AsOf           uint64
}

// PendingMembership is an auto generated low-level Go binding around an user-defined struct.
type PendingMembership struct {
	Members     []common.Address
	Threshold   uint8
	ActivatesAt uint64
}

// IAnchorCouncilMetaData contains all meta data concerning the IAnchorCouncil contract.
var IAnchorCouncilMetaData = bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"activateMembership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"anchor\",\"inputs\":[{\"name\":\"header\",\"type\":\"tuple\",\"internalType\":\"structEpochHeader\",\"components\":[{\"name\":\"number\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"previousDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"total\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"leafCount\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"asOf\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"signatures\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"approvalDigest\",\"inputs\":[{\"name\":\"number\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"digest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"cancelMembership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"epochDigest\",\"inputs\":[{\"name\":\"header\",\"type\":\"tuple\",\"internalType\":\"structEpochHeader\",\"components\":[{\"name\":\"number\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"previousDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"total\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"leafCount\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"asOf\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"isMember\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ledger\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractILiabilityLedger\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"members\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"membershipChangeDelay\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingMembership\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structPendingMembership\",\"components\":[{\"name\":\"members\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"threshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"activatesAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"recoverApprover\",\"inputs\":[{\"name\":\"number\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"digest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"scheduleMembership\",\"inputs\":[{\"name\":\"nextMembers\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"nextThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"threshold\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"EpochApproved\",\"inputs\":[{\"name\":\"number\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"digest\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"approvers\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MembershipActivated\",\"inputs\":[{\"name\":\"members\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"},{\"name\":\"threshold\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MembershipCancelled\",\"inputs\":[{\"name\":\"members\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"},{\"name\":\"threshold\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MembershipScheduled\",\"inputs\":[{\"name\":\"members\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"},{\"name\":\"threshold\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"},{\"name\":\"activatesAt\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ApprovalsOutOfOrder\",\"inputs\":[{\"name\":\"previous\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"next\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"MembersNotAscending\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"MembershipChangePending\",\"inputs\":[{\"name\":\"activatesAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"MembershipChangeTooEarly\",\"inputs\":[{\"name\":\"activatesAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"blockTime\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"NoMembershipChangePending\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotALedger\",\"inputs\":[{\"name\":\"ledger\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"NotAMember\",\"inputs\":[{\"name\":\"recovered\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"NotEnoughApprovals\",\"inputs\":[{\"name\":\"given\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"threshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"SameMembership\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ThresholdTooLow\",\"inputs\":[{\"name\":\"threshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"minimum\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"TooFewMembers\",\"inputs\":[{\"name\":\"count\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"threshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"TooManyMembers\",\"inputs\":[{\"name\":\"count\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maximum\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ZeroDelay\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroMember\",\"inputs\":[]}]",
	ID:  "IAnchorCouncil",
}

// IAnchorCouncil is an auto generated Go binding around an Ethereum contract.
type IAnchorCouncil struct {
	abi abi.ABI
}

// GetABI returns the ABI associated with this contract binding.
func (c *IAnchorCouncil) GetABI() abi.ABI {
	return c.abi
}

// NewIAnchorCouncil creates a new instance of IAnchorCouncil.
func NewIAnchorCouncil() *IAnchorCouncil {
	parsed, err := IAnchorCouncilMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &IAnchorCouncil{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *IAnchorCouncil) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackActivateMembership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x82b9e10a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function activateMembership() returns()
func (iAnchorCouncil *IAnchorCouncil) PackActivateMembership() []byte {
	enc, err := iAnchorCouncil.abi.Pack("activateMembership")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackActivateMembership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x82b9e10a.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function activateMembership() returns()
func (iAnchorCouncil *IAnchorCouncil) TryPackActivateMembership() ([]byte, error) {
	return iAnchorCouncil.abi.Pack("activateMembership")
}

// PackAnchor is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb9e61120.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function anchor((uint64,bytes32,bytes32,uint64,uint64,uint64) header, bytes[] signatures) returns(bytes32)
func (iAnchorCouncil *IAnchorCouncil) PackAnchor(header EpochHeader, signatures [][]byte) []byte {
	enc, err := iAnchorCouncil.abi.Pack("anchor", header, signatures)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAnchor is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb9e61120.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function anchor((uint64,bytes32,bytes32,uint64,uint64,uint64) header, bytes[] signatures) returns(bytes32)
func (iAnchorCouncil *IAnchorCouncil) TryPackAnchor(header EpochHeader, signatures [][]byte) ([]byte, error) {
	return iAnchorCouncil.abi.Pack("anchor", header, signatures)
}

// UnpackAnchor is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb9e61120.
//
// Solidity: function anchor((uint64,bytes32,bytes32,uint64,uint64,uint64) header, bytes[] signatures) returns(bytes32)
func (iAnchorCouncil *IAnchorCouncil) UnpackAnchor(data []byte) ([32]byte, error) {
	out, err := iAnchorCouncil.abi.Unpack("anchor", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackApprovalDigest is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0331e6cc.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function approvalDigest(uint64 number, bytes32 digest) view returns(bytes32)
func (iAnchorCouncil *IAnchorCouncil) PackApprovalDigest(number uint64, digest [32]byte) []byte {
	enc, err := iAnchorCouncil.abi.Pack("approvalDigest", number, digest)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackApprovalDigest is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0331e6cc.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function approvalDigest(uint64 number, bytes32 digest) view returns(bytes32)
func (iAnchorCouncil *IAnchorCouncil) TryPackApprovalDigest(number uint64, digest [32]byte) ([]byte, error) {
	return iAnchorCouncil.abi.Pack("approvalDigest", number, digest)
}

// UnpackApprovalDigest is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x0331e6cc.
//
// Solidity: function approvalDigest(uint64 number, bytes32 digest) view returns(bytes32)
func (iAnchorCouncil *IAnchorCouncil) UnpackApprovalDigest(data []byte) ([32]byte, error) {
	out, err := iAnchorCouncil.abi.Unpack("approvalDigest", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackCancelMembership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x082e847a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function cancelMembership() returns()
func (iAnchorCouncil *IAnchorCouncil) PackCancelMembership() []byte {
	enc, err := iAnchorCouncil.abi.Pack("cancelMembership")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCancelMembership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x082e847a.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function cancelMembership() returns()
func (iAnchorCouncil *IAnchorCouncil) TryPackCancelMembership() ([]byte, error) {
	return iAnchorCouncil.abi.Pack("cancelMembership")
}

// PackEpochDigest is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x91c28c71.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function epochDigest((uint64,bytes32,bytes32,uint64,uint64,uint64) header) pure returns(bytes32)
func (iAnchorCouncil *IAnchorCouncil) PackEpochDigest(header EpochHeader) []byte {
	enc, err := iAnchorCouncil.abi.Pack("epochDigest", header)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackEpochDigest is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x91c28c71.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function epochDigest((uint64,bytes32,bytes32,uint64,uint64,uint64) header) pure returns(bytes32)
func (iAnchorCouncil *IAnchorCouncil) TryPackEpochDigest(header EpochHeader) ([]byte, error) {
	return iAnchorCouncil.abi.Pack("epochDigest", header)
}

// UnpackEpochDigest is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x91c28c71.
//
// Solidity: function epochDigest((uint64,bytes32,bytes32,uint64,uint64,uint64) header) pure returns(bytes32)
func (iAnchorCouncil *IAnchorCouncil) UnpackEpochDigest(data []byte) ([32]byte, error) {
	out, err := iAnchorCouncil.abi.Unpack("epochDigest", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackIsMember is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa230c524.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function isMember(address account) view returns(bool)
func (iAnchorCouncil *IAnchorCouncil) PackIsMember(account common.Address) []byte {
	enc, err := iAnchorCouncil.abi.Pack("isMember", account)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackIsMember is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa230c524.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function isMember(address account) view returns(bool)
func (iAnchorCouncil *IAnchorCouncil) TryPackIsMember(account common.Address) ([]byte, error) {
	return iAnchorCouncil.abi.Pack("isMember", account)
}

// UnpackIsMember is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa230c524.
//
// Solidity: function isMember(address account) view returns(bool)
func (iAnchorCouncil *IAnchorCouncil) UnpackIsMember(data []byte) (bool, error) {
	out, err := iAnchorCouncil.abi.Unpack("isMember", data)
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
func (iAnchorCouncil *IAnchorCouncil) PackLedger() []byte {
	enc, err := iAnchorCouncil.abi.Pack("ledger")
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
func (iAnchorCouncil *IAnchorCouncil) TryPackLedger() ([]byte, error) {
	return iAnchorCouncil.abi.Pack("ledger")
}

// UnpackLedger is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x56397c35.
//
// Solidity: function ledger() view returns(address)
func (iAnchorCouncil *IAnchorCouncil) UnpackLedger(data []byte) (common.Address, error) {
	out, err := iAnchorCouncil.abi.Unpack("ledger", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackMembers is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbdd4d18d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function members() view returns(address[])
func (iAnchorCouncil *IAnchorCouncil) PackMembers() []byte {
	enc, err := iAnchorCouncil.abi.Pack("members")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMembers is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbdd4d18d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function members() view returns(address[])
func (iAnchorCouncil *IAnchorCouncil) TryPackMembers() ([]byte, error) {
	return iAnchorCouncil.abi.Pack("members")
}

// UnpackMembers is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xbdd4d18d.
//
// Solidity: function members() view returns(address[])
func (iAnchorCouncil *IAnchorCouncil) UnpackMembers(data []byte) ([]common.Address, error) {
	out, err := iAnchorCouncil.abi.Unpack("members", data)
	if err != nil {
		return *new([]common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)
	return out0, nil
}

// PackMembershipChangeDelay is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x319f5145.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function membershipChangeDelay() view returns(uint64)
func (iAnchorCouncil *IAnchorCouncil) PackMembershipChangeDelay() []byte {
	enc, err := iAnchorCouncil.abi.Pack("membershipChangeDelay")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMembershipChangeDelay is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x319f5145.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function membershipChangeDelay() view returns(uint64)
func (iAnchorCouncil *IAnchorCouncil) TryPackMembershipChangeDelay() ([]byte, error) {
	return iAnchorCouncil.abi.Pack("membershipChangeDelay")
}

// UnpackMembershipChangeDelay is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x319f5145.
//
// Solidity: function membershipChangeDelay() view returns(uint64)
func (iAnchorCouncil *IAnchorCouncil) UnpackMembershipChangeDelay(data []byte) (uint64, error) {
	out, err := iAnchorCouncil.abi.Unpack("membershipChangeDelay", data)
	if err != nil {
		return *new(uint64), err
	}
	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)
	return out0, nil
}

// PackPendingMembership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x51136bd9.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function pendingMembership() view returns((address[],uint8,uint64))
func (iAnchorCouncil *IAnchorCouncil) PackPendingMembership() []byte {
	enc, err := iAnchorCouncil.abi.Pack("pendingMembership")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPendingMembership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x51136bd9.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function pendingMembership() view returns((address[],uint8,uint64))
func (iAnchorCouncil *IAnchorCouncil) TryPackPendingMembership() ([]byte, error) {
	return iAnchorCouncil.abi.Pack("pendingMembership")
}

// UnpackPendingMembership is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x51136bd9.
//
// Solidity: function pendingMembership() view returns((address[],uint8,uint64))
func (iAnchorCouncil *IAnchorCouncil) UnpackPendingMembership(data []byte) (PendingMembership, error) {
	out, err := iAnchorCouncil.abi.Unpack("pendingMembership", data)
	if err != nil {
		return *new(PendingMembership), err
	}
	out0 := *abi.ConvertType(out[0], new(PendingMembership)).(*PendingMembership)
	return out0, nil
}

// PackRecoverApprover is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xec610c53.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function recoverApprover(uint64 number, bytes32 digest, bytes signature) view returns(address)
func (iAnchorCouncil *IAnchorCouncil) PackRecoverApprover(number uint64, digest [32]byte, signature []byte) []byte {
	enc, err := iAnchorCouncil.abi.Pack("recoverApprover", number, digest, signature)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRecoverApprover is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xec610c53.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function recoverApprover(uint64 number, bytes32 digest, bytes signature) view returns(address)
func (iAnchorCouncil *IAnchorCouncil) TryPackRecoverApprover(number uint64, digest [32]byte, signature []byte) ([]byte, error) {
	return iAnchorCouncil.abi.Pack("recoverApprover", number, digest, signature)
}

// UnpackRecoverApprover is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xec610c53.
//
// Solidity: function recoverApprover(uint64 number, bytes32 digest, bytes signature) view returns(address)
func (iAnchorCouncil *IAnchorCouncil) UnpackRecoverApprover(data []byte) (common.Address, error) {
	out, err := iAnchorCouncil.abi.Unpack("recoverApprover", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackScheduleMembership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb209c6f8.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function scheduleMembership(address[] nextMembers, uint8 nextThreshold) returns()
func (iAnchorCouncil *IAnchorCouncil) PackScheduleMembership(nextMembers []common.Address, nextThreshold uint8) []byte {
	enc, err := iAnchorCouncil.abi.Pack("scheduleMembership", nextMembers, nextThreshold)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackScheduleMembership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb209c6f8.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function scheduleMembership(address[] nextMembers, uint8 nextThreshold) returns()
func (iAnchorCouncil *IAnchorCouncil) TryPackScheduleMembership(nextMembers []common.Address, nextThreshold uint8) ([]byte, error) {
	return iAnchorCouncil.abi.Pack("scheduleMembership", nextMembers, nextThreshold)
}

// PackThreshold is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x42cde4e8.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function threshold() view returns(uint8)
func (iAnchorCouncil *IAnchorCouncil) PackThreshold() []byte {
	enc, err := iAnchorCouncil.abi.Pack("threshold")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackThreshold is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x42cde4e8.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function threshold() view returns(uint8)
func (iAnchorCouncil *IAnchorCouncil) TryPackThreshold() ([]byte, error) {
	return iAnchorCouncil.abi.Pack("threshold")
}

// UnpackThreshold is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x42cde4e8.
//
// Solidity: function threshold() view returns(uint8)
func (iAnchorCouncil *IAnchorCouncil) UnpackThreshold(data []byte) (uint8, error) {
	out, err := iAnchorCouncil.abi.Unpack("threshold", data)
	if err != nil {
		return *new(uint8), err
	}
	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)
	return out0, nil
}

// IAnchorCouncilEpochApproved represents a EpochApproved event raised by the IAnchorCouncil contract.
type IAnchorCouncilEpochApproved struct {
	Number    uint64
	Digest    [32]byte
	Approvers []common.Address
	Raw       *types.Log // Blockchain specific contextual infos
}

const IAnchorCouncilEpochApprovedEventName = "EpochApproved"

// ContractEventName returns the user-defined event name.
func (IAnchorCouncilEpochApproved) ContractEventName() string {
	return IAnchorCouncilEpochApprovedEventName
}

// UnpackEpochApprovedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event EpochApproved(uint64 indexed number, bytes32 indexed digest, address[] approvers)
func (iAnchorCouncil *IAnchorCouncil) UnpackEpochApprovedEvent(log *types.Log) (*IAnchorCouncilEpochApproved, error) {
	event := "EpochApproved"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != iAnchorCouncil.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(IAnchorCouncilEpochApproved)
	if len(log.Data) > 0 {
		if err := iAnchorCouncil.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iAnchorCouncil.abi.Events[event].Inputs {
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

// IAnchorCouncilMembershipActivated represents a MembershipActivated event raised by the IAnchorCouncil contract.
type IAnchorCouncilMembershipActivated struct {
	Members   []common.Address
	Threshold uint8
	Raw       *types.Log // Blockchain specific contextual infos
}

const IAnchorCouncilMembershipActivatedEventName = "MembershipActivated"

// ContractEventName returns the user-defined event name.
func (IAnchorCouncilMembershipActivated) ContractEventName() string {
	return IAnchorCouncilMembershipActivatedEventName
}

// UnpackMembershipActivatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event MembershipActivated(address[] members, uint8 threshold)
func (iAnchorCouncil *IAnchorCouncil) UnpackMembershipActivatedEvent(log *types.Log) (*IAnchorCouncilMembershipActivated, error) {
	event := "MembershipActivated"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != iAnchorCouncil.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(IAnchorCouncilMembershipActivated)
	if len(log.Data) > 0 {
		if err := iAnchorCouncil.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iAnchorCouncil.abi.Events[event].Inputs {
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

// IAnchorCouncilMembershipCancelled represents a MembershipCancelled event raised by the IAnchorCouncil contract.
type IAnchorCouncilMembershipCancelled struct {
	Members   []common.Address
	Threshold uint8
	Raw       *types.Log // Blockchain specific contextual infos
}

const IAnchorCouncilMembershipCancelledEventName = "MembershipCancelled"

// ContractEventName returns the user-defined event name.
func (IAnchorCouncilMembershipCancelled) ContractEventName() string {
	return IAnchorCouncilMembershipCancelledEventName
}

// UnpackMembershipCancelledEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event MembershipCancelled(address[] members, uint8 threshold)
func (iAnchorCouncil *IAnchorCouncil) UnpackMembershipCancelledEvent(log *types.Log) (*IAnchorCouncilMembershipCancelled, error) {
	event := "MembershipCancelled"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != iAnchorCouncil.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(IAnchorCouncilMembershipCancelled)
	if len(log.Data) > 0 {
		if err := iAnchorCouncil.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iAnchorCouncil.abi.Events[event].Inputs {
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

// IAnchorCouncilMembershipScheduled represents a MembershipScheduled event raised by the IAnchorCouncil contract.
type IAnchorCouncilMembershipScheduled struct {
	Members     []common.Address
	Threshold   uint8
	ActivatesAt uint64
	Raw         *types.Log // Blockchain specific contextual infos
}

const IAnchorCouncilMembershipScheduledEventName = "MembershipScheduled"

// ContractEventName returns the user-defined event name.
func (IAnchorCouncilMembershipScheduled) ContractEventName() string {
	return IAnchorCouncilMembershipScheduledEventName
}

// UnpackMembershipScheduledEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event MembershipScheduled(address[] members, uint8 threshold, uint64 activatesAt)
func (iAnchorCouncil *IAnchorCouncil) UnpackMembershipScheduledEvent(log *types.Log) (*IAnchorCouncilMembershipScheduled, error) {
	event := "MembershipScheduled"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != iAnchorCouncil.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(IAnchorCouncilMembershipScheduled)
	if len(log.Data) > 0 {
		if err := iAnchorCouncil.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iAnchorCouncil.abi.Events[event].Inputs {
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
func (iAnchorCouncil *IAnchorCouncil) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], iAnchorCouncil.abi.Errors["ApprovalsOutOfOrder"].ID.Bytes()[:4]) {
		return iAnchorCouncil.UnpackApprovalsOutOfOrderError(raw[4:])
	}
	if bytes.Equal(raw[:4], iAnchorCouncil.abi.Errors["MembersNotAscending"].ID.Bytes()[:4]) {
		return iAnchorCouncil.UnpackMembersNotAscendingError(raw[4:])
	}
	if bytes.Equal(raw[:4], iAnchorCouncil.abi.Errors["MembershipChangePending"].ID.Bytes()[:4]) {
		return iAnchorCouncil.UnpackMembershipChangePendingError(raw[4:])
	}
	if bytes.Equal(raw[:4], iAnchorCouncil.abi.Errors["MembershipChangeTooEarly"].ID.Bytes()[:4]) {
		return iAnchorCouncil.UnpackMembershipChangeTooEarlyError(raw[4:])
	}
	if bytes.Equal(raw[:4], iAnchorCouncil.abi.Errors["NoMembershipChangePending"].ID.Bytes()[:4]) {
		return iAnchorCouncil.UnpackNoMembershipChangePendingError(raw[4:])
	}
	if bytes.Equal(raw[:4], iAnchorCouncil.abi.Errors["NotALedger"].ID.Bytes()[:4]) {
		return iAnchorCouncil.UnpackNotALedgerError(raw[4:])
	}
	if bytes.Equal(raw[:4], iAnchorCouncil.abi.Errors["NotAMember"].ID.Bytes()[:4]) {
		return iAnchorCouncil.UnpackNotAMemberError(raw[4:])
	}
	if bytes.Equal(raw[:4], iAnchorCouncil.abi.Errors["NotEnoughApprovals"].ID.Bytes()[:4]) {
		return iAnchorCouncil.UnpackNotEnoughApprovalsError(raw[4:])
	}
	if bytes.Equal(raw[:4], iAnchorCouncil.abi.Errors["SameMembership"].ID.Bytes()[:4]) {
		return iAnchorCouncil.UnpackSameMembershipError(raw[4:])
	}
	if bytes.Equal(raw[:4], iAnchorCouncil.abi.Errors["ThresholdTooLow"].ID.Bytes()[:4]) {
		return iAnchorCouncil.UnpackThresholdTooLowError(raw[4:])
	}
	if bytes.Equal(raw[:4], iAnchorCouncil.abi.Errors["TooFewMembers"].ID.Bytes()[:4]) {
		return iAnchorCouncil.UnpackTooFewMembersError(raw[4:])
	}
	if bytes.Equal(raw[:4], iAnchorCouncil.abi.Errors["TooManyMembers"].ID.Bytes()[:4]) {
		return iAnchorCouncil.UnpackTooManyMembersError(raw[4:])
	}
	if bytes.Equal(raw[:4], iAnchorCouncil.abi.Errors["ZeroDelay"].ID.Bytes()[:4]) {
		return iAnchorCouncil.UnpackZeroDelayError(raw[4:])
	}
	if bytes.Equal(raw[:4], iAnchorCouncil.abi.Errors["ZeroMember"].ID.Bytes()[:4]) {
		return iAnchorCouncil.UnpackZeroMemberError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// IAnchorCouncilApprovalsOutOfOrder represents a ApprovalsOutOfOrder error raised by the IAnchorCouncil contract.
type IAnchorCouncilApprovalsOutOfOrder struct {
	Previous common.Address
	Next     common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ApprovalsOutOfOrder(address previous, address next)
func IAnchorCouncilApprovalsOutOfOrderErrorID() common.Hash {
	return common.HexToHash("0x0ae564693c3c82af4f7665669daad9cb52931e900ad250dad060d045826b4808")
}

// UnpackApprovalsOutOfOrderError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ApprovalsOutOfOrder(address previous, address next)
func (iAnchorCouncil *IAnchorCouncil) UnpackApprovalsOutOfOrderError(raw []byte) (*IAnchorCouncilApprovalsOutOfOrder, error) {
	out := new(IAnchorCouncilApprovalsOutOfOrder)
	if err := iAnchorCouncil.abi.UnpackIntoInterface(out, "ApprovalsOutOfOrder", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IAnchorCouncilMembersNotAscending represents a MembersNotAscending error raised by the IAnchorCouncil contract.
type IAnchorCouncilMembersNotAscending struct {
	Index *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error MembersNotAscending(uint256 index)
func IAnchorCouncilMembersNotAscendingErrorID() common.Hash {
	return common.HexToHash("0xc3f8d748336d11048bde0479dc91abb1008ecb5d32412f551d6be14b4706823d")
}

// UnpackMembersNotAscendingError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error MembersNotAscending(uint256 index)
func (iAnchorCouncil *IAnchorCouncil) UnpackMembersNotAscendingError(raw []byte) (*IAnchorCouncilMembersNotAscending, error) {
	out := new(IAnchorCouncilMembersNotAscending)
	if err := iAnchorCouncil.abi.UnpackIntoInterface(out, "MembersNotAscending", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IAnchorCouncilMembershipChangePending represents a MembershipChangePending error raised by the IAnchorCouncil contract.
type IAnchorCouncilMembershipChangePending struct {
	ActivatesAt uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error MembershipChangePending(uint64 activatesAt)
func IAnchorCouncilMembershipChangePendingErrorID() common.Hash {
	return common.HexToHash("0x71af53390af28b7532d35d797ad193736472180032ec48826a7464c43fbbd274")
}

// UnpackMembershipChangePendingError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error MembershipChangePending(uint64 activatesAt)
func (iAnchorCouncil *IAnchorCouncil) UnpackMembershipChangePendingError(raw []byte) (*IAnchorCouncilMembershipChangePending, error) {
	out := new(IAnchorCouncilMembershipChangePending)
	if err := iAnchorCouncil.abi.UnpackIntoInterface(out, "MembershipChangePending", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IAnchorCouncilMembershipChangeTooEarly represents a MembershipChangeTooEarly error raised by the IAnchorCouncil contract.
type IAnchorCouncilMembershipChangeTooEarly struct {
	ActivatesAt uint64
	BlockTime   uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error MembershipChangeTooEarly(uint64 activatesAt, uint64 blockTime)
func IAnchorCouncilMembershipChangeTooEarlyErrorID() common.Hash {
	return common.HexToHash("0x28aa9e96bcb4ad370511c5aa8f6e06f67761917b02879c547859ba3dae588df4")
}

// UnpackMembershipChangeTooEarlyError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error MembershipChangeTooEarly(uint64 activatesAt, uint64 blockTime)
func (iAnchorCouncil *IAnchorCouncil) UnpackMembershipChangeTooEarlyError(raw []byte) (*IAnchorCouncilMembershipChangeTooEarly, error) {
	out := new(IAnchorCouncilMembershipChangeTooEarly)
	if err := iAnchorCouncil.abi.UnpackIntoInterface(out, "MembershipChangeTooEarly", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IAnchorCouncilNoMembershipChangePending represents a NoMembershipChangePending error raised by the IAnchorCouncil contract.
type IAnchorCouncilNoMembershipChangePending struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NoMembershipChangePending()
func IAnchorCouncilNoMembershipChangePendingErrorID() common.Hash {
	return common.HexToHash("0x56aac593abb77dbc8df3eb9b505ed28689777af94a7f6cb479dad68a1fc20929")
}

// UnpackNoMembershipChangePendingError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NoMembershipChangePending()
func (iAnchorCouncil *IAnchorCouncil) UnpackNoMembershipChangePendingError(raw []byte) (*IAnchorCouncilNoMembershipChangePending, error) {
	out := new(IAnchorCouncilNoMembershipChangePending)
	if err := iAnchorCouncil.abi.UnpackIntoInterface(out, "NoMembershipChangePending", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IAnchorCouncilNotALedger represents a NotALedger error raised by the IAnchorCouncil contract.
type IAnchorCouncilNotALedger struct {
	Ledger common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotALedger(address ledger)
func IAnchorCouncilNotALedgerErrorID() common.Hash {
	return common.HexToHash("0x44a5b4b2f2697952fce85b89e8be359271d0792b4f0be489d664b4771187455b")
}

// UnpackNotALedgerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotALedger(address ledger)
func (iAnchorCouncil *IAnchorCouncil) UnpackNotALedgerError(raw []byte) (*IAnchorCouncilNotALedger, error) {
	out := new(IAnchorCouncilNotALedger)
	if err := iAnchorCouncil.abi.UnpackIntoInterface(out, "NotALedger", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IAnchorCouncilNotAMember represents a NotAMember error raised by the IAnchorCouncil contract.
type IAnchorCouncilNotAMember struct {
	Recovered common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotAMember(address recovered)
func IAnchorCouncilNotAMemberErrorID() common.Hash {
	return common.HexToHash("0xc74f99ae608608fe36f55415de276becb198d973e61a153582b23d7478cbe28e")
}

// UnpackNotAMemberError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotAMember(address recovered)
func (iAnchorCouncil *IAnchorCouncil) UnpackNotAMemberError(raw []byte) (*IAnchorCouncilNotAMember, error) {
	out := new(IAnchorCouncilNotAMember)
	if err := iAnchorCouncil.abi.UnpackIntoInterface(out, "NotAMember", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IAnchorCouncilNotEnoughApprovals represents a NotEnoughApprovals error raised by the IAnchorCouncil contract.
type IAnchorCouncilNotEnoughApprovals struct {
	Given     *big.Int
	Threshold uint8
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotEnoughApprovals(uint256 given, uint8 threshold)
func IAnchorCouncilNotEnoughApprovalsErrorID() common.Hash {
	return common.HexToHash("0x689c646a47da69184d010de6cb6d79a368ecdc6297936c944e994559952e6b67")
}

// UnpackNotEnoughApprovalsError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotEnoughApprovals(uint256 given, uint8 threshold)
func (iAnchorCouncil *IAnchorCouncil) UnpackNotEnoughApprovalsError(raw []byte) (*IAnchorCouncilNotEnoughApprovals, error) {
	out := new(IAnchorCouncilNotEnoughApprovals)
	if err := iAnchorCouncil.abi.UnpackIntoInterface(out, "NotEnoughApprovals", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IAnchorCouncilSameMembership represents a SameMembership error raised by the IAnchorCouncil contract.
type IAnchorCouncilSameMembership struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SameMembership()
func IAnchorCouncilSameMembershipErrorID() common.Hash {
	return common.HexToHash("0xf823361da3aacad19202d9f0786b4d34b25dc77d9442b4bc523740fb97e7ad00")
}

// UnpackSameMembershipError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SameMembership()
func (iAnchorCouncil *IAnchorCouncil) UnpackSameMembershipError(raw []byte) (*IAnchorCouncilSameMembership, error) {
	out := new(IAnchorCouncilSameMembership)
	if err := iAnchorCouncil.abi.UnpackIntoInterface(out, "SameMembership", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IAnchorCouncilThresholdTooLow represents a ThresholdTooLow error raised by the IAnchorCouncil contract.
type IAnchorCouncilThresholdTooLow struct {
	Threshold uint8
	Minimum   uint8
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ThresholdTooLow(uint8 threshold, uint8 minimum)
func IAnchorCouncilThresholdTooLowErrorID() common.Hash {
	return common.HexToHash("0xaa8349cde43f9a01a0cf7d56428239b5db53ab2eb21b83b53dee2f165de704fc")
}

// UnpackThresholdTooLowError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ThresholdTooLow(uint8 threshold, uint8 minimum)
func (iAnchorCouncil *IAnchorCouncil) UnpackThresholdTooLowError(raw []byte) (*IAnchorCouncilThresholdTooLow, error) {
	out := new(IAnchorCouncilThresholdTooLow)
	if err := iAnchorCouncil.abi.UnpackIntoInterface(out, "ThresholdTooLow", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IAnchorCouncilTooFewMembers represents a TooFewMembers error raised by the IAnchorCouncil contract.
type IAnchorCouncilTooFewMembers struct {
	Count     *big.Int
	Threshold uint8
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TooFewMembers(uint256 count, uint8 threshold)
func IAnchorCouncilTooFewMembersErrorID() common.Hash {
	return common.HexToHash("0x6b4801260737f62062d2e7010fdf0d4495661cd59a1c4ffed4d86c4264d367cb")
}

// UnpackTooFewMembersError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TooFewMembers(uint256 count, uint8 threshold)
func (iAnchorCouncil *IAnchorCouncil) UnpackTooFewMembersError(raw []byte) (*IAnchorCouncilTooFewMembers, error) {
	out := new(IAnchorCouncilTooFewMembers)
	if err := iAnchorCouncil.abi.UnpackIntoInterface(out, "TooFewMembers", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IAnchorCouncilTooManyMembers represents a TooManyMembers error raised by the IAnchorCouncil contract.
type IAnchorCouncilTooManyMembers struct {
	Count   *big.Int
	Maximum *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TooManyMembers(uint256 count, uint256 maximum)
func IAnchorCouncilTooManyMembersErrorID() common.Hash {
	return common.HexToHash("0x3d5600991605612d4fbd5921afde15dadbbdea16fed93c002fe2a8fb1bceda3e")
}

// UnpackTooManyMembersError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TooManyMembers(uint256 count, uint256 maximum)
func (iAnchorCouncil *IAnchorCouncil) UnpackTooManyMembersError(raw []byte) (*IAnchorCouncilTooManyMembers, error) {
	out := new(IAnchorCouncilTooManyMembers)
	if err := iAnchorCouncil.abi.UnpackIntoInterface(out, "TooManyMembers", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IAnchorCouncilZeroDelay represents a ZeroDelay error raised by the IAnchorCouncil contract.
type IAnchorCouncilZeroDelay struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ZeroDelay()
func IAnchorCouncilZeroDelayErrorID() common.Hash {
	return common.HexToHash("0x34c2b9dab051f3b49e0b9121d90ae110f06df318864d2741799dc06f9b92cebc")
}

// UnpackZeroDelayError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ZeroDelay()
func (iAnchorCouncil *IAnchorCouncil) UnpackZeroDelayError(raw []byte) (*IAnchorCouncilZeroDelay, error) {
	out := new(IAnchorCouncilZeroDelay)
	if err := iAnchorCouncil.abi.UnpackIntoInterface(out, "ZeroDelay", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IAnchorCouncilZeroMember represents a ZeroMember error raised by the IAnchorCouncil contract.
type IAnchorCouncilZeroMember struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ZeroMember()
func IAnchorCouncilZeroMemberErrorID() common.Hash {
	return common.HexToHash("0x7553caa6a7b0ddb1dba7e7fe0e24456b98ff20955606dbf7fd658e4702c19fee")
}

// UnpackZeroMemberError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ZeroMember()
func (iAnchorCouncil *IAnchorCouncil) UnpackZeroMemberError(raw []byte) (*IAnchorCouncilZeroMember, error) {
	out := new(IAnchorCouncilZeroMember)
	if err := iAnchorCouncil.abi.UnpackIntoInterface(out, "ZeroMember", raw); err != nil {
		return nil, err
	}
	return out, nil
}
