package helpers

func Interpolate(level int, value_00 int, value_32 int) int {
	return value_00 + level*(value_32-value_00)/32
}

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
