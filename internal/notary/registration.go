package notary

type Registration struct {
	Holder    Address
	Account   Digest
	FromEpoch uint64
}
