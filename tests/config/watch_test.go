package config_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/config"
)

const watchKey = "0x7c0883a69102937d6231471b5dbb6204fe5129617082792ae468d01a3f362318"

func watchValues() map[string]string {
	return map[string]string{
		config.RPCURLKey:         "https://sepolia.base.org",
		config.ChainIDKey:        "84532",
		config.LedgerAddressKey:  "0x4e869a55031de3b0d030aebecc66e523686761a2",
		config.CourtAddressKey:   "0x00000000000000000000000000000000c0c0c0c0",
		config.RPCTimeoutKey:     "10s",
		config.CouncilAddressKey: "0x00000000000000000000000000000000c0c0c0c1",
		config.WatchKeyKey:       watchKey,
		config.PaymentAddressKey: "kinetix-payment-service:50056",
		config.PKIDirKey:         "/pki",
		config.WatchGRPCPortKey:  "50070",
		config.HTTPPortKey:       "8011",
	}
}

func TestAWatchConfigurationLoads(t *testing.T) {
	cfg, err := config.LoadWatch(from(watchValues()))

	require.NoError(t, err)
	require.Equal(t, uint64(84532), cfg.Chain.Ledger.ChainID)
	require.Equal(t, "0x00000000000000000000000000000000c0c0c0c1", cfg.Council.Hex())
	require.Equal(t, watchKey, cfg.WatchKey.Reveal())
	require.Equal(t, "kinetix-payment-service:50056", cfg.PaymentAddress)
	require.Equal(t, []string{"kinetix.local"}, cfg.TrustDomains)
	require.Equal(t, 50070, cfg.GRPCPort)
	require.Equal(t, 8011, cfg.HTTPPort)
	require.NotContains(t, fmt.Sprintf("%+v %#v", cfg, cfg), watchKey[2:])
}

func TestEveryMissingOrBrokenWatchSettingIsNamedAtOnce(t *testing.T) {
	values := watchValues()
	delete(values, config.WatchKeyKey)
	delete(values, config.RPCURLKey)
	values[config.CouncilAddressKey] = "council"
	values[config.PaymentAddressKey] = "no-port"
	values[config.WatchGRPCPortKey] = "0"
	values[config.HTTPPortKey] = "70000"

	_, err := config.LoadWatch(from(values))

	for _, key := range []string{
		config.WatchKeyKey,
		config.RPCURLKey,
		config.CouncilAddressKey,
		config.PaymentAddressKey,
		config.WatchGRPCPortKey,
		config.HTTPPortKey,
	} {
		require.ErrorContains(t, err, key)
	}
}

func TestTheWatchsTwoPortsDiffer(t *testing.T) {
	values := watchValues()
	values[config.HTTPPortKey] = "50070"

	_, err := config.LoadWatch(from(values))

	require.ErrorContains(t, err, "must differ")
}
