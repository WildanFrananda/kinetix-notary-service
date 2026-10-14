package config

import (
	"errors"
	"net"
	"strconv"
	"strings"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

const (
	CouncilAddressKey = "KINETIX_NOTARY_COUNCIL_ADDRESS"
	WatchKeyKey       = "KINETIX_NOTARY_WATCH_PRIVATE_KEY"
	WatchGRPCPortKey  = "KINETIX_NOTARY_WATCH_GRPC_PORT"
)

type Watch struct {
	Chain          Verifier
	Council        notary.Address
	WatchKey       Secret
	PaymentAddress string
	PKIDir         string
	TrustDomains   []string
	GRPCPort       int
	HTTPPort       int
}

func LoadWatch(lookup Lookup) (Watch, error) {
	var problems []string

	chain, err := LoadVerifier(lookup)
	if err != nil {
		problems = append(problems, strings.TrimPrefix(err.Error(), "configuration refused: "))
	}

	cfg := Watch{Chain: chain}
	required := func(key string) string {
		value, ok := lookup(key)

		value = strings.TrimSpace(value)
		if !ok || value == "" {
			problems = append(problems, key+" is not set")
		}

		return value
	}
	port := func(key string) int {
		raw := required(key)
		if raw == "" {
			return 0
		}

		value, err := strconv.Atoi(raw)
		if err != nil || value < 1 || value > 65535 {
			problems = append(problems, key+" must be a port from 1 to 65535")
		}

		return value
	}

	cfg.Council = councilAddress(required(CouncilAddressKey), &problems)

	if raw := required(WatchKeyKey); raw != "" {
		cfg.WatchKey = Secret{value: raw}
	}

	if raw := required(PaymentAddressKey); raw != "" {
		host, port, err := net.SplitHostPort(raw)
		if err != nil || host == "" || port == "" {
			problems = append(problems, PaymentAddressKey+" must be host:port")
		}

		cfg.PaymentAddress = raw
	}

	cfg.PKIDir = required(PKIDirKey)
	cfg.TrustDomains = trustDomains(lookup)
	cfg.GRPCPort = port(WatchGRPCPortKey)
	cfg.HTTPPort = port(HTTPPortKey)

	if cfg.GRPCPort != 0 && cfg.GRPCPort == cfg.HTTPPort {
		problems = append(problems, WatchGRPCPortKey+" and "+HTTPPortKey+" must differ")
	}

	if len(problems) > 0 {
		return Watch{}, errors.New("configuration refused: " + strings.Join(problems, "; "))
	}

	return cfg, nil
}

func councilAddress(raw string, problems *[]string) notary.Address {
	if raw == "" {
		return notary.Address{}
	}

	address, err := notary.ParseAddress(raw)
	if err != nil || address.IsZero() {
		*problems = append(*problems, CouncilAddressKey+" must be a non-zero 0x-prefixed 20-byte address")
	}

	return address
}

func trustDomains(lookup Lookup) []string {
	raw, _ := lookup(TrustDomainKey)
	if raw = strings.TrimSpace(raw); raw == "" {
		return []string{defaultTrustDomain}
	}

	var domains []string

	for _, domain := range strings.Split(raw, ",") {
		if domain = strings.TrimSpace(domain); domain != "" {
			domains = append(domains, domain)
		}
	}

	return domains
}
