package httpapi_test

import (
	"context"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/registering"
)

type registrant struct {
	issued    registering.Issued
	err       error
	principal notary.Holder
	request   notary.RegistrationRequest
	signature notary.Signature
	calls     int
}

func (r *registrant) Register(
	_ context.Context,
	principal notary.Holder,
	request notary.RegistrationRequest,
	signature notary.Signature,
) (registering.Issued, error) {
	r.calls++
	r.principal = principal
	r.request = request
	r.signature = signature

	return r.issued, r.err
}
