package mud

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"

	"container/list"

	"github.com/brianseitel/oasis-mud/helpers"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //
)

var (
	mobList *list.List
)

type mob struct {
	gorm.Model

	//Mob information
	Name        string `json:"name" gorm:"name"`
	Password    string `gorm:"password"`
	Description string `gorm:"type:text"`

	Inventory []*item `gorm:"many2many:player_items;"`
	Equipped  []*item `gorm:"many2many:player_equipped;ForeignKey:item"`
	ItemIds   []int   `json:"items" gorm:"-"`
	Room      *room
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

	Fight   *fight
	FightID uint

	Playable bool
	client   *connection
}

func (m *mob) addItem(item *item) {
	m.Inventory = append(m.Inventory, item)
}

func (m mob) setFight(f *fight) {
	m.Fight = f
}

func (m *mob) attack(target *mob, f *fight) {
	if target.Status != dead {
		damage := dice().Intn(m.damage()) + m.hit()
		target.takeDamage(damage)
		target.notify(fmt.Sprintf("%s attacks you for %d damagel!%s", m.Name, damage, helpers.Newline))
		m.notify(fmt.Sprintf("You strike %s for %d damagel!%s", target.Name, damage, helpers.Newline))

		if target.Status == dead {
			m.notify(fmt.Sprintf("You have KILLED %s to death!!%s", target.Name, helpers.Newline))
			m.Status = standing

			// whisk it away
			target.die()
			m.ShowStatusBar()
			return
		}
		m.Status = fighting
		target.Status = fighting
	}
}

func (m *mob) die() {
	// drop corpse in room
	corpse := &item{itemType: "corpse", Name: "A corpse of " + m.Name, Identifiers: "corpse," + m.Identifiers, Decays: decays, TTL: 1}
	m.Room.Items = append(m.Room.Items, corpse)

	// whisk them away to Nowhere
	for j, mob := range m.Room.Mobs {
		if mob == m {
			m.Room.Mobs = append(m.Room.Mobs[0:j], m.Room.Mobs[j+1:]...)
			break
		}
	}
	m.Room = getRoom(0)
}

func (m *mob) takeDamage(damage int) {
	m.Hitpoints -= damage
	if m.Hitpoints < 0 {
		m.Status = dead
		m.notify(helpers.Red + "You are DEAD!!!" + helpers.Reset)
	}
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

func (m *mob) move(e *exit) {
	if len(m.Room.Mobs) > 0 {
		for i, rm := range m.Room.Mobs {
			if rm.Room == nil {
				rm.Room = getRoom(uint(rm.RoomID))
			}
			if rm == m {
				m.Room.Mobs = append(m.Room.Mobs[0:i], m.Room.Mobs[i+1:]...)
			} else {
				rm.notify(fmt.Sprintf("%s leaves heading %s.\n", m.Name, e.Dir))
			}
		}
	}

	// add mob to new room list
	m.Room = e.Room
	m.Room.Mobs = append(m.Room.Mobs, m)

	for _, rm := range m.Room.Mobs {
		if rm != m {
			rm.notify(fmt.Sprintf("%s arrives in room.\n", m.Name))
		}
	}
}

func (m *mob) ShowStatusBar() {
	if m.client != nil {
		m.client.BufferData(helpers.White + "[" + m.getHitpoints() + helpers.Reset + helpers.Cyan + "hp")
		m.client.BufferData(helpers.White + m.getMana() + helpers.Reset + helpers.Cyan + "mana ")
		m.client.BufferData(helpers.White + m.getMovement() + helpers.Reset + helpers.Cyan + "mv" + helpers.White)
		m.client.BufferData("] >> ")
		m.client.SendBuffer()
	}
}

func (m *mob) notify(message string) {
	if m.client != nil {
		m.client.SendString(message)
	}
}

func (m *mob) regen() {
	if m.Playable && m.client == nil {
		return
	}

	m.regenHitpoints()
	m.regenMana()
	m.regenMovement()
}

func (m *mob) regenHitpoints() *mob {
	if m.Hitpoints >= m.MaxHitpoints {
		return m
	}
	amount := dice().Intn(int(m.MaxHitpoints/20) + (m.Level * m.Constitution))

	multiplier := 1.0
	switch m.Status {
	case fighting:
		multiplier = 0.0
		break
	case sleeping:
		multiplier = 1.5
		break
	}

	amount = int(float64(amount) * multiplier)

	m.Hitpoints += amount
	if m.Hitpoints > m.MaxHitpoints {
		m.Hitpoints = m.MaxHitpoints
	}
	return m
}

func (m *mob) regenMana() *mob {
	if m.Mana >= m.MaxMana {
		return m
	}
	amount := dice().Intn(int(m.MaxMana/20) + (m.Level * m.Intelligence))
	multiplier := 1.0
	switch m.Status {
	case fighting:
		multiplier = 0.0
		break
	case sleeping:
		multiplier = 1.5
		break
	}

	amount = int(float64(amount) * multiplier)

	m.Mana += amount
	if m.Mana > m.MaxMana {
		m.Mana = m.MaxMana
	}
	return m
}

func (m *mob) regenMovement() *mob {
	if m.Movement >= m.MaxMovement {
		return m
	}
	amount := dice().Intn(int(m.MaxMovement/20) + (m.Level * m.Dexterity))

	multiplier := 1.0
	switch m.Status {
	case fighting:
		multiplier = 0.0
		break
	case sleeping:
		multiplier = 1.5
		break
	}

	amount = int(float64(amount) * multiplier)

	m.Movement += amount
	if m.Movement > m.MaxMovement {
		m.Movement = m.MaxMovement
	}
	return m
}

func (m *mob) wander() {
	if m.client != nil {
		return
	}

	if m.Status != standing {
		return
	}
	switch c := len(m.Room.Exits); c {
	case 0:
		return
	case 1:
		m.move(m.Room.Exits[0])
		return
	default:
		for {
			e := m.Room.Exits[dice().Intn(c)]
			m.move(e)
			return
		}
	}
}

func (m *mob) AddItem(item *item) {
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

func newMobDatabase() {
	mobList = list.New()

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
				db.Create(&m)
			}
			if m.Playable == false {
				mobList.PushBack(&mobs)
			}
		}
	}
}
