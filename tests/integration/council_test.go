package integration_test

import (
	"context"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/ethereum"
	"github.com/WildanFrananda/kinetix-notary-service/internal/keccak"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

func TestTheCouncilAdapterReadsTheDeployedCouncil(t *testing.T) {
	c := connect(t)
	ctx := context.Background()
	identity := c.identity(c.deploy(t))
	_, watch := memberKey(t)
	address := c.deployCouncil(t, common.Address(identity.Address), notary.Address(c.from), watch)
	adapter := c.councilAt(identity, address)
	expected := notary.CouncilIdentity{ChainID: identity.ChainID, Address: notary.Address(address)}

	answered, err := adapter.Identity(ctx)
	require.NoError(t, err)
	require.Equal(t, expected, answered)

	for member, wanted := range map[notary.Address]bool{
		notary.Address(c.from): true,
		watch:                  true,
		{0x01}:                 false,
	} {
		is, err := adapter.IsMember(ctx, member)
		require.NoError(t, err)
		require.Equal(t, wanted, is, member.Hex())
	}

	threshold, err := adapter.Threshold(ctx)
	require.NoError(t, err)
	require.Equal(t, uint64(2), threshold)

	sample := notary.Digest{0x42}
	digest, err := adapter.ApprovalDigest(ctx, 7, sample)
	require.NoError(t, err)
	require.Equal(t, keccak.NewScheme().ApprovalDigest(expected, identity.Address, 7, sample), digest)

	anchors, err := adapter.AnchorsLedger(ctx)
	require.NoError(t, err)
	require.False(t, anchors)

	c.handOver(t, common.Address(identity.Address), address)

	anchors, err = adapter.AnchorsLedger(ctx)
	require.NoError(t, err)
	require.True(t, anchors)
}

func TestAnythingButTheCouncilOfThisLedgerIsTheWrongCouncil(t *testing.T) {
	c := connect(t)
	ctx := context.Background()
	identity := c.identity(c.deploy(t))
	other := c.identity(c.deploy(t))
	_, watch := memberKey(t)
	elsewhere := c.deployCouncil(t, common.Address(other.Address), notary.Address(c.from), watch)

	for name, address := range map[string]common.Address{
		"nothing there":            {0x0e, 0x0e},
		"the ledger itself":        common.Address(identity.Address),
		"another ledger's council": elsewhere,
	} {
		t.Run(name, func(t *testing.T) {
			_, err := c.councilAt(identity, address).Identity(ctx)
			require.Equal(t, notary.KindWrongCouncil, notary.KindOf(err), err)
		})
	}

	wrongChain := notary.LedgerIdentity{ChainID: identity.ChainID + 1, Address: identity.Address}
	council := c.deployCouncil(t, common.Address(identity.Address), notary.Address(c.from), watch)
	_, err := c.councilAt(wrongChain, council).Identity(ctx)
	require.Equal(t, notary.KindWrongCouncil, notary.KindOf(err))
}

func TestTheHeadReaderSeesWhatTheWriterSees(t *testing.T) {
	r := newWriterRig(t)
	ctx := context.Background()
	reader := ethereum.NewHeadReader(
		r.chain.client,
		ethereum.Settings{RPCURL: r.chain.url, Ledger: r.ledger, Timeout: 10 * time.Second},
	)

	head, err := r.writer.Head(ctx)
	require.NoError(t, err)

	plan := r.plan(t, head.NextEpoch, head.Head, r.now(t))
	anchor, err := r.prepare(t, plan.Header)
	require.NoError(t, err)
	require.NoError(t, r.writer.Broadcast(ctx, anchor))
	r.mined(t, anchor.Hash)

	written, err := r.writer.Head(ctx)
	require.NoError(t, err)

	read, err := reader.Head(ctx)
	require.NoError(t, err)
	require.Equal(t, written, read)
	require.Equal(t, plan.Digest, read.Head)
}
