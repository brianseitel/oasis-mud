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
	// case cWear:
	//  a.wear()
	//  return
	// case cRemove:
	//  a.remove()
	//  return
	// case cKill:
	//  a.kill()
	// case cFlee:
	//  a.flee()
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
	r := a.mob.getRoom()

	if len(a.args) == 1 {
		a.conn.SendString(
			fmt.Sprintf(
				"%s [ID: %d]\n%s\n%s%s%s",
				r.Name,
				r.ID,
				r.Description,
				exitsString(r.Exits),
				itemsString(r.Items),
				mobsString(r.Mobs),
			),
		)
	} else {
		for _, mob := range r.Mobs {
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

func exitsString(exits []exit) string {
	var output string
	for _, e := range exits {
		output = fmt.Sprintf("%s%s ", output, string(e.Dir))
	}

	return fmt.Sprintf("[%s]%s%s", strings.Trim(output, " "), helpers.Newline, helpers.Newline)
}

func itemsString(items []item) string {
	var output string

	for _, i := range items {
		output = fmt.Sprintf("%s is here.\n%s", i.Name, output)
	}
	return output
}

func mobsString(mobs []mob) string {
	var output string
	output = ""
	for _, m := range mobs {
		output = fmt.Sprintf("%s is here.\n%s", m.Name, output)
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
	room := a.mob.getRoom()
	for _, e := range room.Exits {
		if e.Dir == d {
			a.mob.move(e)
			a.mob.RoomID = e.RoomID
			db.Set("gorm:save_associations", false).Save(&a.mob)
			newAction(a.mob, a.conn, "look")
			return
		}
	}
	a.conn.SendString("Alas, you cannot go that way." + helpers.Newline)
}

func (a *action) drop() {
	room := a.mob.getRoom()
	for _, item := range a.mob.Inventory {
		if a.matchesSubject(item.Identifiers) {
			db.Model(&room).Association("Items").Append(item)
			db.Model(&a.mob).Association("Inventory").Delete(item)
			a.conn.SendString(fmt.Sprintf("You drop %s.", item.Name) + helpers.Newline)
			return
		}
	}
}

func (a *action) get() {
	room := a.mob.getRoom()
	for _, item := range room.Items {
		if a.matchesSubject(item.Identifiers) {
			db.Model(&room).Association("Items").Delete(item)
			db.Model(&a.mob).Association("Inventory").Append(item)
			a.conn.SendString(fmt.Sprintf("You pick up %s.", item.Name) + helpers.Newline)
			return
		}
	}
}

func (a *action) kill() {
	room := a.mob.getRoom()
	for _, m := range room.Mobs {
		if a.matchesSubject(m.Identifiers) {
			newFight(a.mob, &m)
			return
		}
	}

	a.mob.notify("You can't find them.")
}

func (a *action) quit() {
	a.conn.SendString("Seeya!" + helpers.Newline)
	a.conn.end()
}

func (a *action) matchesSubject(s string) bool {
	for _, v := range strings.Split(s, ",") {
		if strings.HasPrefix(v, a.args[1]) {
			return true
		}
	}

	return false
}

func transferItem(i int, from []item, to []item) ([]item, []item) {
	item := from[i]
	from = append(from[0:i], from[i+1:]...)
	to = append(to, item)

	return from, to
}
