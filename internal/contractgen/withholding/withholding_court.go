// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package withholding

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

// OpenRequest is an auto generated low-level Go binding around an user-defined struct.
type OpenRequest struct {
	Active bool
	Epoch  uint64
}

// Registration is an auto generated low-level Go binding around an user-defined struct.
type Registration struct {
	Holder    common.Address
	Account   [32]byte
	FromEpoch uint64
}

// Silence is an auto generated low-level Go binding around an user-defined struct.
type Silence struct {
	Registration [32]byte
	Account      [32]byte
	Holder       common.Address
	DueAt        uint64
	Recorder     common.Address
	RecordedAt   uint64
}

// StatementRequest is an auto generated low-level Go binding around an user-defined struct.
type StatementRequest struct {
	State uint8
	DueAt uint64
}

// SumNode is an auto generated low-level Go binding around an user-defined struct.
type SumNode struct {
	Hash [32]byte
	Sum  uint64
}

// WithholdingCourtMetaData contains all meta data concerning the WithholdingCourt contract.
var WithholdingCourtMetaData = bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"court_\",\"type\":\"address\",\"internalType\":\"contractIStatementCourt\"},{\"name\":\"requestWindow_\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"answerPeriod_\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"answer\",\"inputs\":[{\"name\":\"registration\",\"type\":\"tuple\",\"internalType\":\"structRegistration\",\"components\":[{\"name\":\"holder\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"account\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"fromEpoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"epoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"amount\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"index\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"siblings\",\"type\":\"tuple[]\",\"internalType\":\"structSumNode[]\",\"components\":[{\"name\":\"hash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sum\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"answerPeriod\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"commitmentOf\",\"inputs\":[{\"name\":\"account\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"court\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIStatementCourt\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eip712Domain\",\"inputs\":[],\"outputs\":[{\"name\":\"fields\",\"type\":\"bytes1\",\"internalType\":\"bytes1\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"verifyingContract\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"extensions\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ledger\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractILiabilityLedger\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"openRequestOf\",\"inputs\":[{\"name\":\"registration\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structOpenRequest\",\"components\":[{\"name\":\"active\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"epoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"recordUnanswered\",\"inputs\":[{\"name\":\"registration\",\"type\":\"tuple\",\"internalType\":\"structRegistration\",\"components\":[{\"name\":\"holder\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"account\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"fromEpoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"epoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"recoverSigner\",\"inputs\":[{\"name\":\"registration\",\"type\":\"tuple\",\"internalType\":\"structRegistration\",\"components\":[{\"name\":\"holder\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"account\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"fromEpoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registrationDigest\",\"inputs\":[{\"name\":\"registration\",\"type\":\"tuple\",\"internalType\":\"structRegistration\",\"components\":[{\"name\":\"holder\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"account\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"fromEpoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"requestOf\",\"inputs\":[{\"name\":\"registration\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"epoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structStatementRequest\",\"components\":[{\"name\":\"state\",\"type\":\"uint8\",\"internalType\":\"enumRequestState\"},{\"name\":\"dueAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"requestStatement\",\"inputs\":[{\"name\":\"registration\",\"type\":\"tuple\",\"internalType\":\"structRegistration\",\"components\":[{\"name\":\"holder\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"account\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"fromEpoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"epoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"requestWindow\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"silenceCount\",\"inputs\":[{\"name\":\"epoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"silencesOf\",\"inputs\":[{\"name\":\"epoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple[]\",\"internalType\":\"structSilence[]\",\"components\":[{\"name\":\"registration\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"holder\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"dueAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"recorder\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"recordedAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalSilences\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"EIP712DomainChanged\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SilenceRecorded\",\"inputs\":[{\"name\":\"registration\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"epoch\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"recorder\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StatementAnswered\",\"inputs\":[{\"name\":\"registration\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"epoch\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"commitment\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"amount\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"answerer\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StatementRequested\",\"inputs\":[{\"name\":\"registration\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"epoch\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"holder\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"dueAt\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AlreadyRequested\",\"inputs\":[{\"name\":\"registration\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"epoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"AnswerDoesNotHold\",\"inputs\":[{\"name\":\"epoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"AnswerTooLate\",\"inputs\":[{\"name\":\"dueAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"blockTime\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"BeforeRegistration\",\"inputs\":[{\"name\":\"epoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"fromEpoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"DurationOutOfRange\",\"inputs\":[{\"name\":\"given\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"maximum\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"EpochNotAnchored\",\"inputs\":[{\"name\":\"epoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InvalidShortString\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NoOpenRequest\",\"inputs\":[{\"name\":\"registration\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"epoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"NotACourt\",\"inputs\":[{\"name\":\"court\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"NotOverdue\",\"inputs\":[{\"name\":\"dueAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"blockTime\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"NotSignedBySigner\",\"inputs\":[{\"name\":\"fromEpoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"signer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"recovered\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"NotTheHolder\",\"inputs\":[{\"name\":\"holder\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"RequestStillOpen\",\"inputs\":[{\"name\":\"registration\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"epoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"RequestWindowClosed\",\"inputs\":[{\"name\":\"epoch\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"closedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"SafeCastOverflowedUintDowncast\",\"inputs\":[{\"name\":\"bits\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"StringTooLong\",\"inputs\":[{\"name\":\"str\",\"type\":\"string\",\"internalType\":\"string\"}]}]",
	ID:  "WithholdingCourt",
	Bin: "0x6101e0604052348015610010575f5ffd5b5060405161269338038061269383398101604081905261002f916103b7565b6040518060400160405280601981526020017f4b696e657469782057697468686f6c64696e6720436f75727400000000000000815250604051806040016040528060018152602001603160f81b81525061008e8261018460201b60201c565b6101205261009b81610184565b61014052815160208084019190912060e052815190820120610100524660a05261012760e05161010051604080517f8b73c3c69bb8fe3d512ecc4cf759cc79239f7b179b0ffacaa9a75d522b39400f60208201529081019290925260608201524660808201523060a08201525f9060c00160405160208183030381529060405280519060200120905090565b60805250503060c052610139826101ca565b610142816101ca565b6001600160401b038083166101a05281166101c0526001600160a01b0383166101805261016e83610222565b6001600160a01b0316610160525061046f915050565b5f5f829050601f815111156101b7578260405163305a27a960e01b81526004016101ae91906103f9565b60405180910390fd5b80516101c28261042e565b179392505050565b6001600160401b03811615806101ec57506301e133806001600160401b038216115b1561021f5760405163f4bc090160e01b81526001600160401b03821660048201526301e1338060248201526044016101ae565b50565b5f61023482633693f37160e11b6102c5565b61025f57604051600162b3be9f60e01b031981526001600160a01b03831660048201526024016101ae565b816001600160a01b03166356397c356040518163ffffffff1660e01b8152600401602060405180830381865afa15801561029b573d5f5f3e3d5ffd5b505050506040513d601f19601f820116820180604052508101906102bf9190610454565b92915050565b5f6102cf836102e7565b80156102e057506102e08383610335565b9392505050565b5f6102f9826301ffc9a760e01b610335565b15610329575f80610312846001600160e01b0319610359565b91509150818015610321575080155b949350505050565b505f919050565b919050565b5f80806103428585610359565b915091508180156103505750805b95945050505050565b6301ffc9a760e01b5f818152600483905290819060208260248188617530fa92505f511515601f3d11169150509250929050565b6001600160a01b038116811461021f575f5ffd5b80516001600160401b0381168114610330575f5ffd5b5f5f5f606084860312156103c9575f5ffd5b83516103d48161038d565b92506103e2602085016103a1565b91506103f0604085016103a1565b90509250925092565b602081525f82518060208401528060208501604085015e5f604082850101526040601f19601f83011684010191505092915050565b8051602080830151919081101561044e575f198160200360031b1b821691505b50919050565b5f60208284031215610464575f5ffd5b81516102e08161038d565b60805160a05160c05160e05161010051610120516101405161016051610180516101a0516101c05161217a6105195f395f81816102f10152610a8701525f818161019101526108e301525f818161011d015261121201525f8181610238015281816112e30152818161138c015281816114aa015261157101525f61163101525f61160101525f61180601525f6117de01525f61173901525f61176301525f61178d015261217a5ff3fe608060405234801561000f575f5ffd5b5060043610610114575f3560e01c806374f9214e116100ad5780639060e9491161007d578063d7e1ff0f11610063578063d7e1ff0f146102ec578063debd8cff14610313578063e3ff236d14610326575f5ffd5b80639060e949146102b0578063ba49f58d146102c3575f5ffd5b806374f9214e1461025a5780637d9175121461026257806384b0196e1461028257806386579e951461029d575f5ffd5b806342b2f6ed116100e857806342b2f6ed146101ed57806351722cad1461020057806353622e821461022057806356397c3514610233575f5ffd5b8062971a301461011857806301ffc9a7146101695780631e6363441461018c578063404bd102146101cc575b5f5ffd5b61013f7f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020015b60405180910390f35b61017c610177366004611989565b6103a7565b6040519015158152602001610160565b6101b37f000000000000000000000000000000000000000000000000000000000000000081565b60405167ffffffffffffffff9091168152602001610160565b6101df6101da3660046119f6565b61043f565b604051908152602001610160565b6101df6101fb366004611a71565b6106c4565b61021361020e366004611ad4565b610c22565b6040516101609190611aef565b6101df61022e366004611ba0565b610d08565b61013f7f000000000000000000000000000000000000000000000000000000000000000081565b6005546101df565b610275610270366004611c64565b610f03565b6040516101609190611cb4565b61028a610f88565b6040516101609796959493929190611d5a565b6101df6102ab366004611e19565b610fe6565b6101df6102be366004611e33565b611087565b6101df6102d1366004611ad4565b67ffffffffffffffff165f9081526004602052604090205490565b6101b37f000000000000000000000000000000000000000000000000000000000000000081565b61013f610321366004611e53565b6110e5565b610380610334366004611ea3565b604080518082019091525f8082526020820152505f9081526003602090815260409182902082518084019093525460ff811615158352610100900467ffffffffffffffff169082015290565b6040805182511515815260209283015167ffffffffffffffff169281019290925201610160565b5f7fffffffff0000000000000000000000000000000000000000000000000000000082167f0791a61100000000000000000000000000000000000000000000000000000000148061043957507f01ffc9a7000000000000000000000000000000000000000000000000000000007fffffffff000000000000000000000000000000000000000000000000000000008316145b92915050565b5f5f61044a84610fe6565b90505f6104578285611101565b90505f61046342611185565b825490915067ffffffffffffffff6101009091048116908216116104d45781546040517fbd173e5700000000000000000000000000000000000000000000000000000000815261010090910467ffffffffffffffff9081166004830152821660248201526044015b60405180910390fd5b81547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016600390811783555f8481526020918252604080822080547fffffffffffffffffffffffffffffffffffffffffffffff00000000000000000016905567ffffffffffffffff881682526004835290819020815160c08101835286815289840180359482019490945290929091820190610570908a611edb565b73ffffffffffffffffffffffffffffffffffffffff90811682528554610100900467ffffffffffffffff908116602080850191909152336040808601919091528783166060958601528654600180820189555f9889528389208851600490930201918255928701518184015590860151600282018054968801519186167fffffffff000000000000000000000000000000000000000000000000000000009788161774010000000000000000000000000000000000000000928616830217905560808701516003909201805460a0909801519290951696909516959095179490911690920292909217909155600580549192909161066f908490611f23565b909155505060405133815267ffffffffffffffff8616906020808901359186917f0c7b0834a2ad2109ee5b8b1e1cd086817e9d5bc9d5446b31cab1326dc9cf90e9910160405180910390a45090949350505050565b5f6106d26020860186611edb565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610762576107116020860186611edb565b6040517fb6bd8e8300000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff90911660048201523360248201526044016104cb565b5f61076c86610fe6565b90505f6107876107826060890160408a01611ad4565b6111d8565b90505f610795838888611290565b905073ffffffffffffffffffffffffffffffffffffffff811615806107e657508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1614155b1561085e576107fb6060890160408a01611ad4565b6040517fd7eecde100000000000000000000000000000000000000000000000000000000815267ffffffffffffffff909116600482015273ffffffffffffffffffffffffffffffffffffffff8084166024830152821660448201526064016104cb565b61086e6060890160408a01611ad4565b67ffffffffffffffff168567ffffffffffffffff1610156108e0578461089a60608a0160408b01611ad4565b6040517fd30d092f00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff9283166004820152911660248201526044016104cb565b5f7f000000000000000000000000000000000000000000000000000000000000000067ffffffffffffffff16610915876112e0565b67ffffffffffffffff166109299190611f23565b905080421115610978576040517fa87a071a00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff87166004820152602481018290526044016104cb565b5f8481526003602090815260409182902082518084019093525460ff8116158015845261010090910467ffffffffffffffff1691830191909152610a015760208101516040517fd63a0a6b0000000000000000000000000000000000000000000000000000000081526004810187905267ffffffffffffffff90911660248201526044016104cb565b5f85815260026020908152604080832067ffffffffffffffff8b168452909152812090815460ff166003811115610a3a57610a3a611c87565b14610a84576040517f46cdd0fb0000000000000000000000000000000000000000000000000000000081526004810187905267ffffffffffffffff891660248201526044016104cb565b5f7f0000000000000000000000000000000000000000000000000000000000000000610aaf42611185565b610ab99190611f36565b825490915060019083907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00168280021790555080825f0160016101000a81548167ffffffffffffffff021916908367ffffffffffffffff16021790555060405180604001604052806001151581526020018a67ffffffffffffffff1681525060035f8981526020019081526020015f205f820151815f015f6101000a81548160ff0219169083151502179055506020820151815f0160016101000a81548167ffffffffffffffff021916908367ffffffffffffffff1602179055509050508867ffffffffffffffff168c60200135887f2a1df6d8aab1092d6e23c878cde7a4b2880564ec3c458ec84e2a9737ca35df108f5f016020810190610bdb9190611edb565b6040805173ffffffffffffffffffffffffffffffffffffffff909216825267ffffffffffffffff871660208301520160405180910390a450949a9950505050505050505050565b67ffffffffffffffff81165f908152600460209081526040808320805482518185028101850190935280835260609492939192909184015b82821015610cfd575f8481526020908190206040805160c081018252600486029092018054835260018082015484860152600282015473ffffffffffffffffffffffffffffffffffffffff8082169486019490945267ffffffffffffffff74010000000000000000000000000000000000000000918290048116606087015260039093015493841660808601529092041660a08301529083529092019101610c5a565b505050509050919050565b5f5f610d1389610fe6565b90505f610d20828a611101565b90505f610d2c42611185565b825490915067ffffffffffffffff61010090910481169082161115610d995781546040517f5c3cd68000000000000000000000000000000000000000000000000000000000815261010090910467ffffffffffffffff9081166004830152821660248201526044016104cb565b604080517f0600000000000000000000000000000000000000000000000000000000000000602082810191909152808e0135602183015260418083018d90528351808403909101815260619092019092528051910120610dfd8b828b8b8b8b611532565b610e3f576040517f2fa1f48b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8c1660048201526024016104cb565b82547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660021783555f8481526003602090815260409182902080547fffffffffffffffffffffffffffffffffffffffffffffff000000000000000000169055815183815267ffffffffffffffff8c8116828401523393820193909352918d1691908e01359086907f3adf664ca2a6e37386d93620b94ac7ac196009f275d4438ca9efedc1b0dac6e49060600160405180910390a49b9a5050505050505050505050565b6040805180820182525f80825260208083018290528582526002815283822067ffffffffffffffff86168352905282902082518084019093528054919291829060ff166003811115610f5757610f57611c87565b6003811115610f6857610f68611c87565b81529054610100900467ffffffffffffffff166020909101529392505050565b5f6060805f5f5f6060610f996115fa565b610fa161162a565b604080515f808252602082019092527f0f000000000000000000000000000000000000000000000000000000000000009b939a50919850469750309650945092509050565b5f610439611082610ffc36859003850185611f56565b805160208083015160409384015184517fc877eb1b8555df5f5d3733cc7298d61e52a705cc10a3e6d6678fbb217d9e4da38185015273ffffffffffffffffffffffffffffffffffffffff90941684860152606084019190915267ffffffffffffffff166080808401919091528351808403909101815260a0909201909252805191012090565b611655565b604080517f060000000000000000000000000000000000000000000000000000000000000060208083019190915260218201859052604180830185905283518084039091018152606190920190925280519101205f905b9392505050565b5f6110f96110f285610fe6565b8484611290565b949350505050565b5f82815260026020908152604080832067ffffffffffffffff8516845290915281206001815460ff16600381111561113b5761113b611c87565b146110de576040517f12ef5f2c0000000000000000000000000000000000000000000000000000000081526004810185905267ffffffffffffffff841660248201526044016104cb565b5f67ffffffffffffffff8211156111d457604080517f6dfcc6500000000000000000000000000000000000000000000000000000000081526004810191909152602481018390526044016104cb565b5090565b6040517f22a3808f00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff821660048201525f907f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906322a3808f90602401602060405180830381865afa15801561126c573d5f5f3e3d5ffd5b505050506040513d601f19601f820116820180604052508101906104399190611fe4565b5f5f5f5f61129f87878761169c565b919450925090505f8260038111156112b9576112b9611c87565b1415806112c557508015155b156112d5575f93505050506110de565b509095945050505050565b5f7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663b70e6be66040518163ffffffff1660e01b8152600401602060405180830381865afa15801561134a573d5f5f3e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061136e9190611fff565b67ffffffffffffffff168267ffffffffffffffff16108061142f57507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663aea0e78b6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156113f3573d5f5f3e3d5ffd5b505050506040513d601f19601f820116820180604052508101906114179190611fff565b67ffffffffffffffff168267ffffffffffffffff1610155b15611472576040517ffdd66f0900000000000000000000000000000000000000000000000000000000815267ffffffffffffffff831660048201526024016104cb565b6040517f7d60d35500000000000000000000000000000000000000000000000000000000815267ffffffffffffffff831660048201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1690637d60d3559060240160c060405180830381865afa158015611504573d5f5f3e3d5ffd5b505050506040513d601f19601f82011682018060405250810190611528919061201a565b60a0015192915050565b6040517fc0566abe0000000000000000000000000000000000000000000000000000000081525f9073ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000169063c0566abe906115b0908a908a908a908a908a908a906004016120d8565b602060405180830381865afa1580156115cb573d5f5f3e3d5ffd5b505050506040513d601f19601f820116820180604052508101906115ef919061215b565b979650505050505050565b60606116257f00000000000000000000000000000000000000000000000000000000000000006116e3565b905090565b60606116257f00000000000000000000000000000000000000000000000000000000000000006116e3565b5f610439611661611720565b836040517f19010000000000000000000000000000000000000000000000000000000000008152600281019290925260228201526042902090565b5f808060418490036116d0578435602086013560408701355f1a6116c289828585611856565b9550955095505050506116da565b505f915060029050825b93509350939050565b60605f6116ef83611949565b6040805160208082528183019092529192505f91906020820181803683375050509182525060208101929092525090565b5f3073ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614801561178557507f000000000000000000000000000000000000000000000000000000000000000046145b156117af57507f000000000000000000000000000000000000000000000000000000000000000090565b611625604080517f8b73c3c69bb8fe3d512ecc4cf759cc79239f7b179b0ffacaa9a75d522b39400f60208201527f0000000000000000000000000000000000000000000000000000000000000000918101919091527f000000000000000000000000000000000000000000000000000000000000000060608201524660808201523060a08201525f9060c00160405160208183030381529060405280519060200120905090565b5f80807f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a084111561188f57505f9150600390508261193f565b604080515f808252602082018084528a905260ff891692820192909252606081018790526080810186905260019060a0016020604051602081039080840390855afa1580156118e0573d5f5f3e3d5ffd5b50506040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0015191505073ffffffffffffffffffffffffffffffffffffffff811661193657505f92506001915082905061193f565b92505f91508190505b9450945094915050565b5f60ff8216601f811115610439576040517fb3512b0c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f60208284031215611999575f5ffd5b81357fffffffff00000000000000000000000000000000000000000000000000000000811681146110de575f5ffd5b5f606082840312156119d8575f5ffd5b50919050565b67ffffffffffffffff811681146119f3575f5ffd5b50565b5f5f60808385031215611a07575f5ffd5b611a1184846119c8565b91506060830135611a21816119de565b809150509250929050565b5f5f83601f840112611a3c575f5ffd5b50813567ffffffffffffffff811115611a53575f5ffd5b602083019150836020828501011115611a6a575f5ffd5b9250929050565b5f5f5f5f60a08587031215611a84575f5ffd5b611a8e86866119c8565b9350606085013567ffffffffffffffff811115611aa9575f5ffd5b611ab587828801611a2c565b9094509250506080850135611ac9816119de565b939692955090935050565b5f60208284031215611ae4575f5ffd5b81356110de816119de565b602080825282518282018190525f918401906040840190835b818110156112d5578351805184526020810151602085015273ffffffffffffffffffffffffffffffffffffffff604082015116604085015267ffffffffffffffff606082015116606085015273ffffffffffffffffffffffffffffffffffffffff608082015116608085015267ffffffffffffffff60a08201511660a08501525060c083019250602084019350600181019050611b08565b5f5f5f5f5f5f5f610100888a031215611bb7575f5ffd5b611bc189896119c8565b96506060880135611bd1816119de565b95506080880135945060a0880135611be8816119de565b935060c0880135611bf8816119de565b925060e088013567ffffffffffffffff811115611c13575f5ffd5b8801601f81018a13611c23575f5ffd5b803567ffffffffffffffff811115611c39575f5ffd5b8a60208260061b8401011115611c4d575f5ffd5b602082019350809250505092959891949750929550565b5f5f60408385031215611c75575f5ffd5b823591506020830135611a21816119de565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602160045260245ffd5b8151604082019060048110611cf0577f4e487b71000000000000000000000000000000000000000000000000000000005f52602160045260245ffd5b8083525067ffffffffffffffff602084015116602083015292915050565b5f81518084528060208401602086015e5f6020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b7fff000000000000000000000000000000000000000000000000000000000000008816815260e060208201525f611d9460e0830189611d0e565b8281036040840152611da68189611d0e565b6060840188905273ffffffffffffffffffffffffffffffffffffffff8716608085015260a0840186905283810360c0850152845180825260208087019350909101905f5b81811015611e08578351835260209384019390920191600101611dea565b50909b9a5050505050505050505050565b5f60608284031215611e29575f5ffd5b6110de83836119c8565b5f5f60408385031215611e44575f5ffd5b50508035926020909101359150565b5f5f5f60808486031215611e65575f5ffd5b611e6f85856119c8565b9250606084013567ffffffffffffffff811115611e8a575f5ffd5b611e9686828701611a2c565b9497909650939450505050565b5f60208284031215611eb3575f5ffd5b5035919050565b73ffffffffffffffffffffffffffffffffffffffff811681146119f3575f5ffd5b5f60208284031215611eeb575f5ffd5b81356110de81611eba565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b8082018082111561043957610439611ef6565b67ffffffffffffffff818116838216019081111561043957610439611ef6565b5f6060828403128015611f67575f5ffd5b506040516060810167ffffffffffffffff81118282101715611fb0577f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b6040528235611fbe81611eba565b8152602083810135908201526040830135611fd8816119de565b60408201529392505050565b5f60208284031215611ff4575f5ffd5b81516110de81611eba565b5f6020828403121561200f575f5ffd5b81516110de816119de565b5f60c082840312801561202b575f5ffd5b5060405160c0810167ffffffffffffffff81118282101715612074577f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b60409081528351825260208085015190830152830151612093816119de565b604082015260608301516120a6816119de565b606082015260808301516120b9816119de565b608082015260a08301516120cc816119de565b60a08201529392505050565b67ffffffffffffffff87811682526020820187905285811660408301528416606082015260a06080820181905281018290525f8360c08301825b8581101561214d5782358252602083013561212c816119de565b67ffffffffffffffff16602083015260409283019290910190600101612112565b509998505050505050505050565b5f6020828403121561216b575f5ffd5b815180151581146110de575f5ffd",
}

