package notary

import "context"

type WithholdingScans interface {
	ScanFrom(ctx context.Context, court CourtIdentity) (uint64, bool, error)
	OpenRequests(ctx context.Context, court CourtIdentity) ([]WithholdingRequest, error)
	RecordScan(ctx context.Context, court CourtIdentity, next uint64, found []WithholdingRequest) error
	CloseRequest(ctx context.Context, court CourtIdentity, registration Digest, epoch uint64) error
}
