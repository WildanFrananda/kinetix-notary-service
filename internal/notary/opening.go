package notary

type Opening struct {
	Holder Holder
	Kind   AccountKind
	Salt   Salt
}
