package prooffile_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/prooffile"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

func encoded(t *testing.T) (notary.Claim, string) {
	t.Helper()
	_, claim := notarytest.AnchoredClaim(t, 7, 9, 4)

	var buf bytes.Buffer
	require.NoError(t, prooffile.Encode(&buf, claim, 1))

	return claim, buf.String()
}

func TestAClaimRoundTripsThroughTheDocument(t *testing.T) {
	claim, doc := encoded(t)

	decoded, err := prooffile.Decode(strings.NewReader(doc), 1)

	require.NoError(t, err)
	require.Equal(t, claim, decoded)
	require.Contains(t, doc, `"format": "kinetix.notary.inclusion.v1"`)
	require.Contains(t, doc, `"chainId": "84532"`)
}

func TestNumbersTravelAsDecimalStrings(t *testing.T) {
	_, doc := encoded(t)

	require.Regexp(t, `"amount": "\d+"`, doc)
	require.Regexp(t, `"sum": "\d+"`, doc)
}

func requireMalformed(t *testing.T, doc string) {
	t.Helper()

	_, err := prooffile.Decode(strings.NewReader(doc), 1)
	require.Equal(t, notary.KindInvalidArgument, notary.KindOf(err), "%v", err)
}

func TestMalformedDocumentsAreRefused(t *testing.T) {
	_, doc := encoded(t)

	cases := map[string]string{
		"unknown field":      strings.Replace(doc, `"scheme": 1,`, `"scheme": 1, "extra": true,`, 1),
		"another format":     strings.Replace(doc, "inclusion.v1", "inclusion.v2", 1),
		"another scheme":     strings.Replace(doc, `"scheme": 1`, `"scheme": 2`, 1),
		"negative amount":    strings.Replace(doc, `"amount": "`, `"amount": "-`, 1),
		"numeric amount":     strings.Replace(doc, `"amount": "`, `"amount": 1, "ignored": "`, 1),
		"short ledger":       strings.Replace(doc, `"ledger": "0x`, `"ledger": "0x00`, 1),
		"second document":    doc + doc,
		"trailing garbage":   doc + "garbage",
		"not json":           "verdict: included",
		"empty":              "",
		"chain id overflows": strings.Replace(doc, `"chainId": "84532"`, `"chainId": "18446744073709551616"`, 1),
	}
	for name, body := range cases {
		t.Run(name, func(t *testing.T) {
			requireMalformed(t, body)
		})
	}
}

func TestAProofDeeperThanTheSchemeIsRefused(t *testing.T) {
	_, claim := notarytest.AnchoredClaim(t, 7, 9, 4)
	claim.Proof.Siblings = make([]notary.Node, notary.MaxDepth+1)

	var buf bytes.Buffer
	require.NoError(t, prooffile.Encode(&buf, claim, 1))

	requireMalformed(t, buf.String())
}

func TestAnOversizedDocumentIsRefused(t *testing.T) {
	_, doc := encoded(t)

	requireMalformed(t, strings.Replace(doc, `"format"`, strings.Repeat(" ", 70<<10)+`"format"`, 1))
}

func TestAnAccountOpeningSurvivesTheDocument(t *testing.T) {
	claim := sampleClaim(t)
	claim.Opening = &notary.Opening{
		Holder: notarytest.Holder(t, "9f1d4a3e-1c62-4d0a-9a7b-2f5c8e0b41d7"),
		Kind:   notary.AccountDriverPendingEscrow,
		Salt:   notary.Salt{0x77, 0x01},
	}

	raw, err := prooffile.Marshal(claim, 1)
	require.NoError(t, err)
	require.Contains(t, string(raw), `"kind":"driver_pending_escrow"`)

	decoded, err := prooffile.DecodeAll(bytes.NewReader(raw), 1)
	require.NoError(t, err)
	require.Equal(t, []notary.Claim{claim}, decoded)
}

func TestABundleFromTheAPIDecodesEveryProof(t *testing.T) {
	first := sampleClaim(t)
	second := sampleClaim(t)
	second.Proof.Index = 0

	a, err := prooffile.Marshal(first, 1)
	require.NoError(t, err)
	b, err := prooffile.Marshal(second, 1)
	require.NoError(t, err)

	body := `{"epoch":"3","proofs":[` + string(a) + `,` + string(b) + `]}`

	claims, err := prooffile.DecodeAll(strings.NewReader(body), 1)
	require.NoError(t, err)
	require.Equal(t, []notary.Claim{first, second}, claims)
}

func TestABrokenAccountOrBundleIsRefused(t *testing.T) {
	good, err := prooffile.Marshal(sampleClaim(t), 1)
	require.NoError(t, err)

	withAccount := func(account string) string {
		return strings.Replace(string(good), `"siblings"`, `"account":`+account+`,"siblings"`, 1)
	}

	cases := map[string]string{
		"empty bundle":      `{"proofs":[]}`,
		"proofs not a list": `{"proofs":{}}`,
		"not json":          `proofs`,
		"bad holder": withAccount(
			`{"holder":"two words","kind":"customer_balance","salt":"0x` + strings.Repeat("11", 32) + `"}`,
		),
		"bad kind": withAccount(
			`{"holder":"p-1","kind":"savings","salt":"0x` + strings.Repeat("11", 32) + `"}`,
		),
		"short salt": withAccount(`{"holder":"p-1","kind":"customer_balance","salt":"0x11"}`),
		"unknown field": withAccount(
			`{"holder":"p-1","kind":"customer_balance","salt":"0x` + strings.Repeat("11", 32) + `","pin":"1"}`,
		),
		"bad proof inside": `{"proofs":[{"format":"other"}]}`,
	}
	for name, body := range cases {
		_, err := prooffile.DecodeAll(strings.NewReader(body), 1)

		require.Error(t, err, name)
		require.Equal(t, notary.KindInvalidArgument, notary.KindOf(err), name)
	}
}

