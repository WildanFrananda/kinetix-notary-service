// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package evidence

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

// StreamGenesis is an auto generated low-level Go binding around an user-defined struct.
type StreamGenesis struct {
	Stream             uint8
	Batch              uint64
	Digest             [32]byte
	ClosedAt           uint64
	ThroughTransaction uint64
	ThroughSequence    uint64
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

// EvidenceRegistryMetaData contains all meta data concerning the EvidenceRegistry contract.
var EvidenceRegistryMetaData = bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"admin\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"adminTransferDelay\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"anchorer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"genesis\",\"type\":\"tuple[]\",\"internalType\":\"structStreamGenesis[]\",\"components\":[{\"name\":\"stream\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"batch\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"digest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"closedAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"throughTransaction\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"throughSequence\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"ANCHORER_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"DEFAULT_ADMIN_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MAX_STREAMS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptDefaultAdminTransfer\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"anchor\",\"inputs\":[{\"name\":\"header\",\"type\":\"tuple\",\"internalType\":\"structBatchHeader\",\"components\":[{\"name\":\"stream\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"number\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"previousDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"leafCount\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"throughTransaction\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"throughSequence\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"closedAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"batchAt\",\"inputs\":[{\"name\":\"stream\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"number\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structBatch\",\"components\":[{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"digest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"leafCount\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"throughTransaction\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"throughSequence\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"closedAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"anchoredAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"beginDefaultAdminTransfer\",\"inputs\":[{\"name\":\"newAdmin\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"cancelDefaultAdminTransfer\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"changeDefaultAdminDelay\",\"inputs\":[{\"name\":\"newDelay\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"defaultAdmin\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"defaultAdminDelay\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"defaultAdminDelayIncreaseWait\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRoleAdmin\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"grantRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"hasRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isSealed\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingDefaultAdmin\",\"inputs\":[],\"outputs\":[{\"name\":\"newAdmin\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"schedule\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingDefaultAdminDelay\",\"inputs\":[],\"outputs\":[{\"name\":\"newDelay\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"schedule\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revokeRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"rollbackDefaultAdminDelay\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"scheme\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"seal\",\"inputs\":[{\"name\":\"expectedHeads\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"stateOf\",\"inputs\":[{\"name\":\"stream\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structStreamState\",\"components\":[{\"name\":\"known\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"head\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"genesisBatch\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"nextBatch\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"lastClosedAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"throughTransaction\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"throughSequence\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"streams\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8[]\",\"internalType\":\"uint8[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyInclusion\",\"inputs\":[{\"name\":\"stream\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"number\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"recordDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"index\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"siblings\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"BatchAnchored\",\"inputs\":[{\"name\":\"stream\",\"type\":\"uint8\",\"indexed\":true,\"internalType\":\"uint8\"},{\"name\":\"number\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"digest\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"header\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structBatchHeader\",\"components\":[{\"name\":\"stream\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"number\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"previousDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"leafCount\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"throughTransaction\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"throughSequence\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"closedAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DefaultAdminDelayChangeCanceled\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DefaultAdminDelayChangeScheduled\",\"inputs\":[{\"name\":\"newDelay\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"},{\"name\":\"effectSchedule\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DefaultAdminTransferCanceled\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DefaultAdminTransferScheduled\",\"inputs\":[{\"name\":\"newAdmin\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"acceptSchedule\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleAdminChanged\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"previousAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"newAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleGranted\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleRevoked\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StreamSealed\",\"inputs\":[{\"name\":\"stream\",\"type\":\"uint8\",\"indexed\":true,\"internalType\":\"uint8\"},{\"name\":\"head\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"nextBatch\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"lastClosedAt\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"throughTransaction\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"throughSequence\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AccessControlBadConfirmation\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AccessControlEnforcedDefaultAdminDelay\",\"inputs\":[{\"name\":\"schedule\",\"type\":\"uint48\",\"internalType\":\"uint48\"}]},{\"type\":\"error\",\"name\":\"AccessControlEnforcedDefaultAdminRules\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AccessControlInvalidDefaultAdmin\",\"inputs\":[{\"name\":\"defaultAdmin\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"AccessControlUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"neededRole\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"AnchorerNotGranted\",\"inputs\":[{\"name\":\"anchorer\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"BatchAlreadyAnchored\",\"inputs\":[{\"name\":\"stream\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"number\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"digest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"BatchBeforeGenesis\",\"inputs\":[{\"name\":\"stream\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"genesisBatch\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"number\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"BatchConflict\",\"inputs\":[{\"name\":\"stream\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"number\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"anchored\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"offered\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"BatchOutOfSequence\",\"inputs\":[{\"name\":\"stream\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"expected\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"offered\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"CloseInFuture\",\"inputs\":[{\"name\":\"stream\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"closedAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"blockTime\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"CloseNotAfterPrevious\",\"inputs\":[{\"name\":\"stream\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"closedAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"lastClosedAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"EmptyBatch\",\"inputs\":[{\"name\":\"stream\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"HeadCountMismatch\",\"inputs\":[{\"name\":\"streams\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"offered\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"HeadMismatch\",\"inputs\":[{\"name\":\"stream\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"head\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"InvalidGenesis\",\"inputs\":[{\"name\":\"stream\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"batch\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"digest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"NoStreams\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PositionNotAfterPrevious\",\"inputs\":[{\"name\":\"stream\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"throughTransaction\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"throughSequence\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"previousTransaction\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"previousSequence\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"PreviousDigestMismatch\",\"inputs\":[{\"name\":\"stream\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"head\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"offered\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"RegistryIsSealed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SafeCastOverflowedUintDowncast\",\"inputs\":[{\"name\":\"bits\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"StreamsNotAscending\",\"inputs\":[{\"name\":\"previous\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"stream\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"TooManyLeaves\",\"inputs\":[{\"name\":\"leafCount\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"maximum\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"TooManyStreams\",\"inputs\":[{\"name\":\"streams\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maximum\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"UnknownBatch\",\"inputs\":[{\"name\":\"stream\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"number\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"UnknownStream\",\"inputs\":[{\"name\":\"stream\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"ZeroAnchorer\",\"inputs\":[]}]",
	ID:  "EvidenceRegistry",
	Bin: "0x60c060405234801561000f575f5ffd5b506040516137ed3803806137ed83398101604081905261002e916105dc565b82846001600160a01b03811661005e57604051636116401160e11b81525f60048201526024015b60405180910390fd5b600180546001600160d01b0316600160d01b65ffffffffffff8516021790556100875f826103b2565b5050506001600160a01b0382166100b157604051637110c4a160e01b815260040160405180910390fd5b80515f036100d25760405163667ada5b60e01b815260040160405180910390fd5b602081511115610102578051604051633f0e9b8360e11b8152600481019190915260206024820152604401610055565b5f61010c82610421565b905081518110156101e5575f82828151811061012a5761012a61073c565b602002602001015190505f825f146101685783610148600185610764565b815181106101585761015861073c565b60200260200101515f015161016a565b5f5b90508060ff16825f015160ff16116101a55781516040516343f455cb60e01b815260ff80841660048301529091166024820152604401610055565b81516020830151604080850151905163ad14a26160e01b815260ff90931660048401526001600160401b0390911660248301526044820152606401610055565b5f805b835181101561034a575f8482815181106102045761020461073c565b6020026020010151905081601f61021b9190610764565b610226906008610777565b81516040805160e081018252600180825282860151602080840191825280880180516001600160401b039081168688019081529151811660608088019182528b0151821660808089019182528c0151831660a0808a019182528d0151841660c08a019081529c5160ff9081165f9081526003978890529a909a209851895490151560ff1990911617895595518888015592516002880180549251945196518416600160c01b026001600160c01b03978516600160801b02979097166001600160801b0395851668010000000000000000026001600160801b0319909416928516929092179290921793909316929092179390931790559651929096018054929096166001600160401b0319929092169190911790945516901b9290921791016101e8565b506080819052825160a05261037f7ffd78880d40ebd8f4e5ec2f3d5d5ba6a617d36ab5a96f99f2cae1bd8c669dc90b856103b2565b6103a757604051633012fc8b60e01b81526001600160a01b0385166004820152602401610055565b50505050505061078e565b5f8261040e575f6103cb6002546001600160a01b031690565b6001600160a01b0316146103f257604051631fe1e13d60e11b815260040160405180910390fd5b600280546001600160a01b0319166001600160a01b0384161790555b6104188383610498565b90505b92915050565b5f80805b835181101561048f575f8482815181106104415761044161073c565b602002602001015190508260ff16815f015160ff161115806104775750604081015160208201516001600160401b031615901514155b1561048457509392505050565b519150600101610425565b50509051919050565b5f828152602081815260408083206001600160a01b038516845290915281205460ff16610538575f838152602081815260408083206001600160a01b03861684529091529020805460ff191660011790556104f03390565b6001600160a01b0316826001600160a01b0316847f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a450600161041b565b505f61041b565b80516001600160a01b0381168114610555575f5ffd5b919050565b634e487b7160e01b5f52604160045260245ffd5b60405160c081016001600160401b03811182821017156105905761059061055a565b60405290565b604051601f8201601f191681016001600160401b03811182821017156105be576105be61055a565b604052919050565b80516001600160401b0381168114610555575f5ffd5b5f5f5f5f608085870312156105ef575f5ffd5b6105f88561053f565b9350602085015165ffffffffffff81168114610612575f5ffd5b92506106206040860161053f565b60608601519092506001600160401b0381111561063b575f5ffd5b8501601f8101871361064b575f5ffd5b80516001600160401b038111156106645761066461055a565b61067360208260051b01610596565b80828252602082019150602060c08402850101925089831115610694575f5ffd5b6020840193505b8284101561072e5760c0848b0312156106b2575f5ffd5b6106ba61056e565b845160ff811681146106ca575f5ffd5b81526106d8602086016105c6565b6020820152604085015160408201526106f3606086016105c6565b6060820152610704608086016105c6565b608082015261071560a086016105c6565b60a0820152825260c0939093019260209091019061069b565b969995985093965050505050565b634e487b7160e01b5f52603260045260245ffd5b634e487b7160e01b5f52601160045260245ffd5b8181038181111561041b5761041b610750565b808202811582820484141761041b5761041b610750565b60805160a0516130376107b65f395f8181611713015261177a01525f6117a201526130375ff3fe608060405234801561000f575f5ffd5b50600436106101bb575f3560e01c80638da5cb5b116100f3578063cc8463c811610093578063d547741f1161006e578063d547741f14610441578063d602b9fd14610454578063e283c0441461045c578063f53f43b01461046f575f5ffd5b8063cc8463c8146103e5578063cefc1429146103ed578063cf6eefb7146103f5575f5ffd5b8063a1eda53c116100ce578063a1eda53c1461038f578063a217fddf146103b6578063a615d225146103bd578063c7991ec7146103c5575f5ffd5b80638da5cb5b1461033157806391d148541461033957806397413ab01461037c575f5ffd5b8063631f98521161015e5780637cfd9a51116101395780637cfd9a51146102a957806384ef8ffc146102d05780638518b0ad1461030f5780638960b0851461031e575f5ffd5b8063631f985214610278578063634e93da14610283578063649a5ec714610296575f5ffd5b8063248a9ca311610199578063248a9ca31461020d5780632f2ff15d1461023d578063324030401461025057806336568abe14610265575f5ffd5b806301ffc9a7146101bf578063022d63fb146101e75780630aa6220b14610203575b5f5ffd5b6101d26101cd366004612896565b610509565b60405190151581526020015b60405180910390f35b620697805b60405165ffffffffffff90911681526020016101de565b61020b610564565b005b61022f61021b3660046128d5565b5f9081526020819052604090206001015490565b6040519081526020016101de565b61020b61024b366004612914565b610579565b6102586105be565b6040516101de919061293e565b61020b610273366004612914565b6105cd565b60055460ff166101d2565b61020b610291366004612983565b6106d7565b61020b6102a436600461299c565b6106ea565b61022f7ffd78880d40ebd8f4e5ec2f3d5d5ba6a617d36ab5a96f99f2cae1bd8c669dc90b81565b60025473ffffffffffffffffffffffffffffffffffffffff165b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016101de565b604051600181526020016101de565b61020b61032c366004612a09565b6106fd565b6102ea610a08565b6101d2610347366004612914565b5f9182526020828152604080842073ffffffffffffffffffffffffffffffffffffffff93909316845291905290205460ff1690565b61022f61038a366004612a48565b610a28565b61039761120a565b6040805165ffffffffffff9384168152929091166020830152016101de565b61022f5f81565b61022f602081565b6103d86103d3366004612a72565b611284565b6040516101de9190612a8b565b6101ec61135d565b61020b6113fa565b6001546040805173ffffffffffffffffffffffffffffffffffffffff831681527401000000000000000000000000000000000000000090920465ffffffffffff166020830152016101de565b61020b61044f366004612914565b611456565b61020b611497565b6101d261046a366004612b2b565b6114a9565b61048261047d366004612bae565b6115a3565b6040516101de91905f60e082019050825182526020830151602083015267ffffffffffffffff604084015116604083015267ffffffffffffffff606084015116606083015267ffffffffffffffff608084015116608083015267ffffffffffffffff60a08401511660a083015267ffffffffffffffff60c08401511660c083015292915050565b5f7fffffffff0000000000000000000000000000000000000000000000000000000082167f1a430fb900000000000000000000000000000000000000000000000000000000148061055e575061055e82611679565b92915050565b5f61056e816116ce565b6105766116d8565b50565b816105b0576040517f3fc3c27a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6105ba82826116e4565b5050565b60606105c861170e565b905090565b811580156105f5575060025473ffffffffffffffffffffffffffffffffffffffff8281169116145b156106cd5760015473ffffffffffffffffffffffffffffffffffffffff81169074010000000000000000000000000000000000000000900465ffffffffffff1681151580610649575065ffffffffffff8116155b8061065c57504265ffffffffffff821610155b156106a2576040517f19ca5ebb00000000000000000000000000000000000000000000000000000000815265ffffffffffff821660048201526024015b60405180910390fd5b5050600180547fffffffffffff000000000000ffffffffffffffffffffffffffffffffffffffff1690555b6105ba828261180b565b5f6106e1816116ce565b6105ba82611869565b5f6106f4816116ce565b6105ba826118e8565b5f610707816116ce565b60055460ff1615610744576040517ff3e6ff3100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f61074d61170e565b805190915083146107975780516040517fbe8ea0e7000000000000000000000000000000000000000000000000000000008152600481019190915260248101849052604401610699565b5f6107a3828686611957565b9050815181101561085b578181815181106107c0576107c0612bd6565b602002602001015160035f8484815181106107dd576107dd612bd6565b602002602001015160ff1660ff1681526020019081526020015f206001015486868481811061080e5761080e612bd6565b6040517ff18e8ed000000000000000000000000000000000000000000000000000000000815260ff9095166004860152602485019390935250602090910201356044820152606401610699565b600580547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660011790555f5b8251811015610a00575f60035f8584815181106108a7576108a7612bd6565b60209081029190910181015160ff90811683528282019390935260409182015f20825160e081018452815490941615158452600181015491840191909152600281015467ffffffffffffffff808216938501939093526801000000000000000081048316606085015270010000000000000000000000000000000081048316608085015278010000000000000000000000000000000000000000000000009004821660a0840152600301541660c0820152845190915084908390811061096f5761096f612bd6565b602002602001015160ff167fd869cd5054cea46a1e8780d76c00c41b0a7a2e638c90a47f23acf030704cc5878260200151836060015184608001518560a001518660c001516040516109ef95949392919094855267ffffffffffffffff938416602086015291831660408501528216606084015216608082015260a00190565b60405180910390a250600101610888565b505050505050565b5f6105c860025473ffffffffffffffffffffffffffffffffffffffff1690565b5f7ffd78880d40ebd8f4e5ec2f3d5d5ba6a617d36ab5a96f99f2cae1bd8c669dc90b610a53816116ce565b60055460ff1615610a90576040517ff3e6ff3100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f610aa6610aa16020860186612a72565b6119d0565b90505f610ac0610abb36879003870187612c30565b611a24565b600283015490915067ffffffffffffffff16610ae26040870160208801612d0a565b67ffffffffffffffff161015610b6d57610aff6020860186612a72565b600283015467ffffffffffffffff16610b1e6040880160208901612d0a565b6040517f0d08789500000000000000000000000000000000000000000000000000000000815260ff909316600484015267ffffffffffffffff9182166024840152166044820152606401610699565b600282015468010000000000000000900467ffffffffffffffff16610b986040870160208801612d0a565b67ffffffffffffffff161015610ceb575f600481610bb96020890189612a72565b60ff1660ff1681526020019081526020015f205f876020016020810190610be09190612d0a565b67ffffffffffffffff1667ffffffffffffffff1681526020019081526020015f20600101549050818103610c7957610c1b6020870187612a72565b610c2b6040880160208901612d0a565b6040517f6035ccca00000000000000000000000000000000000000000000000000000000815260ff909216600483015267ffffffffffffffff16602482015260448101839052606401610699565b610c866020870187612a72565b610c966040880160208901612d0a565b6040517fc071754c00000000000000000000000000000000000000000000000000000000815260ff909216600483015267ffffffffffffffff1660248201526044810182905260648101839052608401610699565b600282015468010000000000000000900467ffffffffffffffff16610d166040870160208801612d0a565b67ffffffffffffffff161115610dad57610d336020860186612a72565b600283015468010000000000000000900467ffffffffffffffff16610d5e6040880160208901612d0a565b6040517fb39c398500000000000000000000000000000000000000000000000000000000815260ff909316600484015267ffffffffffffffff9182166024840152166044820152606401610699565b8160010154856040013514610e1757610dc96020860186612a72565b6001830154604080517f77a6bbb300000000000000000000000000000000000000000000000000000000815260ff909316600484015260248301919091528601356044820152606401610699565b610e218583611a89565b6040805160e0810182526060870135815260208101839052908101610e4c60a0880160808901612d0a565b67ffffffffffffffff168152602001610e6b60c0880160a08901612d0a565b67ffffffffffffffff168152602001610e8a60e0880160c08901612d0a565b67ffffffffffffffff168152602001610eaa610100880160e08901612d0a565b67ffffffffffffffff168152602001610ec242611e53565b67ffffffffffffffff16905260045f610ede6020890189612a72565b60ff1660ff1681526020019081526020015f205f876020016020810190610f059190612d0a565b67ffffffffffffffff908116825260208083019390935260409182015f20845181558484015160018083019190915585840151600283018054606089015160808a015160a08b01519488167fffffffffffffffffffffffffffffffff00000000000000000000000000000000909316929092176801000000000000000091881691909102176fffffffffffffffffffffffffffffffff167001000000000000000000000000000000009187169190910277ffffffffffffffffffffffffffffffffffffffffffffffff161778010000000000000000000000000000000000000000000000009286169290920291909117905560c090950151600390910180547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001691909216179055918401839055611041918701908701612d0a565b61104c906001612d50565b60028301805467ffffffffffffffff9290921668010000000000000000027fffffffffffffffffffffffffffffffff0000000000000000ffffffffffffffff9092169190911790556110a5610100860160e08701612d0a565b60028301805467ffffffffffffffff92909216700100000000000000000000000000000000027fffffffffffffffff0000000000000000ffffffffffffffffffffffffffffffff90921691909117905561110560c0860160a08701612d0a565b60028301805467ffffffffffffffff9290921678010000000000000000000000000000000000000000000000000277ffffffffffffffffffffffffffffffffffffffffffffffff90921691909117905561116560e0860160c08701612d0a565b6003830180547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001667ffffffffffffffff929092169190911790556111b06040860160208701612d0a565b67ffffffffffffffff166111c76020870187612a72565b60ff167eea6aacf26d19d941ee9b298ca89412351a8adcfb827acabf1322cddf2c587b83886040516111fa929190612d70565b60405180910390a3949350505050565b6002545f907a010000000000000000000000000000000000000000000000000000900465ffffffffffff16801515801561124c57504265ffffffffffff821610155b611257575f5f61127c565b60025474010000000000000000000000000000000000000000900465ffffffffffff16815b915091509091565b6040805160e0810182525f80825260208201819052918101829052606081018290526080810182905260a0810182905260c08101919091526112c5826119d0565b6040805160e081018252825460ff161515815260018301546020820152600283015467ffffffffffffffff808216938301939093526801000000000000000081048316606083015270010000000000000000000000000000000081048316608083015278010000000000000000000000000000000000000000000000009004821660a08201526003909201541660c082015292915050565b6002545f907a010000000000000000000000000000000000000000000000000000900465ffffffffffff16801515801561139e57504265ffffffffffff8216105b6113d0576001547a010000000000000000000000000000000000000000000000000000900465ffffffffffff166113f4565b60025474010000000000000000000000000000000000000000900465ffffffffffff165b91505090565b60015473ffffffffffffffffffffffffffffffffffffffff1633811461144e576040517fc22c8022000000000000000000000000000000000000000000000000000000008152336004820152602401610699565b610576611ea6565b8161148d576040517f3fc3c27a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6105ba8282611f97565b5f6114a1816116ce565b610576611fbb565b5f5f6114b58989611fc5565b600281015490915067ffffffffffffffff908116908616106114da575f915050611598565b60028101546114f29067ffffffffffffffff1661208b565b60ff168314611504575f915050611598565b604080517f100000000000000000000000000000000000000000000000000000000000000060208083019190915260f88c901b7fff00000000000000000000000000000000000000000000000000000000000000166021830152602282018a905260428083018a905283518084039091018152606290920190925280519101208154611592828888886120e9565b14925050505b979650505050505050565b6040805160e0810182525f80825260208201819052918101829052606081018290526080810182905260a0810182905260c08101919091526115e58383611fc5565b6040805160e0810182528254815260018301546020820152600283015467ffffffffffffffff808216938301939093526801000000000000000081048316606083015270010000000000000000000000000000000081048316608083015278010000000000000000000000000000000000000000000000009004821660a08201526003909201541660c08201529392505050565b5f7fffffffff0000000000000000000000000000000000000000000000000000000082167f3149878600000000000000000000000000000000000000000000000000000000148061055e575061055e8261221a565b61057681336122b0565b6116e25f5f612335565b565b5f828152602081905260409020600101546116fe816116ce565b611708838361248e565b50505050565b60605f7f000000000000000000000000000000000000000000000000000000000000000067ffffffffffffffff81111561174a5761174a612c03565b604051908082528060200260200182016040528015611773578160200160208202803683370190505b5090505f5b7f0000000000000000000000000000000000000000000000000000000000000000811015611805577f000000000000000000000000000000000000000000000000000000000000000081602081106117d2576117d2612bd6565b1a60f81b60f81c8282815181106117eb576117eb612bd6565b60ff90921660209283029190910190910152600101611778565b50919050565b73ffffffffffffffffffffffffffffffffffffffff8116331461185a576040517f6697b23200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b611864828261254c565b505050565b5f61187261135d565b61187b426125ad565b6118859190612e3c565b905061189182826125f8565b60405165ffffffffffff8216815273ffffffffffffffffffffffffffffffffffffffff8316907f3377dc44241e779dd06afab5b788a35ca5f3b778836e2990bdb26a2a4b2e5ed69060200160405180910390a25050565b5f6118f282612693565b6118fb426125ad565b6119059190612e3c565b90506119118282612335565b6040805165ffffffffffff8085168252831660208201527ff1038c18cf84a56e432fdbfaf746924b7ea511dfe03a6506a0ceba4888788d9b910160405180910390a15050565b5f805b84518110156119c45783838281811061197557611975612bd6565b9050602002013560035f87848151811061199157611991612bd6565b602002602001015160ff1660ff1681526020019081526020015f2060010154146119bc5790506119c9565b60010161195a565b505082515b9392505050565b60ff8082165f9081526003602052604081208054919290911661055e576040517f8b7e1f9700000000000000000000000000000000000000000000000000000000815260ff84166004820152602401610699565b5f601360f81b6001835f015184602001518560400151866060015187608001518860a001518960c001518a60e00151604051602001611a6c9a99989796959493929190612e5a565b604051602081830303815290604052805190602001209050919050565b611a9960a0830160808401612d0a565b67ffffffffffffffff165f03611aee57611ab66020830183612a72565b6040517f2aa62fa900000000000000000000000000000000000000000000000000000000815260ff9091166004820152602401610699565b640100000000611b0460a0840160808501612d0a565b67ffffffffffffffff161115611b6e57611b2460a0830160808401612d0a565b6040517f09e51df900000000000000000000000000000000000000000000000000000000815267ffffffffffffffff90911660048201526401000000006024820152604401610699565b6002810154700100000000000000000000000000000000900467ffffffffffffffff16611ba2610100840160e08501612d0a565b67ffffffffffffffff1611611c3757611bbe6020830183612a72565b611bcf610100840160e08501612d0a565b60028301546040517f5999344d00000000000000000000000000000000000000000000000000000000815260ff909316600484015267ffffffffffffffff91821660248401527001000000000000000000000000000000009004166044820152606401610699565b42611c49610100840160e08501612d0a565b67ffffffffffffffff161115611ccf57611c666020830183612a72565b611c77610100840160e08501612d0a565b611c8042611e53565b6040517fe9f968e700000000000000000000000000000000000000000000000000000000815260ff909316600484015267ffffffffffffffff9182166024840152166044820152606401610699565b60028101545f907801000000000000000000000000000000000000000000000000900467ffffffffffffffff16611d0c60c0850160a08601612d0a565b600284015467ffffffffffffffff9182169290921192505f917801000000000000000000000000000000000000000000000000900416611d5260c0860160a08701612d0a565b67ffffffffffffffff1614905081158015611d9d5750808015611d9b5750600383015467ffffffffffffffff16611d8f60e0860160c08701612d0a565b67ffffffffffffffff16115b155b1561170857611daf6020850185612a72565b611dbf60c0860160a08701612d0a565b611dcf60e0870160c08801612d0a565b600286015460038701546040517f9e8b48b400000000000000000000000000000000000000000000000000000000815260ff909516600486015267ffffffffffffffff93841660248601529183166044850152780100000000000000000000000000000000000000000000000090048216606484015216608482015260a401610699565b5f67ffffffffffffffff821115611ea257604080517f6dfcc650000000000000000000000000000000000000000000000000000000008152600481019190915260248101839052604401610699565b5090565b60015473ffffffffffffffffffffffffffffffffffffffff81169074010000000000000000000000000000000000000000900465ffffffffffff16801580611ef657504265ffffffffffff821610155b15611f37576040517f19ca5ebb00000000000000000000000000000000000000000000000000000000815265ffffffffffff82166004820152602401610699565b611f5f5f611f5a60025473ffffffffffffffffffffffffffffffffffffffff1690565b61254c565b50611f6a5f8361248e565b5050600180547fffffffffffff000000000000000000000000000000000000000000000000000016905550565b5f82815260208190526040902060010154611fb1816116ce565b611708838361254c565b6116e25f5f6125f8565b5f5f611fd0846119d0565b600281015490915067ffffffffffffffff90811690841610806120115750600281015467ffffffffffffffff68010000000000000000909104811690841610155b1561205d576040517f2293f21000000000000000000000000000000000000000000000000000000000815260ff8516600482015267ffffffffffffffff84166024820152604401610699565b505060ff82165f90815260046020908152604080832067ffffffffffffffff85168452909152902092915050565b5f60018267ffffffffffffffff16116120a557505f919050565b5f6120b1600184612fdb565b90505f5b67ffffffffffffffff8216156119c95760018267ffffffffffffffff16901c915080806120e190612ffb565b9150506120b5565b5f8484825b8481101561220e57600182165f0361217e576121778387878481811061211657612116612bd6565b905060200201356040517f1100000000000000000000000000000000000000000000000000000000000000602082015260218101839052604181018290525f9060610160405160208183030381529060405280519060200120905092915050565b92506121f8565b6121f586868381811061219357612193612bd6565b90506020020135846040517f1100000000000000000000000000000000000000000000000000000000000000602082015260218101839052604181018290525f9060610160405160208183030381529060405280519060200120905092915050565b92505b677fffffffffffffff600192831c1691016120ee565b50909695505050505050565b5f7fffffffff0000000000000000000000000000000000000000000000000000000082167f7965db0b00000000000000000000000000000000000000000000000000000000148061055e57507f01ffc9a7000000000000000000000000000000000000000000000000000000007fffffffff0000000000000000000000000000000000000000000000000000000083161461055e565b5f8281526020818152604080832073ffffffffffffffffffffffffffffffffffffffff8516845290915290205460ff166105ba576040517fe2517d3f00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8216600482015260248101839052604401610699565b6002547a010000000000000000000000000000000000000000000000000000900465ffffffffffff168015612409574265ffffffffffff821610156123e0576002546001805479ffffffffffffffffffffffffffffffffffffffffffffffffffff167401000000000000000000000000000000000000000090920465ffffffffffff167a01000000000000000000000000000000000000000000000000000002919091179055612409565b6040517f2b1fa2edafe6f7b9e97c1a9e0c3660e645beb2dcaa2d45bdbf9beaf5472e1ec5905f90a15b506002805473ffffffffffffffffffffffffffffffffffffffff167401000000000000000000000000000000000000000065ffffffffffff9485160279ffffffffffffffffffffffffffffffffffffffffffffffffffff16177a0100000000000000000000000000000000000000000000000000009290931691909102919091179055565b5f82612542575f6124b460025473ffffffffffffffffffffffffffffffffffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff1614612501576040517f3fc3c27a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600280547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff84161790555b6119c983836126e4565b5f82158015612575575060025473ffffffffffffffffffffffffffffffffffffffff8381169116145b156125a357600280547fffffffffffffffffffffffff00000000000000000000000000000000000000001690555b6119c983836127dd565b5f65ffffffffffff821115611ea2576040517f6dfcc6500000000000000000000000000000000000000000000000000000000081526030600482015260248101839052604401610699565b600180547401000000000000000000000000000000000000000065ffffffffffff84811682027fffffffffffff0000000000000000000000000000000000000000000000000000841673ffffffffffffffffffffffffffffffffffffffff881617179093559004168015611864576040517f8886ebfc4259abdbc16601dd8fb5678e54878f47b3c34836cfc51154a9605109905f90a1505050565b5f5f61269d61135d565b90508065ffffffffffff168365ffffffffffff16116126c5576126c08382613019565b6119c9565b6119c965ffffffffffff8416620697805f8282188284100282186119c9565b5f8281526020818152604080832073ffffffffffffffffffffffffffffffffffffffff8516845290915281205460ff166127d6575f8381526020818152604080832073ffffffffffffffffffffffffffffffffffffffff86168452909152902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660011790556127743390565b73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16847f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a450600161055e565b505f61055e565b5f8281526020818152604080832073ffffffffffffffffffffffffffffffffffffffff8516845290915281205460ff16156127d6575f8381526020818152604080832073ffffffffffffffffffffffffffffffffffffffff8616808552925280832080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016905551339286917ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9190a450600161055e565b5f602082840312156128a6575f5ffd5b81357fffffffff00000000000000000000000000000000000000000000000000000000811681146119c9575f5ffd5b5f602082840312156128e5575f5ffd5b5035919050565b803573ffffffffffffffffffffffffffffffffffffffff8116811461290f575f5ffd5b919050565b5f5f60408385031215612925575f5ffd5b82359150612935602084016128ec565b90509250929050565b602080825282518282018190525f918401906040840190835b8181101561297857835160ff16835260209384019390920191600101612957565b509095945050505050565b5f60208284031215612993575f5ffd5b6119c9826128ec565b5f602082840312156129ac575f5ffd5b813565ffffffffffff811681146119c9575f5ffd5b5f5f83601f8401126129d1575f5ffd5b50813567ffffffffffffffff8111156129e8575f5ffd5b6020830191508360208260051b8501011115612a02575f5ffd5b9250929050565b5f5f60208385031215612a1a575f5ffd5b823567ffffffffffffffff811115612a30575f5ffd5b612a3c858286016129c1565b90969095509350505050565b5f610100828403128015612a5a575f5ffd5b509092915050565b803560ff8116811461290f575f5ffd5b5f60208284031215612a82575f5ffd5b6119c982612a62565b5f60e0820190508251151582526020830151602083015267ffffffffffffffff604084015116604083015267ffffffffffffffff606084015116606083015267ffffffffffffffff608084015116608083015267ffffffffffffffff60a08401511660a083015260c0830151612b0d60c084018267ffffffffffffffff169052565b5092915050565b803567ffffffffffffffff8116811461290f575f5ffd5b5f5f5f5f5f5f5f60c0888a031215612b41575f5ffd5b612b4a88612a62565b9650612b5860208901612b14565b95506040880135945060608801359350612b7460808901612b14565b925060a088013567ffffffffffffffff811115612b8f575f5ffd5b612b9b8a828b016129c1565b989b979a50959850939692959293505050565b5f5f60408385031215612bbf575f5ffd5b612bc883612a62565b915061293560208401612b14565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b5f610100828403128015612c42575f5ffd5b50604051610100810167ffffffffffffffff81118282101715612c8c577f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b604052612c9883612a62565b8152612ca660208401612b14565b60208201526040838101359082015260608084013590820152612ccb60808401612b14565b6080820152612cdc60a08401612b14565b60a0820152612ced60c08401612b14565b60c0820152612cfe60e08401612b14565b60e08201529392505050565b5f60208284031215612d1a575f5ffd5b6119c982612b14565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b67ffffffffffffffff818116838216019081111561055e5761055e612d23565b828152610120810160ff612d8384612a62565b16602083015267ffffffffffffffff612d9e60208501612b14565b1660408301525f60408401359050806060840152505f606084013590508060808401525067ffffffffffffffff612dd760808501612b14565b1660a0830152612de960a08401612b14565b67ffffffffffffffff811660c084015250612e0660c08401612b14565b67ffffffffffffffff811660e084015250612e2360e08401612b14565b67ffffffffffffffff8116610100840152509392505050565b65ffffffffffff818116838216019081111561055e5761055e612d23565b7fff000000000000000000000000000000000000000000000000000000000000008b1681527fffff0000000000000000000000000000000000000000000000000000000000008a60f01b1660018201527fff000000000000000000000000000000000000000000000000000000000000008960f81b1660038201527fffffffffffffffff0000000000000000000000000000000000000000000000008860c01b16600482015286600c82015285602c8201527fffffffffffffffff0000000000000000000000000000000000000000000000008560c01b16604c820152612f68605482018560c01b7fffffffffffffffff000000000000000000000000000000000000000000000000169052565b612f99605c82018460c01b7fffffffffffffffff000000000000000000000000000000000000000000000000169052565b612fca606482018360c01b7fffffffffffffffff000000000000000000000000000000000000000000000000169052565b606c019a9950505050505050505050565b67ffffffffffffffff828116828216039081111561055e5761055e612d23565b5f60ff821660ff810361301057613010612d23565b60010192915050565b65ffffffffffff828116828216039081111561055e5761055e612d2356",
}

