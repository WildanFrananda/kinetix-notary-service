package notary

type WithholdingRequest struct {
	Registration Digest
	Account      Digest
	Epoch        uint64
	Holder       Address
	DueAt        UnixTime
}
