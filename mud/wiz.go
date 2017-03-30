package mud

import (
	"fmt"
	"strconv"
	"strings"

	"bytes"

	"github.com/brianseitel/oasis-mud/helpers"
)

func (wiz *mob) help(args []string) {
	// do help
}

func (wiz *mob) at(args []string) {
	if len(args) < 2 {
		wiz.notify("At where what?")
		return
	}

	where, what := args[1], strings.Join(args[2:], " ")

	location := wiz.findLocation([]string{where})
	if location == nil {
		wiz.notify("No such location.")
		return
	}

	if location.isPrivate() {
		wiz.notify("That room is private right now.")
		return
	}

	original := wiz.Room
	wiz.Room = location

	newAction(wiz, wiz.client, what)

	wiz.Room = original
	return
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

func (wiz *mob) goTo(args []string) {
	if len(args) < 1 {
		wiz.notify("Goto where?")
		return
	}

	location := wiz.findLocation(args[1:1])
	if location == nil {
		wiz.notify("No such location.")
		return
	}

	if location.isPrivate() {
		wiz.notify("That room is private right now.")
		return
	}

	if wiz.Fight != nil {
		wiz.stopFighting(true)
	}

	if !helpers.HasBit(wiz.Act, playerWizInvis) {
		act("$n $T.", wiz, nil, wiz.Bamfout, actToRoom)
	}

	wiz.Room.removeMob(wiz)
	wiz.Room = location
	wiz.Room.Mobs = append(wiz.Room.Mobs, wiz)

	newAction(wiz, wiz.client, "look")
}

func (wiz *mob) mfind(args []string) {
	if len(args) < 1 {
		wiz.notify("Mfind whom?")
		return
	}

	all := args[1] == "all"
	found := false

	for e := mobList.Front(); e != nil; e = e.Next() {
		m := e.Value.(*mob)
		if all || helpers.MatchesSubject(m.Name, args[1]) {
			wiz.notify("[%5d] %s", m.ID, strings.Title(m.Description))
			found = true
		}
	}

	if !found {
		wiz.notify("Nothing like that exists.")
	}
	return
}

func (wiz *mob) mstat(args []string) {
	if len(args) < 2 {
		wiz.notify("Mstat whom?")
		return
	}

	var victim *mob
	for e := mobList.Front(); e != nil; e = e.Next() {
		m := e.Value.(*mob)
		if helpers.MatchesSubject(m.Name, args[1]) {
			victim = m
			break
		}
	}

	if victim == nil {
		wiz.notify("They aren't here.")
		return
	}

	wiz.notify("Name: %s. ID: %d. Sex: %d. Room: %d.", victim.Name, victim.ID, victim.Gender, victim.Room.ID)
	wiz.notify("Str: %d. Int: %d. Wis: %d. Dex: %d. Cha: %d. Con: %d.", victim.ModifiedAttributes.Strength, victim.ModifiedAttributes.Intelligence, victim.ModifiedAttributes.Wisdom, victim.ModifiedAttributes.Dexterity, victim.ModifiedAttributes.Charisma, victim.ModifiedAttributes.Constitution)
	wiz.notify("HP: %d/%d. Mana: %d/%d. Movement: %d/%d. Practices: %d.", victim.Hitpoints, victim.MaxHitpoints, victim.Mana, victim.MaxMana, victim.Movement, victim.MaxMovement, victim.Practices)
	wiz.notify("Level: %d. Class: %s. Alignment: %d. AC: %d. Gold: %d. Experience: %d.", victim.Level, victim.Job.Name, victim.Alignment, victim.Armor, victim.Gold, victim.Exp)
	wiz.notify("Hitroll: %d. Damroll: %d. Position: %d. Wimpy: %d", victim.Hitroll, victim.Damroll, victim.Status, 0)

	fighter := "(none)"
	if victim.Fight == nil {
		fighter = victim.Fight.Name
	}
	wiz.notify("Fighting: %s", fighter)

	wiz.notify("Carrying: %d. Carry Weight: %d.", victim.Carrying, victim.CarryWeight)

	wiz.notify("Master: %s. Leader: %s. Affected By: %s", victim.master.Name, victim.leader.Name, affectBitName(int(victim.AffectedBy)))
	wiz.notify("Description: %s", victim.Description)

	for _, af := range victim.Affects {
		wiz.notify("Spell: '%s' modifies %s by %d for %d with bits %s.", af.affectType.Skill.Name, af.location, af.modifier, af.duration, affectBitName(int(af.bitVector)))
	}
	return
}

func (wiz *mob) mwhere(args []string) {
	if len(args) < 1 {
		wiz.notify("Mwhere whom?")
		return
	}

	found := false
	for e := mobList.Front(); e != nil; e = e.Next() {
		m := e.Value.(*mob)
		if m.isNPC() && m.Room != nil && helpers.MatchesSubject(m.Name, args[1]) {
			found = true
			wiz.notify("[%5d] %-28s [%5d] %s", m.index.ID, m.Description, m.Room.ID, m.Room.Name)
		}
	}

	if !found {
		act("You didn't find any $T.", wiz, nil, args[1], actToChar)
	}

}

func (wiz *mob) ofind(args []string) {
	if len(args) < 1 {
		wiz.notify("Ofind whom?")
		return
	}

	all := args[1] == "all"
	found := false

	for e := itemList.Front(); e != nil; e = e.Next() {
		i := e.Value.(*item)
		if all || helpers.MatchesSubject(i.Name, args[1]) {
			wiz.notify("[%5d] %s", i.ID, strings.Title(i.ShortDescription))
			found = true
		}
	}

	if !found {
		wiz.notify("Nothing like that exists.")
	}
	return
}

func (wiz *mob) ostat(args []string) {
	if len(args) < 2 {
		wiz.notify("Ostat what?")
		return
	}

	var obj *item
	for e := itemList.Front(); e != nil; e = e.Next() {
		i := e.Value.(*item)
		if helpers.MatchesSubject(i.Name, args[1]) {
			obj = i
			break
		}
	}

	if obj == nil {
		wiz.notify("Nothing like that exists.")
		return
	}

	wiz.notify("Name: %s", obj.Name)
	wiz.notify("ID: %d. Type: %d.", obj.ID, obj.ItemType)
	wiz.notify("Short Description: %s", obj.ShortDescription)
	wiz.notify("Description: %s", obj.Description)
	wiz.notify("Wear bits: %d. Extra bits: %d.", obj.WearFlags, obj.ExtraFlags)
	wiz.notify("Weight: %d.", obj.Weight)
	wiz.notify("Cost: %d. Timer: %d. Level: %d.", obj.Cost, obj.Timer, obj.Level)

	wiz.notify("In room: %d. Carried by: %s. Wear Location: %d", obj.Room.ID, obj.carriedBy.Name, obj.WearLocation)

	for _, af := range obj.Affected {
		wiz.notify("Affects %s by %d", af.location, af.modifier)
	}

	for _, af := range obj.index.Affected {
		wiz.notify("Affects %s by %d", af.location, af.modifier)
	}
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

func (wiz *mob) reboo(args []string) {
	wiz.notify("If you want to reboot, spell it out.")
	return
}

func (wiz *mob) reboot(args []string) {
	wiz.echo([]string{fmt.Sprintf("Reboot by %s", wiz.Name)})
	gameServer.Up = false
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

func (wiz *mob) rstat(args []string) {
	var location *room
	if len(args) < 2 {
		location = wiz.Room
	} else {
		location = wiz.findLocation(args[1:1])
	}

	if location == nil {
		wiz.notify("No such location.")
		return
	}

	if wiz.Room != location && location.isPrivate() {
		wiz.notify("That room is private right now.")
		return
	}

	wiz.notify("Name: %s.", location.Name)
	wiz.notify("Area: %s.", location.Area.Name)

	wiz.notify("ID: %d", location.ID)
	wiz.notify("Sector: %d", location.SectorType)
	wiz.notify("Light: %d", location.Light)

	wiz.notify("Room Flags: %d.", location.RoomFlags)
	wiz.notify("Description: %s.", location.Description)

	wiz.notify("Characters:")
	var buf bytes.Buffer
	for _, m := range location.Mobs {
		if m.Playable {
			buf.Write([]byte(fmt.Sprintf(" %s", m.Name)))
		}
	}
	wiz.notify(buf.String())

	wiz.notify("Objects:")
	var buf2 bytes.Buffer
	for _, o := range location.Items {
		buf2.Write([]byte(fmt.Sprintf(" %s", o.Name)))
	}
	wiz.notify(buf2.String())

	for _, e := range location.Exits {
		wiz.notify("Door: %s. To: %d. Key: %d. Exit Flags: %d.", e.Dir, e.Room.ID, e.Key, e.Flags)
		wiz.notify("Keyword: %s. Description: %s.", e.Keyword, e.Description)
	}
	return
}

func (wiz *mob) shutdow(args []string) {
	wiz.notify("If you want to SHUTDOWN, spell it out.")
	return
}

func (wiz *mob) shutdown(args []string) {
	wiz.echo([]string{fmt.Sprintf("Shutdown by %s", wiz.Name)})
	gameServer.Up = false
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
