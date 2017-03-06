package mud

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	// "github.com/brianseitel/oasis-mud/helpers"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

const (
	weapon = iota
	armor
	healing
)

type ItemAttributeSet struct{}

type Item struct {
	gorm.Model

	ItemType    string
	Name        string `json:"name"`
	Description string
	Min         int
	Max         int
	Speed       int
	Price       int
	// Attributes  ItemAttributeSet
	Identifiers string
}

type ItemDatabase []Item

// Finds an item in the item Database
// If not found, returns an empty item
func FindItem(i int) Item {
	var item Item
	db.First(&item, i)
	return item
}

// Seeds the item database with data from our items directory
func NewItemDatabase() {
	fmt.Println("Creating items!")
	itemFiles, _ := filepath.Glob("./data/items/*.json")

	for _, itemFile := range itemFiles {
		file, err := ioutil.ReadFile(itemFile)
		if err != nil {
			panic(err)
		}

		var list []Item
		json.Unmarshal(file, &list)

		for _, item := range list {
			var items []Item
			db.Find(&items, item)
			if len(items) == 0 {
				fmt.Println("\tCreating item " + item.Name + "!")
				db.Create(&item)
			} else {
				fmt.Println("\tSkipping item " + item.Name + "!")
			}
		}
	}
}
