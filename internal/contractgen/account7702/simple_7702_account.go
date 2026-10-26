// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package account7702

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

// BaseAccountCall is an auto generated low-level Go binding around an user-defined struct.
type BaseAccountCall struct {
	Target common.Address
	Value  *big.Int
	Data   []byte
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

// Simple7702AccountMetaData contains all meta data concerning the Simple7702Account contract.
var Simple7702AccountMetaData = bind.MetaData{
	ABI: "[{\"type\":\"fallback\",\"stateMutability\":\"payable\"},{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"entryPoint\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEntryPoint\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"execute\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"executeBatch\",\"inputs\":[{\"name\":\"calls\",\"type\":\"tuple[]\",\"internalType\":\"structBaseAccount.Call[]\",\"components\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getNonce\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isValidSignature\",\"inputs\":[{\"name\":\"hash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"magicValue\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"onERC1155BatchReceived\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"onERC1155Received\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"onERC721Received\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"id\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"validateUserOp\",\"inputs\":[{\"name\":\"userOp\",\"type\":\"tuple\",\"internalType\":\"structPackedUserOperation\",\"components\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"nonce\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"initCode\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"callData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"accountGasLimits\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"preVerificationGas\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"gasFees\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"paymasterAndData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"userOpHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"missingAccountFunds\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"validationData\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignature\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignatureLength\",\"inputs\":[{\"name\":\"length\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignatureS\",\"inputs\":[{\"name\":\"s\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ExecuteError\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"error\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]",
	ID:  "Simple7702Account",
	Bin: "0x6080604052348015600e575f5ffd5b506111b98061001c5f395ff3fe6080604052600436106100a7575f3560e01c8063b0d691fe11610062578063bc197c811161004a578063bc197c8114610209578063d087d2881461024d578063f23a6e611461026157005b8063b0d691fe146101bc578063b61d27f6146101ea57005b80631626ba7e116100905780631626ba7e1461015157806319822f7c1461017057806334fcd5be1461019d57005b806301ffc9a7146100a9578063150b7a02146100dd575b005b3480156100b4575f5ffd5b506100c86100c3366004610b56565b6102a5565b60405190151581526020015b60405180910390f35b3480156100e8575f5ffd5b506101206100f7366004610cc3565b7f150b7a0200000000000000000000000000000000000000000000000000000000949350505050565b6040517fffffffff0000000000000000000000000000000000000000000000000000000090911681526020016100d4565b34801561015c575f5ffd5b5061012061016b366004610d27565b610421565b34801561017b575f5ffd5b5061018f61018a366004610d6b565b61047f565b6040519081526020016100d4565b3480156101a8575f5ffd5b506100a76101b7366004610dba565b61049d565b3480156101c7575f5ffd5b50604051734337084d9e255ff0702461cf8895ce9e3b5ff10881526020016100d4565b3480156101f5575f5ffd5b506100a7610204366004610e2b565b6105a6565b348015610214575f5ffd5b50610120610223366004610f2e565b7fbc197c810000000000000000000000000000000000000000000000000000000095945050505050565b348015610258575f5ffd5b5061018f610605565b34801561026c575f5ffd5b5061012061027b366004610fdd565b7ff23a6e610000000000000000000000000000000000000000000000000000000095945050505050565b5f7fffffffff0000000000000000000000000000000000000000000000000000000082167f01ffc9a700000000000000000000000000000000000000000000000000000000148061033757507fffffffff0000000000000000000000000000000000000000000000000000000082167f19822f7c00000000000000000000000000000000000000000000000000000000145b8061038357507fffffffff0000000000000000000000000000000000000000000000000000000082167f1626ba7e00000000000000000000000000000000000000000000000000000000145b806103cf57507fffffffff0000000000000000000000000000000000000000000000000000000082167f4e2312e000000000000000000000000000000000000000000000000000000000145b8061041b57507fffffffff0000000000000000000000000000000000000000000000000000000082167f150b7a0200000000000000000000000000000000000000000000000000000000145b92915050565b5f61042c83836106b4565b610456577fffffffff00000000000000000000000000000000000000000000000000000000610478565b7f1626ba7e000000000000000000000000000000000000000000000000000000005b9392505050565b5f6104886106de565b610492848461075d565b9050610478826107bf565b6104a5610804565b805f5b818110156105a057368484838181106104c3576104c3611031565b90506020028101906104d5919061105e565b90505f6105336104e8602084018461109a565b60208401356104fa60408601866110b3565b8080601f0160208091040260200160405190810160405280939291908181526020018383808284375f920191909152505050505a61088b565b905080610596578360010361054f5761054a6108a1565b610596565b826105595f6108b2565b6040517f5a15467500000000000000000000000000000000000000000000000000000000815260040161058d92919061111b565b60405180910390fd5b50506001016104a8565b50505050565b6105ae610804565b5f6105ef858585858080601f0160208091040260200160405190810160405280939291908181526020018383808284375f920191909152505050505a61088b565b9050806105fe576105fe6108a1565b5050505050565b5f734337084d9e255ff0702461cf8895ce9e3b5ff1086040517f35567e1a0000000000000000000000000000000000000000000000000000000081523060048201525f602482015273ffffffffffffffffffffffffffffffffffffffff91909116906335567e1a90604401602060405180830381865afa15801561068b573d5f5f3e3d5ffd5b505050506040513d601f19601f820116820180604052508101906106af9190611175565b905090565b5f306106c084846108e3565b73ffffffffffffffffffffffffffffffffffffffff16149392505050565b33734337084d9e255ff0702461cf8895ce9e3b5ff1081461075b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f6163636f756e743a206e6f742066726f6d20456e747279506f696e7400000000604482015260640161058d565b565b5f6107a9826107706101008601866110b3565b8080601f0160208091040260200160405190810160405280939291908181526020018383808284375f920191909152506106b492505050565b6107b4576001610478565b505f92915050565b50565b80156107bc576040515f90339083908381818185875af1925050503d805f81146105fe576040519150601f19603f3d011682016040523d82523d5f602084013e6105fe565b33301480610825575033734337084d9e255ff0702461cf8895ce9e3b5ff108145b61075b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601b60248201527f6e6f742066726f6d2073656c66206f7220456e747279506f696e740000000000604482015260640161058d565b5f5f5f845160208601878987f195945050505050565b61075b6108ad5f6108b2565b61090b565b60603d82156108c657828111156108c65750815b604051602082018101604052818152815f602083013e9392505050565b5f5f5f5f6108f18686610913565b925092509250610901828261095c565b5090949350505050565b805160208201fd5b5f5f5f835160410361094a576020840151604085015160608601515f1a61093c88828585610a63565b955095509550505050610955565b505081515f91506002905b9250925092565b5f82600381111561096f5761096f61118c565b03610978575050565b600182600381111561098c5761098c61118c565b036109c3576040517ff645eedf00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60028260038111156109d7576109d761118c565b03610a11576040517ffce698f70000000000000000000000000000000000000000000000000000000081526004810182905260240161058d565b6003826003811115610a2557610a2561118c565b03610a5f576040517fd78bce0c0000000000000000000000000000000000000000000000000000000081526004810182905260240161058d565b5050565b5f80807f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a0841115610a9c57505f91506003905082610b4c565b604080515f808252602082018084528a905260ff891692820192909252606081018790526080810186905260019060a0016020604051602081039080840390855afa158015610aed573d5f5f3e3d5ffd5b50506040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0015191505073ffffffffffffffffffffffffffffffffffffffff8116610b4357505f925060019150829050610b4c565b92505f91508190505b9450945094915050565b5f60208284031215610b66575f5ffd5b81357fffffffff0000000000000000000000000000000000000000000000000000000081168114610478575f5ffd5b803573ffffffffffffffffffffffffffffffffffffffff81168114610bb8575f5ffd5b919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715610c3157610c31610bbd565b604052919050565b5f82601f830112610c48575f5ffd5b813567ffffffffffffffff811115610c6257610c62610bbd565b610c9360207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601610bea565b818152846020838601011115610ca7575f5ffd5b816020850160208301375f918101602001919091529392505050565b5f5f5f5f60808587031215610cd6575f5ffd5b610cdf85610b95565b9350610ced60208601610b95565b925060408501359150606085013567ffffffffffffffff811115610d0f575f5ffd5b610d1b87828801610c39565b91505092959194509250565b5f5f60408385031215610d38575f5ffd5b82359150602083013567ffffffffffffffff811115610d55575f5ffd5b610d6185828601610c39565b9150509250929050565b5f5f5f60608486031215610d7d575f5ffd5b833567ffffffffffffffff811115610d93575f5ffd5b84016101208187031215610da5575f5ffd5b95602085013595506040909401359392505050565b5f5f60208385031215610dcb575f5ffd5b823567ffffffffffffffff811115610de1575f5ffd5b8301601f81018513610df1575f5ffd5b803567ffffffffffffffff811115610e07575f5ffd5b8560208260051b8401011115610e1b575f5ffd5b6020919091019590945092505050565b5f5f5f5f60608587031215610e3e575f5ffd5b610e4785610b95565b935060208501359250604085013567ffffffffffffffff811115610e69575f5ffd5b8501601f81018713610e79575f5ffd5b803567ffffffffffffffff811115610e8f575f5ffd5b876020828401011115610ea0575f5ffd5b949793965060200194505050565b5f82601f830112610ebd575f5ffd5b813567ffffffffffffffff811115610ed757610ed7610bbd565b8060051b610ee760208201610bea565b91825260208185018101929081019086841115610f02575f5ffd5b6020860192505b83831015610f24578235825260209283019290910190610f09565b9695505050505050565b5f5f5f5f5f60a08688031215610f42575f5ffd5b610f4b86610b95565b9450610f5960208701610b95565b9350604086013567ffffffffffffffff811115610f74575f5ffd5b610f8088828901610eae565b935050606086013567ffffffffffffffff811115610f9c575f5ffd5b610fa888828901610eae565b925050608086013567ffffffffffffffff811115610fc4575f5ffd5b610fd088828901610c39565b9150509295509295909350565b5f5f5f5f5f60a08688031215610ff1575f5ffd5b610ffa86610b95565b945061100860208701610b95565b93506040860135925060608601359150608086013567ffffffffffffffff811115610fc4575f5ffd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b5f82357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa1833603018112611090575f5ffd5b9190910192915050565b5f602082840312156110aa575f5ffd5b61047882610b95565b5f5f83357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe18436030181126110e6575f5ffd5b83018035915067ffffffffffffffff821115611100575f5ffd5b602001915036819003821315611114575f5ffd5b9250929050565b828152604060208201525f82518060408401528060208501606085015e5f6060828501015260607fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8301168401019150509392505050565b5f60208284031215611185575f5ffd5b5051919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602160045260245ffd",
}

