package mud

import (
	"fmt"
)

func (m *mob) addItem(item *item) {
	m.Inventory = append(m.Inventory, item)
}

func (m *mob) canDropItem(item *item) bool {
	if !hasBit(item.ExtraFlags, itemNoDrop) {
		return true
	}

	if !m.isNPC() && m.Level >= 99 {
		return true
	}

	return false
}

func (m *mob) canSeeItem(item *item) bool {
	return true
}

func (m *mob) carrying(str string) *item {
	for _, i := range m.Inventory {
		if matchesSubject(i.Name, str) {
			return i
		}
	}
	return nil
}

func (m *mob) equippedName(name string) *item {
	for _, i := range m.Equipped {
		if matchesSubject(i.Name, name) {
			return i
		}
	}

	return nil
}

func (m *mob) get(item *item, container *item) {
	if !item.canWear(itemTake) {
		m.notify("You can't take that.")
		return
	}

	if m.Carrying+1 > m.CarryMax {
		m.notify("You can't carry that many items.")
		return
	}

	if m.CarryWeight+item.Weight > m.CarryWeightMax {
		m.notify("You can't carry that much weight.")
		return
	}

	if container != nil {
		m.notify("You get %s from %s.", item.Name, container.Name)
		m.Room.notify(fmt.Sprintf("%s gets %s from %s.%s", m.Name, item.Name, container.Name, Newline), m)
		container.removeObject(item)
	} else {
		m.notify("You get %s.", item.Name)
		m.Room.notify(fmt.Sprintf("%s gets %s.%s", m.Name, item.Name, Newline), m)
		m.Room.removeObject(item)
	}

	if item.ItemType == itemMoney {
		m.Gold += item.Value
	} else {
		m.Inventory = append(m.Inventory, item)
		item.carriedBy = m
	}
}

func (m *mob) removeItem(i *item) {
	for j, it := range m.Inventory {
		if it.index.ID == i.index.ID {
			m.Inventory = append(m.Inventory[:j], m.Inventory[j+1:]...)
			return
		}
	}
}

func (m *mob) sacrifice(args []string) {
	if len(args) < 1 {
		act("$n offers $mself to the gods, who don't bother to answer.", m, nil, nil, actToRoom)
		m.notify("The gods aren't listening.")
		return
	}

	obj := m.carrying(args[1])
	if obj == nil {
		m.notify("You can't find it.")
		return
	}

	if !obj.canWear(itemTake) {
		act("$p is not an acceptable sacrifice.", m, obj, nil, actToChar)
		return
	}

	m.notify("The gods grant you a single gold coin for your sacrifice.")
	m.Gold++

	act("$n sacrifices $p to the gods.", m, obj, nil, actToRoom)
	m.removeItem(obj)
	return
}

func (m *mob) unwearItem(location int, replace bool) bool {
	obj := m.equippedItem(location)

	if obj == nil {
		return true
	}

	if !replace {
		return false
	}

	if hasBit(obj.ExtraFlags, itemNoRemove) {
		act("You can't remove $p.", m, obj, nil, actToChar)
		return false
	}

	act("$n stops using $p.", m, obj, nil, actToRoom)
	act("You stop using $p.", m, obj, nil, actToChar)
	return true
}

func (m *mob) unequipItem(item *item) {
	if item.WearLocation == wearNone {
		return
	}

	m.Armor += applyAC(item, int(item.WearLocation))
	item.WearLocation = wearNone

	for _, af := range item.index.Affected {
		affectModify(m, af, false)
	}
	for _, af := range item.Affected {
		affectModify(m, af, false)
	}

	if item.ItemType == itemLight && m.Room != nil && m.Room.Light > 0 {
		item.Room.Light--
	}

	return
}

