package ignition

import "github.com/consensys/gnark-crypto/ecc/bn254"

type Powers struct {
	Points []bn254.G1Affine
	Tau    bn254.G2Affine
}
