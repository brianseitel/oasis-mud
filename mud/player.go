package mud

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"strconv"

	"github.com/brianseitel/oasis-mud/helpers"
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
	room      Room
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

// Loads a player and authenticates. If not found or not valid, returns error
// Otherwise, returns a Player.
func LoadPlayer(name string, password string) (*Player, error) {
	playerFile, err := ioutil.ReadFile("./data/players/" + name + ".json")
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

// Saves the player to disk
func (player *Player) Save() error {
	output, err := json.MarshalIndent(player, "", "    ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("./data/players/"+player.Name+".json", output, 0644)
	if err != nil {
		return err
	}

	return nil
}

// Adds an item to the player's inventory
func (p *Player) AddItem(item Item) {
	p.Inventory = append(p.Inventory, item)
}

// Removes an item from the player's inventory
func (p *Player) RemoveItem(item Item) {
	for k, i := range p.Inventory {
		if i.Id == item.Id {
			p.Inventory = append(p.Inventory[:k], p.Inventory[k+1:]...)
			return
		}
	}
}

// Return the inventory, grouped by item. Returns a map[string]int
// where map["Big Sword"]3 means the Big Sword has a qty of 3
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

// Retrieves the player's hit points as a string
func (p Player) getHitpoints() string {
	return strconv.Itoa(p.Hitpoints)
}

// Retrieves the player's mana as a string
func (p Player) getMana() string {
	return strconv.Itoa(p.Mana)
}

// Retrieves the player's movement as a string
func (p Player) getMovement() string {
	return strconv.Itoa(p.Movement)
}

// Displays the player's status bar
func (p Player) ShowStatusBar() {
	p.m_request.BufferData(helpers.White + "[" + p.getHitpoints() + helpers.Reset + helpers.Cyan + "hp")
	p.m_request.BufferData(helpers.White + p.getMana() + helpers.Reset + helpers.Cyan + "mana ")
	p.m_request.BufferData(helpers.White + p.getMovement() + helpers.Reset + helpers.Cyan + "mv" + helpers.White)
	p.m_request.BufferData("] >> ")
	p.m_request.SendBuffer()
}

// Instantiates a new PlayerDatabase
func NewPlayerDatabase() *PlayerDatabase {
	return &PlayerDatabase{}
}
