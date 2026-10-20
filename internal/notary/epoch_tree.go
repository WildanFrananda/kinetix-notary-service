package notary

type EpochTree interface {
	Prove(index uint64) (InclusionProof, error)
}
