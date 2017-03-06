package mud

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"

	// "github.com/brianseitel/oasis-mud/helpers"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //
)

var pid int64

type MobDatabase []Mob

type Mob struct {
	gorm.Model

	//Mob information
	Name      string `gorm:"name"`
	Inventory []Item `json:"items",gorm:"many2many:mob_items;"`
	Room      Room
	RoomId    int `json:"current_room"`
	ExitVerb  string

	Hitpoints    int `json:"hitpoints"`
	MaxHitpoints int `json:"max_hitpoints"`
	Mana         int `json:"mana"`
	MaxMana      int `json:"max_mana"`
	Movement     int `json:"movement"`
	MaxMovement  int `json:"max_movement"`

	Exp   int
	Level int

	Job    Job
	Race   Race
	Gender string

	// MobStats
	client *Connection `gorm:"-"`
}

type MobStats struct {
	Strength     int
	Wisdom       int
	Intelligence int
	Dexterity    int
	Charisma     int
	Constitution int
}

// Loads a player and authenticates. If not found or not valid, returns error
// Otherwise, returns a Mob.
func LoadMob(name string, password string) (*Player, error) {
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

func (m *Mob) move(e Exit) {
	// if m.Room.Id <= 0 {
	// 	m.Room = Registry.rooms[m.Room]
	// }

	// old_room := m.Room
	var new_room Room
	db.First(&new_room, e.RoomId)

	// remove mob from old room list
	// delete(old_room.Mobs, m.pid)
	// for _, rm := range old_room.Mobs {
	// if rm.pid != m.pid {
	// rm.notify(fmt.Sprintf("%s leaves heading %s.\n", m.Name, e.Dir))
	// }
	// } //
	// Registry.rooms[old_room.Id] = old_room

	// add mob to new room list
	m.Room = new_room
	// m.Room = new_room.Id

	// m.Room.Mobs[m.pid] = *m
	// Registry.rooms[new_room.Id] = m.Room

	// for _, rm := range new_room.Mobs {
	// if rm.pid != m.pid {
	// 	rm.notify(fmt.Sprintf("%s arrives in room %d.\n", m.Name, m.Room.ID))
	// }
	// }
}

func (m *Mob) notify(message string) {
	if m.client != nil {
		m.client.SendString(message)
	}
}

func (m *Mob) wander() {
	// if m.isPlayer {
	// 	return
	// }
	switch c := len(m.Room.Exits); c {
	case 0:
		return
	case 1:
		m.move(m.Room.Exits[0])
		return
	default:
		for {
			e := m.Room.Exits[dice().Intn(c)]
			// if m.lastRoom.ID != e.RoomId {
			m.move(e)
			// }
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
		if i.ID == item.ID {
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

// Finds an item in the item Database
// If not found, returns an empty item
func FindMob(i int) Mob {
	var mob Mob
	db.First(&mob, i)
	return mob
}

// Instantiates a new MobDatabase
func NewMobDatabase() {
	fmt.Println("Creating Mobs")
	mobFiles, _ := filepath.Glob("./data/mobs/*.json")

	for _, mobFile := range mobFiles {
		file, err := ioutil.ReadFile(mobFile)
		if err != nil {
			panic(err)
		}

		var list []Mob
		err = json.Unmarshal(file, &list)
		if err != nil {
			panic(err)
		}

		for _, mob := range list {
			var mobs Mob
			db.First(&mobs, mob.ID)
			if db.NewRecord(mobs) {
				fmt.Println("\tCreating mob " + mob.Name + "!")
				db.Create(&mob)
			} else {
				fmt.Println("\tSkipping mob " + mob.Name + "!")
			}
		}
	}
}
