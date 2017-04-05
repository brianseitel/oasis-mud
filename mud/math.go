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

	var playerFlags []int
	playerFlags = append(playerFlags, actIsNPC)
	playerFlags = append(playerFlags, actSentinel)
	playerFlags = append(playerFlags, actScavenger)
	playerFlags = append(playerFlags, actAggressive)
	playerFlags = append(playerFlags, actStayArea)
	playerFlags = append(playerFlags, actWimpy)
	playerFlags = append(playerFlags, actPet)
	playerFlags = append(playerFlags, actTrain)
	playerFlags = append(playerFlags, actPet)

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

	var playerFlags []int

	playerFlags = append(playerFlags, affectBlind)
	playerFlags = append(playerFlags, affectInvisible)
	playerFlags = append(playerFlags, affectDetectEvil)
	playerFlags = append(playerFlags, affectDetectInvisible)
	playerFlags = append(playerFlags, affectDetectMagic)
	playerFlags = append(playerFlags, affectDetectHidden)
	playerFlags = append(playerFlags, affectHold)
	playerFlags = append(playerFlags, affectSanctuary)
	playerFlags = append(playerFlags, affectFaerieFire)
	playerFlags = append(playerFlags, affectInfrared)
	playerFlags = append(playerFlags, affectCurse)
	playerFlags = append(playerFlags, affectFlaming)
	playerFlags = append(playerFlags, affectPoison)
	playerFlags = append(playerFlags, affectProtect)
	playerFlags = append(playerFlags, affectParalysis)
	playerFlags = append(playerFlags, affectSneak)
	playerFlags = append(playerFlags, affectHide)
	playerFlags = append(playerFlags, affectSleep)
	playerFlags = append(playerFlags, affectCharm)
	playerFlags = append(playerFlags, affectFlying)
	playerFlags = append(playerFlags, affectPassDoor)

	for _, flag := range playerFlags {
		if hasBit(flags, flag) {
			results = append(results, flag)
			flags ^= flag
		}
	}

	return results
}
