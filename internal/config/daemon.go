package config

import (
	"errors"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	DatabaseURLKey     = "KINETIX_NOTARY_DATABASE_URL"
	AnchorerKeyKey     = "KINETIX_NOTARY_ANCHORER_PRIVATE_KEY"
	PaymentAddressKey  = "KINETIX_NOTARY_PAYMENT_GRPC_ADDR"
	PKIDirKey          = "KINETIX_PKI_DIR"
	TrustDomainKey     = "KINETIX_TRUST_DOMAIN"
	CurrencyKey        = "KINETIX_NOTARY_CURRENCY"
	MaxFeeKey          = "KINETIX_NOTARY_MAX_FEE_PER_GAS_WEI"
	HTTPPortKey        = "KINETIX_NOTARY_HTTP_PORT"
	EpochIntervalKey   = "KINETIX_NOTARY_EPOCH_INTERVAL"
	TickKey            = "KINETIX_NOTARY_TICK"
	SnapshotTimeoutKey = "KINETIX_NOTARY_SNAPSHOT_TIMEOUT"

	defaultTrustDomain     = "kinetix.local"
	defaultEpochInterval   = time.Hour
	defaultTick            = 15 * time.Second
	defaultSnapshotTimeout = 3 * time.Minute
)

type Daemon struct {
	Chain           Verifier
	DatabaseURL     Secret
	AnchorerKey     Secret
	PaymentAddress  string
	PKIDir          string
	TrustDomains    []string
	Currency        string
	MaxFeePerGasWei uint64
	HTTPPort        int
	EpochInterval   time.Duration
	Tick            time.Duration
	SnapshotTimeout time.Duration
}

func LoadDaemon(lookup Lookup) (Daemon, error) {
	var problems []string

	chain, err := LoadVerifier(lookup)
	if err != nil {
		problems = append(problems, strings.TrimPrefix(err.Error(), "configuration refused: "))
	}

	cfg := Daemon{Chain: chain}
	required := func(key string) string {
		value, ok := lookup(key)

		value = strings.TrimSpace(value)
		if !ok || value == "" {
			problems = append(problems, key+" is not set")
		}

		return value
	}
	optional := func(key string) string {
		value, _ := lookup(key)
		return strings.TrimSpace(value)
	}
	duration := func(key string, fallback, least, most time.Duration) time.Duration {
		raw := optional(key)
		if raw == "" {
			return fallback
		}

		d, err := time.ParseDuration(raw)
		if err != nil || d < least || d > most {
			problems = append(problems, key+" must be a duration from "+least.String()+" to "+most.String())
		}

		return d
	}

	if raw := required(DatabaseURLKey); raw != "" {
		parsed, err := url.Parse(raw)
		if err != nil || (parsed.Scheme != "postgres" && parsed.Scheme != "postgresql") || parsed.Host == "" {
			problems = append(problems, DatabaseURLKey+" must be a postgres:// URL with a host")
		}

		cfg.DatabaseURL = Secret{value: raw}
	}

	if raw := required(AnchorerKeyKey); raw != "" {
		cfg.AnchorerKey = Secret{value: raw}
	}

	if raw := required(PaymentAddressKey); raw != "" {
		host, port, err := net.SplitHostPort(raw)
		if err != nil || host == "" || port == "" {
			problems = append(problems, PaymentAddressKey+" must be host:port")
		}

		cfg.PaymentAddress = raw
	}

	cfg.PKIDir = required(PKIDirKey)

	cfg.TrustDomains = []string{defaultTrustDomain}
	if raw := optional(TrustDomainKey); raw != "" {
		cfg.TrustDomains = nil

		for _, domain := range strings.Split(raw, ",") {
			if domain = strings.TrimSpace(domain); domain != "" {
				cfg.TrustDomains = append(cfg.TrustDomains, domain)
			}
		}
	}

	if raw := required(CurrencyKey); raw != "" {
		if len(raw) != 3 || strings.ToUpper(raw) != raw ||
			strings.ContainsFunc(raw, func(r rune) bool { return r < 'A' || r > 'Z' }) {
			problems = append(problems, CurrencyKey+" must be an ISO 4217 code such as IDR")
		}

		cfg.Currency = raw
	}

	if raw := required(MaxFeeKey); raw != "" {
		fee, err := strconv.ParseUint(raw, 10, 64)
		if err != nil || fee == 0 {
			problems = append(problems, MaxFeeKey+" must be a positive number of wei")
		}

		cfg.MaxFeePerGasWei = fee
	}

	if raw := required(HTTPPortKey); raw != "" {
		port, err := strconv.Atoi(raw)
		if err != nil || port < 1 || port > 65535 {
			problems = append(problems, HTTPPortKey+" must be a port from 1 to 65535")
		}

		cfg.HTTPPort = port
	}

	cfg.EpochInterval = duration(EpochIntervalKey, defaultEpochInterval, time.Minute, 24*time.Hour)
	cfg.Tick = duration(TickKey, defaultTick, time.Second, 5*time.Minute)
	cfg.SnapshotTimeout = duration(SnapshotTimeoutKey, defaultSnapshotTimeout, 10*time.Second, 30*time.Minute)

	if len(problems) > 0 {
		return Daemon{}, errors.New("configuration refused: " + strings.Join(problems, "; "))
	}

	return cfg, nil
}
