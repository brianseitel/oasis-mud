package mud

import (
	"container/list"
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	// "github.com/brianseitel/oasis-mud/helpers"
	"github.com/jinzhu/gorm"
)

var (
	itemList list.List
)

const (
	weapon = iota
	armor
	healing
)

const (
	decays = iota
	permanent
)

type position string

const (
	light     position = "light"
	finger1   position = "finger"
	finger2   position = "finger"
	neck1     position = "neck"
	neck2     position = "neck"
	torso     position = "torso"
	head      position = "head"
	legs      position = "legs"
	feet      position = "feet"
	hands     position = "hands"
	arms      position = "arms"
	shield    position = "shield"
	body      position = "body"
	waist     position = "waist"
	wrist1    position = "wrist"
	wrist2    position = "wrist"
	wield     position = "wield"
	held      position = "held"
	floating  position = "floating"
	secondary position = "secondary"
)

type itemAttributeSet struct{}

type item struct {
	gorm.Model

	itemType    string
	Name        string `json:"name"`
	Description string
	Min         int
	Max         int
	Speed       int
	Price       int
	Position    string
	// Attributes  itemAttributeSet
	Identifiers string
	Decays      uint
	TTL         int
}

func newItemDatabase() {
	itemFiles, _ := filepath.Glob("./data/items/*.json")

	for _, itemFile := range itemFiles {
		file, err := ioutil.ReadFile(itemFile)
		if err != nil {
			panic(err)
		}

		var list []item
		json.Unmarshal(file, &list)

		for _, it := range list {
			var items []item
			db.Find(&items, it)
			if len(items) == 0 {
				db.Create(&it)
			}

			itemList.PushBack(it)
		}

	}
}
