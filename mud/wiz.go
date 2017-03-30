package mud

import (
	"strconv"
	"strings"

	"github.com/brianseitel/oasis-mud/helpers"
)

func (wiz *mob) help(args []string) {
	// do help
}

func (wiz *mob) bamfin(args []string) {

	if len(args) < 1 {
		wiz.notify("Set bamfin to where?")
		return
	}

	if !wiz.isNPC() {
		wiz.Bamfin = args[1]
	}
}

func (wiz *mob) bamfout(args []string) {

	if len(args) < 1 {
		wiz.notify("Set bamfin to where?")
		return
	}

	if !wiz.isNPC() {
		wiz.Bamfout = args[1]
	}
}

func (wiz *mob) deny(args []string) {
	if len(args) < 1 {
		wiz.notify("Deny whom?")
		return
	}

	victim := getPlayerByName(args[1])
	if victim == nil {
		wiz.notify("They aren't here.")
		return
	}

	if victim.isNPC() {
		wiz.notify("You can't do this on NPCs.")
		return
	}

	if victim.Trust > wiz.Trust {
		wiz.notify("You failed.")
		return
	}

	helpers.SetBit(victim.Act, playerDeny)
	victim.notify("You are denied access!")
	wiz.notify("OK.")
	newAction(victim, victim.client, "quit")
}

func (wiz *mob) disconnect(args []string) {
	if len(args) < 1 {
		wiz.notify("Disconnect whom?")
		return
	}

	victim := getPlayerByName(args[1])
	if victim == nil {
		wiz.notify("They aren't here.")
		return
	}

	if victim.client == nil {
		act("$N doesn't have a connection.", wiz, nil, victim, actToChar)
		return
	}

	for _, c := range gameServer.connections {
		if c == *victim.client {
			c.conn.Close()
			wiz.notify("OK")
			return
		}
	}

	wiz.notify("Descriptor not found.")
	return
}

func (wiz *mob) pardon(args []string) {
	if len(args) < 2 {
		wiz.notify("Syntax: pardon <character> <killer|thief>")
		return
	}

	arg1, arg2 := args[1], args[2]

	victim := getPlayerByName(arg1)
	if victim == nil {
		wiz.notify("They aren't here.")
		return
	}

	if victim.isNPC() {
		wiz.notify("You can't pardon NPCs.")
		return
	}

	if arg2 == "killer" {
		if helpers.HasBit(victim.Act, playerKiller) {
			helpers.RemoveBit(victim.Act, playerKiller)
			wiz.notify("Killer flag removed.")
			victim.notify("You are no longer a KILLER.")
		}
		return
	}

	if arg2 == "thief" {
		if helpers.HasBit(victim.Act, playerThief) {
			helpers.RemoveBit(victim.Act, playerThief)
			wiz.notify("Thief flag removed.")
			victim.notify("You are no longer a THIEF.")
		}
		return
	}

	wiz.notify("Syntax: pardon <character> <killer|thief>")
	return
}

func (wiz *mob) echo(args []string) {
	if len(args) < 1 {
		wiz.notify("Echo what?")
		return
	}

	for e := mobList.Front(); e != nil; e = e.Next() {
		m := e.Value.(*mob)
		m.notify(strings.Join(args, " "))
	}
	return
}

func (wiz *mob) recho(args []string) {
	if len(args) < 1 {
		wiz.notify("Recho what?")
		return
	}

	for e := mobList.Front(); e != nil; e = e.Next() {
		m := e.Value.(*mob)
		if m.Room == wiz.Room {
			m.notify(strings.Join(args, " "))
		}
	}
	return
}

func (wiz *mob) findLocation(args []string) *room {
	if len(args) < 1 {
		wiz.notify("Find what location?")
		return nil
	}

	num, err := strconv.Atoi(args[1])
	isNumber := err == nil

	if isNumber {
		return getRoom(uint(num))
	}

	victim := getPlayerByName(args[1])
	if victim != nil {
		return victim.Room
	}

	// get object room

	return nil
}

func (wiz *mob) transfer(args []string) {
	if len(args) < 2 {
		wiz.notify("Transfer whom?")
	}

	whom, where := args[1], args[2]

	if whom == "all" {
		for e := mobList.Front(); e != nil; e = e.Next() {
			m := e.Value.(*mob)
			if m.client != nil && m != wiz && m.Room != nil && wiz.canSee(m) {
				newArgs := []string{m.Name, where}
				m.transfer(newArgs)
			}
		}
		return
	}

	var location *room
	if where == "here" {
		location = wiz.Room
	} else {
		location := wiz.findLocation([]string{where})
		if location == nil {
			wiz.notify("No such location.")
			return
		}

		if location.isPrivate() {
			wiz.notify("That room is private right now.")
			return
		}
	}

	victim := getPlayerByName(whom)
	if victim == nil {
		wiz.notify("They aren't here.")
		return
	}

	if victim.Room == nil {
		wiz.notify("They are in limbo.")
		return
	}

	if victim.Fight != nil || victim.Status == fighting {
		victim.stopFighting(true)
	}

	act("$n disappears in a mushroom cloud.", victim, nil, nil, actToRoom)
	victim.Room.removeMob(victim)
	victim.Room = location
	victim.Room.Mobs = append(victim.Room.Mobs, victim)
	act("$n arrives in an explosion of rainbow sprinkles.", victim, nil, nil, actToRoom)
	if wiz != victim {
		act("$n has transferred you.", wiz, nil, victim, actToVict)
	}

	newAction(victim, victim.client, "look")
	wiz.notify("OK.")
}
