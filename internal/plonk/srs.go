package plonk

import (
	"crypto/sha256"
	"encoding/hex"

	kzgbn254 "github.com/consensys/gnark-crypto/ecc/bn254/kzg"
	"github.com/consensys/gnark/constraint"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

const openingPoints = 3

func Cardinality(system constraint.ConstraintSystem) int {
	needed := system.GetNbConstraints() + system.GetNbPublicVariables()

	size := 1
	for size < needed {
		size <<= 1
	}

	return size
}

func SRSSize(system constraint.ConstraintSystem) int {
	return Cardinality(system) + openingPoints
}

func NewLagrange(canonical *kzgbn254.SRS, system constraint.ConstraintSystem) (*kzgbn254.SRS, error) {
	const op = "plonk.NewLagrange"

	size := Cardinality(system)
	if len(canonical.Pk.G1) < size+openingPoints {
		return nil, notary.Errf(
			notary.KindInvalidArgument,
			op,
			nil,
			"an SRS of %d points is too small for a domain of %d",
			len(canonical.Pk.G1),
			size,
		)
	}

	points, err := kzgbn254.ToLagrangeG1(canonical.Pk.G1[:size])
	if err != nil {
		return nil, notary.Errf(notary.KindInvalidArgument, op, err, "the Lagrange basis")
	}

	return &kzgbn254.SRS{Pk: kzgbn254.ProvingKey{G1: points}, Vk: canonical.Vk}, nil
}

func CircuitSHA256(system constraint.ConstraintSystem) (string, error) {
	hash := sha256.New()
	if _, err := system.WriteTo(hash); err != nil {
		return "", notary.Errf(notary.KindInvalidArgument, "plonk.CircuitSHA256", err, "encoding the circuit")
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
