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
	default:
		return "unknown"
	}
}
