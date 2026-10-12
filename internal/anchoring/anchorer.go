package anchoring

import (
	"context"
	"log/slog"
	"math/big"
	"time"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type Anchorer struct {
	scheme   notary.Scheme
	salts    notary.SaltSource
	source   notary.SnapshotSource
	store    notary.EpochStore
	chain    notary.AnchorChain
	observer Observer
	clock    func() time.Time
	log      *slog.Logger
	policy   Policy
}

func NewAnchorer(
	scheme notary.Scheme,
	salts notary.SaltSource,
	source notary.SnapshotSource,
	store notary.EpochStore,
	chain notary.AnchorChain,
	observer Observer,
	clock func() time.Time,
	log *slog.Logger,
	policy Policy,
) *Anchorer {
	return &Anchorer{
		scheme:   scheme,
		salts:    salts,
		source:   source,
		store:    store,
		chain:    chain,
		observer: observer,
		clock:    clock,
		log:      log,
		policy:   policy,
	}
}

func (a *Anchorer) Tick(ctx context.Context) error {
	lease, held, err := a.store.Lease(ctx, a.policy.Instance, a.policy.LeaseTerm)
	if err != nil {
		return err
	}

	a.observer.LeaseHeld(held)

	if !held {
		return nil
	}

	work := a.step(ctx, lease)

	reported := a.report(ctx)
	if work == nil {
		return reported
	}

	if reported != nil {
		a.log.WarnContext(ctx, "reporting failed after a failed step", "err", reported.Error())
	}

	return work
}

func (a *Anchorer) step(ctx context.Context, lease notary.Lease) error {
	record, found, err := a.store.InFlight(ctx)
	if err != nil {
		return err
	}

	if found {
		return a.advance(ctx, lease, record)
	}

	return a.begin(ctx, lease)
}

func (a *Anchorer) begin(ctx context.Context, lease notary.Lease) error {
	const op = "anchoring.Anchorer.begin"

	last, hasLast, err := a.store.LastFinalized(ctx)
	if err != nil {
		return err
	}

	if hasLast && last.Ledger != a.policy.Ledger {
		return notary.Errf(
			notary.KindLedgerDiverged,
			op,
			nil,
			"this database anchors to %s on chain %d, not %s on chain %d",
			last.Ledger.Address.Hex(),
			last.Ledger.ChainID,
			a.policy.Ledger.Address.Hex(),
			a.policy.Ledger.ChainID,
		)
	}

	if hasLast && a.clock().Before(last.Header.AsOf.Time().Add(a.policy.EpochInterval)) {
		return nil
	}

	head, err := a.chain.Head(ctx)
	if err != nil {
		return err
	}

	if head.Sealed {
		return notary.Errf(notary.KindAnchorRefused, op, nil, "the ledger is sealed; point the notary at its successor")
	}

	if hasLast && (head.NextEpoch != last.Header.Number+1 || head.Head != last.Digest) {
		return notary.Errf(
			notary.KindLedgerDiverged,
			op,
			nil,
			"the ledger is at epoch %d with head %s, this database finalized epoch %d with digest %s",
			head.NextEpoch,
			head.Head.Hex(),
			last.Header.Number,
			last.Digest.Hex(),
		)
	}

	if !hasLast && !head.AtGenesis() {
		return notary.Errf(
			notary.KindLedgerDiverged,
			op,
			nil,
			"the ledger already holds epochs up to %d, and this database built none of them",
			head.NextEpoch,
		)
	}

	snapshot, err := a.source.Snapshot(ctx, a.policy.Currency)
	if err != nil {
		return err
	}

	if snapshot.AsOf().Seconds() <= head.LastAsOf {
		return notary.Errf(
			notary.KindAnchorDeferred,
			op,
			nil,
			"the snapshot of %d is not after the ledger's last at %d",
			snapshot.AsOf().Seconds(),
			head.LastAsOf,
		)
	}

	plan, err := notary.PlanEpoch(a.scheme, a.salts, head.NextEpoch, head.Head, snapshot)
	if err != nil {
		return err
	}

	record, err := a.store.SaveBuilt(ctx, lease, a.policy.Ledger, plan)
	if err != nil {
		return err
	}

	a.log.InfoContext(
		ctx,
		"epoch built",
		"epoch",
		record.Header.Number,
		"accounts",
		record.Header.LeafCount,
		"digest",
		record.Digest.Hex(),
		"as_of",
		record.Header.AsOf.Time(),
	)

	return a.advance(ctx, lease, record)
}

func (a *Anchorer) advance(ctx context.Context, lease notary.Lease, record notary.EpochRecord) error {
	if record.Ledger != a.policy.Ledger {
		return notary.Errf(
			notary.KindLedgerDiverged,
			"anchoring.Anchorer.advance",
			nil,
			"epoch %d was built for %s on chain %d",
			record.Header.Number,
			record.Ledger.Address.Hex(),
			record.Ledger.ChainID,
		)
	}

	if record.State == notary.StateBuilt {
		return a.submit(ctx, lease, record)
	}

	return a.watch(ctx, lease, record)
}

func (a *Anchorer) submit(ctx context.Context, lease notary.Lease, record notary.EpochRecord) error {
	nonce, err := a.chain.Nonce(ctx)
	if err != nil {
		return err
	}

	base, tip, err := a.chain.Fees(ctx)
	if err != nil {
		return err
	}

	quote, err := notary.QuoteFees(base, tip, a.policy.FeeCeiling)
	if err != nil {
		if notary.KindOf(err) == notary.KindAnchorDeferred {
			a.observer.FeeCeilingHit()
		}

		return err
	}

	anchor, err := a.chain.Prepare(ctx, record.Header, nonce, quote)
	switch notary.KindOf(err) {
	case notary.KindUnknown:
	case notary.KindAlreadyAnchored:
		if record.State == notary.StateBuilt {
			return a.store.Transition(ctx, lease, record.Header.Number, notary.StateBuilt, notary.StateSubmitted, "")
		}

		return nil
	case notary.KindAnchorRefused:
		return a.abandon(ctx, lease, record, err)
	default:
		return err
	}

	attempt := notary.Attempt{
		Epoch:   record.Header.Number,
		Anchor:  anchor,
		Outcome: notary.AttemptPending,
		SentAt:  a.clock(),
	}
	if err := a.store.RecordAttempt(ctx, lease, attempt); err != nil {
		return err
	}

	a.log.InfoContext(
		ctx,
		"anchor sent",
		"epoch",
		record.Header.Number,
		"tx",
		anchor.Hash.Hex(),
		"nonce",
		anchor.Nonce,
		"max_fee_per_gas",
		anchor.Fees.MaxFeePerGas,
		"max_priority_fee_per_gas",
		anchor.Fees.MaxPriorityFeePerGas,
	)

	return a.chain.Broadcast(ctx, anchor)
}

func (a *Anchorer) watch(ctx context.Context, lease notary.Lease, record notary.EpochRecord) error {
	number := record.Header.Number

	final, found, err := a.chain.EpochAtFinalized(ctx, number)
	if err != nil {
		return err
	}

	if found {
		if final.Digest != record.Digest {
			return a.abandon(ctx, lease, record, conflict(number, final.Digest, record.Digest))
		}

		return a.finalize(ctx, lease, record)
	}

	latest, found, err := a.chain.EpochAtLatest(ctx, number)
	if err != nil {
		return err
	}

	if found {
		if latest.Digest != record.Digest {
			return a.abandon(ctx, lease, record, conflict(number, latest.Digest, record.Digest))
		}

		if record.State == notary.StateSubmitted {
			return a.store.Transition(ctx, lease, number, notary.StateSubmitted, notary.StateMined, "")
		}

		return nil
	}

	if record.State == notary.StateMined {
		a.log.WarnContext(ctx, "anchored epoch left the chain; resubmitting", "epoch", number)

		if err := a.store.Transition(ctx, lease, number, notary.StateMined, notary.StateSubmitted, ""); err != nil {
			return err
		}

		record.State = notary.StateSubmitted
	}

	return a.keepAlive(ctx, lease, record)
}

func (a *Anchorer) keepAlive(ctx context.Context, lease notary.Lease, record notary.EpochRecord) error {
	const op = "anchoring.Anchorer.keepAlive"

	attempts, err := a.store.Attempts(ctx, record.Header.Number)
	if err != nil {
		return err
	}

	if len(attempts) == 0 {
		return a.submit(ctx, lease, record)
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

		receipt, mined, err := a.chain.Receipt(ctx, attempt.Anchor.Hash)
		if err != nil {
			return err
		}

		if mined && !receipt.Succeeded {
			if err := a.store.SettleAttempt(
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

	nonce, err := a.chain.Nonce(ctx)
	if err != nil {
		return err
	}

	if nonce > latest.Anchor.Nonce {
		if reverted {
			a.log.WarnContext(
				ctx,
				"anchor reverted; preparing a fresh one",
				"epoch",
				record.Header.Number,
				"nonce",
				nonce,
			)

			return a.submit(ctx, lease, record)
		}

		if a.anyMined(ctx, attempts) {
			return nil
		}

		return notary.Errf(
			notary.KindLedgerDiverged,
			op,
			nil,
			"nonce %d was used by a transaction this notary never sent, and epoch %d is not on chain",
			latest.Anchor.Nonce,
			record.Header.Number,
		)
	}

	if a.clock().Sub(latest.SentAt) >= a.policy.ReplaceAfter {
		return a.replace(ctx, lease, latest)
	}

	known, err := a.chain.Known(ctx, latest.Anchor.Hash)
	if err != nil {
		return err
	}

	if !known {
		a.log.InfoContext(ctx, "anchor dropped by the node; broadcasting again", "tx", latest.Anchor.Hash.Hex())
		return a.chain.Broadcast(ctx, latest.Anchor)
	}

	return nil
}

func (a *Anchorer) anyMined(ctx context.Context, attempts []notary.Attempt) bool {
	for _, attempt := range attempts {
		if _, mined, err := a.chain.Receipt(ctx, attempt.Anchor.Hash); err == nil && mined {
			return true
		}
	}

	return false
}

func (a *Anchorer) replace(ctx context.Context, lease notary.Lease, latest notary.Attempt) error {
	bumped, raised := latest.Anchor.Fees.Bumped(a.policy.FeeCeiling)
	if !raised {
		a.observer.FeeCeilingHit()
		a.log.WarnContext(
			ctx,
			"anchor is stuck at the fee ceiling; broadcasting it unchanged",
			"tx",
			latest.Anchor.Hash.Hex(),
			"max_fee_per_gas",
			latest.Anchor.Fees.MaxFeePerGas,
		)

		return a.chain.Broadcast(ctx, latest.Anchor)
	}

	anchor, err := a.chain.Resign(latest.Anchor, bumped)
	if err != nil {
		return err
	}

	attempt := notary.Attempt{Epoch: latest.Epoch, Anchor: anchor, Outcome: notary.AttemptPending, SentAt: a.clock()}
	if err := a.store.RecordAttempt(ctx, lease, attempt); err != nil {
		return err
	}

	a.log.InfoContext(
		ctx,
		"anchor replaced",
		"epoch",
		latest.Epoch,
		"nonce",
		anchor.Nonce,
		"replaces",
		latest.Anchor.Hash.Hex(),
		"tx",
		anchor.Hash.Hex(),
		"max_fee_per_gas",
		anchor.Fees.MaxFeePerGas,
	)

	return a.chain.Broadcast(ctx, anchor)
}

func (a *Anchorer) finalize(ctx context.Context, lease notary.Lease, record notary.EpochRecord) error {
	attempts, err := a.store.Attempts(ctx, record.Header.Number)
	if err != nil {
		return err
	}

	for _, attempt := range attempts {
		if attempt.Outcome != notary.AttemptPending && attempt.Outcome != notary.AttemptIncluded {
			continue
		}

		receipt, mined, err := a.chain.Receipt(ctx, attempt.Anchor.Hash)
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

		if err := a.store.SettleAttempt(ctx, lease, attempt.Anchor.Hash, outcome, settled); err != nil {
			return err
		}
	}

	if err := a.store.Transition(
		ctx,
		lease,
		record.Header.Number,
		record.State,
		notary.StateFinalized,
		"",
	); err != nil {
		return err
	}

	a.log.InfoContext(ctx, "epoch final", "epoch", record.Header.Number, "digest", record.Digest.Hex())

	return nil
}

func (a *Anchorer) abandon(ctx context.Context, lease notary.Lease, record notary.EpochRecord, cause error) error {
	if err := a.store.Transition(
		ctx,
		lease,
		record.Header.Number,
		record.State,
		notary.StateAbandoned,
		cause.Error(),
	); err != nil {
		return err
	}

	a.log.ErrorContext(
		ctx,
		"epoch abandoned; a person must decide what happens next",
		"epoch",
		record.Header.Number,
		"reason",
		cause.Error(),
	)

	return cause
}

func (a *Anchorer) report(ctx context.Context) error {
	last, found, err := a.store.LastFinalized(ctx)
	if err != nil {
		return err
	}

	if found {
		a.observer.LastFinalized(last)
	}

	abandoned, err := a.store.AbandonedCount(ctx)
	if err != nil {
		return err
	}

	a.observer.Abandoned(abandoned)

	balance, err := a.chain.Balance(ctx)
	if err != nil {
		return err
	}

	costs, err := a.store.RecentCosts(ctx, a.policy.CostWindow)
	if err != nil {
		return err
	}

	mean := meanOf(costs)
	runway, known := notary.Runway(balance, mean, a.policy.EpochInterval)
	a.observer.Funds(balance, mean, runway, known)

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

func conflict(number uint64, onChain, ours notary.Digest) error {
	return notary.Errf(
		notary.KindAnchorRefused,
		"anchoring.Anchorer.watch",
		nil,
		"the ledger holds epoch %d as %s, this notary built %s",
		number,
		onChain.Hex(),
		ours.Hex(),
	)
}
