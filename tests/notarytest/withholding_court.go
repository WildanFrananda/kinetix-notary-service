package notarytest

import (
	"context"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type RequestKey struct {
	Registration notary.Digest
	Epoch        uint64
}

type RequestRecord struct {
	State notary.RequestState
	DueAt notary.UnixTime
}

type WithholdingCourt struct {
	IdentityValue notary.CourtIdentity
	IdentityErr   error
	Scheme        notary.WithholdingScheme
	Disagree      bool
	Latest        uint64
	Logged        map[uint64][]notary.WithholdingRequest
	States        map[RequestKey]RequestRecord
	Silences      uint64
	Err           error
	Scanned       [][2]uint64
}

var _ notary.WithholdingCourt = (*WithholdingCourt)(nil)

func NewWithholdingCourt(identity notary.CourtIdentity, scheme notary.WithholdingScheme) *WithholdingCourt {
	return &WithholdingCourt{
		IdentityValue: identity,
		Scheme:        scheme,
		Logged:        map[uint64][]notary.WithholdingRequest{},
		States:        map[RequestKey]RequestRecord{},
	}
}

func (w *WithholdingCourt) Identity(context.Context) (notary.CourtIdentity, error) {
	return w.IdentityValue, w.IdentityErr
}

func (w *WithholdingCourt) RegistrationDigest(
	_ context.Context,
	registration notary.Registration,
) (notary.Digest, error) {
	digest := w.Scheme.RegistrationDigest(w.IdentityValue, registration)
	if w.Disagree {
		digest[0] ^= 0xff
	}

	return digest, nil
}

func (w *WithholdingCourt) LatestBlock(context.Context) (uint64, error) {
	return w.Latest, w.Err
}

func (w *WithholdingCourt) Requests(_ context.Context, from, to uint64) ([]notary.WithholdingRequest, error) {
	if w.Err != nil {
		return nil, w.Err
	}

	w.Scanned = append(w.Scanned, [2]uint64{from, to})

	var out []notary.WithholdingRequest

	for block := from; block <= to; block++ {
		out = append(out, w.Logged[block]...)
	}

	return out, nil
}

func (w *WithholdingCourt) RequestOf(
	_ context.Context,
	registration notary.Digest,
	epoch uint64,
) (notary.RequestState, notary.UnixTime, error) {
	record := w.States[RequestKey{Registration: registration, Epoch: epoch}]

	return record.State, record.DueAt, w.Err
}

func (w *WithholdingCourt) TotalSilences(context.Context) (uint64, error) {
	return w.Silences, w.Err
}
