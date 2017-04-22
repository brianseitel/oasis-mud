package mud

import "testing"

func TestAffectModifyAdd(t *testing.T) {
	player := &mob{Playable: true}
	paf := &affect{modifier: 1, bitVector: 1, location: applyHitroll, affectType: &mobSkill{Skill: &skill{MessageOff: "Test!"}}}

	// Add affect
	affectModify(player, paf, true)

	if player.AffectedBy != 1 {
		t.Error("affectModify did not apply player.AffectedBy bit.")
	}

	found := false
	for _, af := range player.Affects {
		if af == paf {
			found = true
		}
	}

	if !found {
		t.Error("Did not add affect.")
	}
}

func TestAffectModifyRemove(t *testing.T) {
	player := &mob{Playable: true, AffectedBy: 1}
	paf := &affect{modifier: 1, bitVector: 1, location: applyHitroll, affectType: &mobSkill{Skill: &skill{MessageOff: "Test!"}}}
	player.Affects = append(player.Affects, paf)

	// Remove affect
	affectModify(player, paf, false)

	if player.AffectedBy != 0 {
		t.Error("affectModify did not remove player.AffectedBy bit.")
	}

	found := false
	for _, af := range player.Affects {
		if af == paf {
			found = true
		}
	}

	if found {
		t.Error("Did not remove affect.")
	}
}

func TestAffectModifyUnknownLocation(t *testing.T) {
	player := &mob{Playable: true}
	paf := &affect{modifier: 1, bitVector: 1, location: 3265, affectType: &mobSkill{Skill: &skill{MessageOff: "Test!"}}}

	// Remove affect
	affectModify(player, paf, false)

	if player.AffectedBy != 0 {
		t.Error("affectModify did not remove player.AffectedBy bit.")
	}

	found := false
	for _, af := range player.Affects {
		if af == paf {
			found = true
		}
	}

	if found {
		t.Error("Did not remove affect.")
	}
}

func TestAffectModifyNPC(t *testing.T) {
	player := &mob{Playable: false}
	paf := &affect{modifier: 1, bitVector: 1, location: 3265, affectType: &mobSkill{Skill: &skill{MessageOff: "Test!"}}}

	// Remove affect
	affectModify(player, paf, false)

	if player.AffectedBy != 0 {
		t.Error("affectModify did not remove player.AffectedBy bit.")
	}

	found := false
	for _, af := range player.Affects {
		if af == paf {
			found = true
		}
	}

	if found {
		t.Error("Did not remove affect.")
	}
}

func TestAffectModifyArmorClass(t *testing.T) {
	player := &mob{Playable: true}
	paf := &affect{modifier: 1, location: applyArmorClass, affectType: &mobSkill{Skill: &skill{MessageOff: "Test!"}}}

	// Remove affect
	affectModify(player, paf, false)

	if player.Armor != 1 {
		t.Error("Did not apply affect: armor class")
	}
}

func TestAffectModifySex(t *testing.T) {
	player := &mob{Playable: true}
	paf := &affect{modifier: 1, location: applySex, affectType: &mobSkill{Skill: &skill{MessageOff: "Test!"}}}

	// Remove affect
	affectModify(player, paf, true)
}

func TestAffectModifyClass(t *testing.T) {
	player := &mob{Playable: true}
	paf := &affect{modifier: 1, location: applyClass, affectType: &mobSkill{Skill: &skill{MessageOff: "Test!"}}}

	// Remove affect
	affectModify(player, paf, true)
}

func TestAffectModifyGold(t *testing.T) {
	player := &mob{Playable: true}
	paf := &affect{modifier: 1, location: applyGold, affectType: &mobSkill{Skill: &skill{MessageOff: "Test!"}}}

	// Remove affect
	affectModify(player, paf, true)
}

func TestAffectModifyLevel(t *testing.T) {
	player := &mob{Playable: true}
	paf := &affect{modifier: 1, location: applyLevel, affectType: &mobSkill{Skill: &skill{MessageOff: "Test!"}}}

	// Remove affect
	affectModify(player, paf, true)
}

