package notary

type ProposedAccount struct {
	Liability Liability
	Secret    AccountSecret
	Salt      Salt
}
