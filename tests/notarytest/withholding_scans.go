package notarytest

import (
	"context"
	"sort"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type WithholdingScans struct {
	Next map[notary.CourtIdentity]uint64
	Open map[notary.CourtIdentity]map[RequestKey]notary.WithholdingRequest
	Err  error
}

var _ notary.WithholdingScans = (*WithholdingScans)(nil)

func NewWithholdingScans() *WithholdingScans {
	return &WithholdingScans{
		Next: map[notary.CourtIdentity]uint64{},
		Open: map[notary.CourtIdentity]map[RequestKey]notary.WithholdingRequest{},
	}
}

func (w *WithholdingScans) ScanFrom(_ context.Context, court notary.CourtIdentity) (uint64, bool, error) {
	if w.Err != nil {
		return 0, false, w.Err
	}

	next, found := w.Next[court]

	return next, found, nil
}

func (w *WithholdingScans) OpenRequests(
	_ context.Context,
	court notary.CourtIdentity,
) ([]notary.WithholdingRequest, error) {
	if w.Err != nil {
		return nil, w.Err
	}

	open := make([]notary.WithholdingRequest, 0, len(w.Open[court]))
	for _, request := range w.Open[court] {
		open = append(open, request)
	}

	sort.Slice(open, func(i, j int) bool {
		return open[i].DueAt.Seconds() < open[j].DueAt.Seconds()
	})

	return open, nil
}

func (w *WithholdingScans) RecordScan(
	_ context.Context,
	court notary.CourtIdentity,
	next uint64,
	found []notary.WithholdingRequest,
) error {
	if w.Err != nil {
		return w.Err
	}

	w.Next[court] = max(w.Next[court], next)

	if w.Open[court] == nil {
		w.Open[court] = map[RequestKey]notary.WithholdingRequest{}
	}

	for _, request := range found {
		key := RequestKey{Registration: request.Registration, Epoch: request.Epoch}
		if _, known := w.Open[court][key]; !known {
			w.Open[court][key] = request
		}
	}

	return nil
}

func (w *WithholdingScans) CloseRequest(
	_ context.Context,
	court notary.CourtIdentity,
	registration notary.Digest,
	epoch uint64,
) error {
	if w.Err != nil {
		return w.Err
	}

	delete(w.Open[court], RequestKey{Registration: registration, Epoch: epoch})

	return nil
}
