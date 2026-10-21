package ethereum_test

import (
	"context"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/ethereum"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

func dial(t *testing.T, url string, rate int, counts *rpcCounts) *ethclient.Client {
	t.Helper()

	client, err := ethereum.Dial(ethereum.Settings{RPCURL: url, Rate: rate}, counts)
	require.NoError(t, err)
	t.Cleanup(client.Close)

	return client
}

func TestA429IsWaitedOutAndSentAgain(t *testing.T) {
	endpoint := newRPCEndpoint(t)
	endpoint.Refuse(1)

	counts := &rpcCounts{}
	client := dial(t, endpoint.URL, 0, counts)
	started := time.Now()

	chainID, err := client.ChainID(context.Background())

	require.NoError(t, err)
	require.Equal(t, uint64(84532), chainID.Uint64())
	require.GreaterOrEqual(t, time.Since(started), time.Second)
	require.Equal(t, 2, endpoint.Served())

	requests, throttled, _ := counts.snapshot()
	require.Equal(t, 2, requests)
	require.Equal(t, 1, throttled)
}

func TestAnEndpointThatKeepsRefusingIsReportedAfterThreeRetries(t *testing.T) {
	endpoint := newRPCEndpoint(t)
	endpoint.Refuse(100)

	counts := &rpcCounts{}
	client := dial(t, endpoint.URL, 0, counts)

	_, err := client.ChainID(context.Background())

	require.ErrorContains(t, err, "429")
	require.Equal(t, 4, endpoint.Served())

	_, throttled, _ := counts.snapshot()
	require.Equal(t, 4, throttled)
}

func TestA429IsNotWaitedOutPastTheCallersDeadline(t *testing.T) {
	endpoint := newRPCEndpoint(t)
	endpoint.Refuse(100)

	client := dial(t, endpoint.URL, 0, &rpcCounts{})

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	started := time.Now()

	_, err := client.ChainID(ctx)

	require.ErrorContains(t, err, "429")
	require.Less(t, time.Since(started), 500*time.Millisecond)
	require.Equal(t, 1, endpoint.Served())
}

func TestAnyOtherFailureIsNotRetried(t *testing.T) {
	endpoint := newRPCEndpoint(t)
	endpoint.Fail(1)
	client := dial(t, endpoint.URL, 0, &rpcCounts{})

	_, err := client.ChainID(context.Background())

	require.ErrorContains(t, err, "502")
	require.Equal(t, 1, endpoint.Served())
}

func TestRequestsKeepToTheRate(t *testing.T) {
	endpoint := newRPCEndpoint(t)
	counts := &rpcCounts{}
	client := dial(t, endpoint.URL, 20, counts)
	started := time.Now()

	for range 11 {
		_, err := client.ChainID(context.Background())
		require.NoError(t, err)
	}

	require.GreaterOrEqual(t, time.Since(started), 450*time.Millisecond)

	requests, _, paced := counts.snapshot()
	require.Equal(t, 11, requests)
	require.Greater(t, paced, 400*time.Millisecond)
}

func TestOnlyAnHTTPEndpointCanBeHeldToARate(t *testing.T) {
	_, err := ethereum.Dial(ethereum.Settings{RPCURL: "wss://rpc.example.org", Rate: 10}, nil)

	require.Equal(t, notary.KindInvalidArgument, notary.KindOf(err))
}
