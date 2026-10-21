package httpapi

import (
	"context"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/registering"
)

type Registrant interface {
	Domain() notary.TypedDataDomain
	Register(
		ctx context.Context,
		request notary.RegistrationRequest,
		signature notary.Signature,
	) (registering.Issued, error)
}
