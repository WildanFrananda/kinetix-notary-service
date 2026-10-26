// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ledger

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

// LiabilityLedgerMetaData contains all meta data concerning the LiabilityLedger contract.
var LiabilityLedgerMetaData = bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"admin\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"adminTransferDelay\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"anchorer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"genesisEpoch_\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"genesisDigest_\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"genesisAsOf\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"ANCHORER_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"DEFAULT_ADMIN_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptDefaultAdminTransfer\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"anchor\",\"inputs\":[{\"name\":\"header\",\"type\":\"tuple\",\"internalType\":\"structEpochHeader\",\"components\":[{\"name\":\"number\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"previousDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"total\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"leafCount\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"asOf\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"beginDefaultAdminTransfer\",\"inputs\":[{\"name\":\"newAdmin\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"cancelDefaultAdminTransfer\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"changeDefaultAdminDelay\",\"inputs\":[{\"name\":\"newDelay\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"defaultAdmin\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"defaultAdminDelay\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"defaultAdminDelayIncreaseWait\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"epochAt\",\"inputs\":[{\"name\":\"number\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structEpoch\",\"components\":[{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"digest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"total\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"leafCount\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"asOf\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"anchoredAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"genesisDigest\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"genesisEpoch\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRoleAdmin\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"grantRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"hasRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"head\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isSealed\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"lastAsOf\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"nextEpoch\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingDefaultAdmin\",\"inputs\":[],\"outputs\":[{\"name\":\"newAdmin\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"schedule\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingDefaultAdminDelay\",\"inputs\":[],\"outputs\":[{\"name\":\"newDelay\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"schedule\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revokeRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"rollbackDefaultAdminDelay\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"scheme\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"seal\",\"inputs\":[{\"name\":\"expectedHead\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyInclusion\",\"inputs\":[{\"name\":\"number\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"commitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"amount\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"index\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"siblings\",\"type\":\"tuple[]\",\"internalType\":\"structSumNode[]\",\"components\":[{\"name\":\"hash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sum\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"DefaultAdminDelayChangeCanceled\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DefaultAdminDelayChangeScheduled\",\"inputs\":[{\"name\":\"newDelay\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"},{\"name\":\"effectSchedule\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DefaultAdminTransferCanceled\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DefaultAdminTransferScheduled\",\"inputs\":[{\"name\":\"newAdmin\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"acceptSchedule\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EpochAnchored\",\"inputs\":[{\"name\":\"number\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"digest\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"previousDigest\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"root\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"total\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"leafCount\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"asOf\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"LedgerSealed\",\"inputs\":[{\"name\":\"head\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"nextEpoch\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"lastAsOf\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleAdminChanged\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"previousAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"newAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleGranted\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleRevoked\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AccessControlBadConfirmation\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AccessControlEnforcedDefaultAdminDelay\",\"inputs\":[{\"name\":\"schedule\",\"type\":\"uint48\",\"internalType\":\"uint48\"}]},{\"type\":\"error\",\"name\":\"AccessControlEnforcedDefaultAdminRules\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AccessControlInvalidDefaultAdmin\",\"inputs\":[{\"name\":\"defaultAdmin\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"AccessControlUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"neededRole\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"EmptyEpochMismatch\",\"inputs\":[{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"total\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"EpochAlreadyAnchored\",\"inputs\":[{\"name\":\"number\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"digest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"EpochBeforeGenesis\",\"inputs\":[{\"name\":\"genesisEpoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"number\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"EpochConflict\",\"inputs\":[{\"name\":\"number\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"anchored\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"given\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"EpochOutOfSequence\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"given\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"HeadMismatch\",\"inputs\":[{\"name\":\"head\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"InvalidGenesis\",\"inputs\":[{\"name\":\"genesisEpoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"genesisDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"LedgerIsSealed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PreviousDigestMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"given\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"SafeCastOverflowedUintDowncast\",\"inputs\":[{\"name\":\"bits\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"SnapshotInFuture\",\"inputs\":[{\"name\":\"asOf\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"blockTime\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"SnapshotNotAfterPrevious\",\"inputs\":[{\"name\":\"asOf\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"previous\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"TooManyLeaves\",\"inputs\":[{\"name\":\"leafCount\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"maximum\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"UnknownEpoch\",\"inputs\":[{\"name\":\"number\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"ZeroAnchorer\",\"inputs\":[]}]",
	ID:  "LiabilityLedger",
	Bin: "0x60c060405234801561000f575f5ffd5b50604051612a7a380380612a7a83398101604081905261002e916102a9565b84866001600160a01b03811661005e57604051636116401160e11b81525f60048201526024015b60405180910390fd5b600180546001600160d01b0316600160d01b65ffffffffffff8516021790556100875f82610162565b5050506001600160401b038316158215146100c75760405163111f019f60e01b81526001600160401b038416600482015260248101839052604401610055565b6001600160a01b0384166100ee57604051637110c4a160e01b815260040160405180910390fd5b6101187f4a8e1518fc3b69d274d0d6c14a98dd84f4bfb1e39bbbe43f417772dd939cd80885610162565b506001600160401b03928316608081905260a0839052600492909255600580549190931668010000000000000000026001600160801b03199091169091171790555061031e915050565b5f826101be575f61017b6002546001600160a01b031690565b6001600160a01b0316146101a257604051631fe1e13d60e11b815260040160405180910390fd5b600280546001600160a01b0319166001600160a01b0384161790555b6101c883836101d1565b90505b92915050565b5f828152602081815260408083206001600160a01b038516845290915281205460ff16610271575f838152602081815260408083206001600160a01b03861684529091529020805460ff191660011790556102293390565b6001600160a01b0316826001600160a01b0316847f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a45060016101cb565b505f6101cb565b80516001600160a01b038116811461028e575f5ffd5b919050565b80516001600160401b038116811461028e575f5ffd5b5f5f5f5f5f5f60c087890312156102be575f5ffd5b6102c787610278565b9550602087015165ffffffffffff811681146102e1575f5ffd5b94506102ef60408801610278565b93506102fd60608801610293565b92506080870151915061031260a08801610293565b90509295509295509295565b60805160a0516127266103545f395f6102d601525f81816104d101528181610768015281816107a7015261186b01526127265ff3fe608060405234801561000f575f5ffd5b50600436106101d1575f3560e01c80638da5cb5b116100fe578063b07eeda81161009e578063cefc14291161006e578063cefc14291461050e578063cf6eefb714610516578063d547741f14610562578063d602b9fd14610575575f5ffd5b8063b07eeda8146104b9578063b70e6be6146104cc578063c0566abe146104f3578063cc8463c814610506575f5ffd5b8063a1eda53c116100d9578063a1eda53c14610446578063a217fddf1461046d578063a89210cb14610474578063aea0e78b146104a7575f5ffd5b80638da5cb5b146103f35780638f7dcfa3146103fb57806391d1485414610403575f5ffd5b8063631f9852116101745780637cfd9a51116101445780637cfd9a51146102f85780637d60d3551461031f57806384ef8ffc146103a55780638518b0ad146103e4575f5ffd5b8063631f98521461028c578063634e93da146102ab578063649a5ec7146102be57806379ac2661146102d1575f5ffd5b80631ffcf5d4116101af5780631ffcf5d414610223578063248a9ca3146102445780632f2ff15d1461026657806336568abe14610279575f5ffd5b806301ffc9a7146101d5578063022d63fb146101fd5780630aa6220b14610219575b5f5ffd5b6101e86101e3366004612316565b61057d565b60405190151581526020015b60405180910390f35b620697805b60405165ffffffffffff90911681526020016101f4565b6102216105d8565b005b610236610231366004612355565b6105ed565b6040519081526020016101f4565b61023661025236600461236e565b5f9081526020819052604090206001015490565b6102216102743660046123ad565b610cfa565b6102216102873660046123ad565b610d3f565b600554700100000000000000000000000000000000900460ff166101e8565b6102216102b93660046123d7565b610e44565b6102216102cc3660046123f0565b610e57565b6102367f000000000000000000000000000000000000000000000000000000000000000081565b6102367f4a8e1518fc3b69d274d0d6c14a98dd84f4bfb1e39bbbe43f417772dd939cd80881565b61033261032d36600461242c565b610e6a565b6040516101f491905f60c082019050825182526020830151602083015267ffffffffffffffff604084015116604083015267ffffffffffffffff606084015116606083015267ffffffffffffffff608084015116608083015267ffffffffffffffff60a08401511660a083015292915050565b60025473ffffffffffffffffffffffffffffffffffffffff165b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016101f4565b604051600181526020016101f4565b6103bf610f2b565b600454610236565b6101e86104113660046123ad565b5f9182526020828152604080842073ffffffffffffffffffffffffffffffffffffffff93909316845291905290205460ff1690565b61044e610f50565b6040805165ffffffffffff9384168152929091166020830152016101f4565b6102365f81565b60055468010000000000000000900467ffffffffffffffff165b60405167ffffffffffffffff90911681526020016101f4565b60055467ffffffffffffffff1661048e565b6102216104c736600461236e565b610fca565b61048e7f000000000000000000000000000000000000000000000000000000000000000081565b6101e8610501366004612445565b611104565b610202611248565b6102216112e5565b6001546040805173ffffffffffffffffffffffffffffffffffffffff831681527401000000000000000000000000000000000000000090920465ffffffffffff166020830152016101f4565b6102216105703660046123ad565b611341565b610221611382565b5f7fffffffff0000000000000000000000000000000000000000000000000000000082167fb35efc0c0000000000000000000000000000000000000000000000000000000014806105d257506105d282611394565b92915050565b5f6105e2816113e9565b6105ea6113f3565b50565b5f7f4a8e1518fc3b69d274d0d6c14a98dd84f4bfb1e39bbbe43f417772dd939cd808610618816113e9565b600554700100000000000000000000000000000000900460ff1615610669576040517f712a4e2f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f61075b61067c368690038601866124f0565b80516020808301516040808501516060860151608087015160a09097015183517f0300000000000000000000000000000000000000000000000000000000000000818801527e0100000000000000000000000000000000000000000000000000000000000060218201527fffffffffffffffff00000000000000000000000000000000000000000000000060c098891b81166023830152602b820196909652604b81019390935290861b8416606b83015295851b831660738201529490931b16607b840152815180840360630181526083909301909152815191012090565b905067ffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016610795602086018661242c565b67ffffffffffffffff16101561081e577f00000000000000000000000000000000000000000000000000000000000000006107d3602086018661242c565b6040517facc3e9cb00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff9283166004820152911660248201526044015b60405180910390fd5b60055467ffffffffffffffff16610838602086018661242c565b67ffffffffffffffff161015610934575f600381610859602088018861242c565b67ffffffffffffffff1667ffffffffffffffff1681526020019081526020015f206001015490508181036108da57610894602086018661242c565b6040517fab455a2500000000000000000000000000000000000000000000000000000000815267ffffffffffffffff909116600482015260248101839052604401610815565b6108e7602086018661242c565b6040517f6d613c3900000000000000000000000000000000000000000000000000000000815267ffffffffffffffff90911660048201526024810182905260448101839052606401610815565b60055467ffffffffffffffff1661094e602086018661242c565b67ffffffffffffffff1611156109be5760055467ffffffffffffffff16610978602086018661242c565b6040517f7a8d061d00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff928316600482015291166024820152604401610815565b600454846020013514610a0b57600480546040517f0796b55a0000000000000000000000000000000000000000000000000000000081529182015260208501356024820152604401610815565b610a14846113ff565b6040805160c08101825285820135815260208101839052908101610a3e608087016060880161242c565b67ffffffffffffffff168152602001610a5d60a087016080880161242c565b67ffffffffffffffff168152602001610a7c60c0870160a0880161242c565b67ffffffffffffffff168152602001610a94426116a6565b67ffffffffffffffff16905260035f610ab0602088018861242c565b67ffffffffffffffff908116825260208083019390935260409182015f208451815584840151600182015591840151600290920180546060860151608087015160a090970151841678010000000000000000000000000000000000000000000000000277ffffffffffffffffffffffffffffffffffffffffffffffff97851670010000000000000000000000000000000002979097166fffffffffffffffffffffffffffffffff91851668010000000000000000027fffffffffffffffffffffffffffffffff000000000000000000000000000000009093169590941694909417179290921617929092179091556004829055610baf9085018561242c565b610bba9060016125d3565b600580547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001667ffffffffffffffff92909216919091179055610c0360c0850160a0860161242c565b6005805467ffffffffffffffff9290921668010000000000000000027fffffffffffffffffffffffffffffffff0000000000000000ffffffffffffffff90921691909117905580610c57602086018661242c565b67ffffffffffffffff167fb7504cd1e32568b16a96dd8fb0da74c649e425006318bb259eca087e629c542d60208701356040880135610c9c60808a0160608b0161242c565b610cac60a08b0160808c0161242c565b610cbc60c08c0160a08d0161242c565b60408051958652602086019490945267ffffffffffffffff928316858501529082166060850152166080830152519081900360a00190a39392505050565b81610d31576040517f3fc3c27a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610d3b82826116f9565b5050565b81158015610d67575060025473ffffffffffffffffffffffffffffffffffffffff8281169116145b15610e3a5760015473ffffffffffffffffffffffffffffffffffffffff81169074010000000000000000000000000000000000000000900465ffffffffffff1681151580610dbb575065ffffffffffff8116155b80610dce57504265ffffffffffff821610155b15610e0f576040517f19ca5ebb00000000000000000000000000000000000000000000000000000000815265ffffffffffff82166004820152602401610815565b5050600180547fffffffffffff000000000000ffffffffffffffffffffffffffffffffffffffff1690555b610d3b8282611723565b5f610e4e816113e9565b610d3b82611781565b5f610e61816113e9565b610d3b826117f9565b6040805160c0810182525f80825260208201819052918101829052606081018290526080810182905260a0810191909152610ea482611868565b6040805160c081018252825481526001830154602082015260029092015467ffffffffffffffff8082169284019290925268010000000000000000810482166060840152700100000000000000000000000000000000810482166080840152780100000000000000000000000000000000000000000000000090041660a082015292915050565b5f610f4b60025473ffffffffffffffffffffffffffffffffffffffff1690565b905090565b6002545f907a010000000000000000000000000000000000000000000000000000900465ffffffffffff168015158015610f9257504265ffffffffffff821610155b610f9d575f5f610fc2565b60025474010000000000000000000000000000000000000000900465ffffffffffff16815b915091509091565b5f610fd4816113e9565b600554700100000000000000000000000000000000900460ff1615611025576040517f712a4e2f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600454821461106b57600480546040517f6f2de3c70000000000000000000000000000000000000000000000000000000081529182015260248101839052604401610815565b600580547001000000000000000000000000000000007fffffffffffffffffffffffffffffff00ffffffffffffffffffffffffffffffff821617918290556004546040805167ffffffffffffffff9384168152680100000000000000009094049092166020840152917f79b7f4b971d94514f2584fdbd244ea1c7fbac6672ab3e05044cb71c34f9752d991015b60405180910390a25050565b5f5f61110f88611868565b600281015490915067ffffffffffffffff68010000000000000000909104811690861610611140575f91505061123e565b60028101546111649068010000000000000000900467ffffffffffffffff16611919565b60ff168314611176575f91505061123e565b600281015467ffffffffffffffff1661119087868661197e565b1461119e575f91505061123e565b80546040805180820182525f80825260209182018190528251808401845260608101919091527fffffffffffffffff00000000000000000000000000000000000000000000000060c08d811b82166061840152606983018d90528b901b1660898201528251808203607101815260918201909352825192820192909220825267ffffffffffffffff891690820152611238908787876119e5565b51149150505b9695505050505050565b6002545f907a010000000000000000000000000000000000000000000000000000900465ffffffffffff16801515801561128957504265ffffffffffff8216105b6112bb576001547a010000000000000000000000000000000000000000000000000000900465ffffffffffff166112df565b60025474010000000000000000000000000000000000000000900465ffffffffffff165b91505090565b60015473ffffffffffffffffffffffffffffffffffffffff16338114611339576040517fc22c8022000000000000000000000000000000000000000000000000000000008152336004820152602401610815565b6105ea611a9c565b81611378576040517f3fc3c27a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610d3b8282611b8d565b5f61138c816113e9565b6105ea611bb1565b5f7fffffffff0000000000000000000000000000000000000000000000000000000082167f314987860000000000000000000000000000000000000000000000000000000014806105d257506105d282611bbb565b6105ea8133611c51565b6113fd5f5f611cd6565b565b64010000000061141560a083016080840161242c565b67ffffffffffffffff16111561147f5761143560a082016080830161242c565b6040517f09e51df900000000000000000000000000000000000000000000000000000000815267ffffffffffffffff90911660048201526401000000006024820152604401610815565b61148f60a082016080830161242c565b67ffffffffffffffff1615801561152857506040805180820182525f8082526020918201819052825180840184527f020000000000000000000000000000000000000000000000000000000000000060608201528351604181830301815260618201855280519084012080825292015290820135141580611528575061151b608082016060830161242c565b67ffffffffffffffff1615155b15611587576040810135611542608083016060840161242c565b6040517fde718304000000000000000000000000000000000000000000000000000000008152600481019290925267ffffffffffffffff166024820152604401610815565b60055468010000000000000000900467ffffffffffffffff166115b060c0830160a0840161242c565b67ffffffffffffffff1611611626576115cf60c0820160a0830161242c565b6005546040517f81ed60e900000000000000000000000000000000000000000000000000000000815267ffffffffffffffff9283166004820152680100000000000000009091049091166024820152604401610815565b4261163760c0830160a0840161242c565b67ffffffffffffffff1611156105ea5761165760c0820160a0830161242c565b611660426116a6565b6040517f18aaced200000000000000000000000000000000000000000000000000000000815267ffffffffffffffff928316600482015291166024820152604401610815565b5f67ffffffffffffffff8211156116f557604080517f6dfcc650000000000000000000000000000000000000000000000000000000008152600481019190915260248101839052604401610815565b5090565b5f82815260208190526040902060010154611713816113e9565b61171d8383611e2f565b50505050565b73ffffffffffffffffffffffffffffffffffffffff81163314611772576040517f6697b23200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61177c8282611eed565b505050565b5f61178a611248565b61179342611f4e565b61179d91906125f3565b90506117a98282611f99565b60405165ffffffffffff8216815273ffffffffffffffffffffffffffffffffffffffff8316907f3377dc44241e779dd06afab5b788a35ca5f3b778836e2990bdb26a2a4b2e5ed6906020016110f8565b5f61180382612034565b61180c42611f4e565b61181691906125f3565b90506118228282611cd6565b6040805165ffffffffffff8085168252831660208201527ff1038c18cf84a56e432fdbfaf746924b7ea511dfe03a6506a0ceba4888788d9b910160405180910390a15050565b5f7f000000000000000000000000000000000000000000000000000000000000000067ffffffffffffffff168267ffffffffffffffff1610806118bb575060055467ffffffffffffffff90811690831610155b156118fe576040517fa4e21e6c00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff83166004820152602401610815565b5067ffffffffffffffff165f90815260036020526040902090565b5f60018267ffffffffffffffff161161193357505f919050565b5f61193f600184612611565b90505f5b67ffffffffffffffff8216156119775760018267ffffffffffffffff16901c9150808061196f90612631565b915050611943565b9392505050565b5f67ffffffffffffffff8416815b838110156119dc578484828181106119a6576119a661264f565b90506040020160200160208101906119be919061242c565b6119d29067ffffffffffffffff168361267c565b915060010161198c565b50949350505050565b604080518082019091525f808252602082015284845f5b84811015611a9057600182165f03611a4657611a3f83878784818110611a2457611a2461264f565b905060400201803603810190611a3a919061268f565b612085565b9250611a7a565b611a77868683818110611a5b57611a5b61264f565b905060400201803603810190611a71919061268f565b84612085565b92505b677fffffffffffffff600192831c1691016119fc565b50909695505050505050565b60015473ffffffffffffffffffffffffffffffffffffffff81169074010000000000000000000000000000000000000000900465ffffffffffff16801580611aec57504265ffffffffffff821610155b15611b2d576040517f19ca5ebb00000000000000000000000000000000000000000000000000000000815265ffffffffffff82166004820152602401610815565b611b555f611b5060025473ffffffffffffffffffffffffffffffffffffffff1690565b611eed565b50611b605f83611e2f565b5050600180547fffffffffffff000000000000000000000000000000000000000000000000000016905550565b5f82815260208190526040902060010154611ba7816113e9565b61171d8383611eed565b6113fd5f5f611f99565b5f7fffffffff0000000000000000000000000000000000000000000000000000000082167f7965db0b0000000000000000000000000000000000000000000000000000000014806105d257507f01ffc9a7000000000000000000000000000000000000000000000000000000007fffffffff000000000000000000000000000000000000000000000000000000008316146105d2565b5f8281526020818152604080832073ffffffffffffffffffffffffffffffffffffffff8516845290915290205460ff16610d3b576040517fe2517d3f00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8216600482015260248101839052604401610815565b6002547a010000000000000000000000000000000000000000000000000000900465ffffffffffff168015611daa574265ffffffffffff82161015611d81576002546001805479ffffffffffffffffffffffffffffffffffffffffffffffffffff167401000000000000000000000000000000000000000090920465ffffffffffff167a01000000000000000000000000000000000000000000000000000002919091179055611daa565b6040517f2b1fa2edafe6f7b9e97c1a9e0c3660e645beb2dcaa2d45bdbf9beaf5472e1ec5905f90a15b506002805473ffffffffffffffffffffffffffffffffffffffff167401000000000000000000000000000000000000000065ffffffffffff9485160279ffffffffffffffffffffffffffffffffffffffffffffffffffff16177a0100000000000000000000000000000000000000000000000000009290931691909102919091179055565b5f82611ee3575f611e5560025473ffffffffffffffffffffffffffffffffffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff1614611ea2576040517f3fc3c27a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600280547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff84161790555b6119778383612164565b5f82158015611f16575060025473ffffffffffffffffffffffffffffffffffffffff8381169116145b15611f4457600280547fffffffffffffffffffffffff00000000000000000000000000000000000000001690555b611977838361225d565b5f65ffffffffffff8211156116f5576040517f6dfcc6500000000000000000000000000000000000000000000000000000000081526030600482015260248101839052604401610815565b600180547401000000000000000000000000000000000000000065ffffffffffff84811682027fffffffffffff0000000000000000000000000000000000000000000000000000841673ffffffffffffffffffffffffffffffffffffffff88161717909355900416801561177c576040517f8886ebfc4259abdbc16601dd8fb5678e54878f47b3c34836cfc51154a9605109905f90a1505050565b5f5f61203e611248565b90508065ffffffffffff168365ffffffffffff1611612066576120618382612708565b611977565b61197765ffffffffffff8416620697805f828218828410028218611977565b604080518082019091525f80825260208201526040805180820190915283516020808601518551918601517f0100000000000000000000000000000000000000000000000000000000000000606086015260618501939093527fffffffffffffffff00000000000000000000000000000000000000000000000060c091821b8116608186015260898501929092529190911b1660a98201528060b181016040516020818303038152906040528051906020012081526020018360200151856020015161215191906125d3565b67ffffffffffffffff1690529392505050565b5f8281526020818152604080832073ffffffffffffffffffffffffffffffffffffffff8516845290915281205460ff16612256575f8381526020818152604080832073ffffffffffffffffffffffffffffffffffffffff86168452909152902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660011790556121f43390565b73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16847f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a45060016105d2565b505f6105d2565b5f8281526020818152604080832073ffffffffffffffffffffffffffffffffffffffff8516845290915281205460ff1615612256575f8381526020818152604080832073ffffffffffffffffffffffffffffffffffffffff8616808552925280832080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016905551339286917ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9190a45060016105d2565b5f60208284031215612326575f5ffd5b81357fffffffff0000000000000000000000000000000000000000000000000000000081168114611977575f5ffd5b5f60c0828403128015612366575f5ffd5b509092915050565b5f6020828403121561237e575f5ffd5b5035919050565b803573ffffffffffffffffffffffffffffffffffffffff811681146123a8575f5ffd5b919050565b5f5f604083850312156123be575f5ffd5b823591506123ce60208401612385565b90509250929050565b5f602082840312156123e7575f5ffd5b61197782612385565b5f60208284031215612400575f5ffd5b813565ffffffffffff81168114611977575f5ffd5b803567ffffffffffffffff811681146123a8575f5ffd5b5f6020828403121561243c575f5ffd5b61197782612415565b5f5f5f5f5f5f60a0878903121561245a575f5ffd5b61246387612415565b95506020870135945061247860408801612415565b935061248660608801612415565b9250608087013567ffffffffffffffff8111156124a1575f5ffd5b8701601f810189136124b1575f5ffd5b803567ffffffffffffffff8111156124c7575f5ffd5b8960208260061b84010111156124db575f5ffd5b60208201935080925050509295509295509295565b5f60c0828403128015612501575f5ffd5b5060405160c0810167ffffffffffffffff8111828210171561254a577f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b60405261255683612415565b8152602083810135908201526040808401359082015261257860608401612415565b606082015261258960808401612415565b608082015261259a60a08401612415565b60a08201529392505050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b67ffffffffffffffff81811683821601908111156105d2576105d26125a6565b65ffffffffffff81811683821601908111156105d2576105d26125a6565b67ffffffffffffffff82811682821603908111156105d2576105d26125a6565b5f60ff821660ff8103612646576126466125a6565b60010192915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b808201808211156105d2576105d26125a6565b5f60408284031280156126a0575f5ffd5b506040805190810167ffffffffffffffff811182821017156126e9577f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b604052823581526126fc60208401612415565b60208201529392505050565b65ffffffffffff82811682821603908111156105d2576105d26125a656",
}

