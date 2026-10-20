package plonk

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"

	kzgbn254 "github.com/consensys/gnark-crypto/ecc/bn254/kzg"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

func LoadSRS(path, wantSHA256 string) (*kzgbn254.SRS, error) {
	const op = "plonk.LoadSRS"

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

	if _, err := io.Copy(hash, file); err != nil {
		return nil, notary.Errf(notary.KindInvalidArgument, op, err, "reading %s", path)
	}

	if got := hex.EncodeToString(hash.Sum(nil)); got != wantSHA256 {
		return nil, notary.Errf(
			notary.KindInvalidArgument,
			op,
			nil,
			"%s has SHA-256 %s, the pin is %s",
			path,
			got,
			wantSHA256,
		)
	}

	return &srs, nil
}

func LoadProver(depth int, srsPath, srsSHA256 string) (*Prover, error) {
	system, err := Compile(depth)
	if err != nil {
		return nil, err
	}

	srs, err := LoadSRS(srsPath, srsSHA256)
	if err != nil {
		return nil, err
	}

	lagrange, err := NewLagrange(srs, system)
	if err != nil {
		return nil, err
	}

	return NewProver(depth, system, srs, lagrange)
}

func CheckSRSFile(path, wantSHA256 string) error {
	const op = "plonk.CheckSRSFile"

	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		return notary.Errf(notary.KindInvalidArgument, op, err, "%s", path)
	}
	defer func() { _ = file.Close() }()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return notary.Errf(notary.KindInvalidArgument, op, err, "reading %s", path)
	}

	if got := hex.EncodeToString(hash.Sum(nil)); got != wantSHA256 {
		return notary.Errf(
			notary.KindInvalidArgument,
			op,
			nil,
			"%s has SHA-256 %s, the pin is %s",
			path,
			got,
			wantSHA256,
		)
	}

	return nil
}
