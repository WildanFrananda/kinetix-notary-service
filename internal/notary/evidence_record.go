package notary

type EvidenceRecord interface {
	Stream() EvidenceStream
	Position() LogPosition
	Parties() []string
}
