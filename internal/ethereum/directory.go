package ethereum

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/WildanFrananda/kinetix-notary-service/internal/contractgen/directory"
	"github.com/WildanFrananda/kinetix-notary-service/internal/contractgen/directoryinterface"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type Directory struct {
	client      *ethclient.Client
	binding     *directory.LedgerDirectory
	contract    *bind.BoundContract
	expected    notary.DirectoryIdentity
	interfaceID [4]byte
	timeout     time.Duration
}

const maxGenerations = 64

var _ notary.Directory = (*Directory)(nil)

func NewDirectory(client *ethclient.Client, settings Settings) *Directory {
	binding := directory.NewLedgerDirectory()

	return &Directory{
		client:      client,
		binding:     binding,
		contract:    binding.Instance(client, common.Address(settings.Directory.Address)),
		expected:    settings.Directory,
		interfaceID: directoryInterfaceID(),
		timeout:     settings.Timeout,
	}
}

func (d *Directory) Identity(ctx context.Context) (notary.DirectoryIdentity, error) {
	const op = "ethereum.Directory.Identity"

	ctx, cancel := context.WithTimeout(ctx, d.timeout)
	defer cancel()

	chainID, err := d.client.ChainID(ctx)
	if err != nil {
		return notary.DirectoryIdentity{}, unavailable(op, err)
	}

	if !chainID.IsUint64() || chainID.Uint64() != d.expected.ChainID {
		return notary.DirectoryIdentity{}, notary.Errf(
			notary.KindWrongDirectory,
			op,
			nil,
			"the RPC endpoint serves chain %s, the directory is configured on chain %d",
			chainID,
			d.expected.ChainID,
		)
	}

	code, err := d.client.CodeAt(ctx, common.Address(d.expected.Address), nil)
	if err != nil {
		return notary.DirectoryIdentity{}, unavailable(op, err)
	}

	if len(code) == 0 {
		return notary.DirectoryIdentity{}, notary.Errf(
			notary.KindWrongDirectory,
			op,
			nil,
			"chain %d has no contract at %s",
			d.expected.ChainID,
			d.expected.Address.Hex(),
		)
	}

	supported, err := probe(
		ctx,
		d.client,
		d.expected.Address,
		d.binding.PackSupportsInterface(d.interfaceID),
		d.binding.UnpackSupportsInterface,
	)
	if err != nil {
		return notary.DirectoryIdentity{}, d.wrap(op, err)
	}

	if !supported {
		return notary.DirectoryIdentity{}, notary.Errf(
			notary.KindWrongDirectory,
			op,
			nil,
			"%s does not announce the ledger directory interface",
			d.expected.Address.Hex(),
		)
	}

	return d.expected, nil
}

func (d *Directory) Current(ctx context.Context) (notary.Generation, error) {
	const op = "ethereum.Directory.Current"

	ctx, cancel := context.WithTimeout(ctx, d.timeout)
	defer cancel()

	current, err := bind.Call(
		d.contract,
		&bind.CallOpts{Context: ctx},
		d.binding.PackCurrent(),
		d.binding.UnpackCurrent,
	)
	if err != nil {
		return notary.Generation{}, d.classify(op, err)
	}

	return d.generation(op, current)
}

func (d *Directory) GenerationOf(ctx context.Context, epoch uint64) (notary.Generation, error) {
	const op = "ethereum.Directory.GenerationOf"

	ctx, cancel := context.WithTimeout(ctx, d.timeout)
	defer cancel()

	found, err := bind.Call(
		d.contract,
		&bind.CallOpts{Context: ctx},
		d.binding.PackGenerationOf(epoch),
		d.binding.UnpackGenerationOf,
	)
	if err != nil {
		return notary.Generation{}, d.classify(op, err)
	}

	return d.generation(op, found.Arg1)
}

func (d *Directory) Lineage(ctx context.Context) ([]notary.Generation, error) {
	const op = "ethereum.Directory.Lineage"

	ctx, cancel := context.WithTimeout(ctx, d.timeout)
	defer cancel()

	opts := &bind.CallOpts{Context: ctx}

	count, err := bind.Call(d.contract, opts, d.binding.PackGenerationCount(), d.binding.UnpackGenerationCount)
	if err != nil {
		return nil, d.classify(op, err)
	}

	if !count.IsUint64() || count.Uint64() == 0 || count.Uint64() > maxGenerations {
		return nil, notary.Errf(
			notary.KindWrongDirectory,
			op,
			nil,
			"%s counts %s generations",
			d.expected.Address.Hex(),
			count,
		)
	}

	lineage := make([]notary.Generation, 0, count.Uint64())

	for i := range count.Uint64() {
		raw, err := bind.Call(
			d.contract,
			opts,
			d.binding.PackGeneration(new(big.Int).SetUint64(i)),
			d.binding.UnpackGeneration,
		)
		if err != nil {
			return nil, d.classify(op, err)
		}

		generation, err := d.generation(op, raw)
		if err != nil {
			return nil, err
		}

		lineage = append(lineage, generation)
	}

	return lineage, nil
}

func (d *Directory) generation(op string, raw directory.Generation) (notary.Generation, error) {
	registeredAt, err := notary.UnixTimeFromSeconds(raw.RegisteredAt)
	if err != nil || raw.ChainId == 0 {
		return notary.Generation{}, notary.Errf(
			notary.KindWrongDirectory,
			op,
			err,
			"%s answered a generation without a chain or a registration time",
			d.expected.Address.Hex(),
		)
	}

	return notary.Generation{
		ChainID:       raw.ChainId,
		Ledger:        notary.Address(raw.Ledger),
		Court:         notary.Address(raw.Court),
		Withholding:   notary.Address(raw.Withholding),
		GenesisEpoch:  raw.GenesisEpoch,
		GenesisDigest: notary.Digest(raw.GenesisDigest),
		Scheme:        raw.Scheme,
		RegisteredAt:  registeredAt,
	}, nil
}

func (d *Directory) classify(op string, err error) error {
	if errors.Is(err, bind.ErrNoCode) {
		return notary.Errf(
			notary.KindWrongDirectory,
			op,
			nil,
			"chain %d has no contract at %s",
			d.expected.ChainID,
			d.expected.Address.Hex(),
		)
	}

	raw, reverted := revertData(err)
	if reverted && bytes.Equal(
		raw[:selectorLength],
		directory.LedgerDirectoryEpochBeforeDirectoryErrorID().Bytes()[:selectorLength],
	) {
		before, decodeErr := d.binding.UnpackEpochBeforeDirectoryError(raw[selectorLength:])
		if decodeErr == nil {
			return notary.Errf(
				notary.KindEpochNotAnchored,
				op,
				nil,
				"epoch %d comes before the directory's first ledger, which starts at epoch %d",
				before.Epoch,
				before.GenesisEpoch,
			)
		}
	}

	return unavailable(op, fmt.Errorf("directory: %w", err))
}

func (d *Directory) wrap(op string, err error) error {
	if errors.Is(err, errNotALedger) {
		return notary.Errf(
			notary.KindWrongDirectory,
			op,
			nil,
			"%s does not answer as a ledger directory",
			d.expected.Address.Hex(),
		)
	}

	return unavailable(op, err)
}

func directoryInterfaceID() [4]byte {
	var id [4]byte
	for _, method := range directoryinterface.NewILedgerDirectory().GetABI().Methods {
		for i := range id {
			id[i] ^= method.ID[i]
		}
	}

	return id
}
