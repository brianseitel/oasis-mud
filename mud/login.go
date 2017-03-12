package mud

import (
	"strings"

	"github.com/brianseitel/oasis-mud/helpers"
)

func login(c *connection) *mob {

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

	var m mob
	db.First(&m, &mob{Name: name})

	if db.NewRecord(m) {
		m.Name = name
		return register(c, &m)
	}

	for len(password) == 0 {
		c.SendString("Password: ")
		input, err := c.buffer.ReadString('\n')
		if err != nil {
			panic(err)
		}
		password = strings.Trim(input, "\r\n")
	}

	notFound := db.First(&m, &mob{Name: name, Password: password}).RecordNotFound()
	if notFound {
		c.SendString(helpers.Red + "Incorrect login. Please try again." + helpers.Reset + helpers.Newline)
		return login(c)
	}

	mob := &m
	mob.client = c
	mob.Status = standing
	mob.Room = getRoom(uint(m.RoomID))
	mobList.PushBack(m)
	return mob
}
