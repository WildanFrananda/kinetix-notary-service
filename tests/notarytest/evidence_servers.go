package notarytest

import (
	"context"
	"sync"

	commonv1 "github.com/WildanFrananda/kinetix-notary-service/internal/protogen/common/v1"
	fulfillmentv1 "github.com/WildanFrananda/kinetix-notary-service/internal/protogen/fulfillment/v1"
	orderv1 "github.com/WildanFrananda/kinetix-notary-service/internal/protogen/order/v1"
	reviewv1 "github.com/WildanFrananda/kinetix-notary-service/internal/protogen/review/v1"
)

type EvidenceServers struct {
	orderv1.UnimplementedOrderLifecycleServiceServer
	fulfillmentv1.UnimplementedDeliveryEvidenceServiceServer
	reviewv1.UnimplementedReviewEvidenceServiceServer

	mu         sync.Mutex
	Events     []*orderv1.LifecycleEvent
	Deliveries []*fulfillmentv1.DeliveryEvidence
	Reviews    []*reviewv1.ReviewEvidence
	Err        error
	Callers    []string
}

func (s *EvidenceServers) AddEvents(events ...*orderv1.LifecycleEvent) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Events = append(s.Events, events...)
}

func (s *EvidenceServers) AddDeliveries(deliveries ...*fulfillmentv1.DeliveryEvidence) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Deliveries = append(s.Deliveries, deliveries...)
}

func (s *EvidenceServers) AddReviews(reviews ...*reviewv1.ReviewEvidence) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Reviews = append(s.Reviews, reviews...)
}

func (s *EvidenceServers) LifecycleEventsSince(
	ctx context.Context,
	request *orderv1.LifecycleEventsSinceRequest,
) (*orderv1.LifecycleEventsSinceResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Callers = append(s.Callers, callerOf(ctx))
	page, next, more := paginate(
		s.Events,
		request.GetAfter(),
		request.GetLimit(),
		(*orderv1.LifecycleEvent).GetPosition,
	)

	return &orderv1.LifecycleEventsSinceResponse{Events: page, Next: next, HasMore: more}, s.Err
}

func (s *EvidenceServers) DeliveriesSince(
	ctx context.Context,
	request *fulfillmentv1.DeliveriesSinceRequest,
) (*fulfillmentv1.DeliveriesSinceResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Callers = append(s.Callers, callerOf(ctx))
	page, next, more := paginate(
		s.Deliveries,
		request.GetAfter(),
		request.GetLimit(),
		(*fulfillmentv1.DeliveryEvidence).GetPosition,
	)

	return &fulfillmentv1.DeliveriesSinceResponse{Deliveries: page, Next: next, HasMore: more}, s.Err
}

func (s *EvidenceServers) ReviewsSince(
	ctx context.Context,
	request *reviewv1.ReviewsSinceRequest,
) (*reviewv1.ReviewsSinceResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Callers = append(s.Callers, callerOf(ctx))
	page, next, more := paginate(
		s.Reviews,
		request.GetAfter(),
		request.GetLimit(),
		(*reviewv1.ReviewEvidence).GetPosition,
	)

	return &reviewv1.ReviewsSinceResponse{Reviews: page, Next: next, HasMore: more}, s.Err
}

func paginate[T any](
	items []T,
	after *commonv1.LogCursor,
	limit int32,
	position func(T) *commonv1.LogCursor,
) ([]T, *commonv1.LogCursor, bool) {
	const (
		defaultLimit = 100
		maxLimit     = 500
	)

	switch {
	case limit <= 0:
		limit = defaultLimit
	case limit > maxLimit:
		limit = maxLimit
	}

	from := &commonv1.LogCursor{TransactionId: after.GetTransactionId(), Sequence: after.GetSequence()}
	page := make([]T, 0, limit)
	next := from
	more := false

	for _, item := range items {
		at := position(item)
		if at != nil && (at.GetTransactionId() < from.GetTransactionId() ||
			(at.GetTransactionId() == from.GetTransactionId() && at.GetSequence() <= from.GetSequence())) {
			continue
		}

		if int32(len(page)) == limit {
			more = true
			break
		}

		page = append(page, item)

		if at != nil {
			next = at
		}
	}

	return page, next, more
}
