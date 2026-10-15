package verifycli

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/prooffile"
)

func Run(ctx context.Context, path string, scheme uint16, verifier ClaimVerifier, stdout io.Writer) int {
	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		return emit(stdout, report{Verdict: "no_proof", Reason: "the proof file could not be opened"}, ExitNoProof)
	}
	defer func() { _ = file.Close() }()

	claims, err := prooffile.DecodeAll(file, scheme)
	if err != nil {
		return emit(stdout, report{Verdict: "malformed_proof", Reason: err.Error()}, ExitMalformedProof)
	}

	worst := ExitIncluded

	for _, claim := range claims {
		code := verify(ctx, claim, verifier, stdout)
		if worst == ExitIncluded {
			worst = code
		}
	}

	return worst
}

func verify(ctx context.Context, claim notary.Claim, verifier ClaimVerifier, stdout io.Writer) int {
	verification, err := verifier.Verify(ctx, claim)
	if err != nil {
		return Fail(stdout, err)
	}

	out := report{
		Verdict:    verification.Verdict.String(),
		ChainID:    strconv.FormatUint(verification.Ledger.ChainID, 10),
		Ledger:     verification.Ledger.Address.Hex(),
		Epoch:      strconv.FormatUint(verification.Epoch.Number, 10),
		Total:      strconv.FormatUint(verification.Epoch.Total, 10),
		AsOf:       verification.Epoch.AsOf.Time().Format(time.RFC3339),
		AnchoredAt: verification.Epoch.AnchoredAt.Time().Format(time.RFC3339),
		Reason:     verification.Reason,
	}
	if verification.Check != 0 {
		out.Court = verification.Court.Address.Hex()
		out.Statement = verification.Check.String()
		out.Faults = strconv.FormatUint(verification.Faults, 10)
	}

	if opening := verification.Opening; opening != nil {
		out.Holder = opening.Holder.String()
		out.Kind = opening.Kind.String()
		out.Amount = strconv.FormatUint(claim.Entry.Amount, 10)
	}

	if verification.Verdict == notary.VerdictIncluded {
		return emit(stdout, out, ExitIncluded)
	}

	if verification.Provable {
		out.Provable = "bad_proof"
	}

	return emit(stdout, out, ExitRejected)
}

func Fail(stdout io.Writer, err error) int {
	verdict, code := failure(notary.KindOf(err))
	return emit(stdout, report{Verdict: verdict, Reason: err.Error()}, code)
}

func failure(kind notary.Kind) (string, int) {
	switch kind {
	case notary.KindEpochNotAnchored:
		return "not_anchored", ExitNotAnchored
	case notary.KindEpochNotFinal:
		return "not_final", ExitNotFinal
	case notary.KindWrongLedger:
		return "wrong_ledger", ExitWrongLedger
	case notary.KindWrongCourt:
		return "wrong_court", ExitWrongCourt
	case notary.KindWrongDirectory:
		return "wrong_directory", ExitWrongDirectory
	case notary.KindMoved:
		return "moved", ExitMoved
	case notary.KindDisagreement:
		return "disagreement", ExitDisagreement
	default:
		return "unknown", ExitUnknown
	}
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
