package notary

import "encoding/hex"

type AccountSecret [32]byte

func ParseAccountSecret(text string) (AccountSecret, error) {
	raw, err := parseFixedHex(text, len(AccountSecret{}))
	if err != nil {
		return AccountSecret{}, Errf(KindInvalidArgument, "notary.ParseAccountSecret", err, "")
	}

	return AccountSecret(raw), nil
}

func (s AccountSecret) Hex() string {
	return "0x" + hex.EncodeToString(s[:])
}

func (s AccountSecret) IsZero() bool {
	return s == AccountSecret{}
}
