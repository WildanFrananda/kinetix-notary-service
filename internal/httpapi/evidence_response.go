package httpapi

import "encoding/json"

type evidenceResponse struct {
	Proofs  []json.RawMessage `json:"proofs"`
	Pending int               `json:"pending"`
}
