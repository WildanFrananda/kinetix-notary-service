package httpapi

import (
	"context"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/proving"
)

type ProofIssuer interface {
	Proofs(ctx context.Context, holder notary.Holder, ref proving.EpochRef) (proving.Bundle, error)
}
