package ignition

import "github.com/consensys/gnark-crypto/ecc/bn254"

type Opening struct {
	First        bn254.G1Affine
	Tau          bn254.G2Affine
	Contribution bn254.G2Affine
}
