package integration_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"github.com/WildanFrananda/kinetix-notary-service/internal/contractgen/evidence"
	"github.com/WildanFrananda/kinetix-notary-service/internal/keccak"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	fulfillmentv1 "github.com/WildanFrananda/kinetix-notary-service/internal/protogen/fulfillment/v1"
	orderv1 "github.com/WildanFrananda/kinetix-notary-service/internal/protogen/order/v1"
	reviewv1 "github.com/WildanFrananda/kinetix-notary-service/internal/protogen/review/v1"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

func (c *chain) deployEvidenceRegistry(t *testing.T, anchorer common.Address) common.Address {
	t.Helper()

	genesis := []evidence.StreamGenesis{{Stream: 1}, {Stream: 2}, {Stream: 3}}
	input := evidence.NewEvidenceRegistry().PackConstructor(c.from, big.NewInt(3600), anchorer, genesis)
	address, tx, err := bind.DeployContract(
		c.transactor(),
		common.FromHex(evidence.EvidenceRegistryMetaData.Bin),
		c.client,
		input,
	)
	require.NoError(t, err)
	_, err = bind.WaitDeployed(context.Background(), c.client, tx.Hash())
	require.NoError(t, err)

	return address
}

func (e *estate) registerFeeds(server *grpc.Server) {
	orderv1.RegisterOrderLifecycleServiceServer(server, e.evidence)
	fulfillmentv1.RegisterDeliveryEvidenceServiceServer(server, e.evidence)
	reviewv1.RegisterReviewEvidenceServiceServer(server, e.evidence)
}

func (c *chain) anchorBatch(t *testing.T, registry common.Address, header notary.EvidenceBatchHeader) {
	t.Helper()

	binding := evidence.NewEvidenceRegistry()
	packed := binding.PackAnchor(evidence.BatchHeader{
		Stream:             uint8(header.Stream),
		Number:             header.Number,
		PreviousDigest:     header.PreviousDigest,
		Root:               header.Root,
		LeafCount:          header.LeafCount,
		ThroughTransaction: header.Through.Transaction,
		ThroughSequence:    header.Through.Sequence,
		ClosedAt:           header.ClosedAt.Seconds(),
	})
	tx, err := bind.Transact(binding.Instance(c.client, registry), c.transactor(), packed)
	require.NoError(t, err)
	receipt, err := bind.WaitMined(context.Background(), c.client, tx.Hash())
	require.NoError(t, err)
	require.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)
}

func (c *chain) evidenceBatch(
	t *testing.T,
	records []notary.EvidenceRecord,
) (notary.RegistryIdentity, notarytest.EvidenceBatch) {
	t.Helper()

	registry := notary.RegistryIdentity{
		ChainID: c.chainID.Uint64(),
		Address: notary.Address(c.deployEvidenceRegistry(t, c.from)),
	}
	head, err := c.client.HeaderByNumber(context.Background(), nil)
	require.NoError(t, err)

	batch := notarytest.BatchOf(t, keccak.NewEvidence(), registry, 0, notary.Digest{}, head.Time-1, records)
	c.anchorBatch(t, common.Address(registry.Address), batch.Header)

	return registry, batch
}
