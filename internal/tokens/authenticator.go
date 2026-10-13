package tokens

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/time/rate"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

const (
	jwksTimeout         = 10 * time.Second
	clockSkew           = 30 * time.Second
	refreshWaitMax      = 100 * time.Millisecond
	defaultRefreshEvery = 10 * time.Second
)

type Authenticator struct {
	keys   keyfunc.Keyfunc
	parser *jwt.Parser
}

func NewAuthenticator(ctx context.Context, settings Settings, log *slog.Logger) (*Authenticator, error) {
	every := settings.RefreshEvery
	if every <= 0 {
		every = defaultRefreshEvery
	}

	keys, err := keyfunc.NewDefaultOverrideCtx(
		ctx,
		[]string{settings.JWKSURL},
		keyfunc.Override{
			HTTPTimeout:       jwksTimeout,
			RateLimitWaitMax:  refreshWaitMax,
			RefreshUnknownKID: rate.NewLimiter(rate.Every(every), 1),
			RefreshErrorHandlerFunc: func(url string) func(context.Context, error) {
				return func(ctx context.Context, err error) {
					log.WarnContext(
						ctx,
						"identity's signing keys could not be refreshed",
						"url",
						url,
						"err",
						err.Error(),
					)
				}
			},
		},
	)
	if err != nil {
		return nil, notary.Errf(
			notary.KindInvalidArgument,
			"tokens.NewAuthenticator",
			err,
			"reading %s",
			settings.JWKSURL,
		)
	}

	return &Authenticator{
		keys: keys,
		parser: jwt.NewParser(
			jwt.WithValidMethods([]string{"RS256"}),
			jwt.WithIssuer(settings.Issuer),
			jwt.WithAudience(settings.Audience),
			jwt.WithExpirationRequired(),
			jwt.WithIssuedAt(),
			jwt.WithLeeway(clockSkew),
		),
	}, nil
}

func (a *Authenticator) Principal(authorization string) (notary.Holder, error) {
	const op = "tokens.Authenticator.Principal"

	scheme, raw, found := strings.Cut(strings.TrimSpace(authorization), " ")
	if !found || !strings.EqualFold(scheme, "Bearer") || strings.TrimSpace(raw) == "" {
		return notary.Holder{}, notary.Errf(notary.KindUnauthenticated, op, nil, "a bearer token is required")
	}

	claims := jwt.MapClaims{}
	if _, err := a.parser.ParseWithClaims(strings.TrimSpace(raw), claims, a.keys.Keyfunc); err != nil {
		return notary.Holder{}, notary.Errf(notary.KindUnauthenticated, op, err, "the token was refused")
	}

	if use, _ := claims["token_use"].(string); use != "access" {
		return notary.Holder{}, notary.Errf(notary.KindUnauthenticated, op, nil, "the token is not an access token")
	}

	subject, err := claims.GetSubject()
	if err != nil || subject == "" {
		return notary.Holder{}, notary.Errf(notary.KindUnauthenticated, op, nil, "the token names no principal")
	}

	holder, err := notary.NewHolder(subject)
	if err != nil {
		return notary.Holder{}, notary.Errf(notary.KindUnauthenticated, op, err, "the token's principal is unusable")
	}

	return holder, nil
}
