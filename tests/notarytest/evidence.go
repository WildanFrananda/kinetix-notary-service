package notarytest

import (
	"testing"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

func EvidenceTime(t *testing.T, seconds int64, nanos int32) notary.EvidenceTime {
	t.Helper()

	at, err := notary.EvidenceTimeFrom(seconds, nanos)
	if err != nil {
		t.Fatal(err)
	}

	return at
}

func Deliveries(t *testing.T, count int) []notary.DeliveryRecord {
	t.Helper()

	records := make([]notary.DeliveryRecord, count)
	for i := range records {
		photo := make([]byte, 32)
		photo[0] = byte(i)
		photo[31] = 0xab
		records[i] = notary.DeliveryRecord{
			At:                notary.LogPosition{Transaction: 900 + uint64(i/3), Sequence: uint64(i) + 1},
			OrderNumber:       "ORD-FIXTURE-" + string(rune('A'+i%26)),
			DriverPrincipalID: "courier-" + string(rune('a'+i%26)),
			PhotoSHA256:       photo,
			DeliveredAt:       EvidenceTime(t, 1_760_000_000+int64(i), int32(i)*1000),
		}
	}

	return records
}

func Lifecycle(t *testing.T, order, buyer, merchant string, transaction uint64, count int) []notary.LifecycleRecord {
	t.Helper()

	records := make([]notary.LifecycleRecord, count)
	for i := range records {
		records[i] = notary.LifecycleRecord{
			At:                  notary.LogPosition{Transaction: transaction + uint64(i), Sequence: uint64(i) + 1},
			OrderNumber:         order,
			Kind:                1,
			FromStatus:          uint32(i + 1),
			ToStatus:            uint32(i + 2),
			ActorKind:           2,
			ActorPrincipalID:    merchant,
			BuyerPrincipalID:    buyer,
			MerchantPrincipalID: merchant,
			OccurredAt:          EvidenceTime(t, 1_760_000_100+int64(i), int32(i)*100),
		}
	}

	return records
}

func Reviews(t *testing.T, count int) []notary.ReviewRecord {
	t.Helper()

	records := make([]notary.ReviewRecord, count)
	for i := range records {
		comment := make([]byte, 32)
		comment[0] = byte(i + 1)
		records[i] = notary.ReviewRecord{
			At:                  notary.LogPosition{Transaction: 700 + uint64(i), Sequence: uint64(i) + 1},
			ReviewID:            "review-" + string(rune('a'+i%26)),
			Kind:                1,
			OrderID:             "order-" + string(rune('a'+i%26)),
			SubjectID:           "product-" + string(rune('a'+i%26)),
			AuthorPrincipalID:   "buyer-" + string(rune('a'+i%26)),
			MerchantPrincipalID: "merchant-" + string(rune('a'+i%26)),
			Rating:              uint32(i%5 + 1),
			CommentSHA256:       comment,
			CreatedAt:           EvidenceTime(t, 1_760_000_200+int64(i), 0),
		}
	}

	return records
}

func Records[R notary.EvidenceRecord](records []R) []notary.EvidenceRecord {
	out := make([]notary.EvidenceRecord, len(records))
	for i, record := range records {
		out[i] = record
	}

	return out
}
