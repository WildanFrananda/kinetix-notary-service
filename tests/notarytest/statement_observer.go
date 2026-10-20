package notarytest

import (
	"slices"
	"sync"

	"github.com/WildanFrananda/kinetix-notary-service/internal/delivering"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type StatementObserver struct {
	mu       sync.Mutex
	outcomes []delivering.Outcome
	backlog  []notary.StatementBacklog
}

var _ delivering.Observer = (*StatementObserver)(nil)

func (o *StatementObserver) Statement(outcome delivering.Outcome) {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.outcomes = append(o.outcomes, outcome)
}

func (o *StatementObserver) StatementBacklog(backlog notary.StatementBacklog) {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.backlog = append(o.backlog, backlog)
}

func (o *StatementObserver) Outcomes() []delivering.Outcome {
	o.mu.Lock()
	defer o.mu.Unlock()

	return slices.Clone(o.outcomes)
}

func (o *StatementObserver) LastBacklog() (notary.StatementBacklog, bool) {
	o.mu.Lock()
	defer o.mu.Unlock()

	if len(o.backlog) == 0 {
		return notary.StatementBacklog{}, false
	}

	return o.backlog[len(o.backlog)-1], true
}
