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
	skillList     list.List
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
			area.age = 0
		}
	}
}

func bootDB() {
	loadSkills()
	loadJobs()
	loadRaces()
	loadItems()
	loadMobs()
	loadRooms()

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

	if m.isNPC() && m.Room.isDark() {
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
	jobs := make(map[string]string)
	jobs["war"] = "Warrior"
	jobs["mag"] = "Mage"
	jobs["cle"] = "Cleric"
	jobs["thi"] = "Thief"
	jobs["ran"] = "Ranger"
	jobs["bar"] = "Bard"

	for abbr, name := range jobs {
		j := &job{Name: name, Abbr: abbr}

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
	races := make(map[string]string)
	races["hum"] = "Human"
	races["elf"] = "Elf"
	races["dwf"] = "Dwarf"
	races["drw"] = "Dark Elf"
	races["gob"] = "Goblin"
	races["drg"] = "Dragon"

	for abbr, name := range races {
		r := &race{Name: name, Abbr: abbr}

		raceList.PushBack(r)
	}
}

func loadRooms() {
	areaFiles, _ := filepath.Glob("./data/area/*.json")

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

		void := &room{ID: 0, Exits: nil, Items: nil, Mobs: nil, Name: "The Void", Description: "A dark, gaping void lies here."}
		roomList.PushBack(void)

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
