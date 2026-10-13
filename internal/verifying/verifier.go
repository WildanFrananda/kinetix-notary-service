package verifying

import (
	"context"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type Verifier struct {
	ledger     notary.Ledger
	scheme     notary.Scheme
	statements notary.StatementScheme
	court      notary.Court
}

func NewVerifier(
	ledger notary.Ledger,
	scheme notary.Scheme,
	statements notary.StatementScheme,
	court notary.Court,
) *Verifier {
	return &Verifier{ledger: ledger, scheme: scheme, statements: statements, court: court}
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
			return v.attest(
				ctx,
				claim,
				notary.Verification{
					Verdict: notary.VerdictRejected,
					Ledger:  identity,
					Epoch:   epoch,
					Reason:  "the account named in the proof does not open its commitment",
					Opening: opening,
				},
				false,
			)
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

	return v.attest(ctx, claim, verification, local != nil)
}

func (v *Verifier) attest(
	ctx context.Context,
	claim notary.Claim,
	verification notary.Verification,
	inclusionRefused bool,
) (notary.Verification, error) {
	const op = "verifying.Verifier.attest"

	court, err := v.court.Identity(ctx)
	if err != nil {
		return notary.Verification{}, err
	}

	faults, err := v.court.FaultCount(ctx, claim.Epoch)
	if err != nil {
		return notary.Verification{}, err
	}

	verification.Court = court
	verification.Faults = faults
	verification.Check = notary.StatementAbsent

	sealed := claim.Statement
	if sealed == nil {
		return verification, nil
	}

	if sealed.Court != court.Address {
		return notary.Verification{}, notary.Errf(
			notary.KindWrongCourt,
			op,
			nil,
			"the proof's statement was signed for the court %s, this verifier reads %s",
			sealed.Court.Hex(),
			court.Address.Hex(),
		)
	}

	statement := notary.StatementOf(v.statements, claim)

	onChain, err := v.court.StatementDigest(ctx, statement)
	if err != nil {
		return notary.Verification{}, err
	}

	if local := v.statements.StatementDigest(court, statement); local != onChain {
		return notary.Verification{}, notary.Errf(
			notary.KindDisagreement,
			op,
			nil,
			"the court hashes the statement to %s, this verifier to %s",
			onChain.Hex(),
			local.Hex(),
		)
	}

	signer, err := v.court.SignerOf(ctx, claim.Epoch)
	if err != nil {
		return notary.Verification{}, err
	}

	recovered, err := v.court.RecoverSigner(ctx, statement, sealed.Signature)
	if err != nil {
		return notary.Verification{}, err
	}

	verification.Signer = signer
	verification.Check = notary.StatementUnattributable

	if !signer.IsZero() && recovered == signer {
		verification.Check = notary.StatementAttributable
		verification.Provable = inclusionRefused
	}

	return verification, nil
}