func sampleClaim(t *testing.T) notary.Claim {
	t.Helper()

	_, claim := notarytest.AnchoredClaim(t, 3, 5, 2)

	return claim
}

func TestAStatementSurvivesTheDocument(t *testing.T) {
	_, claim := notarytest.AnchoredClaim(t, 7, 9, 4)
	claim.Statement = &notary.StatementSignature{
		Court:     notary.Address{0xc0, 0x07},
		Signature: notary.Signature{0: 0xaa, 64: 27},
	}

	var buf bytes.Buffer
	require.NoError(t, prooffile.Encode(&buf, claim, 1))

	decoded, err := prooffile.Decode(strings.NewReader(buf.String()), 1)

	require.NoError(t, err)
	require.Equal(t, claim, decoded)
	require.Contains(t, buf.String(), `"court": "0xc007000000000000000000000000000000000000"`)
	require.Regexp(t, `"signature": "0xaa0{126}1b"`, buf.String())
}

func TestAProofWithoutAStatementStillDecodes(t *testing.T) {
	claim, doc := encoded(t)

	decoded, err := prooffile.Decode(strings.NewReader(doc), 1)

	require.NoError(t, err)
	require.Nil(t, decoded.Statement)
	require.Equal(t, claim, decoded)
	require.NotContains(t, doc, `"statement"`)
}

func TestABrokenStatementIsRefused(t *testing.T) {
	_, doc := encoded(t)

	cases := map[string]string{
		"a short signature": `"statement": {"court": "0xc007000000000000000000000000000000000000", "signature": "0xaa"}`,
		"a court that is not an address": `"statement": {"court": "0xc0", "signature": "0x` +
			strings.Repeat("aa", 65) + `"}`,
	}
	for name, block := range cases {
		t.Run(name, func(t *testing.T) {
			broken := strings.Replace(doc, `"format":`, block+`, "format":`, 1)

			_, err := prooffile.Decode(strings.NewReader(broken), 1)

			require.Equal(t, notary.KindInvalidArgument, notary.KindOf(err))
		})
	}
}

func TestASecondDerivationAccountTravelsAsFormatTwo(t *testing.T) {
	claim := sampleClaim(t)
	claim.Opening = &notary.Opening{
		Holder: notarytest.Holder(t, "9f1d4a3e-1c62-4d0a-9a7b-2f5c8e0b41d7"),
		Kind:   notary.AccountCustomerBalance,
		Salt:   notary.Salt{0x77, 0x01},
		Secret: notary.AccountSecret{0x55, 0x02},
	}

	raw, err := prooffile.Marshal(claim, 1)
	require.NoError(t, err)
	require.Contains(t, string(raw), `"format":"kinetix.notary.inclusion.v2"`)
	require.Contains(t, string(raw), `"secret":"0x5502`)

	decoded, err := prooffile.DecodeAll(bytes.NewReader(raw), 1)
	require.NoError(t, err)
	require.Equal(t, []notary.Claim{claim}, decoded)
}

func TestAFirstDerivationAccountStillTravelsAsFormatOne(t *testing.T) {
	claim := sampleClaim(t)
	claim.Opening = &notary.Opening{
		Holder: notarytest.Holder(t, "p-1"),
		Kind:   notary.AccountCustomerBalance,
		Salt:   notary.Salt{0x77},
	}

	raw, err := prooffile.Marshal(claim, 1)
	require.NoError(t, err)
	require.Contains(t, string(raw), `"format":"kinetix.notary.inclusion.v1"`)
	require.NotContains(t, string(raw), `"secret"`)
}

func TestAFormatAndItsSecretMustAgree(t *testing.T) {
	claim := sampleClaim(t)
	claim.Opening = &notary.Opening{
		Holder: notarytest.Holder(t, "p-1"),
		Kind:   notary.AccountCustomerBalance,
		Salt:   notary.Salt{0x77},
		Secret: notary.AccountSecret{0x55},
	}
	raw, err := prooffile.Marshal(claim, 1)
	require.NoError(t, err)

	v2 := string(raw)
	require.Contains(t, v2, `,"secret":"0x55`+strings.Repeat("0", 62)+`"`)
	zero := `"secret":"0x` + strings.Repeat("0", 64) + `"`
	withZeroSecret := strings.Replace(v2, `"secret":"0x55`+strings.Repeat("0", 62)+`"`, zero, 1)

	requireMalformed(t, strings.Replace(v2, "inclusion.v2", "inclusion.v1", 1))
	requireMalformed(t, strings.Replace(v2, `,"secret":"0x55`+strings.Repeat("0", 62)+`"`, "", 1))
	requireMalformed(t, withZeroSecret)
	requireMalformed(t, strings.Replace(v2, "inclusion.v2", "inclusion.v3", 1))

	plain, err := prooffile.Marshal(sampleClaim(t), 1)
	require.NoError(t, err)
	requireMalformed(t, strings.Replace(string(plain), "inclusion.v1", "inclusion.v2", 1))
}
