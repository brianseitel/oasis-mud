package mud

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"bytes"

	"github.com/brianseitel/oasis-mud/helpers"
)

type action struct {
	mob   *mob
	rooms []*room
	conn  *connection
	args  []string
}

func newAction(m *mob, c *connection, i string) {
	newActionWithInput(&action{mob: m, conn: c, args: strings.Split(i, " ")})
}

func newActionWithInput(a *action) error {

	switch a.getCommand() {
	case cLook:
		a.look()
		return nil
	case cNorth:
		a.move("north")
		return nil
	case cSouth:
		a.move("south")
		return nil
	case cEast:
		a.move("east")
		return nil
	case cWest:
		a.move("west")
		return nil
	case cUp:
		a.move("up")
		return nil
	case cDown:
		a.move("down")
		return nil
	case cQuit:
		a.quit()
		return errors.New("Done")
	case cDrop:
		a.drop()
		return nil
	case cGet:
		a.get()
		return nil
	case cInventory:
		a.inventory()
		return nil
	case cScore:
		a.score()
		return nil
	case cKill:
		a.kill()
		return nil
	case cFlee:
		a.flee()
		return nil
	case cWear:
		a.wear()
		return nil
	case cRemove:
		a.remove()
		return nil
	case cEquipment:
		a.equipment()
		return nil
	case cScan:
		a.scan()
		return nil
	case cRecall:
		a.recall()
		return nil
	case cSkill:
		a.skills()
		return nil
	case cTrip:
		a.trip()
		return nil
	case cTrain:
		a.train()
		return nil
	case cCast:
		a.cast()
		return nil
	case cAffect:
		a.affect()
		return nil
	case cChat:
		a.mob.chatDefault(a.args)
		return nil
	case cChatAuction:
		a.mob.chatAuction(a.args)
		return nil
	case cChatMusic:
		a.mob.chatMusic(a.args)
		return nil
	case cChatQuestion:
		a.mob.chatQuestion(a.args)
		return nil
	case cChatAnswer:
		a.mob.chatAnswer(a.args)
		return nil
	case cChatImmtalk:
		a.mob.chatImmtalk(a.args)
		return nil
	case cSay:
		a.mob.say(a.args)
		return nil
	case cTell:
		a.mob.tell(a.args)
		return nil
	case cReply:
		a.mob.reply(a.args)
		return nil
	case cPut:
		a.put()
		return nil
	case cGive:
		a.give()
		return nil
	case cConsider:
		a.consider()
		return nil
	case cSneak:
		a.sneak()
		return nil
	default:
		a.conn.SendString("Eh?")
	}
	return nil
}

func (a *action) getCommand() command {
	for _, c := range commands {
		if isCommand(c, a.args[0]) == true {
			return c
		}
	}

	return cNoop
}

func isCommand(c command, p string) bool {
	return strings.HasPrefix(string(c), p)
}

func (a *action) skills() {
	const (
		width int = 40
	)
	a.conn.SendString("Skill                              Level\n")
	a.conn.SendString("----------------------------------------\n")
	for _, skill := range a.mob.Skills {
		name := skill.Skill.Name
		level := skill.Level
		spaces := width - len(name) - len(strconv.Itoa(int(level)))
		a.conn.SendString(fmt.Sprintf("%s%s%d", name, strings.Repeat(" ", spaces), level))
	}
	a.conn.SendString("----------------------------------------\n")
}

func (a *action) look() {
	if a.mob == nil {
		return
	}

	if len(a.args) == 1 {
		a.conn.SendString(
			fmt.Sprintf(
				"%s [ID: %d]\n%s\n%s%s%s",
				a.mob.Room.Name,
				a.mob.Room.ID,
				a.mob.Room.Description,
				exitsString(a.mob.Room.Exits),
				itemsString(a.mob.Room.Items),
				mobsString(a.mob.Room.Mobs, a.mob),
			),
		)
	} else {
		for _, mob := range a.mob.Room.Mobs {
			if a.matchesSubject(mob.Identifiers) {
				a.conn.SendString(fmt.Sprintf("You look at %s.", mob.Name))
				a.conn.SendString(helpers.WordWrap(mob.Description, 50))
				return
			}
		}

		for _, item := range a.mob.Room.Items {
			if a.matchesSubject(item.Name) {
				a.conn.SendString(fmt.Sprintf("You look at %s.", item.Name))
				a.conn.SendString(helpers.WordWrap(item.Description, 50))
				return
			}
		}
		a.conn.SendString("Look at what?")
	}
}

