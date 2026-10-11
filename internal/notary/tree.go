package notary

type Tree struct {
	epoch     uint64
	leafCount uint64
	levels    [][]Node
}

func BuildTree(scheme Scheme, epoch uint64, entries []Entry) (*Tree, error) {
	const op = "notary.BuildTree"
	count := uint64(len(entries))
	if count > MaxLeaves {
		return nil, Errf(KindTooManyLeaves, op, nil, "%d accounts exceed the scheme's %d", count, MaxLeaves)
	}
	depth := DepthFor(count)
	base := make([]Node, 1<<depth)
	for i := range base {
		if i < len(entries) {
			base[i] = Node{Hash: scheme.LeafHash(epoch, entries[i]), Sum: entries[i].Amount}
			continue
		}
		base[i] = Node{Hash: scheme.EmptyHash()}
	}
	levels := make([][]Node, 0, depth+1)
	levels = append(levels, base)
	for level := 1; level <= depth; level++ {
		below := levels[level-1]
		above := make([]Node, len(below)/2)
		for i := range above {
			node, err := Parent(scheme, below[2*i], below[2*i+1])
			if err != nil {
				return nil, Errf(KindOverflow, op, err, "the total of epoch %d", epoch)
			}
			above[i] = node
		}
		levels = append(levels, above)
	}
	return &Tree{epoch: epoch, leafCount: count, levels: levels}, nil
}

func (t *Tree) Epoch() uint64 {
	return t.epoch
}

func (t *Tree) LeafCount() uint64 {
	return t.leafCount
}

func (t *Tree) Depth() int {
	return len(t.levels) - 1
}

func (t *Tree) Root() Node {
	return t.levels[len(t.levels)-1][0]
}

func (t *Tree) Prove(index uint64) (InclusionProof, error) {
	if index >= t.leafCount {
		return InclusionProof{}, Errf(
			KindInvalidArgument,
			"notary.Tree.Prove",
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
		siblings[level] = t.levels[level][position^1]
		position >>= 1
	}
	return InclusionProof{Index: index, Siblings: siblings}, nil
}

func (t *Tree) Header(previous Digest, asOf UnixTime) (EpochHeader, error) {
	if asOf.IsZero() {
		return EpochHeader{}, Errf(KindInvalidArgument, "notary.Tree.Header", nil, "epoch %d has no snapshot time", t.epoch)
	}
	root := t.Root()
	return EpochHeader{
		Number:         t.epoch,
		PreviousDigest: previous,
		Root:           root.Hash,
		Total:          root.Sum,
		LeafCount:      t.leafCount,
		AsOf:           asOf,
	}, nil
}
