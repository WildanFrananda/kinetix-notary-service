package notary

type Kind uint8

const (
	KindUnknown Kind = iota
	KindInvalidArgument
	KindOverflow
	KindTooManyLeaves
	KindProofRejected
	KindEpochNotAnchored
	KindEpochNotFinal
	KindWrongLedger
	KindLedgerUnavailable
	KindDisagreement
	KindSourceUnavailable
	KindSourceRefused
	KindStoreUnavailable
	KindLeaseLost
	KindAnchorDeferred
	KindAnchorRefused
	KindAlreadyAnchored
	KindLedgerDiverged
	KindNotAnchorer
	KindEpochAbandoned
	KindUnauthenticated
	KindSignerUnavailable
	KindWrongCourt
)

func (k Kind) String() string {
	switch k {
	case KindInvalidArgument:
		return "invalid_argument"
	case KindOverflow:
		return "overflow"
	case KindTooManyLeaves:
		return "too_many_leaves"
	case KindProofRejected:
		return "proof_rejected"
	case KindEpochNotAnchored:
		return "epoch_not_anchored"
	case KindEpochNotFinal:
		return "epoch_not_final"
	case KindWrongLedger:
		return "wrong_ledger"
	case KindLedgerUnavailable:
		return "ledger_unavailable"
	case KindDisagreement:
		return "disagreement"
	case KindSourceUnavailable:
		return "source_unavailable"
	case KindSourceRefused:
		return "source_refused"
	case KindStoreUnavailable:
		return "store_unavailable"
	case KindLeaseLost:
		return "lease_lost"
	case KindAnchorDeferred:
		return "anchor_deferred"
	case KindAnchorRefused:
		return "anchor_refused"
	case KindAlreadyAnchored:
		return "already_anchored"
	case KindLedgerDiverged:
		return "ledger_diverged"
	case KindNotAnchorer:
		return "not_anchorer"
	case KindEpochAbandoned:
		return "epoch_abandoned"
	case KindUnauthenticated:
		return "unauthenticated"
	case KindSignerUnavailable:
		return "signer_unavailable"
	case KindWrongCourt:
		return "wrong_court"
	default:
		return "unknown"
	}
}
