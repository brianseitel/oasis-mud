package mud

import (
	"fmt"
	"strings"
)

func doBackstab(player *mob, argument string) {
	backstab := player.skill("backstab")
	if backstab == nil {
		player.notify("You don't know how to backstab!")
		return
	}

	if len(argument) < 1 {
		player.notify("Backstab whom?")
		return
	}

	argument, arg1 := oneArgument(argument)

	var victim *mob
	for _, mob := range player.Room.Mobs {
		if matchesSubject(mob.Name, arg1) {
			victim = mob
			break
		}
	}

	if victim == nil {
		player.notify("They aren't here.")
		return
	}

	if victim == player {
		player.notify("How can you sneak up on yourself?")
		return
	}

	if victim.isSafe() {
		return
	}

	wield := player.equippedItem(itemWearWield)

	if wield == nil /* && item.type == "piercing" */ {
		player.notify("You need to wield a piercing weapon.")
		return
	}

	if victim.Fight != nil {
		player.notify("You can't backstab a fighting person.")
		return
	}

	if victim.Hitpoints < victim.MaxHitpoints {
		act("$N is hurt and suspicious ... you can't sneak up.", player, nil, victim, actToChar)
		return
	}

	wait(player, backstab.Skill.Beats)

	chance := dice().Intn(100)
	if !victim.isAwake() || player.isNPC() || chance > int(backstab.Level) {
		multiHit(player, victim, typeBackstab)
	} else {
		player.damage(victim, 0, typeBackstab)
	}

}

func doDisarm(player *mob, argument string) {
	disarm := player.skill("disarm")
	if !player.isNPC() && disarm == nil {
		player.notify("You don't know how to disarm!")
		return
	}

	wield := player.equippedItem(itemWearWield)
	if wield == nil {
		player.notify("You must wield a weapon to disarm.")
		return
	}

	victim := player.Fight
	if victim == nil {
		player.notify("You aren't fighting anyone, fool!")
		return
	}

	victimWield := victim.equippedItem(itemWearWield)
	if victimWield == nil {
		player.notify("Your opponent is not wielding a weapon.")
		return
	}

	wait(player, disarm.Skill.Beats)
	percent := dice().Intn(100) + victim.Level - player.Level
	if player.isNPC() || percent < int(disarm.Level*2/3) {
		player.disarm(victim)
	} else {
		player.notify("You failed.")
	}
}

func doFlee(player *mob, argument string) {
	victim := player.Fight
	if victim == nil || player.Status != fighting {
		player.notify("You aren't fighting anyone, fool.")
		return
	}

	wasIn := player.Room
	if len(player.Room.Exits) == 0 {
		player.notify("There is nowhere to flee!")
		return
	}

	for attempt := 0; attempt < 6; attempt++ {
		count := len(player.Room.Exits)

		number := dice().Intn(count)
		exit := player.Room.Exits[number]
		if exit == nil || exit.isClosed() || (player.isNPC() && hasBit(exit.Room.RoomFlags, roomNoMob)) {
			continue
		}

		player.move(exit)
		nowIn := player.Room
		if nowIn == wasIn {
			continue
		}

		player.Room = wasIn
		act("$n has fled!", player, nil, nil, actToRoom)
		player.Room = nowIn

		if !player.isNPC() {
			player.notify("You flee from combat! You lose 25 experience points!")
			player.gainExp(-25)
		}

		player.stopFighting(true)
		return
	}

	player.notify("You failed! You lose 10 experience points.")
	player.gainExp(-10)
}

func doKick(player *mob, argument string) {
	victim := player.Fight

	kick := player.skill("kick")
	if !player.isNPC() && kick == nil {
		player.notify("You better leave the martial arts to fighters.")
		return
	}

	if victim == nil {
		player.notify("You aren't fighting anyone, fool.")
		return
	}

	wait(player, kick.Skill.Beats)

	if player.isNPC() || dice().Intn(100) < int(kick.Level) {
		player.damage(victim, dice().Intn(player.Level)+1, typeKick)
	} else {
		player.damage(victim, 0, typeKick)
	}
}

