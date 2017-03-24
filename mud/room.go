package mud

import (
	"container/list"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	// "github.com/brianseitel/oasis-mud/helpers"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //
)

var (
	roomList list.List
)

type area struct {
	gorm.Model

	Name  string  `json:"name"`
	Rooms []*room `json:"rooms",gorm:"-"`
}

type room struct {
	gorm.Model

	Area        area
	AreaID      int
	Name        string
	Description string
	Exits       []*exit `gorm:"many2many:room_exits;"`
	ItemIds     []int   `json:"items" gorm:"-"`
	Items       []*item `gorm:"many2many:room_items;"`
	Mobs        []*mob  `gorm:"many2many:room_mobs;"`
	MobIds      []int   `gorm:"-" json:"mobs"`
}

type exit struct {
	gorm.Model
	Dir    string `json:"direction"`
	Room   *room
	RoomID uint `json:"room_id",gorm:"-"`
}

func (r *room) decayItems() {
	for j, item := range r.Items {
		if item.timer == -1 {
			continue
		}

		if item.timer <= 0 {
			r.Items = append(r.Items[0:j], r.Items[j+1:]...)
			for _, m := range r.Mobs {
				if r.ID == m.Room.ID {
					m.notify(fmt.Sprintf("Rats scurry forth and drag away %s!\n", item.name))
				}
			}
			break
		}
		item.timer--
		fmt.Println("Decaying ", item.name, " (", item.timer, " ticks remaining")
	}
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
			db.Create(&ar)
		}

		void := &room{Model: gorm.Model{ID: 0}, Exits: nil, Items: nil, Mobs: nil, Name: "The Void", Description: "A dark, gaping void lies here."}
		roomList.PushBack(void)

		for _, ro := range a.Rooms {
			ro.AreaID = int(a.ID)
			for _, i := range ro.ItemIds {
				var item item
				db.First(&item, i)

				ro.Items = append(ro.Items, &item)
			}

			for _, i := range ro.MobIds {
				var mob mob
				db.First(&mob, i)
			}

			var r room
			db.First(&r, ro.ID)

			if db.NewRecord(&r) {
				db.Create(&ro)
			}

			roomList.PushBack(ro)
		}

		exitsList := list.New()
		for e := roomList.Front(); e != nil; e = e.Next() {
			room := e.Value.(*room)
			for j, x := range room.Exits {
				room.Exits[j] = &exit{Dir: x.Dir, Room: getRoom(x.RoomID), RoomID: x.RoomID}
			}

			for _, x := range room.MobIds {
				fmt.Println("Adding mob", x, " to room ", room.ID)

				mob := getMob(uint(x))
				room.Mobs = append(room.Mobs, mob)
			}

			exitsList.PushBack(room)
		}

		roomList = *exitsList

		for e := mobList.Front(); e != nil; e = e.Next() {
			mob := e.Value.(*mob)
			room := getRoom(uint(mob.RoomID))
			mob.Room = room
		}

	}
}

func getRoom(id uint) *room {
	for e := roomList.Front(); e != nil; e = e.Next() {
		r := e.Value.(*room)
		if r.ID == id {
			return r
		}
	}
	return nil
}

func getMob(id uint) *mob {
	for e := mobList.Front(); e != nil; e = e.Next() {
		m := e.Value.(*mob)
		if m.ID == id {
			return m
		}
	}
	return nil
}

func (r *room) removeMob(m *mob) {
	for j, mob := range r.Mobs {
		if mob == m {
			r.Mobs = append(r.Mobs[0:j], r.Mobs[j+1:]...)
			break
		}
	}
}
