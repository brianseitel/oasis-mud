package mud

type raceStats struct {
	Hitpoints    int
	Mana         int
	Movement     int
	Strength     int
	Wisdom       int
	Dexterity    int
	Charisma     int
	Constitution int
	Intelligence int
}

type race struct {
	ID        int
	Name      string
	Adjective string
	Abbr      string
	Stats     raceStats
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
