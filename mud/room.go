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

type Area struct {
	gorm.Model

	Name  string `json:"name"`
	Rooms []Room `json:"rooms",gorm:"-"`
}

type Room struct {
	gorm.Model

	Area        Area
	AreaId      int
	Name        string
	Description string
	Exits       []Exit `gorm:"many2many:room_exits;"`
	ItemIds     []int  `json:"items" gorm:"-"`
	Items       []Item `gorm:"many2many:room_items;"`
	Mobs        []Mob
	MobIds      []int `gorm:"-" json:"mobs"`
}

type Exit struct {
	gorm.Model
	Dir    string `json:"direction"`
	RoomId int    `json:"room_id",gorm:"-"`
}

type RoomDatabase struct {
	Rooms []Room
}

// Finds a given room in the database. If not found,
// returns a blank room
func FindRoom(r int) Room {
	var room Room
	db.First(&room, r)
	return room
}

// Creates a new room database, seeding it with data from the areas
// directory.
func NewRoomDatabase() {
	areaFiles, _ := filepath.Glob("./data/area/*.json")

	for _, areaFile := range areaFiles {
		file, err := ioutil.ReadFile(areaFile)
		if err != nil {
			panic(err)
		}

		var area Area
		json.Unmarshal(file, &area)
		if err != nil {
			panic(err)
		}

		a := &Area{Name: area.Name}

		db.First(&a)
		if db.NewRecord(&a) {
			fmt.Println("\tCreating area " + a.Name + "!")
			db.Create(&a)
		} else {
			fmt.Println("\tSkipping area " + a.Name + "!")
		}

		for _, room := range area.Rooms {
			room.AreaId = int(a.ID)
			for _, i := range room.ItemIds {
				var item Item
				db.First(&item, i)

				room.Items = append(room.Items, item)
			}

			for _, i := range room.MobIds {
				var mob Mob
				db.First(&mob, i)
				room.Mobs = append(room.Mobs, mob)
			}

			var r Room
			db.First(&r, room.ID)

			if db.NewRecord(&r) {
				fmt.Println("\tCreating room " + room.Name + "!")
				db.Create(&room)
			} else {
				fmt.Println("\tSkipping room " + room.Name + "!")
			}
		}
	}
}
