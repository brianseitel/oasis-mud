package mud

import (
	"fmt"

	"github.com/brianseitel/oasis-mud/helpers"
	// "github.com/brianseitel/oasis-mud/helpers"
	"bytes"
	"strings"
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
	Min              int
	Max              int

	count int
}

type item struct {
	ID               uint
	index            itemIndex
	container        []*item
	inObject         *item
	carriedBy        *mob
	Room             *room
	Affected         []*affect
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
	Min              int
	Max              int
}

func (item *item) canWear(position uint) bool {
	return helpers.HasBit(item.WearFlags, position)
}

func (item *item) hasExtraFlag(flag uint) bool {
	return helpers.HasBit(item.ExtraFlags, flag)
}

func (item *item) isClosed() bool {
	return false
}

func (item *item) removeObject(target *item) {
	for j, it := range item.container {
		if it == target {
			item.container = append(item.container[0:j], item.container[j+1:]...)
			return
		}
	}
}

func affectBitName(vector int) string {
	var buf bytes.Buffer

	if vector&affectBlind == affectBlind {
		buf.Write([]byte(" blind"))
	}
	if vector&affectInvisible == affectInvisible {
		buf.Write([]byte(" invisible"))
	}
	if vector&affectDetectEvil == affectDetectEvil {
		buf.Write([]byte(" detect_evil"))
	}
	if vector&affectDetectInvisible == affectDetectInvisible {
		buf.Write([]byte(" detect_invis"))
	}
	if vector&affectDetectMagic == affectDetectMagic {
		buf.Write([]byte(" detect_magic"))
	}
	if vector&affectDetectHidden == affectDetectHidden {
		buf.Write([]byte(" detect_hidden"))
	}
	if vector&affectHold == affectHold {
		buf.Write([]byte(" hold"))
	}
	if vector&affectSanctuary == affectSanctuary {
		buf.Write([]byte(" sanctuary"))
	}
	if vector&affectFaerieFire == affectFaerieFire {
		buf.Write([]byte(" faerie_fire"))
	}
	if vector&affectInfrared == affectInfrared {
		buf.Write([]byte(" infrared"))
	}
	if vector&affectCurse == affectCurse {
		buf.Write([]byte(" curse"))
	}
	if vector&affectFlaming == affectFlaming {
		buf.Write([]byte(" flaming"))
	}
	if vector&affectPoison == affectPoison {
		buf.Write([]byte(" poison"))
	}
	if vector&affectProtect == affectProtect {
		buf.Write([]byte(" protect"))
	}
	if vector&affectParalysis == affectParalysis {
		buf.Write([]byte(" paralysis"))
	}
	if vector&affectSleep == affectSleep {
		buf.Write([]byte(" sleep"))
	}
	if vector&affectSneak == affectSneak {
		buf.Write([]byte(" sneak"))
	}
	if vector&affectHide == affectHide {
		buf.Write([]byte(" hide"))
	}
	if vector&affectCharm == affectCharm {
		buf.Write([]byte(" charm"))
	}
	if vector&affectFlying == affectFlying {
		buf.Write([]byte(" flying"))
	}
	if vector&affectPassDoor == affectPassDoor {
		buf.Write([]byte(" pass_door"))
	}

	output := buf.String()

	if len(output) == 0 {
		return "none"
	}

	return strings.Trim(output, " ")
}

func affectLocationName(location int) string {
	switch location {
	case applyNone:
		return "none"
	case applyStrength:
		return "strength"
	case applyDexterity:
		return "dexterity"
	case applyIntelligence:
		return "intelligence"
	case applyWisdom:
		return "wisdom"
	case applyConstitution:
		return "constitution"
	case applySex:
		return "sex"
	case applyClass:
		return "class"
	case applyLevel:
		return "level"
	case applyMana:
		return "mana"
	case applyHitpoints:
		return "hitpoints"
	case applyMovement:
		return "movement"
	case applyGold:
		return "gold"
	case applyExp:
		return "experience"
	case applyArmorClass:
		return "armor class"
	case applyHitroll:
		return "hit roll"
	case applyDamroll:
		return "dam roll"
	case applySavingParalysis:
		return "save vs paralysis"
	case applySavingRod:
		return "save vs rod"
	case applySavingPetrify:
		return "save vs petrify"
	case applySavingBreath:
		return "save vs breath"
	case applySavingSpell:
		return "save vs spell"
	}

	return "(unknown)"
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
func createMoney(amount uint) *item {
	if amount <= 0 {
		fmt.Printf("create_money: zero or negative money %d.%s", amount, helpers.Newline)
		amount = 1
	}

	var obj *item
	if amount == 1 {
		obj = createItem(getItem(goldSingle))
	} else {
		obj = createItem(getItem(goldMulti))
		obj.Value = int(amount)
	}

	return obj
}

func extraBitName(flags int) string {

	var buf bytes.Buffer
	if flags&itemGlow == itemGlow {
		buf.Write([]byte(" glow"))
	}
	if flags&itemHum == itemHum {
		buf.Write([]byte(" hum"))
	}
	if flags&itemDark == itemDark {
		buf.Write([]byte(" dark"))
	}
	if flags&itemLock == itemLock {
		buf.Write([]byte(" lock"))
	}
	if flags&itemEvil == itemEvil {
		buf.Write([]byte(" evil"))
	}
	if flags&itemInvis == itemInvis {
		buf.Write([]byte(" invis"))
	}
	if flags&itemMagic == itemMagic {
		buf.Write([]byte(" magic"))
	}
	if flags&itemNoDrop == itemNoDrop {
		buf.Write([]byte(" nodrop"))
	}
	if flags&itemBless == itemBless {
		buf.Write([]byte(" bless"))
	}
	if flags&itemAntiEvil == itemAntiEvil {
		buf.Write([]byte(" anti-evil"))
	}
	if flags&itemAntiGood == itemAntiGood {
		buf.Write([]byte(" anti-good"))
	}
	if flags&itemAntiNeutral == itemAntiNeutral {
		buf.Write([]byte(" anti-neutral"))
	}
	if flags&itemNoRemove == itemNoRemove {
		buf.Write([]byte(" noremove"))
	}
	if flags&itemInventory == itemInventory {
		buf.Write([]byte(" inventory"))
	}

	output := buf.String()

	if len(output) == 0 {
		return "none"
	}

	return strings.Trim(output, " ")

}

func itemTypeName(item *item) string {
	switch item.ItemType {
	case itemLight:
		return "light"
	case itemScroll:
		return "scroll"
	case itemWand:
		return "wand"
	case itemStaff:
		return "staff"
	case itemWeapon:
		return "weapon"
	case itemTreasure:
		return "treasure"
	case itemArmor:
		return "armor"
	case itemPotion:
		return "potion"
	case itemFurniture:
		return "furniture"
	case itemTrash:
		return "trash"
	case itemContainer:
		return "container"
	case itemDrinkContainer:
		return "drink container"
	case itemKey:
		return "key"
	case itemFood:
		return "food"
	case itemMoney:
		return "money"
	case itemBoat:
		return "boat"
	case itemCorpseNPC:
		return "npc corpse"
	case itemCorpsePC:
		return "pc corpse"
	case itemFountain:
		return "fountain"
	case itemPill:
		return "pill"
	}

	return "(unknown)"
}
