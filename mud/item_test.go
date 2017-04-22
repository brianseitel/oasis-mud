package mud

import "testing"

func TestCanWear(t *testing.T) {
	item := &item{WearFlags: wearNone}

	if item.canWear(wearWield) {
		t.Error("Should not be able to wear this item.")
	}

	item.WearFlags = itemWearBody

	if !item.canWear(itemWearBody) {
		t.Error("Should be able to wear this armor")
	}

	if item.canWear(itemWearWield) {
		t.Error("Should not be able to wield armor")
	}
}

func TestHasExtraFlag(t *testing.T) {
	i := &item{ExtraFlags: 1}

	if !i.hasExtraFlag(itemGlow) {
		t.Error("Failed to detect glow flag")
	}

	i = &item{ExtraFlags: 3}
	if !i.hasExtraFlag(itemGlow) && !i.hasExtraFlag(itemHum) {
		t.Error("Failed to properly detect glow and hum flags")
	}
}

func TestItemIsClosed(t *testing.T) {
	item := &item{ClosedFlags: containerClosed}

	if !item.isClosed() {
		t.Error("Failed to detect closed status")
	}

	item.ClosedFlags = 0
	if item.isClosed() {
		t.Error("Failed to detect closed status of item")
	}
}

func TestItemIsClosable(t *testing.T) {
	item := &item{ClosedFlags: containerClosable}

	if !item.isCloseable() {
		t.Error("Failed to detect closable status")
	}

	item.ClosedFlags = 0
	if item.isCloseable() {
		t.Error("Failed to detect closable status of item")
	}
}

func TestRemoveItemFromContainer(t *testing.T) {
	obj := &item{Name: "a test item"}
	container := &item{ItemType: itemContainer}

	container.container = append(container.container, obj)

	found := false
	for _, i := range container.container {
		if i == obj {
			found = true
		}
	}

	if !found {
		t.Error("No item in container")
	}

	container.removeObject(obj)

	found = false
	for _, i := range container.container {
		if i == obj {
			found = true
		}
	}

	if found {
		t.Error("Failed to remove item from container")
	}
}

func TestAffectBitName(t *testing.T) {
	var result string

	result = affectBitName(affectBlind)
	if result != "blind" {
		t.Error("Failed to generate affect name from bit: ", affectBlind)
	}
	result = affectBitName(affectInvisible)
	if result != "invisible" {
		t.Error("Failed to generate affect name from bit: ", affectInvisible)
	}
	result = affectBitName(affectDetectEvil)
	if result != "detect_evil" {
		t.Error("Failed to generate affect name from bit: ", affectDetectEvil)
	}
	result = affectBitName(affectDetectInvisible)
	if result != "detect_invis" {
		t.Error("Failed to generate affect name from bit: ", affectDetectInvisible)
	}
	result = affectBitName(affectDetectMagic)
	if result != "detect_magic" {
		t.Error("Failed to generate affect name from bit: ", affectDetectMagic)
	}
	result = affectBitName(affectDetectHidden)
	if result != "detect_hidden" {
		t.Error("Failed to generate affect name from bit: ", affectDetectHidden)
	}
	result = affectBitName(affectHold)
	if result != "hold" {
		t.Error("Failed to generate affect name from bit: ", affectHold)
	}
	result = affectBitName(affectSanctuary)
	if result != "sanctuary" {
		t.Error("Failed to generate affect name from bit: ", affectSanctuary)
	}
	result = affectBitName(affectFaerieFire)
	if result != "faerie_fire" {
		t.Error("Failed to generate affect name from bit: ", affectFaerieFire)
	}
	result = affectBitName(affectInfrared)
	if result != "infrared" {
		t.Error("Failed to generate affect name from bit: ", affectInfrared)
	}
	result = affectBitName(affectCurse)
	if result != "curse" {
		t.Error("Failed to generate affect name from bit: ", affectCurse)
	}

	result = affectBitName(affectFlaming)
	if result != "flaming" {
		t.Error("Failed to generate affect name from bit: ", affectFlaming)
	}
	result = affectBitName(affectPoison)
	if result != "poison" {
		t.Error("Failed to generate affect name from bit: ", affectPoison)
	}
	result = affectBitName(affectProtect)
	if result != "protect" {
		t.Error("Failed to generate affect name from bit: ", affectProtect)
	}
	result = affectBitName(affectParalysis)
	if result != "paralysis" {
		t.Error("Failed to generate affect name from bit: ", affectParalysis)
	}
	result = affectBitName(affectSleep)
	if result != "sleep" {
		t.Error("Failed to generate affect name from bit: ", affectSleep)
	}
	result = affectBitName(affectSneak)
	if result != "sneak" {
		t.Error("Failed to generate affect name from bit: ", affectSneak)
	}
	result = affectBitName(affectHide)
	if result != "hide" {
		t.Error("Failed to generate affect name from bit: ", affectHide)
	}
	result = affectBitName(affectCharm)
	if result != "charm" {
		t.Error("Failed to generate affect name from bit: ", affectCharm)
	}
	result = affectBitName(affectFlying)
	if result != "flying" {
		t.Error("Failed to generate affect name from bit: ", affectFlying)
	}
	result = affectBitName(affectPassDoor)
	if result != "pass_door" {
		t.Error("Failed to generate affect name from bit: ", affectPassDoor)
	}

	result = affectBitName(0) // fake
	if result != "none" {
		t.Error("Generated name for imaginary affect")
	}

	result = affectBitName(3)
	if result != "blind invisible" {
		t.Error("Failed to generate multiple affect names")
	}
}

