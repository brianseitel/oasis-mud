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

func talkChannel(player *mob, args []string, channel int, verb string) {
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

func chatAuction(player *mob, args []string) {
	talkChannel(player, args, channelAuction, "auction")
}

func chatDefault(player *mob, args []string) {
	talkChannel(player, args, channelChat, "chat")
}

func chatMusic(player *mob, args []string) {
	talkChannel(player, args, channelMusic, "music")
}

func chatQuestion(player *mob, args []string) {
	talkChannel(player, args, channelQuestion, "question")
}

func chatAnswer(player *mob, args []string) {
	talkChannel(player, args, channelQuestion, "answer")
}

func chatImmtalk(player *mob, args []string) {
	talkChannel(player, args, channelImmtalk, "immtalk")
}

func say(player *mob, args []string) {
	if len(args) <= 1 {
		player.notify(fmt.Sprintf("Say what?%s", helpers.Newline))
		return
	}

	message := strings.Join(args, " ")
	player.notify(fmt.Sprintf("You say '%s'%s", message, helpers.Newline))
	player.Room.notify(fmt.Sprintf("%s says '%s'%s", strings.Title(player.Name), message, helpers.Newline), player)
}

func tell(player *mob, args []string) {
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

	victim.reply = player

}

func reply(player *mob, args []string) {
	if player.isNPC() || player.isSilenced() {
		player.notify(fmt.Sprintf("Your message didn't get through.%s", helpers.Newline))
		return
	}

	message := strings.Join(args[1:], " ")

	victim := player.reply
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

	victim.reply = player
}

func emote(player *mob, args []string) {

}
