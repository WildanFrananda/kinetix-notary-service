package fixtures_test

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"math"
	"math/big"
	"os"
	"path/filepath"
	"testing"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/poseidon"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

var confidentialFixturePath = filepath.Join("..", "..", "contracts", "test", "fixtures", "scheme-v2.json")

type confidentialFixture struct {
	Scheme           uint16                 `json:"scheme"`
	Empty            string                 `json:"empty"`
	CompressionCount int                    `json:"compressionCount"`
	Compressions     []compressionCase      `json:"compressions"`
	LeafCount        int                    `json:"leafCount"`
	Leaves           []confidentialLeaf     `json:"leaves"`
	BindingCount     int                    `json:"bindingCount"`
	Bindings         []bindingCase          `json:"bindings"`
	TreeCount        int                    `json:"treeCount"`
	Trees            []confidentialTreeCase `json:"trees"`
	RejectionCount   int                    `json:"rejectionCount"`
	Rejections       []rejection            `json:"rejections"`
	ChainLength      int                    `json:"chainLength"`
	Chain            []chainLink            `json:"chain"`
}

type compressionCase struct {
	Left  string `json:"left"`
	Right string `json:"right"`
	Out   string `json:"out"`
}

type confidentialLeaf struct {
	Epoch      string `json:"epoch"`
	Commitment string `json:"commitment"`
	Amount     string `json:"amount"`
	Hash       string `json:"hash"`
}

type bindingCase struct {
	TreeRoot        string `json:"treeRoot"`
	TotalCommitment string `json:"totalCommitment"`
	Root            string `json:"root"`
}

type confidentialTreeCase struct {
	Epoch       string      `json:"epoch"`
	LeafCount   string      `json:"leafCount"`
	Depth       int         `json:"depth"`
	TreeRoot    string      `json:"treeRoot"`
	Commitments []string    `json:"commitments"`
	Amounts     []string    `json:"amounts"`
	Proofs      []proofCase `json:"proofs"`
}

func fieldElement(label string, i int) notary.Digest {
	seed := sha256.Sum256([]byte(label + "/" + u64(uint64(i))))

	var element fr.Element
	element.SetBytes(seed[:])

	return notary.Digest(element.Bytes())
}

func fieldMaximum() notary.Digest {
	var element fr.Element
	element.SetOne()
	element.Neg(&element)

	return notary.Digest(element.Bytes())
}

func fieldModulus() notary.Digest {
	return notary.Digest(fr.Modulus().FillBytes(make([]byte, 32)))
}

