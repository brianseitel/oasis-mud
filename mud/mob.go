package mud

import (
	"fmt"

	"github.com/brianseitel/oasis-mud/helpers"
)

const (
	actIsNPC      = 1    /* Auto set for mobs */
	actSentinel   = 2    /* Stays in one room */
	actScavenger  = 4    /* picks up objects */
	actAggressive = 8    /* attacks PCs */
	actStayArea   = 64   /* won't leave area */
	actWimpy      = 128  /* flees when hurt */
	actPet        = 256  /* auto set for pets */
	actTrain      = 512  /* can train PCs */
	actPractice   = 1024 /* can practice PCs */
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

type mobIndex struct {
	ID          uint
	Name        string
	Password    string
	Description string
	Affects     []*affect
	AffectedBy  uint
	Act         uint

	Skills      []*mobSkill
	ItemIds     []int `json:"items"`
	EquippedIds []int `json:"equipped"`
	RoomID      int   `json:"current_room"`
	ExitVerb    string

	Hitpoints    int
	MaxHitpoints int
	Mana         int
	MaxMana      int
	Movement     int
	MaxMovement  int

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

	JobID  int `json:"job"`
	RaceID int `json:"race"`
	Gender int

	Attributes         *attributeSet
	ModifiedAttributes *attributeSet

	Status      status
	Identifiers string

	wait uint

	RecallRoomID uint `json:"recall_room_id"`

	Playable bool
}

type mob struct {
	ID uint

	index *mobIndex

	//Mob information
	Name        string `json:"name"`
	Description string

	Affects    []*affect /* list of affects, incl durations */
	AffectedBy uint      /* bit flag */
	Act        uint

	Skills    []*mobSkill
	Inventory []*item
	Equipped  []*item
	Room      *room
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

	Fight *mob

	wait uint

	Timer     int
	WasInRoom *room

	RecallRoomID uint `json:"recall_room_id"`
	replyTarget  *mob

	Playable bool
	client   *connection
}

func (m *mob) addAffect(af *affect) {
	affectModify(m, af, true)
}

func (m *mob) advanceLevel() {

}

func (m *mob) removeAffect(af *affect) {
	affectModify(m, af, false)
}

func (m *mob) canSee(victim *mob) bool {
	if m == victim {
		return true
	}

	if !victim.isNPC() {
		return false
	}

	if m.isAffected(affectBlind) {
		return false
	}

	if m.Room.isDark() && !m.isAffected(affectInfrared) {
		return false
	}

	if victim.isAffected(affectInvisible) && !m.isAffected(affectDetectInvisible) {
		return false
	}

	if victim.isAffected(affectHide) && !m.isAffected(affectDetectHidden) && victim.Fight == nil {
		return false
	}

	return true
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

func (m *mob) regenHitpoints() int {
	var gain int
	if m.isNPC() {
		gain = m.Level * 3 / 2
	} else {
		gain = helpers.Min(5, m.Level)

		switch m.Status {
		case sleeping:
			gain += m.ModifiedAttributes.Constitution
			break
		case resting:
			gain += m.ModifiedAttributes.Constitution / 2
			break
		}
	}

	if helpers.HasBit(m.AffectedBy, affectPoison) {
		gain /= 4
	}

	return helpers.Min(gain, m.MaxHitpoints-m.Hitpoints)
}

func (m *mob) regenMana() int {
	var gain int
	if m.isNPC() {
		gain = m.Level * 3 / 2
	} else {
		gain = helpers.Min(5, m.Level)

		switch m.Status {
		case sleeping:
			gain += m.ModifiedAttributes.Intelligence
			break
		case resting:
			gain += m.ModifiedAttributes.Intelligence / 2
			break
		}
	}

	if helpers.HasBit(m.AffectedBy, affectPoison) {
		gain /= 4
	}

	return helpers.Min(gain, m.MaxMana-m.Mana)
}

func (m *mob) regenMovement() int {
	var gain int
	if m.isNPC() {
		gain = m.Level * 3 / 2
	} else {
		gain = helpers.Min(5, m.Level)

		switch m.Status {
		case sleeping:
			gain += m.ModifiedAttributes.Dexterity
			break
		case resting:
			gain += m.ModifiedAttributes.Dexterity / 2
			break
		}

	}
	if helpers.HasBit(m.AffectedBy, affectPoison) {
		gain /= 4
	}

	return helpers.Min(gain, m.MaxMovement-m.Movement)
}

func (m *mob) stripAffect(name string) {
	for _, af := range m.Affects {
		if af.affectType.Skill.Name == name {
			m.removeAffect(af)
		}
	}
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
		killer.Alignment -= killer.Alignment / 4
		xp = 3 * xp / 4
	}

	xp = helpers.Max(5, int(xp*5/4))
	mod := int(xp * 3 / 4)
	xp = dice().Intn(xp) + mod
	xp = helpers.Max(0, xp)

	return xp
}
