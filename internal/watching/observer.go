package watching

type Observer interface {
	ProvenFaults(total uint64)
}
