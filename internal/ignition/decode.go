package ignition

import (
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

func decodeElement(data []byte, into *fp.Element) error {
	var canonical [fp.Bytes]byte
	for word := range 4 {
		copy(canonical[(3-word)*8:(4-word)*8], data[word*8:(word+1)*8])
	}

	if err := into.SetBytesCanonical(canonical[:]); err != nil {
		return notary.Errf(
			notary.KindProofRejected,
			"ignition.decodeElement",
			err,
			"a coordinate is not a field element",
		)
	}

	return nil
}

func decodeG1(data []byte) (bn254.G1Affine, error) {
	var point bn254.G1Affine
	if err := decodeElement(data[0:32], &point.X); err != nil {
		return point, err
	}

	if err := decodeElement(data[32:64], &point.Y); err != nil {
		return point, err
	}

	if point.IsInfinity() || !point.IsInSubGroup() {
		return point, notary.Errf(
			notary.KindProofRejected,
			"ignition.decodeG1",
			nil,
			"a G1 point is not a point of the group",
		)
	}

	return point, nil
}

func decodeG2(data []byte) (bn254.G2Affine, error) {
	var point bn254.G2Affine
	for i, into := range []*fp.Element{&point.X.A0, &point.X.A1, &point.Y.A0, &point.Y.A1} {
		if err := decodeElement(data[i*32:(i+1)*32], into); err != nil {
			return point, err
		}
	}

	if point.IsInfinity() || !point.IsInSubGroup() {
		return point, notary.Errf(
			notary.KindProofRejected,
			"ignition.decodeG2",
			nil,
			"a G2 point is not a point of the group",
		)
	}

	return point, nil
}
