package ignition

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/consensys/gnark-crypto/ecc/bn254"
	"golang.org/x/crypto/blake2b"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

const (
	maxManifestBytes = 4 << 20
	readBuffer       = 1 << 20
)

type Client struct {
	base     string
	ceremony Ceremony
	http     *http.Client
}

func NewClient(base string, ceremony Ceremony, client *http.Client) (*Client, error) {
	root, err := url.JoinPath(base, ceremony.Name)
	if err != nil {
		return nil, notary.Errf(notary.KindInvalidArgument, "ignition.NewClient", err, "the ceremony URL")
	}

	return &Client{base: root, ceremony: ceremony, http: client}, nil
}

func (c *Client) Manifest(ctx context.Context) (Manifest, error) {
	const op = "ignition.Client.Manifest"

	body, err := c.get(ctx, "manifest.json", "")
	if err != nil {
		return Manifest{}, err
	}
	defer func() { _ = body.Close() }()

	data, err := io.ReadAll(io.LimitReader(body, maxManifestBytes+1))
	if err != nil {
		return Manifest{}, notary.Errf(notary.KindSourceUnavailable, op, err, "reading the manifest")
	}

	if len(data) > maxManifestBytes {
		return Manifest{}, notary.Errf(
			notary.KindProofRejected,
			op,
			nil,
			"the manifest exceeds %d bytes",
			maxManifestBytes,
		)
	}

	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return Manifest{}, notary.Errf(notary.KindProofRejected, op, err, "the manifest is not JSON")
	}

	return manifest, nil
}

func (c *Client) Opening(ctx context.Context, path string) (Opening, error) {
	head, err := c.span(ctx, path, 0, headerSize+g1Size)
	if err != nil {
		return Opening{}, err
	}

	if err := parseFileHeader(head).check(c.ceremony, path); err != nil {
		return Opening{}, err
	}

	first, err := decodeG1(head[headerSize:])
	if err != nil {
		return Opening{}, notary.Errf(notary.KindProofRejected, "ignition.Client.Opening", err, "%s", path)
	}

	offset := int64(headerSize) + int64(c.ceremony.PointsPerTranscript)*g1Size

	tail, err := c.span(ctx, path, offset, 2*g2Size)
	if err != nil {
		return Opening{}, err
	}

	tau, contribution, err := decodeG2Pair(tail, path)
	if err != nil {
		return Opening{}, err
	}

	return Opening{First: first, Tau: tau, Contribution: contribution}, nil
}

func (c *Client) Powers(ctx context.Context, path string, count int) (Powers, error) {
	const op = "ignition.Client.Powers"
	if count < 1 || count > c.ceremony.PointsPerTranscript {
		return Powers{}, notary.Errf(
			notary.KindInvalidArgument,
			op,
			nil,
			"%d points asked of a transcript of %d",
			count,
			c.ceremony.PointsPerTranscript,
		)
	}

	body, err := c.get(ctx, path, "")
	if err != nil {
		return Powers{}, err
	}
	defer func() { _ = body.Close() }()

	hash, err := blake2b.New512(nil)
	if err != nil {
		return Powers{}, notary.Errf(notary.KindInvalidArgument, op, err, "BLAKE2b")
	}

	reader := bufio.NewReaderSize(body, readBuffer)
	hashed := io.TeeReader(reader, hash)

	head := make([]byte, headerSize)
	if err := readFull(hashed, head, path); err != nil {
		return Powers{}, err
	}

	if err := parseFileHeader(head).check(c.ceremony, path); err != nil {
		return Powers{}, err
	}

	points := make([]bn254.G1Affine, count)
	block := make([]byte, g1Size)

	for i := range c.ceremony.PointsPerTranscript {
		if err := readFull(hashed, block, path); err != nil {
			return Powers{}, err
		}

		if i >= count {
			continue
		}

		if points[i], err = decodeG1(block); err != nil {
			return Powers{}, notary.Errf(notary.KindProofRejected, op, err, "%s point %d", path, i)
		}
	}

	pair := make([]byte, 2*g2Size)
	if err := readFull(hashed, pair, path); err != nil {
		return Powers{}, err
	}

	tau, _, err := decodeG2Pair(pair, path)
	if err != nil {
		return Powers{}, err
	}

	checksum := make([]byte, checksumSize)
	if err := readFull(reader, checksum, path); err != nil {
		return Powers{}, err
	}

	if !bytes.Equal(checksum, hash.Sum(nil)) {
		return Powers{}, notary.Errf(
			notary.KindProofRejected,
			op,
			nil,
			"%s does not match its own BLAKE2b checksum",
			path,
		)
	}

	if _, err := reader.ReadByte(); err != io.EOF {
		return Powers{}, notary.Errf(notary.KindProofRejected, op, nil, "%s carries data after its checksum", path)
	}

	return Powers{Points: points, Tau: tau}, nil
}

func (c *Client) span(ctx context.Context, path string, offset int64, length int) ([]byte, error) {
	const op = "ignition.Client.span"

	body, err := c.get(ctx, path, fmt.Sprintf("bytes=%d-%d", offset, offset+int64(length)-1))
	if err != nil {
		return nil, err
	}
	defer func() { _ = body.Close() }()

	data := make([]byte, length)
	if err := readFull(body, data, path); err != nil {
		return nil, err
	}

	if n, _ := body.Read(make([]byte, 1)); n != 0 {
		return nil, notary.Errf(
			notary.KindProofRejected,
			op,
			nil,
			"%s answered more than the %d bytes asked",
			path,
			length,
		)
	}

	return data, nil
}

func (c *Client) get(ctx context.Context, path, byteRange string) (io.ReadCloser, error) {
	const op = "ignition.Client.get"

	target, err := url.JoinPath(c.base, path)
	if err != nil {
		return nil, notary.Errf(notary.KindInvalidArgument, op, err, "%s", path)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	if err != nil {
		return nil, notary.Errf(notary.KindInvalidArgument, op, err, "%s", path)
	}

	want := http.StatusOK

	if byteRange != "" {
		request.Header.Set("Range", byteRange)

		want = http.StatusPartialContent
	}

	response, err := c.http.Do(request)
	if err != nil {
		return nil, notary.Errf(notary.KindSourceUnavailable, op, err, "%s", path)
	}

	if response.StatusCode != want {
		_ = response.Body.Close()

		kind := notary.KindSourceUnavailable
		if response.StatusCode == http.StatusNotFound {
			kind = notary.KindProofRejected
		}

		return nil, notary.Errf(kind, op, nil, "%s answered %s, %d was expected", path, response.Status, want)
	}

	return response.Body, nil
}

func readFull(r io.Reader, into []byte, path string) error {
	if _, err := io.ReadFull(r, into); err != nil {
		return notary.Errf(notary.KindSourceUnavailable, "ignition.readFull", err, "%s ended early", path)
	}

	return nil
}

func decodeG2Pair(data []byte, path string) (bn254.G2Affine, bn254.G2Affine, error) {
	tau, err := decodeG2(data[:g2Size])
	if err != nil {
		return tau, tau, notary.Errf(notary.KindProofRejected, "ignition.decodeG2Pair", err, "%s τ", path)
	}

	contribution, err := decodeG2(data[g2Size:])
	if err != nil {
		return tau, contribution, notary.Errf(
			notary.KindProofRejected,
			"ignition.decodeG2Pair",
			err,
			"%s contribution",
			path,
		)
	}

	return tau, contribution, nil
}
