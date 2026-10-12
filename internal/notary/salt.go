package notary

type Salt [32]byte

func (s Salt) IsZero() bool {
	return s == Salt{}
}
