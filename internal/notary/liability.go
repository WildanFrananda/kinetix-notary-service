package notary

type Liability struct {
	Holder Holder
	Kind   AccountKind
	Amount uint64
}
