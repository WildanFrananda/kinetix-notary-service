package notary

type TypedDataDomain struct {
	Name              string
	Version           string
	ChainID           uint64
	VerifyingContract Address
}
