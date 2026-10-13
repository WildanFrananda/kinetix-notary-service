package prooffile

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

const Format = "kinetix.notary.inclusion.v1"

const maxDocumentBytes = 64 << 10

const maxBundleBytes = 4 << 20

func Decode(r io.Reader, scheme uint16) (notary.Claim, error) {
	const op = "prooffile.Decode"

	decoder := json.NewDecoder(io.LimitReader(r, maxDocumentBytes+1))
	decoder.DisallowUnknownFields()

	var doc document
	if err := decoder.Decode(&doc); err != nil {
		return notary.Claim{}, notary.Errf(notary.KindInvalidArgument, op, err, "the proof is not a valid document")
	}

	if decoder.More() || decoder.InputOffset() > maxDocumentBytes {
		return notary.Claim{}, notary.Errf(
			notary.KindInvalidArgument,
			op,
			nil,
			"the proof carries more than one document",
		)
	}

	if _, err := decoder.Token(); !errors.Is(err, io.EOF) {
		return notary.Claim{}, notary.Errf(notary.KindInvalidArgument, op, nil, "the proof has trailing content")
	}

	if doc.Format != Format {
		return notary.Claim{}, notary.Errf(
			notary.KindInvalidArgument,
			op,
			nil,
			"format %q is not %q",
			doc.Format,
			Format,
		)
	}

	if doc.Scheme != scheme {
		return notary.Claim{}, notary.Errf(
			notary.KindInvalidArgument,
			op,
			nil,
			"the proof uses scheme %d, this verifier speaks scheme %d",
			doc.Scheme,
			scheme,
		)
	}

	if len(doc.Siblings) > notary.MaxDepth {
		return notary.Claim{}, notary.Errf(
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
		return notary.Claim{}, notary.Errf(notary.KindInvalidArgument, op, err, "")
	}

	return claim, nil
}

func Encode(w io.Writer, claim notary.Claim, scheme uint16) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(documentOf(claim, scheme)); err != nil {
		return fmt.Errorf("encode proof document: %w", err)
	}

	return nil
}

func Marshal(claim notary.Claim, scheme uint16) (json.RawMessage, error) {
	raw, err := json.Marshal(documentOf(claim, scheme))
	if err != nil {
		return nil, fmt.Errorf("encode proof document: %w", err)
	}

	return raw, nil
}

func DecodeAll(r io.Reader, scheme uint16) ([]notary.Claim, error) {
	const op = "prooffile.DecodeAll"

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
		claim, err := Decode(bytes.NewReader(raw), scheme)
		if err != nil {
			return nil, err
		}

		return []notary.Claim{claim}, nil
	}

	var documents []json.RawMessage
	if err := json.Unmarshal(proofs, &documents); err != nil {
		return nil, notary.Errf(notary.KindInvalidArgument, op, err, "proofs is not a list")
	}

	if len(documents) == 0 {
		return nil, notary.Errf(notary.KindInvalidArgument, op, nil, "the bundle carries no proof")
	}

	claims := make([]notary.Claim, 0, len(documents))
	for i, doc := range documents {
		claim, err := Decode(bytes.NewReader(doc), scheme)
		if err != nil {
			return nil, notary.Errf(notary.KindInvalidArgument, op, err, "proof %d", i)
		}

		claims = append(claims, claim)
	}

	return claims, nil
}

func documentOf(claim notary.Claim, scheme uint16) document {
	doc := document{
		Format:     Format,
		Scheme:     scheme,
		ChainID:    strconv.FormatUint(claim.Ledger.ChainID, 10),
		Ledger:     claim.Ledger.Address.Hex(),
		Epoch:      strconv.FormatUint(claim.Epoch, 10),
		Commitment: claim.Entry.Commitment.Hex(),
		Amount:     strconv.FormatUint(claim.Entry.Amount, 10),
		Index:      strconv.FormatUint(claim.Proof.Index, 10),
		Siblings:   make([]sibling, len(claim.Proof.Siblings)),
	}
	for i, node := range claim.Proof.Siblings {
		doc.Siblings[i] = sibling{Hash: node.Hash.Hex(), Sum: strconv.FormatUint(node.Sum, 10)}
	}

	if opening := claim.Opening; opening != nil {
		doc.Account = &account{
			Holder: opening.Holder.String(),
			Kind:   opening.Kind.String(),
			Salt:   opening.Salt.Hex(),
		}
	}

	return doc
}

func (d document) claim() (notary.Claim, error) {
	var (
		claim notary.Claim
		err   error
	)
	if claim.Ledger.ChainID, err = number("chainId", d.ChainID); err != nil {
		return notary.Claim{}, err
	}

	if claim.Ledger.Address, err = notary.ParseAddress(d.Ledger); err != nil {
		return notary.Claim{}, fmt.Errorf("ledger: %w", err)
	}

	if claim.Epoch, err = number("epoch", d.Epoch); err != nil {
		return notary.Claim{}, err
	}

	if claim.Entry.Commitment, err = notary.ParseDigest(d.Commitment); err != nil {
		return notary.Claim{}, fmt.Errorf("commitment: %w", err)
	}

	if claim.Entry.Amount, err = number("amount", d.Amount); err != nil {
		return notary.Claim{}, err
	}

	if claim.Proof.Index, err = number("index", d.Index); err != nil {
		return notary.Claim{}, err
	}

	if d.Account != nil {
		opening, err := d.Account.opening()
		if err != nil {
			return notary.Claim{}, err
		}

		claim.Opening = &opening
	}

	claim.Proof.Siblings = make([]notary.Node, len(d.Siblings))
	for i, s := range d.Siblings {
		if claim.Proof.Siblings[i].Hash, err = notary.ParseDigest(s.Hash); err != nil {
			return notary.Claim{}, fmt.Errorf("sibling %d hash: %w", i, err)
		}

		if claim.Proof.Siblings[i].Sum, err = number(fmt.Sprintf("sibling %d sum", i), s.Sum); err != nil {
			return notary.Claim{}, err
		}
	}

	return claim, nil
}

func number(field, text string) (uint64, error) {
	value, err := strconv.ParseUint(text, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%s %q is not an unsigned 64-bit decimal", field, text)
	}

	return value, nil
}

func (a account) opening() (notary.Opening, error) {
	holder, err := notary.NewHolder(a.Holder)
	if err != nil {
		return notary.Opening{}, fmt.Errorf("account holder: %w", err)
	}

	kind, err := notary.AccountKindNamed(a.Kind)
	if err != nil {
		return notary.Opening{}, fmt.Errorf("account kind: %w", err)
	}

	salt, err := notary.ParseSalt(a.Salt)
	if err != nil {
		return notary.Opening{}, fmt.Errorf("account salt: %w", err)
	}

	return notary.Opening{Holder: holder, Kind: kind, Salt: salt}, nil
}
