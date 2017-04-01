package mud

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

func doAffect(player *mob, argument string) {
	for _, af := range player.Affects {
		var duration string

		if af.duration < 30 {
			duration = "a few more seconds"
		} else if af.duration < 60 {
			duration = "less than a minute"
		} else if af.duration < 150 {
			duration = "a few minutes"
		} else if af.duration < 300 {
			duration = "a while"
		} else if af.duration < 600 {
			duration = "a long time"
		} else if af.duration > 900 {
			duration = "practifally forever"
		}

		player.notify("%s for %s.", af.affectType.Skill.Name, duration)
	}
}

func doAreas(player *mob, argument string) {
	for e := areaList.Front(); e != nil; e = e.Next() {
		a := e.Value.(*area)
		player.notify(a.Name)
	}
}

func doCompare(player *mob, argument string) {
	argument, arg1 := oneArgument(argument)
	argument, arg2 := oneArgument(argument)

	if arg1 == "" {
		player.notify("Compare what to what?")
		return
	}

	var obj1 *item
	var obj2 *item
	for _, i := range player.Inventory {
		if matchesSubject(i.Name, arg1) {
			obj1 = i
			break
		}
	}

	if obj1 == nil {
		player.notify("You don't have that.")
		return
	}

	if arg2 == "" {
		for _, i := range player.Inventory {
			if i.WearLocation == wearNone && player.canSeeItem(i) && i.ItemType == obj1.ItemType && (obj1.WearFlags&i.WearFlags & ^itemTake != 0) {
				obj2 = i
				break
			}
		}

		if obj2 == nil {
			player.notify("You don't have anything comparable.")
			return
		}
	} else {
		for _, i := range player.Inventory {
			if matchesSubject(i.Name, arg2) {
				obj2 = i
				break
			}
		}

		if obj2 == nil {
			player.notify("You don't have that.")
			return
		}
	}

	msg := ""
	value1 := 0
	value2 := 0

	if obj1 == obj2 {
		msg = "You compare $p to itself. It looks about the same."
	} else if obj1.ItemType != obj2.ItemType {
		msg = "You can't compare $p to $P."
	} else {
		switch obj1.ItemType {
		default:
			msg = "You can't compare $p to $P."
			break

		case itemArmor:
			value1 = obj1.Min
			value2 = obj2.Min
			break

		case itemWeapon:
			value1 = obj1.Min + obj1.Max
			value2 = obj2.Min + obj2.Max
			break
		}
	}

	if msg == "" {
		if value1 == value2 {
			msg = "$p and $P look about the same."
		} else if value1 > value2 {
			msg = "$p looks better than $P."
		} else {
			msg = "$p looks worse than $P."
		}
	}

	act(msg, player, obj1, obj2, actToChar)
	return
}

func doConsider(player *mob, argument string) {
	if len(argument) == 0 {
		player.notify("Consider killing whom?")
		return
	}

	argument, arg1 := oneArgument(argument)
	var victim *mob
	for _, mob := range player.Room.Mobs {
		if matchesSubject(mob.Name, arg1) {
			victim = mob
			break
		}
	}

	if victim == nil {
		player.notify("They aren't here.")
		return
	}

	diff := victim.Level - player.Level

	var msg string
	if diff <= -10 {
		msg = "You can kill $N naked and weaponless."
	} else if diff <= -5 {
		msg = "$N is no match for you."
	} else if diff <= -2 {
		msg = "$N looks like an easy kill."
	} else if diff <= 1 {
		msg = "The perfect match!"
	} else if diff <= 4 {
		msg = "$N says 'Do you feel lucky, punk?'"
	} else if diff <= 9 {
		msg = "$N laughs at you mercilessly."
	} else {
		msg = "Death will thank you for your gift."
	}

	act(msg, player, nil, victim, actToChar)
	return
}

func doExamine(player *mob, argument string) {
	if len(argument) == 0 {
		player.notify("Examine what?")
		return
	}

	argument, arg1 := oneArgument(argument)

	doLook(player, arg1)

	var obj *item
	for _, i := range player.Inventory {
		if matchesSubject(i.Name, arg1) {
			obj = i
			break
		}
	}

	if obj == nil {
		for _, i := range player.Room.Items {
			if matchesSubject(i.Name, arg1) {
				obj = i
				break
			}
		}
	}

	if obj != nil {
		switch obj.ItemType {
		case itemDrinkContainer:
		case itemContainer:
		case itemCorpseNPC:
		case itemCorpsePC:
			player.notify("When you look inside you see:")
			doLook(player, fmt.Sprintf("in %s", arg1))
			break
		default:
			break
		}
	}
}

