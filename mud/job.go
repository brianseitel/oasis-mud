package mud

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type job struct {
	gorm.Model
	Name string
	Abbr string
}

func newJobDatabase() {
	fmt.Println("Creating jobs")
	jobs := make(map[string]string)
	jobs["war"] = "Warrior"
	jobs["mag"] = "Mage"
	jobs["cle"] = "Cleric"
	jobs["thi"] = "Thief"
	jobs["ran"] = "Ranger"
	jobs["bar"] = "Bard"

	for abbr, name := range jobs {
		j := &job{Name: name, Abbr: abbr}

		var found job
		db.Find(&found, job{Name: name})
		if !db.NewRecord(&found) {
			fmt.Println("\tSkipping job " + j.Name + "!")
		} else {
			fmt.Println("\tCreating job " + j.Name + "!")
			db.Create(&j)
		}
	}
}
