package integration_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/ethereum"
	"github.com/WildanFrananda/kinetix-notary-service/internal/keccak"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

const ceiling = 1_000_000_000_000

type writerRig struct {
	chain  *chain
	writer *ethereum.Writer
	ledger notary.LedgerIdentity
}

func newWriterRig(t *testing.T) *writerRig {
	t.Helper()
	c := connect(t)
	identity := c.identity(c.deploy(t))
	settings := ethereum.Settings{RPCURL: c.url, Ledger: identity, Timeout: 10 * time.Second}

	return &writerRig{chain: c, writer: ethereum.NewWriter(c.client, settings, c.key), ledger: identity}
}

func (r *writerRig) plan(t *testing.T, number uint64, previous notary.Digest, asOf uint64) notary.EpochPlan {
	t.Helper()
	snapshot := notarytest.Snapshot(t, asOf, notarytest.Liabilities(t, 21))
	plan, err := notary.PlanEpoch(
		keccak.NewScheme(),
		&notarytest.Salts{Seed: number},
		&notarytest.Secrets{},
		nil,
		number,
		previous,
		snapshot,
	)
	require.NoError(t, err)

	return plan
}

func (r *writerRig) now(t *testing.T) uint64 {
	t.Helper()

	ctx := context.Background()
	require.NoError(t, r.chain.rpc.CallContext(ctx, nil, "evm_increaseTime", 30))
	require.NoError(t, r.chain.rpc.CallContext(ctx, nil, "evm_mine"))
	head, err := r.chain.client.HeaderByNumber(ctx, nil)
	require.NoError(t, err)

	return head.Time - 1
}

func (r *writerRig) mined(t *testing.T, hash notary.Digest) notary.AnchorReceipt {
	t.Helper()

	var receipt notary.AnchorReceipt

	require.Eventually(
		t,
		func() bool {
			found := false

			var err error

			receipt, found, err = r.writer.Receipt(context.Background(), hash)
			require.NoError(t, err)

			return found
		},
		10*time.Second,
		50*time.Millisecond,
	)

	return receipt
}

func (r *writerRig) prepare(t *testing.T, header notary.EpochHeader) (notary.PreparedAnchor, error) {
	t.Helper()

	ctx := context.Background()
	nonce, err := r.writer.Nonce(ctx)
	require.NoError(t, err)
	base, tip, err := r.writer.Fees(ctx)
	require.NoError(t, err)
	quote, err := notary.QuoteFees(base, tip, ceiling)
	require.NoError(t, err)

	return r.writer.Prepare(ctx, header, nonce, quote)
}

func TestTheWriterAnchorsAnEpochAndWatchesItBecomeFinal(t *testing.T) {
	r := newWriterRig(t)
	ctx := context.Background()
	require.NoError(t, r.writer.VerifyAnchorer(ctx))
	head, err := r.writer.Head(ctx)
	require.NoError(t, err)
	require.True(t, head.AtGenesis())
	require.False(t, head.Sealed)

	balanceBefore, err := r.writer.Balance(ctx)
	require.NoError(t, err)

	plan := r.plan(t, head.NextEpoch, head.Head, r.now(t))
	anchor, err := r.prepare(t, plan.Header)
	require.NoError(t, err)
	require.NoError(t, r.writer.Broadcast(ctx, anchor))

	known, err := r.writer.Known(ctx, anchor.Hash)
	require.NoError(t, err)
	require.True(t, known)
	receipt := r.mined(t, anchor.Hash)
	require.True(t, receipt.Succeeded)
	require.Positive(t, receipt.GasUsed)
	require.LessOrEqual(t, receipt.GasUsed, anchor.GasLimit)
	require.Positive(t, receipt.EffectiveGasPrice)
	require.Zero(t, receipt.L1FeeWei)

	latest, found, err := r.writer.EpochAtLatest(ctx, plan.Header.Number)
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, plan.Digest, latest.Digest)
	_, found, err = r.writer.EpochAtFinalized(ctx, plan.Header.Number)
	require.NoError(t, err)
	require.False(t, found)

	r.chain.finalize(t)
	final, found, err := r.writer.EpochAtFinalized(ctx, plan.Header.Number)
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, plan.Digest, final.Digest)

	head, err = r.writer.Head(ctx)
	require.NoError(t, err)
	require.Equal(t, plan.Header.Number+1, head.NextEpoch)
	require.Equal(t, plan.Digest, head.Head)
	require.Equal(t, plan.Header.AsOf.Seconds(), head.LastAsOf)

	balanceAfter, err := r.writer.Balance(ctx)
	require.NoError(t, err)
	cost, err := receipt.CostWei()
	require.NoError(t, err)
	require.Equal(t, new(big.Int).Sub(balanceBefore, balanceAfter).Uint64(), cost)
}

