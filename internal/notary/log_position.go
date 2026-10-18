package notary

type LogPosition struct {
	Transaction uint64
	Sequence    uint64
}

func (p LogPosition) After(other LogPosition) bool {
	if p.Transaction != other.Transaction {
		return p.Transaction > other.Transaction
	}

	return p.Sequence > other.Sequence
}
