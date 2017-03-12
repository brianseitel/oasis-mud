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
	// Attributes  itemAttributeSet
	Identifiers string
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
