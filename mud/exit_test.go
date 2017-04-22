package mud

import "testing"

func TestHasDoor(t *testing.T) {
	x := &exit{ID: 1, Keyword: "sign", Description: "This is an exit.", Dir: "east", Room: nil, Flags: exitDoor}

	if !x.hasDoor() {
		t.Error("Failed to find door.")
	}

	exitNoDoor := &exit{ID: 1, Keyword: "sign", Description: "This is an exit.", Dir: "east", Room: nil, Flags: 0}

	if exitNoDoor.hasDoor() {
		t.Error("Found a door it shouldn't have.")
	}
}

func TestIsClosed(t *testing.T) {
	x := &exit{ID: 1, Keyword: "sign", Description: "This is an exit.", Dir: "east", Room: nil, Flags: exitClosed}

	if !x.isClosed() {
		t.Error("Failed to find closed door.")
	}

	exitNoDoor := &exit{ID: 1, Keyword: "sign", Description: "This is an exit.", Dir: "east", Room: nil, Flags: 0}

	if exitNoDoor.isClosed() {
		t.Error("Found a closed door it shouldn't have.")
	}
}

func TestIsOpen(t *testing.T) {
	x := &exit{ID: 1, Keyword: "sign", Description: "This is an exit.", Dir: "east", Room: nil, Flags: exitClosed}

	if x.isOpen() {
		t.Error("Found an open door but it should be closed.")
	}

	exitNoDoor := &exit{ID: 1, Keyword: "sign", Description: "This is an exit.", Dir: "east", Room: nil, Flags: 0}

	if !exitNoDoor.isOpen() {
		t.Error("Failed to find an open door.")
	}
}

func TestIsLocked(t *testing.T) {
	x := &exit{ID: 1, Keyword: "sign", Description: "This is an exit.", Dir: "east", Room: nil, Flags: exitLocked}

	if !x.isLocked() {
		t.Error("Failed to find locked door.")
	}

	exitNoDoor := &exit{ID: 1, Keyword: "sign", Description: "This is an exit.", Dir: "east", Room: nil, Flags: 0}

	if exitNoDoor.isLocked() {
		t.Error("Found a locked door it shouldn't have.")
	}
}

func TestIsUnlocked(t *testing.T) {
	x := &exit{ID: 1, Keyword: "sign", Description: "This is an exit.", Dir: "east", Room: nil, Flags: exitLocked}

	if x.isUnlocked() {
		t.Error("Found an locked door but it should be closed.")
	}

	exitNoDoor := &exit{ID: 1, Keyword: "sign", Description: "This is an exit.", Dir: "east", Room: nil, Flags: 0}

	if !exitNoDoor.isUnlocked() {
		t.Error("Failed to find an locked door.")
	}
}

func TestIsPickProof(t *testing.T) {
	x := &exit{ID: 1, Keyword: "sign", Description: "This is an exit.", Dir: "east", Room: nil, Flags: exitPickProof}

	if !x.isPickProof() {
		t.Error("Failed to find pick-proof door.")
	}

	exitNoDoor := &exit{ID: 1, Keyword: "sign", Description: "This is an exit.", Dir: "east", Room: nil, Flags: 0}

	if exitNoDoor.isPickProof() {
		t.Error("Found a pick-proof door it shouldn't have.")
	}
}

func TestIsPickable(t *testing.T) {
	x := &exit{ID: 1, Keyword: "sign", Description: "This is an exit.", Dir: "east", Room: nil, Flags: exitPickProof}

	if x.isPickable() {
		t.Error("Found a pickable door but it should be pick-proof.")
	}

	exitNoDoor := &exit{ID: 1, Keyword: "sign", Description: "This is an exit.", Dir: "east", Room: nil, Flags: 0}

	if !exitNoDoor.isPickable() {
		t.Error("Failed to find a pickable door.")
	}
}

func TestReverseDirection(t *testing.T) {
	if reverseDirection("east") != "west" {
		t.Error("Failed to find reverse direction for east")
	}
	if reverseDirection("west") != "east" {
		t.Error("Failed to find reverse direction for west")
	}
	if reverseDirection("north") != "south" {
		t.Error("Failed to find reverse direction for north")
	}
	if reverseDirection("south") != "north" {
		t.Error("Failed to find reverse direction for south")
	}
	if reverseDirection("up") != "down" {
		t.Error("Failed to find reverse direction for up")
	}
	if reverseDirection("down") != "up" {
		t.Error("Failed to find reverse direction for down")
	}

	if reverseDirection("bad direction") != "oops" {
		t.Error("Failed to default to oops.")
	}
}
