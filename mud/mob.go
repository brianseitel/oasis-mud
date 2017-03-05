package mud

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"

	"github.com/brianseitel/oasis-mud/helpers"
)

var pid int64

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

type MobRank int8

type MobDatabase []Mob

type Mob struct {
	//Mob information
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

	Stats  MobStats `json:"stats"`
	client *Connection

	lastRoom Room
	isPlayer bool
	pid      int64
}

type MobStats struct {
	Strength     int `json:"strength"`
	Wisdom       int `json:"wisdom"`
	Intelligence int `json:"intelligence"`
	Dexterity    int `json:"dexterity"`
	Charisma     int `json:"charisma"`
	Constitution int `json:"constitution"`
}

// Loads a player and authenticates. If not found or not valid, returns error
// Otherwise, returns a Mob.
func LoadMob(name string, password string) (*Mob, error) {
	playerFile, err := ioutil.ReadFile("./data/players/" + name + ".json")
	if err != nil {
		return &Mob{}, errors.New("Mob not found.")
	}

	var player Mob
	json.Unmarshal(playerFile, &player)

	if player.Password != password {
		return &Mob{}, errors.New("Invalid password.")
	}

	return &player, nil
}

// Saves the player to disk
func (player *Mob) Save() error {
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

func (m *Mob) move(e Direction) {
	if m.room.Id <= 0 {
		m.room = Registry.rooms[m.Room]
	}

	old_room := m.room
	new_room := Registry.rooms[e.RoomId]

	// remove mob from old room list
	delete(old_room.Mobs, m.pid)
	for _, rm := range old_room.Mobs {
		if rm.pid != m.pid {
			rm.notify(fmt.Sprintf("%s leaves heading %s.\n", m.Name, e.Dir))
		}
	} //
	Registry.rooms[old_room.Id] = old_room

	// add mob to new room list
	m.room = new_room
	m.Room = new_room.Id
	m.lastRoom = old_room
	m.room.Mobs[m.pid] = *m
	Registry.rooms[new_room.Id] = m.room

	for _, rm := range new_room.Mobs {
		if rm.pid != m.pid {
			rm.notify(fmt.Sprintf("%s arrives in room %d.\n", m.Name, m.room.Id))
		}
	}
}

func (m *Mob) notify(message string) {
	if m.client != nil {
		m.client.SendString(message)
	}
}

func (m *Mob) wander() {
	if m.isPlayer {
		return
	}
	switch c := len(m.room.Exits); c {
	case 0:
		return
	case 1:
		m.move(m.room.Exits[0])
		return
	default:
		for {
			e := m.room.Exits[dice().Intn(c)]
			if m.lastRoom.Id != e.RoomId {
				m.move(e)
			}
			return
		}
	}
}

// Adds an item to the player's inventory
func (p *Mob) AddItem(item Item) {
	p.Inventory = append(p.Inventory, item)
}

// Removes an item from the player's inventory
func (p *Mob) RemoveItem(item Item) {
	for k, i := range p.Inventory {
		if i.Id == item.Id {
			p.Inventory = append(p.Inventory[:k], p.Inventory[k+1:]...)
			return
		}
	}
}

// Return the inventory, grouped by item. Returns a map[string]int
// where map["Big Sword"]3 means the Big Sword has a qty of 3
func (p Mob) getInventory() map[string]int {
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
func (p Mob) getHitpoints() string {
	return strconv.Itoa(p.Hitpoints)
}

// Retrieves the player's mana as a string
func (p Mob) getMana() string {
	return strconv.Itoa(p.Mana)
}

// Retrieves the player's movement as a string
func (p Mob) getMovement() string {
	return strconv.Itoa(p.Movement)
}

// Displays the player's status bar
func (p Mob) ShowStatusBar() {
	p.client.BufferData(helpers.White + "[" + p.getHitpoints() + helpers.Reset + helpers.Cyan + "hp")
	p.client.BufferData(helpers.White + p.getMana() + helpers.Reset + helpers.Cyan + "mana ")
	p.client.BufferData(helpers.White + p.getMovement() + helpers.Reset + helpers.Cyan + "mv" + helpers.White)
	p.client.BufferData("] >> ")
	p.client.SendBuffer()
}

// Finds an item in the item Database
// If not found, returns an empty item
func FindMob(i int) Mob {
	for _, v := range Registry.mobs {
		if v.Id == i {
			return v
		}
	}
	helpers.Dump("Shit!")

	return Mob{}
}

// Instantiates a new MobDatabase
func NewMobDatabase() []Mob {
	mobFiles, _ := filepath.Glob("./data/mobs/*.json")

	var mobs []Mob

	for _, mobFile := range mobFiles {
		file, err := ioutil.ReadFile(mobFile)
		if err != nil {
			panic(err)
		}

		var list []Mob
		json.Unmarshal(file, &list)

		for _, mob := range list {
			mob.isPlayer = false
			pid++
			mob.pid = pid
			mobs = append(mobs, mob)
		}
	}

	return mobs
}
