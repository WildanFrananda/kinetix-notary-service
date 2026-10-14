package notary

type IssuedRegistration struct {
	Registration Registration
	Kind         AccountKind
	Digest       Digest
	Signature    Signature
}
