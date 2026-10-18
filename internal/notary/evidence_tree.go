package notary

type EvidenceTree struct {
	stream    EvidenceStream
	leafCount uint64
	levels    [][]Digest
}

func BuildEvidenceTree(scheme EvidenceScheme, stream EvidenceStream, leaves []Digest) (*EvidenceTree, error) {
	const op = "notary.BuildEvidenceTree"

	count := uint64(len(leaves))
	if count == 0 {
		return nil, Errf(KindInvalidArgument, op, nil, "a %s batch holds at least one record", stream)
	}

	if count > MaxLeaves {
		return nil, Errf(KindTooManyLeaves, op, nil, "%d records exceed the scheme's %d", count, MaxLeaves)
	}

	depth := DepthFor(count)

	base := make([]Digest, 1<<depth)
	for i := range base {
		if i < len(leaves) {
			base[i] = leaves[i]
			continue
		}

		base[i] = scheme.EvidenceEmpty()
	}

	levels := make([][]Digest, 0, depth+1)

	levels = append(levels, base)
	for level := 1; level <= depth; level++ {
		below := levels[level-1]

		above := make([]Digest, len(below)/2)
		for i := range above {
			above[i] = scheme.EvidenceNode(below[2*i], below[2*i+1])
		}

		levels = append(levels, above)
	}

	return &EvidenceTree{stream: stream, leafCount: count, levels: levels}, nil
}

func (t *EvidenceTree) LeafCount() uint64 {
	return t.leafCount
}

func (t *EvidenceTree) Depth() int {
	return len(t.levels) - 1
}

func (t *EvidenceTree) Root() Digest {
	return t.levels[len(t.levels)-1][0]
}

func (t *EvidenceTree) Prove(index uint64) (EvidenceProof, error) {
	if index >= t.leafCount {
		return EvidenceProof{}, Errf(
			KindInvalidArgument,
			"notary.EvidenceTree.Prove",
			nil,
			"index %d is outside the %d records of the %s batch",
			index,
			t.leafCount,
			t.stream,
		)
	}

	siblings := make([]Digest, t.Depth())

	position := index
	for level := range siblings {
		siblings[level] = t.levels[level][position^1]
		position >>= 1
	}

	return EvidenceProof{Index: index, Siblings: siblings}, nil
}