func TestAffectModifyExp(t *testing.T) {
	player := &mob{Playable: true}
	paf := &affect{modifier: 1, location: applyExp, affectType: &mobSkill{Skill: &skill{MessageOff: "Test!"}}}

	// Remove affect
	affectModify(player, paf, true)
}

func TestAffectModifySavingParalysis(t *testing.T) {
	player := &mob{Playable: true}
	paf := &affect{modifier: 1, location: applySavingParalysis, affectType: &mobSkill{Skill: &skill{MessageOff: "Test!"}}}

	// Remove affect
	affectModify(player, paf, true)
}

func TestAffectModifySavingRod(t *testing.T) {
	player := &mob{Playable: true}
	paf := &affect{modifier: 1, location: applySavingRod, affectType: &mobSkill{Skill: &skill{MessageOff: "Test!"}}}

	// Remove affect
	affectModify(player, paf, true)
}

func TestAffectModifySavingPetrify(t *testing.T) {
	player := &mob{Playable: true}
	paf := &affect{modifier: 1, location: applySavingPetrify, affectType: &mobSkill{Skill: &skill{MessageOff: "Test!"}}}

	// Remove affect
	affectModify(player, paf, true)
}

func TestAffectModifySavingBreath(t *testing.T) {
	player := &mob{Playable: true}
	paf := &affect{modifier: 1, location: applySavingBreath, affectType: &mobSkill{Skill: &skill{MessageOff: "Test!"}}}

	// Remove affect
	affectModify(player, paf, true)
}

func TestAffectModifySavingSpell(t *testing.T) {
	player := &mob{Playable: true}
	paf := &affect{modifier: 1, location: applySavingSpell, affectType: &mobSkill{Skill: &skill{MessageOff: "Test!"}}}

	// Remove affect
	affectModify(player, paf, true)
}

func TestAffectModifyNone(t *testing.T) {
	player := &mob{Playable: true}
	paf := &affect{modifier: 1, location: applyNone, affectType: &mobSkill{Skill: &skill{MessageOff: "Test!"}}}

	// Remove affect
	affectModify(player, paf, true)
}

func TestAffectModifyDamroll(t *testing.T) {
	player := &mob{Playable: true}
	paf := &affect{modifier: 1, location: applyDamroll, affectType: &mobSkill{Skill: &skill{MessageOff: "Test!"}}}

	damroll := player.Damroll
	// Remove affect
	affectModify(player, paf, true)

	if player.Damroll != damroll+1 {
		t.Error("Failed to apply damroll bonus")
	}
}

func TestAffectModifyHitroll(t *testing.T) {
	player := &mob{Playable: true}
	paf := &affect{modifier: 1, location: applyHitroll, affectType: &mobSkill{Skill: &skill{MessageOff: "Test!"}}}

	hitroll := player.Hitroll
	// Remove affect
	affectModify(player, paf, true)

	if player.Hitroll != hitroll+1 {
		t.Error("Failed to apply hitroll bonus")
	}
}

func TestAffectModifyMovement(t *testing.T) {
	player := &mob{Playable: true}
	paf := &affect{modifier: 1, location: applyMovement, affectType: &mobSkill{Skill: &skill{MessageOff: "Test!"}}}

	movement := player.Movement
	// Remove affect
	affectModify(player, paf, true)

	if player.Movement != movement+1 {
		t.Error("Failed to apply movement bonus")
	}
}

func TestAffectModifyHitpoints(t *testing.T) {
	player := &mob{Playable: true}
	paf := &affect{modifier: 1, location: applyHitpoints, affectType: &mobSkill{Skill: &skill{MessageOff: "Test!"}}}

	hitpoints := player.Hitpoints
	// Remove affect
	affectModify(player, paf, true)

	if player.Hitpoints != hitpoints+1 {
		t.Error("Failed to apply hitpoints bonus")
	}
}

func TestAffectModifyMana(t *testing.T) {
	player := &mob{Playable: true, ModifiedAttributes: &attributeSet{}}
	paf := &affect{modifier: 1, location: applyMana, affectType: &mobSkill{Skill: &skill{MessageOff: "Test!"}}}

	mana := player.Mana
	// Remove affect
	affectModify(player, paf, true)

	if player.Mana != mana+1 {
		t.Error("Failed to apply mana bonus")
	}
}

