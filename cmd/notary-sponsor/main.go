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
	"github.com/WildanFrananda/kinetix-notary-service/internal/sponsorcli"
	"github.com/WildanFrananda/kinetix-notary-service/internal/sponsoring"
	"github.com/WildanFrananda/kinetix-notary-service/internal/verifying"
	"github.com/WildanFrananda/kinetix-notary-service/internal/wiring"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		_, _ = fmt.Fprintln(stderr, sponsorcli.Usage)
		return sponsorcli.ExitUsage
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	container, err := wiring.Sponsor(os.LookupEnv)
	if err != nil {
		_, _ = fmt.Fprintln(stderr, err)
		return sponsorcli.ExitConfig
	}

	code := sponsorcli.ExitUnknown

	err = container.Invoke(func(
		client *ethclient.Client,
		newest notary.ConfidentialScheme,
		statements notary.StatementScheme,
		directory notary.Directory,
		verifier *verifying.DirectoryVerifier,
		sponsor *sponsoring.Sponsor,
	) {
		defer client.Close()

		code = sponsorcli.Run(
			ctx,
			args,
			sponsorcli.Dependencies{
				Sponsor:    sponsor,
				Directory:  directory,
				Verifier:   verifier,
				Statements: statements,
				Scheme:     newest.Version(),
			},
			stdout,
		)
		if code == sponsorcli.ExitUsage {
			_, _ = fmt.Fprintln(stderr, sponsorcli.Usage)
		}
	})
	if err != nil {
		cause := dig.RootCause(err)
		if notary.KindOf(cause) != notary.KindUnknown {
			return sponsorcli.Fail(stdout, args[0], cause)
		}

		_, _ = fmt.Fprintln(stderr, cause)

		return sponsorcli.ExitConfig
	}

	return code
}
