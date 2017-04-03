package mud

import (
	"fmt"

	"bytes"
	"strings"
)

/* common items */
const (
	startingDagger = 100008
	startingSword  = 100009
	startingMace   = 100010
	startingStaff  = 100011
	startingWhip   = 100012
)

const (
	goldSingle = 100001
	goldMulti  = 100002

	vnumTurd        = 100003
	vnumSeveredHead = 100004
	vnumTornHeart   = 100005
	vnumSlicedArm   = 100006
	vnumSlicedLeg   = 100007

	vnumCorpseNPC = 100014
	vnumCorpsePC  = 100013
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
	wearNone        = 0
	wearLight       = 1
	wearArmor       = 2
	wearFingerLeft  = 3
	wearFingerRight = 4
	wearNeck1       = 5
	wearNeck2       = 6
	wearBody        = 7
	wearHead        = 8
	wearLegs        = 9
	wearFeet        = 10
	wearHands       = 11
	wearArms        = 12
	wearShield      = 13
	wearAbout       = 14
	wearWaist       = 15
	wearWristLeft   = 16
	wearWristRight  = 17
	wearWield       = 18
	wearHold        = 19
	maxWear         = 20
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
	itemLight          = 1
	itemScroll         = 2
	itemWand           = 3
	itemStaff          = 4
	itemWeapon         = 5
	itemTreasure       = 6
	itemArmor          = 7
	itemPotion         = 8
	itemFurniture      = 9
	itemTrash          = 10
	itemContainer      = 11
	itemDrinkContainer = 12
	itemKey            = 13
	itemFood           = 14
	itemMoney          = 15
	itemBoat           = 16
	itemCorpseNPC      = 17
	itemCorpsePC       = 18
	itemFountain       = 19
	itemPill           = 20
)

/* wear flags for items */
const (
	itemTake       = 1
	itemWearNeck   = 2
	itemWearBody   = 4
	itemWearHead   = 8
	itemWearLegs   = 16
	itemWearFinger = 32
	itemWearFeet   = 64
	itemWearHands  = 128
	itemWearArms   = 256
	itemWearShield = 512
	itemWearAbout  = 1024
	itemWearWaist  = 2048
	itemWearWrist  = 4096
	itemWearWield  = 8192
	itemWearHold   = 16384
	itemWearLight  = 32768
)

type itemAttributeSet struct{}

type itemIndex struct {
	ID               int
	Name             string
	Description      string
	ShortDescription string
	ItemType         int   `json:"item_type"`
	ContainedIDs     []int `json:"contained_ids"`
	Affected         []*affect
	ExtraFlags       int `json:"extra_flags"`
	WearFlags        int `json:"wear_flags"`
	Weight           int
	Value            int
	Min              int
	Max              int
	SkillID          int /* items can have skills */
	Charges          int

	count int
}

type item struct {
	ID               int
	index            itemIndex
	container        []*item
	inObject         *item
	carriedBy        *mob
	Room             *room
	Affected         []*affect `json:"affected"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	ShortDescription string    `json:"short_description"`
	ItemType         int       `json:"item_type"`
	ExtraFlags       int       `json:"extra_flags"`
	WearFlags        int       `json:"wear_flags"`
	WearLocation     int       `json:"wear_location"`
	Weight           int       `json:"weight"`
	Cost             int       `json:"cost"`
	Level            int       `json:"level"`
	Timer            int       `json:"timer"`
	Value            int       `json:"value"`
	Min              int       `json:"min"`
	Max              int       `json:"max"`
	Skill            *skill    `json:"skill"` /* items can have skills or spells */
	Charges          int       `json:"charges"`

	ClosedFlags int
}

func (item *item) canWear(position int) bool {
	return hasBit(item.WearFlags, position)
}

func (item *item) hasExtraFlag(flag int) bool {
	return hasBit(item.ExtraFlags, flag)
}

func (item *item) isClosed() bool {
	return hasBit(item.ClosedFlags, containerClosed)
}

func (item *item) isCloseable() bool {
	return hasBit(item.ClosedFlags, containerClosable)
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
func createMoney(amount int) *item {
	if amount <= 0 {
		fmt.Printf("create_money: zero or negative money %d.%s", amount, newline)
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
