package notary

import "context"

type Paymaster interface {
	Identity(ctx context.Context) (PaymasterIdentity, error)
	Funds(ctx context.Context) (PaymasterFunds, error)
	Sponsorship(ctx context.Context, sender Address, callData []byte, paymasterData []byte) (SponsorshipWindow, error)
}
