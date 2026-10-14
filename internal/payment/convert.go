package payment

func unsignedOf(value int64) (uint64, bool) {
	if value < 0 {
		return 0, false
	}

	return uint64(value), true
}

func nanosOf(value int32) (uint32, bool) {
	if value < 0 || value >= 1_000_000_000 {
		return 0, false
	}

	return uint32(value), true
}
