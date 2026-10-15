package proving

import (
	"context"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type DirectoryCheck struct {
	directory notary.Directory
	expected  notary.Generation
}

func NewDirectoryCheck(directory notary.Directory, expected notary.Generation) *DirectoryCheck {
	return &DirectoryCheck{directory: directory, expected: expected}
}

func (c *DirectoryCheck) Check(ctx context.Context) error {
	const op = "proving.DirectoryCheck.Check"

	if _, err := c.directory.Identity(ctx); err != nil {
		return err
	}

	current, err := c.directory.Current(ctx)
	if err != nil {
		return err
	}

	if current.ChainID != c.expected.ChainID ||
		current.Ledger != c.expected.Ledger ||
		current.Court != c.expected.Court ||
		current.Withholding != c.expected.Withholding {
		return notary.Errf(
			notary.KindWrongDirectory,
			op,
			nil,
			"the directory's current generation is the ledger %s with the courts %s and %s on chain %d; "+
				"this notary anchors %s with %s and %s on chain %d",
			current.Ledger.Hex(),
			current.Court.Hex(),
			current.Withholding.Hex(),
			current.ChainID,
			c.expected.Ledger.Hex(),
			c.expected.Court.Hex(),
			c.expected.Withholding.Hex(),
			c.expected.ChainID,
		)
	}

	return nil
}
