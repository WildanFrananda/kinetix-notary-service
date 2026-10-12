package ethereum

import (
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

func Dial(settings Settings) (*ethclient.Client, error) {
	client, err := ethclient.Dial(settings.RPCURL)
	if err != nil {
		return nil, notary.Errf(
			notary.KindLedgerUnavailable,
			"ethereum.Dial",
			nil,
			"the RPC endpoint could not be opened",
		)
	}

	return client, nil
}
