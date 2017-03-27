package mud

import (
	"fmt"

	"strings"

	"github.com/brianseitel/oasis-mud/helpers"
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
		player.notify(fmt.Sprintf("%s what?%s", strings.Title(verb), helpers.Newline))
		return
	}

	if player.isNPC() || player.isSilenced() {
		player.notify(fmt.Sprintf("You can't %s.%s", strings.Title(verb), helpers.Newline))
		return
	}

	message := strings.Join(args[1:], " ")
	switch channel {
	case channelImmtalk:
		player.notify(fmt.Sprintf("You: %s%s", message, helpers.Newline))
		break
	default:
		player.notify(fmt.Sprintf("You %s '%s'%s", verb, message, helpers.Newline))
	}

	for e := mobList.Front(); e != nil; e = e.Next() {
		mob := e.Value.(*mob)
		if mob.client != nil {
			if channel == channelImmtalk && !mob.isImmortal() {
				continue
			}

			if mob != player {
				mob.notify(fmt.Sprintf("%s %ss '%s'%s", player.Name, verb, message, helpers.Newline))
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
		player.notify(fmt.Sprintf("Say what?%s", helpers.Newline))
		return
	}

	message := strings.Join(args, " ")
	player.notify(fmt.Sprintf("You say '%s'%s", message, helpers.Newline))
	player.Room.notify(fmt.Sprintf("%s says '%s'%s", strings.Title(player.Name), message, helpers.Newline), player)
}

func (player *mob) tell(args []string) {
	if player.isNPC() || player.isSilenced() {
		player.notify(fmt.Sprintf("Your message didn't get through.%s", helpers.Newline))
		return
	}

	if len(args) <= 1 {
		player.notify(fmt.Sprintf("Tell whom what?%s", helpers.Newline))
		return
	}

	name := args[1]
	message := strings.Join(args[1:], " ")

	victim := getPlayerByName(name)

	if victim == nil {
		player.notify(fmt.Sprintf("They aren't here.%s", helpers.Newline))
		return
	}

	if !victim.isNPC() && victim.client == nil {
		player.notify(fmt.Sprintf("They aren't here.%s", helpers.Newline))
		return
	}

	player.notify(fmt.Sprintf("You tell %s '%s'%s", strings.Title(name), message, helpers.Newline))
	victim.notify(fmt.Sprintf("%s tells you '%s'%s", strings.Title(player.Name), message, helpers.Newline))

	victim.replyTarget = player

}

func (player *mob) reply(args []string) {
	if player.isNPC() || player.isSilenced() {
		player.notify(fmt.Sprintf("Your message didn't get through.%s", helpers.Newline))
		return
	}

	message := strings.Join(args[1:], " ")

	victim := player.replyTarget
	if victim == nil {
		player.notify(fmt.Sprintf("They aren't here.%s", helpers.Newline))
		return
	}

	if !victim.isNPC() && victim.client == nil {
		player.notify(fmt.Sprintf("They aren't here.%s", helpers.Newline))
		return
	}

	player.notify(fmt.Sprintf("You tell %s '%s'%s", strings.Title(victim.Name), message, helpers.Newline))
	victim.notify(fmt.Sprintf("%s tells you '%s'%s", strings.Title(player.Name), message, helpers.Newline))

	victim.replyTarget = player
}

func (player *mob) emote(args []string) {
	if player.isNPC() || player.isSilenced() {
		player.notify(fmt.Sprintf("You can't show your emotions.%s", helpers.Newline))
		return
	}

	if len(args) <= 1 {
		player.notify(fmt.Sprintf("Emote what?%s", helpers.Newline))
		return
	}

	message := strings.Join(args[1:], " ")

	if !strings.HasSuffix(message, ".") {
		message += "."
	}

	player.notify(fmt.Sprintf("You %s%s", message, helpers.Newline))
	player.Room.notify(fmt.Sprintf("%s %s%s", strings.Title(player.Name), message, helpers.Newline), player)
}
