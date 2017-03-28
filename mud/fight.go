package mud

import "github.com/brianseitel/oasis-mud/helpers"

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
					if !m.isNPC() /* && attacker.GroupedWith(m) */ {
						multiHit(m, attacker, typeUndefined)
					}
				}

				if m.isNPC() {
					if attacker.index == m.index || dBits(3) == 0 {
						number := 0
						var target *mob

						for _, neighbor := range m.Room.Mobs {
							if m.canSee(neighbor) /* && m.GroupedWith(neighbor) */ && dice().Intn(number+1) == 0 {
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
	// victim.makeCorpse() TODO

	if victim.isNPC() {
		// victim.extract(true) TODO
		return
	}

	// victim.extract(false) // TODO
	for _, af := range victim.Affects {
		victim.removeAffect(af)
	}

	victim.AffectedBy = 0
	victim.Armor = 100
	victim.Status = sitting
	victim.Hitpoints = helpers.Max(1, victim.Hitpoints)
	victim.Mana = helpers.Max(1, victim.Mana)
	victim.Movement = helpers.Max(1, victim.Movement)

	// victim.Save() // TODO

	return
}
