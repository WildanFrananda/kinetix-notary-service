package anchoring

import "github.com/WildanFrananda/kinetix-notary-service/internal/notary"

type SchemeTwo struct {
	Scheme    notary.ConfidentialScheme
	Depth     int
	Blindings notary.BlindingSource
	Prover    notary.EpochProver
	Chain     notary.ProofChain
}
