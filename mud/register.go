package mud

import (
	"fmt"
	"strings"
)

func register(c *connection, name string) *mob {

	input := askIfNew(c)

	if len(input) > 0 && strings.ToLower(input) == "n" {
		return login(c)
	}

	password := askForPassword(c)
	job := askForJob(c)
	race := askForRace(c)

	fmt.Println(password)
	var m *mob
	newPlayer := &mob{
		Name:         m.Name,
		Job:          &job,
		Race:         &race,
		Hitpoints:    m.Race.defaultStats("hitpoints"),
		MaxHitpoints: m.Race.defaultStats("hitpoints"),
		Mana:         m.Race.defaultStats("mana"),
		MaxMana:      m.Race.defaultStats("mana"),
		Movement:     m.Race.defaultStats("movement"),
		MaxMovement:  m.Race.defaultStats("movement"),
		Attributes: &attributeSet{
			Strength:     m.Race.defaultStats("strength"),
			Wisdom:       m.Race.defaultStats("wisdom"),
			Dexterity:    m.Race.defaultStats("dexterity"),
			Charisma:     m.Race.defaultStats("charisma"),
			Constitution: m.Race.defaultStats("constitution"),
			Intelligence: m.Race.defaultStats("intelligence"),
		},
		ModifiedAttributes: &attributeSet{
			Strength:     m.Race.defaultStats("strength"),
			Wisdom:       m.Race.defaultStats("wisdom"),
			Dexterity:    m.Race.defaultStats("dexterity"),
			Charisma:     m.Race.defaultStats("charisma"),
			Constitution: m.Race.defaultStats("constitution"),
			Intelligence: m.Race.defaultStats("intelligence"),
		},
		Level:  1,
		Exp:    0,
		Room:   getRoom(1),
		Status: standing,
	}

	fmt.Println(newPlayer)

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
	for e := jobList.Front(); e != nil; e = e.Next() {
		j := e.Value.(*job)
		c.SendString(fmt.Sprintf("[%s] %s%s", j.Abbr, j.Name, newline))
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

		c.SendString("\nInvalid selection. Try again." + newline)
	}
}

func askForRace(c *connection) race {
	var races []race

	for e := raceList.Front(); e != nil; e = e.Next() {
		r := e.Value.(*race)
		c.SendString(fmt.Sprintf("[%s] %s%s", r.Abbr, r.Name, newline))
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
		c.SendString("\nInvalid selection. Try again." + newline)

	}
}
