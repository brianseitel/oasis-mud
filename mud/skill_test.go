package mud

import "testing"

func TestGetSkill(t *testing.T) {
	loadSkills()

	skill := getSkill(1)
	if skill.ID != 1 {
		t.Error("Failed to find skill")
	}

	skill = getSkill(12345)
	if skill != nil {
		t.Error("Found a skill it shouldn't have.")
	}
}

func TestGetSkillByName(t *testing.T) {
	loadSkills()

	skill := getSkillByName("trip")
	if skill == nil {
		t.Error("Failed to find skill")
	}

	skill = getSkillByName("completely fake skill")
	if skill != nil {
		t.Error("Found a skill it shouldn't have.")
	}
}
