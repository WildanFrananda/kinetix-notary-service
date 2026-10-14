package anchoring_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/anchoring"
	"github.com/WildanFrananda/kinetix-notary-service/internal/keccak"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

var ledgerID = notary.LedgerIdentity{ChainID: 84532, Address: notary.Address{0x4e, 0x86}}

type rig struct {
	now      time.Time
	chain    *notarytest.Chain
	store    *notarytest.Store
	source   *notarytest.Source
	observer *notarytest.Observer
	council  *notarytest.Council
	watch    *notarytest.Watch
	signer   notary.ApprovalSigner
	anchorer *anchoring.Anchorer
	policy   anchoring.Policy
}

func newRig(t *testing.T) *rig {
	t.Helper()

	r := &rig{
		now:      time.Unix(1760000000, 0).UTC(),
		chain:    notarytest.NewChain(keccak.NewScheme()),
		observer: &notarytest.Observer{},
	}
	r.store = notarytest.NewStore(r.clock)
	r.signer = notarytest.ApprovalSigner(t, "anchorer")
	r.chain.From = r.signer.Member()
	watcher := notarytest.ApprovalSigner(t, "watch")
	r.watch = &notarytest.Watch{
		Signer:    watcher,
		Scheme:    keccak.NewScheme(),
		Approvals: keccak.NewScheme(),
		Council:   notarytest.CouncilIdentity,
		Ledger:    ledgerID.Address,
	}
	r.council = &notarytest.Council{
		Answers:        notarytest.CouncilIdentity,
		Ledger:         ledgerID.Address,
		Members:        map[notary.Address]bool{r.signer.Member(): true, watcher.Member(): true},
		ThresholdValue: 2,
		Anchors:        true,
		Scheme:         keccak.NewScheme(),
	}
	r.source = &notarytest.Source{Next: func() (notary.Snapshot, error) {
		return notarytest.Snapshot(t, uint64(r.now.Unix())-30, notarytest.Liabilities(t, 12)), nil
	}}
	r.policy = anchoring.Policy{
		Instance:      "notaryd-1",
		Ledger:        ledgerID,
		Council:       notarytest.CouncilIdentity,
		Currency:      "IDR",
		EpochInterval: time.Hour,
		ReplaceAfter:  5 * time.Minute,
		LeaseTerm:     30 * time.Second,
		FeeCeiling:    10_000,
		CostWindow:    24,
	}
	r.build()

	return r
}

func (r *rig) build() {
	r.anchorer = anchoring.NewAnchorer(
		keccak.NewScheme(),
		&notarytest.Salts{Seed: 7},
		&notarytest.Secrets{Seed: 7},
		r.source,
		r.store,
		r.chain,
		anchoring.Quorum{Council: r.council, Watch: r.watch, Signer: r.signer, Scheme: keccak.NewScheme()},
		r.observer,
		r.clock,
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		r.policy,
	)
}

func (r *rig) clock() time.Time { return r.now }

func (r *rig) advance(d time.Duration) {
	r.now = r.now.Add(d)
	r.chain.BlockTime = uint64(r.now.Unix())
}

func (r *rig) tick(t *testing.T) error {
	t.Helper()
	return r.anchorer.Tick(context.Background())
}

func (r *rig) state(t *testing.T, epoch uint64) notary.EpochState {
	t.Helper()

	record, ok := r.store.Records[epoch]
	require.True(t, ok, "epoch %d was never built", epoch)

	return record.State
}

func TestAnEpochIsBuiltSentMinedAndFinalized(t *testing.T) {
	r := newRig(t)

	require.NoError(t, r.tick(t))
	require.Equal(t, notary.StateSubmitted, r.state(t, 0))
	require.Len(t, r.chain.Broadcasts, 1)
	require.Equal(t, []string{"IDR"}, r.source.Asked)
	require.Len(t, r.store.Stored[0], 12)
	require.True(t, r.observer.Held)

	r.chain.Mine()
	require.NoError(t, r.tick(t))
	require.Equal(t, notary.StateMined, r.state(t, 0))

	r.chain.Finalize()
	require.NoError(t, r.tick(t))
	require.Equal(t, notary.StateFinalized, r.state(t, 0))
	attempts := r.store.AttemptsOf(0)
	require.Len(t, attempts, 1)
	require.Equal(t, notary.AttemptIncluded, attempts[0].Outcome)
	require.NotNil(t, attempts[0].Receipt)
	require.Equal(t, uint64(0), r.observer.Last.Header.Number)
	require.True(t, r.observer.RunwayKnown)
	require.Positive(t, r.observer.MeanCost)
	require.Positive(t, r.observer.Runway)

	r.advance(10 * time.Minute)
	require.NoError(t, r.tick(t))
	require.Len(t, r.source.Asked, 1, "a new epoch before the interval")

	r.advance(time.Hour)
	require.NoError(t, r.tick(t))
	require.Equal(t, notary.StateSubmitted, r.state(t, 1))
	require.Equal(t, r.store.Records[0].Digest, r.store.Records[1].Header.PreviousDigest)
}

