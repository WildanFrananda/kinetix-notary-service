package integration_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/contractgen/directory"
	"github.com/WildanFrananda/kinetix-notary-service/internal/contractgen/ledger"
	"github.com/WildanFrananda/kinetix-notary-service/internal/ethereum"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

func (c *chain) deployDirectory(
	t *testing.T,
	ledgerAddress, courtAddress, withholdingAddress common.Address,
) common.Address {
	t.Helper()

	input := directory.NewLedgerDirectory().PackConstructor(
		c.from,
		big.NewInt(3600),
		ledgerAddress,
		courtAddress,
		withholdingAddress,
		7*24*3600,
	)
	address, tx, err := bind.DeployContract(
		c.transactor(),
		common.FromHex(directory.LedgerDirectoryMetaData.Bin),
		c.client,
		input,
	)
	require.NoError(t, err)
	_, err = bind.WaitDeployed(context.Background(), c.client, tx.Hash())
	require.NoError(t, err)

	return address
}

func (c *chain) directoryAt(address common.Address) *ethereum.Directory {
	settings := ethereum.Settings{
		RPCURL:    c.url,
		Directory: notary.DirectoryIdentity{ChainID: c.chainID.Uint64(), Address: notary.Address(address)},
		Timeout:   10 * time.Second,
	}

	return ethereum.NewDirectory(c.client, settings)
}

func (c *chain) seal(t *testing.T, ledgerAddress common.Address) {
	t.Helper()

	ctx := context.Background()
	contract := c.binding.Instance(c.client, ledgerAddress)
	head, err := bind.Call(contract, &bind.CallOpts{Context: ctx}, c.binding.PackHead(), c.binding.UnpackHead)
	require.NoError(t, err)
	c.transact(t, contract, c.binding.PackSeal(head))
}

func (c *chain) deploySuccessor(t *testing.T, predecessor common.Address) common.Address {
	t.Helper()

	ctx := context.Background()
	contract := c.binding.Instance(c.client, predecessor)
	options := &bind.CallOpts{Context: ctx}
	next, err := bind.Call(contract, options, c.binding.PackNextEpoch(), c.binding.UnpackNextEpoch)
	require.NoError(t, err)
	head, err := bind.Call(contract, options, c.binding.PackHead(), c.binding.UnpackHead)
	require.NoError(t, err)
	lastAsOf, err := bind.Call(contract, options, c.binding.PackLastAsOf(), c.binding.UnpackLastAsOf)
	require.NoError(t, err)

	input := c.binding.PackConstructor(c.from, big.NewInt(3600), c.from, next, head, lastAsOf)
	address, tx, err := bind.DeployContract(
		c.transactor(),
		common.FromHex(ledger.LiabilityLedgerMetaData.Bin),
		c.client,
		input,
	)
	require.NoError(t, err)
	_, err = bind.WaitDeployed(ctx, c.client, tx.Hash())
	require.NoError(t, err)

	return address
}

func (c *chain) registerSuccessor(
	t *testing.T,
	directoryAddress, ledgerAddress, courtAddress, withholdingAddress common.Address,
) {
	t.Helper()

	binding := directory.NewLedgerDirectory()
	c.transact(
		t,
		binding.Instance(c.client, directoryAddress),
		binding.PackRegisterSuccessor(ledgerAddress, courtAddress, withholdingAddress),
	)
}

func (c *chain) transact(t *testing.T, contract *bind.BoundContract, calldata []byte) {
	t.Helper()

	ctx := context.Background()
	tx, err := bind.Transact(contract, c.transactor(), calldata)
	require.NoError(t, err)
	receipt, err := bind.WaitMined(ctx, c.client, tx.Hash())
	require.NoError(t, err)
	require.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)
}
