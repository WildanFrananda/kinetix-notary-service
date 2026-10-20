package wiring

import (
	"log/slog"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/dig"

	"github.com/WildanFrananda/kinetix-notary-service/internal/answering"
	"github.com/WildanFrananda/kinetix-notary-service/internal/config"
	"github.com/WildanFrananda/kinetix-notary-service/internal/ethereum"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

const answerLookbackBlocks = 100_000

type answerChainParams struct {
	dig.In
	Client   *ethclient.Client
	Settings ethereum.Settings
	Config   config.Daemon
}

type answerChainResult struct {
	dig.Out
	Writer *ethereum.AnswerWriter
	Chain  notary.AnswerChain
}

func provideAnswerChain(p answerChainParams) (answerChainResult, error) {
	key, err := ethereum.ParseAnswererKey(p.Config.AnswererKey.Reveal())
	if err != nil {
		return answerChainResult{}, err
	}

	writer := ethereum.NewAnswerWriter(p.Client, p.Settings, key, p.Config.MaxFeePerGasWei)

	return answerChainResult{Writer: writer, Chain: writer}, nil
}

type answererParams struct {
	dig.In
	Client        *ethclient.Client
	Settings      ethereum.Settings
	Court         notary.WithholdingCourt
	Chain         notary.AnswerChain
	Writer        *ethereum.AnswerWriter
	Earlier       []notary.Generation `name:"earlier"`
	Source        notary.ProofSource
	Registrations notary.RegistrationStore
	Scheme        notary.Scheme
	Confidential  notary.ConfidentialScheme
	Observer      answering.Observer
	Clock         func() time.Time
	Log           *slog.Logger
	Config        config.Daemon
}

func provideAnswerer(p answererParams) *answering.Answerer {
	venues := []answering.Venue{{Court: p.Court, Chain: p.Chain, FromBlock: p.Config.WithholdingFrom}}

	for _, generation := range p.Earlier {
		venues = append(venues, answering.Venue{
			Court: ethereum.NewWithholdingCourt(p.Client, settingsOf(p.Settings, generation)),
			Chain: p.Writer.At(generation.Withholding),
		})
	}

	return answering.NewAnswerer(
		venues,
		p.Source,
		p.Registrations,
		p.Scheme,
		p.Confidential,
		p.Observer,
		p.Clock,
		p.Log,
		answering.Policy{
			Lookback: answerLookbackBlocks,
			Span:     p.Config.LogRange,
		},
	)
}
