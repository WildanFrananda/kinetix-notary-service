package ethereum

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/WildanFrananda/kinetix-notary-service/internal/contractgen/ledger"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type HeadReader struct {
	binding  *ledger.LiabilityLedger
	contract *bind.BoundContract
	timeout  time.Duration
}

var _ notary.LedgerHeads = (*HeadReader)(nil)

func NewHeadReader(client *ethclient.Client, settings Settings) *HeadReader {
	binding := ledger.NewLiabilityLedger()

	return &HeadReader{
		binding:  binding,
		contract: binding.Instance(client, common.Address(settings.Ledger.Address)),
		timeout:  settings.Timeout,
	}
}

func (r *HeadReader) Head(ctx context.Context) (notary.LedgerHead, error) {
	return readHead(ctx, "ethereum.HeadReader.Head", r.contract, r.binding, r.timeout)
}

func readHead(
	ctx context.Context,
	op string,
	contract *bind.BoundContract,
	binding *ledger.LiabilityLedger,
	timeout time.Duration,
) (notary.LedgerHead, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	opts := &bind.CallOpts{Context: ctx}

	var head notary.LedgerHead

	next, err := bind.Call(contract, opts, binding.PackNextEpoch(), binding.UnpackNextEpoch)
	if err != nil {
		return head, unavailable(op, err)
	}

	digest, err := bind.Call(contract, opts, binding.PackHead(), binding.UnpackHead)
	if err != nil {
		return head, unavailable(op, err)
	}

	lastAsOf, err := bind.Call(contract, opts, binding.PackLastAsOf(), binding.UnpackLastAsOf)
	if err != nil {
		return head, unavailable(op, err)
	}

	genesisEpoch, err := bind.Call(contract, opts, binding.PackGenesisEpoch(), binding.UnpackGenesisEpoch)
	if err != nil {
		return head, unavailable(op, err)
	}

	genesisDigest, err := bind.Call(contract, opts, binding.PackGenesisDigest(), binding.UnpackGenesisDigest)
	if err != nil {
		return head, unavailable(op, err)
	}

	sealed, err := bind.Call(contract, opts, binding.PackIsSealed(), binding.UnpackIsSealed)
	if err != nil {
		return head, unavailable(op, err)
	}

	return notary.LedgerHead{
		NextEpoch:     next,
		Head:          digest,
		LastAsOf:      lastAsOf,
		GenesisEpoch:  genesisEpoch,
		GenesisDigest: genesisDigest,
		Sealed:        sealed,
	}, nil
}
