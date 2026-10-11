package notary

import "math/bits"

func Parent(scheme Scheme, left, right Node) (Node, error) {
	sum, carry := bits.Add64(left.Sum, right.Sum, 0)
	if carry != 0 {
		return Node{}, Errf(KindOverflow, "notary.Parent", nil, "%d + %d does not fit in 64 bits", left.Sum, right.Sum)
	}
	return Node{Hash: scheme.NodeHash(left, right), Sum: sum}, nil
}
