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
	raceList      list.List
	roomList      list.List
	skillList     list.List
)

func bootDB() {
	loadSkills()
	loadJobs()
	loadRaces()
	loadItems()
	loadMobs()
	loadRooms()

	areaUpdate()
}

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

		var list []*mob
		err = json.Unmarshal(file, &list)
		if err != nil {
			panic(err)
		}

		for _, m := range list {

			var skills []*mobSkill
			for _, s := range m.Skills {
				skill := getSkill(s.SkillID)
				skills = append(skills, &mobSkill{Skill: skill, SkillID: s.SkillID, Level: s.Level})
			}
			m.Skills = skills

			mobList.PushBack(m)
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
				item := newItemFromIndex(index)
				ro.Items = append(ro.Items, item)
			}

			for _, i := range ro.MobIds {
				mob := getMob(uint(i))
				if mob == nil {
				}
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
			item := newItemFromIndex(index)
			ro.Items = append(ro.Items, item)
		}

		for _, i := range ro.MobIds {
			mob := getMob(uint(i))
			if mob == nil {
			}
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
