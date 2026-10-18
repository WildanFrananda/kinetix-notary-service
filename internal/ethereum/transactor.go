package ethereum

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	goethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

const gasMarginPercent = 120

type transactor struct {
	client  *ethclient.Client
	key     *ecdsa.PrivateKey
	from    common.Address
	target  common.Address
	chainID uint64
	signer  types.Signer
	timeout time.Duration
	purpose string
}

func newTransactor(
	client *ethclient.Client,
	chainID uint64,
	target notary.Address,
	key *ecdsa.PrivateKey,
	timeout time.Duration,
	purpose string,
) *transactor {
	return &transactor{
		client:  client,
		key:     key,
		from:    crypto.PubkeyToAddress(key.PublicKey),
		target:  common.Address(target),
		chainID: chainID,
		signer:  types.LatestSignerForChainID(new(big.Int).SetUint64(chainID)),
		timeout: timeout,
		purpose: purpose,
	}
}

func (t *transactor) Anchorer() notary.Address {
	return notary.Address(t.from)
}

func (t *transactor) Nonce(ctx context.Context) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, t.timeout)
	defer cancel()

	nonce, err := t.client.NonceAt(ctx, t.from, nil)
	if err != nil {
		return 0, unavailable("ethereum.transactor.Nonce", err)
	}

	return nonce, nil
}

func (t *transactor) Fees(ctx context.Context) (uint64, uint64, error) {
	const op = "ethereum.transactor.Fees"

	ctx, cancel := context.WithTimeout(ctx, t.timeout)
	defer cancel()

	header, err := t.client.HeaderByNumber(ctx, nil)
	if err != nil {
		return 0, 0, unavailable(op, err)
	}

	if header.BaseFee == nil || !header.BaseFee.IsUint64() {
		return 0, 0, notary.Errf(notary.KindWrongLedger, op, nil, "the chain reports no usable EIP-1559 base fee")
	}

	tip, err := t.client.SuggestGasTipCap(ctx)
	if err != nil {
		return 0, 0, unavailable(op, err)
	}

	if !tip.IsUint64() {
		return 0, 0, notary.Errf(
			notary.KindLedgerUnavailable,
			op,
			nil,
			"the suggested tip %s does not fit 64 bits",
			tip,
		)
	}

	return header.BaseFee.Uint64(), tip.Uint64(), nil
}

func (t *transactor) estimate(ctx context.Context, calldata []byte, fees notary.FeeQuote) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, t.timeout)
	defer cancel()

	to := t.target

	gas, err := t.client.EstimateGas(ctx, goethereum.CallMsg{
		From:      t.from,
		To:        &to,
		GasFeeCap: new(big.Int).SetUint64(fees.MaxFeePerGas),
		GasTipCap: new(big.Int).SetUint64(fees.MaxPriorityFeePerGas),
		Data:      calldata,
	})
	if err != nil {
		return 0, fmt.Errorf("estimate gas: %w", err)
	}

	return gas, nil
}

func (t *transactor) Resign(previous notary.PreparedAnchor, fees notary.FeeQuote) (notary.PreparedAnchor, error) {
	const op = "ethereum.transactor.Resign"

	var tx types.Transaction
	if err := tx.UnmarshalBinary(previous.Raw); err != nil {
		return notary.PreparedAnchor{}, notary.Errf(
			notary.KindInvalidArgument,
			op,
			err,
			"decoding %s",
			previous.Hash.Hex(),
		)
	}

	if tx.To() == nil || *tx.To() != t.target || tx.Nonce() != previous.Nonce {
		return notary.PreparedAnchor{}, notary.Errf(
			notary.KindInvalidArgument,
			op,
			nil,
			"%s is not %s",
			previous.Hash.Hex(),
			t.purpose,
		)
	}

	return t.sign(op, tx.Nonce(), tx.Gas(), tx.Data(), fees)
}

func (t *transactor) sign(
	op string,
	nonce, gas uint64,
	calldata []byte,
	fees notary.FeeQuote,
) (notary.PreparedAnchor, error) {
	to := t.target

	tx, err := types.SignNewTx(t.key, t.signer, &types.DynamicFeeTx{
		ChainID:   new(big.Int).SetUint64(t.chainID),
		Nonce:     nonce,
		GasTipCap: new(big.Int).SetUint64(fees.MaxPriorityFeePerGas),
		GasFeeCap: new(big.Int).SetUint64(fees.MaxFeePerGas),
		Gas:       gas,
		To:        &to,
		Data:      calldata,
	})
	if err != nil {
		return notary.PreparedAnchor{}, notary.Errf(notary.KindInvalidArgument, op, err, "signing nonce %d", nonce)
	}

	raw, err := tx.MarshalBinary()
	if err != nil {
		return notary.PreparedAnchor{}, notary.Errf(notary.KindInvalidArgument, op, err, "encoding nonce %d", nonce)
	}

	return notary.PreparedAnchor{Hash: notary.Digest(tx.Hash()), Raw: raw, Nonce: nonce, GasLimit: gas, Fees: fees}, nil
}

func (t *transactor) Broadcast(ctx context.Context, anchor notary.PreparedAnchor) error {
	const op = "ethereum.transactor.Broadcast"

	ctx, cancel := context.WithTimeout(ctx, t.timeout)
	defer cancel()

	var hash common.Hash

	err := t.client.Client().CallContext(ctx, &hash, "eth_sendRawTransaction", hexutil.Encode(anchor.Raw))
	if err == nil {
		return nil
	}

	message := strings.ToLower(err.Error())
	switch {
	case strings.Contains(message, "already known"), strings.Contains(message, "known transaction"):
		return nil
	case strings.Contains(message, "nonce too low"), strings.Contains(message, "underpriced"),
		strings.Contains(message, "insufficient funds"), strings.Contains(message, "fee cap less than block base fee"):
		return notary.Errf(notary.KindAnchorDeferred, op, err, "the node refused %s", anchor.Hash.Hex())
	default:
		return unavailable(op, err)
	}
}

func (t *transactor) Known(ctx context.Context, hash notary.Digest) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, t.timeout)
	defer cancel()

	_, _, err := t.client.TransactionByHash(ctx, common.Hash(hash))
	if errors.Is(err, goethereum.NotFound) {
		return false, nil
	}

	if err != nil {
		return false, unavailable("ethereum.transactor.Known", err)
	}

	return true, nil
}

func (t *transactor) Receipt(ctx context.Context, hash notary.Digest) (notary.AnchorReceipt, bool, error) {
	const op = "ethereum.transactor.Receipt"

	ctx, cancel := context.WithTimeout(ctx, t.timeout)
	defer cancel()

	var raw *rawReceipt
	if err := t.client.Client().CallContext(ctx, &raw, "eth_getTransactionReceipt", common.Hash(hash)); err != nil {
		return notary.AnchorReceipt{}, false, unavailable(op, err)
	}

	if raw == nil {
		return notary.AnchorReceipt{}, false, nil
	}

	receipt, err := raw.decode(hash)
	if err != nil {
		return notary.AnchorReceipt{}, false, notary.Errf(
			notary.KindLedgerUnavailable,
			op,
			err,
			"receipt of %s",
			hash.Hex(),
		)
	}

	return receipt, true, nil
}

func (t *transactor) Balance(ctx context.Context) (*big.Int, error) {
	ctx, cancel := context.WithTimeout(ctx, t.timeout)
	defer cancel()

	balance, err := t.client.BalanceAt(ctx, t.from, nil)
	if err != nil {
		return nil, unavailable("ethereum.transactor.Balance", err)
	}

	return balance, nil
}
