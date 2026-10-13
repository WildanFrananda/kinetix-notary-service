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

	"github.com/WildanFrananda/kinetix-notary-service/internal/contractgen/court"
	"github.com/WildanFrananda/kinetix-notary-service/internal/contractgen/courtinterface"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type Court struct {
	client      *ethclient.Client
	binding     *court.StatementCourt
	contract    *bind.BoundContract
	expected    notary.CourtIdentity
	ledger      notary.Address
	interfaceID [4]byte
	timeout     time.Duration
}

var _ notary.Court = (*Court)(nil)

func NewCourt(client *ethclient.Client, settings Settings) *Court {
	binding := court.NewStatementCourt()

	return &Court{
		client:      client,
		binding:     binding,
		contract:    binding.Instance(client, common.Address(settings.Court)),
		expected:    notary.CourtIdentity{ChainID: settings.Ledger.ChainID, Address: settings.Court},
		ledger:      settings.Ledger.Address,
		interfaceID: courtInterfaceID(),
		timeout:     settings.Timeout,
	}
}

func (c *Court) Identity(ctx context.Context) (notary.CourtIdentity, error) {
	const op = "ethereum.Court.Identity"

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	chainID, err := c.client.ChainID(ctx)
	if err != nil {
		return notary.CourtIdentity{}, unavailable(op, err)
	}

	if !chainID.IsUint64() || chainID.Uint64() != c.expected.ChainID {
		return notary.CourtIdentity{}, notary.Errf(
			notary.KindWrongCourt,
			op,
			nil,
			"the RPC endpoint serves chain %s, the court is configured on chain %d",
			chainID,
			c.expected.ChainID,
		)
	}

	code, err := c.client.CodeAt(ctx, common.Address(c.expected.Address), nil)
	if err != nil {
		return notary.CourtIdentity{}, unavailable(op, err)
	}

	if len(code) == 0 {
		return notary.CourtIdentity{}, notary.Errf(
			notary.KindWrongCourt,
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
		return notary.CourtIdentity{}, c.wrap(op, err)
	}

	if !supported {
		return notary.CourtIdentity{}, notary.Errf(
			notary.KindWrongCourt,
			op,
			nil,
			"%s does not announce the statement court interface",
			c.expected.Address.Hex(),
		)
	}

	judged, err := probe(ctx, c.client, c.expected.Address, c.binding.PackLedger(), c.binding.UnpackLedger)
	if err != nil {
		return notary.CourtIdentity{}, c.wrap(op, err)
	}

	if notary.Address(judged) != c.ledger {
		return notary.CourtIdentity{}, notary.Errf(
			notary.KindWrongCourt,
			op,
			nil,
			"%s judges the ledger %s, not %s",
			c.expected.Address.Hex(),
			judged.Hex(),
			c.ledger.Hex(),
		)
	}

	return c.expected, nil
}

func (c *Court) SignerOf(ctx context.Context, epoch uint64) (notary.Address, error) {
	signer, err := call(c, ctx, "ethereum.Court.SignerOf", c.binding.PackSignerOf(epoch), c.binding.UnpackSignerOf)

	return notary.Address(signer), err
}

func (c *Court) StatementDigest(ctx context.Context, statement notary.Statement) (notary.Digest, error) {
	digest, err := call(
		c,
		ctx,
		"ethereum.Court.StatementDigest",
		c.binding.PackStatementDigest(bound(statement)),
		c.binding.UnpackStatementDigest,
	)

	return notary.Digest(digest), err
}

func (c *Court) RecoverSigner(
	ctx context.Context,
	statement notary.Statement,
	signature notary.Signature,
) (notary.Address, error) {
	signer, err := call(
		c,
		ctx,
		"ethereum.Court.RecoverSigner",
		c.binding.PackRecoverSigner(bound(statement), signature[:]),
		c.binding.UnpackRecoverSigner,
	)

	return notary.Address(signer), err
}

func (c *Court) FaultCount(ctx context.Context, epoch uint64) (uint64, error) {
	count, err := call(c, ctx, "ethereum.Court.FaultCount", c.binding.PackFaultCount(epoch), c.binding.UnpackFaultCount)
	if err != nil {
		return 0, err
	}

	return counted("ethereum.Court.FaultCount", count)
}

func (c *Court) TotalFaults(ctx context.Context) (uint64, error) {
	count, err := call(c, ctx, "ethereum.Court.TotalFaults", c.binding.PackTotalFaults(), c.binding.UnpackTotalFaults)
	if err != nil {
		return 0, err
	}

	return counted("ethereum.Court.TotalFaults", count)
}

func call[T any](c *Court, ctx context.Context, op string, calldata []byte, unpack func([]byte) (T, error)) (T, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	value, err := bind.Call(c.contract, &bind.CallOpts{Context: ctx}, calldata, unpack)
	if err != nil {
		var zero T
		if errors.Is(err, bind.ErrNoCode) {
			return zero, notary.Errf(
				notary.KindWrongCourt,
				op,
				nil,
				"chain %d has no contract at %s",
				c.expected.ChainID,
				c.expected.Address.Hex(),
			)
		}

		return zero, unavailable(op, err)
	}

	return value, nil
}

func counted(op string, count *big.Int) (uint64, error) {
	if !count.IsUint64() {
		return 0, notary.Errf(notary.KindWrongCourt, op, nil, "the court reports %s faults", count)
	}

	return count.Uint64(), nil
}

func bound(statement notary.Statement) court.EpochStatement {
	return court.EpochStatement{
		Ledger:       common.Address(statement.Ledger),
		Epoch:        statement.Epoch,
		Commitment:   statement.Commitment,
		Amount:       statement.Amount,
		Index:        statement.Index,
		SiblingsHash: statement.SiblingsHash,
	}
}

func (c *Court) wrap(op string, err error) error {
	if errors.Is(err, errNotALedger) {
		return notary.Errf(
			notary.KindWrongCourt,
			op,
			nil,
			"%s does not answer as a statement court",
			c.expected.Address.Hex(),
		)
	}

	return unavailable(op, fmt.Errorf("court: %w", err))
}

func courtInterfaceID() [4]byte {
	var id [4]byte
	for _, method := range courtinterface.NewIStatementCourt().GetABI().Methods {
		for i := range id {
			id[i] ^= method.ID[i]
		}
	}

	return id
}
