package metrics

import (
	"context"
	"errors"
	"io"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (e *edge) ClientInterceptors(peer string) []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithChainUnaryInterceptor(func(
			ctx context.Context,
			method string,
			req, reply any,
			cc *grpc.ClientConn,
			invoker grpc.UnaryInvoker,
			opts ...grpc.CallOption,
		) error {
			err := invoker(ctx, method, req, reply, cc, opts...)
			e.grpcCalls.WithLabelValues(peer, method, status.Code(err).String()).Inc()

			return err
		}),
		grpc.WithChainStreamInterceptor(func(
			ctx context.Context,
			desc *grpc.StreamDesc,
			cc *grpc.ClientConn,
			method string,
			streamer grpc.Streamer,
			opts ...grpc.CallOption,
		) (grpc.ClientStream, error) {
			stream, err := streamer(ctx, desc, cc, method, opts...)
			if err != nil {
				e.grpcCalls.WithLabelValues(peer, method, status.Code(err).String()).Inc()
				return nil, err
			}

			return &countedStream{
				ClientStream: stream,
				single:       !desc.ServerStreams,
				count: func(code codes.Code) {
					e.grpcCalls.WithLabelValues(peer, method, code.String()).Inc()
				},
			}, nil
		}),
	}
}

type countedStream struct {
	grpc.ClientStream
	single bool
	once   sync.Once
	count  func(codes.Code)
}

func (s *countedStream) RecvMsg(message any) error {
	err := s.ClientStream.RecvMsg(message)
	switch {
	case errors.Is(err, io.EOF), err == nil && s.single:
		s.once.Do(func() { s.count(codes.OK) })
	case err != nil:
		s.once.Do(func() { s.count(status.Code(err)) })
	}

	return err
}