func TestAWriterWithoutTheLeaseDoesNothing(t *testing.T) {
	r := newRig(t)
	_, held, err := r.store.Lease(context.Background(), "notaryd-2", time.Minute)
	require.NoError(t, err)
	require.True(t, held)

	require.NoError(t, r.tick(t))

	require.False(t, r.observer.Held)
	require.Empty(t, r.source.Asked)
	require.Empty(t, r.store.Records)
}

func TestAStuckAnchorIsOutbidAtTheSameNonce(t *testing.T) {
	r := newRig(t)
	require.NoError(t, r.tick(t))
	first := r.store.AttemptsOf(0)[0]

	r.advance(2 * time.Minute)
	require.NoError(t, r.tick(t))
	require.Len(t, r.store.AttemptsOf(0), 1, "replaced before five minutes")

	r.advance(4 * time.Minute)
	require.NoError(t, r.tick(t))
	attempts := r.store.AttemptsOf(0)
	require.Len(t, attempts, 2)
	second := attempts[1]
	require.Equal(t, first.Anchor.Nonce, second.Anchor.Nonce)
	require.GreaterOrEqual(t, second.Anchor.Fees.MaxFeePerGas*10, first.Anchor.Fees.MaxFeePerGas*11)

	r.chain.Mine()
	r.chain.Finalize()
	require.NoError(t, r.tick(t))
	require.Equal(t, notary.StateFinalized, r.state(t, 0))
	attempts = r.store.AttemptsOf(0)
	require.Equal(t, notary.AttemptSuperseded, attempts[0].Outcome)
	require.Equal(t, notary.AttemptIncluded, attempts[1].Outcome)
}

func TestTheFeeCeilingIsNeverCrossed(t *testing.T) {
	r := newRig(t)
	r.chain.BaseFee = 20_000

	err := r.tick(t)
	require.Equal(t, notary.KindAnchorDeferred, notary.KindOf(err))
	require.Equal(t, notary.StateBuilt, r.state(t, 0))
	require.Empty(t, r.chain.Broadcasts)
	require.Equal(t, 1, r.observer.CeilingHits)

	r.chain.BaseFee = 4_000
	require.NoError(t, r.tick(t))
	require.Equal(t, notary.StateSubmitted, r.state(t, 0))

	for range 3 {
		r.advance(6 * time.Minute)
		require.NoError(t, r.tick(t))
	}

	attempts := r.store.AttemptsOf(0)
	require.Len(t, attempts, 2)
	require.Equal(t, uint64(8_100), attempts[0].Anchor.Fees.MaxFeePerGas)
	require.Equal(t, uint64(9_113), attempts[1].Anchor.Fees.MaxFeePerGas)

	for _, attempt := range attempts {
		require.LessOrEqual(t, attempt.Anchor.Fees.MaxFeePerGas, r.policy.FeeCeiling)
	}

	require.Equal(t, 3, r.observer.CeilingHits)
	last := r.store.AttemptsOf(0)[len(r.store.AttemptsOf(0))-1]
	require.Equal(t, last.Anchor.Hash, r.chain.Broadcasts[len(r.chain.Broadcasts)-1])
}

func TestAReorgSendsTheEpochBackAndItLandsAgain(t *testing.T) {
	r := newRig(t)
	require.NoError(t, r.tick(t))
	r.chain.Mine()
	require.NoError(t, r.tick(t))
	require.Equal(t, notary.StateMined, r.state(t, 0))

	r.chain.Reorg()
	require.NoError(t, r.tick(t))
	require.Equal(t, notary.StateSubmitted, r.state(t, 0))
	require.Len(t, r.chain.Broadcasts, 2)
	require.Equal(t, r.chain.Broadcasts[0], r.chain.Broadcasts[1])

	r.chain.Mine()
	r.chain.Finalize()
	require.NoError(t, r.tick(t))
	require.Equal(t, notary.StateFinalized, r.state(t, 0))
}

func TestAnAnchorTheNodeForgotIsSentAgainUnchanged(t *testing.T) {
	r := newRig(t)
	require.NoError(t, r.tick(t))
	r.chain.Forget()

	r.advance(time.Minute)
	require.NoError(t, r.tick(t))

	require.Len(t, r.store.AttemptsOf(0), 1)
	require.Len(t, r.chain.Broadcasts, 2)
	require.Equal(t, r.chain.Broadcasts[0], r.chain.Broadcasts[1])
}

func TestAConflictingEpochOnChainAbandonsOursAndStopsTheDaemon(t *testing.T) {
	r := newRig(t)
	r.chain.Plant(notary.EpochHeader{Number: 0, Root: notary.Digest{0xbe, 0xef}, AsOf: notarytest.Unix(t, 1759999000)})

	err := r.tick(t)
	require.Equal(t, notary.KindLedgerDiverged, notary.KindOf(err))
	require.Empty(t, r.store.Records)
	require.Empty(t, r.source.Asked)
}

