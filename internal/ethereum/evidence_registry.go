package ethereum

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/WildanFrananda/kinetix-notary-service/internal/contractgen/evidence"
	"github.com/WildanFrananda/kinetix-notary-service/internal/contractgen/evidenceinterface"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type EvidenceRegistry struct {
	client      *ethclient.Client
	binding     *evidence.EvidenceRegistry
	contract    *bind.BoundContract
	expected    notary.RegistryIdentity
	interfaceID [4]byte
	scheme      uint16
	timeout     time.Duration
}

var _ notary.EvidenceRegistry = (*EvidenceRegistry)(nil)

func NewEvidenceRegistry(
	client *ethclient.Client,
	registry notary.RegistryIdentity,
	scheme notary.EvidenceScheme,
	timeout time.Duration,
) *EvidenceRegistry {
	binding := evidence.NewEvidenceRegistry()

	return &EvidenceRegistry{
		client:      client,
		binding:     binding,
		contract:    binding.Instance(client, common.Address(registry.Address)),
		expected:    registry,
		interfaceID: evidenceInterfaceID(),
		scheme:      scheme.EvidenceVersion(),
		timeout:     timeout,
	}
}

func (r *EvidenceRegistry) Identity(ctx context.Context) (notary.RegistryIdentity, error) {
	const op = "ethereum.EvidenceRegistry.Identity"

	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	chainID, err := r.client.ChainID(ctx)
	if err != nil {
		return notary.RegistryIdentity{}, unavailable(op, err)
	}

	if !chainID.IsUint64() || chainID.Uint64() != r.expected.ChainID {
		return notary.RegistryIdentity{}, notary.Errf(
			notary.KindWrongLedger,
			op,
			nil,
			"the RPC endpoint serves chain %s, the evidence registry is configured on chain %d",
			chainID,
			r.expected.ChainID,
		)
	}

	code, err := r.client.CodeAt(ctx, common.Address(r.expected.Address), nil)
	if err != nil {
		return notary.RegistryIdentity{}, unavailable(op, err)
	}

	if len(code) == 0 {
		return notary.RegistryIdentity{}, notary.Errf(
			notary.KindWrongLedger,
			op,
			nil,
			"chain %d has no contract at %s",
			r.expected.ChainID,
			r.expected.Address.Hex(),
		)
	}

	supported, err := probe(
		ctx,
		r.client,
		r.expected.Address,
		r.binding.PackSupportsInterface(r.interfaceID),
		r.binding.UnpackSupportsInterface,
	)
	if err != nil {
		return notary.RegistryIdentity{}, r.wrap(op, err)
	}

	if !supported {
		return notary.RegistryIdentity{}, notary.Errf(
			notary.KindWrongLedger,
			op,
			nil,
			"%s does not announce the evidence registry interface",
			r.expected.Address.Hex(),
		)
	}

	scheme, err := probe(ctx, r.client, r.expected.Address, r.binding.PackScheme(), r.binding.UnpackScheme)
	if err != nil {
		return notary.RegistryIdentity{}, r.wrap(op, err)
	}

	if scheme != r.scheme {
		return notary.RegistryIdentity{}, notary.Errf(
			notary.KindWrongLedger,
			op,
			nil,
			"%s commits evidence with scheme %d, this notary speaks scheme %d",
			r.expected.Address.Hex(),
			scheme,
			r.scheme,
		)
	}

	streams, err := bind.Call(
		r.contract,
		&bind.CallOpts{Context: ctx},
		r.binding.PackStreams(),
		r.binding.UnpackStreams,
	)
	if err != nil {
		return notary.RegistryIdentity{}, unavailable(op, err)
	}

	kept := make(map[notary.EvidenceStream]bool, len(streams))
	for _, stream := range streams {
		kept[notary.EvidenceStream(stream)] = true
	}

	for _, stream := range notary.EvidenceStreams() {
		if !kept[stream] {
			return notary.RegistryIdentity{}, notary.Errf(
				notary.KindWrongLedger,
				op,
				nil,
				"%s keeps no %s stream",
				r.expected.Address.Hex(),
				stream,
			)
		}
	}

	return r.expected, nil
}

func (r *EvidenceRegistry) Batch(
	ctx context.Context,
	stream notary.EvidenceStream,
	number uint64,
) (notary.AnchoredBatch, error) {
	const op = "ethereum.EvidenceRegistry.Batch"

	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	final, err := finalizedBlock(ctx, r.client)
	if err != nil {
		return notary.AnchoredBatch{}, unavailable(op, err)
	}

	batch, found, err := batchAt(ctx, op, r.contract, r.binding, stream, number, final)
	if err != nil {
		return notary.AnchoredBatch{}, err
	}

	if found {
		return batch, nil
	}

	_, found, err = batchAt(ctx, op, r.contract, r.binding, stream, number, nil)
	if err != nil {
		return notary.AnchoredBatch{}, err
	}

	if found {
		return notary.AnchoredBatch{}, notary.Errf(
			notary.KindEpochNotFinal,
			op,
			nil,
			"%s batch %d is anchored but not yet final at block %s",
			stream,
			number,
			final,
		)
	}

	return notary.AnchoredBatch{}, notary.Errf(
		notary.KindEpochNotAnchored,
		op,
		nil,
		"%s has no %s batch %d",
		r.expected.Address.Hex(),
		stream,
		number,
	)
}

