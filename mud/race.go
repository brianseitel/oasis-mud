package mud

import (
	"fmt"

	"github.com/jinzhu/gorm"
	// "github.com/brianseitel/oasis-mud/helpers"
)

type race struct {
	gorm.Model

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

func newRaceDatabase() {
	fmt.Println("Creating races")
	races := make(map[string]string)
	races["hum"] = "Human"
	races["elf"] = "Elf"
	races["dwf"] = "Dwarf"
	races["drw"] = "Dark Elf"
	races["gob"] = "Goblin"
	races["drg"] = "Dragon"

	for abbr, name := range races {
		r := &race{Name: name, Abbr: abbr}

		var found race
		db.Find(&found, race{Name: name})
		if !db.NewRecord(&found) {
			fmt.Println("\tSkipping race " + r.Name + "!")
		} else {
			fmt.Println("\tCreating race " + r.Name + "!")
			db.Create(&r)
		}
	}
}
