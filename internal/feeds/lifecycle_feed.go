package feeds

import (
	"context"
	"time"

	"google.golang.org/grpc"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	orderv1 "github.com/WildanFrananda/kinetix-notary-service/internal/protogen/order/v1"
)

const orderPeer = "order"

type LifecycleFeed struct {
	conn    *grpc.ClientConn
	stub    orderv1.OrderLifecycleServiceClient
	timeout time.Duration
}

var _ notary.EvidenceFeed = (*LifecycleFeed)(nil)

func DialLifecycle(settings Settings, options ...grpc.DialOption) (*LifecycleFeed, error) {
	conn, err := dial("feeds.DialLifecycle", orderPeer, settings, options)
	if err != nil {
		return nil, err
	}

	return &LifecycleFeed{
		conn:    conn,
		stub:    orderv1.NewOrderLifecycleServiceClient(conn),
		timeout: settings.Timeout,
	}, nil
}

func (f *LifecycleFeed) Close() error {
	return closeConnection("feeds.LifecycleFeed.Close", f.conn)
}

func (*LifecycleFeed) Stream() notary.EvidenceStream {
	return notary.StreamOrderLifecycle
}

func (f *LifecycleFeed) Since(ctx context.Context, after notary.LogPosition, limit int32) (notary.EvidencePage, error) {
	const op = "feeds.LifecycleFeed.Since"

	ctx, cancel := context.WithTimeout(ctx, f.timeout)
	defer cancel()

	response, err := f.stub.LifecycleEventsSince(
		ctx,
		&orderv1.LifecycleEventsSinceRequest{After: cursorOf(after), Limit: limit},
	)
	if err != nil {
		return notary.EvidencePage{}, classify(op, orderPeer, err)
	}

	records := make([]notary.EvidenceRecord, 0, len(response.GetEvents()))
	for _, event := range response.GetEvents() {
		record, err := lifecycleRecordOf(op, event)
		if err != nil {
			return notary.EvidencePage{}, err
		}

		records = append(records, record)
	}

	return pageOf(op, orderPeer, records, response.GetNext(), response.GetHasMore())
}

func lifecycleRecordOf(op string, event *orderv1.LifecycleEvent) (notary.LifecycleRecord, error) {
	at, err := positionOf(op, orderPeer, "a lifecycle event", event.GetPosition())
	if err != nil {
		return notary.LifecycleRecord{}, err
	}

	if err := texts(
		op,
		orderPeer,
		field{name: "an order number", value: event.GetOrderNumber()},
		field{name: "a return number", value: event.GetReturnNumber(), optional: true},
		field{name: "an actor", value: event.GetActorPrincipalId(), optional: true},
		field{name: "a buyer", value: event.GetBuyerPrincipalId()},
		field{name: "a merchant", value: event.GetMerchantPrincipalId()},
	); err != nil {
		return notary.LifecycleRecord{}, err
	}

	kind, err := enumOf(op, orderPeer, "event kind", int32(event.GetKind()), true)
	if err != nil {
		return notary.LifecycleRecord{}, err
	}

	actor, err := enumOf(op, orderPeer, "actor kind", int32(event.GetActorKind()), true)
	if err != nil {
		return notary.LifecycleRecord{}, err
	}

	from, err := enumOf(op, orderPeer, "from status", int32(event.GetFromStatus()), false)
	if err != nil {
		return notary.LifecycleRecord{}, err
	}

	to, err := enumOf(op, orderPeer, "to status", int32(event.GetToStatus()), false)
	if err != nil {
		return notary.LifecycleRecord{}, err
	}

	returnStatus, err := enumOf(op, orderPeer, "return status", int32(event.GetReturnStatus()), false)
	if err != nil {
		return notary.LifecycleRecord{}, err
	}

	occurredAt, err := timeOf(op, orderPeer, "an occurrence time", event.GetOccurredAt())
	if err != nil {
		return notary.LifecycleRecord{}, err
	}

	return notary.LifecycleRecord{
		At:                  at,
		OrderNumber:         event.GetOrderNumber(),
		Kind:                kind,
		FromStatus:          from,
		ToStatus:            to,
		ReturnNumber:        event.GetReturnNumber(),
		ReturnStatus:        returnStatus,
		ActorKind:           actor,
		ActorPrincipalID:    event.GetActorPrincipalId(),
		BuyerPrincipalID:    event.GetBuyerPrincipalId(),
		MerchantPrincipalID: event.GetMerchantPrincipalId(),
		OccurredAt:          occurredAt,
	}, nil
}
