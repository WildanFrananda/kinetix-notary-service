package prooffile

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"time"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

const EvidenceFormatV1 = "kinetix.notary.evidence.v1"

func IsEvidence(raw []byte) bool {
	var probe struct {
		Format string            `json:"format"`
		Proofs []json.RawMessage `json:"proofs"`
	}

	if err := json.Unmarshal(raw, &probe); err != nil {
		return false
	}

	if probe.Format == EvidenceFormatV1 {
		return true
	}

	if len(probe.Proofs) == 0 {
		return false
	}

	var first struct {
		Format string `json:"format"`
	}

	return json.Unmarshal(probe.Proofs[0], &first) == nil && first.Format == EvidenceFormatV1
}

func MarshalEvidence(claim notary.EvidenceClaim) (json.RawMessage, error) {
	doc, err := evidenceDocumentOf(claim)
	if err != nil {
		return nil, err
	}

	raw, err := json.Marshal(doc)
	if err != nil {
		return nil, fmt.Errorf("encode evidence document: %w", err)
	}

	return raw, nil
}

func EncodeEvidence(w io.Writer, claim notary.EvidenceClaim) error {
	doc, err := evidenceDocumentOf(claim)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(doc); err != nil {
		return fmt.Errorf("encode evidence document: %w", err)
	}

	return nil
}

func DecodeEvidence(r io.Reader, scheme uint16) (notary.EvidenceClaim, error) {
	const op = "prooffile.DecodeEvidence"

	decoder := json.NewDecoder(io.LimitReader(r, maxDocumentBytes+1))
	decoder.DisallowUnknownFields()

	var doc evidenceDocument
	if err := decoder.Decode(&doc); err != nil {
		return notary.EvidenceClaim{}, notary.Errf(
			notary.KindInvalidArgument,
			op,
			err,
			"the proof is not a valid document",
		)
	}

	if decoder.More() || decoder.InputOffset() > maxDocumentBytes {
		return notary.EvidenceClaim{}, notary.Errf(
			notary.KindInvalidArgument,
			op,
			nil,
			"the proof carries more than one document",
		)
	}

	if _, err := decoder.Token(); !errors.Is(err, io.EOF) {
		return notary.EvidenceClaim{}, notary.Errf(
			notary.KindInvalidArgument,
			op,
			nil,
			"the proof has trailing content",
		)
	}

	if doc.Format != EvidenceFormatV1 {
		return notary.EvidenceClaim{}, notary.Errf(
			notary.KindInvalidArgument,
			op,
			nil,
			"format %q is not %q",
			doc.Format,
			EvidenceFormatV1,
		)
	}

	if doc.Scheme != scheme {
		return notary.EvidenceClaim{}, notary.Errf(
			notary.KindInvalidArgument,
			op,
			nil,
			"the proof uses evidence scheme %d, this verifier speaks scheme %d",
			doc.Scheme,
			scheme,
		)
	}

	if len(doc.Siblings) > notary.MaxDepth {
		return notary.EvidenceClaim{}, notary.Errf(
			notary.KindInvalidArgument,
			op,
			nil,
			"%d siblings exceed the scheme's depth of %d",
			len(doc.Siblings),
			notary.MaxDepth,
		)
	}

	claim, err := doc.claim()
	if err != nil {
		return notary.EvidenceClaim{}, notary.Errf(notary.KindInvalidArgument, op, err, "")
	}

	return claim, nil
}

func DecodeEvidenceAll(r io.Reader, scheme uint16) ([]notary.EvidenceClaim, error) {
	const op = "prooffile.DecodeEvidenceAll"

	raw, err := io.ReadAll(io.LimitReader(r, maxBundleBytes+1))
	if err != nil {
		return nil, notary.Errf(notary.KindInvalidArgument, op, err, "reading the proof")
	}

	if len(raw) > maxBundleBytes {
		return nil, notary.Errf(
			notary.KindInvalidArgument,
			op,
			nil,
			"the proof is larger than %d bytes",
			maxBundleBytes,
		)
	}

	var probe map[string]json.RawMessage
	if err := json.Unmarshal(raw, &probe); err != nil {
		return nil, notary.Errf(notary.KindInvalidArgument, op, err, "the proof is not a JSON object")
	}

	proofs, isBundle := probe["proofs"]
	if !isBundle {
		claim, err := DecodeEvidence(bytes.NewReader(raw), scheme)
		if err != nil {
			return nil, err
		}

		return []notary.EvidenceClaim{claim}, nil
	}

	var documents []json.RawMessage
	if err := json.Unmarshal(proofs, &documents); err != nil {
		return nil, notary.Errf(notary.KindInvalidArgument, op, err, "proofs is not a list")
	}

	if len(documents) == 0 {
		return nil, notary.Errf(notary.KindInvalidArgument, op, nil, "the bundle carries no proof")
	}

	claims := make([]notary.EvidenceClaim, 0, len(documents))
	for i, doc := range documents {
		claim, err := DecodeEvidence(bytes.NewReader(doc), scheme)
		if err != nil {
			return nil, notary.Errf(notary.KindInvalidArgument, op, err, "proof %d", i)
		}

		claims = append(claims, claim)
	}

	return claims, nil
}

