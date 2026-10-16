package bundler

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
