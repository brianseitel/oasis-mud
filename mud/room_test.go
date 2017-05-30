package mud

import "testing"

func TestDecayItems(t *testing.T) {
	item := mockObject("thing", 1)
	item.Timer = 3

	item2 := mockObject("thing2", 2)
	room := mockRoom()
	player := mockPlayer("player")
	room.Items = append(room.Items, item)
	room.Items = append(room.Items, item2)
	room.Mobs = append(room.Mobs, player)

	for i := 0; i < 250; i++ {
		room.decayItems()
	}
}

func TestGetRoom(t *testing.T) {

	getRoom(1)
	getRoom(932342)
}

func TestGetItem(t *testing.T) {
	loadItems()
	getItem(1)
	getItem(1234532431)
}

func TestGetMob(t *testing.T) {
	loadMobs()
	getMob(1)
	getMob(123432653242)
}

func TestIsDark(t *testing.T) {
	room := mockRoom()

	if room.isDark() {
		t.Error("Room should not be dark")
	}

	room.Light = 1234
	if room.isDark() {
		t.Error("Room should not be dark")
	}

	room.Light = 0
	room.RoomFlags = roomDark
	if !room.isDark() {
		t.Error("Room should be dark")
	}

	room.RoomFlags = 0
	room.SectorType = sectorInside
	if !room.isDark() {
		t.Error("Room should be dark")
	}
}

func TestIsPrivate(t *testing.T) {
	room := mockRoom()

	if room.isPrivate() {
		t.Error("Room should not be private.")
	}

	m := mockMob("joe")
	room.Mobs = append(room.Mobs, m)

	room.RoomFlags = roomSolitary
	if !room.isPrivate() {
		t.Error("Room should be private.")
	}

	room.Mobs = append(room.Mobs, m)
	room.Mobs = append(room.Mobs, m)
	room.RoomFlags = roomPrivate
	if !room.isPrivate() {
		t.Error("Room should be private.")
	}
}

func TestRemoveMob(t *testing.T) {
	room := mockRoom()
	m := mockMob("thing")

	room.Mobs = append(room.Mobs, m)

	if len(room.Mobs) != 1 {
		t.Error("Should only have one mob.")
	}

	room.removeMob(m)
	if len(room.Mobs) != 0 {
		t.Error("Failed to remove mob.")
	}
}

func TestRemoveObject(t *testing.T) {
	room := mockRoom()
	i := mockObject("thing", 1234)

	room.Items = append(room.Items, i)

	if len(room.Items) != 1 {
		t.Error("Should only have one item.")
	}

	room.removeObject(i)
	if len(room.Items) != 0 {
		t.Error("Failed to remove item.")
	}
}

func TestFindExit(t *testing.T) {
	r := mockRoom()

	// not existent
	r.findExit("north")
	// not valid
	r.findExit("ass")

	r.Exits = append(r.Exits, &exit{Dir: "east"})
	r.Exits = append(r.Exits, &exit{Dir: "west"})
	r.Exits = append(r.Exits, &exit{Dir: "north"})
	r.Exits = append(r.Exits, &exit{Dir: "south"})
	r.Exits = append(r.Exits, &exit{Dir: "up"})
	r.Exits = append(r.Exits, &exit{Dir: "down"})

	r.findExit("east")
	r.findExit("north")
	r.findExit("south")
	r.findExit("west")
	r.findExit("up")
	r.findExit("down")
}

func TestNotifyRoom(t *testing.T) {
	r := mockRoom()
	player := mockPlayer("player")
	vic := mockPlayer("vic")

	r.Mobs = append(r.Mobs, player)
	r.Mobs = append(r.Mobs, vic)

	r.notify("holla", player)
}
