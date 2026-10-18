package notarytest

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"math/big"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type chainBatch struct {
	header     notary.EvidenceBatchHeader
	digest     notary.Digest
	block      uint64
	anchoredAt uint64
}

type pendingBatch struct {
	anchor notary.PreparedAnchor
	header notary.EvidenceBatchHeader
}

type EvidenceChain struct {
	Scheme     notary.EvidenceScheme
	From       notary.Address
	Authorized bool
	Sealed     bool
	Unknown    map[notary.EvidenceStream]bool
	BaseFee    uint64
	Tip        uint64
	BlockTime  uint64
	BalanceWei *big.Int
	Broadcasts []notary.Digest

	block     uint64
	finalized uint64
	nonce     uint64
	batches   []chainBatch
	headers   map[notary.Digest]notary.EvidenceBatchHeader
	mempool   map[uint64]pendingBatch
	receipts  map[notary.Digest]notary.AnchorReceipt
}

var _ notary.EvidenceChain = (*EvidenceChain)(nil)

func NewEvidenceChain(scheme notary.EvidenceScheme) *EvidenceChain {
	return &EvidenceChain{
		Scheme:     scheme,
		From:       notary.Address{0xe7, 0x1d},
		Authorized: true,
		Unknown:    map[notary.EvidenceStream]bool{},
		BaseFee:    1_000,
		Tip:        100,
		BlockTime:  1760000000,
		BalanceWei: big.NewInt(1_000_000_000_000_000),
		headers:    map[notary.Digest]notary.EvidenceBatchHeader{},
		mempool:    map[uint64]pendingBatch{},
		receipts:   map[notary.Digest]notary.AnchorReceipt{},
	}
}

func (c *EvidenceChain) Anchorer() notary.Address { return c.From }

func (c *EvidenceChain) HoldsAnchorerRole(context.Context) (bool, error) { return c.Authorized, nil }

func (c *EvidenceChain) StreamHead(_ context.Context, stream notary.EvidenceStream) (notary.EvidenceStreamHead, error) {
	blockTime, err := notary.UnixTimeFromSeconds(c.BlockTime)
	if err != nil {
		return notary.EvidenceStreamHead{}, err
	}

	head := notary.EvidenceStreamHead{Known: !c.Unknown[stream], Sealed: c.Sealed, BlockTime: blockTime}
	if !head.Known {
		return head, nil
	}

	for _, b := range c.batches {
		if b.header.Stream == stream && b.block <= c.block {
			head.NextBatch = b.header.Number + 1
			head.Head = b.digest
			head.LastClosedAt = b.header.ClosedAt.Seconds()
			head.Through = b.header.Through
		}
	}

	return head, nil
}

func (c *EvidenceChain) BatchAtLatest(
	_ context.Context,
	stream notary.EvidenceStream,
	number uint64,
) (notary.AnchoredBatch, bool, error) {
	return c.batchAt(stream, number, c.block)
}

func (c *EvidenceChain) BatchAtFinalized(
	_ context.Context,
	stream notary.EvidenceStream,
	number uint64,
) (notary.AnchoredBatch, bool, error) {
	return c.batchAt(stream, number, c.finalized)
}

func (c *EvidenceChain) batchAt(
	stream notary.EvidenceStream,
	number, block uint64,
) (notary.AnchoredBatch, bool, error) {
	for _, b := range c.batches {
		if b.header.Stream == stream && b.header.Number == number && b.block <= block {
			anchoredAt, err := notary.UnixTimeFromSeconds(b.anchoredAt)
			if err != nil {
				return notary.AnchoredBatch{}, false, err
			}

			return notary.AnchoredBatch{
				Stream:     stream,
				Number:     number,
				Root:       b.header.Root,
				Digest:     b.digest,
				LeafCount:  b.header.LeafCount,
				Through:    b.header.Through,
				ClosedAt:   b.header.ClosedAt,
				AnchoredAt: anchoredAt,
			}, true, nil
		}
	}

	return notary.AnchoredBatch{}, false, nil
}

func (c *EvidenceChain) Nonce(context.Context) (uint64, error) { return c.nonce, nil }

func (c *EvidenceChain) Fees(context.Context) (uint64, uint64, error) { return c.BaseFee, c.Tip, nil }

func (c *EvidenceChain) PrepareBatch(
	ctx context.Context,
	header notary.EvidenceBatchHeader,
	nonce uint64,
	fees notary.FeeQuote,
) (notary.PreparedAnchor, error) {
	if err := c.refuse(ctx, header); err != nil {
		return notary.PreparedAnchor{}, err
	}

	return c.sign(header, nonce, 90_000, fees), nil
}

func (c *EvidenceChain) Resign(previous notary.PreparedAnchor, fees notary.FeeQuote) (notary.PreparedAnchor, error) {
	header, ok := c.headers[previous.Hash]
	if !ok {
		return notary.PreparedAnchor{}, notary.Errf(
			notary.KindInvalidArgument,
			"notarytest.EvidenceChain.Resign",
			nil,
			"unknown anchor",
		)
	}

	return c.sign(header, previous.Nonce, previous.GasLimit, fees), nil
}

func (c *EvidenceChain) sign(
	header notary.EvidenceBatchHeader,
	nonce, gas uint64,
	fees notary.FeeQuote,
) notary.PreparedAnchor {
	digest := c.Scheme.BatchDigest(header)

	raw := append([]byte{}, digest[:]...)
	for _, v := range []uint64{nonce, gas, fees.MaxFeePerGas, fees.MaxPriorityFeePerGas} {
		raw = binary.BigEndian.AppendUint64(raw, v)
	}

	anchor := notary.PreparedAnchor{
		Hash:     notary.Digest(sha256.Sum256(raw)),
		Raw:      raw,
		Nonce:    nonce,
		GasLimit: gas,
		Fees:     fees,
	}
	c.headers[anchor.Hash] = header

	return anchor
}

