package payment

import (
	"context"
	"errors"
	"io"
	"math/bits"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/WildanFrananda/kinetix-notary-service/internal/mesh"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	paymentv1 "github.com/WildanFrananda/kinetix-notary-service/internal/protogen/payment/v1"
)

const peer = "payment"

type Source struct {
	conn    *grpc.ClientConn
	stub    paymentv1.LiabilityServiceClient
	timeout time.Duration
}

var _ notary.SnapshotSource = (*Source)(nil)

func Dial(settings Settings, options ...grpc.DialOption) (*Source, error) {
	const op = "payment.Dial"

	creds, err := mesh.ClientCredentials(settings.Identity, peer)
	if err != nil {
		return nil, err
	}

	options = append([]grpc.DialOption{grpc.WithTransportCredentials(creds)}, options...)

	conn, err := grpc.NewClient(settings.Endpoint, options...)
	if err != nil {
		return nil, notary.Errf(notary.KindInvalidArgument, op, err, "payment endpoint %q", settings.Endpoint)
	}

	return &Source{conn: conn, stub: paymentv1.NewLiabilityServiceClient(conn), timeout: settings.Timeout}, nil
}

func (s *Source) Close() error {
	if err := s.conn.Close(); err != nil {
		return notary.Errf(notary.KindUnknown, "payment.Source.Close", err, "")
	}

	return nil
}

func (s *Source) Snapshot(ctx context.Context, currency string) (notary.Snapshot, error) {
	const op = "payment.Source.Snapshot"

	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	stream, err := s.stub.StreamLiabilitySnapshot(ctx, &paymentv1.StreamLiabilitySnapshotRequest{Currency: currency})
	if err != nil {
		return notary.Snapshot{}, classify(op, err)
	}

	reader := newSnapshotReader(currency)

	for {
		part, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return reader.finish(op)
		}

		if err != nil {
			return notary.Snapshot{}, classify(op, err)
		}

		if err := reader.accept(op, part); err != nil {
			return notary.Snapshot{}, err
		}
	}
}

func classify(op string, err error) error {
	switch status.Code(err) {
	case codes.InvalidArgument, codes.FailedPrecondition, codes.PermissionDenied, codes.Unauthenticated,
		codes.Unimplemented, codes.NotFound, codes.OutOfRange:
		return notary.Errf(notary.KindSourceRefused, op, err, "payment refused the snapshot")
	default:
		return notary.Errf(notary.KindSourceUnavailable, op, err, "payment could not deliver the snapshot")
	}
}

func refused(op, format string, args ...any) error {
	return notary.Errf(notary.KindSourceRefused, op, nil, format, args...)
}

func addAmount(total, amount uint64) (uint64, bool) {
	sum, carry := bits.Add64(total, amount, 0)
	return sum, carry == 0
}
