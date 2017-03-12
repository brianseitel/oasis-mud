package mud

import (
	"container/list"
	"encoding/json"
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
	Mobs        []*mob
	MobIds      []int `gorm:"-" json:"mobs"`
}

type exit struct {
	gorm.Model
	Dir    string `json:"direction"`
	Room   *room
	RoomID uint `json:"room_id",gorm:"-"`
}

func (x *exit) getRoom() {
	for e := roomList.Front(); e != nil; e = e.Next() {
		room := e.Value.(*room)
		if room.ID == x.RoomID {
			x.Room = room
			return
		}
	}
	return
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
				ro.Mobs = append(ro.Mobs, &mob)
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
			exitsList.PushBack(room)
		}

		roomList = *exitsList

		for e := mobList.Front(); e != nil; e = e.Next() {
			mob := e.Value.(*mob)
			mob.Room = getRoom(uint(mob.RoomID))
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
