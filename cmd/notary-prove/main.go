package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/consensys/gnark/logger"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/plonk"
	"github.com/WildanFrananda/kinetix-notary-service/internal/proverproc"
)

const (
	srsPathKey   = "KINETIX_NOTARY_SRS_PATH"
	srsSHA256Key = "KINETIX_NOTARY_SRS_SHA256"
)

func main() {
	os.Exit(run())
}

func run() int {
	logger.Disable()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	path, sha := os.Getenv(srsPathKey), os.Getenv(srsSHA256Key)
	if path == "" || sha == "" {
		fmt.Fprintf(os.Stderr, "%s and %s are required\n", srsPathKey, srsSHA256Key)
		return proverproc.ExitMisconfigured
	}

	request, err := proverproc.DecodeRequest(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return proverproc.ExitUnprovable
	}

	prover, err := plonk.LoadProver(request.Depth, path, sha)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return proverproc.ExitMisconfigured
	}

	proof, err := plonk.NewEpochProver(prover).Prove(ctx, request)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())

		if notary.KindOf(err) == notary.KindInvalidArgument {
			return proverproc.ExitUnprovable
		}

		return proverproc.ExitFailed
	}

	if _, err := fmt.Fprintln(os.Stdout, "0x"+hex.EncodeToString(proof)); err != nil {
		return proverproc.ExitFailed
	}

	return proverproc.ExitProved
}
