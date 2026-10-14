package notary

type SignatureRecoverer interface {
	Recover(digest Digest, signature Signature) (Address, error)
}
