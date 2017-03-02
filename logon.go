package main

import (
	"strings"
)

func Login(c *Connection) *Player {

	var name string
	var password string

	for len(name) == 0 {
		c.SendString("Welcome!" + newline + "What's your name? ")
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
		c.SendString(red + err.Error() + reset + newline)
		return Login(c)
	}

	c.SendString("Welcome, " + player.Name + "!" + newline)

	player.m_request = c
	return player
}
