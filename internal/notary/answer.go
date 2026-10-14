package notary

type Answer struct {
	Registration Registration
	Epoch        uint64
	Salt         Salt
	Amount       uint64
	Proof        InclusionProof
}
