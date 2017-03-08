package mud

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	// "github.com/brianseitel/oasis-mud/helpers"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //
)

type area struct {
	gorm.Model

	Name  string `json:"name"`
	Rooms []room `json:"rooms",gorm:"-"`
}

type room struct {
	gorm.Model

	Area        area
	AreaID      int
	Name        string
	Description string
	Exits       []exit `gorm:"many2many:room_exits;"`
	ItemIds     []int  `json:"items" gorm:"-"`
	Items       []item `gorm:"many2many:room_items;"`
	Mobs        []mob
	MobIds      []int `gorm:"-" json:"mobs"`
	Players     []player
}

type exit struct {
	gorm.Model
	Dir    string `json:"direction"`
	RoomID int    `json:"room_id",gorm:"-"`
}

type roomDatabase struct {
	Rooms []room
}

func findRoom(r int) room {
	var (
		room    room
		players []player
	)

	db.First(&room, r).Related(&players)
	return room
}

func newRoomDatabase() {
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

		ar := &area{Name: a.Name}

		db.First(&ar)
		if db.NewRecord(&ar) {
			fmt.Println("\tCreating area " + ar.Name + "!")
			db.Create(&ar)
		} else {
			fmt.Println("\tSkipping area " + ar.Name + "!")
		}

		for _, ro := range a.Rooms {
			ro.AreaID = int(a.ID)
			for _, i := range ro.ItemIds {
				var item item
				db.First(&item, i)

				ro.Items = append(ro.Items, item)
			}

			for _, i := range ro.MobIds {
				var mob mob
				db.First(&mob, i)
				ro.Mobs = append(ro.Mobs, mob)
			}

			var r room
			db.First(&r, ro.ID)

			if db.NewRecord(&r) {
				fmt.Println("\tCreating room " + ro.Name + "!")
				db.Create(&ro)
			} else {
				fmt.Println("\tSkipping room " + ro.Name + "!")
			}
		}
	}
}
