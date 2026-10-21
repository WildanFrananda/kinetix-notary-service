package ethereum_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/ethereum"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

func TestThePaymastersIdentityIsReadOnceAndKept(t *testing.T) {
	endpoint := newRPCEndpoint(t)
	paymaster := ethereum.NewPaymaster(
		dial(t, endpoint.URL, 0, &rpcCounts{}),
		84532,
		notary.Address{0xf0},
		time.Second,
	)

	endpoint.Fail(1)

	_, err := paymaster.Identity(context.Background())
	require.Equal(t, notary.KindLedgerUnavailable, notary.KindOf(err))

	identity, err := paymaster.Identity(context.Background())
	require.NoError(t, err)
	require.Equal(t, notary.Address{19: 0x01}, identity.Court)
	require.Equal(t, big.NewInt(1), identity.Ceiling)

	read := endpoint.Served()

	identity.Ceiling.SetInt64(7)

	again, err := paymaster.Identity(context.Background())
	require.NoError(t, err)
	require.Equal(t, read, endpoint.Served())
	require.Equal(t, big.NewInt(1), again.Ceiling)
}
