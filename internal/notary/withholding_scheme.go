package notary

type WithholdingScheme interface {
	RegistrationDomain(court CourtIdentity) TypedDataDomain
	RegistrationDigest(court CourtIdentity, registration Registration) Digest
	RegistrationRequestDigest(court CourtIdentity, request RegistrationRequest) Digest
}