func doKill(attacker *mob, argument string) {
	argument, arg1 := oneArgument(argument)

	if arg1 == "" {
		attacker.notify("Kill whom?")
		return
	}

	var victim *mob
	for _, m := range attacker.Room.Mobs {
		if matchesSubject(m.Name, arg1) {
			victim = m
			break
		}
	}

	if victim == nil {
		attacker.notify("They aren't here.")
		return
	}

	if victim == attacker {
		attacker.notify("You hit yourself. Ouch!")
		multiHit(attacker, attacker, typeUndefined)
		return
	}

	if !victim.isNPC() {
		attacker.notify("You cannot attack other players.")
		return
	}

	if victim.isSafe() {
		attacker.notify("A voice from the clouds booms, \"THOU SHALT NOT KILL!\"")
		return
	}

	if attacker.Status == fighting || attacker.Fight != nil {
		attacker.notify("You do the best you can!")
	}

	attacker.Fight = victim
	victim.Fight = attacker

	wait(attacker, 1*pulseViolence)
	multiHit(attacker, victim, typeHit)

}

func doRescue(player *mob, argument string) {
	argument, arg1 := oneArgument(argument)

	if arg1 == "" {
		player.notify("Rescue whom?")
		return
	}

	var victim *mob
	for _, m := range player.Room.Mobs {
		if m.Playable && matchesSubject(m.Name, arg1) {
			victim = m
			break
		}
	}

	if victim == nil {
		player.notify("They aren't here.")
		return
	}

	if victim == player {
		player.notify("Maybe you should try running away.")
		return
	}

	if !player.isNPC() && victim.isNPC() {
		player.notify("They don't need your help.")
		return
	}

	if player.Fight == victim {
		player.notify("You're trying to kill them!")
		return
	}

	if victim.Fight == nil {
		player.notify("They aren't fighting right now.")
		return
	}

	rescue := player.skill("rescue")
	if rescue == nil {
		player.notify("You don't know how to rescue!")
		return
	}

	wait(player, rescue.Skill.Beats)

	if !player.isNPC() && dice().Intn(100) > rescue.Level {
		player.notify("You failed the rescue.")
		return
	}

	act("You rescue $N!", player, nil, victim, actToRoom)
	act("$n rescues you!", player, nil, victim, actToVict)
	act("$n rescues $N!", player, nil, victim, actToNotVict)

	attacker := victim.Fight
	victim.stopFighting(false)
	attacker.Fight = player
	player.Fight = attacker
}

