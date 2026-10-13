package httpapi

import "github.com/WildanFrananda/kinetix-notary-service/internal/notary"

type Authenticator interface {
	Principal(authorization string) (notary.Holder, error)
}
