package fixtures_test

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/keccak"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

var operationFixturePath = filepath.Join("..", "..", "contracts", "test", "fixtures", "operation-v1.json")

type operationFixture struct {
	Name           string          `json:"name"`
	Version        string          `json:"version"`
	ChainID        string          `json:"chainId"`
	EntryPoint     string          `json:"entryPoint"`
	OperationCount int             `json:"operationCount"`
	Operations     []operationCase `json:"operations"`
}

type operationCase struct {
	Sender                        string `json:"sender"`
	Nonce                         string `json:"nonce"`
	Delegate                      string `json:"delegate"`
	CallData                      string `json:"callData"`
	VerificationGasLimit          string `json:"verificationGasLimit"`
	CallGasLimit                  string `json:"callGasLimit"`
	PreVerificationGas            string `json:"preVerificationGas"`
	MaxPriorityFeePerGas          string `json:"maxPriorityFeePerGas"`
	MaxFeePerGas                  string `json:"maxFeePerGas"`
	Paymaster                     string `json:"paymaster"`
	PaymasterVerificationGasLimit string `json:"paymasterVerificationGasLimit"`
	PaymasterPostOpGasLimit       string `json:"paymasterPostOpGasLimit"`
	PaymasterData                 string `json:"paymasterData"`
	Hash                          string `json:"hash"`
}

func buildOperationFixture() operationFixture {
	scheme := keccak.NewScheme()
	entryPoint := notary.EntryPointIdentity{ChainID: 84532, Address: labelled("kinetix.notary.fixture-entry-point")}
	delegate := labelled("kinetix.notary.fixture-7702-account")
	paymaster := labelled("kinetix.notary.fixture-paymaster")

	operations := []notary.UserOperation{
		{
			Sender:   labelled("kinetix.notary.fixture-sender-0"),
			CallData: []byte{},
		},
		{
			Sender:                        labelled("kinetix.notary.fixture-sender-1"),
			Nonce:                         7,
			Delegate:                      delegate,
			CallData:                      stretched("call", 300),
			CallGasLimit:                  300_000,
			VerificationGasLimit:          150_000,
			PreVerificationGas:            90_000,
			MaxFeePerGas:                  7_700_000,
			MaxPriorityFeePerGas:          1_100_000,
			Paymaster:                     paymaster,
			PaymasterVerificationGasLimit: 100_000,
		},
		{
			Sender:                        labelled("kinetix.notary.fixture-sender-2"),
			Nonce:                         1 << 40,
			CallData:                      stretched("call", 65),
			CallGasLimit:                  1,
			VerificationGasLimit:          2,
			PreVerificationGas:            3,
			MaxFeePerGas:                  4,
			MaxPriorityFeePerGas:          5,
			Paymaster:                     paymaster,
			PaymasterVerificationGasLimit: 6,
			PaymasterPostOpGasLimit:       7,
			PaymasterData:                 stretched("paymaster-data", 65),
		},
		{
			Sender:                        labelled("kinetix.notary.fixture-sender-3"),
			Nonce:                         math.MaxUint64,
			Delegate:                      delegate,
			CallData:                      stretched("call", 4),
			CallGasLimit:                  math.MaxUint64,
			VerificationGasLimit:          math.MaxUint64,
			PreVerificationGas:            math.MaxUint64,
			MaxFeePerGas:                  math.MaxUint64,
			MaxPriorityFeePerGas:          math.MaxUint64,
			Paymaster:                     paymaster,
			PaymasterVerificationGasLimit: math.MaxUint64,
			PaymasterPostOpGasLimit:       math.MaxUint64,
			PaymasterData:                 stretched("paymaster-data", 1),
		},
	}

	out := operationFixture{
		Name:       keccak.EntryPointName,
		Version:    keccak.EntryPointVersion,
		ChainID:    strconv.FormatUint(entryPoint.ChainID, 10),
		EntryPoint: entryPoint.Address.Hex(),
	}

	for _, op := range operations {
		out.Operations = append(out.Operations, operationCase{
			Sender:                        op.Sender.Hex(),
			Nonce:                         u64(op.Nonce),
			Delegate:                      op.Delegate.Hex(),
			CallData:                      "0x" + hex.EncodeToString(op.CallData),
			VerificationGasLimit:          u64(op.VerificationGasLimit),
			CallGasLimit:                  u64(op.CallGasLimit),
			PreVerificationGas:            u64(op.PreVerificationGas),
			MaxPriorityFeePerGas:          u64(op.MaxPriorityFeePerGas),
			MaxFeePerGas:                  u64(op.MaxFeePerGas),
			Paymaster:                     op.Paymaster.Hex(),
			PaymasterVerificationGasLimit: u64(op.PaymasterVerificationGasLimit),
			PaymasterPostOpGasLimit:       u64(op.PaymasterPostOpGasLimit),
			PaymasterData:                 "0x" + hex.EncodeToString(op.PaymasterData),
			Hash:                          scheme.OperationHash(entryPoint, op).Hex(),
		})
	}

	out.OperationCount = len(out.Operations)

	return out
}

func stretched(label string, size int) []byte {
	out := make([]byte, 0, size+32)
	for i := 0; len(out) < size; i++ {
		block := filled(label, i)
		out = append(out, block[:]...)
	}

	return out[:size]
}

func TestTheSolidityOperationFixtureMatchesTheGoScheme(t *testing.T) {
	encoded, err := json.MarshalIndent(buildOperationFixture(), "", "  ")
	require.NoError(t, err)

	encoded = append(encoded, '\n')

	if *update {
		require.NoError(t, os.WriteFile(operationFixturePath, encoded, 0o644))
		return
	}

	committed, err := os.ReadFile(operationFixturePath)
	require.NoError(t, err, "no fixture yet: run `bin/dev go test ./tests/fixtures -update`")
	require.True(
		t,
		bytes.Equal(committed, encoded),
		"%s is out of date with the Go scheme: run `bin/dev go test ./tests/fixtures -update`",
		operationFixturePath,
	)
}
