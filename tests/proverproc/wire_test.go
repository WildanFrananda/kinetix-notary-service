package proverproc_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/proverproc"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

func TestARequestCrossesTheWireUnchanged(t *testing.T) {
	request := notary.ProofRequest{
		Number:          1<<63 + 5,
		Depth:           8,
		Entries:         notarytest.Entries(4, 7),
		TreeRoot:        notary.Digest{0x01, 0x02},
		TotalCommitment: notary.Digest{0x03},
		Blinding:        notary.Blinding{0x04},
	}
	request.Entries[2].Amount = ^uint64(0)

	var wire bytes.Buffer
	require.NoError(t, proverproc.EncodeRequest(&wire, request))

	decoded, err := proverproc.DecodeRequest(&wire)
	require.NoError(t, err)
	require.Equal(t, request, decoded)
}

func TestAMalformedRequestIsRefused(t *testing.T) {
	cases := map[string]string{
		"not JSON":         "{",
		"an unknown field": `{"number":"1","depth":2,"treeRoot":"0x` + strings.Repeat("00", 32) + `","extra":1}`,
		"a negative epoch": `{"number":"-1"}`,
		"a short root":     `{"number":"1","treeRoot":"0x00"}`,
	}

	for name, text := range cases {
		_, err := proverproc.DecodeRequest(strings.NewReader(text))
		require.Equal(t, notary.KindInvalidArgument, notary.KindOf(err), name)
	}
}
