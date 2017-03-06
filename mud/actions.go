package mud

import (
	"errors"
	"fmt"
	"strings"

	"github.com/brianseitel/oasis-mud/helpers"
)

type action struct {
	player *Player
	rooms  []*Room
	conn   *Connection
	args   []string
}

func newAction(p *Player, c *Connection, i string) {
	newActionWithInput(&action{player: p, conn: c, args: strings.Split(i, " ")})
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
	r := a.player.getRoom()

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
}

func (a *action) inventory() {
	a.conn.SendString(
		fmt.Sprintf("Inventory\n%s\n%s\n%s",
			"-----------------------------------",
			strings.Join(inventoryString(a.player), helpers.Newline),
			"-----------------------------------",
		) + helpers.Newline,
	)
}

func exitsString(exits []Exit) string {
	var output string
	for _, e := range exits {
		output = fmt.Sprintf("%s%s ", output, string(e.Dir))
	}

	return fmt.Sprintf("[%s]%s%s", strings.Trim(output, " "), helpers.Newline, helpers.Newline)
}

func itemsString(items []Item) string {
	var output string

	for _, i := range items {
		output = fmt.Sprintf("%s is here.\n%s", i.Name, output)
	}
	return output
}

func mobsString(mobs []Mob) string {
	var output string
	output = ""
	for _, m := range mobs {
		output = fmt.Sprintf("%s is here.\n%s", m.Name, output)
	}

	return output
}

func inventoryString(p *Player) []string {
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
	room := a.player.getRoom()
	for _, e := range room.Exits {
		if e.Dir == d {
			a.player.move(e)
			a.player.RoomId = int(e.RoomId)
			db.Save(&a.player)
			newAction(a.player, a.conn, "look")
			return
		}
	}
	a.conn.SendString("Alas, you cannot go that way." + helpers.Newline)
}

func (a *action) drop() {
	room := a.player.getRoom()
	for _, item := range a.player.Inventory {
		if a.matchesSubject(item.Identifiers) {
			db.Model(&room).Association("Items").Append(item)
			db.Model(&a.player).Association("Inventory").Delete(item)
			a.conn.SendString(fmt.Sprintf("You drop %s.", item.Name) + helpers.Newline)
			return
		}
	}
}

func (a *action) get() {
	room := a.player.getRoom()
	for _, item := range room.Items {
		if a.matchesSubject(item.Identifiers) {
			db.Model(&room).Association("Items").Delete(item)
			db.Model(&a.player).Association("Inventory").Append(item)
			a.conn.SendString(fmt.Sprintf("You pick up %s.", item.Name) + helpers.Newline)
			return
		}
	}
}

func (a *action) quit() {
	a.conn.SendString("Seeya!" + helpers.Newline)
	a.conn.conn.Close()
}

func (a *action) matchesSubject(s string) bool {
	for _, v := range strings.Split(s, ",") {
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
