package epochwatch

import (
	"math"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	commonv1 "github.com/WildanFrananda/kinetix-notary-service/internal/protogen/common/v1"
	notaryv1 "github.com/WildanFrananda/kinetix-notary-service/internal/protogen/notary/v1"
	paymentv1 "github.com/WildanFrananda/kinetix-notary-service/internal/protogen/payment/v1"
)

const batchSize = 500

func proposalMessage(op string, proposal notary.EpochProposal, currency string) (*notaryv1.EpochProposal, error) {
	header := proposal.Header

	total, ok := signed(header.Total)
	if !ok {
		return nil, malformed(op, "the total of epoch %d does not fit int64", header.Number)
	}

	seconds, ok := signed(header.AsOf.Seconds())
	if !ok {
		return nil, malformed(op, "the snapshot time of epoch %d does not fit int64", header.Number)
	}

	return &notaryv1.EpochProposal{
		SnapshotId:     string(proposal.SnapshotID),
		Number:         header.Number,
		PreviousDigest: header.PreviousDigest[:],
		Root:           header.Root[:],
		Total:          &commonv1.Money{AmountMinor: total, Currency: currency},
		LeafCount:      header.LeafCount,
		AsOf:           &timestamppb.Timestamp{Seconds: seconds},
	}, nil
}

func accountMessage(op string, account notary.ProposedAccount, currency string) (*notaryv1.ProposedAccount, error) {
	amount, ok := signed(account.Liability.Amount)
	if !ok {
		return nil, malformed(op, "a %s balance does not fit int64", account.Liability.Kind)
	}

	return &notaryv1.ProposedAccount{
		Kind:    paymentv1.LiabilityKind(account.Liability.Kind),
		Holder:  account.Liability.Holder.String(),
		Balance: &commonv1.Money{AmountMinor: amount, Currency: currency},
		Secret:  account.Secret[:],
		Salt:    account.Salt[:],
	}, nil
}

func signed(value uint64) (int64, bool) {
	if value > math.MaxInt64 {
		return 0, false
	}

	return int64(value), true
}

func headerOf(op string, message *notaryv1.EpochProposal) (notary.SnapshotID, notary.EpochHeader, error) {
	id, err := notary.NewSnapshotID(message.GetSnapshotId())
	if err != nil {
		return "", notary.EpochHeader{}, malformed(op, "the proposal names no usable snapshot")
	}

	previous, ok := digestOf(message.GetPreviousDigest())
	if !ok {
		return "", notary.EpochHeader{}, malformed(op, "the previous digest is not 32 bytes")
	}

	root, ok := digestOf(message.GetRoot())
	if !ok {
		return "", notary.EpochHeader{}, malformed(op, "the root is not 32 bytes")
	}

	total := message.GetTotal().GetAmountMinor()
	if total < 0 {
		return "", notary.EpochHeader{}, malformed(op, "the total is negative")
	}

	stamp := message.GetAsOf()
	if stamp == nil || !stamp.IsValid() || stamp.GetNanos() != 0 {
		return "", notary.EpochHeader{}, malformed(op, "the snapshot time is not a whole second")
	}

	seconds := stamp.GetSeconds()
	if seconds <= 0 {
		return "", notary.EpochHeader{}, malformed(op, "the snapshot time is not after the Unix epoch")
	}

	asOf, err := notary.UnixTimeFromSeconds(uint64(seconds))
	if err != nil {
		return "", notary.EpochHeader{}, malformed(op, "the snapshot time is unusable")
	}

	return id, notary.EpochHeader{
		Number:         message.GetNumber(),
		PreviousDigest: previous,
		Root:           root,
		Total:          uint64(total),
		LeafCount:      message.GetLeafCount(),
		AsOf:           asOf,
	}, nil
}

func accountOf(op string, message *notaryv1.ProposedAccount) (notary.ProposedAccount, error) {
	code := int32(message.GetKind())
	if code < 1 || code > 255 {
		return notary.ProposedAccount{}, malformed(op, "an account has no kind")
	}

	kind, err := notary.AccountKindOf(uint8(code))
	if err != nil {
		return notary.ProposedAccount{}, malformed(op, "an account has kind %d", code)
	}

	holder, err := notary.NewHolder(message.GetHolder())
	if err != nil {
		return notary.ProposedAccount{}, malformed(op, "a %s account has an unusable holder", kind)
	}

	amount := message.GetBalance().GetAmountMinor()
	if amount < 0 {
		return notary.ProposedAccount{}, malformed(op, "a %s balance is negative", kind)
	}

	secret, ok := digestOf(message.GetSecret())
	if !ok {
		return notary.ProposedAccount{}, malformed(op, "a %s account's secret is not 32 bytes", kind)
	}

	salt, ok := digestOf(message.GetSalt())
	if !ok {
		return notary.ProposedAccount{}, malformed(op, "a %s account's salt is not 32 bytes", kind)
	}

	return notary.ProposedAccount{
		Liability: notary.Liability{Holder: holder, Kind: kind, Amount: uint64(amount)},
		Secret:    notary.AccountSecret(secret),
		Salt:      notary.Salt(salt),
	}, nil
}

func digestOf(raw []byte) (notary.Digest, bool) {
	var digest notary.Digest
	if len(raw) != len(digest) {
		return digest, false
	}

	copy(digest[:], raw)

	return digest, true
}

func malformed(op, format string, args ...any) error {
	return notary.Errf(notary.KindInvalidArgument, op, nil, format, args...)
}
