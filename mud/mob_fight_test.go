package mud

import (
	"testing"
)

func TestDoBackstab(t *testing.T) {
	player := resetTest()
	player.Name = "player"
	victim := mockPlayer("victim")
	victim.Hitpoints = 10000
	victim.MaxHitpoints = 10000
	player.Room = mockRoom()
	player.Room.Mobs = append(player.Room.Mobs, player)
	player.Room.Mobs = append(player.Room.Mobs, victim)

	// can't backstab
	doBackstab(player, "")

	bs := &mobSkill{Skill: &skill{Name: "backstab"}, Level: 50}
	player.Skills = append(player.Skills, bs)

	// has skill, no args
	doBackstab(player, "")

	// not here
	doBackstab(player, "joe")

	// self
	doBackstab(player, "player")

	// no wield
	doBackstab(player, "victim")

	wield := mockObject("sword", 123)
	wield.ItemType = itemWeapon
	wield.WearLocation = itemWearWield

	player.Equipped = append(player.Equipped, wield)

	// already fighting
	victim.Fight = player
	doBackstab(player, "victim")

	// safe
	victim.Act = actGuardian
	doBackstab(player, "victim")
	victim.Act = 0

	// Wounded
	victim.Fight = nil
	victim.Hitpoints--
	doBackstab(player, "victim")
	victim.Hitpoints++

	// sleeping
	victim.Status = sleeping
	doBackstab(player, "victim")

	// standing
	victim.Status = standing
	victim.Hitpoints = victim.MaxHitpoints
	doBackstab(player, "victim")

	bs.Level = 1000
	victim.Hitpoints = victim.MaxHitpoints
	doBackstab(player, "victim")
}

func TestDoDisarm(t *testing.T) {
	player := resetTest()
	victim := mockPlayer("victim")

	// no skill
	doDisarm(player, "")

	// Add skill
	disarm := &mobSkill{Skill: &skill{Name: "disarm"}, Level: 1}
	player.Skills = append(player.Skills, disarm)

	// No wield
	doDisarm(player, "")

	// wield
	wield := mockObject("sword", 123)
	wield.ItemType = itemWeapon
	wield.WearLocation = itemWearWield

	player.Equipped = append(player.Equipped, wield)

	// no target
	doDisarm(player, "")

	// target has no weapon
	player.Fight = victim
	doDisarm(player, "")

	// give them wield
	victim.Equipped = append(victim.Equipped, wield)

	// fail
	disarm.Level = 1
	doDisarm(player, "")

	disarm.Level = 100
	player.Level = 99
	victim.Level = 1
	doDisarm(player, "")
}

func TestDoFlee(t *testing.T) {
	player := resetTest()
	victim := mockPlayer("victim")
	r := mockRoom()

	x := &exit{Dir: "north", Room: mockRoom()}
	player.Room = r
	player.Room.Exits = append(player.Room.Exits, x)

	// not fighting
	doFlee(player, "")

	player.Fight = victim
	player.Status = fighting

	// no exits
	player.Room = mockRoom()
	doFlee(player, "")

	player.Room = r

	// same room
	x.Room = r
	doFlee(player, "")

	x.Room = mockRoom()

	// should work
	doFlee(player, "")
	player.Room = r
	doFlee(player, "")
	player.Room = r
	doFlee(player, "")
	player.Room = r
	doFlee(player, "")

	// door closed
	x.Flags = exitClosed
	player.Room = r
	doFlee(player, "")
	player.Room = r
	doFlee(player, "")
	player.Room = r
	doFlee(player, "")
	player.Room = r
	doFlee(player, "")
}

func TestDoKick(t *testing.T) {
	player := resetTest()
	victim := mockPlayer("victim")

	// no skill
	doKick(player, "")

	kick := mockSkill("kick")
	player.Skills = append(player.Skills, kick)

	// no victim
	doKick(player, "")

	// has victim
	player.Fight = victim
	doKick(player, "")

	// succed
	kick.Level = 100
	doKick(player, "")
}

func TestDoKill(t *testing.T) {
	player := resetTest()
	player.Name = "player"
	victim := mockPlayer("victim")

	player.Room = mockRoom()
	player.Room.Mobs = append(player.Room.Mobs, victim)
	player.Room.Mobs = append(player.Room.Mobs, player)

	// no args
	doKill(player, "")

	// nobody
	doKill(player, "nobody")

	// players can't be attacked
	victim.Playable = true
	doKill(player, "victim")
	victim.Playable = false

	// attack yourself
	doKill(player, "player")

	// victim is safe
	victim.Act = actGuardian
	doKill(player, "victim")
	victim.Act = 0

	// already fighting
	player.Status = fighting
	doKill(player, "victim")

	// should work
	player.Status = standing
	doKill(player, "victim")
}

