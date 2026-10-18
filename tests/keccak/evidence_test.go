package keccak_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/keccak"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

func TestTheEvidencePaddingMatchesCast(t *testing.T) {
	got := keccak.NewEvidence().EvidenceEmpty()

	require.Equal(t, digestOf(t, "0x5fa2358263196dbbf23d1ca7a509451f7a2f64c15837bfbb81298b1e3e24e4fa"), got)
}

func TestAnEvidenceLeafMatchesCast(t *testing.T) {
	got := keccak.NewEvidence().EvidenceLeaf(notary.StreamDelivery, notary.Salt(repeated(0x77)), repeated(0x11))

	require.Equal(t, digestOf(t, "0x78f405e090952ce1bc279e5b3c1c02c06549d859fab8ce551fce0ee85545313a"), got)
}

func TestAnEvidenceNodeMatchesCast(t *testing.T) {
	got := keccak.NewEvidence().EvidenceNode(repeated(0xaa), repeated(0xbb))

	require.Equal(t, digestOf(t, "0xec597859c9074ae5989c8d1d86d0aa5fd09832ccfb20000271bb0de43dfd3517"), got)
}

func TestABatchDigestMatchesCast(t *testing.T) {
	got := keccak.NewEvidence().BatchDigest(notary.EvidenceBatchHeader{
		Stream:         notary.StreamDelivery,
		Number:         5,
		PreviousDigest: repeated(0x22),
		Root:           repeated(0x33),
		LeafCount:      42,
		Through:        notary.LogPosition{Transaction: 777, Sequence: 9},
		ClosedAt:       notarytest.Unix(t, 1_760_000_000),
	})

	require.Equal(t, digestOf(t, "0xffae9f9dd5205a27109ad36341afdc082c963ab22e8085070ed526f576813828"), got)
}

func TestADeliveryRecordMatchesCast(t *testing.T) {
	photo := repeated(0xab)
	got, err := keccak.NewEvidence().RecordDigest(notary.DeliveryRecord{
		At:                notary.LogPosition{Transaction: 745, Sequence: 3},
		OrderNumber:       "ORD-1",
		DriverPrincipalID: "courier-1",
		PhotoSHA256:       photo[:],
		DeliveredAt:       notarytest.EvidenceTime(t, 1_760_000_000, 500),
	})

	require.NoError(t, err)
	require.Equal(t, digestOf(t, "0xa14224c8e9c3726a34a8362c5c9573ffdbffffba07a47a72cc3ddef6f4f3be3a"), got)
}

func TestAReviewRecordMatchesCast(t *testing.T) {
	comment := repeated(0xcd)
	got, err := keccak.NewEvidence().RecordDigest(notary.ReviewRecord{
		At:                  notary.LogPosition{Transaction: 746, Sequence: 1},
		ReviewID:            "6f1d3a2b-0c4e-4b5a-9d8e-7f6a5b4c3d21",
		Kind:                1,
		OrderID:             "order-uuid-1",
		SubjectID:           "TSHIRT-BLK-M",
		AuthorPrincipalID:   "buyer-1",
		MerchantPrincipalID: "merchant-1",
		Rating:              5,
		CommentSHA256:       comment[:],
		CreatedAt:           notarytest.EvidenceTime(t, 1_760_000_100, 0),
	})

	require.NoError(t, err)
	require.Equal(t, digestOf(t, "0x0bc47b2e8f5199e3275fd979228fd9e82331accb4aea6fe85b2ea6caf04d470c"), got)
}

func TestALifecycleRecordMatchesCast(t *testing.T) {
	got, err := keccak.NewEvidence().RecordDigest(notary.LifecycleRecord{
		At:                  notary.LogPosition{Transaction: 744, Sequence: 12},
		OrderNumber:         "ORD-1",
		Kind:                1,
		FromStatus:          2,
		ToStatus:            9,
		ActorKind:           3,
		ActorPrincipalID:    "courier-1",
		BuyerPrincipalID:    "buyer-1",
		MerchantPrincipalID: "merchant-1",
		OccurredAt:          notarytest.EvidenceTime(t, 1_760_000_050, 123_456_000),
	})

	require.NoError(t, err)
	require.Equal(t, digestOf(t, "0xabee0f9954df22be82f404acdcec6538ccb2b48358b341541664f71ad4e629a5"), got)
}

func TestARecordFieldTooLongIsRefusedRatherThanCut(t *testing.T) {
	_, err := keccak.NewEvidence().RecordDigest(notary.DeliveryRecord{
		OrderNumber: string(make([]byte, notary.MaxEvidenceField+1)),
		DeliveredAt: notarytest.EvidenceTime(t, 1_760_000_000, 0),
	})

	require.Equal(t, notary.KindInvalidArgument, notary.KindOf(err))
}

func TestEveryFieldOfARecordIsInItsDigest(t *testing.T) {
	scheme := keccak.NewEvidence()
	base := notarytest.Deliveries(t, 1)[0]
	original, err := scheme.RecordDigest(base)
	require.NoError(t, err)

	changes := map[string]func(notary.DeliveryRecord) notary.DeliveryRecord{
		"transaction": func(r notary.DeliveryRecord) notary.DeliveryRecord { r.At.Transaction++; return r },
		"sequence":    func(r notary.DeliveryRecord) notary.DeliveryRecord { r.At.Sequence++; return r },
		"order":       func(r notary.DeliveryRecord) notary.DeliveryRecord { r.OrderNumber += "x"; return r },
		"driver":      func(r notary.DeliveryRecord) notary.DeliveryRecord { r.DriverPrincipalID += "x"; return r },
		"photo": func(r notary.DeliveryRecord) notary.DeliveryRecord {
			r.PhotoSHA256 = append([]byte(nil), r.PhotoSHA256...)
			r.PhotoSHA256[5] ^= 1

			return r
		},
		"signature": func(r notary.DeliveryRecord) notary.DeliveryRecord { r.SignatureSHA256 = []byte{0}; return r },
		"time": func(r notary.DeliveryRecord) notary.DeliveryRecord {
			r.DeliveredAt = notarytest.EvidenceTime(t, int64(r.DeliveredAt.Seconds()), int32(r.DeliveredAt.Nanos())+1)

			return r
		},
		"field boundary": func(r notary.DeliveryRecord) notary.DeliveryRecord {
			r.OrderNumber, r.DriverPrincipalID = r.OrderNumber+r.DriverPrincipalID[:1], r.DriverPrincipalID[1:]

			return r
		},
	}
	for name, change := range changes {
		changed, err := scheme.RecordDigest(change(base))
		require.NoError(t, err)
		require.NotEqual(t, original, changed, name)
	}
}

func TestTheEvidenceSchemeIsVersionOne(t *testing.T) {
	require.Equal(t, uint16(1), keccak.NewEvidence().EvidenceVersion())
}
