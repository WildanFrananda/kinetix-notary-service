package anchoring

import (
	"bytes"
	"context"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type Quorum struct {
	Council notary.Council
	Watch   notary.EpochWatch
	Signer  notary.ApprovalSigner
	Scheme  notary.ApprovalScheme
}

const gathered = 2

func (a *Anchorer) CheckCouncil(ctx context.Context) error {
	const op = "anchoring.Anchorer.CheckCouncil"

	identity, err := a.quorum.Council.Identity(ctx)
	if err != nil {
		return err
	}

	if identity != a.policy.Council {
		return notary.Errf(
			notary.KindWrongCouncil,
			op,
			nil,
			"the council answers as %s on chain %d, this notary anchors through %s on chain %d",
			identity.Address.Hex(),
			identity.ChainID,
			a.policy.Council.Address.Hex(),
			a.policy.Council.ChainID,
		)
	}

	if member := a.quorum.Signer.Member(); member != a.chain.Anchorer() {
		return notary.Errf(
			notary.KindInvalidArgument,
			op,
			nil,
			"the approval key %s is not the anchorer key %s",
			member.Hex(),
			a.chain.Anchorer().Hex(),
		)
	}

	anchors, err := a.quorum.Council.AnchorsLedger(ctx)
	if err != nil {
		return err
	}

	if !anchors {
		return notary.Errf(
			notary.KindNotAnchorer,
			op,
			nil,
			"the council %s does not hold the anchorer role on %s; the ledger admin grants it first",
			a.policy.Council.Address.Hex(),
			a.policy.Ledger.Address.Hex(),
		)
	}

	member, err := a.quorum.Council.IsMember(ctx, a.quorum.Signer.Member())
	if err != nil {
		return err
	}

	if !member {
		return notary.Errf(
			notary.KindWrongCouncil,
			op,
			nil,
			"%s, the anchorer key, is not a member of the council",
			a.quorum.Signer.Member().Hex(),
		)
	}

	threshold, err := a.quorum.Council.Threshold(ctx)
	if err != nil {
		return err
	}

	if threshold > gathered {
		return notary.Errf(
			notary.KindWrongCouncil,
			op,
			nil,
			"the council needs %d approvals, and this notary gathers %d: its own and the watch's",
			threshold,
			gathered,
		)
	}

	sample := notary.Digest{4}

	onChain, err := a.quorum.Council.ApprovalDigest(ctx, 3, sample)
	if err != nil {
		return err
	}

	if local := a.quorum.Scheme.ApprovalDigest(a.policy.Council, a.policy.Ledger.Address, 3, sample); local != onChain {
		return notary.Errf(
			notary.KindDisagreement,
			op,
			nil,
			"the council hashes an approval to %s, this notary to %s",
			onChain.Hex(),
			local.Hex(),
		)
	}

	return nil
}

func (a *Anchorer) approve(ctx context.Context, record notary.EpochRecord) ([]notary.Signature, error) {
	const op = "anchoring.Anchorer.approve"

	number := record.Header.Number
	if record.SnapshotID.IsZero() {
		return nil, notary.Errf(
			notary.KindApprovalRefused,
			op,
			nil,
			"epoch %d was built without a snapshot id, so the watch cannot check it",
			number,
		)
	}

	accounts, err := a.store.Accounts(ctx, number)
	if err != nil {
		return nil, err
	}

	proposed := make([]notary.ProposedAccount, len(accounts))
	for i, account := range accounts {
		proposed[i] = notary.ProposedAccount{Liability: account.Liability, Secret: account.Secret, Salt: account.Salt}
	}

	proposal := notary.EpochProposal{SnapshotID: record.SnapshotID, Header: record.Header, Accounts: proposed}
	if record.Confidential != nil {
		statement := *record.Confidential
		statement.Proof = nil
		proposal.Confidential = &statement
	}

	watched, err := a.quorum.Watch.Approve(ctx, proposal)
	if err != nil {
		return nil, err
	}

	own := a.quorum.Signer.Member()
	if watched.Member == own {
		return nil, notary.Errf(notary.KindApprovalRefused, op, nil, "the watch approved with the anchorer's own key")
	}

	signature, err := a.quorum.Signer.Sign(
		a.quorum.Scheme.ApprovalDigest(a.policy.Council, a.policy.Ledger.Address, number, record.Digest),
	)
	if err != nil {
		return nil, err
	}

	a.log.InfoContext(ctx, "epoch approved by the watch", "epoch", number, "watch", watched.Member.Hex())

	if bytes.Compare(own[:], watched.Member[:]) < 0 {
		return []notary.Signature{signature, watched.Signature}, nil
	}

	return []notary.Signature{watched.Signature, signature}, nil
}
