package notarytest

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

const (
	IdentityIssuer   = "https://identity.kinetix.test"
	IdentityAudience = "kinetix"
)

type Identity struct {
	URL      string
	key      *rsa.PrivateKey
	kid      string
	Down     atomic.Bool
	Requests atomic.Int32
}

func NewIdentity(t *testing.T) *Identity {
	t.Helper()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	identity := &Identity{key: key, kid: "identity-test-1"}
	server := httptest.NewServer(http.HandlerFunc(identity.serve))
	t.Cleanup(server.Close)
	identity.URL = server.URL + "/.well-known/jwks.json"

	return identity
}

func (i *Identity) serve(w http.ResponseWriter, _ *http.Request) {
	i.Requests.Add(1)

	if i.Down.Load() {
		http.Error(w, "down", http.StatusServiceUnavailable)

		return
	}

	public := i.key.PublicKey

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"keys": []map[string]string{{
			"kty": "RSA",
			"use": "sig",
			"alg": "RS256",
			"kid": i.kid,
			"n":   base64.RawURLEncoding.EncodeToString(public.N.Bytes()),
			"e":   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(public.E)).Bytes()),
		}},
	})
}

func (i *Identity) Claims(subject string) jwt.MapClaims {
	now := time.Now()

	return jwt.MapClaims{
		"iss":       IdentityIssuer,
		"aud":       IdentityAudience,
		"sub":       subject,
		"iat":       now.Unix(),
		"exp":       now.Add(15 * time.Minute).Unix(),
		"token_use": "access",
		"role":      "customer",
	}
}

func (i *Identity) Sign(t *testing.T, claims jwt.MapClaims) string {
	t.Helper()

	return i.SignWithKeyID(t, claims, i.kid)
}

func (i *Identity) SignWithKeyID(t *testing.T, claims jwt.MapClaims, kid string) string {
	t.Helper()

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = kid

	signed, err := token.SignedString(i.key)
	require.NoError(t, err)

	return signed
}

func (i *Identity) Token(t *testing.T, subject string) string {
	t.Helper()

	return i.Sign(t, i.Claims(subject))
}

func (i *Identity) Stranger(t *testing.T, subject string) string {
	t.Helper()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, i.Claims(subject))
	token.Header["kid"] = i.kid

	signed, err := token.SignedString(key)
	require.NoError(t, err)

	return signed
}
