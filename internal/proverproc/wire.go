package proverproc

import (
	"encoding/json"
	"io"
	"strconv"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

const maxRequestBytes = 64 << 20

func EncodeRequest(w io.Writer, request notary.ProofRequest) error {
	wire := wireRequest{
		Number:          strconv.FormatUint(request.Number, 10),
		Depth:           request.Depth,
		TreeRoot:        request.TreeRoot.Hex(),
		TotalCommitment: request.TotalCommitment.Hex(),
		Blinding:        request.Blinding.Hex(),
		Entries:         make([]wireEntry, len(request.Entries)),
	}
	for i, entry := range request.Entries {
		wire.Entries[i] = wireEntry{Commitment: entry.Commitment.Hex(), Amount: strconv.FormatUint(entry.Amount, 10)}
	}

	if err := json.NewEncoder(w).Encode(wire); err != nil {
		return notary.Errf(notary.KindInvalidArgument, "proverproc.EncodeRequest", err, "epoch %d", request.Number)
	}

	return nil
}

func DecodeRequest(r io.Reader) (notary.ProofRequest, error) {
	const op = "proverproc.DecodeRequest"

	var wire wireRequest

	decoder := json.NewDecoder(io.LimitReader(r, maxRequestBytes))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&wire); err != nil {
		return notary.ProofRequest{}, notary.Errf(notary.KindInvalidArgument, op, err, "the request is not JSON")
	}

	number, err := strconv.ParseUint(wire.Number, 10, 64)
	if err != nil {
		return notary.ProofRequest{}, notary.Errf(notary.KindInvalidArgument, op, err, "the epoch number")
	}

	treeRoot, err := notary.ParseDigest(wire.TreeRoot)
	if err != nil {
		return notary.ProofRequest{}, notary.Errf(notary.KindInvalidArgument, op, err, "the tree root")
	}

	commitment, err := notary.ParseDigest(wire.TotalCommitment)
	if err != nil {
		return notary.ProofRequest{}, notary.Errf(notary.KindInvalidArgument, op, err, "the total commitment")
	}

	blinding, err := notary.ParseBlinding(wire.Blinding)
	if err != nil {
		return notary.ProofRequest{}, notary.Errf(notary.KindInvalidArgument, op, err, "the blinding")
	}

	entries := make([]notary.Entry, len(wire.Entries))
	for i, entry := range wire.Entries {
		commitment, err := notary.ParseDigest(entry.Commitment)
		if err != nil {
			return notary.ProofRequest{}, notary.Errf(notary.KindInvalidArgument, op, err, "entry %d's commitment", i)
		}

		amount, err := strconv.ParseUint(entry.Amount, 10, 64)
		if err != nil {
			return notary.ProofRequest{}, notary.Errf(notary.KindInvalidArgument, op, err, "entry %d's amount", i)
		}

		entries[i] = notary.Entry{Commitment: commitment, Amount: amount}
	}

	return notary.ProofRequest{
		Number:          number,
		Depth:           wire.Depth,
		Entries:         entries,
		TreeRoot:        treeRoot,
		TotalCommitment: commitment,
		Blinding:        blinding,
	}, nil
}
