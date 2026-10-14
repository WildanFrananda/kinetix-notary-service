package notary

type SecretSource interface {
	AccountSecret() (AccountSecret, error)
}
