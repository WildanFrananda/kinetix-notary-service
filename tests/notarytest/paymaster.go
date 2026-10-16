package notarytest

import (
	"context"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type Paymaster struct {
	IdentityValue notary.PaymasterIdentity
	FundsValue    notary.PaymasterFunds
	Window        notary.SponsorshipWindow
	Refusal       error
	Err           error
	Asked         []notary.Address
}

var _ notary.Paymaster = (*Paymaster)(nil)

func (p *Paymaster) Identity(context.Context) (notary.PaymasterIdentity, error) {
	return p.IdentityValue, p.Err
}

func (p *Paymaster) Funds(context.Context) (notary.PaymasterFunds, error) {
	return p.FundsValue, p.Err
}

func (p *Paymaster) Sponsorship(
	_ context.Context,
	sender notary.Address,
	_ []byte,
	_ []byte,
) (notary.SponsorshipWindow, error) {
	p.Asked = append(p.Asked, sender)

	if p.Refusal != nil {
		return notary.SponsorshipWindow{}, p.Refusal
	}

	return p.Window, p.Err
}
