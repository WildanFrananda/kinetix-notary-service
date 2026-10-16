package bundler

type gasEstimate struct {
	PreVerificationGas            string `json:"preVerificationGas"`
	VerificationGasLimit          string `json:"verificationGasLimit"`
	CallGasLimit                  string `json:"callGasLimit"`
	PaymasterVerificationGasLimit string `json:"paymasterVerificationGasLimit"`
}
