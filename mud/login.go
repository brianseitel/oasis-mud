package mud

import (
	"strings"

	"github.com/brianseitel/oasis-mud/helpers"
)

func Login(c *Connection) *Player {

	var name string
	var password string

	for len(name) == 0 {
		c.SendString("Welcome!" + helpers.Newline + "What's your name? ")
		input, err := c.buffer.ReadString('\n')
		if err != nil {
			panic(err)
		}
		name = strings.Trim(input, "\r\n")
	}

	for len(password) == 0 {
		c.SendString("Password: ")
		input, err := c.buffer.ReadString('\n')
		if err != nil {
			panic(err)
		}
		password = strings.Trim(input, "\r\n")
	}

	var player Player
	db.First(&player, &Player{Username: name, Password: password})

	if db.NewRecord(player) {
		c.SendString(helpers.Red + "Incorrect login. Please try again." + helpers.Reset + helpers.Newline)
		return Login(c)
	}

	player = getPlayer(player)
	player.client = c
	return &player
}

func getPlayer(p Player) Player {
	var (
		room Room
		// inventory []Item
		job  Job
		race Race
	)

	db.First(&p).Related(&job).Related(&race).Related(&room)

	p.Room = room
	// p.Inventory = inventory
	p.Job = job
	p.Race = race
	return p
}
