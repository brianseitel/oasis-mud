package mud

import (
	"errors"
	"fmt"
	"strings"

	"github.com/brianseitel/oasis-mud/helpers"
)

type action struct {
	mob   *Mob
	rooms []*Room
	conn  *Connection
	args  []string
}

func newAction(p *Mob, c *Connection, i string) {
	newActionWithInput(&action{mob: p, conn: c, args: strings.Split(i, " ")})
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
	// case cWear:
	// 	a.wear()
	// 	return
	// case cRemove:
	// 	a.remove()
	// 	return
	// case cKill:
	// 	a.kill()
	// case cFlee:
	// 	a.flee()
	default:
		a.conn.SendString("Eh?" + helpers.Newline)
	}
	return nil
}

func (i *action) getCommand() command {
	for _, c := range commands {
		if isCommand(c, i.args[0]) == true {
			return c
		}
	}

	return cNoop
}

func isCommand(c command, p string) bool {
	return strings.HasPrefix(string(c), p)
}

func (a *action) look() {
	r := a.mob.room
	a.conn.SendString(
		fmt.Sprintf(
			"%s [ID: %d]\n%s\n%s%s%s",
			r.Name,
			r.Id,
			r.Description,
			exitsString(r),
			itemsString(r),
			mobsString(r, a.mob),
		),
	)
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

func exitsString(r Room) string {
	var exits string

	for _, e := range r.Exits {
		exits = fmt.Sprintf("%s%s ", exits, string(e.Dir))
	}

	return fmt.Sprintf("[%s]%s%s", strings.Trim(exits, " "), helpers.Newline, helpers.Newline)
}

func itemsString(r Room) string {
	var items string

	for _, i := range r.Items {
		items = fmt.Sprintf("%s is here.\n%s", i.Name, items)
	}
	return items
}

func mobsString(r Room, mob *Mob) string {
	var mobs string
	mobs = ""
	for _, m := range r.Mobs {
		if m.Name != mob.Name {
			mobs = fmt.Sprintf("%s is here.\n%s", m.Name, mobs)
		}
	}

	return mobs
}

func inventoryString(p *Mob) []string {
	inventory := make(map[string]int)

	for _, i := range p.Inventory {
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
	for _, e := range a.mob.room.Exits {
		if e.Dir == d {
			a.mob.move(e)
			a.mob.room.Mobs[a.mob.pid] = *a.mob
			Registry.rooms[a.mob.room.Id] = a.mob.room
			newAction(a.mob, a.conn, "look")
			return
		}
	}
	a.conn.SendString("Alas, you cannot go that way." + helpers.Newline)
}

func (a *action) drop() {
	for j, item := range a.mob.Inventory {
		if a.matchesSubject(item.Identifiers) {
			a.mob.Inventory, a.mob.room.Items = transferItem(j, a.mob.Inventory, a.mob.room.Items)
			Registry.rooms[a.mob.room.Id] = a.mob.room
			a.conn.SendString(fmt.Sprintf("You drop %s.", item.Name) + helpers.Newline)
			return
		}
	}
}

func (a *action) get() {
	for j, item := range a.mob.room.Items {
		if a.matchesSubject(item.Identifiers) {
			a.mob.room.Items, a.mob.Inventory = transferItem(j, a.mob.room.Items, a.mob.Inventory)
			Registry.rooms[a.mob.room.Id] = a.mob.room
			a.conn.SendString(fmt.Sprintf("You pick up %s.", item.Name) + helpers.Newline)
			return
		}
	}
}

func (a *action) quit() {
	a.conn.SendString("Seeya!" + helpers.Newline)
	a.conn.conn.Close()
}

func (a *action) matchesSubject(s []string) bool {
	for _, v := range s {
		if strings.HasPrefix(v, a.args[1]) {
			return true
		}
	}

	return false
}

func transferItem(i int, from []Item, to []Item) ([]Item, []Item) {
	item := from[i]
	from = append(from[0:i], from[i+1:]...)
	to = append(to, item)

	return from, to
}
