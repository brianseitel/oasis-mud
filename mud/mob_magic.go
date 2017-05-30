package mud

import (
	"fmt"
)

func doCast(player *mob, argument string) {
	var victim *mob
	var mana int
	var spell *mobSkill

	argument, arg1 := oneArgument(argument)
	argument, arg2 := oneArgument(argument)

	if arg1 == "" {
		player.notify("Cast which what where?\r\n")
		return
	}

	spell = player.skill(arg1)
	if spell == nil {
		player.notify("You can't do that. \r\n")
		return
	}

	mana = 0
	if !player.isNPC() {
		mana = max(spell.Skill.MinMana, 100/(2+player.Level))
	}

	// Find targets
	victim = nil

	switch spell.Skill.Target {
	case targetIgnore:
		break

	case targetCharacterOffensive:
		if arg2 == "" {
			player.notify("Cast the spell on whom?\r\n")
			return
		}

		for _, mob := range player.Room.Mobs {
			if matchesSubject(mob.Name, arg2) {
				victim = mob
				break
			}
		}

		if victim == nil {
			player.notify("They aren't here.\r\n")
			return
		}

		if victim == player {
			player.notify("You can't do that to yourself.\r\n")
			return
		}

	case targetCharacterDefensive:
		if arg2 == "" {
			victim = player
		} else {
			for _, mob := range player.Room.Mobs {
				if matchesSubject(mob.Name, arg2) {
					victim = mob
					break
				}
			}
		}

	case targetCharacterSelf:
		if arg2 != "" {
			player.notify("You cannot cast this spell on another.\r\n")
			return
		}
		victim = player

	case targetObjectInventory:
		break

	default:
		fmt.Printf("cast: bad target for %s\r\n", spell.Skill.Name)
	}

	if !player.isNPC() && player.Mana < mana {
		player.notify("You don't have enough mana.")
		return
	}

	if !player.isNPC() && dice().Intn(100) > int(spell.Level) {
		player.notify("You lost your concentration!")
		player.Mana -= mana / 2
	} else {
		player.Mana -= mana
		doSpell(spell, player, victim)
	}
}
