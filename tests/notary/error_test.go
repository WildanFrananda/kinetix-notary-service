package notary_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

func TestAnErrorReadsTheSameInEveryShape(t *testing.T) {
	cause := errors.New("dial tcp 10.0.0.9:5432")

	require.Equal(
		t,
		"store.Open: store_unavailable",
		notary.Errf(notary.KindStoreUnavailable, "store.Open", nil, "").Error(),
	)
	require.Equal(
		t,
		"store.Open: store_unavailable: dial tcp 10.0.0.9:5432",
		notary.Errf(notary.KindStoreUnavailable, "store.Open", cause, "").Error(),
	)
	require.Equal(
		t,
		"proving.Prover.Proofs: epoch_not_final: epoch 7 is submitted",
		notary.Errf(notary.KindEpochNotFinal, "proving.Prover.Proofs", nil, "epoch %d is submitted", 7).Error(),
	)
	require.Equal(
		t,
		"store.Open: store_unavailable: no connection: dial tcp 10.0.0.9:5432",
		notary.Errf(notary.KindStoreUnavailable, "store.Open", cause, "no connection").Error(),
	)
}

func TestTheCauseStaysReachable(t *testing.T) {
	cause := errors.New("dial tcp 10.0.0.9:5432")
	err := fmt.Errorf("outer: %w", notary.Errf(notary.KindStoreUnavailable, "store.Open", cause, "no connection"))

	require.ErrorIs(t, err, cause)
	require.Equal(t, notary.KindStoreUnavailable, notary.KindOf(err))
}

func TestTheMessageIsTheSentenceWithoutTheOperationOrTheCause(t *testing.T) {
	inner := notary.Errf(
		notary.KindUnauthenticated,
		"tokens.Authenticator.Principal",
		errors.New("token is expired"),
		"the token was refused",
	)
	wrapped := notary.Errf(notary.KindUnauthenticated, "httpapi.ProofHandler", inner, "")

	require.Equal(t, "the token was refused", notary.MessageOf(inner))
	require.Equal(t, "the token was refused", notary.MessageOf(wrapped))
	require.Equal(t, "the token was refused", notary.MessageOf(fmt.Errorf("outer: %w", wrapped)))
	require.Empty(t, notary.MessageOf(notary.Errf(notary.KindStoreUnavailable, "store.Open", nil, "")))
	require.Empty(t, notary.MessageOf(errors.New("not a notary error")))
	require.Empty(t, notary.MessageOf(nil))
}
