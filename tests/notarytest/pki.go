package notarytest

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type PKI struct {
	Root   string
	caCert *x509.Certificate
	caKey  *ecdsa.PrivateKey
	caPEM  []byte
	serial int64
}

func NewPKI(t *testing.T) *PKI {
	t.Helper()

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	template := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "notary test root"},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(time.Hour),
		IsCA:                  true,
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageCertSign,
	}
	der, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	require.NoError(t, err)
	cert, err := x509.ParseCertificate(der)
	require.NoError(t, err)

	return &PKI{
		Root:   t.TempDir(),
		caCert: cert,
		caKey:  key,
		caPEM:  pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}),
		serial: 1,
	}
}

func (p *PKI) Issue(t *testing.T, trustDomain, service string) string {
	t.Helper()

	p.serial++
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	id, err := url.Parse("spiffe://" + trustDomain + "/service/" + service)
	require.NoError(t, err)

	template := &x509.Certificate{
		SerialNumber: big.NewInt(p.serial),
		Subject:      pkix.Name{CommonName: service},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		URIs:         []*url.URL{id},
		DNSNames:     []string{"localhost"},
		IPAddresses:  []net.IP{net.ParseIP("127.0.0.1")},
	}
	der, err := x509.CreateCertificate(rand.Reader, template, p.caCert, &key.PublicKey, p.caKey)
	require.NoError(t, err)
	keyDER, err := x509.MarshalPKCS8PrivateKey(key)
	require.NoError(t, err)

	dir := filepath.Join(p.Root, trustDomain, service)
	require.NoError(t, os.MkdirAll(dir, 0o700))
	require.NoError(
		t,
		os.WriteFile(
			filepath.Join(dir, "tls.crt"),
			pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}),
			0o600,
		),
	)
	require.NoError(
		t,
		os.WriteFile(
			filepath.Join(dir, "tls.key"),
			pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: keyDER}),
			0o600,
		),
	)
	require.NoError(t, os.WriteFile(filepath.Join(dir, "ca.pem"), p.caPEM, 0o600))

	return dir
}
