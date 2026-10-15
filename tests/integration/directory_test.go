package integration_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/config"
	"github.com/WildanFrananda/kinetix-notary-service/internal/contractgen/directory"
	"github.com/WildanFrananda/kinetix-notary-service/internal/keccak"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/verifying"
	"github.com/WildanFrananda/kinetix-notary-service/internal/wiring"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

type generations struct {
	anchored
	withholding common.Address
	directory   common.Address
}

func withDirectory(t *testing.T, a anchored) generations {
	t.Helper()

	withholding := a.chain.deployWithholding(t, a.court)

	return generations{
		anchored:    a,
		withholding: withholding,
		directory:   a.chain.deployDirectory(t, a.address, a.court, withholding),
	}
}

func (g generations) verifier(t *testing.T) *verifying.DirectoryVerifier {
	t.Helper()

	env := map[string]string{
		config.RPCURLKey:           g.chain.url,
		config.ChainIDKey:          strconv.FormatUint(g.identity.ChainID, 10),
		config.DirectoryAddressKey: notary.Address(g.directory).Hex(),
		config.RPCTimeoutKey:       "10s",
	}
	container, err := wiring.Verifier(func(key string) (string, bool) {
		value, ok := env[key]
		return value, ok
	})
	require.NoError(t, err)

	var verifier *verifying.DirectoryVerifier

	err = container.Invoke(func(client *ethclient.Client, v *verifying.DirectoryVerifier) {
		t.Cleanup(client.Close)

		verifier = v
	})
	require.NoError(t, err)

	return verifier
}

func TestTheDirectoryAdapterReadsTheDeployedDirectory(t *testing.T) {
	g := withDirectory(t, anchorEpochZero(t, 4, true))
	adapter := g.chain.directoryAt(g.directory)
	ctx := context.Background()

	identity, err := adapter.Identity(ctx)
	require.NoError(t, err)
	require.Equal(t, notary.Address(g.directory), identity.Address)

	current, err := adapter.Current(ctx)
	require.NoError(t, err)
	require.Equal(t, g.identity, current.LedgerIdentity())
	require.Equal(t, notary.Address(g.court), current.Court)
	require.Equal(t, notary.Address(g.withholding), current.Withholding)
	require.Equal(t, uint64(0), current.GenesisEpoch)
	require.Equal(t, uint16(1), current.Scheme)
	require.False(t, current.RegisteredAt.IsZero())

	found, err := adapter.GenerationOf(ctx, 1_000_000)
	require.NoError(t, err)
	require.Equal(t, current, found)
}

func TestAnythingButADirectoryIsTheWrongDirectory(t *testing.T) {
	g := withDirectory(t, anchorEpochZero(t, 3, true))
	ctx := context.Background()

	for name, address := range map[string]common.Address{
		"no contract": common.HexToAddress("0x00000000000000000000000000000000000d1e00"),
		"a ledger":    g.address,
		"a court":     g.court,
	} {
		t.Run(name, func(t *testing.T) {
			_, err := g.chain.directoryAt(address).Identity(ctx)

			require.Equal(t, notary.KindWrongDirectory, notary.KindOf(err), "%v", err)
		})
	}

	_, err := g.chain.directoryAt(common.HexToAddress("0x00000000000000000000000000000000000d1e00")).Current(ctx)
	require.Equal(t, notary.KindWrongDirectory, notary.KindOf(err), "%v", err)
}

func TestTheVerifierFollowsTheDirectoryAcrossASuccessor(t *testing.T) {
	g := withDirectory(t, anchorEpochZero(t, 6, true))
	c := g.chain
	ctx := context.Background()
	before := g.claim(t, 2)

	c.seal(t, g.address)
	successor := c.deploySuccessor(t, g.address)
	key, signer := notarytest.StatementKey(t, "successor-signer")
	successorCourt := c.deployCourt(t, successor, signer)
	successorWithholding := c.deployWithholding(t, successorCourt)
	c.registerSuccessor(t, g.directory, successor, successorCourt, successorWithholding)

	entries := notarytest.Entries(31, 5)
	tree, err := notary.BuildTree(keccak.NewScheme(), 1, entries)
	require.NoError(t, err)

	head, err := c.ledgerAt(g.identity).Epoch(ctx, 0)
	require.NoError(t, err)
	c.anchor(t, successor, tree, head.Digest)
	c.finalize(t)

	proof, err := tree.Prove(3)
	require.NoError(t, err)

	after := notary.Claim{Ledger: c.identity(successor), Epoch: 1, Entry: entries[3], Proof: proof}
	verifier := g.verifier(t)

	held, err := verifier.Verify(ctx, before)
	require.NoError(t, err)
	require.Equal(t, notary.VerdictIncluded, held.Verdict)
	require.Equal(t, g.identity, held.Ledger)

	moved, err := verifier.Verify(ctx, after)
	require.NoError(t, err)
	require.Equal(t, notary.VerdictIncluded, moved.Verdict)
	require.Equal(t, c.identity(successor), moved.Ledger)
	require.Equal(t, notary.Address(successorCourt), moved.Court.Address)

	successorRig := anchored{chain: c, address: successor, court: successorCourt, identity: c.identity(successor)}
	attributed, err := verifier.Verify(ctx, successorRig.sign(t, after, key, signer))
	require.NoError(t, err)
	require.Equal(t, notary.StatementAttributable, attributed.Check)

	misplaced := after
	misplaced.Ledger = g.identity
	_, err = verifier.Verify(ctx, misplaced)
	require.Equal(t, notary.KindWrongLedger, notary.KindOf(err), "%v", err)
}

func TestTheDirectoryRefusesASuccessorThatDoesNotContinueTheSealedHead(t *testing.T) {
	g := withDirectory(t, anchorEpochZero(t, 3, true))
	c := g.chain
	ctx := context.Background()

	c.seal(t, g.address)
	stray := c.deploy(t)
	_, signer := notarytest.StatementKey(t, "stray-signer")
	strayCourt := c.deployCourt(t, stray, signer)
	strayWithholding := c.deployWithholding(t, strayCourt)
	binding := directory.NewLedgerDirectory()

	_, err := bind.Transact(
		binding.Instance(c.client, g.directory),
		c.transactor(),
		binding.PackRegisterSuccessor(stray, strayCourt, strayWithholding),
	)

	var reverted rpc.DataError
	require.ErrorAs(t, err, &reverted)

	data, ok := reverted.ErrorData().(string)
	require.True(t, ok)

	decoded, err := binding.UnpackError(common.FromHex(data))
	require.NoError(t, err)

	mismatch, ok := decoded.(*directory.LedgerDirectoryGenesisMismatch)
	require.True(t, ok, "%T", decoded)
	require.Equal(t, uint64(1), mismatch.ExpectedEpoch)
	require.Equal(t, uint64(0), mismatch.GivenEpoch)

	current, err := c.directoryAt(g.directory).Current(ctx)
	require.NoError(t, err)
	require.Equal(t, g.identity, current.LedgerIdentity())
}
