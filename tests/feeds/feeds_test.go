package feeds_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/WildanFrananda/kinetix-notary-service/internal/feeds"
	"github.com/WildanFrananda/kinetix-notary-service/internal/mesh"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	fulfillmentv1 "github.com/WildanFrananda/kinetix-notary-service/internal/protogen/fulfillment/v1"
	orderv1 "github.com/WildanFrananda/kinetix-notary-service/internal/protogen/order/v1"
	reviewv1 "github.com/WildanFrananda/kinetix-notary-service/internal/protogen/review/v1"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

type rig struct {
	servers    *notarytest.EvidenceServers
	lifecycle  *feeds.LifecycleFeed
	deliveries *feeds.DeliveryFeed
	reviews    *feeds.ReviewFeed
}

func serve(t *testing.T, pki *notarytest.PKI, service string, servers *notarytest.EvidenceServers) string {
	t.Helper()

	return notarytest.ServeMesh(t, pki.Issue(t, "kinetix.local", service), func(s *grpc.Server) {
		orderv1.RegisterOrderLifecycleServiceServer(s, servers)
		fulfillmentv1.RegisterDeliveryEvidenceServiceServer(s, servers)
		reviewv1.RegisterReviewEvidenceServiceServer(s, servers)
	})
}

func newRig(t *testing.T, as map[string]string) *rig {
	t.Helper()

	pki := notarytest.NewPKI(t)
	servers := &notarytest.EvidenceServers{}
	identity := mesh.Identity{PKIDir: pki.Issue(t, "kinetix.local", "notary"), TrustDomains: []string{"kinetix.local"}}
	settings := func(service string) feeds.Settings {
		name := service
		if other, ok := as[service]; ok {
			name = other
		}

		return feeds.Settings{Endpoint: serve(t, pki, name, servers), Identity: identity, Timeout: 5 * time.Second}
	}

	lifecycle, err := feeds.DialLifecycle(settings("order"))
	require.NoError(t, err)
	t.Cleanup(func() { _ = lifecycle.Close() })

	deliveries, err := feeds.DialDeliveries(settings("matching"))
	require.NoError(t, err)
	t.Cleanup(func() { _ = deliveries.Close() })

	reviews, err := feeds.DialReviews(settings("review"))
	require.NoError(t, err)
	t.Cleanup(func() { _ = reviews.Close() })

	return &rig{servers: servers, lifecycle: lifecycle, deliveries: deliveries, reviews: reviews}
}

func TestEachFeedReadsItsOwnersLogAsRecords(t *testing.T) {
	r := newRig(t, nil)
	lifecycle := notarytest.Lifecycle(t, "ORD-1", "buyer-1", "merchant-1", 10, 3)
	lifecycle[2].ActorKind, lifecycle[2].ActorPrincipalID, lifecycle[2].ReturnNumber = 5, "", "RET-1"
	deliveries := notarytest.Deliveries(t, 2)
	reviews := notarytest.Reviews(t, 2)
	reviews[1].CommentSHA256, reviews[1].MerchantPrincipalID = nil, ""

	for _, e := range lifecycle {
		r.servers.AddEvents(notarytest.WireEvent(e))
	}

	for _, d := range deliveries {
		r.servers.AddDeliveries(notarytest.WireDelivery(d))
	}

	for _, v := range reviews {
		r.servers.AddReviews(notarytest.WireReview(v))
	}

	page, err := r.lifecycle.Since(context.Background(), notary.LogPosition{}, 2)
	require.NoError(t, err)
	require.Equal(t, notarytest.Records(lifecycle[:2]), page.Records)
	require.True(t, page.HasMore)
	require.Equal(t, lifecycle[1].At, page.Next)

	page, err = r.lifecycle.Since(context.Background(), page.Next, 2)
	require.NoError(t, err)
	require.Equal(t, notarytest.Records(lifecycle[2:]), page.Records)
	require.False(t, page.HasMore)

	page, err = r.deliveries.Since(context.Background(), notary.LogPosition{}, 10)
	require.NoError(t, err)
	require.Equal(t, notarytest.Records(deliveries), page.Records)

	page, err = r.reviews.Since(context.Background(), notary.LogPosition{}, 10)
	require.NoError(t, err)
	require.Equal(t, notarytest.Records(reviews), page.Records)

	page, err = r.reviews.Since(context.Background(), reviews[1].At, 10)
	require.NoError(t, err)
	require.Empty(t, page.Records)
	require.Equal(t, reviews[1].At, page.Next)

	require.Equal(t, []string{"notary"}, unique(r.servers.Callers))
}

func unique(values []string) []string {
	seen := map[string]bool{}
	out := []string{}

	for _, v := range values {
		if !seen[v] {
			seen[v] = true
			out = append(out, v)
		}
	}

	return out
}

