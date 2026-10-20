package ethereum

import (
	"context"
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type ProofWriter struct {
	*transactor
	*ProofReader
}

var _ notary.ProofChain = (*ProofWriter)(nil)

func NewProofWriter(client *ethclient.Client, settings Settings, key *ecdsa.PrivateKey) *ProofWriter {
	return &ProofWriter{
		transactor: newTransactor(
			client,
			settings.Ledger.ChainID,
			settings.Ledger.Address,
			key,
			settings.Timeout,
			"a proof for this ledger",
		),
		ProofReader: NewProofReader(client, settings),
	}
}

func (w *ProofWriter) PrepareProof(
	ctx context.Context,
	header notary.EpochHeader,
	plan notary.ConfidentialPlan,
	nonce uint64,
	fees notary.FeeQuote,
) (notary.PreparedAnchor, error) {
	const op = "ethereum.ProofWriter.PrepareProof"

	calldata := w.binding.PackProveEpoch(
		header.Number,
		plan.TreeRoot,
		plan.TotalCommitment,
		header.LeafCount,
		plan.Proof,
	)

	estimate, err := w.estimate(ctx, calldata, fees)
	if err != nil {
		return notary.PreparedAnchor{}, refusedBy(
			op,
			proofRefusals,
			fmt.Sprintf("the ledger refused the proof of epoch %d", header.Number),
			err,
		)
	}

	return w.sign(op, nonce, estimate*gasMarginPercent/100, calldata, fees)
}
