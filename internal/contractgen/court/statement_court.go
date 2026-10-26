// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package court

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

// StatementCourtMetaData contains all meta data concerning the StatementCourt contract.
var StatementCourtMetaData = bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"admin\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"adminTransferDelay\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"ledger_\",\"type\":\"address\",\"internalType\":\"contractILiabilityLedger\"},{\"name\":\"initialSigner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"signerChangeDelay_\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"DEFAULT_ADMIN_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptDefaultAdminTransfer\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"activateSigner\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"beginDefaultAdminTransfer\",\"inputs\":[{\"name\":\"newAdmin\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"cancelDefaultAdminTransfer\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"cancelSigner\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"changeDefaultAdminDelay\",\"inputs\":[{\"name\":\"newDelay\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"defaultAdmin\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"defaultAdminDelay\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"defaultAdminDelayIncreaseWait\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eip712Domain\",\"inputs\":[],\"outputs\":[{\"name\":\"fields\",\"type\":\"bytes1\",\"internalType\":\"bytes1\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"verifyingContract\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"extensions\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"faultCount\",\"inputs\":[{\"name\":\"epoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"faultsOf\",\"inputs\":[{\"name\":\"epoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple[]\",\"internalType\":\"structFault[]\",\"components\":[{\"name\":\"kind\",\"type\":\"uint8\",\"internalType\":\"enumFaultKind\"},{\"name\":\"evidence\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"prover\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"provenAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRoleAdmin\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"grantRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"hasRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isProven\",\"inputs\":[{\"name\":\"evidence\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ledger\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractILiabilityLedger\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingDefaultAdmin\",\"inputs\":[],\"outputs\":[{\"name\":\"newAdmin\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"schedule\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingDefaultAdminDelay\",\"inputs\":[],\"outputs\":[{\"name\":\"newDelay\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"schedule\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingSigner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structPendingSigner\",\"components\":[{\"name\":\"signer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"fromEpoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"activatesAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proveBadProof\",\"inputs\":[{\"name\":\"statement\",\"type\":\"tuple\",\"internalType\":\"structEpochStatement\",\"components\":[{\"name\":\"ledger\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"epoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"commitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"amount\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"index\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"siblingsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"siblings\",\"type\":\"tuple[]\",\"internalType\":\"structSumNode[]\",\"components\":[{\"name\":\"hash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sum\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"proveEquivocation\",\"inputs\":[{\"name\":\"first\",\"type\":\"tuple\",\"internalType\":\"structEpochStatement\",\"components\":[{\"name\":\"ledger\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"epoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"commitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"amount\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"index\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"siblingsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"firstSignature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"second\",\"type\":\"tuple\",\"internalType\":\"structEpochStatement\",\"components\":[{\"name\":\"ledger\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"epoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"commitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"amount\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"index\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"siblingsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"secondSignature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"recoverSigner\",\"inputs\":[{\"name\":\"statement\",\"type\":\"tuple\",\"internalType\":\"structEpochStatement\",\"components\":[{\"name\":\"ledger\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"epoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"commitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"amount\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"index\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"siblingsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revokeRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"rollbackDefaultAdminDelay\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"scheduleSigner\",\"inputs\":[{\"name\":\"signer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"fromEpoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"siblingsHash\",\"inputs\":[{\"name\":\"siblings\",\"type\":\"tuple[]\",\"internalType\":\"structSumNode[]\",\"components\":[{\"name\":\"hash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sum\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"signerChangeDelay\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"signerOf\",\"inputs\":[{\"name\":\"epoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"signerTerms\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple[]\",\"internalType\":\"structSignerTerm[]\",\"components\":[{\"name\":\"signer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"fromEpoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"statementDigest\",\"inputs\":[{\"name\":\"statement\",\"type\":\"tuple\",\"internalType\":\"structEpochStatement\",\"components\":[{\"name\":\"ledger\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"epoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"commitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"amount\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"index\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"siblingsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalFaults\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"DefaultAdminDelayChangeCanceled\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DefaultAdminDelayChangeScheduled\",\"inputs\":[{\"name\":\"newDelay\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"},{\"name\":\"effectSchedule\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DefaultAdminTransferCanceled\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DefaultAdminTransferScheduled\",\"inputs\":[{\"name\":\"newAdmin\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"acceptSchedule\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EIP712DomainChanged\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FaultProven\",\"inputs\":[{\"name\":\"epoch\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"kind\",\"type\":\"uint8\",\"indexed\":true,\"internalType\":\"enumFaultKind\"},{\"name\":\"evidence\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"prover\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleAdminChanged\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"previousAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"newAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleGranted\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleRevoked\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SignerActivated\",\"inputs\":[{\"name\":\"signer\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"fromEpoch\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SignerCancelled\",\"inputs\":[{\"name\":\"signer\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"fromEpoch\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SignerScheduled\",\"inputs\":[{\"name\":\"signer\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"fromEpoch\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"activatesAt\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AccessControlBadConfirmation\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AccessControlEnforcedDefaultAdminDelay\",\"inputs\":[{\"name\":\"schedule\",\"type\":\"uint48\",\"internalType\":\"uint48\"}]},{\"type\":\"error\",\"name\":\"AccessControlEnforcedDefaultAdminRules\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AccessControlInvalidDefaultAdmin\",\"inputs\":[{\"name\":\"defaultAdmin\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"AccessControlUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"neededRole\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"AlreadyProven\",\"inputs\":[{\"name\":\"evidence\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"DifferentEpochs\",\"inputs\":[{\"name\":\"first\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"second\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"EpochNotAnchored\",\"inputs\":[{\"name\":\"epoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"FromEpochTooLow\",\"inputs\":[{\"name\":\"fromEpoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"least\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InvalidShortString\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NoConflict\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NoSignerChangePending\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotALedger\",\"inputs\":[{\"name\":\"ledger\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"NotSignedBySigner\",\"inputs\":[{\"name\":\"epoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"signer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"recovered\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ProofHolds\",\"inputs\":[{\"name\":\"epoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"SafeCastOverflowedUintDowncast\",\"inputs\":[{\"name\":\"bits\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"SameSigner\",\"inputs\":[{\"name\":\"signer\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"SiblingsMismatch\",\"inputs\":[{\"name\":\"signed\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"given\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"SignerChangePending\",\"inputs\":[{\"name\":\"signer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"fromEpoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"SignerChangeStale\",\"inputs\":[{\"name\":\"fromEpoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"nextEpoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"SignerChangeTooEarly\",\"inputs\":[{\"name\":\"activatesAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"blockTime\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"StringTooLong\",\"inputs\":[{\"name\":\"str\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"WrongLedger\",\"inputs\":[{\"name\":\"expected\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"given\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ZeroDelay\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroSigner\",\"inputs\":[]}]",
	ID:  "StatementCourt",
	Bin: "0x6101a0604052348015610010575f5ffd5b50604051613b17380380613b1783398101604081905261002f91610540565b604080518082018252601781527f4b696e657469782053746174656d656e7420436f757274000000000000000000602080830191909152825180840190935260018352603160f81b908301529085876001600160a01b0381166100ac57604051636116401160e11b81525f60048201526024015b60405180910390fd5b600180546001600160d01b0316600160d01b65ffffffffffff8516021790556100d55f826102ff565b5050506100e78261036e60201b60201c565b610120526100f48161036e565b61014052815160208084019190912060e052815190820120610100524660a05261018060e05161010051604080517f8b73c3c69bb8fe3d512ecc4cf759cc79239f7b179b0ffacaa9a75d522b39400f60208201529081019290925260608201524660808201523060a08201525f9060c00160405160208183030381529060405280519060200120905090565b60805250503060c05261019a83632cd7bf0360e21b6103ab565b6101c257604051632252da5960e11b81526001600160a01b03841660048201526024016100a3565b6001600160a01b0382166101e95760405163e5c48ac560e01b815260040160405180910390fd5b806001600160401b03165f0361021257604051631a615ced60e11b815260040160405180910390fd5b6001600160a01b038084166101608190526001600160401b0383166101805260408051808201825292851683528051635b8735f360e11b8152905160059392602080850193919263b70e6be69260048082019392918290030181865afa15801561027e573d5f5f3e3d5ffd5b505050506040513d601f19601f820116820180604052508101906102a291906105b6565b6001600160401b0390811690915282546001810184555f9384526020938490208351910180549490930151909116600160a01b026001600160e01b03199093166001600160a01b03909116179190911790555061062a9350505050565b5f8261035b575f6103186002546001600160a01b031690565b6001600160a01b03161461033f57604051631fe1e13d60e11b815260040160405180910390fd5b600280546001600160a01b0319166001600160a01b0384161790555b61036583836103c6565b90505b92915050565b5f5f829050601f81511115610398578260405163305a27a960e01b81526004016100a391906105cf565b80516103a382610604565b179392505050565b5f6103b58361046d565b8015610365575061036583836104bb565b5f828152602081815260408083206001600160a01b038516845290915281205460ff16610466575f838152602081815260408083206001600160a01b03861684529091529020805460ff1916600117905561041e3390565b6001600160a01b0316826001600160a01b0316847f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a4506001610368565b505f610368565b5f61047f826301ffc9a760e01b6104bb565b156104af575f80610498846001600160e01b03196104df565b915091508180156104a7575080155b949350505050565b505f919050565b919050565b5f80806104c885856104df565b915091508180156104d65750805b95945050505050565b6301ffc9a760e01b5f818152600483905290819060208260248188617530fa92505f511515601f3d11169150509250929050565b6001600160a01b0381168114610527575f5ffd5b50565b80516001600160401b03811681146104b6575f5ffd5b5f5f5f5f5f60a08688031215610554575f5ffd5b855161055f81610513565b602087015190955065ffffffffffff8116811461057a575f5ffd5b604087015190945061058b81610513565b606087015190935061059c81610513565b91506105aa6080870161052a565b90509295509295909350565b5f602082840312156105c6575f5ffd5b6103658261052a565b602081525f82518060208401528060208501604085015e5f604082850101526040601f19601f83011684010191505092915050565b80516020808301519190811015610624575f198160200360031b1b821691505b50919050565b60805160a05160c05160e05161010051610120516101405161016051610180516134506106c75f395f81816103cb0152610c2d01525f818161045a015281816114190152818161187001528181611ba001528181611c8d01528181611cca01528181611d3d0152611de801525f611a8001525f611a5501525f61285e01525f61283601525f61279101525f6127bb01525f6127e501526134505ff3fe608060405234801561000f575f5ffd5b506004361061024f575f3560e01c806384ef8ffc1161013d578063a663d04a116100b8578063cc8463c811610088578063cf6eefb71161006e578063cf6eefb7146105fd578063d547741f1461062b578063d602b9fd1461063e575f5ffd5b8063cc8463c8146105ed578063cefc1429146105f5575f5ffd5b8063a663d04a1461059f578063b745c87b146105a7578063c1c23c65146105c7578063ca9a360d146105da575f5ffd5b8063941c6d0c1161010d578063a1eda53c116100f3578063a1eda53c1461054f578063a217fddf14610576578063a44256c61461057d575f5ffd5b8063941c6d0c1461053457806397326b4a14610547575f5ffd5b806384ef8ffc146104d0578063854209e8146104e15780638da5cb5b146104f657806391d14854146104fe575f5ffd5b8063357c0fb0116101cd57806356397c351161019d578063634e93da11610183578063634e93da1461048f578063649a5ec7146104a257806384b0196e146104b5575f5ffd5b806356397c35146104555780635e8a56041461047c575f5ffd5b8063357c0fb0146103c657806336568abe1461040657806336772c821461041957806351ba300e14610442575f5ffd5b8063181a36a41161022257806322a3808f1161020857806322a3808f14610370578063248a9ca3146103835780632f2ff15d146103b3575f5ffd5b8063181a36a41461033d57806321152eb014610368575f5ffd5b806301ffc9a714610253578063022d63fb1461027b5780630aa6220b1461029757806314d90c45146102a1575b5f5ffd5b610266610261366004612b47565b610646565b60405190151581526020015b60405180910390f35b620697805b60405165ffffffffffff9091168152602001610272565b61029f6106a1565b005b610300604080516060810182525f808252602082018190529181019190915250604080516060810182526006546001600160a01b038116825267ffffffffffffffff600160a01b90910481166020830152600754169181019190915290565b6040805182516001600160a01b0316815260208084015167ffffffffffffffff908116918301919091529282015190921690820152606001610272565b61035061034b366004612be1565b6106b6565b6040516001600160a01b039091168152602001610272565b61029f6106d4565b61035061037e366004612c46565b6107f8565b6103a5610391366004612c61565b5f9081526020819052604090206001015490565b604051908152602001610272565b61029f6103c1366004612c93565b610877565b6103ed7f000000000000000000000000000000000000000000000000000000000000000081565b60405167ffffffffffffffff9091168152602001610272565b61029f610414366004612c93565b6108bc565b6103a5610427366004612c46565b67ffffffffffffffff165f9081526008602052604090205490565b6103a5610450366004612cfe565b61099b565b6103507f000000000000000000000000000000000000000000000000000000000000000081565b61029f61048a366004612d3d565b6109f6565b61029f61049d366004612d72565b610d31565b61029f6104b0366004612d8b565b610d44565b6104bd610d57565b6040516102729796959493929190612dfc565b6002546001600160a01b0316610350565b6104e9610db5565b6040516102729190612eae565b610350610e30565b61026661050c366004612c93565b5f918252602082815260408084206001600160a01b0393909316845291905290205460ff1690565b6103a5610542366004612f04565b610e48565b61029f610e68565b61055761106f565b6040805165ffffffffffff938416815292909116602083015201610272565b6103a55f81565b61026661058b366004612c61565b5f9081526009602052604090205460ff1690565b600a546103a5565b6105ba6105b5366004612c46565b6110d8565b6040516102729190612f4b565b6103a56105d5366004612ff2565b6111b2565b6103a56105e8366004613087565b611362565b610280611556565b61029f6115e2565b600154604080516001600160a01b0383168152600160a01b90920465ffffffffffff16602083015201610272565b61029f610639366004612c93565b611631565b61029f611672565b5f7fffffffff0000000000000000000000000000000000000000000000000000000082167f6d27e6e200000000000000000000000000000000000000000000000000000000148061069b575061069b82611684565b92915050565b5f6106ab816116d9565b6106b36116e3565b50565b5f6106ca6106c385610e48565b84846116ef565b90505b9392505050565b5f6106de816116d9565b604080516060810182526006546001600160a01b038116808352600160a01b90910467ffffffffffffffff9081166020840152600754169282019290925290610753576040517f16d48c8400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600680547fffffffff00000000000000000000000000000000000000000000000000000000169055600780547fffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000169055805160208083015160405167ffffffffffffffff90911681526001600160a01b03909216917f63af16b98e94bb6465101a4b86a469cd7669ca62b406a7e5013bb4a64c7b6a3e91015b60405180910390a25050565b6005545f905b801561086f575f6005610812600184613136565b8154811061082257610822613149565b5f918252602090912001805490915067ffffffffffffffff808616600160a01b909204161161085c57546001600160a01b03169392505050565b508061086781613176565b9150506107fe565b505f92915050565b816108ae576040517f3fc3c27a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6108b8828261173f565b5050565b811580156108d757506002546001600160a01b038281169116145b15610991576001546001600160a01b03811690600160a01b900465ffffffffffff168115158061090d575065ffffffffffff8116155b8061092057504265ffffffffffff821610155b15610966576040517f19ca5ebb00000000000000000000000000000000000000000000000000000000815265ffffffffffff821660048201526024015b60405180910390fd5b5050600180547fffffffffffff000000000000ffffffffffffffffffffffffffffffffffffffff1690555b6108b88282611769565b5f6106cd8383808060200260200160405190810160405280939291908181526020015f905b828210156109ec576109dd604083028601368190038101906131aa565b815260200190600101906109c0565b50505050506117ba565b5f610a00816116d9565b6001600160a01b038316610a40576040517fe5c48ac500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6006546001600160a01b031615610aa8576006546040517fbe9777300000000000000000000000000000000000000000000000000000000081526001600160a01b0382166004820152600160a01b90910467ffffffffffffffff16602482015260440161095d565b600580545f9190610abb90600190613136565b81548110610acb57610acb613149565b5f91825260209091200180549091506001600160a01b0390811690851603610b2a576040517f684a329e0000000000000000000000000000000000000000000000000000000081526001600160a01b038516600482015260240161095d565b80547fffffffffffffffffffffffffffffffffffffffffffffffff0000000000000001600160a01b90910467ffffffffffffffff1601610ba8576040517f2a9c655b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8085166004830152602482015260440161095d565b80545f90610bc690600160a01b900467ffffffffffffffff1661186c565b90508067ffffffffffffffff168467ffffffffffffffff161015610c2a576040517f2a9c655b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff80861660048301528216602482015260440161095d565b5f7f0000000000000000000000000000000000000000000000000000000000000000610c5542611921565b610c5f9190613225565b604080516060810182526001600160a01b03891680825267ffffffffffffffff8981166020808501829052918616938501849052600680547fffffffff00000000000000000000000000000000000000000000000000000000168417600160a01b8302179055600780547fffffffffffffffffffffffffffffffffffffffffffffffff000000000000000016851790558451908152908101929092529293507f0cf637f2ce896a14f33631ebd092d4ae9272fa33634d9a0800780b04e75d7afb910160405180910390a2505050505050565b5f610d3b816116d9565b6108b882611974565b5f610d4e816116d9565b6108b8826119df565b5f6060805f5f5f6060610d68611a4e565b610d70611a79565b604080515f808252602082019092527f0f000000000000000000000000000000000000000000000000000000000000009b939a50919850469750309650945092509050565b60606005805480602002602001604051908101604052809291908181526020015f905b82821015610e27575f84815260209081902060408051808201909152908401546001600160a01b0381168252600160a01b900467ffffffffffffffff1681830152825260019092019101610dd8565b50505050905090565b5f610e436002546001600160a01b031690565b905090565b5f61069b610e63610e5e36859003850185613245565b611aa4565b611b56565b604080516060810182526006546001600160a01b038116808352600160a01b90910467ffffffffffffffff9081166020840152600754169282019290925290610edd576040517f16d48c8400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f610ee742611921565b9050816040015167ffffffffffffffff168167ffffffffffffffff161015610f555760408083015190517f2a2d471000000000000000000000000000000000000000000000000000000000815267ffffffffffffffff9182166004820152908216602482015260440161095d565b610f628260200151611b9d565b60408051808201825283516001600160a01b0390811682526020808601805167ffffffffffffffff908116838601908152600580546001810182555f9190915295517f036b6384b5eca791c62761152d0c79bb0604c104a5fb6f4eb0703f3154bb3db0909601805491518316600160a01b027fffffffff00000000000000000000000000000000000000000000000000000000928316978716979097179690961790955560068054909516909455600780547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001690558651905194519490931684529116917fc2ca66597d24290c92605c30793ab14553b9d669cc5829b7fecb21807578cef391016107ec565b6002545f907a010000000000000000000000000000000000000000000000000000900465ffffffffffff1680151580156110b157504265ffffffffffff821610155b6110bc575f5f6110d0565b600254600160a01b900465ffffffffffff16815b915091509091565b67ffffffffffffffff81165f908152600860209081526040808320805482518185028101850190935280835260609492939192909184015b828210156111a7575f848152602090206040805160808101909152600384029091018054829060ff16600181111561114a5761114a612f1e565b600181111561115b5761115b612f1e565b81526001828101546020808401919091526002909301546001600160a01b0381166040840152600160a01b900467ffffffffffffffff1660609092019190915291835292019101611110565b505050509050919050565b5f6111c36040850160208601612c46565b67ffffffffffffffff166111dd6040890160208a01612c46565b67ffffffffffffffff1614611252576111fc6040880160208901612c46565b61120c6040860160208701612c46565b6040517fbfb240a200000000000000000000000000000000000000000000000000000000815267ffffffffffffffff92831660048201529116602482015260440161095d565b5f61125e888888611c82565b90505f61126c868686611c82565b90508082148061129f575061129d611289368b90038b018b613245565b61129836899003890189613245565b611fa4565b155b156112d6576040517f7412124e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f81831061130d57604080516020810184905290810184905260600160405160208183030381529060405280519060200120611338565b6040805160208101859052908101839052606001604051602081830303815290604052805190602001205b905061135561134d60408c0160208d01612c46565b600183611fd2565b9998505050505050505050565b5f5f61136f878787611c82565b90505f6113c28585808060200260200160405190810160405280939291908181526020015f905b828210156109ec576113b3604083028601368190038101906131aa565b81526020019060010190611396565b90508760a00135811461140e576040517f7e68dade00000000000000000000000000000000000000000000000000000000815260a089013560048201526024810182905260440161095d565b5f6001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001663c0566abe61144e60408c0160208d01612c46565b60408c013561146360808e0160608f01612c46565b8d60800160208101906114769190612c46565b8b8b6040518763ffffffff1660e01b815260040161149996959493929190613303565b602060405180830381865afa1580156114b4573d5f5f3e3d5ffd5b505050506040513d601f19601f820116820180604052508101906114d89190613386565b9050801561152f576114f060408a0160208b01612c46565b6040517f64212abd00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff909116600482015260240161095d565b61154961154260408b0160208c01612c46565b5f85611fd2565b5090979650505050505050565b6002545f907a010000000000000000000000000000000000000000000000000000900465ffffffffffff16801515801561159757504265ffffffffffff8216105b6115c9576001547a010000000000000000000000000000000000000000000000000000900465ffffffffffff166115dc565b600254600160a01b900465ffffffffffff165b91505090565b6001546001600160a01b0316338114611629576040517fc22c802200000000000000000000000000000000000000000000000000000000815233600482015260240161095d565b6106b36121e8565b81611668576040517f3fc3c27a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6108b882826122ae565b5f61167c816116d9565b6106b36122d2565b5f7fffffffff0000000000000000000000000000000000000000000000000000000082167f3149878600000000000000000000000000000000000000000000000000000000148061069b575061069b826122dc565b6106b38133612372565b6116ed5f5f6123dd565b565b5f5f5f5f6116fe878787612507565b919450925090505f82600381111561171857611718612f1e565b14158061172457508015155b15611734575f93505050506106cd565b509095945050505050565b5f82815260208190526040902060010154611759816116d9565b611763838361254e565b50505050565b6001600160a01b03811633146117ab576040517f6697b23200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6117b582826125e5565b505050565b604080515f808252602082019092525f5b835181101561185d57818482815181106117e7576117e7613149565b60200260200101515f015185838151811061180457611804613149565b60200260200101516020015160c01b604051602001611825939291906133a5565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905291506001016117cb565b50805160209091012092915050565b5f5f7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663aea0e78b6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156118ca573d5f5f3e3d5ffd5b505050506040513d601f19601f820116820180604052508101906118ee91906133e6565b905067ffffffffffffffff8116611906846001613225565b67ffffffffffffffff16111561069b576106cd836001613225565b5f67ffffffffffffffff82111561197057604080517f6dfcc65000000000000000000000000000000000000000000000000000000000815260048101919091526024810183905260440161095d565b5090565b5f61197d611556565b61198642612639565b6119909190613401565b905061199c8282612684565b60405165ffffffffffff821681526001600160a01b038316907f3377dc44241e779dd06afab5b788a35ca5f3b778836e2990bdb26a2a4b2e5ed6906020016107ec565b5f6119e982612701565b6119f242612639565b6119fc9190613401565b9050611a0882826123dd565b6040805165ffffffffffff8085168252831660208201527ff1038c18cf84a56e432fdbfaf746924b7ea511dfe03a6506a0ceba4888788d9b910160405180910390a15050565b6060610e437f0000000000000000000000000000000000000000000000000000000000000000612748565b6060610e437f0000000000000000000000000000000000000000000000000000000000000000612748565b5f7f04c8b72ccd2c82dcda824abdea614e146a2f3142f70ad3acc2d3c4eff4c0bf1c825f015183602001518460400151856060015186608001518760a00151604051602001611b3997969594939291909687526001600160a01b0395909516602087015267ffffffffffffffff93841660408701526060860192909252821660808501521660a083015260c082015260e00190565b604051602081830303815290604052805190602001209050919050565b5f61069b611b62612785565b836040517f19010000000000000000000000000000000000000000000000000000000000008152600281019290925260228201526042902090565b5f7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663aea0e78b6040518163ffffffff1660e01b8152600401602060405180830381865afa158015611bfa573d5f5f3e3d5ffd5b505050506040513d601f19601f82011682018060405250810190611c1e91906133e6565b90508067ffffffffffffffff168267ffffffffffffffff1610156108b8576040517f6d5c81aa00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff80841660048301528216602482015260440161095d565b5f6001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016611cba6020860186612d72565b6001600160a01b031614611d3b577f0000000000000000000000000000000000000000000000000000000000000000611cf66020860186612d72565b6040517f185f957c0000000000000000000000000000000000000000000000000000000081526001600160a01b0392831660048201529116602482015260440161095d565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663b70e6be66040518163ffffffff1660e01b8152600401602060405180830381865afa158015611d97573d5f5f3e3d5ffd5b505050506040513d601f19601f82011682018060405250810190611dbb91906133e6565b67ffffffffffffffff16611dd56040860160208701612c46565b67ffffffffffffffff161080611e8d57507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663aea0e78b6040518163ffffffff1660e01b8152600401602060405180830381865afa158015611e42573d5f5f3e3d5ffd5b505050506040513d601f19601f82011682018060405250810190611e6691906133e6565b67ffffffffffffffff16611e806040860160208701612c46565b67ffffffffffffffff1610155b15611ee157611ea26040850160208601612c46565b6040517ffdd66f0900000000000000000000000000000000000000000000000000000000815267ffffffffffffffff909116600482015260240161095d565b5f611eeb85610e48565b90505f611f0161037e6040880160208901612c46565b90505f611f0f8387876116ef565b90506001600160a01b0381161580611f395750816001600160a01b0316816001600160a01b031614155b1561173457611f4e6040880160208901612c46565b6040517fd7eecde100000000000000000000000000000000000000000000000000000000815267ffffffffffffffff90911660048201526001600160a01b0380841660248301528216604482015260640161095d565b5f8160400151836040015114806106cd57505060809081015191015167ffffffffffffffff90811691161490565b5f8181526009602052604090205460ff161561201d576040517fea0336840000000000000000000000000000000000000000000000000000000081526004810182905260240161095d565b5f81815260096020908152604080832080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016600190811790915567ffffffffffffffff8716845260089092529182902082516080810190935291908190859081111561208d5761208d612f1e565b8152602081018490523360408201526060016120a842611921565b67ffffffffffffffff1690528154600181810184555f9384526020909320825160039092020180549293909283917fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0090911690838181111561210c5761210c612f1e565b0217905550602082015160018083019190915560408301516002909201805460609094015167ffffffffffffffff16600160a01b027fffffffff000000000000000000000000000000000000000000000000000000009094166001600160a01b0390931692909217929092179055600a80545f9061218b90849061341f565b909155508190508260018111156121a4576121a4612f1e565b60405133815267ffffffffffffffff8616907f2f7afa658d77d8c7174d53beea0e7ae9234846d771751668d5aa6c737a30bc6a9060200160405180910390a4505050565b6001546001600160a01b03811690600160a01b900465ffffffffffff1680158061221a57504265ffffffffffff821610155b1561225b576040517f19ca5ebb00000000000000000000000000000000000000000000000000000000815265ffffffffffff8216600482015260240161095d565b6122765f6122716002546001600160a01b031690565b6125e5565b506122815f8361254e565b5050600180547fffffffffffff000000000000000000000000000000000000000000000000000016905550565b5f828152602081905260409020600101546122c8816116d9565b61176383836125e5565b6116ed5f5f612684565b5f7fffffffff0000000000000000000000000000000000000000000000000000000082167f7965db0b00000000000000000000000000000000000000000000000000000000148061069b57507f01ffc9a7000000000000000000000000000000000000000000000000000000007fffffffff0000000000000000000000000000000000000000000000000000000083161461069b565b5f828152602081815260408083206001600160a01b038516845290915290205460ff166108b8576040517fe2517d3f0000000000000000000000000000000000000000000000000000000081526001600160a01b03821660048201526024810183905260440161095d565b6002547a010000000000000000000000000000000000000000000000000000900465ffffffffffff1680156124a0574265ffffffffffff82161015612477576002546001805479ffffffffffffffffffffffffffffffffffffffffffffffffffff16600160a01b90920465ffffffffffff167a010000000000000000000000000000000000000000000000000000029190911790556124a0565b6040517f2b1fa2edafe6f7b9e97c1a9e0c3660e645beb2dcaa2d45bdbf9beaf5472e1ec5905f90a15b50600280546001600160a01b0316600160a01b65ffffffffffff9485160279ffffffffffffffffffffffffffffffffffffffffffffffffffff16177a0100000000000000000000000000000000000000000000000000009290931691909102919091179055565b5f8080604184900361253b578435602086013560408701355f1a61252d898285856128ae565b955095509550505050612545565b505f915060029050825b93509350939050565b5f826125db575f6125676002546001600160a01b031690565b6001600160a01b0316146125a7576040517f3fc3c27a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600280547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0384161790555b6106cd8383612994565b5f8215801561260157506002546001600160a01b038381169116145b1561262f57600280547fffffffffffffffffffffffff00000000000000000000000000000000000000001690555b6106cd8383612a59565b5f65ffffffffffff821115611970576040517f6dfcc650000000000000000000000000000000000000000000000000000000008152603060048201526024810183905260440161095d565b60018054600160a01b65ffffffffffff84811682027fffffffffffff000000000000000000000000000000000000000000000000000084166001600160a01b038816171790935590041680156117b5576040517f8886ebfc4259abdbc16601dd8fb5678e54878f47b3c34836cfc51154a9605109905f90a1505050565b5f5f61270b611556565b90508065ffffffffffff168365ffffffffffff16116127335761272e8382613432565b6106cd565b6106cd65ffffffffffff841662069780612af8565b60605f61275483612b07565b6040805160208082528183019092529192505f91906020820181803683375050509182525060208101929092525090565b5f306001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000161480156127dd57507f000000000000000000000000000000000000000000000000000000000000000046145b1561280757507f000000000000000000000000000000000000000000000000000000000000000090565b610e43604080517f8b73c3c69bb8fe3d512ecc4cf759cc79239f7b179b0ffacaa9a75d522b39400f60208201527f0000000000000000000000000000000000000000000000000000000000000000918101919091527f000000000000000000000000000000000000000000000000000000000000000060608201524660808201523060a08201525f9060c00160405160208183030381529060405280519060200120905090565b5f80807f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a08411156128e757505f9150600390508261298a565b604080515f808252602082018084528a905260ff891692820192909252606081018790526080810186905260019060a0016020604051602081039080840390855afa158015612938573d5f5f3e3d5ffd5b50506040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe001519150506001600160a01b03811661298157505f92506001915082905061298a565b92505f91508190505b9450945094915050565b5f828152602081815260408083206001600160a01b038516845290915281205460ff16612a52575f838152602081815260408083206001600160a01b0386168452909152902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055612a0a3390565b6001600160a01b0316826001600160a01b0316847f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a450600161069b565b505f61069b565b5f828152602081815260408083206001600160a01b038516845290915281205460ff1615612a52575f838152602081815260408083206001600160a01b038616808552925280832080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016905551339286917ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9190a450600161069b565b5f8282188284100282186106cd565b5f60ff8216601f81111561069b576040517fb3512b0c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f60208284031215612b57575f5ffd5b81357fffffffff00000000000000000000000000000000000000000000000000000000811681146106cd575f5ffd5b5f60c08284031215612b96575f5ffd5b50919050565b5f5f83601f840112612bac575f5ffd5b50813567ffffffffffffffff811115612bc3575f5ffd5b602083019150836020828501011115612bda575f5ffd5b9250929050565b5f5f5f60e08486031215612bf3575f5ffd5b612bfd8585612b86565b925060c084013567ffffffffffffffff811115612c18575f5ffd5b612c2486828701612b9c565b9497909650939450505050565b67ffffffffffffffff811681146106b3575f5ffd5b5f60208284031215612c56575f5ffd5b81356106cd81612c31565b5f60208284031215612c71575f5ffd5b5035919050565b80356001600160a01b0381168114612c8e575f5ffd5b919050565b5f5f60408385031215612ca4575f5ffd5b82359150612cb460208401612c78565b90509250929050565b5f5f83601f840112612ccd575f5ffd5b50813567ffffffffffffffff811115612ce4575f5ffd5b6020830191508360208260061b8501011115612bda575f5ffd5b5f5f60208385031215612d0f575f5ffd5b823567ffffffffffffffff811115612d25575f5ffd5b612d3185828601612cbd565b90969095509350505050565b5f5f60408385031215612d4e575f5ffd5b612d5783612c78565b91506020830135612d6781612c31565b809150509250929050565b5f60208284031215612d82575f5ffd5b6106cd82612c78565b5f60208284031215612d9b575f5ffd5b813565ffffffffffff811681146106cd575f5ffd5b5f81518084528060208401602086015e5f6020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b7fff000000000000000000000000000000000000000000000000000000000000008816815260e060208201525f612e3660e0830189612db0565b8281036040840152612e488189612db0565b606084018890526001600160a01b038716608085015260a0840186905283810360c0850152845180825260208087019350909101905f5b81811015612e9d578351835260209384019390920191600101612e7f565b50909b9a5050505050505050505050565b602080825282518282018190525f918401906040840190835b8181101561173457835180516001600160a01b0316845260209081015167ffffffffffffffff168185015290930192604090920191600101612ec7565b5f60c08284031215612f14575f5ffd5b6106cd8383612b86565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602160045260245ffd5b602080825282518282018190525f918401906040840190835b81811015611734578351805160028110612fa5577f4e487b71000000000000000000000000000000000000000000000000000000005f52602160045260245ffd5b80855250602081015160208501526001600160a01b03604082015116604085015267ffffffffffffffff606082015116606085015250608083019250602084019350600181019050612f64565b5f5f5f5f5f5f6101c08789031215613008575f5ffd5b6130128888612b86565b955060c087013567ffffffffffffffff81111561302d575f5ffd5b61303989828a01612b9c565b909650945061304d90508860e08901612b86565b92506101a087013567ffffffffffffffff811115613069575f5ffd5b61307589828a01612b9c565b979a9699509497509295939492505050565b5f5f5f5f5f610100868803121561309c575f5ffd5b6130a68787612b86565b945060c086013567ffffffffffffffff8111156130c1575f5ffd5b6130cd88828901612b9c565b90955093505060e086013567ffffffffffffffff8111156130ec575f5ffd5b6130f888828901612cbd565b969995985093965092949392505050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b8181038181111561069b5761069b613109565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b5f8161318457613184613109565b507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0190565b5f60408284031280156131bb575f5ffd5b506040805190810167ffffffffffffffff81118282101715613204577f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b60405282358152602083013561321981612c31565b60208201529392505050565b67ffffffffffffffff818116838216019081111561069b5761069b613109565b5f60c0828403128015613256575f5ffd5b5060405160c0810167ffffffffffffffff8111828210171561329f577f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b6040526132ab83612c78565b815260208301356132bb81612c31565b60208201526040838101359082015260608301356132d881612c31565b606082015260808301356132eb81612c31565b608082015260a0928301359281019290925250919050565b67ffffffffffffffff87811682526020820187905285811660408301528416606082015260a06080820181905281018290525f8360c08301825b858110156133785782358252602083013561335781612c31565b67ffffffffffffffff1660208301526040928301929091019060010161333d565b509998505050505050505050565b5f60208284031215613396575f5ffd5b815180151581146106cd575f5ffd5b5f84518060208701845e91909101928352507fffffffffffffffff000000000000000000000000000000000000000000000000166020820152602801919050565b5f602082840312156133f6575f5ffd5b81516106cd81612c31565b65ffffffffffff818116838216019081111561069b5761069b613109565b8082018082111561069b5761069b613109565b65ffffffffffff828116828216039081111561069b5761069b61310956",
}

