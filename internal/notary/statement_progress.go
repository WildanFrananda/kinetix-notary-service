package notary

import "time"

type StatementProgress struct {
	Epoch uint64
	Day   time.Time
}
