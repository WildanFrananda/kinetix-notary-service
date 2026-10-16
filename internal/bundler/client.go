package bundler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

const (
	maxResponseBytes   = 1 << 20
	codeMethodNotFound = -32601
	codeInternal       = -32603
	codeThrottled      = -32504
	codeLimitExceeded  = -32005
)

type Client struct {
	url     string
	http    *http.Client
	timeout time.Duration
}

var _ notary.Bundler = (*Client)(nil)

func NewClient(url string, timeout time.Duration) *Client {
	return &Client{url: url, http: &http.Client{}, timeout: timeout}
}

func (c *Client) Estimate(
	ctx context.Context,
	op notary.UserOperation,
	entryPoint notary.Address,
	delegation *notary.Delegation,
) (notary.OperationGas, error) {
	const method = "eth_estimateUserOperationGas"

	var estimate gasEstimate
	if err := c.call(ctx, method, []any{operationOf(op, delegation), entryPoint.Hex()}, &estimate); err != nil {
		return notary.OperationGas{}, err
	}

	gas := notary.OperationGas{}
	for _, field := range []struct {
		text string
		into *uint64
		name string
	}{
		{estimate.PreVerificationGas, &gas.PreVerificationGas, "preVerificationGas"},
		{estimate.VerificationGasLimit, &gas.VerificationGasLimit, "verificationGasLimit"},
		{estimate.CallGasLimit, &gas.CallGasLimit, "callGasLimit"},
		{estimate.PaymasterVerificationGasLimit, &gas.PaymasterVerificationGasLimit, "paymasterVerificationGasLimit"},
	} {
		value, err := quantity(field.text)
		if err != nil {
			return notary.OperationGas{}, notary.Errf(
				notary.KindSourceRefused,
				"bundler.Client.Estimate",
				err,
				"the bundler's estimate carries no usable %s",
				field.name,
			)
		}

		*field.into = value
	}

	return gas, nil
}

func (c *Client) Send(
	ctx context.Context,
	op notary.UserOperation,
	entryPoint notary.Address,
	delegation *notary.Delegation,
) (notary.Digest, error) {
	var hash string
	if err := c.call(
		ctx,
		"eth_sendUserOperation",
		[]any{operationOf(op, delegation), entryPoint.Hex()},
		&hash,
	); err != nil {
		return notary.Digest{}, err
	}

	digest, err := notary.ParseDigest(hash)
	if err != nil {
		return notary.Digest{}, notary.Errf(
			notary.KindSourceRefused,
			"bundler.Client.Send",
			err,
			"the bundler answered no operation hash",
		)
	}

	return digest, nil
}

func (c *Client) Receipt(ctx context.Context, hash notary.Digest) (notary.OperationReceipt, bool, error) {
	const op = "bundler.Client.Receipt"

	var answer *receipt
	if err := c.call(ctx, "eth_getUserOperationReceipt", []any{hash.Hex()}, &answer); err != nil {
		return notary.OperationReceipt{}, false, err
	}

	if answer == nil {
		return notary.OperationReceipt{}, false, nil
	}

	transaction, err := notary.ParseDigest(answer.Receipt.TransactionHash)
	if err != nil {
		return notary.OperationReceipt{}, false, notary.Errf(
			notary.KindSourceRefused,
			op,
			err,
			"the receipt names no transaction",
		)
	}

	cost, ok := new(big.Int).SetString(strings.TrimPrefix(answer.ActualGasCost, "0x"), 16)
	if !ok {
		return notary.OperationReceipt{}, false, notary.Errf(
			notary.KindSourceRefused,
			op,
			nil,
			"the receipt carries no gas cost",
		)
	}

	return notary.OperationReceipt{
		Hash:            hash,
		Success:         answer.Success,
		TransactionHash: transaction,
		ActualGasCost:   cost,
		Reason:          answer.Reason,
	}, true, nil
}

func (c *Client) call(ctx context.Context, method string, params []any, out any) error {
	op := "bundler.Client." + method

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	body, err := json.Marshal(request{JSONRPC: "2.0", ID: 1, Method: method, Params: params})
	if err != nil {
		return notary.Errf(notary.KindInvalidArgument, op, err, "")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, bytes.NewReader(body))
	if err != nil {
		return notary.Errf(notary.KindInvalidArgument, op, err, "the bundler URL is not usable")
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return notary.Errf(notary.KindSourceUnavailable, op, nil, "the bundler could not be reached")
	}
	defer func() { _ = resp.Body.Close() }()

	raw, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseBytes))
	if err != nil {
		return notary.Errf(notary.KindSourceUnavailable, op, err, "the bundler's answer could not be read")
	}

	var answer response
	if err := json.Unmarshal(raw, &answer); err != nil {
		return notary.Errf(
			notary.KindSourceUnavailable,
			op,
			nil,
			"the bundler answered HTTP %d without JSON-RPC",
			resp.StatusCode,
		)
	}

	if answer.Error != nil {
		return refusal(op, answer.Error)
	}

	if err := json.Unmarshal(answer.Result, out); err != nil {
		return notary.Errf(notary.KindSourceRefused, op, err, "the bundler's answer does not have the expected shape")
	}

	return nil
}

func refusal(op string, failure *rpcError) error {
	switch failure.Code {
	case codeMethodNotFound, codeInternal, codeThrottled, codeLimitExceeded:
		return notary.Errf(
			notary.KindSourceUnavailable,
			op,
			nil,
			"the bundler could not take it now (code %d): %s",
			failure.Code,
			failure.Message,
		)
	default:
		return notary.Errf(
			notary.KindNotSponsored,
			op,
			nil,
			"the bundler refused it (code %d): %s",
			failure.Code,
			failure.Message,
		)
	}
}

func quantity(text string) (uint64, error) {
	if !strings.HasPrefix(text, "0x") || len(text) < 3 {
		return 0, fmt.Errorf("%q is not a hex quantity", text)
	}

	value, err := strconv.ParseUint(text[2:], 16, 64)
	if err != nil {
		return 0, fmt.Errorf("%q is not a hex quantity: %w", text, err)
	}

	return value, nil
}
