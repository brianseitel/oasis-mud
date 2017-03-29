package mud

import (
	"github.com/brianseitel/oasis-mud/helpers"
)

const (
	exitDoor = 1 << iota
	exitClosed
	exitLocked
	exitPickProof
)

type exit struct {
	ID          uint
	Keyword     string
	Description string
	Dir         string `json:"direction"`
	Room        *room
	RoomID      uint `json:"room_id"`
	Key         int
	Flags       uint
}

func (e *exit) hasDoor() bool {
	return helpers.HasBit(e.Flags, exitDoor)
}

func (e *exit) isClosed() bool {
	return helpers.HasBit(e.Flags, exitClosed)
}

func (e *exit) isOpen() bool {
	return !e.isClosed()
}

func (e *exit) isLocked() bool {
	return helpers.HasBit(e.Flags, exitLocked)
}

func (e *exit) isUnlocked() bool {
	return !e.isLocked()
}

func (e *exit) isPickProof() bool {
	return helpers.HasBit(e.Flags, exitPickProof)
}

func (e *exit) isPickable() bool {
	return !e.isPickProof()
}

func reverseDirection(dir string) string {
	switch dir {
	case "east":
		return "west"
	case "west":
		return "east"
	case "up":
		return "down"
	case "down":
		return "up"
	case "north":
		return "south"
	case "south":
		return "north"
	default:
		return "oops"
	}
}
