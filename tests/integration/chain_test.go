package integration_test

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/contractgen/ledger"
	"github.com/WildanFrananda/kinetix-notary-service/internal/ethereum"
	"github.com/WildanFrananda/kinetix-notary-service/internal/keccak"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

const rpcURLKey = "KINETIX_NOTARY_TEST_RPC_URL"

const finalityLag = 64

type chain struct {
	url     string
	rpc     *rpc.Client
	client  *ethclient.Client
	key     *ecdsa.PrivateKey
	from    common.Address
	chainID *big.Int
	binding *ledger.LiabilityLedger
}

func connect(t *testing.T) *chain {
	t.Helper()

	url := os.Getenv(rpcURLKey)
	if url == "" {
		t.Skip(rpcURLKey + " is not set; these tests need a node, and CI must set it")
	}

	ctx := context.Background()
	rpcClient, err := rpc.DialContext(ctx, url)
	require.NoError(t, err)
	t.Cleanup(rpcClient.Close)
	client := ethclient.NewClient(rpcClient)
	key, err := crypto.GenerateKey()
	require.NoError(t, err)

	from := crypto.PubkeyToAddress(key.PublicKey)
	funds := new(big.Int).Mul(big.NewInt(100), big.NewInt(1e18))
	require.NoError(t, rpcClient.CallContext(ctx, nil, "anvil_setBalance", from, hexutil.EncodeBig(funds)))
	chainID, err := client.ChainID(ctx)
	require.NoError(t, err)

	return &chain{
		url:     url,
		rpc:     rpcClient,
		client:  client,
		key:     key,
		from:    from,
		chainID: chainID,
		binding: ledger.NewLiabilityLedger(),
	}
}

func (c *chain) transactor() *bind.TransactOpts {
	return bind.NewKeyedTransactor(c.key, c.chainID)
}

func (c *chain) deploy(t *testing.T) common.Address {
	t.Helper()

	input := c.binding.PackConstructor(c.from, big.NewInt(3600), c.from, 0, [32]byte{}, 0)
	address, tx, err := bind.DeployContract(
		c.transactor(),
		common.FromHex(ledger.LiabilityLedgerMetaData.Bin),
		c.client,
		input,
	)
	require.NoError(t, err)
	_, err = bind.WaitDeployed(context.Background(), c.client, tx.Hash())
	require.NoError(t, err)

	return address
}

func (c *chain) deployRaw(t *testing.T, initCode string) common.Address {
	t.Helper()

	address, tx, err := bind.DeployContract(c.transactor(), common.FromHex(initCode), c.client, nil)
	require.NoError(t, err)
	_, err = bind.WaitDeployed(context.Background(), c.client, tx.Hash())
	require.NoError(t, err)

	return address
}

func (c *chain) anchor(t *testing.T, address common.Address, tree *notary.Tree, previous notary.Digest) notary.Digest {
	t.Helper()

	ctx := context.Background()
	require.NoError(t, c.rpc.CallContext(ctx, nil, "evm_increaseTime", 30))
	require.NoError(t, c.rpc.CallContext(ctx, nil, "evm_mine"))
	head, err := c.client.HeaderByNumber(ctx, nil)
	require.NoError(t, err)
	asOf, err := notary.UnixTimeFromSeconds(head.Time - 1)
	require.NoError(t, err)
	header, err := tree.Header(previous, asOf)
	require.NoError(t, err)

	packed := c.binding.PackAnchor(ledger.EpochHeader{
		Number:         header.Number,
		PreviousDigest: header.PreviousDigest,
		Root:           header.Root,
		Total:          header.Total,
		LeafCount:      header.LeafCount,
		AsOf:           header.AsOf.Seconds(),
	})
	tx, err := bind.Transact(c.binding.Instance(c.client, address), c.transactor(), packed)
	require.NoError(t, err)
	receipt, err := bind.WaitMined(ctx, c.client, tx.Hash())
	require.NoError(t, err)
	require.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)

	return keccak.NewScheme().EpochDigest(header)
}

func (c *chain) finalize(t *testing.T) {
	t.Helper()
	require.NoError(t, c.rpc.CallContext(context.Background(), nil, "anvil_mine", hexutil.EncodeUint64(finalityLag+1)))
}

func (c *chain) identity(address common.Address) notary.LedgerIdentity {
	return notary.LedgerIdentity{ChainID: c.chainID.Uint64(), Address: notary.Address(address)}
}

func (c *chain) ledgerAt(identity notary.LedgerIdentity) *ethereum.Ledger {
	settings := ethereum.Settings{RPCURL: c.url, Ledger: identity, Timeout: 10 * time.Second}
	return ethereum.NewLedger(c.client, settings, keccak.NewScheme())
}
