package notary

import "encoding/hex"

type Blinding [31]byte

func ParseBlinding(text string) (Blinding, error) {
	raw, err := parseFixedHex(text, len(Blinding{}))
	if err != nil {
		return Blinding{}, Errf(KindInvalidArgument, "notary.ParseBlinding", err, "")
	}

	return Blinding(raw), nil
}

func (b Blinding) Hex() string {
	return "0x" + hex.EncodeToString(b[:])
}

func (b Blinding) IsZero() bool {
	return b == Blinding{}
}
