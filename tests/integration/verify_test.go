package integration_test

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/config"
	"github.com/WildanFrananda/kinetix-notary-service/internal/contractgen/court"
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
	court    common.Address
	signer   notary.Address
	key      *ecdsa.PrivateKey
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

	key, signer := notarytest.StatementKey(t, "integration-signer")
	courtAddress := c.deployCourt(t, address, signer)

	if final {
		c.finalize(t)
	}

	return anchored{
		chain:    c,
		address:  address,
		court:    courtAddress,
		signer:   signer,
		key:      key,
		identity: c.identity(address),
		tree:     tree,
		entries:  entries,
	}
}

func (a anchored) courtIdentity() notary.CourtIdentity {
	return notary.CourtIdentity{ChainID: a.identity.ChainID, Address: notary.Address(a.court)}
}

func (a anchored) sign(t *testing.T, claim notary.Claim, key *ecdsa.PrivateKey, signer notary.Address) notary.Claim {
	t.Helper()

	scheme := keccak.NewScheme()
	digest := scheme.StatementDigest(a.courtIdentity(), notary.StatementOf(scheme, claim))

	signature, err := ethereum.NewStatementSigner([]*ecdsa.PrivateKey{key}).Sign(signer, digest)
	require.NoError(t, err)

	claim.Statement = &notary.StatementSignature{Court: notary.Address(a.court), Signature: signature}

	return claim
}

func (a anchored) claim(t *testing.T, index uint64) notary.Claim {
	t.Helper()

	proof, err := a.tree.Prove(index)
	require.NoError(t, err)

	return notary.Claim{Ledger: a.identity, Epoch: 0, Entry: a.entries[index], Proof: proof}
}

func (a anchored) verifier() *verifying.Verifier {
	scheme := keccak.NewScheme()

	return verifying.NewVerifier(a.chain.ledgerAt(a.identity), scheme, scheme, a.chain.courtAt(a.identity, a.court))
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

	settings := ethereum.Settings{
		RPCURL:  "http://127.0.0.1:1",
		Ledger:  a.identity,
		Court:   notary.Address(a.court),
		Timeout: 2 * time.Second,
	}
	scheme := keccak.NewScheme()
	verifier := verifying.NewVerifier(
		ethereum.NewLedger(client, settings, scheme),
		scheme,
		scheme,
		ethereum.NewCourt(client, settings),
	)

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

	g := withDirectory(t, a)
	env := map[string]string{
		config.RPCURLKey:           a.chain.url,
		config.ChainIDKey:          strconv.FormatUint(a.identity.ChainID, 10),
		config.DirectoryAddressKey: notary.Address(g.directory).Hex(),
		config.RPCTimeoutKey:       "10s",
	}
	container, err := wiring.Verifier(func(key string) (string, bool) {
		value, ok := env[key]
		return value, ok
	})
	require.NoError(t, err)

	var out bytes.Buffer

	code := -1
	err = container.Invoke(func(client *ethclient.Client, scheme notary.Scheme, v *verifying.DirectoryVerifier) {
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

func TestAStatementSignedByTheCourtsSignerIsAttributableOnARealChain(t *testing.T) {
	a := anchorEpochZero(t, 7, true)
	claim := a.sign(t, a.claim(t, 3), a.key, a.signer)

	got, err := a.verifier().Verify(context.Background(), claim)

	require.NoError(t, err)
	require.Equal(t, notary.VerdictIncluded, got.Verdict)
	require.Equal(t, notary.StatementAttributable, got.Check)
	require.Equal(t, a.signer, got.Signer)
	require.Equal(t, uint64(0), got.Faults)
	require.Equal(t, a.courtIdentity(), got.Court)
}

func TestAStatementSignedByAnyOtherKeyIsNotAttributableOnARealChain(t *testing.T) {
	a := anchorEpochZero(t, 7, true)
	key, impostor := notarytest.StatementKey(t, "integration-impostor")

	got, err := a.verifier().Verify(context.Background(), a.sign(t, a.claim(t, 3), key, impostor))

	require.NoError(t, err)
	require.Equal(t, notary.StatementUnattributable, got.Check)
}

func TestAFalseSignedProofIsReportedProvableAndTheCourtConvictsIt(t *testing.T) {
	a := anchorEpochZero(t, 7, true)
	claim := a.claim(t, 3)
	claim.Entry.Amount++
	claim = a.sign(t, claim, a.key, a.signer)

	before, err := a.verifier().Verify(context.Background(), claim)
	require.NoError(t, err)
	require.Equal(t, notary.VerdictRejected, before.Verdict)
	require.True(t, before.Provable)
	require.Equal(t, uint64(0), before.Faults)

	scheme := keccak.NewScheme()
	statement := notary.StatementOf(scheme, claim)

	siblings := make([]court.SumNode, len(claim.Proof.Siblings))
	for i, sibling := range claim.Proof.Siblings {
		siblings[i] = court.SumNode{Hash: sibling.Hash, Sum: sibling.Sum}
	}

	binding := court.NewStatementCourt()
	packed := binding.PackProveBadProof(
		court.EpochStatement{
			Ledger:       common.Address(statement.Ledger),
			Epoch:        statement.Epoch,
			Commitment:   statement.Commitment,
			Amount:       statement.Amount,
			Index:        statement.Index,
			SiblingsHash: statement.SiblingsHash,
		},
		claim.Statement.Signature[:],
		siblings,
	)
	tx, err := bind.Transact(binding.Instance(a.chain.client, a.court), a.chain.transactor(), packed)
	require.NoError(t, err)
	receipt, err := bind.WaitMined(context.Background(), a.chain.client, tx.Hash())
	require.NoError(t, err)
	require.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)

	after, err := a.verifier().Verify(context.Background(), claim)
	require.NoError(t, err)
	require.Equal(t, uint64(1), after.Faults)

	total, err := a.chain.courtAt(a.identity, a.court).TotalFaults(context.Background())
	require.NoError(t, err)
	require.Equal(t, uint64(1), total)
}

func TestACourtThatJudgesAnotherLedgerIsTheWrongCourt(t *testing.T) {
	a := anchorEpochZero(t, 5, true)
	other := a.chain.deploy(t)
	_, signer := notarytest.StatementKey(t, "integration-signer")
	foreign := a.chain.deployCourt(t, other, signer)

	_, err := a.chain.courtAt(a.identity, foreign).Identity(context.Background())

	require.Equal(t, notary.KindWrongCourt, notary.KindOf(err))
}

func TestALedgerIsNotACourt(t *testing.T) {
	a := anchorEpochZero(t, 5, true)

	_, err := a.chain.courtAt(a.identity, a.address).Identity(context.Background())

	require.Equal(t, notary.KindWrongCourt, notary.KindOf(err))
}
