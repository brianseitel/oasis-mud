package mud

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

const (
	weapon = iota
	armor
	healing
)

type ItemAttributeSet struct{}

type Item struct {
	Id          int              `json:"id"`
	ItemType    string           `json:"type"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Min         int              `json:"min"`
	Max         int              `json:"max"`
	Speed       int              `json:"speed"`
	Price       int              `json:"price"`
	Attributes  ItemAttributeSet `json:"attributes"`
	Identifiers []string         `json:"identifiers"`
}

type ItemList struct {
	Items []Item `json:"items"`
}
type ItemDatabase []Item

// Finds an item in the item Database
// If not found, returns an empty item
func FindItem(i int) Item {
	for _, v := range Registry.items {
		if v.Id == i {
			return v
		}
	}
	return Item{}
}

// Seeds the item database with data from our items directory
func NewItemDatabase() []Item {
	itemFiles, _ := filepath.Glob("./data/items/*.json")

	var items []Item

	for _, itemFile := range itemFiles {
		file, err := ioutil.ReadFile(itemFile)
		if err != nil {
			panic(err)
		}

		var list ItemList
		json.Unmarshal(file, &list)

		for _, item := range list.Items {
			items = append(items, item)
		}
	}

	return items
}
