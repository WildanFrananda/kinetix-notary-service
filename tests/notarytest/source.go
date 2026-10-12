package notarytest

import (
	"context"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type Source struct {
	Next  func() (notary.Snapshot, error)
	Asked []string
}

var _ notary.SnapshotSource = (*Source)(nil)

func (s *Source) Snapshot(_ context.Context, currency string) (notary.Snapshot, error) {
	s.Asked = append(s.Asked, currency)
	return s.Next()
}
