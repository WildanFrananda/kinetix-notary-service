package notarytest

import (
	"context"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

var CouncilIdentity = notary.CouncilIdentity{ChainID: 84532, Address: notary.Address{0xc0, 0x0c}}

type Council struct {
	Answers        notary.CouncilIdentity
	Ledger         notary.Address
	Members        map[notary.Address]bool
	ThresholdValue uint64
	Anchors        bool
	Scheme         notary.ApprovalScheme
	Err            error
}

var _ notary.Council = (*Council)(nil)

func (c *Council) Identity(context.Context) (notary.CouncilIdentity, error) {
	return c.Answers, c.Err
}

func (c *Council) IsMember(_ context.Context, member notary.Address) (bool, error) {
	return c.Members[member], c.Err
}

func (c *Council) Threshold(context.Context) (uint64, error) {
	return c.ThresholdValue, c.Err
}

func (c *Council) ApprovalDigest(_ context.Context, number uint64, epoch notary.Digest) (notary.Digest, error) {
	return c.Scheme.ApprovalDigest(c.Answers, c.Ledger, number, epoch), c.Err
}

func (c *Council) AnchorsLedger(context.Context) (bool, error) {
	return c.Anchors, c.Err
}
