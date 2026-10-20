package ignition_test

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/blake2b"

	"github.com/WildanFrananda/kinetix-notary-service/internal/ignition"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/plonk"
	"github.com/WildanFrananda/kinetix-notary-service/internal/poseidon"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

const perTranscript = 40

type contributor struct {
	address string
	secret  int64
}

type ceremonyFiles struct {
	ceremony ignition.Ceremony
	manifest map[string]any
	files    map[string][]byte
	tau      *big.Int
	status   map[string]int
}

func TestAnHonestCeremonyBecomesAnSRSOfPowersOfItsTau(t *testing.T) {
	c := newCeremony(t, perTranscript, []contributor{{"0xA1", 3}, {"0xB2", 5}, {"0xC3", 7}}, 11)

	srs, err := ignition.Import(context.Background(), c.client(t), 17)
	require.NoError(t, err)
	require.Len(t, srs.Pk.G1, 17)

	_, _, g1, g2 := bn254.Generators()
	require.True(t, srs.Pk.G1[0].Equal(&g1))
	require.True(t, srs.Vk.G2[0].Equal(&g2))

	power := new(big.Int).Set(c.tau)

	for k := 1; k < 17; k++ {
		var want bn254.G1Affine
		want.ScalarMultiplication(&g1, power)
		require.True(t, srs.Pk.G1[k].Equal(&want), "power %d", k)

		power.Mul(power, c.tau).Mod(power, fr.Modulus())
	}
}

func TestAnImportedSRSProvesAnEpoch(t *testing.T) {
	system, err := plonk.Compile(0)
	require.NoError(t, err)

	size := plonk.SRSSize(system)
	c := newCeremony(t, size, []contributor{{"0xA1", 3}, {"0xB2", 5}}, 11)

	srs, err := ignition.Import(context.Background(), c.client(t), size)
	require.NoError(t, err)

	lagrange, err := plonk.NewLagrange(srs, system)
	require.NoError(t, err)

	prover, err := plonk.NewProver(0, system, srs, lagrange)
	require.NoError(t, err)

	scheme := poseidon.NewScheme()
	entries := notarytest.Entries(5, 1)

	tree, err := notary.BuildConfidentialTree(scheme, 4, 0, entries)
	require.NoError(t, err)

	blinding := notary.Blinding{0x11}

	header, err := tree.Header(scheme, notary.Digest{}, notarytest.Unix(t, 1_790_000_000), blinding)
	require.NoError(t, err)

	assignment, err := plonk.NewAssignment(4, 0, entries, header, blinding)
	require.NoError(t, err)

	proof, err := prover.Prove(assignment)
	require.NoError(t, err)
	require.NotEmpty(t, proof)
}

func TestACeremonyThatDoesNotHoldIsRefused(t *testing.T) {
	cases := map[string]func(*testing.T, *ceremonyFiles){
		"a contribution not built on the previous one": func(t *testing.T, c *ceremonyFiles) {
			c.rewriteG2(t, "002_0xb2/transcript00.dat", 1, big.NewInt(6))
		},
		"a τ in G2 that is not the τ in G1": func(t *testing.T, c *ceremonyFiles) {
			c.rewriteG2(t, "001_0xa1/transcript00.dat", 0, big.NewInt(4))
		},
		"a sealed τ that is not the manifest's": func(t *testing.T, c *ceremonyFiles) {
			c.manifest["crs"] = map[string]any{"t2": g2Hex(big.NewInt(2))}
		},
		"a point that is not the next power": func(t *testing.T, c *ceremonyFiles) {
			c.rewriteG1(t, ignition.SealedTranscript, 5, big.NewInt(9))
		},
		"a checksum that does not match": func(t *testing.T, c *ceremonyFiles) {
			file := c.files[ignition.SealedTranscript]
			file[len(file)-1] ^= 0x01
		},
		"a coordinate aliased by the modulus": func(t *testing.T, c *ceremonyFiles) {
			c.aliasFirstCoordinate(t, "001_0xa1/transcript00.dat")
		},
		"a point at infinity": func(t *testing.T, c *ceremonyFiles) {
			copy(c.files["002_0xb2/transcript00.dat"][28:92], make([]byte, 64))
		},
		"another ceremony's manifest": func(t *testing.T, c *ceremonyFiles) {
			c.manifest["name"] = "TINY_TEST_5"
		},
		"participants out of order": func(t *testing.T, c *ceremonyFiles) {
			participants, ok := c.manifest["participants"].([]map[string]any)
			require.True(t, ok)

			participants[0], participants[1] = participants[1], participants[0]
		},
		"a transcript the manifest names but nobody published": func(t *testing.T, c *ceremonyFiles) {
			delete(c.files, "002_0xb2/transcript00.dat")
		},
		"a transcript of another size": func(t *testing.T, c *ceremonyFiles) {
			binary.BigEndian.PutUint32(c.files[ignition.SealedTranscript][16:20], perTranscript+1)
		},
	}

	for name, tamper := range cases {
		t.Run(name, func(t *testing.T) {
			c := newCeremony(t, perTranscript, []contributor{{"0xA1", 3}, {"0xB2", 5}}, 11)
			tamper(t, c)

			_, err := ignition.Import(context.Background(), c.client(t), 17)
			require.Equal(t, notary.KindProofRejected, notary.KindOf(err), "%v", err)
		})
	}
}

