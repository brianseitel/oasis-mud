package mud

func interpolate(level int, value00 int, value32 int) int {
	return value00 + level*(value32-value00)/32
}

func uRange(a int, b int, c int) int {

	if b < a {
		return a
	} else if b > c {
		return c
	}
	return b
}

func min(a int, b int) int {
	if a < b {
		return a
	}

	return b
}

func max(a int, b int) int {
	if a > b {
		return a
	}

	return b
}

func hasBit(flag int, bit int) bool {
	return flag&bit == bit
}

func setBit(flag int, bit int) int {
	flag |= bit
	return flag
}

func removeBit(flag int, bit int) int {
	flag ^= bit
	return flag
}

func toggleBit(flag int, bit int) int {
	if hasBit(flag, bit) {
		removeBit(flag, bit)
	} else {
		setBit(flag, bit)
	}
	return flag
}
