package fixtures_test

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark/frontend"
	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/plonk"
	"github.com/WildanFrananda/kinetix-notary-service/internal/poseidon"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

const (
	plonkFixtureDepth    = 3
	plonkFixtureEpoch    = 11
	plonkFixtureAccounts = 6
	plonkFixtureAsOf     = 1_790_000_000
)

var (
	plonkFixturePath  = filepath.Join("..", "..", "contracts", "test", "fixtures", "plonk-v2.json")
	plonkVerifierPath = filepath.Join("..", "..", "contracts", "test", "support", "TestPlonkVerifier.sol")
)

type plonkFixture struct {
	Scheme          uint16      `json:"scheme"`
	Depth           int         `json:"depth"`
	Epoch           string      `json:"epoch"`
	TreeRoot        string      `json:"treeRoot"`
	TotalCommitment string      `json:"totalCommitment"`
	LeafCount       string      `json:"leafCount"`
	PreviousDigest  string      `json:"previousDigest"`
	Root            string      `json:"root"`
	AsOf            string      `json:"asOf"`
	Digest          string      `json:"digest"`
	Commitments     []string    `json:"commitments"`
	Amounts         []string    `json:"amounts"`
	Inclusions      []proofCase `json:"inclusions"`
	PublicInputs    []string    `json:"publicInputs"`
	Proof           string      `json:"proof"`
}

func TestThePlonkFixturesMatchTheGoProver(t *testing.T) {
	prover, assignment, header := plonkFixtureProver(t)
	scheme := poseidon.NewScheme()
	entries := notarytest.Entries(900, plonkFixtureAccounts)

	tree, err := notary.BuildConfidentialTree(scheme, plonkFixtureEpoch, plonkFixtureDepth, entries)
	require.NoError(t, err)

	var exported bytes.Buffer
	require.NoError(t, prover.ExportVerifier(&exported))

	verifier := []byte(strings.Replace(exported.String(), "contract PlonkVerifier", "contract TestPlonkVerifier", 1))
	built := plonkFixture{
		Scheme:          poseidon.NewScheme().Version(),
		Depth:           prover.Depth(),
		Epoch:           u64(plonkFixtureEpoch),
		TreeRoot:        header.TreeRoot.Hex(),
		TotalCommitment: header.TotalCommitment.Hex(),
		LeafCount:       u64(header.Header.LeafCount),
		PreviousDigest:  header.Header.PreviousDigest.Hex(),
		Root:            header.Header.Root.Hex(),
		AsOf:            u64(header.Header.AsOf.Seconds()),
		Digest:          scheme.EpochDigest(header.Header).Hex(),
		Commitments:     make([]string, len(entries)),
		Amounts:         make([]string, len(entries)),
		Inclusions:      make([]proofCase, len(entries)),
		PublicInputs:    publicInputsOf(t, assignment),
	}

	for i, entry := range entries {
		inclusion, err := tree.Prove(uint64(i))
		require.NoError(t, err)

		built.Commitments[i] = entry.Commitment.Hex()
		built.Amounts[i] = u64(entry.Amount)
		built.Inclusions[i] = encodeProof(inclusion)
	}

	if *update {
		proof, err := prover.Prove(assignment)
		require.NoError(t, err)

		built.Proof = "0x" + hex.EncodeToString(proof)

		encoded, err := json.MarshalIndent(built, "", "  ")
		require.NoError(t, err)
		require.NoError(t, os.WriteFile(plonkFixturePath, append(encoded, '\n'), 0o644))
		require.NoError(t, os.WriteFile(plonkVerifierPath, verifier, 0o644))

		return
	}

	committedVerifier, err := os.ReadFile(plonkVerifierPath)
	require.NoError(t, err, "no verifier yet: run `bin/dev go test ./tests/fixtures -update`")
	require.True(
		t,
		bytes.Equal(committedVerifier, verifier),
		"contracts/test/support/TestPlonkVerifier.sol is stale: run `bin/dev go test ./tests/fixtures -update`",
	)

	encoded, err := os.ReadFile(plonkFixturePath)
	require.NoError(t, err, "no fixture yet: run `bin/dev go test ./tests/fixtures -update`")

	var committed plonkFixture
	require.NoError(t, json.Unmarshal(encoded, &committed))
	require.NotEmpty(t, committed.Proof)

	built.Proof = committed.Proof
	require.Equal(
		t,
		built,
		committed,
		"contracts/test/fixtures/plonk-v2.json is stale: run `bin/dev go test ./tests/fixtures -update`",
	)
}

func TestThePlonkFixtureSetupIsDeterministic(t *testing.T) {
	first, _, _ := plonkFixtureProver(t)
	second, _, _ := plonkFixtureProver(t)

	var a, b bytes.Buffer
	require.NoError(t, first.ExportVerifier(&a))
	require.NoError(t, second.ExportVerifier(&b))
	require.Equal(t, a.String(), b.String())
}

func plonkFixtureProver(t *testing.T) (*plonk.Prover, *plonk.LiabilityCircuit, notary.ConfidentialHeader) {
	t.Helper()

	scheme := poseidon.NewScheme()
	entries := notarytest.Entries(900, plonkFixtureAccounts)

	tree, err := notary.BuildConfidentialTree(scheme, plonkFixtureEpoch, plonkFixtureDepth, entries)
	require.NoError(t, err)

	seed := sha256.Sum256([]byte("kinetix.notary.plonk.blinding"))

	var blinding notary.Blinding
	copy(blinding[:], seed[:])

	previous := notary.Digest(sha256.Sum256([]byte("kinetix.notary.plonk.previous")))

	header, err := tree.Header(scheme, previous, notarytest.Unix(t, plonkFixtureAsOf), blinding)
	require.NoError(t, err)

	assignment, err := plonk.NewAssignment(plonkFixtureEpoch, plonkFixtureDepth, entries, header, blinding)
	require.NoError(t, err)

	prover := notarytest.TestProver(t, plonkFixtureDepth)

	return prover, assignment, header
}

func publicInputsOf(t *testing.T, assignment *plonk.LiabilityCircuit) []string {
	t.Helper()

	public, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField(), frontend.PublicOnly())
	require.NoError(t, err)

	vector, ok := public.Vector().(fr.Vector)
	require.True(t, ok)

	out := make([]string, len(vector))
	for i := range vector {
		out[i] = fmt.Sprintf("0x%064x", vector[i].BigInt(new(big.Int)))
	}

	return out
}
