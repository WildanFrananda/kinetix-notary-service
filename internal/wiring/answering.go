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

const (
	answerLookbackBlocks = 100_000
	answerScanSpan       = 2_000
)

type answerChainParams struct {
	dig.In
	Client   *ethclient.Client
	Settings ethereum.Settings
	Config   config.Daemon
}

type answerChainResult struct {
	dig.Out
	Chain notary.AnswerChain
}

func provideAnswerChain(p answerChainParams) (answerChainResult, error) {
	key, err := ethereum.ParseAnswererKey(p.Config.AnswererKey.Reveal())
	if err != nil {
		return answerChainResult{}, err
	}

	return answerChainResult{
		Chain: ethereum.NewAnswerWriter(p.Client, p.Settings, key, p.Config.MaxFeePerGasWei),
	}, nil
}

type answererParams struct {
	dig.In
	Court         notary.WithholdingCourt
	Chain         notary.AnswerChain
	Source        notary.ProofSource
	Registrations notary.RegistrationStore
	Scheme        notary.Scheme
	Observer      answering.Observer
	Clock         func() time.Time
	Log           *slog.Logger
	Config        config.Daemon
}

func provideAnswerer(p answererParams) *answering.Answerer {
	return answering.NewAnswerer(
		p.Court,
		p.Chain,
		p.Source,
		p.Registrations,
		p.Scheme,
		p.Observer,
		p.Clock,
		p.Log,
		answering.Policy{
			FromBlock: p.Config.WithholdingFrom,
			Lookback:  answerLookbackBlocks,
			Span:      answerScanSpan,
		},
	)
}
