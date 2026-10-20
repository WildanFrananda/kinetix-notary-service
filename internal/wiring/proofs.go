package wiring

import (
	"context"
	"log/slog"

	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/dig"

	"github.com/WildanFrananda/kinetix-notary-service/internal/config"
	"github.com/WildanFrananda/kinetix-notary-service/internal/ethereum"
	"github.com/WildanFrananda/kinetix-notary-service/internal/httpapi"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/proving"
	"github.com/WildanFrananda/kinetix-notary-service/internal/tokens"
	"github.com/WildanFrananda/kinetix-notary-service/internal/watching"
)

type statementSignerParams struct {
	dig.In
	Config config.Daemon
}

type statementSignerResult struct {
	dig.Out
	Signer notary.StatementSigner
}

func provideStatementSigner(p statementSignerParams) (statementSignerResult, error) {
	keys, err := ethereum.ParseStatementKeys(p.Config.StatementKeys.Reveal())
	if err != nil {
		return statementSignerResult{}, err
	}

	return statementSignerResult{Signer: ethereum.NewStatementSigner(keys)}, nil
}

type proverParams struct {
	dig.In
	Client       *ethclient.Client
	Settings     ethereum.Settings
	Earlier      []notary.Generation `name:"earlier"`
	Source       notary.ProofSource
	Scheme       notary.Scheme
	Confidential notary.ConfidentialScheme
	Statements   notary.StatementScheme
	Court        notary.Court
	Signer       notary.StatementSigner
	Identity     notary.CourtIdentity
}

func provideProver(p proverParams) *proving.Prover {
	earlier := make([]proving.Venue, 0, len(p.Earlier))
	for _, generation := range p.Earlier {
		earlier = append(earlier, proving.Venue{
			Ledger:   generation.LedgerIdentity(),
			Court:    ethereum.NewCourt(p.Client, settingsOf(p.Settings, generation)),
			Identity: notary.CourtIdentity{ChainID: generation.ChainID, Address: generation.Court},
		})
	}

	return proving.NewProver(
		p.Source,
		p.Scheme,
		p.Confidential,
		p.Statements,
		p.Signer,
		proving.Venue{Ledger: p.Settings.Ledger, Court: p.Court, Identity: p.Identity},
		earlier,
	)
}

type watcherParams struct {
	dig.In
	Client   *ethclient.Client
	Settings ethereum.Settings
	Earlier  []notary.Generation `name:"earlier"`
	Court    notary.Court
	Observer watching.Observer
	Config   config.Daemon
	Log      *slog.Logger
}

func provideWatcher(p watcherParams) *watching.Watcher {
	courts := []notary.Court{p.Court}
	for _, generation := range p.Earlier {
		courts = append(courts, ethereum.NewCourt(p.Client, settingsOf(p.Settings, generation)))
	}

	return watching.NewWatcher(courts, p.Observer, p.Log, p.Config.Tick)
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
