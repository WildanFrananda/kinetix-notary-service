package notarytest

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/holiman/uint256"

	"github.com/WildanFrananda/kinetix-notary-service/internal/contractgen/entrypoint"
)

type LocalBundler struct {
	URL        string
	client     *ethclient.Client
	key        *ecdsa.PrivateKey
	chainID    *big.Int
	entryPoint common.Address
	binding    *entrypoint.EntryPoint
	mu         sync.Mutex
	receipts   map[common.Hash]map[string]any
}

type bundlerCall struct {
	ID     json.RawMessage   `json:"id"`
	Method string            `json:"method"`
	Params []json.RawMessage `json:"params"`
}

type wireOperation struct {
	Sender                        common.Address  `json:"sender"`
	Nonce                         hexutil.Big     `json:"nonce"`
	Factory                       *string         `json:"factory"`
	FactoryData                   *hexutil.Bytes  `json:"factoryData"`
	CallData                      hexutil.Bytes   `json:"callData"`
	CallGasLimit                  hexutil.Big     `json:"callGasLimit"`
	VerificationGasLimit          hexutil.Big     `json:"verificationGasLimit"`
	PreVerificationGas            hexutil.Big     `json:"preVerificationGas"`
	MaxFeePerGas                  hexutil.Big     `json:"maxFeePerGas"`
	MaxPriorityFeePerGas          hexutil.Big     `json:"maxPriorityFeePerGas"`
	Paymaster                     common.Address  `json:"paymaster"`
	PaymasterVerificationGasLimit hexutil.Big     `json:"paymasterVerificationGasLimit"`
	PaymasterPostOpGasLimit       hexutil.Big     `json:"paymasterPostOpGasLimit"`
	PaymasterData                 hexutil.Bytes   `json:"paymasterData"`
	Signature                     hexutil.Bytes   `json:"signature"`
	Authorization                 *wireDelegation `json:"eip7702Auth"`
}

type wireDelegation struct {
	ChainID hexutil.Big    `json:"chainId"`
	Address common.Address `json:"address"`
	Nonce   hexutil.Uint64 `json:"nonce"`
	YParity hexutil.Uint64 `json:"yParity"`
	R       hexutil.Big    `json:"r"`
	S       hexutil.Big    `json:"s"`
}

func StartLocalBundler(
	t *testing.T,
	client *ethclient.Client,
	key *ecdsa.PrivateKey,
	entryPoint common.Address,
) *LocalBundler {
	t.Helper()

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		t.Fatalf("chain id: %v", err)
	}

	b := &LocalBundler{
		client:     client,
		key:        key,
		chainID:    chainID,
		entryPoint: entryPoint,
		binding:    entrypoint.NewEntryPoint(),
		receipts:   map[common.Hash]map[string]any{},
	}
	server := httptest.NewServer(http.HandlerFunc(b.serve))
	t.Cleanup(server.Close)
	b.URL = server.URL

	return b
}

func (b *LocalBundler) serve(w http.ResponseWriter, r *http.Request) {
	var c bundlerCall
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, "not JSON-RPC", http.StatusBadRequest)
		return
	}

	result, code, message := b.answer(r.Context(), c)
	reply := map[string]any{"jsonrpc": "2.0", "id": c.ID}

	if message != "" {
		reply["error"] = map[string]any{"code": code, "message": message}
	} else {
		reply["result"] = result
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(reply)
}

func (b *LocalBundler) answer(ctx context.Context, c bundlerCall) (any, int, string) {
	switch c.Method {
	case "eth_estimateUserOperationGas":
		return map[string]string{
			"preVerificationGas":            "0x186a0",
			"verificationGasLimit":          "0x493e0",
			"callGasLimit":                  "0x7a120",
			"paymasterVerificationGasLimit": "0x30d40",
		}, 0, ""
	case "eth_sendUserOperation":
		var op wireOperation
		if err := json.Unmarshal(c.Params[0], &op); err != nil {
			return nil, -32602, err.Error()
		}

		hash, err := b.include(ctx, op)
		if err != nil {
			return nil, -32500, err.Error()
		}

		return hash.Hex(), 0, ""
	case "eth_getUserOperationReceipt":
		var hash string
		if err := json.Unmarshal(c.Params[0], &hash); err != nil {
			return nil, -32602, err.Error()
		}

		b.mu.Lock()
		defer b.mu.Unlock()

		if receipt, ok := b.receipts[common.HexToHash(hash)]; ok {
			return receipt, 0, ""
		}

		return nil, 0, ""
	default:
		return nil, -32601, "method not found"
	}
}

