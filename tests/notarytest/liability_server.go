package notarytest

import (
	"crypto/tls"
	"crypto/x509"
	"net"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"

	"github.com/WildanFrananda/kinetix-notary-service/internal/mesh"
	paymentv1 "github.com/WildanFrananda/kinetix-notary-service/internal/protogen/payment/v1"
)

type LiabilityServer struct {
	paymentv1.UnimplementedLiabilityServiceServer
	Parts    []*paymentv1.StreamLiabilitySnapshotResponse
	Script   func() []*paymentv1.StreamLiabilitySnapshotResponse
	Err      error
	mu       sync.Mutex
	Callers  []string
	Requests []string
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

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	go func() { _ = server.Serve(listener) }()

	t.Cleanup(server.Stop)

	return listener.Addr().String()
}