func evidenceDocumentOf(claim notary.EvidenceClaim) (evidenceDocument, error) {
	header := claim.Header
	position := claim.Record.Position()

	doc := evidenceDocument{
		Format:   EvidenceFormatV1,
		Scheme:   claim.Scheme,
		ChainID:  strconv.FormatUint(claim.Registry.ChainID, 10),
		Registry: claim.Registry.Address.Hex(),
		Batch: evidenceBatch{
			Stream:             header.Stream.String(),
			Number:             strconv.FormatUint(header.Number, 10),
			PreviousDigest:     header.PreviousDigest.Hex(),
			Root:               header.Root.Hex(),
			LeafCount:          strconv.FormatUint(header.LeafCount, 10),
			ThroughTransaction: strconv.FormatUint(header.Through.Transaction, 10),
			ThroughSequence:    strconv.FormatUint(header.Through.Sequence, 10),
			ClosedAt:           strconv.FormatUint(header.ClosedAt.Seconds(), 10),
			Digest:             claim.Digest.Hex(),
		},
		Record: evidenceRecord{
			Transaction: strconv.FormatUint(position.Transaction, 10),
			Sequence:    strconv.FormatUint(position.Sequence, 10),
		},
		Salt:     claim.Salt.Hex(),
		Index:    strconv.FormatUint(claim.Proof.Index, 10),
		Siblings: make([]string, len(claim.Proof.Siblings)),
	}

	for i, sibling := range claim.Proof.Siblings {
		doc.Siblings[i] = sibling.Hex()
	}

	switch r := claim.Record.(type) {
	case notary.LifecycleRecord:
		doc.Record.Lifecycle = &lifecycleFields{
			OrderNumber:         r.OrderNumber,
			Kind:                r.Kind,
			FromStatus:          r.FromStatus,
			ToStatus:            r.ToStatus,
			ReturnNumber:        r.ReturnNumber,
			ReturnStatus:        r.ReturnStatus,
			ActorKind:           r.ActorKind,
			ActorPrincipalID:    r.ActorPrincipalID,
			BuyerPrincipalID:    r.BuyerPrincipalID,
			MerchantPrincipalID: r.MerchantPrincipalID,
			OccurredAt:          timeText(r.OccurredAt),
		}
	case notary.DeliveryRecord:
		doc.Record.Delivery = &deliveryFields{
			OrderNumber:       r.OrderNumber,
			DriverPrincipalID: r.DriverPrincipalID,
			PhotoSHA256:       hashText(r.PhotoSHA256),
			SignatureSHA256:   hashText(r.SignatureSHA256),
			DeliveredAt:       timeText(r.DeliveredAt),
		}
	case notary.ReviewRecord:
		doc.Record.Review = &reviewFields{
			ReviewID:            r.ReviewID,
			Kind:                r.Kind,
			OrderID:             r.OrderID,
			SubjectID:           r.SubjectID,
			AuthorPrincipalID:   r.AuthorPrincipalID,
			MerchantPrincipalID: r.MerchantPrincipalID,
			Rating:              r.Rating,
			CommentSHA256:       hashText(r.CommentSHA256),
			CreatedAt:           timeText(r.CreatedAt),
		}
	default:
		return evidenceDocument{}, notary.Errf(
			notary.KindInvalidArgument,
			"prooffile.evidenceDocumentOf",
			nil,
			"no document carries a %T",
			claim.Record,
		)
	}

	return doc, nil
}

