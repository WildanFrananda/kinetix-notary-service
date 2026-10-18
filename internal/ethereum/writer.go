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

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/WildanFrananda/kinetix-notary-service/internal/contractgen/council"
	"github.com/WildanFrananda/kinetix-notary-service/internal/contractgen/ledger"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type Writer struct {
	*transactor

	client         *ethclient.Client
	binding        *ledger.LiabilityLedger
	councilBinding *council.AnchorCouncil
	contract       *bind.BoundContract
	timeout        time.Duration
}

var _ notary.AnchorChain = (*Writer)(nil)

func NewWriter(client *ethclient.Client, settings Settings, key *ecdsa.PrivateKey) *Writer {
	binding := ledger.NewLiabilityLedger()

	return &Writer{
		transactor: newTransactor(
			client,
			settings.Ledger.ChainID,
			settings.Council,
			key,
			settings.Timeout,
			"an anchor through this council",
		),
		client:         client,
		binding:        binding,
		councilBinding: council.NewAnchorCouncil(),
		contract:       binding.Instance(client, common.Address(settings.Ledger.Address)),
		timeout:        settings.Timeout,
	}
}

func ParseAnchorerKey(text string) (*ecdsa.PrivateKey, error) {
	return parseKey("ethereum.ParseAnchorerKey", "anchorer", text)
}

func ParseAnswererKey(text string) (*ecdsa.PrivateKey, error) {
	return parseKey("ethereum.ParseAnswererKey", "answerer", text)
}

func parseKey(op, name, text string) (*ecdsa.PrivateKey, error) {
	digits, found := strings.CutPrefix(strings.TrimSpace(text), "0x")
	if !found || len(digits) != 64 {
		return nil, notary.Errf(notary.KindInvalidArgument, op, nil, "the %s key is not 0x and 64 hex digits", name)
	}

	key, err := crypto.HexToECDSA(digits)
	if err != nil {
		return nil, notary.Errf(notary.KindInvalidArgument, op, nil, "the %s key is not a secp256k1 private key", name)
	}

	return key, nil
}

func (w *Writer) Head(ctx context.Context) (notary.LedgerHead, error) {
	return readHead(ctx, "ethereum.Writer.Head", w.contract, w.binding, w.timeout)
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

func (w *Writer) Prepare(
	ctx context.Context,
	header notary.EpochHeader,
	approvals []notary.Signature,
	nonce uint64,
	fees notary.FeeQuote,
) (notary.PreparedAnchor, error) {
	const op = "ethereum.Writer.Prepare"

	signatures := make([][]byte, len(approvals))
	for i := range approvals {
		signatures[i] = approvals[i][:]
	}

	calldata := w.councilBinding.PackAnchor(
		council.EpochHeader{
			Number:         header.Number,
			PreviousDigest: header.PreviousDigest,
			Root:           header.Root,
			Total:          header.Total,
			LeafCount:      header.LeafCount,
			AsOf:           header.AsOf.Seconds(),
		},
		signatures,
	)

	estimate, err := w.estimate(ctx, calldata, fees)
	if err != nil {
		return notary.PreparedAnchor{}, w.refusal(op, header.Number, err)
	}

	return w.sign(op, nonce, estimate*gasMarginPercent/100, calldata, fees)
}

func (w *Writer) refusal(op string, number uint64, err error) error {
	return refusedBy(op, refusals, fmt.Sprintf("the ledger refused epoch %d", number), err)
}

func selectorIs(err error, id common.Hash) bool {
	raw, ok := revertData(err)
	return ok && bytes.Equal(raw[:selectorLength], id.Bytes()[:selectorLength])
}
