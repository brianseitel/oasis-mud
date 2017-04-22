package mud

import (
	"bytes"
	"fmt"
	"strings"
)

func (player *mob) changePosition(pos status) {
	p := player

	switch p.Status {
	case dead:
		p.notify("You can't do that because you are DEAD.")
		return
	case mortal:
		p.notify("You are too busy dying to do that.")
		return
	case incapacitated:
		p.notify("You are too busy bleeding out to do that.")
		return
	case stunned:
		p.notify("You are too stunned to do much of anything.")
		return
	case fighting:
		p.notify("You are too busy fighting to do that.")
		return
	}

	switch pos {
	case resting:
		p.Status = resting
		p.notify("You rest.")
		return
	case sitting:
		p.Status = sitting
		p.notify("You sit.")
		return
	case sleeping:
		p.Status = sleeping
		p.notify("You sleep.")
		return
	case standing:
		p.Status = standing
		p.notify("You stand up.")
		return
	default:
		p.notify("You can't do that.")
		return
	}

}

func doClose(player *mob, argument string) {
	argument, arg1 := oneArgument(argument)

	if arg1 == "" {
		player.notify("Close what?")
		return
	}

	var obj *item
	for _, i := range player.Inventory {
		if matchesSubject(i.Name, arg1) {
			obj = i
			break
		}
	}

	if obj == nil {
		for _, i := range player.Room.Items {
			if matchesSubject(i.Name, arg1) {
				obj = i
				break
			}
		}
	}

	if obj != nil {
		if obj.ItemType != itemContainer {
			player.notify("That isn't a container.")
			return
		}

		if obj.isClosed() {
			player.notify("It's already closed.")
			return
		}

		if !obj.isCloseable() {
			player.notify("You can't do that.")
			return
		}

		obj.ClosedFlags = setBit(obj.ClosedFlags, containerClosed)
		player.notify("Ok.")
		act("$n closes $p.", player, obj, nil, actToRoom)
		return
	}

	for _, e := range player.Room.Exits {
		if matchesSubject(e.Dir, arg1) {

			if e.isClosed() {
				player.notify("It's already closed.")
				return
			}

			e.Flags = setBit(e.Flags, exitClosed)
			act("$n closes the $d.", player, nil, e.Keyword, actToRoom)
			player.notify("Ok.")

			// close other side
			if e.Room != nil {
				for _, ex := range e.Room.Exits {
					if ex.Dir == reverseDirection(e.Dir) {
						ex.Flags = setBit(ex.Flags, exitClosed)
						for _, m := range e.Room.Mobs {
							act("The $d closes.", m, nil, ex.Keyword, actToChar)
						}
					}
				}
			}
			return
		}
	}

	player.notify("That isn't here.")
}

