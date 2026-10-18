package prooffile_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/keccak"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/prooffile"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

var evidenceRegistry = notary.RegistryIdentity{ChainID: 84532, Address: notary.Address{0xe0, 0x0e}}

func evidenceClaims(t *testing.T) []notary.EvidenceClaim {
	t.Helper()

	scheme := keccak.NewEvidence()
	lifecycle := notarytest.Lifecycle(t, "ORD-1", "buyer-1", "merchant-1", 10, 3)
	lifecycle[1].ReturnNumber, lifecycle[1].ReturnStatus = "RET-1", 2
	reviews := notarytest.Reviews(t, 2)
	reviews[1].CommentSHA256 = nil
	deliveries := notarytest.Deliveries(t, 2)
	deliveries[0].SignatureSHA256 = bytes.Repeat([]byte{0x5e}, 32)

	claims := notarytest.BatchOf(
		t,
		scheme,
		evidenceRegistry,
		0,
		notary.Digest{},
		1_760_000_000,
		notarytest.Records(lifecycle),
	).Claims
	claims = append(
		claims,
		notarytest.BatchOf(
			t,
			scheme,
			evidenceRegistry,
			2,
			notary.Digest{0x02},
			1_760_000_001,
			notarytest.Records(deliveries),
		).Claims...,
	)

	return append(
		claims,
		notarytest.BatchOf(
			t,
			scheme,
			evidenceRegistry,
			9,
			notary.Digest{0x09},
			1_760_000_002,
			notarytest.Records(reviews),
		).Claims...,
	)
}

func TestEveryKindOfEvidenceSurvivesItsDocument(t *testing.T) {
	for _, claim := range evidenceClaims(t) {
		var buf bytes.Buffer
		require.NoError(t, prooffile.EncodeEvidence(&buf, claim))
		require.True(t, prooffile.IsEvidence(buf.Bytes()))

		decoded, err := prooffile.DecodeEvidence(&buf, 1)
		require.NoError(t, err)
		require.Equal(t, claim, decoded)
	}
}

func TestAnEvidenceDocumentKeepsEveryNanosecond(t *testing.T) {
	claim := evidenceClaims(t)[3]
	record, ok := claim.Record.(notary.DeliveryRecord)
	require.True(t, ok)
	record.DeliveredAt = notarytest.EvidenceTime(t, 1_760_000_000, 123_456_789)
	claim.Record = record

	raw, err := prooffile.MarshalEvidence(claim)
	require.NoError(t, err)
	require.Contains(t, string(raw), "2025-10-09T08:53:20.123456789Z")

	decoded, err := prooffile.DecodeEvidence(bytes.NewReader(raw), 1)
	require.NoError(t, err)

	kept, ok := decoded.Record.(notary.DeliveryRecord)
	require.True(t, ok)
	require.Equal(t, uint32(123_456_789), kept.DeliveredAt.Nanos())
}

func TestABundleOfEvidenceDecodesEveryProof(t *testing.T) {
	claims := evidenceClaims(t)
	documents := make([]json.RawMessage, len(claims))

	for i, claim := range claims {
		raw, err := prooffile.MarshalEvidence(claim)
		require.NoError(t, err)

		documents[i] = raw
	}

	bundle, err := json.Marshal(map[string]any{"proofs": documents, "pending": 3})
	require.NoError(t, err)
	require.True(t, prooffile.IsEvidence(bundle))

	decoded, err := prooffile.DecodeEvidenceAll(bytes.NewReader(bundle), 1)
	require.NoError(t, err)
	require.Equal(t, claims, decoded)
}

func TestALiabilitiesProofIsNotEvidence(t *testing.T) {
	_, claim := notarytest.AnchoredClaim(t, 3, 5, 1)

	var buf bytes.Buffer
	require.NoError(t, prooffile.Encode(&buf, claim, 1))
	require.False(t, prooffile.IsEvidence(buf.Bytes()))
	require.False(t, prooffile.IsEvidence([]byte("not json")))
}

func TestAMalformedEvidenceDocumentIsRefused(t *testing.T) {
	raw, err := prooffile.MarshalEvidence(evidenceClaims(t)[0])
	require.NoError(t, err)

	text := string(raw)

	for name, mutated := range map[string]string{
		"unknown field":      strings.Replace(text, `"salt":`, `"extra":1,"salt":`, 1),
		"wrong format":       strings.Replace(text, prooffile.EvidenceFormatV1, "kinetix.notary.evidence.v9", 1),
		"wrong scheme":       strings.Replace(text, `"scheme":1`, `"scheme":2`, 1),
		"unknown stream":     strings.Replace(text, `"stream":"order_lifecycle"`, `"stream":"payments"`, 1),
		"fields of another":  strings.Replace(text, `"lifecycle":`, `"review":`, 1),
		"negative number":    strings.Replace(text, `"number":"0"`, `"number":"-1"`, 1),
		"zero salt":          replaceSalt(t, text, "0x"+strings.Repeat("0", 64)),
		"bad time":           strings.Replace(text, `"occurredAt":"2025`, `"occurredAt":"yesterday 2025`, 1),
		"trailing document":  text + `{}`,
		"short sibling":      strings.Replace(text, `"siblings":["0x`, `"siblings":["0x12","0x`, 1),
		"closed at zero":     strings.Replace(text, `"closedAt":"1760000000"`, `"closedAt":"0"`, 1),
		"missing photo hash": "",
	} {
		if name == "missing photo hash" {
			delivery, err := prooffile.MarshalEvidence(evidenceClaims(t)[3])
			require.NoError(t, err)

			mutated = strings.Replace(string(delivery), `"photoSha256":"0x`, `"photoSha256":"","x":"`, 1)
		}

		require.NotEqual(t, text, mutated, name)

		_, err := prooffile.DecodeEvidence(strings.NewReader(mutated), 1)
		require.Error(t, err, name)
		require.Equal(t, notary.KindInvalidArgument, notary.KindOf(err), name)
	}
}

func replaceSalt(t *testing.T, text, salt string) string {
	t.Helper()

	var doc map[string]any
	require.NoError(t, json.Unmarshal([]byte(text), &doc))

	doc["salt"] = salt
	raw, err := json.Marshal(doc)
	require.NoError(t, err)

	return string(raw)
}
