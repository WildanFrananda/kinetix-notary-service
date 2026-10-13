package notary

type StatementSigner interface {
	Holds(signer Address) bool
	Sign(signer Address, digest Digest) (Signature, error)
}
