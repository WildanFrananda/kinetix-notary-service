package ethereum

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"errors"
	"math/big"
	"strings"
	"time"

	goethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/WildanFrananda/kinetix-notary-service/internal/contractgen/ledger"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

const gasMarginPercent = 120

type Writer struct {
	client   *ethclient.Client
	binding  *ledger.LiabilityLedger
	contract *bind.BoundContract
	ledger   notary.LedgerIdentity
	key      *ecdsa.PrivateKey
	from     common.Address
	signer   types.Signer
	timeout  time.Duration
}

var _ notary.AnchorChain = (*Writer)(nil)

func NewWriter(client *ethclient.Client, settings Settings, key *ecdsa.PrivateKey) *Writer {
	binding := ledger.NewLiabilityLedger()
	chainID := new(big.Int).SetUint64(settings.Ledger.ChainID)

	return &Writer{
		client:   client,
		binding:  binding,
		contract: binding.Instance(client, common.Address(settings.Ledger.Address)),
		ledger:   settings.Ledger,
		key:      key,
		from:     crypto.PubkeyToAddress(key.PublicKey),
		signer:   types.LatestSignerForChainID(chainID),
		timeout:  settings.Timeout,
	}
}

func ParseAnchorerKey(text string) (*ecdsa.PrivateKey, error) {
	digits, found := strings.CutPrefix(strings.TrimSpace(text), "0x")
	if !found || len(digits) != 64 {
		return nil, notary.Errf(
			notary.KindInvalidArgument,
			"ethereum.ParseAnchorerKey",
			nil,
			"the anchorer key is not 0x and 64 hex digits",
		)
	}

	key, err := crypto.HexToECDSA(digits)
	if err != nil {
		return nil, notary.Errf(
			notary.KindInvalidArgument,
			"ethereum.ParseAnchorerKey",
			nil,
			"the anchorer key is not a secp256k1 private key",
		)
	}

	return key, nil
}

func (w *Writer) Anchorer() notary.Address {
	return notary.Address(w.from)
}

func (w *Writer) VerifyAnchorer(ctx context.Context) error {
	const op = "ethereum.Writer.VerifyAnchorer"

	ctx, cancel := context.WithTimeout(ctx, w.timeout)
	defer cancel()

	opts := &bind.CallOpts{Context: ctx}

	role, err := bind.Call(w.contract, opts, w.binding.PackANCHORERROLE(), w.binding.UnpackANCHORERROLE)
	if err != nil {
		return unavailable(op, err)
	}

	holds, err := bind.Call(w.contract, opts, w.binding.PackHasRole(role, w.from), w.binding.UnpackHasRole)
	if err != nil {
		return unavailable(op, err)
	}

	if !holds {
		return notary.Errf(
			notary.KindNotAnchorer,
			op,
			nil,
			"%s does not hold the anchorer role on %s",
			w.from.Hex(),
			w.ledger.Address.Hex(),
		)
	}

	return nil
}

func (w *Writer) Head(ctx context.Context) (notary.LedgerHead, error) {
	const op = "ethereum.Writer.Head"

	ctx, cancel := context.WithTimeout(ctx, w.timeout)
	defer cancel()

	opts := &bind.CallOpts{Context: ctx}

	var head notary.LedgerHead

	next, err := bind.Call(w.contract, opts, w.binding.PackNextEpoch(), w.binding.UnpackNextEpoch)
	if err != nil {
		return head, unavailable(op, err)
	}

	digest, err := bind.Call(w.contract, opts, w.binding.PackHead(), w.binding.UnpackHead)
	if err != nil {
		return head, unavailable(op, err)
	}

	lastAsOf, err := bind.Call(w.contract, opts, w.binding.PackLastAsOf(), w.binding.UnpackLastAsOf)
	if err != nil {
		return head, unavailable(op, err)
	}

	genesisEpoch, err := bind.Call(w.contract, opts, w.binding.PackGenesisEpoch(), w.binding.UnpackGenesisEpoch)
	if err != nil {
		return head, unavailable(op, err)
	}

	genesisDigest, err := bind.Call(w.contract, opts, w.binding.PackGenesisDigest(), w.binding.UnpackGenesisDigest)
	if err != nil {
		return head, unavailable(op, err)
	}

	sealed, err := bind.Call(w.contract, opts, w.binding.PackIsSealed(), w.binding.UnpackIsSealed)
	if err != nil {
		return head, unavailable(op, err)
	}

	return notary.LedgerHead{
		NextEpoch:     next,
		Head:          digest,
		LastAsOf:      lastAsOf,
		GenesisEpoch:  genesisEpoch,
		GenesisDigest: genesisDigest,
		Sealed:        sealed,
	}, nil
}

