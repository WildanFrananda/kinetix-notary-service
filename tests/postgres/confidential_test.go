package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/keccak"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/poseidon"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

func confidentialPlan(t *testing.T, number uint64, previous notary.Digest, accounts int) notary.EpochPlan {
	t.Helper()

	snapshot := notarytest.Snapshot(t, 1760000000+number, notarytest.Liabilities(t, accounts))
	p, err := notary.PlanConfidentialEpoch(
		keccak.NewScheme(),
		poseidon.NewScheme(),
		4,
		&notarytest.Salts{Seed: number},
		&notarytest.Secrets{},
		&notarytest.Blindings{Seed: number},
		nil,
		number,
		previous,
		snapshot,
	)
	require.NoError(t, err)

	p.Confidential.Proof = []byte{0xde, 0xad, byte(number)}

	return p
}

func TestAConfidentialEpochKeepsItsStatementBlindingAndProof(t *testing.T) {
	d := open(t)
	ctx := context.Background()
	lease := d.lease(t, "writer")
	p := confidentialPlan(t, 0, notary.Digest{}, 6)

	saved, err := d.store.SaveBuilt(ctx, lease, ledgerID, p)
	require.NoError(t, err)
	require.Equal(t, p.Confidential, saved.Confidential)

	read, found, err := d.store.InFlight(ctx)
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, uint64(0), read.Header.Total)
	require.Equal(t, p.Digest, read.Digest)
	require.Equal(t, *p.Confidential, *read.Confidential)
	require.Equal(t, int64(1), d.count(t, "SELECT count(*) FROM epochs WHERE scheme = 2"))
}

func TestASchemeOneEpochReadsBackWithoutAConfidentialPart(t *testing.T) {
	d := open(t)
	ctx := context.Background()
	_, err := d.store.SaveBuilt(ctx, d.lease(t, "writer"), ledgerID, plan(t, 0, notary.Digest{}, 3))
	require.NoError(t, err)

	read, _, err := d.store.InFlight(ctx)
	require.NoError(t, err)
	require.Nil(t, read.Confidential)
}

func TestAProofAttemptLeavesTheEpochBuiltAndItsCostCountsWithTheAnchor(t *testing.T) {
	d := open(t)
	ctx := context.Background()
	lease := d.lease(t, "writer")
	_, err := d.store.SaveBuilt(ctx, lease, ledgerID, confidentialPlan(t, 0, notary.Digest{}, 6))
	require.NoError(t, err)

	sent := time.Date(2026, 10, 20, 9, 0, 0, 0, time.UTC)
	proof := attemptOf(notary.AttemptForProof, notary.Digest{1}, 7, sent)
	require.NoError(t, d.store.RecordAttempt(ctx, lease, proof))

	inFlight, _, err := d.store.InFlight(ctx)
	require.NoError(t, err)
	require.Equal(t, notary.StateBuilt, inFlight.State)

	anchor := attemptOf(notary.AttemptForAnchor, notary.Digest{2}, 8, sent.Add(time.Minute))
	require.NoError(t, d.store.RecordAttempt(ctx, lease, anchor))

	inFlight, _, err = d.store.InFlight(ctx)
	require.NoError(t, err)
	require.Equal(t, notary.StateSubmitted, inFlight.State)

	attempts, err := d.store.Attempts(ctx, 0)
	require.NoError(t, err)
	require.Len(t, attempts, 2)
	require.Equal(t, notary.AttemptForProof, attempts[0].Kind)
	require.Equal(t, notary.AttemptForAnchor, attempts[1].Kind)

	for _, settled := range []struct {
		attempt notary.Attempt
		gas     uint64
	}{{proof, 390_000}, {anchor, 180_000}} {
		receipt := notary.AnchorReceipt{
			Hash:              settled.attempt.Anchor.Hash,
			Succeeded:         true,
			BlockNumber:       48_400_000,
			GasUsed:           settled.gas,
			EffectiveGasPrice: 1_000,
			L1FeeWei:          5,
		}
		require.NoError(
			t,
			d.store.SettleAttempt(ctx, lease, settled.attempt.Anchor.Hash, notary.AttemptIncluded, &receipt),
		)
	}

	require.NoError(t, d.store.Transition(ctx, lease, 0, notary.StateSubmitted, notary.StateFinalized, ""))

	costs, err := d.store.RecentCosts(ctx, 24)
	require.NoError(t, err)
	require.Equal(t, []uint64{390_000*1_000 + 5 + 180_000*1_000 + 5}, costs)
}

func TestTheDatabaseRefusesAHalfConfidentialEpoch(t *testing.T) {
	d := open(t)
	ctx := context.Background()

	_, err := d.pool.Exec(
		ctx,
		`INSERT INTO epochs (number, state, chain_id, ledger_address, scheme, previous_digest, root, digest, total,
			leaf_count, as_of, currency, tree_root)
		VALUES (0, 'built', 84532, $1, 2, $2, $2, $3, 0, 1, 1760000000, 'IDR', $2)`,
		make([]byte, 20),
		make([]byte, 32),
		append(make([]byte, 31), 1),
	)
	require.ErrorContains(t, err, "epochs_scheme_two_is_whole")

	_, err = d.pool.Exec(
		ctx,
		`INSERT INTO epochs (number, state, chain_id, ledger_address, scheme, previous_digest, root, digest, total,
			leaf_count, as_of, currency, tree_root, total_commitment, blinding, depth, proof)
		VALUES (0, 'built', 84532, $1, 2, $2, $2, $3, 7, 1, 1760000000, 'IDR', $2, $2, $4, 3, $5)`,
		make([]byte, 20),
		make([]byte, 32),
		append(make([]byte, 31), 1),
		make([]byte, 31),
		[]byte{1},
	)
	require.ErrorContains(t, err, "epochs_scheme_two_hides_its_total")
}

func attemptOf(kind notary.AttemptKind, hash notary.Digest, nonce uint64, sent time.Time) notary.Attempt {
	return notary.Attempt{
		Epoch: 0,
		Kind:  kind,
		Anchor: notary.PreparedAnchor{
			Hash:     hash,
			Raw:      []byte{0xaa, byte(nonce)},
			Nonce:    nonce,
			GasLimit: 450_000,
			Fees:     notary.FeeQuote{MaxFeePerGas: 1_000, MaxPriorityFeePerGas: 10},
		},
		SentAt: sent,
	}
}
