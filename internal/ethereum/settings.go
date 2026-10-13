package ethereum

import (
	"time"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type Settings struct {
	RPCURL  string
	Ledger  notary.LedgerIdentity
	Court   notary.Address
	Timeout time.Duration
}
