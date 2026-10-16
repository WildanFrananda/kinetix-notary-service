package ethereum

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/WildanFrananda/kinetix-notary-service/internal/contractgen/paymaster"
	"github.com/WildanFrananda/kinetix-notary-service/internal/contractgen/paymasterbase"
	"github.com/WildanFrananda/kinetix-notary-service/internal/contractgen/paymasterinterface"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type Paymaster struct {
	client      *ethclient.Client
	binding     *paymaster.NotaryPaymaster
	contract    *bind.BoundContract
	chainID     uint64
	address     notary.Address
	interfaceID [4]byte
	timeout     time.Duration
}

var _ notary.Paymaster = (*Paymaster)(nil)

func NewPaymaster(client *ethclient.Client, chainID uint64, address notary.Address, timeout time.Duration) *Paymaster {
	binding := paymaster.NewNotaryPaymaster()

	return &Paymaster{
		client:      client,
		binding:     binding,
		contract:    binding.Instance(client, common.Address(address)),
		chainID:     chainID,
		address:     address,
		interfaceID: paymasterInterfaceID(),
		timeout:     timeout,
	}
}

func (p *Paymaster) Identity(ctx context.Context) (notary.PaymasterIdentity, error) {
	const op = "ethereum.Paymaster.Identity"

	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	served, err := p.client.ChainID(ctx)
	if err != nil {
		return notary.PaymasterIdentity{}, unavailable(op, err)
	}

	if !served.IsUint64() || served.Uint64() != p.chainID {
		return notary.PaymasterIdentity{}, notary.Errf(
			notary.KindNotSponsored,
			op,
			nil,
			"the RPC endpoint serves chain %s, the paymaster is configured on chain %d",
			served,
			p.chainID,
		)
	}

	supported, err := probe(
		ctx,
		p.client,
		p.address,
		p.binding.PackSupportsInterface(p.interfaceID),
		p.binding.UnpackSupportsInterface,
	)
	if err != nil {
		return notary.PaymasterIdentity{}, p.wrap(op, err)
	}

	if !supported {
		return notary.PaymasterIdentity{}, notary.Errf(
			notary.KindNotSponsored,
			op,
			nil,
			"%s does not announce the notary paymaster interface",
			p.address.Hex(),
		)
	}

	identity := notary.PaymasterIdentity{ChainID: p.chainID, Address: p.address}
	for _, read := range []struct {
		into *notary.Address
		pack []byte
	}{
		{&identity.EntryPoint, p.binding.PackEntryPoint()},
		{&identity.Account, p.binding.PackAccount()},
		{&identity.Court, p.binding.PackCourt()},
		{&identity.Withholding, p.binding.PackWithholding()},
		{&identity.Signer, p.binding.PackStatementSigner()},
	} {
		address, err := bind.Call(p.contract, &bind.CallOpts{Context: ctx}, read.pack, p.unpackAddress)
		if err != nil {
			return notary.PaymasterIdentity{}, p.wrap(op, err)
		}

		*read.into = notary.Address(address)
	}

	ceiling, err := bind.Call(
		p.contract,
		&bind.CallOpts{Context: ctx},
		p.binding.PackMaxCostPerOperation(),
		p.binding.UnpackMaxCostPerOperation,
	)
	if err != nil {
		return notary.PaymasterIdentity{}, p.wrap(op, err)
	}

	identity.Ceiling = ceiling

	return identity, nil
}

func (p *Paymaster) Funds(ctx context.Context) (notary.PaymasterFunds, error) {
	const op = "ethereum.Paymaster.Funds"

	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	info, err := bind.Call(
		p.contract,
		&bind.CallOpts{Context: ctx},
		p.binding.PackGetDepositInfo(),
		p.binding.UnpackGetDepositInfo,
	)
	if err != nil {
		return notary.PaymasterFunds{}, p.wrap(op, err)
	}

	return notary.PaymasterFunds{Deposit: info.Deposit, Stake: info.Stake, Staked: info.Staked}, nil
}

func (p *Paymaster) Sponsorship(
	ctx context.Context,
	sender notary.Address,
	callData []byte,
	paymasterData []byte,
) (notary.SponsorshipWindow, error) {
	const op = "ethereum.Paymaster.Sponsorship"

	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	window, err := bind.Call(
		p.contract,
		&bind.CallOpts{Context: ctx},
		p.binding.PackSponsorship(common.Address(sender), callData, paymasterData),
		p.binding.UnpackSponsorship,
	)
	if err != nil {
		if raw, reverted := revertData(err); reverted {
			return notary.SponsorshipWindow{}, p.refusal(op, raw)
		}

		return notary.SponsorshipWindow{}, p.wrap(op, err)
	}

	after, err := windowEdge(op, window.ValidAfter)
	if err != nil {
		return notary.SponsorshipWindow{}, err
	}

	until, err := windowEdge(op, window.ValidUntil)
	if err != nil {
		return notary.SponsorshipWindow{}, err
	}

	return notary.SponsorshipWindow{After: after, Until: until}, nil
}

func windowEdge(op string, seconds *big.Int) (notary.UnixTime, error) {
	if seconds == nil || seconds.Sign() == 0 {
		return notary.UnixTime{}, nil
	}

	if !seconds.IsUint64() {
		return notary.UnixTime{}, notary.Errf(
			notary.KindSourceRefused,
			op,
			nil,
			"the paymaster named the time %s",
			seconds,
		)
	}

	edge, err := notary.UnixTimeFromSeconds(seconds.Uint64())
	if err != nil {
		return notary.UnixTime{}, notary.Errf(
			notary.KindSourceRefused,
			op,
			err,
			"the paymaster named the time %s",
			seconds,
		)
	}

	return edge, nil
}

func (p *Paymaster) unpackAddress(data []byte) (common.Address, error) {
	return p.binding.UnpackAccount(data)
}

func (p *Paymaster) wrap(op string, err error) error {
	if errors.Is(err, errNotALedger) || errors.Is(err, bind.ErrNoCode) {
		return notary.Errf(
			notary.KindNotSponsored,
			op,
			nil,
			"chain %d has no notary paymaster at %s",
			p.chainID,
			p.address.Hex(),
		)
	}

	return unavailable(op, fmt.Errorf("paymaster: %w", err))
}

func paymasterInterfaceID() [4]byte {
	declared := paymasterinterface.NewINotaryPaymaster().GetABI().Methods
	inherited := paymasterbase.NewIPaymaster().GetABI().Methods

	var id [4]byte

	for name, method := range declared {
		if _, ok := inherited[name]; ok {
			continue
		}

		for i := range id {
			id[i] ^= method.ID[i]
		}
	}

	return id
}
