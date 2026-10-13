package wiring

import (
	"context"
	"log/slog"

	"go.uber.org/dig"

	"github.com/WildanFrananda/kinetix-notary-service/internal/config"
	"github.com/WildanFrananda/kinetix-notary-service/internal/httpapi"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/proving"
	"github.com/WildanFrananda/kinetix-notary-service/internal/tokens"
)

type proverParams struct {
	dig.In
	Source notary.ProofSource
	Scheme notary.Scheme
}

func provideProver(p proverParams) *proving.Prover {
	return proving.NewProver(p.Source, p.Scheme)
}

type authenticatorParams struct {
	dig.In
	Config config.Daemon
	Log    *slog.Logger
}

type authenticatorResult struct {
	dig.Out
	Authenticator httpapi.Authenticator
}

func provideAuthenticator(p authenticatorParams) (authenticatorResult, error) {
	authenticator, err := tokens.NewAuthenticator(
		context.Background(),
		tokens.Settings{
			JWKSURL:  p.Config.JWKSURL,
			Issuer:   p.Config.JWTIssuer,
			Audience: p.Config.JWTAudience,
		},
		p.Log,
	)
	if err != nil {
		return authenticatorResult{}, err
	}

	return authenticatorResult{Authenticator: authenticator}, nil
}

type proofHandlerParams struct {
	dig.In
	Authenticator httpapi.Authenticator
	Prover        *proving.Prover
	Scheme        notary.Scheme
	Log           *slog.Logger
}

func provideProofHandler(p proofHandlerParams) *httpapi.ProofHandler {
	return httpapi.NewProofHandler(p.Authenticator, p.Prover, p.Scheme, p.Log)
}
