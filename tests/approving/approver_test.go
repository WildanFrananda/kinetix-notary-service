package approving_test

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/approving"
	"github.com/WildanFrananda/kinetix-notary-service/internal/keccak"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/poseidon"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

const epoch = 4

var ledger = notary.Address{0x4e, 0x86}

type rig struct {
	t        *testing.T
	records  *notarytest.SnapshotRecords
	heads    *notarytest.Heads
	council  *notarytest.Council
	signer   notary.ApprovalSigner
	observer *notarytest.ApprovalObserver
	plan     notary.EpochPlan
	proposal notary.EpochProposal
	proofs   *notarytest.Proofs
	two      notary.ConfidentialScheme
}

func newRig(t *testing.T) *rig {
	t.Helper()

	scheme := keccak.NewScheme()
	snapshot := notarytest.Snapshot(t, 1760000000, notarytest.Liabilities(t, 9))
	previous := notary.Digest{0x0a}
	plan, err := notary.PlanEpoch(
		scheme,
		&notarytest.Salts{Seed: 3},
		&notarytest.Secrets{Seed: 3},
		nil,
		epoch,
		previous,
		snapshot,
	)
	require.NoError(t, err)

	proposed := make([]notary.ProposedAccount, len(plan.Accounts))
	for i, account := range plan.Accounts {
		proposed[i] = notary.ProposedAccount{Liability: account.Liability, Secret: account.Secret, Salt: account.Salt}
	}

	signer := notarytest.ApprovalSigner(t, "watch")

	return &rig{
		t: t,
		records: &notarytest.SnapshotRecords{
			Kept: map[notary.SnapshotID]notary.SnapshotRecord{snapshot.Record().ID: snapshot.Record()},
		},
		heads: &notarytest.Heads{Value: notary.LedgerHead{NextEpoch: epoch, Head: previous, LastAsOf: 1759990000}},
		council: &notarytest.Council{
			Answers:        notarytest.CouncilIdentity,
			Ledger:         ledger,
			Members:        map[notary.Address]bool{signer.Member(): true},
			ThresholdValue: 2,
			Anchors:        true,
			Scheme:         scheme,
		},
		signer:   signer,
		observer: &notarytest.ApprovalObserver{},
		plan:     plan,
		proposal: notary.EpochProposal{SnapshotID: plan.SnapshotID, Header: plan.Header, Accounts: proposed},
		proofs:   &notarytest.Proofs{},
		two:      poseidon.NewScheme(),
	}
}

func (r *rig) tb() *testing.T {
	return r.t
}

func (r *rig) approver() *approving.Approver {
	return approving.NewApprover(
		r.records,
		r.heads,
		r.council,
		r.signer,
		keccak.NewScheme(),
		r.two,
		r.proofs,
		keccak.NewScheme(),
		notarytest.CouncilIdentity,
		ledger,
		r.observer,
		slog.New(slog.NewTextHandler(io.Discard, nil)),
	)
}

func (r *rig) approve(t *testing.T) (notary.Approval, error) {
	t.Helper()

	return r.approver().Approve(context.Background(), r.proposal)
}

func TestAnHonestProposalIsApprovedWithASignatureTheCouncilWouldCount(t *testing.T) {
	r := newRig(t)

	approval, err := r.approve(t)

	require.NoError(t, err)
	require.Equal(t, r.signer.Member(), approval.Member)

	digest := keccak.NewScheme().ApprovalDigest(notarytest.CouncilIdentity, ledger, epoch, r.plan.Digest)

	raw := append([]byte(nil), approval.Signature[:]...)
	raw[64] -= 27
	public, err := crypto.SigToPub(digest[:], raw)
	require.NoError(t, err)
	require.Equal(t, r.signer.Member(), notary.Address(crypto.PubkeyToAddress(*public)))
	require.Equal(t, []string{"approved"}, r.observer.Outcomes)
}

func TestTheProposalsOwnOrderDoesNotMatter(t *testing.T) {
	r := newRig(t)

	accounts := r.proposal.Accounts
	for i, j := 0, len(accounts)-1; i < j; i, j = i+1, j-1 {
		accounts[i], accounts[j] = accounts[j], accounts[i]
	}

	_, err := r.approve(t)

	require.NoError(t, err)
}

