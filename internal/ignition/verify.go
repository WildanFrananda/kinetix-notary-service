package ignition

import (
	"runtime"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

func VerifyChain(openings []Opening) error {
	const op = "ignition.VerifyChain"

	_, _, g1, g2 := bn254.Generators()
	previous := g1

	for i, opening := range openings {
		follows, err := sameRatio(opening.First, previous, g2, opening.Contribution)
		if err != nil {
			return err
		}

		if !follows {
			return notary.Errf(
				notary.KindProofRejected,
				op,
				nil,
				"contribution %d is not built on contribution %d",
				i,
				i-1,
			)
		}

		consistent, err := sameRatio(opening.First, g1, g2, opening.Tau)
		if err != nil {
			return err
		}

		if !consistent {
			return notary.Errf(notary.KindProofRejected, op, nil, "contribution %d's τ in G2 is not its τ in G1", i)
		}

		previous = opening.First
	}

	return nil
}

func VerifyPowers(powers Powers) error {
	const op = "ignition.VerifyPowers"

	_, _, g1, g2 := bn254.Generators()
	sequence := append([]bn254.G1Affine{g1}, powers.Points...)

	scalars := make([]fr.Element, len(sequence)-1)
	for i := range scalars {
		if _, err := scalars[i].SetRandom(); err != nil {
			return notary.Errf(notary.KindInvalidArgument, op, err, "drawing the challenge")
		}
	}

	config := ecc.MultiExpConfig{NbTasks: runtime.NumCPU()}

	var lower, upper bn254.G1Affine
	if _, err := lower.MultiExp(sequence[:len(sequence)-1], scalars, config); err != nil {
		return notary.Errf(notary.KindInvalidArgument, op, err, "the lower combination")
	}

	if _, err := upper.MultiExp(sequence[1:], scalars, config); err != nil {
		return notary.Errf(notary.KindInvalidArgument, op, err, "the upper combination")
	}

	powersOfTau, err := sameRatio(lower, upper, powers.Tau, g2)
	if err != nil {
		return err
	}

	if !powersOfTau {
		return notary.Errf(
			notary.KindProofRejected,
			op,
			nil,
			"the %d points are not successive powers of τ",
			len(powers.Points),
		)
	}

	return nil
}

func sameRatio(a1, b1 bn254.G1Affine, a2, b2 bn254.G2Affine) (bool, error) {
	var negated bn254.G2Affine
	negated.Neg(&a2)

	ok, err := bn254.PairingCheck([]bn254.G1Affine{a1, b1}, []bn254.G2Affine{negated, b2})
	if err != nil {
		return false, notary.Errf(notary.KindInvalidArgument, "ignition.sameRatio", err, "the pairing")
	}

	return ok, nil
}
