package approving_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/keccak"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/poseidon"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

const confidentialDepth = 4

func newConfidentialRig(t *testing.T) *rig {
	t.Helper()

	r := newRig(t)
	snapshot := notarytest.Snapshot(t, 1760000000, notarytest.Liabilities(t, 9))
	plan, err := notary.PlanConfidentialEpoch(
		keccak.NewScheme(),
		poseidon.NewScheme(),
		confidentialDepth,
		&notarytest.Salts{Seed: 3},
		&notarytest.Secrets{Seed: 3},
		&notarytest.Blindings{Seed: 3},
		nil,
		epoch,
		notary.Digest{0x0a},
		snapshot,
	)
	require.NoError(t, err)

	r.records.Kept[snapshot.Record().ID] = snapshot.Record()
	r.plan = plan
	r.proposal.SnapshotID = plan.SnapshotID
	r.proposal.Header = plan.Header
	r.proposal.Accounts = make([]notary.ProposedAccount, len(plan.Accounts))

	for i, account := range plan.Accounts {
		r.proposal.Accounts[i] = notary.ProposedAccount{
			Liability: account.Liability,
			Secret:    account.Secret,
			Salt:      account.Salt,
		}
	}

	statement := *plan.Confidential
	r.proposal.Confidential = &statement
	r.proofs.DepthValue = confidentialDepth
	r.proofs.Proven = true

	return r
}

func TestTheWatchApprovesAProvenSchemeTwoEpochOverItsSchemeTwoDigest(t *testing.T) {
	r := newConfidentialRig(t)

	approval, err := r.approve(t)
	require.NoError(t, err)
	require.Equal(t, r.signer.Member(), approval.Member)
	require.Equal(t, 1, r.proofs.Asked)

	digest := keccak.NewScheme().ApprovalDigest(
		notarytest.CouncilIdentity,
		ledger,
		r.plan.Header.Number,
		poseidon.NewScheme().EpochDigest(r.plan.Header),
	)

	signature := approval.Signature
	if signature[64] >= 27 {
		signature[64] -= 27
	}

	public, err := crypto.SigToPub(digest[:], signature[:])
	require.NoError(t, err)
	require.Equal(t, r.signer.Member(), notary.Address(crypto.PubkeyToAddress(*public)))
}

func TestTheWatchWaitsForAStatementNotYetProven(t *testing.T) {
	r := newConfidentialRig(t)
	r.proofs.Proven = false

	_, err := r.approve(t)
	require.Equal(t, notary.KindAnchorDeferred, notary.KindOf(err))
}

func TestTheWatchRefusesASchemeTwoProposalThatIsNotPaymentsTree(t *testing.T) {
	cases := map[string]func(*rig){
		"another tree root":   func(r *rig) { r.proposal.Confidential.TreeRoot[0] ^= 0x01 },
		"another commitment":  func(r *rig) { r.proposal.Confidential.TotalCommitment[0] ^= 0x01 },
		"another blinding":    func(r *rig) { r.proposal.Confidential.Blinding[0] ^= 0x01 },
		"another depth":       func(r *rig) { r.proposal.Confidential.Depth = confidentialDepth + 1 },
		"a total shown":       func(r *rig) { r.proposal.Header.Total = 5 },
		"the ledger's depth":  func(r *rig) { r.proofs.DepthValue = confidentialDepth + 1 },
		"a watch without two": func(r *rig) { r.two = nil },
	}

	for name, tamper := range cases {
		r := newConfidentialRig(t)
		tamper(r)

		_, err := r.approve(t)
		require.Equal(t, notary.KindApprovalRefused, notary.KindOf(err), name)
		require.Zero(t, r.proofs.Asked, name)
	}
}
