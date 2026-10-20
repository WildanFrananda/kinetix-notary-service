package notary

import "context"

type ProofChain interface {
	ProofReader
	PrepareProof(
		ctx context.Context,
		header EpochHeader,
		plan ConfidentialPlan,
		nonce uint64,
		fees FeeQuote,
	) (PreparedAnchor, error)
	Resign(previous PreparedAnchor, fees FeeQuote) (PreparedAnchor, error)
}