// EvidenceRegistry is an auto generated Go binding around an Ethereum contract.
type EvidenceRegistry struct {
	abi abi.ABI
}

// GetABI returns the ABI associated with this contract binding.
func (c *EvidenceRegistry) GetABI() abi.ABI {
	return c.abi
}

// NewEvidenceRegistry creates a new instance of EvidenceRegistry.
func NewEvidenceRegistry() *EvidenceRegistry {
	parsed, err := EvidenceRegistryMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &EvidenceRegistry{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *EvidenceRegistry) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackConstructor is the Go binding used to pack the parameters required for
// contract deployment.
//
// Solidity: constructor(address admin, uint48 adminTransferDelay, address anchorer, (uint8,uint64,bytes32,uint64,uint64,uint64)[] genesis) returns()
func (evidenceRegistry *EvidenceRegistry) PackConstructor(admin common.Address, adminTransferDelay *big.Int, anchorer common.Address, genesis []StreamGenesis) []byte {
	enc, err := evidenceRegistry.abi.Pack("", admin, adminTransferDelay, anchorer, genesis)
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
func (evidenceRegistry *EvidenceRegistry) PackANCHORERROLE() []byte {
	enc, err := evidenceRegistry.abi.Pack("ANCHORER_ROLE")
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
func (evidenceRegistry *EvidenceRegistry) TryPackANCHORERROLE() ([]byte, error) {
	return evidenceRegistry.abi.Pack("ANCHORER_ROLE")
}

// UnpackANCHORERROLE is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x7cfd9a51.
//
// Solidity: function ANCHORER_ROLE() view returns(bytes32)
func (evidenceRegistry *EvidenceRegistry) UnpackANCHORERROLE(data []byte) ([32]byte, error) {
	out, err := evidenceRegistry.abi.Unpack("ANCHORER_ROLE", data)
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
func (evidenceRegistry *EvidenceRegistry) PackDEFAULTADMINROLE() []byte {
	enc, err := evidenceRegistry.abi.Pack("DEFAULT_ADMIN_ROLE")
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
func (evidenceRegistry *EvidenceRegistry) TryPackDEFAULTADMINROLE() ([]byte, error) {
	return evidenceRegistry.abi.Pack("DEFAULT_ADMIN_ROLE")
}

// UnpackDEFAULTADMINROLE is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (evidenceRegistry *EvidenceRegistry) UnpackDEFAULTADMINROLE(data []byte) ([32]byte, error) {
	out, err := evidenceRegistry.abi.Unpack("DEFAULT_ADMIN_ROLE", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackMAXSTREAMS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa615d225.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function MAX_STREAMS() view returns(uint256)
func (evidenceRegistry *EvidenceRegistry) PackMAXSTREAMS() []byte {
	enc, err := evidenceRegistry.abi.Pack("MAX_STREAMS")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMAXSTREAMS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa615d225.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function MAX_STREAMS() view returns(uint256)
func (evidenceRegistry *EvidenceRegistry) TryPackMAXSTREAMS() ([]byte, error) {
	return evidenceRegistry.abi.Pack("MAX_STREAMS")
}

// UnpackMAXSTREAMS is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa615d225.
//
// Solidity: function MAX_STREAMS() view returns(uint256)
func (evidenceRegistry *EvidenceRegistry) UnpackMAXSTREAMS(data []byte) (*big.Int, error) {
	out, err := evidenceRegistry.abi.Unpack("MAX_STREAMS", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackAcceptDefaultAdminTransfer is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xcefc1429.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function acceptDefaultAdminTransfer() returns()
func (evidenceRegistry *EvidenceRegistry) PackAcceptDefaultAdminTransfer() []byte {
	enc, err := evidenceRegistry.abi.Pack("acceptDefaultAdminTransfer")
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
func (evidenceRegistry *EvidenceRegistry) TryPackAcceptDefaultAdminTransfer() ([]byte, error) {
	return evidenceRegistry.abi.Pack("acceptDefaultAdminTransfer")
}

// PackAnchor is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x97413ab0.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function anchor((uint8,uint64,bytes32,bytes32,uint64,uint64,uint64,uint64) header) returns(bytes32)
func (evidenceRegistry *EvidenceRegistry) PackAnchor(header BatchHeader) []byte {
	enc, err := evidenceRegistry.abi.Pack("anchor", header)
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
func (evidenceRegistry *EvidenceRegistry) TryPackAnchor(header BatchHeader) ([]byte, error) {
	return evidenceRegistry.abi.Pack("anchor", header)
}

// UnpackAnchor is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x97413ab0.
//
// Solidity: function anchor((uint8,uint64,bytes32,bytes32,uint64,uint64,uint64,uint64) header) returns(bytes32)
func (evidenceRegistry *EvidenceRegistry) UnpackAnchor(data []byte) ([32]byte, error) {
	out, err := evidenceRegistry.abi.Unpack("anchor", data)
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
func (evidenceRegistry *EvidenceRegistry) PackBatchAt(stream uint8, number uint64) []byte {
	enc, err := evidenceRegistry.abi.Pack("batchAt", stream, number)
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
func (evidenceRegistry *EvidenceRegistry) TryPackBatchAt(stream uint8, number uint64) ([]byte, error) {
	return evidenceRegistry.abi.Pack("batchAt", stream, number)
}

// UnpackBatchAt is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xf53f43b0.
//
// Solidity: function batchAt(uint8 stream, uint64 number) view returns((bytes32,bytes32,uint64,uint64,uint64,uint64,uint64))
func (evidenceRegistry *EvidenceRegistry) UnpackBatchAt(data []byte) (Batch, error) {
	out, err := evidenceRegistry.abi.Unpack("batchAt", data)
	if err != nil {
		return *new(Batch), err
	}
	out0 := *abi.ConvertType(out[0], new(Batch)).(*Batch)
	return out0, nil
}

// PackBeginDefaultAdminTransfer is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x634e93da.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function beginDefaultAdminTransfer(address newAdmin) returns()
func (evidenceRegistry *EvidenceRegistry) PackBeginDefaultAdminTransfer(newAdmin common.Address) []byte {
	enc, err := evidenceRegistry.abi.Pack("beginDefaultAdminTransfer", newAdmin)
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
func (evidenceRegistry *EvidenceRegistry) TryPackBeginDefaultAdminTransfer(newAdmin common.Address) ([]byte, error) {
	return evidenceRegistry.abi.Pack("beginDefaultAdminTransfer", newAdmin)
}

// PackCancelDefaultAdminTransfer is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd602b9fd.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function cancelDefaultAdminTransfer() returns()
func (evidenceRegistry *EvidenceRegistry) PackCancelDefaultAdminTransfer() []byte {
	enc, err := evidenceRegistry.abi.Pack("cancelDefaultAdminTransfer")
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
func (evidenceRegistry *EvidenceRegistry) TryPackCancelDefaultAdminTransfer() ([]byte, error) {
	return evidenceRegistry.abi.Pack("cancelDefaultAdminTransfer")
}

// PackChangeDefaultAdminDelay is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x649a5ec7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function changeDefaultAdminDelay(uint48 newDelay) returns()
func (evidenceRegistry *EvidenceRegistry) PackChangeDefaultAdminDelay(newDelay *big.Int) []byte {
	enc, err := evidenceRegistry.abi.Pack("changeDefaultAdminDelay", newDelay)
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
func (evidenceRegistry *EvidenceRegistry) TryPackChangeDefaultAdminDelay(newDelay *big.Int) ([]byte, error) {
	return evidenceRegistry.abi.Pack("changeDefaultAdminDelay", newDelay)
}

// PackDefaultAdmin is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x84ef8ffc.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function defaultAdmin() view returns(address)
func (evidenceRegistry *EvidenceRegistry) PackDefaultAdmin() []byte {
	enc, err := evidenceRegistry.abi.Pack("defaultAdmin")
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
func (evidenceRegistry *EvidenceRegistry) TryPackDefaultAdmin() ([]byte, error) {
	return evidenceRegistry.abi.Pack("defaultAdmin")
}

// UnpackDefaultAdmin is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x84ef8ffc.
//
// Solidity: function defaultAdmin() view returns(address)
func (evidenceRegistry *EvidenceRegistry) UnpackDefaultAdmin(data []byte) (common.Address, error) {
	out, err := evidenceRegistry.abi.Unpack("defaultAdmin", data)
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
func (evidenceRegistry *EvidenceRegistry) PackDefaultAdminDelay() []byte {
	enc, err := evidenceRegistry.abi.Pack("defaultAdminDelay")
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
func (evidenceRegistry *EvidenceRegistry) TryPackDefaultAdminDelay() ([]byte, error) {
	return evidenceRegistry.abi.Pack("defaultAdminDelay")
}

// UnpackDefaultAdminDelay is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xcc8463c8.
//
// Solidity: function defaultAdminDelay() view returns(uint48)
func (evidenceRegistry *EvidenceRegistry) UnpackDefaultAdminDelay(data []byte) (*big.Int, error) {
	out, err := evidenceRegistry.abi.Unpack("defaultAdminDelay", data)
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
func (evidenceRegistry *EvidenceRegistry) PackDefaultAdminDelayIncreaseWait() []byte {
	enc, err := evidenceRegistry.abi.Pack("defaultAdminDelayIncreaseWait")
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
func (evidenceRegistry *EvidenceRegistry) TryPackDefaultAdminDelayIncreaseWait() ([]byte, error) {
	return evidenceRegistry.abi.Pack("defaultAdminDelayIncreaseWait")
}

// UnpackDefaultAdminDelayIncreaseWait is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x022d63fb.
//
// Solidity: function defaultAdminDelayIncreaseWait() view returns(uint48)
func (evidenceRegistry *EvidenceRegistry) UnpackDefaultAdminDelayIncreaseWait(data []byte) (*big.Int, error) {
	out, err := evidenceRegistry.abi.Unpack("defaultAdminDelayIncreaseWait", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackGetRoleAdmin is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x248a9ca3.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (evidenceRegistry *EvidenceRegistry) PackGetRoleAdmin(role [32]byte) []byte {
	enc, err := evidenceRegistry.abi.Pack("getRoleAdmin", role)
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
func (evidenceRegistry *EvidenceRegistry) TryPackGetRoleAdmin(role [32]byte) ([]byte, error) {
	return evidenceRegistry.abi.Pack("getRoleAdmin", role)
}

// UnpackGetRoleAdmin is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (evidenceRegistry *EvidenceRegistry) UnpackGetRoleAdmin(data []byte) ([32]byte, error) {
	out, err := evidenceRegistry.abi.Unpack("getRoleAdmin", data)
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
func (evidenceRegistry *EvidenceRegistry) PackGrantRole(role [32]byte, account common.Address) []byte {
	enc, err := evidenceRegistry.abi.Pack("grantRole", role, account)
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
func (evidenceRegistry *EvidenceRegistry) TryPackGrantRole(role [32]byte, account common.Address) ([]byte, error) {
	return evidenceRegistry.abi.Pack("grantRole", role, account)
}

// PackHasRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x91d14854.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (evidenceRegistry *EvidenceRegistry) PackHasRole(role [32]byte, account common.Address) []byte {
	enc, err := evidenceRegistry.abi.Pack("hasRole", role, account)
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
func (evidenceRegistry *EvidenceRegistry) TryPackHasRole(role [32]byte, account common.Address) ([]byte, error) {
	return evidenceRegistry.abi.Pack("hasRole", role, account)
}

// UnpackHasRole is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (evidenceRegistry *EvidenceRegistry) UnpackHasRole(data []byte) (bool, error) {
	out, err := evidenceRegistry.abi.Unpack("hasRole", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackIsSealed is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x631f9852.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function isSealed() view returns(bool)
func (evidenceRegistry *EvidenceRegistry) PackIsSealed() []byte {
	enc, err := evidenceRegistry.abi.Pack("isSealed")
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
func (evidenceRegistry *EvidenceRegistry) TryPackIsSealed() ([]byte, error) {
	return evidenceRegistry.abi.Pack("isSealed")
}

// UnpackIsSealed is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x631f9852.
//
// Solidity: function isSealed() view returns(bool)
func (evidenceRegistry *EvidenceRegistry) UnpackIsSealed(data []byte) (bool, error) {
	out, err := evidenceRegistry.abi.Unpack("isSealed", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackOwner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8da5cb5b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function owner() view returns(address)
func (evidenceRegistry *EvidenceRegistry) PackOwner() []byte {
	enc, err := evidenceRegistry.abi.Pack("owner")
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
func (evidenceRegistry *EvidenceRegistry) TryPackOwner() ([]byte, error) {
	return evidenceRegistry.abi.Pack("owner")
}

// UnpackOwner is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (evidenceRegistry *EvidenceRegistry) UnpackOwner(data []byte) (common.Address, error) {
	out, err := evidenceRegistry.abi.Unpack("owner", data)
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
func (evidenceRegistry *EvidenceRegistry) PackPendingDefaultAdmin() []byte {
	enc, err := evidenceRegistry.abi.Pack("pendingDefaultAdmin")
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
func (evidenceRegistry *EvidenceRegistry) TryPackPendingDefaultAdmin() ([]byte, error) {
	return evidenceRegistry.abi.Pack("pendingDefaultAdmin")
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
func (evidenceRegistry *EvidenceRegistry) UnpackPendingDefaultAdmin(data []byte) (PendingDefaultAdminOutput, error) {
	out, err := evidenceRegistry.abi.Unpack("pendingDefaultAdmin", data)
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
func (evidenceRegistry *EvidenceRegistry) PackPendingDefaultAdminDelay() []byte {
	enc, err := evidenceRegistry.abi.Pack("pendingDefaultAdminDelay")
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
func (evidenceRegistry *EvidenceRegistry) TryPackPendingDefaultAdminDelay() ([]byte, error) {
	return evidenceRegistry.abi.Pack("pendingDefaultAdminDelay")
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
func (evidenceRegistry *EvidenceRegistry) UnpackPendingDefaultAdminDelay(data []byte) (PendingDefaultAdminDelayOutput, error) {
	out, err := evidenceRegistry.abi.Unpack("pendingDefaultAdminDelay", data)
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
func (evidenceRegistry *EvidenceRegistry) PackRenounceRole(role [32]byte, account common.Address) []byte {
	enc, err := evidenceRegistry.abi.Pack("renounceRole", role, account)
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
func (evidenceRegistry *EvidenceRegistry) TryPackRenounceRole(role [32]byte, account common.Address) ([]byte, error) {
	return evidenceRegistry.abi.Pack("renounceRole", role, account)
}

// PackRevokeRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd547741f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (evidenceRegistry *EvidenceRegistry) PackRevokeRole(role [32]byte, account common.Address) []byte {
	enc, err := evidenceRegistry.abi.Pack("revokeRole", role, account)
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
func (evidenceRegistry *EvidenceRegistry) TryPackRevokeRole(role [32]byte, account common.Address) ([]byte, error) {
	return evidenceRegistry.abi.Pack("revokeRole", role, account)
}

// PackRollbackDefaultAdminDelay is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0aa6220b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function rollbackDefaultAdminDelay() returns()
func (evidenceRegistry *EvidenceRegistry) PackRollbackDefaultAdminDelay() []byte {
	enc, err := evidenceRegistry.abi.Pack("rollbackDefaultAdminDelay")
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
func (evidenceRegistry *EvidenceRegistry) TryPackRollbackDefaultAdminDelay() ([]byte, error) {
	return evidenceRegistry.abi.Pack("rollbackDefaultAdminDelay")
}

// PackScheme is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8518b0ad.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function scheme() pure returns(uint16)
func (evidenceRegistry *EvidenceRegistry) PackScheme() []byte {
	enc, err := evidenceRegistry.abi.Pack("scheme")
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
func (evidenceRegistry *EvidenceRegistry) TryPackScheme() ([]byte, error) {
	return evidenceRegistry.abi.Pack("scheme")
}

// UnpackScheme is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8518b0ad.
//
// Solidity: function scheme() pure returns(uint16)
func (evidenceRegistry *EvidenceRegistry) UnpackScheme(data []byte) (uint16, error) {
	out, err := evidenceRegistry.abi.Unpack("scheme", data)
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
func (evidenceRegistry *EvidenceRegistry) PackSeal(expectedHeads [][32]byte) []byte {
	enc, err := evidenceRegistry.abi.Pack("seal", expectedHeads)
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
func (evidenceRegistry *EvidenceRegistry) TryPackSeal(expectedHeads [][32]byte) ([]byte, error) {
	return evidenceRegistry.abi.Pack("seal", expectedHeads)
}

// PackStateOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc7991ec7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function stateOf(uint8 stream) view returns((bool,bytes32,uint64,uint64,uint64,uint64,uint64))
func (evidenceRegistry *EvidenceRegistry) PackStateOf(stream uint8) []byte {
	enc, err := evidenceRegistry.abi.Pack("stateOf", stream)
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
func (evidenceRegistry *EvidenceRegistry) TryPackStateOf(stream uint8) ([]byte, error) {
	return evidenceRegistry.abi.Pack("stateOf", stream)
}

// UnpackStateOf is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc7991ec7.
//
// Solidity: function stateOf(uint8 stream) view returns((bool,bytes32,uint64,uint64,uint64,uint64,uint64))
func (evidenceRegistry *EvidenceRegistry) UnpackStateOf(data []byte) (StreamState, error) {
	out, err := evidenceRegistry.abi.Unpack("stateOf", data)
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
func (evidenceRegistry *EvidenceRegistry) PackStreams() []byte {
	enc, err := evidenceRegistry.abi.Pack("streams")
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
func (evidenceRegistry *EvidenceRegistry) TryPackStreams() ([]byte, error) {
	return evidenceRegistry.abi.Pack("streams")
}

// UnpackStreams is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x32403040.
//
// Solidity: function streams() view returns(uint8[])
func (evidenceRegistry *EvidenceRegistry) UnpackStreams(data []byte) ([]uint8, error) {
	out, err := evidenceRegistry.abi.Unpack("streams", data)
	if err != nil {
		return *new([]uint8), err
	}
	out0 := *abi.ConvertType(out[0], new([]uint8)).(*[]uint8)
	return out0, nil
}

// PackSupportsInterface is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x01ffc9a7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (evidenceRegistry *EvidenceRegistry) PackSupportsInterface(interfaceId [4]byte) []byte {
	enc, err := evidenceRegistry.abi.Pack("supportsInterface", interfaceId)
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
func (evidenceRegistry *EvidenceRegistry) TryPackSupportsInterface(interfaceId [4]byte) ([]byte, error) {
	return evidenceRegistry.abi.Pack("supportsInterface", interfaceId)
}

// UnpackSupportsInterface is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (evidenceRegistry *EvidenceRegistry) UnpackSupportsInterface(data []byte) (bool, error) {
	out, err := evidenceRegistry.abi.Unpack("supportsInterface", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackVerifyInclusion is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe283c044.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function verifyInclusion(uint8 stream, uint64 number, bytes32 salt, bytes32 recordDigest, uint64 index, bytes32[] siblings) view returns(bool)
func (evidenceRegistry *EvidenceRegistry) PackVerifyInclusion(stream uint8, number uint64, salt [32]byte, recordDigest [32]byte, index uint64, siblings [][32]byte) []byte {
	enc, err := evidenceRegistry.abi.Pack("verifyInclusion", stream, number, salt, recordDigest, index, siblings)
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
func (evidenceRegistry *EvidenceRegistry) TryPackVerifyInclusion(stream uint8, number uint64, salt [32]byte, recordDigest [32]byte, index uint64, siblings [][32]byte) ([]byte, error) {
	return evidenceRegistry.abi.Pack("verifyInclusion", stream, number, salt, recordDigest, index, siblings)
}

// UnpackVerifyInclusion is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe283c044.
//
// Solidity: function verifyInclusion(uint8 stream, uint64 number, bytes32 salt, bytes32 recordDigest, uint64 index, bytes32[] siblings) view returns(bool)
func (evidenceRegistry *EvidenceRegistry) UnpackVerifyInclusion(data []byte) (bool, error) {
	out, err := evidenceRegistry.abi.Unpack("verifyInclusion", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// EvidenceRegistryBatchAnchored represents a BatchAnchored event raised by the EvidenceRegistry contract.
type EvidenceRegistryBatchAnchored struct {
	Stream uint8
	Number uint64
	Digest [32]byte
	Header BatchHeader
	Raw    *types.Log // Blockchain specific contextual infos
}

const EvidenceRegistryBatchAnchoredEventName = "BatchAnchored"

// ContractEventName returns the user-defined event name.
func (EvidenceRegistryBatchAnchored) ContractEventName() string {
	return EvidenceRegistryBatchAnchoredEventName
}

// UnpackBatchAnchoredEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event BatchAnchored(uint8 indexed stream, uint64 indexed number, bytes32 digest, (uint8,uint64,bytes32,bytes32,uint64,uint64,uint64,uint64) header)
func (evidenceRegistry *EvidenceRegistry) UnpackBatchAnchoredEvent(log *types.Log) (*EvidenceRegistryBatchAnchored, error) {
	event := "BatchAnchored"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != evidenceRegistry.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(EvidenceRegistryBatchAnchored)
	if len(log.Data) > 0 {
		if err := evidenceRegistry.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range evidenceRegistry.abi.Events[event].Inputs {
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

// EvidenceRegistryDefaultAdminDelayChangeCanceled represents a DefaultAdminDelayChangeCanceled event raised by the EvidenceRegistry contract.
type EvidenceRegistryDefaultAdminDelayChangeCanceled struct {
	Raw *types.Log // Blockchain specific contextual infos
}

const EvidenceRegistryDefaultAdminDelayChangeCanceledEventName = "DefaultAdminDelayChangeCanceled"

// ContractEventName returns the user-defined event name.
func (EvidenceRegistryDefaultAdminDelayChangeCanceled) ContractEventName() string {
	return EvidenceRegistryDefaultAdminDelayChangeCanceledEventName
}

// UnpackDefaultAdminDelayChangeCanceledEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event DefaultAdminDelayChangeCanceled()
func (evidenceRegistry *EvidenceRegistry) UnpackDefaultAdminDelayChangeCanceledEvent(log *types.Log) (*EvidenceRegistryDefaultAdminDelayChangeCanceled, error) {
	event := "DefaultAdminDelayChangeCanceled"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != evidenceRegistry.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(EvidenceRegistryDefaultAdminDelayChangeCanceled)
	if len(log.Data) > 0 {
		if err := evidenceRegistry.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range evidenceRegistry.abi.Events[event].Inputs {
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

// EvidenceRegistryDefaultAdminDelayChangeScheduled represents a DefaultAdminDelayChangeScheduled event raised by the EvidenceRegistry contract.
type EvidenceRegistryDefaultAdminDelayChangeScheduled struct {
	NewDelay       *big.Int
	EffectSchedule *big.Int
	Raw            *types.Log // Blockchain specific contextual infos
}

const EvidenceRegistryDefaultAdminDelayChangeScheduledEventName = "DefaultAdminDelayChangeScheduled"

// ContractEventName returns the user-defined event name.
func (EvidenceRegistryDefaultAdminDelayChangeScheduled) ContractEventName() string {
	return EvidenceRegistryDefaultAdminDelayChangeScheduledEventName
}

// UnpackDefaultAdminDelayChangeScheduledEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event DefaultAdminDelayChangeScheduled(uint48 newDelay, uint48 effectSchedule)
func (evidenceRegistry *EvidenceRegistry) UnpackDefaultAdminDelayChangeScheduledEvent(log *types.Log) (*EvidenceRegistryDefaultAdminDelayChangeScheduled, error) {
	event := "DefaultAdminDelayChangeScheduled"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != evidenceRegistry.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(EvidenceRegistryDefaultAdminDelayChangeScheduled)
	if len(log.Data) > 0 {
		if err := evidenceRegistry.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range evidenceRegistry.abi.Events[event].Inputs {
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

// EvidenceRegistryDefaultAdminTransferCanceled represents a DefaultAdminTransferCanceled event raised by the EvidenceRegistry contract.
type EvidenceRegistryDefaultAdminTransferCanceled struct {
	Raw *types.Log // Blockchain specific contextual infos
}

const EvidenceRegistryDefaultAdminTransferCanceledEventName = "DefaultAdminTransferCanceled"

// ContractEventName returns the user-defined event name.
func (EvidenceRegistryDefaultAdminTransferCanceled) ContractEventName() string {
	return EvidenceRegistryDefaultAdminTransferCanceledEventName
}

// UnpackDefaultAdminTransferCanceledEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event DefaultAdminTransferCanceled()
func (evidenceRegistry *EvidenceRegistry) UnpackDefaultAdminTransferCanceledEvent(log *types.Log) (*EvidenceRegistryDefaultAdminTransferCanceled, error) {
	event := "DefaultAdminTransferCanceled"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != evidenceRegistry.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(EvidenceRegistryDefaultAdminTransferCanceled)
	if len(log.Data) > 0 {
		if err := evidenceRegistry.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range evidenceRegistry.abi.Events[event].Inputs {
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

// EvidenceRegistryDefaultAdminTransferScheduled represents a DefaultAdminTransferScheduled event raised by the EvidenceRegistry contract.
type EvidenceRegistryDefaultAdminTransferScheduled struct {
	NewAdmin       common.Address
	AcceptSchedule *big.Int
	Raw            *types.Log // Blockchain specific contextual infos
}

const EvidenceRegistryDefaultAdminTransferScheduledEventName = "DefaultAdminTransferScheduled"

// ContractEventName returns the user-defined event name.
func (EvidenceRegistryDefaultAdminTransferScheduled) ContractEventName() string {
	return EvidenceRegistryDefaultAdminTransferScheduledEventName
}

// UnpackDefaultAdminTransferScheduledEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event DefaultAdminTransferScheduled(address indexed newAdmin, uint48 acceptSchedule)
func (evidenceRegistry *EvidenceRegistry) UnpackDefaultAdminTransferScheduledEvent(log *types.Log) (*EvidenceRegistryDefaultAdminTransferScheduled, error) {
	event := "DefaultAdminTransferScheduled"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != evidenceRegistry.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(EvidenceRegistryDefaultAdminTransferScheduled)
	if len(log.Data) > 0 {
		if err := evidenceRegistry.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range evidenceRegistry.abi.Events[event].Inputs {
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

// EvidenceRegistryRoleAdminChanged represents a RoleAdminChanged event raised by the EvidenceRegistry contract.
type EvidenceRegistryRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               *types.Log // Blockchain specific contextual infos
}

const EvidenceRegistryRoleAdminChangedEventName = "RoleAdminChanged"

// ContractEventName returns the user-defined event name.
func (EvidenceRegistryRoleAdminChanged) ContractEventName() string {
	return EvidenceRegistryRoleAdminChangedEventName
}

// UnpackRoleAdminChangedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (evidenceRegistry *EvidenceRegistry) UnpackRoleAdminChangedEvent(log *types.Log) (*EvidenceRegistryRoleAdminChanged, error) {
	event := "RoleAdminChanged"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != evidenceRegistry.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(EvidenceRegistryRoleAdminChanged)
	if len(log.Data) > 0 {
		if err := evidenceRegistry.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range evidenceRegistry.abi.Events[event].Inputs {
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

// EvidenceRegistryRoleGranted represents a RoleGranted event raised by the EvidenceRegistry contract.
type EvidenceRegistryRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     *types.Log // Blockchain specific contextual infos
}

const EvidenceRegistryRoleGrantedEventName = "RoleGranted"

// ContractEventName returns the user-defined event name.
func (EvidenceRegistryRoleGranted) ContractEventName() string {
	return EvidenceRegistryRoleGrantedEventName
}

// UnpackRoleGrantedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (evidenceRegistry *EvidenceRegistry) UnpackRoleGrantedEvent(log *types.Log) (*EvidenceRegistryRoleGranted, error) {
	event := "RoleGranted"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != evidenceRegistry.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(EvidenceRegistryRoleGranted)
	if len(log.Data) > 0 {
		if err := evidenceRegistry.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range evidenceRegistry.abi.Events[event].Inputs {
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

// EvidenceRegistryRoleRevoked represents a RoleRevoked event raised by the EvidenceRegistry contract.
type EvidenceRegistryRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     *types.Log // Blockchain specific contextual infos
}

const EvidenceRegistryRoleRevokedEventName = "RoleRevoked"

// ContractEventName returns the user-defined event name.
func (EvidenceRegistryRoleRevoked) ContractEventName() string {
	return EvidenceRegistryRoleRevokedEventName
}

// UnpackRoleRevokedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (evidenceRegistry *EvidenceRegistry) UnpackRoleRevokedEvent(log *types.Log) (*EvidenceRegistryRoleRevoked, error) {
	event := "RoleRevoked"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != evidenceRegistry.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(EvidenceRegistryRoleRevoked)
	if len(log.Data) > 0 {
		if err := evidenceRegistry.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range evidenceRegistry.abi.Events[event].Inputs {
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

// EvidenceRegistryStreamSealed represents a StreamSealed event raised by the EvidenceRegistry contract.
type EvidenceRegistryStreamSealed struct {
	Stream             uint8
	Head               [32]byte
	NextBatch          uint64
	LastClosedAt       uint64
	ThroughTransaction uint64
	ThroughSequence    uint64
	Raw                *types.Log // Blockchain specific contextual infos
}

const EvidenceRegistryStreamSealedEventName = "StreamSealed"

// ContractEventName returns the user-defined event name.
func (EvidenceRegistryStreamSealed) ContractEventName() string {
	return EvidenceRegistryStreamSealedEventName
}

// UnpackStreamSealedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event StreamSealed(uint8 indexed stream, bytes32 head, uint64 nextBatch, uint64 lastClosedAt, uint64 throughTransaction, uint64 throughSequence)
func (evidenceRegistry *EvidenceRegistry) UnpackStreamSealedEvent(log *types.Log) (*EvidenceRegistryStreamSealed, error) {
	event := "StreamSealed"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != evidenceRegistry.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(EvidenceRegistryStreamSealed)
	if len(log.Data) > 0 {
		if err := evidenceRegistry.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range evidenceRegistry.abi.Events[event].Inputs {
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
func (evidenceRegistry *EvidenceRegistry) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], evidenceRegistry.abi.Errors["AccessControlBadConfirmation"].ID.Bytes()[:4]) {
		return evidenceRegistry.UnpackAccessControlBadConfirmationError(raw[4:])
	}
	if bytes.Equal(raw[:4], evidenceRegistry.abi.Errors["AccessControlEnforcedDefaultAdminDelay"].ID.Bytes()[:4]) {
		return evidenceRegistry.UnpackAccessControlEnforcedDefaultAdminDelayError(raw[4:])
	}
	if bytes.Equal(raw[:4], evidenceRegistry.abi.Errors["AccessControlEnforcedDefaultAdminRules"].ID.Bytes()[:4]) {
		return evidenceRegistry.UnpackAccessControlEnforcedDefaultAdminRulesError(raw[4:])
	}
	if bytes.Equal(raw[:4], evidenceRegistry.abi.Errors["AccessControlInvalidDefaultAdmin"].ID.Bytes()[:4]) {
		return evidenceRegistry.UnpackAccessControlInvalidDefaultAdminError(raw[4:])
	}
	if bytes.Equal(raw[:4], evidenceRegistry.abi.Errors["AccessControlUnauthorizedAccount"].ID.Bytes()[:4]) {
		return evidenceRegistry.UnpackAccessControlUnauthorizedAccountError(raw[4:])
	}
	if bytes.Equal(raw[:4], evidenceRegistry.abi.Errors["AnchorerNotGranted"].ID.Bytes()[:4]) {
		return evidenceRegistry.UnpackAnchorerNotGrantedError(raw[4:])
	}
	if bytes.Equal(raw[:4], evidenceRegistry.abi.Errors["BatchAlreadyAnchored"].ID.Bytes()[:4]) {
		return evidenceRegistry.UnpackBatchAlreadyAnchoredError(raw[4:])
	}
	if bytes.Equal(raw[:4], evidenceRegistry.abi.Errors["BatchBeforeGenesis"].ID.Bytes()[:4]) {
		return evidenceRegistry.UnpackBatchBeforeGenesisError(raw[4:])
	}
	if bytes.Equal(raw[:4], evidenceRegistry.abi.Errors["BatchConflict"].ID.Bytes()[:4]) {
		return evidenceRegistry.UnpackBatchConflictError(raw[4:])
	}
	if bytes.Equal(raw[:4], evidenceRegistry.abi.Errors["BatchOutOfSequence"].ID.Bytes()[:4]) {
		return evidenceRegistry.UnpackBatchOutOfSequenceError(raw[4:])
	}
	if bytes.Equal(raw[:4], evidenceRegistry.abi.Errors["CloseInFuture"].ID.Bytes()[:4]) {
		return evidenceRegistry.UnpackCloseInFutureError(raw[4:])
	}
	if bytes.Equal(raw[:4], evidenceRegistry.abi.Errors["CloseNotAfterPrevious"].ID.Bytes()[:4]) {
		return evidenceRegistry.UnpackCloseNotAfterPreviousError(raw[4:])
	}
	if bytes.Equal(raw[:4], evidenceRegistry.abi.Errors["EmptyBatch"].ID.Bytes()[:4]) {
		return evidenceRegistry.UnpackEmptyBatchError(raw[4:])
	}
	if bytes.Equal(raw[:4], evidenceRegistry.abi.Errors["HeadCountMismatch"].ID.Bytes()[:4]) {
		return evidenceRegistry.UnpackHeadCountMismatchError(raw[4:])
	}
	if bytes.Equal(raw[:4], evidenceRegistry.abi.Errors["HeadMismatch"].ID.Bytes()[:4]) {
		return evidenceRegistry.UnpackHeadMismatchError(raw[4:])
	}
	if bytes.Equal(raw[:4], evidenceRegistry.abi.Errors["InvalidGenesis"].ID.Bytes()[:4]) {
		return evidenceRegistry.UnpackInvalidGenesisError(raw[4:])
	}
	if bytes.Equal(raw[:4], evidenceRegistry.abi.Errors["NoStreams"].ID.Bytes()[:4]) {
		return evidenceRegistry.UnpackNoStreamsError(raw[4:])
	}
	if bytes.Equal(raw[:4], evidenceRegistry.abi.Errors["PositionNotAfterPrevious"].ID.Bytes()[:4]) {
		return evidenceRegistry.UnpackPositionNotAfterPreviousError(raw[4:])
	}
	if bytes.Equal(raw[:4], evidenceRegistry.abi.Errors["PreviousDigestMismatch"].ID.Bytes()[:4]) {
		return evidenceRegistry.UnpackPreviousDigestMismatchError(raw[4:])
	}
	if bytes.Equal(raw[:4], evidenceRegistry.abi.Errors["RegistryIsSealed"].ID.Bytes()[:4]) {
		return evidenceRegistry.UnpackRegistryIsSealedError(raw[4:])
	}
	if bytes.Equal(raw[:4], evidenceRegistry.abi.Errors["SafeCastOverflowedUintDowncast"].ID.Bytes()[:4]) {
		return evidenceRegistry.UnpackSafeCastOverflowedUintDowncastError(raw[4:])
	}
	if bytes.Equal(raw[:4], evidenceRegistry.abi.Errors["StreamsNotAscending"].ID.Bytes()[:4]) {
		return evidenceRegistry.UnpackStreamsNotAscendingError(raw[4:])
	}
	if bytes.Equal(raw[:4], evidenceRegistry.abi.Errors["TooManyLeaves"].ID.Bytes()[:4]) {
		return evidenceRegistry.UnpackTooManyLeavesError(raw[4:])
	}
	if bytes.Equal(raw[:4], evidenceRegistry.abi.Errors["TooManyStreams"].ID.Bytes()[:4]) {
		return evidenceRegistry.UnpackTooManyStreamsError(raw[4:])
	}
	if bytes.Equal(raw[:4], evidenceRegistry.abi.Errors["UnknownBatch"].ID.Bytes()[:4]) {
		return evidenceRegistry.UnpackUnknownBatchError(raw[4:])
	}
	if bytes.Equal(raw[:4], evidenceRegistry.abi.Errors["UnknownStream"].ID.Bytes()[:4]) {
		return evidenceRegistry.UnpackUnknownStreamError(raw[4:])
	}
	if bytes.Equal(raw[:4], evidenceRegistry.abi.Errors["ZeroAnchorer"].ID.Bytes()[:4]) {
		return evidenceRegistry.UnpackZeroAnchorerError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// EvidenceRegistryAccessControlBadConfirmation represents a AccessControlBadConfirmation error raised by the EvidenceRegistry contract.
type EvidenceRegistryAccessControlBadConfirmation struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AccessControlBadConfirmation()
func EvidenceRegistryAccessControlBadConfirmationErrorID() common.Hash {
	return common.HexToHash("0x6697b23232a647058342c0724fe7c415cab25915b54e5dbc03f233173d37b41c")
}

// UnpackAccessControlBadConfirmationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AccessControlBadConfirmation()
func (evidenceRegistry *EvidenceRegistry) UnpackAccessControlBadConfirmationError(raw []byte) (*EvidenceRegistryAccessControlBadConfirmation, error) {
	out := new(EvidenceRegistryAccessControlBadConfirmation)
	if err := evidenceRegistry.abi.UnpackIntoInterface(out, "AccessControlBadConfirmation", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EvidenceRegistryAccessControlEnforcedDefaultAdminDelay represents a AccessControlEnforcedDefaultAdminDelay error raised by the EvidenceRegistry contract.
type EvidenceRegistryAccessControlEnforcedDefaultAdminDelay struct {
	Schedule *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AccessControlEnforcedDefaultAdminDelay(uint48 schedule)
func EvidenceRegistryAccessControlEnforcedDefaultAdminDelayErrorID() common.Hash {
	return common.HexToHash("0x19ca5ebb8fb33f00e502c9392eddab1501674629178bf69b853cf037aaf4bb5d")
}

// UnpackAccessControlEnforcedDefaultAdminDelayError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AccessControlEnforcedDefaultAdminDelay(uint48 schedule)
func (evidenceRegistry *EvidenceRegistry) UnpackAccessControlEnforcedDefaultAdminDelayError(raw []byte) (*EvidenceRegistryAccessControlEnforcedDefaultAdminDelay, error) {
	out := new(EvidenceRegistryAccessControlEnforcedDefaultAdminDelay)
	if err := evidenceRegistry.abi.UnpackIntoInterface(out, "AccessControlEnforcedDefaultAdminDelay", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EvidenceRegistryAccessControlEnforcedDefaultAdminRules represents a AccessControlEnforcedDefaultAdminRules error raised by the EvidenceRegistry contract.
type EvidenceRegistryAccessControlEnforcedDefaultAdminRules struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AccessControlEnforcedDefaultAdminRules()
func EvidenceRegistryAccessControlEnforcedDefaultAdminRulesErrorID() common.Hash {
	return common.HexToHash("0x3fc3c27ae3db78c81b8f6e685172134623efa268ee8cd8d54be38ad2a74fc13b")
}

// UnpackAccessControlEnforcedDefaultAdminRulesError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AccessControlEnforcedDefaultAdminRules()
func (evidenceRegistry *EvidenceRegistry) UnpackAccessControlEnforcedDefaultAdminRulesError(raw []byte) (*EvidenceRegistryAccessControlEnforcedDefaultAdminRules, error) {
	out := new(EvidenceRegistryAccessControlEnforcedDefaultAdminRules)
	if err := evidenceRegistry.abi.UnpackIntoInterface(out, "AccessControlEnforcedDefaultAdminRules", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EvidenceRegistryAccessControlInvalidDefaultAdmin represents a AccessControlInvalidDefaultAdmin error raised by the EvidenceRegistry contract.
type EvidenceRegistryAccessControlInvalidDefaultAdmin struct {
	DefaultAdmin common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AccessControlInvalidDefaultAdmin(address defaultAdmin)
func EvidenceRegistryAccessControlInvalidDefaultAdminErrorID() common.Hash {
	return common.HexToHash("0xc22c8022f2a840d6b6a9f113407715f5bbd4e88c1b0dd9434dc00700ba609ed4")
}

// UnpackAccessControlInvalidDefaultAdminError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AccessControlInvalidDefaultAdmin(address defaultAdmin)
func (evidenceRegistry *EvidenceRegistry) UnpackAccessControlInvalidDefaultAdminError(raw []byte) (*EvidenceRegistryAccessControlInvalidDefaultAdmin, error) {
	out := new(EvidenceRegistryAccessControlInvalidDefaultAdmin)
	if err := evidenceRegistry.abi.UnpackIntoInterface(out, "AccessControlInvalidDefaultAdmin", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EvidenceRegistryAccessControlUnauthorizedAccount represents a AccessControlUnauthorizedAccount error raised by the EvidenceRegistry contract.
type EvidenceRegistryAccessControlUnauthorizedAccount struct {
	Account    common.Address
	NeededRole [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AccessControlUnauthorizedAccount(address account, bytes32 neededRole)
func EvidenceRegistryAccessControlUnauthorizedAccountErrorID() common.Hash {
	return common.HexToHash("0xe2517d3fbfae6f8515ef5ff1ccedc3933ab0cbbda0b492c06eb54ad10ef03b3e")
}

// UnpackAccessControlUnauthorizedAccountError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AccessControlUnauthorizedAccount(address account, bytes32 neededRole)
func (evidenceRegistry *EvidenceRegistry) UnpackAccessControlUnauthorizedAccountError(raw []byte) (*EvidenceRegistryAccessControlUnauthorizedAccount, error) {
	out := new(EvidenceRegistryAccessControlUnauthorizedAccount)
	if err := evidenceRegistry.abi.UnpackIntoInterface(out, "AccessControlUnauthorizedAccount", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EvidenceRegistryAnchorerNotGranted represents a AnchorerNotGranted error raised by the EvidenceRegistry contract.
type EvidenceRegistryAnchorerNotGranted struct {
	Anchorer common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AnchorerNotGranted(address anchorer)
func EvidenceRegistryAnchorerNotGrantedErrorID() common.Hash {
	return common.HexToHash("0x3012fc8bb5dd75c0c034d7c58aa1982d269d598072100fabc24d4de862d1e083")
}

// UnpackAnchorerNotGrantedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AnchorerNotGranted(address anchorer)
func (evidenceRegistry *EvidenceRegistry) UnpackAnchorerNotGrantedError(raw []byte) (*EvidenceRegistryAnchorerNotGranted, error) {
	out := new(EvidenceRegistryAnchorerNotGranted)
	if err := evidenceRegistry.abi.UnpackIntoInterface(out, "AnchorerNotGranted", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EvidenceRegistryBatchAlreadyAnchored represents a BatchAlreadyAnchored error raised by the EvidenceRegistry contract.
type EvidenceRegistryBatchAlreadyAnchored struct {
	Stream uint8
	Number uint64
	Digest [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error BatchAlreadyAnchored(uint8 stream, uint64 number, bytes32 digest)
func EvidenceRegistryBatchAlreadyAnchoredErrorID() common.Hash {
	return common.HexToHash("0x6035ccca362de798c967fc5ddf5a627f48878f68f274d1ac779fe416246930c3")
}

// UnpackBatchAlreadyAnchoredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error BatchAlreadyAnchored(uint8 stream, uint64 number, bytes32 digest)
func (evidenceRegistry *EvidenceRegistry) UnpackBatchAlreadyAnchoredError(raw []byte) (*EvidenceRegistryBatchAlreadyAnchored, error) {
	out := new(EvidenceRegistryBatchAlreadyAnchored)
	if err := evidenceRegistry.abi.UnpackIntoInterface(out, "BatchAlreadyAnchored", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EvidenceRegistryBatchBeforeGenesis represents a BatchBeforeGenesis error raised by the EvidenceRegistry contract.
type EvidenceRegistryBatchBeforeGenesis struct {
	Stream       uint8
	GenesisBatch uint64
	Number       uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error BatchBeforeGenesis(uint8 stream, uint64 genesisBatch, uint64 number)
func EvidenceRegistryBatchBeforeGenesisErrorID() common.Hash {
	return common.HexToHash("0x0d0878958571059a7616862a03367cf877644fac9d790878e6167d68f2f289d7")
}

// UnpackBatchBeforeGenesisError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error BatchBeforeGenesis(uint8 stream, uint64 genesisBatch, uint64 number)
func (evidenceRegistry *EvidenceRegistry) UnpackBatchBeforeGenesisError(raw []byte) (*EvidenceRegistryBatchBeforeGenesis, error) {
	out := new(EvidenceRegistryBatchBeforeGenesis)
	if err := evidenceRegistry.abi.UnpackIntoInterface(out, "BatchBeforeGenesis", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EvidenceRegistryBatchConflict represents a BatchConflict error raised by the EvidenceRegistry contract.
type EvidenceRegistryBatchConflict struct {
	Stream   uint8
	Number   uint64
	Anchored [32]byte
	Offered  [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error BatchConflict(uint8 stream, uint64 number, bytes32 anchored, bytes32 offered)
func EvidenceRegistryBatchConflictErrorID() common.Hash {
	return common.HexToHash("0xc071754cfdbb9c05bd1a514ef76b5572b760af6ef162ed46016afd316167776a")
}

// UnpackBatchConflictError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error BatchConflict(uint8 stream, uint64 number, bytes32 anchored, bytes32 offered)
func (evidenceRegistry *EvidenceRegistry) UnpackBatchConflictError(raw []byte) (*EvidenceRegistryBatchConflict, error) {
	out := new(EvidenceRegistryBatchConflict)
	if err := evidenceRegistry.abi.UnpackIntoInterface(out, "BatchConflict", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EvidenceRegistryBatchOutOfSequence represents a BatchOutOfSequence error raised by the EvidenceRegistry contract.
type EvidenceRegistryBatchOutOfSequence struct {
	Stream   uint8
	Expected uint64
	Offered  uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error BatchOutOfSequence(uint8 stream, uint64 expected, uint64 offered)
func EvidenceRegistryBatchOutOfSequenceErrorID() common.Hash {
	return common.HexToHash("0xb39c3985f6b4ab5807374509d65a2c99793b70bd7acb0641f36bce34c38b3faa")
}

// UnpackBatchOutOfSequenceError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error BatchOutOfSequence(uint8 stream, uint64 expected, uint64 offered)
func (evidenceRegistry *EvidenceRegistry) UnpackBatchOutOfSequenceError(raw []byte) (*EvidenceRegistryBatchOutOfSequence, error) {
	out := new(EvidenceRegistryBatchOutOfSequence)
	if err := evidenceRegistry.abi.UnpackIntoInterface(out, "BatchOutOfSequence", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EvidenceRegistryCloseInFuture represents a CloseInFuture error raised by the EvidenceRegistry contract.
type EvidenceRegistryCloseInFuture struct {
	Stream    uint8
	ClosedAt  uint64
	BlockTime uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error CloseInFuture(uint8 stream, uint64 closedAt, uint64 blockTime)
func EvidenceRegistryCloseInFutureErrorID() common.Hash {
	return common.HexToHash("0xe9f968e7bdcd7fc5b23f91fd07ad9a18f73c622fa09f8198abdfdee58064dab5")
}

// UnpackCloseInFutureError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error CloseInFuture(uint8 stream, uint64 closedAt, uint64 blockTime)
func (evidenceRegistry *EvidenceRegistry) UnpackCloseInFutureError(raw []byte) (*EvidenceRegistryCloseInFuture, error) {
	out := new(EvidenceRegistryCloseInFuture)
	if err := evidenceRegistry.abi.UnpackIntoInterface(out, "CloseInFuture", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EvidenceRegistryCloseNotAfterPrevious represents a CloseNotAfterPrevious error raised by the EvidenceRegistry contract.
type EvidenceRegistryCloseNotAfterPrevious struct {
	Stream       uint8
	ClosedAt     uint64
	LastClosedAt uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error CloseNotAfterPrevious(uint8 stream, uint64 closedAt, uint64 lastClosedAt)
func EvidenceRegistryCloseNotAfterPreviousErrorID() common.Hash {
	return common.HexToHash("0x5999344d325d28b495bd9ad3396675712366c905f2fbbcb3621b879e92eb13fc")
}

// UnpackCloseNotAfterPreviousError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error CloseNotAfterPrevious(uint8 stream, uint64 closedAt, uint64 lastClosedAt)
func (evidenceRegistry *EvidenceRegistry) UnpackCloseNotAfterPreviousError(raw []byte) (*EvidenceRegistryCloseNotAfterPrevious, error) {
	out := new(EvidenceRegistryCloseNotAfterPrevious)
	if err := evidenceRegistry.abi.UnpackIntoInterface(out, "CloseNotAfterPrevious", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EvidenceRegistryEmptyBatch represents a EmptyBatch error raised by the EvidenceRegistry contract.
type EvidenceRegistryEmptyBatch struct {
	Stream uint8
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error EmptyBatch(uint8 stream)
func EvidenceRegistryEmptyBatchErrorID() common.Hash {
	return common.HexToHash("0x2aa62fa976a6325b695bf70c580caa4fea6f6fbc4ac0bca15c9492be85d38acf")
}

// UnpackEmptyBatchError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error EmptyBatch(uint8 stream)
func (evidenceRegistry *EvidenceRegistry) UnpackEmptyBatchError(raw []byte) (*EvidenceRegistryEmptyBatch, error) {
	out := new(EvidenceRegistryEmptyBatch)
	if err := evidenceRegistry.abi.UnpackIntoInterface(out, "EmptyBatch", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EvidenceRegistryHeadCountMismatch represents a HeadCountMismatch error raised by the EvidenceRegistry contract.
type EvidenceRegistryHeadCountMismatch struct {
	Streams *big.Int
	Offered *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error HeadCountMismatch(uint256 streams, uint256 offered)
func EvidenceRegistryHeadCountMismatchErrorID() common.Hash {
	return common.HexToHash("0xbe8ea0e772df6ee158b79c8f7156d97dcac0737b032a43f94f9b385463afe2f9")
}

// UnpackHeadCountMismatchError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error HeadCountMismatch(uint256 streams, uint256 offered)
func (evidenceRegistry *EvidenceRegistry) UnpackHeadCountMismatchError(raw []byte) (*EvidenceRegistryHeadCountMismatch, error) {
	out := new(EvidenceRegistryHeadCountMismatch)
	if err := evidenceRegistry.abi.UnpackIntoInterface(out, "HeadCountMismatch", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EvidenceRegistryHeadMismatch represents a HeadMismatch error raised by the EvidenceRegistry contract.
type EvidenceRegistryHeadMismatch struct {
	Stream   uint8
	Head     [32]byte
	Expected [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error HeadMismatch(uint8 stream, bytes32 head, bytes32 expected)
func EvidenceRegistryHeadMismatchErrorID() common.Hash {
	return common.HexToHash("0xf18e8ed0e05571020e598594dc2d046464d499ec16c478bf278c91d7362bde85")
}

// UnpackHeadMismatchError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error HeadMismatch(uint8 stream, bytes32 head, bytes32 expected)
func (evidenceRegistry *EvidenceRegistry) UnpackHeadMismatchError(raw []byte) (*EvidenceRegistryHeadMismatch, error) {
	out := new(EvidenceRegistryHeadMismatch)
	if err := evidenceRegistry.abi.UnpackIntoInterface(out, "HeadMismatch", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EvidenceRegistryInvalidGenesis represents a InvalidGenesis error raised by the EvidenceRegistry contract.
type EvidenceRegistryInvalidGenesis struct {
	Stream uint8
	Batch  uint64
	Digest [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidGenesis(uint8 stream, uint64 batch, bytes32 digest)
func EvidenceRegistryInvalidGenesisErrorID() common.Hash {
	return common.HexToHash("0xad14a26107aa7d20646398124017110e79feed4739667c3bf81dfec51aceb429")
}

// UnpackInvalidGenesisError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidGenesis(uint8 stream, uint64 batch, bytes32 digest)
func (evidenceRegistry *EvidenceRegistry) UnpackInvalidGenesisError(raw []byte) (*EvidenceRegistryInvalidGenesis, error) {
	out := new(EvidenceRegistryInvalidGenesis)
	if err := evidenceRegistry.abi.UnpackIntoInterface(out, "InvalidGenesis", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EvidenceRegistryNoStreams represents a NoStreams error raised by the EvidenceRegistry contract.
type EvidenceRegistryNoStreams struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NoStreams()
func EvidenceRegistryNoStreamsErrorID() common.Hash {
	return common.HexToHash("0x667ada5b150a16e1d46a52935ccb87bcb03cecee4f101cba3897ac4f9a77143f")
}

// UnpackNoStreamsError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NoStreams()
func (evidenceRegistry *EvidenceRegistry) UnpackNoStreamsError(raw []byte) (*EvidenceRegistryNoStreams, error) {
	out := new(EvidenceRegistryNoStreams)
	if err := evidenceRegistry.abi.UnpackIntoInterface(out, "NoStreams", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EvidenceRegistryPositionNotAfterPrevious represents a PositionNotAfterPrevious error raised by the EvidenceRegistry contract.
type EvidenceRegistryPositionNotAfterPrevious struct {
	Stream              uint8
	ThroughTransaction  uint64
	ThroughSequence     uint64
	PreviousTransaction uint64
	PreviousSequence    uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error PositionNotAfterPrevious(uint8 stream, uint64 throughTransaction, uint64 throughSequence, uint64 previousTransaction, uint64 previousSequence)
func EvidenceRegistryPositionNotAfterPreviousErrorID() common.Hash {
	return common.HexToHash("0x9e8b48b4519e7edd0cad5dc17e4e2df6d8c29b0a42851e5994ce09adc55f67b5")
}

// UnpackPositionNotAfterPreviousError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error PositionNotAfterPrevious(uint8 stream, uint64 throughTransaction, uint64 throughSequence, uint64 previousTransaction, uint64 previousSequence)
func (evidenceRegistry *EvidenceRegistry) UnpackPositionNotAfterPreviousError(raw []byte) (*EvidenceRegistryPositionNotAfterPrevious, error) {
	out := new(EvidenceRegistryPositionNotAfterPrevious)
	if err := evidenceRegistry.abi.UnpackIntoInterface(out, "PositionNotAfterPrevious", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EvidenceRegistryPreviousDigestMismatch represents a PreviousDigestMismatch error raised by the EvidenceRegistry contract.
type EvidenceRegistryPreviousDigestMismatch struct {
	Stream  uint8
	Head    [32]byte
	Offered [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error PreviousDigestMismatch(uint8 stream, bytes32 head, bytes32 offered)
func EvidenceRegistryPreviousDigestMismatchErrorID() common.Hash {
	return common.HexToHash("0x77a6bbb39f36ad092c0e2b13b8874b5c34d92fe800a43539f1094ea12480bff2")
}

// UnpackPreviousDigestMismatchError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error PreviousDigestMismatch(uint8 stream, bytes32 head, bytes32 offered)
func (evidenceRegistry *EvidenceRegistry) UnpackPreviousDigestMismatchError(raw []byte) (*EvidenceRegistryPreviousDigestMismatch, error) {
	out := new(EvidenceRegistryPreviousDigestMismatch)
	if err := evidenceRegistry.abi.UnpackIntoInterface(out, "PreviousDigestMismatch", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EvidenceRegistryRegistryIsSealed represents a RegistryIsSealed error raised by the EvidenceRegistry contract.
type EvidenceRegistryRegistryIsSealed struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RegistryIsSealed()
func EvidenceRegistryRegistryIsSealedErrorID() common.Hash {
	return common.HexToHash("0xf3e6ff318836f081aa74168a8619b76f55c7d06015c7eb8be8b0b2fe870504c0")
}

// UnpackRegistryIsSealedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RegistryIsSealed()
func (evidenceRegistry *EvidenceRegistry) UnpackRegistryIsSealedError(raw []byte) (*EvidenceRegistryRegistryIsSealed, error) {
	out := new(EvidenceRegistryRegistryIsSealed)
	if err := evidenceRegistry.abi.UnpackIntoInterface(out, "RegistryIsSealed", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EvidenceRegistrySafeCastOverflowedUintDowncast represents a SafeCastOverflowedUintDowncast error raised by the EvidenceRegistry contract.
type EvidenceRegistrySafeCastOverflowedUintDowncast struct {
	Bits  uint8
	Value *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SafeCastOverflowedUintDowncast(uint8 bits, uint256 value)
func EvidenceRegistrySafeCastOverflowedUintDowncastErrorID() common.Hash {
	return common.HexToHash("0x6dfcc6503a32754ce7a89698e18201fc5294fd4aad43edefee786f88423b1a12")
}

// UnpackSafeCastOverflowedUintDowncastError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SafeCastOverflowedUintDowncast(uint8 bits, uint256 value)
func (evidenceRegistry *EvidenceRegistry) UnpackSafeCastOverflowedUintDowncastError(raw []byte) (*EvidenceRegistrySafeCastOverflowedUintDowncast, error) {
	out := new(EvidenceRegistrySafeCastOverflowedUintDowncast)
	if err := evidenceRegistry.abi.UnpackIntoInterface(out, "SafeCastOverflowedUintDowncast", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EvidenceRegistryStreamsNotAscending represents a StreamsNotAscending error raised by the EvidenceRegistry contract.
type EvidenceRegistryStreamsNotAscending struct {
	Previous uint8
	Stream   uint8
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error StreamsNotAscending(uint8 previous, uint8 stream)
func EvidenceRegistryStreamsNotAscendingErrorID() common.Hash {
	return common.HexToHash("0x43f455cb7c6b84e12de47ffb24d6d2c6328f1d95892d72b2da12dc08c2ead064")
}

// UnpackStreamsNotAscendingError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error StreamsNotAscending(uint8 previous, uint8 stream)
func (evidenceRegistry *EvidenceRegistry) UnpackStreamsNotAscendingError(raw []byte) (*EvidenceRegistryStreamsNotAscending, error) {
	out := new(EvidenceRegistryStreamsNotAscending)
	if err := evidenceRegistry.abi.UnpackIntoInterface(out, "StreamsNotAscending", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EvidenceRegistryTooManyLeaves represents a TooManyLeaves error raised by the EvidenceRegistry contract.
type EvidenceRegistryTooManyLeaves struct {
	LeafCount uint64
	Maximum   uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TooManyLeaves(uint64 leafCount, uint64 maximum)
func EvidenceRegistryTooManyLeavesErrorID() common.Hash {
	return common.HexToHash("0x09e51df967feec9808a68989cf05b6cd9c8c3c92676ae3ea00eda6a12380a9af")
}

// UnpackTooManyLeavesError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TooManyLeaves(uint64 leafCount, uint64 maximum)
func (evidenceRegistry *EvidenceRegistry) UnpackTooManyLeavesError(raw []byte) (*EvidenceRegistryTooManyLeaves, error) {
	out := new(EvidenceRegistryTooManyLeaves)
	if err := evidenceRegistry.abi.UnpackIntoInterface(out, "TooManyLeaves", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EvidenceRegistryTooManyStreams represents a TooManyStreams error raised by the EvidenceRegistry contract.
type EvidenceRegistryTooManyStreams struct {
	Streams *big.Int
	Maximum *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TooManyStreams(uint256 streams, uint256 maximum)
func EvidenceRegistryTooManyStreamsErrorID() common.Hash {
	return common.HexToHash("0x7e1d37064d99a880f4dc620cfffcf10ad387324ee240ad76ca4bb0609b1ba381")
}

// UnpackTooManyStreamsError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TooManyStreams(uint256 streams, uint256 maximum)
func (evidenceRegistry *EvidenceRegistry) UnpackTooManyStreamsError(raw []byte) (*EvidenceRegistryTooManyStreams, error) {
	out := new(EvidenceRegistryTooManyStreams)
	if err := evidenceRegistry.abi.UnpackIntoInterface(out, "TooManyStreams", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EvidenceRegistryUnknownBatch represents a UnknownBatch error raised by the EvidenceRegistry contract.
type EvidenceRegistryUnknownBatch struct {
	Stream uint8
	Number uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error UnknownBatch(uint8 stream, uint64 number)
func EvidenceRegistryUnknownBatchErrorID() common.Hash {
	return common.HexToHash("0x2293f2100b7b5077c8aef5599409d2ed688c284464fee8b993b5a163196812b1")
}

// UnpackUnknownBatchError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error UnknownBatch(uint8 stream, uint64 number)
func (evidenceRegistry *EvidenceRegistry) UnpackUnknownBatchError(raw []byte) (*EvidenceRegistryUnknownBatch, error) {
	out := new(EvidenceRegistryUnknownBatch)
	if err := evidenceRegistry.abi.UnpackIntoInterface(out, "UnknownBatch", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EvidenceRegistryUnknownStream represents a UnknownStream error raised by the EvidenceRegistry contract.
type EvidenceRegistryUnknownStream struct {
	Stream uint8
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error UnknownStream(uint8 stream)
func EvidenceRegistryUnknownStreamErrorID() common.Hash {
	return common.HexToHash("0x8b7e1f972d87d002a7624f1a62232147b6c344d8b0ea791a8a21b9a8044abca8")
}

// UnpackUnknownStreamError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error UnknownStream(uint8 stream)
func (evidenceRegistry *EvidenceRegistry) UnpackUnknownStreamError(raw []byte) (*EvidenceRegistryUnknownStream, error) {
	out := new(EvidenceRegistryUnknownStream)
	if err := evidenceRegistry.abi.UnpackIntoInterface(out, "UnknownStream", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EvidenceRegistryZeroAnchorer represents a ZeroAnchorer error raised by the EvidenceRegistry contract.
type EvidenceRegistryZeroAnchorer struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ZeroAnchorer()
func EvidenceRegistryZeroAnchorerErrorID() common.Hash {
	return common.HexToHash("0x7110c4a1e207dd190eaa2e6c40aeb0e70affea5eb498857bf51f0d9d85807115")
}

// UnpackZeroAnchorerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ZeroAnchorer()
func (evidenceRegistry *EvidenceRegistry) UnpackZeroAnchorerError(raw []byte) (*EvidenceRegistryZeroAnchorer, error) {
	out := new(EvidenceRegistryZeroAnchorer)
	if err := evidenceRegistry.abi.UnpackIntoInterface(out, "ZeroAnchorer", raw); err != nil {
		return nil, err
	}
	return out, nil
}
