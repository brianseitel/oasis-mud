package mud

import (
	"testing"
)

func TestDoAffact(t *testing.T) {
	player := resetTest()

	aff1 := mockAffect("seconds")
	aff1.duration = 1
	aff2 := mockAffect("minute")
	aff2.duration = 31
	aff3 := mockAffect("minutes")
	aff3.duration = 61
	aff4 := mockAffect("a while")
	aff4.duration = 151
	aff5 := mockAffect("a long time")
	aff5.duration = 301
	aff6 := mockAffect("practically forever")
	aff6.duration = 600

	player.Affects = append(player.Affects, aff1)
	player.Affects = append(player.Affects, aff2)
	player.Affects = append(player.Affects, aff3)
	player.Affects = append(player.Affects, aff4)
	player.Affects = append(player.Affects, aff5)
	player.Affects = append(player.Affects, aff6)

	doAffect(player, "")
}

func TestDoAreas(t *testing.T) {
	player := resetTest()
	gameServer.BasePath = "../"
	ar := &area{}

	areaList.PushBack(ar)
	doAreas(player, "")

	areaList.Init()
}

func TestDoCompare(t *testing.T) {
	player := resetTest()

	sword := mockObject("sword", 1)
	sword.ItemType = itemWeapon
	shield := mockObject("shield", 2)
	shield.ItemType = itemWeapon

	player.Inventory = append(player.Inventory, sword)
	player.Inventory = append(player.Inventory, shield)

	doCompare(player, "")

	doCompare(player, "bottle whatever")

	doCompare(player, "sword bottle")

	doCompare(player, "sword sword")

	doCompare(player, "sword shield")

	doCompare(player, "sword")

	sword.ItemType = itemArmor
	shield.ItemType = itemArmor
	doCompare(player, "sword shield")

	sword.ItemType = itemBoat
	shield.ItemType = itemBoat
	doCompare(player, "sword shield")

	shield.ItemType = itemBoat
	doCompare(player, "sword shield")

	sword.ItemType = itemWeapon
	shield.ItemType = itemArmor
	doCompare(player, "sword shield")

	shield.ItemType = itemWeapon

	sword.Min = 1
	shield.Min = 2
	doCompare(player, "sword shield")

	sword.Min = 2
	shield.Min = 1
	doCompare(player, "sword shield")
}

func TestDoConsider(t *testing.T) {
	player := resetTest()
	victim := mockMob("victim")

	player.Room.Mobs = append(player.Room.Mobs, player)
	player.Room.Mobs = append(player.Room.Mobs, victim)

	doConsider(player, "")

	doConsider(player, "foo")

	doConsider(player, "victim")

	for i := 0; i < 25; i++ {
		player.Level = 25 - i
		victim.Level = i
		doConsider(player, "victim")
	}
}

func TestDoExamine(t *testing.T) {
	player := resetTest()
	obj := mockObject("thing", 1)
	obj.ItemType = itemContainer

	obj2 := mockObject("whatever", 2)
	obj2.ItemType = itemWeapon
	player.Room.Items = append(player.Room.Items, obj2)

	player.Inventory = append(player.Inventory, obj)

	doExamine(player, "")

	doExamine(player, "nothing")

	obj.ItemType = itemDrinkContainer
	doExamine(player, "thing")
	obj.ItemType = itemContainer
	doExamine(player, "thing")
	obj.ItemType = itemCorpseNPC
	doExamine(player, "thing")
	obj.ItemType = itemCorpsePC
	doExamine(player, "thing")

	doExamine(player, "whatever")
}

func TestDoHelp(t *testing.T) {
	player := resetTest()
	player.Trust = 100
	h := &help{Keyword: "summary", Level: 10}
	helpList.PushBack(h)

	doHelp(player, "")

	h.Level = 1000
	doHelp(player, "")

	doHelp(player, "foo")
}

func TestDoInventory(t *testing.T) {
	player := resetTest()

	doInventory(player, "")
}

