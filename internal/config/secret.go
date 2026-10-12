package config

type Secret struct {
	value string
}

func (s Secret) Reveal() string {
	return s.value
}

func (s Secret) String() string {
	return "[redacted]"
}

func (s Secret) GoString() string {
	return "config.Secret{[redacted]}"
}
