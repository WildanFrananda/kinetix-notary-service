package mesh

import (
	"crypto/tls"
	"crypto/x509"
	"os"
	"path/filepath"

	"google.golang.org/grpc/credentials"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type Identity struct {
	PKIDir       string
	TrustDomains []string
}

func ClientCredentials(identity Identity, server string) (credentials.TransportCredentials, error) {
	const op = "mesh.ClientCredentials"

	pair, roots, err := load(op, identity.PKIDir)
	if err != nil {
		return nil, err
	}

	domains := append([]string(nil), identity.TrustDomains...)

	return credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{pair},
		RootCAs:      roots,
		MinVersion:   tls.VersionTLS13,
		VerifyConnection: func(state tls.ConnectionState) error {
			if len(state.PeerCertificates) == 0 {
				return notary.Errf(notary.KindSourceRefused, op, nil, "the server presented no certificate")
			}

			if got := ServiceOf(state.PeerCertificates[0], domains); got != server {
				return notary.Errf(notary.KindSourceRefused, op, nil, "dialled %s but the server is %q", server, got)
			}

			return nil
		},
	}), nil
}

func load(op, dir string) (tls.Certificate, *x509.CertPool, error) {
	if dir == "" {
		return tls.Certificate{}, nil, notary.Errf(
			notary.KindInvalidArgument,
			op,
			nil,
			"no PKI directory is configured",
		)
	}

	certificate, err := read(op, dir, "tls.crt")
	if err != nil {
		return tls.Certificate{}, nil, err
	}

	key, err := read(op, dir, "tls.key")
	if err != nil {
		return tls.Certificate{}, nil, err
	}

	authority, err := read(op, dir, "ca.pem")
	if err != nil {
		return tls.Certificate{}, nil, err
	}

	pair, err := tls.X509KeyPair(certificate, key)
	if err != nil {
		return tls.Certificate{}, nil, notary.Errf(
			notary.KindInvalidArgument,
			op,
			err,
			"tls.crt and tls.key in %s are not a pair",
			dir,
		)
	}

	roots := x509.NewCertPool()
	if !roots.AppendCertsFromPEM(authority) {
		return tls.Certificate{}, nil, notary.Errf(
			notary.KindInvalidArgument,
			op,
			nil,
			"%s/ca.pem holds no certificate",
			dir,
		)
	}

	return pair, roots, nil
}

func read(op, dir, name string) ([]byte, error) {
	path := filepath.Join(dir, name)

	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, notary.Errf(notary.KindInvalidArgument, op, err, "reading %s", path)
	}

	if len(data) == 0 {
		return nil, notary.Errf(notary.KindInvalidArgument, op, nil, "%s is empty", path)
	}

	return data, nil
}
