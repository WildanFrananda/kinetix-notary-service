// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ledgertwo

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

// LiabilityLedgerV2MetaData contains all meta data concerning the LiabilityLedgerV2 contract.
var LiabilityLedgerV2MetaData = bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"admin\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"adminTransferDelay\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"anchorer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"genesisEpoch_\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"genesisDigest_\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"genesisAsOf\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"verifier_\",\"type\":\"address\",\"internalType\":\"contractPlonkVerifier\"},{\"name\":\"depth_\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"ANCHORER_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"DEFAULT_ADMIN_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptDefaultAdminTransfer\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"anchor\",\"inputs\":[{\"name\":\"header\",\"type\":\"tuple\",\"internalType\":\"structEpochHeader\",\"components\":[{\"name\":\"number\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"previousDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"total\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"leafCount\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"asOf\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"beginDefaultAdminTransfer\",\"inputs\":[{\"name\":\"newAdmin\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"cancelDefaultAdminTransfer\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"capacity\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"changeDefaultAdminDelay\",\"inputs\":[{\"name\":\"newDelay\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"defaultAdmin\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"defaultAdminDelay\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"defaultAdminDelayIncreaseWait\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"depth\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"epochAt\",\"inputs\":[{\"name\":\"number\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structEpoch\",\"components\":[{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"digest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"total\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"leafCount\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"asOf\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"anchoredAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"genesisDigest\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"genesisEpoch\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRoleAdmin\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"grantRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"hasRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"head\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isProven\",\"inputs\":[{\"name\":\"number\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"leafCount\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isSealed\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"lastAsOf\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"nextEpoch\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingDefaultAdmin\",\"inputs\":[],\"outputs\":[{\"name\":\"newAdmin\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"schedule\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingDefaultAdminDelay\",\"inputs\":[],\"outputs\":[{\"name\":\"newDelay\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"schedule\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proveEpoch\",\"inputs\":[{\"name\":\"number\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"treeRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"totalCommitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"leafCount\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revokeRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"rollbackDefaultAdminDelay\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"scheme\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"seal\",\"inputs\":[{\"name\":\"expectedHead\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalCommitmentOf\",\"inputs\":[{\"name\":\"number\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"treeRootOf\",\"inputs\":[{\"name\":\"number\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifier\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyInclusion\",\"inputs\":[{\"name\":\"number\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"commitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"amount\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"index\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"siblings\",\"type\":\"tuple[]\",\"internalType\":\"structSumNode[]\",\"components\":[{\"name\":\"hash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sum\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"DefaultAdminDelayChangeCanceled\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DefaultAdminDelayChangeScheduled\",\"inputs\":[{\"name\":\"newDelay\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"},{\"name\":\"effectSchedule\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DefaultAdminTransferCanceled\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DefaultAdminTransferScheduled\",\"inputs\":[{\"name\":\"newAdmin\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"acceptSchedule\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EpochAnchored\",\"inputs\":[{\"name\":\"number\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"digest\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"previousDigest\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"root\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"total\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"leafCount\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"asOf\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EpochProven\",\"inputs\":[{\"name\":\"number\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"root\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"treeRoot\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"totalCommitment\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"leafCount\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"LedgerSealed\",\"inputs\":[{\"name\":\"head\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"nextEpoch\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"lastAsOf\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleAdminChanged\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"previousAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"newAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleGranted\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleRevoked\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AccessControlBadConfirmation\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AccessControlEnforcedDefaultAdminDelay\",\"inputs\":[{\"name\":\"schedule\",\"type\":\"uint48\",\"internalType\":\"uint48\"}]},{\"type\":\"error\",\"name\":\"AccessControlEnforcedDefaultAdminRules\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AccessControlInvalidDefaultAdmin\",\"inputs\":[{\"name\":\"defaultAdmin\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"AccessControlUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"neededRole\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"AnchorerNotGranted\",\"inputs\":[{\"name\":\"anchorer\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"DepthOutOfRange\",\"inputs\":[{\"name\":\"depth\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"maximum\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"EmptyEpochMismatch\",\"inputs\":[{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"total\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"EpochAlreadyAnchored\",\"inputs\":[{\"name\":\"number\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"digest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"EpochBeforeGenesis\",\"inputs\":[{\"name\":\"genesisEpoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"number\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"EpochClosed\",\"inputs\":[{\"name\":\"number\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"nextEpoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"EpochConflict\",\"inputs\":[{\"name\":\"number\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"anchored\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"given\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"EpochNotProven\",\"inputs\":[{\"name\":\"number\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"leafCount\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"EpochOutOfSequence\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"given\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"HeadMismatch\",\"inputs\":[{\"name\":\"head\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"InvalidGenesis\",\"inputs\":[{\"name\":\"genesisEpoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"genesisDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"LedgerIsSealed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotAFieldElement\",\"inputs\":[{\"name\":\"value\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"NotAVerifier\",\"inputs\":[{\"name\":\"verifier\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"PreviousDigestMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"given\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ProofRejected\",\"inputs\":[{\"name\":\"number\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"SafeCastOverflowedUintDowncast\",\"inputs\":[{\"name\":\"bits\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"SnapshotInFuture\",\"inputs\":[{\"name\":\"asOf\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"blockTime\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"SnapshotNotAfterPrevious\",\"inputs\":[{\"name\":\"asOf\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"previous\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"StatementAlreadyProven\",\"inputs\":[{\"name\":\"number\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"leafCount\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"TooManyLeaves\",\"inputs\":[{\"name\":\"leafCount\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"maximum\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"TotalNotHidden\",\"inputs\":[{\"name\":\"total\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"UnknownEpoch\",\"inputs\":[{\"name\":\"number\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"ZeroAnchorer\",\"inputs\":[]}]",
	ID:  "LiabilityLedgerV2",
	Bin: "0x610100604052348015610010575f5ffd5b506040516141f43803806141f483398101604081905261002f9161034c565b86886001600160a01b03811661005f57604051636116401160e11b81525f60048201526024015b60405180910390fd5b600180546001600160d01b0316600160d01b65ffffffffffff8516021790556100885f82610204565b5050506001600160401b038516158415146100c85760405163111f019f60e01b81526001600160401b038616600482015260248101859052604401610056565b6001600160a01b0386166100ef57604051637110c4a160e01b815260040160405180910390fd5b816001600160a01b03163b5f0361012457604051631606811f60e21b81526001600160a01b0383166004820152602401610056565b601060ff8216111561015557604051630d34bba760e31b815260ff8216600482015260106024820152604401610056565b61017f7f4a8e1518fc3b69d274d0d6c14a98dd84f4bfb1e39bbbe43f417772dd939cd80887610204565b6101a757604051633012fc8b60e01b81526001600160a01b0387166004820152602401610056565b6001600160401b03948516608081905260a08590526001600160a01b0390921660c05260ff1660e052600692909255600780549190931668010000000000000000026001600160801b0319909116909117179055506103f6915050565b5f82610260575f61021d6002546001600160a01b031690565b6001600160a01b03161461024457604051631fe1e13d60e11b815260040160405180910390fd5b600280546001600160a01b0319166001600160a01b0384161790555b61026a8383610273565b90505b92915050565b5f828152602081815260408083206001600160a01b038516845290915281205460ff16610313575f838152602081815260408083206001600160a01b03861684529091529020805460ff191660011790556102cb3390565b6001600160a01b0316826001600160a01b0316847f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a450600161026d565b505f61026d565b6001600160a01b038116811461032e575f5ffd5b50565b80516001600160401b0381168114610347575f5ffd5b919050565b5f5f5f5f5f5f5f5f610100898b031215610364575f5ffd5b885161036f8161031a565b60208a015190985065ffffffffffff8116811461038a575f5ffd5b60408a015190975061039b8161031a565b95506103a960608a01610331565b9450608089015193506103be60a08a01610331565b60c08a01519093506103cf8161031a565b60e08a015190925060ff811681146103e5575f5ffd5b809150509295985092959890939650565b60805160a05160c05160e051613d8561046f5f395f8181610395015281816103d7015281816112370152818161127a0152818161194201528181611c160152611c7701525f81816103000152611f9101525f61045501525f818161061a015281816109230152818161096201526120390152613d855ff3fe608060405234801561000f575f5ffd5b5060043610610269575f3560e01c806384ef8ffc11610157578063b07eeda8116100d2578063cc8463c811610088578063cf6eefb71161006e578063cf6eefb714610685578063d547741f146106d1578063d602b9fd146106e4575f5ffd5b8063cc8463c814610675578063cefc14291461067d575f5ffd5b8063c01e5dd8116100b8578063c01e5dd81461063c578063c0566abe1461064f578063c9c6958914610662575f5ffd5b8063b07eeda814610602578063b70e6be614610615575f5ffd5b806391d1485411610127578063a217fddf1161010d578063a217fddf146105cb578063a89210cb146105d2578063aea0e78b146105f0575f5ffd5b806391d1485414610561578063a1eda53c146105a4575f5ffd5b806384ef8ffc146105245780638518b0ad146105425780638da5cb5b146105515780638f7dcfa314610559575f5ffd5b806357db9232116101e7578063634e93da116101b757806379ac26611161019d57806379ac2661146104505780637cfd9a51146104775780637d60d3551461049e575f5ffd5b8063634e93da1461042a578063649a5ec71461043d575f5ffd5b806357db92321461037e5780635cfc1a5114610391578063631c56ef146103d2578063631f98521461040b575f5ffd5b8063248a9ca31161023c5780632f2ff15d116102225780632f2ff15d1461034557806336568abe146103585780633dad095b1461036b575f5ffd5b8063248a9ca3146102dc5780632b7ac3f3146102fe575f5ffd5b806301ffc9a71461026d578063022d63fb146102955780630aa6220b146102b15780631ffcf5d4146102bb575b5f5ffd5b61028061027b3660046138c5565b6106ec565b60405190151581526020015b60405180910390f35b620697805b60405165ffffffffffff909116815260200161028c565b6102b9610793565b005b6102ce6102c9366004613904565b6107a8565b60405190815260200161028c565b6102ce6102ea36600461391d565b5f9081526020819052604090206001015490565b7f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200161028c565b6102b961035336600461395c565b611034565b6102b961036636600461395c565b611079565b6102ce61037936600461399d565b61117e565b6102ce61038c366004613a3e565b6115a7565b60017f000000000000000000000000000000000000000000000000000000000000000060ff161b5b60405167ffffffffffffffff909116815260200161028c565b6103f97f000000000000000000000000000000000000000000000000000000000000000081565b60405160ff909116815260200161028c565b600754700100000000000000000000000000000000900460ff16610280565b6102b9610438366004613a57565b6115d0565b6102b961044b366004613a70565b6115e3565b6102ce7f000000000000000000000000000000000000000000000000000000000000000081565b6102ce7f4a8e1518fc3b69d274d0d6c14a98dd84f4bfb1e39bbbe43f417772dd939cd80881565b6104b16104ac366004613a3e565b6115f6565b60405161028c91905f60c082019050825182526020830151602083015267ffffffffffffffff604084015116604083015267ffffffffffffffff606084015116606083015267ffffffffffffffff608084015116608083015267ffffffffffffffff60a08401511660a083015292915050565b60025473ffffffffffffffffffffffffffffffffffffffff16610320565b6040516002815260200161028c565b6103206116b7565b6006546102ce565b61028061056f36600461395c565b5f9182526020828152604080842073ffffffffffffffffffffffffffffffffffffffff93909316845291905290205460ff1690565b6105ac6116dc565b6040805165ffffffffffff93841681529290911660208301520161028c565b6102ce5f81565b60075468010000000000000000900467ffffffffffffffff166103b9565b60075467ffffffffffffffff166103b9565b6102b961061036600461391d565b611756565b6103b97f000000000000000000000000000000000000000000000000000000000000000081565b61028061064a366004613a95565b611893565b61028061065d366004613ace565b611902565b6102ce610670366004613a3e565b6119c7565b61029a6119ed565b6102b9611a8a565b6001546040805173ffffffffffffffffffffffffffffffffffffffff831681527401000000000000000000000000000000000000000090920465ffffffffffff1660208301520161028c565b6102b96106df36600461395c565b611ae6565b6102b9611b27565b5f7fffffffff0000000000000000000000000000000000000000000000000000000082167fb35efc0c00000000000000000000000000000000000000000000000000000000148061077e57507fffffffff0000000000000000000000000000000000000000000000000000000082167f7734dc7500000000000000000000000000000000000000000000000000000000145b8061078d575061078d82611b39565b92915050565b5f61079d81611b8e565b6107a5611b98565b50565b5f7f4a8e1518fc3b69d274d0d6c14a98dd84f4bfb1e39bbbe43f417772dd939cd8086107d381611b8e565b600754700100000000000000000000000000000000900460ff1615610824576040517f712a4e2f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f61091661083736869003860186613b64565b80516020808301516040808501516060860151608087015160a09097015183517f0300000000000000000000000000000000000000000000000000000000000000818801527e0200000000000000000000000000000000000000000000000000000000000060218201527fffffffffffffffff00000000000000000000000000000000000000000000000060c098891b81166023830152602b820196909652604b81019390935290861b8416606b83015295851b831660738201529490931b16607b840152815180840360630181526083909301909152815191012090565b905067ffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000166109506020860186613a3e565b67ffffffffffffffff1610156109d9577f000000000000000000000000000000000000000000000000000000000000000061098e6020860186613a3e565b6040517facc3e9cb00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff9283166004820152911660248201526044015b60405180910390fd5b60075467ffffffffffffffff166109f36020860186613a3e565b67ffffffffffffffff161015610aef575f600381610a146020880188613a3e565b67ffffffffffffffff1667ffffffffffffffff1681526020019081526020015f20600101549050818103610a9557610a4f6020860186613a3e565b6040517fab455a2500000000000000000000000000000000000000000000000000000000815267ffffffffffffffff9091166004820152602481018390526044016109d0565b610aa26020860186613a3e565b6040517f6d613c3900000000000000000000000000000000000000000000000000000000815267ffffffffffffffff909116600482015260248101829052604481018390526064016109d0565b60075467ffffffffffffffff16610b096020860186613a3e565b67ffffffffffffffff161115610b795760075467ffffffffffffffff16610b336020860186613a3e565b6040517f7a8d061d00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff9283166004820152911660248201526044016109d0565b600654846020013514610bc9576006546040517f0796b55a0000000000000000000000000000000000000000000000000000000081526004810191909152602085013560248201526044016109d0565b610bd284611ba4565b5f600581610c42610be66020890189613a3e565b6040890135610bfb60a08b0160808c01613a3e565b6040805167ffffffffffffffff808616602083015291810184905290821660608201525f906080016040516020818303038152906040528051906020012090509392505050565b815260208082019290925260409081015f208151606081018352815481526001820154938101939093526002015460ff161515908201819052909150610cf357610c8f6020860186613a3e565b6040860135610ca460a0880160808901613a3e565b6040517f72deac1800000000000000000000000000000000000000000000000000000000815267ffffffffffffffff9384166004820152602481019290925290911660448201526064016109d0565b6040518060c00160405280866040013581526020018381526020015f67ffffffffffffffff168152602001866080016020810190610d319190613a3e565b67ffffffffffffffff168152602001610d5060c0880160a08901613a3e565b67ffffffffffffffff168152602001610d6842611dbe565b67ffffffffffffffff16905260035f610d846020890189613a3e565b67ffffffffffffffff908116825260208083019390935260409182015f9081208551815585850151600182015592850151600290930180546060870151608088015160a090980151851678010000000000000000000000000000000000000000000000000277ffffffffffffffffffffffffffffffffffffffffffffffff98861670010000000000000000000000000000000002989098166fffffffffffffffffffffffffffffffff91861668010000000000000000027fffffffffffffffffffffffffffffffff0000000000000000000000000000000090931696909516959095171793909316919091179390931790558291600491610e8790890189613a3e565b67ffffffffffffffff16815260208082019290925260409081015f2083518155838301516001820155920151600290920180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016921515929092179091556006839055610ef790860186613a3e565b610f02906001613c47565b600780547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001667ffffffffffffffff92909216919091179055610f4b60c0860160a08701613a3e565b6007805467ffffffffffffffff9290921668010000000000000000027fffffffffffffffffffffffffffffffff0000000000000000ffffffffffffffff90921691909117905581610f9f6020870187613a3e565b67ffffffffffffffff167fb7504cd1e32568b16a96dd8fb0da74c649e425006318bb259eca087e629c542d602088013560408901355f610fe560a08c0160808d01613a3e565b610ff560c08d0160a08e01613a3e565b60408051958652602086019490945267ffffffffffffffff928316858501529082166060850152166080830152519081900360a00190a3509392505050565b8161106b576040517f3fc3c27a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6110758282611e11565b5050565b811580156110a1575060025473ffffffffffffffffffffffffffffffffffffffff8281169116145b156111745760015473ffffffffffffffffffffffffffffffffffffffff81169074010000000000000000000000000000000000000000900465ffffffffffff16811515806110f5575065ffffffffffff8116155b8061110857504265ffffffffffff821610155b15611149576040517f19ca5ebb00000000000000000000000000000000000000000000000000000000815265ffffffffffff821660048201526024016109d0565b5050600180547fffffffffffff000000000000ffffffffffffffffffffffffffffffffffffffff1690555b6110758282611e3b565b6007545f90700100000000000000000000000000000000900460ff16156111d1576040517f712a4e2f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60075467ffffffffffffffff9081169088161015611233576007546040517f2e5386d600000000000000000000000000000000000000000000000000000000815267ffffffffffffffff808a16600483015290911660248201526044016109d0565b60017f000000000000000000000000000000000000000000000000000000000000000060ff161b67ffffffffffffffff168467ffffffffffffffff1611156112e4578360017f000000000000000000000000000000000000000000000000000000000000000060ff161b5b6040517f09e51df900000000000000000000000000000000000000000000000000000000815267ffffffffffffffff9283166004820152911660248201526044016109d0565b7f30644e72e131a029b85045b68181585d2833e84879b9709143e1f593f00000018610611340576040517f1c35841b000000000000000000000000000000000000000000000000000000008152600481018790526024016109d0565b7f30644e72e131a029b85045b68181585d2833e84879b9709143e1f593f0000001851061139c576040517f1c35841b000000000000000000000000000000000000000000000000000000008152600481018690526024016109d0565b604080517f04000000000000000000000000000000000000000000000000000000000000006020808301919091526021820189905260418083018990528351808403909101815260618301845280519082012067ffffffffffffffff808c16608185015260a18401829052881660c1808501919091528451808503909101815260e190930190935281519101205f905f8181526005602052604090206002015490915060ff1615611494576040517f1cf1dcbb00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff808b16600483015260248201849052871660448201526064016109d0565b6114a285858b8b8b8b611e99565b6114e4576040517fea282d2400000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8a1660048201526024016109d0565b60408051606080820183528a825260208083018b815260018486018181525f888152600585528790209551865591519085015551600290930180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169315159390931790925582518b81529182018a905267ffffffffffffffff8981169383019390935284928c16917f272cc6708feee2c8c722412bc22031d0a3a4097c2a2315bb331985b256f3a586910160405180910390a35090505b9695505050505050565b5f6115b182612037565b5067ffffffffffffffff165f9081526004602052604090206001015490565b5f6115da81611b8e565b611075826120cc565b5f6115ed81611b8e565b61107582612144565b6040805160c0810182525f80825260208201819052918101829052606081018290526080810182905260a0810191909152611630826121b3565b6040805160c081018252825481526001830154602082015260029092015467ffffffffffffffff8082169284019290925268010000000000000000810482166060840152700100000000000000000000000000000000810482166080840152780100000000000000000000000000000000000000000000000090041660a082015292915050565b5f6116d760025473ffffffffffffffffffffffffffffffffffffffff1690565b905090565b6002545f907a010000000000000000000000000000000000000000000000000000900465ffffffffffff16801515801561171e57504265ffffffffffff821610155b611729575f5f61174e565b60025474010000000000000000000000000000000000000000900465ffffffffffff16815b915091509091565b5f61176081611b8e565b600754700100000000000000000000000000000000900460ff16156117b1576040517f712a4e2f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60065482146117fa576006546040517f6f2de3c70000000000000000000000000000000000000000000000000000000081526004810191909152602481018390526044016109d0565b600780547001000000000000000000000000000000007fffffffffffffffffffffffffffffff00ffffffffffffffffffffffffffffffff821617918290556006546040805167ffffffffffffffff9384168152680100000000000000009094049092166020840152917f79b7f4b971d94514f2584fdbd244ea1c7fbac6672ab3e05044cb71c34f9752d991015b60405180910390a25050565b5f60055f6118e48686866040805167ffffffffffffffff808616602083015291810184905290821660608201525f906080016040516020818303038152906040528051906020012090509392505050565b815260208101919091526040015f206002015460ff16949350505050565b5f5f61190d886121b3565b600281015490915067ffffffffffffffff6801000000000000000090910481169086161061193e575f91505061159d565b60ff7f0000000000000000000000000000000000000000000000000000000000000000168314611971575f91505061159d565b61197b84846121d8565b611988575f91505061159d565b5f61199f6119978a8a8a612280565b8787876122f6565b67ffffffffffffffff8a165f9081526004602052604090205414925050509695505050505050565b5f6119d182612037565b5067ffffffffffffffff165f9081526004602052604090205490565b6002545f907a010000000000000000000000000000000000000000000000000000900465ffffffffffff168015158015611a2e57504265ffffffffffff8216105b611a60576001547a010000000000000000000000000000000000000000000000000000900465ffffffffffff16611a84565b60025474010000000000000000000000000000000000000000900465ffffffffffff165b91505090565b60015473ffffffffffffffffffffffffffffffffffffffff16338114611ade576040517fc22c80220000000000000000000000000000000000000000000000000000000081523360048201526024016109d0565b6107a561236c565b81611b1d576040517f3fc3c27a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b611075828261245d565b5f611b3181611b8e565b6107a5612481565b5f7fffffffff0000000000000000000000000000000000000000000000000000000082167f3149878600000000000000000000000000000000000000000000000000000000148061078d575061078d8261248b565b6107a58133612521565b611ba25f5f6125a6565b565b611bb46080820160608301613a3e565b67ffffffffffffffff1615611c1257611bd36080820160608301613a3e565b6040517f81d8d68b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff90911660048201526024016109d0565b60017f000000000000000000000000000000000000000000000000000000000000000060ff161b67ffffffffffffffff16611c5360a0830160808401613a3e565b67ffffffffffffffff161115611c9f57611c7360a0820160808301613a3e565b60017f000000000000000000000000000000000000000000000000000000000000000060ff161b61129e565b60075468010000000000000000900467ffffffffffffffff16611cc860c0830160a08401613a3e565b67ffffffffffffffff1611611d3e57611ce760c0820160a08301613a3e565b6007546040517f81ed60e900000000000000000000000000000000000000000000000000000000815267ffffffffffffffff92831660048201526801000000000000000090910490911660248201526044016109d0565b42611d4f60c0830160a08401613a3e565b67ffffffffffffffff1611156107a557611d6f60c0820160a08301613a3e565b611d7842611dbe565b6040517f18aaced200000000000000000000000000000000000000000000000000000000815267ffffffffffffffff9283166004820152911660248201526044016109d0565b5f67ffffffffffffffff821115611e0d57604080517f6dfcc6500000000000000000000000000000000000000000000000000000000081526004810191909152602481018390526044016109d0565b5090565b5f82815260208190526040902060010154611e2b81611b8e565b611e3583836126ff565b50505050565b73ffffffffffffffffffffffffffffffffffffffff81163314611e8a576040517f6697b23200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b611e9482826127c4565b505050565b60408051600480825260a082019092525f91829190602082016080803683370190505090508567ffffffffffffffff16815f81518110611edb57611edb613c67565b602002602001018181525050845f1c81600181518110611efd57611efd613c67565b602002602001018181525050835f1c81600281518110611f1f57611f1f613c67565b6020026020010181815250508267ffffffffffffffff1681600381518110611f4957611f49613c67565b60209081029190910101526040517f7e4f7a8a00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001690637e4f7a8a90611fca908b908b908690600401613c94565b602060405180830381865afa925050508015612021575060408051601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820190925261201e91810190613d2a565b60015b61202e575f91505061159d565b915061159d9050565b7f000000000000000000000000000000000000000000000000000000000000000067ffffffffffffffff168167ffffffffffffffff161080612089575060075467ffffffffffffffff90811690821610155b156107a5576040517fa4e21e6c00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff821660048201526024016109d0565b5f6120d56119ed565b6120de42612825565b6120e89190613d49565b90506120f48282612870565b60405165ffffffffffff8216815273ffffffffffffffffffffffffffffffffffffffff8316907f3377dc44241e779dd06afab5b788a35ca5f3b778836e2990bdb26a2a4b2e5ed690602001611887565b5f61214e8261290b565b61215742612825565b6121619190613d49565b905061216d82826125a6565b6040805165ffffffffffff8085168252831660208201527ff1038c18cf84a56e432fdbfaf746924b7ea511dfe03a6506a0ceba4888788d9b910160405180910390a15050565b5f6121bd82612037565b5067ffffffffffffffff165f90815260036020526040902090565b5f805b82811015612276578383828181106121f5576121f5613c67565b905060400201602001602081019061220d9190613a3e565b67ffffffffffffffff1615158061226057507f30644e72e131a029b85045b68181585d2833e84879b9709143e1f593f000000184848381811061225257612252613c67565b9050604002015f01355f1c10155b1561226e575f91505061078d565b6001016121db565b5060019392505050565b5f6fffffffffffffffff0000000000000000604085901b1667ffffffffffffffff83161770010000000000000000000000000000000017608084901c6fffffffffffffffffffffffffffffffff85166122eb6122e56122df8686612952565b84612952565b82612952565b979650505050505050565b5f8484825b84811015612360575f86868381811061231657612316613c67565b6040029190910135915050600183165f0361233c5761233584826134d6565b9350612349565b61234681856134d6565b93505b50677fffffffffffffff600192831c1691016122fb565b50909695505050505050565b60015473ffffffffffffffffffffffffffffffffffffffff81169074010000000000000000000000000000000000000000900465ffffffffffff168015806123bc57504265ffffffffffff821610155b156123fd576040517f19ca5ebb00000000000000000000000000000000000000000000000000000000815265ffffffffffff821660048201526024016109d0565b6124255f61242060025473ffffffffffffffffffffffffffffffffffffffff1690565b6127c4565b506124305f836126ff565b5050600180547fffffffffffff000000000000000000000000000000000000000000000000000016905550565b5f8281526020819052604090206001015461247781611b8e565b611e3583836127c4565b611ba25f5f612870565b5f7fffffffff0000000000000000000000000000000000000000000000000000000082167f7965db0b00000000000000000000000000000000000000000000000000000000148061078d57507f01ffc9a7000000000000000000000000000000000000000000000000000000007fffffffff0000000000000000000000000000000000000000000000000000000083161461078d565b5f8281526020818152604080832073ffffffffffffffffffffffffffffffffffffffff8516845290915290205460ff16611075576040517fe2517d3f00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff82166004820152602481018390526044016109d0565b6002547a010000000000000000000000000000000000000000000000000000900465ffffffffffff16801561267a574265ffffffffffff82161015612651576002546001805479ffffffffffffffffffffffffffffffffffffffffffffffffffff167401000000000000000000000000000000000000000090920465ffffffffffff167a0100000000000000000000000000000000000000000000000000000291909117905561267a565b6040517f2b1fa2edafe6f7b9e97c1a9e0c3660e645beb2dcaa2d45bdbf9beaf5472e1ec5905f90a15b506002805473ffffffffffffffffffffffffffffffffffffffff167401000000000000000000000000000000000000000065ffffffffffff9485160279ffffffffffffffffffffffffffffffffffffffffffffffffffff16177a0100000000000000000000000000000000000000000000000000009290931691909102919091179055565b5f826127b3575f61272560025473ffffffffffffffffffffffffffffffffffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff1614612772576040517f3fc3c27a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600280547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff84161790555b6127bd83836134e1565b9392505050565b5f821580156127ed575060025473ffffffffffffffffffffffffffffffffffffffff8381169116145b1561281b57600280547fffffffffffffffffffffffff00000000000000000000000000000000000000001690555b6127bd83836135da565b5f65ffffffffffff821115611e0d576040517f6dfcc65000000000000000000000000000000000000000000000000000000000815260306004820152602481018390526044016109d0565b600180547401000000000000000000000000000000000000000065ffffffffffff84811682027fffffffffffff0000000000000000000000000000000000000000000000000000841673ffffffffffffffffffffffffffffffffffffffff881617179093559004168015611e94576040517f8886ebfc4259abdbc16601dd8fb5678e54878f47b3c34836cfc51154a9605109905f90a1505050565b5f5f6129156119ed565b90508065ffffffffffff168365ffffffffffff161161293d576129388382613d67565b6127bd565b6127bd65ffffffffffff841662069780613693565b5f828261295f82826136a2565b90925090506129b082827f1da4d6adfb0d0b494584f763db50a81908580a5f5e295e168b9b8d31770fac4f7f0946129a2e33b4e819707a56a3b3790eab80d0a0c7a0c63451fab2a59acc074c61371f565b9092509050612a0182827f2a39b9d5376afd35580abd6952986570867b87303b07140268794dad4b8f82ea7f27605717d1245c20c546b3c7014e5fa3e4a70e66ecd6de5d32448eea5e3069b261371f565b9092509050612a5282827f24c896cb2594e17b973653193a470ed7e771ebb09a228d638e50809c037736327f0911096c45dd9cda0d61d783957003db6d8701c7869a919ad8042e1b7d597a4961371f565b9092509050612a8282827f26ff6166f1e4b99e27eee6d352a8ce26f1daba01aad28021f11f38a70603afdc61378c565b9092509050612ab182827e8e2faedcf76d08ad6591ff90e50fea79bcb5e18cfb7d954d5b89437bf64b7e61378c565b9092509050612ae182827f19c9da2379b598ace3ad4d1872470830f6184a3cec71eeb632f98820eb353d7861378c565b9092509050612b1182827f0f7c4eb15d8b0b62a8f6a090ec9610a2ab3dfcdb57e2539aa459a40583dfe96b61378c565b9092509050612b4182827f18b99417dc26b5e079750eba282362d1d46900b47dd5ff809171588b08ac398361378c565b9092509050612b7182827f1ee044081160b3eee2d4493feab82141c73f1c054b76320a8848af08a8d91a2661378c565b9092509050612ba182827f29bb95c8763efd3e0e87f5df12ee8a150455b6d7a14780d19122220366c258dc61378c565b9092509050612bd182827f22c23eec9cb13ff8a3ee9a363d740653215e8991f7f9ec12067b4705e9a5c9fb61378c565b9092509050612c0182827f23589e033a31a667680c8b18926c3be09115c7644c4f905cc7deefc1690b42dc61378c565b9092509050612c3182827f304e99b887f2e1e92c9c0cde5f2bdd4764f60b98a219f1f0dd64ec938a6a247c61378c565b9092509050612c6182827f22e817865236ad3a76fbe88bbdf31fcae792f326271d53a3a493b2de7f7d8b4c61378c565b9092509050612c9182827f10c9efe573e86fa5b238a3f5c70a00bf02bc7dcfcc4878e026b77c881dc8b1c961378c565b9092509050612cc182827f0a94f16be920d85f4d6f80e745c6bddbc084f309edffef023ea283c72b89bbb661378c565b9092509050612cf182827f23ed72b4d01d14e3c7888fedb770494abc2c1ea81d6465b5df3da0ebcd69101a61378c565b9092509050612d2182827f17c5115640e4cebeed0e6cbb511ac38498df815a77cb162de8d8f1022eb6bb7461378c565b9092509050612d5182827f2e507fcca290d0d9cf765245750eb04773e09e1fc813959fb99680a82772e4fc61378c565b9092509050612d8182827f0d4a98999f5b39176af6cce65c8d12162433f055cb71d70dfc8803013292bbbf61378c565b9092509050612db182827f238d8022cc09c21ab3c01261a03dc125321d829e22a7a3b7a1bd3c335eccfa2161378c565b9092509050612de182827f010cd8e4c2b7051cb81dc8292e35d8e943ce365613d5b39770e454ed9f4ae16561378c565b9092509050612e1182827f088027e54f2a3604b11178cf0ea3c6aa861173a90fb231238e03539590ecc02761378c565b9092509050612e4182827f1b840f5311a2b1d4b4cd7aa7e5a9a6d161165468daa938108302b73e464819db61378c565b9092509050612e7182827f2bf51a5da1828a1cf9b764b1e16c15929a3a346e44198ea0cb366fcd8db78dc161378c565b9092509050612ea182827f206ad089d8d296ffe68a6a86767a7fe726b8872f9c7beef9d56a3a50f6f2382761378c565b9092509050612ed182827f24d19193171494fa1a54e0a99ac16d05eaec4b6d617c7c199fc07ff88eac550c61378c565b9092509050612f0182827f1dd654a2ca9d9f24f33d88246a40dfb32c40662278b9c0b995d9f9fbaf15213861378c565b9092509050612f3182827f0d171025c925f6e259d20ecbd3a601108c85f05b0fe7793b8acf57f3789785e461378c565b9092509050612f6182827f055bef435a43aec245cd99ccb0f7c8791d9e8cf2b80d95dd98b9908fed877d5561378c565b9092509050612f9182827f10d2ac8c61c8a2e88a2a3f42d9359a14be63d0ad4cfd9f203b75765d0e595f0e61378c565b9092509050612fc182827f103479710e70996982a84057ec3ba6b2254d7966ddc41e85296e3d0790dcfa5661378c565b9092509050612ff182827f2a366f0448fda3c05914ffb48c765da8de96f9aa340db1638225f8372921807b61378c565b909250905061302182827f16be0fb8ef62da17919b6e0378d00594153bb8899aeb686c04626b63285837a461378c565b909250905061305182827f0417038500e9d06c60abbc7f0d0d24c32dec8a0b2aa5a4d52cfd8c78a15bc37061378c565b909250905061308182827f26a6873b43ffd2ccf66ec6f4493ff9b54f4d76480bc486a3e5a0308fdd01381261378c565b90925090506130b182827f0a3314a838f32630a96251914fe5ad262f3db9b2aa8aa9f7391922d36278c49861378c565b90925090506130e182827f0fde0c5429a6beb07f462d4821f48f86aeadb46a090b15a044f4b59399027da461378c565b909250905061311182827f0abc2d5049972a6b9b357e4163793b0bb517e1eb535a984df11a1c89cda2c8a961378c565b909250905061314182827f0dab51d6e3ebfa661d21722fb21e55051b427a5f173f7f17735205dbb77c464e61378c565b909250905061317182827f29c36622598b511d51af6cc37123652fb21be5c2d68fb8efe9b92e70a8c1ae0361378c565b90925090506131a182827f2c03ec80adac2a33ae846bc0b700d0bcc41c4096e53ac6990d6bbe7ea2fbc85c61378c565b90925090506131d182827f0918fdbe9cf3a59fbdb4c6852d065105317303116017d893b8b521e3cebe1e0d61378c565b909250905061320182827f1f19ec22e69ca33f599dd13cd7e495a8176a79df4f7acf79a9d2135acabe235961378c565b909250905061323182827f1c4b037c8ae85ee1eb32b872eb7f76928c4c76b29ceb001346447b291108070461378c565b909250905061326182827f2b68900ed906616d6c826d0bde341766ba3131e04d420de5af0a69c534efd8db61378c565b909250905061329182827f20ca92aa222fcc69448f8dac653c8daaa180ff6dfb397cef723d4f0c782bc7f061378c565b90925090506132c182827f10d22d05bdff6bb375276fc82057db337045a5ab7ac053941f6186289b66b2b661378c565b90925090506132f182827f0b1ffdbb529367bb98f32ba45784cb435aa910b4a00636d1e5ca79e88bdd6cd961378c565b909250905061332182827f2da32b38e7984bc2ed757ec705eccf8282c7b4f82e5e515f6f89bcc33022ce9f61378c565b909250905061335182827f042593ad87403f6d2674b8b55a886725b87eb33958031e94d292cecc6abed1bb61378c565b909250905061338182827f181fa1b4d067783a19d7367bf49b3f01051faedab463a6de9308fbd6e7d419f161378c565b90925090506133b182827f15aaa6cc9b7900b15683c95515c26028a8e35b00ed8a815c34927188c60de66061378c565b909250905061340282827f1bf28a93209084bbbc63234f057254c280b1a636f5a0eced6787320212f75a7a7f1cdb8c8bee5426f02cd9e229776118f156273b312f805c8e6f8c9d81a620cb6f61371f565b909250905061345382827f08299c0abf196d53162e0facb5f1876f516df2505cc387a0f8ea0e8760d5ca7e7f221643d205fe82778a7b7b58cb65c4962d76c0072cabd1124117269d7c710b8a61371f565b90925090506134a482827f2d036a95f81cf49bb7a0143a28c88767f6bd10c3f74b22db487920d43343dbff7f08a50897c06aafe6ea414fb1bceca2267cd4a39486729fbc6d5d1bb7a172ebd261371f565b90925090507f30644e72e131a029b85045b68181585d2833e84879b9709143e1f593f000000184820895945050505050565b5f6127bd8383612952565b5f8281526020818152604080832073ffffffffffffffffffffffffffffffffffffffff8516845290915281205460ff166135d3575f8381526020818152604080832073ffffffffffffffffffffffffffffffffffffffff86168452909152902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660011790556135713390565b73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16847f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a450600161078d565b505f61078d565b5f8281526020818152604080832073ffffffffffffffffffffffffffffffffffffffff8516845290915281205460ff16156135d3575f8381526020818152604080832073ffffffffffffffffffffffffffffffffffffffff8616808552925280832080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016905551339286917ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9190a450600161078d565b5f8282188284100282186127bd565b5f80807f30644e72e131a029b85045b68181585d2833e84879b9709143e1f593f000000184860890507f30644e72e131a029b85045b68181585d2833e84879b9709143e1f593f00000018186087f30644e72e131a029b85045b68181585d2833e84879b9709143e1f593f000000182860892509250509250929050565b5f8061377f6137517f30644e72e131a029b85045b68181585d2833e84879b9709143e1f593f00000015b8689086137cf565b61377a7f30644e72e131a029b85045b68181585d2833e84879b9709143e1f593f0000001613749565b6136a2565b9150915094509492505050565b5f806137c36137bd7f30644e72e131a029b85045b68181585d2833e84879b9709143e1f593f00000018588086137cf565b85613825565b91509150935093915050565b5f807f30644e72e131a029b85045b68181585d2833e84879b9709143e1f593f000000183840990507f30644e72e131a029b85045b68181585d2833e84879b9709143e1f593f00000018082830984099392505050565b5f80807f30644e72e131a029b85045b68181585d2833e84879b9709143e1f593f000000184860890507f30644e72e131a029b85045b68181585d2833e84879b9709143e1f593f00000018186087f30644e72e131a029b85045b68181585d2833e84879b9709143e1f593f0000001827f30644e72e131a029b85045b68181585d2833e84879b9709143e1f593f00000018788080892509250509250929050565b5f602082840312156138d5575f5ffd5b81357fffffffff00000000000000000000000000000000000000000000000000000000811681146127bd575f5ffd5b5f60c0828403128015613915575f5ffd5b509092915050565b5f6020828403121561392d575f5ffd5b5035919050565b803573ffffffffffffffffffffffffffffffffffffffff81168114613957575f5ffd5b919050565b5f5f6040838503121561396d575f5ffd5b8235915061397d60208401613934565b90509250929050565b803567ffffffffffffffff81168114613957575f5ffd5b5f5f5f5f5f5f60a087890312156139b2575f5ffd5b6139bb87613986565b955060208701359450604087013593506139d760608801613986565b9250608087013567ffffffffffffffff8111156139f2575f5ffd5b8701601f81018913613a02575f5ffd5b803567ffffffffffffffff811115613a18575f5ffd5b896020828401011115613a29575f5ffd5b60208201935080925050509295509295509295565b5f60208284031215613a4e575f5ffd5b6127bd82613986565b5f60208284031215613a67575f5ffd5b6127bd82613934565b5f60208284031215613a80575f5ffd5b813565ffffffffffff811681146127bd575f5ffd5b5f5f5f60608486031215613aa7575f5ffd5b613ab084613986565b925060208401359150613ac560408501613986565b90509250925092565b5f5f5f5f5f5f60a08789031215613ae3575f5ffd5b613aec87613986565b955060208701359450613b0160408801613986565b9350613b0f60608801613986565b9250608087013567ffffffffffffffff811115613b2a575f5ffd5b8701601f81018913613b3a575f5ffd5b803567ffffffffffffffff811115613b50575f5ffd5b8960208260061b8401011115613a29575f5ffd5b5f60c0828403128015613b75575f5ffd5b5060405160c0810167ffffffffffffffff81118282101715613bbe577f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b604052613bca83613986565b81526020838101359082015260408084013590820152613bec60608401613986565b6060820152613bfd60808401613986565b6080820152613c0e60a08401613986565b60a08201529392505050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b67ffffffffffffffff818116838216019081111561078d5761078d613c1a565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b60408152826040820152828460608301375f606084830101525f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f85011682016060810160608483030160208501528085518083526080840191506020870193505f92505b80831015613d1e5783518252602082019150602084019350600183019250613cfb565b50979650505050505050565b5f60208284031215613d3a575f5ffd5b815180151581146127bd575f5ffd5b65ffffffffffff818116838216019081111561078d5761078d613c1a565b65ffffffffffff828116828216039081111561078d5761078d613c1a56",
}

