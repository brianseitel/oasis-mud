package mud

import "testing"

func TestAddItem(t *testing.T) {
	player := resetTest()

	sword := mockObject("sword", 123)

	player.addItem(sword)
}

func TestCanDropItem(t *testing.T) {
	player := resetTest()
	sword := mockObject("sword", 123)

	sword.ExtraFlags = itemNoDrop
	if player.canDropItem(sword) == true {
		t.Error("Should not be able to drop")
	}

	player.Level = 99
	if player.canDropItem(sword) == false {
		t.Error("Should be able to drop")
	}
}

func TestCarrying(t *testing.T) {
	player := resetTest()
	sword := mockObject("sword", 123)
	player.addItem(sword)

	if player.carrying("sword") == nil {
		t.Error("Did not find sword")
	}

	if player.carrying("shield") != nil {
		t.Error("found shield when it shouldn't")
	}
}
