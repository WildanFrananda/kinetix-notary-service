package notary

import "context"

type RegistrationStore interface {
	SaveRegistrations(ctx context.Context, principal Holder, issued []IssuedRegistration) error
	RegistrationByDigest(ctx context.Context, digest Digest) (IssuedRegistration, bool, error)
}
