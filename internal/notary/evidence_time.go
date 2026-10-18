package notary

import (
	"math"
	"time"
)

type EvidenceTime struct {
	seconds uint64
	nanos   uint32
}

func EvidenceTimeFrom(seconds int64, nanos int32) (EvidenceTime, error) {
	if seconds <= 0 || nanos < 0 || nanos >= int32(time.Second) {
		return EvidenceTime{}, Errf(
			KindInvalidArgument,
			"notary.EvidenceTimeFrom",
			nil,
			"%d seconds and %d nanoseconds is not a time after the Unix epoch",
			seconds,
			nanos,
		)
	}

	return EvidenceTime{seconds: uint64(seconds), nanos: uint32(nanos)}, nil
}

func (e EvidenceTime) Seconds() uint64 {
	return e.seconds
}

func (e EvidenceTime) Nanos() uint32 {
	return e.nanos
}

func (e EvidenceTime) IsZero() bool {
	return e.seconds == 0
}

func (e EvidenceTime) Time() time.Time {
	if e.seconds == 0 || e.seconds > math.MaxInt64 {
		return time.Time{}
	}

	return time.Unix(int64(e.seconds), int64(e.nanos)).UTC()
}
