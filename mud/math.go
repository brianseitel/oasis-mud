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
		flag = removeBit(flag, bit)
	} else {
		flag = setBit(flag, bit)
	}
	return flag
}

func playerActFlags(flags int) []int {
	var results []int

	playerFlags := []int{actIsNPC, actSentinel, actScavenger, actAggressive, actStayArea, actWimpy, actPet, actTrain}

	for _, flag := range playerFlags {
		if hasBit(flags, flag) {
			results = append(results, flag)
			flags ^= flag
		}
	}

	return results
}

func playerAffectFlags(flags int) []int {
	var results []int

	playerFlags := []int{
		affectBlind, affectInvisible, affectDetectEvil, affectDetectInvisible,
		affectDetectMagic, affectDetectHidden, affectHold, affectSanctuary,
		affectFaerieFire, affectInfrared, affectCurse, affectFlaming,
		affectPoison, affectProtect, affectParalysis, affectSneak, affectHide,
		affectSleep, affectCharm, affectFlying, affectPassDoor,
	}

	for _, flag := range playerFlags {
		if hasBit(flags, flag) {
			results = append(results, flag)
			flags ^= flag
		}
	}

	return results
}
