package answering

import (
	"context"
	"log/slog"
	"sort"
	"time"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type requestKey struct {
	registration notary.Digest
	epoch        uint64
}

type Answerer struct {
	court         notary.WithholdingCourt
	chain         notary.AnswerChain
	source        notary.ProofSource
	registrations notary.RegistrationStore
	scheme        notary.Scheme
	observer      Observer
	clock         func() time.Time
	log           *slog.Logger
	policy        Policy
	next          uint64
	started       bool
	pending       map[requestKey]notary.WithholdingRequest
	unanswerable  map[requestKey]string
}

func NewAnswerer(
	court notary.WithholdingCourt,
	chain notary.AnswerChain,
	source notary.ProofSource,
	registrations notary.RegistrationStore,
	scheme notary.Scheme,
	observer Observer,
	clock func() time.Time,
	log *slog.Logger,
	policy Policy,
) *Answerer {
	return &Answerer{
		court:         court,
		chain:         chain,
		source:        source,
		registrations: registrations,
		scheme:        scheme,
		observer:      observer,
		clock:         clock,
		log:           log,
		policy:        policy,
		pending:       map[requestKey]notary.WithholdingRequest{},
		unanswerable:  map[requestKey]string{},
	}
}

func (a *Answerer) CheckScan(ctx context.Context) error {
	const op = "answering.Answerer.CheckScan"

	latest, err := a.court.LatestBlock(ctx)
	if err != nil {
		return err
	}

	from := uint64(0)
	if latest >= a.policy.Span {
		from = latest - a.policy.Span + 1
	}

	if _, err := a.court.Requests(ctx, from, latest); err != nil {
		return notary.Errf(
			notary.KindOf(err),
			op,
			err,
			"the RPC endpoint does not answer one log query over %d blocks",
			a.policy.Span,
		)
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

	total, err := a.court.TotalSilences(ctx)
	if err != nil {
		return err
	}

	a.observer.Silences(total)

	if balance, err := a.chain.Balance(ctx); err == nil {
		a.observer.AnswererBalance(balance)
	}

	for _, request := range a.ordered() {
		if err := a.settle(ctx, request); err != nil {
			a.observe()

			return err
		}
	}

	a.observe()

	return nil
}

func (a *Answerer) scan(ctx context.Context) error {
	latest, err := a.court.LatestBlock(ctx)
	if err != nil {
		return err
	}

	if !a.started {
		a.next = a.policy.FromBlock
		if latest > a.policy.Lookback && latest-a.policy.Lookback > a.next {
			a.next = latest - a.policy.Lookback
		}

		a.started = true
	}

	for a.next <= latest {
		to := min(a.next+a.policy.Span-1, latest)

		requests, err := a.court.Requests(ctx, a.next, to)
		if err != nil {
			return err
		}

		for _, request := range requests {
			a.pending[requestKey{registration: request.Registration, epoch: request.Epoch}] = request
		}

		a.next = to + 1
	}

	return nil
}

func (a *Answerer) ordered() []notary.WithholdingRequest {
	requests := make([]notary.WithholdingRequest, 0, len(a.pending))
	for _, request := range a.pending {
		requests = append(requests, request)
	}

	sort.Slice(requests, func(i, j int) bool {
		return requests[i].DueAt.Seconds() < requests[j].DueAt.Seconds()
	})

	return requests
}

func (a *Answerer) settle(ctx context.Context, request notary.WithholdingRequest) error {
	key := requestKey{registration: request.Registration, epoch: request.Epoch}

	state, dueAt, err := a.court.RequestOf(ctx, request.Registration, request.Epoch)
	if err != nil {
		return err
	}

	if state != notary.RequestOpen {
		delete(a.pending, key)
		delete(a.unanswerable, key)

		return nil
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
		delete(a.pending, key)
		delete(a.unanswerable, key)

		return nil
	}

	if _, known := a.unanswerable[key]; known {
		return nil
	}

	answer, err := a.prepare(ctx, request)
	if err != nil {
		return a.refuse(ctx, key, request, err)
	}

	hash, err := a.chain.Answer(ctx, answer)
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
	delete(a.pending, key)

	return nil
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

	tree, err := notary.BuildTree(a.scheme, request.Epoch, entries)
	if err != nil {
		return notary.Answer{}, err
	}

	if root := tree.Root(); root.Hash != record.Header.Root || root.Sum != record.Header.Total {
		return notary.Answer{}, notary.Errf(
			notary.KindDisagreement,
			op,
			nil,
			"the accounts stored for epoch %d do not rebuild the root it anchored",
			request.Epoch,
		)
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
