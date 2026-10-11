package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/config"
	"github.com/WildanFrananda/kinetix-notary-service/internal/ethereum"
	"github.com/WildanFrananda/kinetix-notary-service/internal/keccak"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/prooffile"
	"github.com/WildanFrananda/kinetix-notary-service/internal/verifycli"
	"github.com/WildanFrananda/kinetix-notary-service/internal/verifying"
	"github.com/WildanFrananda/kinetix-notary-service/internal/wiring"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

type anchored struct {
	chain    *chain
	address  common.Address
	identity notary.LedgerIdentity
	tree     *notary.Tree
	entries  []notary.Entry
}

func anchorEpochZero(t *testing.T, count int, final bool) anchored {
	t.Helper()
	c := connect(t)
	address := c.deploy(t)
	entries := notarytest.Entries(uint64(count)+7, count)
	tree, err := notary.BuildTree(keccak.NewScheme(), 0, entries)
	require.NoError(t, err)
	c.anchor(t, address, tree, notary.Digest{})
	if final {
		c.finalize(t)
	}
	return anchored{chain: c, address: address, identity: c.identity(address), tree: tree, entries: entries}
}

func (a anchored) claim(t *testing.T, index uint64) notary.Claim {
	t.Helper()
	proof, err := a.tree.Prove(index)
	require.NoError(t, err)
	return notary.Claim{Ledger: a.identity, Epoch: 0, Entry: a.entries[index], Proof: proof}
}

func (a anchored) verifier() *verifying.Verifier {
	return verifying.NewVerifier(a.chain.ledgerAt(a.identity), keccak.NewScheme())
}

func TestEveryAccountOfAFinalEpochVerifiesAgainstARealChain(t *testing.T) {
	a := anchorEpochZero(t, 9, true)

	for i := range a.entries {
		got, err := a.verifier().Verify(context.Background(), a.claim(t, uint64(i)))

		require.NoError(t, err)
		require.Equal(t, notary.VerdictIncluded, got.Verdict, "account %d", i)
		require.Equal(t, a.tree.Root().Hash, got.Epoch.Root)
		require.Equal(t, a.tree.Root().Sum, got.Epoch.Total)
		require.Equal(t, uint64(9), got.Epoch.LeafCount)
		require.NotZero(t, got.Epoch.AnchoredAt)
	}
}

func TestATamperedClaimIsRejectedByTheChainAndTheSchemeAlike(t *testing.T) {
	a := anchorEpochZero(t, 9, true)
	claim := a.claim(t, 4)
	claim.Entry.Amount++

	got, err := a.verifier().Verify(context.Background(), claim)

	require.NoError(t, err)
	require.Equal(t, notary.VerdictRejected, got.Verdict)
	require.Contains(t, got.Reason, "the proof sums to")
}

func TestAnEpochThatIsNotFinalYetIsSaidToBeSo(t *testing.T) {
	a := anchorEpochZero(t, 5, false)

	_, err := a.verifier().Verify(context.Background(), a.claim(t, 1))

	require.Equal(t, notary.KindEpochNotFinal, notary.KindOf(err), "%v", err)
}

func TestAnEpochTheLedgerNeverAnchoredIsSaidToBeSo(t *testing.T) {
	a := anchorEpochZero(t, 5, true)
	claim := a.claim(t, 1)
	claim.Epoch = 5

	_, err := a.verifier().Verify(context.Background(), claim)

	require.Equal(t, notary.KindEpochNotAnchored, notary.KindOf(err), "%v", err)
}

func TestAnAddressWithNoContractIsTheWrongLedger(t *testing.T) {
	c := connect(t)
	identity := c.identity(common.HexToAddress("0x000000000000000000000000000000000000dEaD"))

	_, err := c.ledgerAt(identity).Identity(context.Background())

	require.Equal(t, notary.KindWrongLedger, notary.KindOf(err), "%v", err)
}

func TestAContractThatIsNotALedgerIsTheWrongLedger(t *testing.T) {
	c := connect(t)
	cases := map[string]string{
		"answers nothing": "0x6001600c60003960016000f300",
		"always reverts":  "0x6005600c60003960056000f360006000fd",
	}
	for name, initCode := range cases {
		t.Run(name, func(t *testing.T) {
			identity := c.identity(c.deployRaw(t, initCode))

			_, err := c.ledgerAt(identity).Identity(context.Background())

			require.Equal(t, notary.KindWrongLedger, notary.KindOf(err), "%v", err)
		})
	}
}

func TestAnRPCServingAnotherChainIsTheWrongLedger(t *testing.T) {
	a := anchorEpochZero(t, 3, true)
	identity := a.identity
	identity.ChainID++

	_, err := a.chain.ledgerAt(identity).Identity(context.Background())

	require.Equal(t, notary.KindWrongLedger, notary.KindOf(err), "%v", err)
}

func TestAnUnreachableRPCIsUnknownAndNeverAVerdict(t *testing.T) {
	a := anchorEpochZero(t, 3, true)
	client, err := ethclient.Dial("http://127.0.0.1:1")
	require.NoError(t, err)
	defer client.Close()
	settings := ethereum.Settings{RPCURL: "http://127.0.0.1:1", Ledger: a.identity, Timeout: 2 * time.Second}
	verifier := verifying.NewVerifier(ethereum.NewLedger(client, settings, keccak.NewScheme()), keccak.NewScheme())

	got, err := verifier.Verify(context.Background(), a.claim(t, 0))

	require.Equal(t, notary.KindLedgerUnavailable, notary.KindOf(err), "%v", err)
	require.Equal(t, notary.Verification{}, got)
}

func TestTheComposedVerifierAnswersFromAProofFile(t *testing.T) {
	a := anchorEpochZero(t, 6, true)
	var doc bytes.Buffer
	require.NoError(t, prooffile.Encode(&doc, a.claim(t, 2), keccak.NewScheme().Version()))
	path := filepath.Join(t.TempDir(), "proof.json")
	require.NoError(t, os.WriteFile(path, doc.Bytes(), 0o600))
	env := map[string]string{
		config.RPCURLKey:        a.chain.url,
		config.ChainIDKey:       strconv.FormatUint(a.identity.ChainID, 10),
		config.LedgerAddressKey: a.identity.Address.Hex(),
		config.RPCTimeoutKey:    "10s",
	}
	container, err := wiring.Verifier(func(key string) (string, bool) {
		value, ok := env[key]
		return value, ok
	})
	require.NoError(t, err)

	var out bytes.Buffer
	code := -1
	err = container.Invoke(func(client *ethclient.Client, scheme notary.Scheme, v *verifying.Verifier) {
		defer client.Close()
		code = verifycli.Run(context.Background(), path, scheme.Version(), v, &out)
	})

	require.NoError(t, err)
	require.Equal(t, verifycli.ExitIncluded, code, out.String())
	var report map[string]string
	require.NoError(t, json.Unmarshal(out.Bytes(), &report))
	require.Equal(t, "included", report["verdict"])
	require.Equal(t, a.identity.Address.Hex(), report["ledger"])
}
