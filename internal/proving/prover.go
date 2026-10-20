package proving

import (
	"context"
	"math"
	"sync"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

const cachedEpochs = 4

type Prover struct {
	source       notary.ProofSource
	scheme       notary.Scheme
	confidential notary.ConfidentialScheme
	statements   notary.StatementScheme
	signer       notary.StatementSigner
	current      Venue
	earlier      map[notary.LedgerIdentity]Venue
	mu           sync.Mutex
	cache        map[uint64]proven
	order        []uint64
}

type proven struct {
	accounts []notary.Account
	tree     notary.EpochTree
	signer   notary.Address
}

func NewProver(
	source notary.ProofSource,
	scheme notary.Scheme,
	confidential notary.ConfidentialScheme,
	statements notary.StatementScheme,
	signer notary.StatementSigner,
	current Venue,
	earlier []Venue,
) *Prover {
	byLedger := make(map[notary.LedgerIdentity]Venue, len(earlier))
	for _, venue := range earlier {
		byLedger[venue.Ledger] = venue
	}

	return &Prover{
		source:       source,
		scheme:       scheme,
		confidential: confidential,
		statements:   statements,
		signer:       signer,
		current:      current,
		earlier:      byLedger,
		cache:        map[uint64]proven{},
	}
}

func (p *Prover) CheckCourt(ctx context.Context) error {
	if err := p.check(ctx, p.current); err != nil {
		return err
	}

	for _, venue := range p.earlier {
		if err := p.check(ctx, venue); err != nil {
			return err
		}
	}

	return nil
}

func (p *Prover) check(ctx context.Context, venue Venue) error {
	const op = "proving.Prover.CheckCourt"

	identity, err := venue.Court.Identity(ctx)
	if err != nil {
		return err
	}

	if identity != venue.Identity {
		return notary.Errf(
			notary.KindWrongCourt,
			op,
			nil,
			"the court answers as %s on chain %d, this notary signs for %s on chain %d",
			identity.Address.Hex(),
			identity.ChainID,
			venue.Identity.Address.Hex(),
			venue.Identity.ChainID,
		)
	}

	sample := notary.Statement{
		Epoch:        1,
		Commitment:   notary.Digest{1},
		Amount:       2,
		Index:        3,
		SiblingsHash: notary.Digest{4},
	}

	onChain, err := venue.Court.StatementDigest(ctx, sample)
	if err != nil {
		return err
	}

	if local := p.statements.StatementDigest(venue.Identity, sample); local != onChain {
		return notary.Errf(
			notary.KindDisagreement,
			op,
			nil,
			"the court %s hashes a statement to %s, this notary to %s",
			venue.Identity.Address.Hex(),
			onChain.Hex(),
			local.Hex(),
		)
	}

	newest, err := venue.Court.SignerOf(ctx, math.MaxUint64)
	if err != nil {
		return err
	}

	if !p.signer.Holds(newest) {
		return notary.Errf(
			notary.KindSignerUnavailable,
			op,
			nil,
			"the newest signer of the court %s is %s and this notary holds no key for it",
			venue.Identity.Address.Hex(),
			newest.Hex(),
		)
	}

	return nil
}

func (p *Prover) venueOf(record notary.EpochRecord) (Venue, error) {
	if record.Ledger == p.current.Ledger {
		return p.current, nil
	}

	if venue, ok := p.earlier[record.Ledger]; ok {
		return venue, nil
	}

	return Venue{}, notary.Errf(
		notary.KindDisagreement,
		"proving.Prover.venueOf",
		nil,
		"epoch %d was anchored to %s on chain %d, which is neither this notary's ledger nor one before it",
		record.Header.Number,
		record.Ledger.Address.Hex(),
		record.Ledger.ChainID,
	)
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

	venue, err := p.venueOf(record)
	if err != nil {
		return Bundle{}, err
	}

	built, err := p.load(ctx, venue, record)
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
			Scheme: p.scheme.Version(),
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

		if record.Confidential != nil {
			claim.Scheme = record.Confidential.Scheme
		}

		statement := notary.StatementOf(p.statements, claim)

		signature, err := p.signer.Sign(built.signer, p.statements.StatementDigest(venue.Identity, statement))
		if err != nil {
			return Bundle{}, err
		}

		claim.Statement = &notary.StatementSignature{Court: venue.Identity.Address, Signature: signature}
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

func (p *Prover) load(ctx context.Context, venue Venue, record notary.EpochRecord) (proven, error) {
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

	tree, err := notary.RebuildEpochTree(p.scheme, p.confidential, record, entries)
	if err != nil {
		return proven{}, err
	}

	signer, err := venue.Court.SignerOf(ctx, number)
	if err != nil {
		return proven{}, err
	}

	if signer.IsZero() {
		return proven{}, notary.Errf(
			notary.KindSignerUnavailable,
			op,
			nil,
			"the court %s names no signer for epoch %d",
			venue.Identity.Address.Hex(),
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
