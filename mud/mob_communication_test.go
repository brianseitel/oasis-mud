package mud

import (
	"net"
	"testing"
)

func TestTalkChannel(t *testing.T) {
	player := resetTest()
	vic := mockMob("vic")
	mobList.PushBack(vic)

	s, c := net.Pipe()
	s.Close()
	vic.client = &connection{conn: c}

	// no message
	player.talkChannel("", channelChat, "chats")

	// no NPCs
	player.Playable = false
	player.talkChannel("holla", channelChat, "chats")

	// immtalk
	player.Playable = true
	player.talkChannel("holla", channelImmtalk, "")

	player.talkChannel("holla", channelChat, "holla")
}

func TestAllChannels(t *testing.T) {
	player := resetTest()

	doAuction(player, "holla")
	doChat(player, "holla")
	doMusic(player, "holla")
	doQuestion(player, "holla")
	doAnswer(player, "holla")
	doShout(player, "holla")
	doYell(player, "holla")
	doImmtalk(player, "holla")
}

func TestDoSay(t *testing.T) {
	player := resetTest()

	// no args
	doSay(player, "")

	doSay(player, "holla")
}

func TestDoTell(t *testing.T) {
	player := resetTest()
	vic := mockPlayer("vic")
	mobList.PushBack(vic)

	// no NPCs
	player.Playable = false
	doTell(player, "")

	// no args
	player.Playable = true
	doTell(player, "")

	// tell nobody
	doTell(player, "nobody")

	// tell player with no client
	vic.client = nil
	doTell(player, "vic hi")

	s, c := net.Pipe()
	s.Close()
	vic.client = &connection{conn: c}

	// tell player
	doTell(player, "vic hi")
}

func TestDoReply(t *testing.T) {
	player := resetTest()
	vic := mockPlayer("vic")
	mobList.PushBack(vic)

	// no NPCs
	player.Playable = false
	doReply(player, "")

	// no args
	player.Playable = true
	doReply(player, "")

	// tell nobody
	doReply(player, "nobody")

	player.replyTarget = vic
	// tell player with no client
	vic.client = nil
	doReply(player, "vic hi")

	s, c := net.Pipe()
	s.Close()
	vic.client = &connection{conn: c}

	// tell player
	doReply(player, "vic hi")
}

func TestDoEmote(t *testing.T) {
	player := resetTest()

	player.Playable = false
	doEmote(player, "")

	player.Playable = true
	doEmote(player, "")

	doEmote(player, "bites your face")

	doEmote(player, "bites your face.")
}

func TestDoGroupTell(t *testing.T) {
	player := resetTest()
	vic := mockPlayer("vic")
	mobList.PushBack(vic)

	doGroupTell(player, "")

	player.Act = playerNoTell
	doGroupTell(player, "holla")

	player.Act = 0
	vic.leader = player
	doGroupTell(player, "holla")
}
