package mud

import "testing"

func TestDoSpell(t *testing.T) {
	spell := &mobSkill{Skill: &skill{Name: "armor"}}
	player := &mob{}
	victim := &mob{}

	doSpell(spell, player, victim)

	found := false
	for _, af := range victim.Affects {
		if af.affectType.Skill.Name == "armor" {
			found = true
		}
	}

	if !found {
		t.Error("Failed to apply affect")
	}
}

func TestObjCastSpell(t *testing.T) {
	var spell *skill
	level := 1
	var mobs []*mob
	victim := &mob{}
	mobs = append(mobs, victim)
	player := &mob{Level: 1, Room: &room{Mobs: mobs}}
	obj := &item{}

	// Doesn't do anything
	objCastSpell(spell, level, player, victim, obj)

	// // doesn't do anything either
	spell = &skill{Name: "armor", Target: targetIgnore}
	objCastSpell(spell, level, player, victim, obj)

	// cast attack spell on other player
	spell = &skill{Name: "armor", Target: targetCharacterOffensive}
	objCastSpell(spell, level, player, victim, obj)

	found := false
	for _, af := range victim.Affects {
		if af.affectType.Skill.Name == "armor" {
			found = true
		}
	}

	if !found {
		t.Error("Failed to apply affect")
	}

	// cast attack spell on other player with nil target
	spell = &skill{Name: "armor", Target: targetCharacterOffensive}
	player.Fight = victim
	objCastSpell(spell, level, player, nil, obj)

	found = false
	for _, af := range victim.Affects {
		if af.affectType.Skill.Name == "armor" {
			found = true
		}
	}

	if !found {
		t.Error("Failed to apply affect")
	}

	// cast attack spell on other player who isn't fighting
	spell = &skill{Name: "armor", Target: targetCharacterOffensive}
	player.Fight = nil
	objCastSpell(spell, level, player, nil, obj)

	found = false
	for _, af := range victim.Affects {
		if af.affectType.Skill.Name == "armor" {
			found = true
		}
	}

	if !found {
		t.Error("Failed to apply affect")
	}

	// cast attack spell on other player who is not an NPC
	spell = &skill{Name: "armor", Target: targetCharacterOffensive}
	player.Fight = nil
	victim.Playable = true
	objCastSpell(spell, level, player, nil, obj)

	found = false
	for _, af := range victim.Affects {
		if af.affectType.Skill.Name == "armor" {
			found = true
		}
	}

	if !found {
		t.Error("Failed to apply affect")
	}

	// cast defensive spell on other player
	spell = &skill{Name: "armor", Target: targetCharacterDefensive}
	objCastSpell(spell, level, player, victim, obj)

	found = false
	for _, af := range victim.Affects {
		if af.affectType.Skill.Name == "armor" {
			found = true
		}
	}

	if !found {
		t.Error("Failed to apply affect")
	}

	// where victim is nil (self)
	spell = &skill{Name: "armor", Target: targetCharacterDefensive}
	objCastSpell(spell, level, player, nil, obj)

	found = false
	for _, af := range victim.Affects {
		if af.affectType.Skill.Name == "armor" {
			found = true
		}
	}

	if !found {
		t.Error("Failed to apply affect")
	}

	// cast spell on self
	spell = &skill{Name: "armor", Target: targetCharacterSelf}
	objCastSpell(spell, level, player, victim, obj)

	found = false
	for _, af := range player.Affects {
		if af.affectType.Skill.Name == "armor" {
			found = true
		}
	}

	if !found {
		t.Error("Failed to apply affect")
	}

}
