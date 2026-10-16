package bundler

type authorization struct {
	ChainID string `json:"chainId"`
	Address string `json:"address"`
	Nonce   string `json:"nonce"`
	YParity string `json:"yParity"`
	R       string `json:"r"`
	S       string `json:"s"`
}
