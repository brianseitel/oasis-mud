package mud

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

type saveAffect struct {
	AffectType saveSkill
	Duration   int
	Location   int
	Modifier   int
	Bitvector  int
}

type saveItem struct {
	ID               int
	Name             string
	Description      string
	ShortDescription string
	ItemType         int        `json:"item_type"`
	Contained        []saveItem `json:"contained"`
	Affected         []saveAffect
	ExtraFlags       int
	WearFlags        int
	WearLocation     int
	Weight           int
	Cost             int
	Level            int
	Timer            int
	Value            int
	Min              int
	Max              int
	SkillID          int /* items can have skills */
	Charges          int
}

type saveSkill struct {
	SkillID int
	Level   int
}

type savePlayer struct {
	ID       int    `json:"id"`
	SavedAt  string `json:"saved_at"`
	Name     string `json:"name"`
	Password string `json:"password"`

	Description string `json:"description"`
	Title       string `json:"title"`

	Affects    []saveAffect `json:"affects"`
	AffectedBy int          `json:"affected_by"`
	Act        int          `json:"act"`

	Skills    []saveSkill `json:"skills"`
	Inventory []saveItem  `json:"inventory"`
	Equipped  []saveItem  `json:"equipped"`
	RoomID    int         `json:"current_room"`

	ExitVerb string `json:"exit_verb"`
	Bamfin   string `json:"bamfin"`
	Bamfout  string `json:"bamfout"`

	Hitpoints    int `json:"hitpoints"`
	MaxHitpoints int `json:"max_hitpoints"`
	Mana         int `json:"mana"`
	MaxMana      int `json:"max_mana"`
	Movement     int `json:"movement"`
	MaxMovement  int `json:"max_movement"`

	Armor   int `json:"armor"`
	Hitroll int `json:"hitroll"`
	Damroll int `json:"damroll"`

	Exp       int `json:"exp"`
	Level     int `json:"level"`
	Alignment int `json:"alignment"`
	Practices int `json:"practices"`
	Gold      int `json:"gold"`
	Trust     int `json:"trust"`

	Carrying       int `json:"carrying"`
	CarryMax       int `json:"carry_max"`
	CarryWeight    int `json:"carry_weight"`
	CarryWeightMax int `json:"carry_weight_max"`

	JobID  int `json:"job"`
	RaceID int `json:"race"`
	Gender int `json:"gender"`

	Attributes         *attributeSet `json:"attributes"`
	ModifiedAttributes *attributeSet `json:"modified_attributes"`

	Status       status `json:"status"`
	ShopID       int    `json:"shop_id"`
	RecallRoomID int    `json:"recall_room_id"`

	Playable bool `json:"playable"`
}

func saveCharacter(character *mob) {

	if character.isNPC() { //|| character.Level < 2 {
		fmt.Println("fuck")
		return
	}

	if character.client != nil && character.client.original != nil {
		character = character.client.original
	}

	fmt.Println("Saving", character.Name)
	path := fmt.Sprintf("./data/savefiles/%s.json", character.Name)
	writeCharacter(character, path)

	return
}

func writeCharacter(character *mob, path string) {

	var save savePlayer

	save.SavedAt = time.Now().String()
	save.ID = character.ID
	save.Name = character.Name
	save.Description = character.Description
	save.Gender = character.Gender
	save.JobID = character.Job.ID
	save.RaceID = character.Race.ID
	save.Level = character.Level
	save.Trust = character.Trust
	save.RoomID = character.Room.ID
	save.Hitpoints = character.Hitpoints
	save.MaxHitpoints = character.MaxHitpoints
	save.Mana = character.Mana
	save.MaxMana = character.MaxMana
	save.Movement = character.Movement
	save.MaxMovement = character.MaxMovement
	save.Gold = character.Gold
	save.Exp = character.Exp
	save.Act = character.Act
	save.AffectedBy = character.AffectedBy

	if character.Status == fighting {
		save.Status = standing
	} else {
		save.Status = character.Status
	}

	save.Practices = character.Practices
	save.Alignment = character.Alignment
	save.Hitroll = character.Hitroll
	save.Damroll = character.Damroll
	save.Armor = character.Armor

	if !character.isNPC() {
		save.Password = character.Password
		save.Bamfin = character.Bamfin
		save.Bamfout = character.Bamfout
		save.Title = character.Title
		save.Attributes = character.Attributes
		save.ModifiedAttributes = character.ModifiedAttributes

		var skills []saveSkill
		for _, sk := range character.Skills {
			var skill struct {
				SkillID int
				Level   int
			}

			skill.SkillID = sk.SkillID
			skill.Level = sk.Level

			skills = append(skills, skill)
		}
		save.Skills = skills

	}

	var affects []saveAffect
	for _, af := range character.Affects {
		if af.affectType == nil {
			continue
		}

		var affect saveAffect
		affect.AffectType = saveSkill{SkillID: af.affectType.SkillID, Level: af.affectType.Level}
		affect.Duration = af.duration
		affect.Location = af.location
		affect.Modifier = af.modifier
		affect.Bitvector = af.bitVector
		affects = append(affects, affect)
	}
	save.Affects = affects

	save.Inventory = saveItems(character.Inventory, character)

	save.RecallRoomID = character.RecallRoomID
	save.Playable = true

	results, err := json.Marshal(save)

	err = ioutil.WriteFile(path, results, 0655)

	if err != nil {
		panic(err)
	}

	return
}

func saveItems(items []*item, character *mob) []saveItem {
	var results []saveItem
	for _, i := range items {
		if character.Level < i.Level || i.ItemType == itemKey || i.ItemType == itemPotion {
			continue
		}

		var item saveItem
		item.ID = i.index.ID
		item.Name = i.Name
		item.Description = i.Description
		item.ShortDescription = i.ShortDescription
		item.ItemType = i.ItemType

		item.Contained = saveItems(i.container, character)

		var affects []saveAffect
		for _, af := range i.Affected {
			if af.affectType == nil {
				continue
			}

			var affect saveAffect
			affect.AffectType = saveSkill{SkillID: af.affectType.SkillID, Level: af.affectType.Level}
			affect.Duration = af.duration
			affect.Location = af.location
			affect.Modifier = af.modifier
			affect.Bitvector = af.bitVector
			affects = append(affects, affect)
		}
		item.Affected = affects

		item.ExtraFlags = i.ExtraFlags
		item.WearFlags = i.WearFlags
		item.Weight = i.Weight
		item.Value = i.Value
		item.Min = i.Min
		item.Max = i.Max
		item.Timer = i.Timer
		item.Level = i.Level
		item.Cost = i.Cost

		item.SkillID = 0
		if i.Skill != nil {
			item.SkillID = i.Skill.ID
		}

		item.Charges = i.Charges

		results = append(results, item)
	}
	return results
}
