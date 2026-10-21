package httpapi_test

import (
	"context"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/registering"
)

type registrant struct {
	domain    notary.TypedDataDomain
	issued    registering.Issued
	err       error
	request   notary.RegistrationRequest
	signature notary.Signature
	calls     int
}

func (r *registrant) Domain() notary.TypedDataDomain {
	return r.domain
}

func (r *registrant) Register(
	_ context.Context,
	request notary.RegistrationRequest,
	signature notary.Signature,
) (registering.Issued, error) {
	r.calls++
	r.request = request
	r.signature = signature

	return r.issued, r.err
}
