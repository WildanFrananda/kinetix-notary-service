package notary

import "context"

type EvidenceFeed interface {
	Stream() EvidenceStream
	Since(ctx context.Context, after LogPosition, limit int32) (EvidencePage, error)
}
