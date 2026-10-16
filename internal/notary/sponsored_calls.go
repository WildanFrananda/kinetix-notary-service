package notary

type SponsoredCalls interface {
	Execute(call SponsoredCall) []byte
	RequestStatement(withholding Address, issued IssuedRegistration, epoch uint64) SponsoredCall
	RecordUnanswered(withholding Address, issued IssuedRegistration, epoch uint64) SponsoredCall
	ProveBadProof(court Address, statement Statement, signature Signature, proof InclusionProof) SponsoredCall
	ProveEquivocation(
		court Address,
		first Statement,
		firstSignature Signature,
		second Statement,
		secondSignature Signature,
	) SponsoredCall
}
