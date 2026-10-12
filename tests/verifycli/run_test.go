package verifycli_test

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/prooffile"
	"github.com/WildanFrananda/kinetix-notary-service/internal/verifycli"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

type stubVerifier struct {
	verification notary.Verification
	err          error
}

func (s stubVerifier) Verify(context.Context, notary.Claim) (notary.Verification, error) {
	return s.verification, s.err
}

func proofFile(t *testing.T) (string, notary.AnchoredEpoch) {
	t.Helper()
	epoch, claim := notarytest.AnchoredClaim(t, 3, 5, 1)

	var buf bytes.Buffer
	require.NoError(t, prooffile.Encode(&buf, claim, 1))
	path := filepath.Join(t.TempDir(), "proof.json")
	require.NoError(t, os.WriteFile(path, buf.Bytes(), 0o600))

	return path, epoch
}

func runWith(t *testing.T, path string, verifier verifycli.ClaimVerifier) (int, map[string]string) {
	t.Helper()

	var out bytes.Buffer

	code := verifycli.Run(context.Background(), path, 1, verifier, &out)

	var report map[string]string
	require.NoError(t, json.Unmarshal(out.Bytes(), &report), out.String())

	return code, report
}

func TestAnIncludedClaimReportsTheEpochItWasCheckedAgainst(t *testing.T) {
	path, epoch := proofFile(t)
	verifier := stubVerifier{verification: notary.Verification{
		Verdict: notary.VerdictIncluded,
		Ledger:  notarytest.LedgerIdentity,
		Epoch:   epoch,
	}}

	code, report := runWith(t, path, verifier)

	require.Equal(t, verifycli.ExitIncluded, code)
	require.Equal(t, "included", report["verdict"])
	require.Equal(t, "84532", report["chainId"])
	require.Equal(t, "3", report["epoch"])
	require.Equal(t, "2025-10-09T08:53:20Z", report["asOf"])
	require.Equal(t, "2025-10-09T08:54:20Z", report["anchoredAt"])
}

func TestARejectedClaimCarriesItsReason(t *testing.T) {
	path, epoch := proofFile(t)
	verifier := stubVerifier{verification: notary.Verification{
		Verdict: notary.VerdictRejected,
		Ledger:  notarytest.LedgerIdentity,
		Epoch:   epoch,
		Reason:  "the proof sums to 9, epoch 3 holds 10",
	}}

	code, report := runWith(t, path, verifier)

	require.Equal(t, verifycli.ExitRejected, code)
	require.Equal(t, "rejected", report["verdict"])
	require.Equal(t, "the proof sums to 9, epoch 3 holds 10", report["reason"])
}

func TestEveryFailureHasItsOwnVerdictAndExitCode(t *testing.T) {
	path, _ := proofFile(t)

	cases := []struct {
		kind    notary.Kind
		verdict string
		code    int
	}{
		{notary.KindEpochNotAnchored, "not_anchored", verifycli.ExitNotAnchored},
		{notary.KindEpochNotFinal, "not_final", verifycli.ExitNotFinal},
		{notary.KindLedgerUnavailable, "unknown", verifycli.ExitUnknown},
		{notary.KindWrongLedger, "wrong_ledger", verifycli.ExitWrongLedger},
		{notary.KindDisagreement, "disagreement", verifycli.ExitDisagreement},
		{notary.KindUnknown, "unknown", verifycli.ExitUnknown},
	}
	for _, c := range cases {
		t.Run(c.verdict, func(t *testing.T) {
			code, report := runWith(t, path, stubVerifier{err: notary.Errf(c.kind, "stub", nil, "because")})

			require.Equal(t, c.code, code)
			require.Equal(t, c.verdict, report["verdict"])
			require.Contains(t, report["reason"], "because")
			require.Empty(t, report["epoch"])
		})
	}
}

func TestAMissingProofFileIsNotAVerdict(t *testing.T) {
	code, report := runWith(t, filepath.Join(t.TempDir(), "absent.json"), stubVerifier{})

	require.Equal(t, verifycli.ExitNoProof, code)
	require.Equal(t, "no_proof", report["verdict"])
}

func TestAMalformedProofFileIsNotAVerdict(t *testing.T) {
	path := filepath.Join(t.TempDir(), "proof.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"format":"something else"}`), 0o600))

	code, report := runWith(t, path, stubVerifier{})

	require.Equal(t, verifycli.ExitMalformedProof, code)
	require.Equal(t, "malformed_proof", report["verdict"])
}
