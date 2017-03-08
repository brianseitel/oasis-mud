package mud

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"

	"github.com/brianseitel/oasis-mud/helpers"

	"github.com/jinzhu/gorm"
)

type player struct {
	gorm.Model

	Username string `json:"username",gorm:"username"`
	Password string `gorm:"password"`

	Name      string
	Inventory []item `gorm:"many2many:player_items;"`
	ItemIds   []int  `json:"items" gorm:"-"`
	Room      room
	RoomID    int `json:"current_room" gorm:"room_id"`
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

	client *connection `gorm:"-"`
}

func (p player) TNL() int {
	return (p.Level * 1000) - p.Exp
}

func (p *player) move(e exit) {
	var oldRoom room
	db.First(&oldRoom, p.RoomID)

	var newRoom room
	db.First(&newRoom, e.RoomID)

	for _, rm := range oldRoom.Mobs {
		rm.notify(fmt.Sprintf("%s leaves heading %s\n", p.Name, e.Dir))
	}

	// add mob to new room list
	p.RoomID = int(newRoom.ID)

	for _, rm := range newRoom.Mobs {
		rm.notify(fmt.Sprintf("%s arrives in room %d.\n", p.Name, p.RoomID))
	}
}

// Displays the player's status bar
func (p player) ShowStatusBar() {
	p.client.BufferData(helpers.White + "[" + p.getHitpoints() + helpers.Reset + helpers.Cyan + "hp")
	p.client.BufferData(helpers.White + p.getMana() + helpers.Reset + helpers.Cyan + "mana ")
	p.client.BufferData(helpers.White + p.getMovement() + helpers.Reset + helpers.Cyan + "mv" + helpers.White)
	p.client.BufferData("] >> ")
	p.client.SendBuffer()
}

func (p player) getRoom() room {
	var (
		r room
	)

	db.LogMode(true)
	db.Preload("Exits").Preload("Items").Preload("Mobs").Preload("Players", "active = (?) AND id != (?)", true, p.ID).First(&r, p.RoomID)
	db.LogMode(false)

	return r
}

// Retrieves the player's hit points as a string
func (p player) getHitpoints() string {
	return strconv.Itoa(p.Hitpoints)
}

// Retrieves the player's mana as a string
func (p player) getMana() string {
	return strconv.Itoa(p.Mana)
}

// Retrieves the player's movement as a string
func (p player) getMovement() string {
	return strconv.Itoa(p.Movement)
}

// Instantiates a new playerDatabase
func newPlayerDatabase() {
	fmt.Println("Creating players")
	playerFiles, _ := filepath.Glob("./data/players/*.json")

	for _, playerFile := range playerFiles {
		file, err := ioutil.ReadFile(playerFile)
		if err != nil {
			panic(err)
		}

		var p player
		err = json.Unmarshal(file, &p)
		if err != nil {
			fmt.Println(string(file))
			panic(err)
		}
		for _, i := range p.ItemIds {
			var item item
			db.First(&item, i)

			p.Inventory = append(p.Inventory, item)
		}
		var players player
		db.First(&players, &player{Username: p.Username})
		if db.NewRecord(players) {
			fmt.Println("\tCreating player " + p.Name + "!")
			db.Create(&p)
		} else {
			fmt.Println("\tSkipping player " + p.Name + "!")
		}
	}
}
