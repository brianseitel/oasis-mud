package mud

import (
	"fmt"
	"strings"
)

const (
	typeUndefined = -1
	typeHit       = iota
	typeBackstab
	typeKick
	typePoison
)

func violenceUpdate() {

	for e := mobList.Front(); e != nil; e = e.Next() {
		attacker := e.Value.(*mob)

		victim := attacker.Fight

		if victim == nil {
			continue
		}

		if attacker.isAwake() && attacker.Room == victim.Room {
			multiHit(attacker, victim, typeUndefined)
		} else {
			attacker.stopFighting(false)
		}

		victim = attacker.Fight

		if victim == nil {
			continue
		}

		for _, m := range attacker.Room.Mobs {

			if m.isAwake() && m.Fight == nil {
				// auto-assist other players in group

				if !attacker.isNPC() {
					if !m.isNPC() && isSameGroup(attacker, m) {
						multiHit(m, attacker, typeUndefined)
					}
				}

				if m.isNPC() {
					if attacker.index == m.index || dBits(3) == 0 {
						number := 0
						var target *mob

						for _, neighbor := range m.Room.Mobs {
							if m.canSee(neighbor) && isSameGroup(m, neighbor) && dice().Intn(number+1) == 0 {
								target = neighbor
								number++
							}
						}

						if target != nil {
							multiHit(m, target, typeUndefined)
						}
					}
				}
			}
		}
	}
}

func makeCorpse(victim *mob) {
	var name string
	var corpse *item
	if victim.isNPC() {
		name = victim.Name
		corpse = createItem(getItem(vnumCorpseNPC))
		corpse.Timer = dice().Intn(4)

		if victim.Gold >= 0 {
			corpse.container = append(corpse.container, createMoney(victim.Gold))
			victim.Gold = 0
		}
	} else {
		name = victim.Name
		corpse = createItem(getItem(vnumCorpsePC))
		corpse.Timer = dice().Intn(4)

		if victim.Gold >= 0 {
			corpse.container = append(corpse.container, createMoney(victim.Gold))
			victim.Gold = 0
		}
	}

	corpse.Name = strings.Replace(corpse.ShortDescription, "[name]", name, 1)
	corpse.Description = strings.Replace(corpse.Description, "[name]", name, 1)

	for j := range victim.Inventory {
		victim.Inventory, corpse.container = transferItem(j, victim.Inventory, corpse.container)
	}

	victim.Room.Items = append(victim.Room.Items, corpse)
}

func groupGain(player *mob, victim *mob) {
	if player.isNPC() || !victim.isNPC() || victim == player {
		return
	}

	members := 0
	for _, m := range player.Room.Mobs {
		if isSameGroup(player, m) {
			members++
		}
	}

	if members == 0 {
		members = 1
	}

	leader := player
	if player.leader != nil {
		leader = player.leader
	}

	for _, m := range player.Room.Mobs {
		if !isSameGroup(player, m) {
			continue
		}

		if m.Level-leader.Level > 5 {
			m.notify("You are too high for this group.")
			continue
		}

		if m.Level-leader.Level < -5 {
			m.notify("You are too low for this group.")
			continue
		}

		xp := xpCompute(m, victim) / members
		m.notify("You receive %d experience points.", xp)
		m.gainExp(xp)
		fmt.Println("shit", m.Name)
		for _, i := range m.Inventory {
			if i.WearLocation == wearNone {
				continue
			}

			if (i.hasExtraFlag(itemAntiEvil) && m.isEvil()) || (i.hasExtraFlag(itemAntiGood) && m.isGood()) || (i.hasExtraFlag(itemAntiNeutral) && m.isNeutral()) {
				act("You are zapped by $p.", m, i, nil, actToChar)
				act("$n is zapped by $p.", m, i, nil, actToRoom)
				m.removeItem(i)
				m.Room.Items = append(m.Room.Items, i)
			}
		}
	}
}

func multiHit(attacker *mob, victim *mob, damageType int) {
	var chance int

	attacker.oneHit(victim, damageType)

	if attacker.Fight != victim || damageType == typeBackstab {
		return
	}

	if attacker.isNPC() {
		chance = attacker.Level
	} else {
		secondAttack := attacker.skill("second_attack")
		if secondAttack != nil {
			chance = int(secondAttack.Level / 2)
		} else {
			chance = 0
		}
	}

	if dice().Intn(100) < chance {
		attacker.oneHit(victim, damageType)
		if attacker.Fight != victim {
			return
		}
	}
	if attacker.isNPC() {
		chance = attacker.Level
	} else {
		thirdAttack := attacker.skill("third_attack")
		if thirdAttack != nil {
			chance = int(thirdAttack.Level / 4)
		} else {
			chance = 0
		}
	}

	if dice().Intn(100) < chance {
		attacker.oneHit(victim, damageType)
		if attacker.Fight != victim {
			return
		}
	}

	if attacker.isNPC() {
		chance = attacker.Level
	} else {
		chance = 0
	}

	if dice().Intn(100) < chance {
		attacker.oneHit(victim, damageType)
		if attacker.Fight != victim {
			return
		}
	}

	return
}

func rawKill(victim *mob) {
	victim.stopFighting(false)
	victim.deathCry()
	makeCorpse(victim)

	if victim.isNPC() {
		extractMob(victim, true)
		return
	}

	extractMob(victim, false)
	for _, af := range victim.Affects {
		victim.removeAffect(af)
	}

	victim.AffectedBy = 0
	victim.Armor = 100
	victim.Status = sitting
	victim.Hitpoints = max(1, victim.Hitpoints)
	victim.Mana = max(1, victim.Mana)
	victim.Movement = max(1, victim.Movement)

	saveCharacter(victim)

	return
}
