package ethereum

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

func (l *Ledger) ConfidentialEpoch(ctx context.Context, number uint64) (notary.ConfidentialEpoch, error) {
	const op = "ethereum.Ledger.ConfidentialEpoch"

	epoch, err := l.Epoch(ctx, number)
	if err != nil {
		return notary.ConfidentialEpoch{}, err
	}

	ctx, cancel := context.WithTimeout(ctx, l.timeout)
	defer cancel()

	final, err := l.finalizedBlock(ctx)
	if err != nil {
		return notary.ConfidentialEpoch{}, unavailable(op, err)
	}

	opts := &bind.CallOpts{Context: ctx, BlockNumber: final}

	treeRoot, err := bind.Call(
		l.statements,
		opts,
		l.confidential.PackTreeRootOf(number),
		l.confidential.UnpackTreeRootOf,
	)
	if err != nil {
		return notary.ConfidentialEpoch{}, l.classify(op, err)
	}

	depth, err := bind.Call(l.statements, opts, l.confidential.PackDepth(), l.confidential.UnpackDepth)
	if err != nil {
		return notary.ConfidentialEpoch{}, notary.Errf(
			notary.KindWrongLedger,
			op,
			err,
			"the ledger is not a scheme-2 ledger",
		)
	}

	return notary.ConfidentialEpoch{
		Number:    number,
		TreeRoot:  treeRoot,
		LeafCount: epoch.LeafCount,
		Depth:     int(depth),
	}, nil
}
