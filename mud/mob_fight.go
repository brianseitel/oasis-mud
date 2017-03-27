package mud

import (
	"strings"

	"github.com/brianseitel/oasis-mud/helpers"
)

func (m *mob) attack(target *mob, f *fight) {
	if target.Status != dead {
		damage := m.damage(target)
		target.takeDamage(damage)
		target.notify("%s attacks you for %d damagel!", m.Name, damage)
		m.notify("You strike %s for %d damagel!", target.Name, damage)

		if target.Status == dead {
			m.notify("You have KILLED %s to death!!", target.Name)
			m.Status = standing

			exp := xpCompute(m, target)
			m.Exp += exp
			m.notify("You gain %d experience points!", exp)
			m.checkLevelUp()

			// Cancel fight
			m.Fight = nil
			target.Fight = nil

			// whisk it away
			target.die()
			m.Room.removeMob(target)
			target.Room.removeMob(target)
			m.statusBar()
			return
		}
		m.Status = fighting
		target.Status = fighting
	}
}

func (m *mob) damage(victim *mob) int {
	if m != victim {
		if victim.parry(m) {
			return 0
		}

		if victim.dodge(m) {
			return 0
		}
	}
	return m.oneHit(victim)
}

func (m *mob) damroll() int {
	return 1000
	// return m.Attributes.Strength
}

func (m *mob) deathCry() {
	var msg string
	var drop int
	switch dice().Intn(4) {
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
		item := newItemFromIndex(getItem(uint(drop)))
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

func (m *mob) oneHit(victim *mob) int {
	var dam int
	if m.isNPC() {
		dam = dice().Intn(int(m.Level*3/2)) + int(m.Level/2)
		if m.equippedItem(wearWield) != nil {
			dam += int(dam / 2)
		}
	} else {
		if m.equippedItem(wearWield) != nil {
			wield := m.equippedItem(wearWield)
			dam = dice().Intn(int(wield.Value)) + wield.Value
		} else {
			dam = dice().Intn(4) + 1
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

	return dam
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
	if m.Fight.Mob1 == m {
		victim = m.Fight.Mob2
	} else {
		victim = m.Fight.Mob1
	}

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