func doHelp(player *mob, argument string) {
	if argument == "" {
		argument = "summary"
	}

	for e := helpList.Front(); e != nil; e = e.Next() {
		h := e.Value.(*help)
		if h.Level < player.getTrust() {
			continue
		}

		if matchesSubject(h.Keyword, argument) {
			if h.Level >= 0 && argument != "imotd" {
				player.notify(h.Keyword)
				player.notify("")
			}

			text := strings.TrimLeft(h.Text, ".")
			player.notify(text)
			return
		}
	}

	player.notify("No help on that word.")
}

func doInventory(player *mob, argument string) {
	player.notify("Inventory\n%s\n%s\n%s",
		"-----------------------------------",
		strings.Join(inventoryString(player), newline),
		"-----------------------------------",
	)
}

func doLook(player *mob, argument string) {
	if player.client == nil {
		return
	}

	if player.Status <= sleeping {
		player.notify("You can't see anything but stars.")
		return
	}

	if player.Status == sleeping {
		player.notify("You can't see anything; you're sleeping!")
		return
	}

	if hasBit(player.AffectedBy, affectBlind) {
		return
	}

	if !player.isNPC() && player.Room.isDark() {
		player.notify("It is pitch black...")
		showCharactersToPlayer(player.Room.Mobs, player)
		return
	}

	if len(argument) == 0 {
		// look
		player.notify(player.Room.Name)
		player.Room.showExits(player)
		player.notify(player.Room.Description)

		showItemsToPlayer(player.Room.Items, player)
		showCharactersToPlayer(player.Room.Mobs, player)
		return
	}

	if strings.HasPrefix(argument, "i") {
		// look in

		argument, _ := oneArgument(argument)
		if len(argument) == 0 {
			player.notify("Look in what?")
			return
		}

		var item *item
		for _, i := range player.Inventory {
			if matchesSubject(i.Name, argument) {
				item = i
				break
			}
		}

		if item == nil {
			for _, i := range player.Room.Items {
				if matchesSubject(i.Name, argument) {
					item = i
					break
				}
			}
		}

		if item == nil {
			player.notify("There is nothing like that here.")
			return
		}

		switch item.ItemType {
		default:
			player.notify("That is not a container.")
			break
		case itemContainer:
		case itemCorpseNPC:
		case itemCorpsePC:
			if item.isClosed() {
				player.notify("It is closed.")
				return
			}

			act("$p contains: ", player, item, nil, actToChar)
			showItemsToPlayer(player.Room.Items, player)
			break
		}
		return
	}

	var victim *mob
	for _, m := range player.Room.Mobs {
		if matchesSubject(m.Name, argument) {
			victim = m
			break
		}
	}

	if victim != nil {
		showCharacterToPlayer(victim, player)
	}

	for _, i := range player.Inventory {
		if matchesSubject(i.Name, argument) {
			if player.canSeeItem(i) {
				player.notify(i.Description)
				return
			}
		}
	}

	for _, i := range player.Equipped {
		if matchesSubject(i.Name, argument) {
			if player.canSeeItem(i) {
				player.notify(i.Description)
				return
			}
		}
	}

	for _, i := range player.Room.Items {
		if matchesSubject(i.Name, argument) {
			if player.canSeeItem(i) {
				player.notify(i.Description)
				return
			}
		}
	}

	var door string
	if strings.HasPrefix(argument, "n") || strings.HasPrefix(argument, "north") {
		door = "north"
	} else if strings.HasPrefix(argument, "e") || strings.HasPrefix(argument, "east") {
		door = "east"
	} else if strings.HasPrefix(argument, "s") || strings.HasPrefix(argument, "south") {
		door = "south"
	} else if strings.HasPrefix(argument, "w") || strings.HasPrefix(argument, "west") {
		door = "west"
	} else if strings.HasPrefix(argument, "u") || strings.HasPrefix(argument, "up") {
		door = "up"
	} else if strings.HasPrefix(argument, "d") || strings.HasPrefix(argument, "down") {
		door = "down"
	} else {
		player.notify("You do not see that here.")
		return
	}

	var exit *exit
	for _, e := range player.Room.Exits {
		if e.Dir == door {
			exit = e
			break
		}
	}

	if exit == nil {
		player.notify("Nothing special there.")
		return
	}

	if len(exit.Description) == 0 {
		player.notify(exit.Description)
	} else {
		player.notify("Nothing special there.")
	}

	if len(exit.Keyword) == 0 {
		if exit.isClosed() {
			act("The $d is closed.", player, nil, exit.Keyword, actToChar)
		} else if exit.hasDoor() {
			act("The $d is open.", player, nil, exit.Keyword, actToChar)
		}
	}
	return
}

