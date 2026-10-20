package answering

import "github.com/WildanFrananda/kinetix-notary-service/internal/notary"

type waiting struct {
	venue   int
	request notary.WithholdingRequest
}
