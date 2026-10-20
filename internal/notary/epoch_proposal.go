package notary

type EpochProposal struct {
	SnapshotID   SnapshotID
	Header       EpochHeader
	Accounts     []ProposedAccount
	Confidential *ConfidentialPlan
}
