package notary

type SponsoredCall struct {
	Target        Address
	Data          []byte
	PaymasterData []byte
}