func TestDoRescue(t *testing.T) {
	player := resetTest()
	victim := mockPlayer("victim")
	enemy := mockMob("mob")

	player.Room = mockRoom()
	player.Room.Mobs = append(player.Room.Mobs, victim)
	player.Room.Mobs = append(player.Room.Mobs, player)
	victim.Playable = false

	// no args
	doRescue(player, "")

	// no vic
	doRescue(player, "nobody")

	// self
	doRescue(player, "player")

	// can't rescue NPCs
	doRescue(player, "victim")
	victim.Playable = true

	// can't rescue someone you're fighting
	player.Fight = victim
	doRescue(player, "victim")
	player.Fight = nil

	// can't rescue someone who isn't fighting
	victim.Fight = nil
	doRescue(player, "victim")
	victim.Fight = enemy

	// don't know how to rescue
	doRescue(player, "victim")

	rescue := mockSkill("rescue")
	player.Skills = append(player.Skills, rescue)

	rescue.Level = 1

	// fail
	doRescue(player, "victim")

	// success
	rescue.Level = 100
	doRescue(player, "victim")

}

func TestDamage(t *testing.T) {
	player := resetTest()
	player.Room = mockRoom()

	m := mockMob("mob")
	m.Room = mockRoom()
	m.MaxHitpoints = 10000000
	m.Hitpoints = m.MaxHitpoints

	player.Room.Mobs = append(player.Room.Mobs, player)
	player.Room.Mobs = append(player.Room.Mobs, m)
	// already dead
	m.Status = dead
	player.damage(m, 1, typeUndefined)

	// damage > 1000
	m.Status = standing
	player.damage(m, 100000, typeUndefined)

	m.Status = stunned
	player.damage(m, 1, typeUndefined)

	m.Status = resting
	player.damage(m, 1, typeUndefined)

	// safe
	m.Act = actGuardian
	player.damage(m, 1, typeUndefined)
	m.Act = 0

	// fade them in if invis
	player.AffectedBy = affectInvisible
	player.damage(m, 1, typeUndefined)

	// has sanctuary
	m.AffectedBy = affectSanctuary
	player.damage(m, 1, typeUndefined)

	// has protect
	m.AffectedBy = affectProtect
	player.Alignment = -800
	player.damage(m, 1, typeUndefined)

	// no damage
	player.damage(m, -5, typeUndefined)

	// type > kick
	player.damage(m, 1, typeKick)

	// dodge and parry
	m.Level = 100
	player.Level = 1
	parry := mockSkill("parry")
	parry.Level = 100
	dodge := mockSkill("dodge")
	dodge.Level = 100
	m.Skills = append(m.Skills, parry)
	m.Skills = append(m.Skills, dodge)

	m.Hitpoints = 100
	m.MaxHitpoints = 100

	player.damage(m, 80, typeUndefined)

	m.Hitpoints = 40
	player.damage(m, 30, typeUndefined)

	m.Status = standing
	m.Hitpoints = 8
	m.Playable = true
	m.Exp = m.Level * 1000 * 5
	for i := 0; i < 26; i++ {
		player.damage(m, 1, typeUndefined)
	}

	m.Status = standing
	m.Hitpoints = 8
	m.Playable = false
	m.Act = actWimpy
	for i := 0; i < 7; i++ {
		player.damage(m, 1, typeUndefined)
	}
}

func TestDamageMessage(t *testing.T) {
	player := resetTest()
	victim := mockMob("vic")
	victim.Hitpoints = 10000
	victim.MaxHitpoints = 10000

	player.damageMessage(victim, 0, typeUndefined)
	player.damageMessage(victim, 4, typeUndefined)
	player.damageMessage(victim, 8, typeUndefined)
	player.damageMessage(victim, 12, typeUndefined)
	player.damageMessage(victim, 16, typeUndefined)
	player.damageMessage(victim, 20, typeUndefined)
	player.damageMessage(victim, 24, typeUndefined)
	player.damageMessage(victim, 28, typeUndefined)
	player.damageMessage(victim, 32, typeUndefined)
	player.damageMessage(victim, 36, typeUndefined)
	player.damageMessage(victim, 40, typeUndefined)
	player.damageMessage(victim, 44, typeUndefined)
	player.damageMessage(victim, 48, typeUndefined)
	player.damageMessage(victim, 52, typeUndefined)
	player.damageMessage(victim, 100, typeUndefined)
	player.damageMessage(victim, 1001, typeUndefined)

	player.damageMessage(victim, 1, typeHit)
	player.damageMessage(victim, 1, typeKick)
}

func TestDamroll(t *testing.T) {
	player := resetTest()

	player.damroll()
}

func TestDeathCry(t *testing.T) {
	player := resetTest()
	player.Room = mockRoom()
	x := &exit{Dir: "east", Room: player.Room}
	player.Room.Exits = append(player.Room.Exits, x)

	for i := 0; i < 10000; i++ {
		player.deathCry()
	}

	player.Playable = false
	for i := 0; i < 10000; i++ {
		player.deathCry()
	}

}

