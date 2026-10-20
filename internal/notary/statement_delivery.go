package notary

type StatementDelivery struct {
	Holder   Holder
	Epoch    uint64
	Attempts uint32
}
