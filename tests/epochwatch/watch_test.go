package epochwatch_test

import (
	"context"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/WildanFrananda/kinetix-notary-service/internal/epochwatch"
	"github.com/WildanFrananda/kinetix-notary-service/internal/mesh"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	notaryv1 "github.com/WildanFrananda/kinetix-notary-service/internal/protogen/notary/v1"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

type recordingWatch struct {
	mu        sync.Mutex
	proposals []notary.EpochProposal
	answer    notary.Approval
	err       error
}

func (w *recordingWatch) Approve(_ context.Context, proposal notary.EpochProposal) (notary.Approval, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.proposals = append(w.proposals, proposal)

	return w.answer, w.err
}

type rig struct {
	pki     *notarytest.PKI
	watch   *recordingWatch
	address string
}

func newRig(t *testing.T, serverName string) *rig {
	t.Helper()

	pki := notarytest.NewPKI(t)
	watch := &recordingWatch{
		answer: notary.Approval{Member: notary.Address{0x77}, Signature: notary.Signature{0x01, 64: 27}},
	}
	creds, err := mesh.ServerCredentials(
		mesh.Identity{PKIDir: pki.Issue(t, "kinetix.local", serverName), TrustDomains: []string{"kinetix.local"}},
		[]string{"notary"},
	)
	require.NoError(t, err)

	server := grpc.NewServer(grpc.Creds(creds))
	epochwatch.Register(server, watch)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	go func() { _ = server.Serve(listener) }()

	t.Cleanup(server.Stop)

	return &rig{pki: pki, watch: watch, address: listener.Addr().String()}
}

func (r *rig) client(t *testing.T, caller string) *epochwatch.Client {
	t.Helper()

	client, err := epochwatch.Dial(
		epochwatch.Settings{
			Endpoint: r.address,
			Identity: mesh.Identity{
				PKIDir:       r.pki.Issue(t, "kinetix.local", caller),
				TrustDomains: []string{"kinetix.local"},
			},
			Timeout: 5 * time.Second,
		},
		"IDR",
	)
	require.NoError(t, err)
	t.Cleanup(func() { _ = client.Close() })

	return client
}

func proposal(t *testing.T, count int) notary.EpochProposal {
	t.Helper()

	accounts := make([]notary.ProposedAccount, count)
	for i, liability := range notarytest.Liabilities(t, count) {
		accounts[i] = notary.ProposedAccount{
			Liability: liability,
			Secret:    notary.AccountSecret{byte(i), 0x5e},
			Salt:      notary.Salt{byte(i), 0x5a},
		}
	}

	return notary.EpochProposal{
		SnapshotID: notarytest.SnapshotID(t, 21),
		Header: notary.EpochHeader{
			Number:         41,
			PreviousDigest: notary.Digest{0x0a},
			Root:           notary.Digest{0x0b},
			Total:          123_456,
			LeafCount:      uint64(count),
			AsOf:           notarytest.Unix(t, 1760000000),
		},
		Accounts: accounts,
	}
}

func TestAProposalCrossesTheWireWholeAndTheApprovalComesBack(t *testing.T) {
	r := newRig(t, epochwatch.Service)
	sent := proposal(t, 1201)

	approval, err := r.client(t, "notary").Approve(context.Background(), sent)

	require.NoError(t, err)
	require.Equal(t, r.watch.answer, approval)
	require.Len(t, r.watch.proposals, 1)
	require.Equal(t, sent, r.watch.proposals[0])
}

func TestASchemeTwoProposalCarriesItsStatementAndBlindingButNoProof(t *testing.T) {
	r := newRig(t, epochwatch.Service)
	sent := proposal(t, 7)
	sent.Header.Total = 0
	sent.Confidential = &notary.ConfidentialPlan{
		Scheme:          2,
		TreeRoot:        notary.Digest{0x0c},
		TotalCommitment: notary.Digest{0x0d},
		Blinding:        notary.Blinding{0x0e},
		Depth:           8,
	}

	_, err := r.client(t, "notary").Approve(context.Background(), sent)

	require.NoError(t, err)
	require.Equal(t, sent, r.watch.proposals[0])
}

func TestASchemeTwoProposalWithAZeroBlindingIsRefusedBeforeTheWatchSeesIt(t *testing.T) {
	r := newRig(t, epochwatch.Service)
	sent := proposal(t, 2)
	sent.Confidential = &notary.ConfidentialPlan{Scheme: 2, Depth: 8}

	_, err := r.client(t, "notary").Approve(context.Background(), sent)

	require.Equal(t, notary.KindApprovalRefused, notary.KindOf(err))
	require.Empty(t, r.watch.proposals)
}

func TestAnEmptyEpochIsAProposalWithoutAccounts(t *testing.T) {
	r := newRig(t, epochwatch.Service)
	sent := proposal(t, 0)

	_, err := r.client(t, "notary").Approve(context.Background(), sent)

	require.NoError(t, err)
	require.Empty(t, r.watch.proposals[0].Accounts)
}

func TestTheWatchsAnswerDecidesBetweenAbandoningAndRetrying(t *testing.T) {
	for kind, expected := range map[notary.Kind]notary.Kind{
		notary.KindApprovalRefused:   notary.KindApprovalRefused,
		notary.KindSnapshotUnknown:   notary.KindApprovalRefused,
		notary.KindDisagreement:      notary.KindApprovalRefused,
		notary.KindSourceUnavailable: notary.KindSourceUnavailable,
		notary.KindLedgerUnavailable: notary.KindSourceUnavailable,
		notary.KindUnknown:           notary.KindSourceUnavailable,
	} {
		t.Run(kind.String(), func(t *testing.T) {
			r := newRig(t, epochwatch.Service)
			r.watch.err = notary.Errf(kind, "test", nil, "scripted %s", kind)

			_, err := r.client(t, "notary").Approve(context.Background(), proposal(t, 3))

			require.Equal(t, expected, notary.KindOf(err))
		})
	}
}

func TestOnlyTheNotaryMayAskAndOnlyTheWatchMayAnswer(t *testing.T) {
	r := newRig(t, epochwatch.Service)

	_, err := r.client(t, "catalog").Approve(context.Background(), proposal(t, 3))

	require.Equal(t, notary.KindSourceUnavailable, notary.KindOf(err))
	require.Empty(t, r.watch.proposals)

	impostor := newRig(t, "payment")

	_, err = impostor.client(t, "notary").Approve(context.Background(), proposal(t, 3))

	require.Equal(t, notary.KindSourceUnavailable, notary.KindOf(err))
	require.Empty(t, impostor.watch.proposals)
}

func TestAMalformedProposalIsTheCallersMistake(t *testing.T) {
	r := newRig(t, epochwatch.Service)
	creds, err := mesh.ClientCredentials(
		mesh.Identity{PKIDir: r.pki.Issue(t, "kinetix.local", "notary"), TrustDomains: []string{"kinetix.local"}},
		epochwatch.Service,
	)
	require.NoError(t, err)

	conn, err := grpc.NewClient(r.address, grpc.WithTransportCredentials(creds))
	require.NoError(t, err)
	t.Cleanup(func() { _ = conn.Close() })

	stub := notaryv1.NewEpochWatchServiceClient(conn)

	cases := map[string][]*notaryv1.ApproveEpochRequest{
		"no proposal": {},
		"accounts first": {
			{Part: &notaryv1.ApproveEpochRequest_Accounts{Accounts: &notaryv1.ProposedAccountBatch{}}},
		},
		"a short root": {
			{Part: &notaryv1.ApproveEpochRequest_Proposal{Proposal: &notaryv1.EpochProposal{
				SnapshotId:     string(notarytest.SnapshotID(t, 1)),
				PreviousDigest: make([]byte, 32),
				Root:           make([]byte, 31),
			}}},
		},
		"no snapshot": {
			{Part: &notaryv1.ApproveEpochRequest_Proposal{Proposal: &notaryv1.EpochProposal{
				PreviousDigest: make([]byte, 32),
				Root:           make([]byte, 32),
			}}},
		},
	}
	for name, parts := range cases {
		t.Run(name, func(t *testing.T) {
			stream, err := stub.ApproveEpoch(context.Background())
			require.NoError(t, err)

			for _, part := range parts {
				require.NoError(t, stream.Send(part))
			}

			_, err = stream.CloseAndRecv()

			require.Equal(t, codes.InvalidArgument, status.Code(err), err)
		})
	}

	require.Empty(t, r.watch.proposals)
}
