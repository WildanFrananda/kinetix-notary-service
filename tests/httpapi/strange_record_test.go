package httpapi_test

import "github.com/WildanFrananda/kinetix-notary-service/internal/notary"

type strangeRecord struct{}

func (strangeRecord) Stream() notary.EvidenceStream {
	return notary.StreamOrderLifecycle
}

func (strangeRecord) Position() notary.LogPosition {
	return notary.LogPosition{Transaction: 1, Sequence: 1}
}

func (strangeRecord) Parties() []string {
	return []string{"principal-0005"}
}
