package epochwatch

import (
	"errors"
	"io"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	notaryv1 "github.com/WildanFrananda/kinetix-notary-service/internal/protogen/notary/v1"
)

const Service = "notary-watch"

type Server struct {
	notaryv1.UnimplementedEpochWatchServiceServer
	watch notary.EpochWatch
}

func NewServer(watch notary.EpochWatch) *Server {
	return &Server{watch: watch}
}

func Register(server *grpc.Server, watch notary.EpochWatch) {
	notaryv1.RegisterEpochWatchServiceServer(server, NewServer(watch))
}

func (s *Server) ApproveEpoch(
	stream grpc.ClientStreamingServer[notaryv1.ApproveEpochRequest, notaryv1.ApproveEpochResponse],
) error {
	const op = "epochwatch.Server.ApproveEpoch"

	proposal, err := receive(op, stream)
	if err != nil {
		return toStatus(err)
	}

	approval, err := s.watch.Approve(stream.Context(), proposal)
	if err != nil {
		return toStatus(err)
	}

	if err := stream.SendAndClose(&notaryv1.ApproveEpochResponse{
		Member:    approval.Member[:],
		Signature: approval.Signature[:],
	}); err != nil {
		return notary.Errf(notary.KindUnknown, op, err, "the approval could not be sent back")
	}

	return nil
}

func receive(
	op string,
	stream grpc.ClientStreamingServer[notaryv1.ApproveEpochRequest, notaryv1.ApproveEpochResponse],
) (notary.EpochProposal, error) {
	first, err := stream.Recv()
	if errors.Is(err, io.EOF) {
		return notary.EpochProposal{}, malformed(op, "the stream carried no proposal")
	}

	if err != nil {
		return notary.EpochProposal{}, notary.Errf(notary.KindSourceUnavailable, op, err, "the notary's stream broke")
	}

	message := first.GetProposal()
	if message == nil {
		return notary.EpochProposal{}, malformed(op, "the first part is not a proposal")
	}

	id, header, err := headerOf(op, message)
	if err != nil {
		return notary.EpochProposal{}, err
	}

	confidential, err := confidentialOf(op, message.GetConfidential())
	if err != nil {
		return notary.EpochProposal{}, err
	}

	var accounts []notary.ProposedAccount

	for {
		part, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return notary.EpochProposal{
				SnapshotID:   id,
				Header:       header,
				Accounts:     accounts,
				Confidential: confidential,
			}, nil
		}

		if err != nil {
			return notary.EpochProposal{}, notary.Errf(
				notary.KindSourceUnavailable,
				op,
				err,
				"the notary's stream broke",
			)
		}

		batch := part.GetAccounts()
		if batch == nil {
			return notary.EpochProposal{}, malformed(op, "a part after the proposal is not a batch of accounts")
		}

		for _, message := range batch.GetAccounts() {
			account, err := accountOf(op, message)
			if err != nil {
				return notary.EpochProposal{}, err
			}

			accounts = append(accounts, account)
			if uint64(len(accounts)) > notary.MaxLeaves {
				return notary.EpochProposal{}, malformed(
					op,
					"the proposal carries more than %d accounts",
					notary.MaxLeaves,
				)
			}
		}
	}
}

func toStatus(err error) error {
	if _, isStatus := status.FromError(err); isStatus {
		return err
	}

	kind := notary.KindOf(err)

	code := codes.Internal

	switch kind {
	case notary.KindInvalidArgument:
		code = codes.InvalidArgument
	case notary.KindApprovalRefused, notary.KindSnapshotUnknown, notary.KindDisagreement:
		code = codes.FailedPrecondition
	case notary.KindSourceUnavailable, notary.KindSourceRefused, notary.KindLedgerUnavailable:
		code = codes.Unavailable
	}

	message := notary.MessageOf(err)
	if message == "" {
		message = err.Error()
	}

	return status.Error(code, kind.String()+": "+message)
}
