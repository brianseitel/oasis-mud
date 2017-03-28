package mud

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

	file, err := ioutil.ReadFile(fmt.Sprintf("./data/players/%s.json", name))

	var player *mobIndex
	err = json.Unmarshal(file, &player)
	if err != nil {
		player.Name = name
		return register(c, name)
	}

	for len(password) == 0 {
		c.SendString("Password: ")
		input, err := c.buffer.ReadString('\n')
		if err != nil {
			panic(err)
		}
		password = strings.Trim(input, "\r\n")
	}

	if player.Password != password {
		c.SendString(helpers.Red + "Incorrect login. Please try again." + helpers.Reset + helpers.Newline)
		return login(c)
	}

	p := createMob(player)
	p.client = c
	p.Status = standing
	p.Room.Mobs = append(p.Room.Mobs, p)

	mobList.PushBack(p)
	return p
}