// LiabilityLedger is an auto generated Go binding around an Ethereum contract.
type LiabilityLedger struct {
	abi abi.ABI
}

// GetABI returns the ABI associated with this contract binding.
func (c *LiabilityLedger) GetABI() abi.ABI {
	return c.abi
}

// NewLiabilityLedger creates a new instance of LiabilityLedger.
func NewLiabilityLedger() *LiabilityLedger {
	parsed, err := LiabilityLedgerMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &LiabilityLedger{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *LiabilityLedger) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackConstructor is the Go binding used to pack the parameters required for
// contract deployment.
//
// Solidity: constructor(address admin, uint48 adminTransferDelay, address anchorer, uint64 genesisEpoch_, bytes32 genesisDigest_, uint64 genesisAsOf) returns()
func (liabilityLedger *LiabilityLedger) PackConstructor(admin common.Address, adminTransferDelay *big.Int, anchorer common.Address, genesisEpoch_ uint64, genesisDigest_ [32]byte, genesisAsOf uint64) []byte {
	enc, err := liabilityLedger.abi.Pack("", admin, adminTransferDelay, anchorer, genesisEpoch_, genesisDigest_, genesisAsOf)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackANCHORERROLE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7cfd9a51.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function ANCHORER_ROLE() view returns(bytes32)
func (liabilityLedger *LiabilityLedger) PackANCHORERROLE() []byte {
	enc, err := liabilityLedger.abi.Pack("ANCHORER_ROLE")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackANCHORERROLE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7cfd9a51.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function ANCHORER_ROLE() view returns(bytes32)
func (liabilityLedger *LiabilityLedger) TryPackANCHORERROLE() ([]byte, error) {
	return liabilityLedger.abi.Pack("ANCHORER_ROLE")
}

// UnpackANCHORERROLE is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x7cfd9a51.
//
// Solidity: function ANCHORER_ROLE() view returns(bytes32)
func (liabilityLedger *LiabilityLedger) UnpackANCHORERROLE(data []byte) ([32]byte, error) {
	out, err := liabilityLedger.abi.Unpack("ANCHORER_ROLE", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackDEFAULTADMINROLE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa217fddf.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (liabilityLedger *LiabilityLedger) PackDEFAULTADMINROLE() []byte {
	enc, err := liabilityLedger.abi.Pack("DEFAULT_ADMIN_ROLE")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackDEFAULTADMINROLE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa217fddf.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (liabilityLedger *LiabilityLedger) TryPackDEFAULTADMINROLE() ([]byte, error) {
	return liabilityLedger.abi.Pack("DEFAULT_ADMIN_ROLE")
}

// UnpackDEFAULTADMINROLE is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (liabilityLedger *LiabilityLedger) UnpackDEFAULTADMINROLE(data []byte) ([32]byte, error) {
	out, err := liabilityLedger.abi.Unpack("DEFAULT_ADMIN_ROLE", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackAcceptDefaultAdminTransfer is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xcefc1429.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function acceptDefaultAdminTransfer() returns()
func (liabilityLedger *LiabilityLedger) PackAcceptDefaultAdminTransfer() []byte {
	enc, err := liabilityLedger.abi.Pack("acceptDefaultAdminTransfer")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAcceptDefaultAdminTransfer is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xcefc1429.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function acceptDefaultAdminTransfer() returns()
func (liabilityLedger *LiabilityLedger) TryPackAcceptDefaultAdminTransfer() ([]byte, error) {
	return liabilityLedger.abi.Pack("acceptDefaultAdminTransfer")
}

// PackAnchor is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1ffcf5d4.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function anchor((uint64,bytes32,bytes32,uint64,uint64,uint64) header) returns(bytes32)
func (liabilityLedger *LiabilityLedger) PackAnchor(header EpochHeader) []byte {
	enc, err := liabilityLedger.abi.Pack("anchor", header)
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
func (liabilityLedger *LiabilityLedger) TryPackAnchor(header EpochHeader) ([]byte, error) {
	return liabilityLedger.abi.Pack("anchor", header)
}

// UnpackAnchor is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x1ffcf5d4.
//
// Solidity: function anchor((uint64,bytes32,bytes32,uint64,uint64,uint64) header) returns(bytes32)
func (liabilityLedger *LiabilityLedger) UnpackAnchor(data []byte) ([32]byte, error) {
	out, err := liabilityLedger.abi.Unpack("anchor", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackBeginDefaultAdminTransfer is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x634e93da.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function beginDefaultAdminTransfer(address newAdmin) returns()
func (liabilityLedger *LiabilityLedger) PackBeginDefaultAdminTransfer(newAdmin common.Address) []byte {
	enc, err := liabilityLedger.abi.Pack("beginDefaultAdminTransfer", newAdmin)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackBeginDefaultAdminTransfer is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x634e93da.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function beginDefaultAdminTransfer(address newAdmin) returns()
func (liabilityLedger *LiabilityLedger) TryPackBeginDefaultAdminTransfer(newAdmin common.Address) ([]byte, error) {
	return liabilityLedger.abi.Pack("beginDefaultAdminTransfer", newAdmin)
}

// PackCancelDefaultAdminTransfer is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd602b9fd.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function cancelDefaultAdminTransfer() returns()
func (liabilityLedger *LiabilityLedger) PackCancelDefaultAdminTransfer() []byte {
	enc, err := liabilityLedger.abi.Pack("cancelDefaultAdminTransfer")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCancelDefaultAdminTransfer is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd602b9fd.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function cancelDefaultAdminTransfer() returns()
func (liabilityLedger *LiabilityLedger) TryPackCancelDefaultAdminTransfer() ([]byte, error) {
	return liabilityLedger.abi.Pack("cancelDefaultAdminTransfer")
}

// PackChangeDefaultAdminDelay is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x649a5ec7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function changeDefaultAdminDelay(uint48 newDelay) returns()
func (liabilityLedger *LiabilityLedger) PackChangeDefaultAdminDelay(newDelay *big.Int) []byte {
	enc, err := liabilityLedger.abi.Pack("changeDefaultAdminDelay", newDelay)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackChangeDefaultAdminDelay is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x649a5ec7.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function changeDefaultAdminDelay(uint48 newDelay) returns()
func (liabilityLedger *LiabilityLedger) TryPackChangeDefaultAdminDelay(newDelay *big.Int) ([]byte, error) {
	return liabilityLedger.abi.Pack("changeDefaultAdminDelay", newDelay)
}

// PackDefaultAdmin is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x84ef8ffc.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function defaultAdmin() view returns(address)
func (liabilityLedger *LiabilityLedger) PackDefaultAdmin() []byte {
	enc, err := liabilityLedger.abi.Pack("defaultAdmin")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackDefaultAdmin is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x84ef8ffc.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function defaultAdmin() view returns(address)
func (liabilityLedger *LiabilityLedger) TryPackDefaultAdmin() ([]byte, error) {
	return liabilityLedger.abi.Pack("defaultAdmin")
}

// UnpackDefaultAdmin is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x84ef8ffc.
//
// Solidity: function defaultAdmin() view returns(address)
func (liabilityLedger *LiabilityLedger) UnpackDefaultAdmin(data []byte) (common.Address, error) {
	out, err := liabilityLedger.abi.Unpack("defaultAdmin", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackDefaultAdminDelay is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xcc8463c8.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function defaultAdminDelay() view returns(uint48)
func (liabilityLedger *LiabilityLedger) PackDefaultAdminDelay() []byte {
	enc, err := liabilityLedger.abi.Pack("defaultAdminDelay")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackDefaultAdminDelay is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xcc8463c8.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function defaultAdminDelay() view returns(uint48)
func (liabilityLedger *LiabilityLedger) TryPackDefaultAdminDelay() ([]byte, error) {
	return liabilityLedger.abi.Pack("defaultAdminDelay")
}

// UnpackDefaultAdminDelay is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xcc8463c8.
//
// Solidity: function defaultAdminDelay() view returns(uint48)
func (liabilityLedger *LiabilityLedger) UnpackDefaultAdminDelay(data []byte) (*big.Int, error) {
	out, err := liabilityLedger.abi.Unpack("defaultAdminDelay", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackDefaultAdminDelayIncreaseWait is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x022d63fb.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function defaultAdminDelayIncreaseWait() view returns(uint48)
func (liabilityLedger *LiabilityLedger) PackDefaultAdminDelayIncreaseWait() []byte {
	enc, err := liabilityLedger.abi.Pack("defaultAdminDelayIncreaseWait")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackDefaultAdminDelayIncreaseWait is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x022d63fb.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function defaultAdminDelayIncreaseWait() view returns(uint48)
func (liabilityLedger *LiabilityLedger) TryPackDefaultAdminDelayIncreaseWait() ([]byte, error) {
	return liabilityLedger.abi.Pack("defaultAdminDelayIncreaseWait")
}

// UnpackDefaultAdminDelayIncreaseWait is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x022d63fb.
//
// Solidity: function defaultAdminDelayIncreaseWait() view returns(uint48)
func (liabilityLedger *LiabilityLedger) UnpackDefaultAdminDelayIncreaseWait(data []byte) (*big.Int, error) {
	out, err := liabilityLedger.abi.Unpack("defaultAdminDelayIncreaseWait", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackEpochAt is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7d60d355.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function epochAt(uint64 number) view returns((bytes32,bytes32,uint64,uint64,uint64,uint64))
func (liabilityLedger *LiabilityLedger) PackEpochAt(number uint64) []byte {
	enc, err := liabilityLedger.abi.Pack("epochAt", number)
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
func (liabilityLedger *LiabilityLedger) TryPackEpochAt(number uint64) ([]byte, error) {
	return liabilityLedger.abi.Pack("epochAt", number)
}

// UnpackEpochAt is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x7d60d355.
//
// Solidity: function epochAt(uint64 number) view returns((bytes32,bytes32,uint64,uint64,uint64,uint64))
func (liabilityLedger *LiabilityLedger) UnpackEpochAt(data []byte) (Epoch, error) {
	out, err := liabilityLedger.abi.Unpack("epochAt", data)
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
func (liabilityLedger *LiabilityLedger) PackGenesisDigest() []byte {
	enc, err := liabilityLedger.abi.Pack("genesisDigest")
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
func (liabilityLedger *LiabilityLedger) TryPackGenesisDigest() ([]byte, error) {
	return liabilityLedger.abi.Pack("genesisDigest")
}

// UnpackGenesisDigest is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x79ac2661.
//
// Solidity: function genesisDigest() view returns(bytes32)
func (liabilityLedger *LiabilityLedger) UnpackGenesisDigest(data []byte) ([32]byte, error) {
	out, err := liabilityLedger.abi.Unpack("genesisDigest", data)
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
func (liabilityLedger *LiabilityLedger) PackGenesisEpoch() []byte {
	enc, err := liabilityLedger.abi.Pack("genesisEpoch")
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
func (liabilityLedger *LiabilityLedger) TryPackGenesisEpoch() ([]byte, error) {
	return liabilityLedger.abi.Pack("genesisEpoch")
}

// UnpackGenesisEpoch is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb70e6be6.
//
// Solidity: function genesisEpoch() view returns(uint64)
func (liabilityLedger *LiabilityLedger) UnpackGenesisEpoch(data []byte) (uint64, error) {
	out, err := liabilityLedger.abi.Unpack("genesisEpoch", data)
	if err != nil {
		return *new(uint64), err
	}
	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)
	return out0, nil
}

// PackGetRoleAdmin is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x248a9ca3.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (liabilityLedger *LiabilityLedger) PackGetRoleAdmin(role [32]byte) []byte {
	enc, err := liabilityLedger.abi.Pack("getRoleAdmin", role)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetRoleAdmin is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x248a9ca3.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (liabilityLedger *LiabilityLedger) TryPackGetRoleAdmin(role [32]byte) ([]byte, error) {
	return liabilityLedger.abi.Pack("getRoleAdmin", role)
}

// UnpackGetRoleAdmin is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (liabilityLedger *LiabilityLedger) UnpackGetRoleAdmin(data []byte) ([32]byte, error) {
	out, err := liabilityLedger.abi.Unpack("getRoleAdmin", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackGrantRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2f2ff15d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (liabilityLedger *LiabilityLedger) PackGrantRole(role [32]byte, account common.Address) []byte {
	enc, err := liabilityLedger.abi.Pack("grantRole", role, account)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGrantRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2f2ff15d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (liabilityLedger *LiabilityLedger) TryPackGrantRole(role [32]byte, account common.Address) ([]byte, error) {
	return liabilityLedger.abi.Pack("grantRole", role, account)
}

// PackHasRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x91d14854.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (liabilityLedger *LiabilityLedger) PackHasRole(role [32]byte, account common.Address) []byte {
	enc, err := liabilityLedger.abi.Pack("hasRole", role, account)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackHasRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x91d14854.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (liabilityLedger *LiabilityLedger) TryPackHasRole(role [32]byte, account common.Address) ([]byte, error) {
	return liabilityLedger.abi.Pack("hasRole", role, account)
}

// UnpackHasRole is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (liabilityLedger *LiabilityLedger) UnpackHasRole(data []byte) (bool, error) {
	out, err := liabilityLedger.abi.Unpack("hasRole", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackHead is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8f7dcfa3.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function head() view returns(bytes32)
func (liabilityLedger *LiabilityLedger) PackHead() []byte {
	enc, err := liabilityLedger.abi.Pack("head")
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
func (liabilityLedger *LiabilityLedger) TryPackHead() ([]byte, error) {
	return liabilityLedger.abi.Pack("head")
}

// UnpackHead is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8f7dcfa3.
//
// Solidity: function head() view returns(bytes32)
func (liabilityLedger *LiabilityLedger) UnpackHead(data []byte) ([32]byte, error) {
	out, err := liabilityLedger.abi.Unpack("head", data)
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
func (liabilityLedger *LiabilityLedger) PackIsSealed() []byte {
	enc, err := liabilityLedger.abi.Pack("isSealed")
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
func (liabilityLedger *LiabilityLedger) TryPackIsSealed() ([]byte, error) {
	return liabilityLedger.abi.Pack("isSealed")
}

// UnpackIsSealed is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x631f9852.
//
// Solidity: function isSealed() view returns(bool)
func (liabilityLedger *LiabilityLedger) UnpackIsSealed(data []byte) (bool, error) {
	out, err := liabilityLedger.abi.Unpack("isSealed", data)
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
func (liabilityLedger *LiabilityLedger) PackLastAsOf() []byte {
	enc, err := liabilityLedger.abi.Pack("lastAsOf")
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
func (liabilityLedger *LiabilityLedger) TryPackLastAsOf() ([]byte, error) {
	return liabilityLedger.abi.Pack("lastAsOf")
}

// UnpackLastAsOf is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa89210cb.
//
// Solidity: function lastAsOf() view returns(uint64)
func (liabilityLedger *LiabilityLedger) UnpackLastAsOf(data []byte) (uint64, error) {
	out, err := liabilityLedger.abi.Unpack("lastAsOf", data)
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
func (liabilityLedger *LiabilityLedger) PackNextEpoch() []byte {
	enc, err := liabilityLedger.abi.Pack("nextEpoch")
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
func (liabilityLedger *LiabilityLedger) TryPackNextEpoch() ([]byte, error) {
	return liabilityLedger.abi.Pack("nextEpoch")
}

// UnpackNextEpoch is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xaea0e78b.
//
// Solidity: function nextEpoch() view returns(uint64)
func (liabilityLedger *LiabilityLedger) UnpackNextEpoch(data []byte) (uint64, error) {
	out, err := liabilityLedger.abi.Unpack("nextEpoch", data)
	if err != nil {
		return *new(uint64), err
	}
	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)
	return out0, nil
}

// PackOwner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8da5cb5b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function owner() view returns(address)
func (liabilityLedger *LiabilityLedger) PackOwner() []byte {
	enc, err := liabilityLedger.abi.Pack("owner")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackOwner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8da5cb5b.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function owner() view returns(address)
func (liabilityLedger *LiabilityLedger) TryPackOwner() ([]byte, error) {
	return liabilityLedger.abi.Pack("owner")
}

// UnpackOwner is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (liabilityLedger *LiabilityLedger) UnpackOwner(data []byte) (common.Address, error) {
	out, err := liabilityLedger.abi.Unpack("owner", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackPendingDefaultAdmin is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xcf6eefb7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function pendingDefaultAdmin() view returns(address newAdmin, uint48 schedule)
func (liabilityLedger *LiabilityLedger) PackPendingDefaultAdmin() []byte {
	enc, err := liabilityLedger.abi.Pack("pendingDefaultAdmin")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPendingDefaultAdmin is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xcf6eefb7.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function pendingDefaultAdmin() view returns(address newAdmin, uint48 schedule)
func (liabilityLedger *LiabilityLedger) TryPackPendingDefaultAdmin() ([]byte, error) {
	return liabilityLedger.abi.Pack("pendingDefaultAdmin")
}

// PendingDefaultAdminOutput serves as a container for the return parameters of contract
// method PendingDefaultAdmin.
type PendingDefaultAdminOutput struct {
	NewAdmin common.Address
	Schedule *big.Int
}

// UnpackPendingDefaultAdmin is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xcf6eefb7.
//
// Solidity: function pendingDefaultAdmin() view returns(address newAdmin, uint48 schedule)
func (liabilityLedger *LiabilityLedger) UnpackPendingDefaultAdmin(data []byte) (PendingDefaultAdminOutput, error) {
	out, err := liabilityLedger.abi.Unpack("pendingDefaultAdmin", data)
	outstruct := new(PendingDefaultAdminOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.NewAdmin = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Schedule = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	return *outstruct, nil
}

// PackPendingDefaultAdminDelay is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa1eda53c.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function pendingDefaultAdminDelay() view returns(uint48 newDelay, uint48 schedule)
func (liabilityLedger *LiabilityLedger) PackPendingDefaultAdminDelay() []byte {
	enc, err := liabilityLedger.abi.Pack("pendingDefaultAdminDelay")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPendingDefaultAdminDelay is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa1eda53c.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function pendingDefaultAdminDelay() view returns(uint48 newDelay, uint48 schedule)
func (liabilityLedger *LiabilityLedger) TryPackPendingDefaultAdminDelay() ([]byte, error) {
	return liabilityLedger.abi.Pack("pendingDefaultAdminDelay")
}

// PendingDefaultAdminDelayOutput serves as a container for the return parameters of contract
// method PendingDefaultAdminDelay.
type PendingDefaultAdminDelayOutput struct {
	NewDelay *big.Int
	Schedule *big.Int
}

// UnpackPendingDefaultAdminDelay is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa1eda53c.
//
// Solidity: function pendingDefaultAdminDelay() view returns(uint48 newDelay, uint48 schedule)
func (liabilityLedger *LiabilityLedger) UnpackPendingDefaultAdminDelay(data []byte) (PendingDefaultAdminDelayOutput, error) {
	out, err := liabilityLedger.abi.Unpack("pendingDefaultAdminDelay", data)
	outstruct := new(PendingDefaultAdminDelayOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.NewDelay = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.Schedule = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	return *outstruct, nil
}

// PackRenounceRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x36568abe.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (liabilityLedger *LiabilityLedger) PackRenounceRole(role [32]byte, account common.Address) []byte {
	enc, err := liabilityLedger.abi.Pack("renounceRole", role, account)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRenounceRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x36568abe.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (liabilityLedger *LiabilityLedger) TryPackRenounceRole(role [32]byte, account common.Address) ([]byte, error) {
	return liabilityLedger.abi.Pack("renounceRole", role, account)
}

// PackRevokeRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd547741f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (liabilityLedger *LiabilityLedger) PackRevokeRole(role [32]byte, account common.Address) []byte {
	enc, err := liabilityLedger.abi.Pack("revokeRole", role, account)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRevokeRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd547741f.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (liabilityLedger *LiabilityLedger) TryPackRevokeRole(role [32]byte, account common.Address) ([]byte, error) {
	return liabilityLedger.abi.Pack("revokeRole", role, account)
}

// PackRollbackDefaultAdminDelay is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0aa6220b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function rollbackDefaultAdminDelay() returns()
func (liabilityLedger *LiabilityLedger) PackRollbackDefaultAdminDelay() []byte {
	enc, err := liabilityLedger.abi.Pack("rollbackDefaultAdminDelay")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRollbackDefaultAdminDelay is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0aa6220b.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function rollbackDefaultAdminDelay() returns()
func (liabilityLedger *LiabilityLedger) TryPackRollbackDefaultAdminDelay() ([]byte, error) {
	return liabilityLedger.abi.Pack("rollbackDefaultAdminDelay")
}

// PackScheme is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8518b0ad.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function scheme() pure returns(uint16)
func (liabilityLedger *LiabilityLedger) PackScheme() []byte {
	enc, err := liabilityLedger.abi.Pack("scheme")
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
func (liabilityLedger *LiabilityLedger) TryPackScheme() ([]byte, error) {
	return liabilityLedger.abi.Pack("scheme")
}

// UnpackScheme is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8518b0ad.
//
// Solidity: function scheme() pure returns(uint16)
func (liabilityLedger *LiabilityLedger) UnpackScheme(data []byte) (uint16, error) {
	out, err := liabilityLedger.abi.Unpack("scheme", data)
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
func (liabilityLedger *LiabilityLedger) PackSeal(expectedHead [32]byte) []byte {
	enc, err := liabilityLedger.abi.Pack("seal", expectedHead)
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
func (liabilityLedger *LiabilityLedger) TryPackSeal(expectedHead [32]byte) ([]byte, error) {
	return liabilityLedger.abi.Pack("seal", expectedHead)
}

// PackSupportsInterface is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x01ffc9a7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (liabilityLedger *LiabilityLedger) PackSupportsInterface(interfaceId [4]byte) []byte {
	enc, err := liabilityLedger.abi.Pack("supportsInterface", interfaceId)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSupportsInterface is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x01ffc9a7.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (liabilityLedger *LiabilityLedger) TryPackSupportsInterface(interfaceId [4]byte) ([]byte, error) {
	return liabilityLedger.abi.Pack("supportsInterface", interfaceId)
}

// UnpackSupportsInterface is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (liabilityLedger *LiabilityLedger) UnpackSupportsInterface(data []byte) (bool, error) {
	out, err := liabilityLedger.abi.Unpack("supportsInterface", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackVerifyInclusion is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc0566abe.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function verifyInclusion(uint64 number, bytes32 commitment, uint64 amount, uint64 index, (bytes32,uint64)[] siblings) view returns(bool)
func (liabilityLedger *LiabilityLedger) PackVerifyInclusion(number uint64, commitment [32]byte, amount uint64, index uint64, siblings []SumNode) []byte {
	enc, err := liabilityLedger.abi.Pack("verifyInclusion", number, commitment, amount, index, siblings)
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
func (liabilityLedger *LiabilityLedger) TryPackVerifyInclusion(number uint64, commitment [32]byte, amount uint64, index uint64, siblings []SumNode) ([]byte, error) {
	return liabilityLedger.abi.Pack("verifyInclusion", number, commitment, amount, index, siblings)
}

// UnpackVerifyInclusion is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc0566abe.
//
// Solidity: function verifyInclusion(uint64 number, bytes32 commitment, uint64 amount, uint64 index, (bytes32,uint64)[] siblings) view returns(bool)
func (liabilityLedger *LiabilityLedger) UnpackVerifyInclusion(data []byte) (bool, error) {
	out, err := liabilityLedger.abi.Unpack("verifyInclusion", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// LiabilityLedgerDefaultAdminDelayChangeCanceled represents a DefaultAdminDelayChangeCanceled event raised by the LiabilityLedger contract.
type LiabilityLedgerDefaultAdminDelayChangeCanceled struct {
	Raw *types.Log // Blockchain specific contextual infos
}

const LiabilityLedgerDefaultAdminDelayChangeCanceledEventName = "DefaultAdminDelayChangeCanceled"

// ContractEventName returns the user-defined event name.
func (LiabilityLedgerDefaultAdminDelayChangeCanceled) ContractEventName() string {
	return LiabilityLedgerDefaultAdminDelayChangeCanceledEventName
}

// UnpackDefaultAdminDelayChangeCanceledEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event DefaultAdminDelayChangeCanceled()
func (liabilityLedger *LiabilityLedger) UnpackDefaultAdminDelayChangeCanceledEvent(log *types.Log) (*LiabilityLedgerDefaultAdminDelayChangeCanceled, error) {
	event := "DefaultAdminDelayChangeCanceled"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != liabilityLedger.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(LiabilityLedgerDefaultAdminDelayChangeCanceled)
	if len(log.Data) > 0 {
		if err := liabilityLedger.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range liabilityLedger.abi.Events[event].Inputs {
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

// LiabilityLedgerDefaultAdminDelayChangeScheduled represents a DefaultAdminDelayChangeScheduled event raised by the LiabilityLedger contract.
type LiabilityLedgerDefaultAdminDelayChangeScheduled struct {
	NewDelay       *big.Int
	EffectSchedule *big.Int
	Raw            *types.Log // Blockchain specific contextual infos
}

const LiabilityLedgerDefaultAdminDelayChangeScheduledEventName = "DefaultAdminDelayChangeScheduled"

// ContractEventName returns the user-defined event name.
func (LiabilityLedgerDefaultAdminDelayChangeScheduled) ContractEventName() string {
	return LiabilityLedgerDefaultAdminDelayChangeScheduledEventName
}

// UnpackDefaultAdminDelayChangeScheduledEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event DefaultAdminDelayChangeScheduled(uint48 newDelay, uint48 effectSchedule)
func (liabilityLedger *LiabilityLedger) UnpackDefaultAdminDelayChangeScheduledEvent(log *types.Log) (*LiabilityLedgerDefaultAdminDelayChangeScheduled, error) {
	event := "DefaultAdminDelayChangeScheduled"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != liabilityLedger.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(LiabilityLedgerDefaultAdminDelayChangeScheduled)
	if len(log.Data) > 0 {
		if err := liabilityLedger.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range liabilityLedger.abi.Events[event].Inputs {
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

// LiabilityLedgerDefaultAdminTransferCanceled represents a DefaultAdminTransferCanceled event raised by the LiabilityLedger contract.
type LiabilityLedgerDefaultAdminTransferCanceled struct {
	Raw *types.Log // Blockchain specific contextual infos
}

const LiabilityLedgerDefaultAdminTransferCanceledEventName = "DefaultAdminTransferCanceled"

// ContractEventName returns the user-defined event name.
func (LiabilityLedgerDefaultAdminTransferCanceled) ContractEventName() string {
	return LiabilityLedgerDefaultAdminTransferCanceledEventName
}

// UnpackDefaultAdminTransferCanceledEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event DefaultAdminTransferCanceled()
func (liabilityLedger *LiabilityLedger) UnpackDefaultAdminTransferCanceledEvent(log *types.Log) (*LiabilityLedgerDefaultAdminTransferCanceled, error) {
	event := "DefaultAdminTransferCanceled"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != liabilityLedger.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(LiabilityLedgerDefaultAdminTransferCanceled)
	if len(log.Data) > 0 {
		if err := liabilityLedger.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range liabilityLedger.abi.Events[event].Inputs {
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

// LiabilityLedgerDefaultAdminTransferScheduled represents a DefaultAdminTransferScheduled event raised by the LiabilityLedger contract.
type LiabilityLedgerDefaultAdminTransferScheduled struct {
	NewAdmin       common.Address
	AcceptSchedule *big.Int
	Raw            *types.Log // Blockchain specific contextual infos
}

const LiabilityLedgerDefaultAdminTransferScheduledEventName = "DefaultAdminTransferScheduled"

// ContractEventName returns the user-defined event name.
func (LiabilityLedgerDefaultAdminTransferScheduled) ContractEventName() string {
	return LiabilityLedgerDefaultAdminTransferScheduledEventName
}

// UnpackDefaultAdminTransferScheduledEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event DefaultAdminTransferScheduled(address indexed newAdmin, uint48 acceptSchedule)
func (liabilityLedger *LiabilityLedger) UnpackDefaultAdminTransferScheduledEvent(log *types.Log) (*LiabilityLedgerDefaultAdminTransferScheduled, error) {
	event := "DefaultAdminTransferScheduled"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != liabilityLedger.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(LiabilityLedgerDefaultAdminTransferScheduled)
	if len(log.Data) > 0 {
		if err := liabilityLedger.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range liabilityLedger.abi.Events[event].Inputs {
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

// LiabilityLedgerEpochAnchored represents a EpochAnchored event raised by the LiabilityLedger contract.
type LiabilityLedgerEpochAnchored struct {
	Number         uint64
	Digest         [32]byte
	PreviousDigest [32]byte
	Root           [32]byte
	Total          uint64
	LeafCount      uint64
	AsOf           uint64
	Raw            *types.Log // Blockchain specific contextual infos
}

const LiabilityLedgerEpochAnchoredEventName = "EpochAnchored"

// ContractEventName returns the user-defined event name.
func (LiabilityLedgerEpochAnchored) ContractEventName() string {
	return LiabilityLedgerEpochAnchoredEventName
}

// UnpackEpochAnchoredEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event EpochAnchored(uint64 indexed number, bytes32 indexed digest, bytes32 previousDigest, bytes32 root, uint64 total, uint64 leafCount, uint64 asOf)
func (liabilityLedger *LiabilityLedger) UnpackEpochAnchoredEvent(log *types.Log) (*LiabilityLedgerEpochAnchored, error) {
	event := "EpochAnchored"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != liabilityLedger.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(LiabilityLedgerEpochAnchored)
	if len(log.Data) > 0 {
		if err := liabilityLedger.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range liabilityLedger.abi.Events[event].Inputs {
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

// LiabilityLedgerLedgerSealed represents a LedgerSealed event raised by the LiabilityLedger contract.
type LiabilityLedgerLedgerSealed struct {
	Head      [32]byte
	NextEpoch uint64
	LastAsOf  uint64
	Raw       *types.Log // Blockchain specific contextual infos
}

const LiabilityLedgerLedgerSealedEventName = "LedgerSealed"

// ContractEventName returns the user-defined event name.
func (LiabilityLedgerLedgerSealed) ContractEventName() string {
	return LiabilityLedgerLedgerSealedEventName
}

// UnpackLedgerSealedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event LedgerSealed(bytes32 indexed head, uint64 nextEpoch, uint64 lastAsOf)
func (liabilityLedger *LiabilityLedger) UnpackLedgerSealedEvent(log *types.Log) (*LiabilityLedgerLedgerSealed, error) {
	event := "LedgerSealed"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != liabilityLedger.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(LiabilityLedgerLedgerSealed)
	if len(log.Data) > 0 {
		if err := liabilityLedger.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range liabilityLedger.abi.Events[event].Inputs {
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

// LiabilityLedgerRoleAdminChanged represents a RoleAdminChanged event raised by the LiabilityLedger contract.
type LiabilityLedgerRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               *types.Log // Blockchain specific contextual infos
}

const LiabilityLedgerRoleAdminChangedEventName = "RoleAdminChanged"

// ContractEventName returns the user-defined event name.
func (LiabilityLedgerRoleAdminChanged) ContractEventName() string {
	return LiabilityLedgerRoleAdminChangedEventName
}

// UnpackRoleAdminChangedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (liabilityLedger *LiabilityLedger) UnpackRoleAdminChangedEvent(log *types.Log) (*LiabilityLedgerRoleAdminChanged, error) {
	event := "RoleAdminChanged"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != liabilityLedger.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(LiabilityLedgerRoleAdminChanged)
	if len(log.Data) > 0 {
		if err := liabilityLedger.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range liabilityLedger.abi.Events[event].Inputs {
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

// LiabilityLedgerRoleGranted represents a RoleGranted event raised by the LiabilityLedger contract.
type LiabilityLedgerRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     *types.Log // Blockchain specific contextual infos
}

const LiabilityLedgerRoleGrantedEventName = "RoleGranted"

// ContractEventName returns the user-defined event name.
func (LiabilityLedgerRoleGranted) ContractEventName() string {
	return LiabilityLedgerRoleGrantedEventName
}

// UnpackRoleGrantedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (liabilityLedger *LiabilityLedger) UnpackRoleGrantedEvent(log *types.Log) (*LiabilityLedgerRoleGranted, error) {
	event := "RoleGranted"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != liabilityLedger.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(LiabilityLedgerRoleGranted)
	if len(log.Data) > 0 {
		if err := liabilityLedger.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range liabilityLedger.abi.Events[event].Inputs {
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

// LiabilityLedgerRoleRevoked represents a RoleRevoked event raised by the LiabilityLedger contract.
type LiabilityLedgerRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     *types.Log // Blockchain specific contextual infos
}

const LiabilityLedgerRoleRevokedEventName = "RoleRevoked"

// ContractEventName returns the user-defined event name.
func (LiabilityLedgerRoleRevoked) ContractEventName() string {
	return LiabilityLedgerRoleRevokedEventName
}

// UnpackRoleRevokedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (liabilityLedger *LiabilityLedger) UnpackRoleRevokedEvent(log *types.Log) (*LiabilityLedgerRoleRevoked, error) {
	event := "RoleRevoked"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != liabilityLedger.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(LiabilityLedgerRoleRevoked)
	if len(log.Data) > 0 {
		if err := liabilityLedger.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range liabilityLedger.abi.Events[event].Inputs {
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
func (liabilityLedger *LiabilityLedger) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], liabilityLedger.abi.Errors["AccessControlBadConfirmation"].ID.Bytes()[:4]) {
		return liabilityLedger.UnpackAccessControlBadConfirmationError(raw[4:])
	}
	if bytes.Equal(raw[:4], liabilityLedger.abi.Errors["AccessControlEnforcedDefaultAdminDelay"].ID.Bytes()[:4]) {
		return liabilityLedger.UnpackAccessControlEnforcedDefaultAdminDelayError(raw[4:])
	}
	if bytes.Equal(raw[:4], liabilityLedger.abi.Errors["AccessControlEnforcedDefaultAdminRules"].ID.Bytes()[:4]) {
		return liabilityLedger.UnpackAccessControlEnforcedDefaultAdminRulesError(raw[4:])
	}
	if bytes.Equal(raw[:4], liabilityLedger.abi.Errors["AccessControlInvalidDefaultAdmin"].ID.Bytes()[:4]) {
		return liabilityLedger.UnpackAccessControlInvalidDefaultAdminError(raw[4:])
	}
	if bytes.Equal(raw[:4], liabilityLedger.abi.Errors["AccessControlUnauthorizedAccount"].ID.Bytes()[:4]) {
		return liabilityLedger.UnpackAccessControlUnauthorizedAccountError(raw[4:])
	}
	if bytes.Equal(raw[:4], liabilityLedger.abi.Errors["EmptyEpochMismatch"].ID.Bytes()[:4]) {
		return liabilityLedger.UnpackEmptyEpochMismatchError(raw[4:])
	}
	if bytes.Equal(raw[:4], liabilityLedger.abi.Errors["EpochAlreadyAnchored"].ID.Bytes()[:4]) {
		return liabilityLedger.UnpackEpochAlreadyAnchoredError(raw[4:])
	}
	if bytes.Equal(raw[:4], liabilityLedger.abi.Errors["EpochBeforeGenesis"].ID.Bytes()[:4]) {
		return liabilityLedger.UnpackEpochBeforeGenesisError(raw[4:])
	}
	if bytes.Equal(raw[:4], liabilityLedger.abi.Errors["EpochConflict"].ID.Bytes()[:4]) {
		return liabilityLedger.UnpackEpochConflictError(raw[4:])
	}
	if bytes.Equal(raw[:4], liabilityLedger.abi.Errors["EpochOutOfSequence"].ID.Bytes()[:4]) {
		return liabilityLedger.UnpackEpochOutOfSequenceError(raw[4:])
	}
	if bytes.Equal(raw[:4], liabilityLedger.abi.Errors["HeadMismatch"].ID.Bytes()[:4]) {
		return liabilityLedger.UnpackHeadMismatchError(raw[4:])
	}
	if bytes.Equal(raw[:4], liabilityLedger.abi.Errors["InvalidGenesis"].ID.Bytes()[:4]) {
		return liabilityLedger.UnpackInvalidGenesisError(raw[4:])
	}
	if bytes.Equal(raw[:4], liabilityLedger.abi.Errors["LedgerIsSealed"].ID.Bytes()[:4]) {
		return liabilityLedger.UnpackLedgerIsSealedError(raw[4:])
	}
	if bytes.Equal(raw[:4], liabilityLedger.abi.Errors["PreviousDigestMismatch"].ID.Bytes()[:4]) {
		return liabilityLedger.UnpackPreviousDigestMismatchError(raw[4:])
	}
	if bytes.Equal(raw[:4], liabilityLedger.abi.Errors["SafeCastOverflowedUintDowncast"].ID.Bytes()[:4]) {
		return liabilityLedger.UnpackSafeCastOverflowedUintDowncastError(raw[4:])
	}
	if bytes.Equal(raw[:4], liabilityLedger.abi.Errors["SnapshotInFuture"].ID.Bytes()[:4]) {
		return liabilityLedger.UnpackSnapshotInFutureError(raw[4:])
	}
	if bytes.Equal(raw[:4], liabilityLedger.abi.Errors["SnapshotNotAfterPrevious"].ID.Bytes()[:4]) {
		return liabilityLedger.UnpackSnapshotNotAfterPreviousError(raw[4:])
	}
	if bytes.Equal(raw[:4], liabilityLedger.abi.Errors["TooManyLeaves"].ID.Bytes()[:4]) {
		return liabilityLedger.UnpackTooManyLeavesError(raw[4:])
	}
	if bytes.Equal(raw[:4], liabilityLedger.abi.Errors["UnknownEpoch"].ID.Bytes()[:4]) {
		return liabilityLedger.UnpackUnknownEpochError(raw[4:])
	}
	if bytes.Equal(raw[:4], liabilityLedger.abi.Errors["ZeroAnchorer"].ID.Bytes()[:4]) {
		return liabilityLedger.UnpackZeroAnchorerError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// LiabilityLedgerAccessControlBadConfirmation represents a AccessControlBadConfirmation error raised by the LiabilityLedger contract.
type LiabilityLedgerAccessControlBadConfirmation struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AccessControlBadConfirmation()
func LiabilityLedgerAccessControlBadConfirmationErrorID() common.Hash {
	return common.HexToHash("0x6697b23232a647058342c0724fe7c415cab25915b54e5dbc03f233173d37b41c")
}

// UnpackAccessControlBadConfirmationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AccessControlBadConfirmation()
func (liabilityLedger *LiabilityLedger) UnpackAccessControlBadConfirmationError(raw []byte) (*LiabilityLedgerAccessControlBadConfirmation, error) {
	out := new(LiabilityLedgerAccessControlBadConfirmation)
	if err := liabilityLedger.abi.UnpackIntoInterface(out, "AccessControlBadConfirmation", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// LiabilityLedgerAccessControlEnforcedDefaultAdminDelay represents a AccessControlEnforcedDefaultAdminDelay error raised by the LiabilityLedger contract.
type LiabilityLedgerAccessControlEnforcedDefaultAdminDelay struct {
	Schedule *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AccessControlEnforcedDefaultAdminDelay(uint48 schedule)
func LiabilityLedgerAccessControlEnforcedDefaultAdminDelayErrorID() common.Hash {
	return common.HexToHash("0x19ca5ebb8fb33f00e502c9392eddab1501674629178bf69b853cf037aaf4bb5d")
}

// UnpackAccessControlEnforcedDefaultAdminDelayError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AccessControlEnforcedDefaultAdminDelay(uint48 schedule)
func (liabilityLedger *LiabilityLedger) UnpackAccessControlEnforcedDefaultAdminDelayError(raw []byte) (*LiabilityLedgerAccessControlEnforcedDefaultAdminDelay, error) {
	out := new(LiabilityLedgerAccessControlEnforcedDefaultAdminDelay)
	if err := liabilityLedger.abi.UnpackIntoInterface(out, "AccessControlEnforcedDefaultAdminDelay", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// LiabilityLedgerAccessControlEnforcedDefaultAdminRules represents a AccessControlEnforcedDefaultAdminRules error raised by the LiabilityLedger contract.
type LiabilityLedgerAccessControlEnforcedDefaultAdminRules struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AccessControlEnforcedDefaultAdminRules()
func LiabilityLedgerAccessControlEnforcedDefaultAdminRulesErrorID() common.Hash {
	return common.HexToHash("0x3fc3c27ae3db78c81b8f6e685172134623efa268ee8cd8d54be38ad2a74fc13b")
}

// UnpackAccessControlEnforcedDefaultAdminRulesError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AccessControlEnforcedDefaultAdminRules()
func (liabilityLedger *LiabilityLedger) UnpackAccessControlEnforcedDefaultAdminRulesError(raw []byte) (*LiabilityLedgerAccessControlEnforcedDefaultAdminRules, error) {
	out := new(LiabilityLedgerAccessControlEnforcedDefaultAdminRules)
	if err := liabilityLedger.abi.UnpackIntoInterface(out, "AccessControlEnforcedDefaultAdminRules", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// LiabilityLedgerAccessControlInvalidDefaultAdmin represents a AccessControlInvalidDefaultAdmin error raised by the LiabilityLedger contract.
type LiabilityLedgerAccessControlInvalidDefaultAdmin struct {
	DefaultAdmin common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AccessControlInvalidDefaultAdmin(address defaultAdmin)
func LiabilityLedgerAccessControlInvalidDefaultAdminErrorID() common.Hash {
	return common.HexToHash("0xc22c8022f2a840d6b6a9f113407715f5bbd4e88c1b0dd9434dc00700ba609ed4")
}

// UnpackAccessControlInvalidDefaultAdminError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AccessControlInvalidDefaultAdmin(address defaultAdmin)
func (liabilityLedger *LiabilityLedger) UnpackAccessControlInvalidDefaultAdminError(raw []byte) (*LiabilityLedgerAccessControlInvalidDefaultAdmin, error) {
	out := new(LiabilityLedgerAccessControlInvalidDefaultAdmin)
	if err := liabilityLedger.abi.UnpackIntoInterface(out, "AccessControlInvalidDefaultAdmin", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// LiabilityLedgerAccessControlUnauthorizedAccount represents a AccessControlUnauthorizedAccount error raised by the LiabilityLedger contract.
type LiabilityLedgerAccessControlUnauthorizedAccount struct {
	Account    common.Address
	NeededRole [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AccessControlUnauthorizedAccount(address account, bytes32 neededRole)
func LiabilityLedgerAccessControlUnauthorizedAccountErrorID() common.Hash {
	return common.HexToHash("0xe2517d3fbfae6f8515ef5ff1ccedc3933ab0cbbda0b492c06eb54ad10ef03b3e")
}

// UnpackAccessControlUnauthorizedAccountError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AccessControlUnauthorizedAccount(address account, bytes32 neededRole)
func (liabilityLedger *LiabilityLedger) UnpackAccessControlUnauthorizedAccountError(raw []byte) (*LiabilityLedgerAccessControlUnauthorizedAccount, error) {
	out := new(LiabilityLedgerAccessControlUnauthorizedAccount)
	if err := liabilityLedger.abi.UnpackIntoInterface(out, "AccessControlUnauthorizedAccount", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// LiabilityLedgerEmptyEpochMismatch represents a EmptyEpochMismatch error raised by the LiabilityLedger contract.
type LiabilityLedgerEmptyEpochMismatch struct {
	Root  [32]byte
	Total uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error EmptyEpochMismatch(bytes32 root, uint64 total)
func LiabilityLedgerEmptyEpochMismatchErrorID() common.Hash {
	return common.HexToHash("0xde718304bb7fe1b2763221798581b76b906780448b2dcd8bad16597d786f1080")
}

// UnpackEmptyEpochMismatchError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error EmptyEpochMismatch(bytes32 root, uint64 total)
func (liabilityLedger *LiabilityLedger) UnpackEmptyEpochMismatchError(raw []byte) (*LiabilityLedgerEmptyEpochMismatch, error) {
	out := new(LiabilityLedgerEmptyEpochMismatch)
	if err := liabilityLedger.abi.UnpackIntoInterface(out, "EmptyEpochMismatch", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// LiabilityLedgerEpochAlreadyAnchored represents a EpochAlreadyAnchored error raised by the LiabilityLedger contract.
type LiabilityLedgerEpochAlreadyAnchored struct {
	Number uint64
	Digest [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error EpochAlreadyAnchored(uint64 number, bytes32 digest)
func LiabilityLedgerEpochAlreadyAnchoredErrorID() common.Hash {
	return common.HexToHash("0xab455a25015fe0023a0a76c842f75125738bc3b7b3167d58a41d09451f0baf74")
}

// UnpackEpochAlreadyAnchoredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error EpochAlreadyAnchored(uint64 number, bytes32 digest)
func (liabilityLedger *LiabilityLedger) UnpackEpochAlreadyAnchoredError(raw []byte) (*LiabilityLedgerEpochAlreadyAnchored, error) {
	out := new(LiabilityLedgerEpochAlreadyAnchored)
	if err := liabilityLedger.abi.UnpackIntoInterface(out, "EpochAlreadyAnchored", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// LiabilityLedgerEpochBeforeGenesis represents a EpochBeforeGenesis error raised by the LiabilityLedger contract.
type LiabilityLedgerEpochBeforeGenesis struct {
	GenesisEpoch uint64
	Number       uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error EpochBeforeGenesis(uint64 genesisEpoch, uint64 number)
func LiabilityLedgerEpochBeforeGenesisErrorID() common.Hash {
	return common.HexToHash("0xacc3e9cb2b490f2ca39781764307e7ead97b1740e688ec48ce82eb74ab3cad45")
}

// UnpackEpochBeforeGenesisError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error EpochBeforeGenesis(uint64 genesisEpoch, uint64 number)
func (liabilityLedger *LiabilityLedger) UnpackEpochBeforeGenesisError(raw []byte) (*LiabilityLedgerEpochBeforeGenesis, error) {
	out := new(LiabilityLedgerEpochBeforeGenesis)
	if err := liabilityLedger.abi.UnpackIntoInterface(out, "EpochBeforeGenesis", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// LiabilityLedgerEpochConflict represents a EpochConflict error raised by the LiabilityLedger contract.
type LiabilityLedgerEpochConflict struct {
	Number   uint64
	Anchored [32]byte
	Given    [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error EpochConflict(uint64 number, bytes32 anchored, bytes32 given)
func LiabilityLedgerEpochConflictErrorID() common.Hash {
	return common.HexToHash("0x6d613c397e2d7ea70b114e132c49a28c0fa5de67dbb8170276880d74174393a9")
}

// UnpackEpochConflictError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error EpochConflict(uint64 number, bytes32 anchored, bytes32 given)
func (liabilityLedger *LiabilityLedger) UnpackEpochConflictError(raw []byte) (*LiabilityLedgerEpochConflict, error) {
	out := new(LiabilityLedgerEpochConflict)
	if err := liabilityLedger.abi.UnpackIntoInterface(out, "EpochConflict", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// LiabilityLedgerEpochOutOfSequence represents a EpochOutOfSequence error raised by the LiabilityLedger contract.
type LiabilityLedgerEpochOutOfSequence struct {
	Expected uint64
	Given    uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error EpochOutOfSequence(uint64 expected, uint64 given)
func LiabilityLedgerEpochOutOfSequenceErrorID() common.Hash {
	return common.HexToHash("0x7a8d061d59616552e34e07c4e10cd96598f8a51489f859f21055dba983f36722")
}

// UnpackEpochOutOfSequenceError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error EpochOutOfSequence(uint64 expected, uint64 given)
func (liabilityLedger *LiabilityLedger) UnpackEpochOutOfSequenceError(raw []byte) (*LiabilityLedgerEpochOutOfSequence, error) {
	out := new(LiabilityLedgerEpochOutOfSequence)
	if err := liabilityLedger.abi.UnpackIntoInterface(out, "EpochOutOfSequence", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// LiabilityLedgerHeadMismatch represents a HeadMismatch error raised by the LiabilityLedger contract.
type LiabilityLedgerHeadMismatch struct {
	Head     [32]byte
	Expected [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error HeadMismatch(bytes32 head, bytes32 expected)
func LiabilityLedgerHeadMismatchErrorID() common.Hash {
	return common.HexToHash("0x6f2de3c75673389a279dc32b221eac892abdabf6b2d716499a572ff532f5ec50")
}

// UnpackHeadMismatchError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error HeadMismatch(bytes32 head, bytes32 expected)
func (liabilityLedger *LiabilityLedger) UnpackHeadMismatchError(raw []byte) (*LiabilityLedgerHeadMismatch, error) {
	out := new(LiabilityLedgerHeadMismatch)
	if err := liabilityLedger.abi.UnpackIntoInterface(out, "HeadMismatch", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// LiabilityLedgerInvalidGenesis represents a InvalidGenesis error raised by the LiabilityLedger contract.
type LiabilityLedgerInvalidGenesis struct {
	GenesisEpoch  uint64
	GenesisDigest [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidGenesis(uint64 genesisEpoch, bytes32 genesisDigest)
func LiabilityLedgerInvalidGenesisErrorID() common.Hash {
	return common.HexToHash("0x111f019fa2344bb16ed38186e3c5142be0c96d972db2a225af2e8f4835689530")
}

// UnpackInvalidGenesisError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidGenesis(uint64 genesisEpoch, bytes32 genesisDigest)
func (liabilityLedger *LiabilityLedger) UnpackInvalidGenesisError(raw []byte) (*LiabilityLedgerInvalidGenesis, error) {
	out := new(LiabilityLedgerInvalidGenesis)
	if err := liabilityLedger.abi.UnpackIntoInterface(out, "InvalidGenesis", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// LiabilityLedgerLedgerIsSealed represents a LedgerIsSealed error raised by the LiabilityLedger contract.
type LiabilityLedgerLedgerIsSealed struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error LedgerIsSealed()
func LiabilityLedgerLedgerIsSealedErrorID() common.Hash {
	return common.HexToHash("0x712a4e2f2cfa8aa62d1c3b9f2e3bce545dad83b3813e02d35f0b2a51811808e4")
}

// UnpackLedgerIsSealedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error LedgerIsSealed()
func (liabilityLedger *LiabilityLedger) UnpackLedgerIsSealedError(raw []byte) (*LiabilityLedgerLedgerIsSealed, error) {
	out := new(LiabilityLedgerLedgerIsSealed)
	if err := liabilityLedger.abi.UnpackIntoInterface(out, "LedgerIsSealed", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// LiabilityLedgerPreviousDigestMismatch represents a PreviousDigestMismatch error raised by the LiabilityLedger contract.
type LiabilityLedgerPreviousDigestMismatch struct {
	Expected [32]byte
	Given    [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error PreviousDigestMismatch(bytes32 expected, bytes32 given)
func LiabilityLedgerPreviousDigestMismatchErrorID() common.Hash {
	return common.HexToHash("0x0796b55ac9b13adadc488175cf6ef6fbf1565c2537f4e021d7268b5ab5a335d9")
}

// UnpackPreviousDigestMismatchError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error PreviousDigestMismatch(bytes32 expected, bytes32 given)
func (liabilityLedger *LiabilityLedger) UnpackPreviousDigestMismatchError(raw []byte) (*LiabilityLedgerPreviousDigestMismatch, error) {
	out := new(LiabilityLedgerPreviousDigestMismatch)
	if err := liabilityLedger.abi.UnpackIntoInterface(out, "PreviousDigestMismatch", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// LiabilityLedgerSafeCastOverflowedUintDowncast represents a SafeCastOverflowedUintDowncast error raised by the LiabilityLedger contract.
type LiabilityLedgerSafeCastOverflowedUintDowncast struct {
	Bits  uint8
	Value *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SafeCastOverflowedUintDowncast(uint8 bits, uint256 value)
func LiabilityLedgerSafeCastOverflowedUintDowncastErrorID() common.Hash {
	return common.HexToHash("0x6dfcc6503a32754ce7a89698e18201fc5294fd4aad43edefee786f88423b1a12")
}

// UnpackSafeCastOverflowedUintDowncastError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SafeCastOverflowedUintDowncast(uint8 bits, uint256 value)
func (liabilityLedger *LiabilityLedger) UnpackSafeCastOverflowedUintDowncastError(raw []byte) (*LiabilityLedgerSafeCastOverflowedUintDowncast, error) {
	out := new(LiabilityLedgerSafeCastOverflowedUintDowncast)
	if err := liabilityLedger.abi.UnpackIntoInterface(out, "SafeCastOverflowedUintDowncast", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// LiabilityLedgerSnapshotInFuture represents a SnapshotInFuture error raised by the LiabilityLedger contract.
type LiabilityLedgerSnapshotInFuture struct {
	AsOf      uint64
	BlockTime uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SnapshotInFuture(uint64 asOf, uint64 blockTime)
func LiabilityLedgerSnapshotInFutureErrorID() common.Hash {
	return common.HexToHash("0x18aaced2679cfd7d188264841ef0e83939ee734887344b8d78b84cd47abff282")
}

// UnpackSnapshotInFutureError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SnapshotInFuture(uint64 asOf, uint64 blockTime)
func (liabilityLedger *LiabilityLedger) UnpackSnapshotInFutureError(raw []byte) (*LiabilityLedgerSnapshotInFuture, error) {
	out := new(LiabilityLedgerSnapshotInFuture)
	if err := liabilityLedger.abi.UnpackIntoInterface(out, "SnapshotInFuture", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// LiabilityLedgerSnapshotNotAfterPrevious represents a SnapshotNotAfterPrevious error raised by the LiabilityLedger contract.
type LiabilityLedgerSnapshotNotAfterPrevious struct {
	AsOf     uint64
	Previous uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SnapshotNotAfterPrevious(uint64 asOf, uint64 previous)
func LiabilityLedgerSnapshotNotAfterPreviousErrorID() common.Hash {
	return common.HexToHash("0x81ed60e995f8fcd4753ca10053fc7f547e1dd0c250e9bc06e427fee7bea5773f")
}

// UnpackSnapshotNotAfterPreviousError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SnapshotNotAfterPrevious(uint64 asOf, uint64 previous)
func (liabilityLedger *LiabilityLedger) UnpackSnapshotNotAfterPreviousError(raw []byte) (*LiabilityLedgerSnapshotNotAfterPrevious, error) {
	out := new(LiabilityLedgerSnapshotNotAfterPrevious)
	if err := liabilityLedger.abi.UnpackIntoInterface(out, "SnapshotNotAfterPrevious", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// LiabilityLedgerTooManyLeaves represents a TooManyLeaves error raised by the LiabilityLedger contract.
type LiabilityLedgerTooManyLeaves struct {
	LeafCount uint64
	Maximum   uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TooManyLeaves(uint64 leafCount, uint64 maximum)
func LiabilityLedgerTooManyLeavesErrorID() common.Hash {
	return common.HexToHash("0x09e51df967feec9808a68989cf05b6cd9c8c3c92676ae3ea00eda6a12380a9af")
}

// UnpackTooManyLeavesError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TooManyLeaves(uint64 leafCount, uint64 maximum)
func (liabilityLedger *LiabilityLedger) UnpackTooManyLeavesError(raw []byte) (*LiabilityLedgerTooManyLeaves, error) {
	out := new(LiabilityLedgerTooManyLeaves)
	if err := liabilityLedger.abi.UnpackIntoInterface(out, "TooManyLeaves", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// LiabilityLedgerUnknownEpoch represents a UnknownEpoch error raised by the LiabilityLedger contract.
type LiabilityLedgerUnknownEpoch struct {
	Number uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error UnknownEpoch(uint64 number)
func LiabilityLedgerUnknownEpochErrorID() common.Hash {
	return common.HexToHash("0xa4e21e6c8f881eff51b5764f98a67304231c1cb1ba8af636f941098552d2bd6c")
}

// UnpackUnknownEpochError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error UnknownEpoch(uint64 number)
func (liabilityLedger *LiabilityLedger) UnpackUnknownEpochError(raw []byte) (*LiabilityLedgerUnknownEpoch, error) {
	out := new(LiabilityLedgerUnknownEpoch)
	if err := liabilityLedger.abi.UnpackIntoInterface(out, "UnknownEpoch", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// LiabilityLedgerZeroAnchorer represents a ZeroAnchorer error raised by the LiabilityLedger contract.
type LiabilityLedgerZeroAnchorer struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ZeroAnchorer()
func LiabilityLedgerZeroAnchorerErrorID() common.Hash {
	return common.HexToHash("0x7110c4a1e207dd190eaa2e6c40aeb0e70affea5eb498857bf51f0d9d85807115")
}

// UnpackZeroAnchorerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ZeroAnchorer()
func (liabilityLedger *LiabilityLedger) UnpackZeroAnchorerError(raw []byte) (*LiabilityLedgerZeroAnchorer, error) {
	out := new(LiabilityLedgerZeroAnchorer)
	if err := liabilityLedger.abi.UnpackIntoInterface(out, "ZeroAnchorer", raw); err != nil {
		return nil, err
	}
	return out, nil
}
