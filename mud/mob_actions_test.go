package mud

import (
	"net"
	"testing"
)

func TestChangePosition(t *testing.T) {
	player := mockPlayer("player")

	player.changePosition(standing)
	if player.Status != standing {
		t.Error("Failed to change status to standing")
	}
	player.changePosition(sleeping)
	if player.Status != sleeping {
		t.Error("Failed to change status to sleeping")
	}
	player.changePosition(sitting)
	if player.Status != sitting {
		t.Error("Failed to change status to sitting")
	}
	player.changePosition(resting)
	if player.Status != resting {
		t.Error("Failed to change status to resting")
	}
	player.changePosition(dead)

	player.Status = dead
	player.changePosition(standing)

	player.Status = mortal
	player.changePosition(standing)

	player.Status = incapacitated
	player.changePosition(standing)

	player.Status = stunned
	player.changePosition(standing)

	player.Status = fighting
	player.changePosition(standing)

}

func TestDoClose(t *testing.T) {
	player := resetTest()
	thing := mockObject("thing1", 1)
	room := &room{Mobs: []*mob{}}
	player.Room = room

	// no args
	doClose(player, "")

	// inventory
	player.Inventory = append(player.Inventory, thing)
	doClose(player, "thing1")

	// floor
	player.Inventory = []*item{}
	player.Room.Items = append(player.Room.Items, thing)
	doClose(player, "thing1")

	// not found
	doClose(player, "thing2")

	// gotta be a container
	thing.ItemType = itemContainer
	doClose(player, "thing1")

	// can't already be closed
	thing.ClosedFlags = containerClosed
	doClose(player, "thing1")

	// gotta be closeable
	thing.ClosedFlags = 0
	doClose(player, "thing1")

	// make it closeable; should work
	thing.ClosedFlags = containerClosable
	doClose(player, "thing1")

	x := &exit{Dir: "east", Key: 1}
	x2 := &exit{Dir: "west"}
	room.Exits = append(room.Exits, x)
	room.Exits = append(room.Exits, x2)

	// Already closed
	x.Flags = exitClosed
	doClose(player, "east")

	x.Flags = exitDoor
	x.Room = room
	x.Room.Mobs = append(x.Room.Mobs, mockMob("foo"))
	doClose(player, "east")
}

func TestDoCommands(t *testing.T) {
	loadCommands()
	player := resetTest()

	doCommands(player, "")
}

func TestDoGroup(t *testing.T) {
	_ = resetTest()

	player := mockPlayer("jemaclus")
	player.Room = mockRoom()

	s, c := net.Pipe()
	s.Close()

	player2 := mockPlayer("bob")
	player2.Room = mockRoom()
	player2.client = &connection{conn: c}
	mobList.PushBack(player2)

	player3 := mockPlayer("bubba")
	player3.Room = mockRoom()
	player3.client = &connection{conn: c}
	mobList.PushBack(player3)

	player.Room.Mobs = append(player.Room.Mobs, player2)
	player.Room.Mobs = append(player.Room.Mobs, player3)

	player2.leader = player
	doGroup(player, "")

	doGroup(player, "rick")

	player2.Playable = false
	doGroup(player, "bob")
	player2.Playable = true

	player2.leader = nil
	player.master = player3
	player.leader = player3
	doGroup(player, "bob")

	player.master = nil
	player.leader = nil
	player2.master = player3
	doGroup(player, "bob")

	player2.leader = player
	player2.master = nil
	doGroup(player, "bob")

	player2.Level = 50
	player.Level = 1
	doGroup(player, "bob")

	player.Level = 50
	doGroup(player, "bob")
}

func TestDoHide(t *testing.T) {
	player := mockPlayer("player")
	player.AffectedBy = affectHide

	hide := &skill{Name: "hide"}
	playerSkill := &mobSkill{Level: 100, Skill: hide}
	player.Skills = append(player.Skills, playerSkill)

	doHide(player, "")

	player.Playable = false
	doHide(player, "")
}

