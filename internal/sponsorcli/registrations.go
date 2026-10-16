package sponsorcli

import (
	"encoding/json"
	"io"
	"strconv"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

func decodeRegistrations(r io.Reader) (notary.Address, []notary.IssuedRegistration, error) {
	const op = "sponsorcli.decodeRegistrations"

	decoder := json.NewDecoder(r)
	decoder.DisallowUnknownFields()

	var file registrationsFile
	if err := decoder.Decode(&file); err != nil {
		return notary.Address{}, nil, notary.Errf(notary.KindInvalidArgument, op, err, "not a registrations document")
	}

	withholding, err := notary.ParseAddress(file.Court.Address)
	if err != nil {
		return notary.Address{}, nil, err
	}

	issued := make([]notary.IssuedRegistration, 0, len(file.Registrations))
	for _, field := range file.Registrations {
		registration, err := registrationOf(field)
		if err != nil {
			return notary.Address{}, nil, err
		}

		issued = append(issued, registration)
	}

	if len(issued) == 0 {
		return notary.Address{}, nil, notary.Errf(
			notary.KindInvalidArgument,
			op,
			nil,
			"the document holds no registration",
		)
	}

	return withholding, issued, nil
}

func registrationOf(field registrationField) (notary.IssuedRegistration, error) {
	holder, err := notary.ParseAddress(field.Holder)
	if err != nil {
		return notary.IssuedRegistration{}, err
	}

	account, err := notary.ParseDigest(field.Account)
	if err != nil {
		return notary.IssuedRegistration{}, err
	}

	fromEpoch, err := strconv.ParseUint(field.FromEpoch, 10, 64)
	if err != nil {
		return notary.IssuedRegistration{}, notary.Errf(
			notary.KindInvalidArgument,
			"sponsorcli.registrationOf",
			err,
			"fromEpoch %q is not an epoch",
			field.FromEpoch,
		)
	}

	kind, err := notary.AccountKindNamed(field.Kind)
	if err != nil {
		return notary.IssuedRegistration{}, err
	}

	digest, err := notary.ParseDigest(field.Digest)
	if err != nil {
		return notary.IssuedRegistration{}, err
	}

	signature, err := notary.ParseSignature(field.Signature)
	if err != nil {
		return notary.IssuedRegistration{}, err
	}

	return notary.IssuedRegistration{
		Registration: notary.Registration{Holder: holder, Account: account, FromEpoch: fromEpoch},
		Kind:         kind,
		Digest:       digest,
		Signature:    signature,
	}, nil
}