func TestDisarm(t *testing.T) {
	player := resetTest()
	victim := mockMob("victim")

	player.Room = mockRoom()
	victim.Room = mockRoom()

	player.disarm(victim)

	sword := mockObject("sword", 123)
	sword.ItemType = itemWeapon
	sword.WearLocation = itemWearWield

	victim.Equipped = append(victim.Equipped, sword)
	player.disarm(victim)
	player.disarm(victim)
	player.disarm(victim)
	player.disarm(victim)
	player.disarm(victim)
	player.disarm(victim)

	player.Equipped = append(player.Equipped, sword)
	player.disarm(victim)
	player.disarm(victim)
	player.disarm(victim)
	player.disarm(victim)
	player.disarm(victim)

}

func TestDodge(t *testing.T) {
	player := resetTest()
	victim := mockMob("victim")

	victim.Status = sleeping
	victim.dodge(player)

	victim.Status = standing
	victim.dodge(player)
	victim.dodge(player)

	victim.Level = 100
	victim.dodge(player)
	victim.dodge(player)
	victim.dodge(player)

	victim.Playable = true
	dodge := mockSkill("dodge")
	dodge.Level = 100
	victim.Skills = append(victim.Skills, dodge)

	victim.dodge(player)
	victim.dodge(player)
	victim.dodge(player)
}

func TestParry(t *testing.T) {
	player := resetTest()
	victim := mockMob("victim")

	victim.Status = sleeping
	victim.parry(player)

	victim.Status = standing
	victim.parry(player)
	victim.parry(player)

	victim.Level = 100
	victim.parry(player)
	victim.parry(player)
	victim.parry(player)

	victim.Playable = true
	parry := mockSkill("parry")
	parry.Level = 100
	victim.Skills = append(victim.Skills, parry)
	player.Skills = append(player.Skills, parry)

	sword := mockObject("sword", 123)
	sword.ItemType = itemWeapon
	sword.WearLocation = itemWearWield

	victim.Equipped = append(victim.Equipped, sword)

	victim.parry(player)
}

func TestStopFighting(t *testing.T) {
	player := resetTest()
	victim := mockMob("victim")

	mobList.PushBack(player)
	mobList.PushBack(victim)

	player.stopFighting(false)
	player.stopFighting(true)
}

func TestTakeDamage(t *testing.T) {
	player := resetTest()
	player.Hitpoints = 100

	player.takeDamage(5)
	if player.Hitpoints != 95 {
		t.Error("Failed to take damage.")
	}

	player.takeDamage(99)
}

func TestTrip(t *testing.T) {
	player := resetTest()
	victim := mockMob("what")

	// not fighting
	player.trip()

	// victim
	player.Fight = victim
	player.trip()
	player.trip()
	player.trip()

	player.Level = 100
	player.trip()
	player.trip()
	player.trip()

	trip := mockSkill("trip")
	player.Skills = append(player.Skills, trip)

	player.Playable = false
	victim.wait = 0
	player.trip()
	player.trip()
	player.trip()
	player.trip()

	player.Playable = true
	victim.wait = 0
	player.trip()
	player.trip()
	player.trip()

}

func TestUpdateStatus(t *testing.T) {
	player := resetTest()

	player.Hitpoints = 10
	player.Status = incapacitated
	player.updateStatus()
	if player.Status != standing {
		t.Error("Player should be standing")
	}

	player.Hitpoints = -11
	player.updateStatus()
	if player.Status != dead {
		t.Error("Player should be dead")
	}

	player.Hitpoints = -6
	player.updateStatus()
	if player.Status != mortal {
		t.Error("Player should be mortally wounded")
	}

	player.Hitpoints = -3
	player.updateStatus()
	if player.Status != incapacitated {
		t.Error("Player should be incapacitated")
	}

	player.Hitpoints = -1
	player.updateStatus()
	if player.Status != stunned {
		t.Error("Player should be stunned")
	}
}

func TestOneHit(t *testing.T) {
	player := resetTest()
	victim := mockMob("victim")

	victim.Status = dead
	player.oneHit(victim, typeHit)

	victim.Status = standing

	sword := mockObject("sword", 123)
	sword.ItemType = itemWeapon
	sword.WearLocation = itemWearWield

	player.Equipped = append(player.Equipped, sword)
	player.oneHit(victim, typeUndefined)

	player.Playable = false
	player.oneHit(victim, typeHit)
	player.oneHit(victim, typeHit)
	player.oneHit(victim, typeHit)
	player.oneHit(victim, typeHit)

	player.Playable = true
	victim.AffectedBy = affectDetectInvisible
	player.oneHit(victim, typeUndefined)

	player.oneHit(victim, typeUndefined)
	player.oneHit(victim, typeUndefined)
	player.oneHit(victim, typeUndefined)
	player.oneHit(victim, typeUndefined)
	player.oneHit(victim, typeUndefined)

}
