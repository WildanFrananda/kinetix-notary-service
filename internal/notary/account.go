package notary

type Account struct {
	Position   uint64
	Liability  Liability
	Salt       Salt
	Secret     AccountSecret
	ID         Digest
	Commitment Digest
}
