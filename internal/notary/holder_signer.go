package notary

type HolderSigner interface {
	Holder() Address
	SignOperation(hash Digest) (Signature, error)
	Authorize(chainID uint64, delegate Address, nonce uint64) (Delegation, error)
}
