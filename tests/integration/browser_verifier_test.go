package integration_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"maps"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/config"
	"github.com/WildanFrananda/kinetix-notary-service/internal/keccak"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/prooffile"
	"github.com/WildanFrananda/kinetix-notary-service/internal/verifycli"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

type answer struct {
	exit   int
	stdout string
	stderr string
}

type verifiers struct {
	native   string
	module   string
	wasmExec string
	harness  string
	bun      string
}

func buildVerifiers(t *testing.T) verifiers {
	t.Helper()

	bun, err := exec.LookPath("bun")
	require.NoError(t, err, "the WebAssembly verifier runs under Bun: Go turns fetch off under Node")

	root := strings.TrimSpace(commandOutput(t, "", "go", "list", "-m", "-f", "{{.Dir}}"))
	out := t.TempDir()

	commandOutput(t, root, "sh", "bin/build-verifier-wasm", "integration", out)
	commandOutput(t, root, "go", "build", "-o", filepath.Join(out, "notary-verify"), "./cmd/notary-verify")

	sums, err := os.ReadFile(filepath.Join(out, "SHA256SUMS"))
	require.NoError(t, err)
	require.Contains(t, string(sums), "notary-verify.wasm")
	require.Contains(t, string(sums), "wasm_exec.js")

	return verifiers{
		native:   filepath.Join(out, "notary-verify"),
		module:   filepath.Join(out, "notary-verify.wasm"),
		wasmExec: filepath.Join(out, "wasm_exec.js"),
		harness:  filepath.Join(root, "tests", "integration", "testdata", "verify-wasm.mjs"),
		bun:      bun,
	}
}

func commandOutput(t *testing.T, dir string, name string, args ...string) string {
	t.Helper()

	command := exec.Command(name, args...)
	command.Dir = dir

	var stderr bytes.Buffer

	command.Stderr = &stderr

	out, err := command.Output()
	require.NoError(t, err, "%s %v: %s", name, args, stderr.String())

	return string(out)
}

func (v verifiers) cli(t *testing.T, document []byte, settings map[string]string) answer {
	t.Helper()

	path := filepath.Join(t.TempDir(), "proof.json")
	require.NoError(t, os.WriteFile(path, document, 0o600))

	command := exec.Command(v.native, path)
	for key, value := range settings {
		command.Env = append(command.Env, key+"="+value)
	}

	return answered(t, command)
}

func (v verifiers) browser(t *testing.T, document []byte, settings map[string]string) answer {
	t.Helper()

	path := filepath.Join(t.TempDir(), "proof.json")
	require.NoError(t, os.WriteFile(path, document, 0o600))

	encoded, err := json.Marshal(settings)
	require.NoError(t, err)

	return answered(t, exec.Command(v.bun, v.harness, v.wasmExec, v.module, path, string(encoded)))
}

func answered(t *testing.T, command *exec.Cmd) answer {
	t.Helper()

	var stdout, stderr bytes.Buffer

	command.Stdout = &stdout
	command.Stderr = &stderr

	err := command.Run()

	var exit *exec.ExitError
	if errors.As(err, &exit) {
		return answer{exit: exit.ExitCode(), stdout: stdout.String(), stderr: stderr.String()}
	}

	require.NoError(t, err)

	return answer{stdout: stdout.String(), stderr: stderr.String()}
}

func encoded(t *testing.T, claims ...notary.Claim) []byte {
	t.Helper()

	documents := make([]json.RawMessage, len(claims))

	for i, claim := range claims {
		var doc bytes.Buffer
		require.NoError(t, prooffile.Encode(&doc, claim, keccak.NewScheme().Version()))
		documents[i] = doc.Bytes()
	}

	if len(documents) == 1 {
		return documents[0]
	}

	bundle, err := json.Marshal(map[string][]json.RawMessage{"proofs": documents})
	require.NoError(t, err)

	return bundle
}

