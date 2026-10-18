package ethereum

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/WildanFrananda/kinetix-notary-service/internal/contractgen/evidence"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type EvidenceWriter struct {
	*transactor

	client   *ethclient.Client
	binding  *evidence.EvidenceRegistry
	contract *bind.BoundContract
	registry notary.RegistryIdentity
	timeout  time.Duration
}

var _ notary.EvidenceChain = (*EvidenceWriter)(nil)

func NewEvidenceWriter(
	client *ethclient.Client,
	registry notary.RegistryIdentity,
	timeout time.Duration,
	key *ecdsa.PrivateKey,
) *EvidenceWriter {
	binding := evidence.NewEvidenceRegistry()

	return &EvidenceWriter{
		transactor: newTransactor(
			client,
			registry.ChainID,
			registry.Address,
			key,
			timeout,
			"an anchor into this evidence registry",
		),
		client:   client,
		binding:  binding,
		contract: binding.Instance(client, common.Address(registry.Address)),
		registry: registry,
		timeout:  timeout,
	}
}

func ParseEvidenceAnchorerKey(text string) (*ecdsa.PrivateKey, error) {
	return parseKey("ethereum.ParseEvidenceAnchorerKey", "evidence anchorer", text)
}

func (w *EvidenceWriter) HoldsAnchorerRole(ctx context.Context) (bool, error) {
	const op = "ethereum.EvidenceWriter.HoldsAnchorerRole"

	ctx, cancel := context.WithTimeout(ctx, w.timeout)
	defer cancel()

	opts := &bind.CallOpts{Context: ctx}

	role, err := bind.Call(w.contract, opts, w.binding.PackANCHORERROLE(), w.binding.UnpackANCHORERROLE)
	if err != nil {
		return false, unavailable(op, err)
	}

	held, err := bind.Call(w.contract, opts, w.binding.PackHasRole(role, w.from), w.binding.UnpackHasRole)
	if err != nil {
		return false, unavailable(op, err)
	}

	return held, nil
}

func (w *EvidenceWriter) StreamHead(
	ctx context.Context,
	stream notary.EvidenceStream,
) (notary.EvidenceStreamHead, error) {
	const op = "ethereum.EvidenceWriter.StreamHead"

	ctx, cancel := context.WithTimeout(ctx, w.timeout)
	defer cancel()

	header, err := w.client.HeaderByNumber(ctx, nil)
	if err != nil {
		return notary.EvidenceStreamHead{}, unavailable(op, err)
	}

	blockTime, err := notary.UnixTimeFromSeconds(header.Time)
	if err != nil {
		return notary.EvidenceStreamHead{}, notary.Errf(
			notary.KindLedgerUnavailable,
			op,
			err,
			"the latest block's time",
		)
	}

	opts := &bind.CallOpts{Context: ctx, BlockNumber: header.Number}

	sealed, err := bind.Call(w.contract, opts, w.binding.PackIsSealed(), w.binding.UnpackIsSealed)
	if err != nil {
		return notary.EvidenceStreamHead{}, unavailable(op, err)
	}

	state, err := bind.Call(w.contract, opts, w.binding.PackStateOf(uint8(stream)), w.binding.UnpackStateOf)
	if err != nil {
		if selectorIs(err, evidence.EvidenceRegistryUnknownStreamErrorID()) {
			return notary.EvidenceStreamHead{Sealed: sealed, BlockTime: blockTime}, nil
		}

		return notary.EvidenceStreamHead{}, unavailable(op, err)
	}

	return notary.EvidenceStreamHead{
		Known:        state.Known,
		Sealed:       sealed,
		Head:         state.Head,
		GenesisBatch: state.GenesisBatch,
		NextBatch:    state.NextBatch,
		LastClosedAt: state.LastClosedAt,
		Through:      notary.LogPosition{Transaction: state.ThroughTransaction, Sequence: state.ThroughSequence},
		BlockTime:    blockTime,
	}, nil
}

func (w *EvidenceWriter) BatchAtLatest(
	ctx context.Context,
	stream notary.EvidenceStream,
	number uint64,
) (notary.AnchoredBatch, bool, error) {
	ctx, cancel := context.WithTimeout(ctx, w.timeout)
	defer cancel()

	return batchAt(ctx, "ethereum.EvidenceWriter.BatchAtLatest", w.contract, w.binding, stream, number, nil)
}

func (w *EvidenceWriter) BatchAtFinalized(
	ctx context.Context,
	stream notary.EvidenceStream,
	number uint64,
) (notary.AnchoredBatch, bool, error) {
	const op = "ethereum.EvidenceWriter.BatchAtFinalized"

	ctx, cancel := context.WithTimeout(ctx, w.timeout)
	defer cancel()

	final, err := finalizedBlock(ctx, w.client)
	if err != nil {
		return notary.AnchoredBatch{}, false, unavailable(op, err)
	}

	return batchAt(ctx, op, w.contract, w.binding, stream, number, final)
}

func (w *EvidenceWriter) PrepareBatch(
	ctx context.Context,
	header notary.EvidenceBatchHeader,
	nonce uint64,
	fees notary.FeeQuote,
) (notary.PreparedAnchor, error) {
	const op = "ethereum.EvidenceWriter.PrepareBatch"

	calldata := w.binding.PackAnchor(evidence.BatchHeader{
		Stream:             uint8(header.Stream),
		Number:             header.Number,
		PreviousDigest:     header.PreviousDigest,
		Root:               header.Root,
		LeafCount:          header.LeafCount,
		ThroughTransaction: header.Through.Transaction,
		ThroughSequence:    header.Through.Sequence,
		ClosedAt:           header.ClosedAt.Seconds(),
	})

	estimate, err := w.estimate(ctx, calldata, fees)
	if err != nil {
		return notary.PreparedAnchor{}, refusedBy(
			op,
			evidenceRefusals,
			fmt.Sprintf("the evidence registry refused %s batch %d", header.Stream, header.Number),
			err,
		)
	}

	return w.sign(op, nonce, estimate*gasMarginPercent/100, calldata, fees)
}
