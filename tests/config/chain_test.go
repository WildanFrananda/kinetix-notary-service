package config_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/config"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

const mixedCaseLedger = "0x5FbDB2315678afecb367f032d93F642f64180aa3"

func lookup(values map[string]string) config.Lookup {
	return func(key string) (string, bool) {
		value, ok := values[key]
		return value, ok
	}
}

func complete() map[string]string {
	return map[string]string{
		config.RPCURLKey:        "https://rpc.example.org/v2/secret-api-key",
		config.ChainIDKey:       "84532",
		config.LedgerAddressKey: mixedCaseLedger,
		config.CourtAddressKey:  "0x00000000000000000000000000000000C0C0c0c0",
		config.RPCTimeoutKey:    "10s",
		config.RPCRateKey:       "12",
	}
}

func TestACompleteEnvironmentLoads(t *testing.T) {
	cfg, err := config.LoadChain(lookup(complete()))

	require.NoError(t, err)
	address, err := notary.ParseAddress("0x5fbdb2315678afecb367f032d93f642f64180aa3")
	require.NoError(t, err)
	require.Equal(
		t,
		config.Chain{
			RPCURL:     "https://rpc.example.org/v2/secret-api-key",
			Ledger:     notary.LedgerIdentity{ChainID: 84532, Address: address},
			Court:      notary.Address{16: 0xc0, 17: 0xc0, 18: 0xc0, 19: 0xc0},
			RPCTimeout: 10 * time.Second,
			RPCRate:    12,
		},
		cfg,
	)
}

func TestEveryMissingKeyIsNamedAtOnce(t *testing.T) {
	_, err := config.LoadChain(lookup(map[string]string{}))

	require.Error(t, err)

	keys := []string{
		config.RPCURLKey,
		config.ChainIDKey,
		config.LedgerAddressKey,
		config.CourtAddressKey,
		config.RPCTimeoutKey,
		config.RPCRateKey,
	}
	for _, key := range keys {
		require.Contains(t, err.Error(), key+" is not set")
	}
}

func TestABlankValueIsMissing(t *testing.T) {
	values := complete()
	values[config.ChainIDKey] = "   "

	_, err := config.LoadChain(lookup(values))

	require.ErrorContains(t, err, config.ChainIDKey+" is not set")
}

func TestInvalidValuesAreRefused(t *testing.T) {
	cases := map[string]string{
		config.RPCURLKey:        "ftp://rpc.example.org",
		config.ChainIDKey:       "0",
		config.LedgerAddressKey: "0x0000000000000000000000000000000000000000",
		config.CourtAddressKey:  "0xc0",
		config.RPCTimeoutKey:    "10m",
		config.RPCRateKey:       "0",
	}
	for key, value := range cases {
		t.Run(key, func(t *testing.T) {
			values := complete()
			values[key] = value

			_, err := config.LoadChain(lookup(values))

			require.ErrorContains(t, err, key)
		})
	}
}

func TestARateAboveTheCeilingIsRefused(t *testing.T) {
	values := complete()
	values[config.RPCRateKey] = "1001"

	_, err := config.LoadChain(lookup(values))

	require.ErrorContains(t, err, config.RPCRateKey)
}

func TestAWebSocketEndpointIsRefusedBecauseOnlyHTTPIsPaced(t *testing.T) {
	values := complete()
	values[config.RPCURLKey] = "wss://rpc.example.org/v2/secret-api-key"

	_, err := config.LoadChain(lookup(values))

	require.ErrorContains(t, err, config.RPCURLKey+" must be http or https")
	require.NotContains(t, err.Error(), "secret-api-key")
}

func TestTheRPCURLNeverAppearsInAnError(t *testing.T) {
	values := complete()
	values[config.RPCURLKey] = "gopher://rpc.example.org/v2/secret-api-key"
	values[config.ChainIDKey] = "nope"

	_, err := config.LoadChain(lookup(values))

	require.Error(t, err)
	require.NotContains(t, err.Error(), "secret-api-key")
}
