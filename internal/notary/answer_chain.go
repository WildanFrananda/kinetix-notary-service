package notary

import (
	"context"
	"math/big"
)

type AnswerChain interface {
	Answerer() Address
	Answer(ctx context.Context, answer Answer) (Digest, error)
	Balance(ctx context.Context) (*big.Int, error)
}
