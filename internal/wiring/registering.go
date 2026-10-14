package wiring

import (
	"log/slog"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/dig"

	"github.com/WildanFrananda/kinetix-notary-service/internal/ethereum"
	"github.com/WildanFrananda/kinetix-notary-service/internal/httpapi"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/registering"
)

type withholdingCourtParams struct {
	dig.In
	Client   *ethclient.Client
	Settings ethereum.Settings
}

type withholdingCourtResult struct {
	dig.Out
	Court     notary.WithholdingCourt
	Identity  notary.CourtIdentity `name:"withholding"`
	Recoverer notary.SignatureRecoverer
}

func provideWithholdingCourt(p withholdingCourtParams) withholdingCourtResult {
	return withholdingCourtResult{
		Court:     ethereum.NewWithholdingCourt(p.Client, p.Settings),
		Identity:  notary.CourtIdentity{ChainID: p.Settings.Ledger.ChainID, Address: p.Settings.Withholding},
		Recoverer: ethereum.NewRecoverer(),
	}
}

type registrarParams struct {
	dig.In
	Source      notary.ProofSource
	Scheme      notary.Scheme
	Withholding notary.WithholdingScheme
	Court       notary.Court
	Chain       notary.WithholdingCourt
	Signer      notary.StatementSigner
	Recoverer   notary.SignatureRecoverer
	Store       notary.RegistrationStore
	Identity    notary.CourtIdentity `name:"withholding"`
	Clock       func() time.Time
}

func provideRegistrar(p registrarParams) *registering.Registrar {
	return registering.NewRegistrar(
		p.Source,
		p.Scheme,
		p.Withholding,
		p.Court,
		p.Chain,
		p.Signer,
		p.Recoverer,
		p.Store,
		p.Identity,
		p.Clock,
	)
}

type registrationHandlerParams struct {
	dig.In
	Authenticator httpapi.Authenticator
	Registrar     *registering.Registrar
	Log           *slog.Logger
}

func provideRegistrationHandler(p registrationHandlerParams) *httpapi.RegistrationHandler {
	return httpapi.NewRegistrationHandler(p.Authenticator, p.Registrar, p.Log)
}
