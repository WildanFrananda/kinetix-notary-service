package notary

import "time"

type Lease struct {
	Instance  string
	Token     int64
	ExpiresAt time.Time
}
