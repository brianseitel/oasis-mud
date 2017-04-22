package mud

import "testing"

func TestGetRace(t *testing.T) {
	loadRaces()

	race := getRace(1)

	if race.ID != 1 || race.Name != "Human" {
		t.Error("Did not find race.")
	}

	race = getRace(99)

	if race != nil {
		t.Error("Found fake race when it shouldn't")
	}
}

func TestDefaultStats(t *testing.T) {
	r := race{}

	if r.defaultStats("hitpoints") != 100 {
		t.Error("Did not find valid default stat: hitpoints")
	}
	if r.defaultStats("mana") != 0 {
		t.Error("Did not find valid default stat: mana")
	}
	if r.defaultStats("movement") != 100 {
		t.Error("Did not find valid default stat: movement")
	}
	if r.defaultStats("strength") != 12 {
		t.Error("Did not find valid default stat: strength")
	}
	if r.defaultStats("intelligence") != 12 {
		t.Error("Did not find valid default stat: intelligence")
	}
	if r.defaultStats("dexterity") != 12 {
		t.Error("Did not find valid default stat: dexterity")
	}
	if r.defaultStats("wisdom") != 12 {
		t.Error("Did not find valid default stat: wisdom")
	}
	if r.defaultStats("charisma") != 12 {
		t.Error("Did not find valid default stat: charisma")
	}
	if r.defaultStats("constitution") != 12 {
		t.Error("Did not find valid default stat: constitution")
	}
	if r.defaultStats("foobar") != 0 {
		t.Error("Did not find valid default stat: hitpoints")
	}
}
