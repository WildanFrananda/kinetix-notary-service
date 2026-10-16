package bundler_test

import (
	"context"
	"encoding/json"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/bundler"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type call struct {
	Method string            `json:"method"`
	Params []json.RawMessage `json:"params"`
}

func serve(t *testing.T, answer func(call) (int, string)) (*bundler.Client, *[]call) {
	t.Helper()

	var seen []call

	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			raw, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			require.Equal(t, "application/json", r.Header.Get("Content-Type"))

			var c call
			require.NoError(t, json.Unmarshal(raw, &c))

			seen = append(seen, c)
			status, body := answer(c)
			w.WriteHeader(status)
			_, _ = w.Write([]byte(body))
		}),
	)
	t.Cleanup(server.Close)

	return bundler.NewClient(server.URL+"/v3/84532/secret-bundler-key", 2*time.Second), &seen
}

func operation() notary.UserOperation {
	return notary.UserOperation{
		Sender:                        notary.Address{0x01},
		Nonce:                         9,
		Delegate:                      notary.Address{0x77},
		CallData:                      []byte{0xb6, 0x1d},
		CallGasLimit:                  300_000,
		VerificationGasLimit:          150_000,
		PreVerificationGas:            90_000,
		MaxFeePerGas:                  7_700_000,
		MaxPriorityFeePerGas:          1_100_000,
		Paymaster:                     notary.Address{0xfa},
		PaymasterVerificationGasLimit: 100_000,
		PaymasterData:                 []byte{0x0a},
		Signature:                     []byte{0x5e},
	}
}

func TestAFirstOperationCarriesTheDelegationMarkerAndTheAuthorization(t *testing.T) {
	client, seen := serve(t, func(call) (int, string) {
		return http.StatusOK, `{"jsonrpc":"2.0","id":1,"result":{"preVerificationGas":"0x15f90",` +
			`"verificationGasLimit":"0x249f0","callGasLimit":"0x493e0","paymasterVerificationGasLimit":"0x186a0"}}`
	})
	delegation := &notary.Delegation{
		ChainID:  84532,
		Delegate: notary.Address{0x77},
		Nonce:    3,
		YParity:  1,
		R:        notary.Digest{31: 0x04},
		S:        notary.Digest{0x05},
	}

	gas, err := client.Estimate(context.Background(), operation(), notary.Address{0x43}, delegation)

	require.NoError(t, err)
	require.Equal(
		t,
		notary.OperationGas{
			PreVerificationGas:            90_000,
			VerificationGasLimit:          150_000,
			CallGasLimit:                  300_000,
			PaymasterVerificationGasLimit: 100_000,
		},
		gas,
	)
	require.Equal(t, "eth_estimateUserOperationGas", (*seen)[0].Method)

	var sent map[string]any
	require.NoError(t, json.Unmarshal((*seen)[0].Params[0], &sent))
	require.Equal(t, "0x7702", sent["factory"])
	require.Equal(t, "0x9", sent["nonce"])
	require.Equal(t, "0x493e0", sent["callGasLimit"])
	require.Equal(t, "0xb61d", sent["callData"])
	require.Equal(t, "0x0a", sent["paymasterData"])
	require.Equal(t, notary.Address{0xfa}.Hex(), sent["paymaster"])

	auth, ok := sent["eip7702Auth"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "0x14a34", auth["chainId"])
	require.Equal(t, notary.Address{0x77}.Hex(), auth["address"])
	require.Equal(t, "0x3", auth["nonce"])
	require.Equal(t, "0x1", auth["yParity"])
	require.Equal(t, "0x4", auth["r"])
	require.Equal(t, "0x5"+strings.Repeat("0", 62), auth["s"])

	var entryPoint string
	require.NoError(t, json.Unmarshal((*seen)[0].Params[1], &entryPoint))
	require.Equal(t, notary.Address{0x43}.Hex(), entryPoint)
}

func TestAnAlreadyDelegatedAccountSendsNeitherMarkerNorAuthorization(t *testing.T) {
	client, seen := serve(t, func(call) (int, string) {
		return http.StatusOK, `{"jsonrpc":"2.0","id":1,"result":"` + notary.Digest{0xab}.Hex() + `"}`
	})
	op := operation()
	op.Delegate = notary.Address{}

	hash, err := client.Send(context.Background(), op, notary.Address{0x43}, nil)

	require.NoError(t, err)
	require.Equal(t, notary.Digest{0xab}, hash)

	var sent map[string]any
	require.NoError(t, json.Unmarshal((*seen)[0].Params[0], &sent))
	require.NotContains(t, sent, "factory")
	require.NotContains(t, sent, "eip7702Auth")
}

func TestAReceiptIsReadAndAMissingOneIsNotYetAFailure(t *testing.T) {
	answers := []string{
		`{"jsonrpc":"2.0","id":1,"result":null}`,
		`{"jsonrpc":"2.0","id":1,"result":{"success":false,"actualGasCost":"0x77359400","reason":"0x08c379a0",` +
			`"receipt":{"transactionHash":"` + notary.Digest{0x7a}.Hex() + `"}}}`,
	}
	client, _ := serve(t, func(call) (int, string) {
		answer := answers[0]
		answers = answers[1:]

		return http.StatusOK, answer
	})

	_, found, err := client.Receipt(context.Background(), notary.Digest{0xab})
	require.NoError(t, err)
	require.False(t, found)

	receipt, found, err := client.Receipt(context.Background(), notary.Digest{0xab})
	require.NoError(t, err)
	require.True(t, found)
	require.False(t, receipt.Success)
	require.Equal(t, notary.Digest{0x7a}, receipt.TransactionHash)
	require.Equal(t, big.NewInt(2_000_000_000), receipt.ActualGasCost)
	require.Equal(t, "0x08c379a0", receipt.Reason)
}

func TestARefusalIsNotSponsoredAndAnOutageIsUnavailable(t *testing.T) {
	for _, c := range []struct {
		status int
		body   string
		kind   notary.Kind
	}{
		{http.StatusOK, `{"jsonrpc":"2.0","id":1,"error":{"code":-32500,"message":"AA33 reverted NotSignedByTheNotary"}}`, notary.KindNotSponsored},
		{http.StatusOK, `{"jsonrpc":"2.0","id":1,"error":{"code":-32501,"message":"paymaster validation"}}`, notary.KindNotSponsored},
		{http.StatusOK, `{"jsonrpc":"2.0","id":1,"error":{"code":-32504,"message":"paymaster throttled"}}`, notary.KindSourceUnavailable},
		{http.StatusTooManyRequests, `Too Many Requests`, notary.KindSourceUnavailable},
		{http.StatusOK, `{"jsonrpc":"2.0","id":1,"result":"not a hash"}`, notary.KindSourceRefused},
	} {
		client, _ := serve(t, func(call) (int, string) { return c.status, c.body })

		_, err := client.Send(context.Background(), operation(), notary.Address{0x43}, nil)

		require.Equal(t, c.kind, notary.KindOf(err), c.body)
		require.NotContains(t, err.Error(), "secret-bundler-key")
	}
}

func TestAnUnreachableBundlerIsUnavailableAndNeverNamesItsURL(t *testing.T) {
	client := bundler.NewClient("http://127.0.0.1:1/v3/84532/secret-bundler-key", time.Second)

	_, err := client.Send(context.Background(), operation(), notary.Address{0x43}, nil)

	require.Equal(t, notary.KindSourceUnavailable, notary.KindOf(err))
	require.NotContains(t, err.Error(), "secret-bundler-key")
}
