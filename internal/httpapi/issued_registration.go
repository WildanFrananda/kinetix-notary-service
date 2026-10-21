package httpapi

type issuedRegistration struct {
	Holder    string `json:"holder"`
	Account   string `json:"account"`
	FromEpoch string `json:"fromEpoch"`
	Kind      string `json:"kind"`
	Digest    string `json:"digest"`
	Signature string `json:"signature"`
}
