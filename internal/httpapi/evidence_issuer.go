package httpapi

import (
	"context"

	"github.com/WildanFrananda/kinetix-notary-service/internal/evidencing"
)

type EvidenceIssuer interface {
	OrderProofs(ctx context.Context, principal, orderNumber string) (evidencing.Bundle, error)
	ReviewProofs(ctx context.Context, principal, reviewID string) (evidencing.Bundle, error)
}
