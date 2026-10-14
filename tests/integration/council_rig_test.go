package integration_test

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"math/big"
	"sort"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/contractgen/council"
	"github.com/WildanFrananda/kinetix-notary-service/internal/ethereum"
	"github.com/WildanFrananda/kinetix-notary-service/internal/keccak"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

func (c *chain) deployCouncil(t *testing.T, ledgerAddress common.Address, members ...notary.Address) common.Address {
	t.Helper()

	addresses := make([]common.Address, len(members))
	for i, member := range members {
		addresses[i] = common.Address(member)
	}

	sort.Slice(addresses, func(a, b int) bool { return bytes.Compare(addresses[a][:], addresses[b][:]) < 0 })

	input := council.NewAnchorCouncil().
		PackConstructor(c.from, big.NewInt(3600), ledgerAddress, addresses, 2, 7*24*3600)
	address, tx, err := bind.DeployContract(
		c.transactor(),
		common.FromHex(council.AnchorCouncilMetaData.Bin),
		c.client,
		input,
	)
	require.NoError(t, err)
	_, err = bind.WaitDeployed(context.Background(), c.client, tx.Hash())
	require.NoError(t, err)

	return address
}

func (c *chain) handOver(t *testing.T, ledgerAddress, to common.Address) {
	t.Helper()

	ctx := context.Background()
	contract := c.binding.Instance(c.client, ledgerAddress)
	role, err := bind.Call(
		contract,
		&bind.CallOpts{Context: ctx},
		c.binding.PackANCHORERROLE(),
		c.binding.UnpackANCHORERROLE,
	)
	require.NoError(t, err)

	tx, err := bind.Transact(contract, c.transactor(), c.binding.PackGrantRole(role, to))
	require.NoError(t, err)
	receipt, err := bind.WaitMined(ctx, c.client, tx.Hash())
	require.NoError(t, err)
	require.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)
}

func (c *chain) councilAt(identity notary.LedgerIdentity, address common.Address) *ethereum.Council {
	settings := ethereum.Settings{
		RPCURL:  c.url,
		Ledger:  identity,
		Council: notary.Address(address),
		Timeout: 10 * time.Second,
	}

	return ethereum.NewCouncil(c.client, settings)
}

func memberKey(t *testing.T) (*ecdsa.PrivateKey, notary.Address) {
	t.Helper()

	key, err := crypto.GenerateKey()
	require.NoError(t, err)

	return key, notary.Address(crypto.PubkeyToAddress(key.PublicKey))
}

func approvals(
	t *testing.T,
	identity notary.CouncilIdentity,
	ledger notary.Address,
	header notary.EpochHeader,
	keys ...*ecdsa.PrivateKey,
) []notary.Signature {
	t.Helper()

	scheme := keccak.NewScheme()
	digest := scheme.ApprovalDigest(identity, ledger, header.Number, scheme.EpochDigest(header))

	sorted := append([]*ecdsa.PrivateKey(nil), keys...)
	sort.Slice(sorted, func(a, b int) bool {
		first := crypto.PubkeyToAddress(sorted[a].PublicKey)
		second := crypto.PubkeyToAddress(sorted[b].PublicKey)

		return bytes.Compare(first[:], second[:]) < 0
	})

	signatures := make([]notary.Signature, len(sorted))
	for i, key := range sorted {
		signature, err := ethereum.NewApprovalSigner(key).Sign(digest)
		require.NoError(t, err)

		signatures[i] = signature
	}

	return signatures
}