func (m *mob) damage(victim *mob, dam int, damageType int) {
	if victim.Status == dead {
		return
	}

	if dam > 1000 {
		dam = 1000
	}

	if victim != m {
		if victim.isSafe() {
			return
		}

		if victim.Status == stunned {
			victim.Fight = m
			victim.Status = fighting
		}

		if victim.Status > stunned {
			if m.Fight == nil {
				m.Fight = victim
			}
		}

		/*
		 * If they're invisible, fade them in
		 */
		if hasBit(m.AffectedBy, affectInvisible) {
			m.AffectedBy = removeBit(m.AffectedBy, affectInvisible)
			m.stripAffect("invis")
			act("$n fades into existence.", m, nil, nil, actToRoom)
		}

		if hasBit(victim.AffectedBy, affectSanctuary) {
			dam /= 2
		}

		if hasBit(victim.AffectedBy, affectProtect) && m.isEvil() {
			dam -= dam / 4
		}

		if dam < 0 {
			dam = 0
		}

		if damageType >= typeHit {
			// if m.isNPC() && dice().Intn(100) < m.Level/2 {
			// 	m.disarm(victim)
			// }
			// if m.isNPC() && dice().Intn(100) < m.Level/2 {
			// 	m.trip()
			// }
			if m.parry(victim) {
				return
			}
			if m.dodge(victim) {
				return
			}
		}

		m.damageMessage(victim, dam, damageType)
	}

	// hurt the victim
	victim.Hitpoints -= dam

	victim.updateStatus()

	switch victim.Status {
	case mortal:
		act("$n is mortally wounded and will die soon if not aided.", victim, nil, nil, actToRoom)
		victim.notify("You are mortally wounded and will die soon if not aided.")

	case incapacitated:
		act("$n is incapacitated and will slowly die if not aided.", victim, nil, nil, actToRoom)
		victim.notify("You are incapacitated and will slowly die if not aided.")

	case stunned:
		act("$n is stunned, but will probably recover", victim, nil, nil, actToRoom)
		victim.notify("You are stunned, but will probably recover.")

	case dead:
		act("$n is DEAD!!", victim, 0, 0, actToRoom)
		victim.notify("You have been KILLED!!\r\n")

	default:
		if dam > victim.MaxHitpoints/4 {
			victim.notify("That really did HURT!")
		}
		if victim.Hitpoints < victim.MaxHitpoints/4 {
			victim.notify("%sYou really are BLEEDING!%s", red, reset)
		}
	}

	if !victim.isAwake() {
		victim.stopFighting(false)
	}

	if victim.Status == dead {
		groupGain(m, victim)

		if !victim.isNPC() {
			if victim.Exp > 1000*victim.Level {
				victim.gainExp((1000*victim.Level - victim.Exp) / 2)
			}
		}

		rawKill(victim)

		// if !m.isNPC() && victim.isNPC() {
		// 	// TODO: autoloot, autosacrifice
		// }
	}

	if victim == m {
		return
	}

	if !victim.isNPC() && victim.client == nil {
		wait := max(1, int(victim.wait))
		if dice().Intn(wait) == 0 {
			interpret(victim, "recall")
			return
		}
	}

	if victim.isNPC() && dam > 0 {
		if hasBit(victim.Act, actWimpy) && dBits(1) == 0 && victim.Hitpoints < victim.MaxHitpoints/2 {
			doFlee(victim, "")
		}
	}

	if !victim.isNPC() && victim.Hitpoints < 0 && victim.wait == 0 {
		doFlee(victim, "")
	}
}

func (m *mob) damageMessage(victim *mob, dam int, damageType int) {
	attackTable := []string{"hit", "slice", "stab", "slash", "whip", "claw", "blast", "pound", "crush", "bite", "pierce"}

	var (
		vs     string
		vp     string
		attack string
		punct  string
	)

	if dam == 0 {
		vs = "miss"
		vp = "misses"
	} else if dam <= 4 {
		vs = "scratch"
		vp = "scratches"
	} else if dam <= 8 {
		vs = "graze"
		vp = "grazes"
	} else if dam <= 12 {
		vs = "hit"
		vp = "hits"
	} else if dam <= 16 {
		vs = "injure"
		vp = "injures"
	} else if dam <= 20 {
		vs = "wound"
		vp = "wounds"
	} else if dam <= 24 {
		vs = "maul"
		vp = "mauls"
	} else if dam <= 28 {
		vs = "decimate"
		vp = "decimates"
	} else if dam <= 32 {
		vs = "devastate"
		vp = "devastates"
	} else if dam <= 36 {
		vs = "maim"
		vp = "maims"
	} else if dam <= 40 {
		vs = "MUTLIATE"
		vp = "MUTILATES"
	} else if dam <= 44 {
		vs = "DISEMBOWEL"
		vp = "DISEMBOWELS"
	} else if dam <= 48 {
		vs = "EVISCERATE"
		vp = "EVISCERATES"
	} else if dam <= 52 {
		vs = "MASSACRE"
		vp = "MASSACRES"
	} else if dam <= 100 {
		vs = "*** DEMOLISH ***"
		vp = "*** DEMOLISHES ***"
	} else {
		vs = "*** ANNIHILATE ***"
		vp = "*** ANIHILIATES ***"
	}

	punct = "!"
	if dam <= 24 {
		punct = "."
	}

	var buf1 string
	var buf2 string
	var buf3 string
	if damageType == typeHit {
		buf1 = fmt.Sprintf("$n %s $N%s (%d)", vp, punct, dam)
		buf2 = fmt.Sprintf("You %s $N%s (%d)", vs, punct, dam)
		buf3 = fmt.Sprintf("$n %s you%s (%d)", vp, punct, dam)
	} else {
		// if damageType >= 0 && damageType < MAX_SKILL {
		// 	attack = skillTable[damageType].nounDamage
		// } else
		if damageType >= typeHit && damageType <= typeHit+len(attackTable) {
			attack = attackTable[damageType-typeHit]
		}

		buf1 = fmt.Sprintf("$n's %s %s $N%s (%d)", attack, vp, punct, dam)
		buf2 = fmt.Sprintf("Your %s %s $N%s (%d)", attack, vs, punct, dam)
		buf3 = fmt.Sprintf("$n's %s %s you%s (%d)", attack, vp, punct, dam)
	}

	act(buf1, m, nil, victim, actToNotVict)
	act(buf2, m, nil, victim, actToChar)
	act(buf3, m, nil, victim, actToVict)
}

