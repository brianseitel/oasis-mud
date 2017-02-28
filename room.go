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

func (rdb *RoomDatabase) FindRoom(r int) Room {
	for _, v := range rdb.Rooms {
		if v.Id == r {
			return v
		}
	}

	return Room{}
}

func (r *Room) Display(c Connection) {
	for _, v := range r.ItemIds {
		r.Items = append(r.Items, FindItem(v))
	}

	tmpl := `{{ .Name }}
{{ .Description }}
{{ if .Items }}
{{ range .Items }}There is a {{ .Name }} on the floor.
{{end}}
{{ end }}
{{ if .Exits }}
Exits: {{ range .Exits }}{{ .Dir }} {{ end }}
{{ end }}`

	t := template.Must(template.New("room").Parse(tmpl))

	_ = t.Execute(c.conn, r)
}

func NewRoomDatabase() *RoomDatabase {
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
			rooms.Rooms = append(rooms.Rooms, room)
		}
	}
	return rooms
}
