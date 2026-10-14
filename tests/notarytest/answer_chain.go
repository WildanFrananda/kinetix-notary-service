package notarytest

import (
	"context"
	"math/big"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type AnswerChain struct {
	Address      notary.Address
	Answers      []notary.Answer
	Err          error
	BalanceValue *big.Int
}

var _ notary.AnswerChain = (*AnswerChain)(nil)

func (c *AnswerChain) Answerer() notary.Address {
	return c.Address
}

func (c *AnswerChain) Answer(_ context.Context, answer notary.Answer) (notary.Digest, error) {
	if c.Err != nil {
		return notary.Digest{}, c.Err
	}

	c.Answers = append(c.Answers, answer)

	return notary.Digest{byte(len(c.Answers))}, nil
}

func (c *AnswerChain) Balance(context.Context) (*big.Int, error) {
	return c.BalanceValue, nil
}
