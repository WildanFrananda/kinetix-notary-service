package notarytest

import (
	"testing"

	"github.com/WildanFrananda/kinetix-notary-service/internal/ethereum"
)

func ApprovalSigner(t *testing.T, label string) *ethereum.ApprovalSigner {
	t.Helper()

	key, _ := StatementKey(t, label)

	return ethereum.NewApprovalSigner(key)
}
