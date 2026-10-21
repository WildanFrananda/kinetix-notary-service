package answering

import "github.com/WildanFrananda/kinetix-notary-service/internal/notary"

type Venue struct {
	Identity  notary.CourtIdentity
	Court     notary.WithholdingCourt
	Chain     notary.AnswerChain
	FromBlock uint64
}