func TestAffectLocationName(t *testing.T) {
	var result string

	result = affectLocationName(applyNone)
	if result != "none" {
		t.Error("Failed to generate affect location name from bit: ", applyNone)
	}
	result = affectLocationName(applyStrength)
	if result != "strength" {
		t.Error("Failed to generate affect location name from bit: ", applyStrength)
	}
	result = affectLocationName(applyDexterity)
	if result != "dexterity" {
		t.Error("Failed to generate affect location name from bit: ", applyDexterity)
	}
	result = affectLocationName(applyIntelligence)
	if result != "intelligence" {
		t.Error("Failed to generate affect location name from bit: ", applyIntelligence)
	}
	result = affectLocationName(applyWisdom)
	if result != "wisdom" {
		t.Error("Failed to generate affect location name from bit: ", applyWisdom)
	}
	result = affectLocationName(applyConstitution)
	if result != "constitution" {
		t.Error("Failed to generate affect location name from bit: ", applyConstitution)
	}
	result = affectLocationName(applySex)
	if result != "sex" {
		t.Error("Failed to generate affect location name from bit: ", applySex)
	}
	result = affectLocationName(applyClass)
	if result != "class" {
		t.Error("Failed to generate affect location name from bit: ", applyClass)
	}
	result = affectLocationName(applyLevel)
	if result != "level" {
		t.Error("Failed to generate affect location name from bit: ", applyLevel)
	}
	result = affectLocationName(applyMana)
	if result != "mana" {
		t.Error("Failed to generate affect location name from bit: ", applyMana)
	}
	result = affectLocationName(applyHitpoints)
	if result != "hitpoints" {
		t.Error("Failed to generate affect location name from bit: ", applyHitpoints)
	}
	result = affectLocationName(applyMovement)
	if result != "movement" {
		t.Error("Failed to generate affect location name from bit: ", applyMovement)
	}
	result = affectLocationName(applyGold)
	if result != "gold" {
		t.Error("Failed to generate affect location name from bit: ", applyGold)
	}
	result = affectLocationName(applyExp)
	if result != "experience" {
		t.Error("Failed to generate affect location name from bit: ", applyExp)
	}
	result = affectLocationName(applyArmorClass)
	if result != "armor class" {
		t.Error("Failed to generate affect location name from bit: ", applyArmorClass)
	}
	result = affectLocationName(applyHitroll)
	if result != "hit roll" {
		t.Error("Failed to generate affect location name from bit: ", applyHitroll)
	}
	result = affectLocationName(applyDamroll)
	if result != "dam roll" {
		t.Error("Failed to generate affect location name from bit: ", applyDamroll)
	}
	result = affectLocationName(applySavingParalysis)
	if result != "save vs paralysis" {
		t.Error("Failed to generate affect location name from bit: ", applySavingParalysis)
	}
	result = affectLocationName(applySavingRod)
	if result != "save vs rod" {
		t.Error("Failed to generate affect location name from bit: ", applySavingRod)
	}
	result = affectLocationName(applySavingPetrify)
	if result != "save vs petrify" {
		t.Error("Failed to generate affect location name from bit: ", applySavingPetrify)
	}
	result = affectLocationName(applySavingBreath)
	if result != "save vs breath" {
		t.Error("Failed to generate affect location name from bit: ", applySavingBreath)
	}
	result = affectLocationName(applySavingSpell)
	if result != "save vs spell" {
		t.Error("Failed to generate affect location name from bit: ", applySavingSpell)
	}

	if affectLocationName(9999) != "(unknown)" {
		t.Error("Failed to generate affect location name from nothing")
	}
}

