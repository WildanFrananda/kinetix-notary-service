package notarytest

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	commonv1 "github.com/WildanFrananda/kinetix-notary-service/internal/protogen/common/v1"
	fulfillmentv1 "github.com/WildanFrananda/kinetix-notary-service/internal/protogen/fulfillment/v1"
	orderv1 "github.com/WildanFrananda/kinetix-notary-service/internal/protogen/order/v1"
	reviewv1 "github.com/WildanFrananda/kinetix-notary-service/internal/protogen/review/v1"
)

func cursorOf(position notary.LogPosition) *commonv1.LogCursor {
	return &commonv1.LogCursor{TransactionId: position.Transaction, Sequence: position.Sequence}
}

func stampOf(at notary.EvidenceTime) *timestamppb.Timestamp {
	return timestamppb.New(at.Time())
}

func WireEvent(r notary.LifecycleRecord) *orderv1.LifecycleEvent {
	return &orderv1.LifecycleEvent{
		Position:            cursorOf(r.At),
		OrderNumber:         r.OrderNumber,
		Kind:                orderv1.LifecycleEventKind(r.Kind),
		FromStatus:          commonv1.OrderStatus(r.FromStatus),
		ToStatus:            commonv1.OrderStatus(r.ToStatus),
		ReturnNumber:        r.ReturnNumber,
		ReturnStatus:        orderv1.ReturnStatus(r.ReturnStatus),
		ActorKind:           orderv1.LifecycleActorKind(r.ActorKind),
		ActorPrincipalId:    r.ActorPrincipalID,
		BuyerPrincipalId:    r.BuyerPrincipalID,
		MerchantPrincipalId: r.MerchantPrincipalID,
		OccurredAt:          stampOf(r.OccurredAt),
	}
}

func WireDelivery(r notary.DeliveryRecord) *fulfillmentv1.DeliveryEvidence {
	return &fulfillmentv1.DeliveryEvidence{
		Position:          cursorOf(r.At),
		OrderNumber:       r.OrderNumber,
		DriverPrincipalId: r.DriverPrincipalID,
		PhotoSha256:       r.PhotoSHA256,
		SignatureSha256:   r.SignatureSHA256,
		DeliveredAt:       stampOf(r.DeliveredAt),
	}
}

func WireReview(r notary.ReviewRecord) *reviewv1.ReviewEvidence {
	return &reviewv1.ReviewEvidence{
		Position:            cursorOf(r.At),
		ReviewId:            r.ReviewID,
		Kind:                reviewv1.ReviewKind(r.Kind),
		OrderId:             r.OrderID,
		SubjectId:           r.SubjectID,
		AuthorPrincipalId:   r.AuthorPrincipalID,
		MerchantPrincipalId: r.MerchantPrincipalID,
		Rating:              int32(r.Rating),
		CommentSha256:       r.CommentSHA256,
		CreatedAt:           stampOf(r.CreatedAt),
	}
}
