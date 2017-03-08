package mud

import (
	"fmt"
	"strings"

	"github.com/brianseitel/oasis-mud/helpers"
)

func register(c *connection, p player) *player {

	input := askIfNew(c)

	if len(input) > 0 && strings.ToLower(input) == "n" {
		return login(c)
	}

	password := askForPassword(c)
	job := askForJob(c)
	race := askForRace(c)

	newPlayer := &player{
		Username:     p.Username,
		Password:     password,
		Job:          job,
		Race:         race,
		Hitpoints:    p.Race.defaultStats("hitpoints"),
		MaxHitpoints: p.Race.defaultStats("hitpoints"),
		Mana:         p.Race.defaultStats("mana"),
		MaxMana:      p.Race.defaultStats("mana"),
		Movement:     p.Race.defaultStats("movement"),
		MaxMovement:  p.Race.defaultStats("movement"),
		Strength:     p.Race.defaultStats("strength"),
		Wisdom:       p.Race.defaultStats("wisdom"),
		Dexterity:    p.Race.defaultStats("dexterity"),
		Charisma:     p.Race.defaultStats("charisma"),
		Constitution: p.Race.defaultStats("constitution"),
		Intelligence: p.Race.defaultStats("intelligence"),
		Level:        1,
		Exp:          0,
		RoomID:       1,
	}

	db.Save(&newPlayer)

	p = getplayer(*newPlayer)
	p.client = c
	return &p
}

func askIfNew(c *connection) string {
	c.SendString("Are you a new user? [Y/n]")
	input, _ := c.buffer.ReadString('\n')
	return strings.Trim(input, "\r\n")
}

func askForPassword(c *connection) string {
	c.SendString("Select a password? ")
	input, err := c.buffer.ReadString('\n')
	if err != nil {
		panic(err)
	}
	return strings.Trim(input, "\r\n")
}

func askForJob(c *connection) job {
	var jobs []job
	db.Find(&jobs)
	for _, j := range jobs {
		c.SendString(fmt.Sprintf("[%s] %s%s", j.Abbr, j.Name, helpers.Newline))
	}
	for {
		c.SendString("Select a class from above: ")
		input, _ := c.buffer.ReadString('\n')
		input = strings.Trim(input, "\r\n")

		for _, j := range jobs {
			if strings.ToLower(input) == j.Abbr {
				return j
			}
		}

		c.SendString("\nInvalid selection. Try again." + helpers.Newline)
	}
}

func askForRace(c *connection) race {
	var races []race
	db.Find(&races)
	for _, r := range races {
		c.SendString(fmt.Sprintf("[%s] %s%s", r.Abbr, r.Name, helpers.Newline))
	}

	for {
		c.SendString("Select a race from above: ")
		input, _ := c.buffer.ReadString('\n')

		input = strings.Trim(input, "\r\n")
		for _, r := range races {
			if strings.ToLower(input) == r.Abbr {
				return r
			}
		}
		c.SendString("\nInvalid selection. Try again." + helpers.Newline)

	}
}
