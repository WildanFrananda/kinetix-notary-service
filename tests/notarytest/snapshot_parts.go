package notarytest

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	commonv1 "github.com/WildanFrananda/kinetix-notary-service/internal/protogen/common/v1"
	paymentv1 "github.com/WildanFrananda/kinetix-notary-service/internal/protogen/payment/v1"
)

func SnapshotStream(
	t *testing.T,
	id notary.SnapshotID,
	asOf *timestamppb.Timestamp,
	currency string,
	accounts []*paymentv1.LiabilityAccount,
) ([]*paymentv1.StreamLiabilitySnapshotResponse, *paymentv1.GetLiabilitySnapshotRecordResponse) {
	t.Helper()

	liabilities := make([]notary.Liability, len(accounts))

	var total uint64

	for i, account := range accounts {
		name := account.GetPrincipalId()
		if name == "" {
			name = account.GetSuspensePurpose()
		}

		liabilities[i] = notary.Liability{
			Holder: Holder(t, name),
			Kind:   notary.AccountKind(account.GetKind()),
			Amount: uint64(account.GetBalance().GetAmountMinor()),
		}
		total += liabilities[i].Amount
	}

	record := notary.SnapshotRecord{
		AsOf:     Unix(t, uint64(asOf.GetSeconds())),
		Nanos:    uint32(asOf.GetNanos()),
		Currency: currency,
		Count:    uint64(len(accounts)),
		Total:    total,
	}
	digest, err := notary.SnapshotDigest(record, liabilities)
	require.NoError(t, err)

	money := &commonv1.Money{AmountMinor: int64(total), Currency: currency}

	parts := []*paymentv1.StreamLiabilitySnapshotResponse{
		{Part: &paymentv1.StreamLiabilitySnapshotResponse_Header{Header: &paymentv1.LiabilitySnapshotHeader{
			AsOf:       asOf,
			Currency:   currency,
			SnapshotId: string(id),
		}}},
	}
	if len(accounts) > 0 {
		parts = append(parts, &paymentv1.StreamLiabilitySnapshotResponse{
			Part: &paymentv1.StreamLiabilitySnapshotResponse_Accounts{
				Accounts: &paymentv1.LiabilityAccountBatch{Accounts: accounts},
			},
		})
	}

	parts = append(parts, &paymentv1.StreamLiabilitySnapshotResponse{
		Part: &paymentv1.StreamLiabilitySnapshotResponse_Trailer{Trailer: &paymentv1.LiabilitySnapshotTrailer{
			AccountCount: int64(len(accounts)),
			Total:        money,
			Digest:       digest[:],
		}},
	})
	require.NotEmpty(t, parts)

	return parts, &paymentv1.GetLiabilitySnapshotRecordResponse{
		SnapshotId:   string(id),
		AsOf:         asOf,
		Currency:     currency,
		AccountCount: int64(len(accounts)),
		Total:        money,
		Digest:       digest[:],
	}
}
