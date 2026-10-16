package notary

type OperationGas struct {
	PreVerificationGas            uint64
	VerificationGasLimit          uint64
	CallGasLimit                  uint64
	PaymasterVerificationGasLimit uint64
}
