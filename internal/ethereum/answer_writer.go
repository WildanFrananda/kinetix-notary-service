package ethereum

import (
	"bytes"
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

	"github.com/WildanFrananda/kinetix-notary-service/internal/contractgen/withholding"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

const (
	answerWait = 90 * time.Second
	answerPoll = 2 * time.Second
)

type AnswerWriter struct {
	client  *ethclient.Client
	binding *withholding.WithholdingCourt
	court   common.Address
	chainID *big.Int
	key     *ecdsa.PrivateKey
	from    common.Address
	signer  types.Signer
	timeout time.Duration
	ceiling uint64
	wait    time.Duration
	poll    time.Duration
	sent    *answerSequence
}

type sentAnswer struct {
	pending bool
	nonce   uint64
	fees    notary.FeeQuote
}

var _ notary.AnswerChain = (*AnswerWriter)(nil)

func NewAnswerWriter(client *ethclient.Client, settings Settings, key *ecdsa.PrivateKey, ceiling uint64) *AnswerWriter {
	chainID := new(big.Int).SetUint64(settings.Ledger.ChainID)

	return &AnswerWriter{
		client:  client,
		binding: withholding.NewWithholdingCourt(),
		court:   common.Address(settings.Withholding),
		chainID: chainID,
		key:     key,
		from:    crypto.PubkeyToAddress(key.PublicKey),
		signer:  types.LatestSignerForChainID(chainID),
		timeout: settings.Timeout,
		ceiling: ceiling,
		wait:    answerWait,
		poll:    answerPoll,
		sent:    &answerSequence{},
	}
}

func (w *AnswerWriter) At(withholding notary.Address) *AnswerWriter {
	sibling := *w
	sibling.court = common.Address(withholding)

	return &sibling
}

func (w *AnswerWriter) Answerer() notary.Address {
	return notary.Address(w.from)
}

func (w *AnswerWriter) Balance(ctx context.Context) (*big.Int, error) {
	ctx, cancel := context.WithTimeout(ctx, w.timeout)
	defer cancel()

	balance, err := w.client.BalanceAt(ctx, w.from, nil)
	if err != nil {
		return nil, unavailable("ethereum.AnswerWriter.Balance", err)
	}

	return balance, nil
}

func (w *AnswerWriter) Answer(ctx context.Context, answer notary.Answer) (notary.Digest, error) {
	const op = "ethereum.AnswerWriter.Answer"

	w.sent.mu.Lock()
	defer w.sent.mu.Unlock()

	siblings := make([]withholding.SumNode, len(answer.Proof.Siblings))
	for i, node := range answer.Proof.Siblings {
		siblings[i] = withholding.SumNode{Hash: node.Hash, Sum: node.Sum}
	}

	calldata := w.binding.PackAnswer(
		boundRegistration(answer.Registration),
		answer.Epoch,
		answer.Salt,
		answer.Amount,
		answer.Proof.Index,
		siblings,
	)

	tx, err := w.prepare(ctx, op, answer.Epoch, calldata)
	if err != nil {
		return notary.Digest{}, err
	}

	if err := w.send(ctx, op, tx); err != nil {
		return notary.Digest{}, err
	}

	w.sent.last = sentAnswer{
		pending: true,
		nonce:   tx.Nonce(),
		fees:    notary.FeeQuote{MaxFeePerGas: tx.GasFeeCap().Uint64(), MaxPriorityFeePerGas: tx.GasTipCap().Uint64()},
	}

	return w.await(ctx, op, tx.Hash())
}

func (w *AnswerWriter) prepare(
	ctx context.Context,
	op string,
	epoch uint64,
	calldata []byte,
) (*types.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, w.timeout)
	defer cancel()

	nonce, err := w.client.NonceAt(ctx, w.from, nil)
	if err != nil {
		return nil, unavailable(op, err)
	}

	header, err := w.client.HeaderByNumber(ctx, nil)
	if err != nil {
		return nil, unavailable(op, err)
	}

	if header.BaseFee == nil || !header.BaseFee.IsUint64() {
		return nil, notary.Errf(notary.KindWrongLedger, op, nil, "the chain reports no usable EIP-1559 base fee")
	}

	tip, err := w.client.SuggestGasTipCap(ctx)
	if err != nil {
		return nil, unavailable(op, err)
	}

	if !tip.IsUint64() {
		return nil, notary.Errf(notary.KindLedgerUnavailable, op, nil, "the suggested tip %s does not fit 64 bits", tip)
	}

	fees, err := notary.QuoteFees(header.BaseFee.Uint64(), tip.Uint64(), w.ceiling)
	if err != nil {
		return nil, err
	}

	if w.sent.last.pending && w.sent.last.nonce == nonce {
		if bumped, ok := w.sent.last.fees.Bumped(w.ceiling); ok {
			fees.MaxFeePerGas = max(fees.MaxFeePerGas, bumped.MaxFeePerGas)
			fees.MaxPriorityFeePerGas = max(fees.MaxPriorityFeePerGas, bumped.MaxPriorityFeePerGas)
		}
	}

	feeCap := new(big.Int).SetUint64(fees.MaxFeePerGas)
	tipCap := new(big.Int).SetUint64(fees.MaxPriorityFeePerGas)

	estimate, err := w.client.EstimateGas(ctx, goethereum.CallMsg{
		From:      w.from,
		To:        &w.court,
		GasFeeCap: feeCap,
		GasTipCap: tipCap,
		Data:      calldata,
	})
	if err != nil {
		return nil, answerRefusal(op, epoch, err)
	}

	tx, err := types.SignNewTx(w.key, w.signer, &types.DynamicFeeTx{
		ChainID:   w.chainID,
		Nonce:     nonce,
		GasTipCap: tipCap,
		GasFeeCap: feeCap,
		Gas:       estimate * gasMarginPercent / 100,
		To:        &w.court,
		Data:      calldata,
	})
	if err != nil {
		return nil, notary.Errf(notary.KindInvalidArgument, op, err, "signing nonce %d", nonce)
	}

	return tx, nil
}

