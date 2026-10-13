package config_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/config"
)

const anchorerKey = "0x4c0883a69102937d6231471b5dbb6204fe5129617082792ae468d01a3f362318"

func daemonValues() map[string]string {
	return map[string]string{
		config.RPCURLKey:         "https://sepolia.base.org",
		config.ChainIDKey:        "84532",
		config.LedgerAddressKey:  "0x4e869a55031de3b0d030aebecc66e523686761a2",
		config.RPCTimeoutKey:     "10s",
		config.DatabaseURLKey:    "postgres://notary@kinetix-notary-db:5432/notary",
		config.AnchorerKeyKey:    anchorerKey,
		config.PaymentAddressKey: "kinetix-payment-service:50056",
		config.PKIDirKey:         "/pki",
		config.CurrencyKey:       "IDR",
		config.MaxFeeKey:         "50000000000",
		config.HTTPPortKey:       "8010",
		config.JWKSURLKey:        "http://kinetix-identity-service:5000/.well-known/jwks.json",
		config.JWTIssuerKey:      "https://identity.kinetix.local",
		config.JWTAudienceKey:    "kinetix",
	}
}

func from(values map[string]string) config.Lookup {
	return func(key string) (string, bool) {
		value, ok := values[key]
		return value, ok
	}
}

func TestADaemonConfigurationLoadsWithTheEstatesDefaults(t *testing.T) {
	cfg, err := config.LoadDaemon(from(daemonValues()))

	require.NoError(t, err)
	require.Equal(t, uint64(84532), cfg.Chain.Ledger.ChainID)
	require.Equal(t, anchorerKey, cfg.AnchorerKey.Reveal())
	require.Equal(t, []string{"kinetix.local"}, cfg.TrustDomains)
	require.Equal(t, uint64(50_000_000_000), cfg.MaxFeePerGasWei)
	require.Equal(t, 8010, cfg.HTTPPort)
	require.Equal(t, time.Hour, cfg.EpochInterval)
	require.Equal(t, 15*time.Second, cfg.Tick)
	require.Equal(t, 3*time.Minute, cfg.SnapshotTimeout)
}

func TestSecretsNeverPrint(t *testing.T) {
	cfg, err := config.LoadDaemon(from(daemonValues()))
	require.NoError(t, err)

	for _, rendered := range []string{fmt.Sprint(cfg), fmt.Sprintf("%+v", cfg), fmt.Sprintf("%#v", cfg), cfg.AnchorerKey.String()} {
		require.NotContains(t, rendered, anchorerKey[2:])
		require.NotContains(t, rendered, "kinetix-notary-db")
	}
}

func TestEveryMissingOrBrokenSettingIsNamedAtOnce(t *testing.T) {
	values := daemonValues()
	delete(values, config.AnchorerKeyKey)
	delete(values, config.RPCURLKey)
	values[config.DatabaseURLKey] = "mysql://x"
	values[config.PaymentAddressKey] = "no-port"
	values[config.CurrencyKey] = "idr"
	values[config.MaxFeeKey] = "0"
	values[config.HTTPPortKey] = "70000"
	values[config.EpochIntervalKey] = "10s"
	values[config.TickKey] = "1h"
	values[config.JWKSURLKey] = "ftp://identity/keys"
	delete(values, config.JWTIssuerKey)
	delete(values, config.JWTAudienceKey)

	_, err := config.LoadDaemon(from(values))

	require.Error(t, err)

	for _, key := range []string{
		config.AnchorerKeyKey,
		config.RPCURLKey,
		config.DatabaseURLKey,
		config.PaymentAddressKey,
		config.CurrencyKey,
		config.MaxFeeKey,
		config.HTTPPortKey,
		config.EpochIntervalKey,
		config.TickKey,
		config.JWKSURLKey,
		config.JWTIssuerKey,
		config.JWTAudienceKey,
	} {
		require.ErrorContains(t, err, key)
	}
}

func TestATrustDomainListAllowsACutover(t *testing.T) {
	values := daemonValues()
	values[config.TrustDomainKey] = " kinetix.local , prod.kinetix ,"

	cfg, err := config.LoadDaemon(from(values))

	require.NoError(t, err)
	require.Equal(t, []string{"kinetix.local", "prod.kinetix"}, cfg.TrustDomains)
}

func TestTheMigratorNeedsOnlyTheDatabase(t *testing.T) {
	url, err := config.LoadMigrator(from(map[string]string{config.DatabaseURLKey: "postgres://h/db"}))
	require.NoError(t, err)
	require.Equal(t, "postgres://h/db", url.Reveal())

	_, err = config.LoadMigrator(from(map[string]string{}))
	require.ErrorContains(t, err, config.DatabaseURLKey)
}
