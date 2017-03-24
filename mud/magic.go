package mud


func doSpell(spell *mobSkill, player *mob, victim *mob) {
	var af affect

	af.affectType = spell
	af.duration = 24
	af.modifier = -20
	af.location = applyAC

	victim.addAffect(&af)
	victim.notify("You feel someone protecting you.\r\n")
	if player != victim {
		player.notify("Ok\r\n")
	}
}
