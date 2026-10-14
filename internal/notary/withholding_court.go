package notary

import "context"

type WithholdingCourt interface {
	Identity(ctx context.Context) (CourtIdentity, error)
	RegistrationDigest(ctx context.Context, registration Registration) (Digest, error)
	LatestBlock(ctx context.Context) (uint64, error)
	Requests(ctx context.Context, from, to uint64) ([]WithholdingRequest, error)
	RequestOf(ctx context.Context, registration Digest, epoch uint64) (RequestState, UnixTime, error)
	TotalSilences(ctx context.Context) (uint64, error)
}
