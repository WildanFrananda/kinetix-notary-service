package notary

type Account struct {
	Position   uint64
	Liability  Liability
	Salt       Salt
	Commitment Digest
}