// Simple7702Account is an auto generated Go binding around an Ethereum contract.
type Simple7702Account struct {
	abi abi.ABI
}

// GetABI returns the ABI associated with this contract binding.
func (c *Simple7702Account) GetABI() abi.ABI {
	return c.abi
}

// NewSimple7702Account creates a new instance of Simple7702Account.
func NewSimple7702Account() *Simple7702Account {
	parsed, err := Simple7702AccountMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &Simple7702Account{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *Simple7702Account) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackEntryPoint is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb0d691fe.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function entryPoint() pure returns(address)
func (simple7702Account *Simple7702Account) PackEntryPoint() []byte {
	enc, err := simple7702Account.abi.Pack("entryPoint")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackEntryPoint is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb0d691fe.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function entryPoint() pure returns(address)
func (simple7702Account *Simple7702Account) TryPackEntryPoint() ([]byte, error) {
	return simple7702Account.abi.Pack("entryPoint")
}

// UnpackEntryPoint is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb0d691fe.
//
// Solidity: function entryPoint() pure returns(address)
func (simple7702Account *Simple7702Account) UnpackEntryPoint(data []byte) (common.Address, error) {
	out, err := simple7702Account.abi.Unpack("entryPoint", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackExecute is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb61d27f6.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function execute(address target, uint256 value, bytes data) returns()
func (simple7702Account *Simple7702Account) PackExecute(target common.Address, value *big.Int, data []byte) []byte {
	enc, err := simple7702Account.abi.Pack("execute", target, value, data)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackExecute is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb61d27f6.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function execute(address target, uint256 value, bytes data) returns()
func (simple7702Account *Simple7702Account) TryPackExecute(target common.Address, value *big.Int, data []byte) ([]byte, error) {
	return simple7702Account.abi.Pack("execute", target, value, data)
}

// PackExecuteBatch is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x34fcd5be.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function executeBatch((address,uint256,bytes)[] calls) returns()
func (simple7702Account *Simple7702Account) PackExecuteBatch(calls []BaseAccountCall) []byte {
	enc, err := simple7702Account.abi.Pack("executeBatch", calls)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackExecuteBatch is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x34fcd5be.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function executeBatch((address,uint256,bytes)[] calls) returns()
func (simple7702Account *Simple7702Account) TryPackExecuteBatch(calls []BaseAccountCall) ([]byte, error) {
	return simple7702Account.abi.Pack("executeBatch", calls)
}

// PackGetNonce is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd087d288.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getNonce() view returns(uint256)
func (simple7702Account *Simple7702Account) PackGetNonce() []byte {
	enc, err := simple7702Account.abi.Pack("getNonce")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetNonce is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd087d288.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getNonce() view returns(uint256)
func (simple7702Account *Simple7702Account) TryPackGetNonce() ([]byte, error) {
	return simple7702Account.abi.Pack("getNonce")
}

// UnpackGetNonce is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xd087d288.
//
// Solidity: function getNonce() view returns(uint256)
func (simple7702Account *Simple7702Account) UnpackGetNonce(data []byte) (*big.Int, error) {
	out, err := simple7702Account.abi.Unpack("getNonce", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackIsValidSignature is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1626ba7e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function isValidSignature(bytes32 hash, bytes signature) view returns(bytes4 magicValue)
func (simple7702Account *Simple7702Account) PackIsValidSignature(hash [32]byte, signature []byte) []byte {
	enc, err := simple7702Account.abi.Pack("isValidSignature", hash, signature)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackIsValidSignature is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1626ba7e.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function isValidSignature(bytes32 hash, bytes signature) view returns(bytes4 magicValue)
func (simple7702Account *Simple7702Account) TryPackIsValidSignature(hash [32]byte, signature []byte) ([]byte, error) {
	return simple7702Account.abi.Pack("isValidSignature", hash, signature)
}

// UnpackIsValidSignature is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x1626ba7e.
//
// Solidity: function isValidSignature(bytes32 hash, bytes signature) view returns(bytes4 magicValue)
func (simple7702Account *Simple7702Account) UnpackIsValidSignature(data []byte) ([4]byte, error) {
	out, err := simple7702Account.abi.Unpack("isValidSignature", data)
	if err != nil {
		return *new([4]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([4]byte)).(*[4]byte)
	return out0, nil
}

// PackOnERC1155BatchReceived is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbc197c81.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function onERC1155BatchReceived(address , address , uint256[] , uint256[] , bytes ) returns(bytes4)
func (simple7702Account *Simple7702Account) PackOnERC1155BatchReceived(arg0 common.Address, arg1 common.Address, arg2 []*big.Int, arg3 []*big.Int, arg4 []byte) []byte {
	enc, err := simple7702Account.abi.Pack("onERC1155BatchReceived", arg0, arg1, arg2, arg3, arg4)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackOnERC1155BatchReceived is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbc197c81.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function onERC1155BatchReceived(address , address , uint256[] , uint256[] , bytes ) returns(bytes4)
func (simple7702Account *Simple7702Account) TryPackOnERC1155BatchReceived(arg0 common.Address, arg1 common.Address, arg2 []*big.Int, arg3 []*big.Int, arg4 []byte) ([]byte, error) {
	return simple7702Account.abi.Pack("onERC1155BatchReceived", arg0, arg1, arg2, arg3, arg4)
}

// UnpackOnERC1155BatchReceived is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xbc197c81.
//
// Solidity: function onERC1155BatchReceived(address , address , uint256[] , uint256[] , bytes ) returns(bytes4)
func (simple7702Account *Simple7702Account) UnpackOnERC1155BatchReceived(data []byte) ([4]byte, error) {
	out, err := simple7702Account.abi.Unpack("onERC1155BatchReceived", data)
	if err != nil {
		return *new([4]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([4]byte)).(*[4]byte)
	return out0, nil
}

// PackOnERC1155Received is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf23a6e61.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function onERC1155Received(address , address , uint256 , uint256 , bytes ) returns(bytes4)
func (simple7702Account *Simple7702Account) PackOnERC1155Received(arg0 common.Address, arg1 common.Address, arg2 *big.Int, arg3 *big.Int, arg4 []byte) []byte {
	enc, err := simple7702Account.abi.Pack("onERC1155Received", arg0, arg1, arg2, arg3, arg4)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackOnERC1155Received is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf23a6e61.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function onERC1155Received(address , address , uint256 , uint256 , bytes ) returns(bytes4)
func (simple7702Account *Simple7702Account) TryPackOnERC1155Received(arg0 common.Address, arg1 common.Address, arg2 *big.Int, arg3 *big.Int, arg4 []byte) ([]byte, error) {
	return simple7702Account.abi.Pack("onERC1155Received", arg0, arg1, arg2, arg3, arg4)
}

// UnpackOnERC1155Received is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xf23a6e61.
//
// Solidity: function onERC1155Received(address , address , uint256 , uint256 , bytes ) returns(bytes4)
func (simple7702Account *Simple7702Account) UnpackOnERC1155Received(data []byte) ([4]byte, error) {
	out, err := simple7702Account.abi.Unpack("onERC1155Received", data)
	if err != nil {
		return *new([4]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([4]byte)).(*[4]byte)
	return out0, nil
}

// PackOnERC721Received is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x150b7a02.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function onERC721Received(address , address , uint256 , bytes ) returns(bytes4)
func (simple7702Account *Simple7702Account) PackOnERC721Received(arg0 common.Address, arg1 common.Address, arg2 *big.Int, arg3 []byte) []byte {
	enc, err := simple7702Account.abi.Pack("onERC721Received", arg0, arg1, arg2, arg3)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackOnERC721Received is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x150b7a02.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function onERC721Received(address , address , uint256 , bytes ) returns(bytes4)
func (simple7702Account *Simple7702Account) TryPackOnERC721Received(arg0 common.Address, arg1 common.Address, arg2 *big.Int, arg3 []byte) ([]byte, error) {
	return simple7702Account.abi.Pack("onERC721Received", arg0, arg1, arg2, arg3)
}

// UnpackOnERC721Received is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x150b7a02.
//
// Solidity: function onERC721Received(address , address , uint256 , bytes ) returns(bytes4)
func (simple7702Account *Simple7702Account) UnpackOnERC721Received(data []byte) ([4]byte, error) {
	out, err := simple7702Account.abi.Unpack("onERC721Received", data)
	if err != nil {
		return *new([4]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([4]byte)).(*[4]byte)
	return out0, nil
}

// PackSupportsInterface is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x01ffc9a7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function supportsInterface(bytes4 id) pure returns(bool)
func (simple7702Account *Simple7702Account) PackSupportsInterface(id [4]byte) []byte {
	enc, err := simple7702Account.abi.Pack("supportsInterface", id)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSupportsInterface is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x01ffc9a7.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function supportsInterface(bytes4 id) pure returns(bool)
func (simple7702Account *Simple7702Account) TryPackSupportsInterface(id [4]byte) ([]byte, error) {
	return simple7702Account.abi.Pack("supportsInterface", id)
}

// UnpackSupportsInterface is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 id) pure returns(bool)
func (simple7702Account *Simple7702Account) UnpackSupportsInterface(data []byte) (bool, error) {
	out, err := simple7702Account.abi.Unpack("supportsInterface", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackValidateUserOp is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x19822f7c.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function validateUserOp((address,uint256,bytes,bytes,bytes32,uint256,bytes32,bytes,bytes) userOp, bytes32 userOpHash, uint256 missingAccountFunds) returns(uint256 validationData)
func (simple7702Account *Simple7702Account) PackValidateUserOp(userOp PackedUserOperation, userOpHash [32]byte, missingAccountFunds *big.Int) []byte {
	enc, err := simple7702Account.abi.Pack("validateUserOp", userOp, userOpHash, missingAccountFunds)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackValidateUserOp is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x19822f7c.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function validateUserOp((address,uint256,bytes,bytes,bytes32,uint256,bytes32,bytes,bytes) userOp, bytes32 userOpHash, uint256 missingAccountFunds) returns(uint256 validationData)
func (simple7702Account *Simple7702Account) TryPackValidateUserOp(userOp PackedUserOperation, userOpHash [32]byte, missingAccountFunds *big.Int) ([]byte, error) {
	return simple7702Account.abi.Pack("validateUserOp", userOp, userOpHash, missingAccountFunds)
}

// UnpackValidateUserOp is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x19822f7c.
//
// Solidity: function validateUserOp((address,uint256,bytes,bytes,bytes32,uint256,bytes32,bytes,bytes) userOp, bytes32 userOpHash, uint256 missingAccountFunds) returns(uint256 validationData)
func (simple7702Account *Simple7702Account) UnpackValidateUserOp(data []byte) (*big.Int, error) {
	out, err := simple7702Account.abi.Unpack("validateUserOp", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// UnpackError attempts to decode the provided error data using user-defined
// error definitions.
func (simple7702Account *Simple7702Account) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], simple7702Account.abi.Errors["ECDSAInvalidSignature"].ID.Bytes()[:4]) {
		return simple7702Account.UnpackECDSAInvalidSignatureError(raw[4:])
	}
	if bytes.Equal(raw[:4], simple7702Account.abi.Errors["ECDSAInvalidSignatureLength"].ID.Bytes()[:4]) {
		return simple7702Account.UnpackECDSAInvalidSignatureLengthError(raw[4:])
	}
	if bytes.Equal(raw[:4], simple7702Account.abi.Errors["ECDSAInvalidSignatureS"].ID.Bytes()[:4]) {
		return simple7702Account.UnpackECDSAInvalidSignatureSError(raw[4:])
	}
	if bytes.Equal(raw[:4], simple7702Account.abi.Errors["ExecuteError"].ID.Bytes()[:4]) {
		return simple7702Account.UnpackExecuteErrorError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// Simple7702AccountECDSAInvalidSignature represents a ECDSAInvalidSignature error raised by the Simple7702Account contract.
type Simple7702AccountECDSAInvalidSignature struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ECDSAInvalidSignature()
func Simple7702AccountECDSAInvalidSignatureErrorID() common.Hash {
	return common.HexToHash("0xf645eedf0193584640b6b90cb9477e4c95b98636c148a891d4c0a146dc46e75a")
}

// UnpackECDSAInvalidSignatureError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ECDSAInvalidSignature()
func (simple7702Account *Simple7702Account) UnpackECDSAInvalidSignatureError(raw []byte) (*Simple7702AccountECDSAInvalidSignature, error) {
	out := new(Simple7702AccountECDSAInvalidSignature)
	if err := simple7702Account.abi.UnpackIntoInterface(out, "ECDSAInvalidSignature", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// Simple7702AccountECDSAInvalidSignatureLength represents a ECDSAInvalidSignatureLength error raised by the Simple7702Account contract.
type Simple7702AccountECDSAInvalidSignatureLength struct {
	Length *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ECDSAInvalidSignatureLength(uint256 length)
func Simple7702AccountECDSAInvalidSignatureLengthErrorID() common.Hash {
	return common.HexToHash("0xfce698f7e8e5342cd615f641317bc45fe7e1e4a8b0a14dd1383ff8dc9c41917f")
}

// UnpackECDSAInvalidSignatureLengthError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ECDSAInvalidSignatureLength(uint256 length)
func (simple7702Account *Simple7702Account) UnpackECDSAInvalidSignatureLengthError(raw []byte) (*Simple7702AccountECDSAInvalidSignatureLength, error) {
	out := new(Simple7702AccountECDSAInvalidSignatureLength)
	if err := simple7702Account.abi.UnpackIntoInterface(out, "ECDSAInvalidSignatureLength", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// Simple7702AccountECDSAInvalidSignatureS represents a ECDSAInvalidSignatureS error raised by the Simple7702Account contract.
type Simple7702AccountECDSAInvalidSignatureS struct {
	S [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ECDSAInvalidSignatureS(bytes32 s)
func Simple7702AccountECDSAInvalidSignatureSErrorID() common.Hash {
	return common.HexToHash("0xd78bce0cccb935155ed6428d1c13e50b7f3550fd2b66b9fe266006fea4a5e1eb")
}

// UnpackECDSAInvalidSignatureSError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ECDSAInvalidSignatureS(bytes32 s)
func (simple7702Account *Simple7702Account) UnpackECDSAInvalidSignatureSError(raw []byte) (*Simple7702AccountECDSAInvalidSignatureS, error) {
	out := new(Simple7702AccountECDSAInvalidSignatureS)
	if err := simple7702Account.abi.UnpackIntoInterface(out, "ECDSAInvalidSignatureS", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// Simple7702AccountExecuteError represents a ExecuteError error raised by the Simple7702Account contract.
type Simple7702AccountExecuteError struct {
	Index *big.Int
	Error []byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ExecuteError(uint256 index, bytes error)
func Simple7702AccountExecuteErrorErrorID() common.Hash {
	return common.HexToHash("0x5a1546750bc344ece535156ebbf17c8a06ed4f348e2e1ed293a251976e2c7047")
}

// UnpackExecuteErrorError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ExecuteError(uint256 index, bytes error)
func (simple7702Account *Simple7702Account) UnpackExecuteErrorError(raw []byte) (*Simple7702AccountExecuteError, error) {
	out := new(Simple7702AccountExecuteError)
	if err := simple7702Account.abi.UnpackIntoInterface(out, "ExecuteError", raw); err != nil {
		return nil, err
	}
	return out, nil
}
