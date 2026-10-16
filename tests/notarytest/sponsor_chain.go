package notarytest

import (
	"context"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type SponsorChain struct {
	Delegate        notary.Address
	AccountNonces   uint64
	OperationNonces uint64
	MaxFee          uint64
	Priority        uint64
	Now             notary.UnixTime
	Err             error
}

var _ notary.SponsorChain = (*SponsorChain)(nil)

func (c *SponsorChain) OperationNonce(context.Context, notary.Address, notary.Address) (uint64, error) {
	return c.OperationNonces, c.Err
}

func (c *SponsorChain) AccountNonce(context.Context, notary.Address) (uint64, error) {
	return c.AccountNonces, c.Err
}

func (c *SponsorChain) DelegateOf(context.Context, notary.Address) (notary.Address, error) {
	return c.Delegate, c.Err
}

func (c *SponsorChain) Fees(context.Context) (uint64, uint64, error) {
	return c.MaxFee, c.Priority, c.Err
}

func (c *SponsorChain) BlockTime(context.Context) (notary.UnixTime, error) {
	return c.Now, c.Err
}
