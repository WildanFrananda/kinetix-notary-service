package payment

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/WildanFrananda/kinetix-notary-service/internal/mesh"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	paymentv1 "github.com/WildanFrananda/kinetix-notary-service/internal/protogen/payment/v1"
)

type Records struct {
	conn    *grpc.ClientConn
	stub    paymentv1.LiabilityRecordServiceClient
	timeout time.Duration
}

var _ notary.SnapshotRecords = (*Records)(nil)

func DialRecords(settings Settings, options ...grpc.DialOption) (*Records, error) {
	const op = "payment.DialRecords"

	creds, err := mesh.ClientCredentials(settings.Identity, peer)
	if err != nil {
		return nil, err
	}

	options = append([]grpc.DialOption{grpc.WithTransportCredentials(creds)}, options...)

	conn, err := grpc.NewClient(settings.Endpoint, options...)
	if err != nil {
		return nil, notary.Errf(notary.KindInvalidArgument, op, err, "payment endpoint %q", settings.Endpoint)
	}

	return &Records{conn: conn, stub: paymentv1.NewLiabilityRecordServiceClient(conn), timeout: settings.Timeout}, nil
}

func (r *Records) Close() error {
	if err := r.conn.Close(); err != nil {
		return notary.Errf(notary.KindUnknown, "payment.Records.Close", err, "")
	}

	return nil
}

func (r *Records) Record(ctx context.Context, id notary.SnapshotID) (notary.SnapshotRecord, error) {
	const op = "payment.Records.Record"

	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	answer, err := r.stub.GetLiabilitySnapshotRecord(
		ctx,
		&paymentv1.GetLiabilitySnapshotRecordRequest{SnapshotId: string(id)},
	)
	if err != nil {
		return notary.SnapshotRecord{}, classifyRecord(op, id, err)
	}

	return recordOf(op, id, answer)
}

func recordOf(
	op string,
	id notary.SnapshotID,
	answer *paymentv1.GetLiabilitySnapshotRecordResponse,
) (notary.SnapshotRecord, error) {
	if answer.GetSnapshotId() != string(id) {
		return notary.SnapshotRecord{}, refused(
			op,
			"asked for snapshot %s, payment answered %q",
			id,
			answer.GetSnapshotId(),
		)
	}

	stamp := answer.GetAsOf()
	seconds, positive := unsignedOf(stamp.GetSeconds())
	nanos, inRange := nanosOf(stamp.GetNanos())

	if stamp == nil || !positive || !inRange {
		return notary.SnapshotRecord{}, refused(op, "the record of snapshot %s carries no valid time", id)
	}

	asOf, err := notary.UnixTimeFromSeconds(seconds)
	if err != nil {
		return notary.SnapshotRecord{}, refused(op, "the record of snapshot %s has an unusable time", id)
	}

	total := answer.GetTotal()
	count, countable := unsignedOf(answer.GetAccountCount())
	owed, owing := unsignedOf(total.GetAmountMinor())

	if !countable || !owing || total.GetCurrency() != answer.GetCurrency() {
		return notary.SnapshotRecord{}, refused(op, "the record of snapshot %s has an unusable count or total", id)
	}

	digest, ok := digestOf(answer.GetDigest())
	if !ok {
		return notary.SnapshotRecord{}, refused(op, "the record of snapshot %s carries no 32-byte digest", id)
	}

	return notary.SnapshotRecord{
		ID:       id,
		AsOf:     asOf,
		Nanos:    nanos,
		Currency: answer.GetCurrency(),
		Count:    count,
		Total:    owed,
		Digest:   digest,
	}, nil
}

func classifyRecord(op string, id notary.SnapshotID, err error) error {
	switch status.Code(err) {
	case codes.NotFound:
		return notary.Errf(notary.KindSnapshotUnknown, op, err, "payment never took snapshot %s", id)
	case codes.InvalidArgument, codes.PermissionDenied, codes.Unauthenticated, codes.Unimplemented:
		return notary.Errf(notary.KindSourceRefused, op, err, "payment refused the record of snapshot %s", id)
	default:
		return notary.Errf(
			notary.KindSourceUnavailable,
			op,
			err,
			"payment could not serve the record of snapshot %s",
			id,
		)
	}
}