func (a *action) inventory() {
	a.conn.SendString(
		fmt.Sprintf("Inventory\n%s\n%s\n%s",
			"-----------------------------------",
			strings.Join(inventoryString(a.mob), helpers.Newline),
			"-----------------------------------",
		),
	)
}

func (a *action) equipment() {
	a.conn.SendString(
		fmt.Sprintf("Equipment\n%s\n%s\n%s",
			"-----------------------------------",
			strings.Join(equippedString(a.mob), helpers.Newline),
			"-----------------------------------",
		),
	)
}

func (a *action) score() {
	const (
		width int = 50
	)

	username := a.mob.Name
	id := fmt.Sprintf("Level %d %s %s", a.mob.Level, a.mob.Race.Name, a.mob.Job.Name)
	spaces := width - len(username) - len(id)

	title := fmt.Sprintf("%s%s%s", username, strings.Repeat(" ", spaces), id)

	strength := fmt.Sprintf("%s%s%d%s%s%s%d", "Strength", strings.Repeat(" ", 8), a.mob.Attributes.Strength, strings.Repeat(" ", 11), "Experience", strings.Repeat(" ", 11-len(strconv.Itoa(a.mob.Exp))), a.mob.Exp)
	wisdom := fmt.Sprintf("%s%s%d%s%s%s%d", "Wisdom", strings.Repeat(" ", 10), a.mob.Attributes.Wisdom, strings.Repeat(" ", 11), "TNL", strings.Repeat(" ", 18-len(strconv.Itoa(a.mob.TNL()))), a.mob.TNL())
	intel := fmt.Sprintf("%s%s%d%s%s%s%d", "Intelligence", strings.Repeat(" ", 4), a.mob.Attributes.Intelligence, strings.Repeat(" ", 11), "Alignment", strings.Repeat(" ", 12-len(strconv.Itoa(a.mob.Alignment))), a.mob.Alignment)
	dexterity := fmt.Sprintf("%s%s%d%s%s%s%d", "Dexterity", strings.Repeat(" ", 7), a.mob.Attributes.Dexterity, strings.Repeat(" ", 11), "Practices", strings.Repeat(" ", 12-len(strconv.Itoa(int(a.mob.Practices)))), a.mob.Practices)
	constitution := fmt.Sprintf("%s%s%d", "Constitution", strings.Repeat(" ", 4), a.mob.Attributes.Constitution)
	charisma := fmt.Sprintf("%s%s%d", "Charisma", strings.Repeat(" ", 8), a.mob.Attributes.Charisma)
	a.conn.SendString(
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

func exitsString(exits []*exit) string {
	var output string
	for _, e := range exits {
		output = fmt.Sprintf("%s%s ", output, string(e.Dir))
	}

	return fmt.Sprintf("[%s]%s", strings.Trim(output, " "), helpers.Newline)
}

func itemsString(items []*item) string {
	var output string

	for _, i := range items {
		output = fmt.Sprintf("%s is here.%s%s", i.Name, helpers.Newline, output)
	}
	return output
}

func mobsString(mobs []*mob, player *mob) string {
	var output string
	output = ""
	for _, m := range mobs {
		if m != player {
			output = fmt.Sprintf("%s is here.%s%s", m.Name, helpers.Newline, output)
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

func (a *action) affect() {
	for _, af := range a.mob.Affects {
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

		a.mob.notify("%s for %s.", af.affectType.Skill.Name, duration)
	}
}

func (a *action) cast() {
	var victim *mob
	var player *mob
	var mana int
	var spell *mobSkill

	player = a.mob

	if len(a.args) < 2 {
		a.mob.notify("Cast which what where?\r\n")
		return
	}

	spell = player.skill(a.args[1])
	if spell == nil {
		a.mob.notify("You can't do that. \r\n")
		return
	}

	mana = 0
	if !player.isNPC() {
		mana = helpers.Max(spell.Skill.MinMana, 100/(2+player.Level))
	}

	// Find targets
	victim = nil

	switch spell.Skill.Target {
	case "ignore":
		break

	case "offensive":
		if len(a.args) < 3 {
			a.mob.notify("Cast the spell on whom?\r\n")
			return
		}

		arg := a.args[2]
		for _, mob := range player.Room.Mobs {
			if strings.HasPrefix(mob.Name, arg) {
				victim = mob
				break
			}
		}

		if victim == nil {
			a.mob.notify("They aren't here.\r\n")
			return
		}

		if victim == player {
			a.mob.notify("You can't do that to yourself.\r\n")
			return
		}
		break

	case "defensive":
		if len(a.args) < 3 {
			victim = player
		} else {
			arg := a.args[2]
			for _, mob := range player.Room.Mobs {
				if strings.HasPrefix(mob.Name, arg) {
					victim = mob
					break
				}
			}
		}
		break

	case "self":
		if len(a.args) > 2 {
			a.mob.notify("You cannot cast this spell on another.\r\n")
			return
		}
		victim = player
		break

	case "object":
		break

	default:
		fmt.Printf("cast: bad target for %s\r\n", spell.Skill.Name)
	}

	if !player.isNPC() && player.Mana < mana {
		a.mob.notify("You don't have enough mana.")
		return
	}

	if !player.isNPC() && dice().Intn(100) > int(spell.Level) {
		player.notify("You lost your concentration!")
		player.Mana -= mana / 2
	} else {
		player.Mana -= mana
		doSpell(spell, player, victim)
	}

	return
}

func (a *action) consider() {
	player := a.mob

	if len(a.args) <= 1 {
		player.notify("Consider killing whom?")
		return
	}

	arg1 := a.args[1]
	var victim *mob
	for _, mob := range player.Room.Mobs {
		if helpers.MatchesSubject(mob.Name, arg1) {
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

func (a *action) drop() {
	player := a.mob
	if len(a.args) <= 1 {
		player.notify("Drop what?")
		return
	}

	arg1 := a.args[1]

	num, err := strconv.Atoi(arg1)
	isNumber := err == nil

	if isNumber {
		amount := uint(num)
		if len(a.args) < 2 || amount <= 0 || !strings.HasPrefix(a.args[2], "gold") {
			player.notify("Sorry, you can't do that.")
			return
		}

		if player.Gold < uint(amount) {
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
			if helpers.MatchesSubject(i.Name, arg1) {
				item = i
				break
			}
		}

		if item == nil {
			player.notify("You don't have that item.")
			return
		}

		// TODO: canDropObj

		for j, i := range player.Inventory {
			if i == item {
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

		fmt.Println(name)
		for j := 0; j < len(player.Inventory); j++ {
			item := player.Inventory[j]
			if arg1 == "all" || helpers.MatchesSubject(item.Name, name) {
				found = true

				player.Inventory = append(player.Inventory[0:j], player.Inventory[j+1:]...)
				j--
				player.Room.Items = append(player.Room.Items, item)
				act("$n drops $p.", player, item, nil, actToRoom)
				act("You drop $p.", player, item, nil, actToChar)
				// player.notify("You drop %s.", item.Name)
				// player.Room.notify(fmt.Sprintf("%s drops %s.", player.Name, item.Name), player)
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

func (a *action) flee() {
	if a.mob.Status != fighting {
		a.conn.SendString("You can't flee if you're not fighting, fool.")
		return
	}

	roll := dice().Intn(36 - a.mob.Attributes.Dexterity) // higher the dexterity, better the chance of fleeing successfully
	if roll == 1 {
		a.mob.Fight.Mob1.Status = standing
		a.mob.Fight.Mob2.Status = standing

		a.mob.wander()
		a.conn.SendString("You flee!")
	} else {
		a.conn.SendString("You tried to flee, but failed!")
	}
}

func (a *action) get() {

	player := a.mob

	if len(a.args) <= 1 {
		player.notify("Get what?")
		return
	}
	arg1, args := a.args[1], a.args[2:]

	if len(a.args) == 2 {
		if arg1 != "all" && !strings.HasPrefix(arg1, "all.") {
			// get obj
			for _, i := range player.Room.Items {
				if helpers.MatchesSubject(i.Name, arg1) {
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
					if helpers.MatchesSubject(i.Name, name) || len(name) == 0 {
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
		arg2 := args[0]

		if arg2 == "all" || strings.HasPrefix(arg2, "all.") {
			player.notify("You can't do that.")
			return
		}

		var container *item
		for _, i := range player.Room.Items {
			if strings.HasPrefix(i.Name, arg2) {
				container = i
				break
			}
		}

		if container == nil {
			// try from inventory
			for _, i := range player.Inventory {
				if strings.HasPrefix(i.Name, arg1) {
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

		if helpers.HasBit(uint(container.Value), containerClosed) {
			player.notify("The %s is closed.", container.Name)
			return
		}

		if arg1 != "all" && !strings.HasPrefix(arg1, "all.") {
			// get obj container
			for _, i := range container.container {
				if helpers.MatchesSubject(i.Name, arg1) {
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
				if helpers.MatchesSubject(i.Name, name) || len(name) == 0 {
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

func (a *action) give() {
	player := a.mob
	if len(a.args) <= 2 {
		player.notify("Give what to whom?")
		return
	}

	arg1, arg2 := a.args[1], a.args[2]

	var victim *mob
	for _, mob := range player.Room.Mobs {
		if helpers.MatchesSubject(mob.Name, arg2) {
			victim = mob
			break
		}
	}

	num, err := strconv.Atoi(arg1)
	isNumeric := err == nil

	var amount uint
	if isNumeric {
		if num <= 0 || !strings.HasPrefix(arg2, "coin") {
			player.notify("You can't do that.")
			return
		}
		amount = uint(num)

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
		if helpers.MatchesSubject(i.Name, arg1) {
			item = i
			break
		}
	}

	if item == nil {
		player.notify("You do not have that item.")
		return
	}

	// TODO:
	// if item.WearLocation != wearNone {
	// 	player.notify("You must remove it first.")
	// 	return
	// }

	// TODO:
	// if victim.Room.ID != player.Room.ID {
	// 	player.notify("They are not here.")
	// 	return
	// }

	// TODO:
	// if !item.canDrop(player) {
	// 	player.notify("You can't let go of it.")
	// 	return
	// }

	if victim.Carrying+1 > victim.CarryMax {
		player.notify("%s has their hands full.", victim.Name)
		return
	}

	if victim.CarryWeight+item.Weight > victim.CarryWeightMax {
		player.notify("%s can't carry that much weight.", victim.Name)
		return
	}

	// TODO:
	// if item.canSee(victim) {
	// 	player.notify("%s can't see it.", victim.Name)
	// }

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

func (a *action) kill() {
	for _, m := range a.mob.Room.Mobs {
		if a.matchesSubject(m.Name) {
			newFight(a.mob, m)
			return
		}
	}

	a.mob.notify("You can't find them.")
}

func (a *action) move(d string) {
	if a.mob.Status != standing {
		switch a.mob.Status {
		case fighting:
			a.conn.SendString("You can't move while fighting!")
			break
		}
		return
	}

	for _, e := range a.mob.Room.Exits {
		if e.Dir == d {
			a.mob.move(e)
			newAction(a.mob, a.conn, "look")
			return
		}
	}
	a.conn.SendString("Alas, you cannot go that way.")
}

func (a *action) put() {

	player := a.mob
	if len(a.args) <= 2 {
		player.notify("Put what in what?")
		return
	}

	arg1, arg2 := a.args[1], a.args[2]

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

	if helpers.HasBit(uint(container.Value), containerClosed) {
		player.notify("The %s is closed.", container.Name)
		return
	}

	if arg1 != "all" && !strings.HasPrefix(arg1, "all.") {
		// put obj container
		var item *item
		for _, i := range player.Inventory {
			if helpers.MatchesSubject(i.Name, arg1) {
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

		// TODO
		// if !player.canDropObj(item) {
		// 	player.notify("You can't let go of it.")
		// 	return
		// }

		if item.Weight+container.Weight > uint(container.Value) {
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
			if (arg1 == "all" || strings.HasPrefix(item.Name, name)) && item.WearLocation == wearNone && item != container && item.Weight+container.Weight > uint(container.Value) {
				player.Inventory = append(player.Inventory[0:j], player.Inventory[j+1:]...)
				container.container = append(container.container, item)
				player.notify("You put %s in %s.", item.Name, container.Name)
				player.Room.notify(fmt.Sprintf("%s puts %s in %s.", player.Name, item.Name, container.Name), player)
			}
		}
	}
}

func (a *action) quit() {
	if a.mob.Status == fighting {
		a.conn.SendString("You can't quit now. You're fighting!")
	} else {
		a.conn.SendString("Seeya!")
		a.conn.end()
	}
}

func (a *action) recall() {
	if len(a.args) == 1 {
		room := getRoom(a.mob.RecallRoomID)
		a.mob.Room = room
		a.look()
		return
	}

	if a.args[1] == "set" {
		a.mob.RecallRoomID = a.mob.Room.ID
		a.conn.SendString("Recall set!")
		return
	}

	a.conn.SendString("Recall what?")
}

func (a *action) remove() {
	for j, item := range a.mob.Equipped {
		if a.matchesSubject(item.Name) {
			a.mob.Equipped, a.mob.Inventory = transferItem(j, a.mob.Equipped, a.mob.Inventory)
			a.mob.notify("You remove %s.", item.Name)
			return
		}
	}

	a.mob.notify("You aren't wearing that.")
}

func (a *action) sneak() {
	sneak := a.mob.skill("sneak")
	if sneak == nil {
		a.mob.notify("You don't know how to sneak.")
		return
	}
	for _, affect := range a.mob.Affects {
		if affect.affectType == sneak {
			a.mob.removeAffect(affect)
			break
		}
	}

	var af affect
	af.affectType = sneak
	af.duration = uint(helpers.Max(5, a.mob.Level))
	af.modifier = 0
	af.location = applyNone
	af.bitVector = affectSneak

	a.mob.addAffect(&af)
	a.mob.notify("You are now sneaking.")
}

func (a *action) scan() {
	room := getRoom(a.mob.Room.ID)
	for _, x := range room.Exits {
		a.mob.notify("[%s]", x.Dir)

		if len(x.Room.Mobs) > 0 {
			mobs := x.Room.Mobs
			for _, m := range mobs {
				a.mob.notify("    %s", m.Name)
			}
		} else {
			a.mob.notify("    %s(nothing)%s", helpers.Blue, helpers.Reset)
		}
	}
}

func (a *action) train() {
	player := a.mob
	if player.isNPC() {
		return
	}

	var trainer *mob
	for _, mob := range player.Room.Mobs {
		if mob.isTrainer() {
			trainer = mob
			break
		}
	}

	if trainer == nil {
		player.notify("You can't do that here.")
		return
	}

	if len(a.args) == 1 {
		player.notify("You have %d practice sessions.", player.Practices)
		return
	}

	var cost uint

	costmap := []uint{5, 6, 7, 9, 12, 13, 15}

	var playerAbility int
	var playerOutput string

	if strings.HasPrefix(a.args[1], "str") {
		playerAbility = player.Attributes.Strength
		playerOutput = "strength"
	} else if strings.HasPrefix(a.args[1], "int") {
		playerAbility = player.Attributes.Intelligence
		playerOutput = "intelligence"
	} else if strings.HasPrefix(a.args[1], "wis") {
		playerAbility = player.Attributes.Wisdom
		playerOutput = "wisdom"
	} else if strings.HasPrefix(a.args[1], "dex") {
		playerAbility = player.Attributes.Dexterity
		playerOutput = "dexterity"
	} else if strings.HasPrefix(a.args[1], "cha") {
		playerAbility = player.Attributes.Charisma
		playerOutput = "charisma"
	} else if strings.HasPrefix(a.args[1], "con") {
		playerAbility = player.Attributes.Constitution
		playerOutput = "constitution"
	} else {
		var buf bytes.Buffer

		buf.WriteString("You can train:\r\n")
		if player.Attributes.Strength < 18 {
			buf.WriteString(fmt.Sprintf("Strength      %d\r\n", costmap[player.Attributes.Strength-12]))
		}
		if player.Attributes.Intelligence < 18 {
			buf.WriteString(fmt.Sprintf("Intelligence  %d\r\n", costmap[player.Attributes.Intelligence-12]))
		}
		if player.Attributes.Wisdom < 18 {
			buf.WriteString(fmt.Sprintf("Wisdom        %d\r\n", costmap[player.Attributes.Wisdom-12]))
		}
		if player.Attributes.Dexterity < 18 {
			buf.WriteString(fmt.Sprintf("Dexterity     %d\r\n", costmap[player.Attributes.Dexterity-12]))
		}
		if player.Attributes.Charisma < 18 {
			buf.WriteString(fmt.Sprintf("Charisma      %d\r\n", costmap[player.Attributes.Charisma-12]))
		}
		if player.Attributes.Constitution < 18 {
			buf.WriteString(fmt.Sprintf("Constitution  %d\r\n", costmap[player.Attributes.Constitution-12]))
		}

		message := buf.String()
		if !strings.HasSuffix(message, ":") {
			buf.WriteString(".\r\n")
			player.notify(buf.String())
		} else {
			player.notify("You have nothing left to train, you badass!")
		}

		return
	}

	cost = costmap[playerAbility-12]
	if playerAbility >= 18 {
		player.notify("Your %s is already at maximum.", playerOutput)
		return
	}

	if cost > player.Practices {
		player.notify("You don't have enough practices.")
		return
	}

	player.Practices -= cost
	switch playerOutput {
	case "strength":
		player.Attributes.Strength++
		break
	case "intelligence":
		player.Attributes.Intelligence++
		break
	case "wisdom":
		player.Attributes.Wisdom++
		break
	case "dexterity":
		player.Attributes.Dexterity++
		break
	case "charisma":
		player.Attributes.Charisma++
		break
	case "constitution":
		player.Attributes.Constitution++
		break
	}

	player.notify("Your %s increases for %d practice points!", playerOutput, cost)
	return
}

func (a *action) trip() {

	if a.mob.skill("trip") == nil {
		a.mob.notify("You don't know how to do this.")
		return
	}

	if len(a.args) > 1 {
		for _, m := range a.mob.Room.Mobs {
			if a.matchesSubject(m.Identifiers) {
				newFight(a.mob, m)
				break
			}
		}

		if a.mob.Fight == nil {
			a.mob.notify("Trip who?")
			return
		}
	}
	a.mob.trip()
}

func (a *action) wear() {
	var wearable *item
	for _, item := range a.mob.Inventory {
		if a.matchesSubject(item.Name) {
			wearable = item
			break
		}
	}

	if wearable == nil {
		a.mob.notify("You can't find that.")
		return
	}

	if a.mob.Level < wearable.Level {
		a.mob.notify("You must be level %d to wear this.", wearable.Level)
		return
	}

	if wearable.ItemType == itemLight {
		a.mob.notify("You light up %s and hold it.", wearable.Name)
		a.mob.Room.notify(fmt.Sprintf("%s lights up %s and holds it.", a.mob.Name, wearable.Name), a.mob)
		a.mob.equipItem(wearable, wearLight)
		return
	}

	if wearable.canWear(itemWearFinger) {
		if a.mob.equippedItem(wearFingerLeft) != nil && a.mob.equippedItem(wearFingerRight) != nil {
			return
		}

		if a.mob.equippedItem(wearFingerLeft) == nil {
			a.mob.notify("You wear %s on your left finger.", wearable.Name)
			a.mob.Room.notify(fmt.Sprintf("%s wears %s on their left finger.", a.mob.Name, wearable.Name), a.mob)
			a.mob.equipItem(wearable, wearFingerLeft)
			return
		}

		if a.mob.equippedItem(wearFingerRight) == nil {
			a.mob.notify("You wear %s on your right finger.", wearable.Name)
			a.mob.Room.notify(fmt.Sprintf("%s wears %s on their right finger.", a.mob.Name, wearable.Name), a.mob)
			a.mob.equipItem(wearable, wearFingerRight)
			return
		}

		a.mob.notify("You already wear two rings!")
		return
	}

	if wearable.canWear(itemWearNeck) {
		if a.mob.equippedItem(wearNeck1) != nil && a.mob.equippedItem(wearNeck2) != nil {
			return
		}

		if a.mob.equippedItem(wearNeck1) == nil {
			a.mob.notify("You wear %s on your neck.", wearable.Name)
			a.mob.Room.notify(fmt.Sprintf("%s wears %s on their neck.", a.mob.Name, wearable.Name), a.mob)
			a.mob.equipItem(wearable, wearNeck1)
			return
		}

		if a.mob.equippedItem(wearNeck2) == nil {
			a.mob.notify("You wear %s on your neck.", wearable.Name)
			a.mob.Room.notify(fmt.Sprintf("%s wears %s on their neck.", a.mob.Name, wearable.Name), a.mob)
			a.mob.equipItem(wearable, wearNeck2)
			return
		}

		a.mob.notify("You already wear two neck items!")
		return
	}
	if wearable.canWear(itemWearWrist) {
		if a.mob.equippedItem(wearWristLeft) != nil && a.mob.equippedItem(wearWristRight) != nil {
			return
		}

		if a.mob.equippedItem(wearWristLeft) == nil {
			a.mob.notify("You wear %s on your left wrist.", wearable.Name)
			a.mob.Room.notify(fmt.Sprintf("%s wears %s on their left wrist.", a.mob.Name, wearable.Name), a.mob)
			a.mob.equipItem(wearable, wearWristLeft)
			return
		}

		if a.mob.equippedItem(wearWristRight) == nil {
			a.mob.notify("You wear %s on your right wrist.", wearable.Name)
			a.mob.Room.notify(fmt.Sprintf("%s wears %s on their right wrist.", a.mob.Name, wearable.Name), a.mob)
			a.mob.equipItem(wearable, wearWristRight)
			return
		}

		a.mob.notify("You already wear two wrist items!")
		return
	}

	if wearable.canWear(itemWearBody) {
		a.mob.notify("You wear %s on your body.", wearable.Name)
		a.mob.Room.notify(fmt.Sprintf("%s wears %s on their body.", a.mob.Name, wearable.Name), a.mob)
		a.mob.equipItem(wearable, wearBody)
		return
	}

	if wearable.canWear(itemWearHead) {
		a.mob.notify("You wear %s on your head.", wearable.Name)
		a.mob.Room.notify(fmt.Sprintf("%s wears %s on their head.", a.mob.Name, wearable.Name), a.mob)
		a.mob.equipItem(wearable, wearHead)
		return
	}

	if wearable.canWear(itemWearLegs) {
		a.mob.notify("You wear %s on your legs.", wearable.Name)
		a.mob.Room.notify(fmt.Sprintf("%s wears %s on their legs.", a.mob.Name, wearable.Name), a.mob)
		a.mob.equipItem(wearable, wearLegs)
		return
	}

	if wearable.canWear(itemWearFeet) {
		a.mob.notify("You wear %s on your feet.", wearable.Name)
		a.mob.Room.notify(fmt.Sprintf("%s wears %s on their feet.", a.mob.Name, wearable.Name), a.mob)
		a.mob.equipItem(wearable, wearFeet)
		return
	}

	if wearable.canWear(itemWearHands) {
		a.mob.notify("You wear %s on your hands.", wearable.Name)
		a.mob.Room.notify(fmt.Sprintf("%s wears %s on their hands.", a.mob.Name, wearable.Name), a.mob)
		a.mob.equipItem(wearable, wearHands)
		return
	}

	if wearable.canWear(itemWearWaist) {
		a.mob.notify("You wear %s on your waist.", wearable.Name)
		a.mob.Room.notify(fmt.Sprintf("%s wears %s on their waist.", a.mob.Name, wearable.Name), a.mob)
		a.mob.equipItem(wearable, wearWaist)
		return
	}

	if wearable.canWear(itemWearShield) {
		a.mob.notify("You wear %s as your shield.", wearable.Name)
		a.mob.Room.notify(fmt.Sprintf("%s wears %s as their shield.", a.mob.Name, wearable.Name), a.mob)
		a.mob.equipItem(wearable, wearShield)
		return
	}

	if wearable.canWear(itemWearHold) {
		a.mob.notify("You hold %s.", wearable.Name)
		a.mob.Room.notify(fmt.Sprintf("%s holds %s.", a.mob.Name, wearable.Name), a.mob)
		a.mob.equipItem(wearable, wearHold)
		return
	}

	if wearable.canWear(itemWearWield) {
		if wearable.Weight > uint(a.mob.ModifiedAttributes.Strength) {
			a.mob.notify("It is too heavy for you to wield.")
			return
		}

		a.mob.notify("You wield %s.", wearable.Name)
		a.mob.Room.notify(fmt.Sprintf("%s wields %s.", a.mob.Name, wearable.Name), a.mob)
		a.mob.equipItem(wearable, wearWield)
		return
	}

	a.mob.notify("You can't wear, wield, or hold that.")
}

func (a *action) matchesSubject(s string) bool {
	for _, v := range strings.Split(strings.ToLower(s), " ") {
		if strings.HasPrefix(v, a.args[1]) {
			return true
		}
	}

	return false
}

func transferItem(i int, from []*item, to []*item) ([]*item, []*item) {
	item := from[i]
	from = append(from[0:i], from[i+1:]...)
	to = append(to, item)

	return from, to
}
