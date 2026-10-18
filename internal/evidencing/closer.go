package evidencing

import (
	"context"
	"log/slog"
	"math/big"
	"time"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type Closer struct {
	scheme   notary.EvidenceScheme
	store    notary.EvidenceStore
	chain    notary.EvidenceChain
	observer Observer
	clock    func() time.Time
	log      *slog.Logger
	policy   Policy
}

func NewCloser(
	scheme notary.EvidenceScheme,
	store notary.EvidenceStore,
	chain notary.EvidenceChain,
	observer Observer,
	clock func() time.Time,
	log *slog.Logger,
	policy Policy,
) *Closer {
	return &Closer{
		scheme:   scheme,
		store:    store,
		chain:    chain,
		observer: observer,
		clock:    clock,
		log:      log,
		policy:   policy,
	}
}

func (c *Closer) Close(ctx context.Context, lease notary.Lease) error {
	record, found, err := c.store.InFlightBatch(ctx)
	if err != nil {
		return err
	}

	if found {
		err := c.advance(ctx, lease, record)
		if err != nil && !fatal(err) {
			logFailure(ctx, c.log, "evidence anchoring", record.Header.Stream, err)
			return nil
		}

		return err
	}

	for _, stream := range notary.EvidenceStreams() {
		started, err := c.begin(ctx, lease, stream)
		if err != nil && fatal(err) {
			return err
		}

		if err != nil {
			logFailure(ctx, c.log, "evidence anchoring", stream, err)
		}

		if started {
			return nil
		}
	}

	return nil
}

func (c *Closer) begin(ctx context.Context, lease notary.Lease, stream notary.EvidenceStream) (bool, error) {
	const op = "evidencing.Closer.begin"

	last, hasLast, err := c.store.LastBatch(ctx, stream)
	if err != nil {
		return false, err
	}

	if hasLast && last.State == notary.StateAbandoned {
		return false, notary.Errf(
			notary.KindEpochAbandoned,
			op,
			nil,
			"%s batch %d was abandoned (%s); a person must decide what happens next",
			stream,
			last.Header.Number,
			last.AbandonedReason,
		)
	}

	if hasLast && last.Registry != c.policy.Registry {
		return false, notary.Errf(
			notary.KindLedgerDiverged,
			op,
			nil,
			"this database anchors evidence to %s on chain %d, not %s on chain %d",
			last.Registry.Address.Hex(),
			last.Registry.ChainID,
			c.policy.Registry.Address.Hex(),
			c.policy.Registry.ChainID,
		)
	}

	pending, err := c.store.UnbatchedCount(ctx, stream)
	if err != nil {
		return false, err
	}

	c.observer.EvidencePending(stream, pending)

	if pending == 0 {
		return false, nil
	}

	full := pending >= int64(c.policy.MaxBatch)
	if hasLast && !full && c.clock().Before(last.Header.ClosedAt.Time().Add(c.policy.BatchInterval)) {
		return false, nil
	}

	head, err := c.chain.StreamHead(ctx, stream)
	if err != nil {
		return false, err
	}

	if head.Sealed {
		return false, notary.Errf(
			notary.KindAnchorRefused,
			op,
			nil,
			"the evidence registry is sealed; point the notary at its successor",
		)
	}

	if !head.Known {
		return false, notary.Errf(notary.KindWrongLedger, op, nil, "the evidence registry keeps no %s stream", stream)
	}

	if hasLast && (head.NextBatch != last.Header.Number+1 || head.Head != last.Digest) {
		return false, notary.Errf(
			notary.KindLedgerDiverged,
			op,
			nil,
			"the registry's %s stream is at batch %d with head %s, this database finalized batch %d with digest %s",
			stream,
			head.NextBatch,
			head.Head.Hex(),
			last.Header.Number,
			last.Digest.Hex(),
		)
	}

	if !hasLast && !head.AtGenesis() {
		return false, notary.Errf(
			notary.KindLedgerDiverged,
			op,
			nil,
			"the registry already holds %s batches up to %d, and this database built none of them",
			stream,
			head.NextBatch,
		)
	}

	if head.BlockTime.Seconds() <= head.LastClosedAt {
		return false, notary.Errf(
			notary.KindAnchorDeferred,
			op,
			nil,
			"the chain's time %d is not after the %s stream's last close at %d",
			head.BlockTime.Seconds(),
			stream,
			head.LastClosedAt,
		)
	}

	entries, err := c.store.Unbatched(ctx, stream, c.policy.MaxBatch)
	if err != nil {
		return false, err
	}

	if len(entries) == 0 {
		return false, nil
	}

	if first := entries[0].Record.Position(); !first.After(head.Through) {
		return false, notary.Errf(
			notary.KindLedgerDiverged,
			op,
			nil,
			"the first unbatched %s record at %d/%d does not follow the registry's %d/%d",
			stream,
			first.Transaction,
			first.Sequence,
			head.Through.Transaction,
			head.Through.Sequence,
		)
	}

	plan, err := notary.PlanEvidenceBatch(c.scheme, stream, head.NextBatch, head.Head, head.BlockTime, entries)
	if err != nil {
		return false, err
	}

	record, err := c.store.SaveBatch(ctx, lease, c.policy.Registry, plan)
	if err != nil {
		return false, err
	}

	c.log.InfoContext(
		ctx,
		"evidence batch built",
		"stream",
		stream.String(),
		"batch",
		record.Header.Number,
		"records",
		record.Header.LeafCount,
		"digest",
		record.Digest.Hex(),
		"closed_at",
		record.Header.ClosedAt.Time(),
	)

	return true, c.advance(ctx, lease, record)
}

func (c *Closer) advance(ctx context.Context, lease notary.Lease, record notary.EvidenceBatchRecord) error {
	if record.Registry != c.policy.Registry {
		return notary.Errf(
			notary.KindLedgerDiverged,
			"evidencing.Closer.advance",
			nil,
			"%s batch %d was built for %s on chain %d",
			record.Header.Stream,
			record.Header.Number,
			record.Registry.Address.Hex(),
			record.Registry.ChainID,
		)
	}

	if record.State == notary.StateBuilt {
		return c.submit(ctx, lease, record)
	}

	return c.watch(ctx, lease, record)
}

func (c *Closer) submit(ctx context.Context, lease notary.Lease, record notary.EvidenceBatchRecord) error {
	stream, number := record.Header.Stream, record.Header.Number

	latest, found, err := c.chain.BatchAtLatest(ctx, stream, number)
	if err != nil {
		return err
	}

	if found {
		if latest.Digest != record.Digest {
			return c.abandon(ctx, lease, record, conflict(record, latest.Digest))
		}

		if record.State == notary.StateBuilt {
			return c.store.TransitionBatch(ctx, lease, stream, number, notary.StateBuilt, notary.StateSubmitted, "")
		}

		return nil
	}

	nonce, err := c.chain.Nonce(ctx)
	if err != nil {
		return err
	}

	base, tip, err := c.chain.Fees(ctx)
	if err != nil {
		return err
	}

	quote, err := notary.QuoteFees(base, tip, c.policy.FeeCeiling)
	if err != nil {
		if notary.KindOf(err) == notary.KindAnchorDeferred {
			c.observer.EvidenceFeeCeilingHit()
		}

		return err
	}

	anchor, err := c.chain.PrepareBatch(ctx, record.Header, nonce, quote)
	switch notary.KindOf(err) {
	case notary.KindUnknown:
	case notary.KindAlreadyAnchored:
		if record.State == notary.StateBuilt {
			return c.store.TransitionBatch(ctx, lease, stream, number, notary.StateBuilt, notary.StateSubmitted, "")
		}

		return nil
	case notary.KindAnchorRefused:
		return c.abandon(ctx, lease, record, err)
	default:
		return err
	}

	attempt := notary.BatchAttempt{
		Stream:  stream,
		Batch:   number,
		Anchor:  anchor,
		Outcome: notary.AttemptPending,
		SentAt:  c.clock(),
	}
	if err := c.store.RecordBatchAttempt(ctx, lease, attempt); err != nil {
		return err
	}

	c.log.InfoContext(
		ctx,
		"evidence anchor sent",
		"stream",
		stream.String(),
		"batch",
		number,
		"tx",
		anchor.Hash.Hex(),
		"nonce",
		anchor.Nonce,
		"max_fee_per_gas",
		anchor.Fees.MaxFeePerGas,
		"max_priority_fee_per_gas",
		anchor.Fees.MaxPriorityFeePerGas,
	)

	return c.chain.Broadcast(ctx, anchor)
}

func (c *Closer) watch(ctx context.Context, lease notary.Lease, record notary.EvidenceBatchRecord) error {
	stream, number := record.Header.Stream, record.Header.Number

	final, found, err := c.chain.BatchAtFinalized(ctx, stream, number)
	if err != nil {
		return err
	}

	if found {
		if final.Digest != record.Digest {
			return c.abandon(ctx, lease, record, conflict(record, final.Digest))
		}

		return c.finalize(ctx, lease, record)
	}

	latest, found, err := c.chain.BatchAtLatest(ctx, stream, number)
	if err != nil {
		return err
	}

	if found {
		if latest.Digest != record.Digest {
			return c.abandon(ctx, lease, record, conflict(record, latest.Digest))
		}

		if record.State == notary.StateSubmitted {
			return c.store.TransitionBatch(ctx, lease, stream, number, notary.StateSubmitted, notary.StateMined, "")
		}

		return nil
	}

	if record.State == notary.StateMined {
		c.log.WarnContext(
			ctx,
			"anchored evidence batch left the chain; resubmitting",
			"stream",
			stream.String(),
			"batch",
			number,
		)

		if err := c.store.TransitionBatch(
			ctx,
			lease,
			stream,
			number,
			notary.StateMined,
			notary.StateSubmitted,
			"",
		); err != nil {
			return err
		}

		record.State = notary.StateSubmitted
	}

	return c.keepAlive(ctx, lease, record)
}

func (c *Closer) keepAlive(ctx context.Context, lease notary.Lease, record notary.EvidenceBatchRecord) error {
	const op = "evidencing.Closer.keepAlive"

	stream, number := record.Header.Stream, record.Header.Number

	attempts, err := c.store.BatchAttempts(ctx, stream, number)
	if err != nil {
		return err
	}

	if len(attempts) == 0 {
		return c.submit(ctx, lease, record)
	}

	reverted := false

	for _, attempt := range attempts {
		if attempt.Outcome == notary.AttemptReverted {
			reverted = true
			continue
		}

		if attempt.Outcome != notary.AttemptPending {
			continue
		}

		receipt, mined, err := c.chain.Receipt(ctx, attempt.Anchor.Hash)
		if err != nil {
			return err
		}

		if mined && !receipt.Succeeded {
			if err := c.store.SettleBatchAttempt(
				ctx,
				lease,
				attempt.Anchor.Hash,
				notary.AttemptReverted,
				&receipt,
			); err != nil {
				return err
			}

			reverted = true
		}
	}

	latest := attempts[len(attempts)-1]

	nonce, err := c.chain.Nonce(ctx)
	if err != nil {
		return err
	}

	if nonce > latest.Anchor.Nonce {
		if reverted {
			c.log.WarnContext(
				ctx,
				"evidence anchor reverted; preparing a fresh one",
				"stream",
				stream.String(),
				"batch",
				number,
				"nonce",
				nonce,
			)

			return c.submit(ctx, lease, record)
		}

		if c.anyMined(ctx, attempts) {
			return nil
		}

		return notary.Errf(
			notary.KindLedgerDiverged,
			op,
			nil,
			"nonce %d was used by a transaction this notary never sent, and %s batch %d is not on chain",
			latest.Anchor.Nonce,
			stream,
			number,
		)
	}

	if c.clock().Sub(latest.SentAt) >= c.policy.ReplaceAfter {
		return c.replace(ctx, lease, latest)
	}

	known, err := c.chain.Known(ctx, latest.Anchor.Hash)
	if err != nil {
		return err
	}

	if !known {
		c.log.InfoContext(
			ctx,
			"evidence anchor dropped by the node; broadcasting again",
			"tx",
			latest.Anchor.Hash.Hex(),
		)

		return c.chain.Broadcast(ctx, latest.Anchor)
	}

	return nil
}

func (c *Closer) anyMined(ctx context.Context, attempts []notary.BatchAttempt) bool {
	for _, attempt := range attempts {
		if _, mined, err := c.chain.Receipt(ctx, attempt.Anchor.Hash); err == nil && mined {
			return true
		}
	}

	return false
}

func (c *Closer) replace(ctx context.Context, lease notary.Lease, latest notary.BatchAttempt) error {
	bumped, raised := latest.Anchor.Fees.Bumped(c.policy.FeeCeiling)
	if !raised {
		c.observer.EvidenceFeeCeilingHit()
		c.log.WarnContext(
			ctx,
			"evidence anchor is stuck at the fee ceiling; broadcasting it unchanged",
			"tx",
			latest.Anchor.Hash.Hex(),
			"max_fee_per_gas",
			latest.Anchor.Fees.MaxFeePerGas,
		)

		return c.chain.Broadcast(ctx, latest.Anchor)
	}

	anchor, err := c.chain.Resign(latest.Anchor, bumped)
	if err != nil {
		return err
	}

	attempt := notary.BatchAttempt{
		Stream:  latest.Stream,
		Batch:   latest.Batch,
		Anchor:  anchor,
		Outcome: notary.AttemptPending,
		SentAt:  c.clock(),
	}
	if err := c.store.RecordBatchAttempt(ctx, lease, attempt); err != nil {
		return err
	}

	c.log.InfoContext(
		ctx,
		"evidence anchor replaced",
		"stream",
		latest.Stream.String(),
		"batch",
		latest.Batch,
		"nonce",
		anchor.Nonce,
		"replaces",
		latest.Anchor.Hash.Hex(),
		"tx",
		anchor.Hash.Hex(),
		"max_fee_per_gas",
		anchor.Fees.MaxFeePerGas,
	)

	return c.chain.Broadcast(ctx, anchor)
}

func (c *Closer) finalize(ctx context.Context, lease notary.Lease, record notary.EvidenceBatchRecord) error {
	stream, number := record.Header.Stream, record.Header.Number

	attempts, err := c.store.BatchAttempts(ctx, stream, number)
	if err != nil {
		return err
	}

	for _, attempt := range attempts {
		if attempt.Outcome != notary.AttemptPending && attempt.Outcome != notary.AttemptIncluded {
			continue
		}

		receipt, mined, err := c.chain.Receipt(ctx, attempt.Anchor.Hash)
		if err != nil {
			return err
		}

		outcome := notary.AttemptSuperseded

		var settled *notary.AnchorReceipt

		if mined {
			outcome = notary.AttemptReverted
			if receipt.Succeeded {
				outcome = notary.AttemptIncluded
			}

			settled = &receipt
		}

		if err := c.store.SettleBatchAttempt(ctx, lease, attempt.Anchor.Hash, outcome, settled); err != nil {
			return err
		}
	}

	if err := c.store.TransitionBatch(ctx, lease, stream, number, record.State, notary.StateFinalized, ""); err != nil {
		return err
	}

	c.log.InfoContext(
		ctx,
		"evidence batch final",
		"stream",
		stream.String(),
		"batch",
		number,
		"digest",
		record.Digest.Hex(),
	)

	return nil
}

func (c *Closer) abandon(
	ctx context.Context,
	lease notary.Lease,
	record notary.EvidenceBatchRecord,
	cause error,
) error {
	stream, number := record.Header.Stream, record.Header.Number

	if err := c.store.TransitionBatch(
		ctx,
		lease,
		stream,
		number,
		record.State,
		notary.StateAbandoned,
		cause.Error(),
	); err != nil {
		return err
	}

	c.log.ErrorContext(
		ctx,
		"evidence batch abandoned; a person must decide what happens next",
		"stream",
		stream.String(),
		"batch",
		number,
		"reason",
		cause.Error(),
	)

	return cause
}

func (c *Closer) Report(ctx context.Context) error {
	for _, stream := range notary.EvidenceStreams() {
		last, found, err := c.store.LastFinalizedBatch(ctx, stream)
		if err != nil {
			return err
		}

		if found {
			c.observer.EvidenceFinalized(last)
		}

		pending, err := c.store.UnbatchedCount(ctx, stream)
		if err != nil {
			return err
		}

		c.observer.EvidencePending(stream, pending)
	}

	abandoned, err := c.store.AbandonedBatches(ctx)
	if err != nil {
		return err
	}

	c.observer.EvidenceAbandoned(abandoned)

	balance, err := c.chain.Balance(ctx)
	if err != nil {
		return err
	}

	costs, err := c.store.RecentBatchCosts(ctx, c.policy.CostWindow)
	if err != nil {
		return err
	}

	mean := meanOf(costs)
	streams := time.Duration(len(notary.EvidenceStreams()))
	runway, known := notary.Runway(balance, mean, c.policy.BatchInterval/streams)
	c.observer.EvidenceFunds(balance, mean, runway, known)

	return nil
}

func meanOf(costs []uint64) uint64 {
	if len(costs) == 0 {
		return 0
	}

	sum := new(big.Int)
	for _, cost := range costs {
		sum.Add(sum, new(big.Int).SetUint64(cost))
	}

	return sum.Quo(sum, big.NewInt(int64(len(costs)))).Uint64()
}

func conflict(record notary.EvidenceBatchRecord, onChain notary.Digest) error {
	return notary.Errf(
		notary.KindAnchorRefused,
		"evidencing.Closer.watch",
		nil,
		"the registry holds %s batch %d as %s, this notary built %s",
		record.Header.Stream,
		record.Header.Number,
		onChain.Hex(),
		record.Digest.Hex(),
	)
}
