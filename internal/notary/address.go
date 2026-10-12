package notary

import "encoding/hex"

type Address [20]byte

func ParseAddress(text string) (Address, error) {
	raw, err := parseFixedHex(text, len(Address{}))
	if err != nil {
		return Address{}, Errf(KindInvalidArgument, "notary.ParseAddress", err, "")
	}

	return Address(raw), nil
}

func (a Address) Hex() string {
	return "0x" + hex.EncodeToString(a[:])
}

func (a Address) IsZero() bool {
	return a == Address{}
}
