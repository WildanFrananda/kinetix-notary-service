package bundler

import "encoding/json"

type response struct {
	Result json.RawMessage `json:"result"`
	Error  *rpcError       `json:"error"`
}
