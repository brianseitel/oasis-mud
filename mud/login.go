package mud

import (
	"strings"

	"github.com/brianseitel/oasis-mud/helpers"
)

func login(c *connection) *player {

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

	var p player
	db.First(&p, &player{Username: name})

	if db.NewRecord(p) {
		p.Username = name
		return register(c, p)
	}

	for len(password) == 0 {
		c.SendString("Password: ")
		input, err := c.buffer.ReadString('\n')
		if err != nil {
			panic(err)
		}
		password = strings.Trim(input, "\r\n")
	}

	notFound := db.First(&p, &player{Username: name, Password: password}).RecordNotFound()
	if notFound {
		c.SendString(helpers.Red + "Incorrect login. Please try again." + helpers.Reset + helpers.Newline)
		return login(c)
	}

	player := getplayer(p)
	player.client = c
	return &player
}

func getplayer(p player) player {
	var (
		room      room
		inventory []item
		job       job
		race      race
	)

	db.First(&p).Related(&job).Related(&race).Related(&room).Related(&inventory)

	p.Room = room
	p.Inventory = inventory
	p.Job = job
	p.Race = race
	return p
}
