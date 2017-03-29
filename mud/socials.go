package mud

import "fmt"
import "github.com/brianseitel/oasis-mud/helpers"

type social struct {
	Name        string `json:"name"`
	CharNoArg   string `json:"char_no_arg"`
	OthersNoArg string `json:"others_no_arg"`
	CharFound   string `json:"char_found"`
	OthersFound string `json:"others_found"`
	VictimFound string `json:"victim_found"`
	CharAuto    string `json:"char_auto"`
	OthersAuto  string `json:"others_auto"`
}

func checkSocial(player *mob, command string, args []string) bool {

	fmt.Println(command, args)
	found := false
	var action *social
	for e := socialList.Front(); e != nil; e = e.Next() {
		soc := e.Value.(*social)
		if helpers.MatchesSubject(soc.Name, command) {
			found = true
			action = soc
			break
		}
	}

	if !found {
		return false
	}

	social := action
	switch player.Status {
	case dead:
		player.notify("Lie still; you are DEAD")
		return true
	case incapacitated:
	case mortal:
		player.notify("You are too stunned to do that.")
		return true

	case sleeping:
		if social.Name != "snore" {
			player.notify("In your dreams, or what?")
			return true
		}
	}

	var victim *mob

	if len(args) < 1 {
		act(social.OthersNoArg, player, nil, victim, actToRoom)
		act(social.CharNoArg, player, nil, victim, actToChar)
		return true
	} else {
		for _, m := range player.Room.Mobs {
			if helpers.MatchesSubject(m.Name, args[0]) {
				victim = m
				break
			}
		}
	}

	if victim == nil {
		player.notify("They aren't here.")
	} else if victim == player {
		act(social.OthersAuto, player, nil, victim, actToRoom)
		act(social.CharAuto, player, nil, victim, actToChar)
	} else {
		act(social.OthersFound, player, nil, victim, actToNotVict)
		act(social.CharFound, player, nil, victim, actToChar)
		act(social.VictimFound, player, nil, victim, actToVict)

		if !player.isNPC() && victim.isNPC() && victim.isAwake() {
			switch dBits(4) {
			case 0:
				multiHit(victim, player, typeUndefined)
				break
			case 1:
			case 2:
			case 3:
			case 4:
			case 5:
			case 6:
			case 7:
			case 8:
				act(social.OthersFound, victim, nil, player, actToNotVict)
				act(social.CharFound, victim, nil, player, actToChar)
				act(social.VictimFound, victim, nil, player, actToVict)
				break
			case 9:
			case 10:
			case 11:
			case 12:
				act("$n slaps $N.", victim, nil, player, actToNotVict)
				act("You slap $N.", victim, nil, player, actToChar)
				act("$n slaps you.", victim, nil, player, actToVict)
				break
			}
		}
	}

	return true
}

func doSocial(player *mob, args string) {
	col := 0
	var buf string
	for e := socialList.Front(); e != nil; e = e.Next() {
		social := e.Value.(*social)
		buf = fmt.Sprintf("%-12s%s", social.Name, buf)
		col++
		if col%6 == 0 {
			player.notify(buf)
		}
	}

	if col%6 != 0 {
		player.notify("")
	}
}
