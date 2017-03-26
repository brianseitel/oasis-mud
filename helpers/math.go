package helpers

func Range(a int, b int, c int) int {

	if b < a {
		return a
	} else if b > c {
		return c
	}
	return b
}

func Min(a int, b int) int {
	if a < b {
		return a
	}

	return b
}

func Max(a int, b int) int {
	if a > b {
		return a
	}

	return b
}

func HasBit(flag uint, bit uint) bool {
	return flag&bit == bit
}

func SetBit(flag uint, bit uint) uint {
	flag |= bit
	return flag
}

func RemoveBit(flag uint, bit uint) uint {
	flag ^= bit
	return flag
}
