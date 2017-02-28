package main

import "strconv"

const PLAYERITEMS = 16

const (
	DEFAULT_HITPOINTS = 100
	DEFAULT_MANA      = 0
	DEFAULT_MOVEMENT  = 100
	DEFAULT_EXP       = 1000
)

const (
	REGULAR = iota
	GOD
	ADMIN
)

type PlayerRank int8

type PlayerDatabase []Player

type Player struct {
	//Player information
	Entity
	inventory []int
	Room      int
	exitVerb  string

	hitpoints int
	mana      int
	movement  int
	exp       int

	m_request *Connection
}

func NewPlayer(c *Connection) Player {
	p := Player{
		m_request: c,

		hitpoints: DEFAULT_HITPOINTS,
		mana:      DEFAULT_MANA,
		movement:  DEFAULT_MOVEMENT,

		Room:     3,
		exitVerb: "jaunt",

		inventory: []int{1, 2},
	}
	return p
}

func (p *Player) SetRoom(room int) {
	p.Room = room
}

func (p Player) exitMessage(direction string) string {
	switch direction {
	case "up", "down":
		return "You " + p.exitVerb + " " + direction + "." + newline
	default:
		return "You " + p.exitVerb + " to the " + direction + "." + newline
	}
}

func (p Player) getHitpoints() string {
	return strconv.Itoa(p.hitpoints)
}
func (p Player) getMana() string {
	return strconv.Itoa(p.mana)
}
func (p Player) getMovement() string {
	return strconv.Itoa(p.movement)
}

func (p Player) ShowStatusBar() {
	p.m_request.SendString(white + "[" + p.getHitpoints() + reset + cyan + "hp " + white + p.getMana() + reset + cyan + "mana " + white + p.getMovement() + reset + cyan + "mv" + white + "] >> ")
}

func NewPlayerDatabase() *PlayerDatabase {
	return &PlayerDatabase{}
}
