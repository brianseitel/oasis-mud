package mud

import (
	"fmt"
	"net"
	"testing"
)

func TestDoAdvance(t *testing.T) {
	loadBonuses()
	loadMobs()
	wiz := mockPlayer("Test")
	wiz.Trust = 90
	victim := mockPlayer("bar")
	mobList.PushBack(victim)

	victim = mockPlayer("baz")
	victim.Level = 50
	mobList.PushBack(victim)

	victim = mockPlayer("butt")
	victim.Level = 91
	mobList.PushBack(victim)

	// No argument
	doAdvance(wiz, "")

	// Not a valid level
	doAdvance(wiz, "foo what")

	// Not a valid target
	doAdvance(wiz, "foo 12")

	// Should work
	doAdvance(wiz, "bar 13")

	// Level out of range
	doAdvance(wiz, "baz 1000")

	// Demotion
	doAdvance(wiz, "baz 3")

	// Not enough trust
	doAdvance(wiz, "butt 92")
}

func TestDoAllow(t *testing.T) {
	wiz := mockPlayer("test")

	banList.PushBack("yourmom")
	banList.PushBack("somesite")
	// no argument
	doAllow(wiz, "")

	// allow
	doAllow(wiz, "somesite")

	found := false
	for e := banList.Front(); e != nil; e = e.Next() {
		b := e.Value.(string)

		if b == "somesite" {
			found = true
		}
	}

	if found {
		t.Error("Found site that should have been removed")
	}
}

func TestDoAt(t *testing.T) {
	loadRooms()
	wiz := mockPlayer("wiz")

	room1 := &room{Name: "private", RoomFlags: roomPrivate}
	room2 := &room{Name: "secret", RoomFlags: 0}
	// no args
	doAt(wiz, "")

	// not a real location
	doAt(wiz, "hahaha not real")

	// private room
	doAt(wiz, fmt.Sprintf("%s %s", room1.Name, "look"))

	// should work
	doAt(wiz, fmt.Sprintf("%s %s", room2.Name, "look"))
}

func TestDoBamfin(t *testing.T) {
	wiz := mockPlayer("wiz")

	// no argument
	doBamfin(wiz, "")

	// set
	doBamfin(wiz, "sup")

	// not if NPC
	wiz.Playable = false
	doBamfin(wiz, "hahaha")
}

func TestDoBamfout(t *testing.T) {
	wiz := mockPlayer("wiz")

	// no argument
	doBamfout(wiz, "")

	// set
	doBamfout(wiz, "sup")

	// not if NPC
	wiz.Playable = false
	doBamfout(wiz, "hahaha")
}

func TestDoBan(t *testing.T) {
	wiz := mockPlayer("wiz")

	// npc
	wiz.Playable = false
	doBan(wiz, "something")
	wiz.Playable = true

	// no arg
	doBan(wiz, "")

	// should work
	doBan(wiz, "something")

	// already banned
	doBan(wiz, "something")
}

func TestDoDeny(t *testing.T) {
	wiz := mockPlayer("wiz")
	wiz.Trust = 90
	victim := mockPlayer("jack")
	victim.Trust = 0
	mobList.PushBack(victim)

	// no argument
	doDeny(wiz, "")

	// not a player
	doDeny(wiz, "john")

	// not enough trust
	victim.Trust = 91
	doDeny(wiz, "jack")

	// works
	victim.Trust = 0
	doDeny(wiz, "jack")
}

func TestDoDisconnect(t *testing.T) {
	wiz := mockPlayer("wiz")
	vic := mockPlayer("vic")
	_, client := net.Pipe()
	_, client2 := net.Pipe()

	conn := &connection{conn: client}
	mobList.PushBack(vic)

	// no arg
	doDisconnect(wiz, "")

	// not a player
	doDisconnect(wiz, "flubber")

	// try it
	gameServer.connections = append(gameServer.connections, *conn)
	vic.client = &connection{conn: client}
	doDisconnect(wiz, "vic")

	vic.client = &connection{conn: client2}
	doDisconnect(wiz, "vic")

	// no connection
	vic.client = nil
	doDisconnect(wiz, "vic")

}