func (d evidenceDocument) claim() (notary.EvidenceClaim, error) {
	var (
		claim notary.EvidenceClaim
		err   error
	)

	claim.Scheme = d.Scheme

	if claim.Registry.ChainID, err = number("chainId", d.ChainID); err != nil {
		return notary.EvidenceClaim{}, err
	}

	if claim.Registry.Address, err = notary.ParseAddress(d.Registry); err != nil {
		return notary.EvidenceClaim{}, fmt.Errorf("registry: %w", err)
	}

	if claim.Header, claim.Digest, err = d.Batch.header(); err != nil {
		return notary.EvidenceClaim{}, err
	}

	if claim.Record, err = d.Record.record(claim.Header.Stream); err != nil {
		return notary.EvidenceClaim{}, err
	}

	if claim.Salt, err = notary.ParseSalt(d.Salt); err != nil {
		return notary.EvidenceClaim{}, fmt.Errorf("salt: %w", err)
	}

	if claim.Salt.IsZero() {
		return notary.EvidenceClaim{}, errors.New("salt: a salt is never all zero")
	}

	if claim.Proof.Index, err = number("index", d.Index); err != nil {
		return notary.EvidenceClaim{}, err
	}

	claim.Proof.Siblings = make([]notary.Digest, len(d.Siblings))
	for i, sibling := range d.Siblings {
		if claim.Proof.Siblings[i], err = notary.ParseDigest(sibling); err != nil {
			return notary.EvidenceClaim{}, fmt.Errorf("sibling %d: %w", i, err)
		}
	}

	return claim, nil
}

func (b evidenceBatch) header() (notary.EvidenceBatchHeader, notary.Digest, error) {
	var (
		header notary.EvidenceBatchHeader
		err    error
	)

	if header.Stream, err = notary.EvidenceStreamNamed(b.Stream); err != nil {
		return notary.EvidenceBatchHeader{}, notary.Digest{}, fmt.Errorf("batch stream: %w", err)
	}

	if header.Number, err = number("batch number", b.Number); err != nil {
		return notary.EvidenceBatchHeader{}, notary.Digest{}, err
	}

	if header.LeafCount, err = number("batch leafCount", b.LeafCount); err != nil {
		return notary.EvidenceBatchHeader{}, notary.Digest{}, err
	}

	if header.Through.Transaction, err = number("batch throughTransaction", b.ThroughTransaction); err != nil {
		return notary.EvidenceBatchHeader{}, notary.Digest{}, err
	}

	if header.Through.Sequence, err = number("batch throughSequence", b.ThroughSequence); err != nil {
		return notary.EvidenceBatchHeader{}, notary.Digest{}, err
	}

	closedAt, err := number("batch closedAt", b.ClosedAt)
	if err != nil {
		return notary.EvidenceBatchHeader{}, notary.Digest{}, err
	}

	if header.ClosedAt, err = notary.UnixTimeFromSeconds(closedAt); err != nil {
		return notary.EvidenceBatchHeader{}, notary.Digest{}, fmt.Errorf("batch closedAt: %w", err)
	}

	digests := make([]notary.Digest, 3)

	for i, field := range []struct{ name, text string }{
		{"batch previousDigest", b.PreviousDigest},
		{"batch root", b.Root},
		{"batch digest", b.Digest},
	} {
		if digests[i], err = notary.ParseDigest(field.text); err != nil {
			return notary.EvidenceBatchHeader{}, notary.Digest{}, fmt.Errorf("%s: %w", field.name, err)
		}
	}

	header.PreviousDigest, header.Root = digests[0], digests[1]

	return header, digests[2], nil
}

func (r evidenceRecord) record(stream notary.EvidenceStream) (notary.EvidenceRecord, error) {
	var (
		at  notary.LogPosition
		err error
	)

	if at.Transaction, err = number("record transaction", r.Transaction); err != nil {
		return nil, err
	}

	if at.Sequence, err = number("record sequence", r.Sequence); err != nil {
		return nil, err
	}

	carried := 0

	for _, present := range []bool{r.Lifecycle != nil, r.Delivery != nil, r.Review != nil} {
		if present {
			carried++
		}
	}

	if carried != 1 {
		return nil, fmt.Errorf("the record carries %d sets of fields, a %s record carries exactly one", carried, stream)
	}

	switch {
	case r.Lifecycle != nil && stream == notary.StreamOrderLifecycle:
		return r.Lifecycle.record(at)
	case r.Delivery != nil && stream == notary.StreamDelivery:
		return r.Delivery.record(at)
	case r.Review != nil && stream == notary.StreamReview:
		return r.Review.record(at)
	default:
		return nil, fmt.Errorf("the record's fields do not belong to the %s stream", stream)
	}
}

