package httpapi

import (
	"context"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/registering"
)

type Registrant interface {
	Register(
		ctx context.Context,
		principal notary.Holder,
		request notary.RegistrationRequest,
		signature notary.Signature,
	) (registering.Issued, error)
}
