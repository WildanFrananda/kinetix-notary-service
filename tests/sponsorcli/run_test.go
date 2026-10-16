package sponsorcli_test

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/ethereum"
	"github.com/WildanFrananda/kinetix-notary-service/internal/keccak"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/prooffile"
	"github.com/WildanFrananda/kinetix-notary-service/internal/sponsorcli"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

var (
	holder   = notary.Address{0x40}
	stranger = notary.Address{0x41}
	identity = notary.PaymasterIdentity{
		ChainID:     notarytest.LedgerIdentity.ChainID,
		Address:     notary.Address{0xfa},
		EntryPoint:  notary.Address{0x43},
		Account:     notary.Address{0x77},
		Court:       notarytest.CourtIdentity.Address,
		Withholding: notary.Address{0xd0},
		Signer:      notary.Address{0x51},
		Ceiling:     big.NewInt(1),
	}
)

type sponsorship struct {
	submitted []notary.SponsoredCall
	receipt   notary.OperationReceipt
	err       error
}

func (s *sponsorship) Holder() notary.Address {
	return holder
}

func (s *sponsorship) Calls() notary.SponsoredCalls {
	return ethereum.NewSponsoredCalls()
}

func (s *sponsorship) Identity(context.Context) (notary.PaymasterIdentity, error) {
	return identity, nil
}

func (s *sponsorship) Submit(
	_ context.Context,
	_ notary.PaymasterIdentity,
	call notary.SponsoredCall,
) (notary.OperationReceipt, error) {
	s.submitted = append(s.submitted, call)

	return s.receipt, s.err
}

type verifier struct {
	answers map[uint64]notary.Verification
	err     error
}

func (v verifier) Verify(_ context.Context, claim notary.Claim) (notary.Verification, error) {
	return v.answers[claim.Epoch], v.err
}

func landed() notary.OperationReceipt {
	return notary.OperationReceipt{
		Hash:            notary.Digest{0xab},
		Success:         true,
		TransactionHash: notary.Digest{0x7a},
		ActualGasCost:   big.NewInt(1_234),
	}
}

func directory(court, withholding notary.Address) *notarytest.Directory {
	return &notarytest.Directory{
		IdentityValue: notarytest.DirectoryIdentity,
		Generations: []notary.Generation{{
			ChainID:     identity.ChainID,
			Ledger:      notarytest.LedgerIdentity.Address,
			Court:       court,
			Withholding: withholding,
		}},
	}
}

func deps(s *sponsorship, d notary.Directory, v sponsorcli.ClaimVerifier) sponsorcli.Dependencies {
	return sponsorcli.Dependencies{
		Sponsor:    s,
		Directory:  d,
		Verifier:   v,
		Statements: keccak.NewScheme(),
		Scheme:     keccak.NewScheme().Version(),
	}
}

func registrations(t *testing.T, withholding notary.Address, holders ...notary.Address) string {
	t.Helper()

	type entry struct {
		Holder    string `json:"holder"`
		Account   string `json:"account"`
		FromEpoch string `json:"fromEpoch"`
		Kind      string `json:"kind"`
		Digest    string `json:"digest"`
		Signature string `json:"signature"`
	}

	entries := make([]entry, len(holders))
	for i, h := range holders {
		entries[i] = entry{
			Holder:    h.Hex(),
			Account:   notary.Digest{byte(i + 1)}.Hex(),
			FromEpoch: "4",
			Kind:      notary.AccountCustomerBalance.String(),
			Digest:    notary.Digest{0xdd}.Hex(),
			Signature: notary.Signature{0x01, 64: 27}.Hex(),
		}
	}

	document := map[string]any{
		"court":         map[string]string{"chainId": "84532", "address": withholding.Hex()},
		"epoch":         "6",
		"registrations": entries,
	}

	return written(t, "registrations.json", document)
}

func written(t *testing.T, name string, document any) string {
	t.Helper()

	encoded, err := json.Marshal(document)
	require.NoError(t, err)

	path := filepath.Join(t.TempDir(), name)
	require.NoError(t, os.WriteFile(path, encoded, 0o600))

	return path
}

func proofFile(t *testing.T, claims ...notary.Claim) string {
	t.Helper()

	documents := make([]json.RawMessage, len(claims))
	for i, claim := range claims {
		document, err := prooffile.Marshal(claim, keccak.NewScheme().Version())
		require.NoError(t, err)

		documents[i] = document
	}

	return written(t, "proof.json", map[string]any{"proofs": documents})
}

