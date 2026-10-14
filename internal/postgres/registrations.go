package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/storegen"
)

var _ notary.RegistrationStore = (*Store)(nil)

func (s *Store) SaveRegistrations(
	ctx context.Context,
	principal notary.Holder,
	issued []notary.IssuedRegistration,
) error {
	const op = "postgres.Store.SaveRegistrations"

	params := storegen.InsertRegistrationsParams{}

	for _, one := range issued {
		fromEpoch, err := signed(op, one.Registration.FromEpoch)
		if err != nil {
			return err
		}

		params.Digests = append(params.Digests, one.Digest[:])
		params.HolderAddresses = append(params.HolderAddresses, one.Registration.Holder[:])
		params.Accounts = append(params.Accounts, one.Registration.Account[:])
		params.FromEpochs = append(params.FromEpochs, fromEpoch)
		params.Principals = append(params.Principals, principal.String())
		params.Kinds = append(params.Kinds, int16(one.Kind))
		params.Signatures = append(params.Signatures, one.Signature[:])
	}

	if err := s.queries.InsertRegistrations(ctx, params); err != nil {
		return unavailable(op, err)
	}

	return nil
}

func (s *Store) RegistrationByDigest(
	ctx context.Context,
	digest notary.Digest,
) (notary.IssuedRegistration, bool, error) {
	const op = "postgres.Store.RegistrationByDigest"

	row, err := s.queries.RegistrationByDigest(ctx, digest[:])
	if errors.Is(err, pgx.ErrNoRows) {
		return notary.IssuedRegistration{}, false, nil
	}

	if err != nil {
		return notary.IssuedRegistration{}, false, unavailable(op, err)
	}

	issued, err := registrationOf(op, row)
	if err != nil {
		return notary.IssuedRegistration{}, false, err
	}

	return issued, true, nil
}

func registrationOf(op string, row storegen.RegistrationByDigestRow) (notary.IssuedRegistration, error) {
	digest, err := digestOf(op, row.Digest)
	if err != nil {
		return notary.IssuedRegistration{}, err
	}

	holder, err := addressOf(op, row.HolderAddress)
	if err != nil {
		return notary.IssuedRegistration{}, err
	}

	account, err := digestOf(op, row.Account)
	if err != nil {
		return notary.IssuedRegistration{}, err
	}

	fromEpoch, err := unsigned(op, row.FromEpoch)
	if err != nil {
		return notary.IssuedRegistration{}, err
	}

	if row.Kind < 0 || row.Kind > 255 {
		return notary.IssuedRegistration{}, notary.Errf(
			notary.KindStoreUnavailable,
			op,
			nil,
			"the database holds kind %d",
			row.Kind,
		)
	}

	kind, err := notary.AccountKindOf(uint8(row.Kind))
	if err != nil {
		return notary.IssuedRegistration{}, notary.Errf(notary.KindStoreUnavailable, op, err, "")
	}

	var signature notary.Signature
	if len(row.Signature) != len(signature) {
		return notary.IssuedRegistration{}, notary.Errf(
			notary.KindStoreUnavailable,
			op,
			nil,
			"the database holds a %d-byte signature",
			len(row.Signature),
		)
	}

	copy(signature[:], row.Signature)

	return notary.IssuedRegistration{
		Registration: notary.Registration{Holder: holder, Account: account, FromEpoch: fromEpoch},
		Kind:         kind,
		Digest:       digest,
		Signature:    signature,
	}, nil
}
