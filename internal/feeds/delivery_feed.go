package feeds

import (
	"context"
	"time"

	"google.golang.org/grpc"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	fulfillmentv1 "github.com/WildanFrananda/kinetix-notary-service/internal/protogen/fulfillment/v1"
)

const matchingPeer = "matching"

type DeliveryFeed struct {
	conn    *grpc.ClientConn
	stub    fulfillmentv1.DeliveryEvidenceServiceClient
	timeout time.Duration
}

var _ notary.EvidenceFeed = (*DeliveryFeed)(nil)

func DialDeliveries(settings Settings, options ...grpc.DialOption) (*DeliveryFeed, error) {
	conn, err := dial("feeds.DialDeliveries", matchingPeer, settings, options)
	if err != nil {
		return nil, err
	}

	return &DeliveryFeed{
		conn:    conn,
		stub:    fulfillmentv1.NewDeliveryEvidenceServiceClient(conn),
		timeout: settings.Timeout,
	}, nil
}

func (f *DeliveryFeed) Close() error {
	return closeConnection("feeds.DeliveryFeed.Close", f.conn)
}

func (*DeliveryFeed) Stream() notary.EvidenceStream {
	return notary.StreamDelivery
}

func (f *DeliveryFeed) Since(ctx context.Context, after notary.LogPosition, limit int32) (notary.EvidencePage, error) {
	const op = "feeds.DeliveryFeed.Since"

	ctx, cancel := context.WithTimeout(ctx, f.timeout)
	defer cancel()

	response, err := f.stub.DeliveriesSince(
		ctx,
		&fulfillmentv1.DeliveriesSinceRequest{After: cursorOf(after), Limit: limit},
	)
	if err != nil {
		return notary.EvidencePage{}, classify(op, matchingPeer, err)
	}

	records := make([]notary.EvidenceRecord, 0, len(response.GetDeliveries()))
	for _, delivery := range response.GetDeliveries() {
		record, err := deliveryRecordOf(op, delivery)
		if err != nil {
			return notary.EvidencePage{}, err
		}

		records = append(records, record)
	}

	return pageOf(op, matchingPeer, records, response.GetNext(), response.GetHasMore())
}

func deliveryRecordOf(op string, delivery *fulfillmentv1.DeliveryEvidence) (notary.DeliveryRecord, error) {
	at, err := positionOf(op, matchingPeer, "a delivery", delivery.GetPosition())
	if err != nil {
		return notary.DeliveryRecord{}, err
	}

	if err := texts(
		op,
		matchingPeer,
		field{name: "an order number", value: delivery.GetOrderNumber(), optional: true},
		field{name: "a driver", value: delivery.GetDriverPrincipalId()},
	); err != nil {
		return notary.DeliveryRecord{}, err
	}

	photo, err := hashOf(op, matchingPeer, "photo hash", delivery.GetPhotoSha256(), true)
	if err != nil {
		return notary.DeliveryRecord{}, err
	}

	signature, err := hashOf(op, matchingPeer, "signature hash", delivery.GetSignatureSha256(), false)
	if err != nil {
		return notary.DeliveryRecord{}, err
	}

	deliveredAt, err := timeOf(op, matchingPeer, "a delivery time", delivery.GetDeliveredAt())
	if err != nil {
		return notary.DeliveryRecord{}, err
	}

	return notary.DeliveryRecord{
		At:                at,
		OrderNumber:       delivery.GetOrderNumber(),
		DriverPrincipalID: delivery.GetDriverPrincipalId(),
		PhotoSHA256:       photo,
		SignatureSHA256:   signature,
		DeliveredAt:       deliveredAt,
	}, nil
}