func TestAnEpochOverwrittenWhileInFlightIsAbandoned(t *testing.T) {
	r := newRig(t)
	require.NoError(t, r.tick(t))
	r.chain.Forget()
	r.chain.Plant(notary.EpochHeader{Number: 0, Root: notary.Digest{0xbe, 0xef}, AsOf: notarytest.Unix(t, 1759999000)})

	err := r.tick(t)
	require.Equal(t, notary.KindAnchorRefused, notary.KindOf(err))
	require.Equal(t, notary.StateAbandoned, r.state(t, 0))
	require.Contains(t, r.store.Records[0].AbandonedReason, "holds epoch 0")
	require.Equal(t, int64(1), r.observer.AbandonedSeen)
}

func TestARevertedAnchorIsPreparedAgainAndAPermanentReasonAbandonsIt(t *testing.T) {
	r := newRig(t)
	require.NoError(t, r.tick(t))
	r.chain.Sealed = true
	r.chain.Mine()

	err := r.tick(t)

	require.Equal(t, notary.KindAnchorRefused, notary.KindOf(err))
	require.Equal(t, notary.StateAbandoned, r.state(t, 0))
	require.Equal(t, notary.AttemptReverted, r.store.AttemptsOf(0)[0].Outcome)
	require.Contains(t, r.store.Records[0].AbandonedReason, "LedgerIsSealed")
}

func TestANonceSomeoneElseUsedStopsTheDaemon(t *testing.T) {
	r := newRig(t)
	require.NoError(t, r.tick(t))
	r.chain.Foreign()

	err := r.tick(t)

	require.Equal(t, notary.KindLedgerDiverged, notary.KindOf(err))
	require.Len(t, r.store.AttemptsOf(0), 1)
	require.Equal(t, notary.StateSubmitted, r.state(t, 0))
}

func TestNothingIsBuiltWhenTheSnapshotCannotBeUsed(t *testing.T) {
	r := newRig(t)
	unavailable := notary.Errf(notary.KindSourceUnavailable, "test", errors.New("connection refused"), "")
	r.source.Next = func() (notary.Snapshot, error) { return notary.Snapshot{}, unavailable }

	err := r.tick(t)
	require.Equal(t, notary.KindSourceUnavailable, notary.KindOf(err))
	require.Empty(t, r.store.Records)

	r.source.Next = func() (notary.Snapshot, error) {
		return notarytest.Snapshot(t, 1759990000, nil), nil
	}
	r.chain.Plant(notary.EpochHeader{Number: 0, AsOf: notarytest.Unix(t, 1759995000)})
	r.store.Records[0] = notary.EpochRecord{
		Ledger: ledgerID,
		Header: notary.EpochHeader{Number: 0, AsOf: notarytest.Unix(t, 1759995000)},
		Digest: keccak.NewScheme().
			EpochDigest(notary.EpochHeader{Number: 0, AsOf: notarytest.Unix(t, 1759995000)}),
		State: notary.StateFinalized,
	}
	r.advance(2 * time.Hour)
	err = r.tick(t)
	require.Equal(t, notary.KindAnchorDeferred, notary.KindOf(err))
	require.Len(t, r.store.Records, 1)
}

func TestTheDaemonRefusesALedgerItDidNotBuild(t *testing.T) {
	r := newRig(t)
	r.store.Records[0] = notary.EpochRecord{
		Ledger: notary.LedgerIdentity{ChainID: 1, Address: notary.Address{9}},
		Header: notary.EpochHeader{Number: 0, AsOf: notarytest.Unix(t, 1759000000)},
		State:  notary.StateFinalized,
	}

	err := r.tick(t)
	require.Equal(t, notary.KindLedgerDiverged, notary.KindOf(err))

	r = newRig(t)
	r.chain.Sealed = true
	err = r.tick(t)
	require.Equal(t, notary.KindAnchorRefused, notary.KindOf(err))
	require.Empty(t, r.store.Records)
}

func TestAnEpochAlreadyOnChainIsWatchedNotResent(t *testing.T) {
	r := newRig(t)
	require.NoError(t, r.tick(t))
	header := r.store.Records[0].Header
	r.store.Records[0] = notary.EpochRecord{
		Ledger:     ledgerID,
		Header:     header,
		Digest:     r.store.Records[0].Digest,
		Currency:   "IDR",
		State:      notary.StateBuilt,
		SnapshotID: r.store.Records[0].SnapshotID,
	}
	r.chain.Forget()
	r.chain.Plant(header)
	r.watch.Err = notary.Errf(notary.KindApprovalRefused, "test", nil, "the ledger has moved past epoch 0")

	require.NoError(t, r.tick(t))
	require.Equal(t, notary.StateSubmitted, r.state(t, 0))
	require.Len(t, r.chain.Broadcasts, 1)
	require.Len(t, r.watch.Proposals, 1)

	require.NoError(t, r.tick(t))
	require.Equal(t, notary.StateMined, r.state(t, 0))
}
