package main

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

type Item struct {
	Id          int          `json:"id"`
	ItemType    string       `json:"type"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Min         int          `json:"min"`
	Max         int          `json:"max"`
	Speed       int          `json:"speed"`
	Price       int          `json:"price"`
	Attributes  AttributeSet `json:"attributes"`
}

type ItemList struct {
	Items []Item `json:"items"`
}
type ItemDatabase []Item

func (i *Item) GetAttr(attribute int) int {
	return i.Attributes[attribute]
}

func FindItem(i int) Item {
	for _, v := range *dbItems {
		if v.Id == i {
			return v
		}
	}
	return Item{}
}

func NewItemDatabase() *ItemDatabase {
	itemFiles, _ := filepath.Glob("./items/*.json")

	items := &ItemDatabase{}

	for _, itemFile := range itemFiles {
		file, err := ioutil.ReadFile(itemFile)
		if err != nil {
			panic(err)
		}

		var list ItemList
		json.Unmarshal(file, &list)

		for _, item := range list.Items {
			*items = append(*items, item)
		}
	}

	return items
}
