package sponsorcli

import (
	"context"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type Sponsorship interface {
	Holder() notary.Address
	Calls() notary.SponsoredCalls
	Identity(ctx context.Context) (notary.PaymasterIdentity, error)
	Submit(
		ctx context.Context,
		identity notary.PaymasterIdentity,
		call notary.SponsoredCall,
	) (notary.OperationReceipt, error)
}