func (m *mob) damroll() int {
	return m.Damroll + m.currentStrength()
}

func (m *mob) deathCry() {
	var msg string
	var drop int
	switch dBits(4) {
	case 0:
		msg = "$n hits the ground ... DEAD."
	case 1:
		msg = "$n splatters blood on your armor."
	case 2:
		msg = "You smell $n's sphincter releasing in death."
		drop = vnumTurd
	case 3:
		msg = "$n's severed head plops on the ground."
		drop = vnumSeveredHead
	case 4:
		msg = "$n's heart is torn from $s chest."
		drop = vnumTornHeart
	case 5:
		msg = "$n's arm is sliced from $s dead body."
		drop = vnumSlicedArm
	case 6:
		msg = "$n's leg is sliced from $s dead body."
		drop = vnumSlicedLeg
	default:
		msg = "You hear $n's death cry."
	}

	act(msg, m, nil, nil, actToRoom)

	if drop > 0 {
		item := createItem(getItem(drop))
		itemList.PushBack(item)

		item.Name = strings.Replace(item.Name, "[name]", m.Name, -1)
		item.Description = strings.Replace(item.Description, "[name]", m.Name, -1)
		item.Timer = dice().Intn(5)
		item.Room = m.Room
		m.Room.Items = append(m.Room.Items, item)
	}

	if m.isNPC() {
		msg = "You hear something's death cry."
	} else {
		msg = "You hear someone's death cry."
	}

	oldRoom := m.Room
	for _, exit := range m.Room.Exits {
		if exit.Room != nil {
			m.Room = exit.Room
			act(msg, m, nil, nil, actToRoom)
		}
	}
	m.Room = oldRoom
}

func (m *mob) disarm(victim *mob) {
	victimWield := victim.equippedItem(itemWearWield)
	if victimWield == nil {
		return
	}

	wield := m.equippedItem(itemWearWield)
	if wield == nil && dBits(1) == 0 {
		return
	}

	act("$n disarms you!", m, nil, victim, actToVict)
	act("You disarm $N!", m, nil, victim, actToChar)
	act("$n disarms $N!", m, nil, victim, actToNotVict)

	for j, item := range m.Equipped {
		if wield == item {
			item.WearLocation = wearNone
			m.Equipped, m.Room.Items = transferItem(j, m.Equipped, m.Room.Items)
			break
		}
	}
}

func (m *mob) dodge(attacker *mob) bool {
	var chance int

	if !m.isAwake() {
		return false // can't dodge if you're asleep!
	}

	if m.isNPC() {
		chance = min(60, 2*m.Level)
	} else {
		mobSkill := m.skill("dodge")
		if mobSkill != nil {
			chance = int(mobSkill.Level / 2)
		}
	}

	if dice().Intn(100) >= chance+m.Level-attacker.Level {
		return false
	}

	m.notify("You dodge %s's attack.", attacker.Name)
	attacker.notify("%s dodges your attack.", attacker.Name)

	return true
}