func TestDoFollow(t *testing.T) {
	player := resetTest()
	player2 := mockPlayer("bob")

	player.Room = mockRoom()
	player2.Room = mockRoom()
	player.Room.Mobs = append(player.Room.Mobs, player)
	player.Room.Mobs = append(player.Room.Mobs, player2)

	doFollow(player, "")

	doFollow(player, "rick")

	doFollow(player, "player")

	player.master = player
	doFollow(player, "player")

	player.Level = 50
	player2.Level = 1
	doFollow(player, "bob")

	player.Level = 1
	player.master = player2
	doFollow(player, "bob")
}

func TestDoLock(t *testing.T) {
	player := resetTest()

	player.Room = mockRoom()

	nextRoom := mockRoom()

	x := &exit{Dir: "east", Room: nextRoom, Flags: 0, Key: 1}
	x2 := &exit{Dir: "west", Room: player.Room, Flags: 0, Key: 1}
	nextRoom.Exits = append(nextRoom.Exits, x2)
	player.Room.Exits = append(player.Room.Exits, x)

	// no arg
	doLock(player, "")

	// Closed
	doLock(player, "east")

	// closed but not locked
	x.Flags = exitClosed
	doLock(player, "east")

	// No key, so can't be locked
	x.Key = -1
	doLock(player, "east")

	// doesn't have key
	x.Key = 1
	doLock(player, "east")

	// Give em key
	player.Inventory = append(player.Inventory, &item{ID: 1})

	// already locked
	x.Flags = exitLocked
	doLock(player, "east")

	x.Flags = exitClosed
	doLock(player, "east")
}

func TestDoMove(t *testing.T) {
	player := resetTest()

	player.Room = mockRoom()
	x := &exit{Dir: "east", Room: mockRoom()}
	player.Room.Exits = append(player.Room.Exits, x)

	player.Status = fighting
	doMove(player, "east")

	player.Status = standing
	doMove(player, "east")

	doMove(player, "nowhere")
}

func TestDoPick(t *testing.T) {
	player := resetTest()

	player.Room = mockRoom()

	nextRoom := mockRoom()

	x := &exit{Dir: "east", Room: nextRoom, Flags: 0, Key: -1}
	x2 := &exit{Dir: "west", Room: player.Room, Flags: 0, Key: 1}
	nextRoom.Exits = append(nextRoom.Exits, x2)
	player.Room.Exits = append(player.Room.Exits, x)

	guard := mockMob("guard")

	player.Room.Mobs = append(player.Room.Mobs, guard)
	player.Level = 10
	guard.Level = player.Level + 10

	// no arg
	doPick(player, "")

	// no skill
	doPick(player, "east")

	pick := &mobSkill{Skill: &skill{Name: "pick"}}
	player.Skills = append(player.Skills, pick)

	// too close to guard
	doPick(player, "east")
	guard.Level = 1

	// Fails random check
	doPick(player, "east")

	pick.Level = 100

	// Closed
	x.Flags = 0
	doPick(player, "east")

	// Closed and no key
	x.Flags = exitClosed
	x.Key = -1
	doPick(player, "east")

	// Unlocked with key
	x.Key = 1
	doPick(player, "east")

	x.Flags = exitClosed + exitLocked
	doPick(player, "east")
}

