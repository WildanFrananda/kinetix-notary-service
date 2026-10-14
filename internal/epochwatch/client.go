package epochwatch

import (
	"context"
	"errors"
	"io"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/WildanFrananda/kinetix-notary-service/internal/mesh"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	notaryv1 "github.com/WildanFrananda/kinetix-notary-service/internal/protogen/notary/v1"
)

type Client struct {
	conn     *grpc.ClientConn
	stub     notaryv1.EpochWatchServiceClient
	currency string
	timeout  time.Duration
}

var _ notary.EpochWatch = (*Client)(nil)

func Dial(settings Settings, currency string, options ...grpc.DialOption) (*Client, error) {
	const op = "epochwatch.Dial"

	creds, err := mesh.ClientCredentials(settings.Identity, Service)
	if err != nil {
		return nil, err
	}

	options = append([]grpc.DialOption{grpc.WithTransportCredentials(creds)}, options...)

	conn, err := grpc.NewClient(settings.Endpoint, options...)
	if err != nil {
		return nil, notary.Errf(notary.KindInvalidArgument, op, err, "watch endpoint %q", settings.Endpoint)
	}

	return &Client{
		conn:     conn,
		stub:     notaryv1.NewEpochWatchServiceClient(conn),
		currency: currency,
		timeout:  settings.Timeout,
	}, nil
}

func (c *Client) Close() error {
	if err := c.conn.Close(); err != nil {
		return notary.Errf(notary.KindUnknown, "epochwatch.Client.Close", err, "")
	}

	return nil
}

func (c *Client) Approve(ctx context.Context, proposal notary.EpochProposal) (notary.Approval, error) {
	const op = "epochwatch.Client.Approve"

	message, err := proposalMessage(op, proposal, c.currency)
	if err != nil {
		return notary.Approval{}, err
	}

	parts := []*notaryv1.ApproveEpochRequest{
		{Part: &notaryv1.ApproveEpochRequest_Proposal{Proposal: message}},
	}

	for start := 0; start < len(proposal.Accounts); start += batchSize {
		batch := &notaryv1.ProposedAccountBatch{}

		for _, account := range proposal.Accounts[start:min(start+batchSize, len(proposal.Accounts))] {
			wire, err := accountMessage(op, account, c.currency)
			if err != nil {
				return notary.Approval{}, err
			}

			batch.Accounts = append(batch.Accounts, wire)
		}

		parts = append(
			parts,
			&notaryv1.ApproveEpochRequest{Part: &notaryv1.ApproveEpochRequest_Accounts{Accounts: batch}},
		)
	}

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	stream, err := c.stub.ApproveEpoch(ctx)
	if err != nil {
		return notary.Approval{}, classify(op, err)
	}

	for _, part := range parts {
		if err := stream.Send(part); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return notary.Approval{}, classify(op, err)
		}
	}

	answer, err := stream.CloseAndRecv()
	if err != nil {
		return notary.Approval{}, classify(op, err)
	}

	return approvalOf(op, answer)
}

func approvalOf(op string, answer *notaryv1.ApproveEpochResponse) (notary.Approval, error) {
	var approval notary.Approval
	if len(answer.GetMember()) != len(approval.Member) || len(answer.GetSignature()) != len(approval.Signature) {
		return notary.Approval{}, notary.Errf(
			notary.KindDisagreement,
			op,
			nil,
			"the watch answered with a %d-byte member and a %d-byte signature",
			len(answer.GetMember()),
			len(answer.GetSignature()),
		)
	}

	copy(approval.Member[:], answer.GetMember())
	copy(approval.Signature[:], answer.GetSignature())

	return approval, nil
}

func classify(op string, err error) error {
	switch status.Code(err) {
	case codes.FailedPrecondition, codes.InvalidArgument:
		return notary.Errf(
			notary.KindApprovalRefused,
			op,
			err,
			"the watch refused the epoch: %s",
			status.Convert(err).Message(),
		)
	default:
		return notary.Errf(notary.KindSourceUnavailable, op, err, "the watch could not be asked")
	}
}
