package mud

import (
	"fmt"

	"github.com/brianseitel/oasis-mud/helpers"
)

func (m *mob) addItem(item *item) {
	m.Inventory = append(m.Inventory, item)
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

func (m *mob) removeItem(i *item) {
	for j, it := range m.Inventory {
		if it == i {
			m.Inventory = append(m.Inventory[:j], m.Inventory[j+1:]...)
			return
		}
	}
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
