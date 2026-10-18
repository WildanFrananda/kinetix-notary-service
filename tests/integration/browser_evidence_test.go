package integration_test

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/config"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/prooffile"
	"github.com/WildanFrananda/kinetix-notary-service/internal/verifycli"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

func evidenceSettings(c *chain, registry notary.RegistryIdentity) map[string]string {
	return map[string]string{
		config.RPCURLKey:           c.url,
		config.ChainIDKey:          strconv.FormatUint(registry.ChainID, 10),
		config.EvidenceRegistryKey: registry.Address.Hex(),
		config.RPCTimeoutKey:       "10s",
	}
}

func evidenceEncoded(t *testing.T, claims ...notary.EvidenceClaim) []byte {
	t.Helper()

	documents := make([]json.RawMessage, len(claims))
	for i, claim := range claims {
		raw, err := prooffile.MarshalEvidence(claim)
		require.NoError(t, err)

		documents[i] = raw
	}

	if len(documents) == 1 {
		return documents[0]
	}

	bundle, err := json.Marshal(map[string]any{"proofs": documents, "pending": 0})
	require.NoError(t, err)

	return bundle
}

func TestTheWebAssemblyVerifierAnswersEvidenceExactlyAsTheCLI(t *testing.T) {
	v := buildVerifiers(t)
	c := connect(t)

	registry, batch := c.evidenceBatch(t, notarytest.Records(notarytest.Deliveries(t, 4)))
	c.finalize(t)
	settings := evidenceSettings(c, registry)

	pendingRegistry, pending := c.evidenceBatch(t, notarytest.Records(notarytest.Reviews(t, 2)))

	forged := batch.Claims[1]
	delivery, ok := forged.Record.(notary.DeliveryRecord)
	require.True(t, ok)

	delivery.DriverPrincipalID = "someone-else"
	forged.Record = delivery

	unanchored := batch.Claims[0]
	unanchored.Header.Number = 5

	cases := []struct {
		name     string
		document []byte
		settings map[string]string
		exit     int
	}{
		{"evidence in a final batch", evidenceEncoded(t, batch.Claims[2]), settings, verifycli.ExitIncluded},
		{"a bundle of evidence", evidenceEncoded(t, batch.Claims...), settings, verifycli.ExitIncluded},
		{"evidence that was changed", evidenceEncoded(t, batch.Claims[0], forged), settings, verifycli.ExitRejected},
		{"a batch never anchored", evidenceEncoded(t, unanchored), settings, verifycli.ExitNotAnchored},
		{
			"a batch not final yet",
			evidenceEncoded(t, pending.Claims[0]),
			evidenceSettings(c, pendingRegistry),
			verifycli.ExitNotFinal,
		},
		{
			"a verifier reading another registry",
			evidenceEncoded(t, batch.Claims[0]),
			evidenceSettings(c, pendingRegistry),
			verifycli.ExitWrongLedger,
		},
		{
			"a verifier without a registry",
			evidenceEncoded(t, batch.Claims[0]),
			settingsWith(settings, config.EvidenceRegistryKey, ""),
			verifycli.ExitConfig,
		},
		{
			"an evidence document that is not a proof",
			[]byte(`{"format":"kinetix.notary.evidence.v1","scheme":1}`),
			settings,
			verifycli.ExitMalformedProof,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cli := v.cli(t, c.document, c.settings)
			browser := v.browser(t, c.document, c.settings)

			require.Equal(t, c.exit, cli.exit, "%s%s", cli.stdout, cli.stderr)
			require.Equal(t, cli, browser)
		})
	}
}
