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

type mobSkill struct {
	Skill   *skill `json:"-"`
	SkillID uint   `json:"skill_id"`
	Level   uint   `json:"level"`
}

type mob struct {
	gorm.Model

	//Mob information
	Name        string `json:"name" gorm:"name"`
	Password    string `gorm:"password"`
	Description string `gorm:"type:text"`

	Skills    []*mobSkill
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

	Exp       int
	Level     int
	Alignment int

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

	RecallRoomID uint `json:"recall_room_id"`

	Playable bool
	client   *connection
}

func (m *mob) addItem(item *item) {
	m.Inventory = append(m.Inventory, item)
}

func (m mob) setFight(f *fight) {
	m.Fight = f
}

func (m *mob) checkLevelUp() {
	if m.Exp > (1000 * m.Level) {
		m.Level++
		m.notify(fmt.Sprintf("You have LEVELED UP! You are now Level %d!%s", m.Level, helpers.Newline))
	}
}

func (m *mob) isAwake() bool {
	return m.Status > sleeping
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

func (m *mob) equipped(position position) string {
	equipped := m.equippedItem(position)

	if equipped == nil {
		return "<empty>"
	}

	return equipped.Name
}

func (m *mob) equippedItem(position position) *item {
	for _, i := range m.Equipped {
		if i.Position == string(position) {
			return i
		}
	}
	return nil
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

func (m *mob) skill(name string) *mobSkill {
	for _, s := range m.Skills {
		if s.Skill.Name == name {
			return s
		}
	}
	return nil
}

func (m *mob) loadSkills() {
	var skills []*mobSkill
	for _, s := range m.Skills {
		skill := getSkill(s.SkillID)
		skills = append(skills, &mobSkill{Skill: skill, SkillID: s.SkillID, Level: s.Level})
	}
	m.Skills = skills
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

			var skills []*mobSkill
			for _, s := range m.Skills {
				skill := getSkill(s.SkillID)
				skills = append(skills, &mobSkill{Skill: skill, SkillID: s.SkillID, Level: s.Level})
			}
			mobs.Skills = skills

			if m.Playable == false {
				mobList.PushBack(&mobs)
			}
		}
	}

}

func xpCompute(killer *mob, target *mob) int {
	var xp int

	xp = 300 - helpers.Range(-3, killer.Level-target.Level, 6)*50

	// do align check
	align := killer.Alignment - target.Alignment

	if align > 500 {
		killer.Alignment = helpers.Min(killer.Alignment+(align-500)/4, 1000)
		xp = 5 * xp / 4
	} else if align < -500 {
		killer.Alignment = helpers.Max(killer.Alignment+(align+500)/4, -1000)
	} else {
		xp = 3 * xp / 4
	}

	xp = helpers.Max(5, int(xp*5/4))
	mod := int(xp * 3 / 4)
	xp = dice().Intn(xp) + mod
	xp = helpers.Max(0, xp)

	return xp
}
