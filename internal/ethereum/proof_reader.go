package ethereum

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/WildanFrananda/kinetix-notary-service/internal/contractgen/ledgertwo"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type ProofReader struct {
	binding  *ledgertwo.LiabilityLedgerV2
	contract *bind.BoundContract
	timeout  time.Duration
}

var _ notary.ProofReader = (*ProofReader)(nil)

func NewProofReader(client *ethclient.Client, settings Settings) *ProofReader {
	binding := ledgertwo.NewLiabilityLedgerV2()

	return &ProofReader{
		binding:  binding,
		contract: binding.Instance(client, common.Address(settings.Ledger.Address)),
		timeout:  settings.Timeout,
	}
}

func (r *ProofReader) Depth(ctx context.Context) (int, error) {
	const op = "ethereum.ProofReader.Depth"

	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	depth, err := bind.Call(r.contract, &bind.CallOpts{Context: ctx}, r.binding.PackDepth(), r.binding.UnpackDepth)
	if err != nil {
		return 0, notary.Errf(notary.KindWrongLedger, op, err, "the ledger does not answer as a scheme-2 ledger")
	}

	return int(depth), nil
}

func (r *ProofReader) IsProven(ctx context.Context, number uint64, root notary.Digest, leafCount uint64) (bool, error) {
	const op = "ethereum.ProofReader.IsProven"

	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	proven, err := bind.Call(
		r.contract,
		&bind.CallOpts{Context: ctx},
		r.binding.PackIsProven(number, root, leafCount),
		r.binding.UnpackIsProven,
	)
	if err != nil {
		return false, unavailable(op, err)
	}

	return proven, nil
}
