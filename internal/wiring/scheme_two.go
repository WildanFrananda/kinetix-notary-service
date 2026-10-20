package wiring

import (
	"context"

	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/dig"

	"github.com/WildanFrananda/kinetix-notary-service/internal/anchoring"
	"github.com/WildanFrananda/kinetix-notary-service/internal/config"
	"github.com/WildanFrananda/kinetix-notary-service/internal/entropy"
	"github.com/WildanFrananda/kinetix-notary-service/internal/ethereum"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/plonk"
	"github.com/WildanFrananda/kinetix-notary-service/internal/proverproc"
)

type schemeTwoParams struct {
	dig.In
	Config       config.Daemon
	Client       *ethclient.Client
	Settings     ethereum.Settings
	Confidential notary.ConfidentialScheme
}

func provideSchemeTwo(p schemeTwoParams) (*anchoring.SchemeTwo, error) {
	const op = "wiring.provideSchemeTwo"

	if p.Config.Scheme != p.Confidential.Version() {
		return nil, nil
	}

	if err := plonk.CheckSRSFile(p.Config.SRSPath, p.Config.SRSSHA256); err != nil {
		return nil, err
	}

	key, err := ethereum.ParseAnchorerKey(p.Config.AnchorerKey.Reveal())
	if err != nil {
		return nil, err
	}

	writer := ethereum.NewProofWriter(p.Client, p.Settings, key)

	ctx, cancel := context.WithTimeout(context.Background(), p.Settings.Timeout)
	defer cancel()

	depth, err := writer.Depth(ctx)
	if err != nil {
		return nil, notary.Errf(notary.KindOf(err), op, err, "reading the ledger's depth")
	}

	return &anchoring.SchemeTwo{
		Scheme:    p.Confidential,
		Depth:     depth,
		Blindings: entropy.NewBlindings(),
		Prover: proverproc.NewClient([]string{
			config.SRSPathKey + "=" + p.Config.SRSPath,
			config.SRSSHA256Key + "=" + p.Config.SRSSHA256,
		}),
		Chain: writer,
	}, nil
}
