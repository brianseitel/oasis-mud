package mud

import (
	"fmt"

	"github.com/brianseitel/oasis-mud/helpers"
)

func (m *mob) parry(attacker *mob) bool {
	var chance int

	if !attacker.isAwake() {
		chance = helpers.Min(60, 2*attacker.Level)
	} else {
		// if attacker.equipped("wield") == "<empty>" {
		// 	return false
		// }
		fmt.Println("Checking parry...")

		mobSkill := attacker.skill("parry")
		if mobSkill != nil {
			chance = int(mobSkill.Level / 2)
			fmt.Println("CHANCE: ", mobSkill.Level, chance)
		}
	}

	fmt.Println("Chance: ", chance+attacker.Level-m.Level)
	if dice().Intn(100) >= chance+attacker.Level-m.Level {
		fmt.Println("No Parry :(")
		return false
	}

	attacker.notify(fmt.Sprintf("You parry %s's attack.", m.Name))
	m.notify(fmt.Sprintf("%s parries your attack.", attacker.Name))

	fmt.Println("PARRIED!")
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
	if m != victim {
		if victim.parry(m) {
			return 0
		}
	}
	return m.oneHit(victim)
}

func (m *mob) oneHit(victim *mob) int {
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
