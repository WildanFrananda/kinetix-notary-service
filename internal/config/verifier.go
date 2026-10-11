package config

import (
	"errors"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

const (
	RPCURLKey        = "KINETIX_NOTARY_RPC_URL"
	ChainIDKey       = "KINETIX_NOTARY_CHAIN_ID"
	LedgerAddressKey = "KINETIX_NOTARY_LEDGER_ADDRESS"
	RPCTimeoutKey    = "KINETIX_NOTARY_RPC_TIMEOUT"

	maxRPCTimeout = 5 * time.Minute
)

type Verifier struct {
	RPCURL     string
	Ledger     notary.LedgerIdentity
	RPCTimeout time.Duration
}

func LoadVerifier(lookup Lookup) (Verifier, error) {
	var problems []string
	read := func(key string) string {
		value, ok := lookup(key)
		value = strings.TrimSpace(value)
		if !ok || value == "" {
			problems = append(problems, key+" is not set")
		}
		return value
	}
	var cfg Verifier

	if raw := read(RPCURLKey); raw != "" {
		parsed, err := url.Parse(raw)
		switch {
		case err != nil:
			problems = append(problems, RPCURLKey+" is not a URL")
		case !supportedScheme(parsed.Scheme) || parsed.Host == "":
			problems = append(problems, RPCURLKey+" must be an http, https, ws or wss URL with a host")
		default:
			cfg.RPCURL = raw
		}
	}
	if raw := read(ChainIDKey); raw != "" {
		chainID, err := strconv.ParseUint(raw, 10, 64)
		if err != nil || chainID == 0 {
			problems = append(problems, ChainIDKey+" must be a positive integer")
		}
		cfg.Ledger.ChainID = chainID
	}
	if raw := read(LedgerAddressKey); raw != "" {
		address, err := notary.ParseAddress(raw)
		if err != nil || address.IsZero() {
			problems = append(problems, LedgerAddressKey+" must be a non-zero 0x-prefixed 20-byte address")
		}
		cfg.Ledger.Address = address
	}
	if raw := read(RPCTimeoutKey); raw != "" {
		timeout, err := time.ParseDuration(raw)
		if err != nil || timeout <= 0 || timeout > maxRPCTimeout {
			problems = append(problems, RPCTimeoutKey+" must be a duration above zero and at most 5m")
		}
		cfg.RPCTimeout = timeout
	}

	if len(problems) > 0 {
		return Verifier{}, errors.New("configuration refused: " + strings.Join(problems, "; "))
	}
	return cfg, nil
}

func supportedScheme(scheme string) bool {
	switch scheme {
	case "http", "https", "ws", "wss":
		return true
	default:
		return false
	}
}
