package notarytest

import (
	"context"
	"sync"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type Bundler struct {
	mu          sync.Mutex
	Gas         notary.OperationGas
	EstimateErr error
	SendErr     error
	Answer      *notary.Digest
	Hasher      func(notary.UserOperation) notary.Digest
	Landed      *notary.OperationReceipt
	Estimated   []notary.UserOperation
	Sent        []notary.UserOperation
	Delegations []*notary.Delegation
}

var _ notary.Bundler = (*Bundler)(nil)

func (b *Bundler) Estimate(
	_ context.Context,
	op notary.UserOperation,
	_ notary.Address,
	delegation *notary.Delegation,
) (notary.OperationGas, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.Estimated = append(b.Estimated, op)
	b.Delegations = append(b.Delegations, delegation)

	return b.Gas, b.EstimateErr
}

func (b *Bundler) Send(
	_ context.Context,
	op notary.UserOperation,
	_ notary.Address,
	_ *notary.Delegation,
) (notary.Digest, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.Sent = append(b.Sent, op)

	if b.Answer != nil {
		return *b.Answer, b.SendErr
	}

	if b.Hasher != nil {
		return b.Hasher(op), b.SendErr
	}

	return notary.Digest{}, b.SendErr
}

func (b *Bundler) Receipt(_ context.Context, hash notary.Digest) (notary.OperationReceipt, bool, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.Landed == nil {
		return notary.OperationReceipt{}, false, nil
	}

	receipt := *b.Landed
	receipt.Hash = hash

	return receipt, true, nil
}
