package notary

type SnapshotRecord struct {
	ID       SnapshotID
	AsOf     UnixTime
	Nanos    uint32
	Currency string
	Count    uint64
	Total    uint64
	Digest   Digest
}
