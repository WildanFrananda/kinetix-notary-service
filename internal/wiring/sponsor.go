package wiring

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/dig"

	"github.com/WildanFrananda/kinetix-notary-service/internal/bundler"
	"github.com/WildanFrananda/kinetix-notary-service/internal/config"
	"github.com/WildanFrananda/kinetix-notary-service/internal/ethereum"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/sponsoring"
)

const receiptPoll = 2 * time.Second

func Sponsor(lookup config.Lookup) (*dig.Container, error) {
	container := dig.New()

	providers := []any{
		func() config.Lookup { return lookup },
		provideSponsorConfig,
		provideSponsorSettings,
		provideChainClient,
		provideScheme,
		provideDirectory,
		provideGenerationReaders,
		provideVerifier,
		provideSponsorPorts,
		provideSponsor,
	}
	for _, provider := range providers {
		if err := container.Provide(provider); err != nil {
			return nil, fmt.Errorf("provide %T: %w", provider, err)
		}
	}

	return container, nil
}

type sponsorConfigParams struct {
	dig.In
	Lookup config.Lookup
}

type sponsorConfigResult struct {
	dig.Out
	Sponsor  config.Sponsor
	Verifier config.Verifier
}

func provideSponsorConfig(p sponsorConfigParams) (sponsorConfigResult, error) {
	cfg, err := config.LoadSponsor(p.Lookup)
	if err != nil {
		return sponsorConfigResult{}, err
	}

	return sponsorConfigResult{Sponsor: cfg, Verifier: cfg.Verifier}, nil
}

type sponsorSettingsParams struct {
	dig.In
	Config config.Verifier
}

func provideSponsorSettings(p sponsorSettingsParams) ethereum.Settings {
	return ethereum.Settings{
		RPCURL:    p.Config.RPCURL,
		Directory: p.Config.Directory,
		Timeout:   p.Config.RPCTimeout,
	}
}

type sponsorPortsParams struct {
	dig.In
	Config config.Sponsor
	Client *ethclient.Client
}

type sponsorPortsResult struct {
	dig.Out
	Chain     notary.SponsorChain
	Paymaster notary.Paymaster
	Bundler   notary.Bundler
	Signer    notary.HolderSigner
	Calls     notary.SponsoredCalls
}

func provideSponsorPorts(p sponsorPortsParams) (sponsorPortsResult, error) {
	signer, err := ethereum.OpenHolderKeystore(p.Config.HolderKeystore, p.Config.HolderPassword)
	if err != nil {
		return sponsorPortsResult{}, err
	}

	timeout := p.Config.Verifier.RPCTimeout

	return sponsorPortsResult{
		Chain:     ethereum.NewSponsorChain(p.Client, timeout),
		Paymaster: ethereum.NewPaymaster(p.Client, p.Config.Verifier.Directory.ChainID, p.Config.Paymaster, timeout),
		Bundler:   bundler.NewClient(p.Config.BundlerURL, timeout),
		Signer:    signer,
		Calls:     ethereum.NewSponsoredCalls(),
	}, nil
}

type sponsorParams struct {
	dig.In
	Config    config.Sponsor
	Chain     notary.SponsorChain
	Paymaster notary.Paymaster
	Bundler   notary.Bundler
	Signer    notary.HolderSigner
	Scheme    notary.OperationScheme
	Calls     notary.SponsoredCalls
}

func provideSponsor(p sponsorParams) *sponsoring.Sponsor {
	return sponsoring.NewSponsor(
		p.Chain,
		p.Paymaster,
		p.Bundler,
		p.Signer,
		p.Scheme,
		p.Calls,
		receiptPoll,
		p.Config.ReceiptWait,
	)
}
