package tokens

import "time"

type Settings struct {
	JWKSURL      string
	Issuer       string
	Audience     string
	RefreshEvery time.Duration
}