func TestDoLook(t *testing.T) {
	player := resetTest()

	doLook(player, "")

	player.Status = sleeping
	doLook(player, "")

	player.Status = standing
	player.AffectedBy = affectBlind
	doLook(player, "")
	player.AffectedBy = 0

	player.Room.Light = 0
	player.Room.SectorType = sectorInside
	doLook(player, "")

	// reset room
	player.Room = mockRoom()
	doLook(player, "")

	// look in
	doLook(player, "in")

	box := mockObject("box", 123)
	player.Inventory = append(player.Inventory, box)

	// not a thing
	doLook(player, "in bag")

	// not a container
	doLook(player, "in box")

	box.ItemType = itemCorpsePC
	box.ClosedFlags = containerClosed
	doLook(player, "in box")

	box.ItemType = itemCorpseNPC
	doLook(player, "in box")

	box.ItemType = itemContainer
	doLook(player, "in box")

	orc := mockMob("orc")
	player.Room.Mobs = append(player.Room.Mobs, player)
	player.Room.Mobs = append(player.Room.Mobs, orc)

	// look at mob
	doLook(player, "orc")

	// look at item in room
	sword := mockObject("sword", 543)
	sword.Name = "sword"
	player.Room.Items = append(player.Room.Items, sword)
	doLook(player, "sword")

	// look at exit
	x := &exit{Dir: "east", Keyword: "car", Room: mockRoom()}
	player.Room.Exits = append(player.Room.Exits, x)
	doLook(player, "car")
	doLook(player, "east")
	doLook(player, "north")
	doLook(player, "west")
	doLook(player, "south")
	doLook(player, "down")
	doLook(player, "up")

	x.Description = "stuff"
	doLook(player, "east")

	x.Key = 1
	doLook(player, "east")
	x.Flags = exitClosed
	doLook(player, "east")
	x.Flags = exitDoor
	doLook(player, "east")
}

func TestDoScan(t *testing.T) {
	player := resetTest()
	player.Room = mockRoom()

	doScan(player, "")

	x := &exit{Dir: "east", Room: mockRoom()}
	player.Room.Exits = append(player.Room.Exits, x)
	doScan(player, "")

	x.Dir = ""
	doScan(player, "")

	x.Dir = "west"
	x.Flags = exitClosed
	doScan(player, "")
	m := mockMob("foo")
	x.Flags = 0
	x.Room.Mobs = append(x.Room.Mobs, m)
	doScan(player, "")
}

func TestDoScore(t *testing.T) {
	player := resetTest()

	doScore(player, "")
}

func TestDoSkills(t *testing.T) {
	player := resetTest()
	trip := mockSkill("trip")

	player.Skills = append(player.Skills, trip)
	doSkills(player, "")
}

func TestDoWhere(t *testing.T) {
	player := resetTest()

	doWhere(player, "")

	orc := mockPlayer("orc")
	orc.Room = player.Room
	orc.Room.Area = player.Room.Area
	mobList.PushBack(orc)

	player.Room.Mobs = append(player.Room.Mobs, orc)
	doWhere(player, "")

	doWhere(player, "orc")
	doWhere(player, "bob")
}

func TestDoWho(t *testing.T) {
	player := resetTest()
	bob := mockPlayer("bob")

	loadJobs()

	mobList.PushBack(player)
	mobList.PushBack(bob)

	doWho(player, "")

	doWho(player, "50")
	doWho(player, "1 foot")
	doWho(player, "foot 1")
	doWho(player, "war 1")

	doWho(player, "1 25")
	bob.AffectedBy = affectInvisible
	doWho(player, "1 25")
	doWho(player, "25 35")

	bob.AffectedBy = 0
	bob.Level = 99
	doWho(player, "")
	bob.Level = 98
	doWho(player, "")
	bob.Level = 97
	doWho(player, "")
	bob.Level = 96
	doWho(player, "")

}

func TestExitsString(t *testing.T) {
	var exits []*exit
	exits = append(exits, &exit{Dir: "east"})
	exitsString(exits)
}

func TestItemsString(t *testing.T) {
	var items []*item
	items = append(items, mockObject("stick", 123))
	itemsString(items)
}

func TestMobsString(t *testing.T) {
	player := resetTest()
	var mobs []*mob
	orc := mockMob("orc")
	mobs = append(mobs, orc)
	mobsString(mobs, player)

	orc.AffectedBy = affectInvisible
	mobsString(mobs, player)
}

func TestInventoryString(t *testing.T) {
	player := resetTest()

	sword := mockObject("sword", 123)
	box := mockObject("box", 234)
	player.Inventory = append(player.Inventory, sword)
	player.Inventory = append(player.Inventory, sword)
	player.Inventory = append(player.Inventory, box)
	inventoryString(player)
}

func TestEquippedString(t *testing.T) {
	player := resetTest()
	equippedString(player)
}
