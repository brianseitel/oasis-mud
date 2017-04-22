package mud

import "testing"

func TestDice(t *testing.T) {
	dice := dice()

	if dice == nil {
		t.Error("Failed to create a dice set.")
	}
}

func TestD20(t *testing.T) {
	result := d20()

	if result > 20 || result < 0 {
		t.Errorf("Result of d20() is out of bounds: %d", result)
	}
}

func TestDInt(t *testing.T) {
	result := dInt(15)

	if result > 15 || result < 0 {
		t.Error("Result of dInt(15) is out of bounds: %d", result)
	}
}

func TestDBits(t *testing.T) {
	result := dBits(1)

	if result >= 30 {
		t.Error("Result of dBits(1) is out of bounds: %d", result)
	}
}
