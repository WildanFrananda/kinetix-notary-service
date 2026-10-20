package ignition

import (
	"encoding/binary"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

const (
	headerSize   = 28
	g1Size       = 64
	g2Size       = 128
	checksumSize = 64
)

type fileHeader struct {
	number      uint32
	transcripts uint32
	totalG1     uint32
	totalG2     uint32
	g1          uint32
	g2          uint32
	startFrom   uint32
}

func parseFileHeader(data []byte) fileHeader {
	return fileHeader{
		number:      binary.BigEndian.Uint32(data[0:4]),
		transcripts: binary.BigEndian.Uint32(data[4:8]),
		totalG1:     binary.BigEndian.Uint32(data[8:12]),
		totalG2:     binary.BigEndian.Uint32(data[12:16]),
		g1:          binary.BigEndian.Uint32(data[16:20]),
		g2:          binary.BigEndian.Uint32(data[20:24]),
		startFrom:   binary.BigEndian.Uint32(data[24:28]),
	}
}

func (h fileHeader) check(ceremony Ceremony, path string) error {
	if h.number != 0 ||
		int(h.transcripts) != ceremony.Transcripts ||
		int(h.totalG1) != ceremony.Points ||
		h.totalG2 != 1 ||
		int(h.g1) != ceremony.PointsPerTranscript ||
		h.g2 != 2 ||
		h.startFrom != 0 {
		return notary.Errf(
			notary.KindProofRejected,
			"ignition.fileHeader.check",
			nil,
			"%s opens with %+v, not the first transcript of %s",
			path,
			h,
			ceremony.Name,
		)
	}

	return nil
}