func (w *AnswerWriter) send(ctx context.Context, op string, tx *types.Transaction) error {
	ctx, cancel := context.WithTimeout(ctx, w.timeout)
	defer cancel()

	err := w.client.SendTransaction(ctx, tx)
	if err == nil {
		return nil
	}

	message := strings.ToLower(err.Error())
	switch {
	case strings.Contains(message, "already known"), strings.Contains(message, "known transaction"):
		return nil
	case strings.Contains(message, "nonce too low"), strings.Contains(message, "underpriced"),
		strings.Contains(message, "insufficient funds"), strings.Contains(message, "fee cap less than block base fee"):
		return notary.Errf(notary.KindAnchorDeferred, op, err, "the node refused the answer %s", tx.Hash().Hex())
	default:
		return unavailable(op, err)
	}
}

func (w *AnswerWriter) await(ctx context.Context, op string, hash common.Hash) (notary.Digest, error) {
	deadline := time.NewTimer(w.wait)
	defer deadline.Stop()

	ticker := time.NewTicker(w.poll)
	defer ticker.Stop()

	for {
		receipt, err := w.receipt(ctx, hash)
		if err != nil && !errors.Is(err, goethereum.NotFound) {
			return notary.Digest(hash), unavailable(op, err)
		}

		if receipt != nil {
			w.sent.last.pending = false

			if receipt.Status != types.ReceiptStatusSuccessful {
				return notary.Digest(
						hash,
					), notary.Errf(
						notary.KindAnswerRefused,
						op,
						nil,
						"the answer %s reverted",
						hash.Hex(),
					)
			}

			return notary.Digest(hash), nil
		}

		select {
		case <-ctx.Done():
			return notary.Digest(hash), unavailable(op, ctx.Err())
		case <-deadline.C:
			return notary.Digest(hash), notary.Errf(
				notary.KindAnchorDeferred,
				op,
				nil,
				"the answer %s is not mined after %s",
				hash.Hex(),
				w.wait,
			)
		case <-ticker.C:
		}
	}
}

func (w *AnswerWriter) receipt(ctx context.Context, hash common.Hash) (*types.Receipt, error) {
	ctx, cancel := context.WithTimeout(ctx, w.timeout)
	defer cancel()

	receipt, err := w.client.TransactionReceipt(ctx, hash)
	if err != nil {
		return nil, fmt.Errorf("receipt of %s: %w", hash.Hex(), err)
	}

	return receipt, nil
}

var answerRefusals = []struct {
	id   common.Hash
	name string
}{
	{withholding.WithholdingCourtNoOpenRequestErrorID(), "NoOpenRequest"},
	{withholding.WithholdingCourtAnswerTooLateErrorID(), "AnswerTooLate"},
	{withholding.WithholdingCourtAnswerDoesNotHoldErrorID(), "AnswerDoesNotHold"},
}

func answerRefusal(op string, epoch uint64, err error) error {
	raw, reverted := revertData(err)
	if !reverted {
		if isBareRevert(err) {
			return notary.Errf(notary.KindAnswerRefused, op, err, "the court refused the answer for epoch %d", epoch)
		}

		return unavailable(op, err)
	}

	for _, rule := range answerRefusals {
		if bytes.Equal(raw[:selectorLength], rule.id.Bytes()[:selectorLength]) {
			return notary.Errf(
				notary.KindAnswerRefused,
				op,
				nil,
				"the court refused the answer for epoch %d: %s",
				epoch,
				rule.name,
			)
		}
	}

	return notary.Errf(
		notary.KindAnswerRefused,
		op,
		nil,
		"the court refused the answer for epoch %d with %s",
		epoch,
		hexutil.Encode(raw[:selectorLength]),
	)
}
