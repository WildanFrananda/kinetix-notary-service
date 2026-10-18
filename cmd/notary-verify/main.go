package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/dig"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/prooffile"
	"github.com/WildanFrananda/kinetix-notary-service/internal/verifycli"
	"github.com/WildanFrananda/kinetix-notary-service/internal/verifying"
	"github.com/WildanFrananda/kinetix-notary-service/internal/wiring"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	if len(args) != 1 {
		_, _ = fmt.Fprintln(stderr, "usage: notary-verify <proof.json>")
		return verifycli.ExitUsage
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	document, err := verifycli.ReadProof(args[0])
	if err != nil {
		return verifycli.NoProof(stdout)
	}

	if prooffile.IsEvidence(document) {
		return evidence(ctx, document, stdout, stderr)
	}

	container, err := wiring.Verifier(os.LookupEnv)
	if err != nil {
		_, _ = fmt.Fprintln(stderr, err)
		return verifycli.ExitConfig
	}

	code := verifycli.ExitUnknown

	err = container.Invoke(func(client *ethclient.Client, scheme notary.Scheme, verifier *verifying.DirectoryVerifier) {
		defer client.Close()

		code = verifycli.RunDocument(ctx, bytes.NewReader(document), scheme.Version(), verifier, stdout)
	})

	return settle(code, err, stdout, stderr)
}

func evidence(ctx context.Context, document []byte, stdout, stderr io.Writer) int {
	container, err := wiring.EvidenceVerifier(os.LookupEnv)
	if err != nil {
		_, _ = fmt.Fprintln(stderr, err)
		return verifycli.ExitConfig
	}

	code := verifycli.ExitUnknown

	err = container.Invoke(func(
		client *ethclient.Client,
		scheme notary.EvidenceScheme,
		verifier *verifying.EvidenceVerifier,
	) {
		defer client.Close()

		code = verifycli.RunEvidence(ctx, bytes.NewReader(document), scheme.EvidenceVersion(), verifier, stdout)
	})

	return settle(code, err, stdout, stderr)
}

func settle(code int, err error, stdout, stderr io.Writer) int {
	if err == nil {
		return code
	}

	cause := dig.RootCause(err)
	if notary.KindOf(cause) != notary.KindUnknown {
		return verifycli.Fail(stdout, cause)
	}

	_, _ = fmt.Fprintln(stderr, cause)

	return verifycli.ExitConfig
}
