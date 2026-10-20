package notary

type BlindingSource interface {
	Blinding() (Blinding, error)
}
