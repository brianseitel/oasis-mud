package mud

import (
	"container/list"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	// "github.com/brianseitel/oasis-mud/helpers"
)

var (
	roomList list.List
)

type area struct {
	ID    uint
	Name  string  `json:"name"`
	Rooms []*room `json:"rooms",gorm:"-"`
}

type room struct {
	ID   uint
	Name string

	Area        area
	AreaID      int
	Description string
	Exits       []*exit `gorm:"many2many:room_exits;"`
	ItemIds     []int   `json:"items" gorm:"-"`
	Items       []*item `gorm:"many2many:room_items;"`
	Mobs        []*mob  `gorm:"many2many:room_mobs;"`
	MobIds      []int   `gorm:"-" json:"mobs"`
}

type exit struct {
	ID     uint
	Dir    string `json:"direction"`
	Room   *room
	RoomID uint `json:"room_id",gorm:"-"`
}

func (r *room) decayItems() {
	for j, item := range r.Items {
		if item.Timer == -1 {
			continue
		}

		if item.Timer <= 0 {
			r.Items = append(r.Items[0:j], r.Items[j+1:]...)
			for _, m := range r.Mobs {
				if r.ID == m.Room.ID {
					m.notify(fmt.Sprintf("Rats scurry forth and drag away %s!\n", item.Name))
				}
			}
			break
		}
		item.Timer--
		fmt.Println("Decaying ", item.Name, " (", item.Timer, " ticks remaining")
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

		void := &room{ID: 0, Exits: nil, Items: nil, Mobs: nil, Name: "The Void", Description: "A dark, gaping void lies here."}
		roomList.PushBack(void)

		for _, ro := range a.Rooms {
			ro.AreaID = int(a.ID)
			for _, i := range ro.ItemIds {
				index := getItem(uint(i))
				item := newItemFromIndex(index)
				ro.Items = append(ro.Items, item)
			}

			for _, i := range ro.MobIds {
				mob := getMob(uint(i))
				if mob == nil {
				}
				ro.Mobs = append(ro.Mobs, mob)
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

func getItem(id uint) *itemIndex {
	for e := itemIndexList.Front(); e != nil; e = e.Next() {
		i := e.Value.(itemIndex)
		if i.ID == id {
			return &i
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

func (r *room) notify(message string, except *mob) {
	for _, mob := range r.Mobs {
		if mob != except {
			mob.notify(message)
		}
	}
}

func (r *room) removeMob(m *mob) {
	for j, mob := range r.Mobs {
		if mob == m {
			r.Mobs = append(r.Mobs[0:j], r.Mobs[j+1:]...)
			break
		}
	}
}

func (r *room) removeObject(i *item) {
	for j, item := range r.Items {
		if item == i {
			r.Items = append(r.Items[0:j], r.Items[j+1:]...)
			return
		}
	}
}
