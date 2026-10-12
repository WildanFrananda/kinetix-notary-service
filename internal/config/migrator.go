package config

import (
	"errors"
	"strings"
)

func LoadMigrator(lookup Lookup) (Secret, error) {
	value, ok := lookup(DatabaseURLKey)

	value = strings.TrimSpace(value)
	if !ok || value == "" {
		return Secret{}, errors.New("configuration refused: " + DatabaseURLKey + " is not set")
	}

	return Secret{value: value}, nil
}
