package verifycli

import (
	"context"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type EvidenceClaimVerifier interface {
	Verify(ctx context.Context, claim notary.EvidenceClaim) (notary.EvidenceVerification, error)
}
