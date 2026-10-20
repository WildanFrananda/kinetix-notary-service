package proverproc_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/consensys/gnark/test/unsafekzg"
	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/plonk"
	"github.com/WildanFrananda/kinetix-notary-service/internal/poseidon"
	"github.com/WildanFrananda/kinetix-notary-service/internal/proverproc"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

const depth = 2

type harness struct {
	srs string
	sha string
}

func TestTheChildProverProvesAnEpoch(t *testing.T) {
	h := newHarness(t)

	proof, err := proverproc.NewClient(h.env(h.sha)).Prove(context.Background(), request(t))
	require.NoError(t, err)
	require.NotEmpty(t, proof)
}

func TestTheChildProverRefusesAStatementItsEntriesDoNotMake(t *testing.T) {
	h := newHarness(t)
	forged := request(t)
	forged.Entries[0].Amount++

	_, err := proverproc.NewClient(h.env(h.sha)).Prove(context.Background(), forged)
	require.Equal(t, notary.KindInvalidArgument, notary.KindOf(err))
	require.Contains(t, err.Error(), "cannot be proven")
}

func TestTheChildProverRefusesAnSRSThatIsNotThePinnedOne(t *testing.T) {
	h := newHarness(t)

	_, err := proverproc.NewClient(h.env(hex.EncodeToString(make([]byte, 32)))).
		Prove(context.Background(), request(t))
	require.Equal(t, notary.KindInvalidArgument, notary.KindOf(err))
	require.Contains(t, err.Error(), "misconfigured")
}

func TestTheChildProverRefusesToRunUnconfigured(t *testing.T) {
	newHarness(t)

	_, err := proverproc.NewClient(nil).Prove(context.Background(), request(t))
	require.Equal(t, notary.KindInvalidArgument, notary.KindOf(err))
	require.Contains(t, err.Error(), "KINETIX_NOTARY_SRS_PATH")
}

func TestAChildThatDiesIsAFailureNotAVerdict(t *testing.T) {
	h := newHarness(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := proverproc.NewClient(h.env(h.sha)).Prove(ctx, request(t))
	require.Equal(t, notary.KindUnknown, notary.KindOf(err))
}

func newHarness(t *testing.T) harness {
	t.Helper()

	out := t.TempDir()
	binary := filepath.Join(out, proverproc.Binary)
	root := filepath.Join("..", "..")

	build := exec.Command("go", "build", "-buildvcs=false", "-o", binary, "./cmd/notary-prove")
	build.Dir = root
	output, err := build.CombinedOutput()
	require.NoError(t, err, string(output))

	system, err := plonk.Compile(depth)
	require.NoError(t, err)

	canonical, _, err := unsafekzg.NewSRS(system)
	require.NoError(t, err)

	path := filepath.Join(out, "test.srs")
	file, err := os.Create(path)
	require.NoError(t, err)

	hash := sha256.New()
	_, err = canonical.WriteTo(file)
	require.NoError(t, err)
	require.NoError(t, file.Close())

	written, err := os.ReadFile(path)
	require.NoError(t, err)
	hash.Write(written)

	t.Setenv("PATH", out+string(os.PathListSeparator)+os.Getenv("PATH"))

	return harness{srs: path, sha: hex.EncodeToString(hash.Sum(nil))}
}

func (h harness) env(sha string) []string {
	return []string{"KINETIX_NOTARY_SRS_PATH=" + h.srs, "KINETIX_NOTARY_SRS_SHA256=" + sha}
}

func request(t *testing.T) notary.ProofRequest {
	t.Helper()

	scheme := poseidon.NewScheme()
	entries := notarytest.Entries(5, 3)

	tree, err := notary.BuildConfidentialTree(scheme, 4, depth, entries)
	require.NoError(t, err)

	blinding := notary.Blinding{0x21}

	header, err := tree.Header(scheme, notary.Digest{}, notarytest.Unix(t, 1_790_000_000), blinding)
	require.NoError(t, err)

	return notary.ProofRequest{
		Number:          4,
		Depth:           depth,
		Entries:         entries,
		TreeRoot:        header.TreeRoot,
		TotalCommitment: header.TotalCommitment,
		Blinding:        blinding,
	}
}
