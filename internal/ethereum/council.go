package ethereum

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/WildanFrananda/kinetix-notary-service/internal/contractgen/council"
	"github.com/WildanFrananda/kinetix-notary-service/internal/contractgen/councilinterface"
	"github.com/WildanFrananda/kinetix-notary-service/internal/contractgen/ledger"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type Council struct {
	client         *ethclient.Client
	binding        *council.AnchorCouncil
	contract       *bind.BoundContract
	ledgerBinding  *ledger.LiabilityLedger
	ledgerContract *bind.BoundContract
	expected       notary.CouncilIdentity
	ledger         notary.Address
	interfaceID    [4]byte
	timeout        time.Duration
}

var _ notary.Council = (*Council)(nil)

func NewCouncil(client *ethclient.Client, settings Settings) *Council {
	binding := council.NewAnchorCouncil()
	ledgerBinding := ledger.NewLiabilityLedger()

	return &Council{
		client:         client,
		binding:        binding,
		contract:       binding.Instance(client, common.Address(settings.Council)),
		ledgerBinding:  ledgerBinding,
		ledgerContract: ledgerBinding.Instance(client, common.Address(settings.Ledger.Address)),
		expected:       notary.CouncilIdentity{ChainID: settings.Ledger.ChainID, Address: settings.Council},
		ledger:         settings.Ledger.Address,
		interfaceID:    councilInterfaceID(),
		timeout:        settings.Timeout,
	}
}

func (c *Council) Identity(ctx context.Context) (notary.CouncilIdentity, error) {
	const op = "ethereum.Council.Identity"

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	chainID, err := c.client.ChainID(ctx)
	if err != nil {
		return notary.CouncilIdentity{}, unavailable(op, err)
	}

	if !chainID.IsUint64() || chainID.Uint64() != c.expected.ChainID {
		return notary.CouncilIdentity{}, notary.Errf(
			notary.KindWrongCouncil,
			op,
			nil,
			"the RPC endpoint serves chain %s, the council is configured on chain %d",
			chainID,
			c.expected.ChainID,
		)
	}

	code, err := c.client.CodeAt(ctx, common.Address(c.expected.Address), nil)
	if err != nil {
		return notary.CouncilIdentity{}, unavailable(op, err)
	}

	if len(code) == 0 {
		return notary.CouncilIdentity{}, notary.Errf(
			notary.KindWrongCouncil,
			op,
			nil,
			"chain %d has no contract at %s",
			c.expected.ChainID,
			c.expected.Address.Hex(),
		)
	}

	supported, err := probe(
		ctx,
		c.client,
		c.expected.Address,
		c.binding.PackSupportsInterface(c.interfaceID),
		c.binding.UnpackSupportsInterface,
	)
	if err != nil {
		return notary.CouncilIdentity{}, c.wrap(op, err)
	}

	if !supported {
		return notary.CouncilIdentity{}, notary.Errf(
			notary.KindWrongCouncil,
			op,
			nil,
			"%s does not announce the anchor council interface",
			c.expected.Address.Hex(),
		)
	}

	anchors, err := probe(ctx, c.client, c.expected.Address, c.binding.PackLedger(), c.binding.UnpackLedger)
	if err != nil {
		return notary.CouncilIdentity{}, c.wrap(op, err)
	}

	if notary.Address(anchors) != c.ledger {
		return notary.CouncilIdentity{}, notary.Errf(
			notary.KindWrongCouncil,
			op,
			nil,
			"%s anchors the ledger %s, not %s",
			c.expected.Address.Hex(),
			anchors.Hex(),
			c.ledger.Hex(),
		)
	}

	return c.expected, nil
}

func (c *Council) IsMember(ctx context.Context, member notary.Address) (bool, error) {
	return councilCall(
		c,
		ctx,
		"ethereum.Council.IsMember",
		c.contract,
		c.binding.PackIsMember(common.Address(member)),
		c.binding.UnpackIsMember,
	)
}

func (c *Council) Threshold(ctx context.Context) (uint64, error) {
	threshold, err := councilCall(
		c,
		ctx,
		"ethereum.Council.Threshold",
		c.contract,
		c.binding.PackThreshold(),
		c.binding.UnpackThreshold,
	)

	return uint64(threshold), err
}

func (c *Council) ApprovalDigest(ctx context.Context, number uint64, epoch notary.Digest) (notary.Digest, error) {
	digest, err := councilCall(
		c,
		ctx,
		"ethereum.Council.ApprovalDigest",
		c.contract,
		c.binding.PackApprovalDigest(number, epoch),
		c.binding.UnpackApprovalDigest,
	)

	return notary.Digest(digest), err
}

func (c *Council) AnchorsLedger(ctx context.Context) (bool, error) {
	const op = "ethereum.Council.AnchorsLedger"

	role, err := councilCall(
		c,
		ctx,
		op,
		c.ledgerContract,
		c.ledgerBinding.PackANCHORERROLE(),
		c.ledgerBinding.UnpackANCHORERROLE,
	)
	if err != nil {
		return false, err
	}

	return councilCall(
		c,
		ctx,
		op,
		c.ledgerContract,
		c.ledgerBinding.PackHasRole(role, common.Address(c.expected.Address)),
		c.ledgerBinding.UnpackHasRole,
	)
}

func councilCall[T any](
	c *Council,
	ctx context.Context,
	op string,
	contract *bind.BoundContract,
	calldata []byte,
	unpack func([]byte) (T, error),
) (T, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	value, err := bind.Call(contract, &bind.CallOpts{Context: ctx}, calldata, unpack)
	if err != nil {
		var zero T
		if errors.Is(err, bind.ErrNoCode) {
			return zero, notary.Errf(
				notary.KindWrongCouncil,
				op,
				nil,
				"chain %d has no contract where the council or its ledger should be",
				c.expected.ChainID,
			)
		}

		return zero, unavailable(op, err)
	}

	return value, nil
}

func (c *Council) wrap(op string, err error) error {
	if errors.Is(err, errNotALedger) {
		return notary.Errf(
			notary.KindWrongCouncil,
			op,
			nil,
			"%s does not answer as an anchor council",
			c.expected.Address.Hex(),
		)
	}

	return unavailable(op, fmt.Errorf("council: %w", err))
}

func councilInterfaceID() [4]byte {
	var id [4]byte
	for _, method := range councilinterface.NewIAnchorCouncil().GetABI().Methods {
		for i := range id {
			id[i] ^= method.ID[i]
		}
	}

	return id
}
