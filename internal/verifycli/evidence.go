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

const maxProofBytes = 4 << 20

func ReadProof(path string) ([]byte, error) {
	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, notary.Errf(notary.KindInvalidArgument, "verifycli.ReadProof", err, "opening the proof")
	}
	defer func() { _ = file.Close() }()

	raw, err := io.ReadAll(io.LimitReader(file, maxProofBytes+1))
	if err != nil {
		return nil, notary.Errf(notary.KindInvalidArgument, "verifycli.ReadProof", err, "reading the proof")
	}

	return raw, nil
}

func NoProof(stdout io.Writer) int {
	return emit(stdout, report{Verdict: "no_proof", Reason: "the proof file could not be opened"}, ExitNoProof)
}

func RunEvidence(
	ctx context.Context,
	document io.Reader,
	scheme uint16,
	verifier EvidenceClaimVerifier,
	stdout io.Writer,
) int {
	claims, err := prooffile.DecodeEvidenceAll(document, scheme)
	if err != nil {
		return emit(stdout, report{Verdict: "malformed_proof", Reason: err.Error()}, ExitMalformedProof)
	}

	worst := ExitIncluded

	for _, claim := range claims {
		code := verifyEvidence(ctx, claim, verifier, stdout)
		if worst == ExitIncluded {
			worst = code
		}
	}

	return worst
}

func verifyEvidence(
	ctx context.Context,
	claim notary.EvidenceClaim,
	verifier EvidenceClaimVerifier,
	stdout io.Writer,
) int {
	verification, err := verifier.Verify(ctx, claim)
	if err != nil {
		return Fail(stdout, err)
	}

	position := verification.Record.Position()
	out := evidenceReport{
		Verdict:     verification.Verdict.String(),
		ChainID:     strconv.FormatUint(verification.Registry.ChainID, 10),
		Registry:    verification.Registry.Address.Hex(),
		Stream:      verification.Batch.Stream.String(),
		Batch:       strconv.FormatUint(verification.Batch.Number, 10),
		ClosedAt:    verification.Batch.ClosedAt.Time().Format(time.RFC3339),
		AnchoredAt:  verification.Batch.AnchoredAt.Time().Format(time.RFC3339),
		Transaction: strconv.FormatUint(position.Transaction, 10),
		Sequence:    strconv.FormatUint(position.Sequence, 10),
		Reason:      verification.Reason,
	}

	code := ExitRejected
	if verification.Verdict == notary.VerdictIncluded {
		code = ExitIncluded
	}

	encoded, err := json.Marshal(out)
	if err != nil {
		return ExitUnknown
	}

	if _, err := stdout.Write(append(encoded, '\n')); err != nil {
		return ExitUnknown
	}

	return code
}
