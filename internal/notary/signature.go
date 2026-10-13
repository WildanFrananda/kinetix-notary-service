package notary

import "encoding/hex"

type Signature [65]byte

func ParseSignature(text string) (Signature, error) {
	raw, err := parseFixedHex(text, len(Signature{}))
	if err != nil {
		return Signature{}, Errf(KindInvalidArgument, "notary.ParseSignature", err, "")
	}

	return Signature(raw), nil
}

func (s Signature) Hex() string {
	return "0x" + hex.EncodeToString(s[:])
}