func TestAFeedRefusesARecordTheNotaryCannotCommitTo(t *testing.T) {
	delivery := func(t *testing.T) *fulfillmentv1.DeliveryEvidence {
		t.Helper()
		return notarytest.WireDelivery(notarytest.Deliveries(t, 1)[0])
	}
	event := func(t *testing.T) *orderv1.LifecycleEvent {
		t.Helper()
		return notarytest.WireEvent(notarytest.Lifecycle(t, "ORD-1", "buyer-1", "merchant-1", 10, 1)[0])
	}
	review := func(t *testing.T) *reviewv1.ReviewEvidence {
		t.Helper()
		return notarytest.WireReview(notarytest.Reviews(t, 1)[0])
	}

	for name, c := range map[string]struct {
		event    func(*orderv1.LifecycleEvent)
		delivery func(*fulfillmentv1.DeliveryEvidence)
		review   func(*reviewv1.ReviewEvidence)
	}{
		"an event without a position":   {event: func(e *orderv1.LifecycleEvent) { e.Position = nil }},
		"an event without an order":     {event: func(e *orderv1.LifecycleEvent) { e.OrderNumber = "" }},
		"an event without a buyer":      {event: func(e *orderv1.LifecycleEvent) { e.BuyerPrincipalId = "" }},
		"an event of no kind":           {event: func(e *orderv1.LifecycleEvent) { e.Kind = 0 }},
		"an event by no actor kind":     {event: func(e *orderv1.LifecycleEvent) { e.ActorKind = 0 }},
		"an event without a time":       {event: func(e *orderv1.LifecycleEvent) { e.OccurredAt = nil }},
		"an event with a NUL":           {event: func(e *orderv1.LifecycleEvent) { e.OrderNumber = "ORD\x001" }},
		"an event with a huge field":    {event: func(e *orderv1.LifecycleEvent) { e.ReturnNumber = strings.Repeat("R", 1<<16+1) }},
		"an event with a negative enum": {event: func(e *orderv1.LifecycleEvent) { e.ReturnStatus = -1 }},
		"a delivery without a driver":   {delivery: func(d *fulfillmentv1.DeliveryEvidence) { d.DriverPrincipalId = "" }},
		"a delivery without a photo":    {delivery: func(d *fulfillmentv1.DeliveryEvidence) { d.PhotoSha256 = nil }},
		"a short photo hash":            {delivery: func(d *fulfillmentv1.DeliveryEvidence) { d.PhotoSha256 = []byte{1, 2} }},
		"a short signature hash":        {delivery: func(d *fulfillmentv1.DeliveryEvidence) { d.SignatureSha256 = []byte{1} }},
		"a delivery before the epoch": {
			delivery: func(d *fulfillmentv1.DeliveryEvidence) { d.DeliveredAt = &timestamppb.Timestamp{Seconds: 0} },
		},
		"a review without an author": {review: func(v *reviewv1.ReviewEvidence) { v.AuthorPrincipalId = "" }},
		"a review rated below zero":  {review: func(v *reviewv1.ReviewEvidence) { v.Rating = -3 }},
		"a review rated nothing":     {review: func(v *reviewv1.ReviewEvidence) { v.Rating = 0 }},
		"a review of no kind":        {review: func(v *reviewv1.ReviewEvidence) { v.Kind = 0 }},
		"a short comment hash":       {review: func(v *reviewv1.ReviewEvidence) { v.CommentSha256 = []byte{9} }},
		"a review with bad nanos": {
			review: func(v *reviewv1.ReviewEvidence) { v.CreatedAt = &timestamppb.Timestamp{Seconds: 1, Nanos: 1e9} },
		},
	} {
		t.Run(name, func(t *testing.T) {
			r := newRig(t, nil)

			var err error

			switch {
			case c.event != nil:
				e := event(t)
				c.event(e)
				r.servers.AddEvents(e)
				_, err = r.lifecycle.Since(context.Background(), notary.LogPosition{}, 10)
			case c.delivery != nil:
				d := delivery(t)
				c.delivery(d)
				r.servers.AddDeliveries(d)
				_, err = r.deliveries.Since(context.Background(), notary.LogPosition{}, 10)
			default:
				v := review(t)
				c.review(v)
				r.servers.AddReviews(v)
				_, err = r.reviews.Since(context.Background(), notary.LogPosition{}, 10)
			}

			require.Equal(t, notary.KindSourceRefused, notary.KindOf(err), err)
		})
	}
}

func TestAFeedAnswersOnlyFromTheServiceItNames(t *testing.T) {
	r := newRig(t, map[string]string{"order": "review"})
	r.servers.AddEvents(notarytest.WireEvent(notarytest.Lifecycle(t, "ORD-1", "buyer-1", "merchant-1", 10, 1)[0]))

	_, err := r.lifecycle.Since(context.Background(), notary.LogPosition{}, 10)

	require.Error(t, err)
	require.Empty(t, r.servers.Callers)
}

func TestTheOwnersStatusDecidesRetryOrRefusal(t *testing.T) {
	for code, kind := range map[codes.Code]notary.Kind{
		codes.Unavailable:      notary.KindSourceUnavailable,
		codes.DeadlineExceeded: notary.KindSourceUnavailable,
		codes.PermissionDenied: notary.KindSourceRefused,
		codes.Unimplemented:    notary.KindSourceRefused,
		codes.InvalidArgument:  notary.KindSourceRefused,
	} {
		r := newRig(t, nil)
		r.servers.Err = status.Error(code, "owner says no")

		_, err := r.deliveries.Since(context.Background(), notary.LogPosition{}, 10)
		require.Equal(t, kind, notary.KindOf(err), code.String())
	}
}

func TestAFeedNeedsATimeout(t *testing.T) {
	_, err := feeds.DialReviews(feeds.Settings{Endpoint: "127.0.0.1:1"})
	require.Equal(t, notary.KindInvalidArgument, notary.KindOf(err))
}
