package proving

import (
	"context"
	"math"
	"sync"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

const cachedEpochs = 4

type Prover struct {
	source     notary.ProofSource
	scheme     notary.Scheme
	statements notary.StatementScheme
	court      notary.Court
	signer     notary.StatementSigner
	identity   notary.CourtIdentity
	mu         sync.Mutex
	cache      map[uint64]proven
	order      []uint64
}

type proven struct {
	accounts []notary.Account
	tree     *notary.Tree
	signer   notary.Address
}

func NewProver(
	source notary.ProofSource,
	scheme notary.Scheme,
	statements notary.StatementScheme,
	court notary.Court,
	signer notary.StatementSigner,
	identity notary.CourtIdentity,
) *Prover {
	return &Prover{
		source:     source,
		scheme:     scheme,
		statements: statements,
		court:      court,
		signer:     signer,
		identity:   identity,
		cache:      map[uint64]proven{},
	}
}

func (p *Prover) CheckCourt(ctx context.Context) error {
	const op = "proving.Prover.CheckCourt"

	identity, err := p.court.Identity(ctx)
	if err != nil {
		return err
	}

	if identity != p.identity {
		return notary.Errf(
			notary.KindWrongCourt,
			op,
			nil,
			"the court answers as %s on chain %d, this notary signs for %s on chain %d",
			identity.Address.Hex(),
			identity.ChainID,
			p.identity.Address.Hex(),
			p.identity.ChainID,
		)
	}

	sample := notary.Statement{
		Epoch:        1,
		Commitment:   notary.Digest{1},
		Amount:       2,
		Index:        3,
		SiblingsHash: notary.Digest{4},
	}

	onChain, err := p.court.StatementDigest(ctx, sample)
	if err != nil {
		return err
	}

	if local := p.statements.StatementDigest(p.identity, sample); local != onChain {
		return notary.Errf(
			notary.KindDisagreement,
			op,
			nil,
			"the court hashes a statement to %s, this notary to %s",
			onChain.Hex(),
			local.Hex(),
		)
	}

	newest, err := p.court.SignerOf(ctx, math.MaxUint64)
	if err != nil {
		return err
	}

	if !p.signer.Holds(newest) {
		return notary.Errf(
			notary.KindSignerUnavailable,
			op,
			nil,
			"the court's newest signer is %s and this notary holds no key for it",
			newest.Hex(),
		)
	}

	return nil
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

		claim := notary.Claim{
			Ledger: record.Ledger,
			Epoch:  record.Header.Number,
			Entry:  notary.Entry{Commitment: account.Commitment, Amount: account.Liability.Amount},
			Proof:  proof,
			Opening: &notary.Opening{
				Holder: account.Liability.Holder,
				Kind:   account.Liability.Kind,
				Salt:   account.Salt,
				Secret: account.Secret,
			},
		}

		statement := notary.StatementOf(p.statements, claim)

		signature, err := p.signer.Sign(built.signer, p.statements.StatementDigest(p.identity, statement))
		if err != nil {
			return Bundle{}, err
		}

		claim.Statement = &notary.StatementSignature{Court: p.identity.Address, Signature: signature}
		claims = append(claims, claim)
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

		opening := notary.Opening{
			Holder: account.Liability.Holder,
			Kind:   account.Liability.Kind,
			Salt:   account.Salt,
			Secret: account.Secret,
		}
		if opening.Commitment(p.scheme) != account.Commitment {
			return proven{}, notary.Errf(
				notary.KindDisagreement,
				op,
				nil,
				"the account at position %d of epoch %d does not open its commitment",
				i,
				number,
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

	signer, err := p.court.SignerOf(ctx, number)
	if err != nil {
		return proven{}, err
	}

	if signer.IsZero() {
		return proven{}, notary.Errf(
			notary.KindSignerUnavailable,
			op,
			nil,
			"the court names no signer for epoch %d",
			number,
		)
	}

	built := proven{accounts: accounts, tree: tree, signer: signer}

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
