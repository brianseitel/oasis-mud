package mud

import (
	"fmt"
	"strings"

	"github.com/brianseitel/oasis-mud/helpers"
)

func (m *mob) attack(target *mob) {
	// if target.Status != dead {
	// 	damage := m.damage(target)
	// 	target.takeDamage(damage)

	// 	if target.Status == dead {
	// 		m.notify("You have KILLED %s to death!!", target.Name)
	// 		m.Status = standing

	// 		exp := xpCompute(m, target)
	// 		m.notify("You gain %d experience points!", exp)
	// 		m.gainExp(exp)

	// 		// Cancel fight
	// 		m.Fight = nil
	// 		target.Fight = nil

	// 		// whisk it away
	// 		target.die()
	// 		m.Room.removeMob(target)
	// 		target.Room.removeMob(target)
	// 		m.statusBar()
	// 		return
	// 	}
	// 	m.Status = fighting
	// 	target.Status = fighting
	// }
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
		if helpers.HasBit(m.AffectedBy, affectInvisible) {
			helpers.RemoveBit(m.AffectedBy, affectInvisible)
			m.stripAffect("invis")
			act("$n fades into existence.", m, nil, nil, actToRoom)
		}

		if helpers.HasBit(victim.AffectedBy, affectSanctuary) {
			dam /= 2
		}

		if helpers.HasBit(victim.AffectedBy, affectProtect) && m.isEvil() {
			dam -= dam / 4
		}

		if dam < 0 {
			dam = 0
		}

		if damageType >= typeHit {
			// TODO:
			// if m.isNPC() && dice().Intn(100) < m.Level / 2 {
			// 	m.disarm(victim)
			// }
			// if m.isNPC() && dice().Intn(100) < m.Level / 2 {
			// 	m.trip(victim)
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
		break

	case incapacitated:
		act("$n is incapacitated and will slowly die if not aided.", victim, nil, nil, actToRoom)
		victim.notify("You are incapacitated and will slowly die if not aided.")
		break

	case stunned:
		act("$n is stunned, but will probably recover", victim, nil, nil, actToRoom)
		victim.notify("You are stunned, but will probably recover.")
		break

	case dead:
		act("$n is DEAD!!", victim, 0, 0, actToRoom)
		victim.notify("You have been KILLED!!\r\n")
		break

	default:
		if dam > victim.MaxHitpoints/4 {
			victim.notify("That really did HURT!")
		}
		if victim.Hitpoints < victim.MaxHitpoints/4 {
			victim.notify("%sYou really are BLEEDING!%s", helpers.Red, helpers.Reset)
		}
		break
	}

	if !victim.isAwake() {
		victim.stopFighting(false)
	}

	if victim.Status == dead {
		// TODO: groupGain()

		if !victim.isNPC() {
			if victim.Exp > 1000*victim.Level {
				victim.gainExp((1000*victim.Level - victim.Exp) / 2)
			}
		}

		rawKill(victim)

		if !m.isNPC() && victim.isNPC() {
			// TODO: autoloot, autosacrifice
		}
	}

	if victim == m {
		return
	}

	if !victim.isNPC() && victim.client == nil {
		if dice().Intn(int(victim.wait)) == 0 {
			newAction(victim, victim.client, "recall")
			return
		}
	}

	// TODO: wimpy

	// TODO: flee

	return
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
		vp = "scratch"
		vs = "scratches"
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

	punct = ":"
	if dam <= 24 {
		punct = "."
	}

	var buf1 string
	var buf2 string
	var buf3 string
	if damageType == typeHit {
		buf1 = fmt.Sprintf("$n %s $N%s", vp, punct)
		buf2 = fmt.Sprintf("You %s $N%s", vs, punct)
		buf3 = fmt.Sprintf("$n %s you%s", vp, punct)
	} else {
		// if damageType >= 0 && damageType < MAX_SKILL {
		// 	attack = skillTable[damageType].nounDamage
		// } else
		if damageType >= typeHit && damageType <= typeHit+len(attackTable) {
			attack = attackTable[damageType-typeHit]
		}

		buf1 = fmt.Sprintf("$n's %s %s $N%s", attack, vp, punct)
		buf2 = fmt.Sprintf("Your %s %s $N%s", attack, vs, punct)
		buf3 = fmt.Sprintf("$n's %s %s you%s", attack, vp, punct)
	}

	act(buf1, m, nil, victim, actToNotVict)
	act(buf2, m, nil, victim, actToChar)
	act(buf3, m, nil, victim, actToVict)
}