func TestExtraBitName(t *testing.T) {

	if extraBitName(itemGlow) != "glow" {
		t.Error("Failed to generate extra bit name for bit: ", itemGlow)
	}

	if extraBitName(itemHum) != "hum" {
		t.Error("Failed to generate extra bit name for bit: ", itemHum)
	}

	if extraBitName(itemDark) != "dark" {
		t.Error("Failed to generate extra bit name for bit: ", itemDark)
	}

	if extraBitName(itemLock) != "lock" {
		t.Error("Failed to generate extra bit name for bit: ", itemLock)
	}

	if extraBitName(itemEvil) != "evil" {
		t.Error("Failed to generate extra bit name for bit: ", itemEvil)
	}

	if extraBitName(itemInvis) != "invis" {
		t.Error("Failed to generate extra bit name for bit: ", itemInvis)
	}

	if extraBitName(itemMagic) != "magic" {
		t.Error("Failed to generate extra bit name for bit: ", itemMagic)
	}

	if extraBitName(itemNoDrop) != "nodrop" {
		t.Error("Failed to generate extra bit name for bit: ", itemNoDrop)
	}

	if extraBitName(itemBless) != "bless" {
		t.Error("Failed to generate extra bit name for bit: ", itemBless)
	}

	if extraBitName(itemAntiEvil) != "anti-evil" {
		t.Error("Failed to generate extra bit name for bit: ", itemAntiEvil)
	}

	if extraBitName(itemAntiGood) != "anti-good" {
		t.Error("Failed to generate extra bit name for bit: ", itemAntiGood)
	}

	if extraBitName(itemAntiNeutral) != "anti-neutral" {
		t.Error("Failed to generate extra bit name for bit: ", itemGlow)
	}

	if extraBitName(itemNoRemove) != "noremove" {
		t.Error("Failed to generate extra bit name for bit: ", itemNoRemove)
	}

	if extraBitName(itemInventory) != "inventory" {
		t.Error("Failed to generate extra bit name for bit: ", itemInventory)
	}

	if extraBitName(0) != "none" {
		t.Error("Generated incorrect namem for 0 bit")
	}
}

func TestItemTypeName(t *testing.T) {

	if itemTypeName(&item{ItemType: itemLight}) != "light" {
		t.Error("Failed to get item type: ", itemLight)
	}
	if itemTypeName(&item{ItemType: itemScroll}) != "scroll" {
		t.Error("Failed to get item type: ", itemScroll)
	}
	if itemTypeName(&item{ItemType: itemWand}) != "wand" {
		t.Error("Failed to get item type: ", itemWand)
	}
	if itemTypeName(&item{ItemType: itemStaff}) != "staff" {
		t.Error("Failed to get item type: ", itemStaff)
	}
	if itemTypeName(&item{ItemType: itemWeapon}) != "weapon" {
		t.Error("Failed to get item type: ", itemWeapon)
	}
	if itemTypeName(&item{ItemType: itemTreasure}) != "treasure" {
		t.Error("Failed to get item type: ", itemTreasure)
	}
	if itemTypeName(&item{ItemType: itemArmor}) != "armor" {
		t.Error("Failed to get item type: ", itemArmor)
	}
	if itemTypeName(&item{ItemType: itemPotion}) != "potion" {
		t.Error("Failed to get item type: ", itemPotion)
	}
	if itemTypeName(&item{ItemType: itemFurniture}) != "furniture" {
		t.Error("Failed to get item type: ", itemFurniture)
	}
	if itemTypeName(&item{ItemType: itemTrash}) != "trash" {
		t.Error("Failed to get item type: ", itemTrash)
	}
	if itemTypeName(&item{ItemType: itemContainer}) != "container" {
		t.Error("Failed to get item type: ", itemContainer)
	}
	if itemTypeName(&item{ItemType: itemDrinkContainer}) != "drink container" {
		t.Error("Failed to get item type: ", itemDrinkContainer)
	}
	if itemTypeName(&item{ItemType: itemKey}) != "key" {
		t.Error("Failed to get item type: ", itemKey)
	}
	if itemTypeName(&item{ItemType: itemFood}) != "food" {
		t.Error("Failed to get item type: ", itemFood)
	}
	if itemTypeName(&item{ItemType: itemMoney}) != "money" {
		t.Error("Failed to get item type: ", itemMoney)
	}
	if itemTypeName(&item{ItemType: itemBoat}) != "boat" {
		t.Error("Failed to get item type: ", itemBoat)
	}
	if itemTypeName(&item{ItemType: itemCorpseNPC}) != "npc corpse" {
		t.Error("Failed to get item type: ", itemCorpseNPC)
	}
	if itemTypeName(&item{ItemType: itemCorpsePC}) != "pc corpse" {
		t.Error("Failed to get item type: ", itemCorpsePC)
	}
	if itemTypeName(&item{ItemType: itemFountain}) != "fountain" {
		t.Error("Failed to get item type: ", itemFountain)
	}
	if itemTypeName(&item{ItemType: itemPill}) != "pill" {
		t.Error("Failed to get item type: ", itemPill)
	}

	if itemTypeName(&item{ItemType: 19352}) != "(unknown)" {
		t.Error("Found incorrect naem for fake type")
	}
}

