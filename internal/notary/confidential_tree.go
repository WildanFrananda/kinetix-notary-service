package notary

import "math/bits"

type ConfidentialTree struct {
	epoch     uint64
	leafCount uint64
	total     uint64
	levels    [][]Digest
}

func BuildConfidentialTree(
	scheme ConfidentialScheme,
	epoch uint64,
	depth int,
	entries []Entry,
) (*ConfidentialTree, error) {
	const op = "notary.BuildConfidentialTree"
	if depth < 0 || depth > MaxConfidentialDepth {
		return nil, Errf(KindInvalidArgument, op, nil, "depth %d is outside 0..%d", depth, MaxConfidentialDepth)
	}

	capacity := uint64(1) << depth

	count := uint64(len(entries))
	if count > capacity {
		return nil, Errf(KindTooManyLeaves, op, nil, "%d accounts exceed the generation's %d", count, capacity)
	}

	var total uint64

	base := make([]Digest, capacity)
	for i := range base {
		if uint64(i) >= count {
			base[i] = scheme.EmptyHash()
			continue
		}

		sum, carry := bits.Add64(total, entries[i].Amount, 0)
		if carry != 0 {
			return nil, Errf(KindOverflow, op, nil, "the total of epoch %d does not fit in 64 bits", epoch)
		}

		total = sum
		base[i] = scheme.LeafHash(epoch, entries[i])
	}

	levels := make([][]Digest, 0, depth+1)

	levels = append(levels, base)
	for level := 1; level <= depth; level++ {
		below := levels[level-1]

		above := make([]Digest, len(below)/2)
		for i := range above {
			node, err := scheme.NodeHash(below[2*i], below[2*i+1])
			if err != nil {
				return nil, Errf(KindInvalidArgument, op, err, "level %d of epoch %d", level, epoch)
			}

			above[i] = node
		}

		levels = append(levels, above)
	}

	return &ConfidentialTree{epoch: epoch, leafCount: count, total: total, levels: levels}, nil
}

func (t *ConfidentialTree) Epoch() uint64 {
	return t.epoch
}

func (t *ConfidentialTree) LeafCount() uint64 {
	return t.leafCount
}

func (t *ConfidentialTree) Total() uint64 {
	return t.total
}

func (t *ConfidentialTree) Depth() int {
	return len(t.levels) - 1
}

func (t *ConfidentialTree) TreeRoot() Digest {
	return t.levels[len(t.levels)-1][0]
}

func (t *ConfidentialTree) Prove(index uint64) (InclusionProof, error) {
	if index >= t.leafCount {
		return InclusionProof{}, Errf(
			KindInvalidArgument,
			"notary.ConfidentialTree.Prove",
			nil,
			"index %d is outside the %d accounts of epoch %d",
			index,
			t.leafCount,
			t.epoch,
		)
	}

	siblings := make([]Node, t.Depth())

	position := index
	for level := range siblings {
		siblings[level] = Node{Hash: t.levels[level][position^1]}
		position >>= 1
	}

	return InclusionProof{Index: index, Siblings: siblings}, nil
}

func (t *ConfidentialTree) Header(
	scheme ConfidentialScheme,
	previous Digest,
	asOf UnixTime,
	blinding Blinding,
) (ConfidentialHeader, error) {
	const op = "notary.ConfidentialTree.Header"
	if asOf.IsZero() {
		return ConfidentialHeader{}, Errf(KindInvalidArgument, op, nil, "epoch %d has no snapshot time", t.epoch)
	}

	if blinding.IsZero() {
		return ConfidentialHeader{}, Errf(KindInvalidArgument, op, nil, "epoch %d has no blinding", t.epoch)
	}

	treeRoot := t.TreeRoot()
	commitment := scheme.TotalCommitment(t.total, blinding)

	return ConfidentialHeader{
		Header: EpochHeader{
			Number:         t.epoch,
			PreviousDigest: previous,
			Root:           scheme.BindingRoot(treeRoot, commitment),
			LeafCount:      t.leafCount,
			AsOf:           asOf,
		},
		TreeRoot:        treeRoot,
		TotalCommitment: commitment,
	}, nil
}
