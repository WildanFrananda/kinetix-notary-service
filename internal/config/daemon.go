package config

import (
	"errors"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

const (
	DatabaseURLKey     = "KINETIX_NOTARY_DATABASE_URL"
	AnchorerKeyKey     = "KINETIX_NOTARY_ANCHORER_PRIVATE_KEY"
	StatementKeysKey   = "KINETIX_NOTARY_STATEMENT_PRIVATE_KEYS"
	WithholdingKey     = "KINETIX_NOTARY_WITHHOLDING_COURT_ADDRESS"
	WithholdingFromKey = "KINETIX_NOTARY_WITHHOLDING_FROM_BLOCK"
	LogRangeKey        = "KINETIX_NOTARY_RPC_LOG_RANGE"
	AnswererKeyKey     = "KINETIX_NOTARY_ANSWERER_PRIVATE_KEY"
	PaymentAddressKey  = "KINETIX_NOTARY_PAYMENT_GRPC_ADDR"
	WatchAddressKey    = "KINETIX_NOTARY_WATCH_GRPC_ADDR"
	PKIDirKey          = "KINETIX_PKI_DIR"
	TrustDomainKey     = "KINETIX_TRUST_DOMAIN"
	CurrencyKey        = "KINETIX_NOTARY_CURRENCY"
	MaxFeeKey          = "KINETIX_NOTARY_MAX_FEE_PER_GAS_WEI"
	HTTPPortKey        = "KINETIX_NOTARY_HTTP_PORT"
	EpochIntervalKey   = "KINETIX_NOTARY_EPOCH_INTERVAL"
	TickKey            = "KINETIX_NOTARY_TICK"
	SnapshotTimeoutKey = "KINETIX_NOTARY_SNAPSHOT_TIMEOUT"
	JWKSURLKey         = "IDENTITY_JWKS_URL"
	JWTIssuerKey       = "JWT_ISSUER"
	JWTAudienceKey     = "JWT_AUDIENCE"
	OrderAddressKey    = "KINETIX_NOTARY_ORDER_GRPC_ADDR"
	MatchingAddressKey = "KINETIX_NOTARY_MATCHING_GRPC_ADDR"
	ReviewAddressKey   = "KINETIX_NOTARY_REVIEW_GRPC_ADDR"
	EvidenceKeyKey     = "KINETIX_NOTARY_EVIDENCE_ANCHORER_PRIVATE_KEY"
	EvidenceBatchKey   = "KINETIX_NOTARY_EVIDENCE_BATCH_INTERVAL"
	FeedTimeoutKey     = "KINETIX_NOTARY_FEED_TIMEOUT"
	SchemeKey          = "KINETIX_NOTARY_SCHEME"
	SRSPathKey         = "KINETIX_NOTARY_SRS_PATH"
	SRSSHA256Key       = "KINETIX_NOTARY_SRS_SHA256"
	NotificationKey    = "KINETIX_NOTARY_NOTIFICATION_GRPC_ADDR"
	NotifyTimeoutKey   = "KINETIX_NOTARY_NOTIFY_TIMEOUT"
	VerifyURLKey       = "KINETIX_NOTARY_STATEMENT_VERIFY_URL"
	StatementDayKey    = "KINETIX_NOTARY_STATEMENT_UTC_OFFSET"

	defaultTrustDomain     = "kinetix.local"
	defaultEpochInterval   = time.Hour
	defaultTick            = 15 * time.Second
	defaultSnapshotTimeout = 3 * time.Minute
	defaultEvidenceBatch   = 24 * time.Hour
	defaultFeedTimeout     = 30 * time.Second
	defaultNotifyTimeout   = 30 * time.Second
	maxUTCOffset           = 14 * time.Hour
)

type Daemon struct {
	Chain           Chain
	DatabaseURL     Secret
	AnchorerKey     Secret
	StatementKeys   Secret
	Withholding     notary.Address
	WithholdingFrom uint64
	LogRange        uint64
	AnswererKey     Secret
	Council         notary.Address
	Directory       notary.Address
	Paymaster       notary.Address
	WatchAddress    string
	PaymentAddress  string
	PKIDir          string
	TrustDomains    []string
	Currency        string
	MaxFeePerGasWei uint64
	HTTPPort        int
	EpochInterval   time.Duration
	Tick            time.Duration
	SnapshotTimeout time.Duration
	JWKSURL         string
	JWTIssuer       string
	JWTAudience     string
	OrderAddress    string
	MatchingAddress string
	ReviewAddress   string
	Registry        notary.Address
	EvidenceKey     Secret
	EvidenceBatch   time.Duration
	Scheme          uint16
	SRSPath         string
	SRSSHA256       string
	FeedTimeout     time.Duration
	Notification    string
	NotifyTimeout   time.Duration
	VerifyURL       string
	UTCOffset       time.Duration
}

func LoadDaemon(lookup Lookup) (Daemon, error) {
	var problems []string

	chain, err := LoadChain(lookup)
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

	if raw := required(StatementKeysKey); raw != "" {
		cfg.StatementKeys = Secret{value: raw}
	}

	cfg.Withholding = address(WithholdingKey, required(WithholdingKey), &problems)

	if raw := required(WithholdingFromKey); raw != "" {
		block, err := strconv.ParseUint(raw, 10, 64)
		if err != nil {
			problems = append(problems, WithholdingFromKey+" must be the block the withholding court was deployed in")
		}

		cfg.WithholdingFrom = block
	}

	if raw := required(LogRangeKey); raw != "" {
		blocks, err := strconv.ParseUint(raw, 10, 64)
		if err != nil || blocks == 0 {
			problems = append(problems, LogRangeKey+" must be the most blocks one eth_getLogs may span")
		}

		cfg.LogRange = blocks
	}

	if raw := required(AnswererKeyKey); raw != "" {
		cfg.AnswererKey = Secret{value: raw}
	}

	hostPort := func(key string) string {
		raw := required(key)
		if raw == "" {
			return raw
		}

		host, port, err := net.SplitHostPort(raw)
		if err != nil || host == "" || port == "" {
			problems = append(problems, key+" must be host:port")
		}

		return raw
	}

	cfg.PaymentAddress = hostPort(PaymentAddressKey)
	cfg.OrderAddress = hostPort(OrderAddressKey)
	cfg.MatchingAddress = hostPort(MatchingAddressKey)
	cfg.ReviewAddress = hostPort(ReviewAddressKey)
	cfg.Notification = hostPort(NotificationKey)
	cfg.Registry = address(EvidenceRegistryKey, required(EvidenceRegistryKey), &problems)

	if raw := required(EvidenceKeyKey); raw != "" {
		cfg.EvidenceKey = Secret{value: raw}
	}

	cfg.Council = address(CouncilAddressKey, required(CouncilAddressKey), &problems)
	cfg.Directory = address(DirectoryAddressKey, required(DirectoryAddressKey), &problems)
	cfg.Paymaster = address(PaymasterAddressKey, required(PaymasterAddressKey), &problems)

	cfg.WatchAddress = hostPort(WatchAddressKey)

	cfg.PKIDir = required(PKIDirKey)
	cfg.TrustDomains = trustDomains(lookup)

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
	cfg.EvidenceBatch = duration(EvidenceBatchKey, defaultEvidenceBatch, time.Minute, 7*24*time.Hour)
	cfg.FeedTimeout = duration(FeedTimeoutKey, defaultFeedTimeout, time.Second, 5*time.Minute)
	cfg.NotifyTimeout = duration(NotifyTimeoutKey, defaultNotifyTimeout, time.Second, 5*time.Minute)

	cfg.Scheme = 1

	if raw := optional(SchemeKey); raw != "" {
		switch raw {
		case "1":
			cfg.Scheme = 1
		case "2":
			cfg.Scheme = 2
		default:
			problems = append(problems, SchemeKey+" must be 1 or 2")
		}
	}

	if cfg.Scheme == 2 {
		cfg.SRSPath = required(SRSPathKey)
		if raw := required(SRSSHA256Key); raw != "" {
			if len(raw) != 64 || strings.ContainsFunc(raw, func(r rune) bool {
				return (r < '0' || r > '9') && (r < 'a' || r > 'f')
			}) {
				problems = append(problems, SRSSHA256Key+" must be 64 lowercase hex digits")
			}

			cfg.SRSSHA256 = raw
		}
	}

	if raw := required(JWKSURLKey); raw != "" {
		parsed, err := url.Parse(raw)
		if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
			problems = append(problems, JWKSURLKey+" must be an http or https URL with a host")
		}

		cfg.JWKSURL = raw
	}

	cfg.JWTIssuer = required(JWTIssuerKey)
	cfg.JWTAudience = required(JWTAudienceKey)

	if raw := required(VerifyURLKey); raw != "" {
		parsed, err := url.Parse(raw)
		if err != nil || parsed.Scheme != "https" || parsed.Host == "" {
			problems = append(problems, VerifyURLKey+" must be an https URL with a host")
		}

		cfg.VerifyURL = raw
	}

	if raw := required(StatementDayKey); raw != "" {
		offset, err := utcOffset(raw)
		if err != nil {
			problems = append(problems, StatementDayKey+" must be an offset from UTC such as +07:00, at most 14 hours")
		}

		cfg.UTCOffset = offset
	}

	if len(problems) > 0 {
		return Daemon{}, errors.New("configuration refused: " + strings.Join(problems, "; "))
	}

	return cfg, nil
}

func utcOffset(raw string) (time.Duration, error) {
	parsed, err := time.Parse("-07:00", raw)
	if err != nil {
		return 0, errors.New("not an offset such as +07:00")
	}

	_, seconds := parsed.Zone()
	offset := time.Duration(seconds) * time.Second

	if offset > maxUTCOffset || offset < -maxUTCOffset {
		return 0, errors.New("the offset is beyond 14 hours")
	}

	return offset, nil
}
