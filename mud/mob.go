package mud

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"

	"github.com/brianseitel/oasis-mud/helpers"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //
)

var pid int64

type mob struct {
	gorm.Model

	//Mob information
	Name     string `json:"name" gorm:"name"`
	Password string `gorm:"password"`

	Inventory []item `gorm:"many2many:player_items;"`
	ItemIds   []int  `json:"items" gorm:"-"`
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

	Job    job  `json:"-"`
	JobID  int  `json:"job"`
	Race   race `json:"-"`
	RaceID int  `json:"race"`
	Gender string

	Strength     int
	Wisdom       int
	Intelligence int
	Dexterity    int
	Charisma     int
	Constitution int

	Status      status
	Identifiers string

	Fight   fight
	FightID uint

	client *connection
}

func (m mob) setFight(f *fight) {
	m.FightID = f.ID
	db.Save(&m)
}

func (m *mob) attack(target *mob, f *fight) {
	if target.Status != dead {
		damage := dice().Intn(m.damage()) + m.hit()
		target.takeDamage(damage)
		target.notify(fmt.Sprintf("%s attacks you for %d damage!%s", m.Name, damage, helpers.Newline))
		m.notify(fmt.Sprintf("You strike %s for %d damage!%s", target.Name, damage, helpers.Newline))
		if target.Status == dead {
			m.notify(fmt.Sprintf("You have KILLED %s to death!!%s", target.Name, helpers.Newline))
			m.Status = standing
			db.Save(&m)

			// whisk it away
			target.die()

			db.Delete(&f)
			m.ShowStatusBar()
		}
	}
}

func (m *mob) die() {
	// drop corpse in room
	corpse := &item{itemType: "corpse", Name: "A corpse of " + m.Name, Identifiers: "corpse," + m.Identifiers}
	var room room
	db.Find(&room, m.RoomID)
	db.Model(&room).Association("Items").Append(corpse).Error

	// whisk them away to hell
	m.RoomID = 0
	db.Save(&m)
}

func (m *mob) takeDamage(damage int) {
	m.Hitpoints -= damage
	if m.Hitpoints < -5 {
		m.Status = dead
		m.notify(helpers.Red + "You are DEAD!!!" + helpers.Reset)
	}

	db.Save(&m)
}

func (m *mob) damage() int {
	return m.Strength
}

func (m *mob) hit() int {
	return int(m.Dexterity / 3)
}

func (m mob) TNL() int {
	return (m.Level * 1000) - m.Exp
}

func (m *mob) move(e exit) {
	var oldRoom room
	db.First(&oldRoom, m.RoomID)

	var newRoom room
	db.First(&newRoom, e.RoomID)

	for _, rm := range oldRoom.Mobs {
		rm.notify(fmt.Sprintf("%s leaves heading %s\n", m.Name, e.Dir))
	}

	// add mob to new room list
	m.RoomID = int(newRoom.ID)

	for _, rm := range newRoom.Mobs {
		rm.notify(fmt.Sprintf("%s arrives in room %d.\n", m.Name, m.RoomID))
	}
}

func (m mob) ShowStatusBar() {
	m = getMob(m)
	if m.client != nil {
		m.client.BufferData(helpers.White + "[" + m.getHitpoints() + helpers.Reset + helpers.Cyan + "hp")
		m.client.BufferData(helpers.White + m.getMana() + helpers.Reset + helpers.Cyan + "mana ")
		m.client.BufferData(helpers.White + m.getMovement() + helpers.Reset + helpers.Cyan + "mv" + helpers.White)
		m.client.BufferData("] >> ")
		m.client.SendBuffer()
	}
}

func (m mob) getRoom() room {
	var (
		r room
	)

	db.Preload("Exits").Preload("Items").Preload("Mobs", "id != (?)", m.ID).First(&r, m.RoomID)

	return r
}

func (m *mob) notify(message string) {
	mob := getMob(*m)
	m = &mob
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

// Retrieves the mob's hit points as a string
func (m mob) getHitpoints() string {
	return strconv.Itoa(m.Hitpoints)
}

// Retrieves the mob's mana as a string
func (m mob) getMana() string {
	return strconv.Itoa(m.Mana)
}

// Retrieves the mob's movement as a string
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

func getMob(m mob) mob {
	db.Preload("Job").Preload("Race").Preload("Inventory").Preload("Room").First(&m)

	if m.client == nil {
		for _, conn := range getConnections() {
			if conn.mob.ID == m.ID {
				m.client = &conn
				return m
			}
		}
	}
	return m
}
