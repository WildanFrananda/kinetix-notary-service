package feeds

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/WildanFrananda/kinetix-notary-service/internal/mesh"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

func dial(op, peer string, settings Settings, options []grpc.DialOption) (*grpc.ClientConn, error) {
	if settings.Timeout <= 0 {
		return nil, notary.Errf(notary.KindInvalidArgument, op, nil, "a call to %s needs a positive timeout", peer)
	}

	creds, err := mesh.ClientCredentials(settings.Identity, peer)
	if err != nil {
		return nil, err
	}

	options = append([]grpc.DialOption{grpc.WithTransportCredentials(creds)}, options...)

	conn, err := grpc.NewClient(settings.Endpoint, options...)
	if err != nil {
		return nil, notary.Errf(notary.KindInvalidArgument, op, err, "%s endpoint %q", peer, settings.Endpoint)
	}

	return conn, nil
}

func closeConnection(op string, conn *grpc.ClientConn) error {
	if err := conn.Close(); err != nil {
		return notary.Errf(notary.KindUnknown, op, err, "")
	}

	return nil
}

func classify(op, peer string, err error) error {
	switch status.Code(err) {
	case codes.InvalidArgument, codes.FailedPrecondition, codes.PermissionDenied, codes.Unauthenticated,
		codes.Unimplemented, codes.NotFound, codes.OutOfRange:
		return notary.Errf(notary.KindSourceRefused, op, err, "%s refused the evidence feed", peer)
	default:
		return notary.Errf(notary.KindSourceUnavailable, op, err, "%s could not deliver the evidence feed", peer)
	}
}
