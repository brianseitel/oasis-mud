package mud

import (
	"testing"
)

func TestInterpolate(t *testing.T) {
	results := interpolate(1, 0, 32)

	if results != 1 {
		t.Error("Failed to interpolate")
	}

	results = interpolate(15, 0, 99)
	if results != 46 {
		t.Error("Failed to interpolate")
	}
}

func TestURange(t *testing.T) {
	results := uRange(2, 1, 3)

	if results != 2 {
		t.Error("Failed to find range")
	}

	results = uRange(1, 3, 3)
	if results != 3 {
		t.Error("Failed to find range")
	}

	results = uRange(3, 3, 2)
	if results != 2 {
		t.Error("Failed to find range")
	}
}

func TestMin(t *testing.T) {
	if min(1, 2) != 1 {
		t.Error("Failed to find min.")
	}

	if min(2, 1) != 1 {
		t.Error("Failed to find min.")
	}
}

func TestMax(t *testing.T) {
	if max(2, 1) != 2 {
		t.Error("Failed to find max.")
	}

	if max(1, 2) != 2 {
		t.Error("Failed to find max.")
	}
}

func TestHasBit(t *testing.T) {
	if !hasBit(3, 1) {
		t.Error("Cannot find bit: 1 in 3")
	}

	if hasBit(1, 16) {
		t.Error("Incorrectly found bit that shouldn't exist.")
	}
}

func TestSetBit(t *testing.T) {
	flags := setBit(0, 16)

	if flags != 16 {
		t.Error("Failed to set flag")
	}
}

func TestRemoveBit(t *testing.T) {
	flags := removeBit(16, 16)

	if flags != 0 {
		t.Error("Failed to remove flag")
	}
}

func TestToggleBit(t *testing.T) {
	flags := 0

	flags = toggleBit(flags, 16)
	if flags != 16 {
		t.Error("Failed to toggle flag")
	}

	flags = toggleBit(flags, 16)
	if flags != 0 {
		t.Error("Failed to toggle flag")
	}
}

func TestPlayerActFlags(t *testing.T) {
	if len(playerActFlags(129)) != 2 {
		t.Error("Failed to list flags")
	}
}

func TestPlayerAffectFlags(t *testing.T) {
	if len(playerAffectFlags(129)) != 2 {
		t.Error("Failed to list flags")
	}
}
