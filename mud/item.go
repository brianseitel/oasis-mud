package mud

import (
	"container/list"
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/brianseitel/oasis-mud/helpers"
	// "github.com/brianseitel/oasis-mud/helpers"
)

var (
	itemList      list.List
	itemIndexList list.List
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

/* wear locations */
const (
	wearNone  = -1
	wearLight = iota
	wearArmor
	wearFingerLeft
	wearFingerRight
	wearNeck1
	wearNeck2
	wearBody
	wearHead
	wearLegs
	wearFeet
	wearHands
	wearArms
	wearShield
	wearAbout
	wearWaist
	wearWristLeft
	wearWristRight
	wearWield
	wearHold
	maxWear
)

/* Item flags */
const (
	itemGlow = 1 << iota
	itemHum
	itemDark
	itemLock
	itemEvil
	itemInvis
	itemMagic
	itemNoDrop
	itemBless
	itemAntiGood
	itemAntiEvil
	itemAntiNeutral
	itemNoRemove
	itemInventory
)

/* Item types */
const (
	itemLight = iota
	itemScroll
	itemWand
	itemStaff
	itemWeapon
	itemTreasure
	itemArmor
	itemPotion
	itemFurniture
	itemTrash
	itemContainer
	itemDrinkContainer
	itemKey
	itemFood
	itemMoney
	itemBoat
	itemCorpseNPC
	itemCorpsePC
	itemFountain
	itemPill
)

/* wear flags for items */
const (
	itemTake = 1 << iota
	itemWearNeck
	itemWearBody
	itemWearHead
	itemWearLegs
	itemWearFinger
	itemWearFeet
	itemWearHands
	itemWearArms
	itemWearShield
	itemWearAbout
	itemWearWaist
	itemWearWrist
	itemWearWield
	itemWearHold
	itemWearLight
)

type itemAttributeSet struct{}

type itemIndex struct {
	ID               uint
	Name             string
	Description      string
	ShortDescription string
	ItemType         uint `json:"item_type"`
	Affected         []*affect
	ExtraFlags       uint `json:"extra_flags"`
	WearFlags        uint `json:"wear_flags"`
	Weight           uint
	Value            int
}

type item struct {
	ID               uint
	index            itemIndex
	Name             string
	Description      string
	ShortDescription string
	ItemType         uint
	ExtraFlags       uint
	WearFlags        uint
	WearLocation     uint
	Weight           uint
	Cost             int
	Level            int
	Timer            int
	Value            int
}

func newItemDatabase() {
	itemFiles, _ := filepath.Glob("./data/items/*.json")

	for _, itemFile := range itemFiles {
		file, err := ioutil.ReadFile(itemFile)
		if err != nil {
			panic(err)
		}

		var list []itemIndex
		json.Unmarshal(file, &list)

		for _, it := range list {

			itemIndexList.PushBack(it)
		}

	}
}

func newItemFromIndex(index *itemIndex) *item {
	item := &item{Name: index.Name, Description: index.Description, ShortDescription: index.ShortDescription, ItemType: index.ItemType, ExtraFlags: index.ExtraFlags, WearFlags: index.WearFlags, Weight: index.Weight, Value: index.Value, Timer: -1}

	itemList.PushBack(item)
	return item
}

func applyAC(item *item, wear int) int {
	if item.ItemType != itemArmor {
		return 0
	}

	switch wear {
	case wearBody:
		return 3 * item.Value
	case wearHead:
		return 2 * item.Value
	case wearLegs:
		return 2 * item.Value
	case wearFeet:
		return item.Value
	case wearHands:
		return item.Value
	case wearArms:
		return item.Value
	case wearShield:
		return item.Value
	case wearFingerLeft:
		return item.Value
	case wearFingerRight:
		return item.Value
	case wearNeck1:
		return item.Value
	case wearNeck2:
		return item.Value
	case wearArmor:
		return 2 * item.Value
	case wearWaist:
		return item.Value
	case wearWristLeft:
		return item.Value
	case wearWristRight:
		return item.Value
	case wearHold:
		return item.Value
	}

	return 0
}

func (item *item) canWear(position uint) bool {
	return helpers.HasBit(item.WearFlags, position)
}

func (item *item) hasExtraFlag(flag uint) bool {
	return helpers.HasBit(item.ExtraFlags, flag)
}
