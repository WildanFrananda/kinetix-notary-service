package notarytest

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"math/big"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type chainEpoch struct {
	digest notary.Digest
	header notary.EpochHeader
	block  uint64
}

type pendingAnchor struct {
	anchor notary.PreparedAnchor
	header notary.EpochHeader
	proof  *proofTx
}

type Chain struct {
	Scheme        notary.Scheme
	Confidential  notary.ConfidentialScheme
	DepthValue    int
	RefuseProofs  bool
	From          notary.Address
	Authorized    bool
	Sealed        bool
	BaseFee       uint64
	Tip           uint64
	BlockTime     uint64
	BalanceWei    *big.Int
	GenesisEpoch  uint64
	GenesisDigest notary.Digest
	FeesErr       error
	BroadcastErr  error
	Broadcasts    []notary.Digest
	Approvals     [][]notary.Signature

	block     uint64
	finalized uint64
	nonce     uint64
	epochs    []chainEpoch
	headers   map[notary.Digest]notary.EpochHeader
	mempool   map[uint64]pendingAnchor
	receipts  map[notary.Digest]notary.AnchorReceipt
	mined     map[notary.Digest]uint64
	proven    map[statementKey]uint64
	proofTxs  map[notary.Digest]proofTx
}

var _ notary.AnchorChain = (*Chain)(nil)

func NewChain(scheme notary.Scheme) *Chain {
	return &Chain{
		Scheme:     scheme,
		From:       notary.Address{0x51, 0xb2},
		Authorized: true,
		BaseFee:    1_000,
		Tip:        100,
		BlockTime:  1760000000,
		BalanceWei: big.NewInt(1_000_000_000_000_000),
		headers:    map[notary.Digest]notary.EpochHeader{},
		mempool:    map[uint64]pendingAnchor{},
		receipts:   map[notary.Digest]notary.AnchorReceipt{},
		mined:      map[notary.Digest]uint64{},
		proven:     map[statementKey]uint64{},
		proofTxs:   map[notary.Digest]proofTx{},
	}
}

func (c *Chain) digest(header notary.EpochHeader) notary.Digest {
	if c.Confidential != nil {
		return c.Confidential.EpochDigest(header)
	}

	return c.Scheme.EpochDigest(header)
}

func (c *Chain) Anchorer() notary.Address { return c.From }

func (c *Chain) Head(context.Context) (notary.LedgerHead, error) {
	head := notary.LedgerHead{
		NextEpoch:     c.GenesisEpoch,
		Head:          c.GenesisDigest,
		GenesisEpoch:  c.GenesisEpoch,
		GenesisDigest: c.GenesisDigest,
		Sealed:        c.Sealed,
	}
	if n := len(c.epochs); n > 0 {
		last := c.epochs[n-1]
		head.NextEpoch, head.Head, head.LastAsOf = last.header.Number+1, last.digest, last.header.AsOf.Seconds()
	}

	return head, nil
}

func (c *Chain) EpochAtLatest(_ context.Context, number uint64) (notary.AnchoredEpoch, bool, error) {
	return c.epochAt(number, c.block)
}

func (c *Chain) EpochAtFinalized(_ context.Context, number uint64) (notary.AnchoredEpoch, bool, error) {
	return c.epochAt(number, c.finalized)
}

func (c *Chain) epochAt(number, block uint64) (notary.AnchoredEpoch, bool, error) {
	for _, e := range c.epochs {
		if e.header.Number == number && e.block <= block {
			return notary.AnchoredEpoch{
				Number:    number,
				Root:      e.header.Root,
				Digest:    e.digest,
				Total:     e.header.Total,
				LeafCount: e.header.LeafCount,
				AsOf:      e.header.AsOf,
			}, true, nil
		}
	}

	return notary.AnchoredEpoch{}, false, nil
}

func (c *Chain) Nonce(context.Context) (uint64, error) { return c.nonce, nil }

func (c *Chain) Fees(context.Context) (uint64, uint64, error) {
	return c.BaseFee, c.Tip, c.FeesErr
}

func (c *Chain) Prepare(
	ctx context.Context,
	header notary.EpochHeader,
	approvals []notary.Signature,
	nonce uint64,
	fees notary.FeeQuote,
) (notary.PreparedAnchor, error) {
	c.Approvals = append(c.Approvals, approvals)
	if len(approvals) < 2 {
		return notary.PreparedAnchor{}, notary.Errf(
			notary.KindAnchorRefused,
			"notarytest.Chain",
			nil,
			"NotEnoughApprovals",
		)
	}

	if err := c.refuse(ctx, header); err != nil {
		return notary.PreparedAnchor{}, err
	}

	return c.sign(header, nonce, 150_000, fees), nil
}

func (c *Chain) Resign(previous notary.PreparedAnchor, fees notary.FeeQuote) (notary.PreparedAnchor, error) {
	if proof, ok := c.proofTxs[previous.Hash]; ok {
		return c.signProof(proof, previous.Nonce, previous.GasLimit, fees), nil
	}

	header, ok := c.headers[previous.Hash]
	if !ok {
		return notary.PreparedAnchor{}, notary.Errf(
			notary.KindInvalidArgument,
			"notarytest.Chain.Resign",
			nil,
			"unknown anchor",
		)
	}

	return c.sign(header, previous.Nonce, previous.GasLimit, fees), nil
}

