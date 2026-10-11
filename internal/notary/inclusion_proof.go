package notary

type InclusionProof struct {
	Index    uint64
	Siblings []Node
}
