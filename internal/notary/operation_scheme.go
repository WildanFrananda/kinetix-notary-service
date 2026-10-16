package notary

type OperationScheme interface {
	OperationHash(entryPoint EntryPointIdentity, op UserOperation) Digest
}
