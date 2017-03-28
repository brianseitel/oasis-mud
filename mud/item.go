package mud

import (
	"fmt"

	"github.com/brianseitel/oasis-mud/helpers"
	// "github.com/brianseitel/oasis-mud/helpers"
)

const (
	goldSingle = 100001
	goldMulti  = 100002

	vnumTurd        = 100003
	vnumSeveredHead = 100004
	vnumTornHeart   = 100005
	vnumSlicedArm   = 100006
	vnumSlicedLeg   = 100007
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
	containerClosable = 1 << iota
	containerPickproof
	containerClosed
	containerLocked
)

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
	wearNone  = 99999
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
	ItemType         uint   `json:"item_type"`
	ContainedIDs     []uint `json:"contained_ids"`
	Affected         []*affect
	ExtraFlags       uint `json:"extra_flags"`
	WearFlags        uint `json:"wear_flags"`
	Weight           uint
	Value            int
}

type item struct {
	ID               uint
	index            itemIndex
	container        []*item
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

func newItemFromIndex(index *itemIndex) *item {
	var contained []*item
	item := &item{Name: index.Name, Description: index.Description, ShortDescription: index.ShortDescription, ItemType: index.ItemType, ExtraFlags: index.ExtraFlags, WearFlags: index.WearFlags, Weight: index.Weight, Value: index.Value, Timer: -1}
	for _, id := range index.ContainedIDs {
		i := getItem(id)
		contained = append(contained, newItemFromIndex(i))
	}
	item.container = contained
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

func (item *item) removeObject(target *item) {
	for j, it := range item.container {
		if it == target {
			item.container = append(item.container[0:j], item.container[j+1:]...)
			return
		}
	}
}

func (item *item) extract() {

}

func createMoney(amount uint) *item {
	if amount <= 0 {
		fmt.Printf("create_money: zero or negative money %d.%s", amount, helpers.Newline)
		amount = 1
	}

	var obj *item
	if amount == 1 {
		obj = newItemFromIndex(getItem(goldSingle))
	} else {
		obj = newItemFromIndex(getItem(goldMulti))
		obj.Value = int(amount)
	}

	return obj
}