func (f lifecycleFields) record(at notary.LogPosition) (notary.EvidenceRecord, error) {
	occurredAt, err := parseTime("occurredAt", f.OccurredAt)
	if err != nil {
		return nil, err
	}

	if err := fieldLengths(
		f.OrderNumber,
		f.ReturnNumber,
		f.ActorPrincipalID,
		f.BuyerPrincipalID,
		f.MerchantPrincipalID,
	); err != nil {
		return nil, err
	}

	return notary.LifecycleRecord{
		At:                  at,
		OrderNumber:         f.OrderNumber,
		Kind:                f.Kind,
		FromStatus:          f.FromStatus,
		ToStatus:            f.ToStatus,
		ReturnNumber:        f.ReturnNumber,
		ReturnStatus:        f.ReturnStatus,
		ActorKind:           f.ActorKind,
		ActorPrincipalID:    f.ActorPrincipalID,
		BuyerPrincipalID:    f.BuyerPrincipalID,
		MerchantPrincipalID: f.MerchantPrincipalID,
		OccurredAt:          occurredAt,
	}, nil
}

func (f deliveryFields) record(at notary.LogPosition) (notary.EvidenceRecord, error) {
	deliveredAt, err := parseTime("deliveredAt", f.DeliveredAt)
	if err != nil {
		return nil, err
	}

	photo, err := parseHash("photoSha256", f.PhotoSHA256, true)
	if err != nil {
		return nil, err
	}

	signature, err := parseHash("signatureSha256", f.SignatureSHA256, false)
	if err != nil {
		return nil, err
	}

	if err := fieldLengths(f.OrderNumber, f.DriverPrincipalID); err != nil {
		return nil, err
	}

	return notary.DeliveryRecord{
		At:                at,
		OrderNumber:       f.OrderNumber,
		DriverPrincipalID: f.DriverPrincipalID,
		PhotoSHA256:       photo,
		SignatureSHA256:   signature,
		DeliveredAt:       deliveredAt,
	}, nil
}

func (f reviewFields) record(at notary.LogPosition) (notary.EvidenceRecord, error) {
	createdAt, err := parseTime("createdAt", f.CreatedAt)
	if err != nil {
		return nil, err
	}

	comment, err := parseHash("commentSha256", f.CommentSHA256, false)
	if err != nil {
		return nil, err
	}

	if err := fieldLengths(
		f.ReviewID,
		f.OrderID,
		f.SubjectID,
		f.AuthorPrincipalID,
		f.MerchantPrincipalID,
	); err != nil {
		return nil, err
	}

	return notary.ReviewRecord{
		At:                  at,
		ReviewID:            f.ReviewID,
		Kind:                f.Kind,
		OrderID:             f.OrderID,
		SubjectID:           f.SubjectID,
		AuthorPrincipalID:   f.AuthorPrincipalID,
		MerchantPrincipalID: f.MerchantPrincipalID,
		Rating:              f.Rating,
		CommentSHA256:       comment,
		CreatedAt:           createdAt,
	}, nil
}

func timeText(t notary.EvidenceTime) string {
	return t.Time().Format(time.RFC3339Nano)
}

func parseTime(field, text string) (notary.EvidenceTime, error) {
	parsed, err := time.Parse(time.RFC3339Nano, text)
	if err != nil {
		return notary.EvidenceTime{}, fmt.Errorf("%s %q is not an RFC 3339 time", field, text)
	}

	nanos := parsed.Nanosecond()
	if nanos < 0 || nanos > math.MaxInt32 {
		return notary.EvidenceTime{}, fmt.Errorf("%s %q carries no valid nanoseconds", field, text)
	}

	t, err := notary.EvidenceTimeFrom(parsed.Unix(), int32(nanos))
	if err != nil {
		return notary.EvidenceTime{}, fmt.Errorf("%s: %w", field, err)
	}

	return t, nil
}

func hashText(raw []byte) string {
	if len(raw) == 0 {
		return ""
	}

	var d notary.Digest

	copy(d[:], raw)

	return d.Hex()
}

func parseHash(field, text string, required bool) ([]byte, error) {
	if text == "" {
		if required {
			return nil, fmt.Errorf("%s is missing", field)
		}

		return nil, nil
	}

	d, err := notary.ParseDigest(text)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", field, err)
	}

	return d[:], nil
}

func fieldLengths(values ...string) error {
	for _, value := range values {
		if len(value) > notary.MaxEvidenceField {
			return fmt.Errorf("a field of %d bytes is longer than the scheme's %d", len(value), notary.MaxEvidenceField)
		}
	}

	return nil
}
