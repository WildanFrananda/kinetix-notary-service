package keccak

import (
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

const (
	WithholdingCourtName    = "Kinetix Withholding Court"
	WithholdingCourtVersion = "1"

	registrationType        = "Registration(address holder,bytes32 account,uint64 fromEpoch)"
	registrationRequestType = "RegistrationRequest(address holder,string principal,uint64 issuedAt)"
)

var _ notary.WithholdingScheme = Scheme{}

func (Scheme) AccountID(holder notary.Holder, kind notary.AccountKind, secret notary.AccountSecret) notary.Digest {
	name := holder.String()
	packed := make([]byte, 0, 1+1+32+len(name))
	packed = append(packed, accountIDTag, byte(kind))
	packed = append(packed, secret[:]...)
	packed = append(packed, name...)

	return digest(packed)
}

func (Scheme) AccountCommitment(account notary.Digest, salt notary.Salt) notary.Digest {
	packed := make([]byte, 0, 1+32+32)
	packed = append(packed, accountCommitmentTag)
	packed = append(packed, account[:]...)
	packed = append(packed, salt[:]...)

	return digest(packed)
}

func (Scheme) RegistrationDigest(court notary.CourtIdentity, registration notary.Registration) notary.Digest {
	structure := words(
		digest([]byte(registrationType)),
		addressWord(registration.Holder),
		registration.Account,
		uintWord(registration.FromEpoch),
	)

	return typedData(WithholdingCourtName, WithholdingCourtVersion, court, structure)
}

func (Scheme) RegistrationRequestDigest(
	court notary.CourtIdentity,
	request notary.RegistrationRequest,
) notary.Digest {
	structure := words(
		digest([]byte(registrationRequestType)),
		addressWord(request.Holder),
		digest([]byte(request.Principal.String())),
		uintWord(request.IssuedAt.Seconds()),
	)

	return typedData(WithholdingCourtName, WithholdingCourtVersion, court, structure)
}
