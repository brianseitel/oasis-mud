package mud

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Job struct {
	gorm.Model
	Name string
}

func NewJobDatabase() {
	fmt.Println("Creating jobs")
	jobs := []string{"Warrior", "Mage", "Cleric", "Thief", "Ranger", "Bard"}

	for _, v := range jobs {
		job := &Job{Name: v}

		var found Job
		db.Find(&found, Job{Name: v})
		if !db.NewRecord(&found) {
			fmt.Println("\tSkipping job " + job.Name + "!")
		} else {
			fmt.Println("\tCreating job " + job.Name + "!")
			db.Create(&job)
		}
	}
}
