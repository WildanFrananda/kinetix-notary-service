package config

import (
	"errors"
	"strings"
	"time"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

const DirectoryAddressKey = "KINETIX_NOTARY_DIRECTORY_ADDRESS"

type Verifier struct {
	RPCURL     string
	Directory  notary.DirectoryIdentity
	RPCTimeout time.Duration
}

func LoadVerifier(lookup Lookup) (Verifier, error) {
	var problems []string

	read := requiredReader(lookup, &problems)

	cfg := Verifier{
		RPCURL: rpcURL(read(RPCURLKey), &problems),
		Directory: notary.DirectoryIdentity{
			ChainID: chainID(read(ChainIDKey), &problems),
			Address: address(DirectoryAddressKey, read(DirectoryAddressKey), &problems),
		},
		RPCTimeout: rpcTimeout(read(RPCTimeoutKey), &problems),
	}

	if len(problems) > 0 {
		return Verifier{}, errors.New("configuration refused: " + strings.Join(problems, "; "))
	}

	return cfg, nil
}
