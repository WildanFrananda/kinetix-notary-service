package notary

type Statement struct {
	Ledger       Address
	Epoch        uint64
	Commitment   Digest
	Amount       uint64
	Index        uint64
	SiblingsHash Digest
}

func StatementOf(scheme StatementScheme, claim Claim) Statement {
	return Statement{
		Ledger:       claim.Ledger.Address,
		Epoch:        claim.Epoch,
		Commitment:   claim.Entry.Commitment,
		Amount:       claim.Entry.Amount,
		Index:        claim.Proof.Index,
		SiblingsHash: scheme.SiblingsHash(claim.Proof.Siblings),
	}
}