func doScan(player *mob, argument string) {
	room := getRoom(player.Room.ID)
	for _, x := range room.Exits {
		player.notify("[%s]", x.Dir)

		if len(x.Room.Mobs) > 0 {
			mobs := x.Room.Mobs
			for _, m := range mobs {
				player.notify("    %s", m.Name)
			}
		} else {
			player.notify("    %s(nothing)%s", blue, reset)
		}
	}
}

func doScore(player *mob, argument string) {
	const (
		width int = 50
	)

	username := player.Name
	id := fmt.Sprintf("Level %d %s %s", player.Level, player.Race.Name, player.Job.Name)
	spaces := width - len(username) - len(id)

	title := fmt.Sprintf("%s%s%s", username, strings.Repeat(" ", spaces), id)

	strength := fmt.Sprintf("%s%s%d%s%s%s%d", "Strength", strings.Repeat(" ", 8), player.Attributes.Strength, strings.Repeat(" ", 11), "Experience", strings.Repeat(" ", 11-len(strconv.Itoa(player.Exp))), player.Exp)
	wisdom := fmt.Sprintf("%s%s%d%s%s%s%d", "Wisdom", strings.Repeat(" ", 10), player.Attributes.Wisdom, strings.Repeat(" ", 11), "TNL", strings.Repeat(" ", 18-len(strconv.Itoa(player.TNL()))), player.TNL())
	intel := fmt.Sprintf("%s%s%d%s%s%s%d", "Intelligence", strings.Repeat(" ", 4), player.Attributes.Intelligence, strings.Repeat(" ", 11), "Alignment", strings.Repeat(" ", 12-len(strconv.Itoa(player.Alignment))), player.Alignment)
	dexterity := fmt.Sprintf("%s%s%d%s%s%s%d", "Dexterity", strings.Repeat(" ", 7), player.Attributes.Dexterity, strings.Repeat(" ", 11), "Practices", strings.Repeat(" ", 12-len(strconv.Itoa(int(player.Practices)))), player.Practices)
	constitution := fmt.Sprintf("%s%s%d", "Constitution", strings.Repeat(" ", 4), player.Attributes.Constitution)
	charisma := fmt.Sprintf("%s%s%d", "Charisma", strings.Repeat(" ", 8), player.Attributes.Charisma)
	player.notify(
		fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n",
			strings.Repeat("=", width),
			title,
			strings.Repeat("-", width),
			strength,
			wisdom,
			intel,
			dexterity,
			constitution,
			charisma,
			strings.Repeat("=", width),
		),
	)
}

func doSkills(player *mob, argument string) {
	const (
		width int = 40
	)
	player.notify("Skill                              Level\n")
	player.notify("----------------------------------------\n")
	for _, skill := range player.Skills {
		name := skill.Skill.Name
		level := skill.Level
		spaces := width - len(name) - len(strconv.Itoa(int(level)))
		player.notify(fmt.Sprintf("%s%s%d", name, strings.Repeat(" ", spaces), level))
	}
	player.notify("----------------------------------------\n")
}

func doWhere(player *mob, argument string) {
	if len(argument) == 0 {
		player.notify("Players near you:")
		found := false
		for e := mobList.Front(); e != nil; e = e.Next() {
			m := e.Value.(*mob)
			if m.client != nil && !m.isNPC() && m.Room != nil && m.Room.Area.ID == player.Room.Area.ID && player.canSee(m) {
				found = true
				player.notify("%-28s %s", m.Name, m.Room.Name)
			}
		}

		if !found {
			player.notify("None")
		}
	} else {
		_, arg1 := oneArgument(argument)

		found := false
		for e := mobList.Front(); e != nil; e = e.Next() {
			m := e.Value.(*mob)
			if m.client != nil && m.Room.Area.ID == player.Room.Area.ID && !hasBit(m.AffectedBy, affectHide) && !hasBit(m.AffectedBy, affectSneak) && player.canSee(m) && matchesSubject(m.Name, arg1) {
				found = true
				player.notify("%-28s %s", pers(m, player), m.Room.Name)
			}
		}

		if !found {
			act("You didn't find any $T.", player, nil, arg1, actToChar)
		}
	}
}

