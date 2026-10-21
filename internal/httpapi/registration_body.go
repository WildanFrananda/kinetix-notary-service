package httpapi

type registrationBody struct {
	Holder    string `json:"holder"`
	IssuedAt  string `json:"issuedAt"`
	Signature string `json:"signature"`
}