func run(t *testing.T, d sponsorcli.Dependencies, args ...string) (int, []map[string]string) {
	t.Helper()

	var out bytes.Buffer

	code := sponsorcli.Run(context.Background(), args, d, &out)

	var reports []map[string]string

	scanner := bufio.NewScanner(strings.NewReader(out.String()))
	for scanner.Scan() {
		var report map[string]string
		require.NoError(t, json.Unmarshal(scanner.Bytes(), &report))

		reports = append(reports, report)
	}

	return code, reports
}

func signedClaim(t *testing.T, epoch uint64) notary.Claim {
	t.Helper()

	_, claim := notarytest.AnchoredClaim(t, epoch, 5, 2)
	claim.Statement = &notary.StatementSignature{Court: identity.Court, Signature: notary.Signature{0x02, 64: 27}}

	return claim
}

func TestEachOfTheHoldersRegistrationsRequestsAStatementAndOthersAreLeftAlone(t *testing.T) {
	s := &sponsorship{receipt: landed()}
	path := registrations(t, identity.Withholding, holder, stranger, holder)

	code, reports := run(t, deps(s, directory(identity.Court, identity.Withholding), verifier{}), "request", path, "6")

	require.Equal(t, sponsorcli.ExitLanded, code)
	require.Len(t, s.submitted, 2)
	require.Len(t, reports, 2)
	require.Equal(t, "landed", reports[0]["verdict"])
	require.Equal(t, notary.Digest{0x7a}.Hex(), reports[0]["transaction"])
	require.Equal(t, "1234", reports[0]["gasCost"])
	require.Equal(t, identity.Withholding, s.submitted[0].Target)
	require.Empty(t, s.submitted[0].PaymasterData)
}

func TestAnUnansweredRecordCarriesTheRegistrationSignatureForThePaymaster(t *testing.T) {
	s := &sponsorship{receipt: landed()}

	code, _ := run(
		t,
		deps(s, directory(identity.Court, identity.Withholding), verifier{}),
		"unanswered",
		registrations(t, identity.Withholding, holder),
		"6",
	)

	require.Equal(t, sponsorcli.ExitLanded, code)
	require.Equal(t, notary.Signature{0x01, 64: 27}.Hex(), notary.Signature(s.submitted[0].PaymasterData).Hex())
}

func TestRegistrationsForAnotherCourtOrGenerationOrHolderAreNeverSent(t *testing.T) {
	cases := map[string]struct {
		path      func(t *testing.T) string
		directory *notarytest.Directory
	}{
		"another withholding court": {
			path:      func(t *testing.T) string { return registrations(t, notary.Address{0xd9}, holder) },
			directory: directory(identity.Court, identity.Withholding),
		},
		"an epoch of another generation": {
			path:      func(t *testing.T) string { return registrations(t, identity.Withholding, holder) },
			directory: directory(notary.Address{0xc9}, notary.Address{0xd9}),
		},
		"none of the holder's": {
			path:      func(t *testing.T) string { return registrations(t, identity.Withholding, stranger) },
			directory: directory(identity.Court, identity.Withholding),
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			s := &sponsorship{receipt: landed()}

			code, reports := run(t, deps(s, c.directory, verifier{}), "request", c.path(t), "6")

			require.Equal(t, sponsorcli.ExitNotSponsored, code)
			require.Empty(t, s.submitted)
			require.Equal(t, "not_sponsored", reports[len(reports)-1]["verdict"])
		})
	}
}

func TestAnOperationThatLandedButRevertedIsReportedAsSuch(t *testing.T) {
	receipt := landed()
	receipt.Success = false
	receipt.Reason = "0x0cd7f4b7"
	s := &sponsorship{receipt: receipt}

	code, reports := run(
		t,
		deps(s, directory(identity.Court, identity.Withholding), verifier{}),
		"request",
		registrations(t, identity.Withholding, holder),
		"6",
	)

	require.Equal(t, sponsorcli.ExitReverted, code)
	require.Equal(t, "reverted", reports[0]["verdict"])
	require.Equal(t, "0x0cd7f4b7", reports[0]["reason"])
}