func TestCreateMoney(t *testing.T) {
	gameServer.BasePath = "../"
	loadItems()
	money := createMoney(1)

	if money.Value != 1 {
		t.Error("Failed to create money with value: 1")
	}

	if money.index.ID != goldSingle {
		t.Error("Created incorrect type of money. Expected ", goldSingle)
	}

	money = createMoney(0)
	if money.Value != 1 {
		t.Error("Failed to create money with value: 0")
	}

	money = createMoney(10)
	if money.Value != 10 {
		t.Error("Failed to create money with value: 0")
	}

	if money.index.ID != goldMulti {
		t.Error("Created incorrect type of money. Expected ", goldMulti)
	}
}

func TestApplyAC(t *testing.T) {
	armor := &item{ItemType: itemArmor, Value: 1}

	if applyAC(armor, wearBody) != 3 {
		t.Error("Failed to apply AC for body")
	}

	if applyAC(armor, wearHead) != 2 {
		t.Error("Failed to apply AC for head")
	}

	if applyAC(armor, wearLegs) != 2 {
		t.Error("Failed to apply AC for legs")
	}

	if applyAC(armor, wearFeet) != 1 {
		t.Error("Failed to apply AC for feet")
	}

	if applyAC(armor, wearHands) != 1 {
		t.Error("Failed to apply AC for hands")
	}

	if applyAC(armor, wearArms) != 1 {
		t.Error("Failed to apply AC for arms")
	}

	if applyAC(armor, wearShield) != 1 {
		t.Error("Failed to apply AC for shield")
	}

	if applyAC(armor, wearFingerLeft) != 1 {
		t.Error("Failed to apply AC for left finger")
	}

	if applyAC(armor, wearFingerRight) != 1 {
		t.Error("Failed to apply AC for right finger")
	}

	if applyAC(armor, wearNeck1) != 1 {
		t.Error("Failed to apply AC for neck 1")
	}

	if applyAC(armor, wearNeck2) != 1 {
		t.Error("Failed to apply AC for neck 2")
	}

	if applyAC(armor, wearArmor) != 2 {
		t.Error("Failed to apply AC for armor")
	}

	if applyAC(armor, wearWaist) != 1 {
		t.Error("Failed to apply AC for waist")
	}

	if applyAC(armor, wearWristLeft) != 1 {
		t.Error("Failed to apply AC for left wrist")
	}

	if applyAC(armor, wearWristRight) != 1 {
		t.Error("Failed to apply AC for right wrist")
	}
	if applyAC(armor, wearHold) != 1 {
		t.Error("Failed to apply AC for hold")
	}

	if applyAC(armor, 325252) != 0 {
		t.Error("Failed to return 0 for non-existent location")
	}

	armor.ItemType = itemFood
	if applyAC(armor, wearArmor) != 0 {
		t.Error("Failed to return 0 for non-armor item")
	}
}