func doCommands(player *mob, argument string) {
	var buf bytes.Buffer
	col := 0
	for e := commandList.Front(); e != nil; e = e.Next() {
		c := e.Value.(*cmd)
		if c.Trust <= player.getTrust() {
			buf.Write([]byte(fmt.Sprintf("%-12s", c.Name)))
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

func doGroup(player *mob, argument string) {
	argument, arg1 := oneArgument(argument)

	if arg1 == "" {
		leader := player.leader
		if leader == nil {
			leader = player
		}

		player.notify("%s's group:", leader.Name)
		for e := mobList.Front(); e != nil; e = e.Next() {
			m := e.Value.(*mob)

			if m.Playable && m.client != nil {
				if m.leader == player || player.leader == m {
					job := "Mob"
					if !m.isNPC() {
						job = m.Job.Name
					}

					player.notify("[%2d %8s] %-16s %4d/%4d hp %4d/%4d mana %4d/%4d mv %5dxp", m.Level, job, pers(m, player), m.Hitpoints, m.MaxHitpoints, m.Mana, m.MaxMana, m.Movement, m.MaxMovement, m.Exp)
				}
			}
		}
		return
	}

	var victim *mob
	for _, m := range player.Room.Mobs {
		if matchesSubject(m.Name, arg1) {
			victim = m
			break
		}
	}
	if victim == nil || victim.Room.ID != player.Room.ID {
		player.notify("They aren't here.")
		return
	}

	if victim.isNPC() {
		player.notify("You can't group NPCs.")
		return
	}

	if player.master != nil || (player.leader != nil && player.leader != player) {
		player.notify("But you are following someone else!")
		return
	}

	if victim.master != nil && victim.master != player {
		act("$N isn't following you.", player, nil, victim, actToChar)
		return
	}

	if isSameGroup(player, victim) && player != victim {
		victim.leader = nil
		doFollow(victim, victim.Name)
		act("$n removes $N from $s group.", player, nil, victim, actToNotVict)
		act("$n removes you from $s group.", player, nil, victim, actToVict)
		act("You remove $N from your group.", player, nil, victim, actToChar)
		return
	}

	if player.Level-victim.Level < -5 || player.Level-victim.Level > 5 {
		act("$N cannot join $n's group.", player, nil, victim, actToNotVict)
		act("You cannot join $n's group.", player, nil, victim, actToVict)
		act("$N cannot join your group.", player, nil, victim, actToChar)
		return
	}

	victim.leader = player

	act("$N joins $n's group.", player, nil, victim, actToNotVict)
	act("You join $n's group.", player, nil, victim, actToVict)
	act("$N joins your group.", player, nil, victim, actToChar)
	return
}

func doHide(player *mob, argument string) {
	player.notify("You attempt to hide.")

	if hasBit(player.AffectedBy, affectHide) {
		removeBit(player.AffectedBy, affectHide)
	}

	hide := player.skill("hide")

	if player.isNPC() || dice().Intn(100) < int(hide.Level) {
		player.AffectedBy = setBit(player.AffectedBy, affectHide)
	}
	return
}

func doFollow(player *mob, argument string) {
	if len(argument) < 1 {
		player.notify("Follow whom?")
		return
	}

	argument, arg1 := oneArgument(argument)
	var victim *mob
	for _, m := range player.Room.Mobs {
		if matchesSubject(m.Name, arg1) {
			victim = m
			break
		}
	}

	if victim == nil {
		player.notify("They aren't here.")
		return
	}

	if victim == player {
		if player.master == nil {
			player.notify("You already follow yourself.")
			return
		}

		player.stopFollower()
		return
	}

	if player.Level-victim.Level < -5 || player.Level-victim.Level > 5 {
		player.notify("You are not of the right caliber to follow.")
		return
	}

	if player.master != nil {
		player.stopFollower()
	}

	player.addFollower(victim)
	return
}

func doLock(player *mob, argument string) {
	if len(argument) < 1 {
		player.notify("Lock what?")
		return
	}

	// TODO: check items
	argument, arg1 := oneArgument(argument)
	exit := player.Room.findExit(arg1)
	if exit != nil {
		if !exit.isClosed() {
			player.notify("It's not closed, dummy.")
			return
		}

		if exit.Key < 0 {
			player.notify("It can't be locked.")
			return
		}

		if !player.hasKey(exit.Key) {
			player.notify("You don't have the key.")
			return
		}

		if exit.isLocked() {
			player.notify("It's already locked, dummy.")
			return
		}

		exit.Flags = setBit(exit.Flags, exitLocked)
		player.notify("*click*")
		act("$n unlocks the $d.", player, nil, exit.Keyword, actToRoom)

		/* lock the other side */
		if exit.Room != nil {
			reverseExit := exit.Room.findExit(reverseDirection(exit.Dir))
			if reverseExit != nil && reverseExit.Room == player.Room {
				reverseExit.Flags = setBit(reverseExit.Flags, exitLocked)
			}
		}
	}
	return
}

func doMove(player *mob, d string) {
	if player.Status != standing {
		switch player.Status {
		case fighting:
			player.notify("You can't move while fighting!")
			break
		}
		return
	}

	for _, e := range player.Room.Exits {
		if e.Dir == d {
			player.move(e)
			doLook(player, "")
			return
		}
	}
	player.notify("Alas, you cannot go that way.")
}

func doPick(player *mob, argument string) {
	if len(argument) < 1 {
		player.notify("Pick what?")
		return
	}

	pick := player.skill("pick")
	if pick == nil {
		player.notify("You don't know how to do pick locks.")
		return
	}

	wait(player, pick.Skill.Beats)

	/* look for guards */
	for _, m := range player.Room.Mobs {
		if m.isNPC() && m.isAwake() && player.Level+5 < m.Level {
			act("$N is standing too close. It's too risky.", player, nil, m, actToChar)
			return
		}
	}

	if !player.isNPC() && dice().Intn(100) > int(pick.Level) {
		player.notify("You failed.")
		return
	}

	/* check items */
	argument, arg1 := oneArgument(argument)
	exit := player.Room.findExit(arg1)
	if exit != nil {
		if !exit.isClosed() {
			player.notify("It's not closed, dummy.")
			return
		}

		if exit.Key < 0 {
			player.notify("It can't be unlocked.")
			return
		}

		if !exit.isLocked() {
			player.notify("It's already unlocked, dummy.")
			return
		}

		removeBit(exit.Flags, exitLocked)
		player.notify("*click*")
		act("$n unlocks the $d.", player, nil, exit.Keyword, actToRoom)

		/* unlock the other side */

		if exit.Room != nil {
			reverseExit := exit.Room.findExit(reverseDirection(exit.Dir))
			if reverseExit != nil && reverseExit.Room == player.Room {
				removeBit(reverseExit.Flags, exitLocked)
			}
		}
	}
}

func doPractice(player *mob, argument string) {
	if player.isNPC() {
		return
	}

	if player.Level < 3 {
		player.notify("You must be at least Level 3 before you can practice. Go train instead!")
		return
	}

	argument, arg1 := oneArgument(argument)
	if arg1 == "" {
		col := 0
		var buf bytes.Buffer

		var allSkills []*skill
		for e := skillList.Front(); e != nil; e = e.Next() {
			skill := e.Value.(*skill)
			maxLevel := 99
			switch player.Job.Name {
			case "Warrior":
				maxLevel = skill.Levels.Warrior
				break
			case "Mage":
				maxLevel = skill.Levels.Mage
				break
			case "Cleric":
				maxLevel = skill.Levels.Cleric
				break
			case "Thief":
				maxLevel = skill.Levels.Thief
				break
			case "Ranger":
				maxLevel = skill.Levels.Ranger
				break
			case "Bard":
				maxLevel = skill.Levels.Bard
				break
			}
			if player.Level >= maxLevel {
				allSkills = append(allSkills, skill)
			}
		}

		for _, skill := range allSkills {
			pSkill := player.skill(skill.Name)

			skillLevel := 0
			if pSkill != nil {
				skillLevel = pSkill.Level
			}

			buf.Write([]byte(fmt.Sprintf("%-12s %3d  ", skill.Name, skillLevel)))
			col++
			if col%4 == 0 {
				str := buf.String()
				player.notify(str)
			}
		}

		if buf.Len() > 0 {
			str := buf.String()
			player.notify(str)
		}

		player.notify("You have %d practices remaining.", player.Practices)
	} else {
		var adept int

		if !player.isAwake() {
			player.notify("In your dreams, or what?")
			return
		}

		var trainer *mob
		for _, mob := range player.Room.Mobs {
			if mob.isNPC() && hasBit(mob.Act, actPractice) {
				trainer = mob
				break
			}
		}

		if trainer == nil {
			player.notify("You can't do that here.")
			return
		}

		if player.Practices <= 0 {
			player.notify("You have no practices remaining.")
			return
		}

		skill := getSkillByName(arg1)

		if skill == nil {
			player.notify("That's not something you can learn.")
			return
		}
		pSkill := player.skill(skill.Name)

		maxLevel := 99
		switch player.Job.Name {
		case "Warrior":
			maxLevel = skill.Levels.Warrior
			break
		case "Mage":
			maxLevel = skill.Levels.Mage
			break
		case "Cleric":
			maxLevel = skill.Levels.Cleric
			break
		case "Thief":
			maxLevel = skill.Levels.Thief
			break
		case "Ranger":
			maxLevel = skill.Levels.Ranger
			break
		case "Bard":
			maxLevel = skill.Levels.Bard
			break
		}

		if skill == nil || (!player.isNPC() && player.Level < maxLevel) {
			player.notify("You can't practice that.")
			return
		}

		adept = player.Job.SkillAdept

		pSkill = player.skill(arg1)

		if pSkill == nil {
			pSkill = &mobSkill{Skill: skill, Level: 0}
			player.Skills = append(player.Skills, pSkill)
		}

		if pSkill.Level >= adept {
			player.notify("You've already mastered %s.", pSkill.Skill.Name)
			return
		}

		player.Practices--

		pSkill.Level += bonusTableIntelligence[player.currentIntelligence()].learn
		if pSkill.Level > player.Job.SkillAdept {
			pSkill.Level = player.Job.SkillAdept
		}

		if pSkill.Level < adept {
			act("You practice $T.", player, nil, pSkill.Skill.Name, actToChar)
			act("$n practices $T.", player, nil, pSkill.Skill.Name, actToRoom)
		} else {
			pSkill.Level = adept
			act("You have now mastered $T.", player, nil, pSkill.Skill.Name, actToChar)
			act("$n has now mastered $T.", player, nil, pSkill.Skill.Name, actToRoom)
		}
	}
	return
}

func doOpen(player *mob, argument string) {
	argument, arg1 := oneArgument(argument)

	if arg1 == "" {
		player.notify("Open what?")
		return
	}

	var obj *item
	for _, i := range player.Inventory {
		if matchesSubject(i.Name, arg1) {
			obj = i
			break
		}
	}

	if obj == nil {
		for _, i := range player.Room.Items {
			if matchesSubject(i.Name, arg1) {
				obj = i
				break
			}
		}
	}

	if obj != nil {
		if obj.ItemType != itemContainer {
			player.notify("That isn't a container.")
			return
		}

		if !obj.isClosed() {
			player.notify("It's already open.")
			return
		}

		if !obj.isCloseable() {
			player.notify("You can't do that.")
			return
		}

		obj.ClosedFlags = removeBit(obj.ClosedFlags, containerClosed)
		player.notify("Ok.")
		act("$n closes $p.", player, obj, nil, actToRoom)
		return
	}

	for _, e := range player.Room.Exits {
		if matchesSubject(e.Dir, arg1) {

			if e.isLocked() {
				player.notify("This door is locked.")
				return
			}

			if !e.isClosed() {
				player.notify("It's already open.")
				return
			}

			e.Flags = removeBit(e.Flags, exitClosed)
			act("$n open the $d.", player, nil, e.Keyword, actToRoom)
			player.notify("Ok.")

			// close other side
			if e.Room != nil {
				for _, ex := range e.Room.Exits {
					if ex.Dir == reverseDirection(e.Dir) {
						ex.Flags = removeBit(ex.Flags, exitClosed)
						for _, m := range e.Room.Mobs {
							act("The $d opens.", m, nil, ex.Keyword, actToChar)
						}
					}
				}
			}
			return
		}
	}

	player.notify("That isn't here.")
}

func doQui(player *mob, argument string) {
	player.notify("If you want to quit, spell it out.")
	return
}

func doQuit(player *mob, argument string) {
	extractMob(player, true)
	if player.Status == fighting {
		player.notify("You can't quit now. You're fighting!")
	} else {
		player.notify("Seeya!")
		player.client.end()
	}
}

func doRecall(player *mob, argument string) {
	if len(argument) < 1 {
		room := getRoom(player.RecallRoomID)
		player.Room = room
		doLook(player, "")
		return
	}

	argument, arg1 := oneArgument(argument)
	if arg1 == "set" {
		player.RecallRoomID = player.Room.ID
		player.notify("Recall set!")
		return
	}

	player.notify("Recall what?")
}

func doRest(player *mob, argument string) {
	switch player.Status {
	case resting:
		player.notify("You are already resting.")
		break

	case sleeping:
		player.notify("You sit up and rest.")
		act("$n wakes up and rests.", player, nil, nil, actToRoom)
		break

	case standing:
		player.notify("You rest.")
		act("$n rests.", player, nil, nil, actToRoom)
		player.Status = resting
		break

	case fighting:
		player.notify("You're too busy fighting!")
		break
	}
	return
}

func doSave(player *mob, argument string) {
	saveCharacter(player)
}

func doSleep(player *mob, argument string) {
	switch player.Status {
	case sleeping:
		player.notify("You are already sleeping.")
		return

	case resting:
	case standing:
		player.notify("You sleep.")
		act("$n sleeps.", player, nil, nil, actToRoom)
		player.Status = sleeping
		break

	case fighting:
		player.notify("You're too busy fighting!")
		break
	}
	return
}

func doSneak(player *mob, argument string) {
	sneak := player.skill("sneak")
	if sneak == nil {
		player.notify("You don't know how to sneak.")
		return
	}
	for _, affect := range player.Affects {
		if affect.affectType == sneak {
			player.removeAffect(affect)
			break
		}
	}

	var af affect
	af.affectType = sneak
	af.duration = max(5, player.Level)
	af.modifier = 0
	af.location = applyNone
	af.bitVector = affectSneak

	player.addAffect(&af)
	player.notify("You are now sneaking.")
}

func doSteal(player *mob, argument string) {
	steal := player.skill("steal")
	if steal == nil {
		player.notify("You don't know how to steal.")
		return
	}

	if len(argument) < 1 {
		player.notify("Steal what from whom?")
		return
	}
	argument, arg1 := oneArgument(argument)
	argument, arg2 := oneArgument(argument)

	var victim *mob
	for _, mob := range player.Room.Mobs {
		if matchesSubject(mob.Name, arg2) {
			victim = mob
			break
		}
	}

	if victim == nil {
		player.notify("They aren't here.")
		return
	}

	if victim == player {
		player.notify("That's pointless.")
		return
	}

	wait(player, steal.Skill.Beats)
	percent := dice().Intn(100)
	if victim.isAwake() {
		percent += 10
	} else {
		percent -= 50
	}

	if player.Level+5 < victim.Level || victim.Status == fighting || !victim.isNPC() || (!player.isNPC() && percent > int(steal.Level)) {
		// failed //
		player.notify("Oops.")
		act("$n tried to steal from you.", player, nil, victim, actToVict)
		act("$n tried to steal from $N.", player, nil, victim, actToNotVict)

		doShout(player, fmt.Sprintf("%s is a bloody thief!", player.Name))

		if !player.isNPC() {
			if victim.isNPC() {
				multiHit(victim, player, typeUndefined)
			} else {
				player.Act = setBit(player.Act, playerThief)
			}
		}

		return
	}

	if strings.HasPrefix(arg1, "coin") || strings.HasPrefix(arg1, "gold") {
		// steal money!
		amount := int(victim.Gold) * dice().Intn(10) / 100
		if amount <= 0 {
			player.notify("You couldn't get any gold.")
			return
		}

		player.Gold += amount
		victim.Gold -= amount
		player.notify("Bingo! You stole %d gold coins!", amount)
		return
	}

	var item *item
	for _, i := range victim.Inventory {
		if matchesSubject(i.Name, arg1) {
			item = i
			break
		}
	}

	if item == nil {
		player.notify("You can't find it.")
		return
	}

	if !player.canDropItem(item) {
		player.notify("You can't pry it away.")
		return
	}

	if player.Carrying+1 > player.CarryMax {
		player.notify("You have your hands full.")
		return
	}

	if player.CarryWeight+item.Weight > player.CarryWeightMax {
		player.notify("You can't carry that much weight.")
		return
	}

	for j, i := range victim.Inventory {
		if item == i {
			victim.Inventory, player.Inventory = transferItem(j, victim.Inventory, player.Inventory)
			break
		}
	}

	player.notify("Ok.")
	return
}

func doTrain(player *mob, argument string) {
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

	argument, arg1 := oneArgument(argument)
	if arg1 == "" {
		player.notify("You have %d practice sessions.", player.Practices)
		return
	}

	var cost int

	costmap := []int{5, 6, 7, 9, 12, 13, 15, 18, 21, 24, 27, 32, 40}

	var playerAbility int
	var playerOutput string

	if strings.HasPrefix(arg1, "str") {
		playerAbility = player.Attributes.Strength
		playerOutput = "strength"
	} else if strings.HasPrefix(arg1, "int") {
		playerAbility = player.Attributes.Intelligence
		playerOutput = "intelligence"
	} else if strings.HasPrefix(arg1, "wis") {
		playerAbility = player.Attributes.Wisdom
		playerOutput = "wisdom"
	} else if strings.HasPrefix(arg1, "dex") {
		playerAbility = player.Attributes.Dexterity
		playerOutput = "dexterity"
	} else if strings.HasPrefix(arg1, "cha") {
		playerAbility = player.Attributes.Charisma
		playerOutput = "charisma"
	} else if strings.HasPrefix(arg1, "con") {
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

func doUnlock(player *mob, argument string) {
	if len(argument) < 1 {
		player.notify("Unlock what?")
		return
	}

	// check items
	argument, arg1 := oneArgument(argument)
	exit := player.Room.findExit(arg1)
	if exit != nil {
		if !exit.isClosed() {
			player.notify("It's not closed, dummy.")
			return
		}

		if exit.Key < 0 {
			player.notify("It can't be unlocked.")
			return
		}

		if !player.hasKey(exit.Key) {
			player.notify("You don't have the key.")
			return
		}

		if !exit.isLocked() {
			player.notify("It's already unlocked, dummy.")
			return
		}

		exit.Flags = removeBit(exit.Flags, exitLocked)
		player.notify("*click*")
		act("$n unlocks the $d.", player, nil, exit.Keyword, actToRoom)

		/* unlock the other side */

		if exit.Room != nil {
			reverseExit := exit.Room.findExit(reverseDirection(exit.Dir))
			if reverseExit != nil && reverseExit.Room == player.Room {
				reverseExit.Flags = removeBit(reverseExit.Flags, exitLocked)
			}
		}
	}
	return
}

func doVisible(player *mob, argument string) {
	player.stripAffect("invis")
	player.stripAffect("mass_invis")
	player.stripAffect("sneak")
	removeBit(player.AffectedBy, affectHide)
	removeBit(player.AffectedBy, affectInvisible)
	removeBit(player.AffectedBy, affectSneak)
	player.notify("Ok.")
}

func doWake(player *mob, argument string) {
	switch player.Status {
	case standing:
		player.notify("You're already awake!")
		break

	case sleeping:
		player.notify("You wake up, stretch, and climb to your feet.")
		act("$n wakes up.", player, nil, nil, actToRoom)
		player.Status = standing
		break

	case resting:
		player.notify("You wake up and rest.")
		act("$n wakes up and rests.", player, nil, nil, actToRoom)
		player.Status = standing
		break

	case fighting:
		player.notify("You're too busy fighting!")
		break
	}
	return
}