func TestAffectModifyConstitution(t *testing.T) {
	player := &mob{Playable: true, Attributes: &attributeSet{Constitution: 5}, ModifiedAttributes: &attributeSet{Constitution: 1}}
	paf := &affect{modifier: 1, location: applyConstitution, affectType: &mobSkill{Skill: &skill{MessageOff: "Test!"}}}

	constitution := player.currentConstitution()
	// Remove affect
	affectModify(player, paf, true)

	if player.currentConstitution() != constitution+1 {
		t.Error("Failed to apply constitution bonus")
	}
}

func TestAffectModifyStrength(t *testing.T) {
	player := &mob{Playable: true, Attributes: &attributeSet{Strength: 5}, ModifiedAttributes: &attributeSet{Strength: 1}}
	paf := &affect{modifier: 1, location: applyStrength, affectType: &mobSkill{Skill: &skill{MessageOff: "Test!"}}}

	strength := player.currentStrength()
	// Remove affect
	affectModify(player, paf, true)

	if player.currentStrength() != strength+1 {
		t.Error("Failed to apply strength bonus")
	}
}

func TestAffectModifyIntelligence(t *testing.T) {
	player := &mob{Playable: true, Attributes: &attributeSet{Intelligence: 5}, ModifiedAttributes: &attributeSet{Intelligence: 1}}
	paf := &affect{modifier: 1, location: applyIntelligence, affectType: &mobSkill{Skill: &skill{MessageOff: "Test!"}}}

	intelligence := player.currentIntelligence()
	// Remove affect
	affectModify(player, paf, true)

	if player.currentIntelligence() != intelligence+1 {
		t.Error("Failed to apply intelligence bonus")
	}
}

func TestAffectModifyWisdom(t *testing.T) {
	player := &mob{Playable: true, Attributes: &attributeSet{Wisdom: 5}, ModifiedAttributes: &attributeSet{Wisdom: 1}}
	paf := &affect{modifier: 1, location: applyWisdom, affectType: &mobSkill{Skill: &skill{MessageOff: "Test!"}}}

	wisdom := player.currentWisdom()
	// Remove affect
	affectModify(player, paf, true)

	if player.currentWisdom() != wisdom+1 {
		t.Error("Failed to apply wisdom bonus")
	}
}

func TestAffectModifyDexterity(t *testing.T) {
	player := &mob{Playable: true, Attributes: &attributeSet{Dexterity: 5}, ModifiedAttributes: &attributeSet{Dexterity: 1}}
	paf := &affect{modifier: 1, location: applyDexterity, affectType: &mobSkill{Skill: &skill{MessageOff: "Test!"}}}

	dexterity := player.currentDexterity()
	// Remove affect
	affectModify(player, paf, true)

	if player.currentDexterity() != dexterity+1 {
		t.Error("Failed to apply dexterity bonus")
	}
}

func TestAffectModifyCharisma(t *testing.T) {
	player := &mob{Playable: true, Attributes: &attributeSet{Charisma: 5}, ModifiedAttributes: &attributeSet{Charisma: 1}}
	paf := &affect{modifier: 1, location: applyCharisma, affectType: &mobSkill{Skill: &skill{MessageOff: "Test!"}}}

	charisma := player.currentCharisma()
	// Remove affect
	affectModify(player, paf, true)

	if player.currentCharisma() != charisma+1 {
		t.Error("Failed to apply charisma bonus")
	}
}

func TestIsAffected(t *testing.T) {
	player := &mob{Playable: true, Attributes: &attributeSet{Charisma: 5}, ModifiedAttributes: &attributeSet{Charisma: 1}}
	paf := &affect{modifier: 1, location: applyCharisma, affectType: &mobSkill{Skill: &skill{MessageOff: "Test!"}}}

	affectModify(player, paf, true)
	if !isAffected(player, paf) {
		t.Error("Failed to find charisma bonus")
	}
	affectModify(player, paf, false)

	if isAffected(player, paf) {
		t.Error("Failed to remove charisma bonus")
	}
}
