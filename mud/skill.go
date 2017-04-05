package mud

type targetType string

const (
	targetIgnore             targetType = "ignore"    // spell chooses its own target
	targetCharacterOffensive targetType = "offensive" // spell is offensive and starts combat
	targetCharacterDefensive targetType = "defensive" // spell is defensive (any character)
	targetCharacterSelf      targetType = "self"      // only castable on same mob
	targetObjectInventory    targetType = "inventory" // used on an object
)

type skill struct {
	ID     int
	Name   string     `json:"name"`
	Target targetType `json:"target"`
	Levels struct {
		Warrior int
		Mage    int
		Cleric  int
		Thief   int
		Ranger  int
		Bard    int
	}
	MinMana    int    `json:"minMana"`
	Beats      int    `json:"beats"`
	NounDamage string `json:"nounDamage"` // noun containing message for damage, if applicable
	MessageOff string `json:"messageOff"` // when skill/spell wears off
}

func getSkill(id int) *skill {
	for e := skillList.Front(); e != nil; e = e.Next() {
		s := e.Value.(*skill)
		if s.ID == id {
			return s
		}
	}

	return nil
}

func getSkillByName(name string) *skill {
	for e := skillList.Front(); e != nil; e = e.Next() {
		s := e.Value.(*skill)
		if s.Name == name {
			return s
		}
	}

	return nil
}
