package config_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/config"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

func verifierValues() map[string]string {
	return map[string]string{
		config.RPCURLKey:           "https://rpc.example.org/v2/secret-api-key",
		config.ChainIDKey:          "84532",
		config.DirectoryAddressKey: mixedCaseLedger,
		config.RPCTimeoutKey:       "10s",
	}
}

func TestAVerifierNeedsTheDirectoryAndNothingThatNamesALedger(t *testing.T) {
	cfg, err := config.LoadVerifier(lookup(verifierValues()))

	require.NoError(t, err)
	address, err := notary.ParseAddress("0x5fbdb2315678afecb367f032d93f642f64180aa3")
	require.NoError(t, err)
	require.Equal(
		t,
		config.Verifier{
			RPCURL:     "https://rpc.example.org/v2/secret-api-key",
			Directory:  notary.DirectoryIdentity{ChainID: 84532, Address: address},
			RPCTimeout: 10 * time.Second,
		},
		cfg,
	)
}

func TestAVerifierNamesEveryMissingKeyAtOnce(t *testing.T) {
	_, err := config.LoadVerifier(lookup(map[string]string{}))

	for _, key := range []string{config.RPCURLKey, config.ChainIDKey, config.DirectoryAddressKey, config.RPCTimeoutKey} {
		require.ErrorContains(t, err, key+" is not set")
	}

	require.NotContains(t, err.Error(), config.LedgerAddressKey)
	require.NotContains(t, err.Error(), config.CourtAddressKey)
}

func TestAVerifierRefusesAZeroOrMalformedDirectory(t *testing.T) {
	for _, value := range []string{"0x0000000000000000000000000000000000000000", "0xd1"} {
		values := verifierValues()
		values[config.DirectoryAddressKey] = value

		_, err := config.LoadVerifier(lookup(values))

		require.ErrorContains(t, err, config.DirectoryAddressKey+" must be a non-zero 0x-prefixed 20-byte address")
	}
}
