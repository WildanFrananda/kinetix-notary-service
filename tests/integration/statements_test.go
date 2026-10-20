package integration_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sort"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	commonv1 "github.com/WildanFrananda/kinetix-notary-service/internal/protogen/common/v1"
	notificationv1 "github.com/WildanFrananda/kinetix-notary-service/internal/protogen/notification/v1"
)

func (e *estate) statementStates(t *testing.T) map[string]string {
	t.Helper()

	rows, err := e.pool.Query(context.Background(), "SELECT holder, epoch_number, state FROM statement_deliveries")
	require.NoError(t, err)

	defer rows.Close()

	states := map[string]string{}

	for rows.Next() {
		var (
			holder string
			epoch  int64
			state  string
		)
		require.NoError(t, rows.Scan(&holder, &epoch, &state))

		states[fmt.Sprintf("%s@%d", holder, epoch)] = state
	}

	require.NoError(t, rows.Err())

	return states
}

func TestEveryHolderOfAFinalEpochIsToldOnceThatTheirStatementIsSigned(t *testing.T) {
	e := newEstate(t)
	e.count = 4
	unreachable := "principal-2"
	e.notified.Answer = func(request *notificationv1.NotifyRequest) (*notificationv1.NotifyResponse, error) {
		if request.GetRecipientPrincipalId() == unreachable {
			return &notificationv1.NotifyResponse{
				Error: &commonv1.ErrorDetail{ErrorCode: "NO_CHANNEL_AVAILABLE", Message: "no usable address"},
			}, nil
		}

		return &notificationv1.NotifyResponse{
			Accepted:       true,
			NotificationId: "NTF-" + request.GetRecipientPrincipalId(),
		}, nil
	}

	r := e.start(t)
	ctx := context.Background()

	r.Statements.Tick(ctx)
	require.Empty(t, e.notified.Requests(), "nothing is said before an epoch is final")

	e.tickUntil(t, r, 0, "mined")
	e.chain.finalize(t)
	e.tickUntil(t, r, 0, "finalized")

	r.Statements.Tick(ctx)

	requests := e.notified.Requests()
	require.Len(t, requests, 4)

	recipients := make([]string, 0, len(requests))
	for _, request := range requests {
		recipients = append(recipients, request.GetRecipientPrincipalId())
		require.Equal(t, notificationv1.Template_TEMPLATE_LIABILITY_STATEMENT, request.GetTemplate())
		require.Equal(t, map[string]string{"epoch": "0", "verify_url": statementVerifyURL}, request.GetParams())
		require.Equal(
			t,
			"liability-statement:0:"+request.GetRecipientPrincipalId(),
			request.GetIdempotencyKey().GetKey(),
		)
	}

	sort.Strings(recipients)
	require.Equal(t, []string{"principal-0", "principal-1", "principal-2", "principal-3"}, recipients)
	require.Equal(t, []string{"notary", "notary", "notary", "notary"}, e.notified.Callers())
	require.Equal(
		t,
		map[string]string{
			"principal-0@0": "sent",
			"principal-1@0": "sent",
			"principal-2@0": "refused",
			"principal-3@0": "sent",
		},
		e.statementStates(t),
	)

	r.Statements.Tick(ctx)
	require.Len(t, e.notified.Requests(), 4, "a statement that was handed over is not sent again")

	e.advance(2 * time.Hour)
	e.tickUntil(t, r, 1, "mined")
	e.chain.finalize(t)
	e.tickUntil(t, r, 1, "finalized")

	r.Statements.Tick(ctx)
	require.Len(t, e.notified.Requests(), 4, "an unchanged balance is not announced again")

	response := httptest.NewRecorder()
	r.Handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/metrics", nil))
	body := response.Body.String()
	require.Contains(t, body, `kinetix_notary_statements_total{outcome="sent"} 3`)
	require.Contains(t, body, `kinetix_notary_statements_total{outcome="refused"} 1`)
	require.Contains(t, body, `kinetix_notary_statements_total{outcome="postponed"} 0`)
	require.Contains(t, body, "kinetix_notary_statements_pending 0")
	require.Contains(t, body, "kinetix_notary_statements_overdue 0")
	require.Contains(
		t,
		body,
		`kinetix_grpc_client_calls_total{grpc_code="OK",grpc_method="/notification.v1.NotificationService/Notify",peer="notification"} 4`,
	)
	require.NotContains(t, body, "principal-")
}

func TestAStatementWaitsForNotificationRatherThanBeingLost(t *testing.T) {
	e := newEstate(t)
	e.count = 2

	var answering atomic.Bool

	e.notified.Answer = func(*notificationv1.NotifyRequest) (*notificationv1.NotifyResponse, error) {
		if !answering.Load() {
			return &notificationv1.NotifyResponse{
				Error: &commonv1.ErrorDetail{ErrorCode: "SENDER_UNAVAILABLE", Message: "email did not answer"},
			}, nil
		}

		return &notificationv1.NotifyResponse{Accepted: true, NotificationId: "NTF-late"}, nil
	}

	r := e.start(t)
	ctx := context.Background()

	e.tickUntil(t, r, 0, "mined")
	e.chain.finalize(t)
	e.tickUntil(t, r, 0, "finalized")

	r.Statements.Tick(ctx)
	require.Equal(t, map[string]string{"principal-0@0": "pending", "principal-1@0": "pending"}, e.statementStates(t))

	answering.Store(true)

	r.Statements.Tick(ctx)
	require.Len(t, e.notified.Requests(), 2, "a postponed statement waits for its next attempt")

	e.advance(2 * time.Minute)
	r.Statements.Tick(ctx)
	require.Equal(t, map[string]string{"principal-0@0": "sent", "principal-1@0": "sent"}, e.statementStates(t))
}
