package notary

type Opening struct {
	Holder Holder
	Kind   AccountKind
	Salt   Salt
	Secret AccountSecret
}

func (o Opening) Commitment(scheme Scheme) Digest {
	if o.Secret.IsZero() {
		return scheme.Commitment(o.Holder, o.Kind, o.Salt)
	}

	return scheme.AccountCommitment(scheme.AccountID(o.Holder, o.Kind, o.Secret), o.Salt)
}
