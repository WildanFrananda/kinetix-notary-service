package config

import (
	"errors"
	"strings"
	"time"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

const EvidenceRegistryKey = "KINETIX_NOTARY_EVIDENCE_REGISTRY_ADDRESS"

type EvidenceVerifier struct {
	RPCURL     string
	Registry   notary.RegistryIdentity
	RPCTimeout time.Duration
}

func LoadEvidenceVerifier(lookup Lookup) (EvidenceVerifier, error) {
	var problems []string

	read := requiredReader(lookup, &problems)

	cfg := EvidenceVerifier{
		RPCURL: rpcURL(read(RPCURLKey), &problems),
		Registry: notary.RegistryIdentity{
			ChainID: chainID(read(ChainIDKey), &problems),
			Address: address(EvidenceRegistryKey, read(EvidenceRegistryKey), &problems),
		},
		RPCTimeout: rpcTimeout(read(RPCTimeoutKey), &problems),
	}

	if len(problems) > 0 {
		return EvidenceVerifier{}, errors.New("configuration refused: " + strings.Join(problems, "; "))
	}

	return cfg, nil
}