func (m *mob) oneHit(victim *mob, damageType int) {
	if victim.Status == dead || victim.Room != m.Room {
		return
	}

	wield := m.equippedItem(itemWearWield)

	if damageType == typeUndefined {
		damageType = typeHit
		if wield != nil && wield.ItemType == itemWeapon {
			damageType += wield.Min
		}
	}

	var dam int
	var thac0_00 int
	var thac0_32 int
	if m.isNPC() {
		thac0_00 = 20
		thac0_32 = 0
	} else {
		thac0_00 = m.Job.Thac0_00
		thac0_32 = m.Job.Thac0_32
	}

	thac0 := interpolate(m.Level, thac0_00, thac0_32) - m.hitroll()
	victimAC := max(-15, victim.Armor/10)
	if !m.canSee(victim) {
		victimAC -= 4
	}

	diceroll := 9999
	for diceroll >= 20 {
		diceroll = dBits(5)
	}

	if diceroll == 0 || (diceroll != 19 && diceroll < thac0-victimAC) {
		// miss. //
		m.damage(victim, 0, damageType)
		return
	}

	if m.isNPC() {
		dam = dice().Intn(m.Level*5/2) + (m.Level * 3 / 2)
		if wield == nil {
			dam = dam / 2
		}
	} else {
		if wield != nil {
			dam = dice().Intn(wield.Max) + wield.Min
		} else {
			dam = dice().Intn(4)
		}
	}

	dam += m.damroll()
	enhancedDamage := m.skill("enhanced_damage")
	if enhancedDamage != nil {
		if !m.isNPC() && int(m.skill("enhanced_damage").Level) > 0 {
			dam += int(dam * int(m.skill("enhanced_damage").Level) / 100)
		}
	}

	if victim.Status < sitting {
		dam *= 2
	}

	if dam <= 0 {
		dam = 1
	}

	m.damage(victim, dam, typeHit)
}

func (m *mob) parry(attacker *mob) bool {
	var chance int

	if !m.isAwake() {
		return false
	}

	if m.isNPC() {
		chance = min(60, 2*m.Level)
	} else {
		if m.equippedItem(itemWearWield) == nil {
			return false
		}

		mobSkill := m.skill("parry")
		if mobSkill != nil {
			chance = int(mobSkill.Level / 2)
		}
	}

	if dice().Intn(100) >= chance+m.Level-attacker.Level {
		return false
	}

	attacker.notify("You parry %s's attack.", attacker.Name)
	m.notify("%s parries your attack.", m.Name)

	return true
}

func (m *mob) stopFighting(both bool) {
	for e := mobList.Front(); e != nil; e = e.Next() {
		fighter := e.Value.(*mob)

		if fighter == m || (both && fighter.Fight == m) {
			fighter.Fight = nil
			fighter.Status = standing
			fighter.updateStatus()
		}
	}
}

func (m *mob) takeDamage(damage int) {
	m.Hitpoints -= damage
	if m.Hitpoints < 0 {
		m.Status = dead
		m.notify(red + "You are DEAD!!!" + reset)
	}
}

func (m *mob) trip() {

	if m.Fight == nil {
		m.notify("You aren't fighting anyone.\n")
		return
	}

	victim := m.Fight

	if victim.wait == 0 {

		var chance int
		if m.isNPC() {
			chance = min(60, 2*m.Level)
		} else {
			mobSkill := m.skill("trip")
			if mobSkill != nil {
				chance = int(mobSkill.Level / 2)
			}
		}

		if dice().Intn(100) >= chance+m.Level-victim.Level {
			m.notify("You attempt to trip %s but miss!", victim.Name)
			return
		}

		m.notify("You trip %s and %s goes down!", victim.Name, victim.Name)
		victim.notify("%s trips you and you go down!", m.Name)

		m.wait = 2
		victim.wait = 2
		victim.Status = sitting
	} else {
		m.notify("You can't do this again so soon!\n")
	}
}

func (m *mob) updateStatus() {
	if m.Hitpoints > 0 {
		if m.Status <= stunned {
			m.Status = standing
		}
		return
	}

	if m.isNPC() || m.Hitpoints <= -11 {
		m.Status = dead
		return
	}

	if m.Hitpoints <= -6 {
		m.Status = mortal
	} else if m.Hitpoints <= -3 {
		m.Status = incapacitated
	} else {
		m.Status = stunned
	}
}
