package notarytest

import "github.com/WildanFrananda/kinetix-notary-service/internal/notary"

type StatementCall struct {
	Holder notary.Holder
	Epoch  uint64
}
