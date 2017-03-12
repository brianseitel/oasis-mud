package mud

import (
	"fmt"
	"strings"

	"github.com/brianseitel/oasis-mud/helpers"
)

func register(c *connection, m *mob) *mob {

	input := askIfNew(c)

	if len(input) > 0 && strings.ToLower(input) == "n" {
		return login(c)
	}

	password := askForPassword(c)
	job := askForJob(c)
	race := askForRace(c)

	newPlayer := &mob{
		Name:         m.Name,
		Password:     password,
		Job:          job,
		Race:         race,
		Hitpoints:    m.Race.defaultStats("hitpoints"),
		MaxHitpoints: m.Race.defaultStats("hitpoints"),
		Mana:         m.Race.defaultStats("mana"),
		MaxMana:      m.Race.defaultStats("mana"),
		Movement:     m.Race.defaultStats("movement"),
		MaxMovement:  m.Race.defaultStats("movement"),
		Strength:     m.Race.defaultStats("strength"),
		Wisdom:       m.Race.defaultStats("wisdom"),
		Dexterity:    m.Race.defaultStats("dexterity"),
		Charisma:     m.Race.defaultStats("charisma"),
		Constitution: m.Race.defaultStats("constitution"),
		Intelligence: m.Race.defaultStats("intelligence"),
		Level:        1,
		Exp:          0,
		RoomID:       1,
		Status:       standing,
	}

	db.Save(&newPlayer)

	m.client = c
	return m
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