func TestDoPractice(t *testing.T) {
	player := resetTest()
	loadSkills()
	loadBonuses()

	// can't be npc
	player.Playable = false
	doPractice(player, "")

	// must be level 3
	player.Playable = true
	player.Level = 1
	doPractice(player, "")

	// no args
	player.Level = 5
	player.Job = getJob(1)
	doPractice(player, "")
	player.Job = getJob(2)
	doPractice(player, "")
	player.Job = getJob(3)
	doPractice(player, "")
	player.Job = getJob(4)
	doPractice(player, "")
	player.Job = getJob(5)
	doPractice(player, "")
	player.Job = getJob(6)
	doPractice(player, "")

	// has a skill
	trip := &mobSkill{Skill: &skill{Name: "trip"}, Level: 15}
	player.Skills = append(player.Skills, trip)
	doPractice(player, "")

	// practice specific
	player.Job = getJob(1) // warrior

	// can't while sleeping
	player.Status = sleeping
	doPractice(player, "trip")

	// no trainer
	player.Status = standing
	doPractice(player, "trip")

	trainer := mockMob("trainer")
	trainer.Act = actPractice
	player.Room.Mobs = append(player.Room.Mobs, trainer)

	// no practices
	player.Practices = 0
	doPractice(player, "trip")

	player.Practices = 10000

	// not a skill
	doPractice(player, "something")

	// Not high enough level
	doPractice(player, "meteor")
	player.Job = getJob(1)
	doPractice(player, "meteor")
	player.Job = getJob(2)
	doPractice(player, "meteor")
	player.Job = getJob(3)
	doPractice(player, "meteor")
	player.Job = getJob(4)
	doPractice(player, "meteor")
	player.Job = getJob(5)
	doPractice(player, "meteor")
	player.Job = getJob(6)
	doPractice(player, "meteor")

	// Works
	player.Level = 99
	player.Job = getJob(1)
	doPractice(player, "meteor")

	meteor := player.skill("meteor")
	meteor.Level = 100

	// already mastered
	doPractice(player, "meteor")

	player.Attributes.Intelligence = 99
	meteor.Level = 84
	doPractice(player, "meteor")
}

func TestDoOpen(t *testing.T) {
	player := resetTest()
	thing := mockObject("thing1", 1)
	room := &room{Mobs: []*mob{}}
	player.Room = room

	// no args
	doOpen(player, "")

	// inventory
	player.Inventory = append(player.Inventory, thing)
	doOpen(player, "thing1")

	// floor
	player.Inventory = []*item{}
	player.Room.Items = append(player.Room.Items, thing)
	doOpen(player, "thing1")

	// not found
	doOpen(player, "thing2")

	// gotta be a container
	thing.ItemType = itemContainer
	doOpen(player, "thing1")

	// can't already be Opend
	thing.ClosedFlags = 0
	doOpen(player, "thing1")

	// make it Openable; should work
	thing.ClosedFlags = containerClosable
	doOpen(player, "thing1")

	x := &exit{Dir: "east", Key: 1, Flags: exitLocked}
	x2 := &exit{Dir: "west"}
	room.Exits = append(room.Exits, x)
	room.Exits = append(room.Exits, x2)

	// Already Opened
	x.Flags = 0
	doOpen(player, "east")

	x.Flags = exitDoor
	x.Room = room
	x.Room.Mobs = append(x.Room.Mobs, mockMob("foo"))
	doOpen(player, "east")
}

func TestDoQui(t *testing.T) {
	player := resetTest()

	doQui(player, "")
}

func TestDoQuit(t *testing.T) {
	player := resetTest()

	player.Status = fighting
	doQuit(player, "")

	player.Status = standing
	doQuit(player, "")
}

func TestDoRecall(t *testing.T) {
	player := resetTest()
	room := mockRoom()
	player.Room = room
	roomList.PushBack(room)

	doRecall(player, "")

	player.RecallRoomID = 1

	doRecall(player, "set")

	doRecall(player, "what")
}

func TestDoRest(t *testing.T) {
	player := resetTest()

	player.Status = resting
	doRest(player, "")

	player.Status = sleeping
	doRest(player, "")

	player.Status = standing
	doRest(player, "")

	player.Status = fighting
	doRest(player, "")
}

func TestDoSave(t *testing.T) {
	t.Skip()
}

func TestDoSleep(t *testing.T) {
	player := resetTest()

	player.Status = resting
	doSleep(player, "")

	player.Status = sleeping
	doSleep(player, "")

	player.Status = standing
	doSleep(player, "")

	player.Status = fighting
	doSleep(player, "")
}

func TestDoSneak(t *testing.T) {
	player := resetTest()

	// no skill
	doSneak(player, "")

	skill := &mobSkill{Skill: &skill{Name: "sneak"}, Level: 0}
	player.Skills = append(player.Skills, skill)

	doSneak(player, "")
	doSneak(player, "")
}

