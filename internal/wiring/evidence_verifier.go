package wiring

import (
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/dig"

	"github.com/WildanFrananda/kinetix-notary-service/internal/config"
	"github.com/WildanFrananda/kinetix-notary-service/internal/ethereum"
	"github.com/WildanFrananda/kinetix-notary-service/internal/keccak"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/verifying"
)

func EvidenceVerifier(lookup config.Lookup) (*dig.Container, error) {
	container := dig.New()

	providers := []any{
		func() config.Lookup { return lookup },
		provideEvidenceVerifierConfig,
		provideEvidenceVerifierSettings,
		provideChainClient,
		provideEvidenceScheme,
		provideVerifierRegistry,
		provideEvidenceVerifier,
	}
	for _, provider := range providers {
		if err := container.Provide(provider); err != nil {
			return nil, fmt.Errorf("provide %T: %w", provider, err)
		}
	}

	return container, nil
}

type evidenceVerifierConfigParams struct {
	dig.In
	Lookup config.Lookup
}

func provideEvidenceVerifierConfig(p evidenceVerifierConfigParams) (config.EvidenceVerifier, error) {
	return config.LoadEvidenceVerifier(p.Lookup)
}

type evidenceVerifierSettingsParams struct {
	dig.In
	Config config.EvidenceVerifier
}

func provideEvidenceVerifierSettings(p evidenceVerifierSettingsParams) ethereum.Settings {
	return ethereum.Settings{RPCURL: p.Config.RPCURL, Timeout: p.Config.RPCTimeout}
}

type evidenceSchemeResult struct {
	dig.Out
	Evidence notary.EvidenceScheme
}

func provideEvidenceScheme() evidenceSchemeResult {
	return evidenceSchemeResult{Evidence: keccak.NewEvidence()}
}

type verifierRegistryParams struct {
	dig.In
	Client *ethclient.Client
	Config config.EvidenceVerifier
	Scheme notary.EvidenceScheme
}

type verifierRegistryResult struct {
	dig.Out
	Registry notary.EvidenceRegistry
}

func provideVerifierRegistry(p verifierRegistryParams) verifierRegistryResult {
	return verifierRegistryResult{
		Registry: ethereum.NewEvidenceRegistry(p.Client, p.Config.Registry, p.Scheme, p.Config.RPCTimeout),
	}
}

type evidenceVerifierParams struct {
	dig.In
	Registry notary.EvidenceRegistry
	Scheme   notary.EvidenceScheme
}

func provideEvidenceVerifier(p evidenceVerifierParams) *verifying.EvidenceVerifier {
	return verifying.NewEvidenceVerifier(p.Registry, p.Scheme)
}
