package proverproc

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"os/exec"
	"strings"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

const (
	Binary        = "notary-prove"
	maxProofBytes = 1 << 20
)

type Client struct {
	env []string
}

var _ notary.EpochProver = (*Client)(nil)

func NewClient(env []string) *Client {
	return &Client{env: env}
}

func (c *Client) Prove(ctx context.Context, request notary.ProofRequest) ([]byte, error) {
	const op = "proverproc.Client.Prove"

	var input bytes.Buffer
	if err := EncodeRequest(&input, request); err != nil {
		return nil, err
	}

	var output, failure bytes.Buffer

	command := exec.CommandContext(ctx, Binary)
	command.Env = c.env
	command.Stdin = &input
	command.Stdout = &limited{buffer: &output, limit: maxProofBytes}
	command.Stderr = &limited{buffer: &failure, limit: maxProofBytes}

	err := command.Run()

	reason := strings.TrimSpace(failure.String())

	var exit *exec.ExitError
	switch {
	case err == nil:
	case errors.As(err, &exit) && exit.ExitCode() == ExitUnprovable:
		return nil, notary.Errf(
			notary.KindInvalidArgument,
			op,
			nil,
			"epoch %d cannot be proven: %s",
			request.Number,
			reason,
		)
	case errors.As(err, &exit) && exit.ExitCode() == ExitMisconfigured:
		return nil, notary.Errf(notary.KindInvalidArgument, op, nil, "the prover is misconfigured: %s", reason)
	default:
		return nil, notary.Errf(
			notary.KindUnknown,
			op,
			err,
			"the prover of epoch %d failed: %s",
			request.Number,
			reason,
		)
	}

	proof, err := hex.DecodeString(strings.TrimPrefix(strings.TrimSpace(output.String()), "0x"))
	if err != nil || len(proof) == 0 {
		return nil, notary.Errf(notary.KindUnknown, op, err, "the prover of epoch %d wrote no proof", request.Number)
	}

	return proof, nil
}
