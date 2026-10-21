package tokens_test

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/tokens"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

const principal = "9f1d4a3e-1c62-4d0a-9a7b-2f5c8e0b41d7"

func authenticator(t *testing.T, identity *notarytest.Identity) *tokens.Authenticator {
	t.Helper()

	return authenticatorRefreshing(t, identity, 200*time.Millisecond)
}

func authenticatorRefreshing(t *testing.T, identity *notarytest.Identity, every time.Duration) *tokens.Authenticator {
	t.Helper()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	auth, err := tokens.NewAuthenticator(
		ctx,
		tokens.Settings{
			JWKSURL:      identity.URL,
			Issuer:       notarytest.IdentityIssuer,
			Audience:     notarytest.IdentityAudience,
			RefreshEvery: every,
		},
		slog.New(slog.NewTextHandler(io.Discard, nil)),
	)
	require.NoError(t, err)

	return auth
}

func TestAnAccessTokenFromIdentityNamesItsPrincipal(t *testing.T) {
	identity := notarytest.NewIdentity(t)
	auth := authenticator(t, identity)

	holder, err := auth.Principal("Bearer " + identity.Token(t, principal))

	require.NoError(t, err)
	require.Equal(t, principal, holder.String())

	holder, err = auth.Principal("bearer   " + identity.Token(t, principal))
	require.NoError(t, err)
	require.Equal(t, principal, holder.String())
}

func TestEveryOtherTokenIsRefused(t *testing.T) {
	identity := notarytest.NewIdentity(t)
	auth := authenticator(t, identity)
	mutate := func(change func(jwt.MapClaims)) string {
		claims := identity.Claims(principal)
		change(claims)

		return "Bearer " + identity.Sign(t, claims)
	}

	cases := map[string]string{
		"no header":         "",
		"basic":             "Basic dXNlcjpwYXNz",
		"bearer, no token":  "Bearer ",
		"not a jwt":         "Bearer not-a-token",
		"refresh token":     mutate(func(c jwt.MapClaims) { c["token_use"] = "refresh" }),
		"no token use":      mutate(func(c jwt.MapClaims) { delete(c, "token_use") }),
		"another issuer":    mutate(func(c jwt.MapClaims) { c["iss"] = "https://evil.example" }),
		"another audience":  mutate(func(c jwt.MapClaims) { c["aud"] = "someone-else" }),
		"expired":           mutate(func(c jwt.MapClaims) { c["exp"] = time.Now().Add(-time.Hour).Unix() }),
		"no expiry":         mutate(func(c jwt.MapClaims) { delete(c, "exp") }),
		"issued in future":  mutate(func(c jwt.MapClaims) { c["iat"] = time.Now().Add(time.Hour).Unix() }),
		"no subject":        mutate(func(c jwt.MapClaims) { delete(c, "sub") }),
		"unusable subject":  mutate(func(c jwt.MapClaims) { c["sub"] = "two words" }),
		"signed by someone": "Bearer " + identity.Stranger(t, principal),
		"hmac":              "Bearer " + hmacToken(t, identity.Claims(principal)),
		"unsigned":          "Bearer " + unsignedToken(t, identity.Claims(principal)),
	}
	for name, header := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := auth.Principal(header)

			require.Error(t, err)
			require.Equal(t, notary.KindUnauthenticated, notary.KindOf(err))
		})
	}
}

func TestAForgedKeyIDNeverHoldsARequest(t *testing.T) {
	identity := notarytest.NewIdentity(t)
	auth := authenticatorRefreshing(t, identity, 10*time.Second)

	for range 5 {
		started := time.Now()

		_, err := auth.Principal("Bearer " + forgedToken(t, identity, "no-such-key"))

		require.Equal(t, notary.KindUnauthenticated, notary.KindOf(err))
		require.Less(t, time.Since(started), 2*time.Second)
	}
}

func TestIdentityBeingDownAtStartIsAnOutageNotARefusedToken(t *testing.T) {
	identity := notarytest.NewIdentity(t)
	identity.Down.Store(true)
	auth := authenticator(t, identity)

	_, err := auth.Principal("Bearer " + identity.Token(t, principal))
	require.Equal(t, notary.KindIdentityUnavailable, notary.KindOf(err))
	require.Equal(t, "identity's signing keys cannot be read", notary.MessageOf(err))

	_, err = auth.Principal("Bearer not-a-token")
	require.Equal(t, notary.KindUnauthenticated, notary.KindOf(err))

	identity.Down.Store(false)

	require.Eventually(
		t,
		func() bool {
			holder, err := auth.Principal("Bearer " + identity.Token(t, principal))

			return err == nil && holder.String() == principal
		},
		5*time.Second,
		50*time.Millisecond,
	)
}

func TestIdentityGoingDownKeepsTheKeysItServedAndCallsAnUnknownKeyAnOutage(t *testing.T) {
	identity := notarytest.NewIdentity(t)
	auth := authenticator(t, identity)

	_, err := auth.Principal("Bearer " + identity.Token(t, principal))
	require.NoError(t, err)

	identity.Down.Store(true)

	_, err = auth.Principal("Bearer " + forgedToken(t, identity, "rotated-key"))
	require.Equal(t, notary.KindIdentityUnavailable, notary.KindOf(err))

	_, err = auth.Principal("Bearer " + identity.Stranger(t, principal))
	require.Equal(t, notary.KindUnauthenticated, notary.KindOf(err))

	holder, err := auth.Principal("Bearer " + identity.Token(t, principal))
	require.NoError(t, err)
	require.Equal(t, principal, holder.String())

	identity.Down.Store(false)

	require.Eventually(
		t,
		func() bool {
			_, err := auth.Principal("Bearer " + forgedToken(t, identity, "rotated-key"))

			return notary.KindOf(err) == notary.KindUnauthenticated
		},
		5*time.Second,
		50*time.Millisecond,
	)
}

func hmacToken(t *testing.T, claims jwt.MapClaims) string {
	t.Helper()

	signed, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte("guessable"))
	require.NoError(t, err)

	return signed
}

func unsignedToken(t *testing.T, claims jwt.MapClaims) string {
	t.Helper()

	signed, err := jwt.NewWithClaims(jwt.SigningMethodNone, claims).SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	return signed
}

func forgedToken(t *testing.T, identity *notarytest.Identity, kid string) string {
	t.Helper()

	return identity.SignWithKeyID(t, identity.Claims(principal), kid)
}
