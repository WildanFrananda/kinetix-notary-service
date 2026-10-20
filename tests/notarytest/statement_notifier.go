package notarytest

import (
	"context"
	"fmt"
	"slices"
	"sync"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type StatementNotifier struct {
	mu     sync.Mutex
	calls  []StatementCall
	Answer func(holder notary.Holder, epoch uint64) (string, error)
}

var _ notary.StatementNotifier = (*StatementNotifier)(nil)

func (n *StatementNotifier) NotifyStatement(_ context.Context, holder notary.Holder, epoch uint64) (string, error) {
	n.mu.Lock()
	n.calls = append(n.calls, StatementCall{Holder: holder, Epoch: epoch})
	count := len(n.calls)
	answer := n.Answer
	n.mu.Unlock()

	if answer != nil {
		return answer(holder, epoch)
	}

	return fmt.Sprintf("NTF-%d", count), nil
}

func (n *StatementNotifier) Calls() []StatementCall {
	n.mu.Lock()
	defer n.mu.Unlock()

	return slices.Clone(n.calls)
}