func (w *Writer) EpochAtLatest(ctx context.Context, number uint64) (notary.AnchoredEpoch, bool, error) {
	return w.epochAt(ctx, "ethereum.Writer.EpochAtLatest", number, nil)
}

func (w *Writer) EpochAtFinalized(ctx context.Context, number uint64) (notary.AnchoredEpoch, bool, error) {
	const op = "ethereum.Writer.EpochAtFinalized"

	ctx, cancel := context.WithTimeout(ctx, w.timeout)
	defer cancel()

	header, err := w.client.HeaderByNumber(ctx, big.NewInt(int64(rpc.FinalizedBlockNumber)))
	if err != nil {
		return notary.AnchoredEpoch{}, false, unavailable(op, err)
	}

	return w.epochAt(ctx, op, number, header.Number)
}

func (w *Writer) epochAt(
	ctx context.Context,
	op string,
	number uint64,
	block *big.Int,
) (notary.AnchoredEpoch, bool, error) {
	ctx, cancel := context.WithTimeout(ctx, w.timeout)
	defer cancel()

	epoch, err := bind.Call(
		w.contract,
		&bind.CallOpts{Context: ctx, BlockNumber: block},
		w.binding.PackEpochAt(number),
		w.binding.UnpackEpochAt,
	)
	if err != nil {
		if selectorIs(err, ledger.LiabilityLedgerUnknownEpochErrorID()) ||
			(block != nil && errors.Is(err, bind.ErrNoCode)) {
			return notary.AnchoredEpoch{}, false, nil
		}

		return notary.AnchoredEpoch{}, false, unavailable(op, err)
	}

	asOf, err := notary.UnixTimeFromSeconds(epoch.AsOf)
	if err != nil {
		return notary.AnchoredEpoch{}, false, notary.Errf(notary.KindWrongLedger, op, err, "epoch %d", number)
	}

	anchoredAt, err := notary.UnixTimeFromSeconds(epoch.AnchoredAt)
	if err != nil {
		return notary.AnchoredEpoch{}, false, notary.Errf(notary.KindWrongLedger, op, err, "epoch %d", number)
	}

	return notary.AnchoredEpoch{
		Number:     number,
		Root:       epoch.Root,
		Digest:     epoch.Digest,
		Total:      epoch.Total,
		LeafCount:  epoch.LeafCount,
		AsOf:       asOf,
		AnchoredAt: anchoredAt,
	}, true, nil
}

func (w *Writer) Nonce(ctx context.Context) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, w.timeout)
	defer cancel()

	nonce, err := w.client.NonceAt(ctx, w.from, nil)
	if err != nil {
		return 0, unavailable("ethereum.Writer.Nonce", err)
	}

	return nonce, nil
}