func TestAnOutageIsNeverReadAsARefusal(t *testing.T) {
	cases := map[string]func(*testing.T, *ceremonyFiles){
		"a transcript the server cannot serve": func(t *testing.T, c *ceremonyFiles) {
			c.status["002_0xb2/transcript00.dat"] = http.StatusServiceUnavailable
		},
		"a sealed transcript cut short": func(t *testing.T, c *ceremonyFiles) {
			file := c.files[ignition.SealedTranscript]
			c.files[ignition.SealedTranscript] = file[:len(file)-10]
		},
	}

	for name, tamper := range cases {
		t.Run(name, func(t *testing.T) {
			c := newCeremony(t, perTranscript, []contributor{{"0xA1", 3}, {"0xB2", 5}}, 11)
			tamper(t, c)

			_, err := ignition.Import(context.Background(), c.client(t), 17)
			require.Equal(t, notary.KindSourceUnavailable, notary.KindOf(err), "%v", err)
		})
	}
}

func TestAnSRSLargerThanOneTranscriptIsRefused(t *testing.T) {
	c := newCeremony(t, perTranscript, []contributor{{"0xA1", 3}}, 11)

	_, err := ignition.Import(context.Background(), c.client(t), perTranscript+2)
	require.Equal(t, notary.KindInvalidArgument, notary.KindOf(err))
}

func newCeremony(t *testing.T, points int, contributors []contributor, seal int64) *ceremonyFiles {
	t.Helper()

	c := &ceremonyFiles{
		ceremony: ignition.Ceremony{
			Name:                "TINY TEST",
			Points:              2 * points,
			PointsPerTranscript: points,
			Transcripts:         2,
		},
		files:  map[string][]byte{},
		status: map[string]int{},
	}

	tau := big.NewInt(1)
	participants := make([]map[string]any, 0, len(contributors))

	for i, who := range contributors {
		tau = new(big.Int).Mod(new(big.Int).Mul(tau, big.NewInt(who.secret)), fr.Modulus())
		path := fmt.Sprintf("%03d_%s/transcript00.dat", i+1, strings.ToLower(who.address))
		c.files[path] = c.transcript(tau, big.NewInt(who.secret))
		participants = append(participants, map[string]any{"address": who.address, "position": i + 1})
	}

	c.tau = new(big.Int).Mod(new(big.Int).Mul(tau, big.NewInt(seal)), fr.Modulus())
	c.files[ignition.SealedTranscript] = c.transcript(c.tau, big.NewInt(seal))
	c.manifest = map[string]any{
		"name":                c.ceremony.Name,
		"numG1Points":         c.ceremony.Points,
		"numG2Points":         1,
		"pointsPerTranscript": c.ceremony.PointsPerTranscript,
		"participants":        participants,
		"crs":                 map[string]any{"t2": g2Hex(c.tau)},
	}

	return c
}

func (c *ceremonyFiles) transcript(tau, contribution *big.Int) []byte {
	var out bytes.Buffer

	header := []uint32{
		0,
		uint32(c.ceremony.Transcripts),
		uint32(c.ceremony.Points),
		1,
		uint32(c.ceremony.PointsPerTranscript),
		2,
		0,
	}
	for _, field := range header {
		_ = binary.Write(&out, binary.BigEndian, field)
	}

	_, _, g1, g2 := bn254.Generators()
	power := new(big.Int).Set(tau)

	for range c.ceremony.PointsPerTranscript {
		var point bn254.G1Affine
		point.ScalarMultiplication(&g1, power)
		out.Write(encodeG1(point))

		power.Mul(power, tau).Mod(power, fr.Modulus())
	}

	for _, scalar := range []*big.Int{tau, contribution} {
		var point bn254.G2Affine
		point.ScalarMultiplication(&g2, scalar)
		out.Write(encodeG2(point))
	}

	checksum := blake2b.Sum512(out.Bytes())
	out.Write(checksum[:])

	return out.Bytes()
}

