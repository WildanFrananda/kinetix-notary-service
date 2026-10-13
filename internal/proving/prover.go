package proving

import (
	"context"
	"sync"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

const cachedEpochs = 4

type Prover struct {
	source notary.ProofSource
	scheme notary.Scheme
	mu     sync.Mutex
	cache  map[uint64]proven
	order  []uint64
}

type proven struct {
	accounts []notary.Account
	tree     *notary.Tree
}

func NewProver(source notary.ProofSource, scheme notary.Scheme) *Prover {
	return &Prover{source: source, scheme: scheme, cache: map[uint64]proven{}}
}

func (p *Prover) Proofs(ctx context.Context, holder notary.Holder, ref EpochRef) (Bundle, error) {
	const op = "proving.Prover.Proofs"

	record, err := p.resolve(ctx, ref)
	if err != nil {
		return Bundle{}, err
	}

	switch record.State {
	case notary.StateFinalized:
	case notary.StateAbandoned:
		return Bundle{}, notary.Errf(
			notary.KindEpochAbandoned,
			op,
			nil,
			"epoch %d was abandoned and never anchored",
			record.Header.Number,
		)
	default:
		return Bundle{}, notary.Errf(
			notary.KindEpochNotFinal,
			op,
			nil,
			"epoch %d is %s, not final yet",
			record.Header.Number,
			record.State,
		)
	}

	built, err := p.load(ctx, record)
	if err != nil {
		return Bundle{}, err
	}

	claims := make([]notary.Claim, 0, 2)

	for _, account := range built.accounts {
		if account.Liability.Holder != holder {
			continue
		}

		proof, err := built.tree.Prove(account.Position)
		if err != nil {
			return Bundle{}, err
		}

		claims = append(claims, notary.Claim{
			Ledger: record.Ledger,
			Epoch:  record.Header.Number,
			Entry:  notary.Entry{Commitment: account.Commitment, Amount: account.Liability.Amount},
			Proof:  proof,
			Opening: &notary.Opening{
				Holder: account.Liability.Holder,
				Kind:   account.Liability.Kind,
				Salt:   account.Salt,
			},
		})
	}

	return Bundle{Epoch: record, Claims: claims}, nil
}

func (p *Prover) resolve(ctx context.Context, ref EpochRef) (notary.EpochRecord, error) {
	const op = "proving.Prover.resolve"

	if ref.latest {
		record, found, err := p.source.LastFinalized(ctx)
		if err != nil {
			return notary.EpochRecord{}, err
		}

		if !found {
			return notary.EpochRecord{}, notary.Errf(notary.KindEpochNotAnchored, op, nil, "no epoch is final yet")
		}

		return record, nil
	}

	record, found, err := p.source.Epoch(ctx, ref.number)
	if err != nil {
		return notary.EpochRecord{}, err
	}

	if !found {
		return notary.EpochRecord{}, notary.Errf(
			notary.KindEpochNotAnchored,
			op,
			nil,
			"epoch %d was never built",
			ref.number,
		)
	}

	return record, nil
}

func (p *Prover) load(ctx context.Context, record notary.EpochRecord) (proven, error) {
	const op = "proving.Prover.load"

	number := record.Header.Number

	p.mu.Lock()
	cached, ok := p.cache[number]
	p.mu.Unlock()

	if ok {
		return cached, nil
	}

	accounts, err := p.source.Accounts(ctx, number)
	if err != nil {
		return proven{}, err
	}

	if uint64(len(accounts)) != record.Header.LeafCount {
		return proven{}, notary.Errf(
			notary.KindDisagreement,
			op,
			nil,
			"epoch %d committed to %d accounts, the database holds %d",
			number,
			record.Header.LeafCount,
			len(accounts),
		)
	}

	entries := make([]notary.Entry, len(accounts))
	for i, account := range accounts {
		if account.Position != uint64(i) {
			return proven{}, notary.Errf(
				notary.KindDisagreement,
				op,
				nil,
				"epoch %d has no account at position %d",
				number,
				i,
			)
		}

		entries[i] = notary.Entry{Commitment: account.Commitment, Amount: account.Liability.Amount}
	}

	tree, err := notary.BuildTree(p.scheme, number, entries)
	if err != nil {
		return proven{}, err
	}

	if root := tree.Root(); root.Hash != record.Header.Root || root.Sum != record.Header.Total {
		return proven{}, notary.Errf(
			notary.KindDisagreement,
			op,
			nil,
			"the accounts stored for epoch %d do not rebuild the root it anchored",
			number,
		)
	}

	built := proven{accounts: accounts, tree: tree}

	p.mu.Lock()
	defer p.mu.Unlock()

	if _, ok := p.cache[number]; !ok {
		p.cache[number] = built
		p.order = append(p.order, number)

		if len(p.order) > cachedEpochs {
			delete(p.cache, p.order[0])
			p.order = p.order[1:]
		}
	}

	return built, nil
}
