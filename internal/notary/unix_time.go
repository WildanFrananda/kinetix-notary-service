package notary

import (
	"math"
	"time"
)

type UnixTime struct {
	seconds uint64
}

func UnixTimeOf(t time.Time) (UnixTime, error) {
	const op = "notary.UnixTimeOf"
	if t.Nanosecond() != 0 {
		return UnixTime{}, Errf(KindInvalidArgument, op, nil, "%s is not a whole second", t.Format(time.RFC3339Nano))
	}

	seconds := t.Unix()
	if seconds <= 0 {
		return UnixTime{}, Errf(KindInvalidArgument, op, nil, "%s is not after the Unix epoch", t.Format(time.RFC3339))
	}

	return UnixTime{seconds: uint64(seconds)}, nil
}

func UnixTimeFromSeconds(seconds uint64) (UnixTime, error) {
	if seconds == 0 || seconds > math.MaxInt64 {
		return UnixTime{}, Errf(
			KindInvalidArgument,
			"notary.UnixTimeFromSeconds",
			nil,
			"%d seconds is not a time after the Unix epoch",
			seconds,
		)
	}

	return UnixTime{seconds: seconds}, nil
}

func (u UnixTime) Seconds() uint64 {
	return u.seconds
}

func (u UnixTime) IsZero() bool {
	return u.seconds == 0
}

func (u UnixTime) Time() time.Time {
	if u.seconds == 0 || u.seconds > math.MaxInt64 {
		return time.Time{}
	}

	return time.Unix(int64(u.seconds), 0).UTC()
}
