package notary

type PreparedAnchor struct {
	Hash     Digest
	Raw      []byte
	Nonce    uint64
	GasLimit uint64
	Fees     FeeQuote
}
