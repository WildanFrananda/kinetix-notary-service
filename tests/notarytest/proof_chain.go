package notarytest

import (
	"context"
	"crypto/sha256"
	"encoding/binary"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type statementKey struct {
	number    uint64
	root      notary.Digest
	leafCount uint64
}

type proofTx struct {
	key   statementKey
	proof []byte
}

var _ notary.ProofChain = (*Chain)(nil)

func keyOf(header notary.EpochHeader) statementKey {
	return statementKey{number: header.Number, root: header.Root, leafCount: header.LeafCount}
}

func (c *Chain) provenAt(key statementKey, block uint64) bool {
	at, ok := c.proven[key]
	return ok && at <= block
}

func (c *Chain) Depth(context.Context) (int, error) {
	return c.DepthValue, nil
}

func (c *Chain) IsProven(_ context.Context, number uint64, root notary.Digest, leafCount uint64) (bool, error) {
	return c.provenAt(statementKey{number: number, root: root, leafCount: leafCount}, c.block), nil
}

func (c *Chain) PrepareProof(
	_ context.Context,
	header notary.EpochHeader,
	plan notary.ConfidentialPlan,
	nonce uint64,
	fees notary.FeeQuote,
) (notary.PreparedAnchor, error) {
	const op = "notarytest.Chain.PrepareProof"

	key := keyOf(header)
	switch {
	case c.provenAt(key, c.block):
		return notary.PreparedAnchor{}, notary.Errf(notary.KindAlreadyAnchored, op, nil, "StatementAlreadyProven")
	case c.RefuseProofs || len(plan.Proof) == 0:
		return notary.PreparedAnchor{}, notary.Errf(notary.KindAnchorRefused, op, nil, "ProofRejected")
	}

	return c.signProof(proofTx{key: key, proof: plan.Proof}, nonce, 450_000, fees), nil
}

func (c *Chain) signProof(proof proofTx, nonce, gas uint64, fees notary.FeeQuote) notary.PreparedAnchor {
	raw := append([]byte("proof"), proof.key.root[:]...)
	for _, v := range []uint64{proof.key.number, proof.key.leafCount, nonce, gas, fees.MaxFeePerGas, fees.MaxPriorityFeePerGas} {
		raw = binary.BigEndian.AppendUint64(raw, v)
	}

	prepared := notary.PreparedAnchor{
		Hash:     notary.Digest(sha256.Sum256(raw)),
		Raw:      raw,
		Nonce:    nonce,
		GasLimit: gas,
		Fees:     fees,
	}
	c.proofTxs[prepared.Hash] = proof

	return prepared
}
