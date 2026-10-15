package notarytest

import (
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type Readers struct {
	Ledgers map[notary.Address]notary.Ledger
	Courts  map[notary.Address]notary.Court
}

var _ notary.GenerationReaders = (*Readers)(nil)

func (r *Readers) Ledger(generation notary.Generation) notary.Ledger {
	return r.Ledgers[generation.Ledger]
}

func (r *Readers) Court(generation notary.Generation) notary.Court {
	return r.Courts[generation.Court]
}