func (w *Writer) Fees(ctx context.Context) (uint64, uint64, error) {
	const op = "ethereum.Writer.Fees"

	ctx, cancel := context.WithTimeout(ctx, w.timeout)
	defer cancel()

	header, err := w.client.HeaderByNumber(ctx, nil)
	if err != nil {
		return 0, 0, unavailable(op, err)
	}

	if header.BaseFee == nil || !header.BaseFee.IsUint64() {
		return 0, 0, notary.Errf(notary.KindWrongLedger, op, nil, "the chain reports no usable EIP-1559 base fee")
	}

	tip, err := w.client.SuggestGasTipCap(ctx)
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

func (w *Writer) Prepare(
	ctx context.Context,
	header notary.EpochHeader,
	nonce uint64,
	fees notary.FeeQuote,
) (notary.PreparedAnchor, error) {
	const op = "ethereum.Writer.Prepare"

	ctx, cancel := context.WithTimeout(ctx, w.timeout)
	defer cancel()

	calldata := w.binding.PackAnchor(ledger.EpochHeader{
		Number:         header.Number,
		PreviousDigest: header.PreviousDigest,
		Root:           header.Root,
		Total:          header.Total,
		LeafCount:      header.LeafCount,
		AsOf:           header.AsOf.Seconds(),
	})
	to := common.Address(w.ledger.Address)
	feeCap := new(big.Int).SetUint64(fees.MaxFeePerGas)
	tipCap := new(big.Int).SetUint64(fees.MaxPriorityFeePerGas)

	estimate, err := w.client.EstimateGas(ctx, goethereum.CallMsg{
		From:      w.from,
		To:        &to,
		GasFeeCap: feeCap,
		GasTipCap: tipCap,
		Data:      calldata,
	})
	if err != nil {
		return notary.PreparedAnchor{}, w.refusal(op, header.Number, err)
	}

	return w.sign(op, nonce, estimate*gasMarginPercent/100, calldata, fees)
}

func (w *Writer) Resign(previous notary.PreparedAnchor, fees notary.FeeQuote) (notary.PreparedAnchor, error) {
	const op = "ethereum.Writer.Resign"

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

	if tx.To() == nil || *tx.To() != common.Address(w.ledger.Address) || tx.Nonce() != previous.Nonce {
		return notary.PreparedAnchor{}, notary.Errf(
			notary.KindInvalidArgument,
			op,
			nil,
			"%s is not an anchor to this ledger",
			previous.Hash.Hex(),
		)
	}

	return w.sign(op, tx.Nonce(), tx.Gas(), tx.Data(), fees)
}

func (w *Writer) sign(
	op string,
	nonce, gas uint64,
	calldata []byte,
	fees notary.FeeQuote,
) (notary.PreparedAnchor, error) {
	to := common.Address(w.ledger.Address)

	tx, err := types.SignNewTx(w.key, w.signer, &types.DynamicFeeTx{
		ChainID:   new(big.Int).SetUint64(w.ledger.ChainID),
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

func (w *Writer) Broadcast(ctx context.Context, anchor notary.PreparedAnchor) error {
	const op = "ethereum.Writer.Broadcast"

	ctx, cancel := context.WithTimeout(ctx, w.timeout)
	defer cancel()

	var hash common.Hash

	err := w.client.Client().CallContext(ctx, &hash, "eth_sendRawTransaction", hexutil.Encode(anchor.Raw))
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

func (w *Writer) Known(ctx context.Context, hash notary.Digest) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, w.timeout)
	defer cancel()

	_, _, err := w.client.TransactionByHash(ctx, common.Hash(hash))
	if errors.Is(err, goethereum.NotFound) {
		return false, nil
	}

	if err != nil {
		return false, unavailable("ethereum.Writer.Known", err)
	}

	return true, nil
}

func (w *Writer) Receipt(ctx context.Context, hash notary.Digest) (notary.AnchorReceipt, bool, error) {
	const op = "ethereum.Writer.Receipt"

	ctx, cancel := context.WithTimeout(ctx, w.timeout)
	defer cancel()

	var raw *rawReceipt
	if err := w.client.Client().CallContext(ctx, &raw, "eth_getTransactionReceipt", common.Hash(hash)); err != nil {
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

func (w *Writer) Balance(ctx context.Context) (*big.Int, error) {
	ctx, cancel := context.WithTimeout(ctx, w.timeout)
	defer cancel()

	balance, err := w.client.BalanceAt(ctx, w.from, nil)
	if err != nil {
		return nil, unavailable("ethereum.Writer.Balance", err)
	}

	return balance, nil
}

func (w *Writer) refusal(op string, number uint64, err error) error {
	raw, reverted := revertData(err)
	if !reverted {
		if isBareRevert(err) {
			return notary.Errf(
				notary.KindAnchorRefused,
				op,
				err,
				"the ledger refused epoch %d without a reason",
				number,
			)
		}

		return unavailable(op, err)
	}

	for _, rule := range refusals {
		if bytes.Equal(raw[:selectorLength], rule.id.Bytes()[:selectorLength]) {
			return notary.Errf(rule.kind, op, nil, "the ledger refused epoch %d: %s", number, rule.name)
		}
	}

	return notary.Errf(
		notary.KindAnchorRefused,
		op,
		nil,
		"the ledger refused epoch %d with %s",
		number,
		hexutil.Encode(raw[:selectorLength]),
	)
}

func selectorIs(err error, id common.Hash) bool {
	raw, ok := revertData(err)
	return ok && bytes.Equal(raw[:selectorLength], id.Bytes()[:selectorLength])
}
