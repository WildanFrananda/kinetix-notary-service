package anchoring_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/keccak"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

func recovered(t *testing.T, digest notary.Digest, signature notary.Signature) notary.Address {
	t.Helper()

	raw := append([]byte(nil), signature[:]...)
	raw[64] -= 27
	public, err := crypto.SigToPub(digest[:], raw)
	require.NoError(t, err)

	return notary.Address(crypto.PubkeyToAddress(*public))
}

func TestEveryAnchorCarriesTheNotarysAndTheWatchsApprovalInAddressOrder(t *testing.T) {
	r := newRig(t)

	require.NoError(t, r.tick(t))
	require.Equal(t, notary.StateSubmitted, r.state(t, 0))

	record := r.store.Records[0]
	require.Len(t, r.watch.Proposals, 1)

	proposal := r.watch.Proposals[0]
	require.Equal(t, record.SnapshotID, proposal.SnapshotID)
	require.False(t, proposal.SnapshotID.IsZero())
	require.Equal(t, record.Header, proposal.Header)
	require.Len(t, proposal.Accounts, 12)

	for i, account := range r.store.Stored[0] {
		require.Equal(t, account.Liability, proposal.Accounts[i].Liability)
		require.Equal(t, account.Secret, proposal.Accounts[i].Secret)
		require.Equal(t, account.Salt, proposal.Accounts[i].Salt)
	}

	require.Len(t, r.chain.Approvals, 1)
	approvals := r.chain.Approvals[0]
	require.Len(t, approvals, 2)

	digest := keccak.NewScheme().ApprovalDigest(notarytest.CouncilIdentity, ledgerID.Address, 0, record.Digest)
	first := recovered(t, digest, approvals[0])
	second := recovered(t, digest, approvals[1])
	require.Negative(t, bytes.Compare(first[:], second[:]))
	require.ElementsMatch(
		t,
		[]notary.Address{r.signer.Member(), r.watch.Signer.Member()},
		[]notary.Address{first, second},
	)
}

func TestAWatchThatRefusesAbandonsTheEpochWithItsReason(t *testing.T) {
	r := newRig(t)
	r.watch.Err = notary.Errf(notary.KindApprovalRefused, "test", nil, "the rows are not the snapshot payment recorded")

	err := r.tick(t)

	require.Equal(t, notary.KindApprovalRefused, notary.KindOf(err))
	require.Equal(t, notary.StateAbandoned, r.state(t, 0))
	require.Contains(t, r.store.Records[0].AbandonedReason, "not the snapshot payment recorded")
	require.Empty(t, r.chain.Broadcasts)
	require.Empty(t, r.chain.Approvals)
}

func TestAWatchThatCannotBeReachedLeavesTheEpochBuiltUntilItAnswers(t *testing.T) {
	r := newRig(t)
	r.watch.Err = notary.Errf(notary.KindSourceUnavailable, "test", nil, "connection refused")

	err := r.tick(t)

	require.Equal(t, notary.KindSourceUnavailable, notary.KindOf(err))
	require.Equal(t, notary.StateBuilt, r.state(t, 0))
	require.Empty(t, r.chain.Broadcasts)

	r.watch.Err = nil

	require.NoError(t, r.tick(t))
	require.Equal(t, notary.StateSubmitted, r.state(t, 0))
	require.Len(t, r.watch.Proposals, 2)
	require.Equal(t, r.watch.Proposals[0], r.watch.Proposals[1])
}

func TestAWatchThatApprovesWithTheAnchorersKeyIsRefused(t *testing.T) {
	r := newRig(t)
	r.watch.Signer = r.signer

	err := r.tick(t)

	require.Equal(t, notary.KindApprovalRefused, notary.KindOf(err))
	require.Equal(t, notary.StateAbandoned, r.state(t, 0))
}

func TestAnEpochBuiltWithoutASnapshotIDIsNeverSentToTheWatch(t *testing.T) {
	r := newRig(t)
	r.source.Next = func() (notary.Snapshot, error) {
		return notary.NewSnapshot(
			notarytest.Unix(t, uint64(r.now.Unix())-30),
			"IDR",
			notarytest.Liabilities(t, 3),
		)
	}

	err := r.tick(t)

	require.Equal(t, notary.KindApprovalRefused, notary.KindOf(err))
	require.Equal(t, notary.StateAbandoned, r.state(t, 0))
	require.Empty(t, r.watch.Proposals)
}

func TestTheCouncilIsCheckedBeforeTheNotaryStarts(t *testing.T) {
	ctx := context.Background()

	r := newRig(t)
	require.NoError(t, r.anchorer.CheckCouncil(ctx))

	r = newRig(t)
	r.council.Answers = notary.CouncilIdentity{ChainID: 84532, Address: notary.Address{0x99}}
	require.Equal(t, notary.KindWrongCouncil, notary.KindOf(r.anchorer.CheckCouncil(ctx)))

	r = newRig(t)
	r.council.Anchors = false
	require.Equal(t, notary.KindNotAnchorer, notary.KindOf(r.anchorer.CheckCouncil(ctx)))

	r = newRig(t)
	delete(r.council.Members, r.signer.Member())
	require.Equal(t, notary.KindWrongCouncil, notary.KindOf(r.anchorer.CheckCouncil(ctx)))

	r = newRig(t)
	r.council.ThresholdValue = 3
	err := r.anchorer.CheckCouncil(ctx)
	require.Equal(t, notary.KindWrongCouncil, notary.KindOf(err))
	require.ErrorContains(t, err, "needs 3 approvals")

	r = newRig(t)
	r.council.Ledger = notary.Address{0x01}
	require.Equal(t, notary.KindDisagreement, notary.KindOf(r.anchorer.CheckCouncil(ctx)))

	r = newRig(t)
	r.chain.From = notary.Address{0x51, 0xb2}
	require.Equal(t, notary.KindInvalidArgument, notary.KindOf(r.anchorer.CheckCouncil(ctx)))

	r = newRig(t)
	r.council.Err = notary.Errf(notary.KindLedgerUnavailable, "test", nil, "")
	require.Equal(t, notary.KindLedgerUnavailable, notary.KindOf(r.anchorer.CheckCouncil(ctx)))
}
