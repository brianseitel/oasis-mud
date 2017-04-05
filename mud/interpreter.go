package mud

import (
	"strings"
)

type cmd struct {
	Name     string
	Trust    int
	Position status
	Callback func(*mob, string)
}

func interpret(player *mob, argument string) {
	argument = strings.Trim(argument, " ")

	if len(argument) <= 0 {
		return
	}

	void := getRoom(0)
	if player.Room == void && player.WasInRoom != nil {
		player.Room = player.WasInRoom
		player.Room.Mobs = append(player.Room.Mobs, player)
		act("$n reappears out of nowhere.", player, nil, nil, actToRoom)
		player.WasInRoom = nil
	}

	// No hiding
	player.AffectedBy = removeBit(player.AffectedBy, affectHide)

	// Check freeze
	if !player.isNPC() && hasBit(player.Act, playerFreeze) {
		player.notify("You're totally frozen!")
		return
	}

	argument, command := oneArgument(argument)

	found := false
	trust := player.getTrust()

	var cx *cmd
	for e := commandList.Front(); e != nil; e = e.Next() {
		c := e.Value.(*cmd)
		if matchesSubject(c.Name, command) && c.Trust <= trust {
			found = true
			cx = c
			break
		}
	}
	// log and snoop TODO

	if !found {
		// check socials
		if !checkSocial(player, command, argument) {
			player.notify("Eh?")
		}
		return
	}

	// check positions
	if player.Status < cx.Position {
		switch player.Status {
		case dead:
			player.notify("You can't do that because you're DEAD.")
			break
		case mortal:
		case incapacitated:
			player.notify("You are too wounded to do that.")
			break
		case stunned:
			player.notify("You are too stunned to do that.")
			break
		case sleeping:
			player.notify("In your dreams, or what?")
			break
		case fighting:
			player.notify("No way, you are still fighting!")
			break

		}
		return
	}

	cx.Callback(player, argument)
	return
}

func oneArgument(argument string) (string, string) {
	parts := strings.Split(argument, " ")
	if len(parts) <= 1 {
		return "", argument
	}

	return strings.Join(parts[1:], " "), parts[0]
}
