package mud

import (
	"fmt"
	"strconv"
	"strings"

	"bytes"
)

func doAdvance(wiz *mob, argument string) {
	if len(argument) < 1 {
		wiz.notify("Syntax: advance <char> <level>.")
		return
	}

	argument, name := oneArgument(argument)
	argument, arg2 := oneArgument(argument)
	level, err := strconv.Atoi(arg2)
	if err != nil {
		wiz.notify("Syntax: advance <char> <level>.")
		return
	}

	victim := getPlayerByName(name)

	if victim == nil {
		wiz.notify("They aren't here.")
		return
	}

	if level > 99 || level < 0 {
		wiz.notify("Level must be between 1 and 99.")
		return
	}

	if level > wiz.Trust {
		wiz.notify("Limited to your trust level.")
		return
	}

	if level <= victim.Level {
		wiz.notify("Lowering a player's level!")
		victim.notify("*** OOOOHHHHHHH NNNNOOOOO ***")

		victim.Level = 1
		victim.Exp = 1000
		victim.MaxHitpoints = 10
		victim.MaxMana = 100
		victim.MaxMovement = 100
		victim.Skills = nil
		victim.Practices = 0
		victim.Hitpoints = victim.MaxHitpoints
		victim.Mana = victim.MaxMana
		victim.Movement = victim.MaxMovement
		victim.advanceLevel()
	} else {
		wiz.notify("Raising a player's level!")
		victim.notify("*** OOOOHHHHHHH YYYESSSSSS ***")
	}

	for i := victim.Level; i < level; i++ {
		victim.notify("You raise a level!")
		victim.Level++
		victim.advanceLevel()
	}

	victim.Exp = 1000 * max(1, victim.Level)
	victim.Trust = 0
	return
}

func doAllow(wiz *mob, argument string) {
	argument, arg1 := oneArgument(argument)

	if arg1 == "" {
		wiz.notify("Remove which site from the ban list?")
		return
	}

	for e := banList.Front(); e != nil; e = e.Next() {
		b := e.Value.(string)
		if b == arg1 {
			banList.Remove(e)
			break
		}
	}

	wiz.notify("Site is un-banned.")
}

func doAt(wiz *mob, argument string) {
	if len(argument) < 2 {
		wiz.notify("At where what?")
		return
	}

	what, where := oneArgument(argument)

	location := wiz.findLocation(where)
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

	interpret(wiz, what)

	wiz.Room = original
	return
}

func doBamfin(wiz *mob, argument string) {

	if len(argument) < 1 {
		wiz.notify("Set bamfin to where?")
		return
	}

	if !wiz.isNPC() {
		wiz.Bamfin = argument
	}
}

func doBamfout(wiz *mob, argument string) {

	if len(argument) < 1 {
		wiz.notify("Set bamfin to what?")
		return
	}

	if !wiz.isNPC() {
		wiz.Bamfout = argument
	}
}

func doBan(wiz *mob, argument string) {
	if wiz.isNPC() {
		return
	}

	argument, arg1 := oneArgument(argument)

	if arg1 == "" {
		wiz.notify("Banned sites:")
		for e := banList.Front(); e != nil; e = e.Next() {
			b := e.Value.(string)
			wiz.notify(b)
		}
		return
	}

	for e := banList.Front(); e != nil; e = e.Next() {
		b := e.Value.(string)
		if b == arg1 {
			wiz.notify("That site is already banned.")
			return
		}
	}

	banList.PushBack(arg1)
	wiz.notify("Ok.")
}

func doDeny(wiz *mob, argument string) {
	if len(argument) < 1 {
		wiz.notify("Deny whom?")
		return
	}

	argument, arg1 := oneArgument(argument)
	victim := getPlayerByName(arg1)
	if victim == nil {
		wiz.notify("They aren't here.")
		return
	}

	if victim.Trust > wiz.Trust {
		wiz.notify("You failed.")
		return
	}

	victim.Act = setBit(victim.Act, playerDeny)
	victim.notify("You are denied access!")
	wiz.notify("OK.")
	interpret(wiz, "quit")
}

