package approving

import (
	"context"
	"log/slog"
	"math"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type Approver struct {
	records      notary.SnapshotRecords
	heads        notary.LedgerHeads
	council      notary.Council
	signer       notary.ApprovalSigner
	scheme       notary.Scheme
	confidential notary.ConfidentialScheme
	proofs       notary.ProofReader
	approval     notary.ApprovalScheme
	identity     notary.CouncilIdentity
	ledger       notary.Address
	observer     Observer
	log          *slog.Logger
}

func NewApprover(
	records notary.SnapshotRecords,
	heads notary.LedgerHeads,
	council notary.Council,
	signer notary.ApprovalSigner,
	scheme notary.Scheme,
	confidential notary.ConfidentialScheme,
	proofs notary.ProofReader,
	approval notary.ApprovalScheme,
	identity notary.CouncilIdentity,
	ledger notary.Address,
	observer Observer,
	log *slog.Logger,
) *Approver {
	return &Approver{
		records:      records,
		heads:        heads,
		council:      council,
		signer:       signer,
		scheme:       scheme,
		confidential: confidential,
		proofs:       proofs,
		approval:     approval,
		identity:     identity,
		ledger:       ledger,
		observer:     observer,
		log:          log,
	}
}

func (a *Approver) CheckCouncil(ctx context.Context) error {
	const op = "approving.Approver.CheckCouncil"

	identity, err := a.council.Identity(ctx)
	if err != nil {
		return err
	}

	if identity != a.identity {
		return notary.Errf(
			notary.KindWrongCouncil,
			op,
			nil,
			"the council answers as %s on chain %d, this watch approves for %s on chain %d",
			identity.Address.Hex(),
			identity.ChainID,
			a.identity.Address.Hex(),
			a.identity.ChainID,
		)
	}

	member, err := a.council.IsMember(ctx, a.signer.Member())
	if err != nil {
		return err
	}

	if !member {
		return notary.Errf(
			notary.KindWrongCouncil,
			op,
			nil,
			"%s, this watch's key, is not a member of the council",
			a.signer.Member().Hex(),
		)
	}

	sample := notary.Digest{4}

	onChain, err := a.council.ApprovalDigest(ctx, 3, sample)
	if err != nil {
		return err
	}

	if local := a.approval.ApprovalDigest(a.identity, a.ledger, 3, sample); local != onChain {
		return notary.Errf(
			notary.KindDisagreement,
			op,
			nil,
			"the council hashes an approval to %s, this watch to %s",
			onChain.Hex(),
			local.Hex(),
		)
	}

	return nil
}

func (a *Approver) Approve(ctx context.Context, proposal notary.EpochProposal) (notary.Approval, error) {
	approval, err := a.approve(ctx, proposal)

	switch {
	case err == nil:
		a.observer.Approved()
		a.log.InfoContext(
			ctx,
			"epoch approved",
			"epoch",
			proposal.Header.Number,
			"snapshot",
			proposal.SnapshotID,
			"accounts",
			len(proposal.Accounts),
		)
	case refusal(err):
		a.observer.Refused()
		a.log.ErrorContext(
			ctx,
			"epoch refused",
			"epoch",
			proposal.Header.Number,
			"snapshot",
			proposal.SnapshotID,
			"error",
			err,
		)
	default:
		a.observer.Failed()
		a.log.WarnContext(
			ctx,
			"epoch could not be checked",
			"epoch",
			proposal.Header.Number,
			"snapshot",
			proposal.SnapshotID,
			"error",
			err,
		)
	}

	return approval, err
}

func (a *Approver) approve(ctx context.Context, proposal notary.EpochProposal) (notary.Approval, error) {
	const op = "approving.Approver.Approve"

	record, err := a.records.Record(ctx, proposal.SnapshotID)
	if err != nil {
		return notary.Approval{}, err
	}

	liabilities, total, err := liabilitiesOf(op, proposal)
	if err != nil {
		return notary.Approval{}, err
	}

	if record.Count != uint64(len(liabilities)) || record.Total != total {
		return notary.Approval{}, refuse(
			op,
			"payment recorded %d accounts owed %d in snapshot %s, the notary proposes %d owed %d",
			record.Count,
			record.Total,
			record.ID,
			len(liabilities),
			total,
		)
	}

	digest, err := notary.SnapshotDigest(record, liabilities)
	if err != nil {
		return notary.Approval{}, refuse(op, "the proposed accounts cannot be digested: %v", err)
	}

	if digest != record.Digest {
		return notary.Approval{}, refuse(
			op,
			"the accounts the notary proposes digest to %s, payment recorded snapshot %s as %s",
			digest.Hex(),
			record.ID,
			record.Digest.Hex(),
		)
	}

	if proposal.Header.AsOf != record.AsOf {
		return notary.Approval{}, refuse(
			op,
			"the proposal is dated %d, payment took snapshot %s at %d",
			proposal.Header.AsOf.Seconds(),
			record.ID,
			record.AsOf.Seconds(),
		)
	}

	head, err := a.heads.Head(ctx)
	if err != nil {
		return notary.Approval{}, err
	}

	if head.Sealed {
		return notary.Approval{}, refuse(op, "the ledger is sealed")
	}

	if proposal.Header.Number != head.NextEpoch || proposal.Header.PreviousDigest != head.Head {
		return notary.Approval{}, refuse(
			op,
			"the proposal is epoch %d after %s, the ledger expects epoch %d after %s",
			proposal.Header.Number,
			proposal.Header.PreviousDigest.Hex(),
			head.NextEpoch,
			head.Head.Hex(),
		)
	}

	plan, err := a.rebuild(op, proposal, record)
	if err != nil {
		return notary.Approval{}, err
	}

	if plan.Header != proposal.Header {
		return notary.Approval{}, refuse(
			op,
			"the watch rebuilds root %s, total %d and %d leaves; the proposal claims %s, %d and %d",
			plan.Header.Root.Hex(),
			plan.Header.Total,
			plan.Header.LeafCount,
			proposal.Header.Root.Hex(),
			proposal.Header.Total,
			proposal.Header.LeafCount,
		)
	}

	if plan.Confidential != nil {
		if err := a.checkProven(ctx, op, plan); err != nil {
			return notary.Approval{}, err
		}
	}

	signature, err := a.signer.Sign(a.approval.ApprovalDigest(a.identity, a.ledger, plan.Header.Number, plan.Digest))
	if err != nil {
		return notary.Approval{}, err
	}

	return notary.Approval{Member: a.signer.Member(), Signature: signature}, nil
}

func (a *Approver) checkProven(ctx context.Context, op string, plan notary.EpochPlan) error {
	depth, err := a.proofs.Depth(ctx)
	if err != nil {
		return err
	}

	if depth != plan.Confidential.Depth {
		return refuse(op, "the proposal is a tree of depth %d, the ledger's is %d", plan.Confidential.Depth, depth)
	}

	proven, err := a.proofs.IsProven(ctx, plan.Header.Number, plan.Header.Root, plan.Header.LeafCount)
	if err != nil {
		return err
	}

	if !proven {
		return notary.Errf(
			notary.KindAnchorDeferred,
			op,
			nil,
			"epoch %d is not proven on chain yet; the watch approves only a proven statement",
			plan.Header.Number,
		)
	}

	return nil
}

func (a *Approver) rebuild(
	op string,
	proposal notary.EpochProposal,
	record notary.SnapshotRecord,
) (notary.EpochPlan, error) {
	accounts := make([]notary.Account, len(proposal.Accounts))
	for i, proposed := range proposal.Accounts {
		if proposed.Secret.IsZero() || proposed.Salt.IsZero() {
			return notary.EpochPlan{}, refuse(
				op,
				"%s of %s is proposed without a secret or a salt",
				proposed.Liability.Kind,
				proposed.Liability.Holder,
			)
		}

		accounts[i] = notary.OpenAccount(a.scheme, proposed.Liability, proposed.Salt, proposed.Secret)
	}

	if proposal.Confidential != nil {
		return a.rebuildConfidential(op, proposal, record, accounts)
	}

	plan, err := notary.AssembleEpoch(
		a.scheme,
		proposal.Header.Number,
		proposal.Header.PreviousDigest,
		record.AsOf,
		record.Currency,
		accounts,
	)
	if err != nil {
		return notary.EpochPlan{}, refuse(op, "the proposed accounts do not make a tree: %v", err)
	}

	return plan, nil
}

func (a *Approver) rebuildConfidential(
	op string,
	proposal notary.EpochProposal,
	record notary.SnapshotRecord,
	accounts []notary.Account,
) (notary.EpochPlan, error) {
	claimed := proposal.Confidential
	if a.confidential == nil || a.proofs == nil || a.confidential.Version() != claimed.Scheme {
		return notary.EpochPlan{}, refuse(op, "this watch does not approve scheme %d", claimed.Scheme)
	}

	plan, err := notary.AssembleConfidentialEpoch(
		a.confidential,
		claimed.Depth,
		claimed.Blinding,
		proposal.Header.Number,
		proposal.Header.PreviousDigest,
		record.AsOf,
		record.Currency,
		accounts,
	)
	if err != nil {
		return notary.EpochPlan{}, refuse(op, "the proposed accounts do not make a tree: %v", err)
	}

	if plan.Confidential.TreeRoot != claimed.TreeRoot || plan.Confidential.TotalCommitment != claimed.TotalCommitment {
		return notary.EpochPlan{}, refuse(
			op,
			"the watch rebuilds tree root %s and commitment %s; the proposal claims %s and %s",
			plan.Confidential.TreeRoot.Hex(),
			plan.Confidential.TotalCommitment.Hex(),
			claimed.TreeRoot.Hex(),
			claimed.TotalCommitment.Hex(),
		)
	}

	return plan, nil
}

func liabilitiesOf(op string, proposal notary.EpochProposal) ([]notary.Liability, uint64, error) {
	liabilities := make([]notary.Liability, len(proposal.Accounts))

	var total uint64

	for i, proposed := range proposal.Accounts {
		if proposed.Liability.Amount > math.MaxInt64-total {
			return nil, 0, refuse(op, "the proposed accounts sum past %d minor units", uint64(math.MaxInt64))
		}

		total += proposed.Liability.Amount
		liabilities[i] = proposed.Liability
	}

	return liabilities, total, nil
}

func refuse(op, format string, args ...any) error {
	return notary.Errf(notary.KindApprovalRefused, op, nil, format, args...)
}

func refusal(err error) bool {
	switch notary.KindOf(err) {
	case notary.KindApprovalRefused, notary.KindSnapshotUnknown, notary.KindDisagreement:
		return true
	default:
		return false
	}
}