func (c *Chain) sign(header notary.EpochHeader, nonce, gas uint64, fees notary.FeeQuote) notary.PreparedAnchor {
	digest := c.digest(header)

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

func (c *Chain) refuse(ctx context.Context, header notary.EpochHeader) error {
	const op = "notarytest.Chain"

	head, _ := c.Head(ctx)

	digest := c.digest(header)
	switch {
	case !c.Authorized:
		return notary.Errf(notary.KindNotAnchorer, op, nil, "AccessControlUnauthorizedAccount")
	case c.Sealed:
		return notary.Errf(notary.KindAnchorRefused, op, nil, "LedgerIsSealed")
	case header.Number < head.NextEpoch:
		for _, e := range c.epochs {
			if e.header.Number == header.Number && e.digest == digest {
				return notary.Errf(notary.KindAlreadyAnchored, op, nil, "EpochAlreadyAnchored")
			}
		}

		return notary.Errf(notary.KindAnchorRefused, op, nil, "EpochConflict")
	case header.Number != head.NextEpoch:
		return notary.Errf(notary.KindAnchorRefused, op, nil, "EpochOutOfSequence")
	case header.PreviousDigest != head.Head:
		return notary.Errf(notary.KindAnchorRefused, op, nil, "PreviousDigestMismatch")
	case header.AsOf.Seconds() > c.BlockTime:
		return notary.Errf(notary.KindAnchorDeferred, op, nil, "SnapshotInFuture")
	case header.AsOf.Seconds() <= head.LastAsOf:
		return notary.Errf(notary.KindAnchorRefused, op, nil, "SnapshotNotAfterPrevious")
	case c.Confidential != nil && header.Total != 0:
		return notary.Errf(notary.KindAnchorRefused, op, nil, "TotalNotHidden")
	case c.Confidential != nil && !c.provenAt(keyOf(header), c.block):
		return notary.Errf(notary.KindAnchorDeferred, op, nil, "EpochNotProven")
	}

	return nil
}

func (c *Chain) Broadcast(_ context.Context, anchor notary.PreparedAnchor) error {
	c.Broadcasts = append(c.Broadcasts, anchor.Hash)
	if c.BroadcastErr != nil {
		return c.BroadcastErr
	}

	if anchor.Nonce < c.nonce {
		return notary.Errf(notary.KindAnchorDeferred, "notarytest.Chain", nil, "nonce too low")
	}

	if waiting, ok := c.mempool[anchor.Nonce]; ok && waiting.anchor.Hash != anchor.Hash {
		if anchor.Fees.MaxFeePerGas*10 < waiting.anchor.Fees.MaxFeePerGas*11 {
			return notary.Errf(
				notary.KindAnchorDeferred,
				"notarytest.Chain",
				nil,
				"replacement transaction underpriced",
			)
		}
	}

	waiting := pendingAnchor{anchor: anchor, header: c.headers[anchor.Hash]}
	if proof, ok := c.proofTxs[anchor.Hash]; ok {
		waiting.proof = &proof
	}

	c.mempool[anchor.Nonce] = waiting

	return nil
}

func (c *Chain) Known(_ context.Context, hash notary.Digest) (bool, error) {
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

func (c *Chain) Receipt(_ context.Context, hash notary.Digest) (notary.AnchorReceipt, bool, error) {
	receipt, ok := c.receipts[hash]
	return receipt, ok, nil
}

func (c *Chain) Balance(context.Context) (*big.Int, error) {
	return new(big.Int).Set(c.BalanceWei), nil
}

func (c *Chain) Mine() {
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
		GasUsed:           144_000,
		EffectiveGasPrice: c.BaseFee + min(c.Tip, waiting.anchor.Fees.MaxPriorityFeePerGas),
		L1FeeWei:          7,
	}
	if waiting.proof != nil {
		receipt.GasUsed = 390_000
		if !c.RefuseProofs && !c.provenAt(waiting.proof.key, c.block) {
			receipt.Succeeded = true
			c.proven[waiting.proof.key] = c.block
		}
	} else if c.refuse(context.Background(), waiting.header) == nil {
		receipt.Succeeded = true

		c.epochs = append(
			c.epochs,
			chainEpoch{digest: c.digest(waiting.header), header: waiting.header, block: c.block},
		)
	}

	c.receipts[waiting.anchor.Hash] = receipt
	c.mined[waiting.anchor.Hash] = c.block
	cost, _ := receipt.CostWei()
	c.BalanceWei.Sub(c.BalanceWei, new(big.Int).SetUint64(cost))
}

func (c *Chain) Finalize() {
	c.finalized = c.block
}

func (c *Chain) Reorg() {
	for hash, block := range c.mined {
		if block == c.block {
			delete(c.receipts, hash)
			delete(c.mined, hash)
			c.nonce--
		}
	}

	kept := c.epochs[:0]
	for _, e := range c.epochs {
		if e.block != c.block {
			kept = append(kept, e)
		}
	}

	for key, block := range c.proven {
		if block == c.block {
			delete(c.proven, key)
		}
	}

	c.epochs = kept
	c.block--
}

func (c *Chain) Forget() {
	c.mempool = map[uint64]pendingAnchor{}
}

func (c *Chain) Foreign() {
	delete(c.mempool, c.nonce)
	c.nonce++
	c.block++
}

func (c *Chain) Plant(header notary.EpochHeader) {
	c.block++
	c.epochs = append(c.epochs, chainEpoch{digest: c.digest(header), header: header, block: c.block})
}
