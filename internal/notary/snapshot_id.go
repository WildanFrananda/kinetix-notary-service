package notary

type SnapshotID string

func NewSnapshotID(text string) (SnapshotID, error) {
	if len(text) != 36 {
		return "", Errf(KindInvalidArgument, "notary.NewSnapshotID", nil, "%q is not a snapshot id", text)
	}

	for i := range len(text) {
		c := text[i]

		switch i {
		case 8, 13, 18, 23:
			if c != '-' {
				return "", Errf(KindInvalidArgument, "notary.NewSnapshotID", nil, "%q is not a snapshot id", text)
			}
		default:
			if (c < '0' || c > '9') && (c < 'a' || c > 'f') {
				return "", Errf(KindInvalidArgument, "notary.NewSnapshotID", nil, "%q is not a snapshot id", text)
			}
		}
	}

	return SnapshotID(text), nil
}

func (id SnapshotID) IsZero() bool {
	return id == ""
}
