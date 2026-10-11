package fixtures_test

import (
	"bytes"
	"encoding/json"
	"flag"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/keccak"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

var update = flag.Bool("update", false, "rewrite the fixtures the Solidity differential suite reads")

var fixturePath = filepath.Join("..", "..", "contracts", "test", "fixtures", "scheme-v1.json")

type fixture struct {
	Scheme         uint16      `json:"scheme"`
	TreeCount      int         `json:"treeCount"`
	Trees          []treeCase  `json:"trees"`
	RejectionCount int         `json:"rejectionCount"`
	Rejections     []rejection `json:"rejections"`
	ChainLength    int         `json:"chainLength"`
	Chain          []chainLink `json:"chain"`
}

type treeCase struct {
	Epoch       string      `json:"epoch"`
	LeafCount   string      `json:"leafCount"`
	Depth       int         `json:"depth"`
	Root        string      `json:"root"`
	Total       string      `json:"total"`
	Commitments []string    `json:"commitments"`
	Amounts     []string    `json:"amounts"`
	Proofs      []proofCase `json:"proofs"`
}

type proofCase struct {
	Index         string   `json:"index"`
	SiblingHashes []string `json:"siblingHashes"`
	SiblingSums   []string `json:"siblingSums"`
}

type rejection struct {
	Reason     string    `json:"reason"`
	Tree       int       `json:"tree"`
	Commitment string    `json:"commitment"`
	Amount     string    `json:"amount"`
	Proof      proofCase `json:"proof"`
}

type chainLink struct {
	Number         string `json:"number"`
	PreviousDigest string `json:"previousDigest"`
	Root           string `json:"root"`
	Total          string `json:"total"`
	LeafCount      string `json:"leafCount"`
	AsOf           string `json:"asOf"`
	Digest         string `json:"digest"`
}

func u64(v uint64) string {
	return strconv.FormatUint(v, 10)
}

func encodeProof(proof notary.InclusionProof) proofCase {
	out := proofCase{
		Index:         u64(proof.Index),
		SiblingHashes: make([]string, len(proof.Siblings)),
		SiblingSums:   make([]string, len(proof.Siblings)),
	}
	for i, s := range proof.Siblings {
		out.SiblingHashes[i] = s.Hash.Hex()
		out.SiblingSums[i] = u64(s.Sum)
	}
	return out
}

func buildFixture(t *testing.T) fixture {
	t.Helper()
	scheme := keccak.NewScheme()
	out := fixture{Scheme: scheme.Version()}

	shapes := []struct {
		epoch uint64
		count int
	}{
		{0, 0}, {0, 1}, {1, 2}, {2, 3}, {3, 4}, {4, 5}, {5, 7}, {6, 8},
		{7, 9}, {1 << 40, 16}, {math.MaxUint64 - 1, 17}, {12, 33},
	}
	for i, shape := range shapes {
		entries := notarytest.Entries(uint64(i)+100, shape.count)
		tree, err := notary.BuildTree(scheme, shape.epoch, entries)
		require.NoError(t, err)
		root := tree.Root()
		tc := treeCase{
			Epoch:       u64(shape.epoch),
			LeafCount:   u64(tree.LeafCount()),
			Depth:       tree.Depth(),
			Root:        root.Hash.Hex(),
			Total:       u64(root.Sum),
			Commitments: make([]string, len(entries)),
			Amounts:     make([]string, len(entries)),
			Proofs:      make([]proofCase, len(entries)),
		}
		for j, e := range entries {
			proof, err := tree.Prove(uint64(j))
			require.NoError(t, err)
			tc.Commitments[j] = e.Commitment.Hex()
			tc.Amounts[j] = u64(e.Amount)
			tc.Proofs[j] = encodeProof(proof)
		}
		out.Trees = append(out.Trees, tc)
	}
	out.TreeCount = len(out.Trees)

	const target = 5
	entries := notarytest.Entries(uint64(target)+100, shapes[target].count)
	tree, err := notary.BuildTree(scheme, shapes[target].epoch, entries)
	require.NoError(t, err)
	anchored := notary.AnchoredEpoch{
		Number:    shapes[target].epoch,
		Root:      tree.Root().Hash,
		Total:     tree.Root().Sum,
		LeafCount: tree.LeafCount(),
	}
	fresh := func() (notary.Entry, notary.InclusionProof) {
		proof, err := tree.Prove(2)
		require.NoError(t, err)
		siblings := append([]notary.Node(nil), proof.Siblings...)
		return entries[2], notary.InclusionProof{Index: proof.Index, Siblings: siblings}
	}
	tampers := []struct {
		reason string
		apply  func(*notary.Entry, *notary.InclusionProof)
	}{
		{"inflated amount", func(e *notary.Entry, _ *notary.InclusionProof) { e.Amount++ }},
		{"foreign commitment", func(e *notary.Entry, _ *notary.InclusionProof) { e.Commitment[0] ^= 0xff }},
		{"sibling value moved onto the amount", func(e *notary.Entry, p *notary.InclusionProof) {
			e.Amount += p.Siblings[0].Sum
			p.Siblings[0].Sum = 0
		}},
		{"forged sibling hash", func(_ *notary.Entry, p *notary.InclusionProof) { p.Siblings[1].Hash[7] ^= 0x01 }},
		{"wrong index", func(_ *notary.Entry, p *notary.InclusionProof) { p.Index = 3 }},
		{"index in the padding", func(_ *notary.Entry, p *notary.InclusionProof) { p.Index = 6 }},
		{"proof one level short", func(_ *notary.Entry, p *notary.InclusionProof) {
			p.Siblings = p.Siblings[:len(p.Siblings)-1]
		}},
		{"sums that overflow", func(_ *notary.Entry, p *notary.InclusionProof) {
			for i := range p.Siblings {
				p.Siblings[i].Sum = math.MaxUint64
			}
		}},
	}
	for _, tamper := range tampers {
		entry, proof := fresh()
		tamper.apply(&entry, &proof)
		err := notary.VerifyInclusion(scheme, anchored, entry, proof)
		require.Equal(t, notary.KindProofRejected, notary.KindOf(err), tamper.reason)
		out.Rejections = append(out.Rejections, rejection{
			Reason:     tamper.reason,
			Tree:       target,
			Commitment: entry.Commitment.Hex(),
			Amount:     u64(entry.Amount),
			Proof:      encodeProof(proof),
		})
	}
	out.RejectionCount = len(out.Rejections)

	previous := notary.Digest{}
	asOfSeconds := uint64(1760000000)
	for number := uint64(0); number < 6; number++ {
		chainTree, err := notary.BuildTree(scheme, number, notarytest.Entries(number+500, int(number)*3))
		require.NoError(t, err)
		header, err := chainTree.Header(previous, notarytest.Unix(t, asOfSeconds))
		require.NoError(t, err)
		digest := scheme.EpochDigest(header)
		out.Chain = append(out.Chain, chainLink{
			Number:         u64(header.Number),
			PreviousDigest: header.PreviousDigest.Hex(),
			Root:           header.Root.Hex(),
			Total:          u64(header.Total),
			LeafCount:      u64(header.LeafCount),
			AsOf:           u64(header.AsOf.Seconds()),
			Digest:         digest.Hex(),
		})
		previous = digest
		asOfSeconds += 3600
	}
	out.ChainLength = len(out.Chain)
	return out
}

func TestTheSolidityFixturesMatchTheGoScheme(t *testing.T) {
	encoded, err := json.MarshalIndent(buildFixture(t), "", "  ")
	require.NoError(t, err)
	encoded = append(encoded, '\n')

	if *update {
		require.NoError(t, os.WriteFile(fixturePath, encoded, 0o644))
		return
	}
	committed, err := os.ReadFile(fixturePath)
	require.NoError(t, err, "no fixture yet: run `bin/dev go test ./tests/fixtures -update`")
	require.True(t, bytes.Equal(committed, encoded),
		"%s is out of date with the Go scheme: run `bin/dev go test ./tests/fixtures -update`", fixturePath)
}
