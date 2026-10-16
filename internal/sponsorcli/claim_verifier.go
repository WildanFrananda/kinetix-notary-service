package sponsorcli

import (
	"context"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type ClaimVerifier interface {
	Verify(ctx context.Context, claim notary.Claim) (notary.Verification, error)
}