func (m *mob) wear(wearable *item, replace bool) {
	if m.Level < wearable.Level {
		m.notify("You must be level %d to wear this.", wearable.Level)
		return
	}

	if wearable.ItemType == itemLight {
		removed := m.unwearItem(wearLight, replace)
		if !removed {
			return
		}
		m.notify("You light up %s and hold it.", wearable.Name)
		m.Room.notify(fmt.Sprintf("%s lights up %s and holds it.", m.Name, wearable.Name), m)
		m.equipItem(wearable, wearLight)
		return
	}

	if wearable.canWear(itemWearFinger) {
		removedLeft := m.unwearItem(wearFingerLeft, replace)
		removedRight := m.unwearItem(wearFingerRight, replace)
		if m.equippedItem(wearFingerLeft) != nil && m.equippedItem(wearFingerRight) != nil && !removedLeft && !removedRight {
			return
		}

		if m.equippedItem(wearFingerLeft) == nil {
			m.notify("You wear %s on your left finger.", wearable.Name)
			m.Room.notify(fmt.Sprintf("%s wears %s on their left finger.", m.Name, wearable.Name), m)
			m.equipItem(wearable, wearFingerLeft)
			return
		}

		if m.equippedItem(wearFingerRight) == nil {
			m.notify("You wear %s on your right finger.", wearable.Name)
			m.Room.notify(fmt.Sprintf("%s wears %s on their right finger.", m.Name, wearable.Name), m)
			m.equipItem(wearable, wearFingerRight)
			return
		}

		m.notify("You already wear two rings!")
		return
	}

	if wearable.canWear(itemWearNeck) {
		removedNeck1 := m.unwearItem(wearNeck1, replace)
		removedNeck2 := m.unwearItem(wearNeck2, replace)
		if m.equippedItem(wearNeck1) != nil && m.equippedItem(wearNeck2) != nil && !removedNeck1 && !removedNeck2 {
			return
		}

		if m.equippedItem(wearNeck1) == nil {
			m.notify("You wear %s on your neck.", wearable.Name)
			m.Room.notify(fmt.Sprintf("%s wears %s on their neck.", m.Name, wearable.Name), m)
			m.equipItem(wearable, wearNeck1)
			return
		}

		if m.equippedItem(wearNeck2) == nil {
			m.notify("You wear %s on your neck.", wearable.Name)
			m.Room.notify(fmt.Sprintf("%s wears %s on their neck.", m.Name, wearable.Name), m)
			m.equipItem(wearable, wearNeck2)
			return
		}

		m.notify("You already wear two neck items!")
		return
	}
	if wearable.canWear(itemWearWrist) {
		removedWristLeft := m.unwearItem(wearWristLeft, replace)
		removedWristRight := m.unwearItem(wearWristRight, replace)
		if m.equippedItem(wearWristLeft) != nil && m.equippedItem(wearWristRight) != nil && !removedWristLeft && !removedWristRight {
			return
		}

		if m.equippedItem(wearWristLeft) == nil {
			m.notify("You wear %s on your left wrist.", wearable.Name)
			m.Room.notify(fmt.Sprintf("%s wears %s on their left wrist.", m.Name, wearable.Name), m)
			m.equipItem(wearable, wearWristLeft)
			return
		}

		if m.equippedItem(wearWristRight) == nil {
			m.notify("You wear %s on your right wrist.", wearable.Name)
			m.Room.notify(fmt.Sprintf("%s wears %s on their right wrist.", m.Name, wearable.Name), m)
			m.equipItem(wearable, wearWristRight)
			return
		}

		m.notify("You already wear two wrist items!")
		return
	}

	if wearable.canWear(itemWearBody) {
		removed := m.unwearItem(wearBody, replace)
		if !removed {
			return
		}
		m.notify("You wear %s on your body.", wearable.Name)
		m.Room.notify(fmt.Sprintf("%s wears %s on their body.", m.Name, wearable.Name), m)
		m.equipItem(wearable, wearBody)
		return
	}

	if wearable.canWear(itemWearHead) {
		removed := m.unwearItem(wearHead, replace)
		if !removed {
			return
		}
		m.notify("You wear %s on your head.", wearable.Name)
		m.Room.notify(fmt.Sprintf("%s wears %s on their head.", m.Name, wearable.Name), m)
		m.equipItem(wearable, wearHead)
		return
	}

	if wearable.canWear(itemWearLegs) {
		removed := m.unwearItem(wearLegs, replace)
		if !removed {
			return
		}
		m.notify("You wear %s on your legs.", wearable.Name)
		m.Room.notify(fmt.Sprintf("%s wears %s on their legs.", m.Name, wearable.Name), m)
		m.equipItem(wearable, wearLegs)
		return
	}

	if wearable.canWear(itemWearFeet) {
		removed := m.unwearItem(wearFeet, replace)
		if !removed {
			return
		}
		m.notify("You wear %s on your feet.", wearable.Name)
		m.Room.notify(fmt.Sprintf("%s wears %s on their feet.", m.Name, wearable.Name), m)
		m.equipItem(wearable, wearFeet)
		return
	}

	if wearable.canWear(itemWearHands) {
		removed := m.unwearItem(wearHands, replace)
		if !removed {
			return
		}
		m.notify("You wear %s on your hands.", wearable.Name)
		m.Room.notify(fmt.Sprintf("%s wears %s on their hands.", m.Name, wearable.Name), m)
		m.equipItem(wearable, wearHands)
		return
	}

	if wearable.canWear(itemWearWaist) {
		removed := m.unwearItem(wearWaist, replace)
		if !removed {
			return
		}
		m.notify("You wear %s on your waist.", wearable.Name)
		m.Room.notify(fmt.Sprintf("%s wears %s on their waist.", m.Name, wearable.Name), m)
		m.equipItem(wearable, wearWaist)
		return
	}

	if wearable.canWear(itemWearShield) {
		removed := m.unwearItem(wearShield, replace)
		if !removed {
			return
		}
		m.notify("You wear %s as your shield.", wearable.Name)
		m.Room.notify(fmt.Sprintf("%s wears %s as their shield.", m.Name, wearable.Name), m)
		m.equipItem(wearable, wearShield)
		return
	}

	if wearable.canWear(itemWearHold) {
		removed := m.unwearItem(wearHold, replace)
		if !removed {
			return
		}
		m.notify("You hold %s.", wearable.Name)
		m.Room.notify(fmt.Sprintf("%s holds %s.", m.Name, wearable.Name), m)
		m.equipItem(wearable, wearHold)
		return
	}

	if wearable.canWear(itemWearWield) {
		removed := m.unwearItem(wearWield, replace)
		if !removed {
			return
		}
		if wearable.Weight > m.ModifiedAttributes.Strength {
			m.notify("It is too heavy for you to wield.")
			return
		}

		m.notify("You wield %s.", wearable.Name)
		m.Room.notify(fmt.Sprintf("%s wields %s.", m.Name, wearable.Name), m)
		m.equipItem(wearable, wearWield)
		return
	}

	m.notify("You can't wear, wield, or hold that.")
}
