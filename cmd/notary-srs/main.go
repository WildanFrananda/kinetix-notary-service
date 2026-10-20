package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	kzgbn254 "github.com/consensys/gnark-crypto/ecc/bn254/kzg"

	"github.com/WildanFrananda/kinetix-notary-service/internal/ignition"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/plonk"
)

const defaultBase = "https://aztec-ignition.s3.amazonaws.com/"

func main() {
	os.Exit(run())
}

func run() int {
	log := slog.New(slog.NewJSONHandler(os.Stderr, nil))

	var options options

	flag.IntVar(&options.depth, "depth", -1, "the circuit depth the SRS must serve")
	flag.StringVar(&options.srsPath, "srs", "", "import the ceremony and write the canonical SRS here")
	flag.StringVar(&options.fromPath, "from", "", "read an SRS imported earlier instead of importing")
	flag.StringVar(
		&options.verifierPath,
		"verifier",
		"",
		"where to write the Solidity verifier of that depth, if given",
	)
	flag.StringVar(&options.base, "base", defaultBase, "where the ceremony is published")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := runWith(ctx, log, options); err != nil {
		log.Error("notary-srs failed", "kind", notary.KindOf(err).String(), "err", err.Error())
		return 1
	}

	return 0
}

type options struct {
	depth        int
	srsPath      string
	fromPath     string
	verifierPath string
	base         string
}

func runWith(ctx context.Context, log *slog.Logger, options options) error {
	const op = "notary-srs"
	if options.depth < 0 || options.depth > notary.MaxConfidentialDepth ||
		(options.srsPath == "") == (options.fromPath == "") {
		return notary.Errf(
			notary.KindInvalidArgument,
			op,
			nil,
			"-depth 0..%d and exactly one of -srs or -from are required",
			notary.MaxConfidentialDepth,
		)
	}

	system, err := plonk.Compile(options.depth)
	if err != nil {
		return err
	}

	size := plonk.SRSSize(system)
	log.InfoContext(ctx, "circuit", "depth", options.depth, "constraints", system.GetNbConstraints(), "points", size)

	var srs *kzgbn254.SRS
	if options.fromPath != "" {
		srs, err = readSRS(ctx, log, options.fromPath, size)
	} else {
		srs, err = importSRS(ctx, log, options.base, options.srsPath, size)
	}

	if err != nil {
		return err
	}

	if options.verifierPath == "" {
		return nil
	}

	lagrange, err := plonk.NewLagrange(srs, system)
	if err != nil {
		return err
	}

	prover, err := plonk.NewProver(options.depth, system, srs, lagrange)
	if err != nil {
		return err
	}

	key := sha256.New()
	if err := prover.WriteVerifyingKey(key); err != nil {
		return err
	}

	digest, err := writeHashed(options.verifierPath, prover.ExportVerifier)
	if err != nil {
		return err
	}

	log.InfoContext(
		ctx,
		"verifier written",
		"path",
		options.verifierPath,
		"sha256",
		digest,
		"verifyingKeySha256",
		hex.EncodeToString(key.Sum(nil)),
	)

	return nil
}

func importSRS(ctx context.Context, log *slog.Logger, base, path string, size int) (*kzgbn254.SRS, error) {
	client, err := ignition.NewClient(base, ignition.MainIgnition(), &http.Client{Timeout: time.Hour})
	if err != nil {
		return nil, err
	}

	srs, err := ignition.Import(ctx, client, size)
	if err != nil {
		return nil, err
	}

	digest, err := writeHashed(path, func(w io.Writer) error {
		if _, err := srs.WriteTo(w); err != nil {
			return notary.Errf(notary.KindInvalidArgument, "notary-srs.importSRS", err, "encoding the SRS")
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	log.InfoContext(ctx, "SRS written", "path", path, "points", len(srs.Pk.G1), "sha256", digest)

	return srs, nil
}

func readSRS(ctx context.Context, log *slog.Logger, path string, size int) (*kzgbn254.SRS, error) {
	const op = "notary-srs.readSRS"

	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, notary.Errf(notary.KindInvalidArgument, op, err, "%s", path)
	}
	defer func() { _ = file.Close() }()

	hash := sha256.New()

	var srs kzgbn254.SRS
	if _, err := srs.ReadFrom(io.TeeReader(file, hash)); err != nil {
		return nil, notary.Errf(notary.KindInvalidArgument, op, err, "decoding %s", path)
	}

	if len(srs.Pk.G1) < size {
		return nil, notary.Errf(
			notary.KindInvalidArgument,
			op,
			nil,
			"%s holds %d points, the circuit needs %d",
			path,
			len(srs.Pk.G1),
			size,
		)
	}

	log.InfoContext(
		ctx,
		"SRS read",
		"path",
		path,
		"points",
		len(srs.Pk.G1),
		"sha256",
		hex.EncodeToString(hash.Sum(nil)),
	)

	return &srs, nil
}

func writeHashed(path string, write func(io.Writer) error) (string, error) {
	const op = "notary-srs.writeHashed"

	file, err := os.OpenFile(filepath.Clean(path), os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return "", notary.Errf(notary.KindInvalidArgument, op, err, "%s", path)
	}

	hash := sha256.New()
	if err := write(io.MultiWriter(file, hash)); err != nil {
		_ = file.Close()

		return "", notary.Errf(notary.KindInvalidArgument, op, err, "writing %s", path)
	}

	if err := file.Close(); err != nil {
		return "", notary.Errf(notary.KindInvalidArgument, op, err, "closing %s", path)
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
