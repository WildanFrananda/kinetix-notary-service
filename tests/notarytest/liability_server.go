package notarytest

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"github.com/WildanFrananda/kinetix-notary-service/internal/mesh"
	paymentv1 "github.com/WildanFrananda/kinetix-notary-service/internal/protogen/payment/v1"
)

type LiabilityServer struct {
	paymentv1.UnimplementedLiabilityServiceServer
	paymentv1.UnimplementedLiabilityRecordServiceServer
	Parts     []*paymentv1.StreamLiabilitySnapshotResponse
	Script    func() []*paymentv1.StreamLiabilitySnapshotResponse
	Err       error
	Records   map[string]*paymentv1.GetLiabilitySnapshotRecordResponse
	RecordErr error
	Tamper    func(*paymentv1.GetLiabilitySnapshotRecordResponse)
	mu        sync.Mutex
	Callers   []string
	Requests  []string
	Asked     []string
}

func (s *LiabilityServer) GetLiabilitySnapshotRecord(
	_ context.Context,
	request *paymentv1.GetLiabilitySnapshotRecordRequest,
) (*paymentv1.GetLiabilitySnapshotRecordResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Asked = append(s.Asked, request.GetSnapshotId())
	if s.RecordErr != nil {
		return nil, s.RecordErr
	}

	record, ok := s.Records[request.GetSnapshotId()]
	if !ok {
		return nil, status.Error(codes.NotFound, "no such snapshot")
	}

	if s.Tamper != nil {
		altered, ok := proto.Clone(record).(*paymentv1.GetLiabilitySnapshotRecordResponse)
		if !ok {
			return nil, status.Error(codes.Internal, "the record would not clone")
		}

		s.Tamper(altered)

		return altered, nil
	}

	return record, nil
}

func (s *LiabilityServer) Keep(record *paymentv1.GetLiabilitySnapshotRecordResponse) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Records == nil {
		s.Records = map[string]*paymentv1.GetLiabilitySnapshotRecordResponse{}
	}

	s.Records[record.GetSnapshotId()] = record
}

func (s *LiabilityServer) StreamLiabilitySnapshot(
	request *paymentv1.StreamLiabilitySnapshotRequest,
	stream grpc.ServerStreamingServer[paymentv1.StreamLiabilitySnapshotResponse],
) error {
	caller := ""

	if p, ok := peer.FromContext(stream.Context()); ok {
		if info, ok := p.AuthInfo.(credentials.TLSInfo); ok && len(info.State.PeerCertificates) > 0 {
			caller = mesh.ServiceOf(info.State.PeerCertificates[0], []string{"kinetix.local"})
		}
	}

	s.mu.Lock()
	s.Callers = append(s.Callers, caller)
	s.Requests = append(s.Requests, request.GetCurrency())

	parts := s.Parts
	if s.Script != nil {
		parts = s.Script()
	}

	s.remember(parts)
	s.mu.Unlock()

	for _, part := range parts {
		if err := stream.Send(part); err != nil {
			return err
		}
	}

	return s.Err
}

func ServeLiabilities(t *testing.T, serverDir string, service *LiabilityServer) string {
	t.Helper()

	pair, err := tls.LoadX509KeyPair(filepath.Join(serverDir, "tls.crt"), filepath.Join(serverDir, "tls.key"))
	require.NoError(t, err)
	authority, err := os.ReadFile(filepath.Join(serverDir, "ca.pem"))
	require.NoError(t, err)

	clients := x509.NewCertPool()
	require.True(t, clients.AppendCertsFromPEM(authority))
	server := grpc.NewServer(
		grpc.Creds(
			credentials.NewTLS(&tls.Config{
				Certificates: []tls.Certificate{pair},
				ClientCAs:    clients,
				ClientAuth:   tls.RequireAndVerifyClientCert,
				MinVersion:   tls.VersionTLS13,
			}),
		),
	)
	paymentv1.RegisterLiabilityServiceServer(server, service)
	paymentv1.RegisterLiabilityRecordServiceServer(server, service)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	go func() { _ = server.Serve(listener) }()

	t.Cleanup(server.Stop)

	return listener.Addr().String()
}

func (s *LiabilityServer) remember(parts []*paymentv1.StreamLiabilitySnapshotResponse) {
	var (
		header  *paymentv1.LiabilitySnapshotHeader
		trailer *paymentv1.LiabilitySnapshotTrailer
	)

	for _, part := range parts {
		if h := part.GetHeader(); h != nil {
			header = h
		}

		if tr := part.GetTrailer(); tr != nil {
			trailer = tr
		}
	}

	if header == nil || trailer == nil || header.GetSnapshotId() == "" || len(trailer.GetDigest()) != 32 {
		return
	}

	if s.Records == nil {
		s.Records = map[string]*paymentv1.GetLiabilitySnapshotRecordResponse{}
	}

	s.Records[header.GetSnapshotId()] = &paymentv1.GetLiabilitySnapshotRecordResponse{
		SnapshotId:   header.GetSnapshotId(),
		AsOf:         header.GetAsOf(),
		Currency:     header.GetCurrency(),
		AccountCount: trailer.GetAccountCount(),
		Total:        trailer.GetTotal(),
		Digest:       trailer.GetDigest(),
	}
}
