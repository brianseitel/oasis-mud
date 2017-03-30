package mud

import (
	"fmt"

	"github.com/brianseitel/oasis-mud/helpers"
)

func (m *mob) addItem(item *item) {
	m.Inventory = append(m.Inventory, item)
}

func (m *mob) canDropItem(item *item) bool {
	if !helpers.HasBit(item.ExtraFlags, itemNoDrop) {
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
		if helpers.MatchesSubject(i.Name, str) {
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
		m.Room.notify(fmt.Sprintf("%s gets %s from %s.%s", m.Name, item.Name, container.Name, helpers.Newline), m)
		container.removeObject(item)
	} else {
		m.notify("You get %s.", item.Name)
		m.Room.notify(fmt.Sprintf("%s gets %s.%s", m.Name, item.Name, helpers.Newline), m)
		m.Room.removeObject(item)
	}

	if item.ItemType == itemMoney {
		m.Gold += uint(item.Value)
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
