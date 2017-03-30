package mud

import (
	"container/list"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/brianseitel/oasis-mud/helpers"
)

var (
	areaList      list.List
	itemList      list.List
	itemIndexList list.List
	jobList       list.List
	mobList       *list.List
	mobIndexList  list.List
	raceList      list.List
	roomList      list.List
	shopList      list.List
	skillList     list.List
	socialList    list.List
)

func areaUpdate() {
	for e := areaList.Front(); e != nil; e = e.Next() {
		area := e.Value.(*area)
		area.age++
		if area.age < 3 {
			continue
		}

		playerCount := 0
		for _, room := range area.Rooms {
			for _, mob := range room.Mobs {
				if mob.Playable && mob.client != nil {
					playerCount++
				}
			}
		}

		if playerCount == 0 || area.age >= 15 {
			resetArea(area)
			area.age = dice().Intn(3)
		}
	}
}

func bootDB() {
	loadSkills()
	loadSocials()
	loadJobs()
	loadRaces()
	loadItems()
	loadMobs()
	loadRooms()
	loadShops()

	areaUpdate()
}

func createItem(index *itemIndex) *item {
	var contained []*item
	item := &item{Name: index.Name, Description: index.Description, ShortDescription: index.ShortDescription, ItemType: index.ItemType, ExtraFlags: index.ExtraFlags, WearFlags: index.WearFlags, Weight: index.Weight, Value: index.Value, Timer: -1}
	for _, id := range index.ContainedIDs {
		i := getItem(id)
		contained = append(contained, createItem(i))
	}
	item.container = contained
	itemList.PushBack(item)
	return item
}

func createMob(index *mobIndex) *mob {

	m := &mob{}

	m.index = index
	m.Name = index.Name
	m.Title = index.Title
	m.Description = index.Description
	m.Affects = index.Affects
	m.AffectedBy = index.AffectedBy

	m.Skills = index.Skills

	for _, i := range index.ItemIds {
		m.Inventory = append(m.Inventory, createItem(getItem(uint(i))))
	}
	for _, i := range index.EquippedIds {
		m.Equipped = append(m.Equipped, createItem(getItem(uint(i))))
	}

	m.ExitVerb = index.ExitVerb
	m.Room = getRoom(uint(index.RoomID))

	m.Hitpoints = index.Hitpoints
	m.MaxHitpoints = index.MaxHitpoints
	m.Mana = index.Mana
	m.MaxMana = index.MaxMana
	m.Movement = index.Movement
	m.MaxMovement = index.MaxMovement

	m.Armor = index.Armor
	m.Hitroll = index.Hitroll
	m.Damroll = index.Damroll

	m.Exp = index.Exp
	m.Level = index.Level
	m.Alignment = index.Alignment
	m.Practices = index.Practices
	m.Gold = index.Gold

	m.Carrying = index.Carrying
	m.CarryMax = index.CarryMax
	m.CarryWeight = index.CarryWeight
	m.CarryWeightMax = index.CarryWeightMax

	m.Job = getJob(uint(index.JobID))
	m.Race = getRace(uint(index.RaceID))
	m.Gender = index.Gender

	m.Attributes = index.Attributes
	m.ModifiedAttributes = index.ModifiedAttributes

	m.Status = index.Status
	m.Playable = index.Playable

	var skills []*mobSkill
	for _, s := range m.Skills {
		skill := getSkill(s.SkillID)
		skills = append(skills, &mobSkill{Skill: skill, SkillID: s.SkillID, Level: s.Level})
	}
	m.Skills = skills

	if m.isNPC() && (m.Room != nil && m.Room.isDark()) {
		helpers.SetBit(m.AffectedBy, affectInfrared)
	}
	return m
}

