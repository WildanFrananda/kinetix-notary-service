package wiring_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/config"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/verifying"
	"github.com/WildanFrananda/kinetix-notary-service/internal/wiring"
)

func TestTheEvidenceVerifierGraphResolvesFromItsFourSettings(t *testing.T) {
	values := map[string]string{
		config.RPCURLKey:           "http://127.0.0.1:8545",
		config.ChainIDKey:          "84532",
		config.EvidenceRegistryKey: "0x00000000000000000000000000000000e0e0e0e0",
		config.RPCTimeoutKey:       "10s",
	}
	lookup := func(key string) (string, bool) {
		value, ok := values[key]
		return value, ok
	}

	container, err := wiring.EvidenceVerifier(lookup)
	require.NoError(t, err)

	invoked := container.Invoke(func(
		client *ethclient.Client,
		scheme notary.EvidenceScheme,
		verifier *verifying.EvidenceVerifier,
		registry notary.EvidenceRegistry,
	) {
		defer client.Close()

		require.Equal(t, uint16(1), scheme.EvidenceVersion())
		require.NotNil(t, verifier)
		require.NotNil(t, registry)
	})
	require.NoError(t, invoked)

	delete(values, config.EvidenceRegistryKey)

	container, err = wiring.EvidenceVerifier(lookup)
	require.NoError(t, err)
	require.ErrorContains(t, container.Invoke(func(*verifying.EvidenceVerifier) {}), config.EvidenceRegistryKey)
}
