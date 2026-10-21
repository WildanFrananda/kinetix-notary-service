package wiring

import (
	"fmt"
	"net/http"

	"go.uber.org/dig"
	"google.golang.org/grpc"

	"github.com/WildanFrananda/kinetix-notary-service/internal/approving"
	"github.com/WildanFrananda/kinetix-notary-service/internal/config"
	"github.com/WildanFrananda/kinetix-notary-service/internal/epochwatch"
	"github.com/WildanFrananda/kinetix-notary-service/internal/ethereum"
	"github.com/WildanFrananda/kinetix-notary-service/internal/httpapi"
	"github.com/WildanFrananda/kinetix-notary-service/internal/mesh"
	"github.com/WildanFrananda/kinetix-notary-service/internal/metrics"
)

const watchCaller = "notary"

type watchHTTPParams struct {
	dig.In
	Metrics *metrics.Watch
	Heads   *ethereum.HeadReader
}

func provideWatchHTTP(p watchHTTPParams) http.Handler {
	ready := func(r *http.Request) error {
		if err := p.Heads.Reachable(r.Context()); err != nil {
			return fmt.Errorf("the ledger cannot be read: %w", err)
		}

		return nil
	}

	return httpapi.NewWatchRouter(p.Metrics.Instrument, p.Metrics.Handler(), ready)
}

type watchGRPCParams struct {
	dig.In
	Config   config.Watch
	Approver *approving.Approver
	Metrics  *metrics.Watch
}

func provideWatchGRPC(p watchGRPCParams) (*grpc.Server, error) {
	creds, err := mesh.ServerCredentials(
		mesh.Identity{PKIDir: p.Config.PKIDir, TrustDomains: p.Config.TrustDomains},
		[]string{watchCaller},
	)
	if err != nil {
		return nil, err
	}

	options := append([]grpc.ServerOption{grpc.Creds(creds)}, p.Metrics.ServerOptions()...)
	server := grpc.NewServer(options...)
	epochwatch.Register(server, p.Approver)

	return server, nil
}
