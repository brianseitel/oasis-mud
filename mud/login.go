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

	notFound := db.Preload("Job").Preload("Race").First(&m, &mob{Name: name, Password: password}).RecordNotFound()
	if notFound {
		c.SendString(helpers.Red + "Incorrect login. Please try again." + helpers.Reset + helpers.Newline)
		return login(c)
	}

	file, _ := ioutil.ReadFile(fmt.Sprintf("./data/players/%s.json", name))

	var player *mob
	err := json.Unmarshal(file, &player)
	if err != nil {
		panic(err)
	}

	player.client = c
	player.Status = standing
	player.Job = m.Job
	player.Race = m.Race
	player.Room = getRoom(uint(m.RoomID))
	player.loadSkills()
	mobList.PushBack(player)
	return player
}
