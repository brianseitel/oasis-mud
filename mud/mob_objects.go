package mud

import (
	"fmt"
	"strconv"
	"strings"
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

func doBrandish(player *mob, argument string) {
	staff := player.equippedItem(wearHold)
	if staff == nil {
		player.notify("You hold nothing in your hand.")
		return
	}

	if staff.ItemType != itemStaff {
		player.notify("You can only brandish with a staff.")
		return
	}

	brandish := player.skill("brandish")
	if brandish == nil {
		return
	}

	wait(player, 2*pulseViolence)

	if staff.Skill != nil {
		act("$n brandishes $p.", player, staff, nil, actToRoom)
		act("You brandish $p.", player, staff, nil, actToChar)

		for _, victim := range player.Room.Mobs {
			switch staff.Skill.Target {
			case targetIgnore:
				if victim == player {
					continue
				}
				break

			case targetCharacterOffensive:
				if player.isNPC() {
					if victim.isNPC() {
						continue
					} else if !victim.isNPC() {
						continue
					}
				}
				break

			case targetCharacterDefensive:
				if player.isNPC() {
					if !victim.isNPC() {
						continue
					} else if victim.isNPC() {
						continue
					}
				}
				break

			case targetCharacterSelf:
				if victim != player {
					continue
				}
				break
			}

			objCastSpell(staff.Skill, staff.Level, player, victim, nil)
		}
	}

	staff.Charges--
	if staff.Charges <= 0 {
		act("$n's $p blazes brightly and vanishes in a puff of rainbow sprinkles.", player, staff, nil, actToRoom)
		act("Your $p blazes brightly and vanishes in a puff of rainbow sprinkles.", player, staff, nil, actToChar)
		extractObj(staff)
	}
	return
}

func doDrop(player *mob, argument string) {
	if len(argument) <= 1 {
		player.notify("Drop what?")
		return
	}

	argument, arg1 := oneArgument(argument)
	argument, arg2 := oneArgument(argument)

	num, err := strconv.Atoi(arg1)
	isNumber := err == nil

	if isNumber {
		amount := num
		if arg2 != "" || amount <= 0 || !strings.HasPrefix(arg2, "gold") {
			player.notify("Sorry, you can't do that.")
			return
		}

		if player.Gold < amount {
			player.notify("You haven't got that many coins.")
			return
		}

		player.Gold -= amount

		// TODO: see if we already have gold in the room

		player.Room.Items = append(player.Room.Items, createMoney(amount))
		player.notify("OK.")
		player.Room.notify(fmt.Sprintf("%s drops some gold.", player.Name), player)
		return
	}

	if arg1 != "all" && !strings.HasPrefix(arg1, "all.") {
		// drop obj
		var item *item
		for _, i := range player.Inventory {
			if matchesSubject(i.Name, arg1) {
				item = i
				break
			}
		}

		if item == nil {
			player.notify("You don't have that item.")
			return
		}

		if !player.canDropItem(item) {
			player.notify("You can't let go of it.")
			return
		}

		for j, i := range player.Inventory {
			if i == item {
				item.carriedBy = nil
				item.Room = player.Room
				player.Inventory, player.Room.Items = transferItem(j, player.Inventory, player.Room.Items)
				player.notify("You drop %s.", i.Name)
				player.Room.notify(fmt.Sprintf("%s drops %s.", player.Name, i.Name), player)
				break
			}
		}
	} else {
		// drop all
		found := false

		words := strings.SplitN(arg1, ".", 2)
		var name string
		if len(words) > 1 {
			name = words[1]
		}

		for j := 0; j < len(player.Inventory); j++ {
			item := player.Inventory[j]

			if !player.canDropItem(item) || item.WearLocation != wearNone {
				continue
			}
			if arg1 == "all" || matchesSubject(item.Name, name) {
				found = true

				item.carriedBy = nil
				item.Room = player.Room
				player.Inventory = append(player.Inventory[0:j], player.Inventory[j+1:]...)
				j--
				player.Room.Items = append(player.Room.Items, item)
				act("$n drops $p.", player, item, nil, actToRoom)
				act("You drop $p.", player, item, nil, actToChar)
			}
		}

		if !found {
			if len(name) == 0 {
				player.notify("You are not carrying anything.")
			} else {
				player.notify("You are not carrying any %s.", arg1)
			}
		}
	}

	return
}

func doEat(player *mob, argument string) {
	argument, arg1 := oneArgument(argument)

	if arg1 == "" {
		player.notify("Eat what?")
		return
	}

	var obj *item
	for _, i := range player.Inventory {
		if matchesSubject(i.Name, arg1) {
			obj = i
			break
		}
	}

	if obj == nil {
		player.notify("You don't have that item.")
		return
	}

	if !player.isImmortal() {
		if obj.ItemType != itemFood && obj.ItemType != itemPill {
			player.notify("That isn't edible.")
			return
		}
	}

	act("$n eats $p.", player, obj, nil, actToRoom)
	act("You eat $p.", player, obj, nil, actToChar)

	switch obj.ItemType {
	case itemPill:
		objCastSpell(obj.Skill, obj.Min, player, player, nil)
		break
	}

	extractObj(obj)
}

func doGet(player *mob, argument string) {
	if len(argument) == 0 {
		player.notify("Get what?")
		return
	}

	argument, arg1 := oneArgument(argument)
	argument, arg2 := oneArgument(argument)

	if arg2 == "" {
		if arg1 != "all" && !strings.HasPrefix(arg1, "all.") {
			// get obj
			for _, i := range player.Room.Items {
				if matchesSubject(i.Description, arg1) {
					player.get(i, nil)
					return
				}
			}

			player.notify("I see no %s here.", arg1)
		} else {
			// get all or get all.container
			words := strings.SplitN(arg1, ".", 2)
			var name string
			if len(words) > 1 {
				name = words[1]
			}

			found := false
			if len(player.Room.Items) > 0 {
				for _, i := range player.Room.Items {
					if matchesSubject(i.Description, name) || len(name) == 0 {
						player.get(i, nil)
						found = true
					}
				}
			}

			if !found {
				if len(name) == 0 {
					player.notify("I see nothing here.")
				} else {
					player.notify("I see no %s here.", name)
				}
			}
		}
	} else {
		// get ... container
		if arg2 == "all" || strings.HasPrefix(arg2, "all.") {
			player.notify("You can't do that.")
			return
		}

		var container *item
		for _, i := range player.Room.Items {
			if matchesSubject(i.Description, arg2) {
				container = i
				break
			}
		}

		if container == nil {
			// try from inventory
			for _, i := range player.Inventory {
				if matchesSubject(i.Description, arg1) {
					container = i
					break
				}
			}
		}

		if container == nil {
			player.notify("I see no %s here.", arg2)
			return
		}

		switch container.ItemType {
		case itemContainer:
		case itemCorpseNPC:
			break

		case itemCorpsePC:
			player.notify("You can't do that.%s.")
			return
		default:
			player.notify("That's not a container.")
			return
		}

		if hasBit(container.Value, containerClosed) {
			player.notify("The %s is closed.", container.Name)
			return
		}

		if arg1 != "all" && !strings.HasPrefix(arg1, "all.") {
			// get obj container
			for _, i := range container.container {
				if matchesSubject(i.Description, arg1) {
					player.get(i, container)
					return
				}
			}

			player.notify("I see nothing like that in %s.", container.Name)
		} else {
			// get all container or get all.obj container
			words := strings.SplitN(arg2, ".", 2)
			var name string
			if len(words) > 1 {
				name = words[1]
			}
			found := false
			for _, i := range container.container {
				if matchesSubject(i.Description, name) || len(name) == 0 {
					player.get(i, container)
					found = true
				}
			}

			if !found {
				if len(name) == 0 {
					player.notify("I see nothing in the %s.", container.Name)
				} else {
					player.notify("I see nothing like that in %s.", container.Name)
				}
			}
		}
	}
}

func doGive(player *mob, argument string) {
	if len(argument) <= 1 {
		player.notify("Give what to whom?")
		return
	}

	argument, arg1 := oneArgument(argument)
	argument, arg2 := oneArgument(argument)

	var victim *mob
	for _, mob := range player.Room.Mobs {
		if matchesSubject(mob.Name, arg2) {
			victim = mob
			break
		}
	}

	num, err := strconv.Atoi(arg1)
	isNumeric := err == nil

	var amount int
	if isNumeric {
		if num <= 0 || !strings.HasPrefix(arg2, "coin") {
			player.notify("You can't do that.")
			return
		}
		amount = num

		if victim == nil {
			player.notify("They aren't here.")
			return
		}

		if player.Gold < amount {
			player.notify("You don't have that much gold.")
			return
		}

		player.Gold -= amount
		victim.Gold += amount
		player.notify("You give %s some gold.", victim.Name)
		victim.notify("%s gives you some gold.", player.Name)
		return
	}

	var item *item
	for _, i := range player.Inventory {
		if matchesSubject(i.Name, arg1) {
			item = i
			break
		}
	}

	if item == nil {
		player.notify("You do not have that item.")
		return
	}

	if item.WearLocation != wearNone {
		player.notify("You must remove it first.")
		return
	}

	if victim.Room.ID != player.Room.ID {
		player.notify("They are not here.")
		return
	}

	if !player.canDropItem(item) {
		player.notify("You can't let go of it.")
		return
	}

	if victim.Carrying+1 > victim.CarryMax {
		player.notify("%s has their hands full.", victim.Name)
		return
	}

	if victim.CarryWeight+item.Weight > victim.CarryWeightMax {
		player.notify("%s can't carry that much weight.", victim.Name)
		return
	}

	if victim.canSeeItem(item) {
		player.notify("%s can't see it.", victim.Name)
	}

	for j, it := range player.Inventory {
		if it == item {
			player.Inventory = append(player.Inventory[:j], player.Inventory[j+1:]...)
			break
		}
	}
	victim.Inventory = append(victim.Inventory, item)

	player.notify("You give %s to %s.", item.Name, victim.Name)
	victim.notify("%s gives you %s.", player.Name, item.Name)
	return
}

func doPut(player *mob, argument string) {

	if len(argument) <= 1 {
		player.notify("Put what in what?")
		return
	}

	argument, arg1 := oneArgument(argument)
	argument, arg2 := oneArgument(argument)

	if arg2 == "all" || strings.HasPrefix(arg2, "all.") {
		player.notify("You can't do that.")
		return
	}

	var container *item
	for _, i := range player.Inventory {
		if strings.HasPrefix(i.Name, arg2) {
			container = i
			break
		}
	}

	if container == nil {
		// try from room
		for _, i := range player.Room.Items {
			if strings.HasPrefix(i.Name, arg2) {
				container = i
				break
			}
		}

	}

	if container == nil {
		player.notify("I see no %s here.", arg2)
		return
	}

	if hasBit(container.Value, containerClosed) {
		player.notify("The %s is closed.", container.Name)
		return
	}

	if arg1 != "all" && !strings.HasPrefix(arg1, "all.") {
		// put obj container
		var item *item
		for _, i := range player.Inventory {
			if matchesSubject(i.Name, arg1) {
				item = i
				break
			}
		}

		if item == nil {
			player.notify("You do not have that item.")
			return
		}

		if item == container {
			player.notify("You can't fold it into itself!")
			return
		}

		if !player.canDropItem(item) {
			player.notify("You can't let go of it.")
			return
		}

		if item.Weight+container.Weight > container.Value {
			player.notify("It won't fit.")
			return
		}

		for j, it := range player.Inventory {
			if it == item {
				player.Inventory = append(player.Inventory[0:j], player.Inventory[j+1:]...)
				break
			}
		}

		container.container = append(container.container, item)

		player.notify("You put %s in %s.", item.Name, container.Name)
		player.Room.notify(fmt.Sprintf("%s puts %s in %s.", player.Name, item.Name, container.Name), player)
	} else {
		// put all container or put all.object container
		words := strings.SplitN(arg1, ".", 2)
		var name string
		if len(words) > 1 {
			name = words[1]
		}
		for j, item := range player.Inventory {
			if (arg1 == "all" || strings.HasPrefix(item.Name, name)) && item.WearLocation == wearNone && item != container && item.Weight+container.Weight > container.Value {
				player.Inventory = append(player.Inventory[0:j], player.Inventory[j+1:]...)
				container.container = append(container.container, item)
				player.notify("You put %s in %s.", item.Name, container.Name)
				player.Room.notify(fmt.Sprintf("%s puts %s in %s.", player.Name, item.Name, container.Name), player)
			}
		}
	}
}

func doQuaff(player *mob, argument string) {
	argument, arg1 := oneArgument(argument)

	if arg1 == "" {
		player.notify("Quaff what?")
		return
	}

	var obj *item
	for _, i := range player.Inventory {
		if matchesSubject(i.Name, arg1) {
			obj = i
			break
		}
	}

	if obj == nil {
		player.notify("You don't have that item.")
		return
	}

	if obj.ItemType != itemPotion {
		player.notify("You can quaff only potions.")
		return
	}

	act("$n quaffs $p.", player, obj, nil, actToRoom)
	act("You quaff $p.", player, obj, nil, actToChar)

	objCastSpell(obj.Skill, obj.Min, player, player, nil)

	extractObj(obj)
}

func doRecite(player *mob, argument string) {
	argument, arg1 := oneArgument(argument)

	if arg1 == "" {
		player.notify("Recite what?")
		return
	}

	var obj *item
	for _, i := range player.Inventory {
		if matchesSubject(i.Name, arg1) {
			obj = i
			break
		}
	}

	if obj == nil {
		player.notify("You don't have that item.")
		return
	}

	if obj.ItemType != itemScroll {
		player.notify("You can recite only scrolls.")
		return
	}

	act("$n recites $p.", player, obj, nil, actToRoom)
	act("You recite $p.", player, obj, nil, actToChar)

	objCastSpell(obj.Skill, obj.Min, player, player, nil)

	extractObj(obj)
}

func doWear(player *mob, argument string) {
	if len(argument) < 1 {
		player.notify("Wear, wield, or hold what?")
		return
	}

	argument, arg1 := oneArgument(argument)

	if arg1 == "all" {
		for _, i := range player.Inventory {
			if i.WearLocation == wearNone && player.canSeeItem(i) {
				player.wear(i, false)
			}
		}
	} else {
		wearable := player.carrying(arg1)

		if wearable == nil {
			player.notify("You can't find that.")
			return
		}

		player.wear(wearable, true)
	}
}

func doRemove(player *mob, argument string) {

	if len(argument) < 1 {
		player.notify("Remove what?")
		return
	}

	argument, arg1 := oneArgument(argument)
	obj := player.equippedName(arg1)
	if obj == nil {
		player.notify("You aren't wearing that item.")
		return
	}

	player.unwearItem(obj.WearLocation, true)
	return
}

func doZap(player *mob, argument string) {
	argument, arg1 := oneArgument(argument)

	if player.Fight == nil && arg1 == "" {
		player.notify("Zap whom or what?")
		return
	}

	wand := player.equippedItem(wearHold)
	if wand == nil {
		player.notify("You aren't holding anything.")
		return
	}

	if wand.ItemType != itemWand {
		player.notify("You can only zap with a wind.")
		return
	}

	var obj *item
	var victim *mob
	if arg1 == "" {
		if player.Fight != nil {
			victim = player.Fight
		} else {
			player.notify("Zap whom or what?")
			return
		}
	} else {
		for _, m := range player.Room.Mobs {
			if matchesSubject(m.Name, arg1) {
				victim = m
				break
			}
		}

		for _, o := range player.Inventory {
			if matchesSubject(o.Name, arg1) {
				obj = o
				break
			}
		}

		if obj == nil {
			for _, o := range player.Room.Items {
				if matchesSubject(o.Name, arg1) {
					obj = o
					break
				}
			}
		}

		if obj == nil && victim == nil {
			player.notify("You can't find it.")
		}
	}

	wait(player, 2*pulseViolence)

	if wand.Skill != nil {
		if victim != nil {
			act("$n zap $N with $p.", player, wand, victim, actToRoom)
			act("You zap $N with $p.", player, wand, victim, actToChar)
		} else {
			act("$n zaps $P with $p.", player, wand, obj, actToRoom)
			act("You zap $P with $p.", player, wand, obj, actToChar)
		}

		objCastSpell(wand.Skill, wand.Level, player, victim, obj)
	}

	wand.Charges--

	if wand.Charges <= 0 {
		act("$n's $p explodes into a shower of rainbow sprinkles.", player, wand, nil, actToRoom)
		act("Your $p explodes into a shower of rainbow sprinkles.", player, wand, nil, actToChar)
		extractObj(wand)
	}
}

func (m *mob) equippedName(name string) *item {
	for _, i := range m.Inventory {
		if i.WearLocation != wearNone && matchesSubject(i.Name, name) {
			return i
		}
	}

	return nil
}

func doEquipment(player *mob, argument string) {
	player.notify(
		fmt.Sprintf("Equipment\n%s\n%s\n%s",
			"-----------------------------------",
			strings.Join(equippedString(player), newline),
			"-----------------------------------",
		),
	)
}

func (m *mob) get(item *item, container *item) {
	if !hasBit(item.ItemType, itemTake) {
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
		m.Room.notify(fmt.Sprintf("%s gets %s from %s.%s", m.Name, item.Name, container.Name, newline), m)
		container.removeObject(item)
	} else {
		m.notify("You get %s.", item.Name)
		m.Room.notify(fmt.Sprintf("%s gets %s.%s", m.Name, item.Name, newline), m)
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

func doSacrifice(player *mob, argument string) {
	if len(argument) < 1 {
		act("$n offers $mself to the gods, who don't bother to answer.", player, nil, nil, actToRoom)
		player.notify("The gods aren't listening.")
		return
	}

	argument, arg1 := oneArgument(argument)
	obj := player.carrying(arg1)
	if obj == nil {
		for _, i := range player.Room.Items {
			if matchesSubject(i.Name, arg1) {
				obj = i
				break
			}
		}
	}

	if obj == nil {
		player.notify("You can't find it.")
		return
	}

	if !obj.canWear(itemTake) {
		act("$p is not an acceptable sacrifice.", player, obj, nil, actToChar)
		return
	}

	player.notify("The gods grant you a single gold coin for your sacrifice.")
	player.Gold++

	act("$n sacrifices $p to the gods.", player, obj, nil, actToRoom)
	player.removeItem(obj)
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

	if obj.ItemType == itemLight {
		m.Room.Light--
	}
	obj.WearLocation = wearNone

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

	if wearable.ItemType == itemLight && wearable.canWear(itemWearLight) {
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

		if wearable.Weight > m.currentStrength() {
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
