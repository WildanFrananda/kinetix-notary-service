package wiring_test

import (
	"io"
	"log/slog"
	"testing"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
	"go.uber.org/dig"

	"github.com/WildanFrananda/kinetix-notary-service/internal/anchoring"
	"github.com/WildanFrananda/kinetix-notary-service/internal/config"
	"github.com/WildanFrananda/kinetix-notary-service/internal/metrics"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/verifying"
	"github.com/WildanFrananda/kinetix-notary-service/internal/wiring"
)

const anvilLedger = "0x5fbdb2315678afecb367f032d93f642f64180aa3"

func lookup(values map[string]string) config.Lookup {
	return func(key string) (string, bool) {
		value, ok := values[key]
		return value, ok
	}
}

func TestTheGraphResolvesEveryDependencyOfTheVerifier(t *testing.T) {
	container, err := wiring.Verifier(
		lookup(map[string]string{
			config.RPCURLKey:        "http://127.0.0.1:1",
			config.ChainIDKey:       "31337",
			config.LedgerAddressKey: anvilLedger,
			config.RPCTimeoutKey:    "2s",
		}),
	)
	require.NoError(t, err)

	err = container.Invoke(
		func(client *ethclient.Client, scheme notary.Scheme, ledger notary.Ledger, v *verifying.Verifier) {
			defer client.Close()

			require.NotNil(t, v)
			require.NotNil(t, ledger)
			require.Equal(t, uint16(1), scheme.Version())
		},
	)

	require.NoError(t, err)
}

func TestAMissingSettingSurfacesAsTheRootCause(t *testing.T) {
	container, err := wiring.Verifier(lookup(map[string]string{}))
	require.NoError(t, err)

	err = container.Invoke(func(*verifying.Verifier) {})

	require.Error(t, err)
	require.ErrorContains(t, dig.RootCause(err), config.RPCURLKey+" is not set")
}

func TestTheDaemonGraphNamesEveryMissingSettingBeforeTouchingAnything(t *testing.T) {
	container, err := wiring.Daemon(
		lookup(map[string]string{}),
		metrics.Build{Service: "kinetix-notary-service", Version: "test"},
		slog.New(slog.NewTextHandler(io.Discard, nil)),
	)
	require.NoError(t, err)

	err = container.Invoke(func(*anchoring.Runner) {})

	require.Error(t, err)

	cause := dig.RootCause(err)
	for _, key := range []string{config.DatabaseURLKey, config.AnchorerKeyKey, config.PaymentAddressKey, config.MaxFeeKey} {
		require.ErrorContains(t, cause, key)
	}
}
