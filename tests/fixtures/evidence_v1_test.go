package fixtures_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/keccak"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

var evidenceFixturePath = filepath.Join("..", "..", "contracts", "test", "fixtures", "evidence-v1.json")

type evidenceFixture struct {
	Scheme         uint16             `json:"scheme"`
	Empty          string             `json:"empty"`
	TreeCount      int                `json:"treeCount"`
	Trees          []evidenceTree     `json:"trees"`
	RejectionCount int                `json:"rejectionCount"`
	Rejections     []evidenceTampered `json:"rejections"`
	ChainLength    int                `json:"chainLength"`
	Chain          []evidenceLink     `json:"chain"`
}

type evidenceTree struct {
	Stream        int             `json:"stream"`
	LeafCount     string          `json:"leafCount"`
	Depth         int             `json:"depth"`
	Root          string          `json:"root"`
	Salts         []string        `json:"salts"`
	RecordDigests []string        `json:"recordDigests"`
	Leaves        []string        `json:"leaves"`
	Proofs        []evidenceProof `json:"proofs"`
}

type evidenceProof struct {
	Index    string   `json:"index"`
	Siblings []string `json:"siblings"`
}

type evidenceTampered struct {
	Reason       string        `json:"reason"`
	Tree         int           `json:"tree"`
	Stream       int           `json:"stream"`
	Salt         string        `json:"salt"`
	RecordDigest string        `json:"recordDigest"`
	Proof        evidenceProof `json:"proof"`
}

type evidenceLink struct {
	Stream             int    `json:"stream"`
	Number             string `json:"number"`
	PreviousDigest     string `json:"previousDigest"`
	Root               string `json:"root"`
	LeafCount          string `json:"leafCount"`
	ThroughTransaction string `json:"throughTransaction"`
	ThroughSequence    string `json:"throughSequence"`
	ClosedAt           string `json:"closedAt"`
	Digest             string `json:"digest"`
}

func encodeEvidenceProof(proof notary.EvidenceProof) evidenceProof {
	out := evidenceProof{Index: u64(proof.Index), Siblings: make([]string, len(proof.Siblings))}
	for i, sibling := range proof.Siblings {
		out.Siblings[i] = sibling.Hex()
	}

	return out
}

func evidenceTreeOf(
	t *testing.T,
	stream notary.EvidenceStream,
	count int,
) (*notary.EvidenceTree, []notary.Salt, []notary.Digest) {
	t.Helper()

	scheme := keccak.NewEvidence()
	salts := make([]notary.Salt, count)
	records := make([]notary.Digest, count)
	leaves := make([]notary.Digest, count)

	for i, record := range notarytest.Deliveries(t, count) {
		digest, err := scheme.RecordDigest(record)
		require.NoError(t, err)

		salts[i] = notary.Salt(digestOf(byte(stream), byte(count), byte(i)))
		records[i] = digest
		leaves[i] = scheme.EvidenceLeaf(stream, salts[i], digest)
	}

	tree, err := notary.BuildEvidenceTree(scheme, stream, leaves)
	require.NoError(t, err)

	return tree, salts, records
}

func digestOf(parts ...byte) notary.Digest {
	var d notary.Digest
	for i := range d {
		d[i] = parts[i%len(parts)] ^ byte(i*37)
	}

	return d
}

