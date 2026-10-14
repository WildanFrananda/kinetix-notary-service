package notary

type WithholdingScheme interface {
	RegistrationDigest(court CourtIdentity, registration Registration) Digest
	RegistrationRequestDigest(court CourtIdentity, request RegistrationRequest) Digest
}
