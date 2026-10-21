package ethereum_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/ethereum"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

func proofReader(t *testing.T, endpoint *rpcEndpoint) *ethereum.ProofReader {
	t.Helper()

	settings := ethereum.Settings{
		Ledger:  notary.LedgerIdentity{ChainID: 84532, Address: notary.Address{0xd0}},
		Timeout: 500 * time.Millisecond,
	}

	return ethereum.NewProofReader(dial(t, endpoint.URL, 0, &rpcCounts{}), settings)
}

func TestAThrottledDepthReadIsAnOutageNotAWrongLedger(t *testing.T) {
	endpoint := newRPCEndpoint(t)
	endpoint.Refuse(100)

	_, err := proofReader(t, endpoint).Depth(context.Background())

	require.Equal(t, notary.KindLedgerUnavailable, notary.KindOf(err))
}

func TestALedgerThatRevertsOnDepthIsTheWrongLedger(t *testing.T) {
	endpoint := newRPCEndpoint(t)
	endpoint.Revert()

	_, err := proofReader(t, endpoint).Depth(context.Background())

	require.Equal(t, notary.KindWrongLedger, notary.KindOf(err))
}