func TestTheLedgersRefusalsAreReadBeforeAnyGasIsSpent(t *testing.T) {
	r := newWriterRig(t)
	ctx := context.Background()
	head, err := r.writer.Head(ctx)
	require.NoError(t, err)
	plan := r.plan(t, head.NextEpoch, head.Head, r.now(t))
	anchor, err := r.prepare(t, plan.Header)
	require.NoError(t, err)
	require.NoError(t, r.writer.Broadcast(ctx, anchor))
	r.mined(t, anchor.Hash)

	_, err = r.prepare(t, plan.Header)
	require.Equal(t, notary.KindAlreadyAnchored, notary.KindOf(err), err)
	_, err = r.writer.Resign(notary.PreparedAnchor{Raw: []byte{0x02, 0x01}}, notary.FeeQuote{MaxFeePerGas: 1})
	require.Equal(t, notary.KindInvalidArgument, notary.KindOf(err))

	conflicting := plan.Header
	conflicting.Root = notary.Digest{0xde, 0xad}
	_, err = r.prepare(t, conflicting)
	require.Equal(t, notary.KindAnchorRefused, notary.KindOf(err))
	require.ErrorContains(t, err, "EpochConflict")

	next := r.plan(t, plan.Header.Number+1, notary.Digest{0x01}, r.now(t))
	_, err = r.prepare(t, next.Header)
	require.Equal(t, notary.KindAnchorRefused, notary.KindOf(err))
	require.ErrorContains(t, err, "PreviousDigestMismatch")

	future := r.plan(t, plan.Header.Number+1, plan.Digest, r.now(t)+3600)
	_, err = r.prepare(t, future.Header)
	require.Equal(t, notary.KindAnchorDeferred, notary.KindOf(err))
	require.ErrorContains(t, err, "SnapshotInFuture")

	err = r.writer.Broadcast(ctx, anchor)
	require.Equal(t, notary.KindAnchorDeferred, notary.KindOf(err))
}

func TestOnlyTheAnchorerMayAnchor(t *testing.T) {
	r := newWriterRig(t)
	ctx := context.Background()
	stranger, err := crypto.GenerateKey()
	require.NoError(t, err)

	funds := hexutil.EncodeBig(new(big.Int).Mul(big.NewInt(10), big.NewInt(1e18)))
	require.NoError(
		t,
		r.chain.rpc.CallContext(ctx, nil, "anvil_setBalance", crypto.PubkeyToAddress(stranger.PublicKey), funds),
	)
	settings := ethereum.Settings{RPCURL: r.chain.url, Ledger: r.ledger, Timeout: 10 * time.Second}
	intruder := ethereum.NewWriter(r.chain.client, settings, stranger)

	err = intruder.VerifyAnchorer(ctx)
	require.Equal(t, notary.KindNotAnchorer, notary.KindOf(err))

	head, err := intruder.Head(ctx)
	require.NoError(t, err)
	plan := r.plan(t, head.NextEpoch, head.Head, r.now(t))
	base, tip, err := intruder.Fees(ctx)
	require.NoError(t, err)
	quote, err := notary.QuoteFees(base, tip, ceiling)
	require.NoError(t, err)
	_, err = intruder.Prepare(ctx, plan.Header, 0, quote)
	require.Equal(t, notary.KindNotAnchorer, notary.KindOf(err))
}

func TestAStuckTransactionIsReplacedAtTheSameNonce(t *testing.T) {
	r := newWriterRig(t)
	ctx := context.Background()
	head, err := r.writer.Head(ctx)
	require.NoError(t, err)
	plan := r.plan(t, head.NextEpoch, head.Head, r.now(t))
	nonce, err := r.writer.Nonce(ctx)
	require.NoError(t, err)
	base, tip, err := r.writer.Fees(ctx)
	require.NoError(t, err)
	quote, err := notary.QuoteFees(base, tip, ceiling)
	require.NoError(t, err)
	require.NoError(t, r.chain.rpc.CallContext(ctx, nil, "evm_setAutomine", false))
	t.Cleanup(func() { _ = r.chain.rpc.CallContext(context.Background(), nil, "evm_setAutomine", true) })

	first, err := r.writer.Prepare(ctx, plan.Header, nonce, quote)
	require.NoError(t, err)
	require.NoError(t, r.writer.Broadcast(ctx, first))

	bumped, raised := quote.Bumped(ceiling)
	require.True(t, raised)

	second, err := r.writer.Resign(first, bumped)
	require.NoError(t, err)
	require.Equal(t, first.GasLimit, second.GasLimit)
	require.NotEqual(t, first.Hash, second.Hash)
	require.NoError(t, r.writer.Broadcast(ctx, second))
	same, err := r.writer.Resign(
		first,
		notary.FeeQuote{MaxFeePerGas: quote.MaxFeePerGas + 1, MaxPriorityFeePerGas: quote.MaxPriorityFeePerGas},
	)
	require.NoError(t, err)
	err = r.writer.Broadcast(ctx, same)
	require.Equal(t, notary.KindAnchorDeferred, notary.KindOf(err))

	require.NoError(t, r.chain.rpc.CallContext(ctx, nil, "evm_mine"))

	receipt := r.mined(t, second.Hash)
	require.True(t, receipt.Succeeded)

	_, found, err := r.writer.Receipt(ctx, first.Hash)
	require.NoError(t, err)
	require.False(t, found)

	known, err := r.writer.Known(ctx, first.Hash)
	require.NoError(t, err)
	require.False(t, known)

	latest, found, err := r.writer.EpochAtLatest(ctx, plan.Header.Number)
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, plan.Digest, latest.Digest)
}

func TestAKeyThatIsNotAPrivateKeyIsRefused(t *testing.T) {
	for _, text := range []string{"", "0x", "1234", "0x" + string(make([]byte, 64)), "0xzz00000000000000000000000000000000000000000000000000000000000000"} {
		_, err := ethereum.ParseAnchorerKey(text)
		require.Equal(t, notary.KindInvalidArgument, notary.KindOf(err), "%q", text)
	}

	key, err := ethereum.ParseAnchorerKey(" 0x4c0883a69102937d6231471b5dbb6204fe5129617082792ae468d01a3f362318 ")
	require.NoError(t, err)
	require.NotNil(t, key)
}
