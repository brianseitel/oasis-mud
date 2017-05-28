package mud

import (
	"fmt"
	"strings"
	"time"
)

func register(c *connection, name string) *mob {

	input := askIfNew(c)

	if len(input) > 0 && strings.ToLower(input) == "n" {
		return login(c)
	}

	password := askForPassword(c)
	job := askForJob(c)
	race := askForRace(c)

	newPlayer := &mob{
		CreatedAt:    time.Now().String(),
		Name:         name,
		Password:     password,
		Job:          &job,
		Race:         &race,
		Hitpoints:    race.Stats.Hitpoints,
		MaxHitpoints: race.Stats.Hitpoints,
		Mana:         race.Stats.Mana,
		MaxMana:      race.Stats.Mana,
		Movement:     race.Stats.Movement,
		MaxMovement:  race.Stats.Movement,
		Attributes: &attributeSet{
			Strength:     race.Stats.Strength,
			Wisdom:       race.Stats.Wisdom,
			Dexterity:    race.Stats.Dexterity,
			Charisma:     race.Stats.Charisma,
			Constitution: race.Stats.Constitution,
			Intelligence: race.Stats.Intelligence,
		},
		Level:    1,
		Exp:      0,
		Room:     getRoom(1),
		Status:   standing,
		Playable: true,
	}

	c.mob = newPlayer
	newPlayer.client = c
	newPlayer.index = &mobIndex{}
	saveCharacter(newPlayer)
	return newPlayer
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
	for e := jobList.Front(); e != nil; e = e.Next() {
		j := e.Value.(*job)
		c.SendString(fmt.Sprintf("[%s] %s%s", j.Abbr, j.Name, newline))
	}
	for {
		c.SendString("Select a job from above: ")
		input, _ := c.buffer.ReadString('\n')
		input = strings.Trim(input, "\r\n")

		for e := jobList.Front(); e != nil; e = e.Next() {
			j := e.Value.(*job)

			if strings.ToLower(input) == j.Abbr {
				return *j
			}
		}

		c.SendString("\nInvalid selection. Try again." + newline)
	}
}

func askForRace(c *connection) race {
	for e := raceList.Front(); e != nil; e = e.Next() {
		r := e.Value.(*race)
		c.SendString(fmt.Sprintf("[%s] %s%s", r.Abbr, r.Name, newline))
	}

	for {
		c.SendString("Select a country of origin from above: ")
		input, _ := c.buffer.ReadString('\n')

		input = strings.Trim(input, "\r\n")
		for e := raceList.Front(); e != nil; e = e.Next() {
			r := e.Value.(*race)
			if strings.ToLower(input) == r.Abbr {
				return *r
			}
		}
		c.SendString("\nInvalid selection. Try again." + newline)
	}
}
