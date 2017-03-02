package main

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"text/template"
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
func (rdb *RoomDatabase) FindRoom(r int) Room {
	for _, v := range rdb.Rooms {
		if v.Id == r {
			return v
		}
	}
	return Room{}
}

// Displays the current room
func (r Room) Display(c Connection) {
	tmpl := `{{ .Name }}
{{ .Description }}
{{ if .Items }}{{ range .Items }}There is a {{ .Name }} on the floor.
{{ end }}{{ end }}
{{ if .Exits }}Exits: {{ range .Exits }}{{ .Dir }} {{ end }}{{ end }}`

	t := template.Must(template.New("room").Parse(tmpl))

	_ = t.Execute(c.conn, r)
}

// Removes an item from a room
func (rdb *RoomDatabase) RemoveItem(room Room, item Item) {
	for key, r := range rdb.Rooms {
		if r.Id == room.Id {
			for k, i := range room.Items {
				if i.Id == item.Id {
					rdb.Rooms[key].Items = append(room.Items[:k], room.Items[k+1:]...)
					return
				}
			}
		}
	}
}

// Adds an item to a room
func (rdb *RoomDatabase) AddItem(room Room, item Item) {
	for key, r := range rdb.Rooms {
		if r.Id == room.Id {
			rdb.Rooms[key].Items = append(rdb.Rooms[key].Items, item)
			return
		}
	}
}

// Creates a new room database, seeding it with data from the areas
// directory.
func NewRoomDatabase() *RoomDatabase {
	dbItems = NewItemDatabase()
	areaFiles, _ := filepath.Glob("./areas/*.json")

	rooms := &RoomDatabase{}
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
			rooms.Rooms = append(rooms.Rooms, room)
		}
	}
	return rooms
}
