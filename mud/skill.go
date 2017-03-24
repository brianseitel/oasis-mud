package mud

import (
	"container/list"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

type targetType string

const (
	targetIgnore             targetType = "ignore"    // spell chooses its own target
	targetCharacterOffensive targetType = "offensive" // spell is offensive and starts combat
	targetCharacterDefensive targetType = "defensive" // spell is defensive (any character)
	targetCharacterSelf      targetType = "self"      // only castable on same mob
	targetObjectInventory    targetType = "inventory" // used on an object
)

type skillLevel struct {
	ID      uint
	Skill   *skill
	SkillID uint
	Job     *job
	JobID   uint
	Level   uint
}

type skill struct {
	ID         uint
	Name       string        `json:"name"`
	Levels     []*skillLevel `gorm:"ForeignKey:SkillID"`
	Callback   string        `json:"callback"`
	Target     string        `json:"target"`
	MinMana    int           `json:"minMana"`
	Beats      int           `json:"beats"`
	NounDamage string        `json:"nounDamage"` // noun containing message for damage, if applicable
	MessageOff string        `json:"messageOff"` // when skill/spell wears off
}

var (
	skillList list.List
)

func newSkillDatabase() {
	skillFiles, _ := filepath.Glob("./data/skills/*.json")

	for _, skillFile := range skillFiles {
		file, err := ioutil.ReadFile(skillFile)
		if err != nil {
			panic(err)
		}

		var list []*skill
		json.Unmarshal(file, &list)

		for _, sk := range list {
			skillList.PushBack(sk)
		}
	}
}

func getSkill(id uint) *skill {
	for e := skillList.Front(); e != nil; e = e.Next() {
		s := e.Value.(*skill)
		if s.ID == id {
			return s
		}
	}

	return nil
}
