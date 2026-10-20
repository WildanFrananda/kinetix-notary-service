package answering

import "github.com/WildanFrananda/kinetix-notary-service/internal/notary"

type Venue struct {
	Court     notary.WithholdingCourt
	Chain     notary.AnswerChain
	FromBlock uint64
}