func TestDoSteal(t *testing.T) {
	player := resetTest()
	vic := mockMob("vic")
	player.Room.Mobs = append(player.Room.Mobs, player)
	player.Room.Mobs = append(player.Room.Mobs, vic)

	// no skill
	doSteal(player, "")

	skill := &mobSkill{Skill: &skill{Name: "steal"}, Level: 0}
	player.Skills = append(player.Skills, skill)

	// no args
	doSteal(player, "")

	// not here
	doSteal(player, "gold bob")

	// can't steal from yourself
	doSteal(player, "gold player")

	// Not gonna happen
	vic.Level = player.Level + 25
	doSteal(player, "gold vic")

	// sleeping
	vic.Status = sleeping
	doSteal(player, "gold vic")

	// thief is NPC
	player.Playable = false
	doSteal(player, "gold vic")

	// mark as thief
	player.Playable = true
	doSteal(vic, "gold player")

	// Should work
	skill.Level = 100
	vic.Status = standing
	vic.Level = player.Level
	vic.Gold = 10000
	doSteal(player, "gold vic")

	// steal something they don't have
	doSteal(player, "foo vic")

	// steal something they can't drop
	obj := mockObject("fire", 325)
	obj.ExtraFlags = itemNoDrop
	vic.Inventory = append(vic.Inventory, obj)
	doSteal(player, "fire vic")

	// Can't carry
	obj.ExtraFlags = 0
	player.CarryMax = 100
	player.Carrying = 100
	doSteal(player, "fire vic")

	// Can't carry weight
	player.Carrying = 0
	obj.Weight = 100
	player.CarryWeight = 100
	player.CarryWeightMax = 100
	doSteal(player, "fire vic")

	// Do it
	obj.Weight = 1
	player.CarryWeight = 0
	player.CarryWeightMax = 100
	doSteal(player, "fire vic")
}

func TestDoTrain(t *testing.T) {
	player := resetTest()

	// no npcs
	player.Playable = false
	doTrain(player, "")

	player.Playable = true

	// no trainer
	doTrain(player, "")

	trainer := mockMob("trainer")
	player.Room.Mobs = append(player.Room.Mobs, trainer)

	// trainer is present, no args
	doTrain(player, "")

	player.Practices = 10000

	doTrain(player, "str")
	doTrain(player, "wis")
	doTrain(player, "int")
	doTrain(player, "dex")
	doTrain(player, "cha")
	doTrain(player, "con")
	doTrain(player, "something")

	player.Attributes.Strength = 21
	doTrain(player, "str")

	player.Attributes.Wisdom = 21
	doTrain(player, "wis")

	player.Attributes.Intelligence = 21
	doTrain(player, "int")

	player.Attributes.Dexterity = 21
	doTrain(player, "dex")

	player.Attributes.Charisma = 21
	doTrain(player, "cha")

	player.Attributes.Constitution = 21
	doTrain(player, "con")

	player.Attributes.Constitution = 12
	player.Practices = 0
	doTrain(player, "con")
}

func TestDoUnlock(t *testing.T) {
	player := resetTest()

	player.Room = mockRoom()

	nextRoom := mockRoom()

	x := &exit{Dir: "east", Room: nextRoom, Flags: 0, Key: 1}
	x2 := &exit{Dir: "west", Room: player.Room, Flags: 0, Key: 1}
	nextRoom.Exits = append(nextRoom.Exits, x2)
	player.Room.Exits = append(player.Room.Exits, x)

	// no arg
	doUnlock(player, "")

	// Closed
	doUnlock(player, "east")

	// closed but not locked
	x.Flags = exitClosed
	doUnlock(player, "east")

	// No key, so can't be locked
	x.Key = -1
	doUnlock(player, "east")

	// doesn't have key
	x.Key = 1
	doUnlock(player, "east")

	// Give em key
	player.Inventory = append(player.Inventory, &item{ID: 1})

	// already locked
	x.Flags = 0
	doUnlock(player, "east")

}

func TestDoVisible(t *testing.T) {
	player := resetTest()

	doVisible(player, "")
}

func TestDoWake(t *testing.T) {
	player := resetTest()

	player.Status = resting
	doWake(player, "")

	player.Status = sleeping
	doWake(player, "")

	player.Status = standing
	doWake(player, "")

	player.Status = fighting
	doWake(player, "")
}