func TestDoEcho(t *testing.T) {
	mobList.Init()
	wiz := mockPlayer("wiz")
	mob := mockMob("foo")
	mobList.PushBack(mob)

	doEcho(wiz, "")
	doEcho(wiz, "hi")
}

func TestFindLocation(t *testing.T) {
	mobList.Init()
	wiz := mockPlayer("wiz")
	vic := mockPlayer("vic")

	s, c := net.Pipe()
	s.Close()
	vic.client = &connection{conn: c}

	mobList.PushBack(vic)

	wiz.findLocation("")
	wiz.findLocation("1")

	wiz.findLocation(vic.Name)

	wiz.findLocation("nothing")
}

func TestDoForce(t *testing.T) {
	mobList.Init()
	wiz := mockPlayer("wiz")
	vic := mockPlayer("vic")
	mobList.PushBack(vic)

	// no args
	doForce(wiz, "")

	// force all
	wiz.Trust = 90
	vic.Trust = 0
	doForce(wiz, "all score")

	// force someone who isn't here
	doForce(wiz, "jake score")

	// force yourself
	doForce(wiz, "wiz score")

	// no trust
	wiz.Trust = 0
	vic.Trust = 90
	doForce(wiz, "vic score")

	// should work
	wiz.Trust = 90
	vic.Trust = 0
	doForce(wiz, "vic score")
}

func TestDoFreeze(t *testing.T) {
	wiz := resetTest()
	vic := mockPlayer("vic")
	mobList.PushBack(vic)

	// no args
	doFreeze(wiz, "")

	// not a player
	doFreeze(wiz, "bob")

	// not enough trust
	wiz.Trust = 90
	vic.Trust = 95
	doFreeze(wiz, "vic")

	// should work: freeze
	vic.Trust = 0
	doFreeze(wiz, "vic")

	// unfreeze
	doFreeze(wiz, "vic")
}

func TestDoGoto(t *testing.T) {
	wiz := resetTest()
	vic := mockPlayer("vic")
	mobList.PushBack(vic)

	room := mockRoom()
	wiz.Room = room
	roomList.PushBack(room)

	// no args
	doGoto(wiz, "")

	// no location
	doGoto(wiz, "nowhere special")

	// private room
	room.RoomFlags = roomPrivate
	doGoto(wiz, "1")

	// stop fighting
	wiz.Fight = vic
	doGoto(wiz, "1")

	// become visible
	wiz.Act = playerWizInvis
	doGoto(wiz, "1")
}

func TestDoHolyLight(t *testing.T) {
	wiz := resetTest()

	wiz.Playable = false
	doHolylight(wiz, "")

	wiz.Playable = true
	doHolylight(wiz, "")
	doHolylight(wiz, "")
}

func TestDoInvis(t *testing.T) {
	wiz := resetTest()

	wiz.Playable = false
	doInvis(wiz, "")

	wiz.Playable = true
	doInvis(wiz, "")
	doInvis(wiz, "")
}

func TestDoLog(t *testing.T) {
	wiz := resetTest()
	vic := mockPlayer("vic")
	mobList.PushBack(vic)

	// no arg
	doLog(wiz, "")

	// can't do all
	doLog(wiz, "all")

	// not real player
	doLog(wiz, "foo")

	// no npcs
	vic.Playable = false
	doLog(wiz, "vic")

	// should log, then unlog
	vic.Playable = true
	doLog(wiz, "vic")
	doLog(wiz, "vic")
}

func TestDoMemory(t *testing.T) {
	wiz := resetTest()

	doMemory(wiz, "")
}

func TestDoMfind(t *testing.T) {
	wiz := resetTest()
	mob := mockMob("max")
	mobList.PushBack(mob)

	// no args
	doMfind(wiz, "")

	// find real mob
	doMfind(wiz, "max")

	// find all
	doMfind(wiz, "all")

	// don't find one
	doMfind(wiz, "tubber")
}
