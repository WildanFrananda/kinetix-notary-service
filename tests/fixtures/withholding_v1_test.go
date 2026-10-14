package fixtures_test

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/ethereum"
	"github.com/WildanFrananda/kinetix-notary-service/internal/keccak"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

var withholdingFixturePath = filepath.Join("..", "..", "contracts", "test", "fixtures", "withholding-v1.json")

type withholdingFixture struct {
	Name              string             `json:"name"`
	Version           string             `json:"version"`
	ChainID           string             `json:"chainId"`
	Court             string             `json:"court"`
	WithholdingCourt  string             `json:"withholdingCourt"`
	Signer            string             `json:"signer"`
	RegistrationCount int                `json:"registrationCount"`
	Registrations     []registrationCase `json:"registrations"`
	AccountCount      int                `json:"accountCount"`
	Accounts          []accountCase      `json:"accounts"`
}

type registrationCase struct {
	Holder    string `json:"holder"`
	Account   string `json:"account"`
	FromEpoch string `json:"fromEpoch"`
	Digest    string `json:"digest"`
	Signature string `json:"signature"`
}

type accountCase struct {
	Holder     string `json:"holder"`
	Kind       string `json:"kind"`
	Secret     string `json:"secret"`
	Account    string `json:"account"`
	Salt       string `json:"salt"`
	Commitment string `json:"commitment"`
}

func filled(label string, i int) [32]byte {
	var out [32]byte
	copy(out[:], crypto.Keccak256([]byte(label), []byte(strconv.Itoa(i))))

	return out
}

func buildWithholdingFixture(t *testing.T) withholdingFixture {
	t.Helper()

	scheme := keccak.NewScheme()
	key, err := crypto.ToECDSA(crypto.Keccak256([]byte("kinetix.notary.withholding-fixture")))
	require.NoError(t, err)

	signerAddress := notary.Address(crypto.PubkeyToAddress(key.PublicKey))
	signer := ethereum.NewStatementSigner([]*ecdsa.PrivateKey{key})
	court := notary.CourtIdentity{ChainID: 84532, Address: labelled("kinetix.notary.fixture-withholding-court")}

	out := withholdingFixture{
		Name:             keccak.WithholdingCourtName,
		Version:          keccak.WithholdingCourtVersion,
		ChainID:          strconv.FormatUint(court.ChainID, 10),
		Court:            labelled("kinetix.notary.fixture-statement-court").Hex(),
		WithholdingCourt: court.Address.Hex(),
		Signer:           signerAddress.Hex(),
	}

	holders := []string{
		"9f1d4a3e-1c62-4d0a-9a7b-2f5c8e0b41d7",
		"UNASSIGNED_DRIVER_SHIPPING_FEE",
		"p",
		strings.Repeat("x", notary.MaxHolderBytes),
	}
	kinds := []notary.AccountKind{
		notary.AccountCustomerBalance,
		notary.AccountSuspensePendingEscrow,
		notary.AccountMerchantAvailable,
		notary.AccountDriverPendingEscrow,
	}
	fromEpochs := []uint64{0, 3, 1 << 40, math.MaxUint64}

	for i, name := range holders {
		holder := notarytest.Holder(t, name)
		secret := notary.AccountSecret(filled("secret", i))
		salt := notary.Salt(filled("salt", i))
		account := scheme.AccountID(holder, kinds[i], secret)

		out.Accounts = append(out.Accounts, accountCase{
			Holder:     name,
			Kind:       strconv.Itoa(int(kinds[i])),
			Secret:     secret.Hex(),
			Account:    account.Hex(),
			Salt:       salt.Hex(),
			Commitment: scheme.AccountCommitment(account, salt).Hex(),
		})

		registration := notary.Registration{
			Holder:    labelled("kinetix.notary.fixture-holder-" + strconv.Itoa(i)),
			Account:   account,
			FromEpoch: fromEpochs[i],
		}
		digest := scheme.RegistrationDigest(court, registration)
		signature, err := signer.Sign(signerAddress, digest)
		require.NoError(t, err)

		public, err := crypto.SigToPub(digest[:], append(signature[:64:64], signature[64]-27))
		require.NoError(t, err)
		require.Equal(t, signerAddress, notary.Address(crypto.PubkeyToAddress(*public)))

		out.Registrations = append(out.Registrations, registrationCase{
			Holder:    registration.Holder.Hex(),
			Account:   registration.Account.Hex(),
			FromEpoch: u64(registration.FromEpoch),
			Digest:    digest.Hex(),
			Signature: signature.Hex(),
		})
	}

	out.RegistrationCount = len(out.Registrations)
	out.AccountCount = len(out.Accounts)

	return out
}

func TestTheSolidityWithholdingFixtureMatchesTheGoScheme(t *testing.T) {
	encoded, err := json.MarshalIndent(buildWithholdingFixture(t), "", "  ")
	require.NoError(t, err)

	encoded = append(encoded, '\n')

	if *update {
		require.NoError(t, os.WriteFile(withholdingFixturePath, encoded, 0o644))
		return
	}

	committed, err := os.ReadFile(withholdingFixturePath)
	require.NoError(t, err, "no fixture yet: run `bin/dev go test ./tests/fixtures -update`")
	require.True(
		t,
		bytes.Equal(committed, encoded),
		"%s is out of date with the Go scheme: run `bin/dev go test ./tests/fixtures -update`",
		withholdingFixturePath,
	)
}
