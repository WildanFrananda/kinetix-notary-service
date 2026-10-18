package evidencing_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

func TestAnEvidenceAnchorTheNodeDroppedIsBroadcastAgain(t *testing.T) {
	r := newRig(t)

	r.tick(t)
	require.Len(t, r.chain.Broadcasts, 1)

	r.chain.Forget()
	r.tick(t)

	require.Len(t, r.chain.Broadcasts, 2)
	require.Equal(t, r.chain.Broadcasts[0], r.chain.Broadcasts[1])
	require.Contains(t, r.logs.String(), "dropped by the node")

	r.settle(t)
	require.Equal(t, "finalized", r.state(t, notary.StreamOrderLifecycle, 0))
}

func TestARevertedEvidenceAnchorIsPreparedAfreshAtTheNextNonce(t *testing.T) {
	r := newRig(t)

	r.tick(t)
	r.chain.Sealed = true
	r.chain.Mine()
	r.chain.Sealed = false

	r.advance(time.Minute)
	r.tick(t)

	attempts, err := r.store.BatchAttempts(context.Background(), notary.StreamOrderLifecycle, 0)
	require.NoError(t, err)
	require.Len(t, attempts, 2)
	require.Equal(t, notary.AttemptReverted, attempts[0].Outcome)
	require.Equal(t, attempts[0].Anchor.Nonce+1, attempts[1].Anchor.Nonce)
	require.Contains(t, r.logs.String(), "evidence anchor reverted")

	r.settle(t)
	require.Equal(t, "finalized", r.state(t, notary.StreamOrderLifecycle, 0))
}

func TestANonceSpentByAnotherTransactionStopsTheBatch(t *testing.T) {
	r := newRig(t)

	r.tick(t)
	r.chain.Foreign()
	r.tick(t)

	require.Equal(t, "submitted", r.state(t, notary.StreamOrderLifecycle, 0))
	require.Contains(t, r.logs.String(), "used by a transaction this notary never sent")
}

func TestAnEvidenceBatchThatLeftTheChainIsSentAgain(t *testing.T) {
	r := newRig(t)

	r.tick(t)
	r.chain.Mine()
	r.tick(t)
	require.Equal(t, "mined", r.state(t, notary.StreamOrderLifecycle, 0))

	r.chain.Reorg()
	r.tick(t)

	require.Contains(t, r.logs.String(), "left the chain")
	require.Equal(t, "submitted", r.state(t, notary.StreamOrderLifecycle, 0))

	r.settle(t)
	require.Equal(t, "finalized", r.state(t, notary.StreamOrderLifecycle, 0))
}
