package ethereum

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/WildanFrananda/kinetix-notary-service/internal/contractgen/entrypoint"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

var delegationPrefix = []byte{0xef, 0x01, 0x00}

type SponsorChain struct {
	client  *ethclient.Client
	binding *entrypoint.EntryPoint
	timeout time.Duration
}

var _ notary.SponsorChain = (*SponsorChain)(nil)

func NewSponsorChain(client *ethclient.Client, timeout time.Duration) *SponsorChain {
	return &SponsorChain{client: client, binding: entrypoint.NewEntryPoint(), timeout: timeout}
}

func (c *SponsorChain) OperationNonce(ctx context.Context, entryPoint, sender notary.Address) (uint64, error) {
	const op = "ethereum.SponsorChain.OperationNonce"

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	nonce, err := bind.Call(
		c.binding.Instance(c.client, common.Address(entryPoint)),
		&bind.CallOpts{Context: ctx},
		c.binding.PackGetNonce(common.Address(sender), new(big.Int)),
		c.binding.UnpackGetNonce,
	)
	if err != nil {
		return 0, unavailable(op, fmt.Errorf("entry point: %w", err))
	}

	if !nonce.IsUint64() {
		return 0, notary.Errf(
			notary.KindNotSponsored,
			op,
			nil,
			"%s has an operation nonce beyond 64 bits",
			sender.Hex(),
		)
	}

	return nonce.Uint64(), nil
}

func (c *SponsorChain) AccountNonce(ctx context.Context, account notary.Address) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	nonce, err := c.client.NonceAt(ctx, common.Address(account), nil)
	if err != nil {
		return 0, unavailable("ethereum.SponsorChain.AccountNonce", err)
	}

	return nonce, nil
}

func (c *SponsorChain) DelegateOf(ctx context.Context, account notary.Address) (notary.Address, error) {
	const op = "ethereum.SponsorChain.DelegateOf"

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	code, err := c.client.CodeAt(ctx, common.Address(account), nil)
	if err != nil {
		return notary.Address{}, unavailable(op, err)
	}

	if len(code) == 0 {
		return notary.Address{}, nil
	}

	if len(code) != len(delegationPrefix)+len(notary.Address{}) || !bytes.HasPrefix(code, delegationPrefix) {
		return notary.Address{}, notary.Errf(
			notary.KindNotSponsored,
			op,
			nil,
			"%s is a contract, not a key: only an account a key controls can be sponsored",
			account.Hex(),
		)
	}

	return notary.Address(code[len(delegationPrefix):]), nil
}

func (c *SponsorChain) Fees(ctx context.Context) (uint64, uint64, error) {
	const op = "ethereum.SponsorChain.Fees"

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	head, err := c.client.HeaderByNumber(ctx, nil)
	if err != nil {
		return 0, 0, unavailable(op, err)
	}

	tip, err := c.client.SuggestGasTipCap(ctx)
	if err != nil {
		return 0, 0, unavailable(op, err)
	}

	if head.BaseFee == nil {
		return 0, 0, notary.Errf(notary.KindNotSponsored, op, nil, "the chain reports no base fee")
	}

	maxFee := new(big.Int).Add(new(big.Int).Lsh(head.BaseFee, 1), tip)
	if !maxFee.IsUint64() || !tip.IsUint64() {
		return 0, 0, notary.Errf(notary.KindNotSponsored, op, nil, "the chain's fees exceed 64 bits of wei")
	}

	return maxFee.Uint64(), tip.Uint64(), nil
}

func (c *SponsorChain) BlockTime(ctx context.Context) (notary.UnixTime, error) {
	const op = "ethereum.SponsorChain.BlockTime"

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	head, err := c.client.HeaderByNumber(ctx, nil)
	if err != nil {
		return notary.UnixTime{}, unavailable(op, err)
	}

	now, err := notary.UnixTimeFromSeconds(head.Time)
	if err != nil {
		return notary.UnixTime{}, notary.Errf(notary.KindSourceRefused, op, err, "the chain's head carries no time")
	}

	return now, nil
}
