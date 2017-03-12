package mud

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

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
	case cStats:
		a.stats()
		return nil
	case cKill:
		a.kill()
		return nil
	case cFlee:
		a.flee()
		return nil
	// case cWear:
	//  a.wear()
	//  return
	// case cRemove:
	//  a.remove()
	//  return
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

func (a *action) look() {
	if a.mob == nil {
		fmt.Println("shit!")
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
		a.conn.SendString("Who?" + helpers.Newline)
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

func (a *action) stats() {
	const (
		width int = 50
	)

	username := a.mob.Name
	id := fmt.Sprintf("Level %d %s %s", a.mob.Level, a.mob.Race.Name, a.mob.Job.Name)
	spaces := width - len(username) - len(id)

	title := fmt.Sprintf("%s%s%s", username, strings.Repeat(" ", spaces), id)

	strength := fmt.Sprintf("%s%s%s%s%s%s%s", "Strength", strings.Repeat(" ", 8), strconv.Itoa(a.mob.Strength), strings.Repeat(" ", 11), "Experience", strings.Repeat(" ", 11-len(strconv.Itoa(a.mob.Exp))), strconv.Itoa(a.mob.Exp))
	wisdom := fmt.Sprintf("%s%s%s%s%s%s%s", "Wisdom", strings.Repeat(" ", 10), strconv.Itoa(a.mob.Wisdom), strings.Repeat(" ", 11), "TNL", strings.Repeat(" ", 18-len(strconv.Itoa(a.mob.TNL()))), strconv.Itoa(a.mob.TNL()))
	intel := fmt.Sprintf("%s%s%s", "Intelligence", strings.Repeat(" ", 4), strconv.Itoa(a.mob.Intelligence))
	dexterity := fmt.Sprintf("%s%s%s", "Dexterity", strings.Repeat(" ", 7), strconv.Itoa(a.mob.Dexterity))
	constitution := fmt.Sprintf("%s%s%s", "Constitution", strings.Repeat(" ", 4), strconv.Itoa(a.mob.Constitution))
	charisma := fmt.Sprintf("%s%s%s", "Charisma", strings.Repeat(" ", 8), strconv.Itoa(a.mob.Dexterity))
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

func (a *action) move(d string) {
	if a.mob.Status != standing {
		switch a.mob.Status {
		case fighting:
			a.conn.SendString("You can't move while fighting!" + helpers.Newline)
			break
		}
		fmt.Println("fuck")
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
	for j, item := range a.mob.Room.Items {
		if a.matchesSubject(item.Identifiers) {
			a.mob.Room.Items, a.mob.Inventory = transferItem(j, a.mob.Room.Items, a.mob.Inventory)
			// a.mob.Room.removeItem(item)
			// a.mob.addItem(item)
			message := fmt.Sprintf("%s picks up %s.\n", a.mob.Name, item.Name)
			for _, m := range a.mob.Room.Mobs {
				if m.ID == a.mob.ID {
					m.notify(fmt.Sprintf("You pick up %s.\n", item.Name))
				} else {
					m.notify(message)
				}
			}

			helpers.Dump(a.mob.Room.Items)
			return
		}
	}

	a.mob.notify("Get what?" + helpers.Newline)
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
