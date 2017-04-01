package mud

func doSpell(spell *mobSkill, player *mob, victim *mob) {
	var af affect

	af.affectType = spell
	af.duration = 24
	af.modifier = -20
	af.location = applyArmorClass

	victim.addAffect(&af)

	if player != victim {
		player.notify("Ok\r\n")
	}
}

func objCastSpell(spell *skill, level int, player *mob, victim *mob, obj *item) {

	if spell == nil {
		return
	}

	var target *mob
	switch spell.Target {
	default:
		return

	case targetIgnore:
		target = nil
		break

	case targetCharacterOffensive:
		if victim == nil {
			victim = player.Fight
		}

		if victim == nil || !victim.isNPC() {
			player.notify("You can't do that.")
			return
		}

		target = victim
		break

	case targetCharacterDefensive:
		if victim == nil {
			victim = player
		}

		target = victim
		break

	case targetCharacterSelf:
		target = player
		break

	}

	mspell := &mobSkill{Skill: spell, Level: player.Level}

	doSpell(mspell, player, target)

	if spell.Target == targetCharacterOffensive && victim != player && victim.master != player {
		for _, m := range player.Room.Mobs {
			if victim == m && victim.Fight == nil {
				multiHit(victim, player, typeUndefined)
				break
			}
		}
	}
}
