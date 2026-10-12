package notary

import (
	"encoding/hex"
	"fmt"
	"strings"
)

func parseFixedHex(text string, size int) ([]byte, error) {
	digits, found := strings.CutPrefix(text, "0x")
	if !found {
		return nil, fmt.Errorf("%q has no 0x prefix", text)
	}

	if len(digits) != size*2 {
		return nil, fmt.Errorf("%q is not %d bytes", text, size)
	}

	raw, err := hex.DecodeString(digits)
	if err != nil {
		return nil, fmt.Errorf("%q is not hexadecimal", text)
	}

	return raw, nil
}
