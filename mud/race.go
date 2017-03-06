package mud

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	// "github.com/brianseitel/oasis-mud/helpers"
)

type Race struct {
	gorm.Model

	Name string
}

func NewRaceDatabase() {
	fmt.Println("Creating races")
	races := []string{"Human", "Elf", "Dwarf", "Drow", "Goblin", "Dragon"}

	for _, v := range races {
		race := &Race{Name: v}

		var found Race
		db.Find(&found, Race{Name: v})
		if !db.NewRecord(&found) {
			fmt.Println("\tSkipping race " + race.Name + "!")
		} else {
			fmt.Println("\tCreating race " + race.Name + "!")
			db.Create(&race)
		}
	}
}
