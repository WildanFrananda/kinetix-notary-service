package ethereum

import (
	"context"
	"net/http"
	"net/url"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"golang.org/x/time/rate"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

func Dial(settings Settings, observer RPCObserver) (*ethclient.Client, error) {
	const op = "ethereum.Dial"

	if observer == nil {
		observer = unobserved{}
	}

	transport := &pacedTransport{next: http.DefaultTransport, observer: observer}

	if settings.Rate > 0 {
		parsed, err := url.Parse(settings.RPCURL)
		if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") {
			return nil, notary.Errf(
				notary.KindInvalidArgument,
				op,
				nil,
				"only an http or https RPC endpoint can be held to a request rate",
			)
		}

		transport.pace = rate.NewLimiter(rate.Limit(settings.Rate), 1)
	}

	client, err := rpc.DialOptions(
		context.Background(),
		settings.RPCURL,
		rpc.WithHTTPClient(&http.Client{Transport: transport}),
	)
	if err != nil {
		return nil, notary.Errf(
			notary.KindLedgerUnavailable,
			op,
			nil,
			"the RPC endpoint could not be opened",
		)
	}

	return ethclient.NewClient(client), nil
}
