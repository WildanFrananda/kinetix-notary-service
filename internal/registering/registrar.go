package registering

import (
	"context"
	"time"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

const (
	MaxRequestAge  = 5 * time.Minute
	MaxRequestSkew = time.Minute
)

type Registrar struct {
	source      notary.ProofSource
	scheme      notary.Scheme
	withholding notary.WithholdingScheme
	court       notary.Court
	chain       notary.WithholdingCourt
	signer      notary.StatementSigner
	recoverer   notary.SignatureRecoverer
	store       notary.RegistrationStore
	identity    notary.CourtIdentity
	clock       func() time.Time
}

func NewRegistrar(
	source notary.ProofSource,
	scheme notary.Scheme,
	withholding notary.WithholdingScheme,
	court notary.Court,
	chain notary.WithholdingCourt,
	signer notary.StatementSigner,
	recoverer notary.SignatureRecoverer,
	store notary.RegistrationStore,
	identity notary.CourtIdentity,
	clock func() time.Time,
) *Registrar {
	return &Registrar{
		source:      source,
		scheme:      scheme,
		withholding: withholding,
		court:       court,
		chain:       chain,
		signer:      signer,
		recoverer:   recoverer,
		store:       store,
		identity:    identity,
		clock:       clock,
	}
}

func (r *Registrar) CheckWithholding(ctx context.Context) error {
	const op = "registering.Registrar.CheckWithholding"

	identity, err := r.chain.Identity(ctx)
	if err != nil {
		return err
	}

	if identity != r.identity {
		return notary.Errf(
			notary.KindWrongCourt,
			op,
			nil,
			"the withholding court answers as %s on chain %d, this notary registers for %s on chain %d",
			identity.Address.Hex(),
			identity.ChainID,
			r.identity.Address.Hex(),
			r.identity.ChainID,
		)
	}

	sample := notary.Registration{Holder: notary.Address{1}, Account: notary.Digest{2}, FromEpoch: 3}

	onChain, err := r.chain.RegistrationDigest(ctx, sample)
	if err != nil {
		return err
	}

	if local := r.withholding.RegistrationDigest(r.identity, sample); local != onChain {
		return notary.Errf(
			notary.KindDisagreement,
			op,
			nil,
			"the withholding court hashes a registration to %s, this notary to %s",
			onChain.Hex(),
			local.Hex(),
		)
	}

	return nil
}

func (r *Registrar) Register(
	ctx context.Context,
	principal notary.Holder,
	request notary.RegistrationRequest,
	signature notary.Signature,
) (Issued, error) {
	const op = "registering.Registrar.Register"

	if err := r.admit(principal, request, signature); err != nil {
		return Issued{}, err
	}

	record, found, err := r.source.LastFinalized(ctx)
	if err != nil {
		return Issued{}, err
	}

	if !found {
		return Issued{}, notary.Errf(notary.KindRegistrationNotOpen, op, nil, "no epoch is final yet")
	}

	number := record.Header.Number

	accounts, err := r.source.Accounts(ctx, number)
	if err != nil {
		return Issued{}, err
	}

	mine := make([]notary.Account, 0, 2)

	for _, account := range accounts {
		if account.Liability.Holder != principal {
			continue
		}

		if account.Secret.IsZero() {
			return Issued{}, notary.Errf(
				notary.KindRegistrationNotOpen,
				op,
				nil,
				"epoch %d was built before account secrets; registration opens with the next final epoch",
				number,
			)
		}

		if r.scheme.AccountID(account.Liability.Holder, account.Liability.Kind, account.Secret) != account.ID {
			return Issued{}, notary.Errf(
				notary.KindDisagreement,
				op,
				nil,
				"the account at position %d of epoch %d does not derive from its secret",
				account.Position,
				number,
			)
		}

		mine = append(mine, account)
	}

	if len(mine) == 0 {
		return Issued{}, notary.Errf(
			notary.KindNoAccount,
			op,
			nil,
			"epoch %d holds no account of this principal",
			number,
		)
	}

	signer, err := r.court.SignerOf(ctx, number)
	if err != nil {
		return Issued{}, err
	}

	if signer.IsZero() {
		return Issued{}, notary.Errf(
			notary.KindSignerUnavailable,
			op,
			nil,
			"the court names no signer for epoch %d",
			number,
		)
	}

	issued := make([]notary.IssuedRegistration, 0, len(mine))

	for _, account := range mine {
		registration := notary.Registration{Holder: request.Holder, Account: account.ID, FromEpoch: number}
		digest := r.withholding.RegistrationDigest(r.identity, registration)

		sealed, err := r.signer.Sign(signer, digest)
		if err != nil {
			return Issued{}, err
		}

		issued = append(issued, notary.IssuedRegistration{
			Registration: registration,
			Kind:         account.Liability.Kind,
			Digest:       digest,
			Signature:    sealed,
		})
	}

	if err := r.store.SaveRegistrations(ctx, principal, issued); err != nil {
		return Issued{}, err
	}

	return Issued{Court: r.identity, Epoch: number, Registrations: issued}, nil
}

func (r *Registrar) admit(
	principal notary.Holder,
	request notary.RegistrationRequest,
	signature notary.Signature,
) error {
	const op = "registering.Registrar.admit"

	if request.Principal != principal {
		return notary.Errf(notary.KindRegistrationRefused, op, nil, "the request is signed for another principal")
	}

	if request.Holder.IsZero() {
		return notary.Errf(notary.KindInvalidArgument, op, nil, "a registration names the address that will request")
	}

	now := r.clock()
	issuedAt := request.IssuedAt.Time()

	if issuedAt.Before(now.Add(-MaxRequestAge)) || issuedAt.After(now.Add(MaxRequestSkew)) {
		return notary.Errf(
			notary.KindRegistrationRefused,
			op,
			nil,
			"the request was issued at %s; it must be at most %s old",
			issuedAt.UTC().Format(time.RFC3339),
			MaxRequestAge,
		)
	}

	recovered, err := r.recoverer.Recover(r.withholding.RegistrationRequestDigest(r.identity, request), signature)
	if err != nil {
		return err
	}

	if recovered != request.Holder {
		return notary.Errf(
			notary.KindRegistrationRefused,
			op,
			nil,
			"the request is signed by %s, not by %s",
			recovered.Hex(),
			request.Holder.Hex(),
		)
	}

	return nil
}