func buildEvidenceFixture(t *testing.T) evidenceFixture {
	t.Helper()

	scheme := keccak.NewEvidence()
	out := evidenceFixture{Scheme: scheme.EvidenceVersion(), Empty: scheme.EvidenceEmpty().Hex()}

	shapes := []struct {
		stream notary.EvidenceStream
		count  int
	}{
		{notary.StreamOrderLifecycle, 1},
		{notary.StreamDelivery, 2},
		{notary.StreamReview, 3},
		{notary.StreamOrderLifecycle, 5},
		{notary.StreamDelivery, 8},
		{notary.StreamReview, 9},
		{notary.StreamDelivery, 17},
		{notary.StreamOrderLifecycle, 33},
	}
	for _, shape := range shapes {
		tree, salts, records := evidenceTreeOf(t, shape.stream, shape.count)

		tc := evidenceTree{
			Stream:        int(shape.stream),
			LeafCount:     u64(tree.LeafCount()),
			Depth:         tree.Depth(),
			Root:          tree.Root().Hex(),
			Salts:         make([]string, shape.count),
			RecordDigests: make([]string, shape.count),
			Leaves:        make([]string, shape.count),
			Proofs:        make([]evidenceProof, shape.count),
		}
		for i := range shape.count {
			proof, err := tree.Prove(uint64(i))
			require.NoError(t, err)

			tc.Salts[i] = salts[i].Hex()
			tc.RecordDigests[i] = records[i].Hex()
			tc.Leaves[i] = scheme.EvidenceLeaf(shape.stream, salts[i], records[i]).Hex()
			tc.Proofs[i] = encodeEvidenceProof(proof)
		}

		out.Trees = append(out.Trees, tc)
	}

	out.TreeCount = len(out.Trees)

	const target = 3

	tree, salts, records := evidenceTreeOf(t, shapes[target].stream, shapes[target].count)
	batch := notary.AnchoredBatch{
		Stream:    shapes[target].stream,
		Root:      tree.Root(),
		LeafCount: tree.LeafCount(),
	}
	fresh := func() (notary.EvidenceStream, notary.Salt, notary.Digest, notary.EvidenceProof) {
		proof, err := tree.Prove(2)
		require.NoError(t, err)

		siblings := append([]notary.Digest(nil), proof.Siblings...)

		return shapes[target].stream, salts[2], records[2], notary.EvidenceProof{Index: 2, Siblings: siblings}
	}

	tampers := []struct {
		reason string
		apply  func(*notary.EvidenceStream, *notary.Salt, *notary.Digest, *notary.EvidenceProof)
	}{
		{"another salt", func(_ *notary.EvidenceStream, s *notary.Salt, _ *notary.Digest, _ *notary.EvidenceProof) {
			s[0] ^= 0xff
		}},
		{"another record", func(_ *notary.EvidenceStream, _ *notary.Salt, r *notary.Digest, _ *notary.EvidenceProof) {
			r[31] ^= 0x01
		}},
		{"another stream", func(st *notary.EvidenceStream, _ *notary.Salt, _ *notary.Digest, _ *notary.EvidenceProof) {
			*st = notary.StreamReview
		}},
		{"forged sibling", func(_ *notary.EvidenceStream, _ *notary.Salt, _ *notary.Digest, p *notary.EvidenceProof) {
			p.Siblings[1][7] ^= 0x01
		}},
		{"wrong index", func(_ *notary.EvidenceStream, _ *notary.Salt, _ *notary.Digest, p *notary.EvidenceProof) {
			p.Index = 3
		}},
		{
			"index in the padding",
			func(_ *notary.EvidenceStream, _ *notary.Salt, _ *notary.Digest, p *notary.EvidenceProof) {
				p.Index = 6
			},
		},
		{
			"proof one level short",
			func(_ *notary.EvidenceStream, _ *notary.Salt, _ *notary.Digest, p *notary.EvidenceProof) {
				p.Siblings = p.Siblings[:len(p.Siblings)-1]
			},
		},
	}
	for _, tamper := range tampers {
		stream, salt, record, proof := fresh()
		tamper.apply(&stream, &salt, &record, &proof)

		claimed := batch
		claimed.Stream = stream
		err := notary.VerifyEvidence(scheme, claimed, salt, record, proof)
		require.Equal(t, notary.KindProofRejected, notary.KindOf(err), tamper.reason)

		out.Rejections = append(out.Rejections, evidenceTampered{
			Reason:       tamper.reason,
			Tree:         target,
			Stream:       int(stream),
			Salt:         salt.Hex(),
			RecordDigest: record.Hex(),
			Proof:        encodeEvidenceProof(proof),
		})
	}

	out.RejectionCount = len(out.Rejections)

	for _, stream := range notary.EvidenceStreams() {
		previous := notary.Digest{}
		closedAt := uint64(1_760_000_000)
		through := notary.LogPosition{}

		for number := uint64(0); number < 4; number++ {
			chainTree, _, _ := evidenceTreeOf(t, stream, int(number)*2+1)
			through = notary.LogPosition{Transaction: through.Transaction + number, Sequence: through.Sequence + 3}
			header := notary.EvidenceBatchHeader{
				Stream:         stream,
				Number:         number,
				PreviousDigest: previous,
				Root:           chainTree.Root(),
				LeafCount:      chainTree.LeafCount(),
				Through:        through,
				ClosedAt:       notarytest.Unix(t, closedAt),
			}

			digest := scheme.BatchDigest(header)
			out.Chain = append(out.Chain, evidenceLink{
				Stream:             int(stream),
				Number:             u64(number),
				PreviousDigest:     previous.Hex(),
				Root:               header.Root.Hex(),
				LeafCount:          u64(header.LeafCount),
				ThroughTransaction: u64(through.Transaction),
				ThroughSequence:    u64(through.Sequence),
				ClosedAt:           u64(closedAt),
				Digest:             digest.Hex(),
			})
			previous = digest
			closedAt += 86_400
		}
	}

	out.ChainLength = len(out.Chain)

	return out
}

func TestTheSolidityEvidenceFixturesMatchTheGoScheme(t *testing.T) {
	encoded, err := json.MarshalIndent(buildEvidenceFixture(t), "", "  ")
	require.NoError(t, err)

	encoded = append(encoded, '\n')

	if *update {
		require.NoError(t, os.WriteFile(evidenceFixturePath, encoded, 0o644))
		return
	}

	committed, err := os.ReadFile(evidenceFixturePath)
	require.NoError(t, err, "no fixture yet: run `bin/dev go test ./tests/fixtures -update`")
	require.True(
		t,
		bytes.Equal(committed, encoded),
		"%s is out of date with the Go scheme: run `bin/dev go test ./tests/fixtures -update`",
		evidenceFixturePath,
	)
}
