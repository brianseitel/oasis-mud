package mud

import "testing"

func TestViolenceUpdate(t *testing.T) {
	player := resetTest()
	m := mockMob("foo")

	mobList.PushBack(player)
	mobList.PushBack(m)
	violenceUpdate()

	player.Fight = m
	m.Room = player.Room
	player.Room.Mobs = append(player.Room.Mobs, player)
	player.Room.Mobs = append(player.Room.Mobs, m)

	player.Status = sleeping
	violenceUpdate()

	player.Status = standing
	violenceUpdate()

	player.Status = fighting
	violenceUpdate()

	m.Status = sleeping
	violenceUpdate()

	m.Status = standing
	player.Status = standing
	violenceUpdate()
	violenceUpdate()
	violenceUpdate()
	violenceUpdate()
	violenceUpdate()
}

func TestMakeCorpse(t *testing.T) {
	gameServer.BasePath = "../"
	loadItems()
	player := resetTest()

	obj := mockObject("sword", 123)
	player.Inventory = append(player.Inventory, obj)
	player.Gold = 1000

	makeCorpse(player)

	player.Playable = false
	makeCorpse(player)
}

func TestGroupGain(t *testing.T) {
	player := resetTest()
	friend := mockPlayer("friend")

	victim := mockMob("victim")

	player.Playable = false
	groupGain(player, victim)
	player.Playable = true

	player.Room = mockRoom()
	player.Room.Mobs = append(player.Room.Mobs, player)
	player.Room.Mobs = append(player.Room.Mobs, friend)
	player.Room.Mobs = append(player.Room.Mobs, victim)

	groupGain(player, victim)

	friend.leader = player
	groupGain(player, victim)

	player.leader = friend
	friend.leader = nil
	groupGain(player, victim)
	friend.leader = player
	player.leader = nil

	friend.Level = 25
	player.Level = 1
	groupGain(player, victim)

	player.Level = 25
	friend.Level = 1
	groupGain(player, victim)

	sword := mockObject("sword", 1234)
	sword.ItemType = itemWeapon
	sword.ExtraFlags = itemAntiEvil

	player.Equipped = append(player.Equipped, sword)
	player.Alignment = -800

	groupGain(player, victim)
}

func TestMultiHit(t *testing.T) {
	player := resetTest()
	victim := mockMob("victim")

	player.Level = 50
	victim.Level = 2

	multiHit(player, victim, typeBackstab)
	multiHit(player, victim, typeUndefined)

	player.Fight = victim
	player.Playable = false
	multiHit(player, victim, typeHit)
	multiHit(player, victim, typeHit)
	multiHit(player, victim, typeHit)
	multiHit(player, victim, typeHit)
	player.Playable = true

	multiHit(player, victim, typeHit)

	second := mockSkill("second_attack")
	second.Level = 100
	player.Skills = append(player.Skills, second)
	multiHit(player, victim, typeHit)
	multiHit(player, victim, typeHit)
	multiHit(player, victim, typeHit)

	third := mockSkill("third_attack")
	third.Level = 100
	player.Skills = append(player.Skills, third)
	multiHit(player, victim, typeHit)
	multiHit(player, victim, typeHit)
	multiHit(player, victim, typeHit)
	multiHit(player, victim, typeHit)
}

func TestRawKill(t *testing.T) {
	player := resetTest()
	victim := mockMob("victim")
	rawKill(victim)

	player.Affects = append(player.Affects, mockAffect("blind"))
	rawKill(player)
}
