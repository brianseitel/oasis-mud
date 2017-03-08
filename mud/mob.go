package mud

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"

	// "github.com/brianseitel/oasis-mud/helpers"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //
)

var pid int64

type mob struct {
	gorm.Model

	//Mob information
	Name      string `gorm:"name"`
	Inventory []item `json:"items",gorm:"many2many:mob_items;"`
	Room      room
	RoomID    int `json:"current_room"`
	ExitVerb  string

	Hitpoints    int `json:"hitpoints"`
	MaxHitpoints int `json:"max_hitpoints"`
	Mana         int `json:"mana"`
	MaxMana      int `json:"max_mana"`
	Movement     int `json:"movement"`
	MaxMovement  int `json:"max_movement"`

	Exp   int
	Level int

	Job    job
	Race   race
	Gender string

	Strength     int
	Wisdom       int
	Intelligence int
	Dexterity    int
	Charisma     int
	Constitution int

	client *connection
}

func (m *mob) move(e exit) {
	// if m.Room.Id <= 0 {
	//  m.Room = Registry.rooms[m.Room]
	// }

	// old_room := m.Room
	var newRoom room
	db.First(&newRoom, e.RoomID)

	// remove mob from old room list
	// delete(old_room.Mobs, m.pid)
	// for _, rm := range old_room.Mobs {
	// if rm.pid != m.pid {
	// rm.notify(fmt.Sprintf("%s leaves heading %s.\n", m.Name, e.Dir))
	// }
	// } //
	// Registry.rooms[old_room.Id] = old_room

	// add mob to new room list
	m.Room = newRoom
	// m.Room = newRoom.Id

	// m.Room.Mobs[m.pid] = *m
	// Registry.rooms[newRoom.Id] = m.Room

	// for _, rm := range newRoom.Mobs {
	// if rm.pid != m.pid {
	//  rm.notify(fmt.Sprintf("%s arrives in room %d.\n", m.Name, m.Room.ID))
	// }
	// }
}

func (m *mob) notify(message string) {
	if m.client != nil {
		m.client.SendString(message)
	}
}

func (m *mob) wander() {
	switch c := len(m.Room.Exits); c {
	case 0:
		return
	case 1:
		m.move(m.Room.Exits[0])
		return
	default:
		for {
			e := m.Room.Exits[dice().Intn(c)]
			// if m.lastRoom.ID != e.RoomID {
			m.move(e)
			// }
			return
		}
	}
}

func (m *mob) AddItem(item item) {
	m.Inventory = append(m.Inventory, item)
}

func (m *mob) RemoveItem(item item) {
	for k, i := range m.Inventory {
		if i.ID == item.ID {
			m.Inventory = append(m.Inventory[:k], m.Inventory[k+1:]...)
			return
		}
	}
}

// Return the inventory, grouped by item. Returns a map[string]int
// where map["Big Sword"]3 means the Big Sword has a qty of 3
func (m mob) getInventory() map[string]int {
	inventory := make(map[string]int)
	for _, item := range m.Inventory {
		if _, ok := inventory[item.Name]; ok {
			inventory[item.Name]++
		} else {
			inventory[item.Name] = 1
		}
	}
	return inventory
}

// Retrieves the player's hit points as a string
func (m mob) getHitpoints() string {
	return strconv.Itoa(m.Hitpoints)
}

// Retrieves the player's mana as a string
func (m mob) getMana() string {
	return strconv.Itoa(m.Mana)
}

// Retrieves the player's movement as a string
func (m mob) getMovement() string {
	return strconv.Itoa(m.Movement)
}

func findMob(i int) mob {
	var mob mob
	db.First(&mob, i)
	return mob
}

func newMobDatabase() {
	fmt.Println("Creating Mobs")
	mobFiles, _ := filepath.Glob("./data/mobs/*.json")

	for _, mobFile := range mobFiles {
		file, err := ioutil.ReadFile(mobFile)
		if err != nil {
			panic(err)
		}

		var list []mob
		err = json.Unmarshal(file, &list)
		if err != nil {
			panic(err)
		}

		for _, m := range list {
			var mobs mob
			db.First(&mobs, m.ID)
			if db.NewRecord(mobs) {
				fmt.Println("\tCreating mob " + m.Name + "!")
				db.Create(&m)
			} else {
				fmt.Println("\tSkipping mob " + m.Name + "!")
			}
		}
	}
}
