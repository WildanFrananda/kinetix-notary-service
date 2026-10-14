package notarytest

import (
	"context"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type Registrations struct {
	Saved      map[notary.Digest]notary.IssuedRegistration
	Principals map[notary.Digest]notary.Holder
	Err        error
}

var _ notary.RegistrationStore = (*Registrations)(nil)

func NewRegistrations() *Registrations {
	return &Registrations{
		Saved:      map[notary.Digest]notary.IssuedRegistration{},
		Principals: map[notary.Digest]notary.Holder{},
	}
}

func (r *Registrations) SaveRegistrations(
	_ context.Context,
	principal notary.Holder,
	issued []notary.IssuedRegistration,
) error {
	if r.Err != nil {
		return r.Err
	}

	for _, one := range issued {
		r.Saved[one.Digest] = one
		r.Principals[one.Digest] = principal
	}

	return nil
}

func (r *Registrations) RegistrationByDigest(
	_ context.Context,
	digest notary.Digest,
) (notary.IssuedRegistration, bool, error) {
	one, ok := r.Saved[digest]

	return one, ok, r.Err
}