func verifierSettings(g generations) map[string]string {
	return map[string]string{
		config.RPCURLKey:           g.chain.url,
		config.ChainIDKey:          strconv.FormatUint(g.identity.ChainID, 10),
		config.DirectoryAddressKey: notary.Address(g.directory).Hex(),
		config.RPCTimeoutKey:       "10s",
	}
}

func settingsWith(settings map[string]string, key, value string) map[string]string {
	changed := maps.Clone(settings)
	if value == "" {
		delete(changed, key)
	} else {
		changed[key] = value
	}

	return changed
}

func TestTheWebAssemblyVerifierAnswersExactlyAsTheCLI(t *testing.T) {
	v := buildVerifiers(t)

	a := anchorEpochZero(t, 7, true)
	g := withDirectory(t, a)
	settings := verifierSettings(g)

	signed := a.sign(t, a.claim(t, 3), a.key, a.signer)
	forged := a.claim(t, 3)
	forged.Entry.Amount++
	forged = a.sign(t, forged, a.key, a.signer)
	impostorKey, impostor := notarytest.StatementKey(t, "integration-impostor")
	disowned := a.sign(t, a.claim(t, 5), impostorKey, impostor)
	unanchored := a.claim(t, 1)
	unanchored.Epoch = 5

	pending := anchorEpochZero(t, 5, false)
	pendingSettings := verifierSettings(withDirectory(t, pending))

	cases := []struct {
		name     string
		document []byte
		settings map[string]string
		exit     int
	}{
		{"an attributable inclusion", encoded(t, signed), settings, verifycli.ExitIncluded},
		{"a statement another key signed", encoded(t, disowned), settings, verifycli.ExitIncluded},
		{"a provably false statement", encoded(t, forged), settings, verifycli.ExitRejected},
		{
			"a bundle answers claim by claim",
			encoded(t, signed, forged, a.claim(t, 0)),
			settings,
			verifycli.ExitRejected,
		},
		{"an epoch never anchored", encoded(t, unanchored), settings, verifycli.ExitNotAnchored},
		{"an epoch not final yet", encoded(t, pending.claim(t, 1)), pendingSettings, verifycli.ExitNotFinal},
		{
			"a ledger named as the directory",
			encoded(t, signed),
			settingsWith(settings, config.DirectoryAddressKey, a.address.Hex()),
			verifycli.ExitWrongDirectory,
		},
		{"a document that is not a proof", []byte(`{"format":`), settings, verifycli.ExitMalformedProof},
		{
			"a verifier without its timeout",
			encoded(t, signed),
			settingsWith(settings, config.RPCTimeoutKey, ""),
			verifycli.ExitConfig,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cli := v.cli(t, c.document, c.settings)
			browser := v.browser(t, c.document, c.settings)

			require.Equal(t, c.exit, cli.exit, "%s%s", cli.stdout, cli.stderr)
			require.Equal(t, cli, browser)
		})
	}
}

func TestAnEndpointThatDoesNotAnswerIsUnknownInTheBrowserToo(t *testing.T) {
	v := buildVerifiers(t)

	a := anchorEpochZero(t, 3, true)
	settings := settingsWith(verifierSettings(withDirectory(t, a)), config.RPCURLKey, "http://127.0.0.1:1")
	document := encoded(t, a.claim(t, 0))

	cli := v.cli(t, document, settings)
	browser := v.browser(t, document, settings)

	require.Equal(t, verifycli.ExitUnknown, cli.exit, cli.stdout)
	require.Equal(t, cli.exit, browser.exit, browser.stdout+browser.stderr)

	var cliReport, browserReport map[string]string
	require.NoError(t, json.Unmarshal([]byte(cli.stdout), &cliReport))
	require.NoError(t, json.Unmarshal([]byte(browser.stdout), &browserReport))
	require.Equal(t, "unknown", browserReport["verdict"])
	require.Equal(t, cliReport["verdict"], browserReport["verdict"])
}
