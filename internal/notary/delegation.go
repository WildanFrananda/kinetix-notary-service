package notary

type Delegation struct {
	ChainID  uint64
	Delegate Address
	Nonce    uint64
	YParity  uint8
	R        Digest
	S        Digest
}
