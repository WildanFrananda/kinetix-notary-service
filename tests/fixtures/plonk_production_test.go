package fixtures_test

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"testing"

	kzgbn254 "github.com/consensys/gnark-crypto/ecc/bn254/kzg"
	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/plonk"
	"github.com/WildanFrananda/kinetix-notary-service/internal/poseidon"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

const (
	productionEpoch    = 176
	productionAccounts = 5
	ignitionSRSKey     = "KINETIX_NOTARY_IGNITION_SRS"
)

var (
	productionFixturePath  = filepath.Join("..", "..", "contracts", "test", "fixtures", "plonk-production.json")
	productionVerifierPath = filepath.Join("..", "..", "contracts", "src", "confidential", "PlonkVerifier.sol")
)

type productionFixture struct {
	Depth              int         `json:"depth"`
	SRSSHA256          string      `json:"srsSha256"`
	VerifyingKeySHA256 string      `json:"verifyingKeySha256"`
	VerifierSHA256     string      `json:"verifierSha256"`
	Epoch              string      `json:"epoch"`
	TreeRoot           string      `json:"treeRoot"`
	TotalCommitment    string      `json:"totalCommitment"`
	LeafCount          string      `json:"leafCount"`
	Commitments        []string    `json:"commitments"`
	Amounts            []string    `json:"amounts"`
	Inclusions         []proofCase `json:"inclusions"`
	PublicInputs       []string    `json:"publicInputs"`
	Proof              string      `json:"proof"`
}

func TestTheProductionProofFixtureMatchesTheGoProver(t *testing.T) {
	assignment, header, tree, entries := productionStatement(t)
	built := productionFixture{
		Depth:              plonk.ProductionDepth,
		SRSSHA256:          plonk.ProductionSRSSHA256,
		VerifyingKeySHA256: plonk.ProductionVerifyingKeySHA256,
		VerifierSHA256:     plonk.ProductionVerifierSHA256,
		Epoch:              u64(productionEpoch),
		TreeRoot:           header.TreeRoot.Hex(),
		TotalCommitment:    header.TotalCommitment.Hex(),
		LeafCount:          u64(header.Header.LeafCount),
		Commitments:        make([]string, len(entries)),
		Amounts:            make([]string, len(entries)),
		Inclusions:         make([]proofCase, len(entries)),
		PublicInputs:       publicInputsOf(t, assignment),
	}

	for i, entry := range entries {
		inclusion, err := tree.Prove(uint64(i))
		require.NoError(t, err)

		built.Commitments[i] = entry.Commitment.Hex()
		built.Amounts[i] = u64(entry.Amount)
		built.Inclusions[i] = encodeProof(inclusion)
	}

	srsPath := os.Getenv(ignitionSRSKey)
	if *update {
		require.NotEmpty(t, srsPath, "set %s to the SRS notary-srs imported", ignitionSRSKey)

		proof, err := productionProver(t, srsPath).Prove(assignment)
		require.NoError(t, err)

		built.Proof = "0x" + hex.EncodeToString(proof)

		encoded, err := json.MarshalIndent(built, "", "  ")
		require.NoError(t, err)
		require.NoError(t, os.WriteFile(productionFixturePath, append(encoded, '\n'), 0o644))

		return
	}

	if srsPath != "" {
		productionProver(t, srsPath)
	}

	encoded, err := os.ReadFile(productionFixturePath)
	require.NoError(
		t,
		err,
		"no fixture yet: run `bin/dev env %s=<srs> go test ./tests/fixtures -update`",
		ignitionSRSKey,
	)

	var committed productionFixture
	require.NoError(t, json.Unmarshal(encoded, &committed))
	require.NotEmpty(t, committed.Proof)

	built.Proof = committed.Proof
	require.Equal(
		t,
		built,
		committed,
		"contracts/test/fixtures/plonk-production.json is stale: prove it again with the SRS",
	)
}

func productionStatement(
	t *testing.T,
) (*plonk.LiabilityCircuit, notary.ConfidentialHeader, *notary.ConfidentialTree, []notary.Entry) {
	t.Helper()

	scheme := poseidon.NewScheme()
	entries := notarytest.Entries(700, productionAccounts)

	tree, err := notary.BuildConfidentialTree(scheme, productionEpoch, plonk.ProductionDepth, entries)
	require.NoError(t, err)

	seed := sha256.Sum256([]byte("kinetix.notary.plonk.production.blinding"))

	var blinding notary.Blinding
	copy(blinding[:], seed[:])

	header, err := tree.Header(scheme, notary.Digest{0x07}, notarytest.Unix(t, 1_790_000_000), blinding)
	require.NoError(t, err)

	assignment, err := plonk.NewAssignment(productionEpoch, plonk.ProductionDepth, entries, header, blinding)
	require.NoError(t, err)

	return assignment, header, tree, entries
}

func productionProver(t *testing.T, srsPath string) *plonk.Prover {
	t.Helper()

	file, err := os.Open(filepath.Clean(srsPath))
	require.NoError(t, err)

	defer func() { _ = file.Close() }()

	hash := sha256.New()

	var srs kzgbn254.SRS

	_, err = srs.ReadFrom(io.TeeReader(file, hash))
	require.NoError(t, err)
	require.Equal(t, plonk.ProductionSRSSHA256, hex.EncodeToString(hash.Sum(nil)))

	system, err := plonk.Compile(plonk.ProductionDepth)
	require.NoError(t, err)

	lagrange, err := plonk.NewLagrange(&srs, system)
	require.NoError(t, err)

	prover, err := plonk.NewProver(plonk.ProductionDepth, system, &srs, lagrange)
	require.NoError(t, err)

	key := sha256.New()
	require.NoError(t, prover.WriteVerifyingKey(key))
	require.Equal(t, plonk.ProductionVerifyingKeySHA256, hex.EncodeToString(key.Sum(nil)))

	var exported bytes.Buffer
	require.NoError(t, prover.ExportVerifier(&exported))

	committed, err := os.ReadFile(productionVerifierPath)
	require.NoError(t, err)
	require.True(
		t,
		bytes.Equal(committed, exported.Bytes()),
		"src/confidential/PlonkVerifier.sol is not what the SRS gives",
	)

	return prover
}
