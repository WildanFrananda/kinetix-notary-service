package ignition

import (
	"context"
	"math/big"

	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark-crypto/ecc/bn254/kzg"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

const (
	SealedTranscript = "sealed/transcript00.dat"
	sanityDegree     = 64
)

func Import(ctx context.Context, client *Client, size int) (*kzg.SRS, error) {
	const op = "ignition.Import"
	if size < 2 || size-1 > client.ceremony.PointsPerTranscript {
		return nil, notary.Errf(
			notary.KindInvalidArgument,
			op,
			nil,
			"an SRS of %d points needs 2..%d",
			size,
			client.ceremony.PointsPerTranscript+1,
		)
	}

	manifest, err := client.Manifest(ctx)
	if err != nil {
		return nil, err
	}

	sealedTau, err := checkManifest(manifest, client.ceremony)
	if err != nil {
		return nil, err
	}

	openings := make([]Opening, 0, len(manifest.Participants)+1)
	for _, participant := range manifest.Participants {
		opening, err := client.Opening(ctx, participant.firstTranscript())
		if err != nil {
			return nil, err
		}

		openings = append(openings, opening)
	}

	sealed, err := client.Opening(ctx, SealedTranscript)
	if err != nil {
		return nil, err
	}

	openings = append(openings, sealed)
	if err := VerifyChain(openings); err != nil {
		return nil, err
	}

	if !sealed.Tau.Equal(&sealedTau) {
		return nil, notary.Errf(notary.KindProofRejected, op, nil, "the sealed transcript's τ is not the manifest's")
	}

	powers, err := client.Powers(ctx, SealedTranscript, size-1)
	if err != nil {
		return nil, err
	}

	if !powers.Points[0].Equal(&sealed.First) || !powers.Tau.Equal(&sealed.Tau) {
		return nil, notary.Errf(notary.KindProofRejected, op, nil, "the sealed transcript changed between two reads")
	}

	if err := VerifyPowers(powers); err != nil {
		return nil, err
	}

	srs := NewSRS(powers)
	if err := sanityCheck(srs); err != nil {
		return nil, err
	}

	return srs, nil
}

func NewSRS(powers Powers) *kzg.SRS {
	_, _, g1, g2 := bn254.Generators()

	var srs kzg.SRS

	srs.Pk.G1 = append([]bn254.G1Affine{g1}, powers.Points...)
	srs.Vk.G1 = g1
	srs.Vk.G2 = [2]bn254.G2Affine{g2, powers.Tau}
	srs.Vk.Lines[0] = bn254.PrecomputeLines(srs.Vk.G2[0])
	srs.Vk.Lines[1] = bn254.PrecomputeLines(srs.Vk.G2[1])

	return &srs
}

func checkManifest(manifest Manifest, ceremony Ceremony) (bn254.G2Affine, error) {
	const op = "ignition.checkManifest"

	var tau bn254.G2Affine
	if manifest.Name != ceremony.Name ||
		manifest.NumG1Points != ceremony.Points ||
		manifest.NumG2Points != 1 ||
		manifest.PointsPerTranscript != ceremony.PointsPerTranscript {
		return tau, notary.Errf(
			notary.KindProofRejected,
			op,
			nil,
			"the manifest describes %q with %d points, %d per transcript",
			manifest.Name,
			manifest.NumG1Points,
			manifest.PointsPerTranscript,
		)
	}

	if len(manifest.Participants) == 0 {
		return tau, notary.Errf(notary.KindProofRejected, op, nil, "the manifest names no participant")
	}

	for i := 1; i < len(manifest.Participants); i++ {
		if manifest.Participants[i].Position <= manifest.Participants[i-1].Position {
			return tau, notary.Errf(
				notary.KindProofRejected,
				op,
				nil,
				"participant %d is not after participant %d",
				manifest.Participants[i].Position,
				manifest.Participants[i-1].Position,
			)
		}
	}

	if len(manifest.Crs.T2) != 4 {
		return tau, notary.Errf(
			notary.KindProofRejected,
			op,
			nil,
			"the manifest's τ has %d coordinates",
			len(manifest.Crs.T2),
		)
	}

	for i, into := range []*fp.Element{&tau.X.A0, &tau.X.A1, &tau.Y.A0, &tau.Y.A1} {
		value, ok := new(big.Int).SetString(manifest.Crs.T2[i], 0)
		if !ok {
			return tau, notary.Errf(
				notary.KindProofRejected,
				op,
				nil,
				"the manifest's τ coordinate %d is not a number",
				i,
			)
		}

		if err := into.SetBytesCanonical(value.FillBytes(make([]byte, fp.Bytes))); err != nil {
			return tau, notary.Errf(notary.KindProofRejected, op, err, "the manifest's τ coordinate %d", i)
		}
	}

	if tau.IsInfinity() || !tau.IsInSubGroup() {
		return tau, notary.Errf(notary.KindProofRejected, op, nil, "the manifest's τ is not a point of G2")
	}

	return tau, nil
}

func sanityCheck(srs *kzg.SRS) error {
	const op = "ignition.sanityCheck"

	polynomial := make([]fr.Element, min(sanityDegree, len(srs.Pk.G1)))
	for i := range polynomial {
		if _, err := polynomial[i].SetRandom(); err != nil {
			return notary.Errf(notary.KindInvalidArgument, op, err, "drawing the polynomial")
		}
	}

	var point fr.Element
	if _, err := point.SetRandom(); err != nil {
		return notary.Errf(notary.KindInvalidArgument, op, err, "drawing the point")
	}

	digest, err := kzg.Commit(polynomial, srs.Pk)
	if err != nil {
		return notary.Errf(notary.KindProofRejected, op, err, "committing with the SRS")
	}

	opening, err := kzg.Open(polynomial, point, srs.Pk)
	if err != nil {
		return notary.Errf(notary.KindProofRejected, op, err, "opening with the SRS")
	}

	if err := kzg.Verify(&digest, &opening, point, srs.Vk); err != nil {
		return notary.Errf(notary.KindProofRejected, op, err, "a KZG opening made with the SRS does not verify")
	}

	return nil
}
