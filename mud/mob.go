package mud

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"container/list"

	"github.com/brianseitel/oasis-mud/helpers"
)

var (
	mobList *list.List
)

type mobSkill struct {
	Skill   *skill `json:"-"`
	SkillID uint   `json:"skill_id"`
	Level   uint   `json:"level"`
}

type attributeSet struct {
	Strength     int
	Wisdom       int
	Intelligence int
	Dexterity    int
	Charisma     int
	Constitution int
}

type mob struct {
	ID uint

	//Mob information
	Name        string `json:"name"`
	Password    string
	Description string

	Affects    []*affect /* list of affects, incl durations */
	AffectedBy uint      /* bit flag */

	Skills    []*mobSkill
	Inventory []*item
	Equipped  []*item
	ItemIds   []int `json:"items"`
	Room      *room
	RoomID    int `json:"current_room"`
	ExitVerb  string

	Hitpoints    int `json:"hitpoints"`
	MaxHitpoints int `json:"max_hitpoints"`
	Mana         int `json:"mana"`
	MaxMana      int `json:"max_mana"`
	Movement     int `json:"movement"`
	MaxMovement  int `json:"max_movement"`

	Armor   int
	Hitroll int
	Damroll int

	Exp       int
	Level     int
	Alignment int
	Practices uint
	Gold      uint

	Carrying       uint `json:"carrying"`
	CarryMax       uint `json:"carry_max"`
	CarryWeight    uint `json:"carry_weight"`
	CarryWeightMax uint `json:"carry_weight_max"`

	Job    *job  `json:"-"`
	JobID  int   `json:"job"`
	Race   *race `json:"-"`
	RaceID int   `json:"race"`
	Gender int

	Attributes         *attributeSet
	ModifiedAttributes *attributeSet

	Status      status
	Identifiers string

	Fight   *fight
	FightID uint
	wait    uint

	RecallRoomID uint `json:"recall_room_id"`
	replyTarget  *mob

	Playable bool
	client   *connection
}

func (m *mob) addAffect(af *affect) {
	affectModify(m, af, true)
}

func (m *mob) removeAffect(af *affect) {
	affectModify(m, af, false)
}

func (m *mob) advanceLevel() {

}

func (m *mob) gainExp(gain int) {
	if m.isNPC() || m.Level >= 99 {
		return
	}

	m.Exp += helpers.Max(1000, m.Exp+gain)
	for m.Level < 99 && m.Exp >= 1000*(m.Level+1) {
		m.notify("You raise a level!")
		m.Level++
		m.advanceLevel()
	}
}

func (m *mob) isAffected(flag uint) bool {
	return helpers.HasBit(m.AffectedBy, flag)
}

func (m *mob) hasAffect(ms *mobSkill) bool {
	for _, affect := range m.Affects {
		if affect.affectType == ms {
			return true
		}
	}
	return false
}

func (m *mob) isAwake() bool {
	return m.Status > sleeping
}

func (m *mob) isEvil() bool {
	return m.Alignment <= -350
}

func (m *mob) isGood() bool {
	return m.Alignment >= 350
}

func (m *mob) isImmortal() bool {
	return false
}

func (m *mob) isNeutral() bool {
	return !m.isEvil() && !m.isGood()
}

func (m *mob) isNPC() bool {
	return !m.Playable
}

func (m *mob) isSafe() bool {
	return false
}

func (m *mob) isSilenced() bool {
	return false
}

func (m *mob) isTrainer() bool {
	return true
}

func (m *mob) hit() int {
	return int(m.Attributes.Dexterity / 3)
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
				if !rm.isAffected(affectSneak) {
					rm.notify("%s leaves heading %s.", m.Name, e.Dir)
				}
			}
		}
	}

	// add mob to new room list
	m.Room = e.Room
	m.Room.Mobs = append(m.Room.Mobs, m)

	for _, rm := range m.Room.Mobs {
		if rm != m {
			if !rm.isAffected(affectSneak) {
				rm.notify("%s arrives in room.", m.Name)
			}
		}
	}
}

func (m *mob) statusBar() {
	if m.client != nil {
		m.client.BufferData(fmt.Sprintf("%s[%d%s%shp %s%d%s%smana %s%d%s%smv%s] >>",
			helpers.White, m.Hitpoints, helpers.Reset, helpers.Cyan,
			helpers.White, m.Mana, helpers.Reset, helpers.Cyan,
			helpers.White, m.Movement, helpers.Reset, helpers.Cyan,
			helpers.White))
		m.client.SendBuffer()
	}
}

func (m *mob) equipped(position uint) string {
	equipped := m.equippedItem(position)

	if equipped == nil {
		return "<empty>"
	}

	return equipped.Name
}

func (m *mob) equipItem(item *item, position uint) {
	if m.equippedItem(position) == item {
		return
	}

	if (item.hasExtraFlag(itemAntiEvil) && m.isEvil()) || (item.hasExtraFlag(itemAntiGood) && m.isGood()) || (item.hasExtraFlag(itemAntiNeutral) && m.isNeutral()) {
		m.notify("You are zapped by %s and drop it!%s", item.Name, helpers.Newline)
		m.Room.notify(fmt.Sprintf("%s is zapped by %s and drops it!%s", m.Name, item.Name, helpers.Newline), m)
		// TODO: dropItem()
		return
	}

	m.Armor -= applyAC(item, int(position))
	item.WearLocation = position

	m.Equipped = append(m.Equipped, item)
	for j, i := range m.Inventory {
		if i == item {
			m.Inventory = append(m.Inventory[0:j], m.Inventory[j+1:]...)
			break
		}
	}

	// TODO: item effects

	// TODO: light up room if it's a light

	return
}

func (m *mob) equippedItem(position uint) *item {
	for _, i := range m.Equipped {
		if i.WearLocation == position {
			return i
		}
	}
	return nil
}

func (m *mob) notify(message string, a ...interface{}) {
	if m.client != nil {
		message = fmt.Sprintf("%s%s", message, helpers.Newline)
		m.client.SendString(fmt.Sprintf(message, a...))
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

	amount := dice().Intn(int(m.MaxHitpoints/20) + (m.Level * m.Attributes.Constitution))

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
	amount := dice().Intn(int(m.MaxMana/20) + (m.Level * m.Attributes.Intelligence))
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
	amount := dice().Intn(int(m.MaxMovement/20) + (m.Level * m.Attributes.Dexterity))

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

func (m *mob) looksAt(target interface{}) string {
	// try mob first

	var (
		victim *mob
		obj    *item
	)
	switch target.(type) {
	case *mob:
		victim = target.(*mob)
		return victim.Name
	case *item:
		obj = target.(*item)
		return obj.Name
	}
	return ""
}

func getPlayerByName(name string) *mob {
	for e := mobList.Front(); e != nil; e = e.Next() {
		mob := e.Value.(*mob)

		if mob.Playable && mob.Name == name {
			return mob
		}
	}

	return nil
}

func newMobDatabase() {
	mobList = list.New()

	mobFiles, _ := filepath.Glob("./data/mobs/*.json")

	for _, mobFile := range mobFiles {
		file, err := ioutil.ReadFile(mobFile)
		if err != nil {
			panic(err)
		}

		var list []*mob
		err = json.Unmarshal(file, &list)
		if err != nil {
			panic(err)
		}

		for _, m := range list {

			var skills []*mobSkill
			for _, s := range m.Skills {
				skill := getSkill(s.SkillID)
				skills = append(skills, &mobSkill{Skill: skill, SkillID: s.SkillID, Level: s.Level})
			}
			m.Skills = skills

			mobList.PushBack(m)
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
