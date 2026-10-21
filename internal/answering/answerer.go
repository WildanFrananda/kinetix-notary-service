package answering

import (
	"context"
	"log/slog"
	"sort"
	"time"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type Answerer struct {
	venues        []Venue
	cursors       []cursor
	scans         notary.WithholdingScans
	source        notary.ProofSource
	registrations notary.RegistrationStore
	scheme        notary.Scheme
	confidential  notary.ConfidentialScheme
	observer      Observer
	clock         func() time.Time
	log           *slog.Logger
	policy        Policy
	pending       map[requestKey]waiting
	unanswerable  map[requestKey]string
}

func NewAnswerer(
	venues []Venue,
	scans notary.WithholdingScans,
	source notary.ProofSource,
	registrations notary.RegistrationStore,
	scheme notary.Scheme,
	confidential notary.ConfidentialScheme,
	observer Observer,
	clock func() time.Time,
	log *slog.Logger,
	policy Policy,
) *Answerer {
	return &Answerer{
		venues:        venues,
		cursors:       make([]cursor, len(venues)),
		scans:         scans,
		source:        source,
		registrations: registrations,
		scheme:        scheme,
		confidential:  confidential,
		observer:      observer,
		clock:         clock,
		log:           log,
		policy:        policy,
		pending:       map[requestKey]waiting{},
		unanswerable:  map[requestKey]string{},
	}
}

func (a *Answerer) CheckScan(ctx context.Context) error {
	const op = "answering.Answerer.CheckScan"

	for _, venue := range a.venues {
		latest, err := venue.Court.LatestBlock(ctx)
		if err != nil {
			return err
		}

		from := uint64(0)
		if latest >= a.policy.Span {
			from = latest - a.policy.Span + 1
		}

		if _, err := venue.Court.Requests(ctx, from, latest); err != nil {
			return notary.Errf(
				notary.KindOf(err),
				op,
				err,
				"the RPC endpoint does not answer one log query over %d blocks",
				a.policy.Span,
			)
		}
	}

	return nil
}

func (a *Answerer) Tick(ctx context.Context) error {
	err := a.serve(ctx)
	a.observer.Pass(err == nil)

	return err
}

func (a *Answerer) serve(ctx context.Context) error {
	if err := a.scan(ctx); err != nil {
		return err
	}

	var total uint64

	for _, venue := range a.venues {
		silences, err := venue.Court.TotalSilences(ctx)
		if err != nil {
			return err
		}

		total += silences
	}

	a.observer.Silences(total)

	if len(a.venues) > 0 {
		if balance, err := a.venues[0].Chain.Balance(ctx); err == nil {
			a.observer.AnswererBalance(balance)
		}
	}

	for _, open := range a.ordered() {
		if err := a.settle(ctx, open); err != nil {
			a.observe()

			return err
		}
	}

	a.observe()

	return nil
}

func (a *Answerer) scan(ctx context.Context) error {
	for i, venue := range a.venues {
		if err := a.scanVenue(ctx, i, venue); err != nil {
			return err
		}
	}

	return nil
}

func (a *Answerer) scanVenue(ctx context.Context, i int, venue Venue) error {
	latest, err := venue.Court.LatestBlock(ctx)
	if err != nil {
		return err
	}

	position := &a.cursors[i]
	if !position.started {
		if err := a.resume(ctx, i, venue, latest); err != nil {
			return err
		}
	}

	for position.next <= latest {
		to := min(position.next+a.policy.Span-1, latest)

		requests, err := venue.Court.Requests(ctx, position.next, to)
		if err != nil {
			return err
		}

		if err := a.scans.RecordScan(ctx, venue.Identity, to+1, requests); err != nil {
			return err
		}

		for _, request := range requests {
			a.wait(i, request)
		}

		position.next = to + 1
	}

	return nil
}

func (a *Answerer) resume(ctx context.Context, i int, venue Venue, latest uint64) error {
	next, found, err := a.scans.ScanFrom(ctx, venue.Identity)
	if err != nil {
		return err
	}

	if !found || next < venue.FromBlock {
		next = venue.FromBlock
	}

	if latest > a.policy.Lookback && latest-a.policy.Lookback > next {
		next = latest - a.policy.Lookback
	}

	open, err := a.scans.OpenRequests(ctx, venue.Identity)
	if err != nil {
		return err
	}

	for _, request := range open {
		a.wait(i, request)
	}

	a.cursors[i] = cursor{next: next, started: true}

	return nil
}

func (a *Answerer) wait(venue int, request notary.WithholdingRequest) {
	key := requestKey{venue: venue, registration: request.Registration, epoch: request.Epoch}
	a.pending[key] = waiting{venue: venue, request: request}
}

func (a *Answerer) close(ctx context.Context, key requestKey) error {
	delete(a.pending, key)
	delete(a.unanswerable, key)

	return a.scans.CloseRequest(ctx, a.venues[key.venue].Identity, key.registration, key.epoch)
}

func (a *Answerer) ordered() []waiting {
	open := make([]waiting, 0, len(a.pending))
	for _, entry := range a.pending {
		open = append(open, entry)
	}

	sort.Slice(open, func(i, j int) bool {
		return open[i].request.DueAt.Seconds() < open[j].request.DueAt.Seconds()
	})

	return open
}

func (a *Answerer) settle(ctx context.Context, open waiting) error {
	request := open.request
	venue := a.venues[open.venue]
	key := requestKey{venue: open.venue, registration: request.Registration, epoch: request.Epoch}

	state, dueAt, err := venue.Court.RequestOf(ctx, request.Registration, request.Epoch)
	if err != nil {
		return err
	}

	if state != notary.RequestOpen {
		return a.close(ctx, key)
	}

	if a.clock().After(dueAt.Time()) {
		a.log.ErrorContext(
			ctx,
			"a statement request went unanswered past its deadline",
			"epoch",
			request.Epoch,
			"registration",
			request.Registration.Hex(),
			"due_at",
			dueAt.Time(),
		)

		return a.close(ctx, key)
	}

	if _, known := a.unanswerable[key]; known {
		return nil
	}

	answer, err := a.prepare(ctx, request)
	if err != nil {
		return a.refuse(ctx, key, request, err)
	}

	hash, err := venue.Chain.Answer(ctx, answer)
	if err != nil {
		return a.refuse(ctx, key, request, err)
	}

	a.log.InfoContext(
		ctx,
		"statement request answered",
		"epoch",
		request.Epoch,
		"registration",
		request.Registration.Hex(),
		"tx",
		hash.Hex(),
	)

	return a.close(ctx, key)
}

func (a *Answerer) refuse(ctx context.Context, key requestKey, request notary.WithholdingRequest, err error) error {
	switch notary.KindOf(err) {
	case notary.KindUnanswerable, notary.KindAnswerRefused, notary.KindDisagreement:
	default:
		return err
	}

	a.unanswerable[key] = err.Error()
	a.log.ErrorContext(
		ctx,
		"a statement request cannot be answered",
		"epoch",
		request.Epoch,
		"registration",
		request.Registration.Hex(),
		"kind",
		notary.KindOf(err).String(),
		"err",
		err.Error(),
	)

	return nil
}

func (a *Answerer) prepare(ctx context.Context, request notary.WithholdingRequest) (notary.Answer, error) {
	const op = "answering.Answerer.prepare"

	issued, found, err := a.registrations.RegistrationByDigest(ctx, request.Registration)
	if err != nil {
		return notary.Answer{}, err
	}

	if !found {
		return notary.Answer{}, notary.Errf(
			notary.KindUnanswerable,
			op,
			nil,
			"this notary's database holds no registration %s",
			request.Registration.Hex(),
		)
	}

	record, found, err := a.source.Epoch(ctx, request.Epoch)
	if err != nil {
		return notary.Answer{}, err
	}

	if !found {
		return notary.Answer{}, notary.Errf(
			notary.KindUnanswerable,
			op,
			nil,
			"this notary's database holds no epoch %d",
			request.Epoch,
		)
	}

	accounts, err := a.source.Accounts(ctx, request.Epoch)
	if err != nil {
		return notary.Answer{}, err
	}

	entries := make([]notary.Entry, len(accounts))
	position := -1

	for i, account := range accounts {
		entries[i] = notary.Entry{Commitment: account.Commitment, Amount: account.Liability.Amount}

		if account.Secret.IsZero() || account.ID != issued.Registration.Account {
			continue
		}

		if a.scheme.AccountCommitment(account.ID, account.Salt) != account.Commitment {
			return notary.Answer{}, notary.Errf(
				notary.KindDisagreement,
				op,
				nil,
				"the account at position %d of epoch %d does not open its commitment",
				i,
				request.Epoch,
			)
		}

		position = i
	}

	if position < 0 {
		return notary.Answer{}, notary.Errf(
			notary.KindUnanswerable,
			op,
			nil,
			"epoch %d holds no account %s",
			request.Epoch,
			issued.Registration.Account.Hex(),
		)
	}

	tree, err := notary.RebuildEpochTree(a.scheme, a.confidential, record, entries)
	if err != nil {
		return notary.Answer{}, err
	}

	proof, err := tree.Prove(uint64(position))
	if err != nil {
		return notary.Answer{}, err
	}

	account := accounts[position]

	return notary.Answer{
		Registration: issued.Registration,
		Epoch:        request.Epoch,
		Salt:         account.Salt,
		Amount:       account.Liability.Amount,
		Proof:        proof,
	}, nil
}

func (a *Answerer) observe() {
	a.observer.OpenRequests(len(a.pending))
	a.observer.Unanswerable(len(a.unanswerable))
}

func (a *Answerer) Run(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		if err := a.Tick(ctx); err != nil && ctx.Err() == nil {
			a.log.WarnContext(
				ctx,
				"the withholding court could not be served this pass",
				"kind",
				notary.KindOf(err).String(),
				"err",
				err.Error(),
			)
		}

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}
