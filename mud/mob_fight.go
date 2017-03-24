package mud

import (
	"fmt"

	"github.com/brianseitel/oasis-mud/helpers"
)

func (m *mob) attack(target *mob, f *fight) {
	if target.Status != dead {
		damage := m.damage(target)
		target.takeDamage(damage)
		target.notify(fmt.Sprintf("%s attacks you for %d damagel!%s", m.Name, damage, helpers.Newline))
		m.notify(fmt.Sprintf("You strike %s for %d damagel!%s", target.Name, damage, helpers.Newline))

		if target.Status == dead {
			m.notify(fmt.Sprintf("You have KILLED %s to death!!%s", target.Name, helpers.Newline))
			m.Status = standing

			exp := xpCompute(m, target)
			m.Exp += exp
			m.notify(fmt.Sprintf("You gain %d experience points!%s", exp, helpers.Newline))
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
	return m.Attributes.Strength
}

func (m *mob) die() {
	// drop corpse in room
	corpse := &item{itemType: "corpse", Name: "A corpse of " + m.Name, Identifiers: "corpse," + m.Identifiers, Decays: decays, TTL: 1}
	m.Room.Items = append(m.Room.Items, corpse)

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

	m.notify(fmt.Sprintf("You dodge %s's attack.", m.Name))
	attacker.notify(fmt.Sprintf("%s dodges your attack.", attacker.Name))

	return true
}

func (m *mob) oneHit(victim *mob) int {
	var dam int
	if m.isNPC() {
		dam = dice().Intn(int(m.Level*3/2)) + int(m.Level/2)
		if m.equippedItem("wield") != nil {
			dam += int(dam / 2)
		}
	} else {
		if m.equippedItem("wield") != nil {
			wield := m.equippedItem("wield")
			dam = dice().Intn(int(wield.Max)) + wield.Min
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
		if m.equippedItem("wield") == nil {
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

	m.notify(fmt.Sprintf("You parry %s's attack.", m.Name))
	attacker.notify(fmt.Sprintf("%s parries your attack.", attacker.Name))

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
			m.notify(fmt.Sprintf("You attempt to trip %s but miss!%s", victim.Name, helpers.Newline))
			return
		}

		m.notify(fmt.Sprintf("You trip %s and %s goes down!%s", victim.Name, victim.Name, helpers.Newline))
		victim.notify(fmt.Sprintf("%s trips you and you go down!%s", m.Name, helpers.Newline))

		m.wait = 2
		victim.wait = 2
		victim.Status = sitting
	} else {
		m.notify("You can't do this again so soon!\n")
	}
}
