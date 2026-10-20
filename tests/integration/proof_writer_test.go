package integration_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/contractgen/ledgertwo"
	"github.com/WildanFrananda/kinetix-notary-service/internal/contractgen/plonkverifier"
	"github.com/WildanFrananda/kinetix-notary-service/internal/contractgen/testverifier"
	"github.com/WildanFrananda/kinetix-notary-service/internal/ethereum"
	"github.com/WildanFrananda/kinetix-notary-service/internal/keccak"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/plonk"
	"github.com/WildanFrananda/kinetix-notary-service/internal/poseidon"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

const confidentialDepth = 3

type proofRig struct {
	chain  *chain
	writer *ethereum.ProofWriter
	prover *plonk.EpochProver
}

func newProofRig(t *testing.T) *proofRig {
	t.Helper()

	c := connect(t)
	ledger := c.deployConfidentialLedger(t, c.deployTestVerifier(t), confidentialDepth, c.from, 0, notary.Digest{}, 0)
	settings := ethereum.Settings{RPCURL: c.url, Ledger: c.identity(ledger), Timeout: 10 * time.Second}

	return &proofRig{
		chain:  c,
		writer: ethereum.NewProofWriter(c.client, settings, c.key),
		prover: plonk.NewEpochProver(notarytest.TestProver(t, confidentialDepth)),
	}
}

func (c *chain) deployTestVerifier(t *testing.T) common.Address {
	t.Helper()

	return c.deployCode(t, testverifier.TestPlonkVerifierMetaData.Bin)
}

func (c *chain) deployProductionVerifier(t *testing.T) common.Address {
	t.Helper()

	return c.deployCode(t, plonkverifier.PlonkVerifierMetaData.Bin)
}

func (c *chain) deployCode(t *testing.T, code string) common.Address {
	t.Helper()

	address, tx, err := bind.DeployContract(c.transactor(), common.FromHex(code), c.client, nil)
	require.NoError(t, err)
	_, err = bind.WaitDeployed(context.Background(), c.client, tx.Hash())
	require.NoError(t, err)

	return address
}

func (c *chain) deployConfidentialLedger(
	t *testing.T,
	verifier common.Address,
	depth uint8,
	anchorer common.Address,
	genesisEpoch uint64,
	genesisDigest notary.Digest,
	genesisAsOf uint64,
) common.Address {
	t.Helper()

	input := ledgertwo.NewLiabilityLedgerV2().PackConstructor(
		c.from,
		big.NewInt(3600),
		anchorer,
		genesisEpoch,
		genesisDigest,
		genesisAsOf,
		verifier,
		depth,
	)
	address, tx, err := bind.DeployContract(
		c.transactor(),
		common.FromHex(ledgertwo.LiabilityLedgerV2MetaData.Bin),
		c.client,
		input,
	)
	require.NoError(t, err)
	_, err = bind.WaitDeployed(context.Background(), c.client, tx.Hash())
	require.NoError(t, err)

	return address
}

func (r *proofRig) proven(t *testing.T, number uint64, accounts int) notary.EpochPlan {
	t.Helper()

	snapshot := notarytest.Snapshot(t, 1760000000+number, notarytest.Liabilities(t, accounts))
	plan, err := notary.PlanConfidentialEpoch(
		keccak.NewScheme(),
		poseidon.NewScheme(),
		confidentialDepth,
		&notarytest.Salts{Seed: number},
		&notarytest.Secrets{},
		&notarytest.Blindings{Seed: number},
		nil,
		number,
		notary.Digest{},
		snapshot,
	)
	require.NoError(t, err)

	entries := make([]notary.Entry, len(plan.Accounts))
	for i, account := range plan.Accounts {
		entries[i] = notary.Entry{Commitment: account.Commitment, Amount: account.Liability.Amount}
	}

	proof, err := r.prover.Prove(context.Background(), notary.ProofRequest{
		Number:          number,
		Depth:           confidentialDepth,
		Entries:         entries,
		TreeRoot:        plan.Confidential.TreeRoot,
		TotalCommitment: plan.Confidential.TotalCommitment,
		Blinding:        plan.Confidential.Blinding,
	})
	require.NoError(t, err)

	plan.Confidential.Proof = proof

	return plan
}

func (r *proofRig) quote(t *testing.T) notary.FeeQuote {
	t.Helper()

	base, tip, err := r.writer.Fees(context.Background())
	require.NoError(t, err)

	quote, err := notary.QuoteFees(base, tip, ceiling)
	require.NoError(t, err)

	return quote
}

func TestAProofMadeInGoLandsAndTheLedgerRecordsIt(t *testing.T) {
	r := newProofRig(t)
	ctx := context.Background()
	plan := r.proven(t, 0, 5)

	depth, err := r.writer.Depth(ctx)
	require.NoError(t, err)
	require.Equal(t, confidentialDepth, depth)

	proven, err := r.writer.IsProven(ctx, 0, plan.Header.Root, plan.Header.LeafCount)
	require.NoError(t, err)
	require.False(t, proven)

	nonce, err := r.writer.Nonce(ctx)
	require.NoError(t, err)

	prepared, err := r.writer.PrepareProof(ctx, plan.Header, *plan.Confidential, nonce, r.quote(t))
	require.NoError(t, err)
	require.NoError(t, r.writer.Broadcast(ctx, prepared))
	r.chain.waitReceipt(t, prepared.Hash)

	proven, err = r.writer.IsProven(ctx, 0, plan.Header.Root, plan.Header.LeafCount)
	require.NoError(t, err)
	require.True(t, proven)

	nonce, err = r.writer.Nonce(ctx)
	require.NoError(t, err)

	_, err = r.writer.PrepareProof(ctx, plan.Header, *plan.Confidential, nonce, r.quote(t))
	require.Equal(t, notary.KindAlreadyAnchored, notary.KindOf(err))
}

func TestAProofTheChainRefusesIsARefusalNotAnOutage(t *testing.T) {
	r := newProofRig(t)
	ctx := context.Background()
	plan := r.proven(t, 0, 4)
	forged := *plan.Confidential
	forged.Proof = append([]byte(nil), plan.Confidential.Proof...)
	forged.Proof[len(forged.Proof)-1] ^= 0x01

	nonce, err := r.writer.Nonce(ctx)
	require.NoError(t, err)

	_, err = r.writer.PrepareProof(ctx, plan.Header, forged, nonce, r.quote(t))
	require.Equal(t, notary.KindAnchorRefused, notary.KindOf(err))
	require.Contains(t, err.Error(), "ProofRejected")
}

func TestAStuckProofIsResignedForTheSameLedger(t *testing.T) {
	r := newProofRig(t)
	ctx := context.Background()
	plan := r.proven(t, 0, 3)

	nonce, err := r.writer.Nonce(ctx)
	require.NoError(t, err)

	quote := r.quote(t)
	prepared, err := r.writer.PrepareProof(ctx, plan.Header, *plan.Confidential, nonce, quote)
	require.NoError(t, err)

	bumped, raised := quote.Bumped(ceiling)
	require.True(t, raised)

	resigned, err := r.writer.Resign(prepared, bumped)
	require.NoError(t, err)
	require.Equal(t, prepared.Nonce, resigned.Nonce)
	require.Equal(t, prepared.GasLimit, resigned.GasLimit)
	require.NotEqual(t, prepared.Hash, resigned.Hash)
}

func (c *chain) waitReceipt(t *testing.T, hash notary.Digest) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	receipt, err := bind.WaitMined(ctx, c.client, common.Hash(hash))
	require.NoError(t, err)
	require.Equal(t, uint64(1), receipt.Status)
}
