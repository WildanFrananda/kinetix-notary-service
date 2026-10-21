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
	CourtAddressKey  = "KINETIX_NOTARY_COURT_ADDRESS"
	RPCTimeoutKey    = "KINETIX_NOTARY_RPC_TIMEOUT"
	RPCRateKey       = "KINETIX_NOTARY_RPC_RATE"

	maxRPCTimeout = 5 * time.Minute
	maxRPCRate    = 1000
)

type Chain struct {
	RPCURL     string
	Ledger     notary.LedgerIdentity
	Court      notary.Address
	RPCTimeout time.Duration
	RPCRate    int
}

func LoadChain(lookup Lookup) (Chain, error) {
	var problems []string

	read := requiredReader(lookup, &problems)

	cfg := Chain{
		RPCURL:     rpcURL(read(RPCURLKey), &problems),
		Ledger:     notary.LedgerIdentity{ChainID: chainID(read(ChainIDKey), &problems)},
		Court:      address(CourtAddressKey, read(CourtAddressKey), &problems),
		RPCTimeout: rpcTimeout(read(RPCTimeoutKey), &problems),
		RPCRate:    rpcRate(read(RPCRateKey), &problems),
	}
	cfg.Ledger.Address = address(LedgerAddressKey, read(LedgerAddressKey), &problems)

	if parsed, err := url.Parse(cfg.RPCURL); err == nil && (parsed.Scheme == "ws" || parsed.Scheme == "wss") {
		problems = append(problems, RPCURLKey+" must be http or https: "+RPCRateKey+" paces HTTP requests")
	}

	if len(problems) > 0 {
		return Chain{}, errors.New("configuration refused: " + strings.Join(problems, "; "))
	}

	return cfg, nil
}

func requiredReader(lookup Lookup, problems *[]string) func(string) string {
	return func(key string) string {
		value, ok := lookup(key)

		value = strings.TrimSpace(value)
		if !ok || value == "" {
			*problems = append(*problems, key+" is not set")
		}

		return value
	}
}

func rpcURL(raw string, problems *[]string) string {
	if raw == "" {
		return ""
	}

	parsed, err := url.Parse(raw)

	switch {
	case err != nil:
		*problems = append(*problems, RPCURLKey+" is not a URL")
	case !supportedScheme(parsed.Scheme) || parsed.Host == "":
		*problems = append(*problems, RPCURLKey+" must be an http, https, ws or wss URL with a host")
	default:
		return raw
	}

	return ""
}

func chainID(raw string, problems *[]string) uint64 {
	if raw == "" {
		return 0
	}

	id, err := strconv.ParseUint(raw, 10, 64)
	if err != nil || id == 0 {
		*problems = append(*problems, ChainIDKey+" must be a positive integer")
	}

	return id
}

func address(key, raw string, problems *[]string) notary.Address {
	if raw == "" {
		return notary.Address{}
	}

	parsed, err := notary.ParseAddress(raw)
	if err != nil || parsed.IsZero() {
		*problems = append(*problems, key+" must be a non-zero 0x-prefixed 20-byte address")
	}

	return parsed
}

func rpcTimeout(raw string, problems *[]string) time.Duration {
	if raw == "" {
		return 0
	}

	timeout, err := time.ParseDuration(raw)
	if err != nil || timeout <= 0 || timeout > maxRPCTimeout {
		*problems = append(*problems, RPCTimeoutKey+" must be a duration above zero and at most 5m")
	}

	return timeout
}

func rpcRate(raw string, problems *[]string) int {
	if raw == "" {
		return 0
	}

	perSecond, err := strconv.Atoi(raw)
	if err != nil || perSecond < 1 || perSecond > maxRPCRate {
		*problems = append(
			*problems,
			RPCRateKey+" must be the most requests a second this process sends, from 1 to 1000",
		)
	}

	return perSecond
}

func supportedScheme(scheme string) bool {
	switch scheme {
	case "http", "https", "ws", "wss":
		return true
	default:
		return false
	}
}
