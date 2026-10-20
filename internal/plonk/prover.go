package plonk

import (
	"io"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/kzg"
	"github.com/consensys/gnark/backend"
	gnarkplonk "github.com/consensys/gnark/backend/plonk"
	plonkbn254 "github.com/consensys/gnark/backend/plonk/bn254"
	"github.com/consensys/gnark/backend/solidity"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/scs"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

const solidityVersion = "0.8.37"

type Prover struct {
	depth  int
	system constraint.ConstraintSystem
	key    gnarkplonk.ProvingKey
	vk     gnarkplonk.VerifyingKey
}

func Compile(depth int) (constraint.ConstraintSystem, error) {
	system, err := frontend.Compile(ecc.BN254.ScalarField(), scs.NewBuilder, NewLiabilityCircuit(depth))
	if err != nil {
		return nil, notary.Errf(notary.KindInvalidArgument, "plonk.Compile", err, "the circuit of depth %d", depth)
	}

	return system, nil
}

func NewProver(depth int, system constraint.ConstraintSystem, canonical, lagrange kzg.SRS) (*Prover, error) {
	key, vk, err := gnarkplonk.Setup(system, canonical, lagrange)
	if err != nil {
		return nil, notary.Errf(notary.KindInvalidArgument, "plonk.NewProver", err, "the setup of depth %d", depth)
	}

	return &Prover{depth: depth, system: system, key: key, vk: vk}, nil
}

func (p *Prover) Depth() int {
	return p.depth
}

func (p *Prover) Prove(assignment *LiabilityCircuit) ([]byte, error) {
	const op = "plonk.Prover.Prove"

	full, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	if err != nil {
		return nil, notary.Errf(notary.KindInvalidArgument, op, err, "the witness")
	}

	proof, err := gnarkplonk.Prove(p.system, p.key, full, solidity.WithProverTargetSolidityVerifier(backend.PLONK))
	if err != nil {
		return nil, notary.Errf(notary.KindInvalidArgument, op, err, "the epoch does not satisfy the circuit")
	}

	public, err := full.Public()
	if err != nil {
		return nil, notary.Errf(notary.KindInvalidArgument, op, err, "the public witness")
	}

	if err := gnarkplonk.Verify(
		proof,
		p.vk,
		public,
		solidity.WithVerifierTargetSolidityVerifier(backend.PLONK),
	); err != nil {
		return nil, notary.Errf(notary.KindDisagreement, op, err, "a proof just made does not verify")
	}

	concrete, ok := proof.(*plonkbn254.Proof)
	if !ok {
		return nil, notary.Errf(notary.KindDisagreement, op, nil, "the proof is not a BN254 PLONK proof")
	}

	return concrete.MarshalSolidity(), nil
}

func (p *Prover) ExportVerifier(w io.Writer) error {
	if err := p.vk.ExportSolidity(w, solidity.WithPragmaVersion(solidityVersion)); err != nil {
		return notary.Errf(notary.KindInvalidArgument, "plonk.Prover.ExportVerifier", err, "")
	}

	return nil
}

func (p *Prover) WriteVerifyingKey(w io.Writer) error {
	if _, err := p.vk.WriteTo(w); err != nil {
		return notary.Errf(notary.KindInvalidArgument, "plonk.Prover.WriteVerifyingKey", err, "")
	}

	return nil
}
