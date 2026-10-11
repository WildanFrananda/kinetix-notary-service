package ethereum

import (
	"errors"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
)

const (
	selectorLength = 4
	abiWordLength  = 32
)

var errNotALedger = errors.New("the contract does not answer as a liability ledger")

func isBareRevert(err error) bool {
	return strings.Contains(err.Error(), "execution reverted")
}

func revertData(err error) ([]byte, bool) {
	var dataErr rpc.DataError
	if !errors.As(err, &dataErr) {
		return nil, false
	}
	text, ok := dataErr.ErrorData().(string)
	if !ok {
		return nil, false
	}
	raw, decodeErr := hexutil.Decode(text)
	if decodeErr != nil || len(raw) < selectorLength {
		return nil, false
	}
	return raw, true
}
