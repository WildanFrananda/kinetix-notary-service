package wiring

import (
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/dig"

	"github.com/WildanFrananda/kinetix-notary-service/internal/anchoring"
	"github.com/WildanFrananda/kinetix-notary-service/internal/evidencing"
	"github.com/WildanFrananda/kinetix-notary-service/internal/httpapi"
	"github.com/WildanFrananda/kinetix-notary-service/internal/metrics"
)

type httpParams struct {
	dig.In
	Metrics        *metrics.Metrics
	Runner         *anchoring.Runner
	Evidence       *evidencing.Runner
	Pool           *pgxpool.Pool
	Proofs         *httpapi.ProofHandler
	Registrations  *httpapi.RegistrationHandler
	OrderEvidence  *httpapi.EvidenceHandler `name:"order_evidence"`
	ReviewEvidence *httpapi.EvidenceHandler `name:"review_evidence"`
}

func provideHTTP(p httpParams) http.Handler {
	ready := func(r *http.Request) error {
		if alive, since := p.Runner.Alive(); !alive {
			return fmt.Errorf("the anchoring loop has not completed a pass for %s", since)
		}

		if alive, since := p.Evidence.Alive(); !alive {
			return fmt.Errorf("the evidence loop has not completed a pass for %s", since)
		}

		if err := p.Pool.Ping(r.Context()); err != nil {
			return fmt.Errorf("the database does not answer: %w", err)
		}

		return nil
	}

	return httpapi.NewRouter(
		p.Metrics.Instrument,
		p.Metrics.Handler(),
		ready,
		p.Proofs,
		p.Registrations,
		p.OrderEvidence,
		p.ReviewEvidence,
	)
}
