package integration_test

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/consensys/gnark/test/unsafekzg"
	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/ethereum"
	"github.com/WildanFrananda/kinetix-notary-service/internal/keccak"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/plonk"
	"github.com/WildanFrananda/kinetix-notary-service/internal/poseidon"
	"github.com/WildanFrananda/kinetix-notary-service/internal/proverproc"
	"github.com/WildanFrananda/kinetix-notary-service/internal/verifycli"
	"github.com/WildanFrananda/kinetix-notary-service/internal/verifying"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

func testSRSFile(t *testing.T, depth int) (string, string) {
	t.Helper()

	system, err := plonk.Compile(depth)
	require.NoError(t, err)

	canonical, _, err := unsafekzg.NewSRS(system, unsafekzg.WithToxicValue(notarytest.TestSRSToxicValue()))
	require.NoError(t, err)

	var encoded bytes.Buffer

	_, err = canonical.WriteTo(&encoded)
	require.NoError(t, err)

	path := filepath.Join(t.TempDir(), "test.srs")
	require.NoError(t, os.WriteFile(path, encoded.Bytes(), 0o600))

	sum := sha256.Sum256(encoded.Bytes())

	return path, hex.EncodeToString(sum[:])
}

func putProverOnPath(t *testing.T) {
	t.Helper()

	out := t.TempDir()
	build := exec.Command(
		"go",
		"build",
		"-buildvcs=false",
		"-o",
		filepath.Join(out, proverproc.Binary),
		"./cmd/notary-prove",
	)
	build.Dir = filepath.Join("..", "..")
	output, err := build.CombinedOutput()
	require.NoError(t, err, string(output))

	t.Setenv("PATH", out+string(os.PathListSeparator)+os.Getenv("PATH"))
}

func TestTheDaemonProvesAnchorsAndServesASchemeTwoEpochEveryHolderCanVerify(t *testing.T) {
	e := newEstateOf(t, testConfidential())
	r := e.start(t)
	ctx := context.Background()

	require.NoError(t, r.Prover.CheckCourt(ctx))
	require.NoError(t, r.Directory.Check(ctx))

	e.tickUntil(t, r, 0, "mined")
	e.chain.finalize(t)
	e.tickUntil(t, r, 0, "finalized")

	var total, proofs int64
	require.NoError(t, e.pool.QueryRow(ctx, "SELECT total FROM epochs WHERE number = 0").Scan(&total))
	require.NoError(
		t,
		e.pool.QueryRow(
			ctx,
			"SELECT count(*) FROM attempts WHERE epoch_number = 0 AND kind = 'proof' AND outcome = 'included'",
		).Scan(&proofs),
	)
	require.Zero(t, total)
	require.Equal(t, int64(1), proofs)

	settings := ethereum.Settings{RPCURL: e.chain.url, Ledger: e.ledger, Timeout: 10 * time.Second}
	directory := verifying.NewDirectoryVerifier(
		ethereum.NewDirectory(e.chain.client, ethereum.Settings{
			RPCURL:    e.chain.url,
			Ledger:    e.ledger,
			Directory: notary.DirectoryIdentity{ChainID: e.ledger.ChainID, Address: e.directory},
			Timeout:   10 * time.Second,
		}),
		ethereum.NewReaders(e.chain.client, settings),
		keccak.NewScheme(),
		poseidon.NewScheme(),
		keccak.NewScheme(),
	)

	for _, holder := range []string{"principal-0", "principal-3", "principal-5"} {
		response := e.proofs(t, r, "latest", holder)
		require.Equal(t, http.StatusOK, response.Code, response.Body.String())
		require.Contains(t, response.Body.String(), `"scheme":2`)

		path := filepath.Join(t.TempDir(), "proofs.json")
		require.NoError(t, os.WriteFile(path, response.Body.Bytes(), 0o600))

		var out bytes.Buffer

		code := verifycli.Run(ctx, path, poseidon.NewScheme().Version(), directory, &out)
		require.Equal(t, verifycli.ExitIncluded, code, out.String())

		var report map[string]string
		require.NoError(t, json.Unmarshal(bytes.TrimSpace(out.Bytes()), &report))
		require.Equal(t, "included", report["verdict"])
		require.Equal(t, holder, report["holder"])
		require.Equal(t, "0", report["total"])
		require.Equal(t, "attributable", report["statement"])
	}
}
