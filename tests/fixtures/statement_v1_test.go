package fixtures_test

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/ethereum"
	"github.com/WildanFrananda/kinetix-notary-service/internal/keccak"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

var statementFixturePath = filepath.Join("..", "..", "contracts", "test", "fixtures", "statement-v1.json")

type statementFixture struct {
	Name           string          `json:"name"`
	Version        string          `json:"version"`
	ChainID        string          `json:"chainId"`
	Court          string          `json:"court"`
	Signer         string          `json:"signer"`
	StatementCount int             `json:"statementCount"`
	Statements     []statementCase `json:"statements"`
}

type statementCase struct {
	Ledger        string   `json:"ledger"`
	Epoch         string   `json:"epoch"`
	Commitment    string   `json:"commitment"`
	Amount        string   `json:"amount"`
	Index         string   `json:"index"`
	SiblingHashes []string `json:"siblingHashes"`
	SiblingSums   []string `json:"siblingSums"`
	SiblingsHash  string   `json:"siblingsHash"`
	Digest        string   `json:"digest"`
	Signature     string   `json:"signature"`
}

func labelled(label string) notary.Address {
	var address notary.Address
	copy(address[:], crypto.Keccak256([]byte(label))[12:])

	return address
}

func buildStatementFixture(t *testing.T) statementFixture {
	t.Helper()

	scheme := keccak.NewScheme()
	key, err := crypto.ToECDSA(crypto.Keccak256([]byte("kinetix.notary.statement-fixture")))
	require.NoError(t, err)

	signerAddress := notary.Address(crypto.PubkeyToAddress(key.PublicKey))
	signer := ethereum.NewStatementSigner([]*ecdsa.PrivateKey{key})
	court := notary.CourtIdentity{ChainID: 84532, Address: labelled("kinetix.notary.fixture-court")}
	ledger := labelled("kinetix.notary.fixture-ledger")

	out := statementFixture{
		Name:    keccak.CourtName,
		Version: keccak.CourtVersion,
		ChainID: strconv.FormatUint(court.ChainID, 10),
		Court:   court.Address.Hex(),
		Signer:  signerAddress.Hex(),
	}

	add := func(statement notary.Statement, siblings []notary.Node) {
		statement.SiblingsHash = scheme.SiblingsHash(siblings)
		digest := scheme.StatementDigest(court, statement)
		signature, err := signer.Sign(signerAddress, digest)
		require.NoError(t, err)

		public, err := crypto.SigToPub(digest[:], append(signature[:64:64], signature[64]-27))
		require.NoError(t, err)
		require.Equal(t, signerAddress, notary.Address(crypto.PubkeyToAddress(*public)))

		proof := encodeProof(notary.InclusionProof{Index: statement.Index, Siblings: siblings})
		out.Statements = append(out.Statements, statementCase{
			Ledger:        statement.Ledger.Hex(),
			Epoch:         u64(statement.Epoch),
			Commitment:    statement.Commitment.Hex(),
			Amount:        u64(statement.Amount),
			Index:         u64(statement.Index),
			SiblingHashes: proof.SiblingHashes,
			SiblingSums:   proof.SiblingSums,
			SiblingsHash:  statement.SiblingsHash.Hex(),
			Digest:        digest.Hex(),
			Signature:     signature.Hex(),
		})
	}

	shapes := []struct {
		epoch uint64
		count int
	}{
		{0, 1},
		{3, 5},
		{1 << 40, 16},
		{math.MaxUint64 - 1, 17},
	}
	for i, shape := range shapes {
		entries := notarytest.Entries(uint64(i)+700, shape.count)
		tree, err := notary.BuildTree(scheme, shape.epoch, entries)
		require.NoError(t, err)

		for _, index := range []int{0, shape.count / 2, shape.count - 1} {
			proof, err := tree.Prove(uint64(index))
			require.NoError(t, err)

			add(
				notary.Statement{
					Ledger:     ledger,
					Epoch:      shape.epoch,
					Commitment: entries[index].Commitment,
					Amount:     entries[index].Amount,
					Index:      proof.Index,
				},
				proof.Siblings,
			)
		}
	}

	add(
		notary.Statement{
			Ledger:     ledger,
			Epoch:      math.MaxUint64,
			Commitment: notary.Digest{0xff},
			Amount:     math.MaxUint64,
			Index:      math.MaxUint64,
		},
		[]notary.Node{{Hash: notary.Digest{0xff}, Sum: math.MaxUint64}},
	)

	out.StatementCount = len(out.Statements)

	return out
}

func TestTheSolidityStatementFixtureMatchesTheGoScheme(t *testing.T) {
	encoded, err := json.MarshalIndent(buildStatementFixture(t), "", "  ")
	require.NoError(t, err)

	encoded = append(encoded, '\n')

	if *update {
		require.NoError(t, os.WriteFile(statementFixturePath, encoded, 0o644))
		return
	}

	committed, err := os.ReadFile(statementFixturePath)
	require.NoError(t, err, "no fixture yet: run `bin/dev go test ./tests/fixtures -update`")
	require.True(
		t,
		bytes.Equal(committed, encoded),
		"%s is out of date with the Go scheme: run `bin/dev go test ./tests/fixtures -update`",
		statementFixturePath,
	)
}