func doDisconnect(wiz *mob, argument string) {
	if len(argument) < 1 {
		wiz.notify("Disconnect whom?")
		return
	}

	argument, arg1 := oneArgument(argument)
	victim := getPlayerByName(arg1)
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

func doEcho(wiz *mob, argument string) {
	if len(argument) < 1 {
		wiz.notify("Echo what?")
		return
	}

	for e := mobList.Front(); e != nil; e = e.Next() {
		m := e.Value.(*mob)
		m.notify(argument)
	}
	return
}

func (wiz *mob) findLocation(argument string) *room {
	if len(argument) < 1 {
		wiz.notify("Find what location?")
		return nil
	}

	argument, arg1 := oneArgument(argument)

	num, err := strconv.Atoi(arg1)
	isNumber := err == nil

	if isNumber {
		return getRoom(num)
	}

	victim := getPlayerByName(arg1)
	if victim != nil {
		return victim.Room
	}

	// TODO: get object room

	return nil
}

func doForce(wiz *mob, argument string) {
	if len(argument) < 2 {
		wiz.notify("Syntax: force <char> <action>")
		return
	}

	action, name := oneArgument(argument)

	if name == "all" {
		for e := mobList.Front(); e != nil; e = e.Next() {
			m := e.Value.(*mob)
			if m.client != nil && !m.isNPC() && m.getTrust() < wiz.getTrust() {
				act("$n forces you to '$t'.", wiz, action, m, actToVict)
				interpret(m, action)
			}
		}
	} else {
		victim := getPlayerByName(name)

		if victim == nil {
			wiz.notify("They aren't here.")
			return
		}

		if victim == wiz {
			wiz.notify("You're an idiot.")
			return
		}

		if victim.getTrust() >= wiz.getTrust() {
			wiz.notify("You can't make them do anything.")
			return
		}

		act("$n forces you to '$t'.", wiz, action, victim, actToVict)
		interpret(victim, action)
	}

	wiz.notify("Ok")
	return

}

func doFreeze(wiz *mob, argument string) {
	if len(argument) < 1 {
		wiz.notify("Freeze whom?")
		return
	}

	argument, arg1 := oneArgument(argument)

	victim := getPlayerByName(arg1)
	if victim == nil {
		wiz.notify("They aren't here.")
		return
	}

	if victim.getTrust() >= wiz.getTrust() {
		wiz.notify("You failed.")
		return
	}

	if hasBit(victim.Act, playerFreeze) {
		victim.Act = removeBit(victim.Act, playerFreeze)
		victim.notify("You can play again.")
		wiz.notify("FREEZE removed.")
	} else {
		victim.Act = setBit(victim.Act, playerFreeze)
		victim.notify("You can't do ANYthing.")
		wiz.notify("FREEZE set.")
	}

	// TODO: breaks during testing
	// saveCharacter(victim)
	return
}

func doGoto(wiz *mob, argument string) {
	if len(argument) < 1 {
		wiz.notify("Goto where?")
		return
	}

	location := wiz.findLocation(argument)
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

	if !hasBit(wiz.Act, playerWizInvis) {
		act("$n $T.", wiz, nil, wiz.Bamfout, actToRoom)
	}

	wiz.Room.removeMob(wiz)
	wiz.Room = location
	wiz.Room.Mobs = append(wiz.Room.Mobs, wiz)

	interpret(wiz, "look")
}

func doHolylight(wiz *mob, argument string) {
	if wiz.isNPC() {
		return
	}

	if hasBit(wiz.Act, playerHolylight) {
		wiz.Act = removeBit(wiz.Act, playerHolylight)
		wiz.notify("Holy light mode off.")
	} else {
		wiz.Act = setBit(wiz.Act, playerHolylight)
		wiz.notify("Holy light mode on.")
	}
}

func doInvis(wiz *mob, argument string) {
	if wiz.isNPC() {
		return
	}

	if hasBit(wiz.Act, playerWizInvis) {
		wiz.Act = removeBit(wiz.Act, playerWizInvis)
		act("$n slowly fades into existence.", wiz, nil, nil, actToRoom)
		wiz.notify("You slowly fade back into extistence.")
	} else {
		wiz.Act = setBit(wiz.Act, playerWizInvis)
		act("$n slowly fades out of sight.", wiz, nil, nil, actToRoom)
		wiz.notify("You slowly vanish out of sight.")
	}
	return
}

func doLog(wiz *mob, argument string) {
	argument, arg1 := oneArgument(argument)

	if arg1 == "" {
		wiz.notify("Log whom?")
		return
	}

	if arg1 == "all" {
		wiz.notify("You can't log all.")
	}

	victim := getPlayerByName(arg1)
	if victim == nil {
		wiz.notify("They aren't here.")
		return
	}

	if victim.isNPC() {
		wiz.notify("Not on NPCs.")
		return
	}

	if hasBit(victim.Act, playerLog) {
		victim.Act = removeBit(victim.Act, playerLog)
		wiz.notify("LOG removed.")
	} else {
		victim.Act = setBit(victim.Act, playerLog)
		wiz.notify("LOG added.")
	}
}

func doMemory(wiz *mob, argument string) {

	wiz.notify("Areas:    %5d", areaList.Len())
	wiz.notify("Bans:     %5d", banList.Len())
	wiz.notify("Commands: %5d", commandList.Len())
	wiz.notify("Helps:    %5d", helpList.Len())
	wiz.notify("Mobs:     %5d", mobList.Len())
	wiz.notify("Objects:  %5d", itemList.Len())
	wiz.notify("Rooms:    %5d", roomList.Len())
	wiz.notify("Shops:    %5d", shopList.Len())
	wiz.notify("Skills:   %5d", skillList.Len())
	wiz.notify("Socials:  %5d", socialList.Len())

	return
}

func doMfind(wiz *mob, argument string) {
	if len(argument) < 1 {
		wiz.notify("Mfind whom?")
		return
	}

	argument, arg1 := oneArgument(argument)
	all := arg1 == "all"
	found := false

	for e := mobList.Front(); e != nil; e = e.Next() {
		m := e.Value.(*mob)
		if all || matchesSubject(m.Name, argument) {
			wiz.notify("[%5d] %s", m.ID, strings.Title(m.Description))
			found = true
		}
	}

	if !found {
		wiz.notify("Nothing like that exists.")
	}
	return
}

func doMload(wiz *mob, argument string) {
	if len(argument) < 1 {
		wiz.notify("Syntax: mload <id>.")
		return
	}

	argument, arg1 := oneArgument(argument)
	id, err := strconv.Atoi(arg1)
	if err != nil {
		wiz.notify("Syntax: mload <id>.")
		return
	}
	mob := getMob(id)
	if mob == nil {
		wiz.notify("No mob has that ID.")
		return
	}

	victim := createMob(mob)
	victim.Room = wiz.Room
	victim.Room.Mobs = append(victim.Room.Mobs, victim)
	act("$n has created $N!", wiz, nil, victim, actToRoom)
	wiz.notify("Ok.")
	return
}

func doMstat(wiz *mob, argument string) {
	if len(argument) < 1 {
		wiz.notify("Mstat whom?")
		return
	}

	argument, arg1 := oneArgument(argument)

	var victim *mob
	for e := mobList.Front(); e != nil; e = e.Next() {
		m := e.Value.(*mob)
		if matchesSubject(m.Name, arg1) {
			victim = m
			break
		}
	}

	if victim == nil {
		wiz.notify("They aren't here.")
		return
	}

	wiz.notify("Name: %s. ID: %d. Sex: %d. Room: %d.", victim.Name, victim.ID, victim.Gender, victim.Room.ID)
	wiz.notify("Str: %d. Int: %d. Wis: %d. Dex: %d. Cha: %d. Con: %d.", victim.currentStrength(), victim.currentIntelligence(), victim.currentWisdom(), victim.currentDexterity(), victim.currentCharisma(), victim.currentConstitution())
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

func doMwhere(wiz *mob, argument string) {
	if len(argument) < 1 {
		wiz.notify("Mwhere whom?")
		return
	}

	argument, arg1 := oneArgument(argument)
	found := false
	for e := mobList.Front(); e != nil; e = e.Next() {
		m := e.Value.(*mob)
		if m.isNPC() && m.Room != nil && matchesSubject(m.Name, arg1) {
			found = true
			wiz.notify("[%5d] %-28s [%5d] %s", m.index.ID, m.Description, m.Room.ID, m.Room.Name)
		}
	}

	if !found {
		act("You didn't find any $T.", wiz, nil, arg1, actToChar)
	}

}

func doNoEmote(wiz *mob, argument string) {
	argument, arg1 := oneArgument(argument)

	if arg1 == "" {
		wiz.notify("NoEmote whom?")
		return
	}

	victim := getPlayerByName(arg1)
	if victim == nil {
		wiz.notify("They aren't here.")
		return
	}

	if victim.isNPC() {
		wiz.notify("Not on NPCs.")
		return
	}

	if victim.getTrust() >= wiz.getTrust() {
		wiz.notify("You failed.")
	}

	if hasBit(victim.Act, playerNoEmote) {
		victim.Act = removeBit(victim.Act, playerNoEmote)
		victim.notify("You can emote again!")
		wiz.notify("NO EMOTE removed.")
	} else {
		victim.Act = setBit(victim.Act, playerNoEmote)
		victim.notify("You can't emote!")
		wiz.notify("NO EMOTE set.")
	}
}

func doNoTell(wiz *mob, argument string) {
	argument, arg1 := oneArgument(argument)

	if arg1 == "" {
		wiz.notify("NoTell whom?")
		return
	}

	victim := getPlayerByName(arg1)
	if victim == nil {
		wiz.notify("They aren't here.")
		return
	}

	if victim.isNPC() {
		wiz.notify("Not on NPCs.")
		return
	}

	if victim.getTrust() >= wiz.getTrust() {
		wiz.notify("You failed.")
	}

	if hasBit(victim.Act, playerNoTell) {
		victim.Act = removeBit(victim.Act, playerNoTell)
		victim.notify("You can tell again!")
		wiz.notify("NO TELL removed.")
	} else {
		victim.Act = setBit(victim.Act, playerNoTell)
		victim.notify("You can't tell anymore!")
		wiz.notify("NO TELL set.")
	}
}

func doOfind(wiz *mob, argument string) {
	if len(argument) < 1 {
		wiz.notify("Ofind whom?")
		return
	}

	argument, arg1 := oneArgument(argument)
	all := arg1 == "all"
	found := false

	for e := itemList.Front(); e != nil; e = e.Next() {
		i := e.Value.(*item)
		if all || matchesSubject(i.Name, arg1) {
			wiz.notify("[%5d] %s", i.ID, strings.Title(i.ShortDescription))
			found = true
		}
	}

	if !found {
		wiz.notify("Nothing like that exists.")
	}
	return
}

func doOload(wiz *mob, argument string) {
	if len(argument) < 1 {
		wiz.notify("Syntax: oload <id> <level>")
		return
	}

	argument, arg1 := oneArgument(argument)
	argument, arg2 := oneArgument(argument)
	id, err := strconv.Atoi(arg1)
	if err != nil {
		wiz.notify("Syntax: oload <id> <level>")
		return
	}

	level, err := strconv.Atoi(arg2)
	if err != nil {
		wiz.notify("Syntax: oload <id> <level>")
		return
	}

	if level < 0 || level > wiz.Trust {
		wiz.notify("Limited to your trust level.")
		return
	}

	var objIndex *itemIndex
	for e := itemIndexList.Front(); e != nil; e = e.Next() {
		i := e.Value.(*itemIndex)
		if i.ID == id {
			objIndex = i
			break
		}
	}

	if objIndex == nil {
		wiz.notify("No object has that ID.")
		return
	}

	obj := createItem(objIndex)
	if obj.canWear(itemTake) {
		wiz.Inventory = append(wiz.Inventory, obj)
		obj.carriedBy = wiz
	} else {
		obj.Room = wiz.Room
		obj.Room.Items = append(obj.Room.Items, obj)
	}

	wiz.notify("Ok.")
	return
}

func doOstat(wiz *mob, argument string) {
	if len(argument) < 1 {
		wiz.notify("Ostat what?")
		return
	}

	argument, arg1 := oneArgument(argument)

	var obj *item
	for e := itemList.Front(); e != nil; e = e.Next() {
		i := e.Value.(*item)
		if matchesSubject(i.Name, arg1) {
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

	carriedBy := "(none)"
	if obj.carriedBy != nil {
		carriedBy = obj.carriedBy.Name
	}
	room := "(none)"
	if obj.Room != nil {
		room = strconv.Itoa(obj.Room.ID)
	}
	wiz.notify("In room: %s. Carried by: %s. Wear Location: %d", room, carriedBy, obj.WearLocation)

	for _, af := range obj.Affected {
		wiz.notify("Affects %s by %d", af.location, af.modifier)
	}

	for _, af := range obj.index.Affected {
		wiz.notify("Affects %s by %d", af.location, af.modifier)
	}
	return
}

func doPardon(wiz *mob, argument string) {
	if len(argument) < 1 {
		wiz.notify("Syntax: pardon <character> <killer|thief>")
		return
	}

	argument, arg1 := oneArgument(argument)
	argument, arg2 := oneArgument(argument)

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
		if hasBit(victim.Act, playerKiller) {
			victim.Act = removeBit(victim.Act, playerKiller)
			wiz.notify("Killer flag removed.")
			victim.notify("You are no longer a KILLER.")
		}
		return
	}

	if arg2 == "thief" {
		if hasBit(victim.Act, playerThief) {
			victim.Act = removeBit(victim.Act, playerThief)
			wiz.notify("Thief flag removed.")
			victim.notify("You are no longer a THIEF.")
		}
		return
	}

	wiz.notify("Syntax: pardon <character> <killer|thief>")
	return
}

func doPeace(wiz *mob, argument string) {
	for e := mobList.Front(); e != nil; e = e.Next() {
		m := e.Value.(*mob)
		m.stopFighting(true)
	}

	wiz.notify("Ok.")
}

func doPurge(wiz *mob, argument string) {
	if len(argument) < 1 {
		for _, m := range wiz.Room.Mobs {
			if m.isNPC() {
				extractMob(m, true)
			}
		}

		for _, i := range wiz.Room.Items {
			extractObj(i)
		}

		act("$n purges the room!", wiz, nil, nil, actToRoom)
		wiz.notify("OK.")
		return
	}

	argument, arg1 := oneArgument(argument)
	victim := getPlayerByName(arg1)
	if victim == nil {
		wiz.notify("They aren't here.")
		return
	}

	if !victim.isNPC() {
		wiz.notify("Not on PCs.")
		return
	}

	act("$n purges $N.", wiz, nil, victim, actToNotVict)
	extractMob(victim, true)
}

func doReboo(wiz *mob, argument string) {
	wiz.notify("If you want to reboot, spell it out.")
	return
}

func doReboot(wiz *mob, argument string) {
	doEcho(wiz, fmt.Sprintf("Reboot by %s", wiz.Name))
	gameServer.Up = false
	return
}

func doRestore(wiz *mob, argument string) {
	if len(argument) < 1 {
		wiz.notify("Restore whom?")
		return
	}
	argument, arg1 := oneArgument(argument)
	victim := getPlayerByName(arg1)
	if victim == nil {
		wiz.notify("They aren't here.")
		return
	}

	victim.Hitpoints = victim.MaxHitpoints
	victim.Mana = victim.MaxMana
	victim.Movement = victim.MaxMovement
	victim.updateStatus()
	act("$n has restored you.", wiz, nil, victim, actToVict)
	wiz.notify("OK.")
	return
}

func doReturn(wiz *mob, argument string) {
	if wiz.client == nil {
		return
	}

	if wiz.client.original == nil {
		wiz.notify("You aren't switched.")
		return
	}

	wiz.notify("You return to your original body.")
	wiz.client.mob = wiz.client.original
	wiz.client.original = nil
	wiz.client.mob.client = wiz.client
	wiz.client = nil
	return
}

func doRecho(wiz *mob, argument string) {
	if len(argument) < 1 {
		wiz.notify("Recho what?")
		return
	}

	for e := mobList.Front(); e != nil; e = e.Next() {
		m := e.Value.(*mob)
		if m.Room == wiz.Room {
			m.notify(argument)
		}
	}
	return
}

func doRstat(wiz *mob, argument string) {
	var location *room
	if len(argument) < 1 {
		location = wiz.Room
	} else {
		location = wiz.findLocation(argument)
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

func doSilence(wiz *mob, argument string) {
	argument, arg1 := oneArgument(argument)

	if arg1 == "" {
		wiz.notify("Silence whom?")
		return
	}

	victim := getPlayerByName(arg1)
	if victim == nil {
		wiz.notify("They aren't here.")
		return
	}

	if victim.isNPC() {
		wiz.notify("Not on NPCs.")
		return
	}

	if victim.getTrust() >= wiz.getTrust() {
		wiz.notify("You failed.")
	}

	if hasBit(victim.Act, playerSilence) {
		victim.Act = removeBit(victim.Act, playerSilence)
		victim.notify("You can speak again!")
		wiz.notify("SILENCE removed.")
	} else {
		victim.Act = setBit(victim.Act, playerSilence)
		victim.notify("You have been silenced!")
		wiz.notify("SILENCE set.")
	}
}

func doShutdow(wiz *mob, argument string) {
	wiz.notify("If you want to SHUTDOWN, spell it out.")
	return
}

func doShutdown(wiz *mob, argument string) {
	doEcho(wiz, fmt.Sprintf("Shutdown by %s", wiz.Name))
	gameServer.Up = false
}

func doSlayIncomplete(wiz *mob, argument string) {
	wiz.notify("If you want to SLAY someone, spell it out.")
	return
}

func doSlay(wiz *mob, argument string) {
	argument, arg1 := oneArgument(argument)

	if arg1 == "" {
		wiz.notify("Slay whom?")
		return
	}

	victim := getPlayerByName(arg1)
	if victim == nil || victim.Room != wiz.Room {
		wiz.notify("They aren't here.")
		return
	}

	if victim == wiz {
		wiz.notify("Kill yourself the hard way, coward.")
		return
	}

	if !victim.isNPC() && victim.Level >= wiz.Level {
		wiz.notify("You failed.")
		return
	}

	act("You slay $N in cold blood!", wiz, nil, victim, actToChar)
	act("$n slays you in cold blood!", wiz, nil, victim, actToVict)
	act("$n slays $N in cold blood!", wiz, nil, victim, actToNotVict)
	return
}

func doSlookup(wiz *mob, argument string) {
	argument, arg1 := oneArgument(argument)

	if arg1 == "" {
		wiz.notify("Slookup what?")
		return
	}

	if arg1 == "all" {
		for e := skillList.Front(); e != nil; e = e.Next() {
			sk := e.Value.(*skill)
			wiz.notify("ID: %4d Name: %s", sk.ID, sk.Name)
		}
	} else {
		sk := getSkillByName(arg1)

		if sk == nil {
			wiz.notify("No such skill or spell.")
			return
		}

		wiz.notify("ID: %4d Name: %s", sk.ID, sk.Name)
	}
}

func doSnoop(wiz *mob, argument string) {
	if len(argument) < 1 {
		wiz.notify("Snoop whom?")
		return
	}

	argument, arg1 := oneArgument(argument)
	victim := getPlayerByName(arg1)
	if victim == nil {
		wiz.notify("They aren't here.")
		return
	}

	if victim.client == nil {
		wiz.notify("They aren't connected.")
		return
	}

	if victim == wiz {
		wiz.notify("Cancelling all snoops.")
		for _, c := range gameServer.connections {
			if c.snoopBy == wiz.client {
				c.snoopBy = nil
			}
		}
		return
	}

	if victim.client.snoopBy != nil {
		wiz.notify("Busy already.")
		return
	}

	if victim.getTrust() >= wiz.getTrust() {
		wiz.notify("You failed.")
	}

	if wiz.client != nil {
		for _, c := range gameServer.connections {
			if c.mob == victim || c.original == victim {
				wiz.notify("No snoop loops.")
				return
			}
		}
	}

	victim.client.snoopBy = wiz.client
	wiz.notify("Ok.")
	return
}

func doSwitch(wiz *mob, argument string) {
	if len(argument) < 1 {
		wiz.notify("Switch into whom?")
		return
	}

	argument, arg1 := oneArgument(argument)
	victim := getPlayerByName(arg1)
	if wiz.client == nil {
		return
	}

	if wiz.client.original != nil {
		wiz.notify("You are already switched.")
		return
	}

	if victim == nil {
		wiz.notify("They aren't here.")
		return
	}

	if victim == wiz {
		wiz.notify("Ok.")
		return
	}

	if victim.client != nil {
		wiz.notify("Character is in use.")
		return
	}

	wiz.client.mob = victim
	wiz.client.original = wiz
	victim.client = wiz.client
	wiz.client = nil
	victim.notify("Ok.")
	return
}

func doTransfer(wiz *mob, argument string) {
	if len(argument) < 1 {
		wiz.notify("Transfer whom?")
	}

	where, whom := oneArgument(argument)

	if whom == "all" {
		for e := mobList.Front(); e != nil; e = e.Next() {
			m := e.Value.(*mob)
			if m.client != nil && m != wiz && m.Room != nil && wiz.canSee(m) {
				doTransfer(m, where)
			}
		}
		return
	}

	var location *room
	if where == "here" {
		location = wiz.Room
	} else {
		location := wiz.findLocation(where)
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

	interpret(victim, "look")
	wiz.notify("OK.")
}

func doTrust(wiz *mob, argument string) {
	if len(argument) < 1 {
		wiz.notify("Sytax: trust <char> <level>")
		return
	}

	argument, arg1 := oneArgument(argument)
	argument, arg2 := oneArgument(argument)

	name := arg1
	level, err := strconv.Atoi(arg2)
	if err != nil {
		wiz.notify("Syntax: trust <char> <level>")
		return
	}

	victim := getPlayerByName(name)
	if victim == nil {
		wiz.notify("They aren't here.")
		return
	}

	if level < 0 || level > 99 {
		wiz.notify("Level must be 0 (reset) or 1 to 99.")
		return
	}

	if level < wiz.getTrust() {
		wiz.notify("Limited to your trust level.")
		return
	}

	victim.Trust = level
	return
}

func doUsers(wiz *mob, argument string) {
	count := 0

	for _, c := range gameServer.connections {
		if c.mob != nil && wiz.canSee(c.mob) {
			count++

			name := "(none)"
			if c.original != nil {
				name = c.original.Name
			} else if c.mob != nil {
				name = c.mob.Name
			}
			wiz.notify("[%3d] %s@%s", count, name, c.conn.LocalAddr)
		}
	}

	suffix := ""
	if count != 1 {
		suffix = "s"
	}
	wiz.notify("%d user%", count, suffix)

}

func doWizhelp(player *mob, argument string) {
	var buf bytes.Buffer
	col := 0
	for e := helpList.Front(); e != nil; e = e.Next() {
		h := e.Value.(*help)
		if h.Level >= 90 && h.Level <= player.getTrust() {
			buf.Write([]byte(fmt.Sprintf("%-12s", h.Keyword)))
			col++
			if col%6 == 0 {
				output, _ := buf.ReadString('\n')
				player.notify(output)
			}
		}
	}

	if col%6 != 0 {
		player.notify("")
	}
}

func doWizlock(wiz *mob, argument string) {
	gameServer.Wizlock = !gameServer.Wizlock

	if gameServer.Wizlock {
		wiz.notify("Game wizlocked.")
	} else {
		wiz.notify("Game un-wizlocked.")
	}
}
