package ethereum

import (
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type Readers struct {
	client   *ethclient.Client
	settings Settings
}

var _ notary.GenerationReaders = (*Readers)(nil)

func NewReaders(client *ethclient.Client, settings Settings) *Readers {
	return &Readers{client: client, settings: settings}
}

func (r *Readers) Ledger(generation notary.Generation) notary.Ledger {
	return NewLedgerOfScheme(r.client, r.of(generation), generation.Scheme)
}

func (r *Readers) Confidential(generation notary.Generation) notary.ConfidentialLedger {
	return NewLedgerOfScheme(r.client, r.of(generation), generation.Scheme)
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
