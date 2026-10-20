//go:build js && wasm

package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"syscall/js"

	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/dig"

	"github.com/WildanFrananda/kinetix-notary-service/internal/config"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/prooffile"
	"github.com/WildanFrananda/kinetix-notary-service/internal/verifycli"
	"github.com/WildanFrananda/kinetix-notary-service/internal/verifying"
	"github.com/WildanFrananda/kinetix-notary-service/internal/wiring"
)

var version = "unknown"

func main() {
	js.Global().Set("kinetixNotaryVerifierVersion", version)
	js.Global().Set("kinetixNotaryVerify", js.FuncOf(verify))

	select {}
}

func verify(_ js.Value, args []js.Value) any {
	document, lookup := arguments(args)

	var executor js.Func

	executor = js.FuncOf(func(_ js.Value, callbacks []js.Value) any {
		resolve := callbacks[0]

		go func() {
			defer executor.Release()

			var stdout, stderr bytes.Buffer

			code := run(context.Background(), document, lookup, &stdout, &stderr)

			resolve.Invoke(map[string]any{
				"exit":   code,
				"stdout": stdout.String(),
				"stderr": stderr.String(),
			})
		}()

		return nil
	})

	return js.Global().Get("Promise").New(executor)
}

func arguments(args []js.Value) (string, config.Lookup) {
	document := ""
	if len(args) > 0 && args[0].Type() == js.TypeString {
		document = args[0].String()
	}

	settings := js.Undefined()
	if len(args) > 1 && args[1].Type() == js.TypeObject {
		settings = args[1]
	}

	return document, func(key string) (string, bool) {
		if settings.IsUndefined() {
			return "", false
		}

		value := settings.Get(key)
		if value.Type() != js.TypeString {
			return "", false
		}

		return value.String(), true
	}
}

func run(ctx context.Context, document string, lookup config.Lookup, stdout, stderr io.Writer) int {
	if prooffile.IsEvidence([]byte(document)) {
		return evidence(ctx, document, lookup, stdout, stderr)
	}

	container, err := wiring.Verifier(lookup)
	if err != nil {
		_, _ = fmt.Fprintln(stderr, err)
		return verifycli.ExitConfig
	}

	code := verifycli.ExitUnknown

	err = container.Invoke(func(
		client *ethclient.Client,
		newest notary.ConfidentialScheme,
		verifier *verifying.DirectoryVerifier,
	) {
		defer client.Close()

		code = verifycli.RunDocument(ctx, strings.NewReader(document), newest.Version(), verifier, stdout)
	})

	return settle(code, err, stdout, stderr)
}

func evidence(ctx context.Context, document string, lookup config.Lookup, stdout, stderr io.Writer) int {
	container, err := wiring.EvidenceVerifier(lookup)
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

		code = verifycli.RunEvidence(ctx, strings.NewReader(document), scheme.EvidenceVersion(), verifier, stdout)
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