func buildConfidentialFixture(t *testing.T) confidentialFixture {
	t.Helper()

	scheme := poseidon.NewScheme()
	out := confidentialFixture{Scheme: scheme.Version(), Empty: scheme.EmptyHash().Hex()}

	maximum := fieldMaximum()
	pairs := [][2]notary.Digest{{{}, {}}, {maximum, maximum}, {{}, maximum}, {maximum, {}}}

	for i := range 20 {
		pairs = append(pairs, [2]notary.Digest{
			fieldElement("kinetix.notary.compress.left", i),
			fieldElement("kinetix.notary.compress.right", i),
		})
	}

	for _, pair := range pairs {
		node, err := scheme.NodeHash(pair[0], pair[1])
		require.NoError(t, err)

		out.Compressions = append(out.Compressions, compressionCase{
			Left:  pair[0].Hex(),
			Right: pair[1].Hex(),
			Out:   node.Hex(),
		})
	}

	out.CompressionCount = len(out.Compressions)

	epochs := []uint64{0, 1, 7, 1 << 40, math.MaxUint64}
	amounts := []uint64{0, 1, 1_000_000, math.MaxUint64}

	for i, epoch := range epochs {
		for j, amount := range amounts {
			commitment := notary.Digest(sha256.Sum256([]byte("kinetix.notary.leaf/" + u64(uint64(i*len(amounts)+j)))))
			entry := notary.Entry{Commitment: commitment, Amount: amount}

			out.Leaves = append(out.Leaves, confidentialLeaf{
				Epoch:      u64(epoch),
				Commitment: commitment.Hex(),
				Amount:     u64(amount),
				Hash:       scheme.LeafHash(epoch, entry).Hex(),
			})
		}
	}

	out.LeafCount = len(out.Leaves)

	for i := range 6 {
		treeRoot := fieldElement("kinetix.notary.binding.tree", i)
		commitment := fieldElement("kinetix.notary.binding.total", i)

		out.Bindings = append(out.Bindings, bindingCase{
			TreeRoot:        treeRoot.Hex(),
			TotalCommitment: commitment.Hex(),
			Root:            scheme.BindingRoot(treeRoot, commitment).Hex(),
		})
	}

	out.BindingCount = len(out.Bindings)

	shapes := []struct {
		epoch uint64
		count int
		depth int
	}{
		{0, 0, 0},
		{1, 1, 0},
		{2, 1, 2},
		{3, 3, 2},
		{4, 4, 2},
		{5, 5, 3},
		{6, 8, 3},
		{7, 9, 4},
		{1 << 40, 16, 4},
		{math.MaxUint64 - 1, 11, 4},
	}

	trees := make([]*notary.ConfidentialTree, len(shapes))
	treeEntries := make([][]notary.Entry, len(shapes))

	for i, shape := range shapes {
		entries := notarytest.Entries(uint64(i)+300, shape.count)

		tree, err := notary.BuildConfidentialTree(scheme, shape.epoch, shape.depth, entries)
		require.NoError(t, err)

		trees[i] = tree
		treeEntries[i] = entries

		tc := confidentialTreeCase{
			Epoch:       u64(shape.epoch),
			LeafCount:   u64(tree.LeafCount()),
			Depth:       tree.Depth(),
			TreeRoot:    tree.TreeRoot().Hex(),
			Commitments: make([]string, len(entries)),
			Amounts:     make([]string, len(entries)),
			Proofs:      make([]proofCase, len(entries)),
		}
		for j, entry := range entries {
			proof, err := tree.Prove(uint64(j))
			require.NoError(t, err)

			tc.Commitments[j] = entry.Commitment.Hex()
			tc.Amounts[j] = u64(entry.Amount)
			tc.Proofs[j] = encodeProof(proof)
		}

		out.Trees = append(out.Trees, tc)
	}

	out.TreeCount = len(out.Trees)

	const target = 7

	tree, entries := trees[target], treeEntries[target]
	epoch := notary.ConfidentialEpoch{
		Number:    tree.Epoch(),
		TreeRoot:  tree.TreeRoot(),
		LeafCount: tree.LeafCount(),
		Depth:     tree.Depth(),
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
		{"forged sibling hash", func(_ *notary.Entry, p *notary.InclusionProof) { p.Siblings[1].Hash[31] ^= 0x01 }},
		{"wrong index", func(_ *notary.Entry, p *notary.InclusionProof) { p.Index = 3 }},
		{"index in the padding", func(_ *notary.Entry, p *notary.InclusionProof) { p.Index = 9 }},
		{"proof one level short", func(_ *notary.Entry, p *notary.InclusionProof) {
			p.Siblings = p.Siblings[:len(p.Siblings)-1]
		}},
		{"sibling carrying a sum", func(_ *notary.Entry, p *notary.InclusionProof) { p.Siblings[0].Sum = 1 }},
		{"sibling that is not a field element", func(_ *notary.Entry, p *notary.InclusionProof) {
			p.Siblings[2].Hash = fieldModulus()
		}},
		{"sibling aliased by the modulus", func(_ *notary.Entry, p *notary.InclusionProof) {
			value := new(big.Int).SetBytes(p.Siblings[0].Hash[:])
			value.Add(value, fr.Modulus())
			p.Siblings[0].Hash = notary.Digest(value.FillBytes(make([]byte, 32)))
		}},
	}
	for _, tamper := range tampers {
		entry, proof := fresh()
		tamper.apply(&entry, &proof)

		err := notary.VerifyConfidentialInclusion(scheme, epoch, entry, proof)
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
	asOfSeconds := uint64(1_790_000_000)

	for number := uint64(0); number < 6; number++ {
		chainTree, err := notary.BuildConfidentialTree(scheme, number, 3, notarytest.Entries(number+700, int(number)+1))
		require.NoError(t, err)

		seed := sha256.Sum256([]byte("kinetix.notary.blinding/" + u64(number)))

		var blinding notary.Blinding
		copy(blinding[:], seed[:])

		header, err := chainTree.Header(scheme, previous, notarytest.Unix(t, asOfSeconds), blinding)
		require.NoError(t, err)

		digest := scheme.EpochDigest(header.Header)
		out.Chain = append(out.Chain, chainLink{
			Number:         u64(header.Header.Number),
			PreviousDigest: header.Header.PreviousDigest.Hex(),
			Root:           header.Header.Root.Hex(),
			Total:          u64(header.Header.Total),
			LeafCount:      u64(header.Header.LeafCount),
			AsOf:           u64(header.Header.AsOf.Seconds()),
			Digest:         digest.Hex(),
		})
		previous = digest
		asOfSeconds += 3600
	}

	out.ChainLength = len(out.Chain)

	return out
}

func TestTheSolidityFixturesMatchTheGoConfidentialScheme(t *testing.T) {
	encoded, err := json.MarshalIndent(buildConfidentialFixture(t), "", "  ")
	require.NoError(t, err)

	encoded = append(encoded, '\n')

	if *update {
		require.NoError(t, os.WriteFile(confidentialFixturePath, encoded, 0o644))
		return
	}

	committed, err := os.ReadFile(confidentialFixturePath)
	require.NoError(t, err, "no fixture yet: run `bin/dev go test ./tests/fixtures -update`")
	require.True(
		t,
		bytes.Equal(committed, encoded),
		"contracts/test/fixtures/scheme-v2.json is stale: run `bin/dev go test ./tests/fixtures -update`",
	)
}
