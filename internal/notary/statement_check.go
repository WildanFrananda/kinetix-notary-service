package notary

type StatementCheck uint8

const (
	StatementAbsent StatementCheck = iota + 1
	StatementAttributable
	StatementUnattributable
)

func (c StatementCheck) String() string {
	switch c {
	case StatementAbsent:
		return "absent"
	case StatementAttributable:
		return "attributable"
	case StatementUnattributable:
		return "not_attributable"
	default:
		return "unset"
	}
}
