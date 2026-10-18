package evidencing

import (
	"context"
	"slices"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type Prover struct {
	scheme notary.EvidenceScheme
	source notary.EvidenceProofSource
}

func NewProver(scheme notary.EvidenceScheme, source notary.EvidenceProofSource) *Prover {
	return &Prover{scheme: scheme, source: source}
}

func (p *Prover) OrderProofs(ctx context.Context, principal, orderNumber string) (Bundle, error) {
	const op = "evidencing.Prover.OrderProofs"

	if principal == "" || orderNumber == "" {
		return Bundle{}, notary.Errf(notary.KindInvalidArgument, op, nil, "a proof names its caller and its order")
	}

	stored, err := p.source.OrderEvidence(ctx, orderNumber)
	if err != nil {
		return Bundle{}, err
	}

	orderParties := map[string]bool{}

	for _, s := range stored {
		if record, ok := s.Entry.Record.(notary.LifecycleRecord); ok {
			orderParties[record.BuyerPrincipalID] = true
			orderParties[record.MerchantPrincipalID] = true
		}
	}

	visible := make([]notary.StoredEvidence, 0, len(stored))

	for _, s := range stored {
		switch record := s.Entry.Record.(type) {
		case notary.LifecycleRecord:
			if slices.Contains(record.Parties(), principal) {
				visible = append(visible, s)
			}
		case notary.DeliveryRecord:
			if record.DriverPrincipalID == principal || orderParties[principal] {
				visible = append(visible, s)
			}
		}
	}

	return p.prove(ctx, op, visible)
}

func (p *Prover) ReviewProofs(ctx context.Context, principal, reviewID string) (Bundle, error) {
	const op = "evidencing.Prover.ReviewProofs"

	if principal == "" || reviewID == "" {
		return Bundle{}, notary.Errf(notary.KindInvalidArgument, op, nil, "a proof names its caller and its review")
	}

	stored, err := p.source.ReviewEvidence(ctx, reviewID)
	if err != nil {
		return Bundle{}, err
	}

	visible := make([]notary.StoredEvidence, 0, len(stored))

	for _, s := range stored {
		if slices.Contains(s.Entry.Record.Parties(), principal) {
			visible = append(visible, s)
		}
	}

	return p.prove(ctx, op, visible)
}

func (p *Prover) prove(ctx context.Context, op string, visible []notary.StoredEvidence) (Bundle, error) {
	if len(visible) == 0 {
		return Bundle{}, notary.Errf(
			notary.KindNoEvidence,
			op,
			nil,
			"the notary keeps no evidence this caller is a party to",
		)
	}

	batches := map[batchKey]*provenBatch{}
	bundle := Bundle{Claims: make([]notary.EvidenceClaim, 0, len(visible))}

	for _, s := range visible {
		if err := s.Entry.Holds(p.scheme); err != nil {
			return Bundle{}, err
		}

		if !s.Batched {
			bundle.Pending++
			continue
		}

		key := batchKey{stream: s.Entry.Record.Stream(), number: s.Batch}

		batch, seen := batches[key]
		if !seen {
			proven, final, err := p.batch(ctx, op, key)
			if err != nil {
				return Bundle{}, err
			}

			if final {
				batch = &proven
			}

			batches[key] = batch
		}

		if batch == nil {
			bundle.Pending++
			continue
		}

		claim, err := p.claim(op, batch, s)
		if err != nil {
			return Bundle{}, err
		}

		bundle.Claims = append(bundle.Claims, claim)
	}

	return bundle, nil
}

func (p *Prover) batch(ctx context.Context, op string, key batchKey) (provenBatch, bool, error) {
	record, found, err := p.source.EvidenceBatch(ctx, key.stream, key.number)
	if err != nil {
		return provenBatch{}, false, err
	}

	if !found {
		return provenBatch{}, false, notary.Errf(
			notary.KindDisagreement,
			op,
			nil,
			"a %s record names batch %d, which the database does not hold",
			key.stream,
			key.number,
		)
	}

	if record.State != notary.StateFinalized {
		return provenBatch{}, false, nil
	}

	leaves, err := p.source.EvidenceBatchLeaves(ctx, key.stream, key.number)
	if err != nil {
		return provenBatch{}, false, err
	}

	tree, err := notary.BuildEvidenceTree(p.scheme, key.stream, leaves)
	if err != nil {
		return provenBatch{}, false, err
	}

	if tree.LeafCount() != record.Header.LeafCount || tree.Root() != record.Header.Root {
		return provenBatch{}, false, notary.Errf(
			notary.KindDisagreement,
			op,
			nil,
			"%s batch %d was anchored with %d records and root %s, the database rebuilds %d and %s",
			key.stream,
			key.number,
			record.Header.LeafCount,
			record.Header.Root.Hex(),
			tree.LeafCount(),
			tree.Root().Hex(),
		)
	}

	return provenBatch{record: record, tree: tree, leaves: leaves}, true, nil
}

func (p *Prover) claim(op string, batch *provenBatch, s notary.StoredEvidence) (notary.EvidenceClaim, error) {
	if s.Index >= uint64(len(batch.leaves)) || batch.leaves[s.Index] != s.Entry.Leaf {
		return notary.EvidenceClaim{}, notary.Errf(
			notary.KindDisagreement,
			op,
			nil,
			"the %s record at %d/%d is not leaf %d of batch %d",
			s.Entry.Record.Stream(),
			s.Entry.Record.Position().Transaction,
			s.Entry.Record.Position().Sequence,
			s.Index,
			s.Batch,
		)
	}

	proof, err := batch.tree.Prove(s.Index)
	if err != nil {
		return notary.EvidenceClaim{}, err
	}

	return notary.EvidenceClaim{
		Registry: batch.record.Registry,
		Scheme:   batch.record.Scheme,
		Header:   batch.record.Header,
		Digest:   batch.record.Digest,
		Record:   s.Entry.Record,
		Salt:     s.Entry.Salt,
		Proof:    proof,
	}, nil
}
