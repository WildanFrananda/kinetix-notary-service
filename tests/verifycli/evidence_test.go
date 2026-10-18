package verifycli_test

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/keccak"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/prooffile"
	"github.com/WildanFrananda/kinetix-notary-service/internal/verifycli"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

type stubEvidenceVerifier struct {
	answers func(notary.EvidenceClaim) (notary.EvidenceVerification, error)
}

func (s stubEvidenceVerifier) Verify(
	_ context.Context,
	claim notary.EvidenceClaim,
) (notary.EvidenceVerification, error) {
	return s.answers(claim)
}

func evidenceDocument(t *testing.T) (string, notarytest.EvidenceBatch) {
	t.Helper()

	batch := notarytest.BatchOf(
		t,
		keccak.NewEvidence(),
		notary.RegistryIdentity{ChainID: 84532, Address: notary.Address{0xe0, 0x11}},
		7,
		notary.Digest{0x07},
		1_760_000_000,
		notarytest.Records(notarytest.Reviews(t, 2)),
	)

	documents := make([]json.RawMessage, len(batch.Claims))
	for i, claim := range batch.Claims {
		raw, err := prooffile.MarshalEvidence(claim)
		require.NoError(t, err)

		documents[i] = raw
	}

	bundle, err := json.Marshal(map[string]any{"proofs": documents, "pending": 0})
	require.NoError(t, err)

	return string(bundle), batch
}

func included(batch notarytest.EvidenceBatch) func(notary.EvidenceClaim) (notary.EvidenceVerification, error) {
	return func(claim notary.EvidenceClaim) (notary.EvidenceVerification, error) {
		closedAt := batch.Header.ClosedAt

		return notary.EvidenceVerification{
			Verdict:  notary.VerdictIncluded,
			Registry: claim.Registry,
			Batch: notary.AnchoredBatch{
				Stream:     batch.Header.Stream,
				Number:     batch.Header.Number,
				Digest:     batch.Digest,
				ClosedAt:   closedAt,
				AnchoredAt: closedAt,
			},
			Record: claim.Record,
		}, nil
	}
}

func lines(t *testing.T, out string) []map[string]string {
	t.Helper()

	var reports []map[string]string

	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
		var report map[string]string
		require.NoError(t, json.Unmarshal([]byte(line), &report))

		reports = append(reports, report)
	}

	return reports
}

func TestEveryEvidenceProofInABundleGetsAVerdict(t *testing.T) {
	document, batch := evidenceDocument(t)

	var out bytes.Buffer

	code := verifycli.RunEvidence(
		context.Background(),
		strings.NewReader(document),
		1,
		stubEvidenceVerifier{answers: included(batch)},
		&out,
	)

	require.Equal(t, verifycli.ExitIncluded, code)

	reports := lines(t, out.String())
	require.Len(t, reports, 2)
	require.Equal(t, "included", reports[0]["verdict"])
	require.Equal(t, "review", reports[0]["stream"])
	require.Equal(t, "7", reports[0]["batch"])
	require.Equal(t, "700", reports[0]["transaction"])
	require.Equal(t, "2025-10-09T08:53:20Z", reports[0]["closedAt"])
}

func TestTheWorstEvidenceVerdictIsTheExitCode(t *testing.T) {
	document, batch := evidenceDocument(t)
	seen := 0
	answers := func(claim notary.EvidenceClaim) (notary.EvidenceVerification, error) {
		seen++
		if seen == 2 {
			return notary.EvidenceVerification{}, notary.Errf(notary.KindEpochNotFinal, "test", nil, "not yet")
		}

		return included(batch)(claim)
	}

	var out bytes.Buffer

	code := verifycli.RunEvidence(
		context.Background(),
		strings.NewReader(document),
		1,
		stubEvidenceVerifier{answers},
		&out,
	)

	require.Equal(t, verifycli.ExitNotFinal, code)
	require.Equal(t, "not_final", lines(t, out.String())[1]["verdict"])
}

func TestARejectedEvidenceProofExitsRejected(t *testing.T) {
	document, batch := evidenceDocument(t)
	answers := func(claim notary.EvidenceClaim) (notary.EvidenceVerification, error) {
		verification, _ := included(batch)(claim)
		verification.Verdict, verification.Reason = notary.VerdictRejected, "the proof reaches elsewhere"

		return verification, nil
	}

	var out bytes.Buffer

	code := verifycli.RunEvidence(
		context.Background(),
		strings.NewReader(document),
		1,
		stubEvidenceVerifier{answers},
		&out,
	)

	require.Equal(t, verifycli.ExitRejected, code)
	require.Equal(t, "the proof reaches elsewhere", lines(t, out.String())[0]["reason"])
}

func TestAMalformedEvidenceDocumentIsNamedSo(t *testing.T) {
	var out bytes.Buffer

	code := verifycli.RunEvidence(
		context.Background(),
		strings.NewReader(`{"format":"kinetix.notary.evidence.v1"}`),
		1,
		stubEvidenceVerifier{},
		&out,
	)

	require.Equal(t, verifycli.ExitMalformedProof, code)
	require.Equal(t, "malformed_proof", lines(t, out.String())[0]["verdict"])
}

func TestAProofFileIsReadBoundedAndAMissingOneIsNoProof(t *testing.T) {
	document, _ := evidenceDocument(t)
	path := filepath.Join(t.TempDir(), "evidence.json")
	require.NoError(t, os.WriteFile(path, []byte(document), 0o600))

	raw, err := verifycli.ReadProof(path)
	require.NoError(t, err)
	require.True(t, prooffile.IsEvidence(raw))

	_, err = verifycli.ReadProof(filepath.Join(t.TempDir(), "missing.json"))
	require.Error(t, err)

	var out bytes.Buffer
	require.Equal(t, verifycli.ExitNoProof, verifycli.NoProof(&out))
	require.Equal(t, "no_proof", lines(t, out.String())[0]["verdict"])
}
