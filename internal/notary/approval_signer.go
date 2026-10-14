package notary

type ApprovalSigner interface {
	Member() Address
	Sign(digest Digest) (Signature, error)
}
