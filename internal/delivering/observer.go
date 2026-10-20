package delivering

import "github.com/WildanFrananda/kinetix-notary-service/internal/notary"

type Observer interface {
	Statement(outcome Outcome)
	StatementBacklog(backlog notary.StatementBacklog)
}
