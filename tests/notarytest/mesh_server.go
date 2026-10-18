package notarytest

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"

	"github.com/WildanFrananda/kinetix-notary-service/internal/mesh"
)

func ServeMesh(t *testing.T, serverDir string, register func(*grpc.Server)) string {
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
	register(server)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	go func() { _ = server.Serve(listener) }()

	t.Cleanup(server.Stop)

	return listener.Addr().String()
}

func callerOf(ctx context.Context) string {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return ""
	}

	info, ok := p.AuthInfo.(credentials.TLSInfo)
	if !ok || len(info.State.PeerCertificates) == 0 {
		return ""
	}

	return mesh.ServiceOf(info.State.PeerCertificates[0], []string{"kinetix.local"})
}
