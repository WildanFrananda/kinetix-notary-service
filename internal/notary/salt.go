package notary

import "encoding/hex"

type Salt [32]byte

func ParseSalt(text string) (Salt, error) {
	raw, err := parseFixedHex(text, len(Salt{}))
	if err != nil {
		return Salt{}, Errf(KindInvalidArgument, "notary.ParseSalt", err, "")
	}

	return Salt(raw), nil
}

func (s Salt) Hex() string {
	return "0x" + hex.EncodeToString(s[:])
}

func (s Salt) IsZero() bool {
	return s == Salt{}
}
