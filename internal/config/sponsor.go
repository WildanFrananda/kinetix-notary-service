package config

import (
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

const (
	PaymasterAddressKey = "KINETIX_NOTARY_PAYMASTER_ADDRESS"
	BundlerURLKey       = "KINETIX_NOTARY_BUNDLER_URL"
	HolderKeystoreKey   = "KINETIX_NOTARY_HOLDER_KEYSTORE"
	HolderUnlockFileKey = "KINETIX_NOTARY_HOLDER_PASSWORD_FILE"
	ReceiptWaitKey      = "KINETIX_NOTARY_RECEIPT_WAIT"

	defaultReceiptWait = 3 * time.Minute
	maxReceiptWait     = 30 * time.Minute
)

type Sponsor struct {
	Verifier       Verifier
	Paymaster      notary.Address
	BundlerURL     string
	HolderKeystore string
	HolderPassword string
	ReceiptWait    time.Duration
}

func LoadSponsor(lookup Lookup) (Sponsor, error) {
	var problems []string

	verifier, err := LoadVerifier(lookup)
	if err != nil {
		problems = append(problems, strings.TrimPrefix(err.Error(), "configuration refused: "))
	}

	read := requiredReader(lookup, &problems)
	cfg := Sponsor{
		Verifier:       verifier,
		Paymaster:      address(PaymasterAddressKey, read(PaymasterAddressKey), &problems),
		BundlerURL:     bundlerURL(read(BundlerURLKey), &problems),
		HolderKeystore: read(HolderKeystoreKey),
		HolderPassword: read(HolderUnlockFileKey),
		ReceiptWait:    defaultReceiptWait,
	}

	if raw, ok := lookup(ReceiptWaitKey); ok && strings.TrimSpace(raw) != "" {
		wait, err := time.ParseDuration(strings.TrimSpace(raw))
		if err != nil || wait <= 0 || wait > maxReceiptWait {
			problems = append(problems, ReceiptWaitKey+" must be a duration above zero and at most 30m")
		}

		cfg.ReceiptWait = wait
	}

	if len(problems) > 0 {
		return Sponsor{}, errors.New("configuration refused: " + strings.Join(problems, "; "))
	}

	return cfg, nil
}

func bundlerURL(raw string, problems *[]string) string {
	if raw == "" {
		return ""
	}

	parsed, err := url.Parse(raw)
	if err != nil || (parsed.Scheme != "https" && parsed.Scheme != "http") || parsed.Host == "" {
		*problems = append(*problems, BundlerURLKey+" must be an http or https URL with a host")
		return ""
	}

	return raw
}
