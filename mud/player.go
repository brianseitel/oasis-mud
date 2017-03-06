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

type Player struct {
	gorm.Model

	Username string `json:"username",gorm:"username"`
	Password string `gorm:"password"`

	Name      string
	Inventory []Item `gorm:"many2many:player_items;"`
	ItemIds   []int  `json:"items" gorm:"-"`
	Room      Room
	RoomId    int `json:"current_room"`
	ExitVerb  string

	Hitpoints    int `json:"hitpoints"`
	MaxHitpoints int `json:"max_hitpoints"`
	Mana         int `json:"mana"`
	MaxMana      int `json:"max_mana"`
	Movement     int `json:"movement"`
	MaxMovement  int `json:"max_movement"`

	Exp   int
	Level int

	Job    Job
	Race   Race
	Gender string

	// MobStats
	client *Connection `gorm:"-"`
}

func (p *Player) move(e Exit) {
	var old_room Room
	db.Find(&old_room, p.RoomId)

	var new_room Room
	db.First(&new_room, e.RoomId)

	for _, rm := range old_room.Mobs {
		rm.notify(fmt.Sprintf("%s leaves heading %s\n", p.Name, e.Dir))
	}

	// add mob to new room list
	p.RoomId = int(new_room.ID)

	for _, rm := range new_room.Mobs {
		rm.notify(fmt.Sprintf("%s arrives in room %d.\n", p.Name, p.RoomId))
	}
}

// Displays the player's status bar
func (p Player) ShowStatusBar() {
	p.client.BufferData(helpers.White + "[" + p.getHitpoints() + helpers.Reset + helpers.Cyan + "hp")
	p.client.BufferData(helpers.White + p.getMana() + helpers.Reset + helpers.Cyan + "mana ")
	p.client.BufferData(helpers.White + p.getMovement() + helpers.Reset + helpers.Cyan + "mv" + helpers.White)
	p.client.BufferData("] >> ")
	p.client.SendBuffer()
}

func (p Player) getRoom() Room {
	var (
		r     Room
		exits []Exit
		items []Item
		mobs  []Mob
	)

	db.First(&r, p.RoomId).Related(&exits, "Exits").Related(&items, "Items").Related(&mobs, "Mobs")

	r.Exits = exits
	r.Items = items
	r.Mobs = mobs
	return r
}

// Retrieves the player's hit points as a string
func (p Player) getHitpoints() string {
	return strconv.Itoa(p.Hitpoints)
}

// Retrieves the player's mana as a string
func (p Player) getMana() string {
	return strconv.Itoa(p.Mana)
}

// Retrieves the player's movement as a string
func (p Player) getMovement() string {
	return strconv.Itoa(p.Movement)
}

// Instantiates a new PlayerDatabase
func NewPlayerDatabase() {
	fmt.Println("Creating Players")
	playerFiles, _ := filepath.Glob("./data/players/*.json")

	for _, playerFile := range playerFiles {
		file, err := ioutil.ReadFile(playerFile)
		if err != nil {
			panic(err)
		}

		var player Player
		err = json.Unmarshal(file, &player)
		if err != nil {
			fmt.Println(string(file))
			panic(err)
		}
		for _, i := range player.ItemIds {
			var item Item
			db.First(&item, i)

			player.Inventory = append(player.Inventory, item)
		}
		var players Player
		db.First(&players, &Player{Username: player.Username})
		if db.NewRecord(players) {
			fmt.Println("\tCreating player " + player.Name + "!")
			db.Create(&player)
		} else {
			fmt.Println("\tSkipping player " + player.Name + "!")
		}
	}
}