func doWho(player *mob, argument string) {
	var (
		minLevel        = 0
		maxLevel        = 99
		classRestrict   = false
		classToRestrict *job
	)

	argument, arg1 := oneArgument(argument)
	argument, arg2 := oneArgument(argument)

	if arg2 != "" {
		num1, err1 := strconv.Atoi(arg1)
		num2, err2 := strconv.Atoi(arg2)

		if err1 == nil && err2 == nil {
			minLevel = num1
			maxLevel = num2
		} else {
			// must be a string

			if len(arg1) < 3 {
				player.notify("Classes must be longer than that.")
				return
			}

			classRestrict = true
			for e := jobList.Front(); e != nil; e = e.Next() {
				job := e.Value.(*job)
				if matchesSubject(job.Abbr, arg1) {
					classToRestrict = job
					break
				}
			}

			if classToRestrict == nil {
				player.notify("That is not a class.")
				return
			}
		}
	}

	nMatch := 0

	var buf bytes.Buffer

	for e := mobList.Front(); e != nil; e = e.Next() {
		mob := e.Value.(*mob)

		if mob.client == nil || !player.canSee(mob) {
			continue
		}

		if mob.Level < minLevel || mob.Level > maxLevel || (classRestrict && mob.Job != classToRestrict) {
			continue
		}

		nMatch++
		job := mob.Job.Name
		race := mob.Race.Name
		switch mob.Level {
		default:
			break
		case 99:
			job = "GOD"
			break
		case 98:
			job = "SUP"
			break
		case 97:
			job = "DEI"
			break
		case 96:
			job = "ANG"
			break
		}

		buf.Write([]byte(fmt.Sprintf("[%2d %8s %8s] %s %s%s", mob.Level, race, job, mob.Name, mob.Title, newline)))
	}

	suffix := "s"
	if nMatch == 1 {
		suffix = ""
	}
	player.notify("%d player%s.", nMatch, suffix)
	player.notify(buf.String())
}

/* Helpers */

func exitsString(exits []*exit) string {
	var output string
	for _, e := range exits {
		output = fmt.Sprintf("%s%s ", output, string(e.Dir))
	}

	return fmt.Sprintf("[%s]%s", strings.Trim(output, " "), newline)
}

func itemsString(items []*item) string {
	var output string

	for _, i := range items {
		output = fmt.Sprintf("%s is here.%s%s", i.Name, newline, output)
	}
	return output
}

func mobsString(mobs []*mob, player *mob) string {
	var output string
	output = ""
	for _, m := range mobs {
		if m != player {
			if player.canSee(m) {
				output = fmt.Sprintf("%s is here.%s%s", m.Name, newline, output)
			} else {
				output = fmt.Sprintf("You see glowing red eyes watching YOU!%s%s", newline, output)
			}
		}
	}

	return output
}

func inventoryString(m *mob) []string {
	inventory := make(map[string]int)

	for _, i := range m.Inventory {
		if _, ok := inventory[i.Name]; ok {
			inventory[i.Name]++
		} else {
			inventory[i.Name] = 1
		}
	}

	var items []string
	for name, qty := range inventory {
		if qty > 1 {
			items = append(items, fmt.Sprintf("(%d) %s", qty, name))
		} else {
			items = append(items, name)
		}
	}

	return items
}

func equippedString(m *mob) []string {
	var lines []string
	lines = append(lines, fmt.Sprintf("<light>     %s", m.equipped(wearLight)))
	lines = append(lines, fmt.Sprintf("<finger>    %s", m.equipped(wearFingerLeft)))
	lines = append(lines, fmt.Sprintf("<finger>    %s", m.equipped(wearFingerRight)))
	lines = append(lines, fmt.Sprintf("<neck>      %s", m.equipped(wearNeck1)))
	lines = append(lines, fmt.Sprintf("<neck>      %s", m.equipped(wearNeck2)))
	lines = append(lines, fmt.Sprintf("<head>      %s", m.equipped(wearHead)))
	lines = append(lines, fmt.Sprintf("<legs>      %s", m.equipped(wearLegs)))
	lines = append(lines, fmt.Sprintf("<feet>      %s", m.equipped(wearFeet)))
	lines = append(lines, fmt.Sprintf("<hands>     %s", m.equipped(wearHands)))
	lines = append(lines, fmt.Sprintf("<arms>      %s", m.equipped(wearArms)))
	lines = append(lines, fmt.Sprintf("<shield>    %s", m.equipped(wearShield)))
	lines = append(lines, fmt.Sprintf("<body>      %s", m.equipped(wearBody)))
	lines = append(lines, fmt.Sprintf("<waist>     %s", m.equipped(wearWaist)))
	lines = append(lines, fmt.Sprintf("<wrist>     %s", m.equipped(wearWristLeft)))
	lines = append(lines, fmt.Sprintf("<wrist>     %s", m.equipped(wearWristRight)))
	lines = append(lines, fmt.Sprintf("<wield>     %s", m.equipped(wearWield)))
	lines = append(lines, fmt.Sprintf("<held>      %s", m.equipped(wearHold)))

	return lines
}
