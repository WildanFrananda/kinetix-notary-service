package wiring_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
	"go.uber.org/dig"

	"github.com/WildanFrananda/kinetix-notary-service/internal/config"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/sponsoring"
	"github.com/WildanFrananda/kinetix-notary-service/internal/verifying"
	"github.com/WildanFrananda/kinetix-notary-service/internal/wiring"
)

const anvilPaymaster = "0xcf7ed3acca5a467e9e704c703e8d87f634fb0fc9"

func holderFiles(t *testing.T) (string, string, notary.Address) {
	t.Helper()

	key, err := crypto.GenerateKey()
	require.NoError(t, err)

	encrypted, err := keystore.EncryptKey(
		&keystore.Key{Address: crypto.PubkeyToAddress(key.PublicKey), PrivateKey: key},
		"wiring-password",
		keystore.LightScryptN,
		keystore.LightScryptP,
	)
	require.NoError(t, err)

	dir := t.TempDir()
	keystorePath := filepath.Join(dir, "holder.json")
	passwordPath := filepath.Join(dir, "holder.password")

	require.NoError(t, os.WriteFile(keystorePath, encrypted, 0o600))
	require.NoError(t, os.WriteFile(passwordPath, []byte("wiring-password\n"), 0o600))

	return keystorePath, passwordPath, notary.Address(crypto.PubkeyToAddress(key.PublicKey))
}

func sponsorValues(keystorePath, passwordPath string) map[string]string {
	return map[string]string{
		config.RPCURLKey:           "http://127.0.0.1:1",
		config.ChainIDKey:          "31337",
		config.DirectoryAddressKey: anvilDirectory,
		config.RPCTimeoutKey:       "2s",
		config.PaymasterAddressKey: anvilPaymaster,
		config.BundlerURLKey:       "http://127.0.0.1:2",
		config.HolderKeystoreKey:   keystorePath,
		config.HolderUnlockFileKey: passwordPath,
	}
}

func TestTheGraphResolvesEveryDependencyOfTheSponsor(t *testing.T) {
	keystorePath, passwordPath, holder := holderFiles(t)
	container, err := wiring.Sponsor(lookup(sponsorValues(keystorePath, passwordPath)))
	require.NoError(t, err)

	err = container.Invoke(
		func(
			client *ethclient.Client,
			scheme notary.Scheme,
			statements notary.StatementScheme,
			directory notary.Directory,
			verifier *verifying.DirectoryVerifier,
			sponsor *sponsoring.Sponsor,
		) {
			defer client.Close()

			require.NotNil(t, statements)
			require.NotNil(t, directory)
			require.NotNil(t, verifier)
			require.Equal(t, uint16(1), scheme.Version())
			require.Equal(t, holder, sponsor.Holder())
		},
	)

	require.NoError(t, err)
}

func TestTheSponsorGraphNamesEveryMissingSetting(t *testing.T) {
	container, err := wiring.Sponsor(lookup(map[string]string{}))
	require.NoError(t, err)

	err = container.Invoke(func(*sponsoring.Sponsor) {})

	require.Error(t, err)

	cause := dig.RootCause(err)
	for _, key := range []string{
		config.RPCURLKey,
		config.DirectoryAddressKey,
		config.PaymasterAddressKey,
		config.BundlerURLKey,
		config.HolderKeystoreKey,
		config.HolderUnlockFileKey,
	} {
		require.ErrorContains(t, cause, key)
	}
}

func TestAKeystoreThatDoesNotOpenIsRefusedAsMalformed(t *testing.T) {
	keystorePath, _, _ := holderFiles(t)
	wrong := filepath.Join(t.TempDir(), "wrong.password")
	require.NoError(t, os.WriteFile(wrong, []byte("not-the-password\n"), 0o600))

	container, err := wiring.Sponsor(lookup(sponsorValues(keystorePath, wrong)))
	require.NoError(t, err)

	err = container.Invoke(func(*sponsoring.Sponsor) {})

	require.Equal(t, notary.KindInvalidArgument, notary.KindOf(dig.RootCause(err)), "%v", err)
	require.NotContains(t, err.Error(), "not-the-password")
}
