package wiring_test

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
	"go.uber.org/dig"
	"google.golang.org/grpc"

	"github.com/WildanFrananda/kinetix-notary-service/internal/approving"
	"github.com/WildanFrananda/kinetix-notary-service/internal/config"
	"github.com/WildanFrananda/kinetix-notary-service/internal/metrics"
	"github.com/WildanFrananda/kinetix-notary-service/internal/payment"
	"github.com/WildanFrananda/kinetix-notary-service/internal/wiring"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

func TestTheWatchGraphResolvesAndServesOnlyItsOwnMetrics(t *testing.T) {
	pki := notarytest.NewPKI(t)
	container, err := wiring.Watch(
		lookup(map[string]string{
			config.RPCURLKey:         "http://127.0.0.1:1",
			config.ChainIDKey:        "31337",
			config.LedgerAddressKey:  anvilLedger,
			config.CourtAddressKey:   anvilCourt,
			config.RPCTimeoutKey:     "2s",
			config.CouncilAddressKey: "0x9fe46736679d2d9a65f0992f2272de9f3c7fa6e0",
			config.WatchKeyKey:       "0x7c0883a69102937d6231471b5dbb6204fe5129617082792ae468d01a3f362318",
			config.PaymentAddressKey: "127.0.0.1:1",
			config.PKIDirKey:         pki.Issue(t, "kinetix.local", "notary-watch"),
			config.WatchGRPCPortKey:  "50070",
			config.HTTPPortKey:       "8011",
		}),
		metrics.Build{Service: "kinetix-notary-watch-service", Version: "test"},
		slog.New(slog.NewTextHandler(io.Discard, nil)),
	)
	require.NoError(t, err)

	err = container.Invoke(
		func(
			client *ethclient.Client,
			records *payment.Records,
			approver *approving.Approver,
			server *grpc.Server,
			handler http.Handler,
		) {
			defer client.Close()
			defer func() { _ = records.Close() }()

			require.NotNil(t, approver)
			require.Contains(t, server.GetServiceInfo(), "notary.v1.EpochWatchService")

			response := httptest.NewRecorder()
			handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/metrics", nil))
			body := response.Body.String()
			require.Contains(t, body, `kinetix_build_info{service="kinetix-notary-watch-service",version="test"} 1`)
			require.Contains(t, body, `kinetix_notary_watch_approvals_total{outcome="approved"} 0`)
			require.Contains(
				t,
				body,
				`kinetix_grpc_client_calls_total{grpc_code="OK",grpc_method="/payment.v1.LiabilityRecordService/GetLiabilitySnapshotRecord",peer="payment"} 0`,
			)
			require.Contains(
				t,
				body,
				`kinetix_grpc_server_calls_total{grpc_code="OK",grpc_method="/notary.v1.EpochWatchService/ApproveEpoch"} 0`,
			)
			require.False(t, strings.Contains(body, "kinetix_notary_last_finalized_epoch"))

			ready := httptest.NewRecorder()
			handler.ServeHTTP(ready, httptest.NewRequest(http.MethodGet, "/health/ready", nil))
			require.Equal(t, http.StatusServiceUnavailable, ready.Code, "no ledger answers at 127.0.0.1:1")
		},
	)

	require.NoError(t, err)
}

func TestTheWatchGraphNamesEveryMissingSetting(t *testing.T) {
	container, err := wiring.Watch(
		lookup(map[string]string{}),
		metrics.Build{Service: "kinetix-notary-watch-service", Version: "test"},
		slog.New(slog.NewTextHandler(io.Discard, nil)),
	)
	require.NoError(t, err)

	err = container.Invoke(func(*approving.Approver) {})

	cause := dig.RootCause(err)
	for _, key := range []string{config.CouncilAddressKey, config.WatchKeyKey, config.PaymentAddressKey, config.WatchGRPCPortKey} {
		require.ErrorContains(t, cause, key)
	}
}