func loadItems() {
	itemFiles, _ := filepath.Glob("./data/items/*.json")

	for _, itemFile := range itemFiles {
		file, err := ioutil.ReadFile(itemFile)
		if err != nil {
			panic(err)
		}

		var list []itemIndex
		json.Unmarshal(file, &list)

		for _, it := range list {
			itemIndexList.PushBack(it)
		}

	}
}

func loadJobs() {
	var jobs []*job

	jobs = append(jobs, &job{ID: 1, Name: "Warrior", Abbr: "war"})
	jobs = append(jobs, &job{ID: 2, Name: "Mage", Abbr: "mag"})
	jobs = append(jobs, &job{ID: 3, Name: "Cleric", Abbr: "cle"})
	jobs = append(jobs, &job{ID: 4, Name: "Thief", Abbr: "thi"})
	jobs = append(jobs, &job{ID: 5, Name: "Ranger", Abbr: "ran"})
	jobs = append(jobs, &job{ID: 6, Name: "Bard", Abbr: "bar"})

	for _, j := range jobs {
		jobList.PushBack(j)
	}
}

func loadMobs() {
	mobList = list.New()

	mobFiles, _ := filepath.Glob("./data/mobs/*.json")

	for _, mobFile := range mobFiles {
		file, err := ioutil.ReadFile(mobFile)
		if err != nil {
			panic(err)
		}

		var list []*mobIndex
		err = json.Unmarshal(file, &list)
		if err != nil {
			panic(err)
		}

		for _, m := range list {
			mobIndexList.PushBack(m)
		}
	}

}

func loadRaces() {
	raceList.PushBack(&race{ID: 1, Name: "Human", Abbr: "hum"})
	raceList.PushBack(&race{ID: 2, Name: "Elf", Abbr: "elf"})
	raceList.PushBack(&race{ID: 3, Name: "Dwarf", Abbr: "dwf"})
	raceList.PushBack(&race{ID: 4, Name: "Dark Elft", Abbr: "drw"})
	raceList.PushBack(&race{ID: 5, Name: "Goblin", Abbr: "gob"})
	raceList.PushBack(&race{ID: 6, Name: "Dragon", Abbr: "drg"})
}

func loadRooms() {
	areaFiles, _ := filepath.Glob("./data/area/*.json")

	voidArea := &area{ID: 0, Name: "Limbo", age: 0, numPlayers: 0}
	void := &room{ID: 0, Exits: nil, ItemIds: nil, MobIds: nil, Name: "The Void", Description: "A dark, gaping void lies here."}

	voidArea.Rooms = append(voidArea.Rooms, void)
	areaList.PushBack(voidArea)

	for _, areaFile := range areaFiles {
		file, err := ioutil.ReadFile(areaFile)
		if err != nil {
			panic(err)
		}

		var a area
		json.Unmarshal(file, &a)
		if err != nil {
			panic(err)
		}

		area := &area{ID: a.ID, Name: a.Name, age: 0}

		for _, ro := range a.Rooms {
			ro.AreaID = int(a.ID)
			for _, i := range ro.ItemIds {
				index := getItem(uint(i))
				item := createItem(index)
				ro.Items = append(ro.Items, item)
			}

			for _, i := range ro.MobIds {
				mob := createMob(getMob(uint(i)))
				ro.Mobs = append(ro.Mobs, mob)
			}

			roomList.PushBack(ro)

			area.Rooms = append(area.Rooms, ro)
		}

		exitsList := list.New()
		for e := roomList.Front(); e != nil; e = e.Next() {
			room := e.Value.(*room)
			for j, x := range room.Exits {
				room.Exits[j] = &exit{Dir: x.Dir, Room: getRoom(x.RoomID), RoomID: x.RoomID}
			}

			exitsList.PushBack(room)
		}

		roomList = *exitsList

		areaList.PushBack(area)
	}

	for e := roomList.Front(); e != nil; e = e.Next() {
		room := e.Value.(*room)
		for _, mob := range room.Mobs {
			mob.Room = room
		}
	}
}

