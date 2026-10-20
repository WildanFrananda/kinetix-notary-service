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
	secrets  notary.SecretSource
	source   notary.SnapshotSource
	store    notary.EpochStore
	chain    notary.AnchorChain
	quorum   Quorum
	observer Observer
	clock    func() time.Time
	log      *slog.Logger
	policy   Policy
	two      *SchemeTwo
}

func NewAnchorer(
	scheme notary.Scheme,
	salts notary.SaltSource,
	secrets notary.SecretSource,
	source notary.SnapshotSource,
	store notary.EpochStore,
	chain notary.AnchorChain,
	quorum Quorum,
	observer Observer,
	clock func() time.Time,
	log *slog.Logger,
	policy Policy,
	two *SchemeTwo,
) *Anchorer {
	return &Anchorer{
		scheme:   scheme,
		salts:    salts,
		secrets:  secrets,
		source:   source,
		store:    store,
		chain:    chain,
		quorum:   quorum,
		observer: observer,
		clock:    clock,
		log:      log,
		policy:   policy,
		two:      two,
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
		if err := a.succeeds(ctx, last); err != nil {
			return err
		}
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

	liabilities := snapshot.Liabilities()
	a.observer.Occupancy(uint64(len(liabilities)), a.capacity())

	keys := make([]notary.AccountKey, len(liabilities))
	for i, liability := range liabilities {
		keys[i] = notary.AccountKey{Holder: liability.Holder, Kind: liability.Kind}
	}

	known, err := a.store.AccountSecrets(ctx, keys)
	if err != nil {
		return err
	}

	plan, err := a.plan(ctx, known, head, snapshot)
	if err != nil {
		return err
	}

	if plan.Confidential != nil {
		renewed, held, err := a.store.Lease(ctx, a.policy.Instance, a.policy.LeaseTerm)
		if err != nil {
			return err
		}

		if !held {
			return notary.Errf(
				notary.KindLeaseLost,
				op,
				nil,
				"another writer took the lease while epoch %d was proven",
				head.NextEpoch,
			)
		}

		lease = renewed
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

func (a *Anchorer) capacity() uint64 {
	if a.two == nil {
		return 0
	}

	return uint64(1) << a.two.Depth
}

func (a *Anchorer) succeeds(ctx context.Context, last notary.EpochRecord) error {
	const op = "anchoring.Anchorer.succeeds"

	head, err := a.chain.Head(ctx)
	if err != nil {
		return err
	}

	if head.GenesisEpoch == last.Header.Number+1 && head.GenesisDigest == last.Digest {
		return nil
	}

	return notary.Errf(
		notary.KindLedgerDiverged,
		op,
		nil,
		"this database anchors to %s on chain %d, and %s on chain %d does not start where it ends: "+
			"its genesis is epoch %d with %s, this database finalized epoch %d with digest %s",
		last.Ledger.Address.Hex(),
		last.Ledger.ChainID,
		a.policy.Ledger.Address.Hex(),
		a.policy.Ledger.ChainID,
		head.GenesisEpoch,
		head.GenesisDigest.Hex(),
		last.Header.Number,
		last.Digest.Hex(),
	)
}

func (a *Anchorer) plan(
	ctx context.Context,
	known map[notary.AccountKey]notary.AccountSecret,
	head notary.LedgerHead,
	snapshot notary.Snapshot,
) (notary.EpochPlan, error) {
	if a.two == nil {
		return notary.PlanEpoch(a.scheme, a.salts, a.secrets, known, head.NextEpoch, head.Head, snapshot)
	}

	plan, err := notary.PlanConfidentialEpoch(
		a.scheme,
		a.two.Scheme,
		a.two.Depth,
		a.salts,
		a.secrets,
		a.two.Blindings,
		known,
		head.NextEpoch,
		head.Head,
		snapshot,
	)
	if err != nil {
		return notary.EpochPlan{}, err
	}

	entries := make([]notary.Entry, len(plan.Accounts))
	for i, account := range plan.Accounts {
		entries[i] = notary.Entry{Commitment: account.Commitment, Amount: account.Liability.Amount}
	}

	started := a.clock()
	proof, err := a.two.Prover.Prove(ctx, notary.ProofRequest{
		Number:          plan.Header.Number,
		Depth:           plan.Confidential.Depth,
		Entries:         entries,
		TreeRoot:        plan.Confidential.TreeRoot,
		TotalCommitment: plan.Confidential.TotalCommitment,
		Blinding:        plan.Confidential.Blinding,
	})
	a.observer.Proved(a.clock().Sub(started), err != nil)

	if err != nil {
		return notary.EpochPlan{}, err
	}

	plan.Confidential.Proof = proof

	a.log.InfoContext(
		ctx,
		"epoch proven",
		"epoch",
		plan.Header.Number,
		"accounts",
		plan.Header.LeafCount,
		"seconds",
		a.clock().Sub(started).Seconds(),
	)

	return plan, nil
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
	number := record.Header.Number

	latest, found, err := a.chain.EpochAtLatest(ctx, number)
	if err != nil {
		return err
	}

	if found {
		if latest.Digest != record.Digest {
			return a.abandon(ctx, lease, record, conflict(number, latest.Digest, record.Digest))
		}

		if record.State == notary.StateBuilt {
			return a.store.Transition(ctx, lease, number, notary.StateBuilt, notary.StateSubmitted, "")
		}

		return nil
	}

	if a.two != nil {
		proven, err := a.two.Chain.IsProven(ctx, number, record.Header.Root, record.Header.LeafCount)
		if err != nil {
			return err
		}

		if !proven {
			return a.keepAlive(ctx, lease, record, notary.AttemptForProof)
		}
	}

	approvals, err := a.approve(ctx, record)
	switch notary.KindOf(err) {
	case notary.KindUnknown:
	case notary.KindApprovalRefused:
		return a.abandon(ctx, lease, record, err)
	default:
		return err
	}

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

	anchor, err := a.chain.Prepare(ctx, record.Header, approvals, nonce, quote)
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

	return a.keepAlive(ctx, lease, record, notary.AttemptForAnchor)
}

func (a *Anchorer) keepAlive(
	ctx context.Context,
	lease notary.Lease,
	record notary.EpochRecord,
	kind notary.AttemptKind,
) error {
	const op = "anchoring.Anchorer.keepAlive"

	attempts, err := a.attemptsOf(ctx, record.Header.Number, kind)
	if err != nil {
		return err
	}

	if len(attempts) == 0 {
		return a.fresh(ctx, lease, record, kind)
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
				"transaction reverted; preparing a fresh one",
				"epoch",
				record.Header.Number,
				"kind",
				kind.String(),
				"nonce",
				nonce,
			)

			return a.fresh(ctx, lease, record, kind)
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

func (a *Anchorer) attemptsOf(ctx context.Context, epoch uint64, kind notary.AttemptKind) ([]notary.Attempt, error) {
	attempts, err := a.store.Attempts(ctx, epoch)
	if err != nil {
		return nil, err
	}

	kept := make([]notary.Attempt, 0, len(attempts))
	for _, attempt := range attempts {
		if attempt.Kind == kind {
			kept = append(kept, attempt)
		}
	}

	return kept, nil
}

func (a *Anchorer) fresh(
	ctx context.Context,
	lease notary.Lease,
	record notary.EpochRecord,
	kind notary.AttemptKind,
) error {
	if kind == notary.AttemptForProof {
		return a.sendProof(ctx, lease, record)
	}

	return a.submit(ctx, lease, record)
}

func (a *Anchorer) sendProof(ctx context.Context, lease notary.Lease, record notary.EpochRecord) error {
	const op = "anchoring.Anchorer.sendProof"

	number := record.Header.Number
	if a.two == nil || record.Confidential == nil || len(record.Confidential.Proof) == 0 {
		return a.abandon(
			ctx,
			lease,
			record,
			notary.Errf(notary.KindAnchorRefused, op, nil, "epoch %d holds no proof to send", number),
		)
	}

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

	prepared, err := a.two.Chain.PrepareProof(ctx, record.Header, *record.Confidential, nonce, quote)
	switch notary.KindOf(err) {
	case notary.KindUnknown:
	case notary.KindAlreadyAnchored:
		return nil
	case notary.KindAnchorRefused:
		return a.abandon(ctx, lease, record, err)
	default:
		return err
	}

	attempt := notary.Attempt{
		Epoch:   number,
		Kind:    notary.AttemptForProof,
		Anchor:  prepared,
		Outcome: notary.AttemptPending,
		SentAt:  a.clock(),
	}
	if err := a.store.RecordAttempt(ctx, lease, attempt); err != nil {
		return err
	}

	a.log.InfoContext(ctx, "proof sent", "epoch", number, "tx", prepared.Hash.Hex(), "nonce", prepared.Nonce)

	return a.chain.Broadcast(ctx, prepared)
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

	resign := a.chain.Resign
	if latest.Kind == notary.AttemptForProof && a.two != nil {
		resign = a.two.Chain.Resign
	}

	anchor, err := resign(latest.Anchor, bumped)
	if err != nil {
		return err
	}

	attempt := notary.Attempt{
		Epoch:   latest.Epoch,
		Kind:    latest.Kind,
		Anchor:  anchor,
		Outcome: notary.AttemptPending,
		SentAt:  a.clock(),
	}
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
