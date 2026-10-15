package ethereum

import (
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type Readers struct {
	client   *ethclient.Client
	settings Settings
	scheme   notary.Scheme
}

var _ notary.GenerationReaders = (*Readers)(nil)

func NewReaders(client *ethclient.Client, settings Settings, scheme notary.Scheme) *Readers {
	return &Readers{client: client, settings: settings, scheme: scheme}
}

func (r *Readers) Ledger(generation notary.Generation) notary.Ledger {
	return NewLedger(r.client, r.of(generation), r.scheme)
}

func (r *Readers) Court(generation notary.Generation) notary.Court {
	return NewCourt(r.client, r.of(generation))
}

func (r *Readers) of(generation notary.Generation) Settings {
	settings := r.settings
	settings.Ledger = generation.LedgerIdentity()
	settings.Court = generation.Court
	settings.Withholding = generation.Withholding

	return settings
}
