package ethereum

import (
	"bytes"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type refusal struct {
	id   common.Hash
	name string
	kind notary.Kind
}

func refusedBy(op string, rules []refusal, subject string, err error) error {
	raw, reverted := revertData(err)
	if !reverted {
		if isBareRevert(err) {
			return notary.Errf(notary.KindAnchorRefused, op, err, "%s without a reason", subject)
		}

		return unavailable(op, err)
	}

	for _, rule := range rules {
		if bytes.Equal(raw[:selectorLength], rule.id.Bytes()[:selectorLength]) {
			return notary.Errf(rule.kind, op, nil, "%s: %s", subject, rule.name)
		}
	}

	return notary.Errf(notary.KindAnchorRefused, op, nil, "%s with %s", subject, hexutil.Encode(raw[:selectorLength]))
}
