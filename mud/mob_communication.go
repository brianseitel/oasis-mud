package mud

import (
	"fmt"

	"strings"
)

const (
	channelAuction = 1 << iota
	channelChat
	channelHacker
	channelImmtalk
	channelMusic
	channelQuestion
	channelShout
	channelYell
)

func (player *mob) talkChannel(message string, channel int, verb string) {
	if len(message) <= 1 {
		player.notify("%s what?", strings.Title(verb))
		return
	}

	if player.isNPC() || player.isSilenced() {
		player.notify("You can't %s.", strings.Title(verb))
		return
	}

	switch channel {
	case channelImmtalk:
		player.notify("You: %s", message)
		break
	default:
		player.notify("You %s '%s'", verb, message)
	}

	for e := mobList.Front(); e != nil; e = e.Next() {
		mob := e.Value.(*mob)
		if mob.client != nil {
			if channel == channelImmtalk && !mob.isImmortal() {
				continue
			}

			if mob != player {
				mob.notify("%s %ss '%s'", player.Name, verb, message)
			}
		}
	}
}

func doAuction(player *mob, argument string) {
	player.talkChannel(argument, channelAuction, "auction")
}

func doChat(player *mob, argument string) {
	player.talkChannel(argument, channelChat, "chat")
}

func doMusic(player *mob, argument string) {
	player.talkChannel(argument, channelMusic, "music")
}

func doQuestion(player *mob, argument string) {
	player.talkChannel(argument, channelQuestion, "question")
}

func doAnswer(player *mob, argument string) {
	player.talkChannel(argument, channelQuestion, "answer")
}

func doImmtalk(player *mob, argument string) {
	player.talkChannel(argument, channelImmtalk, "immtalk")
}

func doSay(player *mob, argument string) {
	if len(argument) <= 1 {
		player.notify("Say what?")
		return
	}

	act("You say '$T'.", player, nil, argument, actToChar)
	act("$n says '$T'.", player, nil, argument, actToRoom)
}

func doTell(player *mob, argument string) {
	if player.isNPC() || player.isSilenced() {
		player.notify("Your message didn't get through.")
		return
	}

	if len(argument) <= 1 {
		player.notify("Tell whom what?")
		return
	}

	argument, name := oneArgument(argument)
	victim := getPlayerByName(name)

	if victim == nil {
		player.notify("They aren't here.")
		return
	}

	if !victim.isNPC() && victim.client == nil {
		player.notify("They aren't here.")
		return
	}

	player.notify("You tell %s '%s'", strings.Title(name), argument)
	victim.notify("%s tells you '%s'", strings.Title(player.Name), argument)

	victim.replyTarget = player

}

func doReply(player *mob, argument string) {
	if player.isNPC() || player.isSilenced() {
		player.notify("Your message didn't get through.")
		return
	}

	victim := player.replyTarget
	if victim == nil {
		player.notify("They aren't here.")
		return
	}

	if !victim.isNPC() && victim.client == nil {
		player.notify("They aren't here.")
		return
	}

	player.notify("You tell %s '%s'", strings.Title(victim.Name), argument)
	victim.notify("%s tells you '%s'", strings.Title(player.Name), argument)

	victim.replyTarget = player
}

func doEmote(player *mob, argument string) {
	if player.isNPC() || player.isSilenced() {
		player.notify("You can't show your emotions.")
		return
	}

	if len(argument) <= 1 {
		player.notify("Emote what?")
		return
	}

	if !strings.HasSuffix(argument, ".") {
		argument += "."
	}

	player.notify("You %s", argument)
	player.Room.notify(fmt.Sprintf("%s %s", strings.Title(player.Name), argument), player)
}
