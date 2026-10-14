package notary

type ApprovalScheme interface {
	ApprovalDigest(council CouncilIdentity, ledger Address, number uint64, epoch Digest) Digest
}