func (m *mob) damroll() int {
	return 1000
	// return m.Attributes.Strength
}

func (m *mob) deathCry() {
	var msg string
	var drop int
	switch dBits(4) {
	case 0:
		msg = "$n hits the ground ... DEAD."
		break
	case 1:
		msg = "$n splatters blood on your armor."
		break
	case 2:
		msg = "You smell $n's sphincter releasing in death."
		drop = vnumTurd
		break
	case 3:
		msg = "$n's severed head plops on the ground."
		drop = vnumSeveredHead
		break
	case 4:
		msg = "$n's heart is torn from $s chest."
		drop = vnumTornHeart
		break
	case 5:
		msg = "$n's arm is sliced from $s dead body."
		drop = vnumSlicedArm
		break
	case 6:
		msg = "$n's leg is sliced from $s dead body."
		drop = vnumSlicedLeg
		break
	default:
		msg = "You hear $n's death cry."
	}

	act(msg, m, nil, nil, actToRoom)

	if drop > 0 {
		item := createItem(getItem(uint(drop)))
		itemList.PushBack(item)

		item.Name = strings.Replace(item.Name, "[name]", m.Name, -1)
		item.Description = strings.Replace(item.Description, "[name]", m.Name, -1)

		m.Room.Items = append(m.Room.Items, item)
	}

	if m.isNPC() {
		msg = "You hear something's death cry."
	} else {
		msg = "You hear someone's death cry."
	}

	oldRoom := m.Room
	for _, exit := range m.Room.Exits {
		m.Room = exit.Room
		act(msg, m, nil, nil, actToRoom)
	}
	m.Room = oldRoom
}

func (m *mob) die() {
	m.deathCry()
	// drop corpse in room
	corpse := &item{ItemType: 0, Name: "A corpse of " + m.Name, Timer: 1}
	m.Room.Items = append(m.Room.Items, corpse)

	m.Fight = nil
	m.Status = sitting
	// whisk them away to Nowhere
	m.Room = getRoom(0)
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
		chance = helpers.Min(60, 2*m.Level)
	} else {
		mobSkill := m.skill("dodge")
		if mobSkill != nil {
			chance = int(mobSkill.Level / 2)
		}
	}

	if dice().Intn(100) >= chance+m.Level-attacker.Level {
		return false
	}

	m.notify("You dodge %s's attack.", m.Name)
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
		thac0_00 = 32
		thac0_32 = 0
	}

	thac0 := helpers.Interpolate(m.Level, thac0_00, thac0_32) - m.Hitroll

	victimAC := helpers.Max(-15, victim.Armor/10)
	if !m.canSee(victim) {
		victim.Armor -= 4
	}

	diceroll := 99
	for diceroll >= 25 {
		diceroll = dBits(5)
	}

	if diceroll == 0 || (diceroll != 19 && diceroll < thac0-victimAC) {
		// miss. //
		m.damage(victim, 0, damageType)
		return
	}

	if m.isNPC() {
		dam = dice().Intn(m.Level*3/2) + (m.Level / 2)
		if wield == nil {
			dam += dam / 2
		}
	} else {
		if wield != nil {
			dam = dice().Intn(wield.Max) + wield.Min
		} else {
			dam = dice().Intn(4)
		}
	}

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
		chance = helpers.Min(60, 2*m.Level)
	} else {
		if m.equippedItem(wearWield) == nil {
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

	m.notify("You parry %s's attack.", m.Name)
	attacker.notify("%s parries your attack.", attacker.Name)

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
	return
}

func (m *mob) takeDamage(damage int) {
	m.Hitpoints -= damage
	if m.Hitpoints < 0 {
		m.Status = dead
		m.notify(helpers.Red + "You are DEAD!!!" + helpers.Reset)
	}
}

func (m *mob) trip() {

	if m.Fight == nil {
		m.notify("You aren't fighting anyone.\n")
		return
	}

	var victim *mob
	victim = m.Fight

	if victim.wait == 0 {

		var chance int
		if m.isNPC() {
			chance = helpers.Min(60, 2*m.Level)
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

	if m.Hitpoints <= 6 {
		m.Status = mortal
	} else if m.Hitpoints <= -3 {
		m.Status = incapacitated
	} else {
		m.Status = stunned
	}

	return
}
