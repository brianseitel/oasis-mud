package mud

import (
	"container/list"
)

type race struct {
	ID   uint
	Name string
	Abbr string
}

var raceList list.List

func (r race) defaultStats(s string) int {
	defaults := make(map[string]int)
	defaults["hitpoints"] = 100
	defaults["mana"] = 0
	defaults["movement"] = 100

	defaults["strength"] = 12
	defaults["wisdom"] = 12
	defaults["intelligence"] = 12
	defaults["dexterity"] = 12
	defaults["charisma"] = 12
	defaults["constitution"] = 12

	return defaults[s]
}

func newRaceDatabase() {
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

func getRace(id uint) *race {
	for e := raceList.Front(); e != nil; e = e.Next() {
		r := e.Value.(*race)
		if r.ID == id {
			return r
		}
	}
	return nil
}
