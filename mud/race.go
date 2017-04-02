package mud

type race struct {
	ID   int
	Name string
	Abbr string
}

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

func getRace(id int) *race {
	for e := raceList.Front(); e != nil; e = e.Next() {
		r := e.Value.(*race)
		if r.ID == id {
			return r
		}
	}
	return nil
}
