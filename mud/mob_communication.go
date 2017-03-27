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

func (player *mob) talkChannel(args []string, channel int, verb string) {
	if len(args) <= 1 {
		player.notify("%s what?", strings.Title(verb))
		return
	}

	if player.isNPC() || player.isSilenced() {
		player.notify("You can't %s.", strings.Title(verb))
		return
	}

	message := strings.Join(args[1:], " ")
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

func (player *mob) chatAuction(args []string) {
	player.talkChannel(args, channelAuction, "auction")
}

func (player *mob) chatDefault(args []string) {
	player.talkChannel(args, channelChat, "chat")
}

func (player *mob) chatMusic(args []string) {
	player.talkChannel(args, channelMusic, "music")
}

func (player *mob) chatQuestion(args []string) {
	player.talkChannel(args, channelQuestion, "question")
}

func (player *mob) chatAnswer(args []string) {
	player.talkChannel(args, channelQuestion, "answer")
}

func (player *mob) chatImmtalk(args []string) {
	player.talkChannel(args, channelImmtalk, "immtalk")
}

func (player *mob) say(args []string) {
	if len(args) <= 1 {
		player.notify("Say what?")
		return
	}

	message := strings.Join(args[1:], " ")
	act("You say '$T'.", player, nil, message, actToChar)
	act("$n says '$T'.", player, nil, message, actToRoom)
}

func (player *mob) tell(args []string) {
	if player.isNPC() || player.isSilenced() {
		player.notify("Your message didn't get through.")
		return
	}

	if len(args) <= 1 {
		player.notify("Tell whom what?")
		return
	}

	name := args[1]
	message := strings.Join(args[1:], " ")

	victim := getPlayerByName(name)

	if victim == nil {
		player.notify("They aren't here.")
		return
	}

	if !victim.isNPC() && victim.client == nil {
		player.notify("They aren't here.")
		return
	}

	player.notify("You tell %s '%s'", strings.Title(name), message)
	victim.notify("%s tells you '%s'", strings.Title(player.Name), message)

	victim.replyTarget = player

}

func (player *mob) reply(args []string) {
	if player.isNPC() || player.isSilenced() {
		player.notify("Your message didn't get through.")
		return
	}

	message := strings.Join(args[1:], " ")

	victim := player.replyTarget
	if victim == nil {
		player.notify("They aren't here.")
		return
	}

	if !victim.isNPC() && victim.client == nil {
		player.notify("They aren't here.")
		return
	}

	player.notify("You tell %s '%s'", strings.Title(victim.Name), message)
	victim.notify("%s tells you '%s'", strings.Title(player.Name), message)

	victim.replyTarget = player
}

func (player *mob) emote(args []string) {
	if player.isNPC() || player.isSilenced() {
		player.notify("You can't show your emotions.")
		return
	}

	if len(args) <= 1 {
		player.notify("Emote what?")
		return
	}

	message := strings.Join(args[1:], " ")

	if !strings.HasSuffix(message, ".") {
		message += "."
	}

	player.notify("You %s", message)
	player.Room.notify(fmt.Sprintf("%s %s", strings.Title(player.Name), message), player)
}