// LiabilityLedgerV2 is an auto generated Go binding around an Ethereum contract.
type LiabilityLedgerV2 struct {
	abi abi.ABI
}

// GetABI returns the ABI associated with this contract binding.
func (c *LiabilityLedgerV2) GetABI() abi.ABI {
	return c.abi
}

// NewLiabilityLedgerV2 creates a new instance of LiabilityLedgerV2.
func NewLiabilityLedgerV2() *LiabilityLedgerV2 {
	parsed, err := LiabilityLedgerV2MetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &LiabilityLedgerV2{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *LiabilityLedgerV2) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackConstructor is the Go binding used to pack the parameters required for
// contract deployment.
//
// Solidity: constructor(address admin, uint48 adminTransferDelay, address anchorer, uint64 genesisEpoch_, bytes32 genesisDigest_, uint64 genesisAsOf, address verifier_, uint8 depth_) returns()
func (liabilityLedgerV2 *LiabilityLedgerV2) PackConstructor(admin common.Address, adminTransferDelay *big.Int, anchorer common.Address, genesisEpoch_ uint64, genesisDigest_ [32]byte, genesisAsOf uint64, verifier_ common.Address, depth_ uint8) []byte {
	enc, err := liabilityLedgerV2.abi.Pack("", admin, adminTransferDelay, anchorer, genesisEpoch_, genesisDigest_, genesisAsOf, verifier_, depth_)
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
func (liabilityLedgerV2 *LiabilityLedgerV2) PackANCHORERROLE() []byte {
	enc, err := liabilityLedgerV2.abi.Pack("ANCHORER_ROLE")
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
func (liabilityLedgerV2 *LiabilityLedgerV2) TryPackANCHORERROLE() ([]byte, error) {
	return liabilityLedgerV2.abi.Pack("ANCHORER_ROLE")
}

// UnpackANCHORERROLE is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x7cfd9a51.
//
// Solidity: function ANCHORER_ROLE() view returns(bytes32)
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackANCHORERROLE(data []byte) ([32]byte, error) {
	out, err := liabilityLedgerV2.abi.Unpack("ANCHORER_ROLE", data)
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
func (liabilityLedgerV2 *LiabilityLedgerV2) PackDEFAULTADMINROLE() []byte {
	enc, err := liabilityLedgerV2.abi.Pack("DEFAULT_ADMIN_ROLE")
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
func (liabilityLedgerV2 *LiabilityLedgerV2) TryPackDEFAULTADMINROLE() ([]byte, error) {
	return liabilityLedgerV2.abi.Pack("DEFAULT_ADMIN_ROLE")
}

// UnpackDEFAULTADMINROLE is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackDEFAULTADMINROLE(data []byte) ([32]byte, error) {
	out, err := liabilityLedgerV2.abi.Unpack("DEFAULT_ADMIN_ROLE", data)
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
func (liabilityLedgerV2 *LiabilityLedgerV2) PackAcceptDefaultAdminTransfer() []byte {
	enc, err := liabilityLedgerV2.abi.Pack("acceptDefaultAdminTransfer")
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
func (liabilityLedgerV2 *LiabilityLedgerV2) TryPackAcceptDefaultAdminTransfer() ([]byte, error) {
	return liabilityLedgerV2.abi.Pack("acceptDefaultAdminTransfer")
}

// PackAnchor is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1ffcf5d4.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function anchor((uint64,bytes32,bytes32,uint64,uint64,uint64) header) returns(bytes32)
func (liabilityLedgerV2 *LiabilityLedgerV2) PackAnchor(header EpochHeader) []byte {
	enc, err := liabilityLedgerV2.abi.Pack("anchor", header)
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
func (liabilityLedgerV2 *LiabilityLedgerV2) TryPackAnchor(header EpochHeader) ([]byte, error) {
	return liabilityLedgerV2.abi.Pack("anchor", header)
}

// UnpackAnchor is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x1ffcf5d4.
//
// Solidity: function anchor((uint64,bytes32,bytes32,uint64,uint64,uint64) header) returns(bytes32)
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackAnchor(data []byte) ([32]byte, error) {
	out, err := liabilityLedgerV2.abi.Unpack("anchor", data)
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
func (liabilityLedgerV2 *LiabilityLedgerV2) PackBeginDefaultAdminTransfer(newAdmin common.Address) []byte {
	enc, err := liabilityLedgerV2.abi.Pack("beginDefaultAdminTransfer", newAdmin)
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
func (liabilityLedgerV2 *LiabilityLedgerV2) TryPackBeginDefaultAdminTransfer(newAdmin common.Address) ([]byte, error) {
	return liabilityLedgerV2.abi.Pack("beginDefaultAdminTransfer", newAdmin)
}

// PackCancelDefaultAdminTransfer is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd602b9fd.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function cancelDefaultAdminTransfer() returns()
func (liabilityLedgerV2 *LiabilityLedgerV2) PackCancelDefaultAdminTransfer() []byte {
	enc, err := liabilityLedgerV2.abi.Pack("cancelDefaultAdminTransfer")
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
func (liabilityLedgerV2 *LiabilityLedgerV2) TryPackCancelDefaultAdminTransfer() ([]byte, error) {
	return liabilityLedgerV2.abi.Pack("cancelDefaultAdminTransfer")
}

// PackCapacity is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5cfc1a51.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function capacity() view returns(uint64)
func (liabilityLedgerV2 *LiabilityLedgerV2) PackCapacity() []byte {
	enc, err := liabilityLedgerV2.abi.Pack("capacity")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCapacity is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5cfc1a51.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function capacity() view returns(uint64)
func (liabilityLedgerV2 *LiabilityLedgerV2) TryPackCapacity() ([]byte, error) {
	return liabilityLedgerV2.abi.Pack("capacity")
}

// UnpackCapacity is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5cfc1a51.
//
// Solidity: function capacity() view returns(uint64)
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackCapacity(data []byte) (uint64, error) {
	out, err := liabilityLedgerV2.abi.Unpack("capacity", data)
	if err != nil {
		return *new(uint64), err
	}
	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)
	return out0, nil
}

// PackChangeDefaultAdminDelay is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x649a5ec7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function changeDefaultAdminDelay(uint48 newDelay) returns()
func (liabilityLedgerV2 *LiabilityLedgerV2) PackChangeDefaultAdminDelay(newDelay *big.Int) []byte {
	enc, err := liabilityLedgerV2.abi.Pack("changeDefaultAdminDelay", newDelay)
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
func (liabilityLedgerV2 *LiabilityLedgerV2) TryPackChangeDefaultAdminDelay(newDelay *big.Int) ([]byte, error) {
	return liabilityLedgerV2.abi.Pack("changeDefaultAdminDelay", newDelay)
}

// PackDefaultAdmin is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x84ef8ffc.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function defaultAdmin() view returns(address)
func (liabilityLedgerV2 *LiabilityLedgerV2) PackDefaultAdmin() []byte {
	enc, err := liabilityLedgerV2.abi.Pack("defaultAdmin")
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
func (liabilityLedgerV2 *LiabilityLedgerV2) TryPackDefaultAdmin() ([]byte, error) {
	return liabilityLedgerV2.abi.Pack("defaultAdmin")
}

// UnpackDefaultAdmin is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x84ef8ffc.
//
// Solidity: function defaultAdmin() view returns(address)
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackDefaultAdmin(data []byte) (common.Address, error) {
	out, err := liabilityLedgerV2.abi.Unpack("defaultAdmin", data)
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
func (liabilityLedgerV2 *LiabilityLedgerV2) PackDefaultAdminDelay() []byte {
	enc, err := liabilityLedgerV2.abi.Pack("defaultAdminDelay")
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
func (liabilityLedgerV2 *LiabilityLedgerV2) TryPackDefaultAdminDelay() ([]byte, error) {
	return liabilityLedgerV2.abi.Pack("defaultAdminDelay")
}

// UnpackDefaultAdminDelay is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xcc8463c8.
//
// Solidity: function defaultAdminDelay() view returns(uint48)
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackDefaultAdminDelay(data []byte) (*big.Int, error) {
	out, err := liabilityLedgerV2.abi.Unpack("defaultAdminDelay", data)
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
func (liabilityLedgerV2 *LiabilityLedgerV2) PackDefaultAdminDelayIncreaseWait() []byte {
	enc, err := liabilityLedgerV2.abi.Pack("defaultAdminDelayIncreaseWait")
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
func (liabilityLedgerV2 *LiabilityLedgerV2) TryPackDefaultAdminDelayIncreaseWait() ([]byte, error) {
	return liabilityLedgerV2.abi.Pack("defaultAdminDelayIncreaseWait")
}

// UnpackDefaultAdminDelayIncreaseWait is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x022d63fb.
//
// Solidity: function defaultAdminDelayIncreaseWait() view returns(uint48)
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackDefaultAdminDelayIncreaseWait(data []byte) (*big.Int, error) {
	out, err := liabilityLedgerV2.abi.Unpack("defaultAdminDelayIncreaseWait", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackDepth is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x631c56ef.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function depth() view returns(uint8)
func (liabilityLedgerV2 *LiabilityLedgerV2) PackDepth() []byte {
	enc, err := liabilityLedgerV2.abi.Pack("depth")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackDepth is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x631c56ef.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function depth() view returns(uint8)
func (liabilityLedgerV2 *LiabilityLedgerV2) TryPackDepth() ([]byte, error) {
	return liabilityLedgerV2.abi.Pack("depth")
}

// UnpackDepth is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x631c56ef.
//
// Solidity: function depth() view returns(uint8)
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackDepth(data []byte) (uint8, error) {
	out, err := liabilityLedgerV2.abi.Unpack("depth", data)
	if err != nil {
		return *new(uint8), err
	}
	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)
	return out0, nil
}

// PackEpochAt is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7d60d355.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function epochAt(uint64 number) view returns((bytes32,bytes32,uint64,uint64,uint64,uint64))
func (liabilityLedgerV2 *LiabilityLedgerV2) PackEpochAt(number uint64) []byte {
	enc, err := liabilityLedgerV2.abi.Pack("epochAt", number)
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
func (liabilityLedgerV2 *LiabilityLedgerV2) TryPackEpochAt(number uint64) ([]byte, error) {
	return liabilityLedgerV2.abi.Pack("epochAt", number)
}

// UnpackEpochAt is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x7d60d355.
//
// Solidity: function epochAt(uint64 number) view returns((bytes32,bytes32,uint64,uint64,uint64,uint64))
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackEpochAt(data []byte) (Epoch, error) {
	out, err := liabilityLedgerV2.abi.Unpack("epochAt", data)
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
func (liabilityLedgerV2 *LiabilityLedgerV2) PackGenesisDigest() []byte {
	enc, err := liabilityLedgerV2.abi.Pack("genesisDigest")
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
func (liabilityLedgerV2 *LiabilityLedgerV2) TryPackGenesisDigest() ([]byte, error) {
	return liabilityLedgerV2.abi.Pack("genesisDigest")
}

// UnpackGenesisDigest is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x79ac2661.
//
// Solidity: function genesisDigest() view returns(bytes32)
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackGenesisDigest(data []byte) ([32]byte, error) {
	out, err := liabilityLedgerV2.abi.Unpack("genesisDigest", data)
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
func (liabilityLedgerV2 *LiabilityLedgerV2) PackGenesisEpoch() []byte {
	enc, err := liabilityLedgerV2.abi.Pack("genesisEpoch")
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
func (liabilityLedgerV2 *LiabilityLedgerV2) TryPackGenesisEpoch() ([]byte, error) {
	return liabilityLedgerV2.abi.Pack("genesisEpoch")
}

// UnpackGenesisEpoch is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb70e6be6.
//
// Solidity: function genesisEpoch() view returns(uint64)
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackGenesisEpoch(data []byte) (uint64, error) {
	out, err := liabilityLedgerV2.abi.Unpack("genesisEpoch", data)
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
func (liabilityLedgerV2 *LiabilityLedgerV2) PackGetRoleAdmin(role [32]byte) []byte {
	enc, err := liabilityLedgerV2.abi.Pack("getRoleAdmin", role)
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
func (liabilityLedgerV2 *LiabilityLedgerV2) TryPackGetRoleAdmin(role [32]byte) ([]byte, error) {
	return liabilityLedgerV2.abi.Pack("getRoleAdmin", role)
}

// UnpackGetRoleAdmin is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackGetRoleAdmin(data []byte) ([32]byte, error) {
	out, err := liabilityLedgerV2.abi.Unpack("getRoleAdmin", data)
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
func (liabilityLedgerV2 *LiabilityLedgerV2) PackGrantRole(role [32]byte, account common.Address) []byte {
	enc, err := liabilityLedgerV2.abi.Pack("grantRole", role, account)
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
func (liabilityLedgerV2 *LiabilityLedgerV2) TryPackGrantRole(role [32]byte, account common.Address) ([]byte, error) {
	return liabilityLedgerV2.abi.Pack("grantRole", role, account)
}

// PackHasRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x91d14854.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (liabilityLedgerV2 *LiabilityLedgerV2) PackHasRole(role [32]byte, account common.Address) []byte {
	enc, err := liabilityLedgerV2.abi.Pack("hasRole", role, account)
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
func (liabilityLedgerV2 *LiabilityLedgerV2) TryPackHasRole(role [32]byte, account common.Address) ([]byte, error) {
	return liabilityLedgerV2.abi.Pack("hasRole", role, account)
}

// UnpackHasRole is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackHasRole(data []byte) (bool, error) {
	out, err := liabilityLedgerV2.abi.Unpack("hasRole", data)
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
func (liabilityLedgerV2 *LiabilityLedgerV2) PackHead() []byte {
	enc, err := liabilityLedgerV2.abi.Pack("head")
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
func (liabilityLedgerV2 *LiabilityLedgerV2) TryPackHead() ([]byte, error) {
	return liabilityLedgerV2.abi.Pack("head")
}

// UnpackHead is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8f7dcfa3.
//
// Solidity: function head() view returns(bytes32)
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackHead(data []byte) ([32]byte, error) {
	out, err := liabilityLedgerV2.abi.Unpack("head", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackIsProven is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc01e5dd8.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function isProven(uint64 number, bytes32 root, uint64 leafCount) view returns(bool)
func (liabilityLedgerV2 *LiabilityLedgerV2) PackIsProven(number uint64, root [32]byte, leafCount uint64) []byte {
	enc, err := liabilityLedgerV2.abi.Pack("isProven", number, root, leafCount)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackIsProven is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc01e5dd8.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function isProven(uint64 number, bytes32 root, uint64 leafCount) view returns(bool)
func (liabilityLedgerV2 *LiabilityLedgerV2) TryPackIsProven(number uint64, root [32]byte, leafCount uint64) ([]byte, error) {
	return liabilityLedgerV2.abi.Pack("isProven", number, root, leafCount)
}

// UnpackIsProven is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc01e5dd8.
//
// Solidity: function isProven(uint64 number, bytes32 root, uint64 leafCount) view returns(bool)
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackIsProven(data []byte) (bool, error) {
	out, err := liabilityLedgerV2.abi.Unpack("isProven", data)
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
func (liabilityLedgerV2 *LiabilityLedgerV2) PackIsSealed() []byte {
	enc, err := liabilityLedgerV2.abi.Pack("isSealed")
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
func (liabilityLedgerV2 *LiabilityLedgerV2) TryPackIsSealed() ([]byte, error) {
	return liabilityLedgerV2.abi.Pack("isSealed")
}

// UnpackIsSealed is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x631f9852.
//
// Solidity: function isSealed() view returns(bool)
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackIsSealed(data []byte) (bool, error) {
	out, err := liabilityLedgerV2.abi.Unpack("isSealed", data)
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
func (liabilityLedgerV2 *LiabilityLedgerV2) PackLastAsOf() []byte {
	enc, err := liabilityLedgerV2.abi.Pack("lastAsOf")
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
func (liabilityLedgerV2 *LiabilityLedgerV2) TryPackLastAsOf() ([]byte, error) {
	return liabilityLedgerV2.abi.Pack("lastAsOf")
}

// UnpackLastAsOf is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa89210cb.
//
// Solidity: function lastAsOf() view returns(uint64)
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackLastAsOf(data []byte) (uint64, error) {
	out, err := liabilityLedgerV2.abi.Unpack("lastAsOf", data)
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
func (liabilityLedgerV2 *LiabilityLedgerV2) PackNextEpoch() []byte {
	enc, err := liabilityLedgerV2.abi.Pack("nextEpoch")
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
func (liabilityLedgerV2 *LiabilityLedgerV2) TryPackNextEpoch() ([]byte, error) {
	return liabilityLedgerV2.abi.Pack("nextEpoch")
}

// UnpackNextEpoch is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xaea0e78b.
//
// Solidity: function nextEpoch() view returns(uint64)
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackNextEpoch(data []byte) (uint64, error) {
	out, err := liabilityLedgerV2.abi.Unpack("nextEpoch", data)
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
func (liabilityLedgerV2 *LiabilityLedgerV2) PackOwner() []byte {
	enc, err := liabilityLedgerV2.abi.Pack("owner")
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
func (liabilityLedgerV2 *LiabilityLedgerV2) TryPackOwner() ([]byte, error) {
	return liabilityLedgerV2.abi.Pack("owner")
}

// UnpackOwner is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackOwner(data []byte) (common.Address, error) {
	out, err := liabilityLedgerV2.abi.Unpack("owner", data)
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
func (liabilityLedgerV2 *LiabilityLedgerV2) PackPendingDefaultAdmin() []byte {
	enc, err := liabilityLedgerV2.abi.Pack("pendingDefaultAdmin")
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
func (liabilityLedgerV2 *LiabilityLedgerV2) TryPackPendingDefaultAdmin() ([]byte, error) {
	return liabilityLedgerV2.abi.Pack("pendingDefaultAdmin")
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
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackPendingDefaultAdmin(data []byte) (PendingDefaultAdminOutput, error) {
	out, err := liabilityLedgerV2.abi.Unpack("pendingDefaultAdmin", data)
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
func (liabilityLedgerV2 *LiabilityLedgerV2) PackPendingDefaultAdminDelay() []byte {
	enc, err := liabilityLedgerV2.abi.Pack("pendingDefaultAdminDelay")
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
func (liabilityLedgerV2 *LiabilityLedgerV2) TryPackPendingDefaultAdminDelay() ([]byte, error) {
	return liabilityLedgerV2.abi.Pack("pendingDefaultAdminDelay")
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
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackPendingDefaultAdminDelay(data []byte) (PendingDefaultAdminDelayOutput, error) {
	out, err := liabilityLedgerV2.abi.Unpack("pendingDefaultAdminDelay", data)
	outstruct := new(PendingDefaultAdminDelayOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.NewDelay = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.Schedule = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	return *outstruct, nil
}

// PackProveEpoch is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3dad095b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function proveEpoch(uint64 number, bytes32 treeRoot, bytes32 totalCommitment, uint64 leafCount, bytes proof) returns(bytes32)
func (liabilityLedgerV2 *LiabilityLedgerV2) PackProveEpoch(number uint64, treeRoot [32]byte, totalCommitment [32]byte, leafCount uint64, proof []byte) []byte {
	enc, err := liabilityLedgerV2.abi.Pack("proveEpoch", number, treeRoot, totalCommitment, leafCount, proof)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackProveEpoch is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3dad095b.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function proveEpoch(uint64 number, bytes32 treeRoot, bytes32 totalCommitment, uint64 leafCount, bytes proof) returns(bytes32)
func (liabilityLedgerV2 *LiabilityLedgerV2) TryPackProveEpoch(number uint64, treeRoot [32]byte, totalCommitment [32]byte, leafCount uint64, proof []byte) ([]byte, error) {
	return liabilityLedgerV2.abi.Pack("proveEpoch", number, treeRoot, totalCommitment, leafCount, proof)
}

// UnpackProveEpoch is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x3dad095b.
//
// Solidity: function proveEpoch(uint64 number, bytes32 treeRoot, bytes32 totalCommitment, uint64 leafCount, bytes proof) returns(bytes32)
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackProveEpoch(data []byte) ([32]byte, error) {
	out, err := liabilityLedgerV2.abi.Unpack("proveEpoch", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackRenounceRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x36568abe.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (liabilityLedgerV2 *LiabilityLedgerV2) PackRenounceRole(role [32]byte, account common.Address) []byte {
	enc, err := liabilityLedgerV2.abi.Pack("renounceRole", role, account)
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
func (liabilityLedgerV2 *LiabilityLedgerV2) TryPackRenounceRole(role [32]byte, account common.Address) ([]byte, error) {
	return liabilityLedgerV2.abi.Pack("renounceRole", role, account)
}

// PackRevokeRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd547741f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (liabilityLedgerV2 *LiabilityLedgerV2) PackRevokeRole(role [32]byte, account common.Address) []byte {
	enc, err := liabilityLedgerV2.abi.Pack("revokeRole", role, account)
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
func (liabilityLedgerV2 *LiabilityLedgerV2) TryPackRevokeRole(role [32]byte, account common.Address) ([]byte, error) {
	return liabilityLedgerV2.abi.Pack("revokeRole", role, account)
}

// PackRollbackDefaultAdminDelay is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0aa6220b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function rollbackDefaultAdminDelay() returns()
func (liabilityLedgerV2 *LiabilityLedgerV2) PackRollbackDefaultAdminDelay() []byte {
	enc, err := liabilityLedgerV2.abi.Pack("rollbackDefaultAdminDelay")
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
func (liabilityLedgerV2 *LiabilityLedgerV2) TryPackRollbackDefaultAdminDelay() ([]byte, error) {
	return liabilityLedgerV2.abi.Pack("rollbackDefaultAdminDelay")
}

// PackScheme is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8518b0ad.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function scheme() pure returns(uint16)
func (liabilityLedgerV2 *LiabilityLedgerV2) PackScheme() []byte {
	enc, err := liabilityLedgerV2.abi.Pack("scheme")
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
func (liabilityLedgerV2 *LiabilityLedgerV2) TryPackScheme() ([]byte, error) {
	return liabilityLedgerV2.abi.Pack("scheme")
}

// UnpackScheme is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8518b0ad.
//
// Solidity: function scheme() pure returns(uint16)
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackScheme(data []byte) (uint16, error) {
	out, err := liabilityLedgerV2.abi.Unpack("scheme", data)
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
func (liabilityLedgerV2 *LiabilityLedgerV2) PackSeal(expectedHead [32]byte) []byte {
	enc, err := liabilityLedgerV2.abi.Pack("seal", expectedHead)
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
func (liabilityLedgerV2 *LiabilityLedgerV2) TryPackSeal(expectedHead [32]byte) ([]byte, error) {
	return liabilityLedgerV2.abi.Pack("seal", expectedHead)
}

// PackSupportsInterface is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x01ffc9a7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (liabilityLedgerV2 *LiabilityLedgerV2) PackSupportsInterface(interfaceId [4]byte) []byte {
	enc, err := liabilityLedgerV2.abi.Pack("supportsInterface", interfaceId)
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
func (liabilityLedgerV2 *LiabilityLedgerV2) TryPackSupportsInterface(interfaceId [4]byte) ([]byte, error) {
	return liabilityLedgerV2.abi.Pack("supportsInterface", interfaceId)
}

// UnpackSupportsInterface is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackSupportsInterface(data []byte) (bool, error) {
	out, err := liabilityLedgerV2.abi.Unpack("supportsInterface", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackTotalCommitmentOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x57db9232.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function totalCommitmentOf(uint64 number) view returns(bytes32)
func (liabilityLedgerV2 *LiabilityLedgerV2) PackTotalCommitmentOf(number uint64) []byte {
	enc, err := liabilityLedgerV2.abi.Pack("totalCommitmentOf", number)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTotalCommitmentOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x57db9232.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function totalCommitmentOf(uint64 number) view returns(bytes32)
func (liabilityLedgerV2 *LiabilityLedgerV2) TryPackTotalCommitmentOf(number uint64) ([]byte, error) {
	return liabilityLedgerV2.abi.Pack("totalCommitmentOf", number)
}

// UnpackTotalCommitmentOf is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x57db9232.
//
// Solidity: function totalCommitmentOf(uint64 number) view returns(bytes32)
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackTotalCommitmentOf(data []byte) ([32]byte, error) {
	out, err := liabilityLedgerV2.abi.Unpack("totalCommitmentOf", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackTreeRootOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc9c69589.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function treeRootOf(uint64 number) view returns(bytes32)
func (liabilityLedgerV2 *LiabilityLedgerV2) PackTreeRootOf(number uint64) []byte {
	enc, err := liabilityLedgerV2.abi.Pack("treeRootOf", number)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTreeRootOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc9c69589.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function treeRootOf(uint64 number) view returns(bytes32)
func (liabilityLedgerV2 *LiabilityLedgerV2) TryPackTreeRootOf(number uint64) ([]byte, error) {
	return liabilityLedgerV2.abi.Pack("treeRootOf", number)
}

// UnpackTreeRootOf is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc9c69589.
//
// Solidity: function treeRootOf(uint64 number) view returns(bytes32)
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackTreeRootOf(data []byte) ([32]byte, error) {
	out, err := liabilityLedgerV2.abi.Unpack("treeRootOf", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackVerifier is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2b7ac3f3.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function verifier() view returns(address)
func (liabilityLedgerV2 *LiabilityLedgerV2) PackVerifier() []byte {
	enc, err := liabilityLedgerV2.abi.Pack("verifier")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackVerifier is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2b7ac3f3.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function verifier() view returns(address)
func (liabilityLedgerV2 *LiabilityLedgerV2) TryPackVerifier() ([]byte, error) {
	return liabilityLedgerV2.abi.Pack("verifier")
}

// UnpackVerifier is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x2b7ac3f3.
//
// Solidity: function verifier() view returns(address)
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackVerifier(data []byte) (common.Address, error) {
	out, err := liabilityLedgerV2.abi.Unpack("verifier", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackVerifyInclusion is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc0566abe.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function verifyInclusion(uint64 number, bytes32 commitment, uint64 amount, uint64 index, (bytes32,uint64)[] siblings) view returns(bool)
func (liabilityLedgerV2 *LiabilityLedgerV2) PackVerifyInclusion(number uint64, commitment [32]byte, amount uint64, index uint64, siblings []SumNode) []byte {
	enc, err := liabilityLedgerV2.abi.Pack("verifyInclusion", number, commitment, amount, index, siblings)
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
func (liabilityLedgerV2 *LiabilityLedgerV2) TryPackVerifyInclusion(number uint64, commitment [32]byte, amount uint64, index uint64, siblings []SumNode) ([]byte, error) {
	return liabilityLedgerV2.abi.Pack("verifyInclusion", number, commitment, amount, index, siblings)
}

// UnpackVerifyInclusion is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc0566abe.
//
// Solidity: function verifyInclusion(uint64 number, bytes32 commitment, uint64 amount, uint64 index, (bytes32,uint64)[] siblings) view returns(bool)
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackVerifyInclusion(data []byte) (bool, error) {
	out, err := liabilityLedgerV2.abi.Unpack("verifyInclusion", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// LiabilityLedgerV2DefaultAdminDelayChangeCanceled represents a DefaultAdminDelayChangeCanceled event raised by the LiabilityLedgerV2 contract.
type LiabilityLedgerV2DefaultAdminDelayChangeCanceled struct {
	Raw *types.Log // Blockchain specific contextual infos
}

const LiabilityLedgerV2DefaultAdminDelayChangeCanceledEventName = "DefaultAdminDelayChangeCanceled"

// ContractEventName returns the user-defined event name.
func (LiabilityLedgerV2DefaultAdminDelayChangeCanceled) ContractEventName() string {
	return LiabilityLedgerV2DefaultAdminDelayChangeCanceledEventName
}

// UnpackDefaultAdminDelayChangeCanceledEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event DefaultAdminDelayChangeCanceled()
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackDefaultAdminDelayChangeCanceledEvent(log *types.Log) (*LiabilityLedgerV2DefaultAdminDelayChangeCanceled, error) {
	event := "DefaultAdminDelayChangeCanceled"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != liabilityLedgerV2.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(LiabilityLedgerV2DefaultAdminDelayChangeCanceled)
	if len(log.Data) > 0 {
		if err := liabilityLedgerV2.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range liabilityLedgerV2.abi.Events[event].Inputs {
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

// LiabilityLedgerV2DefaultAdminDelayChangeScheduled represents a DefaultAdminDelayChangeScheduled event raised by the LiabilityLedgerV2 contract.
type LiabilityLedgerV2DefaultAdminDelayChangeScheduled struct {
	NewDelay       *big.Int
	EffectSchedule *big.Int
	Raw            *types.Log // Blockchain specific contextual infos
}

const LiabilityLedgerV2DefaultAdminDelayChangeScheduledEventName = "DefaultAdminDelayChangeScheduled"

// ContractEventName returns the user-defined event name.
func (LiabilityLedgerV2DefaultAdminDelayChangeScheduled) ContractEventName() string {
	return LiabilityLedgerV2DefaultAdminDelayChangeScheduledEventName
}

// UnpackDefaultAdminDelayChangeScheduledEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event DefaultAdminDelayChangeScheduled(uint48 newDelay, uint48 effectSchedule)
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackDefaultAdminDelayChangeScheduledEvent(log *types.Log) (*LiabilityLedgerV2DefaultAdminDelayChangeScheduled, error) {
	event := "DefaultAdminDelayChangeScheduled"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != liabilityLedgerV2.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(LiabilityLedgerV2DefaultAdminDelayChangeScheduled)
	if len(log.Data) > 0 {
		if err := liabilityLedgerV2.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range liabilityLedgerV2.abi.Events[event].Inputs {
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

// LiabilityLedgerV2DefaultAdminTransferCanceled represents a DefaultAdminTransferCanceled event raised by the LiabilityLedgerV2 contract.
type LiabilityLedgerV2DefaultAdminTransferCanceled struct {
	Raw *types.Log // Blockchain specific contextual infos
}

const LiabilityLedgerV2DefaultAdminTransferCanceledEventName = "DefaultAdminTransferCanceled"

// ContractEventName returns the user-defined event name.
func (LiabilityLedgerV2DefaultAdminTransferCanceled) ContractEventName() string {
	return LiabilityLedgerV2DefaultAdminTransferCanceledEventName
}

// UnpackDefaultAdminTransferCanceledEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event DefaultAdminTransferCanceled()
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackDefaultAdminTransferCanceledEvent(log *types.Log) (*LiabilityLedgerV2DefaultAdminTransferCanceled, error) {
	event := "DefaultAdminTransferCanceled"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != liabilityLedgerV2.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(LiabilityLedgerV2DefaultAdminTransferCanceled)
	if len(log.Data) > 0 {
		if err := liabilityLedgerV2.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range liabilityLedgerV2.abi.Events[event].Inputs {
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

// LiabilityLedgerV2DefaultAdminTransferScheduled represents a DefaultAdminTransferScheduled event raised by the LiabilityLedgerV2 contract.
type LiabilityLedgerV2DefaultAdminTransferScheduled struct {
	NewAdmin       common.Address
	AcceptSchedule *big.Int
	Raw            *types.Log // Blockchain specific contextual infos
}

const LiabilityLedgerV2DefaultAdminTransferScheduledEventName = "DefaultAdminTransferScheduled"

// ContractEventName returns the user-defined event name.
func (LiabilityLedgerV2DefaultAdminTransferScheduled) ContractEventName() string {
	return LiabilityLedgerV2DefaultAdminTransferScheduledEventName
}

// UnpackDefaultAdminTransferScheduledEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event DefaultAdminTransferScheduled(address indexed newAdmin, uint48 acceptSchedule)
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackDefaultAdminTransferScheduledEvent(log *types.Log) (*LiabilityLedgerV2DefaultAdminTransferScheduled, error) {
	event := "DefaultAdminTransferScheduled"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != liabilityLedgerV2.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(LiabilityLedgerV2DefaultAdminTransferScheduled)
	if len(log.Data) > 0 {
		if err := liabilityLedgerV2.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range liabilityLedgerV2.abi.Events[event].Inputs {
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

// LiabilityLedgerV2EpochAnchored represents a EpochAnchored event raised by the LiabilityLedgerV2 contract.
type LiabilityLedgerV2EpochAnchored struct {
	Number         uint64
	Digest         [32]byte
	PreviousDigest [32]byte
	Root           [32]byte
	Total          uint64
	LeafCount      uint64
	AsOf           uint64
	Raw            *types.Log // Blockchain specific contextual infos
}

const LiabilityLedgerV2EpochAnchoredEventName = "EpochAnchored"

// ContractEventName returns the user-defined event name.
func (LiabilityLedgerV2EpochAnchored) ContractEventName() string {
	return LiabilityLedgerV2EpochAnchoredEventName
}

// UnpackEpochAnchoredEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event EpochAnchored(uint64 indexed number, bytes32 indexed digest, bytes32 previousDigest, bytes32 root, uint64 total, uint64 leafCount, uint64 asOf)
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackEpochAnchoredEvent(log *types.Log) (*LiabilityLedgerV2EpochAnchored, error) {
	event := "EpochAnchored"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != liabilityLedgerV2.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(LiabilityLedgerV2EpochAnchored)
	if len(log.Data) > 0 {
		if err := liabilityLedgerV2.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range liabilityLedgerV2.abi.Events[event].Inputs {
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

// LiabilityLedgerV2EpochProven represents a EpochProven event raised by the LiabilityLedgerV2 contract.
type LiabilityLedgerV2EpochProven struct {
	Number          uint64
	Root            [32]byte
	TreeRoot        [32]byte
	TotalCommitment [32]byte
	LeafCount       uint64
	Raw             *types.Log // Blockchain specific contextual infos
}

const LiabilityLedgerV2EpochProvenEventName = "EpochProven"

// ContractEventName returns the user-defined event name.
func (LiabilityLedgerV2EpochProven) ContractEventName() string {
	return LiabilityLedgerV2EpochProvenEventName
}

// UnpackEpochProvenEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event EpochProven(uint64 indexed number, bytes32 indexed root, bytes32 treeRoot, bytes32 totalCommitment, uint64 leafCount)
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackEpochProvenEvent(log *types.Log) (*LiabilityLedgerV2EpochProven, error) {
	event := "EpochProven"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != liabilityLedgerV2.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(LiabilityLedgerV2EpochProven)
	if len(log.Data) > 0 {
		if err := liabilityLedgerV2.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range liabilityLedgerV2.abi.Events[event].Inputs {
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

// LiabilityLedgerV2LedgerSealed represents a LedgerSealed event raised by the LiabilityLedgerV2 contract.
type LiabilityLedgerV2LedgerSealed struct {
	Head      [32]byte
	NextEpoch uint64
	LastAsOf  uint64
	Raw       *types.Log // Blockchain specific contextual infos
}

const LiabilityLedgerV2LedgerSealedEventName = "LedgerSealed"

// ContractEventName returns the user-defined event name.
func (LiabilityLedgerV2LedgerSealed) ContractEventName() string {
	return LiabilityLedgerV2LedgerSealedEventName
}

// UnpackLedgerSealedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event LedgerSealed(bytes32 indexed head, uint64 nextEpoch, uint64 lastAsOf)
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackLedgerSealedEvent(log *types.Log) (*LiabilityLedgerV2LedgerSealed, error) {
	event := "LedgerSealed"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != liabilityLedgerV2.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(LiabilityLedgerV2LedgerSealed)
	if len(log.Data) > 0 {
		if err := liabilityLedgerV2.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range liabilityLedgerV2.abi.Events[event].Inputs {
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

// LiabilityLedgerV2RoleAdminChanged represents a RoleAdminChanged event raised by the LiabilityLedgerV2 contract.
type LiabilityLedgerV2RoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               *types.Log // Blockchain specific contextual infos
}

const LiabilityLedgerV2RoleAdminChangedEventName = "RoleAdminChanged"

// ContractEventName returns the user-defined event name.
func (LiabilityLedgerV2RoleAdminChanged) ContractEventName() string {
	return LiabilityLedgerV2RoleAdminChangedEventName
}

// UnpackRoleAdminChangedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackRoleAdminChangedEvent(log *types.Log) (*LiabilityLedgerV2RoleAdminChanged, error) {
	event := "RoleAdminChanged"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != liabilityLedgerV2.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(LiabilityLedgerV2RoleAdminChanged)
	if len(log.Data) > 0 {
		if err := liabilityLedgerV2.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range liabilityLedgerV2.abi.Events[event].Inputs {
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

// LiabilityLedgerV2RoleGranted represents a RoleGranted event raised by the LiabilityLedgerV2 contract.
type LiabilityLedgerV2RoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     *types.Log // Blockchain specific contextual infos
}

const LiabilityLedgerV2RoleGrantedEventName = "RoleGranted"

// ContractEventName returns the user-defined event name.
func (LiabilityLedgerV2RoleGranted) ContractEventName() string {
	return LiabilityLedgerV2RoleGrantedEventName
}

// UnpackRoleGrantedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackRoleGrantedEvent(log *types.Log) (*LiabilityLedgerV2RoleGranted, error) {
	event := "RoleGranted"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != liabilityLedgerV2.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(LiabilityLedgerV2RoleGranted)
	if len(log.Data) > 0 {
		if err := liabilityLedgerV2.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range liabilityLedgerV2.abi.Events[event].Inputs {
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

// LiabilityLedgerV2RoleRevoked represents a RoleRevoked event raised by the LiabilityLedgerV2 contract.
type LiabilityLedgerV2RoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     *types.Log // Blockchain specific contextual infos
}

const LiabilityLedgerV2RoleRevokedEventName = "RoleRevoked"

// ContractEventName returns the user-defined event name.
func (LiabilityLedgerV2RoleRevoked) ContractEventName() string {
	return LiabilityLedgerV2RoleRevokedEventName
}

// UnpackRoleRevokedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackRoleRevokedEvent(log *types.Log) (*LiabilityLedgerV2RoleRevoked, error) {
	event := "RoleRevoked"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != liabilityLedgerV2.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(LiabilityLedgerV2RoleRevoked)
	if len(log.Data) > 0 {
		if err := liabilityLedgerV2.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range liabilityLedgerV2.abi.Events[event].Inputs {
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
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], liabilityLedgerV2.abi.Errors["AccessControlBadConfirmation"].ID.Bytes()[:4]) {
		return liabilityLedgerV2.UnpackAccessControlBadConfirmationError(raw[4:])
	}
	if bytes.Equal(raw[:4], liabilityLedgerV2.abi.Errors["AccessControlEnforcedDefaultAdminDelay"].ID.Bytes()[:4]) {
		return liabilityLedgerV2.UnpackAccessControlEnforcedDefaultAdminDelayError(raw[4:])
	}
	if bytes.Equal(raw[:4], liabilityLedgerV2.abi.Errors["AccessControlEnforcedDefaultAdminRules"].ID.Bytes()[:4]) {
		return liabilityLedgerV2.UnpackAccessControlEnforcedDefaultAdminRulesError(raw[4:])
	}
	if bytes.Equal(raw[:4], liabilityLedgerV2.abi.Errors["AccessControlInvalidDefaultAdmin"].ID.Bytes()[:4]) {
		return liabilityLedgerV2.UnpackAccessControlInvalidDefaultAdminError(raw[4:])
	}
	if bytes.Equal(raw[:4], liabilityLedgerV2.abi.Errors["AccessControlUnauthorizedAccount"].ID.Bytes()[:4]) {
		return liabilityLedgerV2.UnpackAccessControlUnauthorizedAccountError(raw[4:])
	}
	if bytes.Equal(raw[:4], liabilityLedgerV2.abi.Errors["AnchorerNotGranted"].ID.Bytes()[:4]) {
		return liabilityLedgerV2.UnpackAnchorerNotGrantedError(raw[4:])
	}
	if bytes.Equal(raw[:4], liabilityLedgerV2.abi.Errors["DepthOutOfRange"].ID.Bytes()[:4]) {
		return liabilityLedgerV2.UnpackDepthOutOfRangeError(raw[4:])
	}
	if bytes.Equal(raw[:4], liabilityLedgerV2.abi.Errors["EmptyEpochMismatch"].ID.Bytes()[:4]) {
		return liabilityLedgerV2.UnpackEmptyEpochMismatchError(raw[4:])
	}
	if bytes.Equal(raw[:4], liabilityLedgerV2.abi.Errors["EpochAlreadyAnchored"].ID.Bytes()[:4]) {
		return liabilityLedgerV2.UnpackEpochAlreadyAnchoredError(raw[4:])
	}
	if bytes.Equal(raw[:4], liabilityLedgerV2.abi.Errors["EpochBeforeGenesis"].ID.Bytes()[:4]) {
		return liabilityLedgerV2.UnpackEpochBeforeGenesisError(raw[4:])
	}
	if bytes.Equal(raw[:4], liabilityLedgerV2.abi.Errors["EpochClosed"].ID.Bytes()[:4]) {
		return liabilityLedgerV2.UnpackEpochClosedError(raw[4:])
	}
	if bytes.Equal(raw[:4], liabilityLedgerV2.abi.Errors["EpochConflict"].ID.Bytes()[:4]) {
		return liabilityLedgerV2.UnpackEpochConflictError(raw[4:])
	}
	if bytes.Equal(raw[:4], liabilityLedgerV2.abi.Errors["EpochNotProven"].ID.Bytes()[:4]) {
		return liabilityLedgerV2.UnpackEpochNotProvenError(raw[4:])
	}
	if bytes.Equal(raw[:4], liabilityLedgerV2.abi.Errors["EpochOutOfSequence"].ID.Bytes()[:4]) {
		return liabilityLedgerV2.UnpackEpochOutOfSequenceError(raw[4:])
	}
	if bytes.Equal(raw[:4], liabilityLedgerV2.abi.Errors["HeadMismatch"].ID.Bytes()[:4]) {
		return liabilityLedgerV2.UnpackHeadMismatchError(raw[4:])
	}
	if bytes.Equal(raw[:4], liabilityLedgerV2.abi.Errors["InvalidGenesis"].ID.Bytes()[:4]) {
		return liabilityLedgerV2.UnpackInvalidGenesisError(raw[4:])
	}
	if bytes.Equal(raw[:4], liabilityLedgerV2.abi.Errors["LedgerIsSealed"].ID.Bytes()[:4]) {
		return liabilityLedgerV2.UnpackLedgerIsSealedError(raw[4:])
	}
	if bytes.Equal(raw[:4], liabilityLedgerV2.abi.Errors["NotAFieldElement"].ID.Bytes()[:4]) {
		return liabilityLedgerV2.UnpackNotAFieldElementError(raw[4:])
	}
	if bytes.Equal(raw[:4], liabilityLedgerV2.abi.Errors["NotAVerifier"].ID.Bytes()[:4]) {
		return liabilityLedgerV2.UnpackNotAVerifierError(raw[4:])
	}
	if bytes.Equal(raw[:4], liabilityLedgerV2.abi.Errors["PreviousDigestMismatch"].ID.Bytes()[:4]) {
		return liabilityLedgerV2.UnpackPreviousDigestMismatchError(raw[4:])
	}
	if bytes.Equal(raw[:4], liabilityLedgerV2.abi.Errors["ProofRejected"].ID.Bytes()[:4]) {
		return liabilityLedgerV2.UnpackProofRejectedError(raw[4:])
	}
	if bytes.Equal(raw[:4], liabilityLedgerV2.abi.Errors["SafeCastOverflowedUintDowncast"].ID.Bytes()[:4]) {
		return liabilityLedgerV2.UnpackSafeCastOverflowedUintDowncastError(raw[4:])
	}
	if bytes.Equal(raw[:4], liabilityLedgerV2.abi.Errors["SnapshotInFuture"].ID.Bytes()[:4]) {
		return liabilityLedgerV2.UnpackSnapshotInFutureError(raw[4:])
	}
	if bytes.Equal(raw[:4], liabilityLedgerV2.abi.Errors["SnapshotNotAfterPrevious"].ID.Bytes()[:4]) {
		return liabilityLedgerV2.UnpackSnapshotNotAfterPreviousError(raw[4:])
	}
	if bytes.Equal(raw[:4], liabilityLedgerV2.abi.Errors["StatementAlreadyProven"].ID.Bytes()[:4]) {
		return liabilityLedgerV2.UnpackStatementAlreadyProvenError(raw[4:])
	}
	if bytes.Equal(raw[:4], liabilityLedgerV2.abi.Errors["TooManyLeaves"].ID.Bytes()[:4]) {
		return liabilityLedgerV2.UnpackTooManyLeavesError(raw[4:])
	}
	if bytes.Equal(raw[:4], liabilityLedgerV2.abi.Errors["TotalNotHidden"].ID.Bytes()[:4]) {
		return liabilityLedgerV2.UnpackTotalNotHiddenError(raw[4:])
	}
	if bytes.Equal(raw[:4], liabilityLedgerV2.abi.Errors["UnknownEpoch"].ID.Bytes()[:4]) {
		return liabilityLedgerV2.UnpackUnknownEpochError(raw[4:])
	}
	if bytes.Equal(raw[:4], liabilityLedgerV2.abi.Errors["ZeroAnchorer"].ID.Bytes()[:4]) {
		return liabilityLedgerV2.UnpackZeroAnchorerError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// LiabilityLedgerV2AccessControlBadConfirmation represents a AccessControlBadConfirmation error raised by the LiabilityLedgerV2 contract.
type LiabilityLedgerV2AccessControlBadConfirmation struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AccessControlBadConfirmation()
func LiabilityLedgerV2AccessControlBadConfirmationErrorID() common.Hash {
	return common.HexToHash("0x6697b23232a647058342c0724fe7c415cab25915b54e5dbc03f233173d37b41c")
}

// UnpackAccessControlBadConfirmationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AccessControlBadConfirmation()
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackAccessControlBadConfirmationError(raw []byte) (*LiabilityLedgerV2AccessControlBadConfirmation, error) {
	out := new(LiabilityLedgerV2AccessControlBadConfirmation)
	if err := liabilityLedgerV2.abi.UnpackIntoInterface(out, "AccessControlBadConfirmation", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// LiabilityLedgerV2AccessControlEnforcedDefaultAdminDelay represents a AccessControlEnforcedDefaultAdminDelay error raised by the LiabilityLedgerV2 contract.
type LiabilityLedgerV2AccessControlEnforcedDefaultAdminDelay struct {
	Schedule *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AccessControlEnforcedDefaultAdminDelay(uint48 schedule)
func LiabilityLedgerV2AccessControlEnforcedDefaultAdminDelayErrorID() common.Hash {
	return common.HexToHash("0x19ca5ebb8fb33f00e502c9392eddab1501674629178bf69b853cf037aaf4bb5d")
}

// UnpackAccessControlEnforcedDefaultAdminDelayError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AccessControlEnforcedDefaultAdminDelay(uint48 schedule)
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackAccessControlEnforcedDefaultAdminDelayError(raw []byte) (*LiabilityLedgerV2AccessControlEnforcedDefaultAdminDelay, error) {
	out := new(LiabilityLedgerV2AccessControlEnforcedDefaultAdminDelay)
	if err := liabilityLedgerV2.abi.UnpackIntoInterface(out, "AccessControlEnforcedDefaultAdminDelay", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// LiabilityLedgerV2AccessControlEnforcedDefaultAdminRules represents a AccessControlEnforcedDefaultAdminRules error raised by the LiabilityLedgerV2 contract.
type LiabilityLedgerV2AccessControlEnforcedDefaultAdminRules struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AccessControlEnforcedDefaultAdminRules()
func LiabilityLedgerV2AccessControlEnforcedDefaultAdminRulesErrorID() common.Hash {
	return common.HexToHash("0x3fc3c27ae3db78c81b8f6e685172134623efa268ee8cd8d54be38ad2a74fc13b")
}

// UnpackAccessControlEnforcedDefaultAdminRulesError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AccessControlEnforcedDefaultAdminRules()
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackAccessControlEnforcedDefaultAdminRulesError(raw []byte) (*LiabilityLedgerV2AccessControlEnforcedDefaultAdminRules, error) {
	out := new(LiabilityLedgerV2AccessControlEnforcedDefaultAdminRules)
	if err := liabilityLedgerV2.abi.UnpackIntoInterface(out, "AccessControlEnforcedDefaultAdminRules", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// LiabilityLedgerV2AccessControlInvalidDefaultAdmin represents a AccessControlInvalidDefaultAdmin error raised by the LiabilityLedgerV2 contract.
type LiabilityLedgerV2AccessControlInvalidDefaultAdmin struct {
	DefaultAdmin common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AccessControlInvalidDefaultAdmin(address defaultAdmin)
func LiabilityLedgerV2AccessControlInvalidDefaultAdminErrorID() common.Hash {
	return common.HexToHash("0xc22c8022f2a840d6b6a9f113407715f5bbd4e88c1b0dd9434dc00700ba609ed4")
}

// UnpackAccessControlInvalidDefaultAdminError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AccessControlInvalidDefaultAdmin(address defaultAdmin)
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackAccessControlInvalidDefaultAdminError(raw []byte) (*LiabilityLedgerV2AccessControlInvalidDefaultAdmin, error) {
	out := new(LiabilityLedgerV2AccessControlInvalidDefaultAdmin)
	if err := liabilityLedgerV2.abi.UnpackIntoInterface(out, "AccessControlInvalidDefaultAdmin", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// LiabilityLedgerV2AccessControlUnauthorizedAccount represents a AccessControlUnauthorizedAccount error raised by the LiabilityLedgerV2 contract.
type LiabilityLedgerV2AccessControlUnauthorizedAccount struct {
	Account    common.Address
	NeededRole [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AccessControlUnauthorizedAccount(address account, bytes32 neededRole)
func LiabilityLedgerV2AccessControlUnauthorizedAccountErrorID() common.Hash {
	return common.HexToHash("0xe2517d3fbfae6f8515ef5ff1ccedc3933ab0cbbda0b492c06eb54ad10ef03b3e")
}

// UnpackAccessControlUnauthorizedAccountError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AccessControlUnauthorizedAccount(address account, bytes32 neededRole)
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackAccessControlUnauthorizedAccountError(raw []byte) (*LiabilityLedgerV2AccessControlUnauthorizedAccount, error) {
	out := new(LiabilityLedgerV2AccessControlUnauthorizedAccount)
	if err := liabilityLedgerV2.abi.UnpackIntoInterface(out, "AccessControlUnauthorizedAccount", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// LiabilityLedgerV2AnchorerNotGranted represents a AnchorerNotGranted error raised by the LiabilityLedgerV2 contract.
type LiabilityLedgerV2AnchorerNotGranted struct {
	Anchorer common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AnchorerNotGranted(address anchorer)
func LiabilityLedgerV2AnchorerNotGrantedErrorID() common.Hash {
	return common.HexToHash("0x3012fc8bb5dd75c0c034d7c58aa1982d269d598072100fabc24d4de862d1e083")
}

// UnpackAnchorerNotGrantedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AnchorerNotGranted(address anchorer)
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackAnchorerNotGrantedError(raw []byte) (*LiabilityLedgerV2AnchorerNotGranted, error) {
	out := new(LiabilityLedgerV2AnchorerNotGranted)
	if err := liabilityLedgerV2.abi.UnpackIntoInterface(out, "AnchorerNotGranted", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// LiabilityLedgerV2DepthOutOfRange represents a DepthOutOfRange error raised by the LiabilityLedgerV2 contract.
type LiabilityLedgerV2DepthOutOfRange struct {
	Depth   uint8
	Maximum uint8
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error DepthOutOfRange(uint8 depth, uint8 maximum)
func LiabilityLedgerV2DepthOutOfRangeErrorID() common.Hash {
	return common.HexToHash("0x69a5dd38589f776120570bc547ad76f3f9b8b1a5746d06d9a07dd4943effbae3")
}

// UnpackDepthOutOfRangeError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error DepthOutOfRange(uint8 depth, uint8 maximum)
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackDepthOutOfRangeError(raw []byte) (*LiabilityLedgerV2DepthOutOfRange, error) {
	out := new(LiabilityLedgerV2DepthOutOfRange)
	if err := liabilityLedgerV2.abi.UnpackIntoInterface(out, "DepthOutOfRange", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// LiabilityLedgerV2EmptyEpochMismatch represents a EmptyEpochMismatch error raised by the LiabilityLedgerV2 contract.
type LiabilityLedgerV2EmptyEpochMismatch struct {
	Root  [32]byte
	Total uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error EmptyEpochMismatch(bytes32 root, uint64 total)
func LiabilityLedgerV2EmptyEpochMismatchErrorID() common.Hash {
	return common.HexToHash("0xde718304bb7fe1b2763221798581b76b906780448b2dcd8bad16597d786f1080")
}

// UnpackEmptyEpochMismatchError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error EmptyEpochMismatch(bytes32 root, uint64 total)
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackEmptyEpochMismatchError(raw []byte) (*LiabilityLedgerV2EmptyEpochMismatch, error) {
	out := new(LiabilityLedgerV2EmptyEpochMismatch)
	if err := liabilityLedgerV2.abi.UnpackIntoInterface(out, "EmptyEpochMismatch", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// LiabilityLedgerV2EpochAlreadyAnchored represents a EpochAlreadyAnchored error raised by the LiabilityLedgerV2 contract.
type LiabilityLedgerV2EpochAlreadyAnchored struct {
	Number uint64
	Digest [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error EpochAlreadyAnchored(uint64 number, bytes32 digest)
func LiabilityLedgerV2EpochAlreadyAnchoredErrorID() common.Hash {
	return common.HexToHash("0xab455a25015fe0023a0a76c842f75125738bc3b7b3167d58a41d09451f0baf74")
}

// UnpackEpochAlreadyAnchoredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error EpochAlreadyAnchored(uint64 number, bytes32 digest)
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackEpochAlreadyAnchoredError(raw []byte) (*LiabilityLedgerV2EpochAlreadyAnchored, error) {
	out := new(LiabilityLedgerV2EpochAlreadyAnchored)
	if err := liabilityLedgerV2.abi.UnpackIntoInterface(out, "EpochAlreadyAnchored", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// LiabilityLedgerV2EpochBeforeGenesis represents a EpochBeforeGenesis error raised by the LiabilityLedgerV2 contract.
type LiabilityLedgerV2EpochBeforeGenesis struct {
	GenesisEpoch uint64
	Number       uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error EpochBeforeGenesis(uint64 genesisEpoch, uint64 number)
func LiabilityLedgerV2EpochBeforeGenesisErrorID() common.Hash {
	return common.HexToHash("0xacc3e9cb2b490f2ca39781764307e7ead97b1740e688ec48ce82eb74ab3cad45")
}

// UnpackEpochBeforeGenesisError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error EpochBeforeGenesis(uint64 genesisEpoch, uint64 number)
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackEpochBeforeGenesisError(raw []byte) (*LiabilityLedgerV2EpochBeforeGenesis, error) {
	out := new(LiabilityLedgerV2EpochBeforeGenesis)
	if err := liabilityLedgerV2.abi.UnpackIntoInterface(out, "EpochBeforeGenesis", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// LiabilityLedgerV2EpochClosed represents a EpochClosed error raised by the LiabilityLedgerV2 contract.
type LiabilityLedgerV2EpochClosed struct {
	Number    uint64
	NextEpoch uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error EpochClosed(uint64 number, uint64 nextEpoch)
func LiabilityLedgerV2EpochClosedErrorID() common.Hash {
	return common.HexToHash("0x2e5386d6d06806d988b0873d7fb2447aa3e8f6c70a9998aee2dc1b4927b781fb")
}

// UnpackEpochClosedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error EpochClosed(uint64 number, uint64 nextEpoch)
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackEpochClosedError(raw []byte) (*LiabilityLedgerV2EpochClosed, error) {
	out := new(LiabilityLedgerV2EpochClosed)
	if err := liabilityLedgerV2.abi.UnpackIntoInterface(out, "EpochClosed", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// LiabilityLedgerV2EpochConflict represents a EpochConflict error raised by the LiabilityLedgerV2 contract.
type LiabilityLedgerV2EpochConflict struct {
	Number   uint64
	Anchored [32]byte
	Given    [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error EpochConflict(uint64 number, bytes32 anchored, bytes32 given)
func LiabilityLedgerV2EpochConflictErrorID() common.Hash {
	return common.HexToHash("0x6d613c397e2d7ea70b114e132c49a28c0fa5de67dbb8170276880d74174393a9")
}

// UnpackEpochConflictError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error EpochConflict(uint64 number, bytes32 anchored, bytes32 given)
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackEpochConflictError(raw []byte) (*LiabilityLedgerV2EpochConflict, error) {
	out := new(LiabilityLedgerV2EpochConflict)
	if err := liabilityLedgerV2.abi.UnpackIntoInterface(out, "EpochConflict", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// LiabilityLedgerV2EpochNotProven represents a EpochNotProven error raised by the LiabilityLedgerV2 contract.
type LiabilityLedgerV2EpochNotProven struct {
	Number    uint64
	Root      [32]byte
	LeafCount uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error EpochNotProven(uint64 number, bytes32 root, uint64 leafCount)
func LiabilityLedgerV2EpochNotProvenErrorID() common.Hash {
	return common.HexToHash("0x72deac18ba28cebc115a29182a0088eaa31dd139d02a2ea57f0ebe3a121b077f")
}

// UnpackEpochNotProvenError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error EpochNotProven(uint64 number, bytes32 root, uint64 leafCount)
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackEpochNotProvenError(raw []byte) (*LiabilityLedgerV2EpochNotProven, error) {
	out := new(LiabilityLedgerV2EpochNotProven)
	if err := liabilityLedgerV2.abi.UnpackIntoInterface(out, "EpochNotProven", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// LiabilityLedgerV2EpochOutOfSequence represents a EpochOutOfSequence error raised by the LiabilityLedgerV2 contract.
type LiabilityLedgerV2EpochOutOfSequence struct {
	Expected uint64
	Given    uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error EpochOutOfSequence(uint64 expected, uint64 given)
func LiabilityLedgerV2EpochOutOfSequenceErrorID() common.Hash {
	return common.HexToHash("0x7a8d061d59616552e34e07c4e10cd96598f8a51489f859f21055dba983f36722")
}

// UnpackEpochOutOfSequenceError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error EpochOutOfSequence(uint64 expected, uint64 given)
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackEpochOutOfSequenceError(raw []byte) (*LiabilityLedgerV2EpochOutOfSequence, error) {
	out := new(LiabilityLedgerV2EpochOutOfSequence)
	if err := liabilityLedgerV2.abi.UnpackIntoInterface(out, "EpochOutOfSequence", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// LiabilityLedgerV2HeadMismatch represents a HeadMismatch error raised by the LiabilityLedgerV2 contract.
type LiabilityLedgerV2HeadMismatch struct {
	Head     [32]byte
	Expected [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error HeadMismatch(bytes32 head, bytes32 expected)
func LiabilityLedgerV2HeadMismatchErrorID() common.Hash {
	return common.HexToHash("0x6f2de3c75673389a279dc32b221eac892abdabf6b2d716499a572ff532f5ec50")
}

// UnpackHeadMismatchError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error HeadMismatch(bytes32 head, bytes32 expected)
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackHeadMismatchError(raw []byte) (*LiabilityLedgerV2HeadMismatch, error) {
	out := new(LiabilityLedgerV2HeadMismatch)
	if err := liabilityLedgerV2.abi.UnpackIntoInterface(out, "HeadMismatch", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// LiabilityLedgerV2InvalidGenesis represents a InvalidGenesis error raised by the LiabilityLedgerV2 contract.
type LiabilityLedgerV2InvalidGenesis struct {
	GenesisEpoch  uint64
	GenesisDigest [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidGenesis(uint64 genesisEpoch, bytes32 genesisDigest)
func LiabilityLedgerV2InvalidGenesisErrorID() common.Hash {
	return common.HexToHash("0x111f019fa2344bb16ed38186e3c5142be0c96d972db2a225af2e8f4835689530")
}

// UnpackInvalidGenesisError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidGenesis(uint64 genesisEpoch, bytes32 genesisDigest)
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackInvalidGenesisError(raw []byte) (*LiabilityLedgerV2InvalidGenesis, error) {
	out := new(LiabilityLedgerV2InvalidGenesis)
	if err := liabilityLedgerV2.abi.UnpackIntoInterface(out, "InvalidGenesis", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// LiabilityLedgerV2LedgerIsSealed represents a LedgerIsSealed error raised by the LiabilityLedgerV2 contract.
type LiabilityLedgerV2LedgerIsSealed struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error LedgerIsSealed()
func LiabilityLedgerV2LedgerIsSealedErrorID() common.Hash {
	return common.HexToHash("0x712a4e2f2cfa8aa62d1c3b9f2e3bce545dad83b3813e02d35f0b2a51811808e4")
}

// UnpackLedgerIsSealedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error LedgerIsSealed()
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackLedgerIsSealedError(raw []byte) (*LiabilityLedgerV2LedgerIsSealed, error) {
	out := new(LiabilityLedgerV2LedgerIsSealed)
	if err := liabilityLedgerV2.abi.UnpackIntoInterface(out, "LedgerIsSealed", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// LiabilityLedgerV2NotAFieldElement represents a NotAFieldElement error raised by the LiabilityLedgerV2 contract.
type LiabilityLedgerV2NotAFieldElement struct {
	Value [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotAFieldElement(bytes32 value)
func LiabilityLedgerV2NotAFieldElementErrorID() common.Hash {
	return common.HexToHash("0x1c35841b2ff12967633c940be49c5ae2acd0b066a6e182741ac014d5677ee168")
}

// UnpackNotAFieldElementError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotAFieldElement(bytes32 value)
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackNotAFieldElementError(raw []byte) (*LiabilityLedgerV2NotAFieldElement, error) {
	out := new(LiabilityLedgerV2NotAFieldElement)
	if err := liabilityLedgerV2.abi.UnpackIntoInterface(out, "NotAFieldElement", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// LiabilityLedgerV2NotAVerifier represents a NotAVerifier error raised by the LiabilityLedgerV2 contract.
type LiabilityLedgerV2NotAVerifier struct {
	Verifier common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotAVerifier(address verifier)
func LiabilityLedgerV2NotAVerifierErrorID() common.Hash {
	return common.HexToHash("0x581a047cf8450836388a60b29f3fa7489b7b3cc76ca1b45b16be79537e16bef3")
}

// UnpackNotAVerifierError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotAVerifier(address verifier)
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackNotAVerifierError(raw []byte) (*LiabilityLedgerV2NotAVerifier, error) {
	out := new(LiabilityLedgerV2NotAVerifier)
	if err := liabilityLedgerV2.abi.UnpackIntoInterface(out, "NotAVerifier", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// LiabilityLedgerV2PreviousDigestMismatch represents a PreviousDigestMismatch error raised by the LiabilityLedgerV2 contract.
type LiabilityLedgerV2PreviousDigestMismatch struct {
	Expected [32]byte
	Given    [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error PreviousDigestMismatch(bytes32 expected, bytes32 given)
func LiabilityLedgerV2PreviousDigestMismatchErrorID() common.Hash {
	return common.HexToHash("0x0796b55ac9b13adadc488175cf6ef6fbf1565c2537f4e021d7268b5ab5a335d9")
}

// UnpackPreviousDigestMismatchError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error PreviousDigestMismatch(bytes32 expected, bytes32 given)
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackPreviousDigestMismatchError(raw []byte) (*LiabilityLedgerV2PreviousDigestMismatch, error) {
	out := new(LiabilityLedgerV2PreviousDigestMismatch)
	if err := liabilityLedgerV2.abi.UnpackIntoInterface(out, "PreviousDigestMismatch", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// LiabilityLedgerV2ProofRejected represents a ProofRejected error raised by the LiabilityLedgerV2 contract.
type LiabilityLedgerV2ProofRejected struct {
	Number uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ProofRejected(uint64 number)
func LiabilityLedgerV2ProofRejectedErrorID() common.Hash {
	return common.HexToHash("0xea282d242b3a014a0df5455d1368c0e4e2b5f6090080510ff8b6aa97cf54fb46")
}

// UnpackProofRejectedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ProofRejected(uint64 number)
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackProofRejectedError(raw []byte) (*LiabilityLedgerV2ProofRejected, error) {
	out := new(LiabilityLedgerV2ProofRejected)
	if err := liabilityLedgerV2.abi.UnpackIntoInterface(out, "ProofRejected", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// LiabilityLedgerV2SafeCastOverflowedUintDowncast represents a SafeCastOverflowedUintDowncast error raised by the LiabilityLedgerV2 contract.
type LiabilityLedgerV2SafeCastOverflowedUintDowncast struct {
	Bits  uint8
	Value *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SafeCastOverflowedUintDowncast(uint8 bits, uint256 value)
func LiabilityLedgerV2SafeCastOverflowedUintDowncastErrorID() common.Hash {
	return common.HexToHash("0x6dfcc6503a32754ce7a89698e18201fc5294fd4aad43edefee786f88423b1a12")
}

// UnpackSafeCastOverflowedUintDowncastError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SafeCastOverflowedUintDowncast(uint8 bits, uint256 value)
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackSafeCastOverflowedUintDowncastError(raw []byte) (*LiabilityLedgerV2SafeCastOverflowedUintDowncast, error) {
	out := new(LiabilityLedgerV2SafeCastOverflowedUintDowncast)
	if err := liabilityLedgerV2.abi.UnpackIntoInterface(out, "SafeCastOverflowedUintDowncast", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// LiabilityLedgerV2SnapshotInFuture represents a SnapshotInFuture error raised by the LiabilityLedgerV2 contract.
type LiabilityLedgerV2SnapshotInFuture struct {
	AsOf      uint64
	BlockTime uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SnapshotInFuture(uint64 asOf, uint64 blockTime)
func LiabilityLedgerV2SnapshotInFutureErrorID() common.Hash {
	return common.HexToHash("0x18aaced2679cfd7d188264841ef0e83939ee734887344b8d78b84cd47abff282")
}

// UnpackSnapshotInFutureError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SnapshotInFuture(uint64 asOf, uint64 blockTime)
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackSnapshotInFutureError(raw []byte) (*LiabilityLedgerV2SnapshotInFuture, error) {
	out := new(LiabilityLedgerV2SnapshotInFuture)
	if err := liabilityLedgerV2.abi.UnpackIntoInterface(out, "SnapshotInFuture", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// LiabilityLedgerV2SnapshotNotAfterPrevious represents a SnapshotNotAfterPrevious error raised by the LiabilityLedgerV2 contract.
type LiabilityLedgerV2SnapshotNotAfterPrevious struct {
	AsOf     uint64
	Previous uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SnapshotNotAfterPrevious(uint64 asOf, uint64 previous)
func LiabilityLedgerV2SnapshotNotAfterPreviousErrorID() common.Hash {
	return common.HexToHash("0x81ed60e995f8fcd4753ca10053fc7f547e1dd0c250e9bc06e427fee7bea5773f")
}

// UnpackSnapshotNotAfterPreviousError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SnapshotNotAfterPrevious(uint64 asOf, uint64 previous)
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackSnapshotNotAfterPreviousError(raw []byte) (*LiabilityLedgerV2SnapshotNotAfterPrevious, error) {
	out := new(LiabilityLedgerV2SnapshotNotAfterPrevious)
	if err := liabilityLedgerV2.abi.UnpackIntoInterface(out, "SnapshotNotAfterPrevious", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// LiabilityLedgerV2StatementAlreadyProven represents a StatementAlreadyProven error raised by the LiabilityLedgerV2 contract.
type LiabilityLedgerV2StatementAlreadyProven struct {
	Number    uint64
	Root      [32]byte
	LeafCount uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error StatementAlreadyProven(uint64 number, bytes32 root, uint64 leafCount)
func LiabilityLedgerV2StatementAlreadyProvenErrorID() common.Hash {
	return common.HexToHash("0x1cf1dcbb77896c20ef7b54234ad069fe200ac992ff1cda6f9863adf68e78a3ba")
}

// UnpackStatementAlreadyProvenError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error StatementAlreadyProven(uint64 number, bytes32 root, uint64 leafCount)
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackStatementAlreadyProvenError(raw []byte) (*LiabilityLedgerV2StatementAlreadyProven, error) {
	out := new(LiabilityLedgerV2StatementAlreadyProven)
	if err := liabilityLedgerV2.abi.UnpackIntoInterface(out, "StatementAlreadyProven", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// LiabilityLedgerV2TooManyLeaves represents a TooManyLeaves error raised by the LiabilityLedgerV2 contract.
type LiabilityLedgerV2TooManyLeaves struct {
	LeafCount uint64
	Maximum   uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TooManyLeaves(uint64 leafCount, uint64 maximum)
func LiabilityLedgerV2TooManyLeavesErrorID() common.Hash {
	return common.HexToHash("0x09e51df967feec9808a68989cf05b6cd9c8c3c92676ae3ea00eda6a12380a9af")
}

// UnpackTooManyLeavesError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TooManyLeaves(uint64 leafCount, uint64 maximum)
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackTooManyLeavesError(raw []byte) (*LiabilityLedgerV2TooManyLeaves, error) {
	out := new(LiabilityLedgerV2TooManyLeaves)
	if err := liabilityLedgerV2.abi.UnpackIntoInterface(out, "TooManyLeaves", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// LiabilityLedgerV2TotalNotHidden represents a TotalNotHidden error raised by the LiabilityLedgerV2 contract.
type LiabilityLedgerV2TotalNotHidden struct {
	Total uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TotalNotHidden(uint64 total)
func LiabilityLedgerV2TotalNotHiddenErrorID() common.Hash {
	return common.HexToHash("0x81d8d68b510fd9927e38393af691de24fd95be21b278fba3330615b4cf2ca6e1")
}

// UnpackTotalNotHiddenError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TotalNotHidden(uint64 total)
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackTotalNotHiddenError(raw []byte) (*LiabilityLedgerV2TotalNotHidden, error) {
	out := new(LiabilityLedgerV2TotalNotHidden)
	if err := liabilityLedgerV2.abi.UnpackIntoInterface(out, "TotalNotHidden", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// LiabilityLedgerV2UnknownEpoch represents a UnknownEpoch error raised by the LiabilityLedgerV2 contract.
type LiabilityLedgerV2UnknownEpoch struct {
	Number uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error UnknownEpoch(uint64 number)
func LiabilityLedgerV2UnknownEpochErrorID() common.Hash {
	return common.HexToHash("0xa4e21e6c8f881eff51b5764f98a67304231c1cb1ba8af636f941098552d2bd6c")
}

// UnpackUnknownEpochError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error UnknownEpoch(uint64 number)
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackUnknownEpochError(raw []byte) (*LiabilityLedgerV2UnknownEpoch, error) {
	out := new(LiabilityLedgerV2UnknownEpoch)
	if err := liabilityLedgerV2.abi.UnpackIntoInterface(out, "UnknownEpoch", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// LiabilityLedgerV2ZeroAnchorer represents a ZeroAnchorer error raised by the LiabilityLedgerV2 contract.
type LiabilityLedgerV2ZeroAnchorer struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ZeroAnchorer()
func LiabilityLedgerV2ZeroAnchorerErrorID() common.Hash {
	return common.HexToHash("0x7110c4a1e207dd190eaa2e6c40aeb0e70affea5eb498857bf51f0d9d85807115")
}

// UnpackZeroAnchorerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ZeroAnchorer()
func (liabilityLedgerV2 *LiabilityLedgerV2) UnpackZeroAnchorerError(raw []byte) (*LiabilityLedgerV2ZeroAnchorer, error) {
	out := new(LiabilityLedgerV2ZeroAnchorer)
	if err := liabilityLedgerV2.abi.UnpackIntoInterface(out, "ZeroAnchorer", raw); err != nil {
		return nil, err
	}
	return out, nil
}
