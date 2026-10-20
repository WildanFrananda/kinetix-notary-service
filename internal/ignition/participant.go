package ignition

import (
	"fmt"
	"strings"
)

type Participant struct {
	Address  string `json:"address"`
	Position int    `json:"position"`
}

func (p Participant) firstTranscript() string {
	return fmt.Sprintf("%03d_%s/transcript00.dat", p.Position, strings.ToLower(p.Address))
}
