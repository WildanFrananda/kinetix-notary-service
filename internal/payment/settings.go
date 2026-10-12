package payment

import (
	"time"

	"github.com/WildanFrananda/kinetix-notary-service/internal/mesh"
)

type Settings struct {
	Endpoint string
	Identity mesh.Identity
	Timeout  time.Duration
}