func (r *EvidenceRegistry) ConfirmsEvidence(
	ctx context.Context,
	batch notary.AnchoredBatch,
	salt notary.Salt,
	record notary.Digest,
	proof notary.EvidenceProof,
) (bool, error) {
	const op = "ethereum.EvidenceRegistry.ConfirmsEvidence"

	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	final, err := finalizedBlock(ctx, r.client)
	if err != nil {
		return false, unavailable(op, err)
	}

	siblings := make([][32]byte, len(proof.Siblings))
	for i, sibling := range proof.Siblings {
		siblings[i] = sibling
	}

	confirmed, err := bind.Call(
		r.contract,
		&bind.CallOpts{Context: ctx, BlockNumber: final},
		r.binding.PackVerifyInclusion(uint8(batch.Stream), batch.Number, salt, record, proof.Index, siblings),
		r.binding.UnpackVerifyInclusion,
	)
	if err != nil {
		if unknownBatch(err) || errors.Is(err, bind.ErrNoCode) {
			return false, notary.Errf(
				notary.KindEpochNotFinal,
				op,
				nil,
				"%s batch %d is not final at block %s",
				batch.Stream,
				batch.Number,
				final,
			)
		}

		return false, unavailable(op, err)
	}

	return confirmed, nil
}

func (r *EvidenceRegistry) wrap(op string, err error) error {
	if errors.Is(err, errNotALedger) {
		return notary.Errf(
			notary.KindWrongLedger,
			op,
			nil,
			"%s does not answer as an evidence registry",
			r.expected.Address.Hex(),
		)
	}

	return unavailable(op, err)
}

func batchAt(
	ctx context.Context,
	op string,
	contract *bind.BoundContract,
	binding *evidence.EvidenceRegistry,
	stream notary.EvidenceStream,
	number uint64,
	block *big.Int,
) (notary.AnchoredBatch, bool, error) {
	batch, err := bind.Call(
		contract,
		&bind.CallOpts{Context: ctx, BlockNumber: block},
		binding.PackBatchAt(uint8(stream), number),
		binding.UnpackBatchAt,
	)
	if err != nil {
		if unknownBatch(err) || (block != nil && errors.Is(err, bind.ErrNoCode)) {
			return notary.AnchoredBatch{}, false, nil
		}

		return notary.AnchoredBatch{}, false, unavailable(op, err)
	}

	closedAt, err := notary.UnixTimeFromSeconds(batch.ClosedAt)
	if err != nil {
		return notary.AnchoredBatch{}, false, notary.Errf(
			notary.KindWrongLedger,
			op,
			err,
			"%s batch %d",
			stream,
			number,
		)
	}

	anchoredAt, err := notary.UnixTimeFromSeconds(batch.AnchoredAt)
	if err != nil {
		return notary.AnchoredBatch{}, false, notary.Errf(
			notary.KindWrongLedger,
			op,
			err,
			"%s batch %d",
			stream,
			number,
		)
	}

	return notary.AnchoredBatch{
		Stream:     stream,
		Number:     number,
		Root:       batch.Root,
		Digest:     batch.Digest,
		LeafCount:  batch.LeafCount,
		Through:    notary.LogPosition{Transaction: batch.ThroughTransaction, Sequence: batch.ThroughSequence},
		ClosedAt:   closedAt,
		AnchoredAt: anchoredAt,
	}, true, nil
}

func unknownBatch(err error) bool {
	return selectorIs(err, evidence.EvidenceRegistryUnknownBatchErrorID()) ||
		selectorIs(err, evidence.EvidenceRegistryUnknownStreamErrorID())
}

func finalizedBlock(ctx context.Context, client *ethclient.Client) (*big.Int, error) {
	header, err := client.HeaderByNumber(ctx, big.NewInt(int64(rpc.FinalizedBlockNumber)))
	if err != nil {
		return nil, fmt.Errorf("finalized block: %w", err)
	}

	return header.Number, nil
}

func evidenceInterfaceID() [4]byte {
	var id [4]byte
	for _, method := range evidenceinterface.NewIEvidenceRegistry().GetABI().Methods {
		for i := range id {
			id[i] ^= method.ID[i]
		}
	}

	return id
}
