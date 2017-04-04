package mud

import (
	"bytes"
	"fmt"
	"unicode"

	"strings"
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

const (
	playerIsNPC         = 1
	playerBoughtPet     = 2
	playerAutoExit      = 4
	playerAutoLook      = 8
	playerAutoSacrifice = 16
	playerBlank         = 32
	playerBrief         = 128
	playerCombine       = 256
	playerPrompt        = 512
	playerTelnetGA      = 1024
	playerHolylight     = 2048
	playerWizInvis      = 4096
	playerSilence       = 8192
	playerNoEmote       = 32768
	playerNoTell        = 65536
	playerLog           = 262144
	playerDeny          = 524288
	playerFreeze        = 1048576
	playerThief         = 2097152
	playerKiller        = 4194304
)

type mobSkill struct {
	Skill   *skill `json:"-"`
	SkillID int    `json:"skill_id"`
	Level   int    `json:"level"`
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
	ID       int
	Name     string
	Password string

	Description     string
	LongDescription string `json:"long_description"`
	Title           string

	Affects    []*affect
	AffectedBy []int
	Act        []int

	Skills    []*mobSkill
	Inventory []*item `json:"inventory"`
	Equipped  []*item `json:"equipped"`
	RoomID    int     `json:"current_room"`

	ExitVerb string
	Bamfin   string
	Bamfout  string

	Hitpoints    int
	MaxHitpoints int `json:"max_hitpoints"`
	Mana         int
	MaxMana      int `json:"max_mana"`
	Movement     int
	MaxMovement  int `json:"max_movement"`

	Armor   int
	Hitroll int
	Damroll int

	Exp       int
	Level     int
	Alignment int
	Practices int
	Gold      int
	Trust     int

	Carrying       int `json:"carrying"`
	CarryMax       int `json:"carry_max"`
	CarryWeight    int `json:"carry_weight"`
	CarryWeightMax int `json:"carry_weight_max"`

	JobID  int `json:"job"`
	RaceID int `json:"race"`
	Gender int

	Attributes         *attributeSet
	ModifiedAttributes *attributeSet `json:"modified_attributes"`

	Status      status
	Identifiers string
	Shop        *shop

	wait  int
	count int

	RecallRoomID int `json:"recall_room_id"`

	Playable bool
}

type mob struct {
	ID      int
	SavedAt string

	index *mobIndex

	//Mob information
	Name            string `json:"name"`
	Password        string
	Description     string
	LongDescription string
	Title           string

	Affects    []*affect /* list of affects, incl durations */
	AffectedBy int       /* bit flag */
	Act        int

	Skills    []*mobSkill
	Inventory []*item
	Equipped  []*item
	Room      *room

	ExitVerb string
	Bamfout  string
	Bamfin   string

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
	Practices int
	Gold      int
	Trust     int

	Carrying       int `json:"carrying"`
	CarryMax       int `json:"carry_max"`
	CarryWeight    int `json:"carry_weight"`
	CarryWeightMax int `json:"carry_weight_max"`

	Job    *job  `json:"-"`
	Race   *race `json:"-"`
	Gender int

	Attributes         *attributeSet
	ModifiedAttributes *attributeSet

	Status       status
	RecallRoomID int `json:"recall_room_id"`

	/* dynamic stuff */
	Fight       *mob
	master      *mob
	leader      *mob
	wait        int
	Timer       int
	Light       int
	WasInRoom   *room
	replyTarget *mob
	Playable    bool
	client      *connection
}

func (m *mob) addAffect(af *affect) {
	affectModify(m, af, true)
}

func (m *mob) addFollower(master *mob) {
	if m.master != nil {
		return
	}

	m.master = master
	m.leader = nil

	if master.canSee(m) {
		act("$n now follows you.", m, nil, master, actToVict)
	}
	act("You now follow $N.", m, nil, master, actToChar)
	return
}

func (m *mob) advanceLevel() {
	var (
		addHP       int
		addMana     int
		addMovement int
		addPracs    int
	)

	job := m.Job
	addHP = bonusTableConstitution[m.currentConstitution()].hitpoints + (dice().Intn(job.MaxHitpoints) + job.MinHitpoints)

	addMana = 0
	if job.GainsMana {
		addMana = dice().Intn((m.currentIntelligence()*2)+(m.currentWisdom()/8)) + 2
	}

	addMovement = dice().Intn(m.currentConstitution()+(m.currentDexterity()/4)) + 5

	addPracs = bonusTableWisdom[m.currentWisdom()].practice + dice().Intn(3) + 1

	addHP = max(1, addHP)
	addMana = max(0, addMana)
	addMovement = max(10, addMovement)

	m.MaxHitpoints += addHP
	m.MaxMana += addMana
	m.MaxMovement += addMovement
	m.Practices += addPracs

	if !m.isNPC() {
		removeBit(m.Act, playerBoughtPet)
	}

	m.notify("Your gain is: %d/%d hp, %d/%d mana, %d/%d movement, and %d/%d practices.", addHP, m.MaxHitpoints, addMana, m.MaxMana, addMovement, m.MaxMovement, addPracs, m.Practices)
}

func (m *mob) currentStrength() int {
	if m.isNPC() {
		return 13
	}

	var max int
	if m.Job != nil && m.Job.PrimeAttribute == applyStrength {
		max = 27
	} else {
		max = 21
	}

	return uRange(3, m.Attributes.Strength+m.ModifiedAttributes.Strength, max)

}

func (m *mob) currentIntelligence() int {
	if m.isNPC() {
		return 13
	}

	var max int
	if m.Job != nil && m.Job.PrimeAttribute == applyIntelligence {
		max = 27
	} else {
		max = 21
	}

	return uRange(3, m.Attributes.Intelligence+m.ModifiedAttributes.Intelligence, max)

}

func (m *mob) currentWisdom() int {
	if m.isNPC() {
		return 13
	}

	var max int
	if m.Job != nil && m.Job.PrimeAttribute == applyWisdom {
		max = 27
	} else {
		max = 21
	}

	return uRange(3, m.Attributes.Wisdom+m.ModifiedAttributes.Wisdom, max)

}

func (m *mob) currentDexterity() int {
	if m.isNPC() {
		return 13
	}

	var max int
	if m.Job != nil && m.Job.PrimeAttribute == applyDexterity {
		max = 27
	} else {
		max = 21
	}

	return uRange(3, m.Attributes.Dexterity+m.ModifiedAttributes.Dexterity, max)

}

func (m *mob) currentConstitution() int {
	if m.isNPC() {
		return 13
	}

	var max int
	if m.Job != nil && m.Job.PrimeAttribute == applyConstitution {
		max = 27
	} else {
		max = 21
	}

	return uRange(3, m.Attributes.Constitution+m.ModifiedAttributes.Constitution, max)

}

func (m *mob) currentCharisma() int {
	if m.isNPC() {
		return 13
	}

	var max int
	if m.Job != nil && m.Job.PrimeAttribute == applyCharisma {
		max = 27
	} else {
		max = 21
	}

	return uRange(3, m.Attributes.Charisma+m.ModifiedAttributes.Charisma, max)

}

func (m *mob) removeAffect(af *affect) {
	affectModify(m, af, false)
}

func (m *mob) canSee(victim *mob) bool {
	if m == victim {
		return true
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

func (m *mob) checkBlind() bool {
	if hasBit(m.AffectedBy, affectBlind) {
		m.notify("You can't see a thing!")
		return false
	}
	return true
}

func (m *mob) dieFollower() {
	if m.master != nil {
		m.stopFollower()
	}

	m.leader = nil

	for e := mobList.Front(); e != nil; e = e.Next() {
		ch := e.Value.(*mob)
		if ch.master == m {
			ch.stopFollower()
		}
		if ch.leader == m {
			ch.leader = ch
		}
	}

	return
}

func (m *mob) gainExp(gain int) {
	if m.isNPC() || m.Level >= 99 {
		return
	}

	m.Exp = max(1000, m.Exp+gain)
	for m.Level < 99 && m.Exp >= 1000*(m.Level+1) {
		m.notify("You raise a level!")
		m.Level++
		m.advanceLevel()
	}
}

func (m *mob) isAffected(flag int) bool {
	return hasBit(m.AffectedBy, flag)
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

func (m *mob) hasKey(key int) bool {
	for _, i := range m.Inventory {
		if i.ID == int(key) {
			return true
		}
	}
	return false
}

func (m *mob) hit() int {
	return int(m.Attributes.Dexterity / 3)
}

func (m mob) TNL() int {
	return ((m.Level + 1) * 1000) - m.Exp
}

func (m *mob) move(e *exit) {

	if e.isClosed() && !m.isAffected(affectPassDoor) {
		act("The $d is closed.", m, nil, e.Keyword, actToChar)
		return
	}

	if e.Room.isPrivate() {
		m.notify("That room is private right now.")
		return
	}

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
	oldRoom := m.Room
	m.Room = e.Room
	m.Room.Mobs = append(m.Room.Mobs, m)

	if m.equippedItem(itemLight) != nil {
		if oldRoom.Light > 0 {
			oldRoom.Light--
		}
		if m.Room != nil {
			m.Room.Light++
		}
	}

	for _, rm := range oldRoom.Mobs {
		if rm.master == m {
			rm.move(e)
			interpret(rm, "look")
		}
	}

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
		m.client.BufferData(fmt.Sprintf("%s[%d/%d%s%shp %s%d/%d%s%smana %s%d/%d%s%smv%s] >>",
			white, m.Hitpoints, m.MaxHitpoints, reset, cyan,
			white, m.Mana, m.MaxMana, reset, cyan,
			white, m.Movement, m.MaxMovement, reset, cyan,
			white))
		m.client.SendBuffer()
	}
}

func (m *mob) stopFollower() {
	if m.master == nil {
		return
	}

	if hasBit(m.AffectedBy, affectCharm) {
		removeBit(m.AffectedBy, affectCharm)
		m.stripAffect("charm")
	}

	if m.master.canSee(m) {
		act("$n stops following you.", m, nil, m.master, actToVict)
	}
	act("You stop following $N.", m, nil, m.master, actToChar)

	m.master = nil
	m.leader = nil
	return
}

func (m *mob) equipped(position int) string {
	equipped := m.equippedItem(position)

	if equipped == nil {
		return "<empty>"
	}

	return equipped.Name
}

func (m *mob) equipItem(item *item, position int) {
	if m.equippedItem(position) == item {
		return
	}

	if (item.hasExtraFlag(itemAntiEvil) && m.isEvil()) || (item.hasExtraFlag(itemAntiGood) && m.isGood()) || (item.hasExtraFlag(itemAntiNeutral) && m.isNeutral()) {
		m.notify("You are zapped by %s and drop it!%s", item.Name, newline)
		m.Room.notify(fmt.Sprintf("%s is zapped by %s and drops it!%s", m.Name, item.Name, newline), m)
		m.Room.Items = append(m.Room.Items, item)
		for j, i := range m.Inventory {
			if i == item {
				m.Inventory = append(m.Inventory[:j], m.Inventory[j+1:]...)
				break
			}
		}
		return
	}

	m.Armor -= applyAC(item, int(position))
	item.WearLocation = position

	// TODO: item effects

	if item.ItemType == itemLight && m.Room != nil {
		m.Room.Light++
	}
	return
}

func (m *mob) equippedItem(position int) *item {
	for _, i := range m.Inventory {
		if i.WearLocation == position {
			return i
		}
	}
	return nil
}

func (m *mob) hitroll() int {
	return m.Hitroll + m.currentStrength()
}

func (m *mob) notify(message string, a ...interface{}) {
	if m.client != nil {
		message = fmt.Sprintf("%s%s%s", reset, message, newline)
		m.client.SendString(fmt.Sprintf(message, a...))
	}
}

func (m *mob) regenHitpoints() int {
	var gain int
	if m.isNPC() {
		gain = m.Level * 3 / 2
	} else {
		gain = min(5, m.Level)

		switch m.Status {
		case sleeping:
			gain += m.currentConstitution()
			break
		case resting:
			gain += m.currentConstitution() / 2
			break
		}
	}

	if hasBit(m.AffectedBy, affectPoison) {
		gain /= 4
	}

	return min(gain, m.MaxHitpoints-m.Hitpoints)
}

func (m *mob) regenMana() int {
	var gain int
	if m.isNPC() {
		gain = m.Level * 3 / 2
	} else {
		gain = min(5, m.Level)

		switch m.Status {
		case sleeping:
			gain += m.currentIntelligence()
			break
		case resting:
			gain += m.currentIntelligence() / 2
			break
		}
	}

	if hasBit(m.AffectedBy, affectPoison) {
		gain /= 4
	}

	return min(gain, m.MaxMana-m.Mana)
}

func (m *mob) regenMovement() int {
	var gain int
	if m.isNPC() {
		gain = m.Level * 3 / 2
	} else {
		gain = min(5, m.Level)

		switch m.Status {
		case sleeping:
			gain += m.currentDexterity()
			break
		case resting:
			gain += m.currentDexterity() / 2
			break
		}

	}
	if hasBit(m.AffectedBy, affectPoison) {
		gain /= 4
	}

	return min(gain, m.MaxMovement-m.Movement)
}

func (m *mob) stripAffect(name string) {
	for _, af := range m.Affects {
		if af.affectType.Skill.Name == name {
			m.removeAffect(af)
		}
	}
}

func (m *mob) getTrust() int {
	if m.client != nil && m.client.original != nil {
		m = m.client.original
	}

	if m.Trust != 0 {
		return m.Trust
	}

	if m.isNPC() && m.Level >= 90 {
		return 89
	}

	return m.Level
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
		if s.Skill != nil && s.Skill.Name == name {
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

func showCharactersToPlayer(chars []*mob, player *mob) {
	for _, char := range chars {
		if char == player {
			continue
		}
		if player.canSee(char) {
			showCharacterToPlayer(char, player, false)
		} else if player.Room.isDark() && player.isAffected(affectInfrared) {
			player.notify("%sYou see glowing red eyes watching YOU!%s", red, reset)
		}
	}
}

func showCharacterToPlayer(victim *mob, player *mob, showLongDesc bool) {
	var buf bytes.Buffer

	if victim == player {
		return
	}

	if victim.isAffected(affectInvisible) {
		buf.Write([]byte("(Invis)"))
	}
	if victim.isAffected(affectHide) {
		buf.Write([]byte("(Hide)"))
	}
	if victim.isAffected(affectCharm) {
		buf.Write([]byte("(Charmed)"))
	}
	if victim.isAffected(affectPassDoor) {
		buf.Write([]byte("(Translucent)"))
	}
	if victim.isAffected(affectFaerieFire) {
		buf.Write([]byte("(Pink Aura)"))
	}
	if victim.isEvil() && player.isAffected(affectDetectEvil) {
		buf.Write([]byte("(Red Aura)"))
	}
	if victim.isAffected(affectSanctuary) {
		buf.Write([]byte("(white Aura)"))
	}

	if victim.Status == standing && len(victim.Description) > 0 {

		if showLongDesc {
			buf.Write([]byte(victim.LongDescription))
		} else {
			buf.Write([]byte(victim.Description))
		}
		buf.Write([]byte(reset))
		player.notify("%s%s%s", cyan, buf.String(), reset)
		return
	}

	buf.Write([]byte(victim.Name))

	if !victim.isNPC() {
		buf.Write([]byte(" "))
		buf.Write([]byte(victim.Title))
	}

	switch victim.Status {
	case dead:
		buf.Write([]byte(" is DEAD!!"))
		break
	case mortal:
		buf.Write([]byte(" is mortally wounded."))
		break
	case incapacitated:
		buf.Write([]byte(" is incapacitated."))
		break
	case stunned:
		buf.Write([]byte(" is lying here stunned."))
		break
	case sleeping:
		buf.Write([]byte(" is sleeping here."))
		break
	case resting:
		buf.Write([]byte(" is resting here."))
		break
	case standing:
		buf.Write([]byte(" is here."))
		break
	case fighting:
		buf.Write([]byte(" is here fighting "))
		if victim.Fight == nil {
			buf.Write([]byte(" thin air?"))
		} else if victim.Fight == player {
			buf.Write([]byte(" YOU!"))
		} else if victim.Room == victim.Fight.Room {
			buf.Write([]byte(pers(victim.Fight, player)))
			buf.Write([]byte("."))
		} else {
			buf.Write([]byte(" someone who left?"))
		}
		break
	}

	output := buf.String()
	a := []rune(output)
	a[0] = unicode.ToLower(a[0])
	output = string(a)
	player.notify("%s%s%s", cyan, output, reset)
	return
}

func showItemsToPlayer(items []*item, player *mob) {
	if player.client == nil {
		return
	}

	nShow := 0

	var itemList []string
	var itemCounts []int
	for _, item := range items {
		if player.canSeeItem(item) {
			itemDesc := formatItemToChar(item, player)
			combine := false
			for iShow := nShow - 1; iShow >= 0; iShow-- {
				if strings.HasPrefix(itemList[iShow], itemDesc) {
					itemCounts[iShow]++
					combine = true
					break
				}
			}

			if !combine {
				itemList = append(itemList, itemDesc)
				itemCounts = append(itemCounts, 1)
				nShow++
			}
		}
	}

	for iShow := 0; iShow < nShow; iShow++ {
		var buf bytes.Buffer
		buf.Write([]byte(cyan))
		if itemCounts[iShow] != 1 {
			buf.Write([]byte(fmt.Sprintf("(%d) ", itemCounts[iShow])))
		}
		buf.Write([]byte(itemList[iShow]))
		buf.Write([]byte(reset))
		player.notify(buf.String())
	}

}

func formatItemToChar(item *item, player *mob) string {
	var buf bytes.Buffer

	if item.hasExtraFlag(itemInvis) {
		buf.Write([]byte("(Invis)"))
	}
	if player.isAffected(affectDetectEvil) && item.hasExtraFlag(itemEvil) {
		buf.Write([]byte("(Red Aura)"))
	}
	if player.isAffected(affectDetectMagic) && item.hasExtraFlag(itemMagic) {
		buf.Write([]byte("(Magical)"))
	}
	if item.hasExtraFlag(itemGlow) {
		buf.Write([]byte("(Glow)"))
	}
	if item.hasExtraFlag(itemHum) {
		buf.Write([]byte("(Humming)"))
	}

	buf.Write([]byte(item.Description))

	return buf.String()
}

func pers(m *mob, looker *mob) string {
	if looker.canSee(m) {
		return m.Name
	}

	return "someone"
}

func xpCompute(killer *mob, target *mob) int {
	var xp int

	xp = 300 - uRange(-3, killer.Level-target.Level, 6)*50

	// do align check
	align := killer.Alignment - target.Alignment

	if align > 500 {
		killer.Alignment = min(killer.Alignment+(align-500)/4, 1000)
		xp = 5 * xp / 4
	} else if align < -500 {
		killer.Alignment = max(killer.Alignment+(align+500)/4, -1000)
	} else {
		killer.Alignment -= killer.Alignment / 4
		xp = 3 * xp / 4
	}

	xp = max(5, int(xp*5/4))
	mod := int(xp * 3 / 4)
	xp = dice().Intn(xp) + mod
	xp = max(0, xp)

	return xp
}
