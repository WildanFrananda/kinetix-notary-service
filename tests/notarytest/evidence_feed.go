package notarytest

import (
	"context"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type EvidenceFeed struct {
	Kept    notary.EvidenceStream
	Records []notary.EvidenceRecord
	MaxPage int32
	Err     error
	Tamper  func(*notary.EvidencePage)
	Asked   []notary.LogPosition
}

var _ notary.EvidenceFeed = (*EvidenceFeed)(nil)

func (f *EvidenceFeed) Stream() notary.EvidenceStream { return f.Kept }

func (f *EvidenceFeed) Since(_ context.Context, after notary.LogPosition, limit int32) (notary.EvidencePage, error) {
	f.Asked = append(f.Asked, after)
	if f.Err != nil {
		return notary.EvidencePage{}, f.Err
	}

	if f.MaxPage > 0 {
		limit = min(limit, f.MaxPage)
	}

	page := notary.EvidencePage{Next: after}

	for _, record := range f.Records {
		if !record.Position().After(after) {
			continue
		}

		if int32(len(page.Records)) == limit {
			page.HasMore = true
			break
		}

		page.Records = append(page.Records, record)
		page.Next = record.Position()
	}

	if f.Tamper != nil {
		f.Tamper(&page)
	}

	return page, nil
}
