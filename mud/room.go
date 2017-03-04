package mud

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

type Area struct {
	Name  string `json:"name"`
	Rooms []Room `json:"rooms"`
}

type Room struct {
	Id          int
	Name        string
	Description string
	Exits       []Direction `json:"exits"`
	ItemIds     []int       `json:"items"`
	Items       []Item
}

type Direction struct {
	Dir    string `json:"direction"`
	RoomId int    `json:"room_id"`
}

type RoomDatabase struct {
	Rooms []Room
}

// Finds a given room in the database. If not found,
// returns a blank room
func FindRoom(r int) Room {
	for _, v := range Registry.rooms {
		if v.Id == r {
			return v
		}
	}
	return Room{}
}

// Creates a new room database, seeding it with data from the areas
// directory.
func NewRoomDatabase() map[int]Room {
	areaFiles, _ := filepath.Glob("./data/area/*.json")

	rooms := make(map[int]Room)
	for _, areaFile := range areaFiles {
		file, err := ioutil.ReadFile(areaFile)
		if err != nil {
			panic(err)
		}

		var area Area
		json.Unmarshal(file, &area)

		for _, room := range area.Rooms {
			for _, v := range room.ItemIds {
				room.Items = append(room.Items, FindItem(v))
			}
			rooms[room.Id] = room
		}
	}

	return rooms
}
