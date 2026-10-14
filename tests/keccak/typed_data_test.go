package keccak_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/keccak"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

var eip712Domain = []apitypes.Type{
	{Name: "name", Type: "string"},
	{Name: "version", Type: "string"},
	{Name: "chainId", Type: "uint256"},
	{Name: "verifyingContract", Type: "address"},
}

func typedDigest(t *testing.T, name string, court notary.CourtIdentity, typed apitypes.TypedData) notary.Digest {
	t.Helper()

	typed.Types["EIP712Domain"] = eip712Domain
	typed.Domain = apitypes.TypedDataDomain{
		Name:              name,
		Version:           "1",
		ChainId:           (*math.HexOrDecimal256)(new(big.Int).SetUint64(court.ChainID)),
		VerifyingContract: common.Address(court.Address).Hex(),
	}

	raw, _, err := apitypes.TypedDataAndHash(typed)
	require.NoError(t, err)

	return notary.Digest(raw)
}

func TestTheRegistrationRequestDigestIsGoEthereumsEIP712(t *testing.T) {
	court := notary.CourtIdentity{ChainID: 84532, Address: notary.Address{0xd0, 0x0d}}
	request := notary.RegistrationRequest{
		Holder:    notary.Address{0xa1, 0x02},
		Principal: notarytest.Holder(t, "9f1d4a3e-1c62-4d0a-9a7b-2f5c8e0b41d7"),
		IssuedAt:  notarytest.Unix(t, 1760000000),
	}

	want := typedDigest(t, keccak.WithholdingCourtName, court, apitypes.TypedData{
		Types: apitypes.Types{
			"RegistrationRequest": {
				{Name: "holder", Type: "address"},
				{Name: "principal", Type: "string"},
				{Name: "issuedAt", Type: "uint64"},
			},
		},
		PrimaryType: "RegistrationRequest",
		Message: apitypes.TypedDataMessage{
			"holder":    common.Address(request.Holder).Hex(),
			"principal": request.Principal.String(),
			"issuedAt":  "1760000000",
		},
	})

	require.Equal(t, want, keccak.NewScheme().RegistrationRequestDigest(court, request))
}

func TestTheRegistrationDigestIsGoEthereumsEIP712(t *testing.T) {
	court := notary.CourtIdentity{ChainID: 84532, Address: notary.Address{0xd0, 0x0d}}
	registration := notary.Registration{Holder: notary.Address{0xa1}, Account: repeated(0x0a), FromEpoch: 1 << 40}

	want := typedDigest(t, keccak.WithholdingCourtName, court, apitypes.TypedData{
		Types: apitypes.Types{
			"Registration": {
				{Name: "holder", Type: "address"},
				{Name: "account", Type: "bytes32"},
				{Name: "fromEpoch", Type: "uint64"},
			},
		},
		PrimaryType: "Registration",
		Message: apitypes.TypedDataMessage{
			"holder":    common.Address(registration.Holder).Hex(),
			"account":   registration.Account.Hex(),
			"fromEpoch": "1099511627776",
		},
	})

	require.Equal(t, want, keccak.NewScheme().RegistrationDigest(court, registration))
}