// WithholdingCourt is an auto generated Go binding around an Ethereum contract.
type WithholdingCourt struct {
	abi abi.ABI
}

// GetABI returns the ABI associated with this contract binding.
func (c *WithholdingCourt) GetABI() abi.ABI {
	return c.abi
}

// NewWithholdingCourt creates a new instance of WithholdingCourt.
func NewWithholdingCourt() *WithholdingCourt {
	parsed, err := WithholdingCourtMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &WithholdingCourt{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *WithholdingCourt) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackConstructor is the Go binding used to pack the parameters required for
// contract deployment.
//
// Solidity: constructor(address court_, uint64 requestWindow_, uint64 answerPeriod_) returns()
func (withholdingCourt *WithholdingCourt) PackConstructor(court_ common.Address, requestWindow_ uint64, answerPeriod_ uint64) []byte {
	enc, err := withholdingCourt.abi.Pack("", court_, requestWindow_, answerPeriod_)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackAnswer is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x53622e82.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function answer((address,bytes32,uint64) registration, uint64 epoch, bytes32 salt, uint64 amount, uint64 index, (bytes32,uint64)[] siblings) returns(bytes32)
func (withholdingCourt *WithholdingCourt) PackAnswer(registration Registration, epoch uint64, salt [32]byte, amount uint64, index uint64, siblings []SumNode) []byte {
	enc, err := withholdingCourt.abi.Pack("answer", registration, epoch, salt, amount, index, siblings)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAnswer is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x53622e82.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function answer((address,bytes32,uint64) registration, uint64 epoch, bytes32 salt, uint64 amount, uint64 index, (bytes32,uint64)[] siblings) returns(bytes32)
func (withholdingCourt *WithholdingCourt) TryPackAnswer(registration Registration, epoch uint64, salt [32]byte, amount uint64, index uint64, siblings []SumNode) ([]byte, error) {
	return withholdingCourt.abi.Pack("answer", registration, epoch, salt, amount, index, siblings)
}

// UnpackAnswer is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x53622e82.
//
// Solidity: function answer((address,bytes32,uint64) registration, uint64 epoch, bytes32 salt, uint64 amount, uint64 index, (bytes32,uint64)[] siblings) returns(bytes32)
func (withholdingCourt *WithholdingCourt) UnpackAnswer(data []byte) ([32]byte, error) {
	out, err := withholdingCourt.abi.Unpack("answer", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackAnswerPeriod is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd7e1ff0f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function answerPeriod() view returns(uint64)
func (withholdingCourt *WithholdingCourt) PackAnswerPeriod() []byte {
	enc, err := withholdingCourt.abi.Pack("answerPeriod")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAnswerPeriod is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd7e1ff0f.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function answerPeriod() view returns(uint64)
func (withholdingCourt *WithholdingCourt) TryPackAnswerPeriod() ([]byte, error) {
	return withholdingCourt.abi.Pack("answerPeriod")
}

// UnpackAnswerPeriod is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xd7e1ff0f.
//
// Solidity: function answerPeriod() view returns(uint64)
func (withholdingCourt *WithholdingCourt) UnpackAnswerPeriod(data []byte) (uint64, error) {
	out, err := withholdingCourt.abi.Unpack("answerPeriod", data)
	if err != nil {
		return *new(uint64), err
	}
	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)
	return out0, nil
}

// PackCommitmentOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9060e949.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function commitmentOf(bytes32 account, bytes32 salt) pure returns(bytes32)
func (withholdingCourt *WithholdingCourt) PackCommitmentOf(account [32]byte, salt [32]byte) []byte {
	enc, err := withholdingCourt.abi.Pack("commitmentOf", account, salt)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCommitmentOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9060e949.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function commitmentOf(bytes32 account, bytes32 salt) pure returns(bytes32)
func (withholdingCourt *WithholdingCourt) TryPackCommitmentOf(account [32]byte, salt [32]byte) ([]byte, error) {
	return withholdingCourt.abi.Pack("commitmentOf", account, salt)
}

// UnpackCommitmentOf is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x9060e949.
//
// Solidity: function commitmentOf(bytes32 account, bytes32 salt) pure returns(bytes32)
func (withholdingCourt *WithholdingCourt) UnpackCommitmentOf(data []byte) ([32]byte, error) {
	out, err := withholdingCourt.abi.Unpack("commitmentOf", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackCourt is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x00971a30.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function court() view returns(address)
func (withholdingCourt *WithholdingCourt) PackCourt() []byte {
	enc, err := withholdingCourt.abi.Pack("court")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCourt is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x00971a30.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function court() view returns(address)
func (withholdingCourt *WithholdingCourt) TryPackCourt() ([]byte, error) {
	return withholdingCourt.abi.Pack("court")
}

// UnpackCourt is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x00971a30.
//
// Solidity: function court() view returns(address)
func (withholdingCourt *WithholdingCourt) UnpackCourt(data []byte) (common.Address, error) {
	out, err := withholdingCourt.abi.Unpack("court", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackEip712Domain is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x84b0196e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (withholdingCourt *WithholdingCourt) PackEip712Domain() []byte {
	enc, err := withholdingCourt.abi.Pack("eip712Domain")
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
func (withholdingCourt *WithholdingCourt) TryPackEip712Domain() ([]byte, error) {
	return withholdingCourt.abi.Pack("eip712Domain")
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
func (withholdingCourt *WithholdingCourt) UnpackEip712Domain(data []byte) (Eip712DomainOutput, error) {
	out, err := withholdingCourt.abi.Unpack("eip712Domain", data)
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

// PackLedger is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x56397c35.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function ledger() view returns(address)
func (withholdingCourt *WithholdingCourt) PackLedger() []byte {
	enc, err := withholdingCourt.abi.Pack("ledger")
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
func (withholdingCourt *WithholdingCourt) TryPackLedger() ([]byte, error) {
	return withholdingCourt.abi.Pack("ledger")
}

// UnpackLedger is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x56397c35.
//
// Solidity: function ledger() view returns(address)
func (withholdingCourt *WithholdingCourt) UnpackLedger(data []byte) (common.Address, error) {
	out, err := withholdingCourt.abi.Unpack("ledger", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackOpenRequestOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe3ff236d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function openRequestOf(bytes32 registration) view returns((bool,uint64))
func (withholdingCourt *WithholdingCourt) PackOpenRequestOf(registration [32]byte) []byte {
	enc, err := withholdingCourt.abi.Pack("openRequestOf", registration)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackOpenRequestOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe3ff236d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function openRequestOf(bytes32 registration) view returns((bool,uint64))
func (withholdingCourt *WithholdingCourt) TryPackOpenRequestOf(registration [32]byte) ([]byte, error) {
	return withholdingCourt.abi.Pack("openRequestOf", registration)
}

// UnpackOpenRequestOf is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe3ff236d.
//
// Solidity: function openRequestOf(bytes32 registration) view returns((bool,uint64))
func (withholdingCourt *WithholdingCourt) UnpackOpenRequestOf(data []byte) (OpenRequest, error) {
	out, err := withholdingCourt.abi.Unpack("openRequestOf", data)
	if err != nil {
		return *new(OpenRequest), err
	}
	out0 := *abi.ConvertType(out[0], new(OpenRequest)).(*OpenRequest)
	return out0, nil
}

// PackRecordUnanswered is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x404bd102.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function recordUnanswered((address,bytes32,uint64) registration, uint64 epoch) returns(bytes32)
func (withholdingCourt *WithholdingCourt) PackRecordUnanswered(registration Registration, epoch uint64) []byte {
	enc, err := withholdingCourt.abi.Pack("recordUnanswered", registration, epoch)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRecordUnanswered is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x404bd102.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function recordUnanswered((address,bytes32,uint64) registration, uint64 epoch) returns(bytes32)
func (withholdingCourt *WithholdingCourt) TryPackRecordUnanswered(registration Registration, epoch uint64) ([]byte, error) {
	return withholdingCourt.abi.Pack("recordUnanswered", registration, epoch)
}

// UnpackRecordUnanswered is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x404bd102.
//
// Solidity: function recordUnanswered((address,bytes32,uint64) registration, uint64 epoch) returns(bytes32)
func (withholdingCourt *WithholdingCourt) UnpackRecordUnanswered(data []byte) ([32]byte, error) {
	out, err := withholdingCourt.abi.Unpack("recordUnanswered", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackRecoverSigner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdebd8cff.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function recoverSigner((address,bytes32,uint64) registration, bytes signature) view returns(address)
func (withholdingCourt *WithholdingCourt) PackRecoverSigner(registration Registration, signature []byte) []byte {
	enc, err := withholdingCourt.abi.Pack("recoverSigner", registration, signature)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRecoverSigner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdebd8cff.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function recoverSigner((address,bytes32,uint64) registration, bytes signature) view returns(address)
func (withholdingCourt *WithholdingCourt) TryPackRecoverSigner(registration Registration, signature []byte) ([]byte, error) {
	return withholdingCourt.abi.Pack("recoverSigner", registration, signature)
}

// UnpackRecoverSigner is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xdebd8cff.
//
// Solidity: function recoverSigner((address,bytes32,uint64) registration, bytes signature) view returns(address)
func (withholdingCourt *WithholdingCourt) UnpackRecoverSigner(data []byte) (common.Address, error) {
	out, err := withholdingCourt.abi.Unpack("recoverSigner", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackRegistrationDigest is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x86579e95.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function registrationDigest((address,bytes32,uint64) registration) view returns(bytes32)
func (withholdingCourt *WithholdingCourt) PackRegistrationDigest(registration Registration) []byte {
	enc, err := withholdingCourt.abi.Pack("registrationDigest", registration)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRegistrationDigest is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x86579e95.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function registrationDigest((address,bytes32,uint64) registration) view returns(bytes32)
func (withholdingCourt *WithholdingCourt) TryPackRegistrationDigest(registration Registration) ([]byte, error) {
	return withholdingCourt.abi.Pack("registrationDigest", registration)
}

// UnpackRegistrationDigest is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x86579e95.
//
// Solidity: function registrationDigest((address,bytes32,uint64) registration) view returns(bytes32)
func (withholdingCourt *WithholdingCourt) UnpackRegistrationDigest(data []byte) ([32]byte, error) {
	out, err := withholdingCourt.abi.Unpack("registrationDigest", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackRequestOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7d917512.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function requestOf(bytes32 registration, uint64 epoch) view returns((uint8,uint64))
func (withholdingCourt *WithholdingCourt) PackRequestOf(registration [32]byte, epoch uint64) []byte {
	enc, err := withholdingCourt.abi.Pack("requestOf", registration, epoch)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRequestOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7d917512.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function requestOf(bytes32 registration, uint64 epoch) view returns((uint8,uint64))
func (withholdingCourt *WithholdingCourt) TryPackRequestOf(registration [32]byte, epoch uint64) ([]byte, error) {
	return withholdingCourt.abi.Pack("requestOf", registration, epoch)
}

// UnpackRequestOf is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x7d917512.
//
// Solidity: function requestOf(bytes32 registration, uint64 epoch) view returns((uint8,uint64))
func (withholdingCourt *WithholdingCourt) UnpackRequestOf(data []byte) (StatementRequest, error) {
	out, err := withholdingCourt.abi.Unpack("requestOf", data)
	if err != nil {
		return *new(StatementRequest), err
	}
	out0 := *abi.ConvertType(out[0], new(StatementRequest)).(*StatementRequest)
	return out0, nil
}

// PackRequestStatement is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x42b2f6ed.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function requestStatement((address,bytes32,uint64) registration, bytes signature, uint64 epoch) returns(bytes32)
func (withholdingCourt *WithholdingCourt) PackRequestStatement(registration Registration, signature []byte, epoch uint64) []byte {
	enc, err := withholdingCourt.abi.Pack("requestStatement", registration, signature, epoch)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRequestStatement is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x42b2f6ed.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function requestStatement((address,bytes32,uint64) registration, bytes signature, uint64 epoch) returns(bytes32)
func (withholdingCourt *WithholdingCourt) TryPackRequestStatement(registration Registration, signature []byte, epoch uint64) ([]byte, error) {
	return withholdingCourt.abi.Pack("requestStatement", registration, signature, epoch)
}

// UnpackRequestStatement is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x42b2f6ed.
//
// Solidity: function requestStatement((address,bytes32,uint64) registration, bytes signature, uint64 epoch) returns(bytes32)
func (withholdingCourt *WithholdingCourt) UnpackRequestStatement(data []byte) ([32]byte, error) {
	out, err := withholdingCourt.abi.Unpack("requestStatement", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackRequestWindow is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1e636344.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function requestWindow() view returns(uint64)
func (withholdingCourt *WithholdingCourt) PackRequestWindow() []byte {
	enc, err := withholdingCourt.abi.Pack("requestWindow")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRequestWindow is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1e636344.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function requestWindow() view returns(uint64)
func (withholdingCourt *WithholdingCourt) TryPackRequestWindow() ([]byte, error) {
	return withholdingCourt.abi.Pack("requestWindow")
}

// UnpackRequestWindow is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x1e636344.
//
// Solidity: function requestWindow() view returns(uint64)
func (withholdingCourt *WithholdingCourt) UnpackRequestWindow(data []byte) (uint64, error) {
	out, err := withholdingCourt.abi.Unpack("requestWindow", data)
	if err != nil {
		return *new(uint64), err
	}
	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)
	return out0, nil
}

// PackSilenceCount is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xba49f58d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function silenceCount(uint64 epoch) view returns(uint256)
func (withholdingCourt *WithholdingCourt) PackSilenceCount(epoch uint64) []byte {
	enc, err := withholdingCourt.abi.Pack("silenceCount", epoch)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSilenceCount is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xba49f58d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function silenceCount(uint64 epoch) view returns(uint256)
func (withholdingCourt *WithholdingCourt) TryPackSilenceCount(epoch uint64) ([]byte, error) {
	return withholdingCourt.abi.Pack("silenceCount", epoch)
}

// UnpackSilenceCount is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xba49f58d.
//
// Solidity: function silenceCount(uint64 epoch) view returns(uint256)
func (withholdingCourt *WithholdingCourt) UnpackSilenceCount(data []byte) (*big.Int, error) {
	out, err := withholdingCourt.abi.Unpack("silenceCount", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackSilencesOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x51722cad.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function silencesOf(uint64 epoch) view returns((bytes32,bytes32,address,uint64,address,uint64)[])
func (withholdingCourt *WithholdingCourt) PackSilencesOf(epoch uint64) []byte {
	enc, err := withholdingCourt.abi.Pack("silencesOf", epoch)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSilencesOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x51722cad.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function silencesOf(uint64 epoch) view returns((bytes32,bytes32,address,uint64,address,uint64)[])
func (withholdingCourt *WithholdingCourt) TryPackSilencesOf(epoch uint64) ([]byte, error) {
	return withholdingCourt.abi.Pack("silencesOf", epoch)
}

// UnpackSilencesOf is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x51722cad.
//
// Solidity: function silencesOf(uint64 epoch) view returns((bytes32,bytes32,address,uint64,address,uint64)[])
func (withholdingCourt *WithholdingCourt) UnpackSilencesOf(data []byte) ([]Silence, error) {
	out, err := withholdingCourt.abi.Unpack("silencesOf", data)
	if err != nil {
		return *new([]Silence), err
	}
	out0 := *abi.ConvertType(out[0], new([]Silence)).(*[]Silence)
	return out0, nil
}

// PackSupportsInterface is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x01ffc9a7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (withholdingCourt *WithholdingCourt) PackSupportsInterface(interfaceId [4]byte) []byte {
	enc, err := withholdingCourt.abi.Pack("supportsInterface", interfaceId)
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
func (withholdingCourt *WithholdingCourt) TryPackSupportsInterface(interfaceId [4]byte) ([]byte, error) {
	return withholdingCourt.abi.Pack("supportsInterface", interfaceId)
}

// UnpackSupportsInterface is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (withholdingCourt *WithholdingCourt) UnpackSupportsInterface(data []byte) (bool, error) {
	out, err := withholdingCourt.abi.Unpack("supportsInterface", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackTotalSilences is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x74f9214e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function totalSilences() view returns(uint256)
func (withholdingCourt *WithholdingCourt) PackTotalSilences() []byte {
	enc, err := withholdingCourt.abi.Pack("totalSilences")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTotalSilences is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x74f9214e.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function totalSilences() view returns(uint256)
func (withholdingCourt *WithholdingCourt) TryPackTotalSilences() ([]byte, error) {
	return withholdingCourt.abi.Pack("totalSilences")
}

// UnpackTotalSilences is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x74f9214e.
//
// Solidity: function totalSilences() view returns(uint256)
func (withholdingCourt *WithholdingCourt) UnpackTotalSilences(data []byte) (*big.Int, error) {
	out, err := withholdingCourt.abi.Unpack("totalSilences", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// WithholdingCourtEIP712DomainChanged represents a EIP712DomainChanged event raised by the WithholdingCourt contract.
type WithholdingCourtEIP712DomainChanged struct {
	Raw *types.Log // Blockchain specific contextual infos
}

const WithholdingCourtEIP712DomainChangedEventName = "EIP712DomainChanged"

// ContractEventName returns the user-defined event name.
func (WithholdingCourtEIP712DomainChanged) ContractEventName() string {
	return WithholdingCourtEIP712DomainChangedEventName
}

// UnpackEIP712DomainChangedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event EIP712DomainChanged()
func (withholdingCourt *WithholdingCourt) UnpackEIP712DomainChangedEvent(log *types.Log) (*WithholdingCourtEIP712DomainChanged, error) {
	event := "EIP712DomainChanged"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != withholdingCourt.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(WithholdingCourtEIP712DomainChanged)
	if len(log.Data) > 0 {
		if err := withholdingCourt.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range withholdingCourt.abi.Events[event].Inputs {
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

// WithholdingCourtSilenceRecorded represents a SilenceRecorded event raised by the WithholdingCourt contract.
type WithholdingCourtSilenceRecorded struct {
	Registration [32]byte
	Account      [32]byte
	Epoch        uint64
	Recorder     common.Address
	Raw          *types.Log // Blockchain specific contextual infos
}

const WithholdingCourtSilenceRecordedEventName = "SilenceRecorded"

// ContractEventName returns the user-defined event name.
func (WithholdingCourtSilenceRecorded) ContractEventName() string {
	return WithholdingCourtSilenceRecordedEventName
}

// UnpackSilenceRecordedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event SilenceRecorded(bytes32 indexed registration, bytes32 indexed account, uint64 indexed epoch, address recorder)
func (withholdingCourt *WithholdingCourt) UnpackSilenceRecordedEvent(log *types.Log) (*WithholdingCourtSilenceRecorded, error) {
	event := "SilenceRecorded"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != withholdingCourt.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(WithholdingCourtSilenceRecorded)
	if len(log.Data) > 0 {
		if err := withholdingCourt.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range withholdingCourt.abi.Events[event].Inputs {
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

// WithholdingCourtStatementAnswered represents a StatementAnswered event raised by the WithholdingCourt contract.
type WithholdingCourtStatementAnswered struct {
	Registration [32]byte
	Account      [32]byte
	Epoch        uint64
	Commitment   [32]byte
	Amount       uint64
	Answerer     common.Address
	Raw          *types.Log // Blockchain specific contextual infos
}

const WithholdingCourtStatementAnsweredEventName = "StatementAnswered"

// ContractEventName returns the user-defined event name.
func (WithholdingCourtStatementAnswered) ContractEventName() string {
	return WithholdingCourtStatementAnsweredEventName
}

// UnpackStatementAnsweredEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event StatementAnswered(bytes32 indexed registration, bytes32 indexed account, uint64 indexed epoch, bytes32 commitment, uint64 amount, address answerer)
func (withholdingCourt *WithholdingCourt) UnpackStatementAnsweredEvent(log *types.Log) (*WithholdingCourtStatementAnswered, error) {
	event := "StatementAnswered"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != withholdingCourt.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(WithholdingCourtStatementAnswered)
	if len(log.Data) > 0 {
		if err := withholdingCourt.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range withholdingCourt.abi.Events[event].Inputs {
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

// WithholdingCourtStatementRequested represents a StatementRequested event raised by the WithholdingCourt contract.
type WithholdingCourtStatementRequested struct {
	Registration [32]byte
	Account      [32]byte
	Epoch        uint64
	Holder       common.Address
	DueAt        uint64
	Raw          *types.Log // Blockchain specific contextual infos
}

const WithholdingCourtStatementRequestedEventName = "StatementRequested"

// ContractEventName returns the user-defined event name.
func (WithholdingCourtStatementRequested) ContractEventName() string {
	return WithholdingCourtStatementRequestedEventName
}

// UnpackStatementRequestedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event StatementRequested(bytes32 indexed registration, bytes32 indexed account, uint64 indexed epoch, address holder, uint64 dueAt)
func (withholdingCourt *WithholdingCourt) UnpackStatementRequestedEvent(log *types.Log) (*WithholdingCourtStatementRequested, error) {
	event := "StatementRequested"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != withholdingCourt.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(WithholdingCourtStatementRequested)
	if len(log.Data) > 0 {
		if err := withholdingCourt.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range withholdingCourt.abi.Events[event].Inputs {
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
func (withholdingCourt *WithholdingCourt) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], withholdingCourt.abi.Errors["AlreadyRequested"].ID.Bytes()[:4]) {
		return withholdingCourt.UnpackAlreadyRequestedError(raw[4:])
	}
	if bytes.Equal(raw[:4], withholdingCourt.abi.Errors["AnswerDoesNotHold"].ID.Bytes()[:4]) {
		return withholdingCourt.UnpackAnswerDoesNotHoldError(raw[4:])
	}
	if bytes.Equal(raw[:4], withholdingCourt.abi.Errors["AnswerTooLate"].ID.Bytes()[:4]) {
		return withholdingCourt.UnpackAnswerTooLateError(raw[4:])
	}
	if bytes.Equal(raw[:4], withholdingCourt.abi.Errors["BeforeRegistration"].ID.Bytes()[:4]) {
		return withholdingCourt.UnpackBeforeRegistrationError(raw[4:])
	}
	if bytes.Equal(raw[:4], withholdingCourt.abi.Errors["DurationOutOfRange"].ID.Bytes()[:4]) {
		return withholdingCourt.UnpackDurationOutOfRangeError(raw[4:])
	}
	if bytes.Equal(raw[:4], withholdingCourt.abi.Errors["EpochNotAnchored"].ID.Bytes()[:4]) {
		return withholdingCourt.UnpackEpochNotAnchoredError(raw[4:])
	}
	if bytes.Equal(raw[:4], withholdingCourt.abi.Errors["InvalidShortString"].ID.Bytes()[:4]) {
		return withholdingCourt.UnpackInvalidShortStringError(raw[4:])
	}
	if bytes.Equal(raw[:4], withholdingCourt.abi.Errors["NoOpenRequest"].ID.Bytes()[:4]) {
		return withholdingCourt.UnpackNoOpenRequestError(raw[4:])
	}
	if bytes.Equal(raw[:4], withholdingCourt.abi.Errors["NotACourt"].ID.Bytes()[:4]) {
		return withholdingCourt.UnpackNotACourtError(raw[4:])
	}
	if bytes.Equal(raw[:4], withholdingCourt.abi.Errors["NotOverdue"].ID.Bytes()[:4]) {
		return withholdingCourt.UnpackNotOverdueError(raw[4:])
	}
	if bytes.Equal(raw[:4], withholdingCourt.abi.Errors["NotSignedBySigner"].ID.Bytes()[:4]) {
		return withholdingCourt.UnpackNotSignedBySignerError(raw[4:])
	}
	if bytes.Equal(raw[:4], withholdingCourt.abi.Errors["NotTheHolder"].ID.Bytes()[:4]) {
		return withholdingCourt.UnpackNotTheHolderError(raw[4:])
	}
	if bytes.Equal(raw[:4], withholdingCourt.abi.Errors["RequestStillOpen"].ID.Bytes()[:4]) {
		return withholdingCourt.UnpackRequestStillOpenError(raw[4:])
	}
	if bytes.Equal(raw[:4], withholdingCourt.abi.Errors["RequestWindowClosed"].ID.Bytes()[:4]) {
		return withholdingCourt.UnpackRequestWindowClosedError(raw[4:])
	}
	if bytes.Equal(raw[:4], withholdingCourt.abi.Errors["SafeCastOverflowedUintDowncast"].ID.Bytes()[:4]) {
		return withholdingCourt.UnpackSafeCastOverflowedUintDowncastError(raw[4:])
	}
	if bytes.Equal(raw[:4], withholdingCourt.abi.Errors["StringTooLong"].ID.Bytes()[:4]) {
		return withholdingCourt.UnpackStringTooLongError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// WithholdingCourtAlreadyRequested represents a AlreadyRequested error raised by the WithholdingCourt contract.
type WithholdingCourtAlreadyRequested struct {
	Registration [32]byte
	Epoch        uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AlreadyRequested(bytes32 registration, uint64 epoch)
func WithholdingCourtAlreadyRequestedErrorID() common.Hash {
	return common.HexToHash("0x46cdd0fb7d2db3112c73591b2011dc63aca93f893280dca515f60dd290122c34")
}

// UnpackAlreadyRequestedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AlreadyRequested(bytes32 registration, uint64 epoch)
func (withholdingCourt *WithholdingCourt) UnpackAlreadyRequestedError(raw []byte) (*WithholdingCourtAlreadyRequested, error) {
	out := new(WithholdingCourtAlreadyRequested)
	if err := withholdingCourt.abi.UnpackIntoInterface(out, "AlreadyRequested", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// WithholdingCourtAnswerDoesNotHold represents a AnswerDoesNotHold error raised by the WithholdingCourt contract.
type WithholdingCourtAnswerDoesNotHold struct {
	Epoch uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AnswerDoesNotHold(uint64 epoch)
func WithholdingCourtAnswerDoesNotHoldErrorID() common.Hash {
	return common.HexToHash("0x2fa1f48b39b9b6b6b9cfebd069c79075983f2ff2db1d8c7dcfd897dc5429dd28")
}

// UnpackAnswerDoesNotHoldError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AnswerDoesNotHold(uint64 epoch)
func (withholdingCourt *WithholdingCourt) UnpackAnswerDoesNotHoldError(raw []byte) (*WithholdingCourtAnswerDoesNotHold, error) {
	out := new(WithholdingCourtAnswerDoesNotHold)
	if err := withholdingCourt.abi.UnpackIntoInterface(out, "AnswerDoesNotHold", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// WithholdingCourtAnswerTooLate represents a AnswerTooLate error raised by the WithholdingCourt contract.
type WithholdingCourtAnswerTooLate struct {
	DueAt     uint64
	BlockTime uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AnswerTooLate(uint64 dueAt, uint64 blockTime)
func WithholdingCourtAnswerTooLateErrorID() common.Hash {
	return common.HexToHash("0x5c3cd680ab4930d0fa30c8646c992355155b81cc12e62635206ad8881bfcfc6c")
}

// UnpackAnswerTooLateError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AnswerTooLate(uint64 dueAt, uint64 blockTime)
func (withholdingCourt *WithholdingCourt) UnpackAnswerTooLateError(raw []byte) (*WithholdingCourtAnswerTooLate, error) {
	out := new(WithholdingCourtAnswerTooLate)
	if err := withholdingCourt.abi.UnpackIntoInterface(out, "AnswerTooLate", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// WithholdingCourtBeforeRegistration represents a BeforeRegistration error raised by the WithholdingCourt contract.
type WithholdingCourtBeforeRegistration struct {
	Epoch     uint64
	FromEpoch uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error BeforeRegistration(uint64 epoch, uint64 fromEpoch)
func WithholdingCourtBeforeRegistrationErrorID() common.Hash {
	return common.HexToHash("0xd30d092f5afe15b07c027869b0e62d560be8c5e7d28f175dd1f43a5bbecfb25b")
}

// UnpackBeforeRegistrationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error BeforeRegistration(uint64 epoch, uint64 fromEpoch)
func (withholdingCourt *WithholdingCourt) UnpackBeforeRegistrationError(raw []byte) (*WithholdingCourtBeforeRegistration, error) {
	out := new(WithholdingCourtBeforeRegistration)
	if err := withholdingCourt.abi.UnpackIntoInterface(out, "BeforeRegistration", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// WithholdingCourtDurationOutOfRange represents a DurationOutOfRange error raised by the WithholdingCourt contract.
type WithholdingCourtDurationOutOfRange struct {
	Given   uint64
	Maximum uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error DurationOutOfRange(uint64 given, uint64 maximum)
func WithholdingCourtDurationOutOfRangeErrorID() common.Hash {
	return common.HexToHash("0xf4bc09016fa12667f401ba3f9820d43ded35635f99b8b39a6a755a97e47db3bb")
}

// UnpackDurationOutOfRangeError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error DurationOutOfRange(uint64 given, uint64 maximum)
func (withholdingCourt *WithholdingCourt) UnpackDurationOutOfRangeError(raw []byte) (*WithholdingCourtDurationOutOfRange, error) {
	out := new(WithholdingCourtDurationOutOfRange)
	if err := withholdingCourt.abi.UnpackIntoInterface(out, "DurationOutOfRange", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// WithholdingCourtEpochNotAnchored represents a EpochNotAnchored error raised by the WithholdingCourt contract.
type WithholdingCourtEpochNotAnchored struct {
	Epoch uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error EpochNotAnchored(uint64 epoch)
func WithholdingCourtEpochNotAnchoredErrorID() common.Hash {
	return common.HexToHash("0xfdd66f09e3e0335148d05ffa812e23b41ab1c93ea73617ca6ddd26710f05879d")
}

// UnpackEpochNotAnchoredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error EpochNotAnchored(uint64 epoch)
func (withholdingCourt *WithholdingCourt) UnpackEpochNotAnchoredError(raw []byte) (*WithholdingCourtEpochNotAnchored, error) {
	out := new(WithholdingCourtEpochNotAnchored)
	if err := withholdingCourt.abi.UnpackIntoInterface(out, "EpochNotAnchored", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// WithholdingCourtInvalidShortString represents a InvalidShortString error raised by the WithholdingCourt contract.
type WithholdingCourtInvalidShortString struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidShortString()
func WithholdingCourtInvalidShortStringErrorID() common.Hash {
	return common.HexToHash("0xb3512b0c6163e5f0bafab72bb631b9d58cd7a731b082f910338aa21c83d5c274")
}

// UnpackInvalidShortStringError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidShortString()
func (withholdingCourt *WithholdingCourt) UnpackInvalidShortStringError(raw []byte) (*WithholdingCourtInvalidShortString, error) {
	out := new(WithholdingCourtInvalidShortString)
	if err := withholdingCourt.abi.UnpackIntoInterface(out, "InvalidShortString", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// WithholdingCourtNoOpenRequest represents a NoOpenRequest error raised by the WithholdingCourt contract.
type WithholdingCourtNoOpenRequest struct {
	Registration [32]byte
	Epoch        uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NoOpenRequest(bytes32 registration, uint64 epoch)
func WithholdingCourtNoOpenRequestErrorID() common.Hash {
	return common.HexToHash("0x12ef5f2cac24af69d330f4a7ae9d2d98b14d43bae36fbacf6df51b7f36fd3c8d")
}

// UnpackNoOpenRequestError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NoOpenRequest(bytes32 registration, uint64 epoch)
func (withholdingCourt *WithholdingCourt) UnpackNoOpenRequestError(raw []byte) (*WithholdingCourtNoOpenRequest, error) {
	out := new(WithholdingCourtNoOpenRequest)
	if err := withholdingCourt.abi.UnpackIntoInterface(out, "NoOpenRequest", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// WithholdingCourtNotACourt represents a NotACourt error raised by the WithholdingCourt contract.
type WithholdingCourtNotACourt struct {
	Court common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotACourt(address court)
func WithholdingCourtNotACourtErrorID() common.Hash {
	return common.HexToHash("0xff4c4161b866d228828099f46f8d1024545fc32c86bb16e694c95f55684edd9e")
}

// UnpackNotACourtError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotACourt(address court)
func (withholdingCourt *WithholdingCourt) UnpackNotACourtError(raw []byte) (*WithholdingCourtNotACourt, error) {
	out := new(WithholdingCourtNotACourt)
	if err := withholdingCourt.abi.UnpackIntoInterface(out, "NotACourt", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// WithholdingCourtNotOverdue represents a NotOverdue error raised by the WithholdingCourt contract.
type WithholdingCourtNotOverdue struct {
	DueAt     uint64
	BlockTime uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotOverdue(uint64 dueAt, uint64 blockTime)
func WithholdingCourtNotOverdueErrorID() common.Hash {
	return common.HexToHash("0xbd173e570140cc2e12f7e46a73691e55b8e11f9cd24654d209d94e6d000e0497")
}

// UnpackNotOverdueError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotOverdue(uint64 dueAt, uint64 blockTime)
func (withholdingCourt *WithholdingCourt) UnpackNotOverdueError(raw []byte) (*WithholdingCourtNotOverdue, error) {
	out := new(WithholdingCourtNotOverdue)
	if err := withholdingCourt.abi.UnpackIntoInterface(out, "NotOverdue", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// WithholdingCourtNotSignedBySigner represents a NotSignedBySigner error raised by the WithholdingCourt contract.
type WithholdingCourtNotSignedBySigner struct {
	FromEpoch uint64
	Signer    common.Address
	Recovered common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotSignedBySigner(uint64 fromEpoch, address signer, address recovered)
func WithholdingCourtNotSignedBySignerErrorID() common.Hash {
	return common.HexToHash("0xd7eecde17baa6c3c5da78ec7d0f325a551b50706d7779d410753e9ea173f5813")
}

// UnpackNotSignedBySignerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotSignedBySigner(uint64 fromEpoch, address signer, address recovered)
func (withholdingCourt *WithholdingCourt) UnpackNotSignedBySignerError(raw []byte) (*WithholdingCourtNotSignedBySigner, error) {
	out := new(WithholdingCourtNotSignedBySigner)
	if err := withholdingCourt.abi.UnpackIntoInterface(out, "NotSignedBySigner", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// WithholdingCourtNotTheHolder represents a NotTheHolder error raised by the WithholdingCourt contract.
type WithholdingCourtNotTheHolder struct {
	Holder common.Address
	Caller common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotTheHolder(address holder, address caller)
func WithholdingCourtNotTheHolderErrorID() common.Hash {
	return common.HexToHash("0xb6bd8e831b87133a951788a62f78aee0a480ce4f96f3cbdb83eec5e8697c61b7")
}

// UnpackNotTheHolderError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotTheHolder(address holder, address caller)
func (withholdingCourt *WithholdingCourt) UnpackNotTheHolderError(raw []byte) (*WithholdingCourtNotTheHolder, error) {
	out := new(WithholdingCourtNotTheHolder)
	if err := withholdingCourt.abi.UnpackIntoInterface(out, "NotTheHolder", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// WithholdingCourtRequestStillOpen represents a RequestStillOpen error raised by the WithholdingCourt contract.
type WithholdingCourtRequestStillOpen struct {
	Registration [32]byte
	Epoch        uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RequestStillOpen(bytes32 registration, uint64 epoch)
func WithholdingCourtRequestStillOpenErrorID() common.Hash {
	return common.HexToHash("0xd63a0a6b70d31a6aa96ef62b378a9e7718d08212759a72a47918dfdb97106c00")
}

// UnpackRequestStillOpenError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RequestStillOpen(bytes32 registration, uint64 epoch)
func (withholdingCourt *WithholdingCourt) UnpackRequestStillOpenError(raw []byte) (*WithholdingCourtRequestStillOpen, error) {
	out := new(WithholdingCourtRequestStillOpen)
	if err := withholdingCourt.abi.UnpackIntoInterface(out, "RequestStillOpen", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// WithholdingCourtRequestWindowClosed represents a RequestWindowClosed error raised by the WithholdingCourt contract.
type WithholdingCourtRequestWindowClosed struct {
	Epoch    uint64
	ClosedAt *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RequestWindowClosed(uint64 epoch, uint256 closedAt)
func WithholdingCourtRequestWindowClosedErrorID() common.Hash {
	return common.HexToHash("0xa87a071a48c134dae493271860218ee72ff911a816b5bc0dbdc9aa522b7979bd")
}

// UnpackRequestWindowClosedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RequestWindowClosed(uint64 epoch, uint256 closedAt)
func (withholdingCourt *WithholdingCourt) UnpackRequestWindowClosedError(raw []byte) (*WithholdingCourtRequestWindowClosed, error) {
	out := new(WithholdingCourtRequestWindowClosed)
	if err := withholdingCourt.abi.UnpackIntoInterface(out, "RequestWindowClosed", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// WithholdingCourtSafeCastOverflowedUintDowncast represents a SafeCastOverflowedUintDowncast error raised by the WithholdingCourt contract.
type WithholdingCourtSafeCastOverflowedUintDowncast struct {
	Bits  uint8
	Value *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SafeCastOverflowedUintDowncast(uint8 bits, uint256 value)
func WithholdingCourtSafeCastOverflowedUintDowncastErrorID() common.Hash {
	return common.HexToHash("0x6dfcc6503a32754ce7a89698e18201fc5294fd4aad43edefee786f88423b1a12")
}

// UnpackSafeCastOverflowedUintDowncastError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SafeCastOverflowedUintDowncast(uint8 bits, uint256 value)
func (withholdingCourt *WithholdingCourt) UnpackSafeCastOverflowedUintDowncastError(raw []byte) (*WithholdingCourtSafeCastOverflowedUintDowncast, error) {
	out := new(WithholdingCourtSafeCastOverflowedUintDowncast)
	if err := withholdingCourt.abi.UnpackIntoInterface(out, "SafeCastOverflowedUintDowncast", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// WithholdingCourtStringTooLong represents a StringTooLong error raised by the WithholdingCourt contract.
type WithholdingCourtStringTooLong struct {
	Str string
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error StringTooLong(string str)
func WithholdingCourtStringTooLongErrorID() common.Hash {
	return common.HexToHash("0x305a27a93f8e33b7392df0a0f91d6fc63847395853c45991eec52dbf24d72381")
}

// UnpackStringTooLongError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error StringTooLong(string str)
func (withholdingCourt *WithholdingCourt) UnpackStringTooLongError(raw []byte) (*WithholdingCourtStringTooLong, error) {
	out := new(WithholdingCourtStringTooLong)
	if err := withholdingCourt.abi.UnpackIntoInterface(out, "StringTooLong", raw); err != nil {
		return nil, err
	}
	return out, nil
}
