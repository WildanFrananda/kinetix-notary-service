package delivering

import "time"

type Policy struct {
	UTCOffset     time.Duration
	EpochsPerTick int32
	Batch         int32
	FirstRetry    time.Duration
	LastRetry     time.Duration
	OverdueAfter  time.Duration
}
