package sponsorcli

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/prooffile"
)

const Usage = "usage: notary-sponsor request|unanswered <registrations.json> <epoch>\n" +
	"       notary-sponsor prove <proof.json>\n" +
	"       notary-sponsor equivocation <proof.json> <proof.json>"

type runner struct {
	deps     Dependencies
	identity notary.PaymasterIdentity
	stdout   io.Writer
}

func Run(ctx context.Context, args []string, deps Dependencies, stdout io.Writer) int {
	if !wellFormed(args) {
		return ExitUsage
	}

	identity, err := deps.Sponsor.Identity(ctx)
	if err != nil {
		return Fail(stdout, args[0], err)
	}

	r := runner{deps: deps, identity: identity, stdout: stdout}

	switch args[0] {
	case "request", "unanswered":
		epoch, err := strconv.ParseUint(args[2], 10, 64)
		if err != nil {
			return emit(
				stdout,
				report{Action: args[0], Verdict: "malformed", Reason: "the epoch is not a number"},
				ExitMalformed,
			)
		}

		return r.registrations(ctx, args[0], args[1], epoch)
	case "prove":
		return r.prove(ctx, args[1])
	default:
		return r.equivocation(ctx, args[1], args[2])
	}
}

func wellFormed(args []string) bool {
	if len(args) == 0 {
		return false
	}

	switch args[0] {
	case "request", "unanswered", "equivocation":
		return len(args) == 3
	case "prove":
		return len(args) == 2
	default:
		return false
	}
}

func (r runner) registrations(ctx context.Context, action, path string, epoch uint64) int {
	file, code := open(r.stdout, action, path)
	if file == nil {
		return code
	}
	defer func() { _ = file.Close() }()

	withholding, issued, err := decodeRegistrations(file)
	if err != nil {
		return emit(r.stdout, report{Action: action, Verdict: "malformed", Reason: err.Error()}, ExitMalformed)
	}

	if withholding != r.identity.Withholding {
		return r.refuse(
			action,
			epoch,
			notary.Errf(
				notary.KindNotSponsored,
				"sponsorcli.registrations",
				nil,
				"these registrations were signed for the withholding court %s, this paymaster sponsors %s",
				withholding.Hex(),
				r.identity.Withholding.Hex(),
			),
		)
	}

	if err := r.generation(ctx, epoch); err != nil {
		return r.refuse(action, epoch, err)
	}

	worst := ExitLanded

	for _, registration := range issued {
		if registration.Registration.Holder != r.deps.Sponsor.Holder() {
			continue
		}

		call := r.deps.Sponsor.Calls().RequestStatement(r.identity.Withholding, registration, epoch)
		if action == "unanswered" {
			call = r.deps.Sponsor.Calls().RecordUnanswered(r.identity.Withholding, registration, epoch)
		}

		out := report{
			Action:  action,
			Epoch:   strconv.FormatUint(epoch, 10),
			Account: registration.Registration.Account.Hex(),
		}
		worst = worse(worst, r.submit(ctx, out, call))
	}

	if worst == ExitLanded && !r.holds(issued) {
		return r.refuse(
			action,
			epoch,
			notary.Errf(
				notary.KindNotSponsored,
				"sponsorcli.registrations",
				nil,
				"none of these registrations belongs to %s, the key this command signs with",
				r.deps.Sponsor.Holder().Hex(),
			),
		)
	}

	return worst
}

func (r runner) prove(ctx context.Context, path string) int {
	claims, code := r.claims("prove", path)
	if claims == nil {
		return code
	}

	worst := ExitLanded
	proved := false

	for _, claim := range claims {
		out := report{Action: "prove", Epoch: strconv.FormatUint(claim.Epoch, 10)}

		verification, err := r.deps.Verifier.Verify(ctx, claim)
		switch {
		case err != nil:
			worst = worse(worst, unlessNothingToDo(failed(r.stdout, out, err)))
			continue
		case claim.Statement == nil || !verification.Provable:
			out.Verdict = "nothing_to_prove"
			out.Reason = "the proof is " + verification.Verdict.String() + " and its statement " +
				verification.Check.String() + "; no court would convict it"
			emit(r.stdout, out, ExitNothingToDo)

			continue
		case verification.Court.Address != r.identity.Court:
			worst = worse(
				worst,
				r.refuse(
					"prove",
					claim.Epoch,
					notary.Errf(
						notary.KindNotSponsored,
						"sponsorcli.prove",
						nil,
						"epoch %d is judged by the court %s, this paymaster sponsors %s",
						claim.Epoch,
						verification.Court.Address.Hex(),
						r.identity.Court.Hex(),
					),
				),
			)

			continue
		}

		proved = true
		call := r.deps.Sponsor.Calls().ProveBadProof(
			r.identity.Court,
			notary.StatementOf(r.deps.Statements, claim),
			claim.Statement.Signature,
			claim.Proof,
		)
		worst = worse(worst, r.submit(ctx, out, call))
	}

	if !proved && worst == ExitLanded {
		return ExitNothingToDo
	}

	return worst
}

func unlessNothingToDo(code int) int {
	if code == ExitNothingToDo {
		return ExitLanded
	}

	return code
}

