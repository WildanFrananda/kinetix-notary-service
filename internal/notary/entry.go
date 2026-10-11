package notary

type Entry struct {
	Commitment Digest
	Amount     uint64
}
