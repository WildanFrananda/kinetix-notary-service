package notary

type RegistrationRequest struct {
	Holder    Address
	Principal Holder
	IssuedAt  UnixTime
}