func TestEveryDisagreementWithPaymentOrTheLedgerIsRefused(t *testing.T) {
	cases := map[string]func(r *rig){
		"an account left out": func(r *rig) {
			r.proposal.Accounts = r.proposal.Accounts[1:]
		},
		"an amount changed": func(r *rig) {
			r.proposal.Accounts[0].Liability.Amount++
		},
		"an amount moved between accounts": func(r *rig) {
			r.proposal.Accounts[0].Liability.Amount++
			r.proposal.Accounts[1].Liability.Amount--
		},
		"a holder swapped": func(r *rig) {
			r.proposal.Accounts[0].Liability.Holder = notarytest.Holder(r.tb(), "someone-else")
		},
		"another snapshot time": func(r *rig) {
			r.proposal.Header.AsOf = notarytest.Unix(r.tb(), 1760000001)
		},
		"not the next epoch": func(r *rig) {
			r.proposal.Header.Number = epoch + 1
		},
		"after another head": func(r *rig) {
			r.proposal.Header.PreviousDigest = notary.Digest{0x0b}
		},
		"a root the accounts do not build": func(r *rig) {
			r.proposal.Header.Root = notary.Digest{0xde, 0xad}
		},
		"a total the accounts do not add up to": func(r *rig) {
			r.proposal.Header.Total++
		},
		"a leaf count of its own": func(r *rig) {
			r.proposal.Header.LeafCount++
		},
		"another salt": func(r *rig) {
			r.proposal.Accounts[2].Salt[0] ^= 0xff
		},
		"another secret": func(r *rig) {
			r.proposal.Accounts[2].Secret[0] ^= 0xff
		},
		"no secret": func(r *rig) {
			r.proposal.Accounts[2].Secret = notary.AccountSecret{}
		},
		"no salt": func(r *rig) {
			r.proposal.Accounts[2].Salt = notary.Salt{}
		},
		"a sealed ledger": func(r *rig) {
			r.heads.Value.Sealed = true
		},
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			r := newRig(t)
			r.t = t
			mutate(r)

			_, err := r.approve(t)

			require.Equal(t, notary.KindApprovalRefused, notary.KindOf(err), err)
			require.Equal(t, []string{"refused"}, r.observer.Outcomes)
		})
	}
}

func TestASnapshotPaymentNeverTookIsRefusedAndAnOutageIsNot(t *testing.T) {
	r := newRig(t)
	r.proposal.SnapshotID = notarytest.SnapshotID(t, 99)

	_, err := r.approve(t)
	require.Equal(t, notary.KindSnapshotUnknown, notary.KindOf(err))
	require.Equal(t, []string{"refused"}, r.observer.Outcomes)

	r = newRig(t)
	r.records.Err = notary.Errf(notary.KindSourceUnavailable, "test", nil, "")

	_, err = r.approve(t)
	require.Equal(t, notary.KindSourceUnavailable, notary.KindOf(err))
	require.Equal(t, []string{"failed"}, r.observer.Outcomes)

	r = newRig(t)
	r.heads.Err = notary.Errf(notary.KindLedgerUnavailable, "test", nil, "")

	_, err = r.approve(t)
	require.Equal(t, notary.KindLedgerUnavailable, notary.KindOf(err))
	require.Equal(t, []string{"failed"}, r.observer.Outcomes)
}

func TestTheWatchChecksItsCouncilBeforeItServes(t *testing.T) {
	ctx := context.Background()

	r := newRig(t)
	require.NoError(t, r.approver().CheckCouncil(ctx))

	r = newRig(t)
	r.council.Answers = notary.CouncilIdentity{ChainID: 1, Address: notarytest.CouncilIdentity.Address}
	require.Equal(t, notary.KindWrongCouncil, notary.KindOf(r.approver().CheckCouncil(ctx)))

	r = newRig(t)
	r.council.Members = map[notary.Address]bool{}
	require.Equal(t, notary.KindWrongCouncil, notary.KindOf(r.approver().CheckCouncil(ctx)))

	r = newRig(t)
	r.council.Ledger = notary.Address{0x01}
	require.Equal(t, notary.KindDisagreement, notary.KindOf(r.approver().CheckCouncil(ctx)))

	r = newRig(t)
	r.council.Err = notary.Errf(notary.KindLedgerUnavailable, "test", nil, "")
	require.Equal(t, notary.KindLedgerUnavailable, notary.KindOf(r.approver().CheckCouncil(ctx)))
}
