package notary

type SaltSource interface {
	Salt() (Salt, error)
}
