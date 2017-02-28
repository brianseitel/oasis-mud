package main

import "fmt"

type SayCommand struct{}

func (command SayCommand) Handle(c *Connection, line string) {
	c.SendString("You say " + line + "\n")
}

type InventoryCommand struct{}

func (command InventoryCommand) Handle(c *Connection) {
	c.SendString("==========================" + newline)
	c.SendString("Inventory" + newline)
	c.SendString("--------------------------" + newline)
	for _, v := range c.player.inventory {
		item := FindItem(v)
		c.SendString("(1) " + item.Name + newline)
	}
	c.SendString("--------------------------" + newline)
	c.SendString(newline)
}

type MoveCommand struct{}

func (command MoveCommand) Handle(c *Connection, direction string, room Room) {
	for _, v := range room.Exits {
		if v.Dir == direction {
			c.player.Room = v.RoomId
			c.SendString(c.player.exitMessage(direction) + newline)
			return
		}
	}
	c.SendString("There is no exit in that direction." + newline)
}