func (r runner) equivocation(ctx context.Context, firstPath, secondPath string) int {
	firsts, code := r.claims("equivocation", firstPath)
	if firsts == nil {
		return code
	}

	seconds, code := r.claims("equivocation", secondPath)
	if seconds == nil {
		return code
	}

	for _, first := range firsts {
		for _, second := range seconds {
			if !r.conflict(first, second) {
				continue
			}

			if err := r.attributable(ctx, first, second); err != nil {
				return failed(r.stdout, report{Action: "equivocation", Epoch: strconv.FormatUint(first.Epoch, 10)}, err)
			}

			call := r.deps.Sponsor.Calls().ProveEquivocation(
				r.identity.Court,
				notary.StatementOf(r.deps.Statements, first),
				first.Statement.Signature,
				notary.StatementOf(r.deps.Statements, second),
				second.Statement.Signature,
			)

			return r.submit(ctx, report{Action: "equivocation", Epoch: strconv.FormatUint(first.Epoch, 10)}, call)
		}
	}

	return emit(
		r.stdout,
		report{
			Action:  "equivocation",
			Verdict: "nothing_to_prove",
			Reason:  "no two signed statements in these files contradict each other",
		},
		ExitNothingToDo,
	)
}

func (r runner) conflict(first, second notary.Claim) bool {
	if first.Statement == nil || second.Statement == nil || first.Ledger != second.Ledger ||
		first.Epoch != second.Epoch {
		return false
	}

	a := notary.StatementOf(r.deps.Statements, first)
	b := notary.StatementOf(r.deps.Statements, second)

	return a != b && (a.Commitment == b.Commitment || a.Index == b.Index)
}

func (r runner) attributable(ctx context.Context, claims ...notary.Claim) error {
	for _, claim := range claims {
		verification, err := r.deps.Verifier.Verify(ctx, claim)
		if err != nil {
			return err
		}

		if verification.Check != notary.StatementAttributable || verification.Court.Address != r.identity.Court {
			return notary.Errf(
				notary.KindNotSponsored,
				"sponsorcli.attributable",
				nil,
				"the statement for epoch %d is not attributable to the signer of the court this paymaster sponsors",
				claim.Epoch,
			)
		}
	}

	return nil
}

func (r runner) generation(ctx context.Context, epoch uint64) error {
	generation, err := r.deps.Directory.GenerationOf(ctx, epoch)
	if err != nil {
		return err
	}

	if generation.ChainID != r.identity.ChainID ||
		generation.Court != r.identity.Court ||
		generation.Withholding != r.identity.Withholding {
		return notary.Errf(
			notary.KindNotSponsored,
			"sponsorcli.generation",
			nil,
			"epoch %d is held by the generation with the courts %s and %s on chain %d; this paymaster sponsors %s and %s",
			epoch,
			generation.Court.Hex(),
			generation.Withholding.Hex(),
			generation.ChainID,
			r.identity.Court.Hex(),
			r.identity.Withholding.Hex(),
		)
	}

	return nil
}

func (r runner) holds(issued []notary.IssuedRegistration) bool {
	for _, registration := range issued {
		if registration.Registration.Holder == r.deps.Sponsor.Holder() {
			return true
		}
	}

	return false
}

func (r runner) submit(ctx context.Context, out report, call notary.SponsoredCall) int {
	receipt, err := r.deps.Sponsor.Submit(ctx, r.identity, call)
	if err != nil {
		return failed(r.stdout, out, err)
	}

	out.Operation = receipt.Hash.Hex()
	out.Transaction = receipt.TransactionHash.Hex()
	out.GasCost = receipt.ActualGasCost.String()

	if !receipt.Success {
		out.Verdict = "reverted"
		out.Reason = receipt.Reason

		return emit(r.stdout, out, ExitReverted)
	}

	out.Verdict = "landed"

	return emit(r.stdout, out, ExitLanded)
}

func (r runner) refuse(action string, epoch uint64, err error) int {
	return failed(r.stdout, report{Action: action, Epoch: strconv.FormatUint(epoch, 10)}, err)
}

func (r runner) claims(action, path string) ([]notary.Claim, int) {
	file, code := open(r.stdout, action, path)
	if file == nil {
		return nil, code
	}
	defer func() { _ = file.Close() }()

	claims, err := prooffile.DecodeAll(file, r.deps.Scheme)
	if err != nil {
		return nil, emit(r.stdout, report{Action: action, Verdict: "malformed", Reason: err.Error()}, ExitMalformed)
	}

	return claims, ExitLanded
}

func open(stdout io.Writer, action, path string) (*os.File, int) {
	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, emit(
			stdout,
			report{Action: action, Verdict: "no_input", Reason: "the file could not be opened"},
			ExitNoInput,
		)
	}

	return file, ExitLanded
}

func Fail(stdout io.Writer, action string, err error) int {
	return failed(stdout, report{Action: action}, err)
}

func failed(stdout io.Writer, out report, err error) int {
	out.Verdict = notary.KindOf(err).String()
	out.Reason = err.Error()

	switch notary.KindOf(err) {
	case notary.KindNotSponsored, notary.KindWrongDirectory, notary.KindMoved, notary.KindWrongLedger,
		notary.KindWrongCourt:
		return emit(stdout, out, ExitNotSponsored)
	case notary.KindEpochNotAnchored, notary.KindEpochNotFinal, notary.KindNotDue, notary.KindCallRefused:
		return emit(stdout, out, ExitNothingToDo)
	case notary.KindInvalidArgument:
		return emit(stdout, out, ExitMalformed)
	default:
		return emit(stdout, out, ExitUnknown)
	}
}

func worse(current, next int) int {
	if current == ExitLanded {
		return next
	}

	return current
}

func emit(stdout io.Writer, out report, code int) int {
	encoded, err := json.Marshal(out)
	if err != nil {
		return ExitUnknown
	}

	if _, err := stdout.Write(append(encoded, '\n')); err != nil {
		return ExitUnknown
	}

	return code
}
