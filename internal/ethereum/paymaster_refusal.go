package ethereum

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/WildanFrananda/kinetix-notary-service/internal/contractgen/paymaster"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

func (p *Paymaster) refusal(op string, raw []byte) error {
	decoded, err := p.binding.UnpackError(raw)
	if err != nil {
		return notary.Errf(
			notary.KindNotSponsored,
			op,
			nil,
			"the paymaster refuses the call with %s",
			hexutil.Encode(raw),
		)
	}

	switch refused := decoded.(type) {
	case *paymaster.NotaryPaymasterEpochNotAnchored:
		return notary.Errf(notary.KindEpochNotAnchored, op, nil, "epoch %d is not anchored", refused.Epoch)
	case *paymaster.NotaryPaymasterBeforeRegistration:
		return refusedByCourt(
			op,
			"epoch %d comes before the registration, which starts at epoch %d",
			refused.Epoch,
			refused.FromEpoch,
		)
	case *paymaster.NotaryPaymasterRequestStillOpen:
		return refusedByCourt(
			op,
			"this registration's request for epoch %d is still open: it is answered, or recorded unanswered, first",
			refused.Epoch,
		)
	case *paymaster.NotaryPaymasterAlreadyRequested:
		return refusedByCourt(op, "this registration already requested a statement for epoch %d", refused.Epoch)
	case *paymaster.NotaryPaymasterNoOpenRequest:
		return refusedByCourt(op, "this registration has no open request for epoch %d", refused.Epoch)
	case *paymaster.NotaryPaymasterWrongLedger:
		return refusedByCourt(
			op,
			"the statement names the ledger %s, the court judges %s",
			refused.Given.Hex(),
			refused.Expected.Hex(),
		)
	case *paymaster.NotaryPaymasterSiblingsMismatch:
		return refusedByCourt(op, "the siblings given are not the ones the statement was signed over")
	case *paymaster.NotaryPaymasterProofHolds:
		return refusedByCourt(op, "the proof for epoch %d holds against the anchored root", refused.Epoch)
	case *paymaster.NotaryPaymasterDifferentEpochs:
		return refusedByCourt(op, "the statements are for epochs %d and %d, not one", refused.First, refused.Second)
	case *paymaster.NotaryPaymasterNoConflict:
		return refusedByCourt(op, "the two statements do not contradict each other")
	case *paymaster.NotaryPaymasterAlreadyProven:
		return refusedByCourt(
			op,
			"the court already holds this fault, evidence %s",
			common.Hash(refused.Evidence).Hex(),
		)
	case *paymaster.NotaryPaymasterAnotherSignersEpoch:
		return notSponsored(
			op,
			"epoch %d is %s's to sign, not the signer this paymaster sponsors",
			refused.Epoch,
			refused.Signer.Hex(),
		)
	case *paymaster.NotaryPaymasterNotTheHolder:
		return notSponsored(
			op,
			"the registration is %s's, the operation is %s's",
			refused.Holder.Hex(),
			refused.Sender.Hex(),
		)
	case *paymaster.NotaryPaymasterNotSignedByTheNotary:
		return notSponsored(op, "the notary did not sign it: the signature recovers to %s", refused.Recovered.Hex())
	case *paymaster.NotaryPaymasterCallNotSponsored:
		return notSponsored(
			op,
			"the paymaster sponsors no call %s on %s",
			hexutil.Encode(refused.Selector[:]),
			refused.Target.Hex(),
		)
	case *paymaster.NotaryPaymasterTargetNotSponsored:
		return notSponsored(op, "the paymaster sponsors no call to %s", refused.Target.Hex())
	default:
		return notSponsored(op, "the paymaster refuses the call with %s", hexutil.Encode(raw))
	}
}

func refusedByCourt(op, format string, args ...any) error {
	return notary.Errf(notary.KindCallRefused, op, nil, "the court would refuse it: "+format, args...)
}

func notSponsored(op, format string, args ...any) error {
	return notary.Errf(notary.KindNotSponsored, op, nil, format, args...)
}
