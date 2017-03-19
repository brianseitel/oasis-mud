package mud

import (
	"fmt"

	"github.com/brianseitel/oasis-mud/helpers"
)

func (m *mob) parry(target *mob) bool {
	var chance int

	if target.isAwake() {
		chance = helpers.Min(60, 2*target.Level)
	} else {
		if target.equipped("wield") == "<empty>" {
			return false
		}
		chance = 50 / 2
	}

	if dice().Intn(100) >= chance+target.Level-m.Level {
		return false
	}

	target.notify(fmt.Sprintf("You parry %s's attack.", m.Name))
	m.notify(fmt.Sprintf("%s parries your attack.", target.Name))

	return true
}

func (m *mob) die() {
	// drop corpse in room
	corpse := &item{itemType: "corpse", Name: "A corpse of " + m.Name, Identifiers: "corpse," + m.Identifiers, Decays: decays, TTL: 1}
	m.Room.Items = append(m.Room.Items, corpse)

	// whisk them away to Nowhere
	m.Room = getRoom(0)
}

func (m *mob) damroll() int {
	return m.Strength
}

func (m *mob) damage(victim *mob) int {
	var dam int
	if m.Playable == false {
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

	// if m.Playable == true && m.getSkill(sEnhancedDamage) > 0 {
	// 	dam += int(dam * m.getSkill(sEnhancedDamage) / 100)
	// }

	if victim.Status < sitting {
		dam *= 2
	}

	if dam <= 0 {
		dam = 1
	}

	return dam
}

func (m *mob) takeDamage(damage int) {
	m.Hitpoints -= damage
	if m.Hitpoints < 0 {
		m.Status = dead
		m.notify(helpers.Red + "You are DEAD!!!" + helpers.Reset)
	}
}

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

			// whisk it away
			target.die()
			m.Room.removeMob(target)
			target.Room.removeMob(target)
			m.ShowStatusBar()
			return
		}
		m.Status = fighting
		target.Status = fighting
	}
}
