package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/dig"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
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

	container, err := wiring.Verifier(os.LookupEnv)
	if err != nil {
		_, _ = fmt.Fprintln(stderr, err)
		return verifycli.ExitConfig
	}

	code := verifycli.ExitUnknown

	err = container.Invoke(func(client *ethclient.Client, scheme notary.Scheme, verifier *verifying.Verifier) {
		defer client.Close()

		code = verifycli.Run(ctx, args[0], scheme.Version(), verifier, stdout)
	})
	if err != nil {
		cause := dig.RootCause(err)
		if notary.KindOf(cause) != notary.KindUnknown {
			return verifycli.Fail(stdout, cause)
		}

		_, _ = fmt.Fprintln(stderr, cause)

		return verifycli.ExitConfig
	}

	return code
}
