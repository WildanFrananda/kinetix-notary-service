package fixtures_test

import (
	"bytes"
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/keccak"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

var approvalFixturePath = filepath.Join("..", "..", "contracts", "test", "fixtures", "approval-v1.json")

type approvalFixture struct {
	Name          string         `json:"name"`
	Version       string         `json:"version"`
	ChainID       string         `json:"chainId"`
	Council       string         `json:"council"`
	Ledger        string         `json:"ledger"`
	Member        string         `json:"member"`
	ApprovalCount int            `json:"approvalCount"`
	Approvals     []approvalCase `json:"approvals"`
}

type approvalCase struct {
	Number         string `json:"number"`
	PreviousDigest string `json:"previousDigest"`
	Root           string `json:"root"`
	Total          string `json:"total"`
	LeafCount      string `json:"leafCount"`
	AsOf           string `json:"asOf"`
	EpochDigest    string `json:"epochDigest"`
	Digest         string `json:"digest"`
	Signature      string `json:"signature"`
}

func buildApprovalFixture(t *testing.T) approvalFixture {
	t.Helper()

	scheme := keccak.NewScheme()
	key, err := crypto.ToECDSA(crypto.Keccak256([]byte("kinetix.notary.approval-fixture")))
	require.NoError(t, err)

	member := notary.Address(crypto.PubkeyToAddress(key.PublicKey))
	council := notary.CouncilIdentity{ChainID: 84532, Address: labelled("kinetix.notary.fixture-council")}
	ledger := labelled("kinetix.notary.fixture-council-ledger")

	out := approvalFixture{
		Name:    keccak.CouncilName,
		Version: keccak.CouncilVersion,
		ChainID: strconv.FormatUint(council.ChainID, 10),
		Council: council.Address.Hex(),
		Ledger:  ledger.Hex(),
		Member:  member.Hex(),
	}

	headers := []notary.EpochHeader{
		{
			Number:    0,
			Root:      scheme.EmptyHash(),
			AsOf:      notarytest.Unix(t, 1),
			LeafCount: 0,
		},
		{
			Number:         41,
			PreviousDigest: notary.Digest(filled("previous", 41)),
			Root:           notary.Digest(filled("root", 41)),
			Total:          50_000_000,
			LeafCount:      7,
			AsOf:           notarytest.Unix(t, 1791961931),
		},
		{
			Number:         1 << 40,
			PreviousDigest: notary.Digest(filled("previous", 2)),
			Root:           notary.Digest(filled("root", 2)),
			Total:          math.MaxInt64,
			LeafCount:      notary.MaxLeaves,
			AsOf:           notarytest.Unix(t, math.MaxInt64),
		},
		{
			Number:         math.MaxUint64,
			PreviousDigest: notary.Digest(filled("previous", 3)),
			Root:           notary.Digest(filled("root", 3)),
			Total:          math.MaxUint64,
			LeafCount:      math.MaxUint64,
			AsOf:           notarytest.Unix(t, 1),
		},
	}

	for _, header := range headers {
		epoch := scheme.EpochDigest(header)
		approval := scheme.ApprovalDigest(council, ledger, header.Number, epoch)
		signature, err := crypto.Sign(approval[:], key)
		require.NoError(t, err)

		signature[64] += 27

		out.Approvals = append(out.Approvals, approvalCase{
			Number:         u64(header.Number),
			PreviousDigest: header.PreviousDigest.Hex(),
			Root:           header.Root.Hex(),
			Total:          u64(header.Total),
			LeafCount:      u64(header.LeafCount),
			AsOf:           u64(header.AsOf.Seconds()),
			EpochDigest:    epoch.Hex(),
			Digest:         approval.Hex(),
			Signature:      notary.Signature(signature).Hex(),
		})
	}

	out.ApprovalCount = len(out.Approvals)

	return out
}

func TestTheSolidityApprovalFixtureMatchesTheGoScheme(t *testing.T) {
	encoded, err := json.MarshalIndent(buildApprovalFixture(t), "", "  ")
	require.NoError(t, err)

	encoded = append(encoded, '\n')

	if *update {
		require.NoError(t, os.WriteFile(approvalFixturePath, encoded, 0o644))
		return
	}

	committed, err := os.ReadFile(approvalFixturePath)
	require.NoError(t, err, "no fixture yet: run `bin/dev go test ./tests/fixtures -update`")
	require.True(
		t,
		bytes.Equal(committed, encoded),
		"%s is out of date with the Go scheme: run `bin/dev go test ./tests/fixtures -update`",
		approvalFixturePath,
	)
}
