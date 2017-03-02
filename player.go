package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"strconv"
)

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
	Id       int    `json:"id"`
	Name     string `json:"username"`
	Password string `json:"password"`

	Inventory []Item `json:"inventory"`
	Room      int    `json:"current_room"`
	ExitVerb  string `json:"exit_verb"`

	Hitpoints int `json:"hitpoints"`
	Mana      int `json:"mana"`
	Movement  int `json:"movement"`
	Exp       int `json:"experience"`

	Level  int    `json:"level"`
	Class  string `json:"class"`
	Race   string `json:"race"`
	Gender string `json:"gender"`

	Stats     PlayerStats `json:"stats"`
	m_request *Connection
}

type PlayerStats struct {
	Strength     int `json:"strength"`
	Wisdom       int `json:"wisdom"`
	Intelligence int `json:"intelligence"`
	Dexterity    int `json:"dexterity"`
	Charisma     int `json:"charisma"`
	Constitution int `json:"constitution"`
}

func NewPlayer(c *Connection) *Player {
	p := &Player{}
	return p
}

func LoadPlayer(name string, password string) (*Player, error) {
	playerFile, err := ioutil.ReadFile("./players/" + name + ".json")
	if err != nil {
		return &Player{}, errors.New("Player not found.")
	}

	var player Player
	json.Unmarshal(playerFile, &player)

	if player.Password != password {
		return &Player{}, errors.New("Invalid password.")
	}

	return &player, nil
}

func (player *Player) Save() error {
	output, err := json.MarshalIndent(player, "", "    ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("./players/"+player.Name+".json", output, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (p *Player) AddItem(item Item) {
	p.Inventory = append(p.Inventory, item)
}

func (p *Player) RemoveItem(item Item) {
	for k, i := range p.Inventory {
		if i.Id == item.Id {
			p.Inventory = append(p.Inventory[:k], p.Inventory[k+1:]...)
			return
		}
	}
}

func (p *Player) SetRoom(room int) {
	p.Room = room
}

func (p Player) exitMessage(direction string) string {
	switch direction {
	case "up", "down":
		return "You " + p.ExitVerb + " " + direction + "." + newline
	default:
		return "You " + p.ExitVerb + " to the " + direction + "." + newline
	}
}

func (p Player) getInventory() map[string]int {
	inventory := make(map[string]int)
	for _, item := range p.Inventory {
		if _, ok := inventory[item.Name]; ok {
			inventory[item.Name]++
		} else {
			inventory[item.Name] = 1
		}
	}
	return inventory
}

func (p Player) getHitpoints() string {
	return strconv.Itoa(p.Hitpoints)
}
func (p Player) getMana() string {
	return strconv.Itoa(p.Mana)
}
func (p Player) getMovement() string {
	return strconv.Itoa(p.Movement)
}

func (p Player) ShowStatusBar() {
	p.m_request.SendString(white + "[" + p.getHitpoints() + reset + cyan + "hp " + white + p.getMana() + reset + cyan + "mana " + white + p.getMovement() + reset + cyan + "mv" + white + "] >> ")
}

func NewPlayerDatabase() *PlayerDatabase {
	return &PlayerDatabase{}
}