func (b *LocalBundler) include(ctx context.Context, op wireOperation) (common.Hash, error) {
	data := b.binding.PackHandleOps(
		[]entrypoint.PackedUserOperation{packed(op)},
		crypto.PubkeyToAddress(b.key.PublicKey),
	)
	from := crypto.PubkeyToAddress(b.key.PublicKey)

	nonce, err := b.client.PendingNonceAt(ctx, from)
	if err != nil {
		return common.Hash{}, err
	}

	price, err := b.client.SuggestGasPrice(ctx)
	if err != nil {
		return common.Hash{}, err
	}

	fee := new(big.Int).Mul(price, big.NewInt(2))

	var tx *types.Transaction

	if op.Authorization != nil {
		tx = types.NewTx(&types.SetCodeTx{
			ChainID:   uint256.MustFromBig(b.chainID),
			Nonce:     nonce,
			GasTipCap: uint256.MustFromBig(price),
			GasFeeCap: uint256.MustFromBig(fee),
			Gas:       3_000_000,
			To:        b.entryPoint,
			Data:      data,
			AuthList: []types.SetCodeAuthorization{{
				ChainID: *uint256.MustFromBig(op.Authorization.ChainID.ToInt()),
				Address: op.Authorization.Address,
				Nonce:   uint64(op.Authorization.Nonce),
				V:       uint8(op.Authorization.YParity),
				R:       *uint256.MustFromBig(op.Authorization.R.ToInt()),
				S:       *uint256.MustFromBig(op.Authorization.S.ToInt()),
			}},
		})
	} else {
		to := b.entryPoint
		tx = types.NewTx(&types.DynamicFeeTx{
			ChainID:   b.chainID,
			Nonce:     nonce,
			GasTipCap: price,
			GasFeeCap: fee,
			Gas:       3_000_000,
			To:        &to,
			Data:      data,
		})
	}

	signed, err := types.SignTx(tx, types.LatestSignerForChainID(b.chainID), b.key)
	if err != nil {
		return common.Hash{}, err
	}

	if err := b.client.SendTransaction(ctx, signed); err != nil {
		return common.Hash{}, err
	}

	receipt, err := bind.WaitMined(ctx, b.client, signed.Hash())
	if err != nil {
		return common.Hash{}, err
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		return common.Hash{}, fmt.Errorf("handleOps reverted in %s", signed.Hash().Hex())
	}

	for _, log := range receipt.Logs {
		event, err := b.binding.UnpackUserOperationEventEvent(log)
		if err != nil {
			continue
		}

		hash := common.Hash(event.UserOpHash)

		b.mu.Lock()
		b.receipts[hash] = map[string]any{
			"userOpHash":    hash.Hex(),
			"success":       event.Success,
			"actualGasCost": "0x" + event.ActualGasCost.Text(16),
			"reason":        "",
			"receipt":       map[string]string{"transactionHash": signed.Hash().Hex()},
		}
		b.mu.Unlock()

		return hash, nil
	}

	return common.Hash{}, fmt.Errorf("no UserOperationEvent in %s", signed.Hash().Hex())
}

func packed(op wireOperation) entrypoint.PackedUserOperation {
	var initCode []byte

	if op.Factory != nil {
		marker := common.FromHex(*op.Factory)
		initCode = append(append([]byte{}, marker...), make([]byte, common.AddressLength-len(marker))...)

		if op.FactoryData != nil {
			initCode = append(initCode, *op.FactoryData...)
		}
	}

	verificationGas := common.LeftPadBytes(op.PaymasterVerificationGasLimit.ToInt().Bytes(), 16)
	postOpGas := common.LeftPadBytes(op.PaymasterPostOpGasLimit.ToInt().Bytes(), 16)
	paymasterAndData := append([]byte{}, op.Paymaster.Bytes()...)
	paymasterAndData = append(paymasterAndData, verificationGas...)
	paymasterAndData = append(paymasterAndData, postOpGas...)
	paymasterAndData = append(paymasterAndData, op.PaymasterData...)

	return entrypoint.PackedUserOperation{
		Sender:             op.Sender,
		Nonce:              op.Nonce.ToInt(),
		InitCode:           initCode,
		CallData:           op.CallData,
		AccountGasLimits:   pair(op.VerificationGasLimit.ToInt(), op.CallGasLimit.ToInt()),
		PreVerificationGas: op.PreVerificationGas.ToInt(),
		GasFees:            pair(op.MaxPriorityFeePerGas.ToInt(), op.MaxFeePerGas.ToInt()),
		PaymasterAndData:   paymasterAndData,
		Signature:          op.Signature,
	}
}

func pair(high, low *big.Int) [32]byte {
	var word [32]byte
	copy(word[:16], common.LeftPadBytes(high.Bytes(), 16))
	copy(word[16:], common.LeftPadBytes(low.Bytes(), 16))

	return word
}
