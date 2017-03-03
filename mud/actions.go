package mud

import (
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

func newActionWithInput(a *action) {

	switch a.getCommand() {
	case cLook:
		a.look()
		return
	case cNorth:
		a.move("north")
		return
	case cSouth:
		a.move("south")
		return
	case cEast:
		a.move("east")
		return
	case cWest:
		a.move("west")
		return
	case cUp:
		a.move("up")
		return
	case cDown:
		a.move("down")
		return
	// case cDrop:
	// 	a.drop()
	// 	return
	// case cGet:
	// 	a.get()
	// 	return
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
		a.conn.SendString("Eh?")
	}
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
	r := a.player.room
	a.conn.SendString(
		fmt.Sprintf(
			"%s\n%s\n%s%s",
			r.Name,
			r.Description,
			exitsString(r),
			itemsString(r),
			// mobsString(r, a.player),
		),
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

func (a *action) move(d string) {
	for _, e := range a.player.room.Exits {
		if e.Dir == d {
			rooms := a.rooms
			a.player.move(e, rooms)
			newAction(a.player, a.conn, "look")
			return
		}
	}
	a.conn.SendString("Alas, you cannot go that way." + helpers.Newline)
}
