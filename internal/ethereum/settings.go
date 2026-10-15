package ethereum

import (
	"time"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type Settings struct {
	RPCURL      string
	Ledger      notary.LedgerIdentity
	Court       notary.Address
	Withholding notary.Address
	Council     notary.Address
	Directory   notary.DirectoryIdentity
	Timeout     time.Duration
}
