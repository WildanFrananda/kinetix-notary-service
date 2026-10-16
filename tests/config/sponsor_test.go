package config_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/config"
)

func sponsorValues() map[string]string {
	values := verifierValues()
	values[config.PaymasterAddressKey] = "0x00000000000000000000000000000000fafafafa"
	values[config.BundlerURLKey] = "https://bundler.example.org/v3/84532/secret-bundler-key"
	values[config.HolderKeystoreKey] = "/keys/holder.json"
	values[config.HolderUnlockFileKey] = "/keys/holder.password"

	return values
}

func TestASponsorReadsTheDirectoryThePaymasterTheBundlerAndTheHoldersKeystore(t *testing.T) {
	cfg, err := config.LoadSponsor(lookup(sponsorValues()))

	require.NoError(t, err)
	require.Equal(t, uint64(84532), cfg.Verifier.Directory.ChainID)
	require.Equal(t, "0x00000000000000000000000000000000fafafafa", cfg.Paymaster.Hex())
	require.Equal(t, "https://bundler.example.org/v3/84532/secret-bundler-key", cfg.BundlerURL)
	require.Equal(t, "/keys/holder.json", cfg.HolderKeystore)
	require.Equal(t, "/keys/holder.password", cfg.HolderPassword)
	require.Equal(t, 3*time.Minute, cfg.ReceiptWait)
}

func TestASponsorNamesEveryMissingKeyAtOnce(t *testing.T) {
	_, err := config.LoadSponsor(lookup(map[string]string{}))

	for _, key := range []string{
		config.RPCURLKey,
		config.ChainIDKey,
		config.DirectoryAddressKey,
		config.RPCTimeoutKey,
		config.PaymasterAddressKey,
		config.BundlerURLKey,
		config.HolderKeystoreKey,
		config.HolderUnlockFileKey,
	} {
		require.ErrorContains(t, err, key+" is not set")
	}
}

func TestASponsorRefusesWhatCannotBeUsed(t *testing.T) {
	for key, value := range map[string]string{
		config.BundlerURLKey:       "ftp://bundler.example.org/secret-bundler-key",
		config.PaymasterAddressKey: "0x0000000000000000000000000000000000000000",
		config.ReceiptWaitKey:      "31m",
	} {
		t.Run(key, func(t *testing.T) {
			values := sponsorValues()
			values[key] = value

			_, err := config.LoadSponsor(lookup(values))

			require.ErrorContains(t, err, key)
			require.NotContains(t, err.Error(), "secret-bundler-key")
		})
	}
}

func TestTheReceiptWaitIsTunable(t *testing.T) {
	values := sponsorValues()
	values[config.ReceiptWaitKey] = "90s"

	cfg, err := config.LoadSponsor(lookup(values))

	require.NoError(t, err)
	require.Equal(t, 90*time.Second, cfg.ReceiptWait)
}
