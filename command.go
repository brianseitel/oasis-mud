package main

import (
	"strconv"
	"strings"
)

type CommandDatabase struct {
	Commands map[string]Command
}

type Command interface {
	Handle(c *Connection, line string)
}

// Allows the player to speak to another player in the room
// Command: say hello
type SayCommand struct{}

func (command SayCommand) Handle(c *Connection, line string) {
	c.BroadcastToRoom(line)
	c.SendString("You say, \"" + line + "\"" + newline)
}

// Displays the player's inventory
// Command: inv
type InventoryCommand struct{}

func (command InventoryCommand) Handle(c *Connection, line string) {
	c.BufferData("==========================" + newline)
	c.BufferData("Inventory" + newline)
	c.BufferData("--------------------------" + newline)
	for item, qty := range c.player.getInventory() {
		c.BufferData("(" + strconv.Itoa(qty) + ") " + item + newline)
	}
	c.BufferData("--------------------------" + newline)
	c.BufferData(newline)

	c.SendBuffer()
}

// Moves the player from the current room to a neighboring room
// Command: move (north|south|east|west|up|down)
type MoveCommand struct {
	Dir string
}

func (command MoveCommand) Handle(c *Connection, line string) {
	room := dbRooms.FindRoom(c.player.Room)
	for _, v := range room.Exits {
		if v.Dir == command.Dir {
			c.player.Room = v.RoomId
			c.SendString(c.player.exitMessage(command.Dir) + newline)

			room := dbRooms.FindRoom(c.player.Room)
			room.Display(*c)
			return
		}
	}
	c.SendString("There is no exit in that direction." + newline)
}

// Transfers an item from the ground to the player's inventory
// Command: get <item> OR get all
type GetCommand struct{}

func (command GetCommand) Handle(c *Connection, line string) {
	if len(line) <= 0 {
		c.SendString("Get what?")
		return
	}
	room := dbRooms.FindRoom(c.player.Room)
	items := room.Items

	if strings.Trim(line, "\r\n") == "all" {
		for _, item := range items {
			c.player.AddItem(item)
			room := dbRooms.FindRoom(c.player.Room)
			dbRooms.RemoveItem(room, item)
			c.SendString("You picked up " + item.Name + "." + newline)
		}
		return
	} else {
		for _, item := range items {
			name := strings.ToLower(item.Name)
			if strings.Contains(name, line) {
				c.player.AddItem(item)
				dbRooms.RemoveItem(room, item)
				c.SendString("You picked up " + item.Name + "." + newline)
				return
			}
		}
	}

	c.SendString("Get what?")
}

// Transfers an item from the player's inventory to the ground
// Command: drop <item> OR drop all
type DropCommand struct{}

func (command DropCommand) Handle(c *Connection, line string) {
	if len(line) <= 0 {
		c.SendString("Drop what?")
		return
	}

	items := c.player.Inventory

	if strings.Trim(line, "\r\n") == "all" {
		for _, item := range items {
			c.player.RemoveItem(item)
			room := dbRooms.FindRoom(c.player.Room)
			dbRooms.AddItem(room, item)
			c.SendString("You dropped " + item.Name + " to the ground." + newline)
		}
		return
	} else {
		for _, item := range items {
			name := strings.ToLower(item.Name)
			if strings.Contains(name, line) {
				c.player.RemoveItem(item)
				room := dbRooms.FindRoom(c.player.Room)
				dbRooms.AddItem(room, item)
				c.SendString("You dropped " + item.Name + " to the ground." + newline)
				return
			}
		}
	}

	c.SendString("Drop what?")
}

// Command to view the room
// Command: l, look
type LookCommand struct{}

func (command LookCommand) Handle(c *Connection, line string) {
	room := dbRooms.FindRoom(c.player.Room)
	room.Display(*c)
}

// Saves the player to disk
// Command: save
type SaveCommand struct{}

func (command SaveCommand) Handle(c *Connection, line string) {
	err := c.player.Save()
	if err != nil {
		c.SendString("Oops! Something went wrong!" + newline)
	}

	c.SendString("Saved!" + newline)
}

// Loads all commands
func NewCommandDatabase() *CommandDatabase {
	db := &CommandDatabase{}

	cmds := make(map[string]Command)
	cmds["say"] = SayCommand{}
	cmds["inv"] = InventoryCommand{}
	cmds["n"] = MoveCommand{Dir: "north"}
	cmds["north"] = MoveCommand{Dir: "north"}
	cmds["west"] = MoveCommand{Dir: "west"}
	cmds["w"] = MoveCommand{Dir: "west"}
	cmds["east"] = MoveCommand{Dir: "east"}
	cmds["e"] = MoveCommand{Dir: "east"}
	cmds["south"] = MoveCommand{Dir: "south"}
	cmds["s"] = MoveCommand{Dir: "s"}
	cmds["up"] = MoveCommand{Dir: "up"}
	cmds["u"] = MoveCommand{Dir: "u"}
	cmds["down"] = MoveCommand{Dir: "down"}
	cmds["d"] = MoveCommand{Dir: "down"}
	cmds["l"] = LookCommand{}
	cmds["look"] = LookCommand{}
	cmds["drop"] = DropCommand{}
	cmds["save"] = SaveCommand{}
	cmds["get"] = GetCommand{}

	db.Commands = cmds
	return db
}

func (cDb CommandDatabase) Lookup(cmd string) (Command, bool) {
	if command, ok := dbCommands.Commands[cmd]; ok {
		return command, true
	}

	return nil, false
}
