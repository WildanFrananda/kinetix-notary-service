package verifying

import (
	"context"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type Verifier struct {
	ledger notary.Ledger
	scheme notary.Scheme
}

func NewVerifier(ledger notary.Ledger, scheme notary.Scheme) *Verifier {
	return &Verifier{ledger: ledger, scheme: scheme}
}

func (v *Verifier) Verify(ctx context.Context, claim notary.Claim) (notary.Verification, error) {
	const op = "verifying.Verifier.Verify"

	identity, err := v.ledger.Identity(ctx)
	if err != nil {
		return notary.Verification{}, err
	}

	if identity != claim.Ledger {
		return notary.Verification{}, notary.Errf(
			notary.KindWrongLedger,
			op,
			nil,
			"the proof was issued against %s on chain %d, this verifier reads %s on chain %d",
			claim.Ledger.Address.Hex(),
			claim.Ledger.ChainID,
			identity.Address.Hex(),
			identity.ChainID,
		)
	}

	epoch, err := v.ledger.Epoch(ctx, claim.Epoch)
	if err != nil {
		return notary.Verification{}, err
	}

	if opening := claim.Opening; opening != nil {
		if v.scheme.Commitment(opening.Holder, opening.Kind, opening.Salt) != claim.Entry.Commitment {
			return notary.Verification{
				Verdict: notary.VerdictRejected,
				Ledger:  identity,
				Epoch:   epoch,
				Reason:  "the account named in the proof does not open its commitment",
				Opening: opening,
			}, nil
		}
	}

	local := notary.VerifyInclusion(v.scheme, epoch, claim.Entry, claim.Proof)
	if local != nil && notary.KindOf(local) != notary.KindProofRejected {
		return notary.Verification{}, local
	}

	confirmed, err := v.ledger.ConfirmsInclusion(ctx, claim)
	if err != nil {
		return notary.Verification{}, err
	}

	if confirmed != (local == nil) {
		return notary.Verification{}, notary.Errf(
			notary.KindDisagreement,
			op,
			local,
			"the ledger answers %t for epoch %d, this verifier's scheme answers %t",
			confirmed,
			claim.Epoch,
			local == nil,
		)
	}

	verification := notary.Verification{
		Verdict: notary.VerdictIncluded,
		Ledger:  identity,
		Epoch:   epoch,
		Opening: claim.Opening,
	}
	if local != nil {
		verification.Verdict = notary.VerdictRejected
		verification.Reason = local.Error()
	}

	return verification, nil
}
