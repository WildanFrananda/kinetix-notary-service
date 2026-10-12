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