func (c *EvidenceChain) refuse(ctx context.Context, header notary.EvidenceBatchHeader) error {
	const op = "notarytest.EvidenceChain"

	head, _ := c.StreamHead(ctx, header.Stream)
	digest := c.Scheme.BatchDigest(header)

	switch {
	case !c.Authorized:
		return notary.Errf(notary.KindNotAnchorer, op, nil, "AccessControlUnauthorizedAccount")
	case c.Sealed:
		return notary.Errf(notary.KindAnchorRefused, op, nil, "RegistryIsSealed")
	case !head.Known:
		return notary.Errf(notary.KindAnchorRefused, op, nil, "UnknownStream")
	case header.Number < head.NextBatch:
		for _, b := range c.batches {
			if b.header.Stream == header.Stream && b.header.Number == header.Number && b.digest == digest {
				return notary.Errf(notary.KindAlreadyAnchored, op, nil, "BatchAlreadyAnchored")
			}
		}

		return notary.Errf(notary.KindAnchorRefused, op, nil, "BatchConflict")
	case header.Number != head.NextBatch:
		return notary.Errf(notary.KindAnchorRefused, op, nil, "BatchOutOfSequence")
	case header.PreviousDigest != head.Head:
		return notary.Errf(notary.KindAnchorRefused, op, nil, "PreviousDigestMismatch")
	case header.LeafCount == 0:
		return notary.Errf(notary.KindAnchorRefused, op, nil, "EmptyBatch")
	case header.ClosedAt.Seconds() > c.BlockTime:
		return notary.Errf(notary.KindAnchorDeferred, op, nil, "CloseInFuture")
	case header.ClosedAt.Seconds() <= head.LastClosedAt:
		return notary.Errf(notary.KindAnchorRefused, op, nil, "CloseNotAfterPrevious")
	case !header.Through.After(head.Through):
		return notary.Errf(notary.KindAnchorRefused, op, nil, "PositionNotAfterPrevious")
	}

	return nil
}

func (c *EvidenceChain) Broadcast(_ context.Context, anchor notary.PreparedAnchor) error {
	c.Broadcasts = append(c.Broadcasts, anchor.Hash)

	if anchor.Nonce < c.nonce {
		return notary.Errf(notary.KindAnchorDeferred, "notarytest.EvidenceChain", nil, "nonce too low")
	}

	c.mempool[anchor.Nonce] = pendingBatch{anchor: anchor, header: c.headers[anchor.Hash]}

	return nil
}

func (c *EvidenceChain) Known(_ context.Context, hash notary.Digest) (bool, error) {
	if _, ok := c.receipts[hash]; ok {
		return true, nil
	}

	for _, waiting := range c.mempool {
		if waiting.anchor.Hash == hash {
			return true, nil
		}
	}

	return false, nil
}

func (c *EvidenceChain) Receipt(_ context.Context, hash notary.Digest) (notary.AnchorReceipt, bool, error) {
	receipt, ok := c.receipts[hash]
	return receipt, ok, nil
}

func (c *EvidenceChain) Balance(context.Context) (*big.Int, error) {
	return new(big.Int).Set(c.BalanceWei), nil
}

func (c *EvidenceChain) Mine() {
	c.block++

	waiting, ok := c.mempool[c.nonce]
	if !ok {
		return
	}

	delete(c.mempool, c.nonce)
	c.nonce++

	receipt := notary.AnchorReceipt{
		Hash:              waiting.anchor.Hash,
		BlockNumber:       c.block,
		GasUsed:           88_000,
		EffectiveGasPrice: c.BaseFee + min(c.Tip, waiting.anchor.Fees.MaxPriorityFeePerGas),
		L1FeeWei:          5,
	}
	if c.refuse(context.Background(), waiting.header) == nil {
		receipt.Succeeded = true

		c.batches = append(c.batches, chainBatch{
			header:     waiting.header,
			digest:     c.Scheme.BatchDigest(waiting.header),
			block:      c.block,
			anchoredAt: c.BlockTime,
		})
	}

	c.receipts[waiting.anchor.Hash] = receipt
	cost, _ := receipt.CostWei()
	c.BalanceWei.Sub(c.BalanceWei, new(big.Int).SetUint64(cost))
}

func (c *EvidenceChain) Finalize() {
	c.finalized = c.block
}

func (c *EvidenceChain) Plant(header notary.EvidenceBatchHeader) {
	c.block++
	c.batches = append(c.batches, chainBatch{
		header:     header,
		digest:     c.Scheme.BatchDigest(header),
		block:      c.block,
		anchoredAt: c.BlockTime,
	})
}

func (c *EvidenceChain) Forget() {
	c.mempool = map[uint64]pendingBatch{}
}

func (c *EvidenceChain) Foreign() {
	delete(c.mempool, c.nonce)
	c.nonce++
	c.block++
}

func (c *EvidenceChain) Reorg() {
	for hash, receipt := range c.receipts {
		if receipt.BlockNumber == c.block {
			delete(c.receipts, hash)
			c.nonce--
		}
	}

	kept := c.batches[:0]
	for _, b := range c.batches {
		if b.block != c.block {
			kept = append(kept, b)
		}
	}

	c.batches = kept
	c.block--
}
