package httpapi

import "encoding/json"

type proofResponse struct {
	Epoch  string            `json:"epoch"`
	AsOf   string            `json:"asOf"`
	Digest string            `json:"digest"`
	Ledger ledgerResponse    `json:"ledger"`
	Proofs []json.RawMessage `json:"proofs"`
}
