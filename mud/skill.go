package mud

type targetType string

const (
	targetIgnore             targetType = "ignore"    // spell chooses its own target
	targetCharacterOffensive targetType = "offensive" // spell is offensive and starts combat
	targetCharacterDefensive targetType = "defensive" // spell is defensive (any character)
	targetCharacterSelf      targetType = "self"      // only castable on same mob
	targetObjectInventory    targetType = "inventory" // used on an object
)

type skillLevel struct {
	ID      int
	Skill   *skill
	SkillID int
	Job     *job
	JobID   int
	Level   int
}

type skill struct {
	ID         int
	Name       string        `json:"name"`
	Levels     []*skillLevel `gorm:"ForeignKey:SkillID"`
	Callback   string        `json:"callback"`
	Target     targetType    `json:"target"`
	MinMana    int           `json:"minMana"`
	Beats      int           `json:"beats"`
	NounDamage string        `json:"nounDamage"` // noun containing message for damage, if applicable
	MessageOff string        `json:"messageOff"` // when skill/spell wears off
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
