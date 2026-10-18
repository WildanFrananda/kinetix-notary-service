package evidencing

import (
	"context"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type RegistryCheck struct {
	registry    notary.EvidenceRegistry
	chain       notary.EvidenceChain
	expected    notary.RegistryIdentity
	liabilities notary.Address
}

func NewRegistryCheck(
	registry notary.EvidenceRegistry,
	chain notary.EvidenceChain,
	expected notary.RegistryIdentity,
	liabilities notary.Address,
) *RegistryCheck {
	return &RegistryCheck{registry: registry, chain: chain, expected: expected, liabilities: liabilities}
}

func (c *RegistryCheck) Check(ctx context.Context) error {
	const op = "evidencing.RegistryCheck.Check"

	identity, err := c.registry.Identity(ctx)
	if err != nil {
		return err
	}

	if identity != c.expected {
		return notary.Errf(
			notary.KindWrongLedger,
			op,
			nil,
			"the evidence registry answers as %s on chain %d, this notary is configured for %s on chain %d",
			identity.Address.Hex(),
			identity.ChainID,
			c.expected.Address.Hex(),
			c.expected.ChainID,
		)
	}

	if c.chain.Anchorer() == c.liabilities {
		return notary.Errf(
			notary.KindInvalidArgument,
			op,
			nil,
			"the evidence anchorer %s is the liabilities anchorer; evidence is anchored with a key of its own",
			c.chain.Anchorer().Hex(),
		)
	}

	held, err := c.chain.HoldsAnchorerRole(ctx)
	if err != nil {
		return err
	}

	if !held {
		return notary.Errf(
			notary.KindNotAnchorer,
			op,
			nil,
			"%s does not hold the anchorer role on the evidence registry %s",
			c.chain.Anchorer().Hex(),
			identity.Address.Hex(),
		)
	}

	return nil
}
