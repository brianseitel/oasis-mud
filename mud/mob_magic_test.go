package mud

import "testing"

func TestDoCast(t *testing.T) {
	player := resetTest()
	player.Mana = 100000
	player.MaxMana = 1000000
	doCast(player, "")

	doCast(player, "spell")

	spell := mockSkill("spell")
	spell.Skill.MinMana = 25
	player.Skills = append(player.Skills, spell)

	doCast(player, "spell")

	spell.Skill.Target = targetIgnore
	doCast(player, "spell")

	spell.Skill.Target = targetCharacterOffensive
	doCast(player, "spell")
	doCast(player, "spell player")

	target := mockMob("target")
	player.Room.Mobs = append(player.Room.Mobs, target)
	doCast(player, "spell target")
	doCast(player, "spell not here")
	doCast(player, "spell player")

	spell.Skill.Target = targetCharacterDefensive
	doCast(player, "spell")
	doCast(player, "spell target")

	spell.Skill.Target = targetCharacterSelf
	doCast(player, "spell target")

	spell.Skill.Target = targetObjectInventory
	doCast(player, "spell foo")

	spell.Skill.Target = targetCharacterSelf
	spell.Level = 100
	doCast(player, "spell")

	player.Mana = 1
	doCast(player, "spell")
}
