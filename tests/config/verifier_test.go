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
		config.RPCTimeoutKey:    "10s",
	}
}

func TestACompleteEnvironmentLoads(t *testing.T) {
	cfg, err := config.LoadVerifier(lookup(complete()))

	require.NoError(t, err)
	address, err := notary.ParseAddress("0x5fbdb2315678afecb367f032d93f642f64180aa3")
	require.NoError(t, err)
	require.Equal(t, config.Verifier{
		RPCURL:     "https://rpc.example.org/v2/secret-api-key",
		Ledger:     notary.LedgerIdentity{ChainID: 84532, Address: address},
		RPCTimeout: 10 * time.Second,
	}, cfg)
}

func TestEveryMissingKeyIsNamedAtOnce(t *testing.T) {
	_, err := config.LoadVerifier(lookup(map[string]string{}))

	require.Error(t, err)
	for _, key := range []string{config.RPCURLKey, config.ChainIDKey, config.LedgerAddressKey, config.RPCTimeoutKey} {
		require.Contains(t, err.Error(), key+" is not set")
	}
}

func TestABlankValueIsMissing(t *testing.T) {
	values := complete()
	values[config.ChainIDKey] = "   "

	_, err := config.LoadVerifier(lookup(values))

	require.ErrorContains(t, err, config.ChainIDKey+" is not set")
}

func TestInvalidValuesAreRefused(t *testing.T) {
	cases := map[string]string{
		config.RPCURLKey:        "ftp://rpc.example.org",
		config.ChainIDKey:       "0",
		config.LedgerAddressKey: "0x0000000000000000000000000000000000000000",
		config.RPCTimeoutKey:    "10m",
	}
	for key, value := range cases {
		t.Run(key, func(t *testing.T) {
			values := complete()
			values[key] = value

			_, err := config.LoadVerifier(lookup(values))

			require.ErrorContains(t, err, key)
		})
	}
}

func TestTheRPCURLNeverAppearsInAnError(t *testing.T) {
	values := complete()
	values[config.RPCURLKey] = "gopher://rpc.example.org/v2/secret-api-key"
	values[config.ChainIDKey] = "nope"

	_, err := config.LoadVerifier(lookup(values))

	require.Error(t, err)
	require.NotContains(t, err.Error(), "secret-api-key")
}