// StatementCourt is an auto generated Go binding around an Ethereum contract.
type StatementCourt struct {
	abi abi.ABI
}

// GetABI returns the ABI associated with this contract binding.
func (c *StatementCourt) GetABI() abi.ABI {
	return c.abi
}

// NewStatementCourt creates a new instance of StatementCourt.
func NewStatementCourt() *StatementCourt {
	parsed, err := StatementCourtMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &StatementCourt{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *StatementCourt) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackConstructor is the Go binding used to pack the parameters required for
// contract deployment.
//
// Solidity: constructor(address admin, uint48 adminTransferDelay, address ledger_, address initialSigner, uint64 signerChangeDelay_) returns()
func (statementCourt *StatementCourt) PackConstructor(admin common.Address, adminTransferDelay *big.Int, ledger_ common.Address, initialSigner common.Address, signerChangeDelay_ uint64) []byte {
	enc, err := statementCourt.abi.Pack("", admin, adminTransferDelay, ledger_, initialSigner, signerChangeDelay_)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackDEFAULTADMINROLE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa217fddf.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (statementCourt *StatementCourt) PackDEFAULTADMINROLE() []byte {
	enc, err := statementCourt.abi.Pack("DEFAULT_ADMIN_ROLE")
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
func (statementCourt *StatementCourt) TryPackDEFAULTADMINROLE() ([]byte, error) {
	return statementCourt.abi.Pack("DEFAULT_ADMIN_ROLE")
}

// UnpackDEFAULTADMINROLE is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (statementCourt *StatementCourt) UnpackDEFAULTADMINROLE(data []byte) ([32]byte, error) {
	out, err := statementCourt.abi.Unpack("DEFAULT_ADMIN_ROLE", data)
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
func (statementCourt *StatementCourt) PackAcceptDefaultAdminTransfer() []byte {
	enc, err := statementCourt.abi.Pack("acceptDefaultAdminTransfer")
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
func (statementCourt *StatementCourt) TryPackAcceptDefaultAdminTransfer() ([]byte, error) {
	return statementCourt.abi.Pack("acceptDefaultAdminTransfer")
}

// PackActivateSigner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x97326b4a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function activateSigner() returns()
func (statementCourt *StatementCourt) PackActivateSigner() []byte {
	enc, err := statementCourt.abi.Pack("activateSigner")
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
func (statementCourt *StatementCourt) TryPackActivateSigner() ([]byte, error) {
	return statementCourt.abi.Pack("activateSigner")
}

// PackBeginDefaultAdminTransfer is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x634e93da.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function beginDefaultAdminTransfer(address newAdmin) returns()
func (statementCourt *StatementCourt) PackBeginDefaultAdminTransfer(newAdmin common.Address) []byte {
	enc, err := statementCourt.abi.Pack("beginDefaultAdminTransfer", newAdmin)
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
func (statementCourt *StatementCourt) TryPackBeginDefaultAdminTransfer(newAdmin common.Address) ([]byte, error) {
	return statementCourt.abi.Pack("beginDefaultAdminTransfer", newAdmin)
}

// PackCancelDefaultAdminTransfer is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd602b9fd.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function cancelDefaultAdminTransfer() returns()
func (statementCourt *StatementCourt) PackCancelDefaultAdminTransfer() []byte {
	enc, err := statementCourt.abi.Pack("cancelDefaultAdminTransfer")
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
func (statementCourt *StatementCourt) TryPackCancelDefaultAdminTransfer() ([]byte, error) {
	return statementCourt.abi.Pack("cancelDefaultAdminTransfer")
}

// PackCancelSigner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x21152eb0.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function cancelSigner() returns()
func (statementCourt *StatementCourt) PackCancelSigner() []byte {
	enc, err := statementCourt.abi.Pack("cancelSigner")
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
func (statementCourt *StatementCourt) TryPackCancelSigner() ([]byte, error) {
	return statementCourt.abi.Pack("cancelSigner")
}

// PackChangeDefaultAdminDelay is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x649a5ec7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function changeDefaultAdminDelay(uint48 newDelay) returns()
func (statementCourt *StatementCourt) PackChangeDefaultAdminDelay(newDelay *big.Int) []byte {
	enc, err := statementCourt.abi.Pack("changeDefaultAdminDelay", newDelay)
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
func (statementCourt *StatementCourt) TryPackChangeDefaultAdminDelay(newDelay *big.Int) ([]byte, error) {
	return statementCourt.abi.Pack("changeDefaultAdminDelay", newDelay)
}

// PackDefaultAdmin is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x84ef8ffc.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function defaultAdmin() view returns(address)
func (statementCourt *StatementCourt) PackDefaultAdmin() []byte {
	enc, err := statementCourt.abi.Pack("defaultAdmin")
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
func (statementCourt *StatementCourt) TryPackDefaultAdmin() ([]byte, error) {
	return statementCourt.abi.Pack("defaultAdmin")
}

// UnpackDefaultAdmin is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x84ef8ffc.
//
// Solidity: function defaultAdmin() view returns(address)
func (statementCourt *StatementCourt) UnpackDefaultAdmin(data []byte) (common.Address, error) {
	out, err := statementCourt.abi.Unpack("defaultAdmin", data)
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
func (statementCourt *StatementCourt) PackDefaultAdminDelay() []byte {
	enc, err := statementCourt.abi.Pack("defaultAdminDelay")
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
func (statementCourt *StatementCourt) TryPackDefaultAdminDelay() ([]byte, error) {
	return statementCourt.abi.Pack("defaultAdminDelay")
}

// UnpackDefaultAdminDelay is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xcc8463c8.
//
// Solidity: function defaultAdminDelay() view returns(uint48)
func (statementCourt *StatementCourt) UnpackDefaultAdminDelay(data []byte) (*big.Int, error) {
	out, err := statementCourt.abi.Unpack("defaultAdminDelay", data)
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
func (statementCourt *StatementCourt) PackDefaultAdminDelayIncreaseWait() []byte {
	enc, err := statementCourt.abi.Pack("defaultAdminDelayIncreaseWait")
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
func (statementCourt *StatementCourt) TryPackDefaultAdminDelayIncreaseWait() ([]byte, error) {
	return statementCourt.abi.Pack("defaultAdminDelayIncreaseWait")
}

// UnpackDefaultAdminDelayIncreaseWait is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x022d63fb.
//
// Solidity: function defaultAdminDelayIncreaseWait() view returns(uint48)
func (statementCourt *StatementCourt) UnpackDefaultAdminDelayIncreaseWait(data []byte) (*big.Int, error) {
	out, err := statementCourt.abi.Unpack("defaultAdminDelayIncreaseWait", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackEip712Domain is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x84b0196e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (statementCourt *StatementCourt) PackEip712Domain() []byte {
	enc, err := statementCourt.abi.Pack("eip712Domain")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackEip712Domain is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x84b0196e.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (statementCourt *StatementCourt) TryPackEip712Domain() ([]byte, error) {
	return statementCourt.abi.Pack("eip712Domain")
}

// Eip712DomainOutput serves as a container for the return parameters of contract
// method Eip712Domain.
type Eip712DomainOutput struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}

// UnpackEip712Domain is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (statementCourt *StatementCourt) UnpackEip712Domain(data []byte) (Eip712DomainOutput, error) {
	out, err := statementCourt.abi.Unpack("eip712Domain", data)
	outstruct := new(Eip712DomainOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Fields = *abi.ConvertType(out[0], new([1]byte)).(*[1]byte)
	outstruct.Name = *abi.ConvertType(out[1], new(string)).(*string)
	outstruct.Version = *abi.ConvertType(out[2], new(string)).(*string)
	outstruct.ChainId = abi.ConvertType(out[3], new(big.Int)).(*big.Int)
	outstruct.VerifyingContract = *abi.ConvertType(out[4], new(common.Address)).(*common.Address)
	outstruct.Salt = *abi.ConvertType(out[5], new([32]byte)).(*[32]byte)
	outstruct.Extensions = *abi.ConvertType(out[6], new([]*big.Int)).(*[]*big.Int)
	return *outstruct, nil
}

// PackFaultCount is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x36772c82.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function faultCount(uint64 epoch) view returns(uint256)
func (statementCourt *StatementCourt) PackFaultCount(epoch uint64) []byte {
	enc, err := statementCourt.abi.Pack("faultCount", epoch)
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
func (statementCourt *StatementCourt) TryPackFaultCount(epoch uint64) ([]byte, error) {
	return statementCourt.abi.Pack("faultCount", epoch)
}

// UnpackFaultCount is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x36772c82.
//
// Solidity: function faultCount(uint64 epoch) view returns(uint256)
func (statementCourt *StatementCourt) UnpackFaultCount(data []byte) (*big.Int, error) {
	out, err := statementCourt.abi.Unpack("faultCount", data)
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
func (statementCourt *StatementCourt) PackFaultsOf(epoch uint64) []byte {
	enc, err := statementCourt.abi.Pack("faultsOf", epoch)
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
func (statementCourt *StatementCourt) TryPackFaultsOf(epoch uint64) ([]byte, error) {
	return statementCourt.abi.Pack("faultsOf", epoch)
}

// UnpackFaultsOf is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb745c87b.
//
// Solidity: function faultsOf(uint64 epoch) view returns((uint8,bytes32,address,uint64)[])
func (statementCourt *StatementCourt) UnpackFaultsOf(data []byte) ([]Fault, error) {
	out, err := statementCourt.abi.Unpack("faultsOf", data)
	if err != nil {
		return *new([]Fault), err
	}
	out0 := *abi.ConvertType(out[0], new([]Fault)).(*[]Fault)
	return out0, nil
}

// PackGetRoleAdmin is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x248a9ca3.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (statementCourt *StatementCourt) PackGetRoleAdmin(role [32]byte) []byte {
	enc, err := statementCourt.abi.Pack("getRoleAdmin", role)
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
func (statementCourt *StatementCourt) TryPackGetRoleAdmin(role [32]byte) ([]byte, error) {
	return statementCourt.abi.Pack("getRoleAdmin", role)
}

// UnpackGetRoleAdmin is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (statementCourt *StatementCourt) UnpackGetRoleAdmin(data []byte) ([32]byte, error) {
	out, err := statementCourt.abi.Unpack("getRoleAdmin", data)
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
func (statementCourt *StatementCourt) PackGrantRole(role [32]byte, account common.Address) []byte {
	enc, err := statementCourt.abi.Pack("grantRole", role, account)
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
func (statementCourt *StatementCourt) TryPackGrantRole(role [32]byte, account common.Address) ([]byte, error) {
	return statementCourt.abi.Pack("grantRole", role, account)
}

// PackHasRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x91d14854.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (statementCourt *StatementCourt) PackHasRole(role [32]byte, account common.Address) []byte {
	enc, err := statementCourt.abi.Pack("hasRole", role, account)
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
func (statementCourt *StatementCourt) TryPackHasRole(role [32]byte, account common.Address) ([]byte, error) {
	return statementCourt.abi.Pack("hasRole", role, account)
}

// UnpackHasRole is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (statementCourt *StatementCourt) UnpackHasRole(data []byte) (bool, error) {
	out, err := statementCourt.abi.Unpack("hasRole", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackIsProven is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa44256c6.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function isProven(bytes32 evidence) view returns(bool)
func (statementCourt *StatementCourt) PackIsProven(evidence [32]byte) []byte {
	enc, err := statementCourt.abi.Pack("isProven", evidence)
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
func (statementCourt *StatementCourt) TryPackIsProven(evidence [32]byte) ([]byte, error) {
	return statementCourt.abi.Pack("isProven", evidence)
}

// UnpackIsProven is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa44256c6.
//
// Solidity: function isProven(bytes32 evidence) view returns(bool)
func (statementCourt *StatementCourt) UnpackIsProven(data []byte) (bool, error) {
	out, err := statementCourt.abi.Unpack("isProven", data)
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
func (statementCourt *StatementCourt) PackLedger() []byte {
	enc, err := statementCourt.abi.Pack("ledger")
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
func (statementCourt *StatementCourt) TryPackLedger() ([]byte, error) {
	return statementCourt.abi.Pack("ledger")
}

// UnpackLedger is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x56397c35.
//
// Solidity: function ledger() view returns(address)
func (statementCourt *StatementCourt) UnpackLedger(data []byte) (common.Address, error) {
	out, err := statementCourt.abi.Unpack("ledger", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackOwner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8da5cb5b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function owner() view returns(address)
func (statementCourt *StatementCourt) PackOwner() []byte {
	enc, err := statementCourt.abi.Pack("owner")
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
func (statementCourt *StatementCourt) TryPackOwner() ([]byte, error) {
	return statementCourt.abi.Pack("owner")
}

// UnpackOwner is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (statementCourt *StatementCourt) UnpackOwner(data []byte) (common.Address, error) {
	out, err := statementCourt.abi.Unpack("owner", data)
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
func (statementCourt *StatementCourt) PackPendingDefaultAdmin() []byte {
	enc, err := statementCourt.abi.Pack("pendingDefaultAdmin")
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
func (statementCourt *StatementCourt) TryPackPendingDefaultAdmin() ([]byte, error) {
	return statementCourt.abi.Pack("pendingDefaultAdmin")
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
func (statementCourt *StatementCourt) UnpackPendingDefaultAdmin(data []byte) (PendingDefaultAdminOutput, error) {
	out, err := statementCourt.abi.Unpack("pendingDefaultAdmin", data)
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
func (statementCourt *StatementCourt) PackPendingDefaultAdminDelay() []byte {
	enc, err := statementCourt.abi.Pack("pendingDefaultAdminDelay")
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
func (statementCourt *StatementCourt) TryPackPendingDefaultAdminDelay() ([]byte, error) {
	return statementCourt.abi.Pack("pendingDefaultAdminDelay")
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
func (statementCourt *StatementCourt) UnpackPendingDefaultAdminDelay(data []byte) (PendingDefaultAdminDelayOutput, error) {
	out, err := statementCourt.abi.Unpack("pendingDefaultAdminDelay", data)
	outstruct := new(PendingDefaultAdminDelayOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.NewDelay = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.Schedule = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	return *outstruct, nil
}

// PackPendingSigner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x14d90c45.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function pendingSigner() view returns((address,uint64,uint64))
func (statementCourt *StatementCourt) PackPendingSigner() []byte {
	enc, err := statementCourt.abi.Pack("pendingSigner")
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
func (statementCourt *StatementCourt) TryPackPendingSigner() ([]byte, error) {
	return statementCourt.abi.Pack("pendingSigner")
}

// UnpackPendingSigner is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x14d90c45.
//
// Solidity: function pendingSigner() view returns((address,uint64,uint64))
func (statementCourt *StatementCourt) UnpackPendingSigner(data []byte) (PendingSigner, error) {
	out, err := statementCourt.abi.Unpack("pendingSigner", data)
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
func (statementCourt *StatementCourt) PackProveBadProof(statement EpochStatement, signature []byte, siblings []SumNode) []byte {
	enc, err := statementCourt.abi.Pack("proveBadProof", statement, signature, siblings)
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
func (statementCourt *StatementCourt) TryPackProveBadProof(statement EpochStatement, signature []byte, siblings []SumNode) ([]byte, error) {
	return statementCourt.abi.Pack("proveBadProof", statement, signature, siblings)
}

// UnpackProveBadProof is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xca9a360d.
//
// Solidity: function proveBadProof((address,uint64,bytes32,uint64,uint64,bytes32) statement, bytes signature, (bytes32,uint64)[] siblings) returns(bytes32)
func (statementCourt *StatementCourt) UnpackProveBadProof(data []byte) ([32]byte, error) {
	out, err := statementCourt.abi.Unpack("proveBadProof", data)
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
func (statementCourt *StatementCourt) PackProveEquivocation(first EpochStatement, firstSignature []byte, second EpochStatement, secondSignature []byte) []byte {
	enc, err := statementCourt.abi.Pack("proveEquivocation", first, firstSignature, second, secondSignature)
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
func (statementCourt *StatementCourt) TryPackProveEquivocation(first EpochStatement, firstSignature []byte, second EpochStatement, secondSignature []byte) ([]byte, error) {
	return statementCourt.abi.Pack("proveEquivocation", first, firstSignature, second, secondSignature)
}

// UnpackProveEquivocation is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc1c23c65.
//
// Solidity: function proveEquivocation((address,uint64,bytes32,uint64,uint64,bytes32) first, bytes firstSignature, (address,uint64,bytes32,uint64,uint64,bytes32) second, bytes secondSignature) returns(bytes32)
func (statementCourt *StatementCourt) UnpackProveEquivocation(data []byte) ([32]byte, error) {
	out, err := statementCourt.abi.Unpack("proveEquivocation", data)
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
func (statementCourt *StatementCourt) PackRecoverSigner(statement EpochStatement, signature []byte) []byte {
	enc, err := statementCourt.abi.Pack("recoverSigner", statement, signature)
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
func (statementCourt *StatementCourt) TryPackRecoverSigner(statement EpochStatement, signature []byte) ([]byte, error) {
	return statementCourt.abi.Pack("recoverSigner", statement, signature)
}

// UnpackRecoverSigner is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x181a36a4.
//
// Solidity: function recoverSigner((address,uint64,bytes32,uint64,uint64,bytes32) statement, bytes signature) view returns(address)
func (statementCourt *StatementCourt) UnpackRecoverSigner(data []byte) (common.Address, error) {
	out, err := statementCourt.abi.Unpack("recoverSigner", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackRenounceRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x36568abe.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (statementCourt *StatementCourt) PackRenounceRole(role [32]byte, account common.Address) []byte {
	enc, err := statementCourt.abi.Pack("renounceRole", role, account)
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
func (statementCourt *StatementCourt) TryPackRenounceRole(role [32]byte, account common.Address) ([]byte, error) {
	return statementCourt.abi.Pack("renounceRole", role, account)
}

// PackRevokeRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd547741f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (statementCourt *StatementCourt) PackRevokeRole(role [32]byte, account common.Address) []byte {
	enc, err := statementCourt.abi.Pack("revokeRole", role, account)
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
func (statementCourt *StatementCourt) TryPackRevokeRole(role [32]byte, account common.Address) ([]byte, error) {
	return statementCourt.abi.Pack("revokeRole", role, account)
}

// PackRollbackDefaultAdminDelay is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0aa6220b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function rollbackDefaultAdminDelay() returns()
func (statementCourt *StatementCourt) PackRollbackDefaultAdminDelay() []byte {
	enc, err := statementCourt.abi.Pack("rollbackDefaultAdminDelay")
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
func (statementCourt *StatementCourt) TryPackRollbackDefaultAdminDelay() ([]byte, error) {
	return statementCourt.abi.Pack("rollbackDefaultAdminDelay")
}

// PackScheduleSigner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5e8a5604.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function scheduleSigner(address signer, uint64 fromEpoch) returns()
func (statementCourt *StatementCourt) PackScheduleSigner(signer common.Address, fromEpoch uint64) []byte {
	enc, err := statementCourt.abi.Pack("scheduleSigner", signer, fromEpoch)
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
func (statementCourt *StatementCourt) TryPackScheduleSigner(signer common.Address, fromEpoch uint64) ([]byte, error) {
	return statementCourt.abi.Pack("scheduleSigner", signer, fromEpoch)
}

// PackSiblingsHash is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x51ba300e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function siblingsHash((bytes32,uint64)[] siblings) pure returns(bytes32)
func (statementCourt *StatementCourt) PackSiblingsHash(siblings []SumNode) []byte {
	enc, err := statementCourt.abi.Pack("siblingsHash", siblings)
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
func (statementCourt *StatementCourt) TryPackSiblingsHash(siblings []SumNode) ([]byte, error) {
	return statementCourt.abi.Pack("siblingsHash", siblings)
}

// UnpackSiblingsHash is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x51ba300e.
//
// Solidity: function siblingsHash((bytes32,uint64)[] siblings) pure returns(bytes32)
func (statementCourt *StatementCourt) UnpackSiblingsHash(data []byte) ([32]byte, error) {
	out, err := statementCourt.abi.Unpack("siblingsHash", data)
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
func (statementCourt *StatementCourt) PackSignerChangeDelay() []byte {
	enc, err := statementCourt.abi.Pack("signerChangeDelay")
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
func (statementCourt *StatementCourt) TryPackSignerChangeDelay() ([]byte, error) {
	return statementCourt.abi.Pack("signerChangeDelay")
}

// UnpackSignerChangeDelay is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x357c0fb0.
//
// Solidity: function signerChangeDelay() view returns(uint64)
func (statementCourt *StatementCourt) UnpackSignerChangeDelay(data []byte) (uint64, error) {
	out, err := statementCourt.abi.Unpack("signerChangeDelay", data)
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
func (statementCourt *StatementCourt) PackSignerOf(epoch uint64) []byte {
	enc, err := statementCourt.abi.Pack("signerOf", epoch)
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
func (statementCourt *StatementCourt) TryPackSignerOf(epoch uint64) ([]byte, error) {
	return statementCourt.abi.Pack("signerOf", epoch)
}

// UnpackSignerOf is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x22a3808f.
//
// Solidity: function signerOf(uint64 epoch) view returns(address)
func (statementCourt *StatementCourt) UnpackSignerOf(data []byte) (common.Address, error) {
	out, err := statementCourt.abi.Unpack("signerOf", data)
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
func (statementCourt *StatementCourt) PackSignerTerms() []byte {
	enc, err := statementCourt.abi.Pack("signerTerms")
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
func (statementCourt *StatementCourt) TryPackSignerTerms() ([]byte, error) {
	return statementCourt.abi.Pack("signerTerms")
}

// UnpackSignerTerms is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x854209e8.
//
// Solidity: function signerTerms() view returns((address,uint64)[])
func (statementCourt *StatementCourt) UnpackSignerTerms(data []byte) ([]SignerTerm, error) {
	out, err := statementCourt.abi.Unpack("signerTerms", data)
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
func (statementCourt *StatementCourt) PackStatementDigest(statement EpochStatement) []byte {
	enc, err := statementCourt.abi.Pack("statementDigest", statement)
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
func (statementCourt *StatementCourt) TryPackStatementDigest(statement EpochStatement) ([]byte, error) {
	return statementCourt.abi.Pack("statementDigest", statement)
}

// UnpackStatementDigest is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x941c6d0c.
//
// Solidity: function statementDigest((address,uint64,bytes32,uint64,uint64,bytes32) statement) view returns(bytes32)
func (statementCourt *StatementCourt) UnpackStatementDigest(data []byte) ([32]byte, error) {
	out, err := statementCourt.abi.Unpack("statementDigest", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackSupportsInterface is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x01ffc9a7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (statementCourt *StatementCourt) PackSupportsInterface(interfaceId [4]byte) []byte {
	enc, err := statementCourt.abi.Pack("supportsInterface", interfaceId)
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
func (statementCourt *StatementCourt) TryPackSupportsInterface(interfaceId [4]byte) ([]byte, error) {
	return statementCourt.abi.Pack("supportsInterface", interfaceId)
}

// UnpackSupportsInterface is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (statementCourt *StatementCourt) UnpackSupportsInterface(data []byte) (bool, error) {
	out, err := statementCourt.abi.Unpack("supportsInterface", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackTotalFaults is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa663d04a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function totalFaults() view returns(uint256)
func (statementCourt *StatementCourt) PackTotalFaults() []byte {
	enc, err := statementCourt.abi.Pack("totalFaults")
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
func (statementCourt *StatementCourt) TryPackTotalFaults() ([]byte, error) {
	return statementCourt.abi.Pack("totalFaults")
}

// UnpackTotalFaults is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa663d04a.
//
// Solidity: function totalFaults() view returns(uint256)
func (statementCourt *StatementCourt) UnpackTotalFaults(data []byte) (*big.Int, error) {
	out, err := statementCourt.abi.Unpack("totalFaults", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// StatementCourtDefaultAdminDelayChangeCanceled represents a DefaultAdminDelayChangeCanceled event raised by the StatementCourt contract.
type StatementCourtDefaultAdminDelayChangeCanceled struct {
	Raw *types.Log // Blockchain specific contextual infos
}

const StatementCourtDefaultAdminDelayChangeCanceledEventName = "DefaultAdminDelayChangeCanceled"

// ContractEventName returns the user-defined event name.
func (StatementCourtDefaultAdminDelayChangeCanceled) ContractEventName() string {
	return StatementCourtDefaultAdminDelayChangeCanceledEventName
}

// UnpackDefaultAdminDelayChangeCanceledEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event DefaultAdminDelayChangeCanceled()
func (statementCourt *StatementCourt) UnpackDefaultAdminDelayChangeCanceledEvent(log *types.Log) (*StatementCourtDefaultAdminDelayChangeCanceled, error) {
	event := "DefaultAdminDelayChangeCanceled"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != statementCourt.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(StatementCourtDefaultAdminDelayChangeCanceled)
	if len(log.Data) > 0 {
		if err := statementCourt.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range statementCourt.abi.Events[event].Inputs {
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

// StatementCourtDefaultAdminDelayChangeScheduled represents a DefaultAdminDelayChangeScheduled event raised by the StatementCourt contract.
type StatementCourtDefaultAdminDelayChangeScheduled struct {
	NewDelay       *big.Int
	EffectSchedule *big.Int
	Raw            *types.Log // Blockchain specific contextual infos
}

const StatementCourtDefaultAdminDelayChangeScheduledEventName = "DefaultAdminDelayChangeScheduled"

// ContractEventName returns the user-defined event name.
func (StatementCourtDefaultAdminDelayChangeScheduled) ContractEventName() string {
	return StatementCourtDefaultAdminDelayChangeScheduledEventName
}

// UnpackDefaultAdminDelayChangeScheduledEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event DefaultAdminDelayChangeScheduled(uint48 newDelay, uint48 effectSchedule)
func (statementCourt *StatementCourt) UnpackDefaultAdminDelayChangeScheduledEvent(log *types.Log) (*StatementCourtDefaultAdminDelayChangeScheduled, error) {
	event := "DefaultAdminDelayChangeScheduled"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != statementCourt.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(StatementCourtDefaultAdminDelayChangeScheduled)
	if len(log.Data) > 0 {
		if err := statementCourt.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range statementCourt.abi.Events[event].Inputs {
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

// StatementCourtDefaultAdminTransferCanceled represents a DefaultAdminTransferCanceled event raised by the StatementCourt contract.
type StatementCourtDefaultAdminTransferCanceled struct {
	Raw *types.Log // Blockchain specific contextual infos
}

const StatementCourtDefaultAdminTransferCanceledEventName = "DefaultAdminTransferCanceled"

// ContractEventName returns the user-defined event name.
func (StatementCourtDefaultAdminTransferCanceled) ContractEventName() string {
	return StatementCourtDefaultAdminTransferCanceledEventName
}

// UnpackDefaultAdminTransferCanceledEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event DefaultAdminTransferCanceled()
func (statementCourt *StatementCourt) UnpackDefaultAdminTransferCanceledEvent(log *types.Log) (*StatementCourtDefaultAdminTransferCanceled, error) {
	event := "DefaultAdminTransferCanceled"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != statementCourt.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(StatementCourtDefaultAdminTransferCanceled)
	if len(log.Data) > 0 {
		if err := statementCourt.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range statementCourt.abi.Events[event].Inputs {
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

// StatementCourtDefaultAdminTransferScheduled represents a DefaultAdminTransferScheduled event raised by the StatementCourt contract.
type StatementCourtDefaultAdminTransferScheduled struct {
	NewAdmin       common.Address
	AcceptSchedule *big.Int
	Raw            *types.Log // Blockchain specific contextual infos
}

const StatementCourtDefaultAdminTransferScheduledEventName = "DefaultAdminTransferScheduled"

// ContractEventName returns the user-defined event name.
func (StatementCourtDefaultAdminTransferScheduled) ContractEventName() string {
	return StatementCourtDefaultAdminTransferScheduledEventName
}

// UnpackDefaultAdminTransferScheduledEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event DefaultAdminTransferScheduled(address indexed newAdmin, uint48 acceptSchedule)
func (statementCourt *StatementCourt) UnpackDefaultAdminTransferScheduledEvent(log *types.Log) (*StatementCourtDefaultAdminTransferScheduled, error) {
	event := "DefaultAdminTransferScheduled"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != statementCourt.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(StatementCourtDefaultAdminTransferScheduled)
	if len(log.Data) > 0 {
		if err := statementCourt.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range statementCourt.abi.Events[event].Inputs {
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

// StatementCourtEIP712DomainChanged represents a EIP712DomainChanged event raised by the StatementCourt contract.
type StatementCourtEIP712DomainChanged struct {
	Raw *types.Log // Blockchain specific contextual infos
}

const StatementCourtEIP712DomainChangedEventName = "EIP712DomainChanged"

// ContractEventName returns the user-defined event name.
func (StatementCourtEIP712DomainChanged) ContractEventName() string {
	return StatementCourtEIP712DomainChangedEventName
}

// UnpackEIP712DomainChangedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event EIP712DomainChanged()
func (statementCourt *StatementCourt) UnpackEIP712DomainChangedEvent(log *types.Log) (*StatementCourtEIP712DomainChanged, error) {
	event := "EIP712DomainChanged"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != statementCourt.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(StatementCourtEIP712DomainChanged)
	if len(log.Data) > 0 {
		if err := statementCourt.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range statementCourt.abi.Events[event].Inputs {
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

// StatementCourtFaultProven represents a FaultProven event raised by the StatementCourt contract.
type StatementCourtFaultProven struct {
	Epoch    uint64
	Kind     uint8
	Evidence [32]byte
	Prover   common.Address
	Raw      *types.Log // Blockchain specific contextual infos
}

const StatementCourtFaultProvenEventName = "FaultProven"

// ContractEventName returns the user-defined event name.
func (StatementCourtFaultProven) ContractEventName() string {
	return StatementCourtFaultProvenEventName
}

// UnpackFaultProvenEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event FaultProven(uint64 indexed epoch, uint8 indexed kind, bytes32 indexed evidence, address prover)
func (statementCourt *StatementCourt) UnpackFaultProvenEvent(log *types.Log) (*StatementCourtFaultProven, error) {
	event := "FaultProven"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != statementCourt.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(StatementCourtFaultProven)
	if len(log.Data) > 0 {
		if err := statementCourt.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range statementCourt.abi.Events[event].Inputs {
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

// StatementCourtRoleAdminChanged represents a RoleAdminChanged event raised by the StatementCourt contract.
type StatementCourtRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               *types.Log // Blockchain specific contextual infos
}

const StatementCourtRoleAdminChangedEventName = "RoleAdminChanged"

// ContractEventName returns the user-defined event name.
func (StatementCourtRoleAdminChanged) ContractEventName() string {
	return StatementCourtRoleAdminChangedEventName
}

// UnpackRoleAdminChangedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (statementCourt *StatementCourt) UnpackRoleAdminChangedEvent(log *types.Log) (*StatementCourtRoleAdminChanged, error) {
	event := "RoleAdminChanged"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != statementCourt.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(StatementCourtRoleAdminChanged)
	if len(log.Data) > 0 {
		if err := statementCourt.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range statementCourt.abi.Events[event].Inputs {
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

// StatementCourtRoleGranted represents a RoleGranted event raised by the StatementCourt contract.
type StatementCourtRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     *types.Log // Blockchain specific contextual infos
}

const StatementCourtRoleGrantedEventName = "RoleGranted"

// ContractEventName returns the user-defined event name.
func (StatementCourtRoleGranted) ContractEventName() string {
	return StatementCourtRoleGrantedEventName
}

// UnpackRoleGrantedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (statementCourt *StatementCourt) UnpackRoleGrantedEvent(log *types.Log) (*StatementCourtRoleGranted, error) {
	event := "RoleGranted"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != statementCourt.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(StatementCourtRoleGranted)
	if len(log.Data) > 0 {
		if err := statementCourt.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range statementCourt.abi.Events[event].Inputs {
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

// StatementCourtRoleRevoked represents a RoleRevoked event raised by the StatementCourt contract.
type StatementCourtRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     *types.Log // Blockchain specific contextual infos
}

const StatementCourtRoleRevokedEventName = "RoleRevoked"

// ContractEventName returns the user-defined event name.
func (StatementCourtRoleRevoked) ContractEventName() string {
	return StatementCourtRoleRevokedEventName
}

// UnpackRoleRevokedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (statementCourt *StatementCourt) UnpackRoleRevokedEvent(log *types.Log) (*StatementCourtRoleRevoked, error) {
	event := "RoleRevoked"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != statementCourt.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(StatementCourtRoleRevoked)
	if len(log.Data) > 0 {
		if err := statementCourt.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range statementCourt.abi.Events[event].Inputs {
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

// StatementCourtSignerActivated represents a SignerActivated event raised by the StatementCourt contract.
type StatementCourtSignerActivated struct {
	Signer    common.Address
	FromEpoch uint64
	Raw       *types.Log // Blockchain specific contextual infos
}

const StatementCourtSignerActivatedEventName = "SignerActivated"

// ContractEventName returns the user-defined event name.
func (StatementCourtSignerActivated) ContractEventName() string {
	return StatementCourtSignerActivatedEventName
}

// UnpackSignerActivatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event SignerActivated(address indexed signer, uint64 fromEpoch)
func (statementCourt *StatementCourt) UnpackSignerActivatedEvent(log *types.Log) (*StatementCourtSignerActivated, error) {
	event := "SignerActivated"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != statementCourt.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(StatementCourtSignerActivated)
	if len(log.Data) > 0 {
		if err := statementCourt.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range statementCourt.abi.Events[event].Inputs {
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

// StatementCourtSignerCancelled represents a SignerCancelled event raised by the StatementCourt contract.
type StatementCourtSignerCancelled struct {
	Signer    common.Address
	FromEpoch uint64
	Raw       *types.Log // Blockchain specific contextual infos
}

const StatementCourtSignerCancelledEventName = "SignerCancelled"

// ContractEventName returns the user-defined event name.
func (StatementCourtSignerCancelled) ContractEventName() string {
	return StatementCourtSignerCancelledEventName
}

// UnpackSignerCancelledEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event SignerCancelled(address indexed signer, uint64 fromEpoch)
func (statementCourt *StatementCourt) UnpackSignerCancelledEvent(log *types.Log) (*StatementCourtSignerCancelled, error) {
	event := "SignerCancelled"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != statementCourt.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(StatementCourtSignerCancelled)
	if len(log.Data) > 0 {
		if err := statementCourt.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range statementCourt.abi.Events[event].Inputs {
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

// StatementCourtSignerScheduled represents a SignerScheduled event raised by the StatementCourt contract.
type StatementCourtSignerScheduled struct {
	Signer      common.Address
	FromEpoch   uint64
	ActivatesAt uint64
	Raw         *types.Log // Blockchain specific contextual infos
}

const StatementCourtSignerScheduledEventName = "SignerScheduled"

// ContractEventName returns the user-defined event name.
func (StatementCourtSignerScheduled) ContractEventName() string {
	return StatementCourtSignerScheduledEventName
}

// UnpackSignerScheduledEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event SignerScheduled(address indexed signer, uint64 fromEpoch, uint64 activatesAt)
func (statementCourt *StatementCourt) UnpackSignerScheduledEvent(log *types.Log) (*StatementCourtSignerScheduled, error) {
	event := "SignerScheduled"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != statementCourt.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(StatementCourtSignerScheduled)
	if len(log.Data) > 0 {
		if err := statementCourt.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range statementCourt.abi.Events[event].Inputs {
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
func (statementCourt *StatementCourt) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], statementCourt.abi.Errors["AccessControlBadConfirmation"].ID.Bytes()[:4]) {
		return statementCourt.UnpackAccessControlBadConfirmationError(raw[4:])
	}
	if bytes.Equal(raw[:4], statementCourt.abi.Errors["AccessControlEnforcedDefaultAdminDelay"].ID.Bytes()[:4]) {
		return statementCourt.UnpackAccessControlEnforcedDefaultAdminDelayError(raw[4:])
	}
	if bytes.Equal(raw[:4], statementCourt.abi.Errors["AccessControlEnforcedDefaultAdminRules"].ID.Bytes()[:4]) {
		return statementCourt.UnpackAccessControlEnforcedDefaultAdminRulesError(raw[4:])
	}
	if bytes.Equal(raw[:4], statementCourt.abi.Errors["AccessControlInvalidDefaultAdmin"].ID.Bytes()[:4]) {
		return statementCourt.UnpackAccessControlInvalidDefaultAdminError(raw[4:])
	}
	if bytes.Equal(raw[:4], statementCourt.abi.Errors["AccessControlUnauthorizedAccount"].ID.Bytes()[:4]) {
		return statementCourt.UnpackAccessControlUnauthorizedAccountError(raw[4:])
	}
	if bytes.Equal(raw[:4], statementCourt.abi.Errors["AlreadyProven"].ID.Bytes()[:4]) {
		return statementCourt.UnpackAlreadyProvenError(raw[4:])
	}
	if bytes.Equal(raw[:4], statementCourt.abi.Errors["DifferentEpochs"].ID.Bytes()[:4]) {
		return statementCourt.UnpackDifferentEpochsError(raw[4:])
	}
	if bytes.Equal(raw[:4], statementCourt.abi.Errors["EpochNotAnchored"].ID.Bytes()[:4]) {
		return statementCourt.UnpackEpochNotAnchoredError(raw[4:])
	}
	if bytes.Equal(raw[:4], statementCourt.abi.Errors["FromEpochTooLow"].ID.Bytes()[:4]) {
		return statementCourt.UnpackFromEpochTooLowError(raw[4:])
	}
	if bytes.Equal(raw[:4], statementCourt.abi.Errors["InvalidShortString"].ID.Bytes()[:4]) {
		return statementCourt.UnpackInvalidShortStringError(raw[4:])
	}
	if bytes.Equal(raw[:4], statementCourt.abi.Errors["NoConflict"].ID.Bytes()[:4]) {
		return statementCourt.UnpackNoConflictError(raw[4:])
	}
	if bytes.Equal(raw[:4], statementCourt.abi.Errors["NoSignerChangePending"].ID.Bytes()[:4]) {
		return statementCourt.UnpackNoSignerChangePendingError(raw[4:])
	}
	if bytes.Equal(raw[:4], statementCourt.abi.Errors["NotALedger"].ID.Bytes()[:4]) {
		return statementCourt.UnpackNotALedgerError(raw[4:])
	}
	if bytes.Equal(raw[:4], statementCourt.abi.Errors["NotSignedBySigner"].ID.Bytes()[:4]) {
		return statementCourt.UnpackNotSignedBySignerError(raw[4:])
	}
	if bytes.Equal(raw[:4], statementCourt.abi.Errors["ProofHolds"].ID.Bytes()[:4]) {
		return statementCourt.UnpackProofHoldsError(raw[4:])
	}
	if bytes.Equal(raw[:4], statementCourt.abi.Errors["SafeCastOverflowedUintDowncast"].ID.Bytes()[:4]) {
		return statementCourt.UnpackSafeCastOverflowedUintDowncastError(raw[4:])
	}
	if bytes.Equal(raw[:4], statementCourt.abi.Errors["SameSigner"].ID.Bytes()[:4]) {
		return statementCourt.UnpackSameSignerError(raw[4:])
	}
	if bytes.Equal(raw[:4], statementCourt.abi.Errors["SiblingsMismatch"].ID.Bytes()[:4]) {
		return statementCourt.UnpackSiblingsMismatchError(raw[4:])
	}
	if bytes.Equal(raw[:4], statementCourt.abi.Errors["SignerChangePending"].ID.Bytes()[:4]) {
		return statementCourt.UnpackSignerChangePendingError(raw[4:])
	}
	if bytes.Equal(raw[:4], statementCourt.abi.Errors["SignerChangeStale"].ID.Bytes()[:4]) {
		return statementCourt.UnpackSignerChangeStaleError(raw[4:])
	}
	if bytes.Equal(raw[:4], statementCourt.abi.Errors["SignerChangeTooEarly"].ID.Bytes()[:4]) {
		return statementCourt.UnpackSignerChangeTooEarlyError(raw[4:])
	}
	if bytes.Equal(raw[:4], statementCourt.abi.Errors["StringTooLong"].ID.Bytes()[:4]) {
		return statementCourt.UnpackStringTooLongError(raw[4:])
	}
	if bytes.Equal(raw[:4], statementCourt.abi.Errors["WrongLedger"].ID.Bytes()[:4]) {
		return statementCourt.UnpackWrongLedgerError(raw[4:])
	}
	if bytes.Equal(raw[:4], statementCourt.abi.Errors["ZeroDelay"].ID.Bytes()[:4]) {
		return statementCourt.UnpackZeroDelayError(raw[4:])
	}
	if bytes.Equal(raw[:4], statementCourt.abi.Errors["ZeroSigner"].ID.Bytes()[:4]) {
		return statementCourt.UnpackZeroSignerError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// StatementCourtAccessControlBadConfirmation represents a AccessControlBadConfirmation error raised by the StatementCourt contract.
type StatementCourtAccessControlBadConfirmation struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AccessControlBadConfirmation()
func StatementCourtAccessControlBadConfirmationErrorID() common.Hash {
	return common.HexToHash("0x6697b23232a647058342c0724fe7c415cab25915b54e5dbc03f233173d37b41c")
}

// UnpackAccessControlBadConfirmationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AccessControlBadConfirmation()
func (statementCourt *StatementCourt) UnpackAccessControlBadConfirmationError(raw []byte) (*StatementCourtAccessControlBadConfirmation, error) {
	out := new(StatementCourtAccessControlBadConfirmation)
	if err := statementCourt.abi.UnpackIntoInterface(out, "AccessControlBadConfirmation", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// StatementCourtAccessControlEnforcedDefaultAdminDelay represents a AccessControlEnforcedDefaultAdminDelay error raised by the StatementCourt contract.
type StatementCourtAccessControlEnforcedDefaultAdminDelay struct {
	Schedule *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AccessControlEnforcedDefaultAdminDelay(uint48 schedule)
func StatementCourtAccessControlEnforcedDefaultAdminDelayErrorID() common.Hash {
	return common.HexToHash("0x19ca5ebb8fb33f00e502c9392eddab1501674629178bf69b853cf037aaf4bb5d")
}

// UnpackAccessControlEnforcedDefaultAdminDelayError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AccessControlEnforcedDefaultAdminDelay(uint48 schedule)
func (statementCourt *StatementCourt) UnpackAccessControlEnforcedDefaultAdminDelayError(raw []byte) (*StatementCourtAccessControlEnforcedDefaultAdminDelay, error) {
	out := new(StatementCourtAccessControlEnforcedDefaultAdminDelay)
	if err := statementCourt.abi.UnpackIntoInterface(out, "AccessControlEnforcedDefaultAdminDelay", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// StatementCourtAccessControlEnforcedDefaultAdminRules represents a AccessControlEnforcedDefaultAdminRules error raised by the StatementCourt contract.
type StatementCourtAccessControlEnforcedDefaultAdminRules struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AccessControlEnforcedDefaultAdminRules()
func StatementCourtAccessControlEnforcedDefaultAdminRulesErrorID() common.Hash {
	return common.HexToHash("0x3fc3c27ae3db78c81b8f6e685172134623efa268ee8cd8d54be38ad2a74fc13b")
}

// UnpackAccessControlEnforcedDefaultAdminRulesError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AccessControlEnforcedDefaultAdminRules()
func (statementCourt *StatementCourt) UnpackAccessControlEnforcedDefaultAdminRulesError(raw []byte) (*StatementCourtAccessControlEnforcedDefaultAdminRules, error) {
	out := new(StatementCourtAccessControlEnforcedDefaultAdminRules)
	if err := statementCourt.abi.UnpackIntoInterface(out, "AccessControlEnforcedDefaultAdminRules", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// StatementCourtAccessControlInvalidDefaultAdmin represents a AccessControlInvalidDefaultAdmin error raised by the StatementCourt contract.
type StatementCourtAccessControlInvalidDefaultAdmin struct {
	DefaultAdmin common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AccessControlInvalidDefaultAdmin(address defaultAdmin)
func StatementCourtAccessControlInvalidDefaultAdminErrorID() common.Hash {
	return common.HexToHash("0xc22c8022f2a840d6b6a9f113407715f5bbd4e88c1b0dd9434dc00700ba609ed4")
}

// UnpackAccessControlInvalidDefaultAdminError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AccessControlInvalidDefaultAdmin(address defaultAdmin)
func (statementCourt *StatementCourt) UnpackAccessControlInvalidDefaultAdminError(raw []byte) (*StatementCourtAccessControlInvalidDefaultAdmin, error) {
	out := new(StatementCourtAccessControlInvalidDefaultAdmin)
	if err := statementCourt.abi.UnpackIntoInterface(out, "AccessControlInvalidDefaultAdmin", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// StatementCourtAccessControlUnauthorizedAccount represents a AccessControlUnauthorizedAccount error raised by the StatementCourt contract.
type StatementCourtAccessControlUnauthorizedAccount struct {
	Account    common.Address
	NeededRole [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AccessControlUnauthorizedAccount(address account, bytes32 neededRole)
func StatementCourtAccessControlUnauthorizedAccountErrorID() common.Hash {
	return common.HexToHash("0xe2517d3fbfae6f8515ef5ff1ccedc3933ab0cbbda0b492c06eb54ad10ef03b3e")
}

// UnpackAccessControlUnauthorizedAccountError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AccessControlUnauthorizedAccount(address account, bytes32 neededRole)
func (statementCourt *StatementCourt) UnpackAccessControlUnauthorizedAccountError(raw []byte) (*StatementCourtAccessControlUnauthorizedAccount, error) {
	out := new(StatementCourtAccessControlUnauthorizedAccount)
	if err := statementCourt.abi.UnpackIntoInterface(out, "AccessControlUnauthorizedAccount", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// StatementCourtAlreadyProven represents a AlreadyProven error raised by the StatementCourt contract.
type StatementCourtAlreadyProven struct {
	Evidence [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AlreadyProven(bytes32 evidence)
func StatementCourtAlreadyProvenErrorID() common.Hash {
	return common.HexToHash("0xea033684b2eb17b3fc70d5009bef3efe803909cba52e3380dd0e081697f487ad")
}

// UnpackAlreadyProvenError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AlreadyProven(bytes32 evidence)
func (statementCourt *StatementCourt) UnpackAlreadyProvenError(raw []byte) (*StatementCourtAlreadyProven, error) {
	out := new(StatementCourtAlreadyProven)
	if err := statementCourt.abi.UnpackIntoInterface(out, "AlreadyProven", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// StatementCourtDifferentEpochs represents a DifferentEpochs error raised by the StatementCourt contract.
type StatementCourtDifferentEpochs struct {
	First  uint64
	Second uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error DifferentEpochs(uint64 first, uint64 second)
func StatementCourtDifferentEpochsErrorID() common.Hash {
	return common.HexToHash("0xbfb240a28f6d744e2666b046f0954afb0cb85dfbc6d6bd027717efc63eed50f4")
}

// UnpackDifferentEpochsError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error DifferentEpochs(uint64 first, uint64 second)
func (statementCourt *StatementCourt) UnpackDifferentEpochsError(raw []byte) (*StatementCourtDifferentEpochs, error) {
	out := new(StatementCourtDifferentEpochs)
	if err := statementCourt.abi.UnpackIntoInterface(out, "DifferentEpochs", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// StatementCourtEpochNotAnchored represents a EpochNotAnchored error raised by the StatementCourt contract.
type StatementCourtEpochNotAnchored struct {
	Epoch uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error EpochNotAnchored(uint64 epoch)
func StatementCourtEpochNotAnchoredErrorID() common.Hash {
	return common.HexToHash("0xfdd66f09e3e0335148d05ffa812e23b41ab1c93ea73617ca6ddd26710f05879d")
}

// UnpackEpochNotAnchoredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error EpochNotAnchored(uint64 epoch)
func (statementCourt *StatementCourt) UnpackEpochNotAnchoredError(raw []byte) (*StatementCourtEpochNotAnchored, error) {
	out := new(StatementCourtEpochNotAnchored)
	if err := statementCourt.abi.UnpackIntoInterface(out, "EpochNotAnchored", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// StatementCourtFromEpochTooLow represents a FromEpochTooLow error raised by the StatementCourt contract.
type StatementCourtFromEpochTooLow struct {
	FromEpoch uint64
	Least     uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error FromEpochTooLow(uint64 fromEpoch, uint64 least)
func StatementCourtFromEpochTooLowErrorID() common.Hash {
	return common.HexToHash("0x2a9c655bfc90a6b2ff4c7c129889535ba67dd8cdd7b2e2387b0e2e444a2035fd")
}

// UnpackFromEpochTooLowError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error FromEpochTooLow(uint64 fromEpoch, uint64 least)
func (statementCourt *StatementCourt) UnpackFromEpochTooLowError(raw []byte) (*StatementCourtFromEpochTooLow, error) {
	out := new(StatementCourtFromEpochTooLow)
	if err := statementCourt.abi.UnpackIntoInterface(out, "FromEpochTooLow", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// StatementCourtInvalidShortString represents a InvalidShortString error raised by the StatementCourt contract.
type StatementCourtInvalidShortString struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidShortString()
func StatementCourtInvalidShortStringErrorID() common.Hash {
	return common.HexToHash("0xb3512b0c6163e5f0bafab72bb631b9d58cd7a731b082f910338aa21c83d5c274")
}

// UnpackInvalidShortStringError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidShortString()
func (statementCourt *StatementCourt) UnpackInvalidShortStringError(raw []byte) (*StatementCourtInvalidShortString, error) {
	out := new(StatementCourtInvalidShortString)
	if err := statementCourt.abi.UnpackIntoInterface(out, "InvalidShortString", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// StatementCourtNoConflict represents a NoConflict error raised by the StatementCourt contract.
type StatementCourtNoConflict struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NoConflict()
func StatementCourtNoConflictErrorID() common.Hash {
	return common.HexToHash("0x7412124e23106f671f0c15e99815ae376e602111b93496444e63ed6071334091")
}

// UnpackNoConflictError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NoConflict()
func (statementCourt *StatementCourt) UnpackNoConflictError(raw []byte) (*StatementCourtNoConflict, error) {
	out := new(StatementCourtNoConflict)
	if err := statementCourt.abi.UnpackIntoInterface(out, "NoConflict", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// StatementCourtNoSignerChangePending represents a NoSignerChangePending error raised by the StatementCourt contract.
type StatementCourtNoSignerChangePending struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NoSignerChangePending()
func StatementCourtNoSignerChangePendingErrorID() common.Hash {
	return common.HexToHash("0x16d48c849add6c39b79e80fa155b4bf87a0a29d7ad1b668b5e80f101e66238f8")
}

// UnpackNoSignerChangePendingError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NoSignerChangePending()
func (statementCourt *StatementCourt) UnpackNoSignerChangePendingError(raw []byte) (*StatementCourtNoSignerChangePending, error) {
	out := new(StatementCourtNoSignerChangePending)
	if err := statementCourt.abi.UnpackIntoInterface(out, "NoSignerChangePending", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// StatementCourtNotALedger represents a NotALedger error raised by the StatementCourt contract.
type StatementCourtNotALedger struct {
	Ledger common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotALedger(address ledger)
func StatementCourtNotALedgerErrorID() common.Hash {
	return common.HexToHash("0x44a5b4b2f2697952fce85b89e8be359271d0792b4f0be489d664b4771187455b")
}

// UnpackNotALedgerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotALedger(address ledger)
func (statementCourt *StatementCourt) UnpackNotALedgerError(raw []byte) (*StatementCourtNotALedger, error) {
	out := new(StatementCourtNotALedger)
	if err := statementCourt.abi.UnpackIntoInterface(out, "NotALedger", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// StatementCourtNotSignedBySigner represents a NotSignedBySigner error raised by the StatementCourt contract.
type StatementCourtNotSignedBySigner struct {
	Epoch     uint64
	Signer    common.Address
	Recovered common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotSignedBySigner(uint64 epoch, address signer, address recovered)
func StatementCourtNotSignedBySignerErrorID() common.Hash {
	return common.HexToHash("0xd7eecde17baa6c3c5da78ec7d0f325a551b50706d7779d410753e9ea173f5813")
}

// UnpackNotSignedBySignerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotSignedBySigner(uint64 epoch, address signer, address recovered)
func (statementCourt *StatementCourt) UnpackNotSignedBySignerError(raw []byte) (*StatementCourtNotSignedBySigner, error) {
	out := new(StatementCourtNotSignedBySigner)
	if err := statementCourt.abi.UnpackIntoInterface(out, "NotSignedBySigner", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// StatementCourtProofHolds represents a ProofHolds error raised by the StatementCourt contract.
type StatementCourtProofHolds struct {
	Epoch uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ProofHolds(uint64 epoch)
func StatementCourtProofHoldsErrorID() common.Hash {
	return common.HexToHash("0x64212abdba47381a8e97a2ef1859b440605c270d7fc490e549b78acd72e28e3d")
}

// UnpackProofHoldsError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ProofHolds(uint64 epoch)
func (statementCourt *StatementCourt) UnpackProofHoldsError(raw []byte) (*StatementCourtProofHolds, error) {
	out := new(StatementCourtProofHolds)
	if err := statementCourt.abi.UnpackIntoInterface(out, "ProofHolds", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// StatementCourtSafeCastOverflowedUintDowncast represents a SafeCastOverflowedUintDowncast error raised by the StatementCourt contract.
type StatementCourtSafeCastOverflowedUintDowncast struct {
	Bits  uint8
	Value *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SafeCastOverflowedUintDowncast(uint8 bits, uint256 value)
func StatementCourtSafeCastOverflowedUintDowncastErrorID() common.Hash {
	return common.HexToHash("0x6dfcc6503a32754ce7a89698e18201fc5294fd4aad43edefee786f88423b1a12")
}

// UnpackSafeCastOverflowedUintDowncastError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SafeCastOverflowedUintDowncast(uint8 bits, uint256 value)
func (statementCourt *StatementCourt) UnpackSafeCastOverflowedUintDowncastError(raw []byte) (*StatementCourtSafeCastOverflowedUintDowncast, error) {
	out := new(StatementCourtSafeCastOverflowedUintDowncast)
	if err := statementCourt.abi.UnpackIntoInterface(out, "SafeCastOverflowedUintDowncast", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// StatementCourtSameSigner represents a SameSigner error raised by the StatementCourt contract.
type StatementCourtSameSigner struct {
	Signer common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SameSigner(address signer)
func StatementCourtSameSignerErrorID() common.Hash {
	return common.HexToHash("0x684a329eced21ee6c2d72bdaef91eb1248ae7b6be68e0850361db232a4782ee7")
}

// UnpackSameSignerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SameSigner(address signer)
func (statementCourt *StatementCourt) UnpackSameSignerError(raw []byte) (*StatementCourtSameSigner, error) {
	out := new(StatementCourtSameSigner)
	if err := statementCourt.abi.UnpackIntoInterface(out, "SameSigner", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// StatementCourtSiblingsMismatch represents a SiblingsMismatch error raised by the StatementCourt contract.
type StatementCourtSiblingsMismatch struct {
	Signed [32]byte
	Given  [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SiblingsMismatch(bytes32 signed, bytes32 given)
func StatementCourtSiblingsMismatchErrorID() common.Hash {
	return common.HexToHash("0x7e68dade1f921625838c9a65f0505353f5cae857ab8fac0fc3a0d07f29ef5fe1")
}

// UnpackSiblingsMismatchError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SiblingsMismatch(bytes32 signed, bytes32 given)
func (statementCourt *StatementCourt) UnpackSiblingsMismatchError(raw []byte) (*StatementCourtSiblingsMismatch, error) {
	out := new(StatementCourtSiblingsMismatch)
	if err := statementCourt.abi.UnpackIntoInterface(out, "SiblingsMismatch", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// StatementCourtSignerChangePending represents a SignerChangePending error raised by the StatementCourt contract.
type StatementCourtSignerChangePending struct {
	Signer    common.Address
	FromEpoch uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SignerChangePending(address signer, uint64 fromEpoch)
func StatementCourtSignerChangePendingErrorID() common.Hash {
	return common.HexToHash("0xbe977730e37b951325d90256c2cad47e6616e166013e1587bc36024d82a30178")
}

// UnpackSignerChangePendingError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SignerChangePending(address signer, uint64 fromEpoch)
func (statementCourt *StatementCourt) UnpackSignerChangePendingError(raw []byte) (*StatementCourtSignerChangePending, error) {
	out := new(StatementCourtSignerChangePending)
	if err := statementCourt.abi.UnpackIntoInterface(out, "SignerChangePending", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// StatementCourtSignerChangeStale represents a SignerChangeStale error raised by the StatementCourt contract.
type StatementCourtSignerChangeStale struct {
	FromEpoch uint64
	NextEpoch uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SignerChangeStale(uint64 fromEpoch, uint64 nextEpoch)
func StatementCourtSignerChangeStaleErrorID() common.Hash {
	return common.HexToHash("0x6d5c81aaef3921214d8a666dce0b18be17a1443922f6414672a46c7cfcff9770")
}

// UnpackSignerChangeStaleError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SignerChangeStale(uint64 fromEpoch, uint64 nextEpoch)
func (statementCourt *StatementCourt) UnpackSignerChangeStaleError(raw []byte) (*StatementCourtSignerChangeStale, error) {
	out := new(StatementCourtSignerChangeStale)
	if err := statementCourt.abi.UnpackIntoInterface(out, "SignerChangeStale", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// StatementCourtSignerChangeTooEarly represents a SignerChangeTooEarly error raised by the StatementCourt contract.
type StatementCourtSignerChangeTooEarly struct {
	ActivatesAt uint64
	BlockTime   uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SignerChangeTooEarly(uint64 activatesAt, uint64 blockTime)
func StatementCourtSignerChangeTooEarlyErrorID() common.Hash {
	return common.HexToHash("0x2a2d47109124377d337480ded5b5da8a818f42b62f77ba12bace133b3888c5ed")
}

// UnpackSignerChangeTooEarlyError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SignerChangeTooEarly(uint64 activatesAt, uint64 blockTime)
func (statementCourt *StatementCourt) UnpackSignerChangeTooEarlyError(raw []byte) (*StatementCourtSignerChangeTooEarly, error) {
	out := new(StatementCourtSignerChangeTooEarly)
	if err := statementCourt.abi.UnpackIntoInterface(out, "SignerChangeTooEarly", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// StatementCourtStringTooLong represents a StringTooLong error raised by the StatementCourt contract.
type StatementCourtStringTooLong struct {
	Str string
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error StringTooLong(string str)
func StatementCourtStringTooLongErrorID() common.Hash {
	return common.HexToHash("0x305a27a93f8e33b7392df0a0f91d6fc63847395853c45991eec52dbf24d72381")
}

// UnpackStringTooLongError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error StringTooLong(string str)
func (statementCourt *StatementCourt) UnpackStringTooLongError(raw []byte) (*StatementCourtStringTooLong, error) {
	out := new(StatementCourtStringTooLong)
	if err := statementCourt.abi.UnpackIntoInterface(out, "StringTooLong", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// StatementCourtWrongLedger represents a WrongLedger error raised by the StatementCourt contract.
type StatementCourtWrongLedger struct {
	Expected common.Address
	Given    common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error WrongLedger(address expected, address given)
func StatementCourtWrongLedgerErrorID() common.Hash {
	return common.HexToHash("0x185f957c1ab2bcd607a1ccf4a1aa06b24f2dc77f61269eb02e5cc03658428df6")
}

// UnpackWrongLedgerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error WrongLedger(address expected, address given)
func (statementCourt *StatementCourt) UnpackWrongLedgerError(raw []byte) (*StatementCourtWrongLedger, error) {
	out := new(StatementCourtWrongLedger)
	if err := statementCourt.abi.UnpackIntoInterface(out, "WrongLedger", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// StatementCourtZeroDelay represents a ZeroDelay error raised by the StatementCourt contract.
type StatementCourtZeroDelay struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ZeroDelay()
func StatementCourtZeroDelayErrorID() common.Hash {
	return common.HexToHash("0x34c2b9dab051f3b49e0b9121d90ae110f06df318864d2741799dc06f9b92cebc")
}

// UnpackZeroDelayError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ZeroDelay()
func (statementCourt *StatementCourt) UnpackZeroDelayError(raw []byte) (*StatementCourtZeroDelay, error) {
	out := new(StatementCourtZeroDelay)
	if err := statementCourt.abi.UnpackIntoInterface(out, "ZeroDelay", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// StatementCourtZeroSigner represents a ZeroSigner error raised by the StatementCourt contract.
type StatementCourtZeroSigner struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ZeroSigner()
func StatementCourtZeroSignerErrorID() common.Hash {
	return common.HexToHash("0xe5c48ac578e050f5a863729c359600666b17bce09ec64f88719c5ce6fc799f6a")
}

// UnpackZeroSignerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ZeroSigner()
func (statementCourt *StatementCourt) UnpackZeroSignerError(raw []byte) (*StatementCourtZeroSigner, error) {
	out := new(StatementCourtZeroSigner)
	if err := statementCourt.abi.UnpackIntoInterface(out, "ZeroSigner", raw); err != nil {
		return nil, err
	}
	return out, nil
}
