package notary

import "math/bits"

const MaxDepth = 32

const MaxLeaves uint64 = 1 << MaxDepth

func DepthFor(leafCount uint64) int {
	if leafCount <= 1 {
		return 0
	}

	return bits.Len64(leafCount - 1)
}
