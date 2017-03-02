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

	player, err := LoadPlayer(name, password)
	if err != nil {
		c.SendString(helpers.Red + err.Error() + helpers.Reset + helpers.Newline)
		return Login(c)
	}

	c.SendString("Welcome, " + player.Name + "!" + helpers.Newline)

	player.m_request = c
	return player
}
