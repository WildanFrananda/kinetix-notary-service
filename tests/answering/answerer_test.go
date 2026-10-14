package answering_test

import (
	"context"
	"io"
	"log/slog"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/answering"
	"github.com/WildanFrananda/kinetix-notary-service/internal/keccak"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

const epoch = 4

var identity = notary.CourtIdentity{ChainID: 84532, Address: notary.Address{0xd0}}

type observed struct {
	silences     uint64
	open         int
	unanswerable int
	balance      *big.Int
	served       int
	failed       int
}

func (o *observed) Silences(total uint64)        { o.silences = total }
func (o *observed) OpenRequests(count int)       { o.open = count }
func (o *observed) Unanswerable(count int)       { o.unanswerable = count }
func (o *observed) AnswererBalance(wei *big.Int) { o.balance = wei }

func (o *observed) Pass(served bool) {
	if served {
		o.served++
		return
	}

	o.failed++
}

type rig struct {
	court         *notarytest.WithholdingCourt
	chain         *notarytest.AnswerChain
	source        *notarytest.ProofSource
	registrations *notarytest.Registrations
	observer      *observed
	record        notary.EpochRecord
	now           time.Time
	policy        answering.Policy
}

func newRig(t *testing.T) *rig {
	t.Helper()

	scheme := keccak.NewScheme()
	snapshot := notarytest.Snapshot(t, 1760000000, notarytest.Liabilities(t, 15))
	plan, err := notary.PlanEpoch(
		scheme,
		&notarytest.Salts{Seed: 1},
		&notarytest.Secrets{Seed: 1},
		nil,
		epoch,
		notary.Digest{},
		snapshot,
	)
	require.NoError(t, err)

	record := notary.EpochRecord{
		Ledger:   notarytest.LedgerIdentity,
		Header:   plan.Header,
		Digest:   plan.Digest,
		Currency: "IDR",
		State:    notary.StateFinalized,
	}

	court := notarytest.NewWithholdingCourt(identity, scheme)
	court.MaxRange = 201

	return &rig{
		court: court,
		chain: &notarytest.AnswerChain{Address: notary.Address{0xa5}, BalanceValue: big.NewInt(42)},
		source: &notarytest.ProofSource{
			Records: map[uint64]notary.EpochRecord{epoch: record},
			Stored:  map[uint64][]notary.Account{epoch: plan.Accounts},
		},
		registrations: notarytest.NewRegistrations(),
		observer:      &observed{},
		record:        record,
		now:           time.Unix(1760003600, 0).UTC(),
		policy:        answering.Policy{FromBlock: 0, Lookback: 100_000, Span: 200},
	}
}

func (r *rig) answerer() *answering.Answerer {
	return answering.NewAnswerer(
		r.court,
		r.chain,
		r.source,
		r.registrations,
		keccak.NewScheme(),
		r.observer,
		func() time.Time { return r.now },
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		r.policy,
	)
}

func (r *rig) request(t *testing.T, account notary.Digest, tag byte, block uint64) notary.WithholdingRequest {
	t.Helper()

	issued := notary.IssuedRegistration{
		Registration: notary.Registration{Holder: notary.Address{tag}, Account: account, FromEpoch: epoch},
		Kind:         notary.AccountCustomerBalance,
		Digest:       notary.Digest{0xdd, tag},
	}
	require.NoError(
		t,
		r.registrations.SaveRegistrations(
			context.Background(),
			notarytest.Holder(t, "principal-0001"),
			[]notary.IssuedRegistration{issued},
		),
	)

	request := notary.WithholdingRequest{
		Registration: issued.Digest,
		Account:      account,
		Epoch:        epoch,
		Holder:       issued.Registration.Holder,
		DueAt:        notarytest.Unix(t, uint64(r.now.Add(time.Hour).Unix())),
	}
	r.court.Logged[block] = append(r.court.Logged[block], request)
	r.court.States[notarytest.RequestKey{Registration: issued.Digest, Epoch: epoch}] = notarytest.RequestRecord{
		State: notary.RequestOpen,
		DueAt: request.DueAt,
	}
	r.court.Latest = max(r.court.Latest, block)

	return request
}

func (r *rig) account(index int) notary.Account {
	return r.source.Stored[epoch][index]
}

func TestAnOpenRequestIsAnsweredWithAProofTheLedgerWouldAccept(t *testing.T) {
	r := newRig(t)
	account := r.account(3)
	r.request(t, account.ID, 1, 50)
	r.court.Silences = 2

	require.NoError(t, r.answerer().Tick(context.Background()))

	require.Len(t, r.chain.Answers, 1)

	answer := r.chain.Answers[0]
	scheme := keccak.NewScheme()

	require.Equal(t, account.ID, answer.Registration.Account)
	require.Equal(t, uint64(epoch), answer.Epoch)
	require.Equal(t, account.Salt, answer.Salt)
	require.Equal(t, account.Liability.Amount, answer.Amount)
	require.NoError(
		t,
		notary.VerifyInclusion(
			scheme,
			notary.AnchoredEpoch{
				Number:    epoch,
				Root:      r.record.Header.Root,
				Total:     r.record.Header.Total,
				LeafCount: r.record.Header.LeafCount,
			},
			notary.Entry{
				Commitment: scheme.AccountCommitment(answer.Registration.Account, answer.Salt),
				Amount:     answer.Amount,
			},
			answer.Proof,
		),
	)
	require.Equal(t, uint64(2), r.observer.silences)
	require.Equal(t, 0, r.observer.open)
	require.Equal(t, 0, r.observer.unanswerable)
	require.Equal(t, big.NewInt(42), r.observer.balance)
}

func TestARequestNoLongerOpenIsLeftAlone(t *testing.T) {
	for _, state := range []notary.RequestState{notary.RequestAnswered, notary.RequestUnanswered} {
		t.Run(state.String(), func(t *testing.T) {
			r := newRig(t)
			request := r.request(t, r.account(3).ID, 1, 50)
			key := notarytest.RequestKey{Registration: request.Registration, Epoch: epoch}
			r.court.States[key] = notarytest.RequestRecord{State: state, DueAt: request.DueAt}

			require.NoError(t, r.answerer().Tick(context.Background()))

			require.Empty(t, r.chain.Answers)
			require.Equal(t, 0, r.observer.open)
		})
	}
}

func TestTheScanStartsALookbackBehindTheHeadButNeverBeforeTheDeployBlock(t *testing.T) {
	r := newRig(t)
	r.court.Latest = 250_000
	a := r.answerer()

	require.NoError(t, a.Tick(context.Background()))
	require.Equal(t, uint64(150_000), r.court.Scanned[0][0])
	require.Equal(t, uint64(250_000), r.court.Scanned[len(r.court.Scanned)-1][1])

	for _, span := range r.court.Scanned {
		require.LessOrEqual(t, span[1]-span[0]+1, r.policy.Span)
	}

	r.court.Scanned = nil
	r.court.Latest = 250_010

	require.NoError(t, a.Tick(context.Background()))
	require.Equal(t, [][2]uint64{{250_001, 250_010}}, r.court.Scanned)

	r = newRig(t)
	r.court.Latest = 250_000
	r.policy.FromBlock = 240_000
	require.NoError(t, r.answerer().Tick(context.Background()))
	require.Equal(t, uint64(240_000), r.court.Scanned[0][0])
}

func TestAnAccountMissingFromItsEpochIsUnanswerableAndHoldsNothingUp(t *testing.T) {
	r := newRig(t)
	r.request(t, notary.Digest{0x0e}, 1, 50)
	r.request(t, r.account(5).ID, 2, 60)

	require.NoError(t, r.answerer().Tick(context.Background()))

	require.Len(t, r.chain.Answers, 1)
	require.Equal(t, r.account(5).ID, r.chain.Answers[0].Registration.Account)
	require.Equal(t, 1, r.observer.unanswerable)
	require.Equal(t, 1, r.observer.open)
}

func TestARegistrationThisDatabaseNeverIssuedIsUnanswerable(t *testing.T) {
	r := newRig(t)
	request := r.request(t, r.account(3).ID, 1, 50)
	delete(r.registrations.Saved, request.Registration)

	require.NoError(t, r.answerer().Tick(context.Background()))

	require.Empty(t, r.chain.Answers)
	require.Equal(t, 1, r.observer.unanswerable)
}

func TestStoredAccountsThatDoNotRebuildTheRootAreNeverAnswered(t *testing.T) {
	r := newRig(t)
	r.request(t, r.account(3).ID, 1, 50)
	r.source.Stored[epoch][7].Liability.Amount++

	require.NoError(t, r.answerer().Tick(context.Background()))

	require.Empty(t, r.chain.Answers)
	require.Equal(t, 1, r.observer.unanswerable)
}

func TestARefusalBecomesUnanswerableButAnOutageIsRetried(t *testing.T) {
	r := newRig(t)
	r.request(t, r.account(3).ID, 1, 50)
	r.chain.Err = notary.Errf(notary.KindLedgerUnavailable, "test", nil, "")
	a := r.answerer()

	require.Equal(t, notary.KindLedgerUnavailable, notary.KindOf(a.Tick(context.Background())))
	require.Equal(t, 1, r.observer.open)
	require.Equal(t, 1, r.observer.failed)
	require.Equal(t, 0, r.observer.served)

	r.chain.Err = nil

	require.NoError(t, a.Tick(context.Background()))
	require.Len(t, r.chain.Answers, 1)
	require.Equal(t, 1, r.observer.served)

	r = newRig(t)
	r.request(t, r.account(3).ID, 1, 50)
	r.chain.Err = notary.Errf(notary.KindAnswerRefused, "test", nil, "AnswerDoesNotHold")

	require.NoError(t, r.answerer().Tick(context.Background()))
	require.Equal(t, 1, r.observer.unanswerable)
}

func TestARequestPastItsDeadlineIsDropped(t *testing.T) {
	r := newRig(t)
	request := r.request(t, r.account(3).ID, 1, 50)
	r.now = request.DueAt.Time().Add(time.Second)

	require.NoError(t, r.answerer().Tick(context.Background()))

	require.Empty(t, r.chain.Answers)
	require.Equal(t, 0, r.observer.open)
}

func TestAScanTheRPCEndpointRefusesFailsEveryPassAndIsCounted(t *testing.T) {
	r := newRig(t)
	r.request(t, r.account(3).ID, 1, 50)
	r.court.Latest = 5_000
	r.policy.Span = 2_000
	a := r.answerer()

	for range 3 {
		require.Equal(t, notary.KindLedgerUnavailable, notary.KindOf(a.Tick(context.Background())))
	}

	require.Empty(t, r.chain.Answers)
	require.Equal(t, 3, r.observer.failed)
	require.Equal(t, 0, r.observer.served)
}

func TestTheStartUpProbeAsksForOneFullSpanEndingAtTheHead(t *testing.T) {
	r := newRig(t)
	r.court.Latest = 48_089_057

	require.NoError(t, r.answerer().CheckScan(context.Background()))
	require.Equal(t, [][2]uint64{{48_088_858, 48_089_057}}, r.court.Scanned)

	r.court.Scanned = nil
	r.court.Latest = 50

	require.NoError(t, r.answerer().CheckScan(context.Background()))
	require.Equal(t, [][2]uint64{{0, 50}}, r.court.Scanned)
}

func TestTheStartUpProbeRefusesALogRangeTheRPCEndpointRejects(t *testing.T) {
	r := newRig(t)
	r.court.Latest = 48_089_057
	r.policy.Span = 2_000

	err := r.answerer().CheckScan(context.Background())

	require.Equal(t, notary.KindLedgerUnavailable, notary.KindOf(err))
	require.ErrorContains(t, err, "2000 blocks")
	require.ErrorContains(t, err, "at most 201 blocks")
	require.Empty(t, r.court.Scanned)
}
