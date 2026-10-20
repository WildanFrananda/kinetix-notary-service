package notary

import "context"

type EpochProver interface {
	Prove(ctx context.Context, request ProofRequest) ([]byte, error)
}
