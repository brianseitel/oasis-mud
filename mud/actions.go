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
	default:
		a.conn.SendString("Eh?" + helpers.Newline)
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
		a.conn.SendString(fmt.Sprintf("%s%s%d%s", name, strings.Repeat(" ", spaces), level, helpers.Newline))
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
				a.conn.SendString(fmt.Sprintf("You look at %s.", mob.Name) + helpers.Newline)
				a.conn.SendString(helpers.WordWrap(mob.Description, 50) + helpers.Newline)
				return
			}
		}

		for _, item := range a.mob.Room.Items {
			if a.matchesSubject(item.Identifiers) {
				a.conn.SendString(fmt.Sprintf("You look at %s.", item.Name) + helpers.Newline)
				a.conn.SendString(helpers.WordWrap(item.Description, 50) + helpers.Newline)
				return
			}
		}
		a.conn.SendString("Look at what?" + helpers.Newline)
	}
}

func (a *action) inventory() {
	a.conn.SendString(
		fmt.Sprintf("Inventory\n%s\n%s\n%s",
			"-----------------------------------",
			strings.Join(inventoryString(a.mob), helpers.Newline),
			"-----------------------------------",
		) + helpers.Newline,
	)
}

func (a *action) equipment() {
	a.conn.SendString(
		fmt.Sprintf("Equipment\n%s\n%s\n%s",
			"-----------------------------------",
			strings.Join(equippedString(a.mob), helpers.Newline),
			"-----------------------------------",
		) + helpers.Newline,
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

	strength := fmt.Sprintf("%s%s%d%s%s%s%d", "Strength", strings.Repeat(" ", 8), a.mob.Strength, strings.Repeat(" ", 11), "Experience", strings.Repeat(" ", 11-len(strconv.Itoa(a.mob.Exp))), a.mob.Exp)
	wisdom := fmt.Sprintf("%s%s%d%s%s%s%d", "Wisdom", strings.Repeat(" ", 10), a.mob.Wisdom, strings.Repeat(" ", 11), "TNL", strings.Repeat(" ", 18-len(strconv.Itoa(a.mob.TNL()))), a.mob.TNL())
	intel := fmt.Sprintf("%s%s%d%s%s%s%d", "Intelligence", strings.Repeat(" ", 4), a.mob.Intelligence, strings.Repeat(" ", 11), "Alignment", strings.Repeat(" ", 12-len(strconv.Itoa(a.mob.Alignment))), a.mob.Alignment)
	dexterity := fmt.Sprintf("%s%s%d%s%s%s%d", "Dexterity", strings.Repeat(" ", 7), a.mob.Dexterity, strings.Repeat(" ", 11), "Practices", strings.Repeat(" ", 12-len(strconv.Itoa(int(a.mob.Practices)))), a.mob.Practices)
	constitution := fmt.Sprintf("%s%s%d", "Constitution", strings.Repeat(" ", 4), a.mob.Constitution)
	charisma := fmt.Sprintf("%s%s%d", "Charisma", strings.Repeat(" ", 8), a.mob.Dexterity)
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

	return fmt.Sprintf("[%s]%s%s", strings.Trim(output, " "), helpers.Newline, helpers.Newline)
}

func itemsString(items []*item) string {
	var output string

	for _, i := range items {
		output = fmt.Sprintf("%s is here.\n%s", i.Name, output)
	}
	return output
}

func mobsString(mobs []*mob, player *mob) string {
	var output string
	output = ""
	for _, m := range mobs {
		if m != player {
			output = fmt.Sprintf("%s is here.\n%s", m.Name, output)
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
	lines = append(lines, fmt.Sprintf("<light>     %s", m.equipped("light")))
	lines = append(lines, fmt.Sprintf("<finger>    %s", m.equipped("finger1")))
	lines = append(lines, fmt.Sprintf("<finger>    %s", m.equipped("finger2")))
	lines = append(lines, fmt.Sprintf("<neck>      %s", m.equipped("neck1")))
	lines = append(lines, fmt.Sprintf("<neck>      %s", m.equipped("neck2")))
	lines = append(lines, fmt.Sprintf("<torso>     %s", m.equipped("torso")))
	lines = append(lines, fmt.Sprintf("<head>      %s", m.equipped("head")))
	lines = append(lines, fmt.Sprintf("<legs>      %s", m.equipped("legs")))
	lines = append(lines, fmt.Sprintf("<feet>      %s", m.equipped("feet")))
	lines = append(lines, fmt.Sprintf("<hands>     %s", m.equipped("hands")))
	lines = append(lines, fmt.Sprintf("<arms>      %s", m.equipped("arms")))
	lines = append(lines, fmt.Sprintf("<shield>    %s", m.equipped("shield")))
	lines = append(lines, fmt.Sprintf("<body>      %s", m.equipped("body")))
	lines = append(lines, fmt.Sprintf("<waist>     %s", m.equipped("waist")))
	lines = append(lines, fmt.Sprintf("<wrist>     %s", m.equipped("wrist1")))
	lines = append(lines, fmt.Sprintf("<wrist>     %s", m.equipped("wrist2")))
	lines = append(lines, fmt.Sprintf("<wield>     %s", m.equipped("wield")))
	lines = append(lines, fmt.Sprintf("<held>      %s", m.equipped("held")))
	lines = append(lines, fmt.Sprintf("<floating>  %s", m.equipped("floating")))
	lines = append(lines, fmt.Sprintf("<secondary> %s", m.equipped("secondary")))

	return lines
}

func (a *action) move(d string) {
	if a.mob.Status != standing {
		switch a.mob.Status {
		case fighting:
			a.conn.SendString("You can't move while fighting!" + helpers.Newline)
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
	a.conn.SendString("Alas, you cannot go that way." + helpers.Newline)
}

func (a *action) drop() {
	if a.args[1] == "all" {
		for _, item := range a.mob.Inventory {
			a.mob.Room.Items = append(a.mob.Room.Items, item)
			message := fmt.Sprintf("%s picks up %s.\n", a.mob.Name, item.Name)
			for _, m := range a.mob.Room.Mobs {
				if m.ID == a.mob.ID {
					m.notify(fmt.Sprintf("You pick up %s.\n", item.Name))
				} else {
					m.notify(message)
				}
			}
		}
		a.mob.Inventory = nil
		return
	}

	for j, item := range a.mob.Inventory {
		if a.matchesSubject(item.Identifiers) {
			a.mob.Inventory, a.mob.Room.Items = transferItem(j, a.mob.Inventory, a.mob.Room.Items)
			message := fmt.Sprintf("%s drops %s.\n", a.mob.Name, item.Name)
			for _, m := range a.mob.Room.Mobs {
				if m.ID == a.mob.ID {
					m.notify(fmt.Sprintf("You drop %s.\n", item.Name))
				} else {
					m.notify(message)
				}
			}
			return
		}
	}

	a.mob.notify("Drop what?")
}

func (a *action) get() {
	if a.args[1] == "all" {
		for _, item := range a.mob.Room.Items {
			a.mob.Inventory = append(a.mob.Inventory, item)
			message := fmt.Sprintf("%s picks up %s.\n", a.mob.Name, item.Name)
			for _, m := range a.mob.Room.Mobs {
				if m.ID == a.mob.ID {
					m.notify(fmt.Sprintf("You pick up %s.\n", item.Name))
				} else {
					m.notify(message)
				}
			}
		}
		a.mob.Room.Items = nil
		return
	}

	for j, item := range a.mob.Room.Items {
		if a.args[1] == "all" || a.matchesSubject(item.Identifiers) {
			a.mob.Room.Items, a.mob.Inventory = transferItem(j, a.mob.Room.Items, a.mob.Inventory)
			message := fmt.Sprintf("%s picks up %s.\n", a.mob.Name, item.Name)
			for _, m := range a.mob.Room.Mobs {
				if m.ID == a.mob.ID {
					m.notify(fmt.Sprintf("You pick up %s.\n", item.Name))
				} else {
					m.notify(message)
				}
			}
			return
		}
	}

	a.mob.notify("Get what?" + helpers.Newline)
}

func (a *action) wear() {
	for j, item := range a.mob.Inventory {
		if a.matchesSubject(item.Identifiers) {
			for k, eq := range a.mob.Equipped {
				if eq.Position == item.Position {
					a.mob.Equipped, a.mob.Inventory = transferItem(k, a.mob.Equipped, a.mob.Inventory)
					a.mob.notify(fmt.Sprintf("You remove %s and put it in your inventory.%s", eq.Name, helpers.Newline))
				}
			}
			a.mob.Inventory, a.mob.Equipped = transferItem(j, a.mob.Inventory, a.mob.Equipped)
			a.mob.notify(fmt.Sprintf("You wear %s.%s", item.Name, helpers.Newline))
			return
		}
	}

	a.mob.notify("You can't find that.")
}

func (a *action) remove() {
	for j, item := range a.mob.Equipped {
		if a.matchesSubject(item.Identifiers) {
			a.mob.Equipped, a.mob.Inventory = transferItem(j, a.mob.Equipped, a.mob.Inventory)
			a.mob.notify(fmt.Sprintf("You remove %s.%s", item.Name, helpers.Newline))
			return
		}
	}

	a.mob.notify(fmt.Sprintf("You aren't wearing that.%s", helpers.Newline))
}

func (a *action) kill() {
	for _, m := range a.mob.Room.Mobs {
		if a.matchesSubject(m.Identifiers) {
			newFight(a.mob, m)
			return
		}
	}

	a.mob.notify("You can't find them." + helpers.Newline)
}

func (a *action) train() {
	player := a.mob
	if player.isNPC() {
		fmt.Println("crap!")
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
		player.notify("You can't do that here\n")
		return
	}

	if len(a.args) == 1 {
		player.notify(fmt.Sprintf("You have %d practice sessions.%s", player.Practices, helpers.Newline))
		return
	}

	var cost uint

	costmap := []uint{5, 6, 7, 9, 12, 13, 15}

	var playerAbility int
	var playerOutput string

	if strings.HasPrefix(a.args[1], "str") {
		playerAbility = player.Strength
		playerOutput = "strength"
	} else if strings.HasPrefix(a.args[1], "int") {
		playerAbility = player.Intelligence
		playerOutput = "intelligence"
	} else if strings.HasPrefix(a.args[1], "wis") {
		playerAbility = player.Wisdom
		playerOutput = "wisdom"
	} else if strings.HasPrefix(a.args[1], "dex") {
		playerAbility = player.Dexterity
		playerOutput = "dexterity"
	} else if strings.HasPrefix(a.args[1], "cha") {
		playerAbility = player.Charisma
		playerOutput = "charisma"
	} else if strings.HasPrefix(a.args[1], "con") {
		playerAbility = player.Constitution
		playerOutput = "constitution"
	} else {
		var buf bytes.Buffer

		buf.WriteString("You can train:\r\n")
		if player.Strength < 18 {
			buf.WriteString(fmt.Sprintf("Strength      %d\r\n", costmap[player.Strength-12]))
		}
		if player.Intelligence < 18 {
			buf.WriteString(fmt.Sprintf("Intelligence  %d\r\n", costmap[player.Intelligence-12]))
		}
		if player.Wisdom < 18 {
			buf.WriteString(fmt.Sprintf("Wisdom        %d\r\n", costmap[player.Wisdom-12]))
		}
		if player.Dexterity < 18 {
			buf.WriteString(fmt.Sprintf("Dexterity     %d\r\n", costmap[player.Dexterity-12]))
		}
		if player.Charisma < 18 {
			buf.WriteString(fmt.Sprintf("Charisma      %d\r\n", costmap[player.Charisma-12]))
		}
		if player.Constitution < 18 {
			buf.WriteString(fmt.Sprintf("Constitution  %d\r\n", costmap[player.Constitution-12]))
		}

		message := buf.String()
		if !strings.HasSuffix(message, ":") {
			buf.WriteString(".\r\n")
			player.notify(buf.String())
		} else {
			player.notify("You have nothing left to train, you badass!\r\n")
		}

		return
	}

	cost = costmap[playerAbility-12]
	if playerAbility >= 18 {
		player.notify(fmt.Sprintf("Your %s is already at maximum.%s", playerOutput, helpers.Newline))
		return
	}

	if cost > player.Practices {
		player.notify("You don't have enough practices.\r\n")
		return
	}

	player.Practices -= cost
	switch playerOutput {
	case "strength":
		player.Strength++
		break
	case "intelligence":
		player.Intelligence++
		break
	case "wisdom":
		player.Wisdom++
		break
	case "dexterity":
		player.Dexterity++
		break
	case "charisma":
		player.Charisma++
		break
	case "constitution":
		player.Constitution++
		break
	}

	player.notify(fmt.Sprintf("Your %s increases for %d practice points!%s", playerOutput, cost, helpers.Newline))
	return
}

func (a *action) trip() {

	if a.mob.skill("trip") == nil {
		a.mob.notify("You don't know how to do this.\n")
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

func (a *action) flee() {
	if a.mob.Status != fighting {
		a.conn.SendString("You can't flee if you're not fighting, fool." + helpers.Newline)
		return
	}

	roll := dice().Intn(36 - a.mob.Dexterity) // higher the dexterity, better the chance of fleeing successfully
	if roll == 1 {
		a.mob.Fight.Mob1.Status = standing
		a.mob.Fight.Mob2.Status = standing

		a.mob.wander()
		a.conn.SendString("You flee!")
	} else {
		a.conn.SendString("You tried to flee, but failed!" + helpers.Newline)
	}
}

func (a *action) scan() {
	room := getRoom(a.mob.Room.ID)
	for _, x := range room.Exits {
		a.conn.SendString(fmt.Sprintf("[%s]%s", x.Dir, helpers.Newline))

		if len(x.Room.Mobs) > 0 {
			mobs := x.Room.Mobs
			for _, m := range mobs {
				a.conn.SendString(fmt.Sprintf("    %s\n", m.Name))
			}
		} else {
			a.conn.SendString(fmt.Sprintf("    %s(nothing)%s\n", helpers.Blue, helpers.Reset))
		}
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
		a.conn.SendString("Recall set!" + helpers.Newline)
		return
	}

	a.conn.SendString("Recall what?" + helpers.Newline)
}

func (a *action) quit() {
	if a.mob.Status == fighting {
		a.conn.SendString("You can't quit now. You're fighting!" + helpers.Newline)
	} else {
		a.conn.SendString("Seeya!" + helpers.Newline)
		a.conn.end()
	}
}

func (a *action) matchesSubject(s string) bool {
	for _, v := range strings.Split(s, ",") {
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