func (c *ceremonyFiles) rewriteG1(t *testing.T, path string, index int, scalar *big.Int) {
	t.Helper()

	_, _, g1, _ := bn254.Generators()

	var point bn254.G1Affine
	point.ScalarMultiplication(&g1, scalar)

	file := c.files[path]
	copy(file[28+index*64:], encodeG1(point))
	c.reseal(path)
}

func (c *ceremonyFiles) rewriteG2(t *testing.T, path string, index int, scalar *big.Int) {
	t.Helper()

	_, _, _, g2 := bn254.Generators()

	var point bn254.G2Affine
	point.ScalarMultiplication(&g2, scalar)

	file := c.files[path]
	copy(file[28+c.ceremony.PointsPerTranscript*64+index*128:], encodeG2(point))
	c.reseal(path)
}

func (c *ceremonyFiles) aliasFirstCoordinate(t *testing.T, path string) {
	t.Helper()

	file := c.files[path]
	value := new(big.Int).SetBytes(reverseWords(file[28:60]))
	value.Add(value, fp.Modulus())

	if value.BitLen() > 256 {
		t.Skip("the coordinate has no alias below 2^256")
	}

	copy(file[28:60], reverseWords(value.FillBytes(make([]byte, 32))))
	c.reseal(path)
}

func (c *ceremonyFiles) reseal(path string) {
	file := c.files[path]
	checksum := blake2b.Sum512(file[:len(file)-64])
	copy(file[len(file)-64:], checksum[:])
}

func (c *ceremonyFiles) client(t *testing.T) *ignition.Client {
	t.Helper()

	manifest, err := json.Marshal(c.manifest)
	require.NoError(t, err)

	prefix := "/" + c.ceremony.Name + "/"
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path, err := url.PathUnescape(r.URL.EscapedPath())
			if err != nil || !strings.HasPrefix(path, prefix) {
				http.NotFound(w, r)
				return
			}

			name := strings.TrimPrefix(path, prefix)
			if status, ok := c.status[name]; ok {
				w.WriteHeader(status)
				return
			}

			if name == "manifest.json" {
				http.ServeContent(w, r, name, time.Time{}, bytes.NewReader(manifest))
				return
			}

			file, ok := c.files[name]
			if !ok {
				http.NotFound(w, r)
				return
			}

			full := c.ceremony.PointsPerTranscript*64 + 28 + 256 + 64
			if len(file) < full && r.Header.Get("Range") == "" {
				w.Header().Set("Content-Length", fmt.Sprint(full))
				_, _ = w.Write(file)

				return
			}

			http.ServeContent(w, r, name, time.Time{}, bytes.NewReader(file))
		}),
	)
	t.Cleanup(server.Close)

	client, err := ignition.NewClient(server.URL, c.ceremony, server.Client())
	require.NoError(t, err)

	return client
}

func encodeG1(point bn254.G1Affine) []byte {
	x, y := point.X.Bytes(), point.Y.Bytes()

	return append(reverseWords(x[:]), reverseWords(y[:])...)
}

func encodeG2(point bn254.G2Affine) []byte {
	var out []byte

	for _, element := range []fp.Element{point.X.A0, point.X.A1, point.Y.A0, point.Y.A1} {
		canonical := element.Bytes()
		out = append(out, reverseWords(canonical[:])...)
	}

	return out
}

func reverseWords(data []byte) []byte {
	out := make([]byte, 32)
	for word := range 4 {
		copy(out[word*8:(word+1)*8], data[(3-word)*8:(4-word)*8])
	}

	return out
}

func g2Hex(scalar *big.Int) []string {
	_, _, _, g2 := bn254.Generators()

	var point bn254.G2Affine
	point.ScalarMultiplication(&g2, scalar)

	out := make([]string, 0, 4)
	for _, element := range []fp.Element{point.X.A0, point.X.A1, point.Y.A0, point.Y.A1} {
		out = append(out, "0x"+element.Text(16))
	}

	return out
}
