package notary

import "encoding/hex"

type Digest [32]byte

func ParseDigest(text string) (Digest, error) {
	raw, err := parseFixedHex(text, len(Digest{}))
	if err != nil {
		return Digest{}, Errf(KindInvalidArgument, "notary.ParseDigest", err, "")
	}

	return Digest(raw), nil
}

func (d Digest) Hex() string {
	return "0x" + hex.EncodeToString(d[:])
}

func (d Digest) IsZero() bool {
	return d == Digest{}
}
