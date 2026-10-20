package notarytest

import (
	"context"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type Digester interface {
	EpochDigest(header notary.EpochHeader) notary.Digest
}

type Watch struct {
	Signer    notary.ApprovalSigner
	Scheme    Digester
	Approvals notary.ApprovalScheme
	Council   notary.CouncilIdentity
	Ledger    notary.Address
	Err       error
	Proposals []notary.EpochProposal
}

var _ notary.EpochWatch = (*Watch)(nil)

func (w *Watch) Approve(_ context.Context, proposal notary.EpochProposal) (notary.Approval, error) {
	w.Proposals = append(w.Proposals, proposal)
	if w.Err != nil {
		return notary.Approval{}, w.Err
	}

	digest := w.Approvals.ApprovalDigest(
		w.Council,
		w.Ledger,
		proposal.Header.Number,
		w.Scheme.EpochDigest(proposal.Header),
	)

	signature, err := w.Signer.Sign(digest)
	if err != nil {
		return notary.Approval{}, err
	}

	return notary.Approval{Member: w.Signer.Member(), Signature: signature}, nil
}