func loadShops() {
	shopFiles, _ := filepath.Glob("./data/shops/*.json")

	for _, shopFile := range shopFiles {
		file, err := ioutil.ReadFile(shopFile)
		if err != nil {
			panic(err)
		}

		var list []*shop
		json.Unmarshal(file, &list)

		for _, sh := range list {
			sh.keeper = getMob(uint(sh.KeeperID))
			shopList.PushBack(sh)
		}
	}
}

func loadSkills() {
	skillFiles, _ := filepath.Glob("./data/skills/*.json")

	for _, skillFile := range skillFiles {
		file, err := ioutil.ReadFile(skillFile)
		if err != nil {
			panic(err)
		}

		var list []*skill
		json.Unmarshal(file, &list)

		for _, sk := range list {
			skillList.PushBack(sk)
		}
	}
}

func loadSocials() {
	file, err := ioutil.ReadFile("./data/socials/socials.json")
	if err != nil {
		panic(err)
	}

	var list []*social
	err = json.Unmarshal(file, &list)
	if err != nil {
		panic(err)
	}

	for _, sk := range list {
		socialList.PushBack(sk)
	}
}

func resetArea(ar *area) {
	filename := helpers.ToSnake(ar.Name)

	areaFile := fmt.Sprintf("./data/areas/%s.json", filename)
	file, err := ioutil.ReadFile(areaFile)
	if err != nil {
		panic(err)
	}

	var a *area
	json.Unmarshal(file, &a)
	if err != nil {
		panic(err)
	}

	var masterArea *area
	for e := areaList.Front(); e != nil; e = e.Next() {
		area := e.Value.(*area)
		if area.ID == a.ID {
			masterArea = area
			break
		}
	}

	var tempRooms []*room
	for _, ro := range a.Rooms {
		ro.AreaID = int(a.ID)
		for _, i := range ro.ItemIds {
			index := getItem(uint(i))
			item := createItem(index)
			ro.Items = append(ro.Items, item)
		}

		for _, i := range ro.MobIds {
			mob := createMob(getMob(uint(i)))
			ro.Mobs = append(ro.Mobs, mob)
		}

		for e := roomList.Front(); e != nil; e = e.Next() {
			room := e.Value.(*room)
			if room.ID == ro.ID {
				room = ro
			}
		}

		tempRooms = append(tempRooms, ro)
	}

	masterArea.Rooms = tempRooms
	return
}

func extractMob(m *mob, pull bool) {
	if m.Room == nil {
		return
	}

	if pull {
		// TODO: m.dieFollower()
	}

	m.stopFighting(true)

	for _, i := range m.Inventory {
		extractObj(i)
	}

	m.Room.removeMob(m)

	if !pull {
		m.Room = getRoom(uint(1))
		return
	}

	if m.isNPC() {
		m.index.count--
	}

	if m.client != nil && m.client.original != nil {
		m.returnTo([]string{""})
	}

	for e := mobList.Front(); e != nil; e = e.Next() {
		wm := e.Value.(*mob)
		if wm.replyTarget == m {
			wm.replyTarget = nil
		}
	}

	for e := mobList.Front(); e != nil; e = e.Next() {
		wm := e.Value.(*mob)
		if wm == m {
			mobList.Remove(e)
		}
	}

	if m.client != nil {
		m.client.mob = nil
	}

	return
}

func extractObj(obj *item) {
	if obj.Room != nil {
		obj.Room.removeObject(obj)
	}
	if obj.carriedBy != nil {
		obj.carriedBy.removeItem(obj)
	}
	if obj.inObject != nil {
		obj.inObject.removeObject(obj)
	}

	for _, i := range obj.container {
		extractObj(i)
	}

	for e := itemList.Front(); e != nil; e = e.Next() {
		i := e.Value.(*item)
		if i == obj {
			itemList.Remove(e)
			break
		}
	}

	obj.Affected = nil
	obj.index.count--

	return
}