func TestACallNotDueYetOrOneTheCourtWouldRefuseIsNothingToDoNow(t *testing.T) {
	for kind, verdict := range map[notary.Kind]string{
		notary.KindNotDue:      "not_due",
		notary.KindCallRefused: "call_refused",
	} {
		t.Run(verdict, func(t *testing.T) {
			s := &sponsorship{err: notary.Errf(kind, "test", nil, "the court would refuse it")}

			code, reports := run(
				t,
				deps(s, directory(identity.Court, identity.Withholding), verifier{}),
				"unanswered",
				registrations(t, identity.Withholding, holder),
				"6",
			)

			require.Equal(t, sponsorcli.ExitNothingToDo, code)
			require.Equal(t, verdict, reports[0]["verdict"])
			require.Equal(t, "6", reports[0]["epoch"])
		})
	}
}

func TestOnlyAProvableStatementIsTakenToTheCourt(t *testing.T) {
	provable := signedClaim(t, 4)
	holding := signedClaim(t, 5)
	v := verifier{answers: map[uint64]notary.Verification{
		4: {
			Verdict:  notary.VerdictRejected,
			Check:    notary.StatementAttributable,
			Provable: true,
			Court:    notary.CourtIdentity{ChainID: identity.ChainID, Address: identity.Court},
		},
		5: {
			Verdict: notary.VerdictIncluded,
			Check:   notary.StatementAttributable,
			Court:   notary.CourtIdentity{ChainID: identity.ChainID, Address: identity.Court},
		},
	}}
	s := &sponsorship{receipt: landed()}

	code, reports := run(
		t,
		deps(s, directory(identity.Court, identity.Withholding), v),
		"prove",
		proofFile(t, provable, holding),
	)

	require.Equal(t, sponsorcli.ExitLanded, code)
	require.Len(t, s.submitted, 1)
	require.Equal(t, identity.Court, s.submitted[0].Target)
	require.Equal(t, "landed", reports[0]["verdict"])
	require.Equal(t, "nothing_to_prove", reports[1]["verdict"])
}

func TestAProofThatHoldsOrCannotBeCheckedIsNeverTakenToTheCourt(t *testing.T) {
	claim := signedClaim(t, 4)
	for name, v := range map[string]verifier{
		"holds": {answers: map[uint64]notary.Verification{4: {
			Verdict: notary.VerdictIncluded,
			Check:   notary.StatementAttributable,
		}}},
		"not final": {err: notary.Errf(notary.KindEpochNotFinal, "test", nil, "not final")},
	} {
		t.Run(name, func(t *testing.T) {
			s := &sponsorship{receipt: landed()}

			code, _ := run(t, deps(s, directory(identity.Court, identity.Withholding), v), "prove", proofFile(t, claim))

			require.Equal(t, sponsorcli.ExitNothingToDo, code)
			require.Empty(t, s.submitted)
		})
	}
}

func TestTwoSignedStatementsForOneLeafAreAnEquivocation(t *testing.T) {
	first := signedClaim(t, 4)
	second := first
	second.Entry.Amount++
	attributable := notary.Verification{
		Check: notary.StatementAttributable,
		Court: notary.CourtIdentity{ChainID: identity.ChainID, Address: identity.Court},
	}
	v := verifier{answers: map[uint64]notary.Verification{4: attributable}}
	s := &sponsorship{receipt: landed()}

	code, _ := run(
		t,
		deps(s, directory(identity.Court, identity.Withholding), v),
		"equivocation",
		proofFile(t, first),
		proofFile(t, second),
	)

	require.Equal(t, sponsorcli.ExitLanded, code)
	require.Len(t, s.submitted, 1)

	s = &sponsorship{receipt: landed()}
	code, _ = run(
		t,
		deps(s, directory(identity.Court, identity.Withholding), v),
		"equivocation",
		proofFile(t, first),
		proofFile(t, first),
	)

	require.Equal(t, sponsorcli.ExitNothingToDo, code)
	require.Empty(t, s.submitted)
}

func TestUsageMissingFilesAndMalformedInputAreTheirOwnAnswers(t *testing.T) {
	d := deps(&sponsorship{}, directory(identity.Court, identity.Withholding), verifier{})

	code, _ := run(t, d, "transfer", "x")
	require.Equal(t, sponsorcli.ExitUsage, code)

	code, _ = run(t, d, "request", filepath.Join(t.TempDir(), "absent.json"), "6")
	require.Equal(t, sponsorcli.ExitNoInput, code)

	code, _ = run(t, d, "request", written(t, "r.json", map[string]any{"court": 1}), "6")
	require.Equal(t, sponsorcli.ExitMalformed, code)

	code, _ = run(t, d, "request", registrations(t, identity.Withholding, holder), "six")
	require.Equal(t, sponsorcli.ExitMalformed, code)
}
